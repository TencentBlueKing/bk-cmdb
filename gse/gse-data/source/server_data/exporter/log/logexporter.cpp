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

#include "logexporter.h"
#include "tools/time.h"
#include "tools/macros.h"
#include "bbx/gse_errno.h"
#include "tools/error.h"
#include "log/log.h"
//#include "log/boost_log.h"

#include "dataserver.h"
#include "muduo_log/muduo_log_helper.h"

namespace gse
{
namespace data
{

LogExporter::LogExporter()
{
    m_defaultLogSink = -1;
    m_needAllSourceIDS = false;
    m_queueLimitMaxCount = 0;
    m_logSink = -1;
}

LogExporter::~LogExporter()
{
}

int LogExporter::Stop()
{
    return GSE_SUCCESS;
}
bool LogExporter::startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig)
{

    return true;
}

bool LogExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    m_name = ptrExporterConf->m_name;
    if (NULL == ptrExporterConf->m_logConf)
    {
        LOG_WARN("fail to start exporter[%s], because there is no log exporter configure", SAFE_CSTR(m_name.c_str()));
        return false;
    }

    m_queueLimitMaxCount = ptrExporterConf->m_logConf->m_maxQueueLimitNum;
    LOG_INFO("the log exporter (%s) data queue limit max is (%d)", SAFE_CSTR(m_name.c_str()), m_queueLimitMaxCount);
    std::size_t max_count = ptrExporterConf->m_logConf->m_sourcedIDS.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string source_id = ptrExporterConf->m_logConf->m_sourcedIDS.at(idx);
        if (source_id.compare("*") == 0)
        {
            m_needAllSourceIDS = true;
            LOG_DEBUG("split all data into a directory by source id");
            break;
        }

        createLogSink(source_id);
    }

    //m_defaultLogSink = BoostLogKeeper::Instance().createSink("DEBUG", ptrExporterConf->m_logConf->m_logPath, ptrExporterConf->m_logConf->m_fileName, ptrExporterConf->m_logConf->m_maxLogSize, ptrExporterConf->m_logConf->m_maxFileNum, 3600);
    std::string logprefix = ptrExporterConf->m_logConf->m_logPath;
    std::string filename = ptrExporterConf->m_logConf->m_logPath + "/" + ptrExporterConf->m_logConf->m_fileName;
    char name[256] = {0};
    size_t kOneM = 1024 * 1024;
    size_t logsize = ptrExporterConf->m_logConf->m_maxLogSize * kOneM;
    m_defaultLogSink = MuduoLogHelper::Instance().CreateDataSink(filename, logsize, ptrExporterConf->m_logConf->m_maxFileNum, 3600);

    LOG_INFO("success to start exporter[%s] with default logsink[%d]", SAFE_CSTR(m_name.c_str()), m_defaultLogSink);
    return true;
}

int LogExporter::Start()
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
        return GSE_SUCCESS;
    case ExporterConfigWrapper::DataFlowConfType:
        m_upConfLock.RLock();
        if (!startWithDataFlow(m_ptrConfWrapper->m_conf.m_exporterConf))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        return GSE_SUCCESS;
    default:
        break;
    }

    LOG_INFO("it is failed to start exporter[%s] with the config type [%d]", SAFE_CSTR(m_name.c_str()), m_ptrConfWrapper->m_exporterConfTypeEnum);

    return GSE_ERROR;
}

LogExporter::LogSink LogExporter::createLogSink(const std::string &source_id)
{
    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
    case ExporterConfigWrapper::ChannelIDConfType:
        return -1;
    case ExporterConfigWrapper::DataFlowConfType:
    {
        LogExporterConf *ptr_log_conf = m_ptrConfWrapper->m_conf.m_exporterConf->m_logConf;
        std::string logprefix = ptr_log_conf->m_logPath;
        std::string filename = ptr_log_conf->m_logPath + "/" + source_id + "/" + ptr_log_conf->m_fileName;
        char name[256] = {0};
        size_t kOneM = 1024 * 1024;
        size_t logsize = ptr_log_conf->m_maxLogSize * kOneM;
        LogSink tmp_log_sink = MuduoLogHelper::Instance().CreateDataSink(filename, logsize, ptr_log_conf->m_maxFileNum, 3600);
        m_sourceIDToSink.Insert(source_id, tmp_log_sink);
        return tmp_log_sink;
    }
    default:
        break;
    }

    return -1;
}

int LogExporter::Write(DataCell *pDataCell)
{
    if (NULL == pDataCell)
    {
        LOG_WARN("the pointer of datacell which want write by logexporter[%s]", SAFE_CSTR(m_name.c_str()));
        return GSE_SYSTEMERROR;
    }

    pDataCell->DealLineBreak();

    std::string source_id = pDataCell->GetExtensionByIndex(0);

    LogSink tmp_log_sink = -1;
    if (!source_id.empty() && m_sourceIDToSink.Find(source_id, tmp_log_sink))
    {
        // 以sourceid 分目录存储
        MUDUO_LOG_DATA(tmp_log_sink, pDataCell->GetDataBuf());
        LOG_DEBUG("source id (%s) data is :%s", SAFE_CSTR(source_id.c_str()), pDataCell->GetDataBuf());
    }
    else if (!source_id.empty() && m_needAllSourceIDS)
    {
        // 所有sourceid 需要分目录，为发现目标sourceid 的 sink 新建一个
        tmp_log_sink = createLogSink(source_id);
        MUDUO_LOG_DATA(tmp_log_sink, pDataCell->GetDataBuf());
    }
    else
    {
        // 默认不以 sourceid 来分目录
        MUDUO_LOG_DATA(m_defaultLogSink, pDataCell->GetDataBuf());
    }

    OPMetric::LogMsgInc();
    std::string exporter_tag = "file";
    pDataCell->SetOutputType(exporter_tag);
    pDataCell->SetOutputAddress("unkown");
    DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));

    return GSE_SUCCESS;
}
}
}
