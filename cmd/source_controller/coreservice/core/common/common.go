/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package common TODO
package common

import (
	"configcenter/cmd/source_controller/coreservice/core"
	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/errors"
	"configcenter/pkg/http/rest"
	"configcenter/pkg/metadata"
	"configcenter/pkg/storage/driver/mongodb"
)

var _ core.CommonOperation = (*commonOperation)(nil)

type commonOperation struct {
}

// New create a new instance manager instance
func New() core.CommonOperation {
	return &commonOperation{}
}

// GetDistinctField TODO
func (c *commonOperation) GetDistinctField(kit *rest.Kit, option *metadata.DistinctFieldOption) ([]interface{}, errors.CCErrorCoder) {

	ret, err := mongodb.Client().Table(option.TableName).Distinct(kit.Ctx, option.Field, option.Filter)
	if err != nil {
		blog.Errorf("get distinct field failed, err: %v, option:%#v, rid: %s", err, *option, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return ret, nil
}

// GetDistinctCount 根据条件获取指定表中满足条件数据的数量
func (c *commonOperation) GetDistinctCount(kit *rest.Kit, option *metadata.DistinctFieldOption) (int64,
	errors.CCErrorCoder) {
	var count int64
	ret, err := mongodb.Client().Table(option.TableName).Distinct(kit.Ctx, option.Field, option.Filter)
	if err != nil {
		blog.Errorf("get distinct count failed, err: %v, option:%#v, rid: %s", err, *option, kit.Rid)
		return count, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	count = int64(len(ret))
	return count, nil
}
