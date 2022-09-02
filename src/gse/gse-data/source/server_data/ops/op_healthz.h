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
#ifndef _GSE_OP_HEALTHZ_H_
#define _GSE_OP_HEALTHZ_H_

#include <tuple>

#include "net/http/http_handler.hpp"
#include "safe/lock.h"

namespace gse {
namespace data {

using namespace gse::net::http;

using HealthzName = std::string;
using HealthMesage = std::tuple<std::string, int>;

class OpHealthZ : public net::http::HTTPHandler
{
public:
    OpHealthZ();
    virtual ~OpHealthZ();
    int OnGet(HTTPMessagePtr message, std::string& response) noexcept;
    int OnPost(HTTPMessagePtr message, std::string& response) noexcept;
    std::string GetURI();

    static void AddInitHealthInfo(const std::string& healthName, const std::string& healthMessae, int status);

private:
    void DiscoverZkCheck();
    void ChannelIdZkCheck();

private:
    static gse::safe::RWLock m_lock;
};

} // namespace data
} // namespace gse
#endif
