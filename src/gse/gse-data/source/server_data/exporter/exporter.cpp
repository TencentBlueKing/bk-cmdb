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

#include "exporter.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"

#include "dataserver.h"

namespace gse {
namespace data {

Exporter::Exporter()
{
    m_ptrConfWrapper = NULL;
    m_next = NULL;
    m_setStopTimestamp = 0;
    m_id = 0;
    m_ptrOPSReport = NULL;
}

Exporter::~Exporter()
{
    if (m_ptrConfWrapper != NULL)
    {
        delete m_ptrConfWrapper;
        m_ptrConfWrapper = NULL;
    }

    if (NULL != m_next)
    {
        delete m_next;
    }

    LOG_INFO("will to delete the exporter '%s'", SAFE_CSTR(m_name.c_str()));

    // ChannelIDFilter 此处仅保留了指针，不需要释放，由出生地管理
    m_filters.Clear();
}

void Exporter::UpdateConf(ExporterConfigWrapper::ExporterConfigTypeEnum type, void *conf)
{
    m_upConfLock.WLock();
    if (NULL != m_ptrConfWrapper)
    {
        delete m_ptrConfWrapper;
    }
    m_ptrConfWrapper = new ExporterConfigWrapper(type, conf);
    m_upConfLock.UnLock();
}

void Exporter::SetNext(Exporter *ptrNext)
{
    if (NULL != m_next)
    {
        m_next->SetNext(ptrNext);
        return;
    }

    m_next = ptrNext;
}

void Exporter::SetNeedStop()
{
    m_setStopTimestamp = gse::tools::time::GetUTCSecond();
}

bool Exporter::IsNeedStop()
{
    // 超过60 秒即认为达到该被清理的标准
    return m_setStopTimestamp == 0 ? false : ((gse::tools::time::GetUTCSecond() - m_setStopTimestamp) > 120);
}
std::string Exporter::GetExporterName()
{
    return m_name;
}
void Exporter::SetExporterId(int id)
{
    m_id = id;
}
int Exporter::GetExporterId()
{
    return m_id;
}

bool Exporter::FilterByDataID(DataCell *pDataCell)
{
    LOG_DEBUG("NOT IMPLEMENT FOR DATA FLOW FILTER FOR THE CHANNELID (%d)", pDataCell->GetChannelID());
    return true;
}

bool Exporter::FilterByChannelID(DataCell *pDataCell)
{
    pDataCell->ClearTableNames();

    bool bret = false;
    LOG_DEBUG("the filter(%s)'s count is (%d)", SAFE_CSTR(m_name.c_str()), m_filters.Size());

    std::vector<ChannelIDFilter *> tmp_filters;
    m_filters.Copy(tmp_filters);
    std::vector<ChannelIDFilter *>::iterator iter = tmp_filters.begin();
    for (; iter != tmp_filters.end(); ++iter)
    {
        ChannelIDFilter *ptr_filters = *iter;
        if (NULL == ptr_filters)
        {
            continue;
        }

        /*
        if(pDataCell->GetExtensionSize() == 0)
        {
            LOG_DEBUG("filter successfuly, the input data extensions is empty, direct pass , and update the topic (%s) for the channelid (%d)", SAFE_CSTR(ptr_filters->m_tableName.c_str()), pDataCell->GetChannelID());
            pDataCell->AddTableName(ptr_filters->m_tableName);
            continue;
        }
        */
        LOG_DEBUG("the stream to name (%s) config name (%s) table    (%s)", SAFE_CSTR(ptr_filters->m_streamToName.c_str()), SAFE_CSTR(ptr_filters->m_configName.c_str()), SAFE_CSTR(ptr_filters->m_tableName.c_str()));
        if (NULL == ptr_filters)
        {
            LOG_DEBUG("not set any filter for the channel id %d", pDataCell->GetChannelID());
            continue;
        }

        if (!ptr_filters->IsValidData(pDataCell))
        {
            continue;
        }

        LOG_DEBUG("filter successfuly , and update the topic (%s) for the channelid (%d)", SAFE_CSTR(ptr_filters->m_tableName.c_str()), pDataCell->GetChannelID());
        std::string topicname;
        ptr_filters->GetTopicName(topicname);
        pDataCell->AddTableName(topicname);
        return true;
    }

    return bret;
}

bool Exporter::FilterByDataFlow(DataCell *pDataCell)
{
    LOG_DEBUG("NOT IMPLEMENT FOR DATA FLOW FILTER FOR THE CHANNELID (%d)", pDataCell->GetChannelID());

    return true;
}

void Exporter::AddFilter(ChannelIDFilter *ptrChannelIDFilter)
{
    if (NULL != ptrChannelIDFilter)
    {
        LOG_DEBUG("add filter (%s) for the exporter '%s'", SAFE_CSTR(ptrChannelIDFilter->m_configName.c_str()), SAFE_CSTR(m_name.c_str()));
        m_filters.PushBack(ptrChannelIDFilter);
    }
}

void Exporter::DumpConfig(string &config)
{
    m_ptrConfWrapper->DumpConfig(config);
}

void Exporter::DumpStatus(string &config)
{
}

bool Exporter::Filter(DataCell *pDataCell)
{
    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
    case ExporterConfigWrapper::ChannelIDConfType:
        return FilterByChannelID(pDataCell);
    case ExporterConfigWrapper::DataFlowConfType:
        return FilterByDataFlow(pDataCell);
    default:
        break;
    }

    LOG_WARN("not found any filter for the channelid (%d)", pDataCell->GetChannelID());
    return false;
}

int Exporter::WriteTo(DataCell *pDataCell)
{

    if (!Filter(pDataCell))
    {
        LOG_DEBUG("it is failed to send the data for the channel id (%u), the data is not fit for the filter", pDataCell->GetChannelID());
        return GSE_ERROR;
    }

    // int delaySec = 0;
    // TIME_DEALY_SECOND(pDataCell->GetCreationTimestamp(), delaySec);
    // OPMetric::AddMsgDelayMsCounter(kdispathMessageDirection, delaySec);

    if (GSE_SUCCESS != Write(pDataCell))
    {
        LOG_WARN("fail to write data into exporter[%s]", SAFE_CSTR(m_name.c_str()));
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
    }
    if (NULL != m_next)
    {
        m_next->WriteTo(pDataCell);
    }

    return GSE_SUCCESS;
}

std::string Exporter::TypeToStrName(int type)
{
    static std::map<int, std::string> exporterTypeName =
        {
            {KAFKA_COMMON, "kafka"},
            {REDIS_SENTINEL_PUB, "redis sentinel"},
            {REDIS_PUB, "redis stand-alone"},
            {EXPORT_PULSAR, "pulsar"}};

    auto it = exporterTypeName.find(type);
    if (it != exporterTypeName.end())
    {
        return it->second;
    }
    return "";
}

} // namespace data
} // namespace gse
