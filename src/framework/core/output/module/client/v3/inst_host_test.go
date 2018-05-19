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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchHost(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_host_innerip").In([]string{"192.168.100.1"})
	rets, err := cli.Host().SearchHost(cond)
	t.Logf("search host result: %v", rets)
	assert.NoError(t, err)
	assert.NotEmpty(t, rets)
}

func TestDeleteHost(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	err := cli.Host().DeleteHostBatch("1")
	assert.NoError(t, err)
}
func TestUpdateHost(t *testing.T) {
	cli := v3.New(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	data := types.MapStr{"bk_host_name": "test_update"}
	err := cli.Host().UpdateHostBatch(data, "5")
	assert.NoError(t, err)
}
