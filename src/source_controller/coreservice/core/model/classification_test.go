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

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCreateOneClassification(t *testing.T) {

	modelMgr := newModel(t)

	// check empty classification ID
	result, err := modelMgr.CreateOneModelClassification(defaultCtx, metadata.CreateOneModelClassification{})
	require.EqualError(t, err, defaultCtx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.ClassFieldClassificationID).Error())

	// check a new classification ID
	classificationID := xid.New().String()
	result, err = modelMgr.CreateOneModelClassification(defaultCtx, metadata.CreateOneModelClassification{Data: metadata.Classification{
		ClassificationID:   "one_" + classificationID,
		ClassificationName: "test_classification_name",
	},
	})
	require.NoError(t, err)
	require.NotEqual(t, 0, result.Created.ID)

	// check the exists ID
	result, err = modelMgr.CreateOneModelClassification(defaultCtx, metadata.CreateOneModelClassification{Data: metadata.Classification{
		ClassificationID:   "one_" + classificationID,
		ClassificationName: "test_classification_name",
	},
	})
	require.EqualError(t, err, defaultCtx.Error.Errorf(common.CCErrCommDuplicateItem, "").Error())

}

func TestSetOneClassification(t *testing.T) {

	modelMgr := newModel(t)

	// check empty classification ID
	result, err := modelMgr.SetOneModelClassification(defaultCtx, metadata.SetOneModelClassification{})
	require.NotNil(t, result)
	require.EqualError(t, err, defaultCtx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.ClassFieldClassificationID).Error())

	// check a new classification ID
	classificationID := xid.New().String()
	result, err = modelMgr.SetOneModelClassification(defaultCtx, metadata.SetOneModelClassification{Data: metadata.Classification{
		ClassificationID:   "one_" + classificationID,
		ClassificationName: "test_classification_name",
	},
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, uint64(0), result.UpdatedCount.Count)
	require.Equal(t, uint64(1), result.CreatedCount.Count)
	require.Equal(t, result.CreatedCount.Count, uint64(len(result.Created)))

	// check the exists ID
	result, err = modelMgr.SetOneModelClassification(defaultCtx, metadata.SetOneModelClassification{Data: metadata.Classification{
		ClassificationID:   "one_" + classificationID,
		ClassificationName: "test_classification_name",
	},
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, uint64(0), result.CreatedCount.Count)
	require.Equal(t, uint64(1), result.UpdatedCount.Count)
	require.Equal(t, result.UpdatedCount.Count, uint64(len(result.Updated)))

}

func TestCreateManyClassification(t *testing.T) {

	modelMgr := newModel(t)

	// check empty classification ID
	result, err := modelMgr.CreateManyModelClassification(defaultCtx, metadata.CreateManyModelClassifiaction{
		Data: []metadata.Classification{},
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, 0, len(result.Created))
	require.Equal(t, 0, len(result.Exceptions))
	require.Equal(t, 0, len(result.Repeated))

	// check create some new instances
	classificationID := xid.New().String()
	result, err = modelMgr.CreateManyModelClassification(defaultCtx, metadata.CreateManyModelClassifiaction{
		Data: []metadata.Classification{
			metadata.Classification{
				ClassificationID:   "many_" + classificationID,
				ClassificationName: "test_classification_name",
			},
			metadata.Classification{
				ClassificationID:   "many_" + classificationID,
				ClassificationName: "test_classification_name",
			},
			metadata.Classification{
				ClassificationName: "test_classification_name",
			},
		},
	})

	require.NotNil(t, result)
	require.NoError(t, err)

	// check created
	require.Equal(t, 1, len(result.Created))
	require.NotEqual(t, uint64(0), result.Created[0].ID)
	require.Equal(t, int64(0), result.Created[0].OriginIndex)

	// check repeated
	require.Equal(t, 1, len(result.Repeated))
	require.Equal(t, int64(1), result.Repeated[0].OriginIndex)

	// check exceptions
	require.Equal(t, 1, len(result.Exceptions))
	require.Equal(t, int64(2), result.Exceptions[0].OriginIndex)
	require.Equal(t, int64(common.CCErrCommParamsNeedSet), result.Exceptions[0].Code)

	t.Log("result created:", result.Created)
	t.Log("result repeated:", result.Repeated)
	t.Log("result exceptions:", result.Exceptions)

}

func TestSetManyClassification(t *testing.T) {

	modelMgr := newModel(t)

	// check empty classification ID
	result, err := modelMgr.SetManyModelClassification(defaultCtx, metadata.SetManyModelClassification{
		Data: []metadata.Classification{},
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, 0, len(result.Created))
	require.Equal(t, 0, len(result.Exceptions))
	require.Equal(t, 0, len(result.Updated))

	// check create some new instances
	classificationID := xid.New().String()
	result, err = modelMgr.SetManyModelClassification(defaultCtx, metadata.SetManyModelClassification{
		Data: []metadata.Classification{
			metadata.Classification{
				ClassificationID:   "many_" + classificationID,
				ClassificationName: "test_classification_name",
			},
			metadata.Classification{
				ClassificationID:   "many_" + classificationID,
				ClassificationName: "test_classification_name",
			},
			metadata.Classification{
				ClassificationName: "test_classification_name",
			},
		},
	})

	require.NotNil(t, result)
	require.NoError(t, err)

	// check created
	require.Equal(t, 1, len(result.Created))
	require.NotEqual(t, uint64(0), result.Created[0].ID)
	require.Equal(t, int64(0), result.Created[0].OriginIndex)

	// check repeated
	require.Equal(t, 1, len(result.Updated))
	require.Equal(t, int64(1), result.Updated[0].OriginIndex)

	// check exceptions
	require.Equal(t, 1, len(result.Exceptions))
	require.Equal(t, int64(2), result.Exceptions[0].OriginIndex)
	require.Equal(t, int64(common.CCErrCommParamsNeedSet), result.Exceptions[0].Code)

	t.Log("result created:", result.Created)
	t.Log("result updated:", result.Updated)
	t.Log("result exceptions:", result.Exceptions)

}

func TestDeleteSearchModelClassification(t *testing.T) {

	modelMgr := newModel(t)
	inputData := []metadata.Classification{
		metadata.Classification{
			ClassificationID:   "delete_" + xid.New().String(),
			ClassificationName: "test_classification_name",
		},
		metadata.Classification{
			ClassificationID:   "delete_" + xid.New().String(),
			ClassificationName: "test_classification_name",
		},
		metadata.Classification{
			ClassificationID:   "delete_" + xid.New().String(),
			ClassificationName: "test_classification_name",
		},
	}
	// check create some new instances
	result, err := modelMgr.SetManyModelClassification(defaultCtx, metadata.SetManyModelClassification{
		Data: inputData,
	})

	require.NoError(t, err)
	require.Equal(t, uint64(3), result.CreatedCount.Count)
	// check search
	queryResult, err := modelMgr.SearchModelClassification(defaultCtx, metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.ClassFieldClassificationID: mapstr.MapStr{
				"$regex": "delete_",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(3), queryResult.Count)
	t.Log("search:", queryResult.Info)

	// delete all classification
	delResult, err := modelMgr.DeleteModelClassificaiton(defaultCtx, metadata.DeleteOption{
		Condition: mapstr.MapStr{
			metadata.ClassFieldClassificationID: mapstr.MapStr{
				"$regex": "delete_",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(3), delResult.Count)
}

func TestUpdateModelClassification(t *testing.T) {

	modelMgr := newModel(t)

	classificationID := xid.New().String()
	// check empty classification ID
	result, err := modelMgr.CreateOneModelClassification(defaultCtx, metadata.CreateOneModelClassification{
		Data: metadata.Classification{
			ClassificationID:   classificationID,
			ClassificationName: "test_classification_name",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uint64(0), result.Created.ID)

	updateResult, err := modelMgr.UpdateModelClassification(defaultCtx, metadata.UpdateOption{
		Data: mapstr.MapStr{
			metadata.ClassFieldClassificationID:   classificationID,
			metadata.ClassFieldClassificationName: "update_classification_name",
		},
		Condition: mapstr.MapStr{
			metadata.ClassFieldClassificationID: mapstr.MapStr{
				"$eq": classificationID,
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, uint64(1), updateResult.Count)

	// check result
	queryResult, err := modelMgr.SearchModelClassification(defaultCtx, metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.ClassFieldClassificationID: mapstr.MapStr{
				"$eq": classificationID,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(queryResult.Info))
	require.Equal(t, queryResult.Count, int64(len(queryResult.Info)))

	for _, item := range queryResult.Info {

		require.Equal(t, item.ClassificationID, classificationID)
		require.Equal(t, item.ClassificationName, "update_classification_name")
	}
}
