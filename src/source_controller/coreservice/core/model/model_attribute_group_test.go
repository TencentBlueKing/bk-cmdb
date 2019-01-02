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
	"testing"

	"configcenter/src/common/mapstr"

	"configcenter/src/common/metadata"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestManagerModelAttributeGroup(t *testing.T) {

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

	// create attribute group
	groupID := xid.New().String()
	createGrpResult, err := modelMgr.CreateModelAttributeGroup(defaultCtx, inputModel.Spec.ObjectID, metadata.CreateModelAttributeGroup{
		Data: metadata.Group{
			ObjectID:  inputModel.Spec.ObjectID,
			GroupID:   groupID,
			GroupName: "create_group_test",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createGrpResult)

	// query attribute group
	searchGrpResult, err := modelMgr.SearchModelAttributeGroup(defaultCtx, inputModel.Spec.ObjectID, metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.GroupFieldGroupID: groupID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, searchGrpResult)
	require.Equal(t, int64(1), searchGrpResult.Count)
	require.Equal(t, 1, len(searchGrpResult.Info))
	require.Equal(t, groupID, searchGrpResult.Info[0].GroupID)

	t.Logf("search grp:%v", searchGrpResult)

	// update attribute group
	updateGrpResult, err := modelMgr.UpdateModelAttributeGroup(defaultCtx, inputModel.Spec.ObjectID, metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.GroupFieldGroupName: "update_test_group",
		},
		Condition: mapstr.MapStr{
			metadata.GroupFieldGroupID: groupID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateGrpResult)
	require.Equal(t, uint64(1), updateGrpResult.Count)

	// delete attribute group
	deleteGrpResult, err := modelMgr.DeleteModelAttributeGroup(defaultCtx, inputModel.Spec.ObjectID, metadata.DeleteOption{
		Condition: mapstr.MapStr{
			metadata.GroupFieldGroupID: groupID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, deleteGrpResult)
	require.Equal(t, uint64(1), deleteGrpResult.Count)
}
