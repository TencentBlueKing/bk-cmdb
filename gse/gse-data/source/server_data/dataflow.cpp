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

#include "dataflow.h"
#include "dataserver.h"
#include "tools/macros.h"

namespace gse {
namespace data {
using namespace std;

DataFlow::DataFlow()
{
    m_flowConf = NULL;
    m_opsReport = NULL;
    m_ptrIDToStorage = NULL;
}

DataFlow::~DataFlow()
{
    Stop();

    if (m_flowConf != NULL)
    {
        delete m_flowConf;
        m_flowConf = NULL;
    }
}

int DataFlow::Start()
{
    if (NULL == m_flowConf)
    {
        LOG_WARN("the pointer of dataflow configure is NULL, can't start dataflow, please check dataflow configure");
        return GSE_SYSTEMERROR;
    }

    if (m_flowConf->m_channelsConf.size() <= 0)
    {
        LOG_WARN("the number of channel in dataflow configure is zero, please check dataflow configure");
        return GSE_SYSTEMERROR;
    }

    // first stop
    Stop();

    // create channel
    map<string, ChannelConf*>::iterator iter = m_flowConf->m_channelsConf.begin();
    for (; iter != m_flowConf->m_channelsConf.end(); iter++)
    {
        ChannelConf* pChannelConf = iter->second;
        if (NULL == pChannelConf)
        {
            continue;
        }

        DataChannel* pDataChannel = new DataChannel();

        pDataChannel->UpdateConf(*pChannelConf);
        pDataChannel->SetOps(m_opsReport);
        pDataChannel->SetGseConf(DataServer::GetConfigPtr());
        pDataChannel->SetConftor((Configurator*)m_ptrIDToStorage);
        if (GSE_SUCCESS != pDataChannel->Start())
        {
            LOG_WARN("fail to start datachannel[%s]", SAFE_CSTR(pChannelConf->m_name.c_str()));
        }
        else
        {
            LOG_INFO("success to start datachannel[%s]", SAFE_CSTR(pChannelConf->m_name.c_str()));
            m_mapDataChannels[pChannelConf->m_name] = pDataChannel;
        }
    }

    if (m_mapDataChannels.size() == 0)
    {
        LOG_WARN("no datachannel started successfully");
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

int DataFlow::Stop()
{
    map<string, DataChannel*>::iterator iter = m_mapDataChannels.begin();
    for (; iter != m_mapDataChannels.end(); iter++)
    {
        DataChannel* pChannel = iter->second;
        pChannel->Stop();
        delete pChannel;
    }

    m_mapDataChannels.clear();

    return GSE_SUCCESS;
}

void DataFlow::Join()
{
    map<string, DataChannel*>::iterator iter = m_mapDataChannels.begin();
    for (; iter != m_mapDataChannels.end(); iter++)
    {
        DataChannel* pChannel = iter->second;
        pChannel->Join();
    }

    LOG_DEBUG("dataflow service joined");
}

void DataFlow::UpdateConf(DataFlowConf* pFlowConf)
{
    if (NULL == pFlowConf)
    {
        LOG_WARN("the pointer of function parameter(dataflowconf) is NULL ");
        return;
    }

    m_upConfLock.WLock();
    if (m_flowConf != NULL)
    {
        delete m_flowConf;
    }

    m_flowConf = pFlowConf;
    m_upConfLock.UnLock();
}

void DataFlow::UpdateIDToStorage(IDToStorage* ptrIDToStorage)
{
    m_ptrIDToStorage = ptrIDToStorage;
}

} // namespace data
} // namespace gse
