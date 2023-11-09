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
)

// createIndex create index if not exists
func (u *upgrader) createIndex(ctx context.Context, index *metadata.ESIndex, rid string) (bool, error) {
	// check if elastic index exists
	exist, err := u.esCli.IndexExists(index.Name()).Do(ctx)
	if err != nil {
		blog.Errorf("check elastic index[%s] existence failed, err: %v, rid: %s", index.Name(), err, rid)
		return false, err
	}

	if exist {
		return true, nil
	}

	// create new index with the target index name
	_, err = u.esCli.CreateIndex(index.Name()).Body(index.Metadata()).Do(ctx)
	if err != nil {
		blog.Errorf("create elastic index[%s] failed, err: %v, rid: %s", index.Name(), err, rid)
		return false, err
	}

	return false, nil
}

// addAlias add index alias name
func (u *upgrader) addAlias(ctx context.Context, index *metadata.ESIndex, rid string) error {
	aliasName := index.AliasName()
	name := index.Name()

	// it's ok if the alias name index is already exist, but the alias name can not be a real index
	if _, err := u.esCli.Alias().Add(name, aliasName).Do(ctx); err != nil {
		blog.Errorf("bind elastic index[%s] alias %s failed, err: %v, rid: %s", name, aliasName, err, rid)
		return err
	}

	return nil
}
