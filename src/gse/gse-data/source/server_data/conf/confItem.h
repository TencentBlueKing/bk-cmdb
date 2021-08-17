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
#ifndef _GSE_DATA_CONFITEM_H_
#define _GSE_DATA_CONFITEM_H_

#include <string>
#include <vector>
#include <map>

#include "tools/macros.h"
#include "dataStruct/safe_map.h"
#include "conf/bkdata_config_v1.h"
#include "conf/channel_id_config.h"
#include "log/log.h"
#include "conf_common.h"

namespace gse { 
namespace dataserver {

#define KAFKA_MAX_QUEUE 100000
#define KAFKA_MAX_PRODUCER 8

#define PULSAR_MAX_PRODUCER 8
#define DEFAULT_MAX_KAFKA_QUEUE_SIZE 10000000
#define DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE 10000000


#ifndef OPS_QUEUE_MAX
#define OPS_QUEUE_MAX (102400)
#endif

typedef int StorageIndex;
typedef uint32_t DataIDType;
typedef uint32_t ChannelIDType;
typedef uint32_t StreamToIDType;
typedef std::vector<StorageConfigType> StorageConfigVector;
typedef gse::datastruct::SafeMap<StorageIndex, StorageConfigVector *> StorageConfigMap;
typedef gse::datastruct::SafeMap<DataIDType, DataID *> DataIDConfigMap;

typedef enum ReceiverProtocolEnum_
{
    R_PROTO_UNKNOWN = 0,
    R_PROTO_TCP = 1,
    R_PROTO_UDP = 2,
    R_PROTO_KCP = 3,
    R_PROTO_HTTP = 4
} ReceiverProtocolEnum;

typedef enum ProtocolStackEnum_
{
    PS_TYPE_UNKNOWN = 0,
    PS_TYPE_GSEDATA = 1,                   // 2.0 通用的动态协议
    PS_TYPE_GSEDATA_V1 = 2,                // 兼容V1.0 版本的DS通信协议
    PS_TYPE_GSEDATA_V1_FOR_TGLOG_PROXY = 3, // 用于兼容 原 数据链路（tdbank->tglog->datamore->transit server)推送的数据
    PS_TYPE_GSEDATA_V2 = 4                  //gse2.0
} ProtocolStackEnum;

typedef enum ExporterTypeEnum_
{
    E_TYPE_UNKNOWN = 0,
    E_TYPE_LOG = 1,
    E_TYPE_TCP = 2,
    E_TYPE_UDP = 3,
    E_TYPE_KCP = 4,
    E_TYPE_HTTP = 5,
    E_TYPE_REDIS = 6,
    E_TYPE_KAFKA = 7,
    E_TYPE_BKDATA = 8,
    E_TYPE_DS_PROXY = 9,
    E_TYPE_PULSAR = 10
} ExporterTypeEnum;

typedef enum DecodeTypeEnum_
{
    D_TYPE_UNKNOWN = 0,
    D_TYPE_TGLOG_PROTOBUF = 1,
    D_TYPE_TDM_PROTOBUF = 2,
    D_TYPE_TDM_PACKAGE = 3,
    D_TYPE_GSEDATA_PACKAGE = 4,
    D_TYPE_GSEDATA_PACKAGE_V1 = 5,         // 兼容V1.0 版本的DS通信协议数据解析
    D_TYPE_GSEDATA_V1_FOR_TGLOG_PROXY = 6, // 用于兼容 原 数据链路（tdbank->tglog->datamore->transit server)推送的数据
    D_TYPE_ONLY_TRANSPORT = 7,              // 不做实质的解析，仅作数据透传
    D_TYPE_GSEDATA_PACKAGE_V2 = 8,
    D_TYPE_GSEDATA_PACKAGE_OPEN = 9
} DecodeTypeEnum;

typedef enum FilterTypeEnum_
{
    FILTER_TYPE_UNKNOW = 0,
    FILTER_TYPE_CHANNELID = 1
} FilterTypeEnum;

typedef struct ReceiverConf_
{
    std::string m_name;
    ReceiverProtocolEnum m_protocol;
    std::string m_bind;
    uint16_t m_port;
    int m_workThreadNum;
    std::string m_certPath;
    uint32_t m_recvBufSize;
    ProtocolStackEnum m_protoStack;

    ReceiverConf_()
    {
        m_protocol = R_PROTO_UNKNOWN;
        m_port = 0;
        m_workThreadNum = 16;
        m_recvBufSize = 33554432;
        m_protoStack = PS_TYPE_UNKNOWN;
    }

    ReceiverConf_ &operator=(const ReceiverConf_ &srcConf)
    {
        this->m_name = srcConf.m_name;
        this->m_protocol = srcConf.m_protocol;
        this->m_bind = srcConf.m_bind;
        this->m_port = srcConf.m_port;
        this->m_workThreadNum = srcConf.m_workThreadNum;
        this->m_certPath = srcConf.m_certPath;
        this->m_protoStack = srcConf.m_protoStack;
        this->m_recvBufSize = srcConf.m_recvBufSize;
        return *this;
    }

    ReceiverConf_(const ReceiverConf_ &srcConf)
    {
        *this = srcConf;
    }

} ReceiverConf;

typedef struct FilterConf_
{
    std::string m_name;
    FilterTypeEnum m_filterType;

    FilterConf_()
    {
        m_name = "";
        m_filterType = FILTER_TYPE_UNKNOW;
    }

    ~FilterConf_()
    {
        //
    }

    FilterConf_ &operator=(const FilterConf_ &srcConf)
    {
        this->m_name = srcConf.m_name;
        this->m_filterType = srcConf.m_filterType;
        return *this;
    }

    FilterConf_(const FilterConf_ &srcConf)
    {
        *this = srcConf;
    }
} FilterConf;

typedef struct LogExporterConf_
{
    std::string m_logPath;
    std::string m_fileName;
    int m_maxLogSize;
    int m_maxFileNum;
    int m_maxQueueLimitNum;
    std::vector<std::string> m_sourcedIDS;
    std::vector<std::string> m_vFilterName;
    std::vector<FilterConf *> m_vFilterConf;

    LogExporterConf_()
    {
        m_logPath = "";
        m_fileName = "";
        m_maxLogSize = 200;
        m_maxFileNum = 100;
        m_maxQueueLimitNum = 20000000;
    }

    ~LogExporterConf_()
    {
        m_vFilterName.clear();
        for (int i = 0; i < m_vFilterConf.size(); i++)
        {
            FilterConf *pFilterConf = m_vFilterConf[i];
            if (pFilterConf != NULL)
            {
                delete pFilterConf;
            }
        }
        m_vFilterConf.clear();
    }

    LogExporterConf_ &operator=(const LogExporterConf_ &srcConf)
    {
        this->m_logPath = srcConf.m_logPath;
        this->m_fileName = srcConf.m_fileName;
        this->m_maxLogSize = srcConf.m_maxLogSize;
        this->m_maxFileNum = srcConf.m_maxFileNum;
        this->m_maxQueueLimitNum = srcConf.m_maxQueueLimitNum;
        std::size_t max_count = srcConf.m_sourcedIDS.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            this->m_sourcedIDS.push_back(srcConf.m_sourcedIDS.at(idx));
        }
        this->m_vFilterName.clear();
        for (int i = 0; i < srcConf.m_vFilterName.size(); i++)
        {
            this->m_vFilterName.push_back(srcConf.m_vFilterName[i]);
        }

        for (int j = 0; j < this->m_vFilterConf.size(); j++)
        {
            FilterConf *pFilterConf = this->m_vFilterConf[j];
            if (pFilterConf != NULL)
            {
                delete pFilterConf;
            }
        }
        this->m_vFilterConf.clear();

        for (int k = 0; k < srcConf.m_vFilterConf.size(); k++)
        {
            FilterConf *pSrcConf = srcConf.m_vFilterConf[k];
            if (NULL == pSrcConf)
            {
                continue;
            }
            FilterConf *pFilterConf = new FilterConf(*pSrcConf);
            this->m_vFilterConf.push_back(pFilterConf);
        }

        return *this;
    }

    LogExporterConf_(const LogExporterConf_ &srcConf)
    {
        *this = srcConf;
    }

} LogExporterConf;

typedef struct KafkaExporterConf_
{
    std::string m_cluster;
    int32_t m_producerNum;
    std::string m_defaultTopicName;
    std::vector<std::string> m_vFilterName;
    std::vector<FilterConf *> m_vFilterConf;

    KafkaConfig m_kafkaConfig;

    KafkaExporterConf_()
    {
        m_cluster = "";
        m_defaultTopicName = "";
        m_producerNum = KAFKA_MAX_PRODUCER;
    }

    ~KafkaExporterConf_()
    {
        m_vFilterName.clear();
        for (int i = 0; i < m_vFilterConf.size(); i++)
        {
            FilterConf *pFilterConf = m_vFilterConf[i];
            if (pFilterConf != NULL)
            {
                delete pFilterConf;
            }
        }
        m_vFilterConf.clear();
    }

    KafkaExporterConf_ &operator=(const KafkaExporterConf_ &srcConf)
    {
        m_cluster = srcConf.m_cluster;
        m_producerNum = srcConf.m_producerNum;
        m_defaultTopicName = srcConf.m_defaultTopicName;
        m_kafkaConfig = srcConf.m_kafkaConfig;

        this->m_vFilterName.clear();
        for (int i = 0; i < srcConf.m_vFilterName.size(); i++)
        {
            this->m_vFilterName.push_back(srcConf.m_vFilterName[i]);
        }

        for (int j = 0; j < this->m_vFilterConf.size(); j++)
        {
            FilterConf *pFileConf = this->m_vFilterConf[j];
            if (pFileConf != NULL)
            {
                delete pFileConf;
            }
        }
        this->m_vFilterConf.clear();
        for (int k = 0; k < srcConf.m_vFilterConf.size(); k++)
        {
            FilterConf *pSrcConf = srcConf.m_vFilterConf[k];
            if (NULL == pSrcConf)
            {
                continue;
            }
            FilterConf *pFilterConf = new FilterConf(*pSrcConf);
            this->m_vFilterConf.push_back(pFilterConf);
        }
        return *this;
    }

    KafkaExporterConf_(const KafkaExporterConf_ &srcConf)
    {
        *this = srcConf;
    }

    void CopyFilterConfs(std::vector<FilterConf *> &vFilterConf)
    {
        for (int i = 0; i < m_vFilterConf.size(); ++i)
        {
            FilterConf *pTmp = m_vFilterConf[i];
            if (NULL == pTmp)
            {
                continue;
            }

            FilterConf *pFilterConf = new FilterConf(*pTmp);
            vFilterConf.push_back(pFilterConf);
        }
    }

} KafkaExporterConf;

typedef struct PulsarExporterConf_
{
    std::string m_serviceUrl;
    std::string m_topicName;
    std::string m_tlsTrustCertsFilePath;
    std::string m_tlsCertFilePath;
    std::string m_tlsKeyFilePath;
    std::string m_token;
    int32_t m_producerNum;

    PulsarExporterConf_()
    {
        m_producerNum = 8;
        m_serviceUrl = "";
        m_token = "";
        m_tlsKeyFilePath = "";
        m_tlsCertFilePath = "";
        m_tlsTrustCertsFilePath = "";
        m_topicName = "";
    }

    ~PulsarExporterConf_()
    {
    }

    PulsarExporterConf_ &operator=(const PulsarExporterConf_ &srcConf)
    {
        this->m_serviceUrl = srcConf.m_serviceUrl;
        this->m_producerNum = srcConf.m_producerNum;
        this->m_topicName = srcConf.m_topicName;
        this->m_tlsTrustCertsFilePath = srcConf.m_tlsTrustCertsFilePath;
        this->m_tlsCertFilePath = srcConf.m_tlsCertFilePath; 
        this->m_tlsKeyFilePath = srcConf.m_tlsKeyFilePath;
        this->m_token = srcConf.m_token;
        return *this;
    }

    PulsarExporterConf_(const PulsarExporterConf_ &srcConf)
    {
        *this = srcConf;
    }

} PulsarExporterConf;

typedef struct RedisExporterConf_
{
    std::string m_cluster;
    std::vector<std::string> m_vFilterName;
    std::vector<FilterConf *> m_vFilterConf;

    RedisExporterConf_()
    {
        m_cluster = "";
    }

    ~RedisExporterConf_()
    {
        m_vFilterName.clear();
        for (int i = 0; i < m_vFilterConf.size(); i++)
        {
            FilterConf *pFilterConf = m_vFilterConf[i];
            if (pFilterConf != NULL)
            {
                delete pFilterConf;
            }
        }
        m_vFilterConf.clear();
    }

    RedisExporterConf_ &operator=(const RedisExporterConf_ &srcConf)
    {
        this->m_cluster = srcConf.m_cluster;

        this->m_vFilterName.clear();
        for (int i = 0; i < srcConf.m_vFilterName.size(); i++)
        {
            this->m_vFilterName.push_back(srcConf.m_vFilterName[i]);
        }

        for (int j = 0; j < this->m_vFilterConf.size(); j++)
        {
            FilterConf *pFileConf = this->m_vFilterConf[j];
            if (pFileConf != NULL)
            {
                delete pFileConf;
            }
        }
        this->m_vFilterConf.clear();
        for (int k = 0; k < srcConf.m_vFilterConf.size(); k++)
        {
            FilterConf *pSrcConf = srcConf.m_vFilterConf[k];
            if (NULL == pSrcConf)
            {
                continue;
            }
            FilterConf *pFilterConf = new FilterConf(*pSrcConf);
            this->m_vFilterConf.push_back(pFilterConf);
        }

        return *this;
    }

    RedisExporterConf_(const RedisExporterConf_ &srcConf)
    {
        *this = srcConf;
    }
} RedisExporterConf;

typedef struct BkDataExporterConf_
{
    std::string m_zkAddrs;
    std::vector<std::string> m_vFilterName;
    std::vector<FilterConf *> m_vFilterConf;

    BkDataExporterConf_()
    {
        m_zkAddrs = "";
    }

    ~BkDataExporterConf_()
    {
        m_vFilterName.clear();
        for (int i = 0; i < m_vFilterConf.size(); i++)
        {
            FilterConf *pFilterConf = m_vFilterConf[i];
            if (pFilterConf != NULL)
            {
                delete pFilterConf;
            }
        }
        m_vFilterConf.clear();
    }

    BkDataExporterConf_ &operator=(const BkDataExporterConf_ &srcConf)
    {
        this->m_zkAddrs = srcConf.m_zkAddrs;
        this->m_vFilterName.clear();
        for (int i = 0; i < srcConf.m_vFilterName.size(); i++)
        {
            this->m_vFilterName.push_back(srcConf.m_vFilterName[i]);
        }

        for (int j = 0; j < m_vFilterConf.size(); j++)
        {
            FilterConf *pFilterConf = m_vFilterConf[j];
            if (pFilterConf != NULL)
            {
                delete pFilterConf;
            }
        }
        m_vFilterConf.clear();

        for (int k = 0; k < srcConf.m_vFilterConf.size(); k++)
        {
            FilterConf *pSrcConf = srcConf.m_vFilterConf[k];
            if (NULL == pSrcConf)
            {
                continue;
            }
            FilterConf *pFilterConf = new FilterConf(*pSrcConf);
            m_vFilterConf.push_back(pFilterConf);
        }
        return *this;
    }

    BkDataExporterConf_(const BkDataExporterConf_ &srcConf)
    {
        m_vFilterConf.clear();
        *this = srcConf;
    }

} BkDataExporterConf;

#ifndef PROXY_PROTOCOL_TCP
#define PROXY_PROTOCOL_TCP "tcp"
#endif

#ifndef PROXY_PROTOCOL_UDP
#define PROXY_PROTOCOL_UDP "udp"
#endif

#ifndef PROXY_VERSION_1
#define PROXY_VERSION_1 "v1"
#endif

#ifndef PROXY_VERSION_2
#define PROXY_VERSION_2 "v2"
#endif

typedef struct DSProxyExporterConf_
{
    std::vector<Address> m_addresses;
    std::string m_certPath;
    // proxy version v1 or v2
    std::string m_proxyVersion;
    // proxy protocol udp or tcp
    std::string m_proxyProtocol;
    int m_connectionNumEachAddress;
    bool m_heartbeat;
    bool m_noblock;
    bool m_fillChannelid;

    DSProxyExporterConf_ &operator=(const DSProxyExporterConf_ &srcConf)
    {
        std::size_t max_count = srcConf.m_addresses.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            m_addresses.push_back(srcConf.m_addresses.at(idx));
        }
        this->m_certPath = srcConf.m_certPath;
        this->m_connectionNumEachAddress = srcConf.m_connectionNumEachAddress;
        this->m_proxyVersion = srcConf.m_proxyVersion;
        this->m_proxyProtocol = srcConf.m_proxyProtocol;
        this->m_heartbeat = srcConf.m_heartbeat;
        this->m_noblock = srcConf.m_noblock;
        this->m_fillChannelid = srcConf.m_fillChannelid;
        return *this;
    }
    
    DSProxyExporterConf_()
    {
        m_connectionNumEachAddress = 5;
        m_heartbeat = false;
        m_noblock = true;
        m_fillChannelid = false;
    }
    DSProxyExporterConf_(const DSProxyExporterConf_ &srcConf)
    {
        m_addresses.clear();
        *this = srcConf;
    }

} DSProxyExporterConf;

typedef struct ExporterConf_
{
    std::string m_name;
    ExporterTypeEnum m_type;
    std::vector<std::string> m_extensions;
    LogExporterConf *m_logConf;
    BkDataExporterConf *m_bkdataConf;
    KafkaExporterConf *m_kafkaConf;
    RedisExporterConf *m_redisConf;
    DSProxyExporterConf *m_dsProxyConf;
    PulsarExporterConf *m_pulsarConf;

    ExporterConf_()
    {
        m_name = "";
        m_type = E_TYPE_UNKNOWN;
        m_logConf = NULL;
        m_bkdataConf = NULL;
        m_kafkaConf = NULL;
        m_redisConf = NULL;
        m_dsProxyConf = NULL;
        m_pulsarConf = NULL;
    }
    ~ExporterConf_()
    {
        if (m_logConf != NULL)
        {
            delete m_logConf;
            m_logConf = NULL;
        }

        if (m_bkdataConf != NULL)
        {
            delete m_bkdataConf;
            m_bkdataConf = NULL;
        }

        if (m_kafkaConf != NULL)
        {
            delete m_kafkaConf;
            m_kafkaConf = NULL;
        }

        if (m_redisConf != NULL)
        {
            delete m_redisConf;
            m_redisConf = NULL;
        }
        if (m_dsProxyConf != NULL)
        {
            delete m_dsProxyConf;
            m_dsProxyConf = NULL;
        }

        if (m_pulsarConf != NULL)
        {
            delete m_pulsarConf;
            m_pulsarConf = NULL;
        }
    }

    ExporterConf_ &operator=(const ExporterConf_ &srcConf)
    {
        this->m_name = srcConf.m_name;
        this->m_type = srcConf.m_type;
        this->m_extensions = srcConf.m_extensions;
        switch (this->m_type)
        {
        case E_TYPE_LOG:
        {
            if (srcConf.m_logConf != NULL)
            {
                this->m_logConf = new LogExporterConf(*srcConf.m_logConf);
            }
        }
        break;
        case E_TYPE_BKDATA:
        {
            if (srcConf.m_bkdataConf != NULL)
            {
                this->m_bkdataConf = new BkDataExporterConf(*srcConf.m_bkdataConf);
            }
        }
        break;
        case E_TYPE_KAFKA:
        {
            if (srcConf.m_kafkaConf != NULL)
            {
                this->m_kafkaConf = new KafkaExporterConf(*srcConf.m_kafkaConf);
            }
        }
        break;
        case E_TYPE_REDIS:
        {
            if (srcConf.m_redisConf != NULL)
            {
                this->m_redisConf = new RedisExporterConf(*srcConf.m_redisConf);
            }
        }
        break;
        case E_TYPE_DS_PROXY:
        {
            if (srcConf.m_dsProxyConf != NULL)
            {
                this->m_dsProxyConf = new DSProxyExporterConf(*srcConf.m_dsProxyConf);
            }
        }
        break;

        case E_TYPE_PULSAR:
        {
            if (srcConf.m_pulsarConf != NULL)
            {
                LOG_DEBUG("new pulsar export conf");
                this->m_pulsarConf = new PulsarExporterConf(*srcConf.m_pulsarConf);
            }           
        }
        break;
        default:
            break;
        }

        return *this;
    }

    ExporterConf_(const ExporterConf_ &srcConf)
    {
        m_type = E_TYPE_UNKNOWN;
        m_logConf = NULL;
        m_bkdataConf = NULL;
        m_kafkaConf = NULL;
        m_redisConf = NULL;
        m_dsProxyConf = NULL;
        m_pulsarConf = NULL;
        *this = srcConf;
    }

} ExporterConf;



typedef struct ChannelConf_
{
    std::string m_name;
    DecodeTypeEnum m_decodeType;
    uint16_t m_numChannelWorker;
    bool m_isKeepTiming;
    std::string m_receiverName;
    ReceiverConf *m_receiverConf;
    std::vector<std::string> m_vExporterName;
    std::vector<ExporterConf *> m_vExporterConf;

    ChannelConf_()
    {
        m_name = "";
        m_decodeType = D_TYPE_UNKNOWN;
        m_receiverConf = NULL;
        m_numChannelWorker = 0;
        m_isKeepTiming = false;
    }

    ~ChannelConf_()
    {
        if (m_receiverConf != NULL)
        {
            delete m_receiverConf;
            m_receiverConf = NULL;
        }

        for (int i = 0; i < m_vExporterConf.size(); i++)
        {
            ExporterConf *pExporterConf = m_vExporterConf[i];
            delete pExporterConf;
        }
        m_vExporterConf.clear();
    }

    ChannelConf_ &operator=(const ChannelConf_ &srcConf)
    {
        this->m_name = srcConf.m_name;
        this->m_decodeType = srcConf.m_decodeType;
        this->m_isKeepTiming = srcConf.m_isKeepTiming;
        this->m_numChannelWorker = srcConf.m_numChannelWorker;
        this->m_receiverName = srcConf.m_receiverName;
        if (this->m_receiverConf != NULL)
        {
            delete this->m_receiverConf;
        }
        if (srcConf.m_receiverConf != NULL)
        {
            this->m_receiverConf = new ReceiverConf(*srcConf.m_receiverConf);
        }
        this->m_vExporterName.clear();
        for (int i = 0; i < srcConf.m_vExporterName.size(); i++)
        {
            this->m_vExporterName.push_back(srcConf.m_vExporterName[i]);
        }

        for (int j = 0; j < m_vExporterConf.size(); j++)
        {
            ExporterConf *pExporterConf = m_vExporterConf[j];
            if (pExporterConf != NULL)
            {
                delete pExporterConf;
            }
        }
        m_vExporterConf.clear();

        for (int k = 0; k < srcConf.m_vExporterConf.size(); k++)
        {
            ExporterConf *pConf = srcConf.m_vExporterConf[k];
            if (NULL == pConf)
            {
                continue;
            }

            ExporterConf *pExporterConf = new ExporterConf(*pConf);
            m_vExporterConf.push_back(pExporterConf);
        }
        return *this;
    }

    ChannelConf_(const ChannelConf_ &srcConf)
    {
        m_decodeType = D_TYPE_UNKNOWN;
        m_receiverConf = NULL;
        m_vExporterConf.clear();
        *this = srcConf;
    }

} ChannelConf;

typedef struct OpsAddress_
{
    OpsAddress_(){
        m_port = 0;
    }
    std::string m_ip;
    uint16_t m_port;
}OpsAddress;

typedef struct OpsConf_
{
    OpsConf_()
    {
        m_maxQueneSize = OPS_QUEUE_MAX;
        m_opsThreadCount = 16;
        m_opsChannelID = 0;
        m_agentOpsChannelID = 0;
    }
    std::vector<OpsAddress> m_opsAddress;
    uint64_t m_opsChannelID;
    uint64_t m_agentOpsChannelID;
    uint64_t m_opsThreadCount;
    uint64_t m_maxQueneSize;

}OpsConf;

typedef struct DataFlowConf_
{
    std::map<std::string, ChannelConf *> m_channelsConf;
    OpsConf m_ops;

    DataFlowConf_()
    {
        //
    }

    void Clear()
    {
        map<string, ChannelConf *>::iterator channelIter;
        for (channelIter = m_channelsConf.begin(); channelIter != m_channelsConf.end(); channelIter++)
        {
            ChannelConf *pChannelConf = channelIter->second;
            if (pChannelConf != NULL)
            {
                delete pChannelConf;
            }
        }
        m_channelsConf.clear();
    }

    ~DataFlowConf_()
    {
        Clear();
    }

    void GetOpsConf(OpsConf& ops)
    {
        ops = m_ops;
    }

    bool IsExitChannel(const std::string &channelName)
    {
        bool bRet = true;
        map<string, ChannelConf *>::iterator iter;
        iter = m_channelsConf.find(channelName);
        if (iter == m_channelsConf.end())
        {
            bRet = false;
        }

        return bRet;
    }

    bool InsertChannel(const std::string &channelName, ChannelConf *pChannelConf)
    {
        bool bRet = false;
        if (IsExitChannel(channelName))
        {
            LOG_WARN("fail to insert channel config into cache, because there is a same channel named [%s]", SAFE_CSTR(channelName.c_str()));
            return bRet;
        }

        m_channelsConf.insert(map<string, ChannelConf *>::value_type(channelName, pChannelConf));

        bRet = true;
        return bRet;
    }

    DataFlowConf_ &operator=(DataFlowConf_ &srcConf)
    {
        Clear();
        // channelconf
        map<string, ChannelConf *>::iterator channelIter = srcConf.m_channelsConf.begin();
        for (; channelIter != srcConf.m_channelsConf.end(); channelIter++)
        {
            ChannelConf *pChannelConf = new ChannelConf(*channelIter->second);
            if (!InsertChannel(channelIter->first, pChannelConf))
            {
                delete pChannelConf;
            }
        }

        m_ops = srcConf.m_ops;
        return *this;
    }

    DataFlowConf_(DataFlowConf_ &srcConf)
    {
        m_channelsConf.clear();
        *this = srcConf;
    }

} DataFlowConf;

// DataID 和ChannelID 与DataStorage 之间的关系是 1:1
class DataStorage
{
public:
    DataStorage()
    {
        m_next = NULL;
        m_isDataID = false;
        m_storage.m_ptrChannelIDStorage = NULL;
        m_storage.m_ptrDataIDConfig = NULL;  // 不需要在此结构中释放内存
        m_storageType = UNKNOWN;
    }

    ~DataStorage()
    {
        if(!m_isDataID)
        {
            delete m_storage.m_ptrChannelIDStorage;
            m_storage.m_ptrChannelIDStorage = NULL;
        }

        if (NULL != m_next)
        {
            delete m_next;
        }
    }

public:
    inline void SetNext(DataStorage *ptrNext)
    {
        if (NULL != m_next)
        {
            m_next->SetNext(ptrNext);
            return;
        }

        m_next = ptrNext;
    }

public:
    DataStorage *m_next;
    StorageType m_storageType;
    bool m_isDataID;

    union {
        // DataID 兼容V1版本而保留， 可以是多个存储，表示此列表里的存储都需要被写入数据
        StorageConfigVector *m_ptrDataIDConfig;

        // channelid config，可以是多个，表示此链表下的存储都需要被写入数据
        ChannelIDStorage *m_ptrChannelIDStorage;

    } m_storage;
};


class IDToStorage
{
public:
    typedef void (*WatchEventFunc)(void* args, int storageIndex, uint32_t channelID);
public:
    IDToStorage(){};
    virtual ~IDToStorage(){};

public:
    virtual DataID *GetStorageByDataID(uint32_t dataID) = 0;
    virtual DataStorage *GetAllStorages() = 0;
    virtual DataStorage *GetStorageByIndex(int storageIndex) = 0;
    virtual DataStorage *GetStorageByChannelID(uint32_t channelID) = 0;
    virtual void WatchUpdateEvent(WatchEventFunc callback, void *args) = 0;
};

}
}
#endif
