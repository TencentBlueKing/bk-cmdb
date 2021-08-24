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

#include "gseEventManager.h"

#include "tools/atomic.h"
#include "tools/macros.h"
#include "bbx/gse_errno.h"
namespace gse { 
namespace dataserver {
namespace rgse
{

GseEventManager::GseEventManager()
{
    m_index = 0;
}

GseEventManager::~GseEventManager()
{
    clear();
}

int GseEventManager::regEvent(EventBase base, evutil_socket_t fd, EVENT_HANDLER handler, void * arg, bool isEdge /*= false*/, bool isPersist /*= true*/)
{
    struct event* ptrEvent = NULL;
    int flag = 0;

    if(base == NULL)
    {
        return GSE_ERROR;
    }

    if(isPersist)
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
        return GSE_EVENT_CREATE_FAIL;
    }

    if (-1 == event_add(ptrEvent, NULL))
    {
        event_free(ptrEvent);
        return GSE_EVENT_REGISTER_FAIL;
    }

    gse::tools::atomic::AtomAdd(&m_index);
    m_eventMap.Push(m_index, ptrEvent);
    return m_index;
}

int GseEventManager::regTimer(EventBase base, EVENT_HANDLER handler, void * arg, int second)
{
    struct event* ptrEvent = NULL;

    if(base == NULL)
    {
        return GSE_ERROR;
    }

    struct timeval timespan;
    timespan.tv_sec = second;
    timespan.tv_usec = 0;

    ptrEvent = event_new(base, -1, EV_PERSIST, handler, arg);
    if (NULL == ptrEvent)
    {
        return GSE_EVENT_CREATE_FAIL;
    }

    if (-1 == event_add(ptrEvent, &timespan))
    {
        event_free(ptrEvent);
        return GSE_EVENT_REGISTER_FAIL;
    }

    gse::tools::atomic::AtomAdd(&m_index);
    m_eventMap.Push(m_index, ptrEvent);
    return m_index;
}

int GseEventManager::unregEvent(EventIndex index)
{
    struct event* ptrEvent =  NULL;
    m_eventMap.Pop(index, ptrEvent);
    if(NULL != ptrEvent)
    {
        event_free(ptrEvent);
        return GSE_SUCCESS;
    }

    return GSE_EVENT_EMPTY;
}

int GseEventManager::resetTimer(EventIndex index, int second)
{
    struct event* ptrEvent = NULL;

    struct timeval timespan;
    timespan.tv_sec = second;
    timespan.tv_usec = 0;

    if( m_eventMap.Find(index, ptrEvent) )
    {
        if( -1 == event_add(ptrEvent, &timespan))
        {
            return GSE_EVENT_REGISTER_FAIL;
        }
        else
        {
            return GSE_SUCCESS;
        }
    }

    return GSE_EVENT_EMPTY;
}

void GseEventManager::clear()
{
    m_eventMap.Foreach(GseEventManager::dropEvent, this);
    m_eventMap.Clear();
}

void GseEventManager::dropEvent(EventIndex index, struct event * event, void * arg)
{
    event_free(event);
}

int GseEventManager::activeEvent(EventIndex index)
{
    struct event* ptrEvent = NULL;
    if (m_eventMap.Find(index, ptrEvent))
    {
        event_active(ptrEvent, EV_TIMEOUT, 1);
        return GSE_SUCCESS;
    }

    return GSE_EVENT_EMPTY;
}

}

}
}
