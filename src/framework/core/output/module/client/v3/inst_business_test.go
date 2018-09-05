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

package v3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client/v3"
	"configcenter/src/framework/core/types"
)

func TestSearchBusiness(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_biz_name").Eq("蓝鲸")
	rets, err := cli.Business().SearchBusiness(cond)
	//t.Logf("search business result: %v", rets)
	assert.NoError(t, err)
	assert.NotEmpty(t, rets)
}

func TestBusiness(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	data := types.MapStr{
		"bk_biz_name":       "testBiz",
		"bk_biz_maintainer": "build_user",
	}
	id, err := cli.Business().CreateBusiness(data)
	//t.Logf("search business result: %v", id)
	require.NoError(t, err)
	require.True(t, id > 0)

	data.Set("bk_biz_maintainer", "test_user")
	err = cli.Business().UpdateBusiness(data, id)
	require.NoError(t, err)

	cond := common.CreateCondition().Field("bk_biz_name").Eq("testBiz")
	rets, err := cli.Business().SearchBusiness(cond)
	require.NoError(t, err)
	require.NotEmpty(t, rets)
	require.Equal(t, "test_user", rets[0]["bk_biz_maintainer"])

	err = cli.Business().DeleteBusiness(id)
	require.NoError(t, err)

}
