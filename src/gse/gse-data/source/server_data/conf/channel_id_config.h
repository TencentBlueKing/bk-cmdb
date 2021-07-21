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

#include <string>
#include <vector>
#include <list>
#include <json/json.h>
#include <event2/event.h>
#include "tools/strings.h"
#include "safe/lock.h"
#include "thread/thread.h"

#include "datacell.h"
#include "conf/bkdata_config_v1.h"
#include "conf/conf_common.h"

namespace gse { 
namespace dataserver {

// config tglog/659 or tdm
#ifndef ZK_CHANNEL_ID_CONFIG_BASE_PATH
#define ZK_CHANNEL_ID_CONFIG_BASE_PATH (std::string("/gse/config/server/dataserver/channelid"))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_PATH
#define ZK_CHANNEL_ID_CONFIG_PATH(id) (std::string("/gse/config/server/dataserver/channelid/") + std::string(id))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_FILTER_PATH
#define ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/filter")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_METADATA_PATH
#define ZK_CHANNEL_ID_CONFIG_METADATA_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/metadata")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_FILTERNAME_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_FILTERNAME_PATH(channelID, filterName) (ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channelID) + "/" + std::string(filterName))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_EXPORTER_PATH
#define ZK_CHANNEL_ID_CONFIG_EXPORTER_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/exporter")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_EXPORTERNAME_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_EXPORTERNAME_PATH(channelID, exporterName) (ZK_CHANNEL_ID_CONFIG_EXPORTER_PATH(channelID) + "/" + std::string(exporterName))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH
#define ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/channel")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_CHANNELNAME_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_CHANNELNAME_PATH(channelID, channelName) (ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channelID) + "/" + std::string(channelName))
#endif



//TGLOG channelid watch path

#ifndef ZK_CHANNEL_ID_TGLOG_NOTIFY_PATH
#define ZK_CHANNEL_ID_TGLOG_NOTIFY_PATH(watch_path) (std::string("/gse/config/server/dataserver/tglog/") + watch_path)
#endif

//--------------------------------------------------
//streamId
#ifndef ZK_STREAM_ID_CONFIG_BASE_PATH
#define ZK_STREAM_ID_CONFIG_BASE_PATH std::string("/gse/config/server/dataserver/streamto")
#endif

#ifndef ZK_STREAM_ID_CONFIG_PATH
#define ZK_STREAM_ID_CONFIG_PATH(streamToId) ZK_STREAM_ID_CONFIG_BASE_PATH + "/" + std::string(streamToId)
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_METADATA_PATH
#define ZK_STREAM_TO_ID_CONFIG_METADATA_PATH(streamToId) (ZK_STREAM_ID_CONFIG_PATH(streamToId) + "/metadata")
#endif


const std::string kRedisSentinel = "sentinel";
const std::string kRedisSingle = "single";
const std::string kRedisCluster = "cluster";

enum StorageTypeEn
{
    EN_KAFKA_TYPE = 0,
    EN_PULSAR_TYPE = 1,
    EN_REDIS_TYPE  = 2
};


enum StorageTableTypeEn
{
    EN_KAFKA_STORAGE_TABLE = 0,
    EN_PULSAR_STORAGE_TABLE = 1,
    EN_REDIS_SOTRAGE_TABLE = 2
};

class Address
{
public:
    std::string m_ip;
    uint16_t m_port;

public:
    Address() : m_ip(""), m_port(0) {}

    std::string ToString()
    {
        std::string strbuff = "Host:{host:%s, port:%d}";
        std::string str_result;
        char buff[1024] = {0};
        snprintf(buff, sizeof(buff), strbuff.c_str(), m_ip.c_str(), m_port);
        str_result.append(buff);
        return str_result;
    }
};

class MetaLabel
{
public:
    std::string m_odm;
    std::string m_bkBizName;
    int m_bizid;
public:
    MetaLabel() : m_odm(""), m_bkBizName(""), m_bizid(0) {}
};

class Metadata
{
public:
    std::string m_name;
    std::string m_version;
    uint32_t m_channelID;
    std::string m_requestPlatName;
    MetaLabel m_label;

public:
    Metadata() : m_name(""), m_version(""), m_channelID(0), m_requestPlatName("") {}
};


class TargetStorageTable
{
public:
    std::string m_name;
    std::string m_tenant;
    std::string m_namespace;
    std::string m_persistent;
    std::string m_type;//kafka,pular,redis....
public:
    TargetStorageTable() : m_name(""),m_tenant(""), m_namespace(""), m_persistent(""),m_type("")  {}
};

//"biz_id":0,"data_set":"0bkmonitor_1500145","partition":1,"topic_name":"0bkmonitor_15001450"
class KafkaStorageTable
{
public:
    std::string m_bizId;
    std::string m_dataSet;
    std::string m_topicName;
    int m_partition;

public:
    KafkaStorageTable() : m_topicName(""),m_bizId(""), m_dataSet(""), m_partition(1)  {}

    inline void CopyFrom(KafkaStorageTable *src)
    {
        if (src == NULL)
        {
            return;
        }
        m_bizId = src->m_bizId;
        m_dataSet = src->m_dataSet;
        m_topicName = src->m_topicName;
        m_partition = src->m_partition;
        if (m_topicName == "")
        {
            m_topicName = src->m_dataSet + src->m_bizId;
        }
    }
};

class PulsarStorageTable
{
public:
    std::string m_tenant;
    std::string m_namespace;
    std::string m_persistent;
    std::string m_bizId;
    std::string m_dataSet;
    std::string m_topicName;
    PulsarStorageTable():m_tenant(""),m_namespace(""), m_persistent(""), m_bizId(""), m_dataSet(""), m_topicName("") {}

    inline void CopyFrom(PulsarStorageTable *src)
    {
        if (src == NULL)
        {
            return;
        }
        m_tenant = src->m_tenant;
        m_namespace = src->m_namespace;
        m_persistent = src->m_persistent;
        m_bizId = src->m_bizId;
        m_dataSet = src->m_dataSet;
        m_topicName = src->m_topicName;


        std::string pulsar_topic_name;
        if (!m_persistent.empty())
        {
            pulsar_topic_name.append(m_persistent);
            pulsar_topic_name.append("://");
        }
        else
        {
            pulsar_topic_name.append("persistent");
            pulsar_topic_name.append("://");
        }

        if (!m_tenant.empty())
        {
            pulsar_topic_name.append(m_tenant);
            pulsar_topic_name.append("/");
        }
        if (!m_namespace.empty())
        {
            pulsar_topic_name.append(m_namespace);
            pulsar_topic_name.append("/");
        }

        if (m_topicName=="")
        {
            m_topicName = m_dataSet + m_bizId;
        }
        m_topicName = pulsar_topic_name.append(m_topicName);
    }
};

class RedisStorageTable
{
public:
    std::string m_channelName;
    std::string m_bizId;
    std::string m_dataSet;
    RedisStorageTable():m_channelName(""),m_bizId(""),m_dataSet("") {}

    inline void CopyFrom(RedisStorageTable *src)
    {
        if (src == NULL)
        {
            return;
        }
        m_channelName = src->m_channelName;
        m_bizId = src->m_bizId;
        m_dataSet = src->m_dataSet;
        if (m_channelName=="")
        {
            m_channelName = m_dataSet + m_bizId;
        }
    }
};

class StreamTo
{
public:
    std::string m_name;
    std::string m_configName;
    int m_storageType;
    union StorageTable{
        KafkaStorageTable *m_kafkaStorageTable;
        PulsarStorageTable *m_pulsarStorageTable;
        RedisStorageTable *m_redisStorageTable;
    }m_storageTable;

    std::vector<std::string> m_filterNameAnd;
    std::vector<std::string> m_filterNameOr;
    int m_streamToId;

public:
    inline StreamTo *Clone()
    {
        StreamTo *ptr_tmp = new StreamTo();
        ptr_tmp->CopyFrom(this);
        return ptr_tmp;
    }
    inline void CopyFrom(StreamTo *src)
    {
        m_name = src->m_name;
        m_configName = src->m_configName;
        m_filterNameAnd.clear();
        m_filterNameOr.clear();
        for (std::vector<std::string>::iterator it = src->m_filterNameAnd.begin(); it != src->m_filterNameAnd.end(); ++it)
        {
            m_filterNameAnd.push_back((*it));
        }
        for (std::vector<std::string>::iterator it = src->m_filterNameOr.begin(); it != src->m_filterNameOr.end(); ++it)
        {
            m_filterNameOr.push_back((*it));
        }
    }

public:
    StreamTo() : m_name(""), m_configName(""), m_streamToId(-1)
    {
        m_storageType = UNKNOWN;
    }

    ~StreamTo()
    {
        switch (m_storageType)
        {
        case KAFKA_COMMON:
            if (m_storageTable.m_kafkaStorageTable != NULL)
            {
                delete m_storageTable.m_kafkaStorageTable;
                m_storageTable.m_kafkaStorageTable = NULL;
            }
            break;
        case EXPORT_PULSAR:
            if (m_storageTable.m_pulsarStorageTable != NULL)
            {
                delete m_storageTable.m_pulsarStorageTable;
                m_storageTable.m_pulsarStorageTable = NULL;
            }
            break;
        case REDIS_SENTINEL_PUB:
            if (m_storageTable.m_redisStorageTable != NULL)
            {
                delete m_storageTable.m_redisStorageTable;
                m_storageTable.m_redisStorageTable = NULL;
            }
            break;
        }
    }
};


class StreamConfig
{
public:
    std::string m_name;
    std::string m_reportMode;
    std::string m_dataLogDir;
    std::string m_dataLogFileName;
    std::string m_pulsarUrl;
    std::vector<Address *> m_addresses;
    std::string m_token;
    KafkaConfig m_kafkaConfig;


public:
    inline StreamConfig *Clone()
    {
        StreamConfig *ptr_tmp = new StreamConfig();
        ptr_tmp->CopyFrom(this);
        return ptr_tmp;
    }

    inline void CopyFrom(StreamConfig *src)
    {
        m_name = src->m_name;
        m_reportMode = src->m_reportMode;
        m_dataLogDir = src->m_dataLogDir;
        m_dataLogFileName = src->m_dataLogFileName;
        m_pulsarUrl = src->m_pulsarUrl;
        m_kafkaConfig = src->m_kafkaConfig;
        m_token = src->m_token;
        if (!m_addresses.empty())
        {
            for (std::vector<Address *>::iterator it = m_addresses.begin(); it != m_addresses.end(); ++it)
            {
                delete (*it);
            }
            m_addresses.clear();
        }

        for (std::vector<Address *>::iterator it = src->m_addresses.begin(); it != src->m_addresses.end(); ++it)
        {
            Address *ptr_tmp = new Address();
            ptr_tmp->m_ip = (*it)->m_ip;
            ptr_tmp->m_port = (*it)->m_port;
            m_addresses.push_back(ptr_tmp);
        }
    }

    std::string ToString();
public:
    StreamConfig() : m_name(""), m_reportMode(""), m_dataLogDir(""), m_dataLogFileName("")
    {
        m_token = "";
    }

    ~StreamConfig()
    {
        for (std::vector<Address *>::iterator it = m_addresses.begin(); it != m_addresses.end(); ++it)
        {
            delete (*it);
        }
    }
};

class StreamFilter
{
public:
    int16_t m_fieldIndex;
    std::string m_name;
    std::string m_fieldDataType;
    std::string m_fieldDataValue;
    std::string m_separator;
    std::string m_fieldIn;

public:
    StreamFilter() : m_fieldIndex(0), m_name(""), m_fieldDataType(""), m_fieldDataValue(""), m_separator(""), m_fieldIn("")
    {
    }

    std::string ToString();
public:
    inline StreamFilter *Clone()
    {
        StreamFilter *ptr_tmp = new StreamFilter();
        ptr_tmp->CopyFrom(this);
        return ptr_tmp;
    }
    inline void CopyFrom(StreamFilter *src)
    {
        m_name = src->m_name;
        m_fieldIndex = src->m_fieldIndex;
        m_fieldDataType = src->m_fieldDataType;
        m_fieldDataValue = src->m_fieldDataValue;
        m_separator = src->m_separator;
        m_fieldIn = src->m_fieldIn;
    }
};

class ChannelIDFilter
{
public:
    ChannelIDFilter() : m_streamToName(""), m_configName(""), m_tableName(""), m_andFilterStringInProtocol(""), m_isNeedCheckAndFilterInData(false) {
    }
    ~ChannelIDFilter()
    {
    }

public:
    bool IsValidData(DataCell *ptrDataCell);
    void GetTopicName(std::string &topicname);
    std::string ToString();
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

class ChannelIDStorage
{

public:
    uint32_t m_channelID;
    std::string m_streamToName;
    StorageType m_storageType;
    StreamConfig *m_ptrConfigInfo;
    ChannelIDFilter *m_ptrChannelIDFilters;

public:
    ChannelIDStorage() : m_streamToName("")
    {
        m_channelID = 0;
        m_ptrConfigInfo = NULL;
        m_ptrChannelIDFilters = new ChannelIDFilter();
        m_storageType = UNKNOWN;
        m_next = NULL;
    }

    ~ChannelIDStorage()
    {
        delete m_ptrChannelIDFilters;
        m_ptrChannelIDFilters = NULL;
        if (NULL != m_next)
        {
            delete m_next;
        }
    }

    std::string ToString();

public:
    inline void SetNext(ChannelIDStorage *ptrNext)
    {
        if (NULL != m_next)
        {
            m_next->SetNext(ptrNext);
            return;
        }
        m_next = ptrNext;
    }

public:
    ChannelIDStorage *m_next;
};

class KafkaStorage
{
public:
    KafkaStorage();
    ~KafkaStorage();
public:
    std::vector<Address> m_addresses;
    KafkaConfig m_kafkaConfig;
    std::string ToString();
};


class PulsarStorage
{
public:
    PulsarStorage(){}
    ~PulsarStorage(){}
public:
    std::string m_token;
    std::vector<Address> m_addresses;
    std::string ToString()
    {
        std::string str_pulsar_storage_addr;
        std::vector<Address>::iterator it;
        for (it = m_addresses.begin(); it != m_addresses.end(); it++)
        {
            str_pulsar_storage_addr.append((*it).m_ip);
            str_pulsar_storage_addr.append(":");
            str_pulsar_storage_addr.append(gse::tools::strings::ToString((*it).m_port));
            if (it != m_addresses.end())
            {
                str_pulsar_storage_addr.append(";");
                break;
            }
        }
        return str_pulsar_storage_addr;
    }
};


class RedisStorage
{
public:
    RedisStorage() {}
    ~RedisStorage(){}
public:
    std::vector<Address> m_addresses;
    std::string m_passwd;
    std::string m_masterName;
    std::string ToString()
    {
        std::string str_redis_storage_addr;
        std::vector<Address>::iterator it;
        for (it = m_addresses.begin(); it != m_addresses.end(); it++)
        {
            str_redis_storage_addr.append((*it).m_ip);
            str_redis_storage_addr.append(":");
            str_redis_storage_addr.append(gse::tools::strings::ToString((*it).m_port));
            if (it != m_addresses.end())
            {
                str_redis_storage_addr.append(";");
                break;
            }

        }
        return str_redis_storage_addr;
    }
};


typedef enum ZkEventType_
{
    ZK_EVENT_LIST = 0,
    ZK_EVENT_CHANGE = 1,
    ZK_EVENT_CREATE = 2,
    ZK_EVENT_DELETE = 3
}ZkEventType;

class ZkEvent
{
public:
    ZkEventType m_eventType;
    void *m_msg;
    ZkEvent() : m_eventType(ZK_EVENT_CHANGE), m_msg(NULL) {
    }
};

class ChannelIdExporterConfig
{
public:
    ChannelIdExporterConfig();
    ~ChannelIdExporterConfig();

    ChannelIdExporterConfig &operator=(const ChannelIdExporterConfig &srcConf);

    ChannelIdExporterConfig(const ChannelIdExporterConfig &srcConf);

    inline void SetNeedDelete()
    {
        m_setDeleteTimestamp = gse::tools::time::GetUTCSecond();
    }
    inline bool IsNeedDelete()
    {
        return m_setDeleteTimestamp == 0 ? false : ((gse::tools::time::GetUTCSecond() - m_setDeleteTimestamp) > 60);
    }

public:
    std::string m_reportMode;
    int m_storageType;
    union Storage
    {
        KafkaStorage*   m_kafkaStorage;
        PulsarStorage*  m_pulsarStorage;
        RedisStorage* m_redisStorage;
    }m_storage;
    int m_streamToId;
    std::string m_name;
    std::string m_rawJsonStr;
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
        m_bMetaInitFinish  = false;
        m_bFilterInitFinish = false;
        m_bStreamToInitFinish = false;
        m_success = true;
    }
    ChannelIDConfig(std::string channelid_str) : m_strChannelId(channelid_str)
    {
        m_ptrMetadata = new Metadata();
        m_setDeleteTimestamp = 0;
        m_success = true;
        m_bMetaInitFinish  = false;
        m_bFilterInitFinish = false;
        m_bStreamToInitFinish = false;
    }
    ~ChannelIDConfig()
    {
        if (NULL != m_ptrMetadata)
        {
            delete m_ptrMetadata;
        }

        for (std::vector<StreamTo *>::iterator it = m_streamTo.begin(); it != m_streamTo.end(); ++it)
        {
            delete (*it);
        }

        for (std::vector<StreamConfig *>::iterator it = m_streamConfig.begin(); it != m_streamConfig.end(); ++it)
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
        bool all_done =  false;
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

//    void InitLevel1RequestCount()
//    {
//        //request  [filter, channels]
//        m_Level1RRCount.m_requestCount = 2;
//    }
public:
    // 返回ChannelID关联的所有存储的信息
    ChannelIDStorage *ToChannelIDStorage();
    StreamFilter *GetFilter(const std::string &filterName);

private:
    StreamConfig *GetConfig(const std::string &configName);

public:
    Metadata *m_ptrMetadata;
    std::vector<StreamTo *> m_streamTo;
    std::vector<StreamConfig *> m_streamConfig;
    std::vector<StreamFilter *> m_streamFilter;
    std::string m_originData;

    std::string m_strChannelId;
//    ZkMessageCount m_Level1RRCount;
//    ZkMessageCount m_channelRRCount;
//    ZkMessageCount m_FilterRRCount;
//    ZkMessageCount m_MetaRRCount;
    ZkMessageCount m_zkReqResponseCount;
    bool m_success;
    int m_setDeleteTimestamp;
private:

    bool m_bMetaInitFinish;
    bool m_bFilterInitFinish;
    bool m_bStreamToInitFinish;
};

class ChannelIDConfigFactory
{
public:
    ChannelIDConfigFactory(){};
    virtual ~ChannelIDConfigFactory(){};

public:
    bool ParseMetadata(const Json::Value &inputJson, Metadata *ptrMetadata, std::string &errorMsg);
    bool ParseStreamTo(const Json::Value &inputJson, std::vector<StreamTo *> &streamTo, std::string &errorMsg);
    bool ParseStreamConfig(const Json::Value &inputJson, std::vector<StreamConfig *> &streamConfig, std::string &errorMsg);
    bool ParseStreamFilter(const Json::Value &inputJson, std::vector<StreamFilter *> &streamFilter, std::string &errorMsg);


    bool ParseExporterConfig(const Json::Value &inputJson, ChannelIdExporterConfig *exporter_cfg, std::string &errorMsg);
    bool ParseStorageAddresses(Json::Value &storage_addresses, std::vector<Address> &addresses);
    void ParseKafkaConfig(const Json::Value &input_json, KafkaStorageTable *kafka_table);
    void ParsePulsarConfig(const Json::Value &input_json, PulsarStorageTable *pulsar_table);
    void ParseRedisConfig(const Json::Value &input_json, RedisStorageTable *redis_table);
};

class ChannelIdStreamConfig
{
public:
    uint32_t m_channelID;
    uint32_t m_streamToId;
    std::string m_channelName;
    ChannelIDFilter *m_ptrChannelIDFilters;
    int m_storageType;
    union StorageTable{
        KafkaStorageTable *m_kafkaStorageTable;
        PulsarStorageTable *m_pulsarStorageTable;
        RedisStorageTable *m_redisStorageTable;
    }m_storageTable;

    std::vector<StreamFilter *> m_streamFilterAnd;
    std::vector<StreamFilter *> m_streamFilterOr;
    ChannelIDFilter m_filter;
public:
    ChannelIdStreamConfig() : m_channelName("")
    {
        m_channelID = 0;
        m_streamToId = 0;
        m_ptrChannelIDFilters = new ChannelIDFilter();
        m_next = NULL;
        m_storageType = UNKNOWN;
        m_setDeleteTimestamp = 0;
    }

    ~ChannelIdStreamConfig()
    {
        delete m_ptrChannelIDFilters;
        m_ptrChannelIDFilters = NULL;
        if (NULL != m_next)
        {
            delete m_next;
        }

        if (m_storageType == KAFKA_COMMON)
        {
            if (m_storageTable.m_kafkaStorageTable != NULL)
            {
                delete m_storageTable.m_kafkaStorageTable;
                m_storageTable.m_kafkaStorageTable = NULL;
            }
        }
        else if (m_storageType == REDIS_SENTINEL_PUB)
        {
            if (m_storageTable.m_redisStorageTable != NULL)
            {
                delete m_storageTable.m_redisStorageTable;
                m_storageTable.m_redisStorageTable = NULL;
            }
        }
        else if (m_storageType == EXPORT_PULSAR)
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
        switch(m_storageType)
        {
            case KAFKA_COMMON:
            table_name = m_storageTable.m_kafkaStorageTable->m_topicName;
            break;
        case EXPORT_PULSAR:
            table_name = m_storageTable.m_pulsarStorageTable->m_topicName;
            break;
        case REDIS_SENTINEL_PUB:
            table_name = m_storageTable.m_redisStorageTable->m_channelName;
            break;
        }
    }

    void CopyTableConfig(StreamTo * ptr_stream_to)
    {
        m_storageType = ptr_stream_to->m_storageType;

        switch(ptr_stream_to->m_storageType)
        {
        case KAFKA_COMMON:
            m_storageTable.m_kafkaStorageTable = new KafkaStorageTable();
            m_storageTable.m_kafkaStorageTable->CopyFrom(ptr_stream_to->m_storageTable.m_kafkaStorageTable);
            break;
        case EXPORT_PULSAR:
            m_storageTable.m_pulsarStorageTable = new PulsarStorageTable();
            m_storageTable.m_pulsarStorageTable->CopyFrom(ptr_stream_to->m_storageTable.m_pulsarStorageTable);
            break;
        case REDIS_SENTINEL_PUB:
            m_storageTable.m_redisStorageTable = new RedisStorageTable();
            m_storageTable.m_redisStorageTable->CopyFrom(ptr_stream_to->m_storageTable.m_redisStorageTable);
            break;
        }
    }

    void AddAndFilters(StreamFilter* filter)
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


    void AddOrFilters(StreamFilter* filter)
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

class ChannelIdExporterManager : public gse::thread::Thread
{
public:
    ChannelIdExporterManager();
    virtual ~ChannelIdExporterManager();

public:
    int Init();
    int run();
    int ThreadFun();
    void stop();

    int UpdateExporterConfig(ZkEvent *event);
    void CleanInvalidExporter();
    Exporter * GetExport(uint32_t exporter_id);
    int DeleteExporter(uint32_t exporter_id);

    void FreeEvent();
private:
    Exporter * CreateExporter(ChannelIdExporterConfig* ptr_stream_to_id_config);
    int Update(ChannelIdExporterConfig* ptr_stream_to_id_config);

private:
    static void ExporterUpdateEvent(int fd, short which, void *v);
    static void InvalidExporterCleanTimerEvent(int fd, short which, void *v);

private:
    gse::safe::RWLock m_rwLock;
    std::map<uint32_t, Exporter*> m_exporters;
    std::list<Exporter*> m_needDeleteExporters;

    //  thread run
    int m_exporterUpdateNotifyFd[2];
    event_base *m_evBase;
    struct event *m_cleanTimerEvent;
    struct event *m_exporterUpdateEvent;
    struct timeval m_cleanTime;
};


class ChannelIdManager : public gse::thread::Thread
{
public:
    ChannelIdManager();
    ChannelIdManager(ChannelIdExporterManager *channelIdExporterManager);
    virtual ~ChannelIdManager();

public:
    int run();
    void stop();
    int ThreadFun();
    //channel id manager
    int Update(ZkEvent *event);
    void FreeEvent();
    //filter
    bool WriteByFilter(DataCell* ptr_datacell);
    void FreeChannelIdPtr(ChannelIDConfig * ptr);

public:
    int m_channelIdUpudateNotifyFd[2];
private:
    static void ChannelIdUpdateEvent(int fd, short which, void *v);
    static void InvalidDataCleanTimerEvent(int fd, short which, void *v);
    ChannelIdStreamConfig *GetChannelStreamConfig(uint32_t channel_id);
    int HandleChannelIdUpdate(ChannelIDConfig* ptr_channel_id_config);
    int HandleChannelIdDelete(uint32_t channelid);
    void CleanInvalidChannelId();
    void CleanChannelIdInvalidPtr();
private:
    gse::safe::RWLock m_rwLock;
    std::map<uint32_t, ChannelIdStreamConfig*> m_channelIds;
    std::list<ChannelIdStreamConfig*> m_needDeleteChannelIds;

    gse::safe::RWLock m_freeChannelidLock;
    std::map<std::string, ChannelIDConfig*> m_needFreeChannelIDconfig;

    ChannelIdExporterManager *m_channelIdExporterManager;

//  thread run
    event_base *m_evBase;
    struct event *m_cleanTimerEvent;
    struct event *m_channelIdUpdateEvent;
    struct timeval m_cleanTime;
};

}
}
#endif
