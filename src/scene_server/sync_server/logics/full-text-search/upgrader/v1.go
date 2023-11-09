/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package upgrader

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
)

func init() {
	RegisterUpgrader(1, upgraderInst.upgradeV1)
}

var v1Indexes = []string{
	"bk_cmdb.bk_biz_set_obj_20210710",
	"bk_cmdb.biz_20210710",
	"bk_cmdb.set_20210710",
	"bk_cmdb.module_20210710",
	"bk_cmdb.host_20210710",
	"bk_cmdb.model_20210710",
	"bk_cmdb.object_instance_20210710",
}

var (
	// tableMappingStr is the json string of table property related index mapping
	tableMappingStr string
)

// upgradeV1 add table property related index mapping
func (u *upgrader) upgradeV1(ctx context.Context, rid string) (*UpgraderFuncResult, error) {
	tableMappings := metadata.ESIndexMetaMappings{Properties: map[string]metadata.ESIndexMetaMappingsProperty{
		metadata.TablePropertyName: {PropertyType: metadata.IndexPropertyTypeObject},
	}}

	var err error
	tableMappingStr, err = json.MarshalToString(tableMappings)
	if err != nil {
		blog.Errorf("marshal table mapping[%+v] failed, err: %v, rid: %s", tableMappings, err, rid)
		return nil, err
	}

	currentIndexMap := make(map[string]struct{})
	for _, indexes := range types.IndexMap {
		for _, index := range indexes {
			currentIndexMap[index.Name()] = struct{}{}
		}
	}

	for _, index := range v1Indexes {
		if _, exists := currentIndexMap[index]; !exists {
			continue
		}

		if err = u.addTablePropertyMapping(ctx, index, rid); err != nil {
			return nil, err
		}

	}

	return &UpgraderFuncResult{Indexes: v1Indexes}, nil
}

// addTablePropertyMapping add table property index mappings if not exists
func (u *upgrader) addTablePropertyMapping(ctx context.Context, name string, rid string) error {
	// check if table property mapping exists in the index
	// table property mapping example: {"mappings":{"properties":{"tables":{"type":"object"}}}}
	IndexMapping, err := u.esCli.GetMapping().
		Index(name).
		Do(ctx)
	if err != nil {
		blog.Errorf("get index[%s] mapping failed, err: %v, rid: %s", name, err, rid)
		return err
	}

	indexMetadata := new(metadata.ESIndexMetadata)
	if err = mapstruct.Decode2StructWithTag(IndexMapping, indexMetadata, "json"); err != nil {
		blog.Errorf("decode index[%s] table mapping %+v failed, err: %v, rid: %s", name, IndexMapping, err, rid)
		return err
	}

	for property := range indexMetadata.Mappings.Properties {
		if property == metadata.TablePropertyName {
			return nil
		}
	}

	// add table property mapping
	_, err = u.esCli.PutMapping().
		BodyString(tableMappingStr).
		Index(name).
		Do(ctx)
	if err != nil {
		blog.Errorf("add index[%s] table mapping %s failed, err: %v, rid: %s", name, tableMappingStr, err, rid)
		return err
	}

	return nil
}
