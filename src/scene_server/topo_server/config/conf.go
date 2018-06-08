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
 
package config

import (
	confHandler "configcenter/src/common/confcenter"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"
)

// ConfCenter discover configure changed. get, update configures
type ConfCenter interface {
	Start() error
	Stop() error
	GetConfigureCxt() []byte
	GetErrorCxt() map[string]errors.ErrorCode
	GetLanguageResCxt() map[string]language.LanguageMap
}

// NewConfCenter create a ConfCenter object
func NewConfCenter(serv string) ConfCenter {
	return confHandler.NewConfCenter(serv, types.CC_MODULE_TOPO)
}
