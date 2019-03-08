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

package auditoplog

// AuditOpType  log type
type AuditOpType int

const (
	// AuditOpTypeAdd  object add log
	AuditOpTypeAdd AuditOpType = iota + 1
	// AuditOpTypeModify object add log
	AuditOpTypeModify AuditOpType = 2
	// AuditOpTypeDel object del log
	AuditOpTypeDel AuditOpType = 3
	// AuditOpTypeHostModule host  change module
	AuditOpTypeHostModule AuditOpType = 100
)

//操作类型代码分两部分， 前2位表示大类入，后1位表示操作类型，1增加，2.修改，3，删除， 列入100
/*const (
	//主机操作
	AuditOpTypeHostAdd AuditOpType = iota * 101
	AuditOpTypeHostModify
	AuditOpTypeHostDel

	//模块操作
	AuditOpTypeModuleAdd = 201
	AuditOpTypeModuleModify
	AuditOpTypeModuleDel

	//集群操作
	AuditOpTypeSetAdd = 301
	AuditOpTypeSetModify
	AuditOpTypeSetDel
	//业务操作
	AuditOpTypeAppAdd = 401
	AuditOpTypeAppModify
	AuditOpTypeAppDel
	//通用对象操作
	AuditOpTypeObjAdd = 501
	AuditOpTypeObjModify
	AuditOpTypeObjDel

	//进程操作
	AuditOpTypeProcAdd = 601
	AuditOpTypeProcModify
	AuditOpTypeProcDel
)*/

type AuditLogHosts struct {
	ID      int64
	Content interface{}
	InnerIP string
}

type AuditLogExt struct {
	ID      int64 //操作实例id
	Content interface{}
	ExtKey  string
}

type AuditLogContext struct {
	ID      int64 //操作实例id
	Content interface{}
}
