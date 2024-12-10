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

package collections

import (
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	registerIndexes(common.BKTableNameBaseHost, commHostBaseIndexes)
}

var commHostBaseIndexes = []types.Index{
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkHostName",
		Keys:                    bson.D{{common.BKHostNameField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkHostInnerIP",
		Keys:                    bson.D{{common.BKHostInnerIPField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkCloudInstID",
		Keys:                    bson.D{{common.BKCloudInstIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkCloudID",
		Keys:                    bson.D{{common.BKCloudIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkOsType",
		Keys:                    bson.D{{common.BKOSTypeField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkHostOuterIP",
		Keys:       bson.D{{common.BKHostOuterIPField, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{common.BKHostOuterIPField: map[string]string{
			common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkCloudInstID_bkCloudVendor",
		Keys:       bson.D{{common.BKCloudInstIDField, 1}, {common.BKCloudVendor, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKCloudInstIDField: map[string]string{common.BKDBType: "string"},
			common.BKCloudVendor:      map[string]string{common.BKDBType: "string"},
		},
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkHostID",
		Keys:                    bson.D{{common.BKHostIDField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkAssetID",
		Keys:                    bson.D{{common.BKAssetIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkAgentID",
		Keys:       bson.D{{common.BKAgentIDField, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{common.BKAgentIDField: map[string]string{common.BKDBType: "string",
			common.BKDBGT: ""}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIP_bkCloudID",
		Keys:       bson.D{{common.BKHostInnerIPField, 1}, {common.BKCloudIDField, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKHostInnerIPField: map[string]string{common.BKDBType: "string"},
			common.BKCloudIDField:     map[string]string{common.BKDBType: "number"},
			common.BKAddressingField:  common.BKAddressingStatic,
		},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIPV6_bkCloudID",
		Keys:       bson.D{{common.BKHostInnerIPv6Field, 1}, {common.BKCloudIDField, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKCloudIDField:       map[string]string{common.BKDBType: "number"},
			common.BKAddressingField:    common.BKAddressingStatic,
			common.BKHostInnerIPv6Field: map[string]string{common.BKDBType: "string"},
		},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkCloudID_bkHostInnerIP",
		Keys:       bson.D{{common.BKCloudIDField, 1}, {common.BKHostInnerIPField, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKCloudIDField:     map[string]string{common.BKDBType: "number"},
			common.BKHostInnerIPField: map[string]string{common.BKDBType: "string"},
			common.BKAddressingField:  common.BKAddressingStatic,
		},
	},
}
