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

#include "datachannel.h"

#include "receiver/httpreceiver.h"
//##include "receiver/kcpreceiver.h"
#include "receiver/udpreceiver.h"
#include "receiver/tcpreceiver.h"
#include "exporter/log/logexporter.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"
////#include "exporter/dsproxy/proxy_exporter.h"
#include "exporter/dsproxy/proxy_exporter.h"
#include "exporter/kafka/kafkaexporter.h"
#include "exporter/redis/redisexporter.h"
#include "exporter/pulsar/pulsarexporter.h"
#include "protocol.h"
//#include "codec/tglog_pb/tglog_pb_codec.h"
//#include "codec/tdm_pkg/tdm_pkg_codec.h"
//#include "codec/tdm_pb/tdm_pb_codec.h"
//#include "codec/gsedata_pkg/gsedata_pkg_codec.h"
//#include "codec/gsedata_v1/gsedata_pkg_codec_v1.h"
//#include "codec/gsedata_tglog_v1/gsedata_tglog_pkg_codec_v1.h"
//#include "codec/transport/transport.h"
//#include "codec/gsedata_pkg_v2/gsedata_pkg_codec_v2.h"
#include "time_center.h"
#include "dataserver.h"

namespace gse { 
namespace dataserver {

#define MAX_WORKER_NUM 64

DataChannel::DataChannel()
{
    m_channelConf = NULL;
    m_receiver = NULL;
    m_codec = NULL;
    m_nextWorkerId = 0;
    m_numWorker = 6;
    m_isKeepTiming = false;
    m_ptrEventThread = NULL;
    m_opsReport = NULL;
    m_conftor = NULL;
    m_ptrEventThread = NULL;
    m_ptrIDToStorage = NULL;
}

DataChannel::~DataChannel()
{
    Stop();

    if (m_channelConf != NULL)
    {
        delete m_channelConf;
        m_channelConf = NULL;
    }

    if (m_codec != NULL)
    {
        delete m_codec;
        m_codec = NULL;
    }
}

void DataChannel::cleanExporter(evutil_socket_t fd, short what, void *args)
{

    DataChannel *ptr_this = reinterpret_cast<DataChannel *>(args);

    int max_count = ptr_this->m_toDeletedExporter.Size();
    LOG_DEBUG("datachannel[%s] clean the invalid exporters by timer, invalid exporter size:%d", SAFE_CSTR(ptr_this->m_channelName.c_str()), max_count);

    while (--max_count >= 0)
    {
        Exporter *ptr_exporter = NULL;
        ptr_this->m_toDeletedExporter.Pop(ptr_exporter);
        if (NULL == ptr_exporter)
        {
            break;
        }

        if (!ptr_exporter->IsNeedStop())
        {
            LOG_DEBUG("push expoter(%p:%s) to delete list", ptr_exporter, ptr_exporter->GetExporterName().c_str());
            ptr_this->m_toDeletedExporter.Push(ptr_exporter);
            continue;
        }

        std::string name = ptr_exporter->GetExporterName();
        LOG_DEBUG("will stop the invalid exporter (%s) by timer", SAFE_CSTR(name.c_str()));
        ptr_exporter->Stop();
        LOG_DEBUG("real delete exporer %p, name:%s", ptr_exporter, ptr_exporter->GetExporterName().c_str());
        delete ptr_exporter;

        LOG_DEBUG("finish stop the exporter (%s) by timer", SAFE_CSTR(name.c_str()));
    }
}

int DataChannel::Start()
{
    if (NULL == m_channelConf)
    {
        LOG_WARN("the pointer of channel configure is NULL, can't start channel, please check channel configure");
        return GSE_SYSTEMERROR;
    }

    Stop();

    // init
    if (GSE_SUCCESS != init())
    {
        LOG_WARN("fail to init datachannel[%s]", SAFE_CSTR(m_channelName.c_str()));
    }

    // channel worker
    if (GSE_SUCCESS != startChannelWorker())
    {
        LOG_WARN("fail to start channelworker for datachannel[%s]", SAFE_CSTR(m_channelName.c_str()));
        Stop();
        return GSE_ERROR;
    }

    // receiver
    if (GSE_SUCCESS != startReceiver())
    {
        LOG_WARN("fail to start receiver[%s] for datachannel[%s]", SAFE_CSTR(m_receiverName.c_str()), SAFE_CSTR(m_channelName.c_str()));
        Stop();
        return GSE_ERROR;
    }
    LOG_INFO("success to start receiver[%s] in channel[%s]", SAFE_CSTR(m_receiverName.c_str()), SAFE_CSTR(m_channelName.c_str()));

    // exporter
    createExporter();

    if (GSE_SUCCESS != createDataflowExporter())
    {
        LOG_WARN("fail to start exporter for datachannel[%s]", SAFE_CSTR(m_channelName.c_str()));
        Stop();
        return GSE_ERROR;
    }

    LOG_DEBUG("datachannel[%s] has start. the size of channelworker is [%d]", SAFE_CSTR(m_channelName.c_str()), m_vChannelWorkers.size());

    m_ptrEventThread = new rgse::GseEventThread();
    int registerResponseVal = m_ptrEventThread->registerTimerPersistEvent(DataChannel::cleanExporter, this, 60); // 1 minutes
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return registerResponseVal;
    }

    int startEventThreadVal = m_ptrEventThread->Start();
    if (GSE_SUCCESS > startEventThreadVal)
    {
        LOG_ERROR("failed to start event thread, error code %d", startEventThreadVal);
        return startEventThreadVal;
    }

    return GSE_SUCCESS;
}

int DataChannel::Stop()
{
    if (NULL != m_ptrEventThread)
    {
        m_ptrEventThread->stop();
        m_ptrEventThread = NULL;
    }

    // stop receiver
    stopReceiver();

    // stop channel worker
    stopChannelWorker();

    // stop exporter
    clearDataflowExporter();

    return GSE_SUCCESS;
}

void DataChannel::Join()
{
    ChannelWorker *pWorker = NULL;
    for (int i = 0; i < m_vChannelWorkers.size(); i++)
    {
        pWorker = m_vChannelWorkers[i];
        if (NULL == pWorker)
        {
            continue;
        }

        pWorker->Join();
    }

    if (m_receiver != NULL)
    {
        m_receiver->Join();
    }
}

void DataChannel::SetOps(OpsCollection* ptr_ops_report)
{
    m_opsReport = ptr_ops_report;
}

void DataChannel::UpdateConf(const ChannelConf &channelConf)
{
    m_upConfLock.WLock();
    if (m_channelConf != NULL)
    {
        delete m_channelConf;
    }

    m_channelConf = new ChannelConf(channelConf);
    m_upConfLock.UnLock();
}

void DataChannel::UpdateIDToStorage(IDToStorage *ptrIDToStorage)
{
    m_ptrIDToStorage = ptrIDToStorage;
    m_ptrIDToStorage->WatchUpdateEvent(DataChannel::handleConfigUpdateEvent, this);
}


void DataChannel::stopExporter(uint32_t channelid)
{
    Exporter *ptr_todelete_exporter = NULL;
    bool bret = m_channelIDToExporter.Find(channelid, ptr_todelete_exporter);
    if (NULL != ptr_todelete_exporter)
    {
        m_channelIDToExporter.Pop(channelid, ptr_todelete_exporter);
        LOG_DEBUG("stop exporter(%s:%p) for the channelid (%d)",  ptr_todelete_exporter->GetExporterName().c_str(), ptr_todelete_exporter, channelid);
        ptr_todelete_exporter->SetNeedStop();
        LOG_DEBUG("push channel_id(%d) expoter(%p:%s) to delete list", channelid, ptr_todelete_exporter, ptr_todelete_exporter->GetExporterName().c_str());
        m_toDeletedExporter.Push(ptr_todelete_exporter);
    }
}

void DataChannel::handleConfigUpdateEvent(void *args, int storageIndex, uint32_t channelID)
{
    DataChannel *ptr_this = reinterpret_cast<DataChannel *>(args);
    // ptr_this->createExporter();
    LOG_DEBUG("update event callback the storage index %d, the dataid %d", storageIndex, channelID);
    DataStorage *ptr_storage = NULL;
    if (storageIndex >= 0)
    {
        ptr_storage = ptr_this->m_ptrIDToStorage->GetStorageByIndex(storageIndex);
        Exporter *ptr_todelete_exporter = NULL;
        bool bret = ptr_this->m_storageIndexToExporter.Pop(storageIndex, ptr_todelete_exporter);
        if (NULL != ptr_todelete_exporter)
        {
            ptr_todelete_exporter->SetNeedStop();
            LOG_DEBUG("push expoter(%p:%s) to delete list, storage index(%d)", ptr_todelete_exporter, ptr_todelete_exporter->GetExporterName().c_str(), storageIndex);
            ptr_this->m_toDeletedExporter.Push(ptr_todelete_exporter);
        }
    }

    if (!is_dataid(channelID))
    {
        ptr_storage = ptr_this->m_ptrIDToStorage->GetStorageByChannelID(channelID);
    }

    if (NULL == ptr_storage)
    {
        return;
    }

    if (NULL == ptr_storage->m_storage.m_ptrChannelIDStorage)
    {
        if (!is_dataid(channelID))
        {
            LOG_WARN("channel[%d] config not exist, stop exporter", channelID);

            ptr_this->stopExporter(channelID);
        }

        return;
    }

    Exporter *ptr_exporter = NULL;
    do
    {
        Exporter *ptr_exporter_tmp = NULL;
        LOG_DEBUG("storage is data:%d, type:%d, ptr:%p", ptr_storage->m_isDataID, ptr_storage->m_storageType, ptr_storage);
        if (ptr_storage->m_isDataID)
        {
            ptr_this->updateDataIDExporter(ptr_storage->m_storage.m_ptrDataIDConfig);
        }
        else
        {
            ptr_this->updateChannelIDExporter(ptr_storage->m_storage.m_ptrChannelIDStorage);
        }

        ptr_storage = ptr_storage->m_next;

    } while (NULL != ptr_storage);

    delete ptr_storage;
}

int DataChannel::DecodeMsg(DataCell *pDataCell)
{
    if (m_codec != NULL)
    {
        int iRet = m_codec->DecodeMsg(pDataCell);
        if (iRet != GSE_SUCCESS)
        {
            LOG_ERROR("Message decoding failed, channeid:%d, client ip:%s, port:%d", 
            pDataCell->GetChannelID(), pDataCell->GetSourceIp().c_str(), pDataCell->GetSourcePort());
            pDataCell->SetErrorMsg("decode failed", OPS_ERROR_DECODE);
            return iRet;
        }
    }

    return GSE_SUCCESS;
}

void DataChannel::dataflowExporterWrite(DataCell *pDataCell)
{
    for (int i = 0; i < m_vDataflowExporter.size(); i++)
    {
        Exporter *pExporter = m_vDataflowExporter[i];
        if (pExporter != NULL)
        {
            if (GSE_SUCCESS != pExporter->Write(pDataCell))
            {
                LOG_WARN("fail to write data into exporter[%s]", SAFE_CSTR(pExporter->GetExporterName().c_str()));
            }
        }
        //LOG_DEBUG("send data for the channel id (%d) by the exporter(%s)", pDataCell->GetChannelID(), SAFE_CSTR(pExporter->GetExporterName().c_str()));
    }
}

void DataChannel::updateDataIDExporter(StorageConfigVector *ptrStorage)
{
    if (NULL == ptrStorage)
    {
        LOG_ERROR("the dataid's storage config is not seted");
        return;
    }

    Exporter *ptr_exporter = NULL;
    std::size_t max_count = ptrStorage->size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        Exporter *ptr_exporter_tmp = NULL;
        StorageConfigType config = ptrStorage->at(idx);
        StorageConfigType *ptr_config = new StorageConfigType(config);

        switch (ptr_config->m_storageType)
        {
        case KAFKA_COMMON:
            ptr_exporter_tmp = new KafkaExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::DataIDConfType, ptr_config);
            LOG_INFO("create a new kafka exporter for the storage %s", SAFE_CSTR(ptr_config->m_host.c_str()));
            break;
        case KAFKA_OP:
            LOG_ERROR("create kafka ops, now paas");
            continue;
        case REDIS_PUB:
            LOG_ERROR("create redis pub client");
        case REDIS_SENTINEL_PUB:
            LOG_ERROR("create redis sentinel client");
            ptr_exporter_tmp = new RedisExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::DataIDConfType, ptr_config);
            break;
        case EXPORT_DSPROXY:
            LOG_ERROR("in the proxy mode, it is not be set exporter for dataid config");
            continue;
        case EXPORT_PULSAR:
            ptr_exporter_tmp = new PulsarExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::DataIDConfType, ptr_config);
            LOG_INFO("create a new pular exporter for the storage %s", SAFE_CSTR(ptr_config->m_host.c_str()));
            break;
        default:
            LOG_ERROR("it is failed to create a exporter for the unkown storage(index %d type %d)", ptr_config->m_clusterIndex, ptr_config->m_storageType);
            continue;
        }

        LOG_DEBUG("new exporter(%p:%s), type:%d", ptr_exporter_tmp, ptr_exporter_tmp->GetExporterName().c_str(), ptr_config->m_storageType);
        if (GSE_SUCCESS != ptr_exporter_tmp->Start())
        {
            LOG_ERROR("it is failed to start the exporter (%s) ", SAFE_CSTR(ptr_exporter_tmp->GetExporterName().c_str()));
            delete ptr_exporter_tmp;
            continue;
        }

        LOG_DEBUG("save the storage by cluster index (%d) (%s)", ptr_config->m_clusterIndex, SAFE_CSTR(ptr_config->m_host.c_str()));
        // 仅在此处插入新的exporter
        Exporter *ptr_todelete_exporter = NULL;
        bool bret = m_storageIndexToExporter.Pop(ptr_config->m_clusterIndex, ptr_todelete_exporter);
        if (NULL != ptr_todelete_exporter)
        {
            ptr_todelete_exporter->SetNeedStop();
            LOG_DEBUG("push expoter(%p:%s) to delete list, update data id", ptr_todelete_exporter, ptr_todelete_exporter->GetExporterName().c_str());
            m_toDeletedExporter.Push(ptr_todelete_exporter);
        }
        m_storageIndexToExporter.Push(ptr_config->m_clusterIndex, ptr_exporter_tmp);
    }

    LOG_DEBUG("valid exporter size:%d invalid exporter size:%d", m_storageIndexToExporter.Size(), m_toDeletedExporter.Size());
}

Exporter *DataChannel::createDataIDExporter(StorageConfigVector *ptrStorage)
{
    if (NULL == ptrStorage)
    {
        LOG_ERROR("the dataid's storage config is not seted");
        return NULL;
    }

    Exporter *ptr_exporter = NULL;
    std::size_t max_count = ptrStorage->size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        Exporter *ptr_exporter_tmp = NULL;
        StorageConfigType config = ptrStorage->at(idx);
        StorageConfigType *ptr_config = new StorageConfigType(config);

        Exporter *ptr_exists_exporter = NULL;
        bool bret = m_storageIndexToExporter.Find(ptr_config->m_clusterIndex, ptr_exists_exporter);
        if (NULL != ptr_exists_exporter)
        {
            LOG_WARN("the storage index (%d) is already exists, will pass this config item, please to check the config", ptr_config->m_clusterIndex);
            continue;
        }

        switch (ptr_config->m_storageType)
        {
        case KAFKA_COMMON:
            ptr_exporter_tmp = new KafkaExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::DataIDConfType, ptr_config);
            LOG_INFO("create a new kafka exporter for the storage %s", SAFE_CSTR(ptr_config->m_host.c_str()));
            break;
        case KAFKA_OP:
            LOG_ERROR("create kafka ops, now paas");
            continue;
        case REDIS_PUB:
            LOG_ERROR("create redis pub client");
        case REDIS_SENTINEL_PUB:
            LOG_ERROR("create redis sentinel client");
            ptr_exporter_tmp = new RedisExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::DataIDConfType, ptr_config);
            break;
        case EXPORT_DSPROXY:
            LOG_ERROR("in the proxy mode, it is not be set exporter for dataid config");
            continue;
        case EXPORT_PULSAR:
            LOG_INFO("create a new pulsar exporter for the storage %s", SAFE_CSTR(ptr_config->m_host.c_str()));
            ptr_exporter_tmp = new PulsarExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::DataIDConfType, ptr_config);
            break;            
        default:
            LOG_ERROR("it is failed to create a exporter for the unkown storage(index %d type %d)", ptr_config->m_clusterIndex, ptr_config->m_storageType);
            continue;
        }

        if (GSE_SUCCESS != ptr_exporter_tmp->Start())
        {
            LOG_ERROR("it is failed to start the exporter (%s) ", SAFE_CSTR(ptr_exporter_tmp->GetExporterName().c_str()));
            delete ptr_exporter_tmp;
            continue;
        }

        LOG_DEBUG("save the storage by cluster index (%d) (%s)", config.m_clusterIndex, SAFE_CSTR(config.m_host.c_str()));
        // 仅在此处插入新的exporter
        m_storageIndexToExporter.Push(config.m_clusterIndex, ptr_exporter_tmp);
    }
    return ptr_exporter;
}

void DataChannel::updateChannelIDExporter(ChannelIDStorage *ptrStorage)
{
    if (NULL == ptrStorage)
    {
        LOG_ERROR("the channelid's storage config is not seted");
        return;
    }

    // ChannelIDStorage 一个 channelid 只有一个 storage 配置

    LOG_DEBUG("update storage exporter for the storage %s ", ptrStorage->m_ptrConfigInfo->m_name.c_str());

    Exporter *ptr_exporter = NULL;

    do
    {
        if (NULL != ptr_exporter)
        {
            if (ptr_exporter->IsExistsAndUpdateFilter(ptrStorage->m_ptrConfigInfo->m_name, ptrStorage->m_ptrChannelIDFilters))
            {
                LOG_WARN("the storage (%s) already association with a exporter", ptrStorage->m_ptrConfigInfo->m_name.c_str());
                ptrStorage = ptrStorage->m_next;
                continue;
            }
        }
        Exporter *ptr_exporter_tmp = NULL;

        switch (ptrStorage->m_storageType)
        {
        case KAFKA_COMMON:
            ptr_exporter_tmp = new KafkaExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            LOG_INFO("create a new kafka exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            break;
        case EXPORT_FILE:
            ptr_exporter_tmp = new LogExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            LOG_INFO("create a new log exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            break;
        case EXPORT_DSPROXY:
            LOG_INFO("create a ds proxy exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            ptr_exporter_tmp = new ProxyExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            break;
        case EXPORT_PULSAR:
            ptr_exporter_tmp = new PulsarExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            LOG_INFO("create a new pulsar exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            break;

        default:
            LOG_ERROR("it is not supported storage type (%d), please the check the config(only suport 'kafka' and 'file' for the config (%s)", (int)(ptrStorage->m_storageType), SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            ptrStorage = ptrStorage->m_next;
            continue;
        }

        if (GSE_SUCCESS != ptr_exporter_tmp->Start())
        {
            LOG_ERROR("it is failed to start the exporter (%s) ", SAFE_CSTR(ptr_exporter_tmp->GetExporterName().c_str()));
            delete ptr_exporter_tmp;
            ptrStorage = ptrStorage->m_next;
            continue;
        }

        ptr_exporter_tmp->AddFilter(ptrStorage->m_ptrChannelIDFilters);
        if (NULL != ptr_exporter)
        {
            ptr_exporter->SetNext(ptr_exporter_tmp);
        }
        else
        {
            ptr_exporter = ptr_exporter_tmp;
            // 一个channelid 关联一组 数据流向，一个数据流向是一个exporter，
            // 此处只需要记录第一个，其余的exporter 可通过exporter 职责链结构递归使用
            LOG_DEBUG("set exporter for the channelid (%d)", ptrStorage->m_channelID);
            Exporter *ptr_todelete_exporter = NULL;
            bool bret = m_channelIDToExporter.Pop(ptrStorage->m_channelID, ptr_todelete_exporter);
            if (NULL != ptr_todelete_exporter)
            {
                ptr_todelete_exporter->SetNeedStop();
                m_toDeletedExporter.Push(ptr_todelete_exporter);
            }
            m_channelIDToExporter.Push(ptrStorage->m_channelID, ptr_exporter);
        }
        ptrStorage = ptrStorage->m_next;
    } while (NULL != ptrStorage);

    LOG_DEBUG("valid channelid exporter size:%d invalid exporter size:%d", m_channelIDToExporter.Size(), m_toDeletedExporter.Size());
    return;
}


Exporter *DataChannel::createChannelIDExporter(ChannelIDStorage *ptrStorage)
{
    if (NULL == ptrStorage)
    {
        LOG_ERROR("create channelid's exporter, the channelid's storage config is not seted");
        return NULL;
    }

    // ChannelIDStorage 一个 channelid 只有一个 storage 配置

    LOG_DEBUG("create exporter for the storage %s ", ptrStorage->m_ptrConfigInfo->m_name.c_str());

    Exporter *ptr_exporter = NULL;

    do
    {

        if (NULL != ptr_exporter)
        {
            if (ptr_exporter->IsExistsAndUpdateFilter(ptrStorage->m_ptrConfigInfo->m_name, ptrStorage->m_ptrChannelIDFilters))
            {
                LOG_WARN("the storage (%s) already association with a exporter", ptrStorage->m_ptrConfigInfo->m_name.c_str());
                ptrStorage = ptrStorage->m_next;
                continue;
            }
        }
        Exporter *ptr_exporter_tmp = NULL;

        switch (ptrStorage->m_storageType)
        {
        case KAFKA_COMMON:
            ptr_exporter_tmp = new KafkaExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            LOG_INFO("create a new kafka exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            break;
        case EXPORT_FILE:
            ptr_exporter_tmp = new LogExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            LOG_INFO("create a new log exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            break;
        case EXPORT_DSPROXY:
            LOG_INFO("create a ds proxy exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            ptr_exporter_tmp = new ProxyExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            break;
        case EXPORT_PULSAR:
            ptr_exporter_tmp = new PulsarExporter();
            ptr_exporter_tmp->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptrStorage);
            LOG_INFO("create a new pulsar exporter for the storage %s", SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            break;

        default:
            LOG_ERROR("it is not supported storage type (%d), please the check the config(only suport 'kafka' and 'file' for the config (%s)", (int)(ptrStorage->m_storageType), SAFE_CSTR(ptrStorage->m_streamToName.c_str()));
            ptrStorage = ptrStorage->m_next;
            continue;
        }

        if (GSE_SUCCESS != ptr_exporter_tmp->Start())
        {
            LOG_ERROR("it is failed to start the exporter (%s) ", SAFE_CSTR(ptr_exporter_tmp->GetExporterName().c_str()));
            delete ptr_exporter_tmp;
            ptrStorage = ptrStorage->m_next;
            continue;
        }

        ptr_exporter_tmp->AddFilter(ptrStorage->m_ptrChannelIDFilters);
        if (NULL != ptr_exporter)
        {
            ptr_exporter->SetNext(ptr_exporter_tmp);
        }
        else
        {
            ptr_exporter = ptr_exporter_tmp;
            // 一个channelid 关联一组 数据流向，一个数据流向是一个exporter，
            // 此处只需要记录第一个，其余的exporter 可通过exporter 职责链结构递归使用
            LOG_DEBUG("set exporter for the channelid (%d)", ptrStorage->m_channelID);
            m_channelIDToExporter.Push(ptrStorage->m_channelID, ptr_exporter);
        }
        ptrStorage = ptrStorage->m_next;
    } while (NULL != ptrStorage);

    return ptr_exporter;
}

//zk channelid exporter
void DataChannel::createExporter()
{
    DataStorage *ptrDataStorage = m_ptrIDToStorage->GetAllStorages();
    if (NULL == ptrDataStorage)
    {
        LOG_ERROR("it is failed to create the exporter chain, because the storage is not seted");
        return;
    }

    Exporter *ptr_exporter = NULL;
    do
    {
        Exporter *ptr_exporter_tmp = NULL;
        if (ptrDataStorage->m_isDataID)
        {
            createDataIDExporter(ptrDataStorage->m_storage.m_ptrDataIDConfig);
        }
        else
        {
            createChannelIDExporter(ptrDataStorage->m_storage.m_ptrChannelIDStorage);
        }

        ptrDataStorage = ptrDataStorage->m_next;

    } while (NULL != ptrDataStorage);

    delete ptrDataStorage;
}

int DataChannel::ExportData(DataCell *pDataCell)
{
    pDataCell->SetOutputTimestamp(TimeCenter::Instance()->GetDateTime());

    dataflowExporterWrite(pDataCell);

    if (pDataCell->GetChannelID() == 0)
    {
        return GSE_SUCCESS;
    }

    if (m_conftor->GetChannelIdManager() != NULL)
    {
        m_conftor->GetChannelIdManager()->WriteByFilter(pDataCell);
        return GSE_SUCCESS;
    }

    return GSE_SUCCESS;
}

int DataChannel::init()
{
    if (NULL == m_channelConf)
    {
        LOG_WARN("the pointer of channel configure is NULL, can't start channel, please check channel configure");
        return GSE_SYSTEMERROR;
    }

    DecodeTypeEnum decodeType = D_TYPE_UNKNOWN;
    m_upConfLock.RLock();
    decodeType = m_channelConf->m_decodeType;
    m_channelName = m_channelConf->m_name;
    m_receiverName = m_channelConf->m_receiverName;
    m_upConfLock.UnLock();

    m_codec = CodecFactory::CreateCodec(decodeType);

    if (NULL != m_codec)
    {
    }

    return GSE_SUCCESS;
}

int DataChannel::startChannelWorker()
{
    if (NULL == m_channelConf)
    {
        LOG_WARN("the pointer of channel configure is NULL, can't start channel, please check channel configure");
        return GSE_SYSTEMERROR;
    }

    m_upConfLock.RLock();
    m_numWorker = m_channelConf->m_numChannelWorker;
    m_upConfLock.UnLock();
    // start channel worker
    if (m_numWorker <= 0)
    {
        m_numWorker = 6;
    }

    if (m_numWorker > MAX_WORKER_NUM)
    {
        LOG_WARN("the config of channel worker number[%d] is more then MAX_WORKER_NUM[%d]. will set channel woker number to MAX_WORKER_NUM", m_numWorker, MAX_WORKER_NUM);
        m_numWorker = MAX_WORKER_NUM;
    }

    uint32_t workerRealNum = m_numWorker;

    for (uint32_t i = 0; i < m_numWorker; i++)
    {
        ChannelWorker *pWorker = new ChannelWorker(i, this);
        if (GSE_SUCCESS != pWorker->Init())
        {
            LOG_WARN("fail to initialize the channel worker[%d]", i);
            delete pWorker;
            workerRealNum -= 1;
            continue;
        }

        if (GSE_SUCCESS != pWorker->Start())
        {
            LOG_WARN("fail to start the channel worker[%d]", i);
            delete pWorker;
            workerRealNum -= 1;
            continue;
        }

        m_vChannelWorkers.push_back(pWorker);
    }

    if (workerRealNum <= 0)
    {
        LOG_WARN("fail to start channel worker, there is zero worker be created");
        return GSE_ERROR;
    }

    m_numWorker = workerRealNum;

    return GSE_SUCCESS;
}

int DataChannel::stopChannelWorker()
{
    ChannelWorker *pWorker = NULL;
    for (int i = 0; i < m_vChannelWorkers.size(); i++)
    {
        pWorker = m_vChannelWorkers[i];
        if (NULL == pWorker)
        {
            continue;
        }

        pWorker->stop();

        delete pWorker;
    }
    m_vChannelWorkers.clear();

    return GSE_SUCCESS;
}

int DataChannel::startReceiver()
{
    ReceiverProtocolEnum recvProto = R_PROTO_UNKNOWN;
    m_upConfLock.RLock();
    if (m_channelConf->m_receiverConf != NULL)
    {
        recvProto = m_channelConf->m_receiverConf->m_protocol;
    }
    switch (recvProto)
    {
    case R_PROTO_TCP:
    {
        m_receiver = new TcpReceiver();
    }
    break;
    case R_PROTO_UDP:
    {
        m_receiver = new UdpReceiver();
    }
    break;
    case R_PROTO_KCP:
    {
        //m_receiver = new KcpReceiver();
    }
    break;
    case R_PROTO_HTTP:
    {
        m_receiver = new HttpReceiver();
    }
    break;
    default:
        LOG_WARN("unknow the protocol[%d] type of receiver[%s] in channel[%s]", recvProto, SAFE_CSTR(m_receiverName.c_str()), SAFE_CSTR(m_channelName.c_str()));
        break;
    }

    if (m_receiver != NULL)
    {
        m_receiver->SetGseConf(m_gseCfg);
        m_receiver->UpdateConf(*m_channelConf->m_receiverConf);
    }

    m_upConfLock.UnLock();
    LOG_DEBUG("start receiver:%s", m_channelConf->m_receiverConf->m_name.c_str());
    if (NULL == m_receiver)
    {
        return GSE_ERROR;
    }

    m_receiver->SetRecvCallBack(DataChannel::handlRecvDataCell, this);

    return m_receiver->Start();
}

int DataChannel::stopReceiver()
{
    if (m_receiver != NULL)
    {
        m_receiver->Stop();
        delete m_receiver;
    }

    return GSE_SUCCESS;
}

int DataChannel::createDataflowExporter()
{
    m_upConfLock.WLock();
    for (int i = 0; i < m_channelConf->m_vExporterConf.size(); i++)
    {
        ExporterConf *pExporterConf = m_channelConf->m_vExporterConf[i];
        Exporter *pExporter = NULL;
        switch (pExporterConf->m_type)
        {
        case E_TYPE_LOG:
        {
            pExporter = new LogExporter();
            LOG_DEBUG("create log exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        case E_TYPE_BKDATA:
        {
            //
        }
        break;
        case E_TYPE_KAFKA:
        {
            pExporter = new KafkaExporter();
            LOG_DEBUG("create kafka exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        case E_TYPE_REDIS:
        {
            pExporter = new RedisExporter();
            LOG_DEBUG("create redis exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        case E_TYPE_DS_PROXY:
        {
            pExporter = new ProxyExporter();
            LOG_DEBUG("create proxy exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;

        case E_TYPE_PULSAR:
        {
            pExporter = new PulsarExporter();
            LOG_DEBUG("create pulsar exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        default:
            LOG_WARN("unknow the type[%d] of exporter[%s] in channel[%s]", pExporterConf->m_type, SAFE_CSTR(pExporterConf->m_name.c_str()), SAFE_CSTR(m_channelName.c_str()));
            break;
        }

        if (pExporter != NULL)
        {
            pExporter->SetGseConf(m_gseCfg);
            pExporter->UpdateConf(ExporterConfigWrapper::DataFlowConfType, pExporterConf);
            if (GSE_SUCCESS != pExporter->Start())
            {
                LOG_WARN("fail to start exporter[%s] in channel[%s]", SAFE_CSTR(pExporterConf->m_name.c_str()), SAFE_CSTR(m_channelName.c_str()));
                delete pExporter;
                continue;
            }
            else
            {
                m_vDataflowExporter.push_back(pExporter);
                LOG_INFO("success to start exporter[%s] in channel[%s]", SAFE_CSTR(pExporterConf->m_name.c_str()), SAFE_CSTR(m_channelName.c_str()));
            }
        }
        else
        {
            LOG_ERROR("failed to create exporter");
        }
    }
    m_upConfLock.UnLock();

    return GSE_SUCCESS;
}

void DataChannel::clearDataflowExporter()
{
    Exporter *pExporter = NULL;
    for (int i = 0; i < m_vDataflowExporter.size(); i++)
    {
        pExporter = m_vDataflowExporter[i];
        if (NULL == pExporter)
        {
            continue;
        }

        pExporter->Stop();
        delete pExporter;
    }

    m_vDataflowExporter.clear();
}

void DataChannel::handlRecvDataCell(DataCell *pDataCell, void *pCaller)
{
    DataChannel *pDataChannel = (DataChannel *)pCaller;
    pDataChannel->localHanlerRecvDataCell(pDataCell);
}

void DataChannel::localHanlerRecvDataCell(DataCell *pDataCell)
{
    if (m_isKeepTiming)
    {
        directExport(pDataCell);
    }
    else
    {
        dispatchChannelWorker(pDataCell);
    }
}

int DataChannel::directExport(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    iRet = DecodeMsg(pDataCell);
    if (iRet != GSE_SUCCESS)
    {
        delete pDataCell;
        return iRet;
    }

    iRet = ExportData(pDataCell);
    if (iRet != GSE_SUCCESS)
    {
        delete pDataCell;
        return iRet;
    }

    // success
    delete pDataCell;
    return GSE_SUCCESS;
}

int DataChannel::dispatchChannelWorker(DataCell *pDataCell)
{
    ChannelWorker *pWorker = NULL;
    int workerIndex = (m_nextWorkerId++) % m_numWorker;
    pWorker = m_vChannelWorkers[workerIndex];
    if (NULL == pWorker)
    {
        LOG_WARN("the worker[%d] is invalid, failed to dispatch datacell[0x%x] to worker", workerIndex, pDataCell);
        delete pDataCell;
        return GSE_ERROR;
    }

    pWorker->Notify(pDataCell);
    return GSE_SUCCESS;
}
}
}
