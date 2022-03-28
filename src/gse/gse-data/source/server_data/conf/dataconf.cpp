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

#include "dataconf.h"
#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/macros.h"
#include <json/json.h>

#include "conf/confItem.h"

namespace gse
{
namespace dataserver
{

using namespace std;

////////////////////////////////////////////////////////////////////////////////////////////////////////
// class DataProcessConfig
////////////////////////////////////////////////////////////////////////////////////////////////////////

DataProcessConfig::DataProcessConfig()
{
    m_enableOps = false;
    m_filenum = 100;
    m_filesize = 100;
    m_promethusPort = 59402;
}

DataProcessConfig::DataProcessConfig(const std::string &configfile)
{
    parseConfigFile(configfile);
}

DataProcessConfig::~DataProcessConfig()
{
}

std::string DataProcessConfig::getDataSvrClusterName()
{
    return m_clusterName;
}

std::string DataProcessConfig::getDataSvrIp()
{
    return m_serverIp;
}

std::string DataProcessConfig::getDataSvrInstanceId()
{
    return m_instanceId;
}

int DataProcessConfig::parseConfigFile(const std::string &filename)
{
    std::string config_content;
    int ret_value = gse::tools::filesystem::ReadFromFile(filename, config_content);
    if (GSE_SUCCESS != ret_value)
    {
        LOG_ERROR("it is failed to read the config by the absolute path ( %s ) ", SAFE_CSTR(filename.c_str()));
        return ret_value;
    }

    LOG_DEBUG("Base config file[%s]", filename.c_str());
    Json::Reader reader(Json::Features::strictMode());
    Json::Value cfgValue;
    if (!reader.parse(config_content, cfgValue, false))
    {
        LOG_ERROR("fail to parse dataflow configure by json. the configure text is [%s]", SAFE_CSTR(config_content.c_str()));
        return GSE_JSON_INVALID;
    }

    m_dataflowfile = cfgValue.get("dataflow", "").asString();
    m_eventZk = cfgValue.get("eventzkhost", "").asString();
    m_configZk = cfgValue.get("zkhost", "").asString();
    m_zkauth = cfgValue.get("zkauth", "").asString();
    m_serverIp = cfgValue.get("datasvrip", "0.0.0.0").asString();
    m_clusterName = cfgValue.get("clustername", "").asString();

    m_channelidZkHost = cfgValue.get("channelidzkhost", "").asString();
    m_channelidZkAuth = cfgValue.get("channelidzkauth", "").asString();

    m_promethusBindIp = cfgValue.get("prometheus_http_svr_ip", "0.0.0.0").asString();
    m_promethusPort = cfgValue.get("prometheus_datasvr_port", 59402).asInt();

    m_regionID = cfgValue.get("dftregid", "").asString();
    m_cityID = cfgValue.get("dftcityid", "").asString();
    m_passwdfile = cfgValue.get("password_keyfile", "").asString();
    m_certPath = cfgValue.get("cert", "").asString();
    m_watchpath = cfgValue.get("watchpath", "").asString();
    m_logPath = cfgValue.get("log", "").asString();
    m_pidFilePath = cfgValue.get("pid", "").asString();
    m_filesize = cfgValue.get("logfilesize", 100).asInt();
    m_filenum = cfgValue.get("logfilenum", 100).asInt();
    m_logLevel = cfgValue.get("level", "debug").asString();


    m_enableOps = cfgValue.get("enableops", false).asBool();
    return GSE_SUCCESS;
}

std::string DataProcessConfig::get_prometheus_http_svr_ip()
{
    return m_promethusBindIp;
}

uint16_t DataProcessConfig::get_prometheus_datasvr_port()
{
    return m_promethusPort;
}

bool DataProcessConfig::getOpsFlag()
{
    return m_enableOps;
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// class DataConf
////////////////////////////////////////////////////////////////////////////////////////////////////////
DataConf::DataConf()
{
    m_dataflowConf = NULL;
}

DataConf::~DataConf()
{
    if (m_dataflowConf != NULL)
    {
        delete m_dataflowConf;
        m_dataflowConf = NULL;
    }

    clearReceiverConfMap();
    clearExporterConfMap();
    clearFilterConfMap();
}

int DataConf::DecodeConf(const string &text)
{
    //string tt = "{\"receiver\":[{\"name\":\"r_bkdata\",\"protocol\":1,\"bind\":\"127.0.0.1\",\"port\":58625,\"cert\":\"etc/cert\",\"protostack\":1},{\"name\":\"r_tglog\",\"protocol\":2,\"bind\":\"127.0.0.1\",\"port\":58625,\"cert\":\"\",\"protostack\":0}],\"exporter\":[{\"name\":\"e_log\",\"type\":1,\"path\":\"logs/data/\",\"filename\":\"%Y%m%d-%5N.log\",\"maxlogsize\":1000,\"maxfilenum\":100},{\"name\":\"e_bkdata_withfilters\",\"type\":8,\"filters\":[\"FilterRegex\"],\"zk\":\"\"},{\"name\":\"e_bkdata\",\"type\":8,\"filters\":[\"FilterRegex\"],\"zk\":\"\"}],\"filters\":[{\"name\":\"FilterRegex\",\"matcher\":{\"conditions\":[{\"delimiter\":\"|\",\"fieldindex\":0,\"word\":\"test.*\"}]}}],\"channel\":[{\"name\":\"bkdata\",\"decode\":0,\"receiver\":\"r_tglog\",\"exporter\":[\"e_log\"]}]}";
    return parseDataFlowConf(text);
}

DataFlowConf *DataConf::GetDataFlowConf()
{
    if (NULL == m_dataflowConf)
    {
        return NULL;
    }

    m_lockDataFlowConf.Lock();
    DataFlowConf *pDataFlowConf = new DataFlowConf(*m_dataflowConf);
    m_lockDataFlowConf.Unlock();

    return pDataFlowConf;
}

int DataConf::parseDataFlowConf(const string &confText)
{
    int iRet = GSE_SUCCESS;

    // begin to decode dataflow config
    LOG_DEBUG("begin to decode dataflow config[%s]", SAFE_CSTR(confText.c_str()));
    Json::Reader reader(Json::Features::strictMode());
    Json::Value cfgValue;
    if (!reader.parse(confText, cfgValue, false))
    {
        LOG_ERROR("fail to parse dataflow configure by json. the configure text is [%s]", SAFE_CSTR(confText.c_str()));
        return GSE_JSON_INVALID;
    }

    DataFlowConf *pDataFlowConf = new DataFlowConf();
    // parse filters
    if (cfgValue.isMember("filters") && cfgValue["filters"].isArray())
    {
        Json::Value filterValue = cfgValue["filters"];
        parseFilterConf(filterValue);
    }
    // parse receiver
    if (cfgValue.isMember("receiver") && cfgValue["receiver"].isArray())
    {
        Json::Value receiverValue = cfgValue["receiver"];
        parseReceiverConf(receiverValue);
    }

    // parse exporter
    if (cfgValue.isMember("exporter") && cfgValue["exporter"].isArray())
    {
        Json::Value exporterValue = cfgValue["exporter"];
        parseExporterConf(exporterValue);
    }

    // parse channel
    if (cfgValue.isMember("channel") && cfgValue["channel"].isArray())
    {
        Json::Value channelValue = cfgValue["channel"];
        parseChannelConf(pDataFlowConf, channelValue);
    }

    if (cfgValue.isMember("ops"))
    {
        Json::Value ops_value = cfgValue["ops"];
        parseOpsConf(pDataFlowConf, ops_value);
    }

    // update dataflow config
    m_lockDataFlowConf.Lock();
    DataFlowConf *pTmp = m_dataflowConf;
    m_dataflowConf = pDataFlowConf;
    pDataFlowConf = pTmp;
    m_lockDataFlowConf.Unlock();

    delete pDataFlowConf;
    pDataFlowConf = NULL;

    return iRet;
}

int DataConf::parseReceiverConf(const Json::Value &cfgValue)
{
    int iRet = GSE_SUCCESS;

    clearReceiverConfMap();

    for (Json::ArrayIndex index = 0; index < cfgValue.size(); ++index)
    {
        string name = cfgValue[index].get("name", "").asString();
        if ("" == name)
        {
            LOG_WARN("the name of receiver in config is empty, please check config");
            continue;
        }

        if (isExitReceiver(name))
        {
            LOG_WARN("there is a same receiver[%s], please check config", SAFE_CSTR(name.c_str()));
            continue;
        }

        ReceiverConf *pReceiverConf = new ReceiverConf();
        pReceiverConf->m_name = name;
        pReceiverConf->m_protocol = ReceiverProtocolEnum(cfgValue[index].get("protocol", int(R_PROTO_UNKNOWN)).asInt());
        pReceiverConf->m_bind = cfgValue[index].get("bind", "").asString();
        pReceiverConf->m_port = cfgValue[index].get("port", 0).asUInt();
        pReceiverConf->m_certPath = cfgValue[index].get("cert", "").asString();
        pReceiverConf->m_protoStack = ProtocolStackEnum(cfgValue[index].get("protostack", int(PS_TYPE_UNKNOWN)).asInt());
        pReceiverConf->m_workThreadNum = cfgValue[index].get("workthreadnum", 4).asInt();
        //for udp receiver
        pReceiverConf->m_recvBufSize = cfgValue[index].get("recvbufsize", 33554432).asInt64();

        if (!insertReceiver(name, pReceiverConf))
        {
            delete pReceiverConf;
            pReceiverConf = NULL;
        }
    }

    return iRet;
}

int DataConf::parseFilterConf(const Json::Value &cfgValue)
{
    int iRet = GSE_SUCCESS;

    clearFilterConfMap();

    for (Json::ArrayIndex index = 0; index < cfgValue.size(); ++index)
    {
        string name = cfgValue[index].get("name", "").asString();
        if ("" == name)
        {
            LOG_WARN("the name of filter in config is empty, please check config");
            continue;
        }

        if (isExitFilter(name))
        {
            LOG_WARN("there is a same filter[%s], please check config", SAFE_CSTR(name.c_str()));
            continue;
        }

        FilterConf *pFilterConf = new FilterConf();
        pFilterConf->m_name = name;
        pFilterConf->m_filterType = FilterTypeEnum(cfgValue[index].get("type", 0).asInt());

        if (!insertFilter(name, pFilterConf))
        {
            delete pFilterConf;
            pFilterConf = NULL;
            LOG_WARN("fail to save filter[%s] into cache", name.c_str());
        }
    }

    return iRet;
}

int DataConf::parseExporterConf(const Json::Value &cfgValue)
{
    int iRet = GSE_SUCCESS;

    clearExporterConfMap();

    for (Json::ArrayIndex index = 0; index < cfgValue.size(); ++index)
    {
        string name = cfgValue[index].get("name", "").asString();
        if ("" == name)
        {
            LOG_WARN("the name of exporter in config is empty, please check config");
            continue;
        }

        if (isExitExporter(name))
        {
            LOG_WARN("there is a same exporter[%s], please check config", SAFE_CSTR(name.c_str()));
            continue;
        }

        ExporterConf *pExporterConf = new ExporterConf();
        pExporterConf->m_name = name;
        if (cfgValue[index].isMember("extensions") && cfgValue[index]["extensions"].isArray())
        {
            Json::Value extensions = cfgValue[index]["extensions"];
            for (Json::ArrayIndex idx = 0; idx < extensions.size(); ++idx)
            {
                pExporterConf->m_extensions.push_back(extensions[idx].asString());
            }
        }

        pExporterConf->m_type = ExporterTypeEnum(cfgValue[index].get("type", int(E_TYPE_UNKNOWN)).asInt());
        switch (pExporterConf->m_type)
        {
        case E_TYPE_LOG:
        {
            LogExporterConf *logExporterConf = parseLogExporterConf(name, cfgValue[index]);
            if (NULL == logExporterConf)
            {
                LOG_WARN("fail to parse log exporter config for exporter[%s]", SAFE_CSTR(name.c_str()));
                continue;
            }
            pExporterConf->m_logConf = logExporterConf;
        }
        break;
        case E_TYPE_BKDATA:
        {
            BkDataExporterConf *pBkDataExporterConf = parseBkDataExporterConf(name, cfgValue[index]);
            if (NULL == pBkDataExporterConf)
            {
                LOG_WARN("fail to parse bkdata exporter config for exporter[%s]", SAFE_CSTR(name.c_str()));
                continue;
            }
            pExporterConf->m_bkdataConf = pBkDataExporterConf;
        }
        break;
        case E_TYPE_KAFKA:
        {
            KafkaExporterConf *pKafkaExporterConf = parseKafkaExporterConf(name, cfgValue[index]);
            if (NULL == pKafkaExporterConf)
            {
                LOG_WARN("fail to parse kafka exporter config for exporter[%s]", SAFE_CSTR(name.c_str()));
                continue;
            }
            pExporterConf->m_kafkaConf = pKafkaExporterConf;
        }
        break;
        case E_TYPE_REDIS:
        {
            RedisExporterConf *pRedisExporterConf = parseRedisExporterConf(name, cfgValue[index]);
            if (NULL == pRedisExporterConf)
            {
                LOG_WARN("fail to parse redis exporter config for exporter[%s]", SAFE_CSTR(name.c_str()));
                continue;
            }
            pExporterConf->m_redisConf = pRedisExporterConf;
        }
        break;
        case E_TYPE_DS_PROXY:
        {
            DSProxyExporterConf *ptr_ds_proxy_exporter_conf = parseDSProxyExporterConf(name, cfgValue[index]);
            if (NULL == ptr_ds_proxy_exporter_conf)
            {
                LOG_WARN("failed to parse the ds proxy exporter config for exporter[%s]", SAFE_CSTR(name.c_str()));
                continue;
            }
            pExporterConf->m_dsProxyConf = ptr_ds_proxy_exporter_conf;
            break;
        }
        case E_TYPE_PULSAR:
            {
                PulsarExporterConf* ptr_pulsar_exporter_conf = parsePulsarExporterConf(name, cfgValue[index]);
                if(NULL == ptr_pulsar_exporter_conf)
                {
                    LOG_WARN("failed to parse the pulsar exporter config for exporter[%s]", SAFE_CSTR(name.c_str()));
                    continue;
                }
                
                pExporterConf->m_pulsarConf = ptr_pulsar_exporter_conf; 
            }
        break;
        default:
            LOG_WARN("unknown the exporter config type");
            delete pExporterConf;
            pExporterConf = NULL;
            break;
        }

        if (pExporterConf != NULL)
        {
            if (!insertExporter(name, pExporterConf))
            {
                delete pExporterConf;
                pExporterConf = NULL;
            }
        }
    }

    return iRet;
}

LogExporterConf *DataConf::parseLogExporterConf(const std::string &name, const Json::Value &cfgValue)
{
    LogExporterConf *logExporterConf = new LogExporterConf();
    logExporterConf->m_logPath = cfgValue.get("path", "").asString();
    logExporterConf->m_fileName = cfgValue.get("filename", "").asString();
    logExporterConf->m_maxLogSize = cfgValue.get("maxlogsize", 200).asInt();
    logExporterConf->m_maxFileNum = cfgValue.get("maxfilenum", 100).asInt();

    // 以单条数据1K计算，最高约 缓存 18GB 的存储空间
    logExporterConf->m_maxQueueLimitNum = cfgValue.get("queue_limit_max", 20000000).asInt();
    if (cfgValue.isMember("source_ids") && cfgValue["source_ids"].isArray())
    {
        for (Json::ArrayIndex idx = 0; idx < cfgValue["source_ids"].size(); ++idx)
        {
            logExporterConf->m_sourcedIDS.push_back(cfgValue["source_ids"][idx].asString());
        }
    }
    if (cfgValue.isMember("filters") && cfgValue["filters"].isArray())
    {
        for (Json::ArrayIndex j = 0; j < cfgValue["filters"].size(); j++)
        {
            string filterName = cfgValue["filters"][j].asString();
            logExporterConf->m_vFilterName.push_back(filterName);
            FilterConf *pConf = findFilterConf(filterName);
            if (pConf != NULL)
            {
                FilterConf *pFilterConf = new FilterConf(*pConf);
                logExporterConf->m_vFilterConf.push_back(pFilterConf);
            }
            else
            {
                LOG_WARN("don't find the filter[%s] for exporter[%s], please check config of dataflow", filterName.c_str(), name.c_str());
            }
        }
    }

    return logExporterConf;
}

KafkaExporterConf *DataConf::parseKafkaExporterConf(const std::string &name, const Json::Value &cfgValue)
{
    KafkaExporterConf* pKafkaExporterConf = new KafkaExporterConf();
    pKafkaExporterConf->m_cluster = cfgValue.get("cluster", "").asString();
    pKafkaExporterConf->m_producerNum = cfgValue.get("producernum", KAFKA_MAX_PRODUCER).asInt();
    pKafkaExporterConf->m_defaultTopicName = cfgValue.get("defaulttopicname", "").asString();

    pKafkaExporterConf->m_kafkaConfig.m_queueBufferingMaxMessages =  cfgValue.get("queue_buffering_max_messages", DEFAULT_MAX_KAFKA_QUEUE_SIZE).asString();
    pKafkaExporterConf->m_kafkaConfig.m_messageMaxBytes = cfgValue.get("message_max_bytes", DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE).asString();
    pKafkaExporterConf->m_kafkaConfig.m_securityProtocol = cfgValue.get("security_protocol", "").asString();
    pKafkaExporterConf->m_kafkaConfig.m_saslMechanisms = cfgValue.get("sasl_mechanisms", "").asString();
    pKafkaExporterConf->m_kafkaConfig.m_saslUserName = cfgValue.get("sasl_username", "").asString();
    pKafkaExporterConf->m_kafkaConfig.m_saslPasswd = cfgValue.get("sasl_passwd", "").asString();
    pKafkaExporterConf->m_kafkaConfig.m_requestRequiredAcks =  cfgValue.get("request_required_acks", "1").asString();
    pKafkaExporterConf->m_kafkaConfig.m_queueBufferingMaxMs =  cfgValue.get("queue_buffering_max_ms", "200").asString();

    if (cfgValue.isMember("filters") && cfgValue["filters"].isArray())
    {
        for (Json::ArrayIndex j = 0; j < cfgValue["filters"].size(); j++)
        {
            string filterName = cfgValue["filters"][j].asString();
            pKafkaExporterConf->m_vFilterName.push_back(filterName);
            FilterConf *pConf = findFilterConf(filterName);
            if (pConf != NULL)
            {
                FilterConf *pFilterConf = new FilterConf(*pConf);
                pKafkaExporterConf->m_vFilterConf.push_back(pFilterConf);
            }
            else
            {
                LOG_WARN("don't find the filter[%s] for exporter[%s], please check config of dataflow", filterName.c_str(), name.c_str());
            }
        }
    }

    return pKafkaExporterConf;
}

RedisExporterConf *DataConf::parseRedisExporterConf(const std::string &name, const Json::Value &cfgValue)
{
    RedisExporterConf *pRedisExporterConf = new RedisExporterConf();
    pRedisExporterConf->m_cluster = cfgValue.get("cluster", "").asString();
    if (cfgValue.isMember("filters") && cfgValue["filters"].isArray())
    {
        for (Json::ArrayIndex j = 0; j < cfgValue["filters"].size(); j++)
        {
            string filterName = cfgValue["filters"][j].asString();
            pRedisExporterConf->m_vFilterName.push_back(filterName);
            FilterConf *pConf = findFilterConf(filterName);
            if (pConf != NULL)
            {
                FilterConf *pFilterConf = new FilterConf(*pConf);
                pRedisExporterConf->m_vFilterConf.push_back(pFilterConf);
            }
            else
            {
                LOG_WARN("don't find the filter[%s] for exporter[%s], please check config of dataflow", filterName.c_str(), name.c_str());
            }
        }
    }

    return pRedisExporterConf;
}

BkDataExporterConf *DataConf::parseBkDataExporterConf(const std::string &name, const Json::Value &cfgValue)
{
    BkDataExporterConf *pBkDataExporterConf = new BkDataExporterConf();
    pBkDataExporterConf->m_zkAddrs = cfgValue.get("zkaddrs", "").asString();
    if (cfgValue.isMember("filters") && cfgValue["filters"].isArray())
    {
        for (Json::ArrayIndex j = 0; j < cfgValue["filters"].size(); j++)
        {
            string filterName = cfgValue["filters"][j].asString();
            pBkDataExporterConf->m_vFilterName.push_back(filterName);
            FilterConf *pConf = findFilterConf(filterName);
            if (pConf != NULL)
            {
                FilterConf *pFilterConf = new FilterConf(*pConf);
                pBkDataExporterConf->m_vFilterConf.push_back(pFilterConf);
            }
            else
            {
                LOG_WARN("don't find the filter[%s] for exporter[%s], please check config of dataflow", filterName.c_str(), name.c_str());
            }
        }
    }

    return pBkDataExporterConf;
}

DSProxyExporterConf *DataConf::parseDSProxyExporterConf(const std::string &name, const Json::Value &cfgValue)
{
    DSProxyExporterConf *ptr_ds_proxy_exporter_conf = new DSProxyExporterConf();

    if (cfgValue.isMember("addresses") && cfgValue["addresses"].isArray())
    {
        for (Json::ArrayIndex j = 0; j < cfgValue["addresses"].size(); j++)
        {
            Address address;
            address.m_ip = cfgValue["addresses"][j].get("ip", "").asString();
            address.m_port = cfgValue["addresses"][j].get("port", 0).asUInt();

            ptr_ds_proxy_exporter_conf->m_addresses.push_back(address);
        }
    }

    ptr_ds_proxy_exporter_conf->m_proxyProtocol = cfgValue.get("proxyprotocol", PROXY_PROTOCOL_TCP).asString();
    ptr_ds_proxy_exporter_conf->m_proxyVersion = cfgValue.get("proxyversion", PROXY_VERSION_1).asString();
    ptr_ds_proxy_exporter_conf->m_certPath = cfgValue.get("cert", "").asString();
    ptr_ds_proxy_exporter_conf->m_connectionNumEachAddress = cfgValue.get("connectionnum", 3).asUInt();
    ptr_ds_proxy_exporter_conf->m_heartbeat = cfgValue.get("heartbeat", false).asBool();
    ptr_ds_proxy_exporter_conf->m_noblock = cfgValue.get("noblock", true).asBool();
    ptr_ds_proxy_exporter_conf->m_fillChannelid = cfgValue.get("fillChannelid", true).asBool();

    return ptr_ds_proxy_exporter_conf;
}

/*
{
    "name":"pulsar_sz_cluster",
    "type":10,
    "service_url":"pulsar://localhost:6650",
    "producernum":8
},
*/

PulsarExporterConf* DataConf::parsePulsarExporterConf(const std::string& name, const Json::Value& cfgValue)
{
    PulsarExporterConf* ptr_pulsar_exporter_conf = new PulsarExporterConf();

    ptr_pulsar_exporter_conf->m_serviceUrl = cfgValue.get("service_url", "").asString();
    ptr_pulsar_exporter_conf->m_producerNum = cfgValue.get("producernum", 8).asInt();
    ptr_pulsar_exporter_conf->m_topicName = cfgValue.get("topicname", "").asString();
    ptr_pulsar_exporter_conf->m_tlsTrustCertsFilePath = cfgValue.get("tlstrustcertsfilepath", "").asString();
    ptr_pulsar_exporter_conf->m_tlsCertFilePath = cfgValue.get("tlscertfilepath", "").asString();
    ptr_pulsar_exporter_conf->m_tlsKeyFilePath = cfgValue.get("tlskeyfilepath", "").asString();
    ptr_pulsar_exporter_conf->m_token = cfgValue.get("token", "").asString();

    
    LOG_DEBUG("pulsar export config, url:%s, producernum:%d, topicname:%s, token:%s",   
        ptr_pulsar_exporter_conf->m_serviceUrl.c_str(),  ptr_pulsar_exporter_conf->m_producerNum, 
        ptr_pulsar_exporter_conf->m_topicName.c_str(), ptr_pulsar_exporter_conf->m_token.c_str());
    return ptr_pulsar_exporter_conf;
}


int DataConf::parseChannelConf(DataFlowConf *pDataFlowConf, const Json::Value &cfgValue)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataFlowConf)
    {
        LOG_WARN("the pointer of dataflow conf is NULL when parse receiver conf");
        iRet = GSE_INVALIDARGS;
        return iRet;
    }

    for (Json::ArrayIndex index = 0; index < cfgValue.size(); index++)
    {
        string name = cfgValue[index].get("name", "").asString();
        if ("" == name)
        {
            LOG_WARN("the name of channel in config is empty, please check config");
            continue;
        }

        if (pDataFlowConf->IsExitChannel(name))
        {
            LOG_WARN("there is a same channel[%s], please check config", SAFE_CSTR(name.c_str()));
            continue;
        }

        ChannelConf* pChannelConf = new ChannelConf();
        pChannelConf->m_name = name;
        pChannelConf->m_decodeType = DecodeTypeEnum(cfgValue[index].get("decode", int(D_TYPE_UNKNOWN)).asInt());
        pChannelConf->m_isKeepTiming = cfgValue[index].get("keeptiming", false).asBool();
        pChannelConf->m_numChannelWorker = cfgValue[index].get("workernum", 8).asUInt();

        pChannelConf->m_receiverName = cfgValue[index].get("receiver", "").asString();

        ReceiverConf *pReceiverConf = findReceiverConf(pChannelConf->m_receiverName);
        if (pReceiverConf != NULL)
        {
            pChannelConf->m_receiverConf = new ReceiverConf(*pReceiverConf);
        }

        if (cfgValue[index].isMember("exporter") && cfgValue[index]["exporter"].isArray())
        {
            for (Json::ArrayIndex j = 0; j < cfgValue[index]["exporter"].size(); j++)
            {
                string exporterName = cfgValue[index]["exporter"][j].asString();
                pChannelConf->m_vExporterName.push_back(exporterName);
                ExporterConf *pConf = findExporterConf(exporterName);
                if (pConf != NULL)
                {
                    ExporterConf *pExporterConf = new ExporterConf(*pConf);
                    pChannelConf->m_vExporterConf.push_back(pExporterConf);
                }
            }
        }

        pDataFlowConf->InsertChannel(name, pChannelConf);
    }

    return iRet;
}

int DataConf::parseOpsConf(DataFlowConf* pDataFlowConf, const Json::Value& cfgValue)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataFlowConf)
    {
        LOG_WARN("the pointer of dataflow conf is NULL when parse ops conf");
        iRet = GSE_INVALIDARGS;
        return iRet;
    }

    OpsConf* ptr_ops_conf = &pDataFlowConf->m_ops;
    ptr_ops_conf->m_agentOpsChannelID = cfgValue.get("agentopschannelid", 0).asInt();
    ptr_ops_conf->m_opsChannelID = cfgValue.get("opschannelid", 0).asInt();
    ptr_ops_conf->m_opsThreadCount = cfgValue.get("ops_threadcount", 16).asInt();
    if (cfgValue.isMember("ops_addresses") && cfgValue["ops_addresses"].isArray())
    {
        for (Json::ArrayIndex j = 0; j < cfgValue["ops_addresses"].size(); j++)
        {
            OpsAddress address;
            address.m_ip = cfgValue["ops_addresses"][j].get("ip","").asString();
            address.m_port = cfgValue["ops_addresses"][j].get("port",0).asUInt();

            ptr_ops_conf->m_opsAddress.push_back(address);
        }
    }

    return iRet;
}

void DataConf::clearReceiverConfMap()
{
    map<string, ReceiverConf *>::iterator receiverIter = m_receiversConf.begin();
    for (; receiverIter != m_receiversConf.end(); receiverIter++)
    {
        ReceiverConf *pReceiverConf = receiverIter->second;
        if (pReceiverConf != NULL)
        {
            delete pReceiverConf;
        }
    }
    m_receiversConf.clear();
}

void DataConf::clearExporterConfMap()
{
    map<string, ExporterConf *>::iterator exporterIter = m_exportersConf.begin();
    for (; exporterIter != m_exportersConf.end(); exporterIter++)
    {
        ExporterConf *pExporterConf = exporterIter->second;
        if (pExporterConf != NULL)
        {
            delete pExporterConf;
        }
    }
    m_exportersConf.clear();
}

void DataConf::clearFilterConfMap()
{
    map<string, FilterConf *>::iterator filterIter = m_filtersConf.begin();
    for (; filterIter != m_filtersConf.end(); filterIter++)
    {
        FilterConf *pFilterConf = filterIter->second;
        if (pFilterConf != NULL)
        {
            delete pFilterConf;
        }
    }
    m_filtersConf.clear();
}

bool DataConf::isExitReceiver(const std::string &receiverName)
{
    bool bRet = true;
    map<string, ReceiverConf *>::iterator iter;
    iter = m_receiversConf.find(receiverName);
    if (iter == m_receiversConf.end())
    {
        bRet = false;
    }

    return bRet;
}

bool DataConf::isExitExporter(const std::string &exporterName)
{
    bool bRet = true;
    map<string, ExporterConf *>::iterator iter;
    iter = m_exportersConf.find(exporterName);
    if (iter == m_exportersConf.end())
    {
        bRet = false;
    }

    return bRet;
}

bool DataConf::isExitFilter(const std::string filterName)
{
    bool bRet = true;
    map<string, FilterConf *>::iterator iter;
    iter = m_filtersConf.find(filterName);
    if (iter == m_filtersConf.end())
    {
        bRet = false;
    }

    return bRet;
}

bool DataConf::insertReceiver(const std::string &receiverName, ReceiverConf *pReceiverConf)
{
    bool bRet = false;
    if (isExitReceiver(receiverName))
    {
        LOG_WARN("fail to insert receiver config into cache, because there is a same receiver named [%s]", SAFE_CSTR(receiverName.c_str()));
        return bRet;
    }

    m_receiversConf.insert(map<string, ReceiverConf *>::value_type(receiverName, pReceiverConf));

    bRet = true;
    return bRet;
}

bool DataConf::insertExporter(const std::string &exporterName, ExporterConf *pExporterConf)
{
    bool bRet = false;
    if (isExitExporter(exporterName))
    {
        LOG_WARN("fail to insert exporter config into cache, because there is a same exporter named [%s]", SAFE_CSTR(exporterName.c_str()));
        return bRet;
    }

    m_exportersConf.insert(map<string, ExporterConf *>::value_type(exporterName, pExporterConf));

    bRet = true;
    return bRet;
}

bool DataConf::insertFilter(const std::string &filterName, FilterConf *pFilterConf)
{
    bool bRet = false;
    if (isExitFilter(filterName))
    {
        LOG_WARN("fail to insert filter config into cache, because there is a same filter named [%s]", SAFE_CSTR(filterName.c_str()));
        return bRet;
    }

    m_filtersConf.insert(map<string, FilterConf *>::value_type(filterName, pFilterConf));
    LOG_INFO("success save filter[%s] into cache", filterName.c_str());

    bRet = true;
    return bRet;
}

FilterConf *DataConf::findFilterConf(std::string &filterName)
{
    FilterConf *pFilterConf = NULL;
    map<string, FilterConf *>::iterator iter;
    iter = m_filtersConf.find(filterName);
    if (iter == m_filtersConf.end())
    {
        return pFilterConf;
    }

    pFilterConf = iter->second;

    return pFilterConf;
}

ExporterConf *DataConf::findExporterConf(std::string &exporterName)
{
    ExporterConf *pExporterConf = NULL;
    map<string, ExporterConf *>::iterator iter;
    iter = m_exportersConf.find(exporterName);
    if (iter == m_exportersConf.end())
    {
        return pExporterConf;
    }

    pExporterConf = iter->second;

    return pExporterConf;
}

ReceiverConf *DataConf::findReceiverConf(std::string &receiverName)
{
    ReceiverConf *pReceiverConf = NULL;
    map<string, ReceiverConf *>::iterator iter;
    iter = m_receiversConf.find(receiverName);
    if (iter == m_receiversConf.end())
    {
        return pReceiverConf;
    }

    pReceiverConf = iter->second;

    return pReceiverConf;
}
}
}
