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
 
package auditlog

import (
	"configcenter/src/common"
	"configcenter/src/source_controller/common/commondata"
	"fmt"
)

//AuditHostsLog  批量新加主机操作日志
func (cli *Client) GetAuditlogs(input commondata.ObjQueryInput) (interface{}, error) {
	url := fmt.Sprintf("%s/audit/v1/search", cli.GetAddress())
	return cli.GetRequestInfo(common.HTTPSelectPost, input, url)

}
