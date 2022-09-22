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

#ifndef _GSE_BKDATA_TIME_CENTER_H__

#define _GSE_BKDATA_TIME_CENTER_H__

#include <event2/event.h>

//#include "gseTime.h"

#include "eventthread/event_thread.h"
#include "safe/lock.h"
#include "tools/thread.h"

namespace gse {
namespace data {
/**
 *@brief 时间中心，统一获取时间的位置避免高频获取时间带来性能的开销和时获取到的时间不准确
 */
class TimeCenter : public gse::tools::thread::Thread
{

public:
    TimeCenter();
    ~TimeCenter();

public:
    /**
     * @brief 获取时间
     */
    time_t GetDateTime();
    /**
     * @brief 获取监控tag 时间戳
     */
    time_t GetTagTime();
    /**
     * @brief 获取标准时间格式， eg：2017:01:02T12:00:00+08:00
     * @param timestamp 输出参数
     */
    void StandardTimestampString(std::string& timestamp);

public:
    /**
     * @brief 单例接口
     */
    static TimeCenter* Instance();

public:
    virtual int Init();
    virtual int UnInit();
    virtual int ThreadFun();

private:
    static void timeUpdateEventHandler(evutil_socket_t fd, short which, void* v);
    void updateTime();

private:
    /**
     * @brief 时间
     */
    time_t m_dataTime;
    /**
     * @brief 监控的tag
     */
    time_t m_tagTime;
    /**
     *@brief 停止标记
     */
    bool m_isNeedStop;
    /**
     *@brief 同步读写锁
     */
    gse::safe::RWLock m_rwlock;
    /**
     * @brief 时间戳字符串
     */
    std::string m_dateTimeStr;

    event_base* m_baseEv;
    event* m_updateTimeEvent;
};

} // namespace data
} // namespace gse
#endif
