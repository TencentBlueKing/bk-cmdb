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

// Package testmodel model for test only
package testmodel

import (
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/base"
	idgenerator "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/id-generator"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/gen"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// Interface for test only
type Interface interface {
	base.Generic[table.TestModel]
	// GetNameByID Get name by id
	GetNameByID(kt *kit.Kit, id string) (string, error)
}

var _ Interface = new(dao)

// dao test model dao.
type dao struct {
	*base.GenericDao[table.TestModel]
}

// NewDao new test model dao
func NewDao(orm orm.Interface, idGen idgenerator.Interface) Interface {
	return &dao{
		GenericDao: base.NewGenericDao[table.TestModel](orm, idGen),
	}
}

// GetNameByID Get name by id
func (d *dao) GetNameByID(kt *kit.Kit, id string) (string, error) {
	m := gen.Use(d.Orm.DB()).TestModel
	model, err := m.WithContext(kt).
		Where(m.BaseID.Eq(id)).
		Select(m.Name).
		First()
	if err != nil {
		log.Error(kt, "get name by id failed", "err", err)
		return "", err
	}
	return model.Name, nil
}
