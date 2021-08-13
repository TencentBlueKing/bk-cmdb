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

#include "gseEventThread.h"
#include "tools/macros.h"
#include "bbx/gse_errno.h"
#include "tools/system.h"
#include <stdio.h>
#include <stdlib.h>
#include <vector>
#include "log/log.h"
namespace gse { 
namespace dataserver {

namespace rgse
{

GseEventThread::GseEventThread()
{
    m_threadId = 0;
    m_ptrEventBase = event_base_new();
    m_isNeedRunning = true;
}
int GseEventThread::ThreadFun()
{
   m_isNeedRunning = true;

    if(m_ptrEventBase == NULL)
    {
        LOG_FATAL("Gse Event Thread base event is NULL");
        return GSE_ERROR;
    }

    do
    {
        //
        // When it is done, event_base_loop() returns 0 if it exited normally,
        // -1 if it exited because of some unhandled error in the backend,
        // and 1 if it exited because there were no more pending or active events.
        //

        if(0 == event_base_dispatch(m_ptrEventBase))
        {
            LOG_INFO("exit from event manager run function");
            break;
        }
        else
        {
            LOG_INFO("exit from event manager fun function , loop continue");
            gse::tools::system::Sleep(1);
        }

    }
    while (m_isNeedRunning);

    return GSE_SUCCESS;
}
int GseEventThread::run()
{
    m_isNeedRunning = true;

    if(m_ptrEventBase == NULL)
    {
        LOG_FATAL("Gse Event Thread base event is NULL");
        return GSE_ERROR;
    }

    do
    {
        //
        // When it is done, event_base_loop() returns 0 if it exited normally,
        // -1 if it exited because of some unhandled error in the backend,
        // and 1 if it exited because there were no more pending or active events.
        //

        if(0 == event_base_dispatch(m_ptrEventBase))
        {
            LOG_INFO("exit from event manager run function");
            break;
        }
        else
        {
            LOG_INFO("exit from event manager fun function , loop continue");
            gse::tools::system::Sleep(1);
        }

    }
    while (m_isNeedRunning);

    return GSE_SUCCESS;
}
int GseEventThread::registerReadPersistEvent(EVENT_HANDLER eventHandler, evutil_socket_t fd, void* args, bool edgeTriger /*= false*/)
{
    return m_eventManager.regEvent(m_ptrEventBase, fd, eventHandler, args, edgeTriger, true);
}
int GseEventThread::registerTimerPersistEvent(EVENT_HANDLER eventHandler, void * args, int second)
{
    return m_eventManager.regTimer(m_ptrEventBase, eventHandler, args, second);
}
int GseEventThread::unRegisterEvent(EventIndex eventIndex)
{
    return m_eventManager.unregEvent(eventIndex);
}
int GseEventThread::resetEvent(EventIndex eventIndex, int second)
{
    return m_eventManager.resetTimer(eventIndex, second);
}
int GseEventThread::syncWait()
{
    return event_base_dispatch(m_ptrEventBase);
}

void GseEventThread::stop()
{
    m_isNeedRunning = false;
    if (NULL != m_ptrEventBase)
    {
        event_base_loopexit(m_ptrEventBase,NULL);
    }
}

void GseEventThread::clear()
{
    m_eventManager.clear();

}
GseEventThread::~GseEventThread()
{
    if (NULL != m_ptrEventBase)
    {
        event_base_free(m_ptrEventBase);
    }
}

int GseEventThread::activeEvent(EventIndex eventIndex)
{
    return m_eventManager.activeEvent(eventIndex);
}

}

}
}
