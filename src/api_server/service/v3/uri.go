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

package v3

import (
	"errors"
	"strings"

	"github.com/emicklei/go-restful"
)

type RequestType string

const (
	UnknownType RequestType = "unknown"
	TopoType    RequestType = "topo"
	HostType    RequestType = "host"
	ProcType    RequestType = "proc"
	EventType   RequestType = "event"
)

type V3URLPath string

func (u V3URLPath) FilterChain(req *restful.Request) (RequestType, error) {
	switch {
	case u.WithTopo(req):
		return TopoType, nil
	case u.WithHost(req):
		return HostType, nil
	case u.WithProc(req):
		return ProcType, nil
	case u.WithEvent(req):
		return EventType, nil
	default:
		return UnknownType, errors.New("unknown requested with backend process")
	}
}

func (u *V3URLPath) WithTopo(req *restful.Request) (isHit bool) {
	topoRoot := "/topo/v3"
	from, to := rootPath, topoRoot
	switch {
	case strings.HasPrefix(string(*u), rootPath+"/audit/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/biz/"):
		from, to, isHit = rootPath+"/biz", topoRoot+"/app", true

	case strings.HasPrefix(string(*u), rootPath+"/topo/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/identifier/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/inst/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/module/"):
		from, to, isHit = rootPath, topoRoot, true

		// Attention:
		// do not change the check sequences.
	case string(*u) == rootPath+"/object":
		from, to, isHit = rootPath, topoRoot, true

	case string(*u) == rootPath+"/objects":
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/object/attr"):
		from, to, isHit = rootPath+"/object/attr", topoRoot+"/objectattr", true

	case strings.HasPrefix(string(*u), rootPath+"/object/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/objects/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/objectatt/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/set/"):
		from, to, isHit = rootPath, topoRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

func (u *V3URLPath) WithHost(req *restful.Request) (isHit bool) {
	hostRoot := "/host/v3"
	from, to := rootPath, hostRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/host/"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/hosts/"):
		from, to, isHit = rootPath, hostRoot, true

	case string(*u) == (rootPath + "/userapi"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/userapi/"):
		from, to, isHit = rootPath, hostRoot, true

	case string(*u) == (rootPath + "/usercustom"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/usercustom/"):
		from, to, isHit = rootPath, hostRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

func (u *V3URLPath) WithEvent(req *restful.Request) (isHit bool) {
	eventRoot := "/event/v3"
	from, to := rootPath, eventRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/event/"):
		from, to, isHit = rootPath+"/event", eventRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

func (u *V3URLPath) WithProc(req *restful.Request) (isHit bool) {
	procRoot := "/process/v3"
	from, to := rootPath, procRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/proc/"):
		from, to, isHit = rootPath+"/proc", procRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

func (u V3URLPath) revise(req *restful.Request, from, to string) {
	req.Request.RequestURI = to + req.Request.RequestURI[len(from):]
	req.Request.URL.Path = to + req.Request.URL.Path[len(from):]
}
