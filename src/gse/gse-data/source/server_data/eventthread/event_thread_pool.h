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

#ifndef _GSE_EVENT_THREAD_POOL_H_
#define _GSE_EVENT_THREAD_POOL_H_

#include "event_thread.h"
namespace gse {
namespace data {

#define Thread_All -1
typedef int32_t ThreadID;

class EventThreadPool
{
public:
    EventThreadPool();
    EventThreadPool(uint32_t threadNum);
    ~EventThreadPool();

    int regReadEvent(ThreadID id, EVENT_HANDLER eventHandler, evutil_socket_t fd, void* args, bool edgeTriger = false);
    int regTimer(ThreadID id, EVENT_HANDLER eventHandler, void* args, int second);
    int unRegEvent(ThreadID id, EventIndex eventIndex);
    int resetTimer(ThreadID id, EventIndex eventIndex, int second);
    void clear(ThreadID id);

    int start();
    void stop();

private:
    void initThreadPool();
    void distroyThreadPool();
    void startThreadPool();
    void stopThreadPool();

protected:
    ThreadID m_threadNum;

private:
    std::vector<EventThread*> m_threadPool;
};

} // namespace data
} // namespace gse
#endif
