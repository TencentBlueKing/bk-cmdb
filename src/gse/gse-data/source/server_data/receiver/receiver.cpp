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

#include "receiver.h"
namespace gse {
namespace data {

Receiver::Receiver()
{
    m_ptrOPSReport = NULL;
    m_recevierConf = NULL;
    m_fnRecvData = NULL;
    m_ptrSSL = NULL;
    m_pCaller = NULL;
    m_servPort = 0;
}

Receiver::~Receiver()
{
    if (m_recevierConf != NULL)
    {
        delete m_recevierConf;
        m_recevierConf = NULL;
    }

    if (NULL != m_ptrSSL)
    {
        delete m_ptrSSL;
    }
}

void Receiver::UpdateConf(const ReceiverConf& receiverConf)
{
    if (m_recevierConf != NULL)
    {
        delete m_recevierConf;
        m_recevierConf = NULL;
    }

    m_recevierConf = new ReceiverConf(receiverConf);
}

void Receiver::SetRecvCallBack(RecvDataCallBack fnRecvData, void* pCaller)
{
    m_fnRecvData = fnRecvData;
    m_pCaller = pCaller;
    LOG_DEBUG("Set receiver(%s) callback:%p", m_recevierConf->m_name.c_str(), fnRecvData);
}
} // namespace data
} // namespace gse
