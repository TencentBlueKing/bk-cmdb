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

# ifndef _GSE_EVENT_THREAD_H_
# define _GSE_EVENT_THREAD_H_

/**
 * @brief ?? libevent ?????
 * @auth  uriwang
 */


#include "thread/thread.h"
#include "gseEventManager.h"
namespace gse { 
namespace dataserver {

namespace rgse{

class GseEventThread:public gse::thread::Thread
{

public:
	GseEventThread(GseEventManager * eventManager);
	GseEventThread();
    ~GseEventThread();

public:
    /*!
     * @brief: register read event
     * @param[in] eventHandler: event hanler
     * @param[in] fd: notify file descriptor
     * @param[in] args: arguments
     * @param[in] edgeTriger: true for edge-triggered, false for level-triggered
     * @return: if fails, return -1; otherwise, return non-negtive interger as event index
     */
    int registerReadPersistEvent(EVENT_HANDLER eventHandler, evutil_socket_t fd, void* args, bool edgeTriger = false);

    /*!
     * @brief: register timer event
     * @param[in] eventHandler: event hanler
     * @param[in] args: arguments
     * @param[in] second: timer interval in seconds
     * @return: if fails, return -1; otherwise, return non-negtive interger as event index
     */
    int registerTimerPersistEvent(EVENT_HANDLER eventHandler, void* args, int second);
    int unRegisterEvent(EventIndex eventIndex);

    int resetEvent(EventIndex eventIndex, int second);
    int activeEvent(EventIndex eventIndex);

public:
	virtual int run();
    virtual int ThreadFun();
    virtual void stop();

public:
    int syncWait();
    void clear();

private:
    rgse::GseEventManager m_eventManager;

private:
    gse::thread::ThreadId m_threadId;
    event_base*           m_ptrEventBase;
    volatile bool         m_isNeedRunning;
};

}

}
}
#endif
