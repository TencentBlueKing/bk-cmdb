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

package match

import (
	"github.com/emicklei/go-restful"
)

type service struct {
	name  string
	match Match
}

var (
	relations = make([]service, 0)
)

const (
	RootPath = "/api/v3"
)

// Match 判断请求是否属于当前服务
// from 来源前缀
// to  目标地址前缀
type Match func(req *restful.Request) (from, to string, isHit bool)

// RegisterService 注册一个服务，name 服务的名字，也是服务发现需要使用的， match函数判断请求是否属于当前服务
func RegisterService(name string, match Match) {
	relations = append(relations, service{name: name, match: match})
}

// FilterMatch 从服务列表中挑选出来对应的服务项目
func FilterMatch(req *restful.Request) (name string, isHit bool) {
	for _, item := range relations {
		if from, to, isHit := item.match(req); isHit {
			// 将来源地址替换成目标地址
			revise(req, from, to)
			return item.name, isHit
		}
	}

	return
}

func revise(req *restful.Request, from, to string) {
	req.Request.RequestURI = to + req.Request.RequestURI[len(from):]
	req.Request.URL.Path = to + req.Request.URL.Path[len(from):]
}
