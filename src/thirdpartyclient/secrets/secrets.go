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
package secrets

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"

	"github.com/prometheus/client_golang/prometheus"
)

type SecretsClient interface {
	// GetCloudAccountSecretKey get cloud account secret key
	GetCloudAccountSecretKey(ctx context.Context, h http.Header) (string, error)
}

// NewSecretsClient new a secrets client
func NewSecretsClient(tls *util.TLSClientConfig, config SecretsConfig, reg prometheus.Registerer) (SecretsClient, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	client, err := util.NewClient(tls)
	if err != nil {
		return nil, err
	}

	c := &util.Capability{
		Client: client,
		Discover: &scDiscovery{
			servers: []string{config.SecretsAddrs},
		},
		Throttle: flowctrl.NewRateLimiter(1000, 1000),
		Mock: util.MockInfo{
			Mocked: false,
		},
		Reg: reg,
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set(common.BKHTTPSecretsToken, config.SecretsToken)
	header.Set(common.BKHTTPSecretsProject, config.SecretsProject)
	header.Set(common.BKHTTPSecretsEnv, config.SecretsEnv)

	return &secretsClient{
		client:      rest.NewRESTClient(c, "/"),
		config:      config,
		basicHeader: header,
	}, nil

}

// secretsClient implement the interface SecretsClient
type secretsClient struct {
	client rest.ClientInterface
	// secrets client config
	config SecretsConfig
	// http header info
	basicHeader http.Header
}

// SecretsConfig the config of secrets client
type SecretsConfig struct {
	// SecretKeyUrl, the url to get secret_key which used to encrypt and decrypt cloud account
	SecretKeyUrl string
	//SecretsAddrs, the addrs of bk-secrets service, start with http:// or https://
	SecretsAddrs string
	// SecretsToken , as a header param for sending the api request to bk-secrets service
	SecretsToken string
	// SecretsProject, as a header param for sending the api request to bk-secrets service
	SecretsProject string
	// SecretsEnv, as a header param for sending the api request to bk-secrets service
	SecretsEnv string
}

// Validate validate the secrets config fields
func (c *SecretsConfig) Validate() error {
	if c.SecretKeyUrl == "" {
		return errors.New("SecretKeyUrl can't be empty")
	}

	if c.SecretsAddrs == "" {
		return errors.New("SecretsAddrs can't be empty")
	}

	if c.SecretsToken == "" {
		return errors.New("SecretsToken can't be empty")
	}

	if c.SecretsProject == "" {
		return errors.New("SecretsProject can't be empty")
	}

	if c.SecretsEnv == "" {
		return errors.New("SecretsEnv can't be empty")
	}

	return nil
}

// scDiscovery implement the servcie discovery inferface
type scDiscovery struct {
	// servers address, must prefixed with http:// or https://
	servers []string
	index   int
	sync.Mutex
}

func (s *scDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
}

func (s *scDiscovery) GetServersChan() chan []string {
	return nil
}
