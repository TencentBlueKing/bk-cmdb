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
package y3_8_202004141131

import (
	"context"
	"fmt"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

/*
	bk_port_enable字段重命名为bk_enable_port
*/
func init() {
	upgrader.RegistUpgrader("y3.7.202004141131", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	blog.Infof("start execute y3.7.202004141131")

	if err := updateEnablePortAttribute(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.7.202004141131] updateEnablePortAttribute failed, error %s", err.Error())
		return fmt.Errorf("updateEnablePortAttribute failed, error %s", err.Error())
	}

	if err := updateProcessAndProcTemplateEnablePortAttribute(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.7.202004141131] updateEnablePortAttributeValue failed, error %s", err.Error())
		return fmt.Errorf("updateEnablePortAttributeValue failed, error %s", err.Error())
	}

	return nil
}
