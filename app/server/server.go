/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package server

import (
	"github.com/miekg/dns"
	"go.osspkg.com/logx"

	"github.com/osspkg/fdns/app/client"
	"github.com/osspkg/fdns/app/record"
	"github.com/osspkg/fdns/app/rules"
)

type Server struct {
	cli     *client.Client
	rec     *record.Records
	adblock *rules.AdBlock
	regexp  *rules.RegexpRules
	static  *rules.StaticRules
}

func New(
	cli *client.Client,
	rec *record.Records,
	ab *rules.AdBlock,
	rr *rules.RegexpRules,
	sr *rules.StaticRules,
) *Server {
	return &Server{
		cli:     cli,
		rec:     rec,
		adblock: ab,
		regexp:  rr,
		static:  sr,
	}
}

func (v *Server) Exchange(q dns.Question) ([]dns.RR, error) {
	if res, ok := v.rec.Get(q.Qtype, q.Name); ok {
		return record.CreateRR(q.Qtype, q.Name, res.Lifetime, res.Value...), nil
	}

	if values, ok := v.static.Search(q.Qtype, q.Name); ok {
		return record.CreateRR(q.Qtype, q.Name, record.DefaultTTL, values...), nil
	}

	if v.adblock.Contain(q.Name) {
		return nil, nil
	}

	if values, ok := v.regexp.Convert(q.Qtype, q.Name); ok {
		return record.CreateRR(q.Qtype, q.Name, record.DefaultTTL, values...), nil
	}

	resp, err := v.cli.Exchange(q)
	if err != nil {
		logx.Error("DNS Exchange", "err", err, "domain", q.Name)
		return nil, err
	}

	value, ttl := record.ParseRR(resp)
	if len(value) > 0 {
		v.rec.Set(q.Qtype, q.Name, ttl, value...)
	}

	return resp, nil
}
