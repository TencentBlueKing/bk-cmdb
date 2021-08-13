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

#include <string>
#include <map>
#include <json/json.h>

#include "safe/lock.h"

#include "conf/confItem.h"
#include "config/configfile.h"
#include "tools/json_property.hpp"

using namespace gse::config;
using namespace gse::tools::json;

namespace gse { 
namespace dataserver {


class DataProcessConfig
{
public:
    std::string m_configZk;
    std::string m_eventZk;
    std::string m_zkauth;
    std::string m_clusterName;
    std::string m_instanceId;
    std::string m_dataflowfile;
    std::string m_regionID;
    std::string m_cityID;
    std::string m_watchpath;
    std::string m_certPath;
    std::string m_passwdfile;
    std::string m_serverIp;

    std::string m_channelidZkHost;
    std::string m_channelidZkAuth;

    std::string m_runtimePath;
    std::string m_logPath;
    std::string m_logLevel;
    std::string m_pidFilePath;
    std::string m_promethusBindIp;
    unsigned int m_promethusPort;
    unsigned int m_filesize;
    unsigned int m_filenum;

    //------------------------
    bool m_enableOps;

public:
    DataProcessConfig();
    DataProcessConfig(const std::string &configfile);
    ~DataProcessConfig();
public:
    //static const MetaType *properties[];

    std::string get_prometheus_http_svr_ip();
    uint16_t get_prometheus_datasvr_port();
    std::string getDataSvrClusterName();
    std::string getDataSvrInstanceId();
    std::string getDataSvrIp();
    bool getOpsFlag();

private:
    int parseConfigFile(const std::string &filename);
};

class DataConf
{
public:
    DataConf();
    virtual ~DataConf();

public:
    int DecodeConf(const string& text);
    DataFlowConf* GetDataFlowConf();
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
    FilterConf* findFilterConf(std::string& filterName);
    ExporterConf* findExporterConf(std::string& exporterName);
    ReceiverConf* findReceiverConf(std::string& receiverName);

private:
    gse::safe::MutexLock m_lockDataFlowConf;
    DataFlowConf* m_dataflowConf;
    std::map<std::string, ReceiverConf*> m_receiversConf;
    std::map<std::string, ExporterConf*> m_exportersConf;
    std::map<std::string, FilterConf*> m_filtersConf;
};

}
}
#endif
