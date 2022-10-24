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

#include "kafkaexporter.h"
#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/error.h"
#include "tools/macros.h"
#include "tools/net.h"

#include "dataserver.h"
#include "utils.h"

#include "filter/channelid_filter.h"
namespace gse {
namespace data {

KafkaExporter::KafkaExporter()
    : m_nextProducerId(0),
      m_kafkaMaxQueue(DEFAULT_MAX_KAFKA_QUEUE_SIZE),
      m_threadStoped(false),
      m_kafkaMaxMessageBytes(DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE),
      m_producerNum(DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE)
{
}

KafkaExporter::~KafkaExporter()
{
    clear();
}

int KafkaExporter::createKafkaProducers()
{
    KafkaProducer::m_runtimeDataDirector = m_configPtr->m_runtimePath;
    if (m_kafkaBrokers.empty())
    {
        LOG_ERROR("not found the kafka brokers address for the exporter (%s)", SAFE_CSTR(m_name.c_str()));
        return GSE_SUCCESS;
    }

    for (int i = 0; i < m_producerNum; i++)
    {
        KafkaProducer *pKafkaProducer = new KafkaProducer();
        if (NULL == pKafkaProducer)
        {
            LOG_WARN("fail to create kafka producer [%d] object for exporter[%s]", i, SAFE_CSTR(m_name.c_str()));
            continue;
        }
        pKafkaProducer->SetKafkaConfig(m_kafkaConfig);
        int iRet = pKafkaProducer->CreateProducer(m_kafkaBrokers);
        if (iRet != GSE_SUCCESS)
        {
            LOG_WARN("fail to start kafka producer[%d] exporter[%s] with kafka brokers[%s]", i, SAFE_CSTR(m_name.c_str()), SAFE_CSTR(m_kafkaBrokers.c_str()));
            delete pKafkaProducer;
            return GSE_ERROR;
        }

        m_vKafkaProducer.push_back(pKafkaProducer);
    }

    LOG_DEBUG("create kafka producer size %d for the exporter %s", m_vKafkaProducer.size(), m_name.c_str());

    return GSE_SUCCESS;
}

bool KafkaExporter::startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig)
{
    StreamToCluster *streamToCluster = &ptrChannelIDConfig->m_streamToCluster;
    m_name = streamToCluster->m_name;
    std::size_t max_count = streamToCluster->m_kafkaCluster->m_addresses.size();

    for (std::size_t kafka_idx = 0; kafka_idx < max_count; ++kafka_idx)
    {
        if (!m_kafkaBrokers.empty())
        {
            m_kafkaBrokers.append(",");
        }
        Address &address = streamToCluster->m_kafkaCluster->m_addresses.at(kafka_idx);
        m_kafkaBrokers.append(address.m_ip);
        m_kafkaBrokers.append(":");
        m_kafkaBrokers.append(gse::tools::strings::ToString(address.m_port));
    }

    m_kafkaConfig = ptrChannelIDConfig->m_streamToCluster.m_kafkaCluster->m_kafkaConfig;

    LOG_DEBUG("start exporter(%s) witch channel id config(address:%s)",
              m_name.c_str(), m_kafkaBrokers.c_str());

    m_producerNum = KAFKA_MAX_PRODUCER;

    return true;
}

bool KafkaExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    // clear first
    clear();

    m_kafkaBrokers = ptrExporterConf->m_kafkaConf->m_cluster;
    m_producerNum = ptrExporterConf->m_kafkaConf->m_producerNum;
    m_defaultTopicName = ptrExporterConf->m_kafkaConf->m_defaultTopicName;
    m_name = ptrExporterConf->m_name;

    if (m_producerNum <= 0)
    {
        m_producerNum = KAFKA_MAX_PRODUCER;
    }

    m_kafkaConfig = ptrExporterConf->m_kafkaConf->m_kafkaConfig;
    return true;
}

int KafkaExporter::Start()
{
    // struct timeval timer_10s = { 10, 0 };
    if (0 > m_eventManager.RegisterMsecTimerPersistEvent(KafkaExporter::KafkaPoll, this, 10))
    {
        LOG_ERROR("failed to register kafka poll timer");
        return GSE_ERROR;
    }

    std::string m_advertiseIP = m_configPtr->GetAdvertiseIp();
    if (m_advertiseIP == "")
    {
        m_advertiseIP = gse::tools::net::GetMachineIp();
    }

    m_eventManager.Start();

    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
    case ExporterConfigWrapper::ChannelIDConfType:
        m_upConfLock.RLock();
        if (!startWithChannelID(m_ptrConfWrapper->m_conf.m_ptrChannelIdExporterConfig))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_DEBUG("start kafka exporter (%s) by channel id config", SAFE_CSTR(m_name.c_str()));
        return createKafkaProducers();
    case ExporterConfigWrapper::DataFlowConfType:
        m_upConfLock.RLock();
        if (!startWithDataFlow(m_ptrConfWrapper->m_conf.m_exporterConf))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_DEBUG("start kafka exporter (%s) by data flow config", SAFE_CSTR(m_name.c_str()));
        return createKafkaProducers();
    default:
        break;
    }

    return GSE_ERROR;
}

int KafkaExporter::Stop()
{
    clear();

    return GSE_SUCCESS;
}

void KafkaExporter::FormatKey(DataCell *pDataCell, std::string &key)
{
    char tagkey[1024] = {0};
    uint32_t now = TimeCenter::Instance()->GetDateTime();
    uint32_t timestampMin = (now / 60) * 60;

    gse::tools::strings::Snprintf(tagkey, sizeof(tagkey), "%u|%lu|%s|%d", pDataCell->GetChannelID(), timestampMin, m_advertiseIP.c_str(), pDataCell->GetServerPort());
    pDataCell->SetOutputTag(std::string(tagkey));

    static uint32_t index = 0;
    index++;
    key = "ds=" + gse::tools::strings::ToString((uint64_t)now) + "&" + "tag=" + tagkey + "&index=" + gse::tools::strings::ToString(index);
}

int KafkaExporter::Write(DataCell *pDataCell)
{

    if (NULL == pDataCell)
    {
        LOG_WARN("the parameter is valid. the pointer of data cell is [%x]", pDataCell);
        return GSE_ERROR;
    }

    pDataCell->DealLineBreak();
    int producerSize = m_vKafkaProducer.size();
    if (producerSize <= 0)
    {
        LOG_WARN("there is no kafka producer in exporter[%s]", SAFE_CSTR(m_name.c_str()));
        return GSE_ERROR;
    }

    int producerIndex = (m_nextProducerId++) % producerSize;
    KafkaProducer *pKafkaProducer = m_vKafkaProducer[producerIndex];
    if (NULL == pKafkaProducer)
    {
        LOG_WARN("the pointer of kafka producer is NULL in exporter[%s]. the index of kafka producer array is [%d]", SAFE_CSTR(m_name.c_str()), producerIndex);
        return GSE_ERROR;
    }

    std::string value(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());
    std::string messageKey;
    FormatKey(pDataCell, messageKey);

    std::vector<std::string> str_topics;
    pDataCell->GetTableName(str_topics);

    if (str_topics.empty())
    {
        if (m_defaultTopicName.empty())
        {
            LOG_INFO("not found the topic for the channel id (%u), defualt topic name not config:%s", pDataCell->GetChannelID(), m_defaultTopicName.c_str());
            pDataCell->SetErrorMsg("no topic", OPS_ERROR_KAFKA_NO_TOPIC_ERROR);
            DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            return GSE_ERROR;
        }
        LOG_INFO("not found the topic for the channel id (%u), use defualt topic name:%s", pDataCell->GetChannelID(), m_defaultTopicName.c_str());
        str_topics.push_back(m_defaultTopicName);
        // return GSE_ERROR;
    }

    std::size_t max_count = str_topics.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string str_topic = str_topics.at(idx);
        LOG_DEBUG("will send the data (%s) to kafka with the topic (%s),index(%d), the channelid (%d) the key (%s) exporter name (%s)",
                  (value.c_str()), (str_topic.c_str()), idx, pDataCell->GetChannelID(), (messageKey.c_str()), (m_name.c_str()));

        if (GSE_SUCCESS != pKafkaProducer->ExcuteProduce(str_topic, pDataCell->GetPartition(), messageKey, value, pDataCell->GetSourceIp()))
        {
            pDataCell->SetErrorMsg("kafka produce failed", OPS_ERROR_KAFKA_ERROR);
            DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            LOG_WARN("it is failed to send the data (%s) to kafka with the topic(%s), the channelid (%d), the key (%s) exporter name (%s)",
                     value.c_str(), str_topic.c_str(), pDataCell->GetChannelID(),
                     messageKey.c_str(), m_name.c_str());
        }
    }

    pDataCell->SetOutputType("kafka");
    pDataCell->SetOutputAddress(m_kafkaBrokers);
    DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
    return GSE_SUCCESS;
}

void KafkaExporter::clear()
{
    if (!m_threadStoped)
    {
        m_eventManager.clear();
        m_eventManager.Stop();
        m_threadStoped = true;
    }

    for (int i = 0; i < m_vKafkaProducer.size(); i++)
    {
        KafkaProducer *pKafkaProducer = m_vKafkaProducer[i];
        if (pKafkaProducer != NULL)
        {
            pKafkaProducer->CloseProducer();
            delete pKafkaProducer;
        }
    }

    m_vKafkaProducer.clear();
}

void KafkaExporter::KafkaPoll(int fd, short what, void *v)
{
    KafkaExporter *ptr_this = reinterpret_cast<KafkaExporter *>(v);

    for (int i = 0; i < ptr_this->m_vKafkaProducer.size(); i++)
    {
        KafkaProducer *pKafkaProducer = ptr_this->m_vKafkaProducer[i];
        if (pKafkaProducer != NULL)
        {
            pKafkaProducer->KafkaPoll();
        }
    }
}
} // namespace data
} // namespace gse
