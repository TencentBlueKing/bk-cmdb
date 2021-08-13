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

#include "log/log.h"
#include "tools/macros.h"
#include "tools/error.h"
#include "tools/time.h"
#include "bbx/gse_errno.h"
#include "tools/pipe.h"

#include "conf/confItem.h"
#include "exporter/exporter.h"
#include "exporter/kafka/kafkaexporter.h"
#include "exporter/redis/redisexporter.h"
#include "exporter/pulsar/pulsarexporter.h"

namespace gse { 
namespace dataserver {

bool ChannelIDFilter::CheckAndFilter( DataCell *ptrDataCell)
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
                //one filter match success,check next
                continue;
            }

            return false;
        }
    }
        // Note:  check in data, this case only check the first filed
    return true;
}

bool ChannelIDFilter::CheckProtocolFilter(DataCell *ptrDataCell, StreamFilter *filter)
{
    std::string extension_str;
    ptrDataCell->GetExtensionString(extension_str);

    LOG_DEBUG("will check or extension's (%s) with the filter value (%s)", SAFE_CSTR(extension_str.c_str()), SAFE_CSTR(filter->m_fieldDataValue.c_str()));
    std::string tagfilter = "[" + filter->m_fieldDataValue + "]";
    if(extension_str.find(tagfilter) != std::string::npos)
    {
        return true;
    }

    return false;
}

std::string StreamFilter::ToString()
{
    std::string strbuff = "{\nName:%s, FieldIndex:%d, FieldDataType:%s, FieldDataValue:%s, Separator:%s, FieldIn:%s\n}";
    std::string str_result;
    char buff[1024] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_name.c_str(), m_fieldIndex, m_fieldDataType.c_str(), m_fieldDataValue.c_str(), m_separator.c_str(), m_fieldIn.c_str());
    str_result.append(buff);
    str_result.append("\n");
    return str_result;
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

    for (std::vector<StreamFilter*>::iterator it = m_streamFilterAnd.begin(); it != m_streamFilterAnd.end(); it++)
    {
        stream_filter_and.append((*it)->ToString());
    }


    std::string stream_filter_or;

    for (std::vector<StreamFilter*>::iterator it = m_streamFilterOr.begin(); it != m_streamFilterOr.end(); it++)
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
    char * pos = strstr(ptrDataCell->GetDataBuf(), filter->m_separator.c_str());
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

StreamConfig *ChannelIDConfig::GetConfig(const std::string &configName)
{
    std::size_t max_count = m_streamConfig.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        StreamConfig *ptr_stream_config = m_streamConfig.at(idx);
        if (ptr_stream_config->m_name.compare(configName) == 0)
        {
            return ptr_stream_config;
        }
    }
    return NULL;
}

ChannelIDStorage *ChannelIDConfig::ToChannelIDStorage()
{
#if 0
    ChannelIDStorage *ptr_channelid_storag = NULL;

    std::size_t max_count = m_streamTo.size();
    LOG_DEBUG("create channel id storage: streamto size:%d", max_count);
    for (std::size_t stream_to_idx = 0; stream_to_idx < max_count; ++stream_to_idx)
    {
        StreamTo *ptr_stream_to = m_streamTo.at(stream_to_idx);
        ChannelIDStorage *ptr_channelid_storage_tmp = new ChannelIDStorage();
        ptr_channelid_storage_tmp->m_streamToName = ptr_stream_to->m_name;
        ptr_channelid_storage_tmp->m_channelID = m_ptrMetadata->m_channelID;
        ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_streamToName = ptr_stream_to->m_name;
        ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_configName = ptr_stream_to->m_configName;
        ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_tableName = ptr_stream_to->m_ptrTargetStorageTable->m_name;

        // get config
        StreamConfig *ptr_stream_config = GetConfig(ptr_stream_to->m_configName);
        if (NULL == ptr_stream_config)
        {
            LOG_ERROR("the config name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(ptr_stream_to->m_configName.c_str()), m_ptrMetadata->m_channelID);
            delete ptr_channelid_storage_tmp;
            continue;
        }
        ptr_channelid_storage_tmp->m_ptrConfigInfo = ptr_stream_config;

        if (ptr_stream_config->m_reportMode == "kafka") // only suport kafka and file
        {
            ptr_channelid_storage_tmp->m_storageType = KAFKA_COMMON;
        }
        else if (ptr_stream_config->m_reportMode == "pulsar")
        {
            ptr_channelid_storage_tmp->m_storageType = EXPORT_PULSAR;
            std::string pulsar_topic_name;
            if (!ptr_stream_to->m_ptrTargetStorageTable->m_persistent.empty())
            {
                pulsar_topic_name.append(ptr_stream_to->m_ptrTargetStorageTable->m_persistent);
                pulsar_topic_name.append("://");
            }
            else
            {
                pulsar_topic_name.append("persistent");
                pulsar_topic_name.append("://");
            }
            if (!ptr_stream_to->m_ptrTargetStorageTable->m_tenant.empty())
            {
                pulsar_topic_name.append(ptr_stream_to->m_ptrTargetStorageTable->m_tenant);
                pulsar_topic_name.append("/");
            }
            if (!ptr_stream_to->m_ptrTargetStorageTable->m_namespace.empty())
            {
                pulsar_topic_name.append(ptr_stream_to->m_ptrTargetStorageTable->m_namespace);
                pulsar_topic_name.append("/");
            }
            pulsar_topic_name.append(ptr_stream_to->m_ptrTargetStorageTable->m_name);
        }
        else
        {
            ptr_channelid_storage_tmp->m_storageType = EXPORT_FILE;
        }

        // get and filter
        std::size_t max_count = m_streamTo.at(stream_to_idx)->m_filterNameAnd.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            std::string filter_name = m_streamTo.at(stream_to_idx)->m_filterNameAnd.at(idx);
            StreamFilter *ptr_stream_filter = GetFilter(filter_name);
            if (NULL == ptr_stream_filter)
            {
                LOG_ERROR("the filter name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(filter_name.c_str()), m_ptrMetadata->m_channelID);
                continue;
            }
            ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_streamFilterAnd.push_back(ptr_stream_filter);

            if (ptr_stream_filter->m_fieldIn.compare("protocol") == 0)
            {
                ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_andFilterStringInProtocol.append(ptr_stream_filter->m_fieldDataValue);
            }
            else
            {
                ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_isNeedCheckAndFilterInData = true;
            }
        }

        //  get or filter
        max_count = m_streamTo.at(stream_to_idx)->m_filterNameOr.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            std::string filter_name = m_streamTo.at(stream_to_idx)->m_filterNameOr.at(idx);
            StreamFilter *ptr_stream_filter = GetFilter(filter_name);
            if (NULL == ptr_stream_filter)
            {
                LOG_ERROR("the filter name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(filter_name.c_str()), m_ptrMetadata->m_channelID);
                continue;
            }
            ptr_channelid_storage_tmp->m_ptrChannelIDFilters->m_streamFilterOr.push_back(ptr_stream_filter);
        }

        if (NULL != ptr_channelid_storag)
        {
            ptr_channelid_storag->SetNext(ptr_channelid_storage_tmp);
        }
        else
        {
            ptr_channelid_storag = ptr_channelid_storage_tmp;
        }
    }

    return ptr_channelid_storag;
#endif
    return NULL;
}

std::string ChannelIDStorage::ToString()
{
    std::string strbuff = "ChannelID:%d, StreamToName:%s, StorageType:%d{\n";
    std::string str_result;
    char buff[1024] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_channelID, m_streamToName.c_str(), m_storageType);
    str_result.append(buff);

    if (m_ptrConfigInfo != NULL)
    {
        str_result.append(m_ptrConfigInfo->ToString());
    }

    if (m_ptrChannelIDFilters != NULL)
    {
        str_result.append(m_ptrChannelIDFilters->ToString());
    }

    str_result.append("\n}\n\n");
    return str_result;
}


std::string StreamConfig::ToString()
{
    std::string host_list;
    for (std::vector<Address*>::iterator it = m_addresses.begin(); it != m_addresses.end(); it++)
    {
        host_list.append((*it)->ToString());
        host_list.append(",");
    }
    host_list.erase(host_list.end()-1);

    std::string strbuff = "StorageConfig:{\nName:%s, ReportMode:%s, LogDir:%s, PulsarURL:%s, Host:%s, Token:%s\n}";
    std::string str_result;
    char buff[4096] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_name.c_str(), m_reportMode.c_str(), m_dataLogDir.c_str(), m_pulsarUrl.c_str(), host_list.c_str(), m_token.c_str());
    str_result.append(buff);
    str_result.append("\n");
    return str_result;
}


bool ChannelIDConfigFactory::ParseMetadata(const Json::Value &metadata, Metadata *ptrMetadata, std::string &errorMsg)
{
    ptrMetadata->m_name = metadata.get("name", "").asString();
    ptrMetadata->m_version = metadata.get("version", "").asString();
    ptrMetadata->m_requestPlatName = metadata.get("plat_name", "").asString();
    ptrMetadata->m_channelID = metadata.get("channel_id", 0).asInt();
    if (metadata.isMember("label"))
    {
        ptrMetadata->m_label.m_bizid = metadata["label"].get("label", 0).asInt();
        ptrMetadata->m_label.m_bkBizName = metadata["label"].get("bk_biz_name", "").asString();
        ptrMetadata->m_label.m_odm = metadata["label"].get("odm", "").asString();
    }
    return true;
}


void ChannelIDConfigFactory::ParseKafkaConfig(const Json::Value &input_json, KafkaStorageTable *kafka_table)
{
    kafka_table->m_topicName = input_json.get("topic_name", "").asString();
    kafka_table->m_bizId = input_json.get("biz_id", "").asString();
    kafka_table->m_dataSet = input_json.get("data_set", "").asString();
    kafka_table->m_partition = input_json.get("partition", -1).asInt();
    return;
}


void ChannelIDConfigFactory::ParsePulsarConfig(const Json::Value &input_json, PulsarStorageTable *pulsar_table)
{
    pulsar_table->m_namespace = input_json.get("namespace", "").asString();
    pulsar_table->m_persistent = input_json.get("persisten", "persistent").asString();
    pulsar_table->m_tenant = input_json.get("tenant", "").asString();
    pulsar_table->m_topicName = input_json.get("topic_name", "").asString();
    pulsar_table->m_bizId = input_json.get("biz_id", "").asString();
    pulsar_table->m_dataSet = input_json.get("data_set", "").asString();
    return;
}

void ChannelIDConfigFactory::ParseRedisConfig(const Json::Value &input_json, RedisStorageTable *redis_table)
{
    redis_table->m_channelName = input_json.get("channel_name", "").asString();
    redis_table->m_bizId = input_json.get("biz_id", "").asString();
    redis_table->m_dataSet = input_json.get("data_set", "").asString();

    return;
}

bool ChannelIDConfigFactory::ParseStreamTo(const Json::Value &inputJson, std::vector<StreamTo *> &streamTo, std::string &errorMsg)
{
    if (!inputJson.isArray())
    {
        errorMsg = "'stream_to' node is not an array";
        return false;
    }

    Json::ArrayIndex maxSize = inputJson.size();
    for (Json::ArrayIndex idx = 0; idx < maxSize; ++idx)
    {
        Json::Value stream_channel = inputJson[idx];
        StreamTo *ptr_stream_to = new StreamTo();
        streamTo.push_back(ptr_stream_to);
        ptr_stream_to->m_name = stream_channel.get("name", "").asString();

        if (!stream_channel.isMember("stream_to"))
        {
            errorMsg = "'stream_to' node is not exist";
            return false;
        }
        Json::Value stream_to_value = stream_channel["stream_to"];
        ptr_stream_to->m_streamToId = stream_to_value.get("stream_to_id", -1).asInt();

        if (stream_to_value.isMember("kafka"))
        {
            Json::Value kafka_config = stream_to_value["kafka"];
            KafkaStorageTable * kafka_table = new KafkaStorageTable();
            ParseKafkaConfig(kafka_config, kafka_table);
            ptr_stream_to->m_storageTable.m_kafkaStorageTable = kafka_table;
            ptr_stream_to->m_storageType = KAFKA_COMMON;
        }
        else if (stream_to_value.isMember("pulsar"))
        {
            Json::Value pulsar_config = stream_to_value["pulsar"];
            PulsarStorageTable * pulsar_table = new PulsarStorageTable();
            ParsePulsarConfig(pulsar_config, pulsar_table);
            ptr_stream_to->m_storageTable.m_pulsarStorageTable = pulsar_table;
            ptr_stream_to->m_storageType = EXPORT_PULSAR;
        }
        else if (stream_to_value.isMember("redis"))
        {
            Json::Value redis_config = stream_to_value["redis"];
            ptr_stream_to->m_storageType = REDIS_SENTINEL_PUB;
            RedisStorageTable * redis_table = new RedisStorageTable();
            ParseRedisConfig(redis_config, redis_table);
            ptr_stream_to->m_storageTable.m_redisStorageTable = redis_table;
        }
        else
        {
            errorMsg = "Invalid configuration, missing one of them (kafka/pulsar/redis)";
            return false;
        }

        if (stream_channel.isMember("filter_name_and"))
        {
            Json::Value filter_name_and = stream_channel["filter_name_and"];
            if (!filter_name_and.isArray())
            {
                errorMsg = "'stream_to/filter_name_and' is not an array";
                return false;
            }

            Json::ArrayIndex filter_name_and_size = filter_name_and.size();
            for (Json::ArrayIndex idx = 0; idx != filter_name_and_size; ++idx)
            {
                if (!filter_name_and[idx].isString())
                {
                    errorMsg = "stream_to/filter_name_and[" + gse::tools::strings::ToString(idx) + "] is not a string value";
                    return false;
                }
                ptr_stream_to->m_filterNameAnd.push_back(filter_name_and[idx].asString());
            }
        }

        if (stream_channel.isMember("filter_name_or"))
        {
            Json::Value filter_name_or = stream_channel["filter_name_or"];
            if (!filter_name_or.isArray())
            {
                errorMsg = "'stream_to/filter_name_or' is not an array";
                return false;
            }

            Json::ArrayIndex filter_name_or_size = filter_name_or.size();
            for (Json::ArrayIndex idx = 0; idx != filter_name_or_size; ++idx)
            {
                if (!filter_name_or[idx].isString())
                {
                    errorMsg = "stream_to/filter_name_or[" + gse::tools::strings::ToString(idx) + "] is not a string value";
                    return false;
                }
                ptr_stream_to->m_filterNameOr.push_back(filter_name_or[idx].asString());
            }
        }
    }
    return true;
}

bool ChannelIDConfigFactory::ParseStreamConfig(const Json::Value &inputJson, std::vector<StreamConfig *> &streamConfig, std::string &errorMsg)
{
    if (!inputJson.isArray())
    {
        errorMsg = "'stream_configs' is not an array";
        return false;
    }

    Json::ArrayIndex stream_config_size = inputJson.size();
    for (Json::ArrayIndex idx = 0; idx < stream_config_size; ++idx)
    {
        StreamConfig *ptr_stream_config = new StreamConfig();
        streamConfig.push_back(ptr_stream_config);
        Json::Value stream_config_value = inputJson[idx];
        ptr_stream_config->m_name = stream_config_value.get("name", "").asString();
        ptr_stream_config->m_reportMode = stream_config_value.get("report_mode", "").asString();
        ptr_stream_config->m_dataLogDir = stream_config_value.get("data_log_path", "").asString();

        // 以单条数据1K计算，最高约 缓存 18GB 的存储空间
        ptr_stream_config->m_dataLogFileName = stream_config_value.get("data_log_filename", "").asString();

        if (ptr_stream_config->m_reportMode == "kafka")
        {
            ptr_stream_config->m_kafkaConfig.m_queueBufferingMaxMessages =  stream_config_value.get("queue_buffering_max_messages", DEFAULT_MAX_KAFKA_QUEUE_SIZE).asString();
            ptr_stream_config->m_kafkaConfig.m_messageMaxBytes =  stream_config_value.get("message_max_bytes", DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE).asString();
            ptr_stream_config->m_kafkaConfig.m_securityProtocol =  stream_config_value.get("security_protocol", "").asString();
            ptr_stream_config->m_kafkaConfig.m_saslMechanisms =  stream_config_value.get("sasl_mechanisms", "").asString();
            ptr_stream_config->m_kafkaConfig.m_saslUserName =  stream_config_value.get("sasl_username", "").asString();
            ptr_stream_config->m_kafkaConfig.m_saslPasswd =  stream_config_value.get("sasl_passwd", "").asString();
            ptr_stream_config->m_kafkaConfig.m_requestRequiredAcks =  stream_config_value.get("request_required_acks", "1").asString();
            ptr_stream_config->m_kafkaConfig.m_queueBufferingMaxMs =  stream_config_value.get("queue_buffering_max_ms", "200").asString();
        }
        else if (ptr_stream_config->m_reportMode == "pulsar")
        {
            ptr_stream_config->m_token = stream_config_value.get("token", "").asString();
        }

        
        if (!stream_config_value.isMember("storage_address"))
        {
            continue;
        }

        Json::Value stream_config_storage_address_value = stream_config_value["storage_address"];
        if (!stream_config_storage_address_value.isArray())
        {
            continue;
        }

        Json::ArrayIndex stream_config_storage_address_value_size = stream_config_storage_address_value.size();
        for (Json::ArrayIndex idx = 0; idx < stream_config_storage_address_value_size; ++idx)
        {
            Address *ptr_address = new Address();
            ptr_stream_config->m_addresses.push_back(ptr_address);
            ptr_address->m_ip = stream_config_storage_address_value[idx].get("ip", "").asString();
            ptr_address->m_port = stream_config_storage_address_value[idx].get("port", 0).asUInt64();
        }
    }
    return true;
}


bool ChannelIDConfigFactory::ParseStorageAddresses(Json::Value &storage_addresses, std::vector<Address> &addresses)
{
    if (!storage_addresses.isArray())
    {
        LOG_ERROR("storage address must bu array, input json(%s)", storage_addresses.toStyledString().c_str());
        return false;
    }

    Json::ArrayIndex storage_address_size = storage_addresses.size();
    for (Json::ArrayIndex idx = 0; idx < storage_address_size; ++idx)
    {
        Address address;
        address.m_ip = storage_addresses[idx].get("ip", "").asString();
        address.m_port = storage_addresses[idx].get("port", 0).asInt();
        addresses.push_back(address);
    }

    return true;
}

bool ChannelIDConfigFactory::ParseExporterConfig(const Json::Value &exporter_cfg_json, ChannelIdExporterConfig * exporter_cfg, std::string &errorMsg)
{
    exporter_cfg->m_name = exporter_cfg_json.get("name", "").asString();
    exporter_cfg->m_reportMode = exporter_cfg_json.get("report_mode", "").asString();
    exporter_cfg->m_rawJsonStr = exporter_cfg_json.toStyledString();
    if (exporter_cfg->m_reportMode == "kafka")
    {
        if (!exporter_cfg_json.isMember("kafka"))
        {
            LOG_ERROR("failed to parse exporter config, the kafka entry is missing (%s)",
                      exporter_cfg_json.toStyledString().c_str());
            return false;
        }

        exporter_cfg->m_storageType = KAFKA_COMMON;
        exporter_cfg->m_storage.m_kafkaStorage = new KafkaStorage();

        KafkaConfig * ptr_kafka_config = &exporter_cfg->m_storage.m_kafkaStorage->m_kafkaConfig;
        ptr_kafka_config->m_queueBufferingMaxMessages =  exporter_cfg_json["kafka"].get("queue_buffering_max_messages", DEFAULT_MAX_KAFKA_QUEUE_SIZE).asString();
        ptr_kafka_config->m_messageMaxBytes =  exporter_cfg_json["kafka"].get("message_max_bytes", DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE).asString();
        ptr_kafka_config->m_securityProtocol =  exporter_cfg_json["kafka"].get("security_protocol", "").asString();
        ptr_kafka_config->m_saslMechanisms =  exporter_cfg_json["kafka"].get("sasl_mechanisms", "").asString();
        ptr_kafka_config->m_saslUserName =  exporter_cfg_json["kafka"].get("sasl_username", "").asString();
        ptr_kafka_config->m_saslPasswd =  exporter_cfg_json["kafka"].get("sasl_passwd", "").asString();
        ptr_kafka_config->m_requestRequiredAcks =  exporter_cfg_json["kafka"].get("request_required_acks", "1").asString();
        ptr_kafka_config->m_queueBufferingMaxMs =  exporter_cfg_json["kafka"].get("queue_buffering_max_ms", "200").asString();

        LOG_DEBUG("exporter name:%s, parse kafka config(securit_protocal:%s, sasl_mechanisms:%s, user_name:%s, buffering_max_ms:%s, buffering_max_message:%s)",
                  exporter_cfg->m_name.c_str(),
                  ptr_kafka_config->m_securityProtocol.c_str(),
                  ptr_kafka_config->m_saslMechanisms.c_str(),
                  ptr_kafka_config->m_saslUserName.c_str(),
                  ptr_kafka_config->m_queueBufferingMaxMs.c_str(),
                  ptr_kafka_config->m_queueBufferingMaxMessages.c_str());

        if (!exporter_cfg_json["kafka"].isMember("storage_address"))
        {
            LOG_ERROR("failed to parse exporter config, the storage_address entry is missing in Kafka JSON(%s)",
                      exporter_cfg_json.toStyledString().c_str());
            return false;
        }

        Json::Value storage_addresses = exporter_cfg_json["kafka"]["storage_address"];
        if (!ParseStorageAddresses(storage_addresses, exporter_cfg->m_storage.m_kafkaStorage->m_addresses))
        {
            LOG_ERROR("failed to parse exporter config, the storage_address is empty",
                      exporter_cfg_json.toStyledString().c_str());
            return false;
        }

    }
    else if (exporter_cfg->m_reportMode == "pulsar")
    {
        if (!exporter_cfg_json.isMember("pulsar"))
        {
            LOG_ERROR("failed to parse exporter config, the pulsar entry is missing (%s)",
                      exporter_cfg_json.toStyledString().c_str());
            return false;
        }

        exporter_cfg->m_storageType = EXPORT_PULSAR;
        exporter_cfg->m_storage.m_pulsarStorage = new PulsarStorage();
        exporter_cfg->m_storage.m_pulsarStorage->m_token = exporter_cfg_json["pulsar"].get("token", "").asString();
        Json::Value storage_addresses = exporter_cfg_json["pulsar"]["storage_address"];
        ParseStorageAddresses(storage_addresses, exporter_cfg->m_storage.m_pulsarStorage->m_addresses);
    }
    else if (exporter_cfg->m_reportMode == "redis")
    {
        if (!exporter_cfg_json.isMember("redis"))
        {
            LOG_ERROR("failed to redis exporter config, the pulsar entry is missing (%s)",
                      exporter_cfg_json.toStyledString().c_str());
            return false;
        }

        exporter_cfg->m_storage.m_redisStorage = new RedisStorage();
        std::string mode = exporter_cfg_json["redis"].get("mode", "sentinel").asString();
        if (mode == kRedisSentinel)
        {
            exporter_cfg->m_storageType = REDIS_SENTINEL_PUB;
        }
        else if (mode == kRedisSingle)
        {
            exporter_cfg->m_storageType = REDIS_PUB;
        }

        exporter_cfg->m_storage.m_redisStorage->m_masterName = exporter_cfg_json["redis"].get("master_name", "").asString();
        exporter_cfg->m_storage.m_redisStorage->m_passwd = exporter_cfg_json["redis"].get("passwd", "").asString();
        Json::Value storage_addresses = exporter_cfg_json["redis"]["storage_address"];
        ParseStorageAddresses(storage_addresses, exporter_cfg->m_storage.m_redisStorage->m_addresses);
    }
    else
    {
        LOG_ERROR("failed to parse exporter config, invalid config (%s)",
                  exporter_cfg_json.toStyledString().c_str());
        return false;
    }
    return true;
}


bool ChannelIDConfigFactory::ParseStreamFilter(const Json::Value &inputJson, std::vector<StreamFilter *> &streamFilter, std::string &errorMsg)
{
    if (!inputJson.isArray())
    {
        errorMsg = "'stream_filters' is not an array";
        return false;
    }

    Json::ArrayIndex stream_filter_size = inputJson.size();
    for (Json::ArrayIndex idx = 0; idx < stream_filter_size; ++idx)
    {
        StreamFilter *ptr_stream_filter = new StreamFilter();
        Json::Value stream_filter_value = inputJson[idx];
        streamFilter.push_back(ptr_stream_filter);
        ptr_stream_filter->m_name = stream_filter_value.get("name", "").asString();
        ptr_stream_filter->m_fieldIndex = stream_filter_value.get("field_index", 0).asInt64();
        ptr_stream_filter->m_fieldDataType = stream_filter_value.get("field_data_type", "").asString();
        ptr_stream_filter->m_fieldDataValue = stream_filter_value.get("field_data_value", "").asString();
        ptr_stream_filter->m_separator = stream_filter_value.get("field_separator", "").asString();
        ptr_stream_filter->m_fieldIn = stream_filter_value.get("field_in", "").asString();
    }
    return true;
}

void ChannelIdExporterManager::CleanInvalidExporter()
{
    std::list<Exporter*>::iterator it = m_needDeleteExporters.begin();
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
    m_storage.m_kafkaStorage = NULL;
    m_storage.m_pulsarStorage = NULL;
    m_storage.m_redisStorage = NULL;
    m_setDeleteTimestamp = 0;
    m_storageType = UNKNOWN;
    m_streamToId = 0;
}
ChannelIdExporterConfig::~ChannelIdExporterConfig()
{
    if (m_storageType == KAFKA_COMMON)
    {
        if (m_storage.m_kafkaStorage != NULL)
        {
            delete m_storage.m_kafkaStorage;
            m_storage.m_kafkaStorage = NULL;
        }
    }
    else if (m_storageType == REDIS_SENTINEL_PUB)
    {
        if (m_storage.m_redisStorage != NULL)
        {
            delete m_storage.m_redisStorage;
            m_storage.m_redisStorage = NULL;
        }
    }
    else if (m_storageType == EXPORT_PULSAR)
    {
        if (m_storage.m_pulsarStorage != NULL)
        {
            delete m_storage.m_pulsarStorage;
            m_storage.m_pulsarStorage = NULL;
        }
    }
}

ChannelIdExporterConfig& ChannelIdExporterConfig::operator=(const ChannelIdExporterConfig &srcConf)
{
    this->m_name = srcConf.m_name;
    this->m_reportMode = srcConf.m_reportMode;
    this->m_streamToId = srcConf.m_streamToId;


    if (srcConf.m_reportMode == "kafka")
    {
        this->m_storageType = KAFKA_COMMON;
        if (this->m_storage.m_kafkaStorage != NULL)
        {
            delete this->m_storage.m_kafkaStorage;
            this->m_storage.m_kafkaStorage = NULL;
        }
        this->m_storage.m_kafkaStorage = new KafkaStorage();
        *this->m_storage.m_kafkaStorage = *srcConf.m_storage.m_kafkaStorage;
        LOG_DEBUG("copy kafka storage config, dest:%s, src:%s",
                  this->m_storage.m_kafkaStorage->m_kafkaConfig.m_saslUserName.c_str(), srcConf.m_storage.m_kafkaStorage->m_kafkaConfig.m_saslUserName.c_str());

    }
    else if (srcConf.m_reportMode == "pulsar")
    {
        this->m_storageType = EXPORT_PULSAR;
        if (this->m_storage.m_pulsarStorage != NULL)
        {
            delete this->m_storage.m_pulsarStorage;
            this->m_storage.m_pulsarStorage = NULL;
        }
        this->m_storage.m_pulsarStorage = new PulsarStorage();
        *this->m_storage.m_pulsarStorage = *srcConf.m_storage.m_pulsarStorage;
    }
    else if (srcConf.m_reportMode == "redis")
    {
        if (this->m_storage.m_redisStorage != NULL)
        {
            delete this->m_storage.m_redisStorage;
            this->m_storage.m_redisStorage = NULL;
        }
        this->m_storageType = srcConf.m_storageType;
        this->m_storage.m_redisStorage = new RedisStorage();
        *this->m_storage.m_redisStorage = *srcConf.m_storage.m_redisStorage;
    }
    return *this;
}

ChannelIdExporterConfig::ChannelIdExporterConfig(const ChannelIdExporterConfig &srcConf)
{
    m_storage.m_kafkaStorage = NULL;
    m_storage.m_pulsarStorage = NULL;
    m_storage.m_redisStorage = NULL;
    m_setDeleteTimestamp = 0;
    *this = srcConf;
}

ChannelIdExporterManager::ChannelIdExporterManager()
{
    m_cleanTime.tv_sec = 10;
    m_cleanTime.tv_usec = 0;
    m_exporterUpdateNotifyFd[0] = -1;
    m_exporterUpdateNotifyFd[1] = -1;
    m_evBase = NULL;
    m_cleanTimerEvent = NULL;
    m_exporterUpdateEvent = NULL;
}

void ChannelIdExporterManager::FreeEvent()
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

ChannelIdExporterManager::~ChannelIdExporterManager()
{
    LOG_DEBUG("delete ChannelIdExporterManager obj(%p)", this);
    FreeEvent();
    std::list<Exporter*>::iterator it = m_needDeleteExporters.begin();
    for (; it != m_needDeleteExporters.end(); it++)
    {
        Exporter *ptr_delete = (*it);
        LOG_DEBUG("stop exporter:(id:%d, name:%s)", ptr_delete->GetExporterId(), ptr_delete->GetExporterName().c_str());
        it = m_needDeleteExporters.erase(it);
        ptr_delete->Stop();
        delete ptr_delete;
    }

    std::map<uint32_t, Exporter*>::iterator m_it = m_exporters.begin();
    for (; m_it != m_exporters.end();)
    {
        Exporter *ptr_delete = m_it->second;
        LOG_DEBUG("stop exporter:(id:%d, name:%s)", ptr_delete->GetExporterId(), ptr_delete->GetExporterName().c_str());
        ptr_delete->Stop();
        delete ptr_delete;
        m_exporters.erase(m_it++);
    }
}

void ChannelIdExporterManager::InvalidExporterCleanTimerEvent(int fd, short which, void *v)
{
    ChannelIdExporterManager *ptr_this = (ChannelIdExporterManager*)v;
    ptr_this->CleanInvalidExporter();
}

int ChannelIdExporterManager::UpdateExporterConfig(ZkEvent *event)
{
    int ret = gse::tools::pipe::WritePipe(m_exporterUpdateNotifyFd[1], &event, sizeof(event));
    return (ret > 0 ? GSE_SUCCESS : GSE_ERROR);
}

void ChannelIdExporterManager::stop()
{
    FreeEvent();
}

void ChannelIdExporterManager::ExporterUpdateEvent(int fd, short which, void *v)
{
    ChannelIdExporterManager *ptr_this = (ChannelIdExporterManager*)v;
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

    switch(event->m_eventType)
    {
    case ZK_EVENT_LIST:
        break;
    case ZK_EVENT_DELETE:
    {
        uint32_t* ptr_channelid = (uint32_t*)event->m_msg;
        ptr_this->DeleteExporter(*ptr_channelid);
        delete ptr_channelid;
        break;
    }
    case ZK_EVENT_CHANGE://no break
    case ZK_EVENT_CREATE:
    {
        exporter_cfg = (ChannelIdExporterConfig *)event->m_msg;
        ptr_this->Update(exporter_cfg);
        LOG_DEBUG("update exporter config:(%s)", exporter_cfg->m_rawJsonStr.c_str());
        delete exporter_cfg;
        break;
    }
    default:
        LOG_ERROR("recv unkown event:%d", event->m_eventType);
    }

    delete event;
    return;

}
int ChannelIdExporterManager::Init()
{
    return 0;
}

int ChannelIdExporterManager::ThreadFun()
{
    int ret = -1;
    gse::tools::pipe::CreatePipe(m_exporterUpdateNotifyFd);
    m_evBase = event_base_new();
    if (NULL == m_evBase)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new event base object for channel id manager thread, errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }
    m_cleanTimerEvent = event_new(m_evBase, -1, EV_PERSIST, ChannelIdExporterManager::InvalidExporterCleanTimerEvent, this);
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
    m_exporterUpdateEvent = event_new(m_evBase, m_exporterUpdateNotifyFd[0], EV_PERSIST|EV_READ, ChannelIdExporterManager::ExporterUpdateEvent, this);
    if (NULL == m_exporterUpdateEvent)
    {
        int tmp_errno = gse_errno;
        std::string tmp_errmsg = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("failed to new channel id update event errno:%d, errmsg:%s", tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
        return GSE_SYSTEMERROR;
    }

    ret = event_add(m_exporterUpdateEvent, NULL);
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
    LOG_FATAL("exporter manager event base loop abnomal break. ret:%d, errno:%d, errmsg:%s", ret, tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
    return GSE_SUCCESS;
}


Exporter * ChannelIdExporterManager::CreateExporter(ChannelIdExporterConfig* ptr_stream_to_id_config)
{
    Exporter *ptr_exporter = NULL;

    switch (ptr_stream_to_id_config->m_storageType)
    {
    case KAFKA_COMMON:
        LOG_DEBUG("create exporter (id:%d, name:%s, type:%d, address:%s) ",
                  ptr_stream_to_id_config->m_streamToId,
                  ptr_stream_to_id_config->m_name.c_str(),
                  ptr_stream_to_id_config->m_storageType,
                  ptr_stream_to_id_config->m_storage.m_kafkaStorage->ToString().c_str()
                  );
        ptr_exporter = new KafkaExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        LOG_INFO("create a new kafka exporter for the storage %s", SAFE_CSTR(ptr_stream_to_id_config->m_storage.m_kafkaStorage->ToString().c_str()));
        break;
    case REDIS_PUB:
        LOG_DEBUG("create redis client for the storage %s", ptr_stream_to_id_config->m_storage.m_redisStorage->ToString().c_str());
        ptr_exporter = new RedisExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        break;
    case REDIS_SENTINEL_PUB:
        LOG_DEBUG("create redis sentinel client for the storage %s", ptr_stream_to_id_config->m_storage.m_redisStorage->ToString().c_str());
        ptr_exporter = new RedisExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        break;
    case EXPORT_DSPROXY:
        LOG_ERROR("unsupport ds proxy exporter config");
        return NULL;
    case EXPORT_PULSAR:
        LOG_INFO("create a new pulsar exporter for the storage %s", SAFE_CSTR(ptr_stream_to_id_config->m_storage.m_pulsarStorage->ToString().c_str()));
        ptr_exporter = new PulsarExporter();
        ptr_exporter->UpdateConf(ExporterConfigWrapper::ChannelIDConfType, ptr_stream_to_id_config);
        break;
    default:
        LOG_ERROR("it is failed to create a exporter for the unkown storage(index %d type %s)", ptr_stream_to_id_config->m_streamToId, ptr_stream_to_id_config->m_reportMode.c_str());
        return NULL;
    }

    if (GSE_SUCCESS != ptr_exporter->Start())
    {
        LOG_ERROR("it is failed to start the exporter (%s) ", SAFE_CSTR(ptr_exporter->GetExporterName().c_str()));
        delete ptr_exporter;
        return NULL;
    }

    return ptr_exporter;

}

Exporter * ChannelIdExporterManager::GetExport(uint32_t exporter_id)
{
    Exporter *ptr_exporter = NULL;
    m_rwLock.RLock();
    std::map<uint32_t, Exporter*>::iterator it = m_exporters.find(exporter_id);
    if (it != m_exporters.end())
    {
        ptr_exporter = it->second;
    }
    m_rwLock.UnLock();
    return ptr_exporter;
}

int ChannelIdExporterManager::DeleteExporter(uint32_t exporter_id)
{
    m_rwLock.WLock();
    std::map<uint32_t, Exporter*>::iterator it = m_exporters.find(exporter_id);
    if (it != m_exporters.end())
    {
        if (it->second != NULL)
        {
            it->second->SetNeedStop();
            m_needDeleteExporters.push_back(it->second);
        }
        LOG_DEBUG("delete expoerter id(%d) exporter config, will stop old exporter(%p))", exporter_id, it->second);
        m_exporters.erase(it);
    }
    m_rwLock.UnLock();
    return GSE_SUCCESS;
}
int ChannelIdExporterManager::Update(ChannelIdExporterConfig* ptr_exporter_cfg)
{
    Exporter* exporter = NULL;
    exporter = CreateExporter(ptr_exporter_cfg);
    if (exporter == NULL)
    {
        LOG_ERROR("failed to update exporter, exporter(name:%s, id:%d, type:%d)",
                  ptr_exporter_cfg->m_name.c_str(), ptr_exporter_cfg->m_streamToId, ptr_exporter_cfg->m_storageType);
        return GSE_ERROR;
    }
    exporter->SetExporterId(ptr_exporter_cfg->m_streamToId);
    m_rwLock.WLock();
    std::map<uint32_t, Exporter*>::iterator it = m_exporters.find(ptr_exporter_cfg->m_streamToId);
    if (it != m_exporters.end())
    {
        if (it->second != NULL)
        {
            it->second->SetNeedStop();
            m_needDeleteExporters.push_back(it->second);
        }
        LOG_DEBUG("update streamToId(%d) exporter config, will stop old exporter(%p)", ptr_exporter_cfg->m_streamToId, it->second);
        m_exporters.erase(it);
    }
    m_exporters.insert(std::make_pair(ptr_exporter_cfg->m_streamToId, exporter));
    m_rwLock.UnLock();
    LOG_DEBUG("push exporter(%d:%p) into manager list", ptr_exporter_cfg->m_streamToId, exporter);
    return GSE_SUCCESS;
}

//channel id config manager
bool ChannelIdManager::WriteByFilter(DataCell* ptr_datacell)
{
    uint32_t channel_id = ptr_datacell->GetChannelID();
    ChannelIdStreamConfig *ptr_channel_id_stream_cfg = GetChannelStreamConfig(channel_id);
    while (ptr_channel_id_stream_cfg != NULL)
    {
        std::string table_name;
        ptr_channel_id_stream_cfg->GetTableName(table_name);
        LOG_DEBUG("write date cell(channleid:%d) filter by channel name(%s), table name(%s)", channel_id, ptr_channel_id_stream_cfg->m_channelName.c_str(), table_name.c_str());
        if (ptr_channel_id_stream_cfg->m_filter.IsValidData(ptr_datacell))
        {
            std::string table_name;
            ptr_channel_id_stream_cfg->GetTableName(table_name);
            LOG_DEBUG("data cell(channleid:%d, stream_to_id:%d) filter success, will send by exporter. table name[%s]", ptr_datacell->GetChannelID(), ptr_channel_id_stream_cfg->m_streamToId, table_name.c_str());
            ptr_datacell->AddTableName(table_name);
            Exporter * exporter = m_channelIdExporterManager->GetExport(ptr_channel_id_stream_cfg->m_streamToId);
            if (exporter != NULL)
            {
                exporter->Write(ptr_datacell);
                LOG_DEBUG("data cell(channleid:%d) filter success, send by exporter[%s]", ptr_datacell->GetChannelID(), exporter->GetExporterName().c_str());
            }
        }
        ptr_channel_id_stream_cfg = ptr_channel_id_stream_cfg->m_next;
    }
    return true;
}

ChannelIdStreamConfig *ChannelIdManager::GetChannelStreamConfig(uint32_t channel_id)
{
    std::map<uint32_t, ChannelIdStreamConfig*>::iterator it;
    ChannelIdStreamConfig* ptr_channel_id_stream_cfg = NULL;
    m_rwLock.RLock();
    it = m_channelIds.find(channel_id);
    if (it != m_channelIds.end())
    {
        ptr_channel_id_stream_cfg = it->second;
    }
    m_rwLock.UnLock();
    return ptr_channel_id_stream_cfg;
}

ChannelIdManager::ChannelIdManager()
{
    m_channelIdExporterManager = NULL;
    m_evBase = NULL;
    m_cleanTimerEvent = NULL;
    m_channelIdUpdateEvent = NULL;
    m_cleanTime.tv_sec = 20;
    m_cleanTime.tv_usec = 0;
}

ChannelIdManager::ChannelIdManager(ChannelIdExporterManager *channelIdExporterManager)
{
    m_cleanTime.tv_sec = 20;
    m_cleanTime.tv_usec = 0;
    m_channelIdExporterManager = channelIdExporterManager;
    m_evBase = NULL;
    m_cleanTimerEvent = NULL;
    m_channelIdUpdateEvent = NULL;
}

ChannelIdManager::~ChannelIdManager()
{
    FreeEvent();
    std::list<ChannelIdStreamConfig*>::iterator it = m_needDeleteChannelIds.begin();
    for (; it != m_needDeleteChannelIds.end();)
    {
        ChannelIdStreamConfig *ptr_delete = (*it);
        if (ptr_delete != NULL)
        {
            delete ptr_delete;
        }

        it = m_needDeleteChannelIds.erase(it);
    }

    std::map<uint32_t, ChannelIdStreamConfig*>::iterator m_it = m_channelIds.begin();
    for (; m_it != m_channelIds.end();)
    {
        ChannelIdStreamConfig *ptr_delete = m_it->second;
        if (ptr_delete != NULL)
        {
            delete ptr_delete;
        }

        m_channelIds.erase(m_it++);
    }
}

void ChannelIdManager::InvalidDataCleanTimerEvent(int fd, short which, void *v)
{
    ChannelIdManager *ptr_this = (ChannelIdManager *)v;

    ptr_this->CleanInvalidChannelId();
    ptr_this->CleanChannelIdInvalidPtr();
}

void ChannelIdManager::FreeChannelIdPtr(ChannelIDConfig * ptr)
{
    if (ptr == NULL)
    {
        return;
    }
    ptr->SetNeedDelete();
    std::string str_key = ptr->m_strChannelId + "|" + gse::tools::strings::ToString((uint64_t)ptr);
    m_freeChannelidLock.WLock();
    std::map<std::string, ChannelIDConfig*>::iterator it = m_needFreeChannelIDconfig.find(str_key);
    if (it == m_needFreeChannelIDconfig.end())
    {
        m_needFreeChannelIDconfig.insert(std::make_pair(str_key, ptr));
        LOG_DEBUG("Will free channelid(%s) mem ptr:%p, key:%s", ptr->m_strChannelId.c_str(), ptr, str_key.c_str());
    }
    m_freeChannelidLock.UnLock();
}
void ChannelIdManager::CleanChannelIdInvalidPtr()
{
    LOG_DEBUG("Free channelid mem, map size:%d", m_needFreeChannelIDconfig.size());
    m_freeChannelidLock.WLock();
    std::map<std::string, ChannelIDConfig*>::iterator it = m_needFreeChannelIDconfig.begin();
    for (; it != m_needFreeChannelIDconfig.end();)
    {
        ChannelIDConfig *ptr_delete = it->second;
        if (ptr_delete != NULL)
        {
            LOG_DEBUG("Free channelid mem, key:%s, %p", it->first.c_str(), it->second);
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
    std::list<ChannelIdStreamConfig*>::iterator it = m_needDeleteChannelIds.begin();
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
    ChannelIDConfig * ptr_channel_id_cfg;
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

    switch(event->m_eventType)
    {
    case ZK_EVENT_LIST:
        break;
    case ZK_EVENT_DELETE:
    {
        uint32_t* ptr_channelid = (uint32_t*)event->m_msg;
        ptr_this->HandleChannelIdDelete(*ptr_channelid);
        delete ptr_channelid;
        break;
    }
    case ZK_EVENT_CHANGE:
    case ZK_EVENT_CREATE:
    {
        ptr_channel_id_cfg = (ChannelIDConfig *)event->m_msg;
        ptr_this->HandleChannelIdUpdate(ptr_channel_id_cfg);
        LOG_DEBUG("update channel id config:(%s)", ptr_channel_id_cfg->m_originData.c_str());
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
        event_base_loopexit(m_evBase, NULL);
        event_base_free(m_evBase);
        m_evBase = NULL;
    }
}
void ChannelIdManager::stop()
{
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
    m_channelIdUpdateEvent = event_new(m_evBase, m_channelIdUpudateNotifyFd[0], EV_PERSIST|EV_READ, ChannelIdManager::ChannelIdUpdateEvent, this);
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
    LOG_FATAL("channel id manager event base loop abnomal break. ret:%d, errno:%d, errmsg:%s", ret, tmp_errno, SAFE_CSTR(tmp_errmsg.c_str()));
    return GSE_SUCCESS;
}
int ChannelIdManager::Update(ZkEvent *event)
{
    int ret = gse::tools::pipe::WritePipe(m_channelIdUpudateNotifyFd[1], &event, sizeof(event));
    return ((ret < 0 ) ? GSE_ERROR : GSE_SUCCESS);
}
int ChannelIdManager::HandleChannelIdDelete(uint32_t channelid)
{
    m_rwLock.WLock();
    std::map<uint32_t, ChannelIdStreamConfig*>::iterator it = m_channelIds.find(channelid);
    if (it != m_channelIds.end())
    {
        ChannelIdStreamConfig* ptr_channelid_config = it->second;
        if (ptr_channelid_config != NULL)
        {
            ptr_channelid_config->SetNeedDelete();
            m_needDeleteChannelIds.push_back(ptr_channelid_config);
        }
        m_channelIds.erase(it);
    }
    m_rwLock.UnLock();
    return GSE_SUCCESS;
}
int ChannelIdManager::HandleChannelIdUpdate(ChannelIDConfig* ptr_channel_id_config)
{
    ChannelIdStreamConfig *ptr_channel_id_stream_cfg = NULL;
    if (ptr_channel_id_config->m_ptrMetadata == NULL)
    {
        LOG_ERROR("channel id config invalid");
        return GSE_ERROR;
    }
    uint32_t m_channelID = ptr_channel_id_config->m_ptrMetadata->m_channelID;
    std::size_t max_count = ptr_channel_id_config->m_streamTo.size();
    if (max_count <= 0)
    {
        LOG_ERROR("channel id config invalid, config json(%s)", ptr_channel_id_config->m_originData.c_str());
        return GSE_ERROR;
    }
    LOG_DEBUG("channelid :%d, channel count:%d, start create correlation channelid config", ptr_channel_id_config->m_ptrMetadata->m_channelID, max_count);
    for (std::size_t stream_to_idx = 0; stream_to_idx < max_count; ++stream_to_idx)
    {
        ChannelIdStreamConfig *ptr_channelid_config_tmp = new ChannelIdStreamConfig();
        StreamTo *ptr_stream_to = ptr_channel_id_config->m_streamTo.at(stream_to_idx);
        int stream_to_id = ptr_stream_to->m_streamToId;
        ptr_channelid_config_tmp->m_channelName = ptr_stream_to->m_name;
        ptr_channelid_config_tmp->m_streamToId = ptr_stream_to->m_streamToId;
        for (std::size_t idx = 0; idx < ptr_stream_to->m_filterNameAnd.size(); ++idx)
        {
            std::string filter_name = ptr_stream_to->m_filterNameAnd.at(idx);
            StreamFilter *ptr_stream_filter = ptr_channel_id_config->GetFilter(filter_name);
            if (NULL == ptr_stream_filter)
            {
                LOG_WARN("the filter name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(filter_name.c_str()), ptr_channel_id_config->m_ptrMetadata->m_channelID);
                continue;
            }
            LOG_DEBUG("channelid(%d) add Andfilter(%s)", ptr_channel_id_config->m_ptrMetadata->m_channelID, ptr_stream_filter->m_name.c_str());
            ptr_channelid_config_tmp->AddAndFilters(ptr_stream_filter);
        }
        for (std::size_t idx = 0; idx < ptr_stream_to->m_filterNameOr.size(); ++idx)
        {
            std::string filter_name = ptr_stream_to->m_filterNameOr.at(idx);
            StreamFilter *ptr_stream_filter = ptr_channel_id_config->GetFilter(filter_name);
            if (NULL == ptr_stream_filter)
            {
                LOG_WARN("the filter name (%s) is invalid, please to check the channel id (%u) config in zk", SAFE_CSTR(filter_name.c_str()), ptr_channel_id_config->m_ptrMetadata->m_channelID);
                continue;
            }
            LOG_DEBUG("channelid(%d) add or filter(%s)", ptr_channel_id_config->m_ptrMetadata->m_channelID, ptr_stream_filter->m_name.c_str());
            ptr_channelid_config_tmp->AddOrFilters(ptr_stream_filter);
        }
        ptr_channelid_config_tmp->CopyTableConfig(ptr_stream_to);
        ptr_channelid_config_tmp->m_channelID = m_channelID;
        std::string table_name;
        if (ptr_channel_id_stream_cfg == NULL)
        {
            ptr_channel_id_stream_cfg = ptr_channelid_config_tmp;
            ptr_channel_id_stream_cfg->GetTableName(table_name);
            LOG_DEBUG("Create channelId first stream config(id:%d, table name:%s, exporter id:%d, ptr:%p)",
                      ptr_channel_id_stream_cfg->m_channelID, table_name.c_str(), ptr_channel_id_stream_cfg->m_streamToId, ptr_channelid_config_tmp);
        }
        else
        {
            ptr_channel_id_stream_cfg->SetNext(ptr_channelid_config_tmp);
            ptr_channelid_config_tmp->GetTableName(table_name);
            LOG_DEBUG("Create ChannelId next stream config(id:%d, table name:%s, exporter id:%d, ptr:%p)",
                      ptr_channelid_config_tmp->m_channelID, table_name.c_str(), ptr_channelid_config_tmp->m_streamToId, ptr_channelid_config_tmp);
        }
    }

    m_rwLock.WLock();
    std::map<uint32_t, ChannelIdStreamConfig*>::iterator it = m_channelIds.find(ptr_channel_id_config->m_ptrMetadata->m_channelID);
    if (it != m_channelIds.end())
    {
        ChannelIdStreamConfig* ptr_channel_id_delete = NULL;
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
    return GSE_SUCCESS;
}
KafkaStorage::KafkaStorage()
{}
KafkaStorage::~KafkaStorage()
{}
std::string KafkaStorage::ToString()
{
    std::string kafka_string;
    std::vector<Address>::iterator it;
    for (it = m_addresses.begin(); it != m_addresses.end(); it++)
    {
        kafka_string.append((*it).m_ip);
        kafka_string.append(":");
        kafka_string.append(gse::tools::strings::ToString((*it).m_port));
        if (it != m_addresses.end())
        {
            kafka_string.append(";");
            break;
        }
    }
    kafka_string.append(", username:");
    kafka_string.append(m_kafkaConfig.m_saslUserName);
    kafka_string.append(", saslmechanisms:");
    kafka_string.append(m_kafkaConfig.m_saslMechanisms);
    return kafka_string;
}
}
}
