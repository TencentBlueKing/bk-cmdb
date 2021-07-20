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

#ifndef _GSE_DATA_RECEIVER_H_
#define _GSE_DATA_RECEIVER_H_

#include <string>
#include "conf/confItem.h"
#include "safe/ssl.h"
#include "safe/lock.h"
#include "datacell.h"
#include "opscollection/ops.h"

namespace gse { 
namespace dataserver {

typedef void (*RecvDataCallBack)(DataCell *pDataCell, void* pCaller);

class Receiver
{
public:
    Receiver();
    virtual ~Receiver();

public:
    virtual int Start() = 0;
    virtual int Stop() = 0;
    virtual void Join() = 0;

public:
    inline void SetGseConf(DataProcessConfig& cfg)
    {
        m_gseCfg = cfg;
    }

    inline void SetOPS(OpsCollection*  ptrOPSReport)
    {
        m_ptrOPSReport = ptrOPSReport;
    }

public:
    void UpdateConf(const ReceiverConf& receiverConf);
    void SetRecvCallBack(RecvDataCallBack fnRecvData, void* pCaller);

protected:
    ReceiverConf* m_recevierConf;
    RecvDataCallBack m_fnRecvData;
    void* m_pCaller;
    std::string m_servIp;
    uint16_t m_servPort;
    SSL*  m_ptrSSL;

protected:
    DataProcessConfig  m_gseCfg;

protected:
    OpsCollection*  m_ptrOPSReport;

private:
   gse::safe::RWLock m_upConfLock;
};

}
}
#endif
