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

package middleware

import (
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func GetAuditLogHeader(clientSet apimachinery.ClientSetInterface, header http.Header, objID string) ([]metadata.Header, error) {
	supplierAccount := util.GetOwnerID(header)
	rid := util.GetHTTPCCRequestID(header)
	ctx := util.NewContextFromHTTPHeader(header)

	filter := map[string]interface{}{
		metadata.AttributeFieldObjectID:        objID,
		metadata.AttributeFieldSupplierAccount: supplierAccount,
		metadata.AttributeFieldIsSystem: map[string]interface{}{
			common.BKDBNE: true,
		},
		metadata.AttributeFieldIsAPI: map[string]interface{}{
			common.BKDBNE: true,
		},
	}
	qc := &metadata.QueryCondition{
		Condition: filter,
	}
	rsp, err := clientSet.CoreService().Model().ReadModelAttr(ctx, header, objID, qc)
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), rid)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", objID, rsp.ErrMsg, rid)
		return nil, err
	}

	headers := make([]metadata.Header, 0)
	if nil != err {
		blog.Errorf("[audit]failed to get the object(%s)' attribute, error info is %s, rid: %s", objID, err.Error(), rid)
		return headers, nil
	}
	for _, attr := range rsp.Data.Info {
		headers = append(headers, metadata.Header{
			PropertyID:   attr.PropertyID,
			PropertyName: attr.PropertyName,
		})
	}
	return headers, nil
}

func DeepCopyToMap(v interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}
