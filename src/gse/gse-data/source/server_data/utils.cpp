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

#include <string.h>


#include <unistd.h>
#include <sys/sysinfo.h>
#include <sys/statfs.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

#include <sys/socket.h>
#include <netinet/in.h>
#include <linux/sockios.h>
#include <net/if.h>
#include <sys/ioctl.h>

#include <linux/ethtool.h>

#include "utils.h"
#include "log/log.h"
#include "tools/net.h"
//#include "warn_report.h"
namespace gse { 
namespace data {

void alarm(const string &ip, const string content, int warnID)
{
    LOG_ERROR("send warn failed , ip:%s, content:%s, warnid:%d", ip.c_str(), content.c_str(), warnID);
}


int splitString(char* target, string delimiter, vector<int>& outVec)
{
    if (NULL == target)
    {
        return 0;
    }

    char* ptr = strtok(target, delimiter.c_str());
    while (NULL != ptr)
    {
        outVec.push_back(::atoi(ptr));
        ptr = strtok(NULL, delimiter.c_str());
    }

    return outVec.size();
}


int GetNetDevSpeed(const char* devName)
{
    struct ifreq ifr, *ifrp;  // 接口请求结构
    int fd;  // to  access socket  通过socket访问网卡的 文件描述符号fd

    /* Setup our control structures. */
    memset(&ifr, 0, sizeof(ifr));
    strcpy(ifr.ifr_name, devName);

    /* Open control socket. */
    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0)
    {
        return -1;
    }

    int err;
    struct ethtool_cmd ep;
    ep.cmd = ETHTOOL_GSET; // ethtool-copy.h:380:#define ETHTOOL_GSET         0x00000001 /* Get settings. */
    ifr.ifr_data = (caddr_t) & ep; //   caddr_t 是void类型，而这句话是什么意思
    err = ioctl(fd, SIOCETHTOOL, &ifr); //  int ioctl(int handle, int cmd,[int *argdx, int argcx]);
    close(fd);
    if (err != 0)
    {
        return -1;
    }

    return ep.speed;
}


uint32_t stringHash(const char *str, size_t len)
{
    /*
      * djb2 is one of the best string hash functions
      */

    uint32_t hash = 5381;
    int c;

    for(int i = 0; i < len; i++)
    {
        c = str[i];
        hash = ((hash << 5) + hash) + c; /* hash * 33 + c */
    }

    return hash;
}


}
}
