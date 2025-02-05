/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package client

import (
	"sync"

	"github.com/miekg/dns"
	"go.osspkg.com/goppy/v2/xdns"

	"github.com/osspkg/fdns/app/zone"
)

type Client struct {
	cli *xdns.Client
	mux sync.RWMutex
}

func New(z *zone.Zone, cli *xdns.Client) *Client {
	c := &Client{
		cli: cli,
	}

	c.cli.SetZoneResolver(z)

	return c
}

func (v *Client) Exchange(question dns.Question) ([]dns.RR, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	return v.cli.Exchange(question)
}
