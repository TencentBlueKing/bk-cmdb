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

/*
    @discrip: async msg dispather 
    @author:robychen
    @time:20160405
*/

#ifndef _GSE_THREAD_MSG_DISPATCHER_H_
#define _GSE_THREAD_MSG_DISPATCHER_H_

#include "gseEventThreadPool.h"
namespace gse { 
namespace dataserver {

namespace rgse
{

class GseMsgHandlerIf
{
public:
    virtual ~GseMsgHandlerIf(){}
    virtual uint32_t getWorkerNum() = 0;
    virtual uint32_t getNextWorker() = 0;
    virtual void msgHandler(char * pMsg, int msgLen, void * arg) = 0;
};

class GseMsgDispatcher
{
public:
    GseMsgDispatcher(GseMsgHandlerIf * handler);
    ~GseMsgDispatcher();

    int notify(char * pMsg, uint32_t len, void * arg);

    int start();
    void stop();

    class MsgDispItem
    {
    public:
        MsgDispItem();
        ~MsgDispItem();
        
    public:
        uint32_t m_len;
        char * m_msg;
        void * m_custom;
    };

private:
    static void notifyHandler(int fd, short witch, void* arg);
    void localNotifyHandler(int fd, short witch);
    int init();
    void unInit();

private:
    GseEventThreadPool m_threadPool;
    vector<int> m_notifyWriteFd;
    vector<int> m_notifyReadFd;
    ThreadID m_nextThread;

    GseMsgHandlerIf * m_handler;
};

}

}
}
#endif

