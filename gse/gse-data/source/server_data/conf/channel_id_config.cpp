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

#include "channel_id_config.h"
#include <string.h>

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/error.h"
#include "tools/macros.h"
#include "tools/pipe.h"
#include "tools/time.h"

#include "conf/confItem.h"
#include "dataserver.h"
#include "exporter/dsproxy/proxy_exporter.h"
#include "exporter/exporter.h"
#include "exporter/kafka/kafkaexporter.h"
#include "exporter/pulsar/pulsarexporter.h"
#include "exporter/redis/redisexporter.h"
#include "ops/op_healthz.h"
#include "ops/op_metric.h"

namespace gse {
namespace data {

bool ChannelIDFilter::CheckAndFilter(DataCell *ptrDataCell)
{
    if (m_streamFilterAnd.empty())
    {
        return true;
    }

    std::string extension_str;
    ptrDataCell->GetExtensionString(extension_str);
    LOG_DEBUG("will to check and the extensions's '%s' with the filter '%s'", SAFE_CSTR(extension_str.c_str()), SAFE_CSTR(m_andFilterStringInProtocol.c_str()));
    // 因为 extension 包含的内容，一定是 andFilterString 包含内容的父集，所以此处相等时可以直接比对
    // 此处是为了优化检测速度，因为绝大多数情况下 不需要循环遍历所有条件
    if (!m_isNeedCheckAndFilterInData && extension_str.size() == m_andFilterStringInProtocol.size())
    {
        return extension_str.find(m_andFilterStringInProtocol) != std::string::npos;
    }

    std::size_t max_cnt = m_streamFilterAnd.size();

    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        StreamFilter *ptr_stream_filter = m_streamFilterAnd.at(idx);

        // Note: check in protocol
        if (ptr_stream_filter->m_fieldIn.compare("protocol") == 0)
        {
            if (CheckProtocolFilter(ptrDataCell, ptr_stream_filter))
            {
                continue;
            }

            return false;
        }

        if (ptr_stream_filter->m_fieldIn.compare("data") == 0)
        {
            if (CheckDataFilter(ptrDataCell, ptr_stream_filter))
            {
                // one filter match success,check next
                continue;
            }

            return false;
        }
    }

    return true;
}

bool ChannelIDFilter::CheckProtocolFilter(DataCell *ptrDataCell, StreamFilter *filter)
{
    std::string extension_str;
    ptrDataCell->GetExtensionString(extension_str);

    LOG_DEBUG("will check or extension's (%s) with the filter value (%s)", SAFE_CSTR(extension_str.c_str()), SAFE_CSTR(filter->m_fieldDataValue.c_str()));
    std::string tagfilter = "[" + filter->m_fieldDataValue + "]";
    if (extension_str.find(tagfilter) != std::string::npos)
    {
        return true;
    }

    return false;
}

std::string ChannelIDFilter::ToString()
{
    /*
    std::string m_streamToName;
    std::string m_configName;
    std::string m_tableName;
    std::vector<StreamFilter *> m_streamFilterAnd;
    std::vector<StreamFilter *> m_streamFilterOr;
    std::string m_andFilterStringInProtocol;
    bool m_isNeedCheckAndFilterInData;
    */

    std::string stream_filter_and;

    for (std::vector<StreamFilter *>::iterator it = m_streamFilterAnd.begin(); it != m_streamFilterAnd.end(); it++)
    {
        stream_filter_and.append((*it)->ToString());
    }

    std::string stream_filter_or;

    for (std::vector<StreamFilter *>::iterator it = m_streamFilterOr.begin(); it != m_streamFilterOr.end(); it++)
    {
        stream_filter_or.append((*it)->ToString());
    }

    std::string strbuff = "StreamTo:{\nStreamToName:%s, StorageConfigName:%s, TopicName:%s, StreamFilterAnd:%s, StreamFilterOr:%s\n}";
    std::string str_result;
    char buff[4096] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_streamToName.c_str(), m_configName.c_str(), m_tableName.c_str(), stream_filter_and.c_str(), stream_filter_or.c_str());
    str_result.append(buff);
    return str_result;
}

// Note:  check in data, this case only check the first filed
bool ChannelIDFilter::CheckDataFilter(DataCell *ptrDataCell, StreamFilter *filter)
{
    string table_name;
    char *pos = strstr(ptrDataCell->GetDataBuf(), filter->m_separator.c_str());
    if (pos != NULL)
    {
        table_name.assign(ptrDataCell->GetDataBuf(), (pos - ptrDataCell->GetDataBuf()));
    }
    else
    {
        table_name.assign(ptrDataCell->GetDataBuf());
    }

    if (table_name.compare(filter->m_fieldDataValue) != 0)
    {
        LOG_DEBUG("Filter rules do not match, logs:%s, streamfilter:%s, table name:%s", ptrDataCell->GetDataBuf(), filter->m_fieldDataValue.c_str(), table_name.c_str());
        return false;
    }
    LOG_DEBUG("Filter rules match success, logs:%s, streamfilter:%s", ptrDataCell->GetDataBuf(), filter->m_fieldDataValue.c_str());

    return true;
}

bool ChannelIDFilter::checkOrFilter(DataCell *ptrDataCell)
{
    if (m_streamFilterOr.size() == 0)
    {
        return true;
    }

    std::size_t max_cnt = m_streamFilterOr.size();
    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        StreamFilter *ptr_stream_filter = m_streamFilterOr.at(idx);

        // Note: check filter in protocol
        if (ptr_stream_filter->m_fieldIn.compare("protocol") == 0)
        {
            if (CheckProtocolFilter(ptrDataCell, ptr_stream_filter))
            {
                return true;
            }
            continue;
        }

        // Note: check filter in data, ignore the filed index ,only check first charactors
        if (ptr_stream_filter->m_fieldIn.compare("data") == 0)
        {
            if (CheckDataFilter(ptrDataCell, ptr_stream_filter))
            {
                return true;
            }
            continue;
        }
    }

    return false;
}

void ChannelIDFilter::GetTopicName(std::string &topicname)
{
    topicname.append(m_tableName);
}

uint32_t ChannelIDFilter::GetMemSize()
{
    uint32_t memSize = sizeof(ChannelIDFilter);
    memSize += m_streamToName.length();
    memSize += m_configName.length();
    memSize += m_tableName.length();
    memSize += m_andFilterStringInProtocol.length();

    for (auto filter : m_streamFilterAnd)
    {
        memSize += filter->GetMemSize();
    }

    for (auto filter : m_streamFilterOr)
    {
        memSize += filter->GetMemSize();
    }

    return memSize;
}

bool ChannelIDFilter::IsValidData(DataCell *ptrDataCell)
{
    if (m_streamFilterOr.empty())
    {
        return CheckAndFilter(ptrDataCell);
    }

    if (checkOrFilter(ptrDataCell))
    {
        return CheckAndFilter(ptrDataCell);
    }

    return false;
}

StreamFilter *ChannelIDConfig::GetFilter(const std::string &filterName)
{
    std::size_t max_count = m_streamFilter.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        StreamFilter *ptr_stream_filter = m_streamFilter.at(idx);
        if (ptr_stream_filter->m_name.compare(filterName) == 0)
        {
            return ptr_stream_filter;
        }
    }

    return NULL;
}

void ChannelIdStreamExporterManager::CleanInvalidExporter()
{
    std::list<Exporter *>::iterator it = m_needDeleteExporters.begin();
    for (; it != m_needDeleteExporters.end();)
    {
        if ((*it) != NULL)
        {
            if ((*it)->IsNeedStop())
            {
                Exporter *ptr_delete = (*it);
                it = m_needDeleteExporters.erase(it);
                LOG_DEBUG("stop exporter:(id:%d, name:%s)", ptr_delete->GetExporterId(), ptr_delete->GetExporterName().c_str());
                ptr_delete->Stop();
                delete ptr_delete;
            }
            else
            {
                ++it;
            }
        }
        else
        {
            it = m_needDeleteExporters.erase(it);
        }
    }
}

ChannelIdExporterConfig::ChannelIdExporterConfig()
{
    m_setDeleteTimestamp = 0;
    m_streamToId = 0;
}
ChannelIdExporterConfig::~ChannelIdExporterConfig()
{
}

ChannelIdExporterConfig &ChannelIdExporterConfig::operator=(ChannelIdExporterConfig &srcConf)
{
    this->m_streamToId = srcConf.m_streamToId;
    this->m_streamToCluster.CopyFrom(&srcConf.m_streamToCluster);
}

ChannelIdExporterConfig::ChannelIdExporterConfig(ChannelIdExporterConfig &srcConf)
{
    m_setDeleteTimestamp = 0;
    *this = srcConf;
}

ChannelIdStreamExporterManager::ChannelIdStreamExporterManager()
{
    m_cleanTime.tv_sec = 10;
    m_cleanTime.tv_usec = 0;
    m_exporterUpdateNotifyFd[0] = -1;
    m_exporterUpdateNotifyFd[1] = -1;
    m_evBase = NULL;
    m_cleanTimerEvent = NULL;
    m_exporterUpdateEvent = NULL;
    m_stoped = false;

    LOG_DEBUG("ChannelIdStreamExporterManager init");
}

void ChannelIdStreamExporterManager::FreeEvent()
{
    if (m_cleanTimerEvent != NULL)
    {
        event_free(m_cleanTimerEvent);
        m_cleanTimerEvent = NULL;
    }

    if (m_exporterUpdateEvent != NULL)
    {
        event_free(m_exporterUpdateEvent);
        m_exporterUpdateEvent = NULL;
    }

    for (int i = 0; i < 2; i++)
    {
        if (m_exporterUpdateNotifyFd[i] != -1)
        {
            gse::tools::pipe::ClosePipe(m_exporterUpdateNotifyFd[i]);
            m_exporterUpdateNotifyFd[i] = -1;
        }
    }

    if (m_evBase != NULL)
    {
        event_base_loopexit(m_evBase, NULL);
        event_base_free(m_evBase);
        m_evBase = NULL;
    }
}

ChannelIdStreamExporterManager::~ChannelIdStreamExporterManager()
{
    LOG_DEBUG("delete ChannelIdExporterManager obj(%p)", this);
    FreeEvent();
    std::list<Exporter *>::iterator it = m_needDeleteExporters.begin();
    for (; it != m_needDeleteExporters.end(); it++)
    {
        Exporter *ptr_delete = (*it);
        LOG_DEBUG("stop exporter:(id:%d, name:%s)", ptr_delete->GetExporterId(), ptr_delete->GetExporterName().c_str());
        it = m_needDeleteExporters.erase(it);
        ptr_delete->Stop();
        delete ptr_delete;
    }

    std::map<uint32_t, Exporter *>::iterator m_it = m_exporters.begin();
    for (; m_it != m_exporters.end();)
    {
        Exporter *ptr_delete = m_it->second;
        LOG_DEBUG("stop exporter:(id:%d, name:%s)", ptr_delete->GetExporterId(), ptr_delete->GetExporterName().c_str());
        ptr_delete->Stop();
        delete ptr_delete;
        m_exporters.erase(m_it++);
    }
}

void ChannelIdStreamExporterManager::InvalidExporterCleanTimerEvent(int fd, short which, void *v)
{
    ChannelIdStreamExporterManager *ptr_this = (ChannelIdStreamExporterManager *)v;
    ptr_this->CleanInvalidExporter();
}

int ChannelIdStreamExporterManager::UpdateExporterConfig(ZkEvent *event)
{
    int ret = gse::tools::pipe::WritePipe(m_exporterUpdateNotifyFd[1], &event, sizeof(event));
    return (ret > 0 ? GSE_SUCCESS : GSE_ERROR);
}

void ChannelIdStreamExporterManager::Stop()
{
    m_stoped = true;
    FreeEvent();
}

void ChannelIdStreamExporterManager::ExporterUpdateEvent(int fd, short which, void *v)
{
    ChannelIdStreamExporterManager *ptr_this = (ChannelIdStreamExporterManager *)v;
    ChannelIdExporterConfig *exporter_cfg = NULL;

    ZkEvent *event = NULL;
    if (gse::tools::pipe::ReadPipe(fd, &event, sizeof(event)) < 0)
    {
        int tmp_errno = gse_socket_errno;
        string tmp_errormsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to read pipe fd (%d), errno(%d), error message (%s)", fd, tmp_errno, SAFE_CSTR(tmp_errormsg.c_str()));
        return;
    }

    if (event == NULL)
    {
        LOG_ERROR("event msg invalid(null)");
        return;
    }

    switch (event->m_eventType)
    {
    case ZK_EVENT_LIST:
        break;
    case ZK_EVENT_DELETE:
    {
        uint32_t *ptr_channelid = (uint32_t *)event->m_msg;
        ptr_this->DeleteExporter(*ptr_channelid);
        delete ptr_channelid;
        break;
    }
    case ZK_EVENT_CHANGE: // no break
    case ZK_EVENT_CREATE:
    {
        exporter_cfg = (ChannelIdExporterConfig *)event->m_msg;
        ptr_this->StreamExporterUpdate(exporter_cfg);
        delete exporter_cfg;
        break;
    }
    default:
        LOG_ERROR("recv unkown event:%d", event->m_eventType);
    }

    delete event;
    return;
}

int ChannelIdStreamExporterManager::Init()
{
    return 0;
}

int ChannelIdStreamExporterManager::ThreadFun()
{
    LOG_DEBUG("ChannelIdStreamExporterManager running");

    int rc = -1;
    gse::tools::pipe::CreatePipe(m_exporterUpdateNotifyFd);
    m_evBase = event_base_new();
    if (NULL == m_evBase)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new event base object for channel id manager thread, errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    m_cleanTimerEvent = event_new(m_evBase, -1, EV_PERSIST, ChannelIdStreamExporterManager::InvalidExporterCleanTimerEvent, this);
    if (NULL == m_cleanTimerEvent)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new clean timer event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    rc = event_add(m_cleanTimerEvent, &m_cleanTime);
    if (rc < 0)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to start clean invalid channel id timer event. rc:%d, errno:%d, errmsg:%s", rc, tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    m_exporterUpdateEvent = event_new(m_evBase, m_exporterUpdateNotifyFd[0], EV_PERSIST | EV_READ,
                                      ChannelIdStreamExporterManager::ExporterUpdateEvent, this);
    if (NULL == m_exporterUpdateEvent)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new channel id update event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }

    rc = event_add(m_exporterUpdateEvent, NULL);
    if (rc < 0)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to start channel id update event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }

    do
    {
        rc = event_base_loop(m_evBase, 0);
        int tmpErrno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmpErrno);
        LOG_FATAL("exporter manager event base loop abnomal break. ret:%d, errno:%d, errmsg:%s", rc,
                  tmpErrno, SAFE_CSTR(tmp_errmsg.c_str()));
    } while (rc < 0 && !m_stoped);

    return GSE_SUCCESS;
}

Exporter *ChannelIdStreamExporterManager::CreateExporter(ChannelIdExporterConfig *ptr_stream_to_id_config)
{
    Exporter *ptr_exporter = NULL;
    std::string hostAddr;

    StreamToCluster *cluster = &ptr_stream_to_id_config->m_streamToCluster;

    switch (ptr_stream_to_id_config->m_streamToCluster.m_clusterType)
    {
    case KAFKA_COMMON:
        ptr_exporter = new KafkaExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        LOG_INFO("create a new kafka exporter for the storage %s", cluster->m_kafkaCluster->ToString().c_str());
        hostAddr = cluster->m_kafkaCluster->ToString();
        break;
    case REDIS_PUB:
        LOG_DEBUG("create redis client for the storage %s", cluster->m_redisCluster->ToString().c_str());
        ptr_exporter = new RedisExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        hostAddr = cluster->m_redisCluster->ToString();
        break;
    case REDIS_SENTINEL_PUB:
        LOG_DEBUG("create redis sentinel client for the storage %s", cluster->m_redisCluster->ToString().c_str());
        ptr_exporter = new RedisExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        hostAddr = cluster->m_redisCluster->ToString();
        break;
    case EXPORT_DSPROXY:
        LOG_DEBUG("create ds proxy client for the storage %s", cluster->m_dsproxy->ToString().c_str());
        ptr_exporter = new ProxyExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        hostAddr = cluster->m_dsproxy->ToString();
        break;
    case EXPORT_PULSAR:
        LOG_INFO("create a new pulsar exporter for the storage %s", cluster->m_pulsarCluster->ToString().c_str());
        ptr_exporter = new PulsarExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        hostAddr = cluster->m_pulsarCluster->ToString();

        break;
    default:
        LOG_ERROR("it is failed to create a exporter for the unkown storage(index: %d, type: %s)", ptr_stream_to_id_config->m_streamToId, cluster->m_reportMode.c_str());
        return NULL;
    }

    ptr_exporter->SetGseConf(DataServer::GetConfigPtr());
    if (GSE_SUCCESS != ptr_exporter->Start())
    {
        LOG_ERROR("it is failed to start the exporter (%s) ", SAFE_CSTR(ptr_exporter->GetExporterName().c_str()));
        delete ptr_exporter;
        // std::string errMsg = "failed to start exporter, host address: " + hostAddr;
        //  OpHealthZ::AddInitHealthInfo(ptr_exporter->TypeToStrName(cluster->m_clusterType) + "_" + gse::tools::strings::ToString(ptr_stream_to_id_config->m_streamToId), errMsg, -1);
        return NULL;
    }

    // OpHealthZ::AddInitHealthInfo(ptr_exporter->TypeToStrName(cluster->m_clusterType) + "_" + gse::tools::strings::ToString(ptr_stream_to_id_config->m_streamToId), "ok", 0);
    return ptr_exporter;
}

Exporter *ChannelIdStreamExporterManager::GetExport(uint32_t exporter_id)
{
    Exporter *ptr_exporter = NULL;
    m_rwLock.RLock();
    std::map<uint32_t, Exporter *>::iterator it = m_exporters.find(exporter_id);
    if (it != m_exporters.end())
    {
        ptr_exporter = it->second;
    }
    m_rwLock.UnLock();
    return ptr_exporter;
}

bool ChannelIdStreamExporterManager::Find(uint32_t exporterId)
{
    m_rwLock.RLock();
    auto it = m_exporters.find(exporterId);
    if (it != m_exporters.end())
    {
        m_rwLock.UnLock();
        return true;
    }
    m_rwLock.UnLock();
    return false;
}

int ChannelIdStreamExporterManager::DeleteExporter(uint32_t exporter_id)
{
    m_rwLock.WLock();
    std::map<uint32_t, Exporter *>::iterator it = m_exporters.find(exporter_id);
    if (it != m_exporters.end())
    {
        if (it->second != NULL)
        {
            it->second->SetNeedStop();
            m_needDeleteExporters.push_back(it->second);
        }
        LOG_DEBUG("delete exporter-id(%d), will stop old exporter(%p))", exporter_id, it->second);
        m_exporters.erase(it);
    }
    m_rwLock.UnLock();
    return GSE_SUCCESS;
}

int ChannelIdStreamExporterManager::StreamExporterUpdate(ChannelIdExporterConfig *ptr_exporter_cfg)
{
    Exporter *exporter = NULL;
    exporter = CreateExporter(ptr_exporter_cfg);
    if (exporter == NULL)
    {
        LOG_ERROR("failed to update exporter, exporter(name:%s, id:%d, type:%d)",
                  ptr_exporter_cfg->m_streamToCluster.m_name.c_str(), ptr_exporter_cfg->m_streamToId, ptr_exporter_cfg->m_streamToCluster.m_clusterType);
        return GSE_ERROR;
    }
    exporter->SetExporterId(ptr_exporter_cfg->m_streamToId);
    m_rwLock.WLock();
    std::map<uint32_t, Exporter *>::iterator it = m_exporters.find(ptr_exporter_cfg->m_streamToId);
    if (it != m_exporters.end())
    {
        if (it->second != NULL)
        {
            it->second->SetNeedStop();
            m_needDeleteExporters.push_back(it->second);
        }
        LOG_DEBUG("update streamto_id(%d) exporter config, will stop old exporter(%p)", ptr_exporter_cfg->m_streamToId, it->second);
        m_exporters.erase(it);
    }

    m_exporters.insert(std::make_pair(ptr_exporter_cfg->m_streamToId, exporter));
    m_rwLock.UnLock();
    LOG_DEBUG("push exporter(%d:%p) into manager list", ptr_exporter_cfg->m_streamToId, exporter);
    return GSE_SUCCESS;
}

// channel id config manager
bool ChannelIdManager::WriteByFilter(DataCell *ptr_datacell)
{
    uint32_t channel_id = 0;

    if (!m_isPlatId)
    {
        channel_id = ptr_datacell->GetChannelID();
    }
    else
    {
        channel_id = GetPlatNum(ptr_datacell->GetChannelID());
    }

    ChannelIdStreamConfig *ptr_channel_id_stream_cfg = GetChannelStreamConfig(channel_id);
    LOG_DEBUG("get channel route by id:%d, config:%p, isplat:%d", channel_id, ptr_channel_id_stream_cfg, m_isPlatId);
    bool has_exporter = false;
    while (ptr_channel_id_stream_cfg != NULL)
    {
        has_exporter = true;
        std::string table_name;
        ptr_channel_id_stream_cfg->GetTableName(table_name);
        LOG_DEBUG("write date cell(channleid:%d) filter by channel name(%s), table name(%s)",
                  channel_id, ptr_channel_id_stream_cfg->m_channelName.c_str(), table_name.c_str());
        if (ptr_channel_id_stream_cfg->m_filter.IsValidData(ptr_datacell))
        {
            LOG_DEBUG("data cell(channleid:%d, stream_to_id:%d) filter success, will send by exporter. table name[%s]", ptr_datacell->GetChannelID(), ptr_channel_id_stream_cfg->m_streamToId, table_name.c_str());
            ptr_datacell->ClearTableNames();
            ptr_datacell->AddTableName(table_name);
            Exporter *exporter = DataServer::GetStreamExporterManagerInst()->GetExport(ptr_channel_id_stream_cfg->m_streamToId);
            if (exporter != NULL)
            {
                exporter->Write(ptr_datacell);
                LOG_DEBUG("data cell(channleid:%d) filter success, send by exporter[%s]", ptr_datacell->GetChannelID(), exporter->GetExporterName().c_str());
            }
            else
            {
                LOG_INFO("failed to send data cell(channleid:%d), no exporter(%d) config", channel_id, ptr_channel_id_stream_cfg->m_streamToId);
            }
        }
        ptr_channel_id_stream_cfg = ptr_channel_id_stream_cfg->m_next;
        if (ptr_channel_id_stream_cfg != NULL)
        {
            LOG_DEBUG("get channel next route by id:%d, config:%p, isplat:%d", channel_id, ptr_channel_id_stream_cfg, m_isPlatId);
        }
    }

    return has_exporter;
}

ChannelIdStreamConfig *ChannelIdManager::GetChannelStreamConfig(uint32_t channel_id)
{
    ChannelIdStreamConfig *ptr_channel_id_stream_cfg = NULL;
    m_rwLock.RLock();
    auto it = m_channelIds.find(channel_id);
    if (it != m_channelIds.end())
    {
        ptr_channel_id_stream_cfg = it->second;
    }
    m_rwLock.UnLock();
    return ptr_channel_id_stream_cfg;
}

void ChannelIdManager::GetChannelIdListByStreamId(uint32_t streamid, std::vector<uint32_t> &channelidList)
{
    m_rwLock.RLock();
    for (auto it : m_channelIds)
    {
        LOG_DEBUG("find channelid by streamid:%d, config(channelid:%d, streamid:%d)", streamid, it.second->m_channelID, it.second->m_streamToId);
        if ((it.second != nullptr))
        {
            ChannelIdStreamConfig *pNext = it.second;
            do
            {
                if (pNext->m_streamToId == streamid)
                {
                    channelidList.push_back(pNext->m_channelID);
                }
                pNext = pNext->m_next;
            } while (pNext != nullptr);
        }
    }
    m_rwLock.UnLock();
}

ChannelIdManager::ChannelIdManager(bool is_platid)
{
    m_evBase = NULL;
    m_cleanTimerEvent = NULL;
    m_channelIdUpdateEvent = NULL;
    m_cleanTime.tv_sec = 20;
    m_cleanTime.tv_usec = 0;
    m_isPlatId = is_platid;
    m_stoped = false;
}

ChannelIdManager::~ChannelIdManager()
{
    FreeEvent();
    std::list<ChannelIdStreamConfig *>::iterator it = m_needDeleteChannelIds.begin();
    for (; it != m_needDeleteChannelIds.end();)
    {
        ChannelIdStreamConfig *ptr_delete = (*it);
        if (ptr_delete != NULL)
        {
            delete ptr_delete;
        }

        it = m_needDeleteChannelIds.erase(it);
    }

    auto itChannel = m_channelIds.begin();
    for (; itChannel != m_channelIds.end();)
    {
        ChannelIdStreamConfig *ptr_delete = itChannel->second;
        if (ptr_delete != NULL)
        {
            delete ptr_delete;
        }

        m_channelIds.erase(itChannel++);
    }
}

void ChannelIdManager::InvalidDataCleanTimerEvent(int fd, short which, void *v)
{
    ChannelIdManager *ptr_this = (ChannelIdManager *)v;

    ptr_this->CleanInvalidChannelId();
    ptr_this->CleanChannelIdInvalidPtr();
}

void ChannelIdManager::FreeChannelIdPtr(ChannelIDConfig *ptr)
{
    if (ptr == NULL)
    {
        return;
    }
    ptr->SetNeedDelete();
    std::string str_key = ptr->m_strChannelId + "|" + gse::tools::strings::ToString((uint64_t)ptr);
    m_freeChannelidLock.WLock();
    std::map<std::string, ChannelIDConfig *>::iterator it = m_needFreeChannelIDconfig.find(str_key);
    if (it == m_needFreeChannelIDconfig.end())
    {
        m_needFreeChannelIDconfig.insert(std::make_pair(str_key, ptr));
        LOG_DEBUG("Will free channelid(%s) mem ptr:%p, key:%s", ptr->m_strChannelId.c_str(), ptr, str_key.c_str());
    }
    m_freeChannelidLock.UnLock();
}
void ChannelIdManager::CleanChannelIdInvalidPtr()
{
    LOG_DEBUG("free channelid mem, map size:%d", m_needFreeChannelIDconfig.size());
    m_freeChannelidLock.WLock();
    std::map<std::string, ChannelIDConfig *>::iterator it = m_needFreeChannelIDconfig.begin();
    for (; it != m_needFreeChannelIDconfig.end();)
    {
        ChannelIDConfig *ptr_delete = it->second;
        if (ptr_delete != NULL)
        {
            LOG_DEBUG("free channelid mem, key:%s, %p", it->first.c_str(), it->second);
            if (ptr_delete->IsNeedDelete())
            {
                LOG_DEBUG("free channelid config, channlied(%s), request count(%d), reponse count(%d), create time(%d), now(%d)", ptr_delete->m_strChannelId.c_str(),
                          ptr_delete->m_zkReqResponseCount.m_requestCount,
                          ptr_delete->m_zkReqResponseCount.m_responseCount, ptr_delete->m_setDeleteTimestamp, gse::tools::time::GetUTCSecond());
                m_needFreeChannelIDconfig.erase(it++);
                delete ptr_delete;
            }
            else
            {
                ++it;
            }
        }
        else
        {
            m_needFreeChannelIDconfig.erase(it++);
        }
    }
    m_freeChannelidLock.UnLock();
}
void ChannelIdManager::CleanInvalidChannelId()
{
    std::list<ChannelIdStreamConfig *>::iterator it = m_needDeleteChannelIds.begin();
    for (; it != m_needDeleteChannelIds.end();)
    {
        ChannelIdStreamConfig *ptr_delete = (*it);
        if (ptr_delete != NULL)
        {
            if (ptr_delete->IsNeedDelete())
            {
                std::string table_name;
                ptr_delete->GetTableName(table_name);
                LOG_DEBUG("delete invalid channelid[%d][%s], set delete time(%d), now(%d), free addr(%p)", ptr_delete->m_channelID, table_name.c_str(),
                          ptr_delete->m_setDeleteTimestamp, gse::tools::time::GetUTCSecond(), ptr_delete);
                delete ptr_delete;
                it = m_needDeleteChannelIds.erase(it);
            }
            else
            {
                ++it;
            }
        }
        else
        {
            it = m_needDeleteChannelIds.erase(it);
            LOG_ERROR("Scanning the channelid config's list found invalid ptr, list size:%d", m_needDeleteChannelIds.size());
        }
    }
}

void ChannelIdManager::ChannelIdUpdateEvent(int fd, short which, void *v)
{
    ChannelIdManager *ptr_this = (ChannelIdManager *)v;

    ZkEvent *event = NULL;
    ChannelIDConfig *ptr_channel_id_cfg;
    if (gse::tools::pipe::ReadPipe(fd, &event, sizeof(event)) < 0)
    {
        int tmp_errno = gse_socket_errno;
        string tmp_errormsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to read pipe fd (%d), errno(%d), error message (%s)", fd, tmp_errno, SAFE_CSTR(tmp_errormsg.c_str()));
        return;
    }

    if (event == NULL)
    {
        LOG_ERROR("event msg invalid(null)");
        return;
    }

    switch (event->m_eventType)
    {
    case ZK_EVENT_LIST:
        break;
    case ZK_EVENT_DELETE:
    {
        uint32_t *ptr_channelid = (uint32_t *)event->m_msg;
        ptr_this->HandleChannelIdDelete(*ptr_channelid);
        delete ptr_channelid;
        break;
    }
    case ZK_EVENT_CHANGE:
    case ZK_EVENT_CREATE:
    {
        ptr_channel_id_cfg = (ChannelIDConfig *)event->m_msg;
        ptr_this->HandleChannelIdUpdate(ptr_channel_id_cfg);
        delete ptr_channel_id_cfg;
        break;
    }
    default:
        LOG_ERROR("recv unkown event:%d", event->m_eventType);
    }

    delete event;
    return;
}

void ChannelIdManager::FreeEvent()
{
    if (m_cleanTimerEvent != NULL)
    {
        event_free(m_cleanTimerEvent);
        m_cleanTimerEvent = NULL;
    }

    if (m_channelIdUpdateEvent != NULL)
    {
        event_free(m_channelIdUpdateEvent);
        m_channelIdUpdateEvent = NULL;
    }

    for (int i = 0; i < 2; i++)
    {
        if (m_channelIdUpudateNotifyFd[i] != -1)
        {
            gse::tools::pipe::ClosePipe(m_channelIdUpudateNotifyFd[i]);
            m_channelIdUpudateNotifyFd[i] = -1;
        }
    }

    if (m_evBase != NULL)
    {
        if (m_stoped)
        {
            event_base_free(m_evBase);
            m_evBase = NULL;
        }
        else
        {
            event_base_loopbreak(m_evBase);
            LOG_DEBUG("stop event base loop");
        }
    }
}
void ChannelIdManager::Stop()
{
    LOG_DEBUG("channelid manager stopping, channelId:%d", m_isPlatId);
    m_stoped = true;
    FreeEvent();
}

int ChannelIdManager::ThreadFun()
{
    int ret = -1;
    gse::tools::pipe::CreatePipe(m_channelIdUpudateNotifyFd);
    m_evBase = event_base_new();
    if (NULL == m_evBase)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new event base object for channel id manager thread, errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    m_cleanTimerEvent = event_new(m_evBase, -1, EV_PERSIST, ChannelIdManager::InvalidDataCleanTimerEvent, this);
    if (NULL == m_cleanTimerEvent)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new clean timer event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    ret = event_add(m_cleanTimerEvent, &m_cleanTime);
    if (ret < 0)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to start clean invalid channel id timer event. ret:%d, errno:%d, errmsg:%s", ret, tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    m_channelIdUpdateEvent = event_new(m_evBase, m_channelIdUpudateNotifyFd[0], EV_PERSIST | EV_READ, ChannelIdManager::ChannelIdUpdateEvent, this);
    if (NULL == m_channelIdUpdateEvent)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new channel id update event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    ret = event_add(m_channelIdUpdateEvent, NULL);
    if (ret < 0)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to start channel id update event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    ret = event_base_loop(m_evBase, 0);
    int tmp_errno = gse_errno;
    std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
    LOG_ERROR("channel id manager thread exit. ret:%d, errno:%d, errmsg:%s", ret, tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
    return GSE_SUCCESS;
}

int ChannelIdManager::Update(ZkEvent *event)
{
    int ret = gse::tools::pipe::WritePipe(m_channelIdUpudateNotifyFd[1], &event, sizeof(event));
    return ((ret < 0) ? GSE_ERROR : GSE_SUCCESS);
}

int ChannelIdManager::HandleChannelIdDelete(uint32_t channelid)
{
    m_rwLock.WLock();
    auto it = m_channelIds.find(channelid);
    if (it != m_channelIds.end())
    {
        ChannelIdStreamConfig *ptr_channelid_config = it->second;
        if (ptr_channelid_config != NULL)
        {
            LOG_DEBUG("delete channelid(%d), channelid  count(%d)", channelid, m_channelIds.size());
            OPMetric::AddChannelIdMetrics("channelid", -1);
            OPMetric::AddChannelIdUsageMemoryMetrics("channelid", -1 * (it->second->GetUseMemSize()));
            ptr_channelid_config->SetNeedDelete();
            m_needDeleteChannelIds.push_back(ptr_channelid_config);
        }
        m_channelIds.erase(it);
    }
    m_rwLock.UnLock();
    return GSE_SUCCESS;
}

bool ChannelIdManager::Find(uint32_t channelId)
{
    m_rwLock.RLock();
    auto it = m_channelIds.find(channelId);
    if (it != m_channelIds.end())
    {
        m_rwLock.UnLock();
        return true;
    }
    m_rwLock.UnLock();
    return false;
}

int ChannelIdManager::HandleChannelIdUpdate(ChannelIDConfig *ptr_channel_id_config)
{
    ChannelIdStreamConfig *ptr_channel_id_stream_cfg = NULL;
    if (ptr_channel_id_config->m_ptrMetadata == NULL)
    {
        LOG_ERROR("channel id config invalid");
        return GSE_ERROR;
    }
    uint32_t m_channelID = ptr_channel_id_config->m_ptrMetadata->m_channelID;
    std::size_t max_count = ptr_channel_id_config->m_channels.size();
    if (max_count <= 0)
    {
        LOG_ERROR("channel id config invalid, config json(%s)", ptr_channel_id_config->m_originData.c_str());
        return GSE_ERROR;
    }

    LOG_DEBUG("channelid :%d, channel count:%d, start create correlation channelid config", ptr_channel_id_config->m_ptrMetadata->m_channelID, max_count);
    for (std::size_t stream_to_idx = 0; stream_to_idx < max_count; ++stream_to_idx)
    {
        ChannelIdStreamConfig *ptrChannelidConfigTmp = new ChannelIdStreamConfig();
        Channel *ptrChannel = ptr_channel_id_config->m_channels.at(stream_to_idx);
        StreamTo *ptrStreamTo = &ptrChannel->m_streamTo;
        int stream_to_id = ptrStreamTo->m_streamToId;
        ptrChannelidConfigTmp->m_channelName = ptrChannel->m_name;
        ptrChannelidConfigTmp->m_streamToId = ptrStreamTo->m_streamToId;

        for (std::size_t idx = 0; idx < ptrChannel->m_filterAndName.size(); ++idx)
        {
            std::string filter_name = ptrChannel->m_filterAndName.at(idx);
            StreamFilter *ptr_stream_filter = ptr_channel_id_config->GetFilter(filter_name);
            if (NULL == ptr_stream_filter)
            {
                LOG_WARN("the filter name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(filter_name.c_str()), ptr_channel_id_config->m_ptrMetadata->m_channelID);
                continue;
            }
            LOG_DEBUG("channelid(%d) add Andfilter(%s)", ptr_channel_id_config->m_ptrMetadata->m_channelID, ptr_stream_filter->m_name.c_str());
            ptrChannelidConfigTmp->AddAndFilters(ptr_stream_filter);
        }

        for (std::size_t idx = 0; idx < ptrChannel->m_filterOrName.size(); ++idx)
        {
            std::string filter_name = ptrChannel->m_filterOrName.at(idx);
            StreamFilter *ptr_stream_filter = ptr_channel_id_config->GetFilter(filter_name);
            if (NULL == ptr_stream_filter)
            {
                LOG_WARN("the filter name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(filter_name.c_str()), ptr_channel_id_config->m_ptrMetadata->m_channelID);
                continue;
            }
            LOG_DEBUG("channelid(%d) add or filter(%s)", ptr_channel_id_config->m_ptrMetadata->m_channelID, ptr_stream_filter->m_name.c_str());
            ptrChannelidConfigTmp->AddOrFilters(ptr_stream_filter);
        }

        ptrChannelidConfigTmp->CopyTableConfig(ptrStreamTo);
        ptrChannelidConfigTmp->m_channelID = m_channelID;
        std::string table_name;
        if (ptr_channel_id_stream_cfg == NULL)
        {
            ptr_channel_id_stream_cfg = ptrChannelidConfigTmp;
            ptr_channel_id_stream_cfg->GetTableName(table_name);
            LOG_DEBUG("Create channelId first stream config(id:%d, table name:%s, exporter id:%d, ptr:%p)",
                      ptr_channel_id_stream_cfg->m_channelID, table_name.c_str(), ptr_channel_id_stream_cfg->m_streamToId, ptrChannelidConfigTmp);
        }
        else
        {
            ptr_channel_id_stream_cfg->SetNext(ptrChannelidConfigTmp);
            ptrChannelidConfigTmp->GetTableName(table_name);
            LOG_DEBUG("Create ChannelId next stream config(id:%d, table name:%s, exporter id:%d, ptr:%p)",
                      m_channelID, table_name.c_str(), ptr_channel_id_stream_cfg->m_streamToId, ptrChannelidConfigTmp);
        }
    }

    m_rwLock.WLock();

    auto it = m_channelIds.find(ptr_channel_id_config->m_ptrMetadata->m_channelID);
    int count = 1;
    if (it != m_channelIds.end())
    {
        ChannelIdStreamConfig *ptr_channel_id_delete = NULL;
        count -= 1;
        LOG_DEBUG("delete old exist channelid:%d", ptr_channel_id_config->m_ptrMetadata->m_channelID);
        ptr_channel_id_delete = it->second;
        if (ptr_channel_id_delete != NULL)
        {
            ptr_channel_id_delete->SetNeedDelete();
            m_needDeleteChannelIds.push_back(ptr_channel_id_delete);
        }
        m_channelIds.erase(it);
    }

    m_channelIds.insert(std::make_pair(m_channelID, ptr_channel_id_stream_cfg));
    m_rwLock.UnLock();
    OPMetric::AddChannelIdMetrics("channelid", count);
    OPMetric::AddChannelIdUsageMemoryMetrics("channelid", count * (ptr_channel_id_stream_cfg->GetUseMemSize()));
    return GSE_SUCCESS;
}

void ChannelIdStreamConfig::Dump(string &strDump)
{
    Json::Value root;
    Json::Value channels;
    Json::Value filters;

    root["channelid"] = Json::Value(m_channelID);

    ChannelIdStreamConfig *pNext = this;
    do
    {
        Json::Value channel;
        channel["streamid"] = Json::Value(pNext->m_streamToId);
        channel["name"] = Json::Value(pNext->m_channelName);

        switch (m_storageType)
        {
        case KAFKA_TOPIC:
            if (pNext->m_storageTable.m_kafkaStorageTable != nullptr)
            {
                channel["kafka"] = pNext->m_storageTable.m_kafkaStorageTable->ToJSONValue();
            }
            break;
        case REDIS_CHANNEL:
            if (pNext->m_storageTable.m_redisStorageTable != nullptr)
            {
                channel["redis"] = pNext->m_storageTable.m_redisStorageTable->ToJSONValue();
            }
            break;
        case PULSAR_TOPIC:
            if (pNext->m_storageTable.m_pulsarStorageTable != nullptr)
            {
                channel["pulsar"] = pNext->m_storageTable.m_pulsarStorageTable->ToJSONValue();
            }
            break;
        case DSPROXY:
            break;

        default:
            LOG_WARN("unkown table type:%d", m_storageType);
        }

        Json::Value andFilter;

        for (auto filter : pNext->m_filter.m_streamFilterAnd)
        {
            andFilter.append(Json::Value(filter->m_name));
        }

        if (andFilter.size() > 0)
        {
            channel["filter_name_and"] = andFilter;
        }

        Json::Value orFilter;
        for (auto filter : pNext->m_filter.m_streamFilterOr)
        {
            orFilter.append(Json::Value(filter->m_name));
        }

        if (orFilter.size() > 0)
        {
            channel["filter_name_or"] = orFilter;
        }

        channels.append(channel);

        // filter

        for (auto filter : pNext->m_filter.m_streamFilterAnd)
        {
            filters.append(filter->ToJSONValue());
        }

        pNext = pNext->m_next;
    } while (pNext != nullptr);

    root["channels"] = channels;
    if (filters.size() > 0)
    {
        root["filters"] = filters;
    }

    strDump = root.toStyledString();
    return;
}

} // namespace data
} // namespace gse
