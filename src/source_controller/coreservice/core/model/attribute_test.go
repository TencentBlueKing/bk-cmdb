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

func TestAttributeMaintainces(t *testing.T) {

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

	// crate a new model
	objectID := xid.New().String()
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

	// create attribute
	propertyID := xid.New().String()
	createAttrResult, err := modelMgr.CreateModelAttributes(defaultCtx, objectID, metadata.CreateModelAttributes{
		Attributes: []metadata.Attribute{
			metadata.Attribute{
				PropertyID:   propertyID,
				PropertyName: "create_attribute",
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, createAttrResult)
	require.Equal(t, uint64(1), uint64(len(createAttrResult.Created)))
	require.NotEqual(t, uint64(0), createAttrResult.Created[0].ID)
	require.Equal(t, uint64(0), uint64(len(createAttrResult.Repeated)))
	require.Equal(t, uint64(0), uint64(len(createAttrResult.Exceptions)))

	// update attribute
	updateResult, err := modelMgr.UpdateModelAttributes(defaultCtx, objectID, metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.AttributeFieldPropertyName: "create_attribute_to_updated",
		},
		Condition: mapstr.MapStr{
			metadata.AttributeFieldPropertyID: mapstr.MapStr{
				"$eq": propertyID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Equal(t, uint64(1), updateResult.Count, propertyID)

	// search attribute
	searchResult, err := modelMgr.SearchModelAttributes(defaultCtx, objectID, metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.AttributeFieldPropertyID: mapstr.MapStr{
				"$eq": propertyID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, searchResult)
	require.Equal(t, uint64(1), uint64(len(searchResult.Info)))
	require.Equal(t, int64(1), searchResult.Count)

	for _, attr := range searchResult.Info {
		require.Equal(t, "create_attribute_to_updated", attr.PropertyName)
	}

	// delete the attribues
	deleteResult, err := modelMgr.DeleteModelAttributes(defaultCtx, objectID, metadata.DeleteOption{
		Condition: mapstr.MapStr{
			metadata.AttributeFieldPropertyID: mapstr.MapStr{
				"$eq": propertyID,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, deleteResult)
	require.Equal(t, uint64(1), deleteResult.Count, propertyID)

}
