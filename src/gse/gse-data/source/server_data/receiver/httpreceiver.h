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

#ifndef _GSE_DATA_HTTP_RECEIVER_H_
#define _GSE_DATA_HTTP_RECEIVER_H_

#include "receiver.h"
#include "net/http/http_server.h"
namespace gse { 
namespace dataserver {

#ifndef SEND_DATA_BY_DATAID
#define SEND_DATA_BY_DATAID "/gse/v1/senddata/"
#endif

class HttpReceiver: public Receiver
{
public:
    HttpReceiver();
    virtual ~HttpReceiver();

public:
    int Start();
    int Stop();
    void Join();

private:

private:
    uint32_t getChannelIDFromURI(const char* uri);
    void OnHttpMessageHandler(gse::net::http::HttpMessagePtr message, std::string &response);
    void makeResponse(int errorCode, const std::string &message);

private:
    gse::net::http::HttpServer *m_httpserver;
    std::thread m_listenThread;
};

}
}
#endif
