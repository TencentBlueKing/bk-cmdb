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

// GseApiServerClient TODO
type GseApiServerClient struct {
	endpoints []string
	tlsConf   *util.TLSClientConfig
	index     int
	sync.Mutex
}

// NewGseApiServerClient new gse api server client
func NewGseApiServerClient(endpoints []string, conf *util.TLSClientConfig) (*GseApiServerClient, error) {
	return &GseApiServerClient{
		endpoints: endpoints,
		tlsConf:   conf,
	}, nil
}

// GetAgentStatus get agent status
func (g *GseApiServerClient) GetAgentStatus(ctx context.Context,
	requestInfo *apiserver.AgentStatusRequest) (*apiserver.AgentStatusResponse, error) {
	client, err := g.getClient()
	if err != nil {
		return nil, err
	}
	return client.getAgentStatus(ctx, requestInfo)
}

func (g *GseApiServerClient) getClient() (*apiClient, error) {
	g.Lock()
	defer g.Unlock()
	g.index++
	if g.index >= len(g.endpoints) {
		g.index = 0
	}

	client, err := createGseApiServerClient(g.endpoints[g.index], g.tlsConf)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GseTaskServerClient TODO
type GseTaskServerClient struct {
	endpoints []string
	tlsConf   *util.TLSClientConfig
	index     int
	sync.Mutex
}

// NewGseTaskServerClient new gse task server client
func NewGseTaskServerClient(endpoints []string, conf *util.TLSClientConfig) (*GseTaskServerClient, error) {
	return &GseTaskServerClient{
		endpoints: endpoints,
		tlsConf:   conf,
	}, nil
}

// PushFileV2 push host identifier to gse agent
func (g *GseTaskServerClient) PushFileV2(ctx context.Context,
	fileList []*taskserver.API_FileInfoV2) (*taskserver.API_CommRsp, error) {
	client, err := g.getClient()
	if err != nil {
		return nil, err
	}
	return client.pushFileV2(ctx, fileList)
}

// GetPushFileRst get push file task result
func (g *GseTaskServerClient) GetPushFileRst(ctx context.Context, seqno string) (*taskserver.API_MapRsp, error) {
	client, err := g.getClient()
	if err != nil {
		return nil, err
	}
	return client.getPushFileRst(ctx, seqno)
}

func (g *GseTaskServerClient) getClient() (*taskClient, error) {
	g.Lock()
	defer g.Unlock()
	g.index++
	if g.index >= len(g.endpoints) {
		g.index = 0
	}

	client, err := createGseTaskServerClient(g.endpoints[g.index], g.tlsConf)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// createGseApiServerClient create thrift client for gse apiServer
func createGseApiServerClient(endpoint string, conf *util.TLSClientConfig) (*apiClient, error) {
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
	return &apiClient{client: client, trans: trans}, nil
}

// createGseTaskServerClient create thrift client for gse taskServer
func createGseTaskServerClient(endpoint string, conf *util.TLSClientConfig) (*taskClient, error) {
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
	return &taskClient{client: client, trans: trans}, nil
}

type apiClient struct {
	client *apiserver.CacheAPIClient
	trans  thrift.TTransport
}

func (a *apiClient) getAgentStatus(ctx context.Context,
	requestInfo *apiserver.AgentStatusRequest) (*apiserver.AgentStatusResponse, error) {
	defer a.trans.Close()
	return a.client.GetAgentStatus(ctx, requestInfo)
}

type taskClient struct {
	client *taskserver.DoSomeCmdClient
	trans  thrift.TTransport
}

func (t *taskClient) pushFileV2(ctx context.Context, fileList []*taskserver.API_FileInfoV2) (*taskserver.API_CommRsp,
	error) {
	defer t.trans.Close()
	return t.client.PushFileV2(ctx, fileList)
}

func (t *taskClient) getPushFileRst(ctx context.Context, seqno string) (*taskserver.API_MapRsp, error) {
	defer t.trans.Close()
	return t.client.GetPushFileRst(ctx, seqno)
}
