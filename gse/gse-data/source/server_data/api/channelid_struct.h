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

#ifndef _GSE_CHANNEL_ID_STRUCT_H_
#define _GSE_CHANNEL_ID_STRUCT_H_

#include <list>
#include <memory>
#include <string>
#include <vector>

#include <event2/event.h>
#include <json/json.h>

#include "conf/conf_common.h"
#include "datacell.h"
#include "error_code.h"
#include "safe/lock.h"
#include "tools/thread.h"
#include "tools/strings.h"

namespace gse {
namespace data {

const std::string kRedisSentinel = "sentinel";
const std::string kRedisSingle = "single";
const std::string kRedisCluster = "cluster";

#define DEFAULT_MAX_KAFKA_QUEUE_SIZE 200000
#define DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE 10000000
enum ClusterType
{
    UNKNOWN = 0,            // UNKONW
    KAFKA_COMMON = 1,       // sys:1, type:2
    KAFKA_OP = 2,           // sys:1, type:2
    REDIS_SENTINEL_PUB = 3, // sys:3
    REDIS_PUB = 4,          // sys:4
    EXPORT_FILE = 5,        // only suport for channelid,
    EXPORT_PULSAR = 7,
    EXPORT_DSPROXY = 6 // only suport for channelid

};

class Address
{
public:
    Address();

    std::string ToString();
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &jsonValue, ApiError &error);
    bool Valid(ApiError &error);

public:
    std::string m_ip;
    uint16_t m_port;
};

class Label
{
public:
    Label();
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &jsonValue);

public:
    bool m_valid;
    int m_bkBizId;
    std::string m_bkBizName;
    std::string m_odm;
};

class Metadata
{

public:
    Metadata();
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool ParseMetadata(const Json::Value &metaJson, ApiError &error);
    bool Valid(ApiError &error);

public:
    std::string m_version;
    std::string m_requestPlatName;
    Label m_label;
    uint32_t m_channelID;
    bool m_isPlatId;
};

class KafkaTopic
{
public:
    KafkaTopic();

    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &kafkaJson, ApiError &error);
    void CopyFrom(KafkaTopic *src);
    uint32_t GetMemSize();

public:
    std::string m_topicName;
    std::string m_dataSet;
    int m_bizId;
    int m_partition;
};

//"biz_id":0,"data_set":"0bkmonitor_1500145","partition":1,"topic_name":"0bkmonitor_15001450"

class PulsarTopic
{
public:
    PulsarTopic();

    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &pulsarJson, ApiError &error);
    void CopyFrom(PulsarTopic *src);
    uint32_t GetMemSize();

public:
    std::string m_dataSet;
    std::string m_topicName;
    std::string m_tenant;
    std::string m_namespace;
    std::string m_persistent;
    int m_bizId;
};

class RedisChannel
{
public:
    RedisChannel();

    Json::Value ToJSONValue();

    std::string ToJSON();
    bool Parse(const Json::Value &redisJson, ApiError &error);
    void CopyFrom(RedisChannel *src);
    uint32_t GetMemSize();

public:
    std::string m_channelName;
    std::string m_dataSet;
    int m_bizId;
};

class DsProxyTable
{
public:
    std::string m_name;
    DsProxyTable()
        : m_name("") {}
};

class StreamTo
{
public:
    StreamTo();
    StreamTo *Clone();
    void CopyFrom(StreamTo *src);
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &channel, ApiError &error);
    bool Valid(ApiError &error);

public:
    uint32_t m_streamToId;
    std::string m_reportMode;
    KafkaTopic m_kafkaTopic;
    PulsarTopic m_pulsarTopic;
    RedisChannel m_redisChannel;
    DsProxyTable m_proxy;
    std::string m_logfilePath;
};

class Channel
{
public:
    bool Parse(const Json::Value &routeJson, ApiError &error);
    bool Valid(ApiError &error);
    Json::Value ToJSONValue();

    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    std::shared_ptr<Channel> Clone();
    void CopyFrom(Channel *src);

public:
    std::string m_name;
    StreamTo m_streamTo;
    uint32_t m_streamToId;
    std::vector<std::string> m_filterAndName;
    std::vector<std::string> m_filterOrName;
};

class StreamFilter
{
public:
    StreamFilter();
    std::shared_ptr<StreamFilter> Clone();
    void CopyFrom(StreamFilter *src);
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &streamFilterJson, ApiError &error);
    bool Valid(ApiError &error);
    std::string ToString();

    uint32_t GetMemSize();

public:
    int16_t m_fieldIndex;
    std::string m_name;
    std::string m_fieldDataType;
    std::string m_fieldDataValue;
    std::string m_separator;
    std::string m_fieldIn;
};

class KafkaCluster
{
public:
    KafkaCluster();
    ~KafkaCluster();

    bool Parse(const Json::Value &jsonValue, ApiError &error);
    std::string ToString();
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Valid(ApiError &error);

public:
    std::vector<Address> m_addresses;
    KafkaConfig m_kafkaConfig;
};

class PulsarCluster
{
public:
    PulsarCluster();
    ~PulsarCluster();
    bool Parse(const Json::Value &jsonValue, ApiError &api);
    std::string ToString();
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Valid(ApiError &error);

public:
    std::string m_token;
    std::vector<Address> m_addresses;
};

class DSProxyExporter
{
public:
    std::vector<Address> m_addresses;
    //支持SSL
    std::string m_certPath;
    std::string m_passwdFilePath;
    std::string m_proxyVersion;
    std::string m_proxyProtocol;

    std::string m_httpURI;
    std::string m_thirdPartyCertPasswd;
    std::string m_thirdPartyCertFile;
    std::string m_thirdPartyKeyFile;
    bool m_isThirdPartyCert;
    int m_connectionNumEachAddress;
    bool m_heartbeat;
    bool m_noblock;
    bool m_fillChannelid;

    DSProxyExporter &operator=(const DSProxyExporter &srcConf)
    {
        std::size_t max_count = srcConf.m_addresses.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            m_addresses.push_back(srcConf.m_addresses.at(idx));
        }
        this->m_certPath = srcConf.m_certPath;
        this->m_passwdFilePath = srcConf.m_passwdFilePath;
        this->m_connectionNumEachAddress = srcConf.m_connectionNumEachAddress;
        this->m_proxyVersion = srcConf.m_proxyVersion;
        this->m_proxyProtocol = srcConf.m_proxyProtocol;
        this->m_heartbeat = srcConf.m_heartbeat;
        this->m_noblock = srcConf.m_noblock;
        this->m_fillChannelid = srcConf.m_fillChannelid;
        this->m_thirdPartyCertFile = srcConf.m_thirdPartyCertFile;
        this->m_thirdPartyCertPasswd = srcConf.m_thirdPartyCertPasswd;
        this->m_thirdPartyKeyFile = srcConf.m_thirdPartyKeyFile;
        this->m_httpURI = srcConf.m_httpURI;
        this->m_isThirdPartyCert = srcConf.m_isThirdPartyCert;
        return *this;
    }

    DSProxyExporter()
    {
        m_connectionNumEachAddress = 5;
        m_heartbeat = false;
        m_noblock = true;
        m_fillChannelid = false;
        m_isThirdPartyCert = false;
    }
    DSProxyExporter(const DSProxyExporter &srcConf)
    {
        m_addresses.clear();
        *this = srcConf;
    }

    std::string ToString()
    {
        std::string strDsProxyAddr;
        std::vector<Address>::iterator it;
        for (it = m_addresses.begin(); it != m_addresses.end(); it++)
        {
            strDsProxyAddr.append((*it).m_ip);
            strDsProxyAddr.append(":");
            strDsProxyAddr.append(gse::tools::strings::ToString((*it).m_port));
            if (it != m_addresses.end())
            {
                strDsProxyAddr.append(";");
                break;
            }
        }
        return strDsProxyAddr;
    }

    bool Parse(const Json::Value &jsValue, ApiError &error);
    bool Valid(ApiError &error);
    std::string ToJSON();
    Json::Value ToJSONValue();
};

class RedisCluster
{
public:
    RedisCluster() {}
    ~RedisCluster() {}
    bool Parse(const Json::Value &jsonValue, ApiError &error);
    std::string ToString();
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Valid(ApiError &error);

public:
    std::vector<Address> m_addresses;
    std::string m_passwd;
    std::string m_masterName;
    std::string m_sentinelPasswd;
    std::string m_mode;
};

class DsProxy
{
public:
    DsProxy()
        : m_connectionNumEachAddress(2), m_heartbeat(false), m_fillChannelid(true), m_isThirdPartyCert(false),
          m_noblock(true)
    {
    }

    std::string ToJSON();
    Json::Value ToJSONValue();
    std::string ToString() {}

public:
    std::vector<std::string> m_extentions;
    std::string m_certPath;
    std::string m_proxyVersion;
    std::string m_proxyProtocol;
    int m_connectionNumEachAddress;
    bool m_heartbeat;
    bool m_noblock;
    bool m_fillChannelid;
    bool m_isThirdPartyCert;

    std::vector<Address> m_storagAddress;
    std::string m_httpURI;
    std::string m_thirdPartyKeyFile;
    std::string m_thirdPartyCertPasswd;
    std::string m_thirdPartyCertFile;
};

class StreamToCluster
{
public:
    StreamToCluster();
    ~StreamToCluster();

public:
    StreamToCluster *Clone();
    void CopyFrom(StreamToCluster *src);
    Json::Value ToJSONValue();
    std::string ToJSON();
    bool Parse(const Json::Value &jsonValue, ApiError &error);
    bool Valid(ApiError &error);

public:
    std::string m_name;
    std::string m_reportMode;
    int m_clusterType; //将reportMode转为集群类型

    std::shared_ptr<KafkaCluster> m_kafkaCluster;
    std::shared_ptr<PulsarCluster> m_pulsarCluster;
    std::shared_ptr<RedisCluster> m_redisCluster;
    std::shared_ptr<DSProxyExporter> m_dsproxy;
    std::string m_logfilePath;
};

} // namespace data
} // namespace gse

#endif // CHANNEL_ID_STRUCT_H
