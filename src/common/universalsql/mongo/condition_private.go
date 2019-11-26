/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mongo

import (
	"fmt"
	"reflect"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql"
)

func parseConditionFromMapStr(inputCond *mongoCondition, inputKey string, inputCondMapStr mapstr.MapStr) (outputCond *mongoCondition, err error) {

	outputCond = inputCond
	err = inputCondMapStr.ForEach(func(operatorKey string, val interface{}) error {

		switch operatorKey {
		case universalsql.AND:
			vals, err := inputCondMapStr.MapStrArray(operatorKey)
			if nil != err {
				return err
			}

			if outputCond, err = parseAnd(outputCond, inputKey, vals); nil != err {
				return err
			}

		case universalsql.OR:
			vals, err := inputCondMapStr.MapStrArray(operatorKey)
			if nil != err {
				return err
			}
			if outputCond, err = parseOr(outputCond, inputKey, vals); nil != err {
				return err
			}

		case universalsql.EQ, universalsql.NEQ,
			universalsql.GT, universalsql.GTE, universalsql.LTE, universalsql.LT,
			universalsql.IN, universalsql.NIN, universalsql.EXISTS:
			ele, err := convertToElement(inputKey, operatorKey, val)
			if nil != err {
				return err
			}
			outputCond.Element(ele)
		default:

			tmpType := reflect.TypeOf(val)
			if nil == tmpType {
				// val is nil , return origin by $eq
				outputCond.Element(&Eq{Key: operatorKey, Val: val})
				return nil
			}

			// TODO optimization
			switch tmpType.Kind() {
			case reflect.String,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Float32, reflect.Float64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool, reflect.Slice, reflect.Struct:
				// Compatible with older versions of mongodb equal syntax
				if 0 != len(inputKey) && 0 != len(operatorKey) {
					outputCond.Element(&Eq{Key: inputKey, Val: (&Eq{Key: operatorKey, Val: val}).ToMapStr()})
				} else {
					outputCond.Element(&Eq{Key: operatorKey, Val: val})
				}

			default:
				operatorVal, err := inputCondMapStr.MapStr(operatorKey)
				if nil != err {
					return err
				}

				hit := false
				operatorVal.ForEach(func(key string, val interface{}) error {
					// $regexp operator maybe associate with a $options operator, so it need to skip out
					// when it's a regex operator, do not to step into another embed parse operation.
					if key == universalsql.REGEX {
						hit = true
						return nil
					}
					return nil
				})

				if hit {
					// if hit, then add the element with the follow ways,
					// which is a key:value element, just like the mongodb's
					// original usage.
					outputCond.Element(&KV{Key: operatorKey, Val: val})
					return nil
				}

				tmpCond := newCondition()
				tmpCond, err = parseConditionFromMapStr(tmpCond, operatorKey, operatorVal)
				if nil != err {
					return err
				}
				if 0 != len(inputKey) {
					// ATTENTION: check embed condition, maybe is not a good way
					tmp, ok := outputCond.embed[inputKey]
					if !ok {
						outputCond.embed[inputKey] = tmpCond
						return nil
					}
					tmp.merge(tmpCond)
					return nil
				}

				outputCond.merge(tmpCond)
			}

		} // end operatorKey switch

		return nil
	})

	return outputCond, err

}

func convertToElement(key, operator string, val interface{}) (universalsql.ConditionElement, error) {

	switch operator {
	case universalsql.EQ:
		return &Eq{Key: key, Val: val}, nil
	case universalsql.NEQ:
		return &Neq{Key: key, Val: val}, nil
	case universalsql.GT:
		return &Gt{Key: key, Val: val}, nil
	case universalsql.GTE:
		return &Gte{Key: key, Val: val}, nil
	case universalsql.LTE:
		return &Lte{Key: key, Val: val}, nil
	case universalsql.LT:
		return &Lt{Key: key, Val: val}, nil
	case universalsql.IN:
		return &In{Key: key, Val: val}, nil
	case universalsql.NIN:
		return &Nin{Key: key, Val: val}, nil
	case universalsql.REGEX:
		return &Regex{Key: key, Val: val}, nil
	case universalsql.EXISTS:
		return &Exists{Key: key, Val: val}, nil
	default:
		// deal embed condition
		return nil, fmt.Errorf("not support the operator '%s'", operator)
	}

}

func parseAnd(targetCond *mongoCondition, embedName string, vals []mapstr.MapStr) (outputCond *mongoCondition, err error) {

	outputCond = targetCond
	for _, targetValMapStr := range vals {
		andCond := newCondition()
		andCond, err := parseConditionFromMapStr(andCond, "", targetValMapStr)
		if nil != err {
			return outputCond, err
		}
		if 0 == len(embedName) {
			// ATTENTION: maybe it is not a good way , to check embed condition
			for key, val := range andCond.embed {
				andCond.Element(&Eq{Key: key, Val: val.ToMapStr()})
			}
			outputCond.and = append(outputCond.and, andCond)
		} else {
			andCond.and = append(andCond.and, andCond.elements...)
			andCond.elements = []universalsql.ConditionElement{}
			tmp, ok := outputCond.embed[embedName]
			if !ok {
				outputCond.embed[embedName] = andCond
				continue
			}
			tmp.merge(andCond)
		}
	}

	return outputCond, nil
}

func parseOr(targetCond *mongoCondition, embedName string, vals []mapstr.MapStr) (outputCond *mongoCondition, err error) {

	outputCond = targetCond
	for _, targetValMapStr := range vals {
		orCond := newCondition()
		orCond, err := parseConditionFromMapStr(orCond, "", targetValMapStr)
		if nil != err {
			return outputCond, err
		}

		if 0 == len(embedName) {
			// ATTENTION: maybe it is not a good way , to check embed condition
			for key, val := range orCond.embed {
				orCond.Element(&Eq{Key: key, Val: val.ToMapStr()})
			}
			outputCond.or = append(outputCond.or, orCond)
		} else {
			orCond.or = append(orCond.or, orCond.elements...)
			orCond.elements = []universalsql.ConditionElement{}
			tmp, ok := outputCond.embed[embedName]
			if !ok {
				outputCond.embed[embedName] = orCond
				continue
			}
			tmp.merge(orCond)
		}
	}

	return outputCond, nil
}
