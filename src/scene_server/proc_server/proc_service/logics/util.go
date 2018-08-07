/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"fmt"
	"math"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/util"
	//"configcenter/src/common/blog"
)

func ParseProcInstMatchCondition(str string, isString bool) (data interface{}, notParse bool, err error) {
	notParse = false
	if "*" == (str) {
		return nil, false, nil
	} else if -1 != strings.Index(str, "?") {
		if isString {
			parttern := strings.Replace(str, "?", ".", -1)
			return common.KvMap{common.BKDBLIKE: fmt.Sprintf("^%s$", parttern)}, false, nil
		} else {
			pattern := strings.Split(str, "?")
			if 2 == len(pattern) {
				var p1, p2 int64
				var isHeader = false
				if 0 != len(pattern[1]) {
					p2, err = util.GetInt64ByInterface(pattern[1])
					if nil != err {
						return str, false, fmt.Errorf("%s not integer", pattern[1])
					}
				}

				if len(pattern[0]) == 0 {
					isHeader = true
				} else {
					p1, err = util.GetInt64ByInterface(pattern[0])
					if nil != err {
						return str, false, fmt.Errorf("%s not integer", pattern[1])
					}
				}
				base := int64(math.Pow10(len(pattern[1])))
				first := base * 10
				var num []int64
				var i int64 = 0
				if isHeader {
					i = 1
				}
				for ; i < 10; i++ {
					num = append(num, first*p1+i*base+p2)
				}
				return common.KvMap{common.BKDBIN: num}, false, nil
			} else {
				return str, false, nil
			}
		}

	} else {
		splitRange := strings.Split(str, "[")
		switch len(splitRange) {
		case 0:
			return str, false, nil
		case 2:
			if strings.HasSuffix(splitRange[1], "]") {
				if -1 < strings.Index(str, ",") && -1 == strings.Index(str, "-") {
					enumKey := strings.Split(strings.TrimRight(splitRange[1], "]"), ",")
					if isString {
						var strs []string
						for _, s := range enumKey {
							strs = append(strs, fmt.Sprintf("%s%s", splitRange[0], s))
						}
						return common.KvMap{common.BKDBIN: strs}, false, nil
					} else {
						p1, err := util.GetInt64ByInterface(splitRange[0])
						if nil != err {
							return str, false, fmt.Errorf("%s not integer", splitRange[0])
						}
						var nums []int64
						for _, s := range enumKey {
							p2, err := util.GetInt64ByInterface(s)
							if nil != err {
								return str, false, fmt.Errorf("%s not integer", s)
							}
							base := int64(math.Pow10(len(s)))
							nums = append(nums, p1*base+p2)
						}
						return common.KvMap{common.BKDBIN: nums}, false, nil
					}
				} else {
					return nil, true, nil
				}
			} else {
				return str, false, nil
			}
		default:
			return str, false, nil

		}
	}
	return nil, true, nil
}
