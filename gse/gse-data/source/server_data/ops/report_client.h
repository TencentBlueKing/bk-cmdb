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

#ifndef SOURCE_GSE_DATA_SERVER_DATA_DSCLIENT_OPSCLIENT_H_
#define SOURCE_GSE_DATA_SERVER_DATA_DSCLIENT_OPSCLIENT_H_

#include <sys/socket.h>
#include <arpa/inet.h>
#include <string>
#include <vector>
#include "net/udp/udp_client.h"

//#include "codec/gsedata_pkg/gsedata_package.h"
//#include "codec/gsedata_pkg/gsedata_extension.h"
#include "protocol.h"
namespace gse { 
namespace data {

class ReportClient
{

public:
    ReportClient();
    ~ReportClient();

public:
    int Open( std::string& ip, uint16_t port);
    int SendTo(int dataid, const char* ptrData, int dataLen);
    int SendToByV2(int dataID, const char* ptrData, int dataLen);
    int Close();

private:
    gse::net::udp::UDPClient*          m_ptrUDPClient;

};

}
}
#endif /* SOURCE_GSE_DATA_SERVER_DATA_DSCLIENT_OPSCLIENT_H_ */

