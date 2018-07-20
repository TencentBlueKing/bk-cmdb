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
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	httpClient "configcenter/src/source_controller/api/client"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient.Client
}

// NewClient 创建审计日志操作接口
// log OpDesc The multi-language version is used as the translation key.
// There is no corresponding value in the language pack. The current content is displayed.
func NewClient(address string, header http.Header) *Client {

	cli := &Client{}
	cli.SetAddress(address)

	cli.Base = base.BaseLogic{}
	cli.Base.CreateHttpClient()
	for key := range header {
		cli.Base.HttpCli.SetHeader(key, header.Get(key))
	}

	return cli
}

var bk_inst_id_fields string = "inst_id"

// AuditHostLog  新加主机操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditHostLog(id interface{}, Content interface{}, OpDesc string, InnerIP, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKHostInnerIPField: InnerIP, common.BKOpTypeField: OpType, bk_inst_id_fields: id}
	url := fmt.Sprintf("%s/audit/v1/host/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)

	return cli.GetRequestInfo(common.HTTPCreate, data, url)

}

// AuditHostsLog  批量新加主机操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditHostsLog(Content []auditoplog.AuditLogExt, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType}
	url := fmt.Sprintf("%s/audit/v1/hosts/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)

}

// AuditAppLog  新加业务操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditAppLog(id interface{}, Content interface{}, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType, bk_inst_id_fields: id}
	url := fmt.Sprintf("%s/audit/v1/app/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditSetLog  新加集群操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditSetLog(id interface{}, Content interface{}, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType, bk_inst_id_fields: id}
	url := fmt.Sprintf("%s/audit/v1/set/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditSetsLog  批量新加集群操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditSetsLog(Content []auditoplog.AuditLogContext, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType}
	url := fmt.Sprintf("%s/audit/v1/sets/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditModuleLog  新加模块操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditModuleLog(id interface{}, Content interface{}, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType, bk_inst_id_fields: id}
	url := fmt.Sprintf("%s/audit/v1/module/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditModulesLog  批量新加模块操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditModulesLog(Content []auditoplog.AuditLogContext, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType}
	url := fmt.Sprintf("%s/audit/v1/modules/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditObjLog  新加对象操作日志
// OpDesc The multi-language version is used as the translation key.
// There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditObjLog(id interface{}, Content interface{}, OpDesc, opTarget string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType, common.BKOpTargetField: opTarget, bk_inst_id_fields: id}
	blog.InfoJSON("add audit log: %s", data)
	url := fmt.Sprintf("%s/audit/v1/obj/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditObjsLog  批量新加对象操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditObjsLog(Content []auditoplog.AuditLogContext, OpDesc, opTarget string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType, common.BKOpTargetField: opTarget}
	url := fmt.Sprintf("%s/audit/v1/objs/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditProcLog  新加进程操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditProcLog(id interface{}, Content interface{}, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType, bk_inst_id_fields: id}
	blog.InfoJSON("add audit log: %s", data)
	url := fmt.Sprintf("%s/audit/v1/proc/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}

// AuditProcsLog  批量新加进程操作日志 OpDesc The multi-language version is used as the translation key. There is no corresponding value in the language pack. The current content is displayed.
func (cli *Client) AuditProcsLog(Content []auditoplog.AuditLogContext, OpDesc string, ownerID, appID, user string, OpType auditoplog.AuditOpType) (interface{}, error) {
	data := common.KvMap{common.BKContentField: Content, common.BKOpDescField: OpDesc, common.BKOpTypeField: OpType}
	url := fmt.Sprintf("%s/audit/v1/procs/%s/%s/%s", cli.GetAddress(), ownerID, appID, user)
	return cli.GetRequestInfo(common.HTTPCreate, data, url)
}
