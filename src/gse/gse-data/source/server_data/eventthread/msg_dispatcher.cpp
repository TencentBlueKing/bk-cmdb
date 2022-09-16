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

#include "msg_dispatcher.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/error.h"
#include "tools/macros.h"
#include "tools/pipe.h"
#include <sys/types.h>
#include <sys/uio.h>
#include <unistd.h>
namespace gse {
namespace data {

namespace rgse {

MsgDispatcher::MsgDispatcher(GseMsgHandlerIf *handler)
    : m_threadPool(handler->getWorkerNum()), m_handler(handler), m_nextThread(0)
{
}

MsgDispatcher::~MsgDispatcher()
{
    if (m_handler != NULL)
        delete m_handler;
}

void MsgDispatcher::localNotifyHandler(int fd, short witch)
{
    MsgDispItem *pItem = NULL;

    int nreadcnt = ::read(fd, &pItem, sizeof(pItem));
    if (nreadcnt < 0)
    {
        int iErrorNo = gse_errno;
        LOG_ERROR("worker failed to read notify message. ret[%d], errno[%d], errmsg[%s]",
                  nreadcnt, iErrorNo, SAFE_CSTR(gse::tools::error::ErrnoToStr(iErrorNo).c_str()));
        return;
    }

    if (pItem == NULL)
    {
        LOG_ERROR("worker read notify NULL");
        return;
    }

    if (m_handler != NULL)
    {
        m_handler->msgHandler(pItem->m_msg, pItem->m_len, pItem->m_custom);
    }

    delete pItem;
}

void MsgDispatcher::notifyHandler(int fd, short witch, void *arg)
{
    MsgDispatcher *pDispatcher = reinterpret_cast<MsgDispatcher *>(arg);
    pDispatcher->localNotifyHandler(fd, witch);
}

int MsgDispatcher::notify(char *pMsg, uint32_t len, void *arg)
{
    MsgDispItem *pItem = new MsgDispItem();
    pItem->m_msg = pMsg;
    pItem->m_len = len;
    pItem->m_custom = arg;

    ThreadID id = m_handler->getNextWorker();
    m_nextThread = id % m_handler->getWorkerNum();

    if (::write(m_notifyWriteFd[id], &pItem, sizeof(pItem)) < 0)
    {
        LOG_ERROR("notify worker [%d] error", m_nextThread);
        delete pItem;
        return GSE_SYSTEMERROR;
    }

    return GSE_SUCCESS;
}

int MsgDispatcher::start()
{
    int ret;
    ret = init();
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("init fail");
        return ret;
    }

    ret = m_threadPool.start();
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("start fail");
        return ret;
    }

    return GSE_SUCCESS;
}

void MsgDispatcher::stop()
{
    m_threadPool.stop();
    unInit();
}

int MsgDispatcher::init()
{
    uint32_t workerNum = m_handler->getWorkerNum();

    for (int i = 0; i < workerNum; i++)
    {
        gse::tools::pipe::pipe_t fd[2];
        if (gse::tools::pipe::CreatePipe(fd) < 0)
        {
            LOG_ERROR("create pipe err %d", errno);
            unInit();
            return GSE_SYSTEMERROR;
        }

        m_notifyReadFd.push_back(fd[0]);
        m_notifyWriteFd.push_back(fd[1]);

        m_threadPool.regReadEvent(i, reinterpret_cast<EVENT_HANDLER>(MsgDispatcher::notifyHandler), fd[0], this, false);
    }

    return GSE_SUCCESS;
}

void MsgDispatcher::unInit()
{
    for (int i = 0; i < m_notifyReadFd.size(); i++)
    {
        gse::tools::pipe::ClosePipe(m_notifyReadFd[i]);
        gse::tools::pipe::ClosePipe(m_notifyWriteFd[i]);
    }

    m_threadPool.clear(Thread_All);
}

MsgDispatcher::MsgDispItem::MsgDispItem()
    : m_custom(NULL), m_len(0), m_msg(NULL)
{
}

MsgDispatcher::MsgDispItem::~MsgDispItem()
{
}

} // namespace rgse

} // namespace data
} // namespace gse
