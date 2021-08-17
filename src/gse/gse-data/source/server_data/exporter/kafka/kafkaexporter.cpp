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
#include "log/log.h"
#include "tools/error.h"
#include "tools/macros.h"
#include "bbx/gse_errno.h"
#include "tools/net.h"

#include "utils.h"
#include "dataserver.h"

#include "filter/channelid_filter.h"
namespace gse { 
namespace dataserver {

KafkaExporter::KafkaExporter()
{
    m_nextProducerId = 0;
    m_kafkaMaxQueue = DEFAULT_MAX_KAFKA_QUEUE_SIZE;
    m_kafkaMaxMessageBytes = DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE;
    m_producerNum = KAFKA_MAX_PRODUCER;

    m_selfIp = gse::tools::net::GetMachineIp();
}

KafkaExporter::~KafkaExporter()
{
    clear();
}

int KafkaExporter::createKafkaProducers()
{
    KafkaProducer::m_runtimeDataDirector = m_gseCfg.m_runtimePath;
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
        int iRet = pKafkaProducer->createProducer(m_kafkaBrokers);
        if (iRet != GSE_SUCCESS)
        {
            LOG_WARN("fail to start kafka producer[%d] exporter[%s] with kafka brokers[%s]", i, SAFE_CSTR(m_name.c_str()), SAFE_CSTR(m_kafkaBrokers.c_str()));
            delete pKafkaProducer;
            continue;
        }

        m_vKafkaProducer.push_back(pKafkaProducer);
    }

    LOG_DEBUG("create kafka producer size %d for the exporter %s", m_vKafkaProducer.size(), m_name.c_str());

    return GSE_SUCCESS;
}

bool KafkaExporter::startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig)
{
    m_name = ptrChannelIDConfig->m_name;
    std::size_t max_count = ptrChannelIDConfig->m_storage.m_kafkaStorage->m_addresses.size();
    for (std::size_t kafka_idx = 0; kafka_idx < max_count; ++kafka_idx)
    {
        if (!m_kafkaBrokers.empty())
        {
            m_kafkaBrokers.append(",");
        }
        Address &address = ptrChannelIDConfig->m_storage.m_kafkaStorage->m_addresses.at(kafka_idx);
        m_kafkaBrokers.append(address.m_ip);
        m_kafkaBrokers.append(":");
        m_kafkaBrokers.append(gse::tools::strings::ToString(address.m_port));
    }
    m_kafkaConfig = ptrChannelIDConfig->m_storage.m_kafkaStorage->m_kafkaConfig;

    LOG_DEBUG("start exporter(%s) witch channel id config(address:%s, kafka config[security protocol:%s, sasl mechanisms:%s, user name:%s, buffering max meesage:%s])",
              ptrChannelIDConfig->m_name.c_str(), m_kafkaBrokers.c_str(),
              ptrChannelIDConfig->m_storage.m_kafkaStorage->m_kafkaConfig.m_securityProtocol.c_str(),
              ptrChannelIDConfig->m_storage.m_kafkaStorage->m_kafkaConfig.m_saslMechanisms.c_str(),
              ptrChannelIDConfig->m_storage.m_kafkaStorage->m_kafkaConfig.m_saslUserName.c_str(),
              ptrChannelIDConfig->m_storage.m_kafkaStorage->m_kafkaConfig.m_queueBufferingMaxMessages.c_str());

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
bool KafkaExporter::startWithDataID(StorageConfigType *ptrStorageConfig)
{
    m_kafkaBrokers.append(ptrStorageConfig->m_host);
    m_kafkaBrokers.append(":");
    m_kafkaBrokers.append(gse::tools::strings::ToString(ptrStorageConfig->m_port));

    if (ptrStorageConfig->m_maxKafkaMaxQueue <= 0)
    {
        m_kafkaMaxQueue = DEFAULT_MAX_KAFKA_QUEUE_SIZE;
    }
    else
    {
        m_kafkaMaxQueue = ptrStorageConfig->m_maxKafkaMaxQueue;
    }

    if (ptrStorageConfig->m_maxKafkaMessageBytes <= 0)
    {
        m_kafkaMaxMessageBytes = DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE;
    }
    else
    {
        m_kafkaMaxMessageBytes = ptrStorageConfig->m_maxKafkaMessageBytes;
    }

    m_producerNum = KAFKA_MAX_PRODUCER;
    m_name = m_kafkaBrokers;
    m_kafkaConfig.m_messageMaxBytes = gse::tools::strings::ToString(m_kafkaMaxMessageBytes);
    m_kafkaConfig.m_queueBufferingMaxMessages = gse::tools::strings::ToString(m_kafkaMaxQueue);
    return true;
}

int KafkaExporter::Start()
{
    //struct timeval timer_10s = { 10, 0 };
    if (0 > m_eventManager.registerTimerPersistEvent(KafkaExporter::KafkaPoll, this, 60))
    {
        LOG_ERROR("FAIL TO REGISTER TIMER , REPORT LOAD BALANCE FUNCTION");
        return GSE_ERROR;
    }

    m_eventManager.Start();

    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
    case ExporterConfigWrapper::DataIDConfType:
        m_upConfLock.RLock();
        if (!startWithDataID(m_ptrConfWrapper->m_conf.m_ptrDataIDConfig))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_DEBUG("start kafka exporter (%s) by data flow config", SAFE_CSTR(m_name.c_str()));
        return createKafkaProducers();
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

    char tagkey[1024] = {0};
    uint32_t now = TimeCenter::Instance()->GetDateTime();
    uint32_t timestamp_min = now / 60 * 60;
    std::string str_server_ip;
    uint32_t server_ip  = pDataCell->GetServerIp();
    if (server_ip > 0)
    {
        str_server_ip = pDataCell->GetStrServerIp();
    }
    else
    {
        str_server_ip = m_selfIp;
    }

    gse::tools::strings::Snprintf(tagkey, sizeof(tagkey), "%u|%lu|%s|%d", pDataCell->GetChannelID(), timestamp_min, SAFE_CSTR(str_server_ip.c_str()), pDataCell->GetServerPort());
    static uint32_t index  = 0;
    index++;
    std::string str_key = "ds=" + gse::tools::strings::ToString((uint64_t)now) + "&" + "tag=" + tagkey + "&index="+gse::tools::strings::ToString(index);
    pDataCell->SetOutputTag(std::string(tagkey));

    std::vector<std::string> str_topics;
    pDataCell->GetTableName(str_topics);

    if (str_topics.empty())
    {
        if (m_defaultTopicName.empty())
        {
            LOG_INFO("not found the topic for the channel id (%u), defualt topic name not config:%s", pDataCell->GetChannelID(), m_defaultTopicName.c_str());
            pDataCell->SetErrorMsg("no topic", OPS_ERROR_KAFKA_NO_TOPIC_ERROR);
            DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            return GSE_ERROR;
        }
        LOG_INFO("not found the topic for the channel id (%u), use defualt topic name:%s", pDataCell->GetChannelID(), m_defaultTopicName.c_str());
        str_topics.push_back(m_defaultTopicName);
        //return GSE_ERROR;
    }

    std::size_t max_count = str_topics.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string str_topic = str_topics.at(idx);
        LOG_DEBUG("will send the data (%s) to kafka with the topic (%s), the channelid (%d) the key (%s) exporter name (%s)", SAFE_CSTR(value.c_str()), SAFE_CSTR(str_topic.c_str()), pDataCell->GetChannelID(), SAFE_CSTR(str_key.c_str()), SAFE_CSTR(m_name.c_str()));
        if (GSE_SUCCESS != pKafkaProducer->excuteProduce(str_topic, pDataCell->GetPartition(), str_key, value, pDataCell->GetSourceIp()))
        {
            pDataCell->SetErrorMsg("kafka produce failed", OPS_ERROR_KAFKA_ERROR);
            DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            LOG_WARN("it is failed to send the data (%s) to kafka with the topic(%s), the channelid (%d), the key (%s) exporter name (%s)", SAFE_CSTR(value.c_str()), SAFE_CSTR(str_topic.c_str()), pDataCell->GetChannelID(), SAFE_CSTR(str_key.c_str()), SAFE_CSTR(m_name.c_str()));
        }
        OPMetric::KafkaMsgInc();
    }

    pDataCell->SetOutputType("kafka");
    pDataCell->SetOutputAddress(m_kafkaBrokers);
    DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
    return GSE_SUCCESS;
}

void KafkaExporter::clear()
{
    m_eventManager.clear();
    m_eventManager.stop();
    for (int i = 0; i < m_vKafkaProducer.size(); i++)
    {
        KafkaProducer *pKafkaProducer = m_vKafkaProducer[i];
        if (pKafkaProducer != NULL)
        {
            pKafkaProducer->closeProducer();
            delete pKafkaProducer;
        }
    }

    m_vKafkaProducer.clear();
}

void KafkaExporter::KafkaPoll(int fd, short what, void* v)
{
    KafkaExporter* ptr_this = reinterpret_cast<KafkaExporter*>(v);

    for (int i = 0; i < ptr_this->m_vKafkaProducer.size(); i++)
    {
        KafkaProducer *pKafkaProducer = ptr_this->m_vKafkaProducer[i];
        if (pKafkaProducer != NULL)
        {
            pKafkaProducer->KafkaPoll();
        }
    }
}
}
}
