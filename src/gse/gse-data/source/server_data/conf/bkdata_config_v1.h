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

//
// 兼容V1 版本DS的配置管理逻辑
//

#ifndef _GSE_DATA_CONFIG_BKDATA_V1_H_
#define _GSE_DATA_CONFIG_BKDATA_V1_H_

#include "log/log.h"
#include "safe/lock.h"
#include "tools/atomic.h"
#include "tools/strings.h"
#include "tools/time.h"
#include "json/json.h"
#include <string>
#include <vector>
namespace gse {
namespace data {

#define VERSION_DATAID 6696
#define DATASERVER_FLOW_DATAID 238
#define DATASERVER_EVENT_DATAID 293
#define COLLOCTOR_OP_TYPE_AGENT_FLOW 294
#define DATASERVER_MONITOR_TAG 295
#define DATASERVER_MONITOR_TAG_FLOW_DATAID 296

/**
 *@brief DataID 对象定义
 */
class DataID
{
public:
    DataID()
        : m_storeSys(-1), m_bizId(-1), m_partitions(1), m_dataSet(""), m_keyTopic("")
    {
        m_setDeleteTimestamp = m_nextPartition = m_dataId = m_clusterIndex = 0;
        m_next = NULL;
        m_tenant = "";
        m_namespace = "";
        m_persistent = "persistent";
        m_optype = 0;
        m_type = 0;
    }

    ~DataID()
    {
        // free linked DataId
        if (m_next != NULL)
        {
            delete m_next;
        }
        m_next = NULL;
    }
    std::string ToString();

public:
    inline int nextPartion()
    {
        int ret = 0;
        if (m_partitions > 0)
        {
            ret = abs(gse::tools::atomic::AtomAddAfter(&m_nextPartition) % m_partitions);
            return ret;
        }
        LOG_WARN("data id[%d] partion invalid:%d", m_dataId, m_partitions);
        return 0;
    }
    inline void SetNeedDelete()
    {
        m_setDeleteTimestamp = gse::tools::time::GetUTCSecond();
    }
    inline bool IsNeedDelete()
    {
        return m_setDeleteTimestamp == 0 ? false : ((gse::tools::time::GetUTCSecond() - m_setDeleteTimestamp) > 60);
    }

public:
    /**
     *@brief 类型
     */
    int m_type;
    /**
     *@brief ops 使用
     */
    int m_optype; // use for opsserver
    /**
     *@brief 数据写入的存储，暂时仅支持kafka  类型为 0
     */
    int m_storeSys;

    // 当一个 dataid 与多个存储集群关联的时候，next 会被设置
    DataID *m_next;

    /**
     *@brief 业务id
     */
    int m_bizId;
    /**
     *@brief 分配的partion 数量
     */
    int m_partitions;
    /**
     *@brief 数据写入存储的集群编号
     */
    int m_clusterIndex;

    /**
     * @brief dataid 值
     *   6bit userid + 8bit check+ 18bit dataid
     */
    uint32_t m_dataId;

    /**
     * @brief 数据集
     */
    std::string m_dataSet;
    /**
     *@brief topic ，dataset+bizid
     */
    std::string m_keyTopic;

    std::string m_tenant;
    std::string m_namespace;
    std::string m_persistent;

private:
    /**
     * @brief 下一个备选的partition
     */
    int m_nextPartition;

    /**
     * @brief 读写锁
     */
    gse::safe::MutexLock m_mutex;

    int m_setDeleteTimestamp;
};

/**
 * @brief 定义 地址类型, 兼容 DS 1.0  的配置，需要被废弃
 */
typedef struct _AddressIP
{
    std::string m_ip;
    int m_port;
} AddressIP;
typedef std::vector<AddressIP> ADDRESS_IP;

/**
 * @brief server的基础配置信息， 兼容DS 1.0 的配置，需要被废弃
 */
typedef struct _BaseCfg
{
    /**
     * @brief tnm2 告警id
     */
    int m_warnId;
    /**
     * @ kafka 队列上限
     */
    int m_kafkaQueueMax;
    /**
     * @brief 单个日志文件大小上限
     */
    int m_logfileSize;
    /**
     * @brief 日志文件数量上限，达到此上限则回绕
     */
    int m_logfileNum;
    /**
     * @brief 服务器启动的线程数
     */
    int m_alliothread;
    /**
     * @brief 启动远程流水日志
     */
    bool m_enableRemoteStream;
    /**
     * @brief 设置是否启动兼容模式，如果设置为true 则表示仅使用逻辑分区，服务器不会向区域+城市节点注册自身信息
     */
    bool m_onlyUseLogicalSetting;

    /**
     * @brief 复合id
     */
    int m_composeId;

    // server config
    /**
     * @brief 服务器监听的IP地址
     */
    string m_listenIp;
    int m_listenPort;

    // vector<Remote_Addr> m_servers;
    int m_svrNum;
    int m_workerNum;
    bool m_serverUseSsl;
    bool m_clientUseSsl;

    /**
     * @brief 服务器所在的逻辑区域
     */
    std::string m_logicalId;
    /**
     * @brief 服务器所在的物理分区（大区）
     */
    std::string m_regId;
    /**
     * @brief 服务器所在的城市
     */
    std::string m_cityId;
    /**
     * @brief 日志级别
     */
    std::string m_logLevel;
    /**
     * @brief 地址信息
     */
    ADDRESS_IP m_addressIP;
    /**
     * @brief 有kakfa 信息的zk地址
     */
    Json::Value m_zookeeper;
    /**
     * @brief 有kafka 信息的zk 地址字符串
     */
    std::string m_zookeeperStr;

    // 1: get dataid and storage from /config/leaf/data/ and /config/leaf/kafka/
    // 2: get dataid and storage from /gse/config/etc/dataserver/data/ and /gse/config/etc/dataserver/storage/
#define DS_ZK_CONFIG_V1 1
#define DS_ZK_CONFIG_V2 2
    int m_zk_config_version;
} BaseCfg;

/*
 * parse zk node string to struct
 *
 * storage zk node format
 * [
 *     {
 *         "host":"127.0.0.1",
 *         "port":9092,
 *         "passwd":"",
 *         "type":1 // StorageType
 *     }
 * ]
 */
// int parseStorageNode(int clusterIndex, const std::string &nodeStr, std::vector<StorageConfigType> &storageConfigs);

/*
{
    "type":1,
    "biz_id":0,
    "cluster_index":2,
    "data_set":"name",
    "msg_system":1,
    "partition":1,
    "optype":1
}
*/
DataID *parseToDataID(const Json::Value &root);

/*
 * parse zk node string to struct
 * DataId should be free by users
 *
 * dataid zk node
 * {
 *     "type":1,
 *     "biz_id":0,
 *     "cluster_index":1,
 *     "data_set":"set",
 *     "msg_system":1,
 *     "partition":1
 * }
 *
 * */
DataID *parseToDataID(const std::string &dataInfo);

} // namespace data
} // namespace gse
#endif
