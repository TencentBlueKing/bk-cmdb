/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"configcenter/src/common"
	"testing"

	"github.com/rs/xid"

	"configcenter/src/common/metadata"

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
	require.EqualError(t, err, defaultCtx.Error.Error(common.CCErrCommDuplicateItem).Error())

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

	// result, err = modelMgr.CreateManyModelClassification(defaultCtx, metadata.CreateManyModelClassifiaction{})
	// require.NoError(t, err)
	// //require.NotEqual(t, 0, result.Created.ID)

	// // check the exists ID
	// result, err = modelMgr.CreateManyModelClassification(defaultCtx, metadata.CreateManyModelClassifiaction{})
	// require.EqualError(t, err, defaultCtx.Error.Error(common.CCErrCommDuplicateItem).Error())

}
