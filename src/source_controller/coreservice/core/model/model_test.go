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

package model_test

import (
	"encoding/json"
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCreateOneModel(t *testing.T) {

	modelMgr := newModel(t)

	inputModel := metadata.CreateModel{}
	// create a new empty
	dataResult, err := modelMgr.CreateModel(defaultCtx, inputModel)

	require.NotNil(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, uint64(0), dataResult.Created.ID)
	tmpErr, ok := err.(errors.CCErrorCoder)
	require.True(t, ok, "err must be the errors of the cmdb")
	require.Equal(t, common.CCErrCommParamsNeedSet, tmpErr.GetCode())

	// create a valid model with a invalid classificationID
	inputModel.Spec = metadata.Object{
		ObjectID: xid.New().String(),
		ObjCls:   xid.New().String(),
	}
	inputModel.Attributes = []metadata.Attribute{}
	dataResult, err = modelMgr.CreateModel(defaultCtx, inputModel)

	require.NotNil(t, err)
	require.NotNil(t, dataResult)
	tmpErr, ok = err.(errors.CCErrorCoder)
	require.True(t, ok, "err must be the errors of the cmdb")
	require.Equal(t, common.CCErrCommParamsIsInvalid, tmpErr.GetCode())

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

	inputModel.Spec.ObjCls = classificationID
	inputModel.Spec.ObjectName = "test_create_model"
	inputModel.Spec.ObjectID = xid.New().String()
	inputModel.Attributes = []metadata.Attribute{
		metadata.Attribute{
			ObjectID:     inputModel.Spec.ObjectID,
			PropertyID:   xid.New().String(),
			PropertyName: xid.New().String(),
		},
	}

	dataResult, err = modelMgr.CreateModel(defaultCtx, inputModel)
	require.NoError(t, err)
	require.NotNil(t, dataResult)
	require.NotEqual(t, uint64(0), dataResult.Created.ID)

}

func TestSetOneModel(t *testing.T) {

	modelMgr := newModel(t)

	inputModel := metadata.SetModel{}
	// create a new empty
	dataResult, err := modelMgr.SetModel(defaultCtx, inputModel)

	require.NotNil(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, 0, len(dataResult.Created))
	require.Equal(t, dataResult.CreatedCount.Count, uint64(len(dataResult.Created)))
	require.Equal(t, dataResult.UpdatedCount.Count, uint64(len(dataResult.Updated)))
	require.Equal(t, 0, len(dataResult.Exceptions))

	tmpErr, ok := err.(errors.CCErrorCoder)
	require.True(t, ok, "err must be the errors of the cmdb")
	require.Equal(t, common.CCErrCommParamsNeedSet, tmpErr.GetCode())

	// create a valid model with a invalid classificationID
	inputModel.Spec = metadata.Object{
		ObjectID: xid.New().String(),
		ObjCls:   xid.New().String(),
	}
	inputModel.Attributes = []metadata.Attribute{}
	dataResult, err = modelMgr.SetModel(defaultCtx, inputModel)

	require.NotNil(t, err)
	require.NotNil(t, dataResult)
	tmpErr, ok = err.(errors.CCErrorCoder)
	require.True(t, ok, "err must be the errors of the cmdb")
	require.Equal(t, common.CCErrCommParamsIsInvalid, tmpErr.GetCode())

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

	inputModel.Spec.ObjCls = classificationID
	inputModel.Spec.ObjectName = "test_create_model"
	inputModel.Spec.ObjectID = xid.New().String()
	inputModel.Attributes = []metadata.Attribute{
		metadata.Attribute{
			ObjectID:     inputModel.Spec.ObjectID,
			PropertyID:   xid.New().String(),
			PropertyName: xid.New().String(),
		},
	}

	dataResult, err = modelMgr.SetModel(defaultCtx, inputModel)
	require.NoError(t, err)
	require.NotNil(t, dataResult)
	require.Equal(t, 1, len(dataResult.Created))
	require.NotEqual(t, uint64(0), dataResult.Created[0].ID)
	//require.NotEqual(t, uint64(0), dataResult.Created)

}

func TestSearchAndDeleteModel(t *testing.T) {

	modelMgr := newModel(t)
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

	inputModel.Spec.ObjCls = classificationID
	inputModel.Spec.ObjectName = "delete_create_model"
	inputModel.Spec.ObjectID = xid.New().String()
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

	// search the created one
	searchResult, err := modelMgr.SearchModelWithAttribute(defaultCtx, metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.ModelFieldObjectID: mapstr.MapStr{
				"$regex": inputModel.Spec.ObjectID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, searchResult)
	require.Equal(t, int64(1), searchResult.Count)
	require.Equal(t, searchResult.Count, int64(len(searchResult.Info)))
	resultStr, _ := json.Marshal(searchResult)
	t.Logf("the query result:%s", resultStr)

	// search delete the one
	deleteResult, err := modelMgr.DeleteModel(defaultCtx, metadata.DeleteOption{
		Condition: mapstr.MapStr{
			metadata.ModelFieldObjectID: mapstr.MapStr{
				"$regex": inputModel.Spec.ObjectID,
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.Equal(t, uint64(1), deleteResult.Count)

	// search the created one
	searchResult, err = modelMgr.SearchModelWithAttribute(defaultCtx, metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.ModelFieldObjectID: mapstr.MapStr{
				"$regex": inputModel.Spec.ObjectID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, searchResult)
	require.Equal(t, int64(0), searchResult.Count)
	require.Equal(t, searchResult.Count, int64(len(searchResult.Info)))
}
