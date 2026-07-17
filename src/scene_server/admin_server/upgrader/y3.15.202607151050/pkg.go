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

package y3_15_202607151050

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegisterUpgrade("y3.15.202607151050", upgrade)
}

func upgrade(kit *rest.Kit, db dal.Dal) error {
	blog.Infof("start re-execute y3.14.202502101200 for y3.15.202607151050")
	if err := addHostOsKernelVersionField(kit, db); err != nil {
		blog.Errorf("re-upgrade y3.14.202502101200 add host os kernel version field failed, err: %v", err)
		return err
	}
	blog.Infof("re-execute y3.14.202502101200 for y3.15.202607151050, add host os kernel version field success!")

	blog.Infof("start re-execute y3.14.202603231000 for y3.15.202607151050")
	if err := addCustomResourceCollection(kit, db); err != nil {
		blog.Errorf("re-upgrade y3.14.202603231000 add custom resource collection failed, err: %v", err)
		return err
	}
	blog.Infof("re-execute y3.14.202603231000 for y3.15.202607151050, init custom resource workload collection success")

	blog.Infof("start re-execute y3.14.202607141200 for y3.15.202607151050")
	if err := updateObjAttDesBoolDefaultValue(kit, db); err != nil {
		blog.Errorf("re-upgrade y3.14.202607141200 update bool attribute default value failed, err: %v", err)
		return err
	}
	blog.Infof("re-execute y3.14.202607141200 for y3.15.202607151050, update bool attribute default value success!")

	blog.Infof("start re-execute y3.14.202607141510 for y3.15.202607151050")
	if err := addObjAttIsHiddenField(kit, db); err != nil {
		blog.Errorf("re-upgrade y3.14.202607141510 add objattr is_hidden field failed, err: %v", err)
		return err
	}
	blog.Infof("re-execute y3.14.202607141510 for y3.15.202607151050, add objattr is_hidden field success!")

	blog.Infof("start re-execute y3.14.202607141812 for y3.15.202607151050")
	if err := updateHostBkCPUArchitectureAttr(kit, db); err != nil {
		blog.Errorf("re-upgrade y3.14.202607141812 update host bk_cpu_architecture attribute failed, err: %v", err)
		return err
	}
	blog.Infof("re-execute y3.14.202607141812 for y3.15.202607151050, update host bk_cpu_architecture attr success!")

	return nil
}
