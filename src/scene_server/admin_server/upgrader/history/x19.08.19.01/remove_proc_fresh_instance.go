/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package x19_08_19_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const (
	tableNameSubscription = "cc_Subscription"
	subscriptionNameField = "subscription_name"
)

func removeProcFreshInstance(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := tableNameSubscription
	SubscriptionName := "process instance refresh [Do not remove it]"
	filter := map[string]interface{}{
		subscriptionNameField:  SubscriptionName,
		common.BKOperatorField: conf.User,
	}
	err := db.Table(tableName).Delete(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
