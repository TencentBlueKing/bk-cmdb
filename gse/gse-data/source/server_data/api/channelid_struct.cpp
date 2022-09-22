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
#include "channelid_struct.h"

#include <string>

#include "channelid_zk_api.h"
#include "conf/conf_common.h"
#include "safe_json.h"
#include <json/json.h>

namespace gse {
namespace data {

Address::Address()
    : m_ip(""), m_port(0) {}

std::string Address::ToString()
{
    std::string strbuff = "Host:{host:%s, port:%d}";
    std::string strResult;
    char buff[1024] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_ip.c_str(), m_port);
    strResult.append(buff);
    return strResult;
}

Json::Value Address::ToJSONValue()
{
    Json::Value meta;
    meta["ip"] = m_ip;
    meta["port"] = m_port;
    return meta;
}
std::string Address::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool Address::Parse(const Json::Value &jsonValue, ApiError &error)
{
    m_ip = SafeJson::GetString(jsonValue, "ip", "");
    m_port = SafeJson::GetInt(jsonValue, "port", 0);
    return true;
}

bool Address::Valid(ApiError &error)
{
    if (m_ip == "")
    {
        error = AddressInvalid;
        return false;
    }

    if (m_port == 0)
    {
        error = PortInvalid;
        return false;
    }

    return true;
}

Label::Label()
    : m_valid(false), m_bkBizId(-1)
{
}

Json::Value Label::ToJSONValue()
{
    Json::Value data;
    if (m_valid)
    {
        data["odm"] = Json::Value(m_odm);
        data["bk_biz_id"] = Json::Value(m_bkBizId);
        data["bk_biz_name"] = Json::Value(m_bkBizName);
    }

    return data;
}

std::string Label::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool Label::Parse(const Json::Value &jsonValue)
{
    m_valid = true;
    m_bkBizId = SafeJson::GetInt(jsonValue, "bk_biz_id", 0);
    m_bkBizName = SafeJson::GetString(jsonValue, "bk_biz_name", "");
    m_odm = SafeJson::GetString(jsonValue, "odm", "");
    return true;
}

Metadata::Metadata()
    : m_version(""), m_channelID(0), m_requestPlatName(""), m_isPlatId(false) {}

Json::Value Metadata::ToJSONValue()
{
    Json::Value meta;

    meta["version"] = Json::Value(m_version);
    meta["channel_id"] = Json::Value(m_channelID);
    meta["plat_name"] = Json::Value(m_requestPlatName);
    meta["label"] = m_label.ToJSONValue();
    return meta;
}

std::string Metadata::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool Metadata::ParseMetadata(const Json::Value &metaJson, ApiError &error)
{
    try
    {
        m_version = SafeJson::GetString(metaJson, "version", "");
        m_requestPlatName = SafeJson::GetString(metaJson, "plat_name", "");
        m_channelID = SafeJson::GetInt(metaJson, "channel_id", 0);

        if (metaJson.isMember("label"))
        {
            m_label.Parse(metaJson["label"]);
        }
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse metadata, error:%s", e.what());
        return false;
    }

    return true;
}

bool Metadata::Valid(ApiError &error)
{
    if (m_requestPlatName.empty())
    {
        error = PlatnameNotSet;
        return false;
    }

    std::string platid;
    std::string errorMsg;
    if (!ChannelIdZkApiInst()->ReadPlatNumber(m_requestPlatName, platid, errorMsg))
    {
        LOG_ERROR("plat_name(%s) invalid, error:%s", m_requestPlatName.c_str(), errorMsg.c_str());
        error = PlatnameNotSupport;
        return false;
    }

    return true;
}

KafkaTopic::KafkaTopic()
    : m_bizId(0), m_partition(1)
{
}

Json::Value KafkaTopic::ToJSONValue()
{
    Json::Value data;
    data["topic_name"] = m_topicName;
    data["data_set"] = m_dataSet;
    data["biz_id"] = m_bizId;
    data["partition"] = m_partition;
    return data;
}

std::string KafkaTopic::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool KafkaTopic::Parse(const Json::Value &kafkaJson, ApiError &error)
{
    m_topicName = SafeJson::GetString(kafkaJson, "topic_name", "");
    if (m_topicName == "")
    {
        error = TopNameInvalid;
        return false;
    }
    //兼容老版本临时存在
    m_dataSet = SafeJson::GetString(kafkaJson, "data_set", "");
    m_bizId = SafeJson::GetInt(kafkaJson, "biz_id", 0);
    m_partition = SafeJson::GetInt(kafkaJson, "partition", 1);

    return false;
}

uint32_t KafkaTopic::GetMemSize()
{
    uint32_t memsize = sizeof(KafkaTopic);
    memsize += m_dataSet.length();
    memsize += m_topicName.length();
    return memsize;
}
void KafkaTopic::CopyFrom(KafkaTopic *src)
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
        m_topicName = src->m_dataSet + gse::tools::strings::ToString(src->m_bizId);
    }
}

PulsarTopic::PulsarTopic()
    : m_bizId(0)
{
}

Json::Value PulsarTopic::ToJSONValue()
{
    Json::Value data;

    data["topic_name"] = m_topicName;
    data["tenant"] = m_tenant;
    data["namespace"] = m_namespace;
    data["data_set"] = m_dataSet;
    data["biz_id"] = Json::Value(m_bizId);
    data["persistent"] = m_persistent;
    return data;
}

std::string PulsarTopic::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool PulsarTopic::Parse(const Json::Value &pulsarJson, ApiError &error)
{
    m_topicName = SafeJson::GetString(pulsarJson, "topic_name", "");
    if (m_topicName == "")
    {
        error = TopNameInvalid;
        return false;
    }

    m_namespace = SafeJson::GetString(pulsarJson, "namespace", "");
    m_tenant = SafeJson::GetString(pulsarJson, "tenant", "");
    return true;
}

uint32_t PulsarTopic::GetMemSize()
{
    uint32_t memsize = sizeof(PulsarTopic);
    memsize += m_dataSet.length();
    memsize += m_topicName.length();
    memsize += m_tenant.length();
    memsize += m_namespace.length();
    memsize += m_persistent.length();
    return memsize;
}

void PulsarTopic::CopyFrom(PulsarTopic *src)
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

    std::string pulsarTopicName;
    if (!m_persistent.empty())
    {
        pulsarTopicName.append(m_persistent);
        pulsarTopicName.append("://");
    }
    else
    {
        pulsarTopicName.append("persistent");
        pulsarTopicName.append("://");
    }

    if (!m_tenant.empty())
    {
        pulsarTopicName.append(m_tenant);
        pulsarTopicName.append("/");
    }
    if (!m_namespace.empty())
    {
        pulsarTopicName.append(m_namespace);
        pulsarTopicName.append("/");
    }

    if (m_topicName == "")
    {
        m_topicName = m_dataSet + gse::tools::strings::ToString(m_bizId);
    }
    m_topicName = pulsarTopicName.append(m_topicName);
}

RedisChannel::RedisChannel()
    : m_bizId(0) {}

Json::Value RedisChannel::ToJSONValue()
{
    Json::Value data;

    data["channel_name"] = m_channelName;
    data["data_set"] = m_dataSet;
    data["biz_id"] = m_bizId;
    return data;
}

uint32_t RedisChannel::GetMemSize()
{
    uint32_t memsize = sizeof(RedisChannel);
    memsize += m_channelName.length();
    memsize += m_dataSet.length();
    return memsize;
}

std::string RedisChannel::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool RedisChannel::Parse(const Json::Value &redisJson, ApiError &error)
{
    m_channelName = SafeJson::GetString(redisJson, "channel_name", "");
    if (m_channelName.empty())
    {
        error = ChannelNameInvalid;
        return false;
    }

    m_dataSet = SafeJson::GetString(redisJson, "data_set", "");
    m_bizId = SafeJson::GetInt(redisJson, "biz_id", 0);

    return true;
}

void RedisChannel::CopyFrom(RedisChannel *src)
{
    m_channelName = src->m_channelName;
    m_dataSet = src->m_dataSet;
    m_bizId = src->m_bizId;
}

StreamTo::StreamTo()
    : m_streamToId(0) {}

StreamTo *StreamTo::Clone()
{
    StreamTo *streamTo = new StreamTo();
    streamTo->CopyFrom(this);
    return streamTo;
}

void StreamTo::CopyFrom(StreamTo *src)
{
    m_reportMode = src->m_reportMode;
    m_streamToId = src->m_streamToId;

    if (src->m_reportMode == ChannelIDOperationRequestMethodKafka)
    {
        m_kafkaTopic.m_bizId = src->m_kafkaTopic.m_bizId;
        m_kafkaTopic.m_dataSet = src->m_kafkaTopic.m_dataSet;
        m_kafkaTopic.m_topicName = src->m_kafkaTopic.m_topicName;
        m_kafkaTopic.m_partition = src->m_kafkaTopic.m_partition;
    }
    else if (src->m_reportMode == ChannelIDOperationRequestMethodRedis)
    {
        m_redisChannel.m_bizId = src->m_redisChannel.m_bizId;
        m_redisChannel.m_dataSet = src->m_redisChannel.m_dataSet;
        m_redisChannel.m_channelName = src->m_redisChannel.m_channelName;
    }
    else if (src->m_reportMode == ChannelIDOperationRequestMethodPulsar)
    {
        m_pulsarTopic.m_topicName = src->m_pulsarTopic.m_topicName;
        m_pulsarTopic.m_tenant = src->m_pulsarTopic.m_tenant;
        m_pulsarTopic.m_namespace = src->m_pulsarTopic.m_namespace;
    }
}

//这个完整的json
Json::Value StreamTo::ToJSONValue()
{
    Json::Value data;
    Json::Value stream_to_json;
    LOG_DEBUG("StreamTo to json value, report mode:%s", m_reportMode.c_str());
    if (m_reportMode == ChannelIDOperationRequestMethodKafka)
    {
        stream_to_json[ChannelIDOperationRequestMethodKafka] = m_kafkaTopic.ToJSONValue();
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodRedis)
    {
        stream_to_json[ChannelIDOperationRequestMethodRedis] = m_redisChannel.ToJSONValue();
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodPulsar)
    {
        stream_to_json[ChannelIDOperationRequestMethodPulsar] = m_pulsarTopic.ToJSONValue();
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodDsProxy)
    {
        stream_to_json[ChannelIDOperationRequestMethodDsProxy] = Json::Value(Json::objectValue);
    }

    stream_to_json["stream_to_id"] = Json::Value(m_streamToId);

    return stream_to_json;
}

std::string StreamTo::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool StreamTo::Parse(const Json::Value &channel, ApiError &error)
{
    bool parseResult = false;

    m_streamToId = SafeJson::GetInt(channel, "stream_to_id", 0);

    if (channel.isMember(ChannelIDOperationRequestMethodKafka))
    {
        m_reportMode = ChannelIDOperationRequestMethodKafka;
        parseResult = m_kafkaTopic.Parse(channel[ChannelIDOperationRequestMethodKafka], error);
    }
    else if (channel.isMember(ChannelIDOperationRequestMethodPulsar))
    {
        m_reportMode = ChannelIDOperationRequestMethodPulsar;
        parseResult = m_pulsarTopic.Parse(channel[ChannelIDOperationRequestMethodPulsar], error);
    }
    else if (channel.isMember(ChannelIDOperationRequestMethodRedis))
    {
        m_reportMode = ChannelIDOperationRequestMethodRedis;
        parseResult = m_redisChannel.Parse(channel[ChannelIDOperationRequestMethodRedis], error);
    }
    else if (channel.isMember(ChannelIDOperationRequestMethodDsProxy))
    {
        m_reportMode = ChannelIDOperationRequestMethodDsProxy;

        parseResult = true;
    }
    else
    {
        LOG_ERROR("failed to parse streamto, unkown mode, json(%s)", channel.toStyledString().c_str());
        return false;
    }

    return parseResult;
}

bool StreamTo::Valid(ApiError &error)
{
    if (m_reportMode.compare(ChannelIDOperationRequestMethodKafka) == 0)
    {
        if (m_kafkaTopic.m_topicName.empty())
        {
            error = TopNameInvalid;
            LOG_ERROR("stream_to/kafka/topic_name is not seted");
            return false;
        }
    }
    else if (m_reportMode.compare(ChannelIDOperationRequestMethodPulsar) == 0)
    {
        if (m_pulsarTopic.m_topicName.empty())
        {
            error = TopNameInvalid;
            LOG_ERROR("stream_to/pulsar/topic_name is not seted");
            return false;
        }
    }
    else if (m_reportMode.compare(ChannelIDOperationRequestMethodRedis) == 0)
    {
        if (m_redisChannel.m_channelName.empty())
        {
            error = ChannelNameInvalid;
            LOG_ERROR("stream_to/redis/channel_name is not seted");
            return false;
        }
    }
    else if (m_reportMode.compare(ChannelIDOperationRequestMethodDsProxy) == 0)
    {
    }
    else
    {
        // 代码走到这里，说明 m_reportMode 是无效的
        error = UnsupportReportMode;
        LOG_ERROR("the stream_configs/report_mode(%s) is invalid only be seted(file or kafka)", m_reportMode.c_str());
        return false;
    }

    if (!ChannelIdZkApiInst()->ExistStreamToClusterId(gse::tools::strings::ToString(m_streamToId)))
    {
        error = StreamToIdNotExist;
        LOG_ERROR("stream_to_id %d not exist, error:%s", m_streamToId);
        return false;
    }

    return true;
}

bool Channel::Parse(const Json::Value &routeJson, ApiError &error)
{
    try
    {
        m_name = SafeJson::GetString(routeJson, "name", "");

        if (!routeJson.isMember("stream_to"))
        {
            error = StreamToNotSet;
            return false;
        }

        m_streamTo.Parse(routeJson["stream_to"], error);

        if (routeJson.isMember("filter_name_and"))
        {
            Json::Value filterNameAnd = routeJson["filter_name_and"];
            if (!filterNameAnd.isArray())
            {
                error = FilterNotARRAY;
                return false;
            }

            Json::ArrayIndex filterNameAndSize = filterNameAnd.size();
            for (Json::ArrayIndex idx = 0; idx != filterNameAndSize; ++idx)
            {
                if (!filterNameAnd[idx].isString())
                {
                    error = FilterNameNotString;
                    return false;
                }
                m_filterAndName.push_back(filterNameAnd[idx].asString());
            }
        }

        if (routeJson.isMember("filter_name_or"))
        {
            Json::Value filterNameOr = routeJson["filter_name_or"];
            if (!filterNameOr.isArray())
            {
                error = FilterNotARRAY;
                return false;
            }

            Json::ArrayIndex filter_name_or_size = filterNameOr.size();
            for (Json::ArrayIndex idx = 0; idx != filter_name_or_size; ++idx)
            {
                if (!filterNameOr[idx].isString())
                {
                    error = FilterNameNotString;
                    return false;
                }

                m_filterOrName.push_back(filterNameOr[idx].asString());
            }
        }
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse route, error:%s", e.what());
        return false;
    }

    return true;
}

Json::Value Channel::ToJSONValue()
{
    Json::Value data;
    data["stream_to"] = m_streamTo.ToJSONValue();
    data["name"] = m_name;

    data["filter_name_and"].resize(0);
    data["filter_name_or"].resize(0);

    Json::Value filters;
    filters.resize(0);
    for (auto filterName : m_filterAndName)
    {
        filters.append(Json::Value(filterName));
    }

    data["filter_name_and"] = filters;

    filters.resize(0);
    for (auto filterName : m_filterOrName)
    {
        filters.append(Json::Value(filterName));
    }

    data["filter_name_or"] = filters;
    return data;
}

bool Channel::Valid(ApiError &error)
{
    if (m_name.empty())
    {
        error = ChannelRouteNameNotSet;
        return false;
    }

    if (!m_streamTo.Valid(error))
    {
        LOG_ERROR("failed to check stream, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    return true;
}

void Channel::CopyFrom(Channel *src)
{
    m_name = src->m_name;

    m_streamTo.CopyFrom(&src->m_streamTo);

    m_filterAndName.clear();
    m_filterOrName.clear();

    for (auto itFilterAndName : src->m_filterAndName)
    {
        m_filterAndName.push_back(itFilterAndName);
    }

    for (auto itFilterOrName : src->m_filterOrName)
    {
        m_filterOrName.push_back(itFilterOrName);
    }
}

std::shared_ptr<Channel> Channel::Clone()
{
    auto channel = std::make_shared<Channel>();
    channel->CopyFrom(this);
    return channel;
}

StreamFilter::StreamFilter()
    : m_fieldIndex(0)
{
}

std::shared_ptr<StreamFilter> StreamFilter::Clone()
{
    auto ptrStreamFilter = std::make_shared<StreamFilter>();
    ptrStreamFilter->CopyFrom(this);
    return ptrStreamFilter;
}

void StreamFilter::CopyFrom(StreamFilter *src)
{
    m_name = src->m_name;
    m_fieldIndex = src->m_fieldIndex;
    m_fieldDataType = src->m_fieldDataType;
    m_fieldDataValue = src->m_fieldDataValue;
    m_separator = src->m_separator;
    m_fieldIn = src->m_fieldIn;
}

Json::Value StreamFilter::ToJSONValue()
{
    Json::Value data;

    data["name"] = m_name;
    data["field_index"] = m_fieldIndex;
    data["field_data_type"] = m_fieldDataType;
    data["field_data_value"] = m_fieldDataValue;
    data["field_separator"] = m_separator;
    data["field_in"] = m_fieldIn;

    return data;
}

std::string StreamFilter::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool StreamFilter::Parse(const Json::Value &streamFilterJson, ApiError &error)
{
    try
    {
        m_name = SafeJson::GetString(streamFilterJson, "name", "");

        if (m_name.empty())
        {
            error = FilterNameInvalid;
            return false;
        }

        m_fieldIndex = SafeJson::GetInt(streamFilterJson, "field_index", 0);
        m_fieldDataType = SafeJson::GetString(streamFilterJson, "field_data_type", "");
        m_fieldDataValue = SafeJson::GetString(streamFilterJson, "field_data_value", "");
        m_separator = SafeJson::GetString(streamFilterJson, "field_separator", "");
        m_fieldIn = SafeJson::GetString(streamFilterJson, "field_in", "");
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse streamfilter json, error:%s", e.what());
        error = ReadFromConfigJsonInvalid;
        return false;
    }
    return true;
}

bool StreamFilter::Valid(ApiError &error)
{
    std::string errorMsg;
    if (m_name.empty())
    {
        error = StreamFilterNameNotSet;
        LOG_ERROR("the stream_filters/name is not seted, for the stream_filter(%s)", m_name.c_str());
        return false;
    }

    if (m_fieldIndex < 0)
    {
        error = StreamFilterIndexInvalid;
        LOG_ERROR("index(%d) invalid for the stream_filter(%s)", m_fieldIndex, m_name.c_str());
        return false;
    }

    if (m_fieldIn.compare("protocol") != 0 && m_fieldIn.compare("data") != 0)
    {
        errorMsg = "stream_filters/field_in(" + m_fieldIn + ") is invalid, only be seted (protocol or data)";
        LOG_ERROR("%s, for the stream_filter(%s)", errorMsg.c_str(), m_name.c_str());
        return false;
    }

    if (m_fieldDataType.compare("int") == 0)
    {
        if (!gse::tools::strings::IsNumber(m_fieldDataValue))
        {
            error = StreamFilterIndexInvalid;
            errorMsg = "stream_filters/field_data_value(" + m_fieldDataValue + ") should be number type";
            LOG_ERROR("%s, for the stream_filter(%s)", errorMsg.c_str(), m_name.c_str());
            return false;
        }
    }

    if (m_fieldDataType.compare("byte") == 0)
    {
        if (!gse::tools::strings::IsNumber(m_fieldDataValue))
        {
            error = StreamFilterIndexInvalid;
            errorMsg = "stream_filters/field_data_value(" + m_fieldDataValue + ") should be byte type";
            LOG_ERROR("%s, for the stream_filter(%s)", errorMsg.c_str(), m_name.c_str());
            return false;
        }
    }

    if (m_fieldDataType.compare("string") != 0)
    {
        error = StreamFilterIndexInvalid;
        errorMsg = "stream_filters/field_data_type(" + m_fieldDataType + ") is invalid, only be seted (int、string、byte)";
        LOG_ERROR("%s, for the stream_filter(%s)", errorMsg.c_str(), m_name.c_str());
        return false;
    }

    return true;
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

uint32_t StreamFilter::GetMemSize()
{
    uint32_t memSize = sizeof(StreamFilter);
    memSize += m_name.length();
    memSize += m_fieldDataType.length();
    memSize += m_fieldDataValue.length();
    memSize += m_separator.length();
    memSize += m_fieldIn.length();
    return memSize;
}

StreamToCluster::StreamToCluster()
    : m_clusterType(KAFKA_COMMON)
{
}

StreamToCluster::~StreamToCluster()
{
}

void StreamToCluster::CopyFrom(StreamToCluster *src)
{
    m_name = src->m_name;
    m_reportMode = src->m_reportMode;
    m_clusterType = src->m_clusterType;

    if (src->m_reportMode == ChannelIDOperationRequestMethodKafka)
    {
        if (m_kafkaCluster == nullptr)
        {
            m_kafkaCluster = std::make_shared<KafkaCluster>();
        }
        *m_kafkaCluster = *src->m_kafkaCluster;
    }
    else if (src->m_reportMode == ChannelIDOperationRequestMethodRedis)
    {
        if (m_redisCluster == nullptr)
        {
            m_redisCluster = std::make_shared<RedisCluster>();
        }
        *m_redisCluster = *src->m_redisCluster;
    }
    else if (src->m_reportMode == ChannelIDOperationRequestMethodPulsar)
    {
        if (m_pulsarCluster == nullptr)
        {
            m_pulsarCluster = std::make_shared<PulsarCluster>();
        }
        *m_pulsarCluster = *src->m_pulsarCluster;
    }
    else if (src->m_reportMode == ChannelIDOperationRequestMethodDsProxy)
    {
        if (m_dsproxy == nullptr)
        {
            m_dsproxy = std::make_shared<DSProxyExporter>();
        }
        *m_dsproxy = *src->m_dsproxy;
    }
}

Json::Value StreamToCluster::ToJSONValue()
{
    Json::Value data;

    data["name"] = Json::Value(m_name);
    data["report_mode"] = Json::Value(m_reportMode);

    if (m_reportMode == ChannelIDOperationRequestMethodKafka)
    {
        data[ChannelIDOperationRequestMethodKafka] = m_kafkaCluster->ToJSONValue();
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodPulsar)
    {
        data[ChannelIDOperationRequestMethodPulsar] = m_pulsarCluster->ToJSONValue();
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodRedis)
    {
        data[ChannelIDOperationRequestMethodRedis] = m_redisCluster->ToJSONValue();
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodDsProxy)
    {
        data[ChannelIDOperationRequestMethodDsProxy] = m_dsproxy->ToJSONValue();
    }
    return data;
}

std::string StreamToCluster::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

StreamToCluster *StreamToCluster::Clone()
{
    StreamToCluster *ptr_tmp = new StreamToCluster();
    ptr_tmp->CopyFrom(this);
    return ptr_tmp;
}

bool StreamToCluster::Parse(const Json::Value &jsonValue, ApiError &error)
{
    m_name = SafeJson::GetString(jsonValue, "name", "");
    m_reportMode = SafeJson::GetString(jsonValue, "report_mode", "");

    if (m_reportMode == ChannelIDOperationRequestMethodKafka)
    {
        if (!jsonValue.isMember(ChannelIDOperationRequestMethodKafka))
        {
            error = KafkaNotSet;
            LOG_ERROR("failed to parse kafka json, no kafka object");
            return false;
        }
        m_kafkaCluster = std::make_shared<KafkaCluster>();
        if (!m_kafkaCluster->Parse(jsonValue[ChannelIDOperationRequestMethodKafka], error))
        {
            return false;
        }

        m_clusterType = KAFKA_COMMON;
        return true;
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodPulsar)
    {
        if (!jsonValue.isMember(ChannelIDOperationRequestMethodPulsar))
        {
            error = PulsarNotSet;
            LOG_ERROR("failed to parse pulsar json, no pulsar object");
            return false;
        }

        m_pulsarCluster = std::make_shared<PulsarCluster>();
        if (!m_pulsarCluster->Parse(jsonValue[ChannelIDOperationRequestMethodPulsar], error))
        {
            return false;
        }

        m_clusterType = EXPORT_PULSAR;
        return true;
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodRedis)
    {
        if (!jsonValue.isMember(ChannelIDOperationRequestMethodRedis))
        {
            error = RedisNotSet;
            return false;
        }

        m_redisCluster = std::make_shared<RedisCluster>();
        if (!m_redisCluster->Parse(jsonValue[ChannelIDOperationRequestMethodRedis], error))
        {
            return false;
        }

        if (m_redisCluster->m_mode == kRedisSentinel)
        {
            m_clusterType = REDIS_SENTINEL_PUB;
        }
        else if (m_redisCluster->m_mode == kRedisSingle)
        {
            m_clusterType = REDIS_PUB;
        }
        else
        {
            error = UnsupportRedisClusterMode;
            return false;
        }

        return true;
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodDsProxy)
    {
        if (!jsonValue.isMember(ChannelIDOperationRequestMethodDsProxy))
        {
            error = ProxyNotSet;
            return false;
        }

        m_dsproxy = std::make_shared<DSProxyExporter>();
        if (!m_dsproxy->Parse(jsonValue[ChannelIDOperationRequestMethodDsProxy], error))
        {
            return false;
        }
        m_clusterType = EXPORT_DSPROXY;
    }
    else
    {
        error = UnsupportReportMode;
        LOG_ERROR("failed to parse stream_to, unkown report_mode, json(%s)", jsonValue.toStyledString().c_str());
        return false;
    }

    return true;
}

bool StreamToCluster::Valid(ApiError &error)
{
    if (!ChannelIdCommApiInst()->IsSupportClusterType(m_reportMode))
    {
        error = UnsupportReportMode;
        return false;
    }

    if (m_reportMode == ChannelIDOperationRequestMethodKafka)
    {
        return m_kafkaCluster->Valid(error);
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodPulsar)
    {
        return m_pulsarCluster->Valid(error);
    }

    else if (m_reportMode == ChannelIDOperationRequestMethodRedis)
    {
        return m_redisCluster->Valid(error);
    }
    else if (m_reportMode == ChannelIDOperationRequestMethodDsProxy)
    {
        return m_dsproxy->Valid(error);
    }
    else
    {
        error = UnsupportReportMode;
        return false;
    }

    return false;
}

KafkaCluster::KafkaCluster()
{
}
KafkaCluster::~KafkaCluster()
{
}
std::string KafkaCluster::ToString()
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
    if (m_kafkaConfig.m_saslUserName != "")
    {
        kafka_string.append(", username:");
        kafka_string.append(m_kafkaConfig.m_saslUserName);
    }

    if (m_kafkaConfig.m_saslMechanisms != "")
    {
        kafka_string.append(", saslmechanisms:");
        kafka_string.append(m_kafkaConfig.m_saslMechanisms);
    }

    return kafka_string;
}

Json::Value KafkaCluster::ToJSONValue()
{
    Json::Value data;
    Json::Value addresses;
    for (auto address : m_addresses)
    {
        addresses.append(address.ToJSONValue());
    }

    data["storage_address"] = addresses;

    if (m_kafkaConfig.m_securityProtocol != "")
    {
        data["security_protocol"] = Json::Value(m_kafkaConfig.m_securityProtocol);
    }

    if (m_kafkaConfig.m_saslMechanisms != "")
    {
        data["sasl_mechanisms"] = Json::Value(m_kafkaConfig.m_saslMechanisms);
    }

    if (m_kafkaConfig.m_saslUserName != "")
    {
        data["sasl_username"] = Json::Value(m_kafkaConfig.m_saslUserName);
    }

    if (m_kafkaConfig.m_saslPasswd != "")
    {
        data["sasl_passwd"] = Json::Value(m_kafkaConfig.m_saslPasswd);
    }

    LOG_DEBUG("kafka cluster json:%s", data.toStyledString().c_str());
    return data;
}

std::string KafkaCluster::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool KafkaCluster::Valid(ApiError &error)
{
    if (m_addresses.size() == 0)
    {
        error = AddressNotExist;
        return false;
    }

    for (auto addr : m_addresses)
    {
        if (!addr.Valid(error))
        {
            return false;
        }
    }

    return true;
}

bool KafkaCluster::Parse(const Json::Value &jsonValue, ApiError &error)
{
    try
    {
        if (!jsonValue.isMember("storage_address"))
        {
            LOG_ERROR("failed to parse kafka cluster, need storage_address");
            error = AddressNotExist;
            return false;
        }

        if (!jsonValue["storage_address"].isArray())
        {
            LOG_ERROR("failed to parse kafka cluster, storage_address not array");
            error = AddressNotARRAY;
            return false;
        }

        for (int i = 0; i < jsonValue["storage_address"].size(); i++)
        {
            Address address;
            if (!address.Parse(jsonValue["storage_address"][i], error))
            {
                return false;
            }
            m_addresses.push_back(address);
        }

        m_kafkaConfig.m_queueBufferingMaxMessages = SafeJson::GetInt(jsonValue, "queue_buffering_max_messages", DEFAULT_MAX_KAFKA_QUEUE_SIZE);
        m_kafkaConfig.m_messageMaxBytes = SafeJson::GetInt(jsonValue, "message_max_bytes", DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE);
        m_kafkaConfig.m_securityProtocol = SafeJson::GetString(jsonValue, "security_protocol", "");
        m_kafkaConfig.m_saslMechanisms = SafeJson::GetString(jsonValue, "sasl_mechanisms", "");
        m_kafkaConfig.m_saslUserName = SafeJson::GetString(jsonValue, "sasl_username", "");
        m_kafkaConfig.m_saslPasswd = SafeJson::GetString(jsonValue, "sasl_passwd", "");
        m_kafkaConfig.m_requestRequiredAcks = SafeJson::GetString(jsonValue, "request_required_acks", "1");
        m_kafkaConfig.m_queueBufferingMaxMs = SafeJson::GetInt(jsonValue, "queue_buffering_max_ms", -1);
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse kafka, error:%s", e.what());
        return false;
    }

    return true;
}

PulsarCluster::PulsarCluster()
{
}

PulsarCluster::~PulsarCluster()
{
}

bool PulsarCluster::Parse(const Json::Value &jsonValue, ApiError &error)
{
    try
    {
        if (!jsonValue.isMember("storage_address"))
        {
            LOG_ERROR("failed to parse pulsar cluster, need storage_address");
            error = AddressNotExist;
            return false;
        }

        if (!jsonValue["storage_address"].isArray())
        {
            LOG_ERROR("failed to parse kafka cluster, storage_address not array");
            error = AddressNotARRAY;
            return false;
        }

        for (int i = 0; i < jsonValue["storage_address"].size(); i++)
        {
            Address address;
            if (!address.Parse(jsonValue["storage_address"][i], error))
            {
                return false;
            }
            m_addresses.push_back(address);
        }

        m_token = SafeJson::GetString(jsonValue, "token", "");
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse pulsar, error:%s", e.what());
        return false;
    }

    return true;
}

bool PulsarCluster::Valid(ApiError &error)
{
    if (m_addresses.size() == 0)
    {
        error = AddressNotExist;
        return false;
    }

    for (auto addr : m_addresses)
    {
        if (!addr.Valid(error))
        {
            return false;
        }
    }

    return true;
}

std::string PulsarCluster::ToString()
{
    std::string strPulsarStorageAddr;
    std::vector<Address>::iterator it;
    for (it = m_addresses.begin(); it != m_addresses.end(); it++)
    {
        strPulsarStorageAddr.append((*it).m_ip);
        strPulsarStorageAddr.append(":");
        strPulsarStorageAddr.append(gse::tools::strings::ToString((*it).m_port));
        if (it != m_addresses.end())
        {
            strPulsarStorageAddr.append(";");
            break;
        }
    }
    return strPulsarStorageAddr;
}

Json::Value PulsarCluster::ToJSONValue()
{
    Json::Value data;

    Json::Value addresses;
    for (auto address : m_addresses)
    {
        addresses.append(address.ToJSONValue());
    }

    data["storage_address"] = addresses;
    data["token"] = Json::Value(m_token);
    return data;
}

std::string PulsarCluster::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool RedisCluster::Parse(const Json::Value &jsonValue, ApiError &error)
{
    try
    {
        if (!jsonValue.isMember("storage_address"))
        {
            LOG_ERROR("failed to parse kafka cluster, need storage_address");
            error = AddressNotExist;
            return false;
        }

        if (!jsonValue["storage_address"].isArray())
        {
            LOG_ERROR("failed to parse kafka cluster, storage_address not array");
            error = AddressNotARRAY;
            return false;
        }

        for (int i = 0; i < jsonValue["storage_address"].size(); i++)
        {
            Address address;
            if (!address.Parse(jsonValue["storage_address"][i], error))
            {
                return false;
            }
            m_addresses.push_back(address);
        }

        m_masterName = SafeJson::GetString(jsonValue, "master_name", "");
        m_passwd = SafeJson::GetString(jsonValue, "passwd", "");
        m_sentinelPasswd = SafeJson::GetString(jsonValue, "sentinel_passwd", "");
        m_mode = SafeJson::GetString(jsonValue, "mode", "sentinel");
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse redis,json:%s, error:%s", jsonValue.toStyledString().c_str(), e.what());
        return false;
    }
    catch (...)
    {
        LOG_ERROR("failed to parse redis,json:%s, unkonw exception", jsonValue.toStyledString().c_str());
        return false;
    }

    return true;
}

std::string RedisCluster::ToString()
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

Json::Value RedisCluster::ToJSONValue()
{
    Json::Value data;
    Json::Value addresses;
    for (auto address : m_addresses)
    {
        addresses.append(address.ToJSONValue());
    }

    data["storage_address"] = addresses;
    data["passwd"] = Json::Value(m_passwd);
    data["master_name"] = Json::Value(m_masterName);
    data["sentinel_passwd"] = Json::Value(m_sentinelPasswd);
    data["mode"] = Json::Value(m_mode);
    return data;
}

std::string RedisCluster::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool RedisCluster::Valid(ApiError &error)
{
    if (m_addresses.size() == 0)
    {
        error = AddressNotExist;
        return false;
    }

    for (auto addr : m_addresses)
    {
        if (!addr.Valid(error))
        {
            return false;
        }
    }

    if (!(m_mode == kRedisSentinel || m_mode == kRedisSingle))
    {
        error = UnsupportRedisClusterMode;
        return false;
    }

    return true;
}

Json::Value DSProxyExporter::ToJSONValue()
{
    Json::Value data;
    Json::Value addresses;
    for (auto address : m_addresses)
    {
        addresses.append(address.ToJSONValue());
    }

    data["storage_address"] = addresses;
    data["http_request_uri"] = Json::Value(m_httpURI);
    data["certpath"] = Json::Value(m_certPath);
    data["proxyprotocol"] = Json::Value(m_proxyProtocol);
    data["proxyversion"] = Json::Value(m_proxyVersion);
    data["connectionnum"] = Json::Value(m_connectionNumEachAddress);
    data["heartbeat"] = Json::Value(m_heartbeat);
    data["fillchannelid"] = Json::Value(m_fillChannelid);
    data["thirdparty_keyfile"] = Json::Value(m_thirdPartyKeyFile);
    data["thirdparty_cert_passwd"] = Json::Value(m_thirdPartyCertPasswd);
    data["thirdparty_certfile"] = Json::Value(m_thirdPartyCertFile);
    data["is_thirdparty_cert"] = Json::Value(m_isThirdPartyCert);
    return data;
}

std::string DSProxyExporter::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

bool DSProxyExporter::Valid(ApiError &error)
{
    if (m_addresses.size() == 0)
    {
        error = AddressNotExist;
        return false;
    }

    for (auto addr : m_addresses)
    {
        if (!addr.Valid(error))
        {
            return false;
        }
    }

    return true;
}

bool DSProxyExporter::Parse(const Json::Value &jsValue, ApiError &error)
{
    if (!jsValue.isMember("storage_address"))
    {
        error = AddressNotExist;
        LOG_ERROR("failed to parse stream to config, storage_address not set, input json(%s)", jsValue.toStyledString().c_str());
        return false;
    }

    if (!jsValue["storage_address"].isArray())
    {
        LOG_ERROR("failed to parse proxy cluster, storage_address not array");
        error = AddressNotARRAY;
        return false;
    }

    for (int i = 0; i < jsValue["storage_address"].size(); i++)
    {
        Address address;
        if (!address.Parse(jsValue["storage_address"][i], error))
        {
            return false;
        }
        m_addresses.push_back(address);
    }

    m_certPath = SafeJson::GetString(jsValue, "certpath", "");

    m_proxyProtocol = SafeJson::GetString(jsValue, "proxyprotocol", "tcp");
    m_proxyVersion = SafeJson::GetString(jsValue, "proxyversion", "v1");
    m_connectionNumEachAddress = SafeJson::GetInt(jsValue, "connectionnum", 2);
    m_heartbeat = SafeJson::GetBool(jsValue, "fillchannelid", true);
    m_fillChannelid = SafeJson::GetBool(jsValue, "fillchannelid", true);

    m_httpURI = SafeJson::GetBool(jsValue, "http_request_uri", "");

    m_thirdPartyCertPasswd = SafeJson::GetString(jsValue, "thirdparty_cert_passwd", "");
    m_thirdPartyCertFile = SafeJson::GetString(jsValue, "thirdparty_certfile", "");
    m_thirdPartyKeyFile = SafeJson::GetString(jsValue, "thirdparty_keyfile", "");
    m_isThirdPartyCert = SafeJson::GetBool(jsValue, "is_thirdparty_cert", false);

    LOG_DEBUG("finish to parse dsproxy config(%s)", ToJSON().c_str());
    return true;
}

Json::Value DsProxy::ToJSONValue()
{
    Json::Value data;
    Json::Value addresses;
    for (auto address : m_storagAddress)
    {
        addresses.append(address.ToJSONValue());
    }

    data["storage_address"] = addresses;
    data["http_request_uri"] = Json::Value(m_httpURI);
    data["certpath"] = Json::Value(m_certPath);
    data["proxyprotocol"] = Json::Value(m_proxyProtocol);
    data["proxyversion"] = Json::Value(m_proxyVersion);
    data["connectionnum"] = Json::Value(m_connectionNumEachAddress);
    data["heartbeat"] = Json::Value(m_heartbeat);
    data["fillchannelid"] = Json::Value(m_fillChannelid);
    data["thirdparty_keyfile"] = Json::Value(m_thirdPartyKeyFile);
    data["thirdparty_cert_passwd"] = Json::Value(m_thirdPartyCertPasswd);
    data["thirdparty_certfile"] = Json::Value(m_thirdPartyCertFile);
    data["is_thirdparty_cert"] = Json::Value(m_isThirdPartyCert);
    return data;
}

std::string DsProxy::ToJSON()
{
    Json::FastWriter writer;
    return writer.write(ToJSONValue());
}

} // namespace data
} // namespace gse
