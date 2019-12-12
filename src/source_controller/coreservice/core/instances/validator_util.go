/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instances

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// FillLostedFieldValue fill the value in inst map data
func FillLostedFieldValue(ctx context.Context, valData mapstr.MapStr, propertys []metadata.Attribute) {
	rid := util.ExtractRequestIDFromContext(ctx)
	for _, field := range propertys {
		if field.PropertyID == common.BKChildStr || field.PropertyID == common.BKParentStr {
			continue
		}
		_, ok := valData[field.PropertyID]
		if !ok {
			switch field.PropertyType {
			case common.FieldTypeSingleChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeLongChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeInt:
				valData[field.PropertyID] = nil
			case common.FieldTypeEnum:
				enumOptions, err := metadata.ParseEnumOption(ctx, field.Option)
				if err != nil {
					blog.Warnf("ParseEnumOption failed: %v, rid: %s", err, rid)
					valData[field.PropertyID] = nil
					continue
				}
				if len(enumOptions) > 0 {
					var defaultOption *metadata.EnumVal
					for _, k := range enumOptions {
						if k.IsDefault {
							defaultOption = &k
							break
						}
					}
					if nil != defaultOption {
						valData[field.PropertyID] = defaultOption.ID
					} else {
						valData[field.PropertyID] = nil
					}
				} else {
					valData[field.PropertyID] = nil
				}
			case common.FieldTypeDate:
				valData[field.PropertyID] = nil
			case common.FieldTypeTime:
				valData[field.PropertyID] = nil
			case common.FieldTypeUser:
				valData[field.PropertyID] = nil
			case common.FieldTypeTimeZone:
				valData[field.PropertyID] = nil
			case common.FieldTypeBool:
				valData[field.PropertyID] = false
			default:
				valData[field.PropertyID] = nil
			}
		}
	}
}

func isEmpty(value interface{}) bool {
	return value == nil || value == ""
}
