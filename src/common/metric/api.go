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
	"encoding/json"
	"errors"
	"fmt"

	"configcenter/src/common/http/httpclient"
)

func NewMetricController(conf Config, healthFunc HealthFunc, collectors ...*Collector) []Action {
	return newMetricController(conf, healthFunc, collectors...)
}

type RunModeType string

// used when your module running with Master_Slave_Mode mode
type RoleType string

// metric const define
const (
	MetricPort = 60060
)

// Config define metric's define
type Config struct {
	// name of your module
	ModuleName string
	// server address
	ServerAddress string
	// self defined info labeled on your metrics.
	Labels map[string]string
	// metric http server's ssl configuration
	SvrCaFile   string
	SvrCertFile string
	SvrKeyFile  string
	CertPasswd  string
}

// HealthFunc returns HealthMeta
type HealthFunc func() HealthMeta

// HealthMeta define the HealthMeta that shows whether this server healthy
type HealthMeta struct {
	// if this module is healthy
	IsHealthy bool `json:"healthy"`
	// messages which describes the health status
	Message string `json:"message"`

	Items []HealthItem `json:"items"`
}

// HealthItem define
type HealthItem struct {
	// item name
	Name string `json:"name"`
	// if this module is healthy
	IsHealthy bool `json:"healthy"`
	// messages which describes the health status
	Message string `json:"message"`
}

// MetricMeta define the MetricMeta that shows the named metric
type MetricMeta struct {
	// metric's name
	Name string `json:"name"`
	// metric's help info, which should be short and briefly.
	Help string `json:"help"`
}

type MetricInterf interface {
	GetMeta() *MetricMeta
	GetValue() (*FloatOrString, error)
	GetExtension() (*MetricExtension, error)
}

type MetricExtension struct{}

type CollectInter interface {
	Collect() []MetricInterf
}

func NewCollector(name string, collector CollectInter) *Collector {
	return &Collector{
		Name:      CollectorName(name),
		Collector: collector,
	}
}

func CheckHealthy(address string) error {
	if "" == address {
		return errors.New("address not found")
	}
	out, err := httpclient.NewHttpClient().GET(address+"/healthz", nil, nil)
	if err != nil {
		return err
	}
	resp := HealthResponse{}
	err = json.Unmarshal(out, &resp)
	if err != nil {
		fmt.Printf("healthz return %s", out)
		return err
	}
	if !resp.Result {
		return errors.New(resp.Message)
	}
	return nil
}

// NewHealthItem build the HealthItem depend on checkHealthFuc return
func NewHealthItem(name string, err error) HealthItem {
	mongoHealthy := HealthItem{Name: name}
	if err != nil {
		mongoHealthy.IsHealthy = false
		mongoHealthy.Message = err.Error()
	} else {
		mongoHealthy.IsHealthy = true
	}
	return mongoHealthy
}
