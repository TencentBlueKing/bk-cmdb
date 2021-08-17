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

#include "redisexporter.h"
#include "log/log.h"
#include "tools/macros.h"
#include "bbx/gse_errno.h"
#include "tools/strings.h"

#include "dataserver.h"

namespace gse { 
namespace dataserver {

RedisExporter::RedisExporter()
{
    m_ptrSentinelPubliser = NULL;
    m_ptrPubliser = NULL;
}

RedisExporter::~RedisExporter()
{
    //
}

bool RedisExporter::startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig)
{

    switch (ptrChannelIDConfig->m_storageType)
    {
    case REDIS_SENTINEL_PUB:
    {
        std::string sentinal_host = ptrChannelIDConfig->m_storage.m_redisStorage->m_addresses[0].m_ip;
        uint16_t sentinal_port = ptrChannelIDConfig->m_storage.m_redisStorage->m_addresses[0].m_port;
        std::string passwd = ptrChannelIDConfig->m_storage.m_redisStorage->m_passwd;
        std::string master_name = ptrChannelIDConfig->m_storage.m_redisStorage->m_masterName;
        RedisSentinelPublisher *producer = new RedisSentinelPublisher(sentinal_host, sentinal_port, passwd);
         LOG_INFO("create a new redis sentinel publish producer, the address is  (%s:%d), master name(%s)", SAFE_CSTR(sentinal_host.c_str()), sentinal_port, master_name.c_str());
        producer->setMasterName(master_name);
        if (producer->init() == GSE_SUCCESS)
        {
            m_ptrSentinelPubliser = producer;
        }
        else
        {
            LOG_ERROR("it is failed to create a new redis sentinel publish producer, the address is  (%s:%d)", SAFE_CSTR(sentinal_host.c_str()), sentinal_port);
            delete producer;
            return false;
        }
    }
    break;
    case REDIS_PUB:
    {

        std::string host = ptrChannelIDConfig->m_storage.m_redisStorage->m_addresses[0].m_ip;
        uint16_t port = ptrChannelIDConfig->m_storage.m_redisStorage->m_addresses[0].m_port;
        std::string passwd = ptrChannelIDConfig->m_storage.m_redisStorage->m_passwd;
        RedisPublishProducer *producer = new RedisPublishProducer(host, port, passwd);
        LOG_INFO("create a new redis publish producer, the address is  (%s:%d)", SAFE_CSTR(host.c_str()), port);
        if (producer->init() == GSE_SUCCESS)
        {
            m_ptrPubliser = producer;
        }
        else
        {
            LOG_ERROR("it is failed to create a new redis publish producer, the address is  (%s:%d)", SAFE_CSTR(host.c_str()), port);
            delete producer;
            return false;
        }
        return true;
    }
    break;
    default:
        LOG_ERROR("it is failed to create a redis producer, because the storage  type (%d) is invalid", ptrChannelIDConfig->m_storageType);
        return false;
    }
    return true; // success
}

bool RedisExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    return false;
}

bool RedisExporter::startWithDataID(StorageConfigType *ptrStorageConfig)
{
    switch (ptrStorageConfig->m_storageType)
    {
    case REDIS_SENTINEL_PUB:
    {
        RedisSentinelPublisher *producer = new RedisSentinelPublisher(ptrStorageConfig->m_host, ptrStorageConfig->m_port, ptrStorageConfig->m_passwd);
         LOG_INFO("create a new redis sentinel publish producer, the address is  (%s:%d)", SAFE_CSTR(ptrStorageConfig->m_host.c_str()), ptrStorageConfig->m_port);
        producer->setMasterName(ptrStorageConfig->m_masterName);
        if (producer->init() == GSE_SUCCESS)
        {
            m_ptrSentinelPubliser = producer;
        }
        else
        {
            LOG_ERROR("it is failed to create a new redis sentinel publish producer, the address is  (%s:%d)", SAFE_CSTR(ptrStorageConfig->m_host.c_str()), ptrStorageConfig->m_port);
            delete producer;
            return false;
        }
    }
    break;
    case REDIS_PUB:
    {
        RedisPublishProducer *producer = new RedisPublishProducer(ptrStorageConfig->m_host, ptrStorageConfig->m_port, ptrStorageConfig->m_passwd);
        LOG_INFO("create a new redis publish producer, the address is  (%s:%d)", SAFE_CSTR(ptrStorageConfig->m_host.c_str()), ptrStorageConfig->m_port);
        if (producer->init() == GSE_SUCCESS)
        {
            m_ptrPubliser = producer;
        }
        else
        {
            LOG_ERROR("it is failed to create a new redis publish producer, the address is  (%s:%d)", SAFE_CSTR(ptrStorageConfig->m_host.c_str()), ptrStorageConfig->m_port);
            delete producer;
            return false;
        }
    }
    break;
    default:
        LOG_ERROR("it is failed to create a redis producer, because the storage  (%s:%d) type (%d) is invalid", SAFE_CSTR(ptrStorageConfig->m_host.c_str()), ptrStorageConfig->m_port, ptrStorageConfig->m_storageType);
        return false;
    }
    return true; // success
}

int RedisExporter::Start()
{
    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
    case ExporterConfigWrapper::DataIDConfType:
        m_upConfLock.RLock();
        if (!startWithDataID(m_ptrConfWrapper->m_conf.m_ptrDataIDConfig))
        {
            m_upConfLock.UnLock();
            LOG_ERROR("start redis exporter (%s) by dataid config failed", SAFE_CSTR(m_name.c_str()));
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_INFO("start redis exporter (%s) by dataid config", SAFE_CSTR(m_name.c_str()));
        return GSE_SUCCESS;
    case ExporterConfigWrapper::ChannelIDConfType:
        m_upConfLock.RLock();
        if (!startWithChannelID(m_ptrConfWrapper->m_conf.m_ptrChannelIdExporterConfig))
        {
            m_upConfLock.UnLock();
            LOG_ERROR("start redis exporter (%s) by channelid config failed", SAFE_CSTR(m_name.c_str()));
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_INFO("start redis exporter (%s) by channelid config", SAFE_CSTR(m_name.c_str()));
        return GSE_SUCCESS;
    case ExporterConfigWrapper::DataFlowConfType:
        m_upConfLock.RLock();
        if (!startWithDataFlow(m_ptrConfWrapper->m_conf.m_exporterConf))
        {
            m_upConfLock.UnLock();
            LOG_ERROR("start redis exporter (%s) by dataflow config failed", SAFE_CSTR(m_name.c_str()));
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_INFO("start redis exporter (%s) by data flow config", SAFE_CSTR(m_name.c_str()));
        return GSE_SUCCESS;
    default:
        LOG_ERROR("the config type (%d) used to redis exporter (%s) is invalid", m_ptrConfWrapper->m_exporterConfTypeEnum, SAFE_CSTR(m_name.c_str()));
        break;
    }

    return GSE_ERROR;
}

int RedisExporter::Stop()
{
    if(NULL != m_ptrSentinelPubliser)
    {
        delete m_ptrSentinelPubliser;
    }
    if(NULL != m_ptrPubliser)
    {
        delete m_ptrPubliser;
    }
    return GSE_SUCCESS;
}

int RedisExporter::Write(DataCell *pDataCell)
{

    if(NULL == pDataCell)
    {
        LOG_WARN("the parameter is valid. the pointer of data cell is [%x]", pDataCell);
        return GSE_ERROR;
    }

    pDataCell->DealLineBreak();
    //m_ptrOPSReport->PutOpsData(pDataCell->ToOPS(kOutputState + "_redis"));

    std::string value(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());

    std::vector<std::string> str_topics;
    pDataCell->GetTableName(str_topics);

    if(str_topics.empty())
    {
        LOG_INFO("not found the topic for the channel id (%u)", pDataCell->GetChannelID());
        return GSE_ERROR;
    }

    std::string master_host;
    int master_port;

    std::size_t max_count = str_topics.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string str_topic = str_topics.at(idx);
        LOG_DEBUG("will send the data (%s) to redis with the topic (%s), the channelid (%d) exporter name (%s)", SAFE_CSTR(value.c_str()), SAFE_CSTR(str_topic.c_str()), pDataCell->GetChannelID(),SAFE_CSTR(m_name.c_str()));
        if (NULL != m_ptrSentinelPubliser)
        {
            OPMetric::RedisMsgInc();
            m_ptrSentinelPubliser->getRedisMasterHostAndPort(master_host, master_port);
            int ret = m_ptrSentinelPubliser->produce(str_topic, value);

            if (ret == GSE_SUCCESS)
            {
                std::string redis_host = master_host + "|" + gse::tools::strings::ToString(master_port);
                pDataCell->SetOutputType("redis_sentinel");
                pDataCell->SetOutputAddress(redis_host);
                DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
            }
            else
            {
                pDataCell->SetErrorMsg("failed to send to redis", OPS_ERROR_REDIS_ERROR);
                DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            }
        }

        if (NULL != m_ptrPubliser)
        {
            OPMetric::RedisMsgInc();
            int ret = m_ptrPubliser->produce(str_topic, value);

            if (ret == GSE_SUCCESS)
            {
                std::string redis_host;
                int port = 0;
                std::string address = redis_host + "|" + gse::tools::strings::ToString(port);
                m_ptrPubliser->GetHost(redis_host, port);
                pDataCell->SetOutputType("redis_sentinel");
                pDataCell->SetOutputAddress(address);
                DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
            }
            else
            {
                pDataCell->SetErrorMsg("failed to end to redis", OPS_ERROR_REDIS_ERROR);
                DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            }
        }
    }

    LOG_WARN("it is not set any redis writer for the data (%d)", pDataCell->GetChannelID());
    return GSE_SUCCESS;
}

}
}
