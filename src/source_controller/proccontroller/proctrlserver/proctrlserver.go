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

package proctrlserver

import (
    "configcenter/src/common/backbone"
    "github.com/emicklei/go-restful"
    "net/http"
    "configcenter/src/storage"
    "configcenter/src/storage/mgoclient"
    "configcenter/src/storage/redisclient"
    cfnc "configcenter/src/common/backbone/configcenter"
)

type ProctrlServer struct {
    Core *backbone.Engine
    DbInstance storage.DI
    CacheDI storage.DI
    MongoCfg mgoclient.MongoConfig
    RedisCfg redisclient.RedisConfig
}

func (ps *ProctrlServer) WebService(filter restful.FilterFunction) http.Handler {
    
    container := new(restful.Container)
    // v3
    v3WS := new(restful.WebService)
    v3WS.Path("/process/v3").Filter(filter).Produces(restful.MIME_JSON)
    
    v3WS.Route(v3WS.DELETE("/module").To(ps.DeleteProc2Module))
    v3WS.Route(v3WS.POST("/module").To(ps.CreateProc2Module))
    v3WS.Route(v3WS.POST("/module/search").To(ps.GetProc2Module))
    
    container.Add(v3WS)
    
    return container
}

func (ps *ProctrlServer) OnProcessConfUpdate(previous, current cfnc.ProcessConfig) {
    prefix := storage.DI_MONGO
    ps.MongoCfg = mgoclient.MongoConfig{
        Address:      current.ConfigMap[prefix+".host"],
        User:         current.ConfigMap[prefix+".user"],
        Password:     current.ConfigMap[prefix+".pwd"],
        Database:     current.ConfigMap[prefix+".database"],
        Port:         current.ConfigMap[prefix+".port"],
        MaxOpenConns: current.ConfigMap[prefix+".maxOpenConns"],
        MaxIdleConns: current.ConfigMap[prefix+".maxIDleConns"],
        Mechanism:    current.ConfigMap[prefix+".mechanism"],
    }
    
    prefix = storage.DI_REDIS
    ps.RedisCfg = redisclient.RedisConfig{
        Address:  current.ConfigMap[prefix+".host"],
        User:     current.ConfigMap[prefix+".user"],
        Password: current.ConfigMap[prefix+".pwd"],
        Database: current.ConfigMap[prefix+".database"],
        Port:     current.ConfigMap[prefix+".port"],
    }
}

