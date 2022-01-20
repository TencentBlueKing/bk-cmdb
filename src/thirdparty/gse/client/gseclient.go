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

package client

import (
	"context"
	"sync"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/ssl"
	apiserver "configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"
	taskserver "configcenter/src/thirdparty/gse/push_file_forsyncdata"

	"github.com/apache/thrift/lib/go/thrift"
)

type GseApiServerClient struct {
	clients []*apiserver.CacheAPIClient
	index   int
	sync.Mutex
}

// NewGseApiServerClient new gse api server client
func NewGseApiServerClient(endpoints []string, conf *util.TLSClientConfig) (*GseApiServerClient, error) {
	var clients []*apiserver.CacheAPIClient
	for _, endpoint := range endpoints {
		client, err := createGseApiServerClient(endpoint, conf)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return &GseApiServerClient{
		clients: clients,
	}, nil
}

// GetAgentStatus get agent status
func (g *GseApiServerClient) GetAgentStatus(ctx context.Context,
	requestInfo *apiserver.AgentStatusRequest) (*apiserver.AgentStatusResponse, error) {
	g.Lock()
	defer g.Unlock()
	return g.getClient().GetAgentStatus(ctx, requestInfo)
}

func (g *GseApiServerClient) getClient() *apiserver.CacheAPIClient {
	g.index++
	if g.index >= len(g.clients) {
		g.index = 0
	}
	return g.clients[g.index]
}

type GseTaskServerClient struct {
	clients []*taskserver.DoSomeCmdClient
	index   int
	sync.Mutex
}

// NewGseTaskServerClient new gse task server client
func NewGseTaskServerClient(endpoints []string, conf *util.TLSClientConfig) (*GseTaskServerClient, error) {
	var clients []*taskserver.DoSomeCmdClient
	for _, endpoint := range endpoints {
		client, err := createGseTaskServerClient(endpoint, conf)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return &GseTaskServerClient{
		clients: clients,
	}, nil
}

// PushFileV2 push host identifier to gse agent
func (g *GseTaskServerClient) PushFileV2(ctx context.Context,
	fileList []*taskserver.API_FileInfoV2) (*taskserver.API_CommRsp, error) {
	g.Lock()
	defer g.Unlock()
	return g.getClient().PushFileV2(ctx, fileList)
}

// GetPushFileRst get push file task result
func (g *GseTaskServerClient) GetPushFileRst(ctx context.Context, seqno string) (*taskserver.API_MapRsp, error) {
	g.Lock()
	defer g.Unlock()
	return g.getClient().GetPushFileRst(ctx, seqno)
}

func (g *GseTaskServerClient) getClient() *taskserver.DoSomeCmdClient {
	g.index++
	if g.index >= len(g.clients) {
		g.index = 0
	}
	return g.clients[g.index]
}

// createGseApiServerClient create thrift client for gse apiServer
func createGseApiServerClient(endpoint string, conf *util.TLSClientConfig) (*apiserver.CacheAPIClient, error) {
	var trans thrift.TTransport
	cfg, err := ssl.ClientTLSConfVerity(conf.CAFile, conf.CertFile, conf.KeyFile, conf.Password)
	if err != nil {
		return nil, err
	}
	cfg.InsecureSkipVerify = conf.InsecureSkipVerify

	trans, err = thrift.NewTSSLSocket(endpoint, cfg)
	if err != nil {
		return nil, err
	}

	trans = thrift.NewTFramedTransport(trans)
	clientProtocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	iprot := clientProtocolFactory.GetProtocol(trans)
	oprot := clientProtocolFactory.GetProtocol(trans)
	clientInner := thrift.NewTStandardClient(iprot, oprot)
	client := apiserver.NewCacheAPIClient(clientInner)

	if err = trans.Open(); err != nil {
		return nil, err
	}
	return client, nil
}

// createGseTaskServerClient create thrift client for gse taskServer
func createGseTaskServerClient(endpoint string, conf *util.TLSClientConfig) (*taskserver.DoSomeCmdClient, error) {
	cfg, err := ssl.ClientTLSConfVerity(conf.CAFile, conf.CertFile, conf.KeyFile, conf.Password)
	if err != nil {
		return nil, err
	}
	cfg.InsecureSkipVerify = conf.InsecureSkipVerify

	trans, err := thrift.NewTSSLSocket(endpoint, cfg)
	if err != nil {
		return nil, err
	}

	clientPRotocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	iprot := clientPRotocolFactory.GetProtocol(trans)
	oprot := clientPRotocolFactory.GetProtocol(trans)
	clientInner := thrift.NewTStandardClient(iprot, oprot)
	client := taskserver.NewDoSomeCmdClient(clientInner)

	if err = trans.Open(); err != nil {
		return nil, err
	}
	return client, nil
}
