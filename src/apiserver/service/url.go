/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"configcenter/src/apiserver/service/match"
	"github.com/emicklei/go-restful"
)

// URLPath url path filter
type URLPath string

// FilterChain url path filter
func (u URLPath) FilterChain(req *restful.Request) (RequestType, error) {
	var serverType RequestType
	var err error

	switch {
	case u.WithCache(req):
		serverType = CacheType
	case u.WithTopo(req):
		serverType = TopoType
	case u.WithHost(req):
		serverType = HostType
	case u.WithProc(req):
		serverType = ProcType
	case u.WithEvent(req):
		serverType = EventType
	case u.WithDataCollect(req):
		return DataCollectType, nil
	case u.WithOperation(req):
		return OperationType, nil
	case u.WithTask(req):
		return TaskType, nil
	case u.WithAdmin(req):
		return AdminType, nil
	case u.WithCloud(req):
		return CloudType, nil
	default:
		if server, isHit := match.FilterMatch(req); isHit {
			return RequestType(server), nil
		}
		serverType = UnknownType
		err = errors.New("unknown requested with backend process")
	}

	return serverType, err
}

var topoURLRegexp = regexp.MustCompile(fmt.Sprintf("^/api/v3/(%s)/(inst|object|objects|topo|biz|module|set|resource)/.*$", verbs))

// WithTopo parse topo api's url
func (u *URLPath) WithTopo(req *restful.Request) (isHit bool) {
	topoRoot := "/topo/v3"
	from, to := rootPath, topoRoot
	switch {
	case strings.HasPrefix(string(*u), rootPath+"/biz/"):
		from, to, isHit = rootPath+"/biz", topoRoot+"/app", true

	case strings.HasPrefix(string(*u), rootPath+"/topo/"):
		from, to, isHit = rootPath, topoRoot, true

	case topoURLRegexp.MatchString(string(*u)):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/identifier/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/inst/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/module/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/object/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/set/"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/objectclassification"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/classificationobject"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/objectattr"):
		from, to, isHit = rootPath, topoRoot, true
	case strings.Contains(string(*u), "/objectunique"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/objectattgroup"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/objectattgroupproperty"):
		from, to, isHit = rootPath, topoRoot, true

	// TODO remove it
	case strings.Contains(string(*u), "/objectattgroupasst"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/objecttopo"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/topomodelmainline"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/topoinst"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/topopath"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/topoassociationtype"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/objectassociation"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/instassociation"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/insttopo"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/instance"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/instassociationdetail"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/associationtype"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/find/full_text"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/find/audit_dict"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.Contains(string(*u), "/findmany/audit_list"):
		from, to, isHit = rootPath, topoRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/find/audit"):
		from, to, isHit = rootPath, topoRoot, true

	case topoURLRegexp.MatchString(string(*u)):
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

// hostCloudAreaURLRegexp host server operator cloud area api regex
var hostCloudAreaURLRegexp = regexp.MustCompile(fmt.Sprintf("^/api/v3/(%s)/(cloudarea|cloudarea/.*)$", verbs))
var hostURLRegexp = regexp.MustCompile(fmt.Sprintf("^/api/v3/(%s)/(host|hosts|host_apply_rule|host_apply_plan)/.*$", verbs))

// WithHost transform the host's url
func (u *URLPath) WithHost(req *restful.Request) (isHit bool) {
	hostRoot := "/host/v3"
	from, to := rootPath, hostRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/host/"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/hosts/"):
		from, to, isHit = rootPath, hostRoot, true

	// dynamic grouping URL matching, and proxy to host server.
	case string(*u) == (rootPath + "/dynamicgroup"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/dynamicgroup/"):
		from, to, isHit = rootPath, hostRoot, true

	case string(*u) == (rootPath + "/usercustom"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/usercustom/"):
		from, to, isHit = rootPath, hostRoot, true

	case hostCloudAreaURLRegexp.MatchString(string(*u)):
		from, to, isHit = rootPath, hostRoot, true

	case hostURLRegexp.MatchString(string(*u)):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/system/config"):
		from, to, isHit = rootPath, hostRoot, true

	case strings.HasPrefix(string(*u), rootPath+"/findmany/module_relation/bk_biz_id/"):
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

// WithEvent transform event's url
func (u *URLPath) WithEvent(req *restful.Request) (isHit bool) {
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

const verbs = "create|createmany|update|updatemany|delete|deletemany|find|findmany"

var procUrlRegexp = regexp.MustCompile(fmt.Sprintf("^/api/v3/(%s)/proc/.*$", verbs))

// WithProc transform the proc's url
func (u *URLPath) WithProc(req *restful.Request) (isHit bool) {
	procRoot := "/process/v3"
	from, to := rootPath, procRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/proc/"):
		from, to, isHit = rootPath+"/proc", procRoot, true
	case procUrlRegexp.MatchString(string(*u)):
		from, to, isHit = rootPath, procRoot, true
	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

// WithDataCollect transform DataCollect's url
func (u *URLPath) WithDataCollect(req *restful.Request) (isHit bool) {
	dataCollectRoot := "/collector/v3"
	from, to := rootPath, dataCollectRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/collector/"):
		from, to, isHit = rootPath+"/collector", dataCollectRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

var operationUrlRegexp = regexp.MustCompile(fmt.Sprintf("^/api/v3/(%s)/operation/.*$", verbs))

// WithOperation transform OperationStatistic's url
func (u *URLPath) WithOperation(req *restful.Request) (isHit bool) {
	statisticsRoot := "/operation/v3"
	from, to := rootPath, statisticsRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/operation/"):
		from, to, isHit = rootPath, statisticsRoot, true
	case operationUrlRegexp.MatchString(string(*u)):
		from, to, isHit = rootPath, statisticsRoot, true
	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

// WithTask transform task server  url
func (u *URLPath) WithTask(req *restful.Request) (isHit bool) {
	statisticsRoot := "/task/v3"
	from, to := rootPath, statisticsRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/task/"):
		from, to, isHit = rootPath, statisticsRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

// WithAdmin transform admin server url
func (u *URLPath) WithAdmin(req *restful.Request) (isHit bool) {
	adminRoot := "/migrate/v3"
	from, to := rootPath, adminRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/admin/"):
		from, to, isHit = rootPath+"/admin", adminRoot, true

	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

var cloudUrlRegexp = regexp.MustCompile(fmt.Sprintf("^/api/v3/(%s)/cloud/.*$", verbs))

// WithCloud transform cloud's url
func (u *URLPath) WithCloud(req *restful.Request) (isHit bool) {
	cloudRoot := "/cloud/v3"
	from, to := rootPath, cloudRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/cloud/"):
		from, to, isHit = rootPath, cloudRoot, true
	case cloudUrlRegexp.MatchString(string(*u)):
		from, to, isHit = rootPath, cloudRoot, true
	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}

func (u URLPath) revise(req *restful.Request, from, to string) {
	req.Request.RequestURI = to + req.Request.RequestURI[len(from):]
	req.Request.URL.Path = to + req.Request.URL.Path[len(from):]
}

// WithCache transform cache service's url
func (u *URLPath) WithCache(req *restful.Request) (isHit bool) {
	cacheRoot := "/cache/v3"
	from, to := rootPath, cacheRoot

	switch {
	case strings.HasPrefix(string(*u), rootPath+"/cache/"):
		from, to, isHit = rootPath+"/cache", cacheRoot, true
	default:
		isHit = false
	}

	if isHit {
		u.revise(req, from, to)
		return true
	}
	return false
}
