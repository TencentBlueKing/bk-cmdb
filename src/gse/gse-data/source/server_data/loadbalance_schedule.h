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

#ifndef _GSE_DATA_LOADBALANCE_SCHEDULE_H_
#define _GSE_DATA_LOADBALANCE_SCHEDULE_H_

#include <list>
#include <map>
#include <memory>

#include "concurrentqueue/blockingconcurrentqueue.h"
#include "concurrentqueue/concurrentqueue.h"

#include "eventloop/event_loop.h"
#include "eventloop/loop_thread.h"
#include "log/log.h"
#include "net/tcp_connection.h"
#include "tools/json_helper.hpp"

#include <rapidjson/document.h>     // rapidjson's DOM-style API
#include <rapidjson/prettywriter.h> // for stringify JSON
#include <rapidjson/stringbuffer.h>

#include "bbx/gse_errno.h"
#include "tools/rapidjson_macro.h"

#include "conftor/conftor.h"
namespace gse {
namespace data {

static const std::string kDefaultClusterName = "default";
static const std::string kDefaultZoneId = "default";
static const std::string kDefaultCityId = "default";

static const int kDefaultMaxServerNode = 5;

static const int kDefaultDequeueBulkSize = 100;

static const char *kMetaData = "metadata";
static const char *kZoneId = "zone_id";
static const char *kCityId = "city_id";
static const char *kClusterName = "cluster_name";
static const char *kClusterId = "cluster_id";
static const char *kServiceId = "service_id";
static const char *kNodeIp = "node_ip";

static const char *kAgentServiceName = "agent_service";
static const char *kServiceIp = "service_ip";
static const char *kServicePort = "port";
static const char *kProtocol = "protocol";
static const char *kSSL = "ssl";

static const char *kLoad = "load";
static const char *kAgentCountMax = "agent_count_max";
static const char *kAgentCount = "agent_count";

static const char *kExtend = "extend";
static const char *kCpuRate = "cpu_rate";
static const char *kMemoryRate = "memory_rate";
static const char *kNetworkRate = "network_rate";

static const char *kBlanceConfigNodePath = "/gse/v2/service/data";

static const int kMaxAccessReqQueueSize = 100000;

static const int defaultConnectionLoadWeight = 40;
static const int defaultCpuUsageLoadWeight = 20;
static const int defaultNetUsageLoadWeight = 20;
static const int defaultMemUsageLoadWeight = 20;

enum EnBalanceRequestMsgType
{
    EnAccessReq,
    EnAccessResp,
    EnRedirectAccess
};

class AccessReq
{
public:
    std::string m_zoneId;
    std::string m_cityId;
    std::string m_clusterName;
    std::string m_clusterId;
    std::string m_agentId;

    bool Unmarshal(const char *msg)
    {

        std::string errMsg;
        rj::Document doc;
        std::string strMsg(msg);
        if (!gse::tools::json::LoadDocument(strMsg, doc, errMsg))
        {
            LOG_ERROR("the json is invalid, json(%s), errmsg:%s", msg, errMsg.c_str());
            return false;
        }

        //这里需要整体补充errormsg，当前仅实现逻辑

        auto zoneId = gse::tools::json::JsonHelper<std::string>::GetValue(doc, kZoneId, "");

        if (zoneId.m_isOK)
        {
            m_zoneId = zoneId.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kZoneId);
        }

        auto cityId = gse::tools::json::JsonHelper<std::string>::GetValue(doc, kCityId, "");

        if (cityId.m_isOK)
        {
            m_cityId = cityId.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kCityId);
        }

        auto clusterName = gse::tools::json::JsonHelper<std::string>::GetValue(doc, kClusterName, "");

        if (clusterName.m_isOK)
        {
            m_clusterName = clusterName.m_value;
        }
        else
        {
            LOG_ERROR("%s invalid", kClusterName);
        }
        return true;
    };
};

class LoadWeight
{
public:
    LoadWeight()
        : m_connectionLoadWeight(defaultConnectionLoadWeight),
          m_cpuUsageLoadWeight(defaultCpuUsageLoadWeight),
          m_netUsageLoadWeight(defaultNetUsageLoadWeight),
          m_memUsageLoadWeight(defaultMemUsageLoadWeight)
    {
    }

public:
    int m_connectionLoadWeight;
    int m_cpuUsageLoadWeight;
    int m_memUsageLoadWeight;
    int m_netUsageLoadWeight;
};

class ServerBalanceInfo
{
public:
    std::string m_clusterId;   //没有配置,default
    std::string m_clusterName; //没有配置,default
    std::string m_zoneId;      //如果没有配置,default
    std::string m_cityId;      //如果没有配置,default
    std::string m_hostAddr;
    std::string m_nodeIp;
    std::string m_serviceId; //
    std::string m_protocol;
    uint16_t m_hostPort;
    bool m_ssl;

    float m_cpuUsage;
    float m_netUsage;
    float m_memUsage;

    int m_maxAgentCount;
    int m_agentCount;
    float m_loadWeight;

public:
    bool Equal(ServerBalanceInfo *balanceInfo);
    void Update(const ServerBalanceInfo *balanceInfo);

    bool UnmarshalLoadBalance(const std::string &data);

    float CalcLoadWeight(const LoadWeight &loadWeightCfg);
};

class ClusterList
{
public:
    std::string m_zoneId;
    std::string m_cityId;
    std::string m_clusterName;
    std::string m_clusterId;
    std::string m_agentId;
    std::string m_endpoint;
    float m_nice;
};

class AccessResp
{
public:
    void Pack(const std::vector<ServerBalanceInfo> &serverBalanceInfo, int errCode, const std::string &errmsg)
    {
        // m_clusterList;
        for (auto &iterInfo : serverBalanceInfo)
        {
            ClusterList clusterList;
            clusterList.m_cityId = iterInfo.m_cityId;
            clusterList.m_clusterId = iterInfo.m_clusterId;
            clusterList.m_clusterName = iterInfo.m_clusterName;
            clusterList.m_endpoint = iterInfo.m_nodeIp + ":" + gse::tools::strings::ToString(iterInfo.m_hostPort);
            clusterList.m_zoneId = iterInfo.m_zoneId;
            clusterList.m_nice = iterInfo.m_loadWeight;
            m_clusterList.push_back(clusterList);
        }

        m_errorCode = errCode;
        m_errMsg = errmsg;
    }

    std::string ToJson()
    {
        rapidjson::StringBuffer strBuf;
        rapidjson::Writer<rapidjson::StringBuffer> writer(strBuf);

        writer.StartObject();
        writer.Key("cluster_list");
        writer.StartArray();

        for (auto &iter : m_clusterList)
        {
            writer.StartObject();
            RAPIDJSON_SET_STRING(writer, "endpoint", iter.m_endpoint);
            RAPIDJSON_SET_STRING(writer, "zone_id", iter.m_zoneId);
            RAPIDJSON_SET_STRING(writer, "city_id", iter.m_cityId);
            RAPIDJSON_SET_STRING(writer, "cluster_name", iter.m_clusterName);
            RAPIDJSON_SET_FLOAT(writer, "nice", iter.m_nice);
            writer.EndObject();
        }

        writer.EndArray();

        RAPIDJSON_SET_INT32(writer, "error_code", m_errorCode);
        RAPIDJSON_SET_STRING(writer, "error_message", m_errMsg);
        writer.EndObject();

        std::string data(strBuf.GetString(), strBuf.GetSize());

        return data;
    }

public:
    std::vector<ClusterList> m_clusterList;
    int m_errorCode;
    std::string m_errMsg;
};

using ClusterId = std::string;
using ClusterName = std::string;

using ZoneId = std::string;
using CityId = std::string;

// 调度策略:
// 1. 优先按集群维度，进行调度，如果集群不存在，则查找default集群
// 2. 集群内按照zone/city 维度进行调度
// 3. 如果当前city不匹配，则选择zone下其他city
// 4. 如果zone下无可用city，选择其他zone下city，优选default
// 5. 如果当前集群下无可用，选择default集群下，2->4 的调度选择策略

//每个Zone 下的City List

using CallBackFunc = std::function<void(const char *, int)>;

class Task
{
public:
    std::shared_ptr<gse::net::Message> m_message;
    CallBackFunc m_callback;
};

class CityPreferServerList
{
public:
    virtual ~CityPreferServerList();
    void Push(ServerBalanceInfo *serverBalanceInfo);
    uint64_t Size();
    ServerBalanceInfo *operator[](int index);

public:
    //各城市的负载信息列表，CityId维护
    std::list<ServerBalanceInfo *> m_cityAllNodes;
};

class CityEndPoint
{
public:
    explicit CityEndPoint(const std::string &cityId);
    virtual ~CityEndPoint();
    bool Push(ServerBalanceInfo *serverBalanceInfo);

    bool FindPreferServerEndPoint(std::vector<ServerBalanceInfo> &serverBanlanceInfo);

public:
    CityPreferServerList m_serverList;
    std::string m_cityId;
};

class ZoneEndPoint
{
public:
    explicit ZoneEndPoint(const std::string &zoneId);
    ~ZoneEndPoint();
    bool FindCity(const std::string &cityId, CityEndPoint *&cityEndPoint);
    bool PushCity(const std::string &cityId, CityEndPoint *cityEndPoint);

    bool FindPreferCityEndPoint(const std::string &cityId, CityEndPoint *&cityEndPoint);
    bool GetAllCity(std::vector<std::string> &cityIds);
    size_t CitySize();

public:
    std::string m_zoneId;
    std::map<CityId, CityEndPoint *> m_zonePreferCityList;
};

//每个cluster 下的Zone 列表，当前cluster下的所有Zone信息
class ZonePreferServerList
{
public:
    std::map<ZoneId, ZoneEndPoint *> m_zonePreferServerList;
    std::string ZoneListName();
};

class ClusterEndPoint
{
public:
    explicit ClusterEndPoint(const std::string &clusterName);
    virtual ~ClusterEndPoint();
    bool PushZone(const std::string &zoneId, ZoneEndPoint *zoneEndPoint);
    bool FindZone(const std::string &zoneId, ZoneEndPoint *&zoneEndPoint);
    bool RemoveZone(const std::string &zoneId);
    std::string ZoneListName();
    size_t ZoneSize();

    bool FindPreferZoneEndPoint(const std::string &zoneId, ZoneEndPoint *&zoneEndPoint);

public:
    std::map<ZoneId, ZoneEndPoint *> m_preferZoneList;
    std::string m_clusterName;
};

class ClusterPreferServerList
{

public:
    ClusterPreferServerList();
    virtual ~ClusterPreferServerList();

    bool Find(const std::string &clusterName, ClusterEndPoint *&clusterEndPoint);
    bool Push(const std::string &clusterName, ClusterEndPoint *clusterEndPoint);
    void Swap(ClusterPreferServerList &clusterPreferServerList);

    bool FindPreferClusterEndPoint(const std::string &clusterName, ClusterEndPoint *&clusterEndPoint);

public:
    //集群下所有server负载信息
    // 优先集群，agent的集群名称和server 的集群名称一致
    // default -> zone1,zone2

    std::map<ClusterName, ClusterEndPoint *> m_clusterPreferServerList;
};

class LoadBalanceSchedule
{
public:
    LoadBalanceSchedule();
    LoadBalanceSchedule(const std::string &zkHost, const std::string &zkauth);
    virtual ~LoadBalanceSchedule();

    bool Start();
    void Stop();

    bool EnqueuAccessRequest(CallBackFunc callback, std::shared_ptr<gse::net::Message> msg);

private:
    // 10s get zk all service list
    void UpdateLoadBalance();

    void ScheduleAccessRequest(const AccessReq &accessReq, std::vector<ServerBalanceInfo> &nodes);
    void ProcessAccessRequest(const Task &task);
    void ProcessQueueMsg();

    bool EnqueuAccessRequest(const Task &task);

private:
    //队列处理请求
    moodycamel::ConcurrentQueue<Task> *m_accessRequestMessageQueuePtr;

    // schedule config parameter
    int m_maxServerNode; // max available  server node

    int m_timeoutSecond;
    std::shared_ptr<gse::eventloop::EventLoopThread> m_evLoopThread;
    std::shared_ptr<gse::eventloop::TimerEventWatcher> m_timerUpdateLoadBalance; // get load
    std::shared_ptr<gse::eventloop::TimerEventWatcher> m_timerProcessAccess;     // get load

    Conftor *m_zkClient;
    ClusterPreferServerList m_preferClusterList;
    std::string m_clusterName;
    std::string m_clusterId;
    std::string m_zkHost;
    std::string m_zkAuth;
};

} // namespace data
} // namespace gse

#endif // _GSE_DATA_LOADBALANCE_SCHEDULE_H_
