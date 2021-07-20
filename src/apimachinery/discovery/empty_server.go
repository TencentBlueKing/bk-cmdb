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

package discovery

var (
	emptyServerInst = &emptyServer{}
)

// emptyServer 适配服务不存在的情况， 当服务不存在的时候，返回空的服务
type emptyServer struct {
}

func (es *emptyServer) GetServers() ([]string, error) {
	return []string{}, nil
}

func (es *emptyServer) GetServersChan() chan []string {
	return make(chan []string, 20)
}
