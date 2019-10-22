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

package backbone

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"
)

func NewMockBackbone(c *Config) (*Engine, error) {
	engine := &Engine{
		CoreAPI:  apimachinery.NewMockApiMachinery(),
		SvcDisc:  &mockDisc{},
		Language: language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:    errors.NewFromCtx(errors.EmptyErrorsSetting),
	}

	return engine, nil
}

type mockDisc struct{}

func (*mockDisc) Ping() error {
	return nil
}

func (*mockDisc) Stop() error {
	return nil
}

func (*mockDisc) Register(path string, c types.ServerInfo) error {
	return nil
}

func (*mockDisc) Cancel() {

}

func (*mockDisc) ClearRegisterPath() error {
	return nil
}
