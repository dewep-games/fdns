/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package record

import (
	"fmt"
	"hash"
	"time"

	"github.com/cespare/xxhash/v2"
	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/ioutils/pool"
	"go.osspkg.com/xc"
)

type Record struct {
	Value    []string
	Lifetime uint32
}

type Records struct {
	data cache.TCacheTTL[string, *Record]
	pool *pool.Pool[hash.Hash]
}

func New(ctx xc.Context) *Records {
	return &Records{
		data: cache.NewWithTTL[string, *Record](ctx.Context(), 15*time.Minute),
		pool: pool.New[hash.Hash](func() hash.Hash { return xxhash.New() }),
	}
}

func (v *Records) key(qtype uint16, name string) string {
	h := v.pool.Get()
	defer func() { v.pool.Put(h) }()

	fmt.Fprintf(h, "%d %s", qtype, name) //nolint: errcheck

	return string(h.Sum(nil))
}

func (v *Records) Set(qtype uint16, name string, ttl uint32, values ...string) {
	v.data.SetWithTTL(
		v.key(qtype, name),
		&Record{Value: values, Lifetime: ttl},
		time.Unix(int64(ttl), 0),
	)
}

func (v *Records) Has(qtype uint16, name string) bool {
	return v.data.Has(v.key(qtype, name))
}

func (v *Records) Get(qtype uint16, name string) (*Record, bool) {
	return v.data.Get(v.key(qtype, name))
}

func (v *Records) Del(qtype uint16, name string) {
	v.data.Del(v.key(qtype, name))
}
