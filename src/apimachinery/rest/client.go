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

package rest

import (
	"strings"
	"time"

	"configcenter/src/apimachinery/util"
	"github.com/prometheus/client_golang/prometheus"
)

type ClientInterface interface {
	Verb(verb VerbType) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
	Patch() *Request
}

func NewRESTClient(c *util.Capability, baseUrl string) ClientInterface {
	if baseUrl != "/" {
		baseUrl = strings.Trim(baseUrl, "/")
		baseUrl = "/" + baseUrl + "/"
	}

	if c.ToleranceLatencyTime <= 0 {
		// set default tolerance latency time
		c.ToleranceLatencyTime = 2 * time.Second
	}

	client := &RESTClient{
		baseUrl:    baseUrl,
		capability: c,
	}

	if c.MetricOpts.Register != nil {

		var buckets []float64
		if len(c.MetricOpts.DurationBuckets) == 0 {
			// set default buckets
			buckets = []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000}
		} else {
			// use user defined buckets
			buckets = c.MetricOpts.DurationBuckets
		}

		client.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cmdb_apimachinary_requests_duration_millisecond",
			Help:    "third party api request duration millisecond.",
			Buckets: buckets,
		}, []string{"handler", "status_code", "dimension"})

		if err := c.MetricOpts.Register.Register(client.requestDuration); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				client.requestDuration = are.ExistingCollector.(*prometheus.HistogramVec)
			} else {
				panic(err)
			}
		}
	}

	return client
}

type RESTClient struct {
	baseUrl    string
	capability *util.Capability

	requestDuration *prometheus.HistogramVec
}

func (r *RESTClient) Verb(verb VerbType) *Request {
	return &Request{
		parent:     r,
		verb:       verb,
		baseURL:    r.baseUrl,
		capability: r.capability,
	}
}

func (r *RESTClient) Post() *Request {
	return r.Verb(POST)
}

func (r *RESTClient) Put() *Request {
	return r.Verb(PUT)
}

func (r *RESTClient) Get() *Request {
	return r.Verb(GET)
}

func (r *RESTClient) Delete() *Request {
	return r.Verb(DELETE)
}

func (r *RESTClient) Patch() *Request {
	return r.Verb(PATCH)
}
