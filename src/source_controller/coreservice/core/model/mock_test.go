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
	"context"
	"testing"
	"time"

	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/model"
	"configcenter/src/storage/dal/mongo/local"

	"github.com/stretchr/testify/require"
)

type mockDependences struct {
}

// HasInstance used to check if the model has some instances
func (s *mockDependences) HasInstance(ctx core.ContextParams, objIDS []string) (exists bool, err error) {
	return false, nil
}

// HasAssociation used to check if the model has some associations
func (s *mockDependences) HasAssociation(ctx core.ContextParams, objIDS []string) (exists bool, err error) {
	return false, nil
}

// CascadeDeleteAssociation cascade delete all associated data (included instances, model association, instance association) associated with modelObjID
func (s *mockDependences) CascadeDeleteAssociation(ctx core.ContextParams, objIDS []string) error {
	return nil
}

// CascadeDeleteInstances cascade delete all instances(included instances, instance association) associated with modelObjID
func (s *mockDependences) CascadeDeleteInstances(ctx core.ContextParams, objIDS []string) error {
	return nil
}

func newModel(t *testing.T) core.ModelOperation {

	db, err := local.NewMgo("mongodb://cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013/cmdb", time.Minute)
	require.NoError(t, err)
	return model.New(db, &mockDependences{}, nil)
}

var defaultCtx = func() core.ContextParams {
	err, _ := errors.NewFactory("../../../../../resources/errors/")
	lan, _ := language.New("../../../../../resources/language/")
	return core.ContextParams{
		Context:         context.Background(),
		ReqID:           "test_req_id",
		SupplierAccount: "test_owner",
		User:            "test_user",
		Error:           err.CreateDefaultCCErrorIf("en"),
		Lang:            lan.CreateDefaultCCLanguageIf("en"),
	}
}()
