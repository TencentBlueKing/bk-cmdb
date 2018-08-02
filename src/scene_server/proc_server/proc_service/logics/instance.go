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
	"configcenter/src/common/metadata"
)

func (lgc *Logics) RefreshProcInstance(eventData *metadata.EventInst) {
	switch eventData.ObjType {
	case metadata.EventObjTypeProcModule:

	}
}

func (lgc *Logics) refreshProcInstByProcModule(eventData *metadata.EventInst) {
	if metadata.EventTypeRelation != eventData.EventType {
		return
	}
	if metadata.EventActionDelete == eventData.Action {
		// delete process bind module relation, unregister process info
	} else {
		// compare  pre-change data with the current data and find the newly added data to register process info

	}

}

func (lgc *Logics) refreshProcInstByProcess(eventData *metadata.EventInst) {
	switch eventData.Action {
	case metadata.EventActionCreate:
		// create proccess not refresh process instance , because not bind module
	case metadata.EventActionUpdate:
		// refresh process instance register again
	case metadata.EventActionDelete:
		// delete all register process
	}
}

func (lgc *Logics) refreshProcInstByHostInfo(eventData *metadata.EventInst) {
	if metadata.EventTypeRelation == eventData.EventType {
		if metadata.EventActionDelete == eventData.Action {
			// delete host from module , unregister module bind all process info
		} else {
			// compare pre-change data with the current data and find the newly added data to register process info
		}
	} else {
		// host fields supperid, cloud id, innerip   change , register process info  agin
	}
}

func (lgc *Logics) refreshProcInstModuleByHostInfo(eventData *metadata.EventInst) {
	if metadata.EventActionUpdate == eventData.EventType {
		// module change name, unregister pre-change module name bind process , then register current module name bind process
	}
}
