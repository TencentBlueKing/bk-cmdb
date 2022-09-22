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

#ifndef _GSE_DATA_OP_STACK_H_
#define _GSE_DATA_OP_STACK_H_

#include <map>

#include "net/http/http_handler.hpp"

namespace gse {
namespace data {

class Stack : public net::http::HTTPHandler
{
    typedef int (Stack::*SubCmdFuncPtr)(gse::net::http::HTTPMessagePtr message, const std::string& param, std::string& result);
    typedef int (Stack::*CmdFuncPtr)(gse::net::http::HTTPMessagePtr message, const std::string& param, std::string& result);

public:
    Stack();
    virtual ~Stack();

    std::string GetURI();
    int OnGet(gse::net::http::HTTPMessagePtr message, std::string& response) noexcept;
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string& response) noexcept;

private:
    int HandleListCmd(net::http::HTTPMessagePtr message, std::string& queryResult);
    int QueryChannelID(gse::net::http::HTTPMessagePtr message, const std::string& channelID, std::string& queryResult);
    int QueryStreamExporterConfig(gse::net::http::HTTPMessagePtr message, const std::string& streamID, std::string& queryResult);
    int Sync(gse::net::http::HTTPMessagePtr message, const std::string& channelID, std::string& queryResult);

    int SyncChannelId(gse::net::http::HTTPMessagePtr message, const std::string& channelID, std::string& queryResult);

    void PackErrorResult(const std::string& errMsg, std::string& queryResult);

private:
    std::map<std::string, CmdFuncPtr> m_cmdList;
    std::map<std::string, SubCmdFuncPtr> m_syncCmdList;
};

} // namespace data
} // namespace gse
#endif // STACK_H
