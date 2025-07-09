/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package errors

// EmptyErrorsSetting empty errors setting
var EmptyErrorsSetting = map[string]ErrorCode{}

// NewFactory create new CCErrorIf instance,
// dir is directory of errors description resource
func NewFactory(dir string) (CCErrorIf, error) {

	tmp := &ccErrorHelper{errCode: make(map[string]ErrorCode)}

	errcode, err := LoadErrorResourceFromDir(dir)
	if err != nil {
		// blog.Errorf("failed to load the error resource, error info is %s", err.Error())
		return nil, err
	}
	tmp.Load(errcode)

	return tmp, nil
}

// NewFromCtx TODO
func NewFromCtx(errcode map[string]ErrorCode) CCErrorIf {
	tmp := &ccErrorHelper{}
	tmp.Load(errcode)
	return tmp
}

// New 根据response返回的信息产生错误
func New(errCode int, errMsg string) CCErrorCoder {
	return &ccError{code: errCode, callback: func() string {
		return errMsg
	}}
}
