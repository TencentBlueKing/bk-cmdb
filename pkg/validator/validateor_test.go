/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package validator

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	EnvID    string `json:"uid" validate:"required"`
	Force    bool   `json:"-" validate:"required"`
	Operator string `validate:"required"`
}

type testStruct2 struct {
	EnvID    string `json:"uid" validate:"required"`
	Force    bool   `json:"-"`
	Operator string `validate:"required"`
}

type testStruct3 struct {
	Name  string `json:"name" validate:"required,gt=1"`
	Count int    `json:"count" validate:"required,gt=0"`
}

func (s *testStruct2) Validate(ctx context.Context) error {
	return fmt.Errorf("hit validate")
}

func TestValidate(t *testing.T) {
	d := testStruct{}
	err := Struct(context.Background(), d)
	assert.Equal(t, err.Error(), "uid is a required field")

	t2 := &testStruct2{
		EnvID:    "abc",
		Force:    false,
		Operator: "ab",
	}

	err = Struct(context.Background(), t2)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "hit validate")
	}

	t3 := testStruct3{
		Name:  "dd_dabc",
		Count: 1,
	}
	err = Struct(context.Background(), t3)
	assert.NoError(t, err)

	name := "testValidate"
	err = Struct(context.Background(), name)
	assert.Equal(t, err.Error(), "validator: (nil string)")
	// assert.Equal(t, err.Error(), "Force is required")
	// assert.Equal(t, err.Error(), "Operator is required")
}
