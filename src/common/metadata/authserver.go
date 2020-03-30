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

package metadata

type PullResourceParam struct {
	Collection string                 `json:"collection"`
	Condition  map[string]interface{} `json:"condition"`
	Fields     []string               `json:"fields"`
	Limit      int64                  `json:"limit"`
	Offset     int64                  `json:"offset"`
}

type PullResourceResponse struct {
	BaseResp `json:",inline"`
	Data     PullResourceResult `json:"data"`
}

type PullResourceResult struct {
	Count int64                    `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}
