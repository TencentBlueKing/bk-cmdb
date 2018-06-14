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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	meta "configcenter/src/common/metadata"
	"configcenter/src/scene_server/host_server/service"
	"configcenter/src/source_controller/api/metadata"
)

type Logics struct {
	*backbone.Engine
}

func (lgc *Logics) GetHostAttributes(ownerID string, header http.Header) ([]metadata.Header, error) {
	searchOp := service.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), header, searchOp)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	headers := make([]metadata.Header, 0)
	for _, p := range result.Data {
		if p.PropertyID == common.BKChildStr {
			continue
		}
		headers = append(headers, metadata.Header{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
		})
	}

	return headers, nil
}

func (lgc *Logics) GenerateHostLogs(ownerID string, hostID string, logHeaders []metadata.Header, pheader http.Header) (*metadata.Content, error) {
	ctnt := new(metadata.Content)
	ctnt.Headers = logHeaders

	// get host details, pre data
	result, err := lgc.CoreAPI.HostController().Host().GetHostByID(context.Background(), hostID, pheader)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get host pre data failed, err, %v, %v", err, result.ErrMsg)
	}

	hostInfo, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid host info data")
	}

	attributes, err := lgc.GetObjectAsst(ownerID, pheader)
	if err != nil {
		return nil, err
	}

	for key, val := range attributes {
		if item, ok := hostInfo[key]; ok {
			if item == nil {
				continue
			}

			strItem, ok := item.(string)
			if !ok {
				return nil, errors.New("invalid parameter")
			}
			ids := make([]int64, 0)
			for _, strID := range strings.Split(strItem, ",") {
				id, err := strconv.ParseInt(strID, 10, 64)
				if err != nil {
					return nil, err
				}
				ids = append(ids, id)
			}

			cond := make(map[string]interface{})
			cond[common.BKHostIDField] = map[string]interface{}{"$in": ids}
			q := meta.QueryInput{
				Start:     0,
				Limit:     common.BKNoLimit,
				Sort:      "",
				Condition: cond,
			}

		}
	}

}
