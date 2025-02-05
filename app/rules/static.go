/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"fmt"
	"hash"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/lib/pq"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/ioutils/pool"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/xc"

	"github.com/osspkg/fdns/app/db"
)

type StaticRules struct {
	db   db.DB
	data cache.TCacheReplace[string, []string]
	pool *pool.Pool[hash.Hash]
}

func NewStaticRules(dbc db.DB) *StaticRules {
	return &StaticRules{
		db:   dbc,
		data: cache.NewWithReplace[string, []string](),
		pool: pool.New[hash.Hash](func() hash.Hash { return xxhash.New() }),
	}
}

func (v *StaticRules) Up(ctx xc.Context) error {
	routine.Interval(ctx.Context(), time.Hour, func(ctx context.Context) {
		if err := v.Reload(ctx); err != nil {
			logx.Error("StaticRules reload", "err", err)
		}
	})
	return nil
}

func (v *StaticRules) Down() error {
	return nil
}

func (v *StaticRules) key(qtype uint16, name string) string {
	h := v.pool.Get()
	defer func() { v.pool.Put(h) }()

	fmt.Fprintf(h, "%d %s", qtype, name) //nolint: errcheck

	return string(h.Sum(nil))
}

func (v *StaticRules) Reload(ctx context.Context) error {
	tmp := make(map[string][]string, 10)

	err := v.db.Slave().Query(ctx, "load_static_rules", func(q orm.Querier) {
		q.SQL(`SELECT "qtype","name","data" FROM "static_list" WHERE "deleted_at" IS NULL;`)
		q.Bind(func(bind orm.Scanner) error {
			var (
				name  string
				qtype uint16
				data  pq.StringArray
			)
			if err := bind.Scan(&qtype, &name, &data); err != nil {
				return err
			}
			tmp[v.key(qtype, name)] = data
			return nil
		})
	})
	if err != nil {
		return err
	}

	v.data.Replace(tmp)

	return nil
}

func (v *StaticRules) Search(qtype uint16, domain string) ([]string, bool) {
	return v.data.Get(v.key(qtype, domain))
}
