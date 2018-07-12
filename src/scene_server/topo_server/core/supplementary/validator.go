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

package supplementary

import (
	"net/http"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ValidatorInterface the validator methods
type ValidatorInterface interface {
	ValidatorCreate(header http.Header, obj *metadata.Object, attr *metadata.Attribute, data mapstr.MapStr) (bool, error)
	ValidatorUpdate(header http.Header, obj *metadata.Object, attr *metadata.Attribute, data mapstr.MapStr, instID int64) (bool, error)
}

type validator struct {
}

func (v *validator) ValidatorCreate(header http.Header, obj *metadata.Object, attr *metadata.Attribute, data mapstr.MapStr) (bool, error) {
	return false, nil
}
func (v *validator) ValidatorUpdate(header http.Header, obj *metadata.Object, attr *metadata.Attribute, data mapstr.MapStr, instID int64) (bool, error) {
	return true, nil
}
