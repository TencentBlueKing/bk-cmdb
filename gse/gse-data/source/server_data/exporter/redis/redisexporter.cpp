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

#include "bbx/gse_errno.h"
#include "db/redisapi/errors.hpp"
#include "log/log.h"
#include "ops/op_metric.h"
#include "tools/finally.hpp"
#include "tools/macros.h"
#include "tools/strings.h"

#include "dataserver.h"

namespace gse {
namespace data {

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
    StreamToCluster *streamToCluster = &ptrChannelIDConfig->m_streamToCluster;
    std::string host = streamToCluster->m_redisCluster->m_addresses[0].m_ip;
    uint16_t port = streamToCluster->m_redisCluster->m_addresses[0].m_port;
    std::string passwd = streamToCluster->m_redisCluster->m_passwd;
    std::string masterName = streamToCluster->m_redisCluster->m_masterName;
    std::string sentinelPasswd = streamToCluster->m_redisCluster->m_sentinelPasswd;

    switch (streamToCluster->m_clusterType)
    {
    case REDIS_SENTINEL_PUB: {

        RedisSentinelPublisher *producer = new RedisSentinelPublisher(host, port, passwd);
        LOG_INFO("create a new redis sentinel publish producer, the address is  (%s:%d), master name(%s)", SAFE_CSTR(host.c_str()), port, masterName.c_str());
        producer->SetMasterName(masterName);
        producer->SetSentinelPasswd(sentinelPasswd);
        if (producer->Init() == GSE_SUCCESS)
        {
            m_ptrSentinelPubliser = producer;
        }
        else
        {
            LOG_ERROR("it is failed to create a new redis sentinel publish producer, the address is  (%s:%d)", SAFE_CSTR(host.c_str()), port);
            delete producer;
            return false;
        }
    }
    break;
    case REDIS_PUB: {
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
        LOG_ERROR("it is failed to create a redis producer, because the storage  type (%d) is invalid", streamToCluster->m_clusterType);
        return false;
    }
    return true; // success
}

bool RedisExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    return false;
}

int RedisExporter::Start()
{
    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
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
    if (NULL != m_ptrSentinelPubliser)
    {
        delete m_ptrSentinelPubliser;
        m_ptrSentinelPubliser = NULL;
    }
    if (NULL != m_ptrPubliser)
    {
        delete m_ptrPubliser;
        m_ptrPubliser = NULL;
    }
    return GSE_SUCCESS;
}

int RedisExporter::Write(DataCell *pDataCell)
{

    if (NULL == pDataCell)
    {
        LOG_WARN("the parameter is valid. the pointer of data cell is [%x]", pDataCell);
        return GSE_ERROR;
    }

    pDataCell->DealLineBreak();
    // m_ptrOPSReport->PutOpsData(pDataCell->ToOPS(kOutputState + "_redis"));

    std::string value(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());

    std::vector<std::string> str_topics;
    pDataCell->GetTableName(str_topics);

    if (str_topics.empty())
    {
        LOG_INFO("not found the topic for the channel id (%u)", pDataCell->GetChannelID());
        return GSE_ERROR;
    }

    std::string master_host;
    int master_port;

    int errcode = 0;
    uint64_t beginTimestamp = gse::tools::time::GetUTCMillsec();
    auto _ = gse::tools::defer::finally([&]() {
        uint64_t endTimestamp = gse::tools::time::GetUTCMillsec();
        uint64_t costTime = (endTimestamp >= beginTimestamp) ? (endTimestamp >= beginTimestamp) : 0;
        OPMetric::AddMessageQueneRequestMillsecondsMetrics("redis", costTime);
        OPMetric::AddSendMsgBytesCounter(kOutputRedis, kOutputRedis, value.size());
        OPMetric::AddSendMsgCounter(kOutputRedis, kOutputRedis, errcode, 1);
    });

    std::size_t max_count = str_topics.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string str_topic = str_topics.at(idx);
        LOG_DEBUG("will send the data (%s) to redis with the channelname (%s), the channelid (%d) exporter name (%s)", SAFE_CSTR(value.c_str()), SAFE_CSTR(str_topic.c_str()), pDataCell->GetChannelID(), SAFE_CSTR(m_name.c_str()));
        if (NULL != m_ptrSentinelPubliser)
        {
            m_ptrSentinelPubliser->GetHost(master_host, master_port);
            gse::redis::RedisErrorCode ret = m_ptrSentinelPubliser->Produce(str_topic, value);

            if (ret == gse::redis::RedisErrorCode::E_OK)
            {
                std::string redis_host = master_host + "|" + gse::tools::strings::ToString(master_port);
                pDataCell->SetOutputType("redis_sentinel");
                pDataCell->SetOutputAddress(redis_host);
                DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
            }
            else
            {
                pDataCell->SetErrorMsg("failed to send to redis", OPS_ERROR_REDIS_ERROR);
                DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
                LOG_ERROR("failed to send msg to redis, error:%d", ret);
                errcode = (int)ret;
            }
        }

        if (NULL != m_ptrPubliser)
        {
            gse::redis::RedisErrorCode ret = m_ptrPubliser->produce(str_topic, value);

            if (ret == gse::redis::RedisErrorCode::E_OK)
            {
                std::string redis_host;
                int port = 0;
                std::string address = redis_host + "|" + gse::tools::strings::ToString(port);
                m_ptrPubliser->GetHost(redis_host, port);
                pDataCell->SetOutputType("redis_sentinel");
                pDataCell->SetOutputAddress(address);
                DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
            }
            else
            {
                pDataCell->SetErrorMsg("failed to end to redis", OPS_ERROR_REDIS_ERROR);
                DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
                LOG_ERROR("failed to send msg to redis, error:%d", ret);
                errcode = (int)ret;
            }
        }
    }

    return GSE_SUCCESS;
}

} // namespace data
} // namespace gse
