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
	"fmt"
	"github.com/rs/xid"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

func getInlineProcInstKey(hostID, moduleID int64) string {
	return fmt.Sprintf("%d-%d", hostID, moduleID)
}

func getGseProcNameSpace(appID, moduleID int64) string {
	return fmt.Sprintf("%d.%d", appID, moduleID)
}

func getTaskID() string {
	return fmt.Sprintf("cc:task:gse:%s:%s", common.BKSTRIDPrefix, xid.New().String())
}

func getGseOpInstKey(moduleID, procID int64) string {
	return fmt.Sprintf("%d-%d", moduleID, procID)
}

func GetProcInstModel(appID, setID, moduleID, hostID, procID, funcID, procNum int64, maxInstID uint64) []*metadata.ProcInstanceModel {
	if 0 >= procNum {
		procNum = 1
	}
	instProc := make([]*metadata.ProcInstanceModel, 0)
	for numIdx := int64(1); numIdx < procNum+1; numIdx++ {
		procIdx := (maxInstID-1)*uint64(procNum) + uint64(numIdx)
		item := new(metadata.ProcInstanceModel)
		item.ApplicationID = appID
		item.SetID = setID
		item.ModuleID = moduleID
		item.FuncID = funcID
		item.HostID = hostID
		item.HostInstanID = maxInstID
		item.ProcInstanceID = procIdx
		item.ProcID = procID
		item.HostProcID = uint64(numIdx)
		instProc = append(instProc, item)
	}
	return instProc
}
