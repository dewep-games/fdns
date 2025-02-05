/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package api

import (
	"go.osspkg.com/goppy/v2/web"

	"github.com/osspkg/fdns/app/db"
)

type Api struct {
	router web.Router
	db     db.DB
}

func NewApi(r web.RouterPool, dbc db.DB) *Api {
	return &Api{
		router: r.Main(),
		db:     dbc,
	}
}

func (v *Api) Up() error {
	api := v.router.Collection("/api")
	api.Get("/blacklist/adblock/list", v.AdblockList)
	return nil
}

func (v *Api) Down() error {
	return nil
}
