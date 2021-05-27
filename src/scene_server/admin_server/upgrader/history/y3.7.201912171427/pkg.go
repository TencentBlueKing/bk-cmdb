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
package y3_7_201912171427

import (
	"context"
	"fmt"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

/*
	set_template 新增 version 字段，用于支持加速diff集群与模板

	背景：set_template 层面需要支持与它的所有实例（集群）进行diff，提示用户进行集群同步，
	如果不通过扩展数据结构支持，新的diff需求将需要大量的数据查询才能实现
*/
func init() {
	upgrader.RegistUpgrader("y3.7.201912171427", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	blog.Infof("start execute y3.7.201912171427")

	if err := addSetVersionField(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.7.201912171427] addSetVersionField failed, error %s", err.Error())
		return fmt.Errorf("addSetVersionField failed, error %s", err.Error())
	}

	if err := addSetDefaultVersion(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.7.201912171427] addSetDefaultVersion failed, error %s", err.Error())
		return fmt.Errorf("addSetDefaultVersion failed, error %s", err.Error())
	}

	if err := addSetTemplateDefaultVersion(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.7.201912171427] addSetTemplateDefaultVersion failed, error %s", err.Error())
		return fmt.Errorf("addSetTemplateDefaultVersion failed, error %s", err.Error())
	}

	return nil
}
