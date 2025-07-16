/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package audit

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/audit/config"

	"go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type auditClient struct {
	client rest.ClientInterface
	token  string
}

// ReportAuditData report audit data to audit center
func (a *auditClient) ReportAuditData(ctx context.Context, h http.Header, opt *v1.ExportLogsServiceRequest) error {
	h.Set("Content-Type", "application/json")
	h.Set("X-BK-TOKEN", a.token)

	body, err := protojson.Marshal(opt)
	if err != nil {
		return fmt.Errorf("marshal audit request failed, err: %v", err)
	}

	res := a.client.Post().
		WithContext(ctx).
		Body(body).
		SubResourcef("/v1/logs").
		WithHeaders(h).
		Do()

	if res.StatusCode >= 300 || res.Err != nil {
		return fmt.Errorf("http request failed, err: %v, status %s", res.Err, res.Status)
	}

	return nil
}

// newAuditClient new audit center client
func newAuditClient(conf *config.Config) (*auditClient, error) {
	client, err := util.NewClient(nil)
	if err != nil {
		blog.Errorf("new http client failed, err: %v", err)
		return nil, err
	}

	c := &util.Capability{
		Client:   client,
		Discover: &discovery{endpoint: conf.Endpoint},
	}

	restCli := rest.NewRESTClient(c, "/")

	return &auditClient{
		client: restCli,
		token:  conf.Token,
	}, nil
}

type discovery struct {
	endpoint string
}

// GetServers get servers
func (s *discovery) GetServers() ([]string, error) {
	return []string{s.endpoint}, nil
}

// GetServersChan get servers chan
func (s *discovery) GetServersChan() chan []string {
	return nil
}
