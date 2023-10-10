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

package y3_10_202304012150

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.10.202304012150", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute y3.10.202304012150")

	if err = addModelQuoteRelationCollection(ctx, db); err != nil {
		blog.Errorf("upgrade y3.10.202304012150 add model quote relation collection failed, err: %v", err)
		return err
	}

	if err = addModelQuoteRelationIndex(ctx, db); err != nil {
		blog.Errorf("upgrade y3.10.202304012150 add model quote relation index failed, err: %v", err)
		return err
	}

	if err = addHiddenClassification(ctx, db, conf); err != nil {
		blog.Errorf("upgrade y3.10.202304012150 add hidden table classification failed, err: %v", err)
		return err
	}

	blog.Infof("upgrade y3.10.202304012150 add model quote relation collection success")
	return nil
}
