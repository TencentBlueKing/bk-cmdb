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

package association_test

import (
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCreateAssociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	//create association kind without ID
	inputParams := metadata.CreateAssociationKind{}
	inputParams.Data.AssociationKindName = xid.New().String()
	inputParams.Data.DestinationToSourceNote = "test"
	inputParams.Data.Direction = "test"
	inputParams.Data.SourceToDestinationNote = "test"
	dataResult, err := asstMgr.CreateAssociationKind(defaultCtx, inputParams)
	require.NotNil(t, err)

	//create a valid association kind
	inputParams.Data.AssociationKindID = xid.New().String()
	dataResult, err = asstMgr.CreateAssociationKind(defaultCtx, inputParams)
	require.Nil(t, err)
	require.NotNil(t, dataResult)
	require.NotEqual(t, uint64(0), dataResult.Created.ID)

	//create a new association kind with dup kind ID
	dataResult, err = asstMgr.CreateAssociationKind(defaultCtx, inputParams)
	require.NotNil(t, err)

	//create a new association kind with dup kind Name
	inputParams.Data.AssociationKindName = xid.New().String()
	dataResult, err = asstMgr.CreateAssociationKind(defaultCtx, inputParams)
	require.Nil(t, err)
}

func TestCreateManyAssociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	//create two association kind
	inputParams := metadata.CreateManyAssociationKind{}
	asstKind1 := metadata.AssociationKind{}
	asstKind2 := metadata.AssociationKind{}
	kindID := xid.New().String()

	asstKind1.AssociationKindID = kindID
	asstKind1.AssociationKindName = "test1"
	asstKind1.DestinationToSourceNote = "test"
	asstKind1.Direction = "test"
	asstKind1.SourceToDestinationNote = "test"
	inputParams.Datas = append(inputParams.Datas, asstKind1)

	asstKind2.AssociationKindID = kindID
	asstKind2.AssociationKindName = xid.New().String()
	asstKind2.DestinationToSourceNote = "test"
	asstKind2.Direction = "test"
	asstKind2.SourceToDestinationNote = "test"
	inputParams.Datas = append(inputParams.Datas, asstKind2)

	// create a new association kind without bk_asset_id
	dataResult, err := asstMgr.CreateManyAssociationKind(defaultCtx, inputParams)
	require.Nil(t, err)
	require.Equal(t, 1, len(dataResult.Created))
	require.Equal(t, 1, len(dataResult.Repeated))
}

func TestSetAssociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	kindID := xid.New().String()
	inputParams := metadata.SetAssociationKind{}
	inputParams.Data.AssociationKindName = kindID
	inputParams.Data.DestinationToSourceNote = "test"
	inputParams.Data.Direction = "test"
	inputParams.Data.SourceToDestinationNote = "test"

	//create set association kind with create result
	dataResult, err := asstMgr.SetAssociationKind(defaultCtx, inputParams)
	require.Nil(t, err)
	require.Equal(t, 1, len(dataResult.Created))

	//create set association kind with update result
	dataResult, err = asstMgr.SetAssociationKind(defaultCtx, inputParams)
	require.Nil(t, err)
	require.Equal(t, 1, len(dataResult.Updated))

}

func TestSetManyAssociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	//create two association kind  in the same kind id
	inputParams := metadata.SetManyAssociationKind{}
	asstKind1 := metadata.AssociationKind{}
	asstKind2 := metadata.AssociationKind{}
	kindID := xid.New().String()

	asstKind1.AssociationKindID = kindID
	asstKind1.AssociationKindName = "test1"
	asstKind1.DestinationToSourceNote = "test"
	asstKind1.Direction = "test"
	asstKind1.SourceToDestinationNote = "test"
	inputParams.Datas = append(inputParams.Datas, asstKind1)

	asstKind2.AssociationKindID = kindID
	asstKind2.AssociationKindName = xid.New().String()
	asstKind2.DestinationToSourceNote = "test"
	asstKind2.Direction = "test"
	asstKind2.SourceToDestinationNote = "test"
	inputParams.Datas = append(inputParams.Datas, asstKind2)

	dataResult, err := asstMgr.SetManyAssociationKind(defaultCtx, inputParams)
	require.Nil(t, err)
	require.Equal(t, 1, len(dataResult.Created))
	require.Equal(t, 1, len(dataResult.Updated))

}

func TestUpdateAssociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	//create association kind without ID
	asstKind := metadata.CreateAssociationKind{}
	updateKind := metadata.UpdateOption{}
	kindID := xid.New().String()

	asstKind.Data.AssociationKindID = kindID
	asstKind.Data.AssociationKindName = "test1"
	asstKind.Data.DestinationToSourceNote = "test"
	asstKind.Data.Direction = "test"
	asstKind.Data.SourceToDestinationNote = "test"

	// create a new association kind
	dataResult, err := asstMgr.CreateAssociationKind(defaultCtx, asstKind)
	require.Nil(t, err)
	require.NotNil(t, dataResult.Created)

	updateKind.Condition.Set(common.AssociationKindIDField, kindID)
	updateKind.Data.Set(common.AssociationKindNameField, "test2")

	// update the created association kind
	updateResult, err := asstMgr.UpdateAssociationKind(defaultCtx, updateKind)
	require.Nil(t, err)
	require.Equal(t, 1, updateResult.Count)

}

func TestSearchAssociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	//create association kind
	asstKind := metadata.CreateAssociationKind{}
	kindID := xid.New().String()

	asstKind.Data.AssociationKindID = kindID
	asstKind.Data.AssociationKindName = "test1"
	asstKind.Data.DestinationToSourceNote = "test"
	asstKind.Data.Direction = "test"
	asstKind.Data.SourceToDestinationNote = "test"

	dataResult, err := asstMgr.CreateAssociationKind(defaultCtx, asstKind)
	require.Nil(t, err)
	require.NotNil(t, dataResult.Created)
	require.NotEqual(t, 0, dataResult.Created.ID)

	//search association kind
	searchCond := metadata.QueryCondition{}
	searchCond.Condition = mapstr.MapStr{common.AssociationKindIDField: asstKind}
	searchResult, err := asstMgr.SearchAssociationKind(defaultCtx, searchCond)
	require.Nil(t, err)
	require.NotEqual(t, 0, searchResult.Count)

}

func TestDeleteassociationKind(t *testing.T) {
	asstMgr := newAssociation(t)

	//create association kind
	asstKind := metadata.CreateAssociationKind{}
	kindID := xid.New().String()

	asstKind.Data.AssociationKindID = kindID
	asstKind.Data.AssociationKindName = "test1"
	asstKind.Data.DestinationToSourceNote = "test"
	asstKind.Data.Direction = "test"
	asstKind.Data.SourceToDestinationNote = "test"

	dataResult, err := asstMgr.CreateAssociationKind(defaultCtx, asstKind)
	require.Nil(t, err)
	require.NotNil(t, dataResult.Created)
	require.NotEqual(t, 0, dataResult.Created.ID)

	//delete the created association kind by condition
	searchCond := metadata.DeleteOption{}
	searchCond.Condition = mapstr.MapStr{common.AssociationKindIDField: asstKind}
	searchResult, err := asstMgr.DeleteAssociationKind(defaultCtx, searchCond)
	require.Nil(t, err)
	require.NotEqual(t, 0, searchResult.Count)

}
