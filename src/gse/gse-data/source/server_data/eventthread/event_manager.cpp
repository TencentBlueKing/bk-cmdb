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

#include "event_manager.h"

#include "bbx/gse_errno.h"
#include "tools/atomic.h"
#include "tools/macros.h"

namespace gse {
namespace data {

EventManager::EventManager()
{
    m_index = 0;
}

EventManager::~EventManager()
{
    clear();
}

int EventManager::regEvent(EventBase base, evutil_socket_t fd, EVENT_HANDLER handler, void* arg, bool isEdge /*= false*/, bool isPersist /*= true*/)
{
    struct event* ptrEvent = NULL;
    int flag = 0;

    if (base == NULL)
    {
        return GSE_ERROR;
    }

    if (isPersist)
    {
        flag = EV_PERSIST;
    }

    if (!isEdge)
    {
        ptrEvent = event_new(base, fd, EV_READ | flag, handler, arg);
    }
    else
    {
        ptrEvent = event_new(base, fd, EV_READ | flag | EV_ET, handler, arg);
    }

    if (NULL == ptrEvent)
    {
        return GSE_ERROR;
    }

    if (-1 == event_add(ptrEvent, NULL))
    {
        event_free(ptrEvent);
        return GSE_ERROR;
    }

    gse::tools::atomic::AtomAdd(&m_index);
    m_eventMap.Push(m_index, ptrEvent);
    return m_index;
}

int EventManager::regMsTimer(EventBase base, EVENT_HANDLER handler, void* arg, int ms)
{
    struct event* ptrEvent = NULL;

    if (base == NULL)
    {
        return GSE_ERROR;
    }

    struct timeval timespan;
    int sec = ms / 1000;
    int usec = (ms % 1000) * 1000;

    timespan.tv_sec = sec;
    timespan.tv_usec = usec;

    ptrEvent = event_new(base, -1, EV_PERSIST, handler, arg);
    if (NULL == ptrEvent)
    {
        return GSE_ERROR;
    }

    if (-1 == event_add(ptrEvent, &timespan))
    {
        event_free(ptrEvent);
        return GSE_ERROR;
    }

    gse::tools::atomic::AtomAdd(&m_index);
    m_eventMap.Push(m_index, ptrEvent);
    return m_index;
}

int EventManager::regTimer(EventBase base, EVENT_HANDLER handler, void* arg, int second)
{
    return regMsTimer(base, handler, arg, second);
}

int EventManager::unregEvent(EventIndex index)
{
    struct event* ptrEvent = NULL;
    m_eventMap.Pop(index, ptrEvent);
    if (NULL != ptrEvent)
    {
        event_free(ptrEvent);
        return GSE_SUCCESS;
    }

    return GSE_ERROR;
}

int EventManager::resetTimer(EventIndex index, int second)
{
    struct event* ptrEvent = NULL;

    struct timeval timespan;
    timespan.tv_sec = second;
    timespan.tv_usec = 0;

    if (m_eventMap.Find(index, ptrEvent))
    {
        if (-1 == event_add(ptrEvent, &timespan))
        {
            return GSE_ERROR;
        }
        else
        {
            return GSE_SUCCESS;
        }
    }

    return GSE_ERROR;
}

void EventManager::clear()
{
    m_eventMap.Foreach(EventManager::dropEvent, this);
    m_eventMap.Clear();
}

void EventManager::dropEvent(EventIndex index, struct event* event, void* arg)
{
    event_free(event);
}

int EventManager::activeEvent(EventIndex index)
{
    struct event* ptrEvent = NULL;
    if (m_eventMap.Find(index, ptrEvent))
    {
        event_active(ptrEvent, EV_TIMEOUT, 1);
        return GSE_SUCCESS;
    }

    return GSE_ERROR;
}

} // namespace data
} // namespace gse
