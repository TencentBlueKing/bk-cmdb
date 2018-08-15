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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func getInlineProcInstKey(hostID, procID int64) string {
	return fmt.Sprintf("%d-%d", hostID, procID)
}

func getGseProcNameSpace(appID, moduleID int64) string {
	return fmt.Sprintf("%d.%d", appID, moduleID)
}

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
				return parseProcInstMatchIntQuestionMark(str, pattern[0], pattern[1])
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
					return parseProcInstMatchIntEnum(str, splitRange[0], splitRange[1], isString)
				} else if isString && 0 < len(splitRange[0]) {
					return common.KvMap{common.BKDBLIKE: fmt.Sprintf("^%s", splitRange[0])}, true, nil
				} else {
					return str, false, nil
				}
			} else {
				return str, false, nil
			}
		default:
			return str, false, nil

		}
	}
	return str, true, nil
}

func parseProcInstMatchIntQuestionMark(rawMatch, part1, part2 string) (data interface{}, notParse bool, err error) {
	var p1, p2 int64
	var isHeader = false
	if 0 != len(part2) {
		p2, err = util.GetInt64ByInterface(part2)
		if nil != err {
			return nil, false, fmt.Errorf("%s not integer", part2)
		}
	}

	if len(part1) == 0 {
		isHeader = true
	} else {
		p1, err = util.GetInt64ByInterface(part1)
		if nil != err {
			return nil, false, fmt.Errorf("%s not integer", part1)
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
	return common.KvMap{common.BKDBIN: num}, false, nil
}

func parseProcInstMatchIntEnum(rawString, part1, part2 string, isString bool) (data interface{}, notParse bool, err error) {
	enumKey := strings.Split(strings.TrimRight(part2, "]"), ",")
	if isString {
		var strs []string
		for _, s := range enumKey {
			strs = append(strs, fmt.Sprintf("%s%s", part1, s))
		}
		return common.KvMap{common.BKDBIN: strs}, false, nil
	} else {
		p1, err := util.GetInt64ByInterface(part1)
		if nil != err {
			return rawString, false, fmt.Errorf("%s not integer", part1)
		}
		var nums []int64
		for _, s := range enumKey {
			p2, err := util.GetInt64ByInterface(s)
			if nil != err {
				return rawString, false, fmt.Errorf("%s not integer", s)
			}
			base := int64(math.Pow10(len(s)))
			nums = append(nums, p1*base+p2)
		}
		return common.KvMap{common.BKDBIN: nums}, false, nil
	}
}

func GetProcInstModel(appID, setID, moduleID, hostID, procID, funcID, procNum int64, maxInstID uint64) []*metadata.ProcInstanceModel {
	if 0 >= procNum {
		procNum = 1
	}
	instProc := make([]*metadata.ProcInstanceModel, 0)
	for numIdx := int64(1); numIdx < procNum+1; numIdx++ {
		procIdx := (maxInstID-1)*uint64(procNum) + uint64(numIdx)
		item := new(metadata.ProcInstanceModel)
		item.ApplicationID = appID
		item.SetID = setID
		item.ModuleID = moduleID
		item.FuncID = funcID
		item.HostID = hostID
		item.HostInstanID = maxInstID
		item.ProcInstanceID = procIdx
		item.ProcID = procID
		item.HostProcID = uint64(numIdx)
		instProc = append(instProc, item)
	}
	return instProc
}

func NewRegexRole(reg string, isString bool) (*RegexRole, error) {
	ret := &RegexRole{
		//isString: isString,
		rawRegex: reg,
	}

	return ret, ret.splitRegexRange()
}

func (r *RegexRole) MatchStr(str string) bool {
	for _, mixed := range r.mixed {
		if mixed == str {
			return true
		}
	}
	strHead, intFoot, isIntFoot := r.splitStrRightN(str, len(r.prefix))
	if strHead != r.prefix || false == isIntFoot {
		return false
	}

	for _, rangeItem := range r.ranges {
		if intFoot >= rangeItem.min && intFoot <= rangeItem.max {
			return true
		}
	}
	return false
}

func (r *RegexRole) MatchInt64(id int64) bool {
	strID := strconv.FormatInt(id, 10)
	for _, mixed := range r.mixed {
		if mixed == strID {
			return true
		}
	}
	strHead, intFoot, hasFoot := r.splitIntRightN(id, len(r.prefix))
	if strHead != r.prefix || false == hasFoot {
		return false
	}

	for _, rangeItem := range r.ranges {
		if intFoot >= rangeItem.min && intFoot <= rangeItem.max {
			return true
		}
	}
	return false
}

func (r *RegexRole) splitStrRightN(str string, n int) (head string, foot int64, isIntFoot bool) {
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

func (r *RegexRole) splitIntRightN(num int64, n int) (head string, foot int64, hasFoot bool) {
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

func (r *RegexRole) splitRegexRange() error {
	splitRange := strings.Split(r.rawRegex, "[")
	if 2 != len(splitRange) {
		return fmt.Errorf("not found regex role ")
	}
	secPart := strings.TrimRight(splitRange[1], "]")
	r.prefix = splitRange[0]
	rangeArr := strings.Split(secPart, ",")
	for _, item := range rangeArr {
		itemSplit := strings.Split(item, "-")
		switch len(itemSplit) {
		case 1:
			r.mixed = append(r.mixed, fmt.Sprintf("%s%s", r.prefix, item))
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
			r.ranges = append(r.ranges, regexRange{min: min, max: max})
		}
	}

	return nil
}

type regexRange struct {
	min int64
	max int64
}
type RegexRole struct {
	//isString bool
	rawRegex string
	prefix   string
	mixed    []string
	ranges   []regexRange
}
