/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package record

import (
	"net"

	"github.com/miekg/dns"
)

const (
	DefaultTTL = 3600
)

func CreateRR(qtype uint16, name string, ttl uint32, values ...string) []dns.RR {
	result := make([]dns.RR, 0, len(values))
	for _, value := range values {
		switch qtype {
		case dns.TypeA:
			result = append(result, &dns.A{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				A: net.ParseIP(value),
			})
		case dns.TypeAAAA:
			result = append(result, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				AAAA: net.ParseIP(value),
			})
		}
	}
	return result
}

func ParseRR(v []dns.RR) ([]string, uint32) {
	result := make([]string, 0, len(v))
	var ttl uint32 = 0
	for _, rr := range v {
		switch vv := rr.(type) {
		case *dns.A:
			result = append(result, vv.A.String())
		case *dns.AAAA:
			result = append(result, vv.AAAA.String())
		}
		if ttl < rr.Header().Ttl {
			ttl = rr.Header().Ttl
		}
	}

	return result, ttl
}
