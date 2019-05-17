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

func TestAssociationInstance(t *testing.T) {

	asstMgr := newAssociation(t)
	instMgr := newInstances(t)

	kindID := xid.New().String()
	objID := "bk_switch"
	asstObjID := "bk_router"
	associationName := objID + "_" + kindID + "_" + asstObjID

	//create association kind
	asstKindParams := metadata.CreateAssociationKind{}
	asstKindParams.Data.AssociationKindID = kindID
	asstKindParams.Data.AssociationKindName = xid.New().String()
	asstKindParams.Data.DestinationToSourceNote = "test"
	asstKindParams.Data.Direction = "test"
	asstKindParams.Data.SourceToDestinationNote = "test"
	createasstKindResult, err := asstMgr.CreateAssociationKind(defaultCtx, asstKindParams)
	require.Nil(t, err)
	require.NotNil(t, createasstKindResult)
	require.NotEqual(t, 0, createasstKindResult.Created.ID)

	//create association model
	asstModelParams := metadata.CreateModelAssociation{}
	asstModelParams.Spec.AsstKindID = kindID
	asstModelParams.Spec.ObjectID = objID
	asstModelParams.Spec.AsstObjID = asstObjID
	asstModelParams.Spec.AssociationName = associationName
	asstModelParams.Spec.AssociationAliasName = "sw2router"
	createasstModelResult, err := asstMgr.CreateModelAssociation(defaultCtx, metadata.CreateModelAssociation{})
	require.Nil(t, err)
	require.NotNil(t, createasstModelResult)

	//create bk_switch instance
	objInstanceParams := metadata.CreateModelInstance{}
	objInstanceParams.Data.Set(common.BKInstNameField, xid.New().String())
	objInstanceParams.Data.Set(common.BKAssetIDField, xid.New().String())
	instResult, err := instMgr.CreateModelInstance(defaultCtx, objID, objInstanceParams)
	require.Nil(t, err)
	require.NotNil(t, instResult)
	require.NotEqual(t, uint64(0), instResult.Created.ID)
	instID := instResult.Created.ID

	//create bk_router instance
	objAsstInstanceParams := metadata.CreateModelInstance{}
	objAsstInstanceParams.Data.Set(common.BKInstNameField, xid.New().String())
	objAsstInstanceParams.Data.Set(common.BKAssetIDField, xid.New().String())
	instAsstResult, err := instMgr.CreateModelInstance(defaultCtx, asstObjID, objAsstInstanceParams)
	require.Nil(t, err)
	require.NotNil(t, instAsstResult)
	require.NotEqual(t, uint64(0), instAsstResult.Created.ID)
	asstInstID := instAsstResult.Created.ID

	//create association instance
	inputParams := metadata.CreateOneInstanceAssociation{}
	inputParams.Data.AssociationKindID = kindID
	inputParams.Data.AsstInstID = int64(asstInstID)
	inputParams.Data.AsstObjectID = asstObjID
	inputParams.Data.InstID = int64(instID)
	inputParams.Data.ObjectAsstID = associationName
	inputParams.Data.ObjectID = objID
	asstInstResult, err := asstMgr.CreateOneInstanceAssociation(defaultCtx, inputParams)
	require.Nil(t, err)
	require.NotNil(t, asstInstResult)

	//create dup association instance
	asstInstResult, err = asstMgr.CreateOneInstanceAssociation(defaultCtx, inputParams)
	require.NotNil(t, err)

	//create bk_switch instance 1
	objInstanceParams = metadata.CreateModelInstance{}
	objInstanceParams.Data.Set(common.BKInstNameField, xid.New().String())
	objInstanceParams.Data.Set(common.BKAssetIDField, xid.New().String())
	instResult, err = instMgr.CreateModelInstance(defaultCtx, objID, objInstanceParams)
	require.Nil(t, err)
	require.NotNil(t, instResult)
	require.NotEqual(t, uint64(0), instResult.Created.ID)
	instID1 := instResult.Created.ID

	//create bk_switch instance 2
	objInstanceParams = metadata.CreateModelInstance{}
	objInstanceParams.Data.Set(common.BKInstNameField, xid.New().String())
	objInstanceParams.Data.Set(common.BKAssetIDField, xid.New().String())
	instResult, err = instMgr.CreateModelInstance(defaultCtx, objID, objInstanceParams)
	require.Nil(t, err)
	require.NotNil(t, instResult)
	require.NotEqual(t, uint64(0), instResult.Created.ID)
	instID2 := instResult.Created.ID

	//create bk_router instance1
	objAsstInstanceParams = metadata.CreateModelInstance{}
	objAsstInstanceParams.Data.Set(common.BKInstNameField, xid.New().String())
	objAsstInstanceParams.Data.Set(common.BKAssetIDField, xid.New().String())
	instAsstResult, err = instMgr.CreateModelInstance(defaultCtx, asstObjID, objAsstInstanceParams)
	require.Nil(t, err)
	require.NotNil(t, instAsstResult)
	require.NotEqual(t, uint64(0), instAsstResult.Created.ID)
	asstInstID1 := instAsstResult.Created.ID

	//create bk_router instance2
	objAsstInstanceParams = metadata.CreateModelInstance{}
	objAsstInstanceParams.Data.Set(common.BKInstNameField, xid.New().String())
	objAsstInstanceParams.Data.Set(common.BKAssetIDField, xid.New().String())
	instAsstResult, err = instMgr.CreateModelInstance(defaultCtx, asstObjID, objAsstInstanceParams)
	require.Nil(t, err)
	require.NotNil(t, instAsstResult)
	require.NotEqual(t, uint64(0), instAsstResult.Created.ID)
	asstInstID2 := instAsstResult.Created.ID

	//create many association instance
	asstInst1 := metadata.InstAsst{}
	asstInst2 := metadata.InstAsst{}
	asstInst3 := metadata.InstAsst{}

	asstInst1.AssociationKindID = kindID
	asstInst1.AsstInstID = int64(asstInstID1)
	asstInst1.AsstObjectID = asstObjID
	asstInst1.InstID = int64(instID1)
	asstInst1.ObjectAsstID = associationName
	asstInst1.ObjectID = objID

	asstInst2.AssociationKindID = kindID
	asstInst2.AsstInstID = int64(asstInstID2)
	asstInst2.AsstObjectID = asstObjID
	asstInst2.InstID = int64(instID2)
	asstInst2.ObjectAsstID = associationName
	asstInst2.ObjectID = objID

	asstInst3.AssociationKindID = kindID
	asstInst3.AsstInstID = int64(asstInstID)
	asstInst3.AsstObjectID = asstObjID
	asstInst3.InstID = int64(instID)
	asstInst3.ObjectAsstID = associationName
	asstInst3.ObjectID = objID

	createManyParams := metadata.CreateManyInstanceAssociation{}
	createManyParams.Datas = append(createManyParams.Datas, asstInst1, asstInst2)
	instManyAsstResult, err := asstMgr.CreateManyInstanceAssociation(defaultCtx, createManyParams)
	require.Nil(t, err)
	require.NotNil(t, instManyAsstResult)
	require.NotEqual(t, uint64(0), len(instManyAsstResult.Created))
	require.NotEqual(t, uint64(0), len(instManyAsstResult.Repeated))

	//search association instance
	searchCond := metadata.QueryCondition{}
	searchCond.Condition = mapstr.MapStr{common.BKObjIDField: objID}
	asstSearchResult, err := asstMgr.SearchInstanceAssociation(defaultCtx, searchCond)
	require.Nil(t, err)
	require.NotNil(t, asstSearchResult)
	require.NotEqual(t, 0, asstSearchResult.Count)

	//delete association instance
	deleteCond := metadata.DeleteOption{}
	deleteCond.Condition = mapstr.MapStr{common.BKObjIDField: objID}
	deleteResult, err := asstMgr.DeleteInstanceAssociation(defaultCtx, deleteCond)
	require.Nil(t, err)
	require.NotNil(t, deleteResult)
	require.NotEqual(t, 0, deleteResult.Count)
}
