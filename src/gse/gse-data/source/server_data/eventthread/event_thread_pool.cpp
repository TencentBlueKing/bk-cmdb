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

#include "event_thread_pool.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/hostinfo.h"
#include "tools/macros.h"
namespace gse {
namespace data {

EventThreadPool::EventThreadPool(uint32_t threadNum)
{
    m_threadNum = threadNum;

    initThreadPool();
}

EventThreadPool::EventThreadPool()
{
    m_threadNum = gse::tools::hostinfo::GetCpuNum();
    initThreadPool();
}

EventThreadPool::~EventThreadPool()
{
    distroyThreadPool();
}

int EventThreadPool::regReadEvent(ThreadID id, EVENT_HANDLER eventHandler, evutil_socket_t fd, void* args, bool edgeTriger /*= false*/)
{
    if (m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_ERROR;
    }

    if (id < 0)
    {
        for (int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->RegisterReadPersistEvent(eventHandler, fd, args, edgeTriger);
        }
    }
    else
    {
        m_threadPool[id]->RegisterReadPersistEvent(eventHandler, fd, args, edgeTriger);
    }

    return GSE_SUCCESS;
}

int EventThreadPool::regTimer(ThreadID id, EVENT_HANDLER eventHandler, void* args, int second)
{
    if (m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_ERROR;
    }

    if (id < 0)
    {
        for (int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->RegisterTimerPersistEvent(eventHandler, args, second);
        }
    }
    else
    {
        m_threadPool[id]->RegisterTimerPersistEvent(eventHandler, args, second);
    }

    return GSE_SUCCESS;
}

int EventThreadPool::resetTimer(ThreadID id, EventIndex eventIndex, int second)
{
    if (m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_ERROR;
    }

    if (id < 0)
    {
        for (int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->ResetEvent(eventIndex, second);
        }
    }
    else
    {
        m_threadPool[id]->ResetEvent(eventIndex, second);
    }

    return GSE_SUCCESS;
}

int EventThreadPool::unRegEvent(ThreadID id, EventIndex eventIndex)
{
    if (m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return GSE_ERROR;
    }

    if (id < 0)
    {
        for (int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->UnRegisterEvent(eventIndex);
        }
    }
    else
    {
        m_threadPool[id]->UnRegisterEvent(eventIndex);
    }

    return GSE_SUCCESS;
}

void EventThreadPool::clear(ThreadID id)
{
    if (m_threadPool.size() <= id)
    {
        LOG_ERROR("invalid thread id %d >= pool size %d", id, m_threadPool.size());
        return;
    }

    if (id < 0)
    {
        for (int i = 0; i < m_threadPool.size(); i++)
        {
            m_threadPool[i]->clear();
        }
    }
    else
    {
        m_threadPool[id]->clear();
    }
}

int EventThreadPool::start()
{
    startThreadPool();

    return GSE_SUCCESS;
}

void EventThreadPool::stop()
{
    stopThreadPool();
}

void EventThreadPool::initThreadPool()
{
    for (int i = 0; i < m_threadNum; i++)
    {
        EventThread* pEventThread = new EventThread();
        if (pEventThread == NULL)
        {
            LOG_ERROR("create gse event thread fail");
            ::abort();
        }

        m_threadPool.push_back(pEventThread);
    }
}

void EventThreadPool::distroyThreadPool()
{
    for (int i = 0; i < m_threadPool.size(); i++)
    {
        EventThread* pEventThread = m_threadPool[i];
        delete pEventThread;
    }
}

void EventThreadPool::startThreadPool()
{
    for (int i = 0; i < m_threadNum; i++)
    {
        EventThread* pEventThread = m_threadPool[i];
        pEventThread->Start();
    }
}

void EventThreadPool::stopThreadPool()
{
    for (int i = 0; i < m_threadNum; i++)
    {
        EventThread* pEventThread = m_threadPool[i];
        pEventThread->Stop();
        pEventThread->Join();
    }
}

} // namespace data
} // namespace gse
