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

#include "event_thread.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/macros.h"
#include "tools/system.h"
#include <stdio.h>
#include <stdlib.h>
#include <vector>
namespace gse {
namespace data {

EventThread::EventThread()
    : m_threadId(0), m_isNeedRunning(true)
{
    m_ptrEventBase = event_base_new();
}

int EventThread::ThreadFun()
{
    m_isNeedRunning = true;

    if (m_ptrEventBase == NULL)
    {
        LOG_FATAL("Gse Event Thread base event is NULL");
        return GSE_ERROR;
    }

    m_threadId = gse::tools::thread::GetThreadID();

    LOG_DEBUG("event thread(%d) start", m_threadId);
    do
    {
        //
        // When it is done, event_base_loop() returns 0 if it exited normally,
        // -1 if it exited because of some unhandled error in the backend,
        // and 1 if it exited because there were no more pending or active events.
        //

        if (0 == event_base_dispatch(m_ptrEventBase))
        {
            LOG_INFO("exit from event manager run function");
            break;
        }
        else
        {
            LOG_INFO("exit from event manager fun function , loop continue");
            gse::tools::system::Sleep(1);
        }

    } while (m_isNeedRunning);

    LOG_INFO("event thread(%d) exit", m_threadId);
    //    std::unique_lock<std::mutex> lck(m_syncWaitLock);
    //    m_syncWait.notify_all();

    return GSE_SUCCESS;
}

int EventThread::RegisterReadPersistEvent(EVENT_HANDLER eventHandler, evutil_socket_t fd, void* args, bool edgeTriger /*= false*/)
{
    return m_eventManager.regEvent(m_ptrEventBase, fd, eventHandler, args, edgeTriger, true);
}
int EventThread::RegisterTimerPersistEvent(EVENT_HANDLER eventHandler, void* args, int second)
{
    return m_eventManager.regTimer(m_ptrEventBase, eventHandler, args, second * 1000);
}

int EventThread::RegisterMsecTimerPersistEvent(EVENT_HANDLER eventHandler, void* args, int msec)
{
    return m_eventManager.regMsTimer(m_ptrEventBase, eventHandler, args, msec);
}

int EventThread::UnRegisterEvent(EventIndex eventIndex)
{
    return m_eventManager.unregEvent(eventIndex);
}
int EventThread::ResetEvent(EventIndex eventIndex, int second)
{
    return m_eventManager.resetTimer(eventIndex, second);
}
int EventThread::SyncWait()
{
    return event_base_dispatch(m_ptrEventBase);
}

void EventThread::Stop()
{
    m_isNeedRunning = false;
    if (NULL != m_ptrEventBase)
    {
        event_base_loopexit(m_ptrEventBase, NULL);
        //        std::unique_lock<std::mutex> lck(m_syncWaitLock);
        //        while (m_syncWait.wait_for(lck, std::chrono::seconds(5)) == std::cv_status::timeout)
        //        {
        //            LOG_WARN("wait stop event thread timeout");
        //            break;
        //        }
    }
}

void EventThread::clear()
{
    m_eventManager.clear();
}

EventThread::~EventThread()
{
    if (NULL != m_ptrEventBase)
    {
        event_base_free(m_ptrEventBase);
        m_ptrEventBase = NULL;
    }
}

int EventThread::ActiveEvent(EventIndex eventIndex)
{
    return m_eventManager.activeEvent(eventIndex);
}

} // namespace data
} // namespace gse
