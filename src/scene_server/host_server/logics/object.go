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

package logics

import (
	"context"
	"fmt"
	"net/http"

	meta "configcenter/src/common/metadata"
	"configcenter/src/scene_server/host_server/service"
)

// get the object attributes
func (lgc *Logics) GetObjectAttributes(ownerID, objID string, pheader http.Header, page meta.BasePage) ([]meta.Attribute, error) {
	opt := service.NewOperation().WithOwnerID(ownerID).WithPage(page).WithObjID(objID).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get object attribute failed, err: %v, %v", err, result.ErrMsg)
	}

	return result.Data, nil
}
