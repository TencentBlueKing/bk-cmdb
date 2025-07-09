/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package common

import (
	"fmt"
)

const (
	// PatternIP regular pattern for ip
	PatternIP = `^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(
1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$`
	// PatternMultipleIP regular pattern for Multiple ip
	PatternMultipleIP = `^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(
1\\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)(,(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.
((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d))*$`
	// PatternPort regular pattern for port range
	PatternPort = `(([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))`
)

// PatternMultiplePortRange regular pattern for multiple port range
var PatternMultiplePortRange = fmt.Sprintf(`^((%s-%s)|(%s))(,((%s)|(%s-%s)))*$`,
	PatternPort, PatternPort, PatternPort, PatternPort, PatternPort, PatternPort)
