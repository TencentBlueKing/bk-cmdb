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

package x19_02_15_10

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const (
	tableNameSubscription = "cc_Subscription"
	subscriptionNameField = "subscription_name"
)

func fixEventSubscribeLastTime(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(common.BKDefaultOwnerID)
	cond.Field(subscriptionNameField).Like("process instance refresh")

	data := mapstr.MapStr{
		common.LastTimeField: metadata.Now(),
	}

	err := db.Table(tableNameSubscription).Update(ctx, cond.ToMapStr(), data)
	if err != nil {
		return err
	}
	return nil
}
