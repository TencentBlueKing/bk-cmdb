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

#ifndef _GSE_DATA_DATACONF_H_
#define _GSE_DATA_DATACONF_H_

#include <json/json.h>
#include <map>
#include <string>

#include "safe/lock.h"

#include "conf/confItem.h"
#include "tools/json_property.hpp"

using namespace gse::tools::json;

namespace gse {
namespace data {

static const int kDefaultZkTimeoutMs = 300000;
static const int kDefaultNetDevMaxSpeed = 10000;
static const int kDefaultMaxAgentCount = 50000;
static const std::string kDefaultDevName = "eth1";

static const int kDefaultAgentCountLoadWeight = 40;
static const int kDefaultCpuUsageLoadWeightLoadWeight = 20;
static const int kDefaultNetUsageLoadWeightLoadWeight = 20;
static const int kDefaultMemUsageLoadWeightLoadWeight = 20;

class PulsarConfig
{
public:
    int m_maxPendingMessages;
    int m_sendTimeoutMs;
    int m_compressionType;
    int m_maxPendingMessagesAcrossPartitions;
};

class LoggerConfig
{
public:
    LoggerConfig()
        : m_fileCount(100), m_fileSize(100), m_rotate(0), m_flushIntervalMs(DFT_LOG_FLUSH_MSEC) {}

public:
    std::string m_path;
    std::string m_level;
    int m_fileSize;
    int m_fileCount;
    int m_rotate;
    int m_flushIntervalMs;
};

class ZooKeeperConfig
{
public:
    ZooKeeperConfig()
        : m_timeout(kDefaultZkTimeoutMs)
    {
    }

    std::string m_serviceDiscoverZkHost;
    std::string m_serviceDiscoverZkAuth;
    int m_timeout;

    std::string m_channelIdConfigZkHost;
    std::string m_channelIdConfigZkAuth;
};

class MetricsConfig
{
public:
    MetricsConfig()
        : m_promethusListenPort(59402), m_threadNum(1) {}

public:
    std::string m_promethusBindIP;
    uint16_t m_promethusListenPort;
    int m_threadNum;
};

class OpsConfig
{
public:
    OpsConfig()
        : m_enableOps(false) {}

public:
    bool m_enableOps;

    std::vector<std::string> m_opsServerAddrs;
    int m_port;
};

class ClusterInfoConfig
{
public:
    ClusterInfoConfig()
        : m_serviceName("data")
    {
    }

public:
    std::string m_clusterName;
    std::string m_clusterId;
    std::string m_serviceName;
    std::string m_instanceId;
    std::string m_advertiseIp;
    std::string m_zoneId;
    std::string m_cityId;
};

class BalanceLoadWeightConfig
{
public:
    BalanceLoadWeightConfig()
        : m_agentCountLoadWeight(kDefaultAgentCountLoadWeight),
          m_cpuUsageLoadWeight(kDefaultCpuUsageLoadWeightLoadWeight),
          m_netUsageLoadWeight(kDefaultNetUsageLoadWeightLoadWeight),
          m_memUsageLoadWeight(kDefaultMemUsageLoadWeightLoadWeight),
          m_netDevName(kDefaultDevName),
          m_netDevMaxSpeed(kDefaultNetDevMaxSpeed),
          m_maxAgentCount(kDefaultMaxAgentCount)
    {
    }

    int m_agentCountLoadWeight;
    int m_cpuUsageLoadWeight;
    int m_memUsageLoadWeight;
    int m_netUsageLoadWeight;
    int m_maxAgentCount;
    std::string m_netDevName;
    int m_netDevMaxSpeed;
};

class TGLogConfig
{
public:
    std::string m_watchPath;
};

class ChannelIdApiServiceConfig
{
public:
    ChannelIdApiServiceConfig()
        : m_valid(false), m_threadNum(1)
    {
    }

public:
    std::string m_bindIp;
    int m_bindPort;
    int m_threadNum;
    bool m_valid;
};

class DataIdZKConfig
{
public:
    std::string m_zkhost;
    std::string m_zkAuth;
    std::string m_dataIdPath;
    std::string m_storagePath;
};

class DataProcessConfig
{
public:
    std::string m_pidFilePath;
    LoggerConfig m_loggerConfig;
    ZooKeeperConfig m_zooKeeperConfig;
    MetricsConfig m_metricsConfig;
    OpsConfig m_opsConfig;
    ClusterInfoConfig m_clusterInfoConfig;
    TGLogConfig m_tglogConfig;
    ChannelIdApiServiceConfig m_channelIdApiServiceConfig;
    std::vector<DataIdZKConfig> m_dataIdZk;
    BalanceLoadWeightConfig m_balanceConfig;
    std::string m_runtimePath;

    //------------------------
    bool m_enableOps;
    std::string m_configFilePath;

public:
    DataProcessConfig();
    DataProcessConfig(const std::string& configfile);
    bool LoadConfig(const std::string& configfile);
    ~DataProcessConfig();

public:
    // static const MetaType *properties[];

    std::string GetPrometheusBindIp();
    int GetPrometheusListenerPort();
    std::string GetDataSvrClusterName();
    std::string GetDataSvrInstanceId();
    std::string GetAdvertiseIp();
    bool GetOpsFlag();

    ZooKeeperConfig* GetZookeeperConfig();
    std::vector<DataIdZKConfig>* GetDataIdZkConfig();

    ChannelIdApiServiceConfig* GetChannelIdApiServiceConfig();

public:
    int ParseConfigFile(const std::string& filename);
};

class DataConf
{
public:
    DataConf();
    virtual ~DataConf();

public:
    int DecodeConf(const string& text);
    DataFlowConf* GetDataFlowConf();

public:
    FilterConf* findFilterConf(std::string& filterName);
    ExporterConf* findExporterConf(std::string& exporterName);
    ReceiverConf* findReceiverConf(std::string& receiverName);

protected:
private:
    int parseDataFlowConf(const string& confText);
    int parseChannelConf(DataFlowConf* pDataFlowConf, const Json::Value& cfgValue);
    int parseOpsConf(DataFlowConf* pDataFlowConf, const Json::Value& cfgValue);
    int parseReceiverConf(const Json::Value& cfgValue);
    int parseExporterConf(const Json::Value& cfgValue);
    LogExporterConf* parseLogExporterConf(const std::string& name, const Json::Value& cfgValue);
    KafkaExporterConf* parseKafkaExporterConf(const std::string& name, const Json::Value& cfgValue);
    RedisExporterConf* parseRedisExporterConf(const std::string& name, const Json::Value& cfgValue);
    BkDataExporterConf* parseBkDataExporterConf(const std::string& name, const Json::Value& cfgValue);
    DSProxyExporterConf* parseDSProxyExporterConf(const std::string& name, const Json::Value& cfgValue);
    PulsarExporterConf* parsePulsarExporterConf(const std::string& name, const Json::Value& cfgValue);
    int parseFilterConf(const Json::Value& cfgValue);
    bool isExitReceiver(const std::string& receiverName);
    bool isExitExporter(const std::string& exporterName);
    bool isExitFilter(const std::string filterName);
    bool insertReceiver(const std::string& receiverName, ReceiverConf* pReceiverConf);
    bool insertExporter(const std::string& exporterName, ExporterConf* pExporterConf);
    bool insertFilter(const std::string& filterName, FilterConf* pFilterConf);
    void clearReceiverConfMap();
    void clearExporterConfMap();
    void clearFilterConfMap();

private:
    gse::safe::MutexLock m_lockDataFlowConf;
    DataFlowConf* m_dataflowConf;
    std::map<std::string, ReceiverConf*> m_receiversConf;
    std::map<std::string, ExporterConf*> m_exportersConf;
    std::map<std::string, FilterConf*> m_filtersConf;
};

} // namespace data
} // namespace gse
#endif
