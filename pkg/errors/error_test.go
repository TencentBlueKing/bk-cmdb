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

package cerr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// 测试error转换成responseErr
// 测试error的detail解析
// 测试join类型解析
// 测试validator解析
func Test_ErrorConv(t *testing.T) {
	errManager := NewErrorManager("cmdb")

	t.Run("error translate test", func(t *testing.T) {
		err := NewError(INVALID_REQUEST, "test invalid request")
		respErr := errManager.ConvToRespError(err)
		assert.Equal(t, "cmdb", respErr.System)
		assert.Equal(t, INVALID_REQUEST, respErr.Code)
		assert.Equal(t, "test invalid request", respErr.Details[0])
		assert.Equal(t, "test invalid request", respErr.Details[0])
	})

	t.Run("error join parse test", func(t *testing.T) {
		testJoinErr_1 := fmt.Errorf("this is error one")
		testJoinErr_2 := fmt.Errorf("this is error two")
		testJoinErr_3 := fmt.Errorf("this is error three")
		resultError := errors.Join(testJoinErr_1, testJoinErr_2, testJoinErr_3)
		details := errManager.UnwrapDetails(resultError)
		assert.Equal(t, 3, len(details))
		assert.Equal(t, "this is error one", details[0])
		assert.Equal(t, "this is error two", details[1])
		assert.Equal(t, "this is error three", details[2])
	})

}

// RegisterModel struct for test validate
type RegisterModel struct {
	Username string `validate:"required,numeric"`
	Password string `validate:"required,numeric"`
	Name     string `validate:"required"`
	Age      int    `validate:"required,gte=0,lte=100,numeric"`
	Gender   string `validate:"required,oneof=男 女"`
}

func Test_validate(t *testing.T) {
	errManager := NewErrorManager("cmdb")
	model := RegisterModel{
		Username: "123",
		Password: "456",
		Name:     "this",
	}

	validate := validator.New()
	err := validate.Struct(model)
	if err != nil {
		validateErr := errManager.WrapValidationErrors(err)
		details := errManager.UnwrapDetails(validateErr)
		assert.Equal(t, 2, len(details))
	}
}
