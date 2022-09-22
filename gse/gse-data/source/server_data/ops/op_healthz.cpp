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
#include "op_healthz.h"

#include <rapidjson/document.h>     // rapidjson's DOM-style API
#include <rapidjson/prettywriter.h> // for stringify JSON
#include <rapidjson/stringbuffer.h>

#include "bbx/bbx.h"
#include "dataserver.h"
#include "discover/zkapi/zk_api.h"

#define HEALTHZ_ZK_SERIVCE_NODE_PATH "/gse/v2/service/data"
#define HEALTHZ_ZK_CHNANELID_NODE_PATH "/gse/config/server/dataserver/channelid"

namespace gse {
namespace data {
static gse::safe::RWLock g_healthLock;
static std::map<std::string, HealthMesage> g_initHealthZInfo;

OpHealthZ::OpHealthZ()
{
}

OpHealthZ::~OpHealthZ()
{
}

void OpHealthZ::AddInitHealthInfo(const std::string &healthName, const std::string &healthMessae, int status)
{
    g_healthLock.WLock();
    auto it = g_initHealthZInfo.find(healthName);

    if (it != g_initHealthZInfo.end())
    {
        g_initHealthZInfo.erase(it);
    }

    g_initHealthZInfo.insert(std::pair<HealthzName, HealthMesage>(healthName, HealthMesage(healthMessae, status)));
    g_healthLock.UnLock();
}

void OpHealthZ::DiscoverZkCheck()
{
    auto zkClient = DataServer::GetDiscoverZkApiClientInst();
    if (zkClient == nullptr)
    {

        OpHealthZ::AddInitHealthInfo("zookeeper", "zookeeper not init", -1);
        return;
    }

    std::vector<std::string> childNode;
    int ret = zkClient->ZkGetChildren(HEALTHZ_ZK_SERIVCE_NODE_PATH, nullptr, nullptr, childNode, nullptr);
    if (ret != GSE_SUCCESS)
    {
        OpHealthZ::AddInitHealthInfo("zookeeper", gse::bbx::errnos::Err2Msg(ret).c_str(), -1);
        return;
    }

    if (childNode.size() == 0)
    {
        OpHealthZ::AddInitHealthInfo("zookeeper", gse::bbx::errnos::Err2Msg(ret).c_str(), -1);
    }

    OpHealthZ::AddInitHealthInfo("zookeeper", "ok", 0);

    return;
}

void OpHealthZ::ChannelIdZkCheck()
{
    auto channelZkApiInst = DataServer::GetChannelIdZkApiInst();
    if (channelZkApiInst == nullptr)
    {
        OpHealthZ::AddInitHealthInfo("channelid_zookeeper", "zookeeper not init", -1);
        return;
    }

    auto zkApiClient = channelZkApiInst->GetZkApiClient();
    if (zkApiClient == nullptr)
    {
        OpHealthZ::AddInitHealthInfo("channelid_zookeeper", "zookeeper not init", -1);
        return;
    }

    std::vector<std::string> childNode;
    int ret = zkApiClient->ZkGetChildren(HEALTHZ_ZK_CHNANELID_NODE_PATH, nullptr, nullptr, childNode, nullptr);
    if (ret != GSE_SUCCESS)
    {
        OpHealthZ::AddInitHealthInfo("channelid_zookeeper", gse::bbx::errnos::Err2Msg(ret).c_str(), -1);
        return;
    }

    OpHealthZ::AddInitHealthInfo("channelid_zookeeper", "ok", 0);

    return;
}

int OpHealthZ::OnGet(HTTPMessagePtr message, std::string &response) noexcept
{
    rapidjson::Document rootNode(rapidjson::kObjectType);
    rapidjson::Document::AllocatorType &allocator = rootNode.GetAllocator();
    rapidjson::Value allStatus(rapidjson::kObjectType);

    ChannelIdZkCheck();
    DiscoverZkCheck();

    bool isOK = true;
    g_healthLock.RLock();
    for (auto &it : g_initHealthZInfo)
    {
        int status;
        std::string &statusMessage = std::get<0>(it.second);
        status = std::get<1>(it.second);

        if (status < 0)
        {
            isOK = false;
        }

        rapidjson::Value strKey;
        rapidjson::Value strValue;
        strKey.SetString(it.first.c_str(), it.first.length(), allocator);
        strValue.SetString(statusMessage.c_str(), statusMessage.length(), allocator);
        allStatus.AddMember(strKey.Move(), strValue.Move(), allocator);
    }
    g_healthLock.UnLock();

    rootNode.AddMember("ok", rapidjson::Value().SetBool(isOK), allocator);
    std::string version;
#ifdef _VERSION_
    version.assign(std::string(_VERSION_));
#else
    version.assign("1.0");
#endif

#ifndef _GSE_SSL_
    version.append("i");
#endif

    rootNode.AddMember("version", version, allocator);
    rootNode.AddMember("data", allStatus, allocator);

    rapidjson::StringBuffer strBuffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(strBuffer);

    rootNode.Accept(writer);

    response.assign(strBuffer.GetString());
    if (!isOK)
    {
        return HTTP_STATUS_CODE_500;
    }
    return HTTP_STATUS_CODE_200;
}

int OpHealthZ::OnPost(HTTPMessagePtr message, std::string &response) noexcept
{
    return OnGet(message, response);
}

std::string OpHealthZ::GetURI()
{
    return std::string("/healthz");
}
} // namespace data
} // namespace gse
