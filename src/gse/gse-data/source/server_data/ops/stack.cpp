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

#include "stack.h"

#include <algorithm>
#include <cctype>
#include <string>

#include "conf/channel_id_config.h"
#include "dataserver.h"

namespace gse {
namespace data {
// pararam
// channel-id={channel-id}
// streamto-id={streamto-id}
// query channel-id's list by streamto-id streamto-id={streamto-id}&list-channelid=true

// dynamic config
// query streamto status, connection, send,recv, base config, error message....

// sync=channelid&channel-id={}
// sync=streamid&channel-id={}
//  force load  channelid's config
//  force load streamto's config

static const std::string kErrorMsgChannelIdNotExsit = "channel-id not exist";
static const std::string kErrorMsgStreamIdNotExsit = "stream-id not exist";

Stack::Stack()
{
    // m_cmdList.insert(std::make_pair("channel-id", &Stack::QueryChannelID));
    m_cmdList = {{"channel-id", &Stack::QueryChannelID},
                 {"stream-id", &Stack::QueryStreamExporterConfig},
                 {"sync", &Stack::Sync}};
}

Stack::~Stack()
{
}

std::string Stack::GetURI()
{
    return std::string("/stack");
}
int Stack::OnGet(net::http::HTTPMessagePtr message, string& response) noexcept
{
    return OnPost(message, response);
}

int Stack::OnPost(gse::net::http::HTTPMessagePtr message, std::string& response) noexcept
{
    LOG_DEBUG("recv http request, URI:%s", message->GetURI().c_str());
    auto cmdIt = m_cmdList.begin();
    for (; cmdIt != m_cmdList.end(); cmdIt++)
    {
        auto httpParam = message->GetURLParameter(cmdIt->first);
        if (httpParam != nullptr)
        {
            std::string param = httpParam->GetValue();
            CmdFuncPtr cmdPtr = cmdIt->second;
            return (this->*cmdPtr)(message, param, response);
        }
    }

    return GSE_ERROR;
}

int Stack::Sync(gse::net::http::HTTPMessagePtr message, const std::string& channelID, std::string& queryResult)
{
    return HTTP_STATUS_CODE_200;
}

int Stack::SyncChannelId(gse::net::http::HTTPMessagePtr message, const std::string& channelID, std::string& queryResult)
{
}

int Stack::QueryChannelID(gse::net::http::HTTPMessagePtr message, const std::string& channelID, std::string& queryResult)
{
    LOG_DEBUG("query channelid(%s)'s config", channelID.c_str());
    uint32_t uChannelId = gse::tools::strings::StringToUint32(channelID);

    if (DataServer::GetChannelIdManagerInst()->Find(uChannelId))
    {
        auto channelIdConfig = DataServer::GetChannelIdManagerInst()->GetChannelStreamConfig(uChannelId);
        // static config
        if (channelIdConfig)
        {
            rapidjson::Document doc;
            std::string strConfig;
            channelIdConfig->Dump(strConfig);
            if (doc.Parse(strConfig.data()).HasParseError())
            {
                PackErrorResult(kErrorMsgStreamIdNotExsit, queryResult);
                return HTTP_STATUS_CODE_200;
            }

            rapidjson::StringBuffer strBuff;
            rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strBuff);
            doc.Accept(writer);

            queryResult.assign(strBuff.GetString());
        }
        else
        {
            PackErrorResult(kErrorMsgStreamIdNotExsit, queryResult);
        }

        return HTTP_STATUS_CODE_200;
    }

    PackErrorResult(kErrorMsgStreamIdNotExsit, queryResult);
    return HTTP_STATUS_CODE_200;
}

void Stack::PackErrorResult(const std::string& errMsg, std::string& queryResult)
{
    rapidjson::StringBuffer strBuf;
    rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strBuf);
    writer.StartObject();
    RAPIDJSON_SET_INT32(writer, "error_code", -1);
    RAPIDJSON_SET_STRING(writer, "error_message", errMsg.c_str());
    writer.EndObject();

    queryResult.assign(strBuf.GetString());
}

// ?streamto-id={streamto-id}&list-channelid=true
int Stack::QueryStreamExporterConfig(gse::net::http::HTTPMessagePtr message, const std::string& streamID, std::string& queryResult)
{
    LOG_DEBUG("query streamID(%s)'s config", streamID.c_str());
    auto listChannelidParam = message->GetURLParameter("list-channelid");
    bool bListChannelId = false;
    if (listChannelidParam != nullptr)
    {
        auto value = listChannelidParam->GetValue();
        if (strncasecmp(value.c_str(), "true", sizeof("true")) == 0)
        {
            bListChannelId = true;
        }

        LOG_DEBUG("input param:%s, isbool:%d", value.c_str(), bListChannelId);
    }

    uint32_t uStreamId = gse::tools::strings::StringToUint32(streamID);

    bool isExist = DataServer::GetStreamExporterManagerInst()->Find(uStreamId);
    if (isExist)
    {
        auto exporter = DataServer::GetStreamExporterManagerInst()->GetExport(uStreamId);

        rapidjson::Document rootJson(rapidjson::kObjectType);
        auto& allocator = rootJson.GetAllocator();
        std::string streamConfig;
        exporter->DumpConfig(streamConfig);

        rapidjson::Document docStream;

        if (docStream.Parse(streamConfig.data()).HasParseError())
        {
            PackErrorResult(kErrorMsgStreamIdNotExsit, queryResult);
            return HTTP_STATUS_CODE_200;
        }

        rootJson.AddMember("stream", docStream, allocator);
        if (bListChannelId)
        {
            std::vector<uint32_t> channelIdList;
            DataServer::GetChannelIdManagerInst()->GetChannelIdListByStreamId(uStreamId, channelIdList);

            rj::Value channelidList(rj::kObjectType);
            rapidjson::Value channellist(rapidjson::kArrayType);
            LOG_DEBUG("get streamid:%d channelid list size:%d", uStreamId, channelIdList.size());

            for (auto id : channelIdList)
            {
                channellist.PushBack(id, allocator);
            }
            rootJson.AddMember("channelid-list", channellist, allocator);
        }

        rapidjson::StringBuffer strBuff;
        rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strBuff);
        rootJson.Accept(writer);

        queryResult.assign(strBuff.GetString());
    }
    else
    {
        PackErrorResult(kErrorMsgStreamIdNotExsit, queryResult);
    }
    return HTTP_STATUS_CODE_200;
}

} // namespace data
} // namespace gse
