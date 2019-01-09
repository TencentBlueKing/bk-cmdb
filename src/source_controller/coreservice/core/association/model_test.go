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
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func searchModelAssociation(t *testing.T, cond universalsql.Condition) (*metadata.QueryResult, error) {
	assoMgr := newAssociation(t)
	return assoMgr.SearchModelAssociation(defaultCtx, metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	})
}

func createModel(t *testing.T, objectID string) {

	modelMgr := newModel(t)
	// create a valid model with a valid classificationID
	inputModel := metadata.CreateModel{}
	// create a valid model with a valid classificationID
	classificationID := xid.New().String()
	result, err := modelMgr.CreateOneModelClassification(defaultCtx, metadata.CreateOneModelClassification{
		Data: metadata.Classification{
			ClassificationID:   classificationID,
			ClassificationName: "test_classification_name_to_test_create_model",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uint64(0), result.Created.ID)

	// crate a new model
	inputModel.Spec.ObjCls = classificationID
	inputModel.Spec.ObjectName = "create_model_for_attribute"
	inputModel.Spec.ObjectID = objectID
	inputModel.Attributes = []metadata.Attribute{
		metadata.Attribute{
			ObjectID:     inputModel.Spec.ObjectID,
			PropertyID:   xid.New().String(),
			PropertyName: xid.New().String(),
		},
	}

	dataResult, err := modelMgr.CreateModel(defaultCtx, inputModel)
	require.NoError(t, err)
	require.NotNil(t, dataResult)
	require.NotEqual(t, uint64(0), dataResult.Created.ID)
}

func createAssociationKind(t *testing.T, asstKindID string, assoMgr core.AssociationOperation) {
	dataResult, err := assoMgr.CreateAssociationKind(defaultCtx, metadata.CreateAssociationKind{
		Data: metadata.AssociationKind{
			AssociationKindID:   asstKindID,
			AssociationKindName: "asst_kind_name_test_model_association",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, dataResult)
	require.NotEqual(t, uint64(0), dataResult.Created.ID)
}

func TestCreateModelAssociationMaintaince(t *testing.T) {

	assoMgr := newAssociation(t)

	// create base data
	firstObjectID := xid.New().String()
	secondObjectID := xid.New().String()
	assoKindID := xid.New().String()
	createModel(t, firstObjectID)
	createModel(t, secondObjectID)
	createAssociationKind(t, assoKindID, assoMgr)

	// create a empty association
	createModelAssoResult, err := assoMgr.CreateModelAssociation(defaultCtx, metadata.CreateModelAssociation{})
	require.NotNil(t, err)
	require.NotNil(t, createModelAssoResult)
	require.Equal(t, common.CCErrCommParamsNeedSet, err.(errors.CCErrorCoder).GetCode())

	// create a  association missing assokind
	createModelAssoResult, err = assoMgr.CreateModelAssociation(defaultCtx, metadata.CreateModelAssociation{
		Spec: metadata.Association{
			AssociationName: xid.New().String(),
			ObjectID:        firstObjectID,
			AsstObjID:       secondObjectID,
		},
	})
	require.NotNil(t, err)
	require.NotNil(t, createModelAssoResult)
	require.Equal(t, common.CCErrCommParamsNeedSet, err.(errors.CCErrorCoder).GetCode())

	// create a valid association
	createModelAssoResult, err = assoMgr.CreateModelAssociation(defaultCtx, metadata.CreateModelAssociation{
		Spec: metadata.Association{
			AssociationName: xid.New().String(),
			ObjectID:        firstObjectID,
			AsstObjID:       secondObjectID,
			AsstKindID:      assoKindID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createModelAssoResult)

	// search
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.AssociationFieldObjectID, Val: firstObjectID})
	queryResult, err := assoMgr.SearchModelAssociation(defaultCtx, metadata.QueryCondition{Condition: cond.ToMapStr()})
	require.NoError(t, err)
	require.NotNil(t, queryResult)

	// update
	updateResult, err := assoMgr.UpdateModelAssociation(defaultCtx, metadata.UpdateOption{
		Data: mapstr.MapStr{
			"bk_obj_asst_name": "test_name_update",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.NotEqual(t, uint64(0), updateResult.Count)
	t.Logf("update result:%v", updateResult)

	// delete
	deleteResult, err := assoMgr.DeleteModelAssociation(defaultCtx, metadata.DeleteOption{Condition: cond.ToMapStr()})
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.NotEqual(t, uint64(0), deleteResult.Count)
}
