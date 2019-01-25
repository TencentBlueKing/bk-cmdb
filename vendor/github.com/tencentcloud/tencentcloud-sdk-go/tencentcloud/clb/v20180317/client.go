// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20180317

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-03-17"

type Client struct {
    common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
    cpf := profile.NewClientProfile()
    client = &Client{}
    client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
    return
}

func NewClient(credential *common.Credential, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
    client = &Client{}
    client.Init(region).
        WithCredential(credential).
        WithProfile(clientProfile)
    return
}


func NewCreateListenerRequest() (request *CreateListenerRequest) {
    request = &CreateListenerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "CreateListener")
    return
}

func NewCreateListenerResponse() (response *CreateListenerResponse) {
    response = &CreateListenerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 在一个负载均衡实例下创建监听器。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) CreateListener(request *CreateListenerRequest) (response *CreateListenerResponse, err error) {
    if request == nil {
        request = NewCreateListenerRequest()
    }
    response = NewCreateListenerResponse()
    err = c.Send(request, response)
    return
}

func NewCreateLoadBalancerRequest() (request *CreateLoadBalancerRequest) {
    request = &CreateLoadBalancerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "CreateLoadBalancer")
    return
}

func NewCreateLoadBalancerResponse() (response *CreateLoadBalancerResponse) {
    response = &CreateLoadBalancerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateLoadBalancer 接口用来创建负载均衡实例。为了使用负载均衡服务，您必须要购买一个或者多个负载均衡实例。通过成功调用该接口，会返回负载均衡实例的唯一 ID。用户可以购买的负载均衡实例类型分为：公网（应用型）、内网（应用型）。可以参考产品说明的产品类型。
// 本接口成功返回后，可使用查询负载均衡实例列表接口DescribeLoadBalancers查询负载均衡实例的状态，以确定是否创建成功。
func (c *Client) CreateLoadBalancer(request *CreateLoadBalancerRequest) (response *CreateLoadBalancerResponse, err error) {
    if request == nil {
        request = NewCreateLoadBalancerRequest()
    }
    response = NewCreateLoadBalancerResponse()
    err = c.Send(request, response)
    return
}

func NewCreateRuleRequest() (request *CreateRuleRequest) {
    request = &CreateRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "CreateRule")
    return
}

func NewCreateRuleResponse() (response *CreateRuleResponse) {
    response = &CreateRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateRule 接口用于在一个已存在的应用型负载均衡七层监听器下创建转发规则，七层监听器中，后端机器必须绑定到规则上而非监听器上。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) CreateRule(request *CreateRuleRequest) (response *CreateRuleResponse, err error) {
    if request == nil {
        request = NewCreateRuleRequest()
    }
    response = NewCreateRuleResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteListenerRequest() (request *DeleteListenerRequest) {
    request = &DeleteListenerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DeleteListener")
    return
}

func NewDeleteListenerResponse() (response *DeleteListenerResponse) {
    response = &DeleteListenerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口用来删除应用型（四层和七层）负载均衡实例下的监听器。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) DeleteListener(request *DeleteListenerRequest) (response *DeleteListenerResponse, err error) {
    if request == nil {
        request = NewDeleteListenerRequest()
    }
    response = NewDeleteListenerResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteLoadBalancerRequest() (request *DeleteLoadBalancerRequest) {
    request = &DeleteLoadBalancerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DeleteLoadBalancer")
    return
}

func NewDeleteLoadBalancerResponse() (response *DeleteLoadBalancerResponse) {
    response = &DeleteLoadBalancerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteLoadBalancer 接口用来删除用户指定的一个负载均衡实例。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) DeleteLoadBalancer(request *DeleteLoadBalancerRequest) (response *DeleteLoadBalancerResponse, err error) {
    if request == nil {
        request = NewDeleteLoadBalancerRequest()
    }
    response = NewDeleteLoadBalancerResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteRuleRequest() (request *DeleteRuleRequest) {
    request = &DeleteRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DeleteRule")
    return
}

func NewDeleteRuleResponse() (response *DeleteRuleResponse) {
    response = &DeleteRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteRule 接口用来删除应用型负载均衡实例七层监听器下的转发规则。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) DeleteRule(request *DeleteRuleRequest) (response *DeleteRuleResponse, err error) {
    if request == nil {
        request = NewDeleteRuleRequest()
    }
    response = NewDeleteRuleResponse()
    err = c.Send(request, response)
    return
}

func NewDeregisterTargetsRequest() (request *DeregisterTargetsRequest) {
    request = &DeregisterTargetsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DeregisterTargets")
    return
}

func NewDeregisterTargetsResponse() (response *DeregisterTargetsResponse) {
    response = &DeregisterTargetsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeregisterTargets 接口用来将一台或多台后端机器从应用型负载均衡的监听器上解绑，对于四层监听器（TCP、UDP），只需指定监听器ID即可，对于七层监听器（HTTP、HTTPS），还需通过LocationId或者Domain+Url指定转发规则。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) DeregisterTargets(request *DeregisterTargetsRequest) (response *DeregisterTargetsResponse, err error) {
    if request == nil {
        request = NewDeregisterTargetsRequest()
    }
    response = NewDeregisterTargetsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeListenersRequest() (request *DescribeListenersRequest) {
    request = &DescribeListenersRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DescribeListeners")
    return
}

func NewDescribeListenersResponse() (response *DescribeListenersResponse) {
    response = &DescribeListenersResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeListeners 接口可根据负载均衡器 ID，监听器的协议或者端口作为过滤条件获取监听器列表。如果不指定任何过滤条件，默认返该负载均衡器下的默认数据长度（20 个）的监听器。
func (c *Client) DescribeListeners(request *DescribeListenersRequest) (response *DescribeListenersResponse, err error) {
    if request == nil {
        request = NewDescribeListenersRequest()
    }
    response = NewDescribeListenersResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeLoadBalancersRequest() (request *DescribeLoadBalancersRequest) {
    request = &DescribeLoadBalancersRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DescribeLoadBalancers")
    return
}

func NewDescribeLoadBalancersResponse() (response *DescribeLoadBalancersResponse) {
    response = &DescribeLoadBalancersResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询负载均衡实例列表
func (c *Client) DescribeLoadBalancers(request *DescribeLoadBalancersRequest) (response *DescribeLoadBalancersResponse, err error) {
    if request == nil {
        request = NewDescribeLoadBalancersRequest()
    }
    response = NewDescribeLoadBalancersResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTargetsRequest() (request *DescribeTargetsRequest) {
    request = &DescribeTargetsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DescribeTargets")
    return
}

func NewDescribeTargetsResponse() (response *DescribeTargetsResponse) {
    response = &DescribeTargetsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTargets 接口用来查询应用型负载均衡实例的某些监听器后端绑定的机器列表。
func (c *Client) DescribeTargets(request *DescribeTargetsRequest) (response *DescribeTargetsResponse, err error) {
    if request == nil {
        request = NewDescribeTargetsRequest()
    }
    response = NewDescribeTargetsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskStatusRequest() (request *DescribeTaskStatusRequest) {
    request = &DescribeTaskStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "DescribeTaskStatus")
    return
}

func NewDescribeTaskStatusResponse() (response *DescribeTaskStatusResponse) {
    response = &DescribeTaskStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口用于查询异步执行任务的状态，对于非查询类的接口（创建/删除负载均衡实例、监听器、规则以及绑定或解绑后端机器等），在调用成功后都需要使用本接口查询任务是否最终执行成功。
func (c *Client) DescribeTaskStatus(request *DescribeTaskStatusRequest) (response *DescribeTaskStatusResponse, err error) {
    if request == nil {
        request = NewDescribeTaskStatusRequest()
    }
    response = NewDescribeTaskStatusResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDomainRequest() (request *ModifyDomainRequest) {
    request = &ModifyDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "ModifyDomain")
    return
}

func NewModifyDomainResponse() (response *ModifyDomainResponse) {
    response = &ModifyDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyDomain接口用来修改应用型负载均衡七层监听器下的域名。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) ModifyDomain(request *ModifyDomainRequest) (response *ModifyDomainResponse, err error) {
    if request == nil {
        request = NewModifyDomainRequest()
    }
    response = NewModifyDomainResponse()
    err = c.Send(request, response)
    return
}

func NewModifyListenerRequest() (request *ModifyListenerRequest) {
    request = &ModifyListenerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "ModifyListener")
    return
}

func NewModifyListenerResponse() (response *ModifyListenerResponse) {
    response = &ModifyListenerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyListener接口用来修改应用型负载均衡监听器的属性，包括监听器名称、健康检查参数、证书信息、转发策略等。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) ModifyListener(request *ModifyListenerRequest) (response *ModifyListenerResponse, err error) {
    if request == nil {
        request = NewModifyListenerRequest()
    }
    response = NewModifyListenerResponse()
    err = c.Send(request, response)
    return
}

func NewModifyLoadBalancerAttributesRequest() (request *ModifyLoadBalancerAttributesRequest) {
    request = &ModifyLoadBalancerAttributesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "ModifyLoadBalancerAttributes")
    return
}

func NewModifyLoadBalancerAttributesResponse() (response *ModifyLoadBalancerAttributesResponse) {
    response = &ModifyLoadBalancerAttributesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 修改负载均衡实例的属性，目前仅用于修改负载均衡实例的名称。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) ModifyLoadBalancerAttributes(request *ModifyLoadBalancerAttributesRequest) (response *ModifyLoadBalancerAttributesResponse, err error) {
    if request == nil {
        request = NewModifyLoadBalancerAttributesRequest()
    }
    response = NewModifyLoadBalancerAttributesResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRuleRequest() (request *ModifyRuleRequest) {
    request = &ModifyRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "ModifyRule")
    return
}

func NewModifyRuleResponse() (response *ModifyRuleResponse) {
    response = &ModifyRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyRule 接口用来修改应用型负载均衡七层监听器下的转发规则的各项属性，包括转发路径、健康检查属性、转发策略等。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) ModifyRule(request *ModifyRuleRequest) (response *ModifyRuleResponse, err error) {
    if request == nil {
        request = NewModifyRuleRequest()
    }
    response = NewModifyRuleResponse()
    err = c.Send(request, response)
    return
}

func NewModifyTargetPortRequest() (request *ModifyTargetPortRequest) {
    request = &ModifyTargetPortRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "ModifyTargetPort")
    return
}

func NewModifyTargetPortResponse() (response *ModifyTargetPortResponse) {
    response = &ModifyTargetPortResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyTargetPort接口用于修改监听器绑定的后端云服务器的端口。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) ModifyTargetPort(request *ModifyTargetPortRequest) (response *ModifyTargetPortResponse, err error) {
    if request == nil {
        request = NewModifyTargetPortRequest()
    }
    response = NewModifyTargetPortResponse()
    err = c.Send(request, response)
    return
}

func NewModifyTargetWeightRequest() (request *ModifyTargetWeightRequest) {
    request = &ModifyTargetWeightRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "ModifyTargetWeight")
    return
}

func NewModifyTargetWeightResponse() (response *ModifyTargetWeightResponse) {
    response = &ModifyTargetWeightResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyTargetWeight 接口用于修改监听器绑定的后端机器的转发权重。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) ModifyTargetWeight(request *ModifyTargetWeightRequest) (response *ModifyTargetWeightResponse, err error) {
    if request == nil {
        request = NewModifyTargetWeightRequest()
    }
    response = NewModifyTargetWeightResponse()
    err = c.Send(request, response)
    return
}

func NewRegisterTargetsRequest() (request *RegisterTargetsRequest) {
    request = &RegisterTargetsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("clb", APIVersion, "RegisterTargets")
    return
}

func NewRegisterTargetsResponse() (response *RegisterTargetsResponse) {
    response = &RegisterTargetsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// RegisterTargets 接口用来将一台或多台后端机器注册到应用型负载均衡的监听器，对于四层监听器（TCP、UDP），只需指定监听器ID即可，对于七层监听器（HTTP、HTTPS），还需通过LocationId或者Domain+Url指定转发规则。
// 本接口为异步接口，本接口返回成功后需以返回的RequestID为入参，调用DescribeTaskStatus接口查询本次任务是否成功。
func (c *Client) RegisterTargets(request *RegisterTargetsRequest) (response *RegisterTargetsResponse, err error) {
    if request == nil {
        request = NewRegisterTargetsRequest()
    }
    response = NewRegisterTargetsResponse()
    err = c.Send(request, response)
    return
}
