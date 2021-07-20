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

#include "gseEventThreadPool.h"
#include "log/log.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"
#include "tools/hostinfo.h"
namespace gse { 
namespace dataserver {
namespace rgse
{

GseEventThreadPool::GseEventThreadPool(uint32_t threadNum)
{
    m_threadNum = threadNum;

    initThreadPool();
}

GseEventThreadPool::GseEventThreadPool()
{
    m_threadNum = gse::tools::hostinfo::GetCpuNum();
    initThreadPool();
}

GseEventThreadPool::~GseEventThreadPool()
{
    distroyThreadPool();
}

int GseEventThreadPool::regReadEvent(ThreadID id, EVENT_HANDLER eventHandler, evutil_socket_t fd, void * args, bool edgeTriger /*= false*/)
{
    if(m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_INVALIDARGS;
    }

    if(id < 0)
    {
        for(int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->registerReadPersistEvent(eventHandler, fd, args, edgeTriger);
        }
    }
    else
    {
        m_threadPool[id]->registerReadPersistEvent(eventHandler, fd, args, edgeTriger);
    }

    return GSE_SUCCESS;
}

int GseEventThreadPool::regTimer(ThreadID id, EVENT_HANDLER eventHandler, void * args, int second)
{
    if(m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_INVALIDARGS;
    }

    if(id < 0)
    {
        for(int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->registerTimerPersistEvent(eventHandler, args, second);
        }
    }
    else
    {
        m_threadPool[id]->registerTimerPersistEvent(eventHandler, args, second);
    }

    return GSE_SUCCESS;
}

int GseEventThreadPool::resetTimer(ThreadID id, EventIndex eventIndex, int second)
{
    if(m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_INVALIDARGS;
    }

    if(id < 0)
    {
        for(int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->resetEvent(eventIndex, second);
        }
    }
    else
    {
        m_threadPool[id]->resetEvent(eventIndex, second);
    }

    return GSE_SUCCESS;
}

int GseEventThreadPool::unRegEvent(ThreadID id, EventIndex eventIndex)
{
    if(m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_INVALIDARGS;
    }

    if(id < 0)
    {
        for(int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->unRegisterEvent(eventIndex);
        }
    }
    else
    {
        m_threadPool[id]->unRegisterEvent(eventIndex);
    }

    return GSE_SUCCESS;
}

void GseEventThreadPool::clear(ThreadID id)
{
    if(m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return;
    }

    if(id < 0)
    {
        for(int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->clear();
        }
    }
    else
    {
        m_threadPool[id]->clear();
    }
}

int GseEventThreadPool::start()
{
    startThreadPool();

    return GSE_SUCCESS;
}

void GseEventThreadPool::stop()
{
    stopThreadPool();
}

void GseEventThreadPool::initThreadPool()
{
    for(int i = 0; i < m_threadNum; i++)
    {
        GseEventThread * pEventThread = new GseEventThread();
        if(pEventThread == NULL)
        {
            LOG_ERROR("create gse event thread fail");
            ::abort();
        }

        m_threadPool.push_back(pEventThread);
    }
}

void GseEventThreadPool::distroyThreadPool()
{
    for(int i = 0; i < m_threadPool.size(); i++)
    {
        GseEventThread * pEventThread = m_threadPool[i];
        delete pEventThread;
    }
}

void GseEventThreadPool::startThreadPool()
{
    for(int i = 0; i < m_threadNum; i++)
    {
        GseEventThread * pEventThread = m_threadPool[i];
        pEventThread->Start();
    }
}

void GseEventThreadPool::stopThreadPool()
{
    for(int i = 0; i < m_threadNum; i++)
    {
        GseEventThread * pEventThread = m_threadPool[i];
        pEventThread->Stop();
        pEventThread->Join();
    }
}

}



}
}
