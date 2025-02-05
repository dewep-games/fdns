/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"go.osspkg.com/goppy/v2"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
	"go.osspkg.com/goppy/v2/xdns"

	"github.com/osspkg/fdns/app/api"
	"github.com/osspkg/fdns/app/client"
	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/fdns/app/record"
	"github.com/osspkg/fdns/app/rules"
	"github.com/osspkg/fdns/app/server"
	"github.com/osspkg/fdns/app/zone"
)

var Version = "v0.0.0-dev"

func main() {
	app := goppy.New("fDNS", Version, "")
	app.Plugins(
		web.WithServer(),
		web.WithClient(),
		orm.WithORM(),
		orm.WithPgsqlClient(),
		orm.WithMigration(),
		xdns.WithServer(),
		xdns.WithClient(),
	)

	app.Plugins(api.Plugins...)
	app.Plugins(client.Plugins...)
	app.Plugins(db.Plugin)
	app.Plugins(record.Plugins...)
	app.Plugins(rules.Plugins...)
	app.Plugins(server.Plugins...)
	app.Plugins(zone.Plugins...)

	app.Run()
}
