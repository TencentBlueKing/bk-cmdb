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

#include "pulsarexporter.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/strings.h"

#include "dataserver.h"
#include "filter/channelid_filter.h"

namespace gse {
namespace data {

PulsarExporter::PulsarExporter()
{
    m_nextProducerId = 0;
    m_producerNum = PULSAR_MAX_PRODUCER;
    m_tlsTrustCertsFilePath = "";
    m_tlsKeyFilePath = "";
    m_tlsCertFilePath = "";
    m_token = "";
}

PulsarExporter::~PulsarExporter()
{
    clear();
}

int PulsarExporter::createPulsarProducers()
{
    for (int i = 0; i < m_producerNum; i++)
    {
        PulsarProducer *producer = new PulsarProducer();
        if (NULL == producer)
        {
            LOG_WARN("fail to create pulsar producer [%d] object for exporter[%s]", i, SAFE_CSTR(m_name.c_str()));
            continue;
        }

        int iRet = producer->createProducer(m_serivce_url, m_tlsCertFilePath, m_tlsKeyFilePath, m_token);
        if (iRet != GSE_SUCCESS)
        {
            LOG_WARN("fail to start pulsar producer[%d] exporter[%s] with service url[%s]", i, SAFE_CSTR(m_name.c_str()), SAFE_CSTR(m_serivce_url.c_str()));
            delete producer;
            return GSE_ERROR;
        }

        m_pulsarPorducers.push_back(producer);
    }

    LOG_DEBUG("create pulsar producer size %d for the exporter %s", m_pulsarPorducers.size(), m_name.c_str());

    return GSE_SUCCESS;
}

bool PulsarExporter::startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig)
{
    StreamToCluster *streamToCluster = &ptrChannelIDConfig->m_streamToCluster;
    m_name = streamToCluster->m_name;

    m_serivce_url = "pulsar://" + streamToCluster->m_pulsarCluster->m_addresses[0].m_ip + ":" +
                    gse::tools::strings::ToString(streamToCluster->m_pulsarCluster->m_addresses[0].m_port);

    LOG_DEBUG("pulsar url:%s", m_serivce_url.c_str());
    m_token = streamToCluster->m_pulsarCluster->m_token;

    if (m_producerNum <= 0 || m_producerNum > PULSAR_MAX_PRODUCER)
    {
        m_producerNum = PULSAR_MAX_PRODUCER;
    }

    return true;
}

bool PulsarExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    m_producerNum = ptrExporterConf->m_pulsarConf->m_producerNum;
    m_topicName = ptrExporterConf->m_pulsarConf->m_topicName;
    m_token = ptrExporterConf->m_pulsarConf->m_token;
    m_tlsTrustCertsFilePath = ptrExporterConf->m_pulsarConf->m_tlsTrustCertsFilePath;
    m_tlsKeyFilePath = ptrExporterConf->m_pulsarConf->m_tlsKeyFilePath;
    m_tlsCertFilePath = ptrExporterConf->m_pulsarConf->m_tlsCertFilePath;
    LOG_DEBUG("start pulsar data flow, producernum:%d, m_token:%s",
              ptrExporterConf->m_pulsarConf->m_producerNum, m_token.c_str());
    m_serivce_url = ptrExporterConf->m_pulsarConf->m_serviceUrl;
    m_name = ptrExporterConf->m_name;
    LOG_DEBUG("start pulsar with dataflow config, producernum:%d, url:%s, name:%s",
              m_producerNum, m_serivce_url.c_str(), m_name.c_str());
    if (m_producerNum <= 0)
    {
        m_producerNum = PULSAR_MAX_PRODUCER;
    }
    return true;
}

int PulsarExporter::Start()
{
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
        LOG_DEBUG("start pulsar exporter (%s) by zk config", SAFE_CSTR(m_name.c_str()));
        return createPulsarProducers();
    case ExporterConfigWrapper::DataFlowConfType:
        m_upConfLock.RLock();
        if (!startWithDataFlow(m_ptrConfWrapper->m_conf.m_exporterConf))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        LOG_DEBUG("start pulsar exporter (%s) by data flow config", SAFE_CSTR(m_name.c_str()));
        return createPulsarProducers();
    default:
        break;
    }

    return GSE_ERROR;
}

int PulsarExporter::Stop()
{
    clear();

    return GSE_SUCCESS;
}

int PulsarExporter::Write(DataCell *pDataCell)
{
    if (NULL == pDataCell)
    {
        LOG_WARN("the parameter is valid. the pointer of data cell is [%x]", pDataCell);
        return GSE_ERROR;
    }

    pDataCell->DealLineBreak();

    int producerSize = m_pulsarPorducers.size();
    if (producerSize <= 0)
    {
        LOG_WARN("there is no pulsar producer in exporter[%s]", SAFE_CSTR(m_name.c_str()));
        return GSE_ERROR;
    }

    int producerIndex = (m_nextProducerId++) % producerSize;
    PulsarProducer *pulsarProducer = m_pulsarPorducers[producerIndex];
    if (NULL == pulsarProducer)
    {
        LOG_WARN("the pointer of pulsar producer is NULL in exporter[%s]. the index of kafka producer array is [%d]", SAFE_CSTR(m_name.c_str()), producerIndex);
        return GSE_ERROR;
    }

    std::string value(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());

    std::string str_key;
    std::vector<std::string> str_topics;
    pDataCell->GetTableName(str_topics);
    pDataCell->GetDataKey(str_key);

    if (str_topics.empty())
    {
        LOG_INFO("not found the topic for the channel id (%u), use dataflow topic(%s)", pDataCell->GetChannelID(), m_topicName.c_str());
        str_topics.push_back(m_topicName);
    }

    std::size_t max_count = str_topics.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string str_topic = str_topics.at(idx);
        LOG_DEBUG("will send the data (%s) to pulsar with the topic (%s), the channelid (%d) the key (%s) exporter name (%s)", SAFE_CSTR(value.c_str()), SAFE_CSTR(str_topic.c_str()), pDataCell->GetChannelID(), SAFE_CSTR(str_key.c_str()), SAFE_CSTR(m_name.c_str()));
        if (GSE_SUCCESS != pulsarProducer->excuteProduce(str_topic, pDataCell->GetPartition(), str_key, value))
        {
            LOG_ERROR("it is failed to send the data (%s) to pulsar with the topic(%s), the channelid (%d), the key (%s) exporter name (%s)", SAFE_CSTR(value.c_str()), SAFE_CSTR(str_topic.c_str()), pDataCell->GetChannelID(), SAFE_CSTR(str_key.c_str()), SAFE_CSTR(m_name.c_str()));
            pDataCell->SetErrorMsg("pulsar produce failed", OPS_ERROR_PULSAR_ERROR);
            DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        }
    }

    pDataCell->SetOutputType("pulsar");
    pDataCell->SetOutputAddress(m_serivce_url);
    DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
    return GSE_SUCCESS;
}

void PulsarExporter::clear()
{
    for (int i = 0; i < m_pulsarPorducers.size(); i++)
    {
        PulsarProducer *producer = m_pulsarPorducers[i];
        if (producer != NULL)
        {
            producer->closeProducer();
            delete producer;
        }
    }

    m_pulsarPorducers.clear();
}

} // namespace data
} // namespace gse
