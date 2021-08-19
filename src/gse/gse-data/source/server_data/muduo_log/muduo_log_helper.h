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

#ifndef _MUDUOHELPER_H_
#define _MUDUOHELPER_H_


#include <string>

#include <muduo/base/AsyncLogging.h>

#include "dataStruct/safe_map.h"
#include "safe/lock.h"

namespace gse {
namespace dataserver {

class MuduoLogHelper
{
public:
	~MuduoLogHelper();

	static MuduoLogHelper& Instance();

	/**
	Describe:
		create the data sink
	Parameters:
		szPrefix		[in] the path where the data log file save
		max_size_file	[in] the max size of one file. its unit is MB
		max_num_file	[in] the max number of the log files
		flush_msec	[in] the time span decide when flush the data from cache to log file. if zero, will flush immediately
								its unit is microseconds
	*******************************************************************************************************/
    int CreateDataSink(const std::string& szPrefix = "./logs/data/", unsigned int max_size_file = 10, unsigned int max_num_file = 20, int flush_msec = 500);
    static void WriteWithoutFmt(int sink,  const char* buff);
    static void asyncOutput(int sink_id, const char* log, int len);

private:
	MuduoLogHelper();

	MuduoLogHelper(const MuduoLogHelper&);
    MuduoLogHelper& operator=(const MuduoLogHelper&);
    muduo::AsyncLogging* GetLogSink(int sink_id);

    gse::safe::RWLock m_rwLock;
    std::map<int, muduo::AsyncLogging *> m_sinks;
    
};

#define STRINGIZE2(s) #s
#define STRINGIZE(s)  STRINGIZE2(s)


#define MUDUO_LOG_DATA(sink, buff)\
    MuduoLogHelper::WriteWithoutFmt(sink, buff)
#endif


}
}
