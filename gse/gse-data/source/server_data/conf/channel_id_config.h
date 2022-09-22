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

#ifndef _GSE_DATA_CONFIG_CHANNELID_V2_H_
#define _GSE_DATA_CONFIG_CHANNELID_V2_H_

#include <list>
#include <string>
#include <unordered_map>
#include <vector>

#include <event2/event.h>
#include <json/json.h>

#include "safe/lock.h"
#include "tools/thread.h"
#include "tools/finally.hpp"
#include "tools/rapidjson_macro.h"
#include "tools/strings.h"

#include "api/channelid_def.h"
#include "api/channelid_struct.h"
#include "conf/conf_common.h"
#include "datacell.h"

namespace gse {
namespace data {

enum TopicType
{
    UNKONW_TOPIC = 0,
    KAFKA_TOPIC = 1,
    REDIS_CHANNEL = 2,
    PULSAR_TOPIC = 3,
    DSPROXY = 4
};

typedef enum ZkEventType_
{
    ZK_EVENT_LIST = 0,
    ZK_EVENT_CHANGE = 1,
    ZK_EVENT_CREATE = 2,
    ZK_EVENT_DELETE = 3
} ZkEventType;

class ZkEvent
{
public:
    ZkEventType m_eventType;
    void *m_msg;
    ZkEvent()
        : m_eventType(ZK_EVENT_CHANGE), m_msg(NULL)
    {
    }
};

class ChannelIDFilter
{
public:
    ChannelIDFilter()
        : m_streamToName(""), m_configName(""), m_tableName(""), m_andFilterStringInProtocol(""), m_isNeedCheckAndFilterInData(false)
    {
    }
    ~ChannelIDFilter()
    {
        for (std::vector<StreamFilter *>::iterator it = m_streamFilterAnd.begin(); it != m_streamFilterAnd.end(); ++it)
        {
            delete (*it);
        }
        m_streamFilterAnd.clear();

        for (std::vector<StreamFilter *>::iterator it = m_streamFilterOr.begin(); it != m_streamFilterOr.end(); ++it)
        {
            delete (*it);
        }
        m_streamFilterOr.clear();
    }

public:
    bool IsValidData(DataCell *ptrDataCell);
    void GetTopicName(std::string &topicname);
    std::string ToString();
    uint32_t GetMemSize();

private:
    bool CheckAndFilter(DataCell *ptrDataCell);
    bool checkOrFilter(DataCell *ptrDataCell);

    bool CheckProtocolFilter(DataCell *ptrDataCell, StreamFilter *filter);
    bool CheckDataFilter(DataCell *ptrDataCell, StreamFilter *filter);

public:
    std::string m_streamToName;
    std::string m_configName;
    std::string m_tableName;
    std::vector<StreamFilter *> m_streamFilterAnd;
    std::vector<StreamFilter *> m_streamFilterOr;
    std::string m_andFilterStringInProtocol;
    bool m_isNeedCheckAndFilterInData;
};

class ChannelIdExporterConfig
{
public:
    ChannelIdExporterConfig();
    ~ChannelIdExporterConfig();

    ChannelIdExporterConfig &operator=(ChannelIdExporterConfig &srcConf);

    ChannelIdExporterConfig(ChannelIdExporterConfig &srcConf);

    inline void SetNeedDelete()
    {
        m_setDeleteTimestamp = gse::tools::time::GetUTCSecond();
    }
    inline bool IsNeedDelete()
    {
        return m_setDeleteTimestamp == 0 ? false : ((gse::tools::time::GetUTCSecond() - m_setDeleteTimestamp) > 60);
    }

public:
    StreamToCluster m_streamToCluster;
    int m_streamToId;

private:
    int m_setDeleteTimestamp;
};

struct ZkMessageCount
{
    int m_requestCount;
    int m_responseCount;
    ZkMessageCount()
    {
        m_requestCount = 0;
        m_responseCount = 0;
    }
    void AddNeedAckCount(int count)
    {
        m_requestCount += count;
    }
    void AddResponse()
    {
        m_responseCount++;
    }
    bool Finished()
    {
        return ((m_requestCount == m_responseCount) ? true : false);
    }
};
class ChannelIDConfig
{
public:
    ChannelIDConfig()
    {
        m_ptrMetadata = new Metadata();
        m_setDeleteTimestamp = 0;
        m_bMetaInitFinish = false;
        m_bFilterInitFinish = false;
        m_bStreamToInitFinish = false;
        m_success = true;
    }
    ChannelIDConfig(std::string channelid_str)
        : m_strChannelId(channelid_str)
    {
        m_ptrMetadata = new Metadata();
        m_setDeleteTimestamp = 0;
        m_success = true;
    }
    ~ChannelIDConfig()
    {
        if (NULL != m_ptrMetadata)
        {
            delete m_ptrMetadata;
        }

        for (std::vector<Channel *>::iterator it = m_channels.begin(); it != m_channels.end(); ++it)
        {
            delete (*it);
        }

        for (std::vector<StreamFilter *>::iterator it = m_streamFilter.begin(); it != m_streamFilter.end(); ++it)
        {
            delete (*it);
        }
    }

    inline void SetNeedDelete()
    {
        m_setDeleteTimestamp = gse::tools::time::GetUTCSecond();
    }
    inline bool IsNeedDelete()
    {
        return m_setDeleteTimestamp == 0 ? false : ((gse::tools::time::GetUTCSecond() - m_setDeleteTimestamp) > 120);
    }

    void AddNeedAckCount(int count)
    {
        m_zkReqResponseCount.AddNeedAckCount(count);
    }

    void AddResponseCount()
    {
        m_zkReqResponseCount.AddResponse();
    }

    bool IsComplete()
    {
        bool all_done = false;
        return all_done;
    }

    bool IsSuccess()
    {
        return m_success;
    }

    void SetError()
    {
        m_success = false;
    }

public:
    StreamFilter *GetFilter(const std::string &filterName);

public:
    Metadata *m_ptrMetadata;
    std::vector<StreamFilter *> m_streamFilter;
    std::vector<Channel *> m_channels;
    std::string m_originData;

    std::string m_strChannelId;

    ZkMessageCount m_zkReqResponseCount;
    bool m_success;
    int m_setDeleteTimestamp;

private:
    bool m_bMetaInitFinish;
    bool m_bFilterInitFinish;
    bool m_bStreamToInitFinish;
};

// inner channelid
class ChannelIdStreamConfig
{
public:
    uint32_t m_channelID;
    uint32_t m_streamToId;
    std::string m_channelName;
    int m_storageType;
    union StorageTable
    {
        KafkaTopic *m_kafkaStorageTable;
        PulsarTopic *m_pulsarStorageTable;
        RedisChannel *m_redisStorageTable;
    } m_storageTable;

    ChannelIDFilter m_filter;

public:
    ChannelIdStreamConfig()
        : m_channelName(""), m_channelID(0), m_streamToId(0),
          m_next(NULL), m_storageType(UNKONW_TOPIC),
          m_setDeleteTimestamp(0)
    {
        m_storageTable.m_kafkaStorageTable = nullptr;
    }

    ~ChannelIdStreamConfig()
    {
        if (NULL != m_next)
        {
            delete m_next;
        }

        if (m_storageType == KAFKA_TOPIC)
        {
            if (m_storageTable.m_kafkaStorageTable != NULL)
            {
                delete m_storageTable.m_kafkaStorageTable;
                m_storageTable.m_kafkaStorageTable = NULL;
            }
        }
        else if (m_storageType == REDIS_CHANNEL)
        {
            if (m_storageTable.m_redisStorageTable != NULL)
            {
                delete m_storageTable.m_redisStorageTable;
                m_storageTable.m_redisStorageTable = NULL;
            }
        }
        else if (m_storageType == PULSAR_TOPIC)
        {
            if (m_storageTable.m_pulsarStorageTable != NULL)
            {
                delete m_storageTable.m_pulsarStorageTable;
                m_storageTable.m_pulsarStorageTable = NULL;
            }
        }
    }

    bool IsNeedDelete()
    {
        return m_setDeleteTimestamp == 0 ? false : ((gse::tools::time::GetUTCSecond() - m_setDeleteTimestamp) > 60);
    }

    void GetTableName(std::string &table_name)
    {
        switch (m_storageType)
        {
        case KAFKA_TOPIC:
            table_name = m_storageTable.m_kafkaStorageTable->m_topicName;
            break;
        case PULSAR_TOPIC:
            table_name = m_storageTable.m_pulsarStorageTable->m_topicName;
            break;
        case REDIS_CHANNEL:
            table_name = m_storageTable.m_redisStorageTable->m_channelName;
            break;
        }
    }

    void CopyTableConfig(StreamTo *ptrStreamTo)
    {

        if (ptrStreamTo->m_reportMode == ChannelIDOperationRequestMethodKafka)
        {
            m_storageTable.m_kafkaStorageTable = new KafkaTopic();
            m_storageTable.m_kafkaStorageTable->CopyFrom(&ptrStreamTo->m_kafkaTopic);
            m_storageType = KAFKA_TOPIC;
        }
        else if (ptrStreamTo->m_reportMode == ChannelIDOperationRequestMethodPulsar)
        {
            m_storageTable.m_pulsarStorageTable = new PulsarTopic();
            m_storageTable.m_pulsarStorageTable->CopyFrom(&ptrStreamTo->m_pulsarTopic);
            m_storageType = PULSAR_TOPIC;
        }
        else if (ptrStreamTo->m_reportMode == ChannelIDOperationRequestMethodRedis)
        {
            m_storageTable.m_redisStorageTable = new RedisChannel();
            m_storageTable.m_redisStorageTable->CopyFrom(&ptrStreamTo->m_redisChannel);
            m_storageType = REDIS_CHANNEL;
        }
    }

    void AddAndFilters(StreamFilter *filter)
    {
        StreamFilter *ptr_filter = new StreamFilter();
        ptr_filter->CopyFrom(filter);
        m_filter.m_streamFilterAnd.push_back(ptr_filter);
        if (filter->m_fieldIn.compare("protocol") == 0)
        {
            m_filter.m_andFilterStringInProtocol.append(filter->m_fieldDataValue);
        }
        else
        {
            m_filter.m_isNeedCheckAndFilterInData = true;
        }
    }

    void AddOrFilters(StreamFilter *filter)
    {
        StreamFilter *ptr_filter = new StreamFilter();
        ptr_filter->CopyFrom(filter);
        m_filter.m_streamFilterOr.push_back(ptr_filter);
    }

    inline void SetNeedDelete()
    {
        m_setDeleteTimestamp = gse::tools::time::GetUTCSecond();
    }
    std::string ToString()
    {
        std::string fmt_str;
        return fmt_str;
    }

    uint32_t GetUseMemSize()
    {
        uint32_t memSize = 0;
        memSize = sizeof(ChannelIdStreamConfig) + m_channelName.length();
        if (m_storageType == KAFKA_TOPIC)
        {
            if (m_storageTable.m_kafkaStorageTable != NULL)
            {
                memSize += m_storageTable.m_kafkaStorageTable->GetMemSize();
            }
        }
        else if (m_storageType == REDIS_CHANNEL)
        {
            if (m_storageTable.m_redisStorageTable != NULL)
            {
                memSize += m_storageTable.m_redisStorageTable->GetMemSize();
            }
        }
        else if (m_storageType == PULSAR_TOPIC)
        {
            if (m_storageTable.m_pulsarStorageTable != NULL)
            {
                memSize += m_storageTable.m_pulsarStorageTable->GetMemSize();
            }
        }
        memSize += m_filter.GetMemSize();
        return memSize;
    }

    void Dump(std::string &strDump);

public:
    inline void SetNext(ChannelIdStreamConfig *ptrNext)
    {
        if (NULL != m_next)
        {
            m_next->SetNext(ptrNext);
            return;
        }
        m_next = ptrNext;
    }

public:
    ChannelIdStreamConfig *m_next;
    int m_setDeleteTimestamp;
};

class Exporter;

class ChannelIdStreamExporterManager : public gse::tools::thread::Thread
{
public:
    ChannelIdStreamExporterManager();
    virtual ~ChannelIdStreamExporterManager();

public:
    int Init();
    int ThreadFun();
    void Stop();

    int UpdateExporterConfig(ZkEvent *event);
    void CleanInvalidExporter();
    Exporter *GetExport(uint32_t exporter_id);
    int DeleteExporter(uint32_t exporter_id);

    void FreeEvent();
    bool Find(uint32_t exporterId);

private:
    Exporter *CreateExporter(ChannelIdExporterConfig *ptr_stream_to_id_config);
    int StreamExporterUpdate(ChannelIdExporterConfig *ptr_stream_to_id_config);

private:
    static void ExporterUpdateEvent(int fd, short which, void *v);
    static void InvalidExporterCleanTimerEvent(int fd, short which, void *v);

private:
    gse::safe::RWLock m_rwLock;
    std::map<uint32_t, Exporter *> m_exporters;
    std::list<Exporter *> m_needDeleteExporters;
    bool m_stoped;
    //  thread run
    int m_exporterUpdateNotifyFd[2];
    event_base *m_evBase;
    struct event *m_cleanTimerEvent;
    struct event *m_exporterUpdateEvent;
    struct timeval m_cleanTime;
};

class ChannelIdManager : public gse::tools::thread::Thread
{
public:
    ChannelIdManager(bool is_platid = false);
    virtual ~ChannelIdManager();

public:
    int ThreadFun();
    void Stop();

    // channel id manager
    int Update(ZkEvent *event);
    void FreeEvent();
    bool WriteByFilter(DataCell *ptr_datacell);
    void FreeChannelIdPtr(ChannelIDConfig *ptr);
    bool Find(uint32_t channelId);
    ChannelIdStreamConfig *GetChannelStreamConfig(uint32_t channel_id);

    // only for stack
    void GetChannelIdListByStreamId(uint32_t streamid, std::vector<uint32_t> &channelidList);

public:
    int m_channelIdUpudateNotifyFd[2];

private:
    static void ChannelIdUpdateEvent(int fd, short which, void *v);
    static void InvalidDataCleanTimerEvent(int fd, short which, void *v);
    int HandleChannelIdUpdate(ChannelIDConfig *ptr_channel_id_config);
    int HandleChannelIdDelete(uint32_t channelid);
    void CleanInvalidChannelId();
    void CleanChannelIdInvalidPtr();

private:
    gse::safe::RWLock m_rwLock;
    std::unordered_map<uint32_t, ChannelIdStreamConfig *> m_channelIds;
    std::list<ChannelIdStreamConfig *> m_needDeleteChannelIds;

    gse::safe::RWLock m_freeChannelidLock;
    std::map<std::string, ChannelIDConfig *> m_needFreeChannelIDconfig;

    bool m_isPlatId;
    //  thread run
    event_base *m_evBase;
    struct event *m_cleanTimerEvent;
    struct event *m_channelIdUpdateEvent;
    struct timeval m_cleanTime;
    bool m_stoped;
};

} // namespace data
} // namespace gse
#endif
