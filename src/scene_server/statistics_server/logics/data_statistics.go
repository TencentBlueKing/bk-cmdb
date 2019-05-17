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
	"context"
	"encoding/json"
	"gopkg.in/redis.v5"
	"net/http"
	"time"

	"configcenter/src/auth/extensions"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/statistics_server/service"
)

type Logics struct {
	*backbone.Engine
	header      http.Header
	rid         string
	ccErr       errors.DefaultCCErrorIf
	ccLang      language.DefaultCCLanguageIf
	user        string
	ownerID     string
	cache       *redis.Client
	AuthManager *extensions.AuthManager
}

func (lgc *Logics) TimerUpdateData(ctx context.Context) {
	go func() {
		timer := time.NewTicker(12 * time.Hour)
		for range timer.C {
			lgc.HostDataStatistics(ctx)
			lgc.InstDataStatistics(ctx)
		}
	}()
}

func (lgc *Logics) HostDataStatistics(ctx context.Context) {
	exist, err := lgc.cache.Exists(service.HostInfoStatistics).Result()
	if err != nil {
		blog.Errorf("check redis key exists fail, err: %v", err)
	}

	if !exist {
		if err := lgc.InitHostCountData(ctx); err != nil {
			blog.Errorf("init host statistics_server data fail, er: %v", err)
		}
		return
	}

	hostInfoStatistics, err := lgc.cache.LRange(service.HostInfoStatistics, 0, -1).Result()
	if err != nil {
		blog.Errorf("get host data from redis fail, error: %v", err)
		return
	}

	if len(hostInfoStatistics) == 0 {
		if err := lgc.InitHostCountData(ctx); err != nil {
			blog.Errorf("init host statistics_server data fail, er: %v", err)
		}
		return
	}

	if err := lgc.UpdateHostCountData(ctx, hostInfoStatistics); err != nil {
		blog.Errorf("update host statistics_server data fail, er: %v", err)
	}

	return

}

func (lgc *Logics) InitHostCountData(ctx context.Context) error {
	lgc.CoreAPI.TopoServer().Instance().SearchApp()
}

func (lgc *Logics) UpdateHostCountData(ctx context.Context, hostInfoStatistics []string) error {
	item := hostInfoStatistics[1]
	hostData := new(metadata.HostStatisticsData)
	if err := json.Unmarshal([]byte(item), hostData); err != nil {
		blog.Errorf("unmarshal host data fail, err: %v", err)
	}

	countFields := hostData.StatisticsFields
	for _, field := range countFields {

	}
}

func (lgc *Logics) InstDataStatistics(ctx context.Context) {
	instInfoStatistics, err := lgc.cache.LRange(service.InstInfoStatistics, 0, -1).Result()
	if err != nil {
		blog.Errorf("get host data from redis fail, error: %v", err)
		return
	}

	if len(instInfoStatistics) == 0 {
		if err := lgc.InitInsTCountData(ctx); err != nil {
			blog.Errorf("init host statistics_server data fail, er: %v", err)
		}
		return
	}

	if err := lgc.UpdateInstCountData(ctx, instInfoStatistics); err != nil {
		blog.Errorf("update host statistics_server data fail, er: %v", err)
	}

	return
}

func (lgc *Logics) InitInsTCountData(ctx context.Context) error {

}

func (lgc *Logics) UpdateInstCountData(ctx context.Context, instInfoStatistics []string) error {

}
