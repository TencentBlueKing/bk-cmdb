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

#ifndef _GSE_DATA_CHANNELID_API_
#define _GSE_DATA_CHANNELID_API_

#include <memory>
#include <set>

#include <json/json.h>

#include "channelid_struct.h"
#include "channelid_def.h"
#include "error_code.h"
#include "log/log.h"
#include "safe_json.h"
#include "tools/json_helper.hpp"

namespace gse {
namespace data {

using RangeStart = int;
using RangeEnd = int;

class StreamToMetadata
{
public:
    std::string m_name;
    std::string m_version;
    std::string m_requestPlatName;
    Label m_label;
    uint32_t m_streamToId;
    StreamToMetadata()
    {
        m_streamToId = 0;
    }

public:
    inline Json::Value ToJSONValue()
    {
        Json::Value meta;

        meta["name"] = Json::Value(m_name);
        meta["version"] = Json::Value(m_version);
        meta["stream_to_id"] = Json::Value(m_streamToId);
        meta["plat_name"] = Json::Value(m_requestPlatName);
        meta["label"] = m_label.ToJSONValue();
        return meta;
    }

    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    bool ParseMetadata(const Json::Value& metaJson, ApiError& error)
    {
        try
        {
            m_name = SafeJson::GetString(metaJson, "name", "");
            m_version = SafeJson::GetString(metaJson, "version", "");
            m_requestPlatName = SafeJson::GetString(metaJson, "plat_name", "");
            m_streamToId = SafeJson::GetInt(metaJson, "stream_to_id", 0);

            if (metaJson.isMember("label"))
            {
                m_label.Parse(metaJson["label"]);
            }
        }
        catch (exception& e)
        {
            error = MetadataJsonInvalid;
            LOG_ERROR("failed to parse streamto metadata, error:%s", e.what());
            return false;
        }

        return true;
    }
};

class Operation
{
public:
    std::string m_operatorName;
    std::string m_method;

public:
    inline Json::Value ToJSONValue()
    {
        Json::Value data;

        data["operator_name"] = m_operatorName;
        return data;
    }
    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    bool ParseOperation(const Json::Value& operation, ApiError& error);
    bool Valid(ApiError& error);
};

class StreamIdCondition
{
public:
    StreamIdCondition();
    bool Parse(const Json::Value& routeJson, ApiError& error);

public:
    uint32_t m_streamToId;
    std::string m_requestPlatName;
    Label m_label;
};

class UpdateStreamToIdSpecification
{
public:
    StreamToMetadata* m_ptrMetadata;
    // StreamToConfig m_streamToConfig;

public:
    UpdateStreamToIdSpecification()
    {
        m_ptrMetadata = new StreamToMetadata();
    }

    ~UpdateStreamToIdSpecification()
    {
        if (NULL != m_ptrMetadata)
        {
            delete m_ptrMetadata;
        }
    }
};

class ChannelIdConfig
{
public:
    ChannelIdConfig()
    {
        m_ptrMetadata = std::make_shared<Metadata>();
    }
    inline Json::Value ToJSONValue()
    {
        Json::Value data;

        data["metadata"] = m_ptrMetadata->ToJSONValue();
        data["route"].resize(0);
        for (auto channel : m_channels)
        {
            data["route"].append(channel->ToJSONValue());
        }

        data["stream_filters"].resize(0);
        for (auto filter : m_streamFilters)
        {
            data["stream_filters"].append(filter->ToJSONValue());
        }
        return data;
    }
    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

public:
    std::shared_ptr<Metadata> m_ptrMetadata;
    std::vector<std::shared_ptr<Channel>> m_channels;
    std::vector<std::shared_ptr<StreamFilter>> m_streamFilters;
};

class StreamToIdConfig
{
public:
    StreamToIdConfig()
    {
        m_ptrMetadata = std::make_shared<StreamToMetadata>();
        m_ptrStreamToCluster = std::make_shared<StreamToCluster>();
    }

    std::shared_ptr<StreamToMetadata> m_ptrMetadata;
    std::shared_ptr<StreamToCluster> m_ptrStreamToCluster;

    inline Json::Value ToJSONValue()
    {
        Json::Value data;

        data["metadata"] = m_ptrMetadata->ToJSONValue();
        data["stream_to"] = m_ptrStreamToCluster->ToJSONValue();
        return data;
    }
    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

private:
};

class StreamToIdUpdate
{
public:
    StreamIdCondition* m_ptrCondition;
    Operation* m_ptrOperation;
    UpdateStreamToIdSpecification* m_ptrSpecification;
    std::string m_originData;

public:
    StreamToIdUpdate()
    {
        m_ptrCondition = new StreamIdCondition();
        m_ptrOperation = new Operation();
        m_ptrSpecification = new UpdateStreamToIdSpecification();
    }

    ~StreamToIdUpdate()
    {
        if (NULL != m_ptrCondition)
        {
            delete m_ptrCondition;
        }

        if (NULL != m_ptrOperation)
        {
            delete m_ptrOperation;
        }

        if (NULL != m_ptrSpecification)
        {
            delete m_ptrSpecification;
        }
    }
};

class StreamToIdQuery
{
public:
    StreamIdCondition* m_ptrCondition;
    Operation* m_ptrOperation;
    std::string m_originData;

public:
    StreamToIdQuery()
    {
        m_ptrCondition = new StreamIdCondition();
        m_ptrOperation = new Operation();
    }

    ~StreamToIdQuery()
    {
        if (NULL != m_ptrCondition)
        {
            delete m_ptrCondition;
        }

        if (NULL != m_ptrOperation)
        {
            delete m_ptrOperation;
        }
    }
};

// class StreamToIdDelete
//{
// public:
//     StreamIdCondition* m_ptrCondition;
//     Operation* m_ptrOperation;
//     std::string m_originData;

// public:
//     StreamToIdDelete()
//     {
//         m_ptrCondition = new StreamIdCondition();
//         m_ptrOperation = new Operation();
//     }
//     ~StreamToIdDelete()
//     {
//         if (NULL != m_ptrCondition)
//         {
//             delete m_ptrCondition;
//         }

//        if (NULL != m_ptrOperation)
//        {
//            delete m_ptrOperation;
//        }
//    }
//};

// typedef StreamToIdRequest StreamToIdConfig;

class Condition
{
public:
    Condition()
    {
        m_channelID = 0;
    }

    bool Parse(const Json::Value& requestJson, ApiError& error);
    bool Valid();

public:
    uint32_t m_channelID;
    std::string m_name;
    std::string m_version;
    std::string m_requestPlatName;
    Label m_label;
};

class StreamToCondition
{
public:
    StreamToCondition()
    {
        m_streamToId = 0;
    }

    bool Parse(const Json::Value& requestJson, ApiError& error);
    bool Valid();

public:
    uint32_t m_streamToId;
    std::string m_name;
    std::string m_version;
    std::string m_requestPlatName;
    Label m_label;
};

class Query
{
public:
    Condition* m_ptrCondition;
    Operation* m_ptrOperation;
    std::string m_originData;

public:
    Query()
    {
        m_ptrCondition = new Condition();
        m_ptrOperation = new Operation();
    }

    ~Query()
    {
        if (NULL != m_ptrCondition)
        {
            delete m_ptrCondition;
        }

        if (NULL != m_ptrOperation)
        {
            delete m_ptrOperation;
        }
    }
};

class UpdateSpecification
{
public:
    std::shared_ptr<Metadata> m_ptrMetadata;
    std::vector<std::shared_ptr<Channel>> m_channels;
    std::vector<std::shared_ptr<StreamFilter>> m_streamFilters;

public:
    UpdateSpecification()
    {
    }

    ~UpdateSpecification()
    {
    }

    bool Parse(const Json::Value& jsonValue, ApiError& error);
};

class Update
{
public:
    Condition* m_ptrCondition;
    Operation* m_ptrOperation;
    UpdateSpecification* m_ptrSpecification;
    std::string m_originData;

public:
    Update()
    {
        m_ptrCondition = new Condition();
        m_ptrOperation = new Operation();
        m_ptrSpecification = new UpdateSpecification();
    }

    ~Update()
    {
        if (NULL != m_ptrCondition)
        {
            delete m_ptrCondition;
        }

        if (NULL != m_ptrOperation)
        {
            delete m_ptrOperation;
        }

        if (NULL != m_ptrSpecification)
        {
            delete m_ptrSpecification;
        }
    }
};

class DeleteSpecification
{
public:
    bool Parse(const Json::Value& jsonValue, ApiError& error);

public:
    std::vector<std::string> m_channelNames;
    std::vector<std::string> m_streamFilterNames;
};

class Delete
{
public:
    Condition* m_ptrCondition;
    Operation* m_ptrOperation;
    DeleteSpecification* m_ptrDeleteSpecification;
    std::string m_originData;

public:
    Delete()
    {
        m_ptrCondition = new Condition();
        m_ptrOperation = new Operation();
        m_ptrDeleteSpecification = new DeleteSpecification();
    }
    ~Delete()
    {
        if (NULL != m_ptrCondition)
        {
            delete m_ptrCondition;
        }

        if (NULL != m_ptrOperation)
        {
            delete m_ptrOperation;
        }

        if (NULL != m_ptrDeleteSpecification)
        {
            delete m_ptrDeleteSpecification;
        }
    }
};

class ChannelIdCommonAPI
{
public:
    ChannelIdCommonAPI();
    bool IsBuiltinChannelId(int channelid);
    static bool IsTglogPlat(const std::string& platName);
    bool IsSupportClusterType(const std::string& clusterType);
    static void MakeResponse(Json::Value& response, ApiError& error, const Json::Value& data);
    static void SplitZkPath(const std::string& nodePath, std::vector<std::string>& nodes);

private:
    std::vector<std::tuple<int, int>> m_buildInChannelIdRange;
    std::set<std::string> m_reportModeTypes;
};

std::shared_ptr<ChannelIdCommonAPI> ChannelIdCommApiInst();

} // namespace data
} // namespace gse
#endif
