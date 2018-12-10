/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package errors

// EmptyErrorsSetting empty errors setting
var EmptyErrorsSetting = map[string]ErrorCode{}

// New create new CCErrorIf instance,
// dir is directory of errors description resource
func New(dir string) (CCErrorIf, error) {

	tmp := &ccErrorHelper{errCode: make(map[string]ErrorCode)}

	errcode, err := LoadErrorResourceFromDir(dir)
	if nil != err {
		//blog.Error("failed to load the error resource, error info is %s", err.Error())
		return nil, err
	}
	tmp.Load(errcode)

	return tmp, nil
}

func NewFromCtx(errcode map[string]ErrorCode) CCErrorIf {
	tmp := &ccErrorHelper{}
	tmp.Load(errcode)
	return tmp
}
