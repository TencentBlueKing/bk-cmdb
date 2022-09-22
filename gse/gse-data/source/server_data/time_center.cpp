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

#include "time_center.h"
#include "tools/time.h"
#include "log/log.h"

#include "bbx/gse_errno.h"
#include "tools/error.h"

#include "common/time.h"

namespace gse {
namespace data {


TimeCenter::TimeCenter()
{
    DateTime datetime;
    m_dataTime = datetime.UTCSecondsFrom1970();
    m_tagTime = m_dataTime;
    m_dateTimeStr = datetime.ToLocalDateTimeString();
    m_isNeedStop = false;
    m_updateTimeEvent = NULL;
    m_baseEv = NULL;
}
TimeCenter::~TimeCenter()
{}
TimeCenter* TimeCenter::Instance()
{
    static TimeCenter _instance;
    return &_instance;
}

void TimeCenter::StandardTimestampString(std::string& timestamp)
{

}

time_t TimeCenter::GetDateTime()
{
    time_t t_date_time;
    m_rwlock.RLock();
    t_date_time = m_dataTime;
    m_rwlock.UnLock();
    return t_date_time;
}
time_t TimeCenter::GetTagTime()
{
    time_t t_date_time;
    m_rwlock.RLock();
    t_date_time = m_tagTime;
    m_rwlock.UnLock();
    return t_date_time;
}

int TimeCenter::Init()
{
    return GSE_SUCCESS;
}

int TimeCenter::UnInit()
{
    m_isNeedStop = true;
    return GSE_SUCCESS;
}

void TimeCenter::updateTime()
{
    m_rwlock.WLock();
    m_dataTime = ::time(NULL);

    static int timeout = 60;
    if( timeout-- <= 0)
    {
        timeout = 60;
        m_tagTime = m_dataTime;
    }
    m_rwlock.UnLock();
    return;
}

void TimeCenter::timeUpdateEventHandler(evutil_socket_t fd, short which, void * v)
{
    TimeCenter *self = (TimeCenter*)v;
    self->updateTime();
    return;
}

int TimeCenter::ThreadFun()
{
    //
    // 独立线程运行，每一秒钟更新一次时间戳
    //
    m_baseEv = event_base_new();
    if (m_baseEv == NULL)
    {
        int tmp_errno = gse_errno;
        string errmsg = gse::tools::error::ErrnoToStr(tmp_errno);

        LOG_ERROR("timecenter failed to init base event, the error code is (%d), the error message is (%s)",
                  tmp_errno, errmsg.c_str());
        return GSE_ERROR;
    }

    m_updateTimeEvent = event_new(m_baseEv, -1, EV_PERSIST, TimeCenter::timeUpdateEventHandler, (void*)this);
    if (m_updateTimeEvent == NULL)
    {
        int tmp_errno = gse_errno;
        string errStr = gse::tools::error::ErrnoToStr(tmp_errno);
        LOG_ERROR("timecenter failed to register event for update time event, the error code is (%d), the error message is (%s)",
                  tmp_errno, errStr.c_str());
        return GSE_ERROR;
    }

    timeval timer_interval;
    timer_interval.tv_sec = 1;
    timer_interval.tv_usec = 0;

    if (-1 == event_add(m_updateTimeEvent, &timer_interval))
    {
        event_free(m_updateTimeEvent);
        return GSE_ERROR;
    }

    int rc = -1;
    while(rc < 0)
    {
        rc = event_base_loop(m_baseEv, 0);
        LOG_ERROR("module (timecenter) event base loop break abnormally, return code is (%d)", rc);
    }

    LOG_ERROR("module (timecenter) event base loop break abnormally, return code is (%d)", rc);
    return GSE_ERROR;
/*
    static int timeout = 60;
    while(!m_isNeedStop)
    {
        sleep(1);// 1 second
        m_dataTime = DateTime().UTCSecondsFrom1970();
        if( timeout-- <= 0)
        {
            timeout = 60;
            m_tagTime = m_dataTime;
        }
        //m_rwlock.wLock();
        //m_dateTimeStr = DateTime().ToLocalDateTimeString();
        //m_rwlock.unLock();
    }
*/
}
}
}


