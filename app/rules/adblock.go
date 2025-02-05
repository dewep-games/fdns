/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.osspkg.com/algorithms/structs/bloom"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/validate"
	"go.osspkg.com/xc"

	"github.com/osspkg/fdns/app/db"
)

const (
	AdBlockStatic  = "static"
	AdBlockDynamic = "dynamic"
)

var rex = regexp.MustCompile(`(?miU)^\|\|([a-z0-9-.]+)\^(\n|\r|\$)`)

type AdBlock struct {
	bloom *bloom.Bloom
	cli   *web.ClientHttp
	db    db.DB
	mux   sync.RWMutex
}

func NewAdBlock(dbc db.DB, cli web.ClientHttpPool) *AdBlock {
	return &AdBlock{
		cli: cli.Create(),
		db:  dbc,
	}
}

func (v *AdBlock) createBloom(size uint64) (*bloom.Bloom, error) {
	return bloom.New(
		bloom.Quantity(size*2, 0.0001),
	)
}

func (v *AdBlock) Up(ctx xc.Context) error {
	var err error

	if v.bloom, err = v.createBloom(10_000); err != nil {
		return fmt.Errorf("create bloom: %w", err)
	}

	if err = v.Reload(ctx.Context()); err != nil {
		logx.Error("AdBlock reload", "err", err)
	}

	go routine.Interval(ctx.Context(), time.Hour*6, func(ctx context.Context) {
		if e := v.Upgrade(ctx); e != nil {
			logx.Error("AdBlock reload", "err", e)
		}
		if e := v.Reload(ctx); e != nil {
			logx.Error("AdBlock update", "err", e)
		}
	})

	return nil
}

func (v *AdBlock) Down() error {
	return nil
}

func (v *AdBlock) Upgrade(ctx context.Context) error {
	list := make(map[uint64]string, 10)
	err := v.db.Slave().Query(ctx, "get_adblock_list", func(q orm.Querier) {
		q.SQL(`SELECT "id", "url" FROM "black_list" WHERE "deleted_at" IS NULL AND "type" = $1`,
			AdBlockDynamic,
		)
		q.Bind(func(bind orm.Scanner) error {
			var (
				id  uint64
				url string
			)
			if err := bind.Scan(&id, &url); err != nil {
				return err
			}
			list[id] = url
			return nil
		})
	})
	if err != nil {
		return err
	}

	for id, uri := range list {
		var count int
		count, err = func() (int, error) {
			var b []byte
			if e := v.cli.Call(ctx, http.MethodGet, uri, nil, &b); e != nil {
				return 0, e
			}

			rexResult := rex.FindAll(b, -1)
			result := make([]string, 0, 100)

			for _, rr := range rexResult {
				rule := strings.Trim(string(rr[2:len(rr)-1]), "\n^") + "."

				result = append(result, rule)
				if len(result) == 100 {
					if e := v.save(ctx, id, result); e != nil {
						return 0, e
					}
					result = result[:0]
				}
			}

			if len(result) > 0 {
				if e := v.save(ctx, id, result); e != nil {
					return 0, e
				}
			}

			return len(rexResult), nil
		}()
		if err != nil {
			logx.Error("AdBlock upgrade", "uri", uri, "err", err)
		} else {
			logx.Info("AdBlock upgrade", "uri", uri, "count", count)
		}
	}

	return nil
}

func (v *AdBlock) save(ctx context.Context, id uint64, data []string) error {
	return v.db.Master().Tx(ctx, "save_adblock_rules", func(v orm.Tx) {
		v.Exec(func(e orm.Executor) {
			e.SQL(`INSERT INTO "black_list_rules" ("list_id", "data", "created_at", "updated_at")
							VALUES ($1, $2, now(), now()) ON CONFLICT ("data") DO NOTHING;`)
			for _, datum := range data {
				e.Params(id, datum)
			}
		})
	})
}

func (v *AdBlock) Reload(ctx context.Context) error {
	count := 0
	err := v.db.Slave().Query(ctx, "count_adblock_rules", func(q orm.Querier) {
		q.SQL(`SELECT COUNT(blr."id") FROM "black_list_rules" blr
                		LEFT JOIN "black_list" bl on bl."id" = blr."list_id"
                		WHERE bl."deleted_at" IS NULL AND blr."deleted_at" IS NULL;`)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&count)
		})
	})
	if err != nil {
		return err
	}
	if count <= 0 {
		return nil
	}

	var bf *bloom.Bloom
	if bf, err = v.createBloom(uint64(count)); err != nil {
		return err
	}

	err = v.db.Slave().Query(ctx, "load_adblock_rules", func(q orm.Querier) {
		q.SQL(`SELECT blr."data" FROM "black_list_rules" blr
                		LEFT JOIN "black_list" bl on bl."id" = blr."list_id"
                		WHERE bl."deleted_at" IS NULL AND blr."deleted_at" IS NULL;`)
		q.Bind(func(bind orm.Scanner) error {
			var data string
			if err0 := bind.Scan(&data); err0 != nil {
				return err0
			}
			bf.Add([]byte(data))
			return nil
		})
	})
	if err != nil {
		return err
	}

	v.mux.Lock()
	v.bloom = bf
	v.mux.Unlock()

	return nil
}

func (v *AdBlock) Contain(name string) bool {
	has := false

	levels := validate.CountDomainLevels(name)
	params := make([]interface{}, 0, levels)

	v.mux.RLock()
	for i := validate.CountDomainLevels(name); i >= 1; i-- {
		subDomain := validate.GetDomainLevel(name, i)
		params = append(params, subDomain)
		has = has || v.bloom.Contain([]byte(subDomain))
	}
	v.mux.Unlock()

	if !has || len(params) == 0 {
		return false
	}

	count := 0
	err := v.db.Slave().Query(context.Background(), "get_one_adblock_rule", func(q orm.Querier) {
		q.SQL(`SELECT COUNT(blr."id") FROM "black_list_rules" blr
                		LEFT JOIN "black_list" bl on bl."id" = blr."list_id"
                		WHERE bl."deleted_at" IS NULL AND blr."deleted_at" IS NULL AND blr."data" = ANY($1);`,
			params,
		)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&count)
		})
	})
	if err != nil {
		logx.Error("Check domain in adblock", "err", err, "domain", name)
	}

	return count > 0
}
