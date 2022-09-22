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

#include "muduo_log_helper.h"
#include <muduo/base/Logging.h>
#include <muduo/base/Timestamp.h>

#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <sys/time.h>
#include <time.h>
#include <sys/stat.h>
#include <unistd.h>
#include <map>

namespace gse {
namespace data {
using namespace std;
MuduoLogHelper::MuduoLogHelper()
{

}

MuduoLogHelper::~MuduoLogHelper()
{
    m_rwLock.WLock();
    std::map<int, muduo::AsyncLogging*>::iterator it =  m_sinks.begin();
    for (; it != m_sinks.end(); it++)
    {
        delete it->second;
    }
    m_rwLock.UnLock();
    m_sinks.clear();
}

MuduoLogHelper& MuduoLogHelper::Instance()
{
	static MuduoLogHelper single;
	return single;
}

void MuduoLogHelper::asyncOutput(int sink_id, const char* log, int len)
{
    muduo::AsyncLogging* ptr_async_log = NULL;
    ptr_async_log = MuduoLogHelper::Instance().GetLogSink(sink_id);

    if (ptr_async_log != NULL)
    {
        ptr_async_log->append(log, len);
    }
}

muduo::AsyncLogging* MuduoLogHelper::GetLogSink(int sink_id)
{
    m_rwLock.RLock();
    std::map<int, muduo::AsyncLogging*>::iterator it = m_sinks.find(sink_id);
    if (it != m_sinks.end())
    {
        m_rwLock.UnLock();
        return it->second;
    }
    m_rwLock.UnLock();
    return NULL;
}

//not
int MuduoLogHelper::CreateDataSink(const std::string& szPrefix/* = "./logs/data/"*/, unsigned int max_size_file/* = 10*/, unsigned int max_num_file/* = 20*/, int flush_msec/* = 500*/)
{
    static int sink_id = 0;
    static int sink_id_count = 0;
    muduo::AsyncLogging *ptr_log = NULL;
    ptr_log = new muduo::AsyncLogging(szPrefix.c_str(), max_size_file, (int)max_num_file, flush_msec, 3);
    ptr_log->start();
    m_rwLock.WLock();
    sink_id++;
    m_sinks.insert(std::pair<int, muduo::AsyncLogging *>(sink_id, ptr_log));
    sink_id_count++;
    m_rwLock.UnLock();
	muduo::Logger::setOutput(&MuduoLogHelper::asyncOutput);
    return sink_id;
}

void MuduoLogHelper::WriteWithoutFmt(int sink, const char* buff)
{
    SINK_LOG(sink) << buff;
}

}
}
