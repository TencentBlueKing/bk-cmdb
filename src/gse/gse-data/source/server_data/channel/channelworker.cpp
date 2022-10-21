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

#include "channelworker.h"
#include "datachannel.h"
#include "dataserver.h"
#include "log/log.h"
#include "tools/error.h"
#include "tools/pipe.h"

namespace gse {
namespace data {

#define MAX_QUEUE_SIZE (40960 * 60)

ChannelWorker::ChannelWorker(int id, void *pParent)
{
    m_workId = id;
    m_pParent = pParent;

    m_notifyFd[0] = -1;
    m_notifyFd[1] = -1;

    m_evBase = NULL;
    m_notifyEvent = NULL;

    // m_ptrDataQueue = new boost::circular_buffer<DataCell *>(MAX_QUEUE_SIZE);
    m_ptrDataQueue = new moodycamel::ConcurrentQueue<DataCell *>();

    m_opsReportTime.tv_sec = 5;
    m_opsReportTime.tv_usec = 0;

    m_msgHandleTime.tv_sec = 0;
    m_msgHandleTime.tv_usec = 1000 * 100; // 100ms

    m_opsTimerEvent = NULL;
    m_msgTimerEvent = NULL;
    m_ptrOPSReport = NULL;
}

ChannelWorker::~ChannelWorker()
{
    if (m_notifyFd[0] != -1)
    {
        gse::tools::pipe::ClosePipe(m_notifyFd[0]);
        m_notifyFd[0] = -1;
    }

    if (m_notifyFd[1] != -1)
    {
        gse::tools::pipe::ClosePipe(m_notifyFd[1]);
        m_notifyFd[1] = -1;
    }

    cleanEvent();
    if (m_ptrDataQueue != NULL)
    {
        delete m_ptrDataQueue;
    }
}

void ChannelWorker::cleanEvent()
{
    if (m_notifyEvent != NULL)
    {
        event_del(m_notifyEvent);
        event_free(m_notifyEvent);
        m_notifyEvent = NULL;
    }

    if (m_evBase != NULL)
    {
        event_base_free(m_evBase);
        m_evBase = NULL;
    }
}

int ChannelWorker::Init()
{
    int ret = gse::tools::pipe::CreatePipe(m_notifyFd);
    if (ret < 0)
    {
        int iErrno = gse_errno;
        LOG_ERROR("channel worker[%d]: fail to create notify pipe. ret=[%d], errno=[%d], errmsg=[%s]",
                  m_workId, ret, iErrno, gse::tools::error::ErrnoToStr(iErrno).c_str());
        return GSE_SYSTEMERROR;
    }

    for (int i = 0; i < 2; i++)
    {
        if (!gse::tools::net::SocketSetnonblocking(m_notifyFd[i]))
        {
            int iErrno = gse_errno;
            LOG_ERROR("channel worker[%d]: fail to set notify pipe to non-blocking. errno=[%d], errmsg=[%s]",
                      m_workId, iErrno, gse::tools::error::ErrnoToStr(iErrno).c_str());
            return GSE_SYSTEMERROR;
        }
    }

    return GSE_SUCCESS;
}

void ChannelWorker::Notify(DataCell *pDataCell)
{
    if (NULL == pDataCell)
    {
        LOG_WARN("channel worker[%d]: the pointer of data cell is NULL", m_workId);
        return;
    }

    if (m_ptrDataQueue->size_approx() >= MAX_QUEUE_SIZE)
    {
        LOG_ERROR("enqueue channel worker queue failed, the queue exceeds the max limit(%d), queue msg count:%d", MAX_QUEUE_SIZE, m_ptrDataQueue->size_approx());
        delete pDataCell;
        pDataCell = NULL;
        return;
    }
    int max_retry_time = 0;
    while (!m_ptrDataQueue->enqueue(pDataCell))
    {
        max_retry_time++;
        if (max_retry_time > 5)
        {
            LOG_ERROR("enqueued failure, exceeding maximum retry times(%d)", max_retry_time);
            delete pDataCell;
            pDataCell = NULL;
            break;
        }
    }
}

int ChannelWorker::ThreadFun()
{
    int ret = -1;

    m_evBase = event_base_new();
    if (NULL == m_evBase)
    {
        LOG_ERROR("channel worker[%d]: fail to new event base object");
        return GSE_SYSTEMERROR;
    }

    m_msgTimerEvent = event_new(m_evBase, -1, EV_PERSIST, ChannelWorker::QueueMsgHandler, this);
    if (NULL == m_msgTimerEvent)
    {
        LOG_ERROR("channel worker[%d]: fail to new msg handle", m_workId);
        cleanEvent();
        return GSE_SYSTEMERROR;
    }

    ret = event_add(m_msgTimerEvent, &m_msgHandleTime);
    if (ret < 0)
    {
        LOG_ERROR("channel worker[%d]: fail to add msg handle vent. ret=[%d]", m_workId, ret);
        return GSE_SYSTEMERROR;
    }

    m_opsTimerEvent = event_new(m_evBase, -1, EV_PERSIST, ChannelWorker::OpsMsgHandler, this);
    if (NULL == m_opsTimerEvent)
    {
        LOG_ERROR("channel worker[%d]: fail to new ops event", m_workId);
        cleanEvent();
        return GSE_SYSTEMERROR;
    }

    ret = event_add(m_opsTimerEvent, &m_opsReportTime);
    if (ret < 0)
    {
        LOG_ERROR("channel worker[%d]: fail to add ops event. ret=[%d]", m_workId, ret);
        return GSE_SYSTEMERROR;
    }

    ret = event_base_loop(m_evBase, 0);
    LOG_FATAL("channel worker[%d]: event base loop abnomal break. ret=[%d]", m_workId, ret);

    return GSE_SUCCESS;
}

void ChannelWorker::stop()
{
    if (m_evBase != NULL)
    {
        event_base_loopexit(m_evBase, NULL);
    }
}

void ChannelWorker::QueueMsgHandler(int fd, short which, void *v)
{
    ChannelWorker *ptr_this = (ChannelWorker *)v;

    DataCell *pDataCell;

    while (ptr_this->m_ptrDataQueue->try_dequeue(pDataCell))
    {
        ptr_this->dealData(pDataCell);
    }
}

void ChannelWorker::ReportQueueSizeMetric()
{
    std::string metric_name = "worker_queune_msg_count";
    OPMetric::Instance().ReportQueneCount(metric_name, m_workId, m_ptrDataQueue->size_approx());
}

void ChannelWorker::OpsMsgHandler(int fd, short which, void *v)
{
    ChannelWorker *ptr_this = (ChannelWorker *)v;
    ptr_this->ReportQueueSizeMetric();
}

void ChannelWorker::dealData(DataCell *pDataCell)
{
    DataChannel *pDataChannel = (DataChannel *)m_pParent;
    int iRet = pDataChannel->DecodeMsg(pDataCell);
    if (GSE_SUCCESS != iRet)
    {
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        delete pDataCell;
        return;
    }

    if (pDataCell->GetErrorCode() == OPS_ERROR_TRANSFER)
    {
        DataServer::GetOpsReportClient()->TransferOpsData(pDataCell);
        delete pDataCell;
        return;
    }

    if (pDataCell->GetErrorCode() == OPS_ERROR_SKIP)
    {
        delete pDataCell;
        return;
    }

    pDataChannel->SetDataCellChannelidByServiceId(pDataCell);

    DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_DEALING_STATE));

    iRet = pDataChannel->ExportData(pDataCell);
    if (GSE_SUCCESS != iRet)
    {
        delete pDataCell;
        return;
    }

    delete pDataCell;
}

} // namespace data
} // namespace gse
