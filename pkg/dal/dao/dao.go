/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package dao ...
package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/TencentBlueKing/bk-cmdb/pkg/client/database"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/base"
	idgenerator "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/id-generator"
	testmodel "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/test-model"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
)

// Dao dao interface
type Dao interface {
	base.DynamicConstructor
	TestModel() testmodel.Interface
}

type dao struct {
	orm orm.Interface
	idg idgenerator.Interface
	base.DynamicConstructor
	testModel testmodel.Interface
}

// NewDao new dao instance
func NewDao(ctx context.Context, config *config.DBConfig) (Dao, error) {
	db, err := database.NewGORMClient(config)
	if err != nil {
		return nil, fmt.Errorf("create gorm client failed: %w", err)
	}
	var opts = make([]orm.Option, 0)

	if config.SlowLogThresholdMS > 0 {
		opts = append(opts, orm.SlowRequest(time.Duration(config.SlowLogThresholdMS)*time.Millisecond))
	}

	if config.Debug {
		opts = append(opts, orm.Debug())
	}

	ormInst, err := orm.New(ctx, db, opts...)
	if err != nil {
		return nil, fmt.Errorf("new orm failed: %w", err)
	}
	idg := idgenerator.New(db)

	dynamicConstructor := base.NewDynamicConstructor(ormInst, idg)

	daoInst := &dao{
		orm:                ormInst,
		idg:                idg,
		DynamicConstructor: dynamicConstructor,
		testModel:          testmodel.NewDao(ormInst, idg),
	}
	return daoInst, nil
}

// TestModel test model dao for example use
func (d *dao) TestModel() testmodel.Interface {
	return d.testModel
}
