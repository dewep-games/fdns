/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package zone

import (
	"context"
	"time"

	"github.com/lib/pq"
	"go.osspkg.com/do"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/logx"
	"go.osspkg.com/network/address"
	"go.osspkg.com/random"
	"go.osspkg.com/routine"
	"go.osspkg.com/validate"
	"go.osspkg.com/xc"

	"github.com/osspkg/fdns/app/db"
)

type Zone struct {
	db   db.DB
	data cache.TCache[string, []string]
}

func NewZone(dbc db.DB) *Zone {
	return &Zone{
		db:   dbc,
		data: cache.NewWithReplace[string, []string](),
	}
}

func (v *Zone) Up(ctx xc.Context) error {
	routine.Interval(ctx.Context(), 15*time.Minute, func(ctx context.Context) {
		if err := v.Reload(ctx); err != nil {
			logx.Error("Reload dns zone", "err", err)
		}
	})
	return nil
}

func (v *Zone) Down() error {
	return nil
}

func (v *Zone) Resolve(zone string) (result []string) {
	for i := 2; i >= 0; i-- {
		vv := validate.GetDomainLevel(zone, i)
		if ips, ok := v.data.Get(vv); ok {
			result = append(result, ips...)
			return random.Shuffle(result)
		}
	}
	return
}

func (v *Zone) Reload(ctx context.Context) error {
	keys := do.ToMap[string](v.data.Keys())

	err := v.db.Slave().Query(ctx, "reload_dns_zone", func(q orm.Querier) {
		q.SQL(`SELECT "name", "data" FROM "dns_zone" WHERE "deleted_at" IS NULL;`)
		q.Bind(func(bind orm.Scanner) error {
			var (
				name string
				data pq.StringArray
			)
			if err := bind.Scan(&name, &data); err != nil {
				return err
			}

			v.data.Set(name, address.Normalize("53", data...))
			delete(keys, name)

			return nil
		})
	})
	if err != nil {
		return err
	}

	for key := range keys {
		v.data.Del(key)
	}

	return nil
}
