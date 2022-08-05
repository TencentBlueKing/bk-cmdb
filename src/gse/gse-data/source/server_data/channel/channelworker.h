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

#ifndef _GSE_DATA_CHANNELWORKER_H_
#define _GSE_DATA_CHANNELWORKER_H_

#include <event.h>

#include "concurrentqueue/concurrentqueue.h"
#include "datacell.h"
#include "datastruct/safe_queue.h"
#include "ops/ops.h"
#include "safe/lock.h"
#include "tools/thread.h"

namespace gse {
namespace data {

class ChannelWorker : public gse::tools::thread::Thread
{
public:
    ChannelWorker(int id, void* pParent);
    virtual ~ChannelWorker();

public:
    int Init();
    int ThreadFun();
    void stop();

    void Notify(DataCell* pDataCell);

public:
    inline void SetOPS(OpsCollection* ptrOPSReport)
    {
        m_ptrOPSReport = ptrOPSReport;
    }

private:
    static void QueueMsgHandler(int fd, short which, void* v);
    static void OpsMsgHandler(int fd, short which, void* v);

private:
    void cleanEvent();
    void dealData(DataCell* pDataCell);
    void ReportQueueSizeMetric();

private:
    int m_notifyFd[2];

    event_base* m_evBase;
    struct event* m_notifyEvent;

    struct event* m_opsTimerEvent;
    struct event* m_msgTimerEvent;

    struct timeval m_msgHandleTime;
    struct timeval m_opsReportTime;

    int m_workId;
    void* m_pParent;
    gse::safe::RWLock m_dataQueueLock;
    moodycamel::ConcurrentQueue<DataCell*>* m_ptrDataQueue;

private:
    OpsCollection* m_ptrOPSReport;
};

} // namespace data
} // namespace gse

#endif
