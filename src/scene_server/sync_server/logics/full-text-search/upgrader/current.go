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
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/types"
)

// initCurrentEsIndex initialize current es indexes, all indices are stored in types.IndexMap
func (u *upgrader) initCurrentEsIndex() {
	if len(types.IndexMap) > 0 {
		return
	}

	for _, name := range types.AllIndexNames {
		meta := &metadata.ESIndexMetadata{
			Settings: u.indexSetting,
			Mappings: metadata.ESIndexMetaMappings{
				Properties: map[string]metadata.ESIndexMetaMappingsProperty{
					metadata.IndexPropertyID:                {PropertyType: metadata.IndexPropertyTypeKeyword},
					metadata.IndexPropertyDataKind:          {PropertyType: metadata.IndexPropertyTypeKeyword},
					metadata.IndexPropertyBKObjID:           {PropertyType: metadata.IndexPropertyTypeKeyword},
					metadata.IndexPropertyBKSupplierAccount: {PropertyType: metadata.IndexPropertyTypeKeyword},
					metadata.IndexPropertyBKBizID:           {PropertyType: metadata.IndexPropertyTypeKeyword},
					metadata.IndexPropertyKeywords:          {PropertyType: metadata.IndexPropertyTypeKeyword},
					metadata.TablePropertyName:              {PropertyType: metadata.IndexPropertyTypeObject},
				},
			},
		}

		for _, field := range types.IndexExcludeFieldsMap[name] {
			delete(meta.Mappings.Properties, field)
		}

		for _, field := range types.IndexExtraFieldsMap[name] {
			meta.Mappings.Properties[field] = metadata.ESIndexMetaMappingsProperty{
				PropertyType: metadata.IndexPropertyTypeKeyword,
			}
		}

		types.IndexMap[name] = []*metadata.ESIndex{metadata.NewESIndex(name, types.IndexVersionMap[name], meta)}
	}
}

// createCurrentEsIndex create current es indexes in es
func (u *upgrader) createCurrentEsIndex(ctx context.Context, rid string) (map[string]struct{}, error) {
	newIndexMap := make(map[string]struct{})

	for name, indexes := range types.IndexMap {
		for _, index := range indexes {
			exists, err := u.createIndex(ctx, index, rid)
			if err != nil {
				return nil, err
			}

			if !exists {
				newIndexMap[name] = struct{}{}
			}

			if err = u.addAlias(ctx, index, rid); err != nil {
				return nil, err
			}
		}
	}

	blog.Infof("finished es index initialization")

	return newIndexMap, nil
}
