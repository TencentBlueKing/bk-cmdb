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
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client/v3"
	"configcenter/src/framework/core/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchGroup(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_obj_id").Eq("biz")
	rets, err := cli.Group().SearchGroups(cond)
	require.NoError(t, err)
	t.Logf("search group result: %v", rets)
	require.NotEmpty(t, rets)
}

func TestGroup(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	data := types.MapStr{"bk_group_id": "d1tupdbszpo1", "bk_group_name": "group1", "bk_group_index": 5, "bk_obj_id": "host", "bk_supplier_account": "0"}
	id, err := cli.Group().CreateGroup(data)
	require.NoError(t, err, "create error")
	t.Logf("create group result: %v", id)
	require.True(t, id > 0, "id unexpect")
	data.Set("id", id)

	cond := common.CreateCondition().Field("bk_group_id").Eq("d1tupdbszpo1").Field("bk_obj_id").Eq("host")
	data.Set("bk_group_name", "updated")
	err = cli.Group().UpdateGroup(data, cond)
	require.NoError(t, err, "update error")

	rets, err := cli.Group().SearchGroups(cond)
	require.NoError(t, err, "search error")
	require.NotEmpty(t, rets)
	require.Equal(t, "updated", rets[0]["bk_group_name"], "unexpected value")

	// cond.Field("id").Eq(id)
	// err = cli.Group().DeleteGroup(cond)
	// require.NoError(t, err, "delete error")

}
