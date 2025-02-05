/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package api

//go:generate easyjson

import (
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
)

type (
	//easyjson:json
	AdblockListModel []AdblockListItemModel

	//easyjson:json
	AdblockListItemModel struct {
		ID        int64   `json:"id"`
		Data      string  `json:"data"`
		Count     int64   `json:"count"`
		DeletedAt *string `json:"deleted_at"`
	}
)

func (v *Api) AdblockList(ctx web.Context) {
	result := AdblockListModel{}
	err := v.db.Slave().Query(ctx.Context(), "", func(q orm.Querier) {
		q.SQL("SELECT abl.`id`,abl.`data`,abl.`deleted_at`,COUNT(abr.`id`) AS cnt " +
			"FROM `blacklist_adblock_list` AS abl " +
			"JOIN `blacklist_adblock_rules` AS abr ON abr.`list_id` = abl.`id` " +
			"GROUP BY abr.`list_id` " +
			"ORDER BY NULL")
		q.Bind(func(bind orm.Scanner) error {
			item := AdblockListItemModel{}
			if err := bind.Scan(&item.ID, &item.Data, &item.DeletedAt, &item.Count); err != nil {
				return err
			}
			result = append(result, item)
			return nil
		})
	})
	if err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.JSON(200, result)
}
