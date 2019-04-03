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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

func NewScopeMatch(reg string, isString bool) *ScopeMatch {
	ret := &ScopeMatch{
		isString: isString,
		rawRegex: reg,
	}

	return ret
}

// ParseConds parse gse kit control conidtion, eg: 11[1-10],11?,*
func (s *ScopeMatch) ParseConds() (data interface{}, err error) {
	if "*" == s.rawRegex {
		return nil, nil
	} else if -1 != strings.Index(s.rawRegex, "?") {
		if s.isString {
			parttern := strings.Replace(s.rawRegex, "?", ".", -1)
			return common.KvMap{common.BKDBLIKE: fmt.Sprintf("^%s$", parttern)}, nil
		} else {
			pattern := strings.Split(s.rawRegex, "?")
			if 2 == len(pattern) {
				return s.parseMatchIntQuestionMark(pattern[0], pattern[1])
			} else {
				return s.getRealVal(s.rawRegex)
			}
		}
	} else {
		splitRange := strings.Split(s.rawRegex, "[")
		switch len(splitRange) {
		case 0:
			return s.getRealVal(s.rawRegex)
		case 2:
			if strings.HasSuffix(splitRange[1], "]") {
				if -1 == strings.Index(s.rawRegex, "-") {
					return s.parseMatchIntEnum(splitRange[0], splitRange[1])
				} else {
					if s.isString {
						s.needExtCompare = true
						err := s.splitIntScope()
						if nil != err {
							return nil, err
						}
						if 0 < len(splitRange[0]) {
							return common.KvMap{common.BKDBLIKE: fmt.Sprintf("^%s", splitRange[0])}, nil
						} else {
							return nil, nil
						}
					} else {
						s.needExtCompare = true
						return nil, s.splitIntScope()
					}
				}

			} else {
				return s.getRealVal(s.rawRegex)
			}
		default:
			return s.getRealVal(s.rawRegex)

		}
	}
	//return s.getRealVal(s.rawRegex)
}

func (s *ScopeMatch) MatchStr(str string) bool {
	if !s.needExtCompare {
		return true
	}
	for _, mixed := range s.mixed {
		if mixed == str {
			return true
		}
	}
	strHead, intFoot, isIntFoot := s.splitStrRightN(str, len(s.prefix))
	if strHead != s.prefix || false == isIntFoot {
		return false
	}

	for _, rangeItem := range s.ranges {
		if intFoot >= rangeItem.min && intFoot <= rangeItem.max {
			return true
		}
	}
	return false
}

func (s *ScopeMatch) MatchInt64(id int64) bool {
	if !s.needExtCompare {
		return true
	}
	strID := strconv.FormatInt(id, 10)
	for _, mixed := range s.mixed {
		if mixed == strID {
			return true
		}
	}
	strHead, intFoot, hasFoot := s.splitIntRightN(id, len(s.prefix))
	if strHead != s.prefix || false == hasFoot {
		return false
	}

	for _, rangeItem := range s.ranges {
		if intFoot >= rangeItem.min && intFoot <= rangeItem.max {
			return true
		}
	}
	return false
}

// parseMatchIntQuestionMark handle ? for int match
// eg 1? to {10,11,12,13,14,15,16,17,18,19}
func (s *ScopeMatch) parseMatchIntQuestionMark(part1, part2 string) (data interface{}, err error) {
	var p1, p2 int64
	var isHeader = false
	if 0 != len(part2) {
		p2, err = util.GetInt64ByInterface(part2)
		if nil != err {
			return nil, fmt.Errorf("%s not integer", part2)
		}
	}

	if len(part1) == 0 {
		isHeader = true
	} else {
		p1, err = util.GetInt64ByInterface(part1)
		if nil != err {
			return nil, fmt.Errorf("%s not integer", part1)
		}
	}
	base := int64(math.Pow10(len(part2)))
	first := base * 10
	var num []int64
	var i int64 = 0
	if isHeader {
		i = 1
	}
	for ; i < 10; i++ {
		num = append(num, first*p1+i*base+p2)
	}
	return common.KvMap{common.BKDBIN: num}, nil
}

// parseMatchIntEnum  handle emnu for int match
// eg 1[2,14,17] to {12,114, 117}
func (s *ScopeMatch) parseMatchIntEnum(part1, part2 string) (data interface{}, err error) {

	enumKey := strings.Split(strings.TrimRight(part2, "]"), ",")
	if s.isString {
		var strs []string
		for _, s := range enumKey {
			strs = append(strs, fmt.Sprintf("%s%s", part1, s))
		}
		return common.KvMap{common.BKDBIN: strs}, nil
	} else {

		p1 := int64(0)
		if 0 < len(part1) {
			p1, err = util.GetInt64ByInterface(part1)
			if nil != err {
				return nil, fmt.Errorf("%s not integer", part1)
			}
		}

		var nums []int64
		for _, s := range enumKey {
			p2, err := util.GetInt64ByInterface(s)
			if nil != err {
				return nil, fmt.Errorf("%s not integer", s)
			}
			base := int64(math.Pow10(len(s)))
			nums = append(nums, p1*base+p2)
		}
		return common.KvMap{common.BKDBIN: nums}, nil
	}
}

// getRealVal get match  real type for  db diff
func (s *ScopeMatch) getRealVal(str string) (interface{}, error) {
	if s.isString {
		return str, nil
	}
	val, err := util.GetInt64ByInterface(str)
	if nil != err {
		return "", fmt.Errorf("%s not integer", str)
	}
	return val, nil
}

func (s *ScopeMatch) splitStrRightN(str string, n int) (head string, foot int64, isIntFoot bool) {
	strLen := len(str)
	if strLen <= n {
		head = str
		return
	}
	head = str[0:n]
	var err error
	foot, err = util.GetInt64ByInterface(str[n:])
	if nil != err {
		foot = 0
		isIntFoot = false
	} else {
		isIntFoot = true
	}
	return
}

func (s *ScopeMatch) splitIntRightN(num int64, n int) (head string, foot int64, hasFoot bool) {
	strNum := strconv.FormatInt(num, 10)
	if len(strNum) <= n {
		head = strNum
		return
	}
	head = strNum[0:n]
	baseMode := int64(math.Pow10(len(strNum) - n))
	foot = num % baseMode
	hasFoot = true
	return
}

func (s *ScopeMatch) splitIntScope() error {
	splitRange := strings.Split(s.rawRegex, "[")
	if 2 != len(splitRange) {
		return fmt.Errorf("matching rule format error ")
	}
	secPart := strings.TrimRight(splitRange[1], "]")
	s.prefix = splitRange[0]
	rangeArr := strings.Split(secPart, ",")
	for _, item := range rangeArr {
		itemSplit := strings.Split(item, "-")
		switch len(itemSplit) {
		case 1:
			s.mixed = append(s.mixed, fmt.Sprintf("%s%s", s.prefix, item))
		case 2:
			min, err := util.GetInt64ByInterface(itemSplit[0])
			if nil != err {
				return fmt.Errorf("%s not integer", itemSplit[0])
			}
			max, err := util.GetInt64ByInterface(itemSplit[1])
			if nil != err {
				return fmt.Errorf("%s not integer", itemSplit[1])
			}
			if min > max {
				continue
			}
			s.ranges = append(s.ranges, scopeItem{min: min, max: max})
		}
	}

	return nil
}

type scopeItem struct {
	min int64
	max int64
}
type ScopeMatch struct {
	isString       bool // Whether it is a string comparison, otherwise it is a string comparison
	needExtCompare bool //With range comparison, you can't use condition to filter out unsatisfied data, you need to use code for secondary comparison.
	rawRegex       string
	prefix         string
	mixed          []string
	ranges         []scopeItem
}
