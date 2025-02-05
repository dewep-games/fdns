/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package db

import (
	"go.osspkg.com/goppy/v2/orm"
)

type (
	DB interface {
		Master() orm.Stmt
		Slave() orm.Stmt
	}

	object struct {
		orm orm.ORM
	}
)

func (v *object) Master() orm.Stmt {
	return v.orm.Tag("master")
}

func (v *object) Slave() orm.Stmt {
	return v.orm.Tag("slave")
}
