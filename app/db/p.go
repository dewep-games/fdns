/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package db

import (
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/plugins"
)

var Plugin = plugins.Plugin{
	Inject: func(v orm.ORM) DB {
		return &object{orm: v}
	},
}
