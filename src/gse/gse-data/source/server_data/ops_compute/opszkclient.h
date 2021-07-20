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

#ifndef _DATASERVER_OPS_COMPUTE_OPSZKCLIENT_
#define _DATASERVER_OPS_COMPUTE_OPSZKCLIENT_

#include <string>
#include <vector>
#include "utilTools.h"
#include "conf/dataconf.h"
#include "zkClient.h"
#include "serverConfig.h"
namespace gse { 
namespace dataserver {
using namespace std;
/**
 * @brief 运营服务器的zk 配置管理器
 */
class OpsZKClient
{
public:
    OpsZKClient();
    ~OpsZKClient();

    /**
     * @brief 默认回调
     * @param type
     * @param state
     * @param path
     * @param wctx
     */
    static void configWatchFunc(int type, int state, const char * path, void * wctx);

    /**
     * @brief 初始化操作
     * @param host zookeeper 地址字符串
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int setup(const string &host);

    /**
     * @brief 同步运营系统的dataid
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int getAllDataIdNode();

    /**
     * @brief 获取dataid节点数据
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int getDataIdNode(int id, const string &path);

    /**
     * @brief 删除某用户的dataid
     * @param userId 用户id
     * @param dataId 数据id
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int deleteUserDataId(int userId, int dataId);

    /**
     * @brief 获取所有kakfa 集群的信息
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int getAllKafkaNode();

    /**
     * @brief 获取某用户的某个集群下的kafka 地址集合
     * @param userId 用户id
     * @param clusterIndex 集群id
     * @param path 节点路径
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int getKafkaNode(int userId, int clusterIndex, const string &nodePath);

    /**
     * @brief 删除某用户的某个集群
     * @param userId 用户id
     * @param clusterIndex 集群索引
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int deleteUserStorage(int userId, int clusterIndex);

private:
    /**
     * @brief 禁止拷贝
     */
    DISALLOW_COPY_AND_ASSIGN(OpsZKClient);

    /**
     * @brief 链接zk
     * @return = 0 操作成功
     *         < 0 操作失败
     */
    int connectToZk();

    ZkClient m_zkclient;
    string m_host;
};

}
}
#endif

