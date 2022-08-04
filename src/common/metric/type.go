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

package metric

import (
	"errors"
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/common/version"
)

// MetricFamily TODO
type MetricFamily struct {
	MetaData     *MetaData                   `json:"metaData"`
	MetricBundle map[CollectorName][]*Metric `json:"metricBundle"`
	ReportTimeMs int64                       `json:"reportTimeMs"`
}

// Metric TODO
type Metric struct {
	*MetricMeta `json:",inline"`
	Value       *FloatOrString   `json:"value"`
	Extension   *MetricExtension `json:"extension"`
}

func newMetric(m MetricInterf) (*Metric, error) {
	if m == nil {
		return nil, errors.New("metric is nil")
	}
	meta := m.GetMeta()
	if len(meta.Name) == 0 {
		return nil, errors.New("metric name is null")
	}

	if len(meta.Help) == 0 {
		return nil, errors.New("metric help is null")
	}

	val, err := m.GetValue()
	if nil != err {
		return nil, err
	}
	if nil == val {
		return nil, errors.New("metric value is nil")
	}

	extension, err := m.GetExtension()
	if nil != err {
		return nil, err
	}
	return &Metric{
		MetricMeta: meta,
		Value:      val,
		Extension:  extension,
	}, nil
}

// CollectorName TODO
type CollectorName string

// Collector TODO
type Collector struct {
	Name      CollectorName
	Collector CollectInter
}

// MetaData TODO
type MetaData struct {
	Module        string            `json:"module"`
	ServerAddress string            `json:"server_address"`
	ClusterID     string            `json:"clusterID"`
	Labels        map[string]string `json:"label"`
}

// HealthResponse TODO
type HealthResponse struct {
	Code    int        `json:"code"`
	OK      bool       `json:"ok"`
	Message string     `json:"message"`
	Data    HealthInfo `json:"data"`
	Result  bool       `json:"result"`
}

// SetCommonResponse TODO
func (h *HealthResponse) SetCommonResponse() {
	// set version in healthz response
	h.Data.Version = map[string]interface{}{
		"version":   version.CCVersion,
		"time":      version.CCBuildTime,
		"commit_id": version.CCGitHash,
	}
}

// HealthInfo TODO
type HealthInfo struct {
	Module     string `json:"module"`
	Address    string `json:"address"`
	HealthMeta `json:",inline"`
	AtTime     metadata.Time          `json:"at_time"`
	Version    map[string]interface{} `json:"version"`
}

// VersionInfo TODO
type VersionInfo struct {
	Module    string `json:"module"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	CommitID  string `json:"commit_id"`
}

// Action TODO
type Action struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}
