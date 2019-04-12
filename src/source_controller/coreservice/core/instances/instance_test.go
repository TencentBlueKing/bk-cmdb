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

package instances_test

import (
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCreateOneInstance(t *testing.T) {

	instMgr := newInstances(t)
	objID := "bk_switch"
	inputParams := metadata.CreateModelInstance{}
	inputParams.Data.Set(common.BKInstNameField, xid.New().String())

	// create a new bk_switch instance without bk_asset_id
	dataResult, err := instMgr.CreateModelInstance(defaultCtx, objID, inputParams)
	require.NotNil(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, uint64(0), dataResult.Created.ID)
	tmpErr, ok := err.(errors.CCErrorCoder)
	require.True(t, ok, "err must be the errors of the cmdb")
	require.Equal(t, common.CCErrCommParamsNeedSet, tmpErr.GetCode())

	// create a valid model  instance with valid params
	inputParams.Data.Set(common.BKAssetIDField, xid.New().String())
	inputParams.Data.Set("bk_sn", "cmdb_sn")
	dataResult, err = instMgr.CreateModelInstance(defaultCtx, objID, inputParams)
	require.Nil(t, err)
	require.NotEqual(t, uint64(0), dataResult.Created.ID)

}

func TestCreateManyInstance(t *testing.T) {

	instMgr := newInstances(t)
	objID := "bk_switch"
	inputParams := metadata.CreateManyModelInstance{}
	inputParams.Datas = append(inputParams.Datas, mapstr.MapStr{
		common.BKInstNameField: xid.New().String(),
		common.BKAssetIDField:  xid.New().String(),
	})

	inputParams.Datas = append(inputParams.Datas, mapstr.MapStr{
		common.BKInstNameField: "test1",
		common.BKAssetIDField:  "assetID1",
	})

	inputParams.Datas = append(inputParams.Datas, mapstr.MapStr{
		common.BKInstNameField: "test1",
		common.BKAssetIDField:  "assetID1",
	})

	inputParams.Datas = append(inputParams.Datas, mapstr.MapStr{
		common.BKInstNameField: xid.New().String(),
		common.BKAssetIDField:  xid.New().String(),
	})

	// create two new bk_switch instance without bk_asset_id
	dataResult, err := instMgr.CreateManyModelInstance(defaultCtx, objID, inputParams)

	require.Nil(t, err)
	require.NotNil(t, dataResult)
	require.NotEqual(t, 0, len(dataResult.Repeated))
	require.NotEqual(t, 0, len(dataResult.Created))
}

func TestUpdateOneInstance(t *testing.T) {

	instMgr := newInstances(t)
	objID := "bk_switch"

	//create one bk_switch instance data
	inputParams := metadata.CreateModelInstance{}
	inputParams.Data.Set(common.BKInstNameField, xid.New().String())
	inputParams.Data.Set(common.BKAssetIDField, xid.New().String())
	inputParams.Data.Set("bk_sn", "cmdb_sn")
	dataResult, err := instMgr.CreateModelInstance(defaultCtx, objID, inputParams)
	require.Nil(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, uint64(0), dataResult.Created.ID)

	//update one bk_switch instance by condition
	updateParams := metadata.UpdateOption{}
	updateParams.Condition = mapstr.MapStr{"bk_sn": "cmdb_sn"}
	updateParams.Data = mapstr.MapStr{"bk_operator": "test"}
	updateResult, err := instMgr.UpdateModelInstance(defaultCtx, objID, updateParams)

	require.Nil(t, err)
	require.NotNil(t, updateResult)
	require.NotEqual(t, uint64(0), updateResult.Count)

}

func TestSearchAndDeleteInstance(t *testing.T) {
	instMgr := newInstances(t)
	objID := "bk_switch"

	//create one bk_switch instance data
	inputParams := metadata.CreateModelInstance{}
	inputParams.Data.Set(common.BKInstNameField, "test_sw1")
	inputParams.Data.Set(common.BKAssetIDField, "test_sw_001")
	inputParams.Data.Set("bk_sn", "cmdb_sn")
	dataResult, err := instMgr.CreateModelInstance(defaultCtx, objID, inputParams)
	require.Nil(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, uint64(0), dataResult.Created.ID)

	//search  this instance
	searchCond := metadata.QueryCondition{}
	searchCond.Condition.Set("bk_sn", "cmdb_sn")
	searchResult, err := instMgr.SearchModelInstance(defaultCtx, objID, searchCond)
	require.Nil(t, err)
	require.NotNil(t, searchResult)
	require.NotEqual(t, uint64(0), searchResult.Count)
	require.NotEqual(t, uint64(0), len(searchResult.Info))

	//delete   this instance
	deleteCond := metadata.DeleteOption{}
	deleteCond.Condition.Set("bk_sn", "cmdb_sn")
	deleteResult, err := instMgr.DeleteModelInstance(defaultCtx, objID, deleteCond)
	require.Nil(t, err)
	require.NotNil(t, deleteResult)
	require.NotEqual(t, uint64(0), deleteResult.Count)
}

func TestCascadeDeleteInstance(t *testing.T) {
	instMgr := newInstances(t)
	objID := "bk_switch"

	//create one bk_switch instance data
	inputParams := metadata.CreateModelInstance{}
	inputParams.Data.Set(common.BKInstNameField, xid.New().String())
	inputParams.Data.Set(common.BKAssetIDField, xid.New().String())
	inputParams.Data.Set("bk_sn", "cmdb_sn")
	dataResult, err := instMgr.CreateModelInstance(defaultCtx, objID, inputParams)
	require.Nil(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, uint64(0), dataResult.Created.ID)

	//delete   this instance
	deleteCond := metadata.DeleteOption{}
	deleteCond.Condition.Set("bk_sn", "cmdb_sn")
	deleteResult, err := instMgr.CascadeDeleteModelInstance(defaultCtx, objID, deleteCond)
	require.Nil(t, err)
	require.NotNil(t, deleteResult)
	require.NotEqual(t, uint64(0), deleteResult.Count)
}
