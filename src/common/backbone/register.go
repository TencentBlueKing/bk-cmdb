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

package backbone

import (
    "time"
    "fmt"

    "configcenter/src/common/types"
    regd "configcenter/src/common/RegisterDiscover"
    "github.com/gin-gonic/gin/json"
    "configcenter/src/framework/core/errors"
)

type ServiceDiscoverInterface interface {
    // Ping to check if this service discovery service is health.
    Ping() error
    
    // stop the service discover service
    Stop() error
    
    // register local server info, it can only be called for once.
    Register(path string, c types.ServerInfo) error
}

func NewServcieDiscovery(zkAddr string)(ServiceDiscoverInterface, error) {
    s := new(serviceDiscovery)
    s.client = regd.NewRegDiscoverEx(zkAddr, 5 * time.Second)
    if err := s.client.Start(); nil != err {
        return nil, fmt.Errorf("start service discovery failed, err: %v", err)
    }
    return s, nil
}

type serviceDiscovery struct {
    client *regd.RegDiscover
}

func (s *serviceDiscovery) Ping() error {
    return s.client.Ping()
}

func (s *serviceDiscovery) Stop() error {
    return s.client.Stop()
}

func (s *serviceDiscovery) Register(path string, c types.ServerInfo) error {
    if c.IP == "0.0.0.0" {
        return errors.New("register ip can not be 0.0.0.0")
    }
    
    js, err := json.Marshal(c)
    if err != nil {
        return err
    }
    
    return s.client.RegisterAndWatchService(path ,js)    
}






