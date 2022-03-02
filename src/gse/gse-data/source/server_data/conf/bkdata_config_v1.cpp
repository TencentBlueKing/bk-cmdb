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

#include "bkdata_config_v1.h"

#include "json/json.h"
#include "log/log.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"

#include "conf/confItem.h"

namespace gse { 
namespace dataserver {

int parseStorageNode(int clusterIndex, const std::string &nodeStr, vector<StorageConfigType> &storageConfigs)
{
    Json::Value root;
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(nodeStr, root, false))
    {
        LOG_ERROR("Json parse err for %s", SAFE_CSTR(nodeStr.c_str()));
        return GSE_JSON_INVALID;
    }

    if (!root.isArray())
    {
        LOG_ERROR("JSON IS INVALID");
        return GSE_JSON_INVALID;
    }

    int size = root.size();
    for (Json::ArrayIndex i = 0; i < size; ++i)
    {
        Json::Value node = root[i];
        StorageConfigType config;
        config.m_clusterIndex = clusterIndex;

        if (!node.isMember("type"))
        {
            // 兼容旧版本，没有type字段的都是kafka节点
            node["type"] = KAFKA_COMMON;
        }

        if (!node["type"].isInt())
        {
            LOG_ERROR("type is not int");
            return GSE_JSON_INVALID;
        }

        if(!node["cluster_index"].isInt())
        {
            LOG_WARN("cluster index is not int, will use the input %d", clusterIndex);
            config.m_clusterIndex = clusterIndex;
        }
        else
        {
            config.m_clusterIndex = node["cluster_index"].asInt();
        }

        if (!node["host"].isString())
        {
            LOG_ERROR("host is not string");
            return GSE_JSON_INVALID;
        }

        if (!node["port"].isInt())
        {
            LOG_ERROR("port is not int");
            return GSE_JSON_INVALID;
        }

        if (!node.isMember("passwd"))
        {
            // passwd 缺省为空串
            node["passwd"] = "";
        }

        if (!node["passwd"].isString())
        {
            LOG_ERROR("passwd is not string");
            return GSE_JSON_INVALID;
        }

        std::string mastername;
        if (node.isMember("mastername"))
        {
            mastername = node.get("mastername", "").asString();
        }
        else
        {
            mastername = "mymaster";
        }

        std::string token;
        if (node.isMember("token"))
        {
            token = node.get("token", "").asString();
        }
        config.m_maxKafkaMaxQueue = node.get("queue_buffering_max_messages", DEFAULT_MAX_KAFKA_QUEUE_SIZE).asInt();
        config.m_maxKafkaMessageBytes = node.get("message_max_bytes", DEFAULT_MAX_KAFKA_MESSAGE_BYTES_SIZE).asInt();

        config.m_storageType = node["type"].asInt();
        config.m_host = node["host"].asString();
        config.m_port = node["port"].asInt();
        config.m_passwd = node["passwd"].asString();
        config.m_token = node.get("token", "").asString();
        config.m_masterName = mastername;
        

        storageConfigs.push_back(config);
    }
    return GSE_SUCCESS;
}

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
DataID *parseToDataID(const Json::Value &root)
{
    DataID *ptrDataid = new DataID();

    if (root.isMember("biz_id"))
    {
        if(!root["biz_id"].isInt())
        {
            LOG_ERROR("json invalid,biz_id is not integer %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_bizId = root.get("biz_id", 0).asInt();
    }
    else
    {
        ptrDataid->m_bizId = 0;
    }

    if (root.isMember("partition"))
    {
        if(!root["partition"].isInt())
        {
            LOG_ERROR("json invalid,partition is not integer %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_partitions = root.get("partition", 1).asInt();
    }
    else
    {
        ptrDataid->m_partitions = 1;
    }

    if (root.isMember("msg_system"))
    {
        if(!root["msg_system"].isInt())
        {
            LOG_ERROR("json invalid,msg_system is not integer %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_storeSys = root.get("msg_system", 1).asInt();
    }
    else
    {
        ptrDataid->m_storeSys = 1;
    }

    if (root.isMember("data_set"))
    {
        if(!root["data_set"].isString())
        {
            LOG_ERROR("json invalid,data_set is not string %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_dataSet = root.get("data_set", "").asString();
    }
    else
    {
        ptrDataid->m_dataSet = "";
    }

    if (root.isMember("cluster_index"))
    {
        if(!root["cluster_index"].isInt())
        {
            LOG_ERROR("json invalid,cluster_index is not integer %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_clusterIndex = root.get("cluster_index", 0).asInt();
    }
    else
    {
        ptrDataid->m_clusterIndex = 0;
    }

    ptrDataid->m_keyTopic = ptrDataid->m_dataSet + gse::tools::strings::ToString(ptrDataid->m_bizId);

    if (root.isMember("type"))
    {
        if(!root["type"].isInt())
        {
            LOG_ERROR("json invalid,type is not integer %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_type = root.get("type", KAFKA_COMMON).asInt();
    }
    else
    {
        ptrDataid->m_type = KAFKA_COMMON;
    }

    if (root.isMember("optype"))
    {
        if(!root["optype"].isInt())
        {
            LOG_ERROR("json invalid,optype is not integer %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_optype = root.get("optype", 0).asInt();
    }
    else
    {
        ptrDataid->m_optype = 0;
    }


    if (root.isMember("tenant"))
    {
        if(!root["tenant"].isString())
        {
            LOG_ERROR("json invalid,tenant is not string %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_tenant = root.get("tenant", "").asString();
    }
    else
    {
        ptrDataid->m_tenant = "";
    }


    if (root.isMember("namespace"))
    {
        if(!root["namespace"].isString())
        {
            LOG_ERROR("json invalid,namespace is not string %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_namespace = root.get("namespace", "").asString();
    }
    else
    {
        ptrDataid->m_namespace = "";
    }


    if (root.isMember("persistent"))
    {
        if(!root["persistent"].isString())
        {
            LOG_ERROR("json invalid,persistent is not string %s", root.toStyledString().c_str());
            goto FAILED;
        }
        ptrDataid->m_persistent = root.get("persistent", "persistent").asString();
    }
    else
    {
        ptrDataid->m_persistent = "persistent";
    }

    // transit to same as storage type
    if (ptrDataid->m_storeSys == KAFKA_COMMON && ptrDataid->m_type == KAFKA_OP)
    {
        ptrDataid->m_storeSys = KAFKA_OP;
    }

    return ptrDataid;

    FAILED:
    if (ptrDataid != NULL)
    {
        delete ptrDataid;
    }
    return NULL;
}

/*
[{
    "type":1,
    "biz_id":0,
    "cluster_index":2,
    "data_set":"name",
    "msg_system":1,
    "partition":1,
    "optype":1
}]
*/
// new format, can output to different storages
// not support op now
// linked list: dataid-->dataid-->NULL
DataID *parseToDataIdArray(const Json::Value &root)
{
    DataID *head = NULL, *node = NULL;

    if (root.size() == 0)
    {
        return NULL;
    }

    for (Json::Value::const_iterator it = root.begin(); it != root.end(); ++it)
    {
        DataID *ptrDataid = parseToDataID(*it);
        if (ptrDataid == NULL)
        {
            return NULL;
        }

        if (head == NULL)
        {
            head = ptrDataid;
            node = head;
        }
        else
        {
            node->m_next = ptrDataid;
            node = ptrDataid;
        }
    }

    return head;
}

// parse dataid node
DataID *parseToDataID(const std::string &dataInfo)
{
    // 解析dataid 生成dataid 配置对象
    Json::Value root;
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(dataInfo, root, false))
    {
        return NULL;
    }

    if (root.isArray())
    {
        return parseToDataIdArray(root);
    }
    else
    {
        return parseToDataID(root);
    }
    return NULL;
}

std::string DataID::ToString()
{
    std::string strbuff = "DataId:\n{\nDataID:%d, StoreSys:%d, BizId:%d, Partitions:%d, ClusterIndex:%d, DataSet:%s, TopicName:%s, Namespace:%s, Tenent:%s\n}\n";
    std::string str_result;
    char buff[1024] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_dataId, m_storeSys, m_bizId, m_partitions, m_clusterIndex, m_dataSet.c_str(),
             toTopicString().c_str(), m_namespace.c_str(), m_tenant.c_str());
    str_result.append(buff);
    return str_result;
}

std::string StorageConfigType::ToString()
{
    std::string strbuff = "StorageConfig:\n{\nClusterIndex:%d, StorageType:%d, Host:%s, Port:%d, MasterName:%s, Password:%s, Token:%s\n}\n";
    std::string str_result;
    char buff[1024] = {0};
    snprintf(buff, sizeof(buff), strbuff.c_str(), m_clusterIndex, m_storageType, m_host.c_str(), m_port, m_masterName.c_str(), m_passwd.c_str(), m_token.c_str());
    str_result.append(buff);

    return str_result;
}

}
}
