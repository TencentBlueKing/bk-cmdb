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

package service

import (
	"encoding/json"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

func (s *Service) GetBizModuleHostAmount(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AddStatisticalChart(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	chartInfo := new(metadata.StatisticChartInfo)
	if err := json.NewDecoder(req.Request.Body).Decode(chartInfo); err != nil {
		blog.Errorf("add statistics chart failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

}

func (s *Service) DeleteStatisticalChart(req *restful.Request, resp *restful.Response) {

}

func (s *Service) SearchStatisticalCharts(req *restful.Request, resp *restful.Response) {

}

func (s *Service) UpdateStatisticalChart(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetChartData(req *restful.Request, resp *restful.Response) {

}
