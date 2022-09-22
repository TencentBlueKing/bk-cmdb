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

#include "channelid_api.h"
#include <json/json.h>
#include <json/reader.h>

#include "channelid_def.h"
#include "channelid_zk_api.h"
#include "dataserver.h"
#include "error_code.h"
#include "log/log.h"
#include "safe_json.h"
#include "tools/json_helper.hpp"
#include "tools/strings.h"

namespace gse {
namespace data {

bool Operation::ParseOperation(const Json::Value &operation, ApiError &error)
{
    try
    {
        m_operatorName = SafeJson::GetString(operation, "operator_name", "");
        m_method = SafeJson::GetString(operation, "method", "");
    }
    catch (exception &e)
    {
        LOG_ERROR("failed to parse operation, exception:%s", e.what());
        error = OperationParamInvalid;
        return false;
    }

    return true;
}

bool Operation::Valid(ApiError &error)
{
    if (m_operatorName.empty())
    {
        error = OperationNotSet;
        return false;
    }

    return true;
}

bool Condition::Parse(const Json::Value &requestJson, ApiError &error)
{
    m_requestPlatName = SafeJson::GetString(requestJson, "plat_name", "");
    m_channelID = SafeJson::GetInt(requestJson, "channel_id", 0);
    if (m_channelID == 0)
    {
        error = ChannelIdNotSet;
        return false;
    }

    if (requestJson.isMember("label"))
    {
        m_label.Parse(requestJson["label"]);
    }

    return true;
}

bool Condition::Valid()
{
    if (m_channelID == 0)
    {
        return false;
    }
    return true;
}

bool StreamToCondition::Parse(const Json::Value &requestJson, ApiError &error)
{
    m_requestPlatName = SafeJson::GetString(requestJson, "plat_name", "");
    m_streamToId = SafeJson::GetInt(requestJson, "stream_to_id", 0);
    if (m_streamToId == 0)
    {
        error = StreamToIdNotSet;
        return false;
    }

    if (requestJson.isMember("label"))
    {
        m_label.Parse(requestJson["label"]);
    }

    return true;
}

bool StreamToCondition::Valid()
{
    if (m_streamToId == 0)
    {
        return false;
    }
    return true;
}

bool UpdateSpecification::Parse(const Json::Value &jsonValue, ApiError &error)
{
    if (jsonValue.isMember("metadata"))
    {
        if (m_ptrMetadata->ParseMetadata(jsonValue["metadata"], error))
        {
            LOG_ERROR("failed to parse metadata, error:%s", GET_ERROR_MESSAGE(error).c_str());
            error = MetadataJsonInvalid;
            return false;
        }
    }

    if (!jsonValue.isMember("route"))
    {
        error = ChannelRouteNotSet;
        return false;
    }

    // parse route array

    if (!jsonValue["route"].isArray())
    {
        error = ChannelRouteNotARRAY;
        LOG_ERROR("failed to parse channel route, the route is not an array, input json:%s", jsonValue.toStyledString().c_str());
        return false;
    }

    for (int i = 0; i < jsonValue["route"].size(); i++)
    {
        auto channel = std::make_shared<Channel>();
        if (!channel->Parse(jsonValue["route"][i], error))
        {
            LOG_ERROR("failed to parse channel route, error:%s", GET_ERROR_MESSAGE(error).c_str());
            return false;
        }

        m_channels.push_back(channel);
    }

    if (jsonValue.isMember("stream_filters"))
    {
        if (!jsonValue["stream_filters"].isArray())
        {
            error = FilterNotARRAY;
            LOG_ERROR("failed to parse stream_filters, stream_filters is not an array, json:%s", jsonValue.toStyledString().c_str());
            return false;
        }

        for (int i = 0; i < jsonValue["stream_filters"].size(); i++)
        {
            auto streamFilter = std::make_shared<StreamFilter>();
            if (!streamFilter->Parse(jsonValue["stream_filters"][i], error))
            {
                LOG_ERROR("failed to parse stream_filter, json:(%s), error:%s", jsonValue["stream_filters"][i].toStyledString().c_str(), GET_ERROR_MESSAGE(error).c_str());
                return false;
            }

            m_streamFilters.push_back(streamFilter);
        }
    }

    return true;
}

bool DeleteSpecification::Parse(const Json::Value &jsonValue, ApiError &error)
{
    if (jsonValue.isMember("route"))
    {
        Json::Value channelValues = jsonValue["route"];
        Json::ArrayIndex channelSize = channelValues.size();
        for (Json::ArrayIndex idx = 0; idx < channelSize; ++idx)
        {
            if (!channelValues[idx].isString())
            {
                error = ChannelNameNotString;
                return false;
            }
            m_channelNames.push_back(channelValues[idx].asString());
        }
    }

    if (jsonValue.isMember("stream_filters"))
    {
        Json::Value streamFilterValue = jsonValue["stream_filters"];
        Json::ArrayIndex streamFilterValueSize = streamFilterValue.size();
        for (Json::ArrayIndex idx = 0; idx < streamFilterValueSize; ++idx)
        {
            if (!streamFilterValue[idx].isString())
            {
                error = FilterNameNotString;
                return false;
            }
            m_streamFilterNames.push_back(streamFilterValue[idx].asString());
        }
    }

    return true;
}

StreamIdCondition::StreamIdCondition()
    : m_streamToId(0) {}

bool StreamIdCondition::Parse(const Json::Value &requestJson, ApiError &error)
{
    m_requestPlatName = SafeJson::GetString(requestJson, "plat_name", "");
    m_streamToId = SafeJson::GetInt(requestJson, "stream_to_id", 0);
    if (m_streamToId == 0)
    {
        error = StreamToIdNotSet;
        return false;
    }

    if (requestJson.isMember("label"))
    {
        m_label.Parse(requestJson["label"]);
    }

    return true;
}

//-----------------------------------------

ChannelIdCommonAPI::ChannelIdCommonAPI()
{
    m_buildInChannelIdRange.push_back(std::make_tuple(MONITOR_RESERVER_CHANNEL_ID_START_1, MONITOR_RESERVER_CHANNEL_ID_END_1));
    m_buildInChannelIdRange.push_back(std::make_tuple(MONITOR_RESERVER_CHANNEL_ID_START_2, MONITOR_RESERVER_CHANNEL_ID_END_2));
    m_buildInChannelIdRange.push_back(std::make_tuple(GSE_RESERVER_CHANNEL_ID_START_1, GSE_RESERVER_CHANNEL_ID_END_1));

    m_reportModeTypes = {{ChannelIDOperationRequestMethodFile},
                         {ChannelIDOperationRequestMethodKafka},
                         {ChannelIDOperationRequestMethodPulsar},
                         {ChannelIDOperationRequestMethodRedis},
                         {ChannelIDOperationRequestMethodDsProxy}};
}

bool ChannelIdCommonAPI::IsBuiltinChannelId(int channelid)
{
    for (auto range : m_buildInChannelIdRange)
    {
        if (channelid >= std::get<0>(range) && channelid <= std::get<1>(range))
        {
            return true;
        }
    }

    return false;
}

bool ChannelIdCommonAPI::IsTglogPlat(const std::string &platName)
{
    return (platName.compare(SUPPORT_PLAT_NAME_TGLOG) == 0);
}

bool ChannelIdCommonAPI::IsSupportClusterType(const std::string &clusterType)
{
    auto it = m_reportModeTypes.find(clusterType);
    if (it != m_reportModeTypes.end())
    {
        return true;
    }
    return false;
}

void ChannelIdCommonAPI::MakeResponse(Json::Value &response, ApiError &error, const Json::Value &data)
{
    int errorCode = GET_ERROR_CODE(error);
    std::string errorMsg = GET_ERROR_MESSAGE(error);
    // int bkErrorCode = BK_ERROR_CODE(errorCode);

    bool result = true;
    if (errorCode != 0)
    {
        response["message"] = Json::Value(errorMsg);
        result = false;
    }
    else
    {
        response["message"] = Json::Value("success");
        result = true;
    }

    response["result"] = result;
    response["code"] = errorCode;
    response["data"] = data;
}

void ChannelIdCommonAPI::SplitZkPath(const std::string &nodePath, std::vector<std::string> &nodes)
{
    std::string::size_type pos2 = nodePath.find('/');
    std::string::size_type pos1 = 0;
    while (std::string::npos != pos2)
    {
        nodes.push_back(nodePath.substr(pos1, pos2 - pos1));
        pos1 = pos2 + 1;
        pos2 = nodePath.find('/', pos1);
    }
    nodes.push_back(nodePath.substr(pos1));
    return;
}

// bool ChannelIdCommonAPI::IsInLimit(int channelid, const std::string &platName)
//{
//     std::map<std::string, std::list<std::pair<int, int>>> buildid_channelid_limit;
//     std::list<std::pair<int, int>> gse_range_list;
//     std::string plat_name = std::string(SUPPORT_PLAT_NAME_GSE);
//     std::pair<int, int> range;
//     range = std::make_pair(GSE_RESERVER_CHANNEL_ID_START_1, GSE_RESERVER_CHANNEL_ID_END_1);
//     gse_range_list.push_back(range);
//     buildid_channelid_limit[plat_name] = gse_range_list;

//    platname = std::string(SUPPORT_PLAT_NAME_BKMONITOR);

//    std::list<std::pair<int, int>> monitor_range_list;
//    monitor_range_list.push_back(std::make_pair(MONITOR_RESERVER_CHANNEL_ID_START_1, MONITOR_RESERVER_CHANNEL_ID_END_1));
//    monitor_range_list.push_back(std::make_pair(MONITOR_RESERVER_CHANNEL_ID_START_2, MONITOR_RESERVER_CHANNEL_ID_END_2));
//    buildid_channelid_limit[plat_name] = monitor_range_list;

//    std::map<std::string, std::list<std::pair<int, int>>>::iterator it;
//    it = buildid_channelid_limit.find(platname);
//    if (it != buildid_channelid_limit.end())
//    {
//        std::list<std::pair<int, int>> &range = it->second;
//        for (std::list<std::pair<int, int>>::iterator it_range = range.begin(); it_range != range.end(); it_range++)
//        {
//            if (channel_id >= GSE_RESERVER_CHANNEL_ID_START_1 && channel_id <= GSE_RESERVER_CHANNEL_ID_END_1)
//            {
//                return false;
//            }
//            else
//            {
//                return true;
//            }
//        }
//    }
//}

std::shared_ptr<ChannelIdCommonAPI> ChannelIdCommApiInst()
{
    return DataServer::GetChannelIdCommonApiInst();
}

} // namespace data
} // namespace gse
