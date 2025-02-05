/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/lib/pq"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/xc"

	"github.com/osspkg/fdns/app/db"
)

type RexRule struct {
	rule  *regexp.Regexp
	qtype uint16
	data  []string
}

func NewRexRule(pattern string, qtype uint16, data []string) (*RexRule, error) {
	rx, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RexRule{
		rule:  rx,
		qtype: qtype,
		data:  data,
	}, nil
}

func (v *RexRule) Compile(name string) []string {
	if v.rule == nil {
		return nil
	}
	result := make([]string, 0, len(v.data))
	matches := v.rule.FindStringSubmatchIndex(name)
	if matches == nil {
		return append(result, v.data...)
	}
	for _, s := range v.data {
		value := v.rule.ExpandString([]byte{}, s, name, matches)
		result = append(result, string(value))
	}
	return result
}

func (v *RexRule) Match(name string) bool {
	if v.rule == nil {
		return false
	}
	return v.rule.MatchString(name)
}

//----------------------------------------------------------------------------------------------------------------------

type RegexpRules struct {
	db   db.DB
	data map[uint16][]*RexRule
	mux  sync.RWMutex
}

func NewRexRules(dbc db.DB) *RegexpRules {
	return &RegexpRules{
		db:   dbc,
		data: make(map[uint16][]*RexRule, 100),
	}
}

func (v *RegexpRules) Up(ctx xc.Context) error {
	routine.Interval(ctx.Context(), time.Hour, func(ctx context.Context) {
		if err := v.Reload(ctx); err != nil {
			logx.Error("RegexpRules reload", "err", err)
		}
	})
	return nil
}

func (v *RegexpRules) Down() error {
	return nil
}

func (v *RegexpRules) Reload(ctx context.Context) error {
	result := make(map[uint16][]*RexRule, 10)
	err := v.db.Master().Query(ctx, "load_regexp_rules", func(q orm.Querier) {
		q.SQL(`SELECT "qtype", "name", "data" FROM "regexp_list" WHERE "deleted_at" IS NULL;`)
		q.Bind(func(bind orm.Scanner) error {
			var (
				name  string
				qtype uint16
				data  pq.StringArray
			)
			if err := bind.Scan(&qtype, &name, &data); err != nil {
				return err
			}

			rr, err := NewRexRule(name, qtype, data)
			if err != nil {
				return err
			}

			val, ok := result[qtype]
			if !ok {
				val = make([]*RexRule, 0, 2)
			}
			val = append(val, rr)
			result[qtype] = val

			return nil
		})
	})
	if err != nil {
		return err
	}
	v.mux.Lock()
	v.data = result
	v.mux.Unlock()
	return nil
}

func (v *RegexpRules) Convert(qtype uint16, name string) ([]string, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if _, ok := v.data[qtype]; !ok {
		return nil, false
	}

	for _, datum := range v.data[qtype] {
		if datum.Match(name) {
			return datum.Compile(name), true
		}
	}

	return nil, false
}
