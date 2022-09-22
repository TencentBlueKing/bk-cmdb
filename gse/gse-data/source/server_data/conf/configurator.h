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

#ifndef _GSE_DATA_CONFIGURATOR_H_
#define _GSE_DATA_CONFIGURATOR_H_

#include <string.h>
#include <string>

#include "conftor/conftor.h"
#include "dataconf.h"
#include "datastruct/safe_map.h"
#include "datastruct/safe_queue.h"
#include "eventthread/event_thread.h"
#include "tools/net.h"
#include "tools/strings.h"
#include "tools/thread.h"

#include "balance_config.h"
#include "conf/channel_id_config.h"
#include "discover/zkapi/zk_api.h"
#include "tools/json_helper.hpp"
#include "tools/json_property.hpp"

namespace gse {
namespace data {

#define DATA_BASE_PATH "/gse/config/server/dataserver"
#define DATA_OPS_SERVICE_CONF_PATH "/gse/v2/config/data/ops_config"
#define ZK_SERIVCE_NODE_PATH "/gse/v2/service/data"

enum ConfItemFlag
{
    CONFITEMFLAG_UNSET = -1,
    CONFITEMFLAG_DATACONF = 0,
    CONFITEMFLAG_CHANNELID_LIST = 1,
    CONFITEMFLAG_CHANNELID_CONFIG = 2,
    CONFITEMFLAG_CHANNELID_FILTER_CONFIG = 3,
    CONFITEMFLAG_CHANNELID_EXPORTER_CONFIG = 4,
    CONFITEMFLAG_CHANNELID_CHANNEL_CONFIG = 5,
    CONFITEMFLAG_STORAGE = 6,
    CONFITEMFLAG_STORAGE_CONFIG = 7,
    CONFITEMFLAG_DATAID = 8,
    CONFITEMFLAG_DATAID_CONFIG = 9,
    CONFITEMFLAG_BALANCE_CONFIG = 10,

    CONFITEMFLAG_SCHEDULE_CONFIG = 11,
    CONFITEMFLAG_DATAID_FROM_BKDATA = 12,
    CONFITEMFLAG_DATAID_CONFIG_FROM_BKDATA = 13,
    CONFITEMFLAG_STORAGE_FROM_BKDATA = 14,
    CONFITEMFLAG_STORAGE_CONFIG_FROM_BKDATA = 15,

    CONFITEMFLAG_STREAMTO_CONFIG_LIST,
    CONFITEMFLAG_STREAMTO_CONFIG_VALUE,

    CONFITEMFLAG_TGLOG_CHANNEL_ID_VALUE,

    CONFITEMFLAG_PLAT_ID_CONFIG_LIST,
    CONFITEMFLAG_PLAT_ID_CONFIG_VALUE,

    CONFITEMFLAG_OPS_SERVICE_CONFIG_LIST,
    CONFITEMFLAG_OPS_SERVICE_CONFIG_VALUE,
    CONFITEMFLAG_MAX_NUM
};

inline int split_storage_index(const char *target)
{

    if (NULL == target)
    {
        return 0;
    }
    std::string data(target);

    char *ptr = strtok(const_cast<char *>(data.c_str()), "_");
    while (NULL != ptr)
    {
        int val = atoi(ptr);
        if (val > 0)
        {
            return val;
        }
        ptr = strtok(NULL, "_");
    }
    return 0;
}

inline uint32_t split_channel_id(const char *target)
{
    if (NULL == target)
    {
        return 0;
    }
    std::string data(target);

    std::size_t pos = data.find_last_of("/");
    if (pos != data.npos)
    {
        return gse::tools::strings::StringToUint32(data.substr(pos + 1));
    }
    return 0;
}

typedef struct _EventCallbackParams
{
    IDToStorage::WatchEventFunc m_eventCallbackFunc;
    void *m_ptrCallbackArgs;
} EventCallbackParams;

class ServiceNodeMeta
{
public:
    ServiceNodeMeta() {}

    rapidjson::Document ToJsonValue()
    {
        rapidjson::Document nodeMeta(rapidjson::kObjectType);
        nodeMeta.AddMember("zone_id", m_zoneId, nodeMeta.GetAllocator());
        nodeMeta.AddMember("city_id", m_cityId, nodeMeta.GetAllocator());
        nodeMeta.AddMember("cluster_name", m_clusterName, nodeMeta.GetAllocator());
        nodeMeta.AddMember("cluster_id", m_clusterId, nodeMeta.GetAllocator());
        nodeMeta.AddMember("service_name", m_serviceName, nodeMeta.GetAllocator());
        nodeMeta.AddMember("service_id", m_serviceId, nodeMeta.GetAllocator());
        nodeMeta.AddMember("node_ip", m_advertiseIp, nodeMeta.GetAllocator());

        return nodeMeta;
    }

public:
    std::string m_zoneId;
    std::string m_cityId;
    std::string m_clusterName;
    std::string m_clusterId;
    std::string m_serviceName;
    std::string m_serviceId;
    std::string m_advertiseIp;

public:
    static const MetaType *properties[];
};

class ServiceNode
{
public:
    ServiceNode()
        : m_ssl(false), m_port(58625) {}

    rapidjson::Document ToJsonValue()
    {
        rapidjson::Document docValue(rapidjson::kObjectType);

        docValue.AddMember("service_ip", m_serviceIP, docValue.GetAllocator());
        docValue.AddMember("protocol", m_protocol, docValue.GetAllocator());
        docValue.AddMember("port", m_port, docValue.GetAllocator());
        docValue.AddMember("ssl", m_ssl, docValue.GetAllocator());
        if (m_proto > 0)
        {
            docValue.AddMember("proto", m_proto, docValue.GetAllocator());
        }

        return docValue;
    }

    std::string ToJsonStr()
    {
        auto jsonDoc = ToJsonValue();

        rapidjson::StringBuffer strBuffer;
        rapidjson::Writer<rapidjson::StringBuffer> writer(strBuffer);

        jsonDoc.Accept(writer);
        LOG_DEBUG("service node json(%s)", strBuffer.GetString());
        return std::string(strBuffer.GetString());
    }

public:
    std::string m_serviceIP;
    uint32_t m_port;
    std::string m_protocol;
    int m_proto;
    bool m_ssl;

public:
    static const MetaType *properties[];
};

class NodeLoadBalance
{
public:
    NodeLoadBalance()
        : m_cpuUsage(0.0),
          m_memUsage(0.0),
          m_netUsage(0.0),
          m_maxConnectionCount(0),
          m_connectionCount(0)
    {
    }
    std::string ToJsonStr()
    {
    }

    rapidjson::Document ToJsonValue()
    {
        rapidjson::Document docExtandValue(rapidjson::kObjectType);
        rapidjson::Document loadValue(rapidjson::kObjectType);

        loadValue.AddMember("agent_count_max", m_maxConnectionCount, loadValue.GetAllocator());
        loadValue.AddMember("agent_count", m_connectionCount, loadValue.GetAllocator());

        docExtandValue.AddMember("cpu_rate", m_cpuUsage, loadValue.GetAllocator());
        docExtandValue.AddMember("memory_rate", m_memUsage, loadValue.GetAllocator());
        docExtandValue.AddMember("network_rate", m_netUsage, loadValue.GetAllocator());

        loadValue.AddMember("extend", docExtandValue, loadValue.GetAllocator());

        return loadValue;
    }

public:
    double m_cpuUsage;
    double m_memUsage;
    double m_netUsage;
    uint64_t m_maxConnectionCount;
    uint64_t m_connectionCount;
};

class ServerNode
{
public:
    ServerNode();
    std::string ToJson()
    {
    }

private:
    ServiceNode m_seviceNode;
    ServiceNodeMeta m_serviceMate;
    NodeLoadBalance m_loadInfo;
};

using ServiceID = int;

class OpsServiceConfig
{
public:
    OpsServiceConfig()
        : m_serivceId(0),
          m_channelId(-1)
    {
    }

    int ParseJsonConfig(const std::string &jsonConfig)
    {
        std::string errMsg;
        rj::Document doc;
        if (!gse::tools::json::LoadDocument(jsonConfig.c_str(), doc, errMsg))
        {
            LOG_ERROR("failed to load json, input json:%s", jsonConfig.c_str());
            return GSE_ERROR;
        }

        auto channelIdVale = gse::tools::json::JsonHelper<int>::GetValue(doc, "channel_id", 0);

        if (channelIdVale.m_isOK)
        {
            m_channelId = channelIdVale.m_value;
        }
        else
        {
            LOG_ERROR("failed to get channel_id value, errmsg:%s", channelIdVale.m_errMsg.c_str());
            return GSE_ERROR;
        }

        auto serviceIdValue = gse::tools::json::JsonHelper<int>::GetValue(doc, "service_id", -1);

        if (serviceIdValue.m_isOK)
        {
            m_serivceId = serviceIdValue.m_value;
        }
        else
        {
            LOG_ERROR("failed to get service_id value, errmsg:%s", serviceIdValue.m_errMsg.c_str());
            return GSE_ERROR;
        }
        return GSE_SUCCESS;
    }

public:
    int m_serivceId;
    int m_channelId;
};

class Configurator : public IDToStorage
{
public:
    typedef void (*PTRConfigWatchFunc)(WatchConfItem &confItem);

public:
    Configurator();
    virtual ~Configurator();

public:
    // Init init self from remote config center
    int Init(std::shared_ptr<DataProcessConfig> configPtr);
    void DeleteEphemeralZkNodes();
    int StartConftor();
    void StopConftor();

public:
    // static functions
    static void reportBalance(evutil_socket_t fd, short what, void *args);
    static uint32_t getDefaultTglogChannelId();
    static void ReportSystemLoad(evutil_socket_t fd, short what, void *args);

public:
    ChannelIdManager *GetPlatIdManager();
    DataFlowConf *GetDataFlowConf();

    void SetChannelIdZkClient(std::shared_ptr<gse::discover::zkapi::ZkApi> zkClient, bool zkAcl);
    void SetDiscoverZkClient(std::shared_ptr<gse::discover::zkapi::ZkApi> zkClient, bool zkAcl);

public:
    void GetChannelsConfigValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig *ptr_channelid_config);

    void GetChannelConfigValueCallBack(std::string &path, int rc, std::string &property_value, void *ptr_callbackobj);
    void ParseMetaConfig(std::string &path, int rc, std::string &property_value, ChannelIDConfig *ptr_channelid_config);
    void GetFiltersValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig *ptr_channelid_config);
    void UpdateFilterValue(std::string &path, int rc, std::string &property_value, void *ptr_callbackobj);

    int CreateV2ServiceNode();
    int CreateOpsServiceConfigNode();
    int GetChannelIdByOpsServiceId(int serviceId);

private:
    // InitFromLocalConfig init self from local config
    // Maybe this function will be replaced by another way
    int initFromLocalConfig(const std::string &dataFlowConf);

    static void watchConfCallBack(WatchConfItem &confItem, void *lpWatcher);

    static void channelNodeChangeCallBack(WatchConfItem &confItem, void *lpWatcher);

    static void updateDataConf(WatchConfItem &confItem, void *lpWatcher);
    static void handleChannelIdChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher);
    static void handlePlatIdChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher);
    static void updateChannelIDConfig(WatchConfItem &confItem, void *lpWatcher);
    static void updatePlatIdConfig(WatchConfItem &confItem, void *lpWatcher);

    static void updateBalanceConfig(WatchConfItem &confItem, void *lpWatcher);
    static void handleStreamChildListChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher);
    static void handleStreamValueChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher);
    static void handleTglogChannelIdChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher);

    static void HandleOpsConfigChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher);
    static void UpdateOpsConfigValue(WatchConfItem &confItem, void *lpWatcher);

    static void GetChannelsListResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void *ptr_callbackobj);
    static void GetFiltersResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void *ptr_callbackobj);

    static void ChannelValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj);
    static void GetMetaValueCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj);
    static void FilterValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj);

    template <typename T>
    void ToJsonObj(rapidjson::Writer<rapidjson::StringBuffer> &writer, T object);

private:
    // update config from zk

    int watchConfigsFromZK();
    int watchDataFlow();
    // watchStorageFromBKData 仅用于兼容存量bkdata 里关于dataid 的配置，未来需要被废弃
    int watchChannelID();
    int watchPlatID();
    int WatchTglogChannelId();
    int watchOpsServiceConfig();

    int watchBalanceConfig();
    void watchChannelIdConfig(WatchConfItem &confItem);
    void watchPlatIdConfig(WatchConfItem &confItem);
    int updateChannelIDConfigFromZK(uint32_t channelID);

    int updateDataFlowConf(const std::string &context);
    int updateChannelID(uint32_t channelId, const std::string &context);
    int DeleteChannelID(uint32_t channelId);
    bool FindChannelID(uint32_t channelId);

    int updatePlatID(uint32_t plat_id);
    int deletePlatID(uint32_t platid);

    // updateDataID compatible with existing configurations about dataid (now it is channel id)
    int updateDataID(uint32_t dataID, const std::string &context);

    int updateDataServerBalanceConfig(const string &path, const std::string &context);
    int updateDataServerLoadBalanceNode();

    void WatchOpsConfigValueChangeEvt(const std::string &path);
    void WatchOpsConfigChildNodeChangeEvt(const std::string &path);
    void updateLocationFromZK();

    void ReportLoadInfo();

    // private ,support async zk api get channelid
    bool CanSendUpdateEventMsg(ChannelIDConfig *ptr_channelid_config);
    void SendUpdateChannelIdConfigEventMsg(ChannelIDConfig *ptr_channelid_config);
    bool IsFinishAndFailed(ChannelIDConfig *ptr_channelid_config);
    void DoUpdateChannelIdConfigRequest(ChannelIDConfig *ptr_channelid_config);
    void DelayFreeChannelIdConfig(ChannelIDConfig *ptr_channelid_config);

    int GetOpsConfig(const std::string &zkNodePath, std::string &configValue);

    void UpdateOpsServiceConfig(const std::string &cfgJson);
    std::string BuildServiceNodeJson();
    bool GetServiceConfig(const std::string &serviceName, ServiceNode &serviceConfig);

private:
    // read channelid config
    bool readMetadata(const std::string &metadataPath, ChannelIDConfig *ptr_channelid_config, std::string &errorMsg);
    bool readFilters(const std::string &fileterPath, ChannelIDConfig *ptr_channelid_config, std::string &errorMsg);
    bool readChannels(const std::string &channelPath, ChannelIDConfig *ptr_channelid_config, std::string &errorMsg);

private:
    // set config into zk

    int createBaseConfItem();

private:
    Conftor *m_conftor;
    DataConf *m_dataConf;

    std::shared_ptr<SystemResourceMonitor> m_systemResourceMonitor;
    std::shared_ptr<SystemConnectionMonitor> m_systemConnectionMonitor;

    std::string m_selfIp;
    std::string m_clusterName;
    std::string m_instanceId;
    std::string m_localDataFlowConfig;
    std::string m_watchPath;

    std::shared_ptr<gse::discover::zkapi::ZkApi> m_channelIdZkClient;
    bool m_channelIdZkAcl;

    std::shared_ptr<gse::discover::zkapi::ZkApi> m_discoverZkClient;
    bool m_discoverZkAcl;

private:
    Json::Value m_cfgLoadBalanceValue;
    gse::datastruct::SafeQueue<ChannelIDConfig *> m_toDeleteChannelIDS;
    gse::datastruct::SafeQueue<ChannelIdExporterConfig *> m_toDeleteStreamToIDS;

    gse::datastruct::SafeMap<ChannelIDType, ChannelIDConfig *> m_channelIDS;
    gse::datastruct::SafeMap<StreamToIDType, ChannelIdExporterConfig *> m_streamToIdConfig;

    gse::datastruct::SafeMap<ServiceID, OpsServiceConfig> m_opsConfig;

    std::string m_zoneID;
    std::string m_cityID;
    std::string m_ethName;
    EventThread *m_ptrEventThread;

private:
    float m_cpuUsage;
    float m_hostPerformance;
    float m_hostLoadBance;
    float m_loadweight;
    float m_dloadweight;

    std::string m_bkdataZK;

    Conftor *m_channelIdZKConftor;

    ChannelIdManager *m_platIdManager;

private:
    void (*m_ptrConfigFunc[CONFITEMFLAG_MAX_NUM])(WatchConfItem &confItem, void *lpWatcher);
    std::vector<EventCallbackParams *> m_callbacks;
};

struct ChannelIdZkCallBackObj
{
    Configurator *m_ptrConfigurator;
    ChannelIDConfig *m_ptrChannelidConfig;

    ChannelIdZkCallBackObj()
    {
        m_ptrConfigurator = NULL;
        m_ptrChannelidConfig = NULL;
    }
};

} // namespace data
} // namespace gse
#endif //_GSE_DATA_CONFIGURATOR_H_
