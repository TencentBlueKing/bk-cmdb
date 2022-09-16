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

#ifndef _GSE_EVENT_THREAD_H_
#define _GSE_EVENT_THREAD_H_

#include <condition_variable>
#include <mutex>

#include "event_manager.h"
#include "tools/thread.h"

namespace gse {
namespace data {

class EventThread final : public gse::tools::thread::Thread
{

public:
    explicit EventThread(EventManager* eventManager);
    EventThread();
    ~EventThread();

public:
    int RegisterReadPersistEvent(EVENT_HANDLER eventHandler, evutil_socket_t fd, void* args, bool edgeTriger = false);
    int RegisterTimerPersistEvent(EVENT_HANDLER eventHandler, void* args, int second);
    int RegisterMsecTimerPersistEvent(EVENT_HANDLER eventHandler, void* args, int msec);
    int UnRegisterEvent(EventIndex eventIndex);

    int ResetEvent(EventIndex eventIndex, int second);
    int ActiveEvent(EventIndex eventIndex);

public:
    virtual int ThreadFun();
    virtual void Stop();

public:
    int SyncWait();
    void clear();

private:
    EventManager m_eventManager;

private:
    gse::tools::thread::ThreadId m_threadId;
    event_base* m_ptrEventBase;
    volatile bool m_isNeedRunning;

    std::condition_variable m_syncWait;
    std::mutex m_syncWaitLock;
};

} // namespace data
} // namespace gse
#endif
