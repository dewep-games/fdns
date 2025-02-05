/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package zone

import (
	"go.osspkg.com/goppy/v2/plugins"
)

var Plugins = plugins.Inject(
	NewZone,
)
