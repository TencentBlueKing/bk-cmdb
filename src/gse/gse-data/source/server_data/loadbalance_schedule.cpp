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

#include "loadbalance_schedule.h"

#include <rapidjson/document.h>     // rapidjson's DOM-style API
#include <rapidjson/prettywriter.h> // for stringify JSON
#include <rapidjson/stringbuffer.h>

#include "net/message.h"
#include "net/protocol/data_protocol.h"

#include "net/tcp_connection.h"
#include "tools/json_helper.hpp"

#include "bbx/gse_errno.h"
#include "conftor/zkconftor.h"
#include "log/log.h"
#include "tools/error.h"

namespace gse {
namespace data {

static const int defaultMaxServerNode = 5;

LoadBalanceSchedule::LoadBalanceSchedule()
    : m_timeoutSecond(60),
      m_maxServerNode(defaultMaxServerNode),
      m_zkClient(NULL)
{
    m_accessRequestMessageQueuePtr = new moodycamel::ConcurrentQueue<Task>();
}

LoadBalanceSchedule::LoadBalanceSchedule(const std::string& zkHost, const std::string& zkauth)
    : m_timeoutSecond(60),
      m_maxServerNode(defaultMaxServerNode),
      m_zkClient(NULL)
{
    m_accessRequestMessageQueuePtr = new moodycamel::ConcurrentQueue<Task>();
    m_zkHost = zkHost;
    m_zkAuth = zkauth;
}

LoadBalanceSchedule::~LoadBalanceSchedule()
{
}

//调度仅仅是城市下的Nice值优选策略

void LoadBalanceSchedule::UpdateLoadBalance()
{
    // Get All child list
    std::vector<std::string> loadBalaceNodes;
    int ret = m_zkClient->GetChildConfItem(std::string(kBlanceConfigNodePath), loadBalaceNodes, NULL, NULL, -1);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to get node balance, %d", ret);
        return;
    }

    ClusterPreferServerList clusterPreferList;

    for (auto it : loadBalaceNodes)
    {
        std::string nodePath = std::string(kBlanceConfigNodePath) + "/" + it;
        std::string nodeValue;
        int ret = m_zkClient->GetConfItem(nodePath, nodeValue, NULL, NULL, -1);
        if (ret != GSE_SUCCESS)
        {
            LOG_ERROR("failed to get zk node(%s), ret:%d", nodePath.c_str(), ret);
            continue;
        }

        ServerBalanceInfo* balanceInfo = new ServerBalanceInfo();
        if (!balanceInfo->UnmarshalLoadBalance(nodeValue))
        {
            LOG_ERROR("failed to parse(%s)'s value(%s)", nodePath.c_str(), nodeValue.c_str());
            continue;
        }

        ClusterEndPoint* clusterEndPoint = nullptr;
        bool find = clusterPreferList.Find(balanceInfo->m_clusterName, clusterEndPoint);
        if (!find)
        {
            clusterEndPoint = new ClusterEndPoint(balanceInfo->m_clusterName);
            clusterPreferList.Push(balanceInfo->m_clusterName, clusterEndPoint);
        }

        ZoneEndPoint* zoneEndPoint = nullptr;
        find = clusterEndPoint->FindZone(balanceInfo->m_zoneId, zoneEndPoint);
        if (!find)
        {
            zoneEndPoint = new ZoneEndPoint(balanceInfo->m_zoneId);
            clusterEndPoint->PushZone(balanceInfo->m_zoneId, zoneEndPoint);
        }

        CityEndPoint* cityEndPoint = nullptr;
        // refresh , only add to list
        find = zoneEndPoint->FindCity(balanceInfo->m_cityId, cityEndPoint);
        if (!find)
        {
            cityEndPoint = new CityEndPoint(balanceInfo->m_cityId);
            LOG_DEBUG("find city(%s)", balanceInfo->m_cityId.c_str());
            zoneEndPoint->PushCity(balanceInfo->m_cityId, cityEndPoint);
        }

        //添加server 负载信息;
        LoadWeight loadWeightCfg;
        balanceInfo->CalcLoadWeight(loadWeightCfg);
        cityEndPoint->Push(balanceInfo);
    }

    m_preferClusterList.Swap(clusterPreferList);
}

bool LoadBalanceSchedule::EnqueuAccessRequest(const Task& task)
{
    if (m_accessRequestMessageQueuePtr->size_approx() >= kMaxAccessReqQueueSize)
    {
        LOG_ERROR("enqueue access request failed, the queue exceeds the max limit(%d), queue msg count:%d",
                  kMaxAccessReqQueueSize,
                  m_accessRequestMessageQueuePtr->size_approx());
        return false;
    }

    int max_retry_time = 0;
    while (!m_accessRequestMessageQueuePtr->enqueue(task))
    {
        max_retry_time++;
        if (max_retry_time > 5)
        {
            LOG_ERROR("enqueued failure, exceeding maximum retry times(%d)", max_retry_time);
            return false;
        }
    }

    return true;
}

void LoadBalanceSchedule::ProcessQueueMsg()
{
    Task tasks[kDefaultDequeueBulkSize];
    int times = 0;
    size_t num = 0;

    const int maxDequeueTimes = 10;
    do
    {
        num = m_accessRequestMessageQueuePtr->try_dequeue_bulk(tasks, kDefaultDequeueBulkSize);
        for (size_t i = 0; i < num; ++i)
        {
            ProcessAccessRequest(tasks[i]);
        }
        times++;
    } while ((num != 0) && times < maxDequeueTimes);

    return;
}

void LoadBalanceSchedule::ProcessAccessRequest(const Task& task)
{
    auto message = std::move(std::string(task.m_message->GetBody(), task.m_message->GetBodyLen()));
    AccessReq req;
    AccessResp resp;
    std::vector<ServerBalanceInfo> nodes;

    bool ret = req.Unmarshal(message.c_str());
    if (ret)
    {
        ScheduleAccessRequest(req, nodes);
        if (nodes.size() > 0)
        {
            resp.Pack(nodes, GSE_SUCCESS, "success");
            std::string respJson = std::move(resp.ToJson());
            task.m_callback(respJson.c_str(), respJson.length());
            LOG_DEBUG("process access request, send response(%s)", respJson.c_str());
        }
        else
        {
            resp.Pack(nodes, GSE_ERROR, "no access node available");
            std::string respJson = std::move(resp.ToJson());
            task.m_callback(respJson.c_str(), respJson.length());
            LOG_ERROR("failed to schedule access request, no available access node, cluster_name:%s, cluster_id:%s ,zone_id:%s, city_id:%s",
                      req.m_clusterName.c_str(), req.m_clusterId.c_str(), req.m_zoneId.c_str(), req.m_cityId.c_str());
        }
    }
    else
    {
        LOG_ERROR("failed to parse access request, request msg(%s)", message.c_str());
        std::string respJson = std::move(resp.ToJson());
        resp.Pack(nodes, GSE_ERROR, "request param invalid");
        task.m_callback(respJson.c_str(), respJson.length());
    }
}

bool LoadBalanceSchedule::EnqueuAccessRequest(CallBackFunc callback, std::shared_ptr<gse::net::Message> msg)
{
    Task task;
    task.m_callback = callback;
    task.m_message = msg;
    return EnqueuAccessRequest(task);
}

void LoadBalanceSchedule::ScheduleAccessRequest(const AccessReq& accessReq, std::vector<ServerBalanceInfo>& nodes)
{
    std::vector<std::string> preferClusterList;
    preferClusterList.push_back(accessReq.m_clusterName);
    if (accessReq.m_clusterName != kDefaultClusterName)
    {
        preferClusterList.push_back(kDefaultClusterName);
    }

    std::vector<std::string> preferZoneList;
    preferZoneList.push_back(accessReq.m_zoneId);
    if (accessReq.m_zoneId != kDefaultZoneId)
    {
        preferZoneList.push_back(kDefaultZoneId);
    }

    std::vector<std::string> preferCityList;
    preferCityList.push_back(accessReq.m_cityId);

    ClusterEndPoint* clusterEndPoint = nullptr;

    for (auto& iterClusterName : preferClusterList)
    {
        LOG_DEBUG("find prefer cluster:%s", iterClusterName.c_str());
        bool find = m_preferClusterList.FindPreferClusterEndPoint(iterClusterName, clusterEndPoint);

        if (!find)
        {
            LOG_DEBUG("not find prefer cluster:%s, continue find next", iterClusterName.c_str());
            continue;
        }

        ZoneEndPoint* zoneEndPoint = nullptr;
        for (auto& iterZoneId : preferZoneList)
        {
            LOG_DEBUG("find prefer zone_id:%s", iterZoneId.c_str());
            find = clusterEndPoint->FindPreferZoneEndPoint(iterZoneId, zoneEndPoint);
            if (!find)
            {
                LOG_DEBUG("not find prefer zone_id:%s, continue find next", iterZoneId.c_str());
                continue;
            }

            CityEndPoint* cityEndPoint = nullptr;
            std::vector<std::string> cityIds;
            zoneEndPoint->GetAllCity(cityIds);
            std::vector<std::string> allCityids;
            allCityids.clear();
            allCityids.insert(allCityids.end(), preferCityList.begin(), preferCityList.end());
            allCityids.insert(allCityids.end(), cityIds.begin(), cityIds.end());

            for (auto& iterCityId : allCityids)
            {
                LOG_DEBUG("find city:%s", iterCityId.c_str());
                find = zoneEndPoint->FindPreferCityEndPoint(iterCityId, cityEndPoint);

                if (!find)
                {
                    continue;
                }

                // std::vector<ServerBalanceInfo> serverBalanceInfo;

                find = cityEndPoint->FindPreferServerEndPoint(nodes);
                if (find)
                {
                    LOG_DEBUG("find prefer server list's count:%d", nodes.size());
                    return;
                }
            }
        }
    }

    return;
}

bool LoadBalanceSchedule::Start()
{
    ZkConftorParam conftorParam;
    conftorParam.m_ZkHost = m_zkHost;
    conftorParam.m_ZkAuth = m_zkAuth;
    conftorParam.m_password = m_zkAuth;

    m_zkClient = new ZkConftor(conftorParam);
    m_zkClient->Start();
    m_evLoopThread = std::make_shared<gse::eventloop::EventLoopThread>();
    bool ret = m_evLoopThread->Start(true, nullptr, nullptr);
    if (!ret)
    {
        LOG_ERROR("failed to start schedule event loop");
        return false;
    }

    //定时获取负载信息
    auto timerEventHandler = std::bind(&LoadBalanceSchedule::UpdateLoadBalance, this);
    m_timerUpdateLoadBalance = std::make_shared<gse::eventloop::TimerEventWatcher>(m_evLoopThread->GetLoop(), timerEventHandler, m_timeoutSecond);
    m_timerUpdateLoadBalance->Init();
    m_timerUpdateLoadBalance->StartTimer();

    UpdateLoadBalance();

    auto processQueueMsghandler = std::bind(&LoadBalanceSchedule::ProcessQueueMsg, this);

    m_timerProcessAccess = std::make_shared<gse::eventloop::TimerEventWatcher>(m_evLoopThread->GetLoop(), processQueueMsghandler, 0, 200, true);
    m_timerProcessAccess->Init();
    m_timerProcessAccess->StartTimer();

    LOG_DEBUG("successfully start load balance scheduler");

    return true;
}

void LoadBalanceSchedule::Stop()
{
    m_timerProcessAccess->StopTimer();
    m_timerUpdateLoadBalance->StopTimer();
    m_evLoopThread->Stop(true);
}

ClusterPreferServerList::ClusterPreferServerList()
{
}
ClusterPreferServerList::~ClusterPreferServerList()
{
    for (auto& it : m_clusterPreferServerList)
    {
        ClusterEndPoint* clusterEndPoint = it.second;
        if (clusterEndPoint != NULL)
        {
            LOG_DEBUG("when update clean old cluster(%s) config, zone size:%d", it.first.c_str(), clusterEndPoint->ZoneSize());
            delete clusterEndPoint;
        }
        else
        {
            LOG_WARN("when update clean old cluster(%s) config, cluster invalid", it.first.c_str());
        }
    }
    m_clusterPreferServerList.clear();
}

bool ClusterPreferServerList::Find(const std::string& clusterName, ClusterEndPoint*& clusterEndPoint)
{
    auto iter = m_clusterPreferServerList.find(clusterName);
    if (iter != m_clusterPreferServerList.end())
    {
        clusterEndPoint = iter->second;
        LOG_DEBUG("find cluter(%s) zone info, zone size:%d", clusterName.c_str(), clusterEndPoint->m_preferZoneList.size());
        return true;
    }

    LOG_DEBUG("can't find cluter(%s) zone info", clusterName.c_str());
    return false;
}

bool ClusterPreferServerList::FindPreferClusterEndPoint(const std::string& clusterName, ClusterEndPoint*& clusterEndPoint)
{
    bool find = false;
    find = Find(clusterName, clusterEndPoint);
    if (find)
    {
        return true;
    }

    return find;
}

void ClusterPreferServerList::Swap(ClusterPreferServerList& clusterPreferServerList)
{
    m_clusterPreferServerList.swap(clusterPreferServerList.m_clusterPreferServerList);
}

bool ClusterPreferServerList::Push(const std::string& clusterName, ClusterEndPoint* clusterEndPoint)
{
    auto ret = m_clusterPreferServerList.insert(std::pair<ClusterName, ClusterEndPoint*>(clusterName, clusterEndPoint));
    if (ret.second == false)
    {
        LOG_ERROR("cluster:%s already exist ,insert failed", clusterName.c_str());
        return false;
    }

    LOG_DEBUG("cluster:%s insert success, cluster count:%d", clusterName.c_str(), m_clusterPreferServerList.size());
    return true;
}

std::string ZonePreferServerList::ZoneListName()
{
    std::string zoneList;
    zoneList.append("[");

    for (auto iter = m_zonePreferServerList.begin(); iter != m_zonePreferServerList.end(); ++iter)
    {
        zoneList.append(iter->first);
        if (iter != m_zonePreferServerList.end())
        {
            zoneList.append(", ");
        }
    }
    zoneList.append("]");

    return zoneList;
}

ClusterEndPoint::~ClusterEndPoint()
{
    for (auto& it : m_preferZoneList)
    {
        ZoneEndPoint* zoneEndPoint = it.second;
        if (zoneEndPoint != NULL)
        {
            LOG_DEBUG("when update clean old zone(%s) config, zone size:%d", it.first.c_str(), zoneEndPoint->CitySize());
            delete zoneEndPoint;
        }
        else
        {
            LOG_DEBUG("when update clean old zone(%s) config, zone config invalid", it.first.c_str());
        }
    }
    m_preferZoneList.clear();
}

ClusterEndPoint::ClusterEndPoint(const std::string& clusterName)
    : m_clusterName(clusterName)
{
}

std::string ClusterEndPoint::ZoneListName()
{
    //    std::string zoneList;
    //    zoneList.append("[");

    //    for (auto iter = m_zonePreferServerList.begin(); iter != m_zonePreferServerList.end(); ++iter)
    //    {
    //        zoneList.append(iter->first);
    //        if (iter != m_zonePreferServerList.end())
    //        {
    //            zoneList.append(", ");
    //        }
    //    }
    //    zoneList.append("]");

    //    return zoneList;
}

size_t ClusterEndPoint::ZoneSize()
{
    return m_preferZoneList.size();
}

bool ClusterEndPoint::PushZone(const std::string& zoneId, ZoneEndPoint* zoneEndPoint)
{
    auto ret = m_preferZoneList.insert(std::pair<ZoneId, ZoneEndPoint*>(zoneId, zoneEndPoint));
    if (ret.second == false)
    {
        LOG_ERROR("zone:%s already exist ,insert failed", zoneId.c_str());
        return false;
    }

    LOG_DEBUG("zone:%s insert success, zone count:%d", zoneId.c_str(), m_preferZoneList.size());
    return true;
}

bool ClusterEndPoint::FindZone(const std::string& zoneId, ZoneEndPoint*& zoneEndPoint)
{
    auto iter = m_preferZoneList.find(zoneId);
    if (iter != m_preferZoneList.end())
    {
        LOG_DEBUG("find zone:%s in cluster:%s", zoneId.c_str(), m_clusterName.c_str());
        zoneEndPoint = iter->second;
        return true;
    }

    LOG_DEBUG("not find zone:%s in cluster:%s", zoneId.c_str(), m_clusterName.c_str());
    return false;
}

bool ClusterEndPoint::RemoveZone(const std::string& zoneId)
{
    return true;
}

bool ClusterEndPoint::FindPreferZoneEndPoint(const std::string& zoneId, ZoneEndPoint*& zoneEndPoint)
{
    auto iter = m_preferZoneList.find(zoneId);
    if (iter != m_preferZoneList.end())
    {
        zoneEndPoint = iter->second;
        LOG_DEBUG("find perfer zone(%s)", zoneId.c_str());
        return true;
    }

    iter = m_preferZoneList.find(kDefaultZoneId);
    if (iter != m_preferZoneList.end())
    {
        zoneEndPoint = iter->second;
        LOG_DEBUG("find default zone(%s)", kDefaultZoneId.c_str());
        return true;
    }

    LOG_DEBUG("can't find zone(%s)", zoneId.c_str());

    return false;
}

ZoneEndPoint::ZoneEndPoint(const std::string& zoneId)
    : m_zoneId(zoneId)
{
}

ZoneEndPoint::~ZoneEndPoint()
{
    for (auto& it : m_zonePreferCityList)
    {
        CityEndPoint* cityEndPoint = it.second;
        if (cityEndPoint != NULL)
        {
            LOG_DEBUG("when update clean old city(%s) config, city size:%d", it.first.c_str(), cityEndPoint->m_serverList.Size());
            delete cityEndPoint;
        }
        else
        {
            LOG_WARN("when update clean old city(%s) config, city invalid", it.first.c_str());
        }
    }
    m_zonePreferCityList.clear();
}

bool ZoneEndPoint::FindCity(const std::string& cityId, CityEndPoint*& cityEndPoint)
{
    auto iter = m_zonePreferCityList.find(cityId);
    if (iter != m_zonePreferCityList.end())
    {
        cityEndPoint = iter->second;
        return true;
    }

    return false;
}

bool ZoneEndPoint::PushCity(const std::string& cityId, CityEndPoint* cityEndPoint)
{
    auto ret = m_zonePreferCityList.insert(std::pair<CityId, CityEndPoint*>(cityId, cityEndPoint));
    if (ret.second == false)
    {
        LOG_WARN("city:%s already exist ,insert failed", cityId.c_str());
        return false;
    }

    LOG_DEBUG("city:%s insert success, city count:%d", cityId.c_str(), m_zonePreferCityList.size());
    return true;
}

bool ZoneEndPoint::FindPreferCityEndPoint(const std::string& cityId, CityEndPoint*& cityEndPoint)
{
    auto iter = m_zonePreferCityList.find(cityId);
    if (iter != m_zonePreferCityList.end())
    {
        cityEndPoint = iter->second;
        LOG_DEBUG("find perfer city(%s) ", cityId.c_str());
        return true;
    }

    LOG_DEBUG("can't find perfer city(%s)", cityId.c_str());

    return false;
}

bool ZoneEndPoint::GetAllCity(std::vector<std::string>& cityIds)
{
    for (auto iter : m_zonePreferCityList)
    {
        cityIds.push_back(iter.first);
        LOG_DEBUG("zone id:%s has city's id:%s", m_zoneId.c_str(), iter.first.c_str());
    }

    return false;
}

size_t ZoneEndPoint::CitySize()
{
    return m_zonePreferCityList.size();
}

CityEndPoint::CityEndPoint(const std::string& cityId)
    : m_cityId(cityId)
{
}

CityEndPoint::~CityEndPoint()
{
}

bool CityEndPoint::FindPreferServerEndPoint(std::vector<ServerBalanceInfo>& serverBanlanceInfo)
{
    auto iter = m_serverList.m_cityAllNodes.begin();
    for (int i = 0; i < m_serverList.Size(); i++)
    {
        serverBanlanceInfo.push_back(*(*iter));
        LOG_DEBUG("city(%s)'s server load info(node_ip:%s, load:%f)", (*iter)->m_cityId.c_str(), (*iter)->m_nodeIp.c_str(), (*iter)->m_loadWeight);
        if (iter != m_serverList.m_cityAllNodes.end())
        {
            iter++;
        }
    }
    return true;
}

bool CityEndPoint::Push(ServerBalanceInfo* serverBalanceInfo)
{
    m_serverList.Push(serverBalanceInfo);
    return true;
}

bool ServerBalanceInfo::Equal(ServerBalanceInfo* balanceInfo)
{
    if (m_clusterName == balanceInfo->m_clusterName &&
        m_zoneId == balanceInfo->m_zoneId &&
        m_cityId == balanceInfo->m_cityId &&
        m_nodeIp == balanceInfo->m_nodeIp)
    {
        return true;
    }

    return false;
}

void ServerBalanceInfo::Update(const ServerBalanceInfo* balanceInfo)
{
    *this = *balanceInfo;
}

bool ServerBalanceInfo::UnmarshalLoadBalance(const std::string& data)
{
    std::string errMsg;
    rj::Document doc;
    if (!gse::tools::json::LoadDocument(data, doc, errMsg))
    {
        LOG_ERROR("the json is invalid, json(%s), errmsg:%s", data.c_str(), errMsg.c_str());
        return false;
    }

    //这里需要整体补充errormsg，当前仅实现逻辑

    rj::Value::ConstMemberIterator iter = doc.FindMember(kMetaData);
    if (iter != doc.MemberEnd() && iter->value.IsObject())
    {
        auto zoneId = gse::tools::json::JsonHelper<std::string>::GetValue(iter->value, kZoneId, "");

        if (zoneId.m_isOK)
        {
            m_zoneId = zoneId.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kZoneId);
            return false;
        }

        auto cityId = gse::tools::json::JsonHelper<std::string>::GetValue(iter->value, kCityId, "");

        if (cityId.m_isOK)
        {
            m_cityId = cityId.m_value;
        }
        else
        {
            LOG_ERROR("city_id invalid");
            return false;
        }

        auto clusterName = gse::tools::json::JsonHelper<std::string>::GetValue(iter->value, kClusterName, "");

        if (clusterName.m_isOK)
        {
            m_clusterName = clusterName.m_value;
        }
        else
        {
            LOG_ERROR("cluster_name invalid");
            return false;
        }

        auto clusterId = gse::tools::json::JsonHelper<std::string>::GetValue(iter->value, kClusterId, "");

        if (clusterId.m_isOK)
        {
            m_clusterId = clusterId.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kClusterId);
            return false;
        }

        auto serviceId = gse::tools::json::JsonHelper<std::string>::GetValue(iter->value, kServiceId, "");

        if (serviceId.m_isOK)
        {
            m_serviceId = serviceId.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kServiceId);
            return false;
        }

        auto nodeIp = gse::tools::json::JsonHelper<std::string>::GetValue(iter->value, kNodeIp, "");
        if (nodeIp.m_isOK)
        {
            m_nodeIp = nodeIp.m_value;
        }
        else
        {
            LOG_ERROR("%s json invalid", kNodeIp);
            return false;
        }
    }

    rj::Value::ConstMemberIterator agentServiceIter = doc.FindMember(kAgentServiceName);
    if (agentServiceIter != doc.MemberEnd() && agentServiceIter->value.IsObject())
    {

        auto serviceIp = gse::tools::json::JsonHelper<std::string>::GetValue(agentServiceIter->value, kServiceIp, "");

        if (serviceIp.m_isOK)
        {
            m_hostAddr = serviceIp.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kServiceIp);
            return false;
        }

        auto servicePort = gse::tools::json::JsonHelper<int>::GetValue(agentServiceIter->value, kServicePort, 0);

        if (servicePort.m_isOK)
        {
            m_hostPort = servicePort.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kServicePort);
            return false;
        }

        auto protocol = gse::tools::json::JsonHelper<std::string>::GetValue(agentServiceIter->value, kProtocol, "TCP");

        if (protocol.m_isOK)
        {
            m_protocol = protocol.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kProtocol);
        }

        auto ssl = gse::tools::json::JsonHelper<bool>::GetValue(agentServiceIter->value, kSSL, 0);

        if (ssl.m_isOK)
        {
            m_ssl = ssl.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kProtocol);
        }
    }

    rj::Value::ConstMemberIterator loadIter = doc.FindMember(kLoad);
    if (loadIter != doc.MemberEnd() && loadIter->value.IsObject())
    {

        auto agentCountMax = gse::tools::json::JsonHelper<int>::GetValue(loadIter->value, kAgentCountMax, 0);

        if (agentCountMax.m_isOK)
        {
            m_maxAgentCount = agentCountMax.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kAgentCountMax);
            return false;
        }

        auto agentCount = gse::tools::json::JsonHelper<int>::GetValue(loadIter->value, kAgentCount, 0);

        if (agentCount.m_isOK)
        {
            m_agentCount = agentCount.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kAgentCount);
            return false;
        }

        rj::Value::ConstMemberIterator extendIter = loadIter->value.FindMember(kExtend);

        if (extendIter != doc.MemberEnd() && extendIter->value.IsObject())
        {
            auto cpuUsage = gse::tools::json::JsonHelper<double>::GetValue(extendIter->value, kCpuRate, 0.0);

            if (cpuUsage.m_isOK)
            {
                m_cpuUsage = cpuUsage.m_value;
            }
            else
            {
                LOG_ERROR("%s invalid", kCpuRate);
                return false;
            }

            auto memUsage = gse::tools::json::JsonHelper<double>::GetValue(extendIter->value, kMemoryRate, 0.0);

            if (memUsage.m_isOK)
            {
                m_memUsage = memUsage.m_value;
            }
            else
            {
                LOG_ERROR("%s invalid", kMemoryRate);
                return false;
            }

            auto netUsage = gse::tools::json::JsonHelper<double>::GetValue(extendIter->value, kNetworkRate, 0.0);

            if (netUsage.m_isOK)
            {
                m_netUsage = netUsage.m_value;
            }
            else
            {
                LOG_ERROR("%s invalid", kNetworkRate);
                return false;
            }
        }
    }
    return true;
}

float ServerBalanceInfo::CalcLoadWeight(const LoadWeight& loadWeightCfg)
{
    float totalLoadWeight = (loadWeightCfg.m_connectionLoadWeight +
                             loadWeightCfg.m_memUsageLoadWeight +
                             loadWeightCfg.m_netUsageLoadWeight +
                             loadWeightCfg.m_cpuUsageLoadWeight);
    if (totalLoadWeight <= 0)
    {
        LOG_ERROR("loadweight config invalid, total value must 100");
        return 0;
    }

    float agentcountLoad = 0;
    float cpuusageLoad = 0;
    float memusageLoad = 0;
    float netusageLoad = 0;

    int agentCount = 0;
    if (m_agentCount >= m_maxAgentCount)
    {
        agentCount = m_maxAgentCount;
    }
    else
    {
        agentCount = m_agentCount;
    }

    agentcountLoad = (100 * (agentCount) / m_maxAgentCount) * ((loadWeightCfg.m_connectionLoadWeight * 1.0) / totalLoadWeight);
    LOG_DEBUG("agentcount loadweight:%f", agentcountLoad);

    cpuusageLoad = m_cpuUsage * (loadWeightCfg.m_cpuUsageLoadWeight * 1.0 / totalLoadWeight);
    LOG_DEBUG("cpuusage loadweight:%f", cpuusageLoad);

    netusageLoad = m_netUsage * (loadWeightCfg.m_netUsageLoadWeight * 1.0 / totalLoadWeight);
    LOG_DEBUG("netusage loadweight:%f", netusageLoad);

    memusageLoad = m_memUsage * (loadWeightCfg.m_memUsageLoadWeight * 1.0 / totalLoadWeight);
    LOG_DEBUG("memusage loadweight:%f", memusageLoad);

    float loadWeight = agentcountLoad + cpuusageLoad + memusageLoad + netusageLoad;
    float nice = 100 - (100.0 * loadWeight) / totalLoadWeight;

    m_loadWeight = nice;

    LOG_DEBUG("cluster(%s) server(%s), agent_count:%f, agent_count_load:%f , cpuusage:%f, cpuusage_load:%f, netusage:%f, netusage_load:%f, memusage:%f,, memusage_load:%f, total use load:%f, nice:%f",
              m_clusterName.c_str(), m_nodeIp.c_str(),
              agentCount, agentcountLoad,
              m_cpuUsage, cpuusageLoad,
              m_netUsage, netusageLoad,
              m_memUsage, memusageLoad,
              loadWeight, nice);

    LOG_DEBUG("nodeip:%s, agent count:%d, cpu usage:%f, net usage:%f, mem usage:%f, connection loadWeight:%d, net usage loadWeight:%d, cpu usage loadWeight:%d, mem usage:%d, totalLoadWeight:%f",
              m_nodeIp.c_str(),
              m_agentCount,
              m_cpuUsage,
              m_netUsage,
              m_memUsage,
              loadWeightCfg.m_connectionLoadWeight,
              loadWeightCfg.m_memUsageLoadWeight,
              loadWeightCfg.m_netUsageLoadWeight,
              loadWeightCfg.m_cpuUsageLoadWeight,
              totalLoadWeight);

    return nice;
}

CityPreferServerList::~CityPreferServerList()
{
    for (auto& it : m_cityAllNodes)
    {
        ServerBalanceInfo* balanceInfo = it;
        LOG_DEBUG("when update clean old balance info(%s) config(cpu:%f, net:%f, mem:%f, connection:%d)",
                  balanceInfo->m_nodeIp.c_str(),
                  balanceInfo->m_cpuUsage, balanceInfo->m_netUsage, balanceInfo->m_memUsage,
                  balanceInfo->m_agentCount);
        delete balanceInfo;
    }

    m_cityAllNodes.clear();
}

void CityPreferServerList::Push(ServerBalanceInfo* serverBalanceInfo)
{
    ///正常情况下不会有重复的，但为了保证准确，所以会遍历查找,data server 集群内节点数量也是可控
    bool find = false;
    LOG_DEBUG("push balance info(cluster:%s, zone:%s, city:%s, node ip:%s)",
              serverBalanceInfo->m_clusterId.c_str(),
              serverBalanceInfo->m_zoneId.c_str(),
              serverBalanceInfo->m_cityId.c_str(), serverBalanceInfo->m_nodeIp.c_str());

    for (auto it = m_cityAllNodes.begin(); it != m_cityAllNodes.end(); it++)
    {
        if ((*it)->Equal(serverBalanceInfo))
        {
            (*it)->Update(serverBalanceInfo);
            find = true;
            delete serverBalanceInfo;
            break;
        }
    }

    if (!find)
    {
        m_cityAllNodes.push_back(serverBalanceInfo);
        LOG_DEBUG("can't find city node list, push balance info(cluster:%s, zone:%s, city:%s, node ip:%s)",
                  serverBalanceInfo->m_clusterId.c_str(),
                  serverBalanceInfo->m_zoneId.c_str(),
                  serverBalanceInfo->m_cityId.c_str(), serverBalanceInfo->m_nodeIp.c_str());
    }

    m_cityAllNodes.sort([=](ServerBalanceInfo* first, ServerBalanceInfo* second) -> bool {
        return first->m_loadWeight > second->m_loadWeight;
    });

    // TODO:print list
}

uint64_t CityPreferServerList::Size()
{
    return m_cityAllNodes.size();
}

ServerBalanceInfo* CityPreferServerList::operator[](int index)
{
    if (0 <= index && index < m_cityAllNodes.size())
    {
        // return m_cityAllNodes[index];
        return NULL;
    }
    else
    {
        return NULL;
    }
}

} // namespace data
} // namespace gse
