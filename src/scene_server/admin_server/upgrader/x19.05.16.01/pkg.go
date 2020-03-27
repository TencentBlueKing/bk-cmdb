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
package x19_05_16_01

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("x19.05.16.01", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("from now on, the cmdb version will be v3.5.x")

	err = changeProcessName(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x19.05.16.01] changeProcessName error  %s", err.Error())
		return err
	}
	err = createServiceTemplateTables(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x19.05.16.01] createServiceTemplateTables error  %s", err.Error())
		return err
	}
	err = addModuleProperty(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x19.05.16.01] addModuleProperty error  %s", err.Error())
		return err
	}
	err = upgradeServiceTemplate(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x19.05.16.01] upgradeServiceTemplate error  %s", err.Error())
		return err
	}
	if err := updateProcessBindIPProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateProcessBindIPProperty error err: %s", err.Error())
		return err
	}
	if err := updateProcessNameProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateProcessNameProperty error, err: %s", err.Error())
		return err
	}
	if err := updateAutoTimeGapProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateAutoTimeGapProperty error, err: %s", err.Error())
		return err
	}
	if err := updateProcNumProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateProcNumProperty error, err: %s", err.Error())
		return err
	}
	if err := updatePriorityProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updatePriorityProperty error, err: %s", err.Error())
		return err
	}
	if err := updateTimeoutProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateTimeoutProperty error, err: %s", err.Error())
		return err
	}
	if err := updateProcessNamePropertyIndex(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateProcessNamePropertyIndex error, err: %s", err.Error())
		return err
	}
	if err := updateFuncNamePropertyIndex(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateFuncNamePropertyIndex error, err: %s", err.Error())
		return err
	}
	if err := deleteProcessUnique(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] deleteProcessUnique error, err: %s", err.Error())
		return err
	}
	if err := addInnerCategory(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] addInnerCategory error, err: %s", err.Error())
		return err
	}
	if err := updateFuncIDProperty(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] updateFuncIDProperty error, err: %s", err.Error())
		return err
	}
	if err := UpdateProcPortPropertyGroupName(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade x19.05.16.01] UpdateProcPortPropertyGroupName error, err: %s", err.Error())
		return err
	}
	return nil
}
