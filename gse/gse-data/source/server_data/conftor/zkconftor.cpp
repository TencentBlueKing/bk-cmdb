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

#include "zkconftor.h"
#include "log/log.h"
#include "tools/macros.h"
#include <zookeeper/zookeeper.h>
//#include "discover/zkapi/zk_client.h"
#include "discover/zkapi/zk_api.h"

namespace gse {
namespace data {

ZkConftor::ZkConftor(const ZkConftorParam &zkParam)
{
    m_zkClient = NULL;
    m_external = false;
    m_zkParam = zkParam;
    m_acl = zkParam.m_ZkAuth.empty() ? false : true;
}

ZkConftor::ZkConftor(gse::discover::zkapi::ZkApi *zkClient, bool acl)
    : m_zkClient(zkClient), m_external(true), m_acl(acl)
{
}

ZkConftor::~ZkConftor()
{
    if (m_zkClient != NULL && !m_external)
    {
        m_zkClient->ZkClose();
        delete m_zkClient;
        m_zkClient = NULL;
    }
}

int ZkConftor::Start()
{
    //使用外的zkClient
    if (m_zkClient == NULL)
    {
        int iRet = connectZkHost();
        if (iRet != GSE_SUCCESS)
        {
            return iRet;
        }
    }

    return GSE_SUCCESS;
}

int ZkConftor::Stop()
{
    closeZkHost();
    return GSE_SUCCESS;
}

int ZkConftor::CreateConfItemWithParents(const std::string &key, std::string &value, bool isEphemeral)
{
    std::vector<std::string> split_values;
    gse::tools::strings::SplitString(split_values, key, "/");
    std::size_t max_count = split_values.size();
    std::string target_path;
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        if (idx != 0)
        {
            target_path.append("/");
        }
        target_path.append(split_values.at(idx));

        if (idx == (max_count - 1))
        {
            // 已经拼接完全部路径
            if (isEphemeral)
            {
                return CreateEphemeralNode(target_path, value);
            }
            else
            {
                return CreateConfItem(target_path, value);
            }
        }

        if (m_zkClient->ZkExists(target_path, NULL, NULL, NULL))
        {
            continue;
        }

        // 创建中间不存在的节点
        CreateConfItem(target_path, split_values.at(idx));
    }
    return GSE_SUCCESS;
}

int ZkConftor::CreateConfItem(const std::string &key, std::string &value)
{
    int iRet = createNode(key, value);
    return iRet;
}

int ZkConftor::CreateEphemeralNode(const std::string &path, const std::string &value)
{
    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to create ephemeral node[%s], because the conf client is NULL", SAFE_CSTR(path.c_str()));
        return GSE_ERROR;
    }

    string strrst;
    int ret = m_zkClient->ZkCreateEphemeral(path, value, strrst, m_acl);
    if (ret != GSE_SUCCESS)
    {
        if (GSE_ZK_NODE_EXIST == ret)
        {
            LOG_INFO("fail to create ephemeral node[%s], because the node is exist", SAFE_CSTR(path.c_str()));
            ret = GSE_SUCCESS;
        }
        else
        {
            LOG_ERROR("fail to create ephemeral node[%s], ret=[%d]", SAFE_CSTR(path.c_str()), ret);
        }

        return ret;
    }

    LOG_INFO("success to create ephemeral node[%s] for value[%s]", SAFE_CSTR(path.c_str()), SAFE_CSTR(value.c_str()));

    return GSE_SUCCESS;
}

int ZkConftor::GetConfItem(const std::string &key, std::string &value, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag)
{
    LOG_INFO("get config item[%s], watch callback[0x%x], wather[0x%x], confItemFlag[%d]", key.c_str(), pFnWatchConf, lpWatcher, confItemFlag);

    WatcherInfo *pWatcher = NULL;
    if (pFnWatchConf != NULL)
    {
        pWatcher = new WatcherInfo();
        pWatcher->m_pFnWatchConf = pFnWatchConf;
        pWatcher->m_lpWatcher = lpWatcher;
        pWatcher->m_watchConfItemFlag = confItemFlag;
    }

    int iRet = getNode(key, value, pWatcher);
    return iRet;
}

int ZkConftor::GetConfItemAsync(const std::string &key, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag, FnZkGetValueCallBack pFnGetValueCallBack, void *ptr_callback)
{
    WatcherInfo *pWatcher = NULL;
    if (pFnWatchConf != NULL)
    {
        pWatcher = new WatcherInfo();
        pWatcher->m_pFnWatchConf = pFnWatchConf;
        pWatcher->m_lpWatcher = lpWatcher;
        pWatcher->m_watchConfItemFlag = confItemFlag;
    }

    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to get config value of key[%s], because didn't create the conftor, check the conftor start or not");
        return GSE_ERROR;
    }

    int iRet = GSE_SUCCESS;
    gse::discover::zkapi::ZkApi::ZK_WATCH_FUN wfun = (pFnWatchConf == NULL ? NULL : ZkConftor::valueWatcher);
    gse::discover::zkapi::ZkApi::ZK_DATA_FUN sfun = (pFnGetValueCallBack == NULL ? NULL : pFnGetValueCallBack);

    iRet = m_zkClient->ZkaGet(key, wfun, this, sfun, ptr_callback);

    if (GSE_SUCCESS != iRet)
    {
        if (iRet == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("fail to get value of node:%s, node not exist", key.c_str());
        }
        else
        {
            LOG_WARN("fail to get value of node:%s, ret:%d", key.c_str(), iRet);
        }
    }
    else
    {
        LOG_INFO("success to get value of the node:%s, ret:%d", key.c_str(), iRet);
        if (pWatcher != NULL)
        {
            updateWatcher(key, pWatcher);
        }
    }
    return iRet;
}

int ZkConftor::GetChildConfItemAsync(const std::string &key, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag, FnZkGetChildCallBack pFnGetChildCallBack, void *ptr_callback)
{
    WatcherInfo *pWatcher = NULL;
    if (pFnWatchConf != NULL)
    {
        pWatcher = new WatcherInfo();
        pWatcher->m_pFnWatchConf = pFnWatchConf;
        pWatcher->m_lpWatcher = lpWatcher;
        pWatcher->m_watchConfItemFlag = confItemFlag;
    }

    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to get config nodes of key[%s], because didn't create the conftor, check the conftor start or not");
        return GSE_ERROR;
    }
    // int32_t ZkNewClient::ZkaGetChildren(const string & path,  ZK_WATCH_FUN wfun, void * wctx, ZK_STRINGS_FUN sfun, const void * data)

    int iRet = GSE_SUCCESS;
    gse::discover::zkapi::ZkApi::ZK_WATCH_FUN wfun = (pFnWatchConf == NULL ? NULL : ZkConftor::childWatcher);
    gse::discover::zkapi::ZkApi::ZK_STRINGS_FUN sfun = (pFnGetChildCallBack == NULL ? NULL : pFnGetChildCallBack);

    iRet = m_zkClient->ZkaGetChildren(key, wfun, this, sfun, ptr_callback);
    if (GSE_SUCCESS != iRet)
    {
        if (iRet == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("fail to get nodes [%s], node not exist", key.c_str());
        }
        else
        {
            LOG_ERROR("fail to get nodes [%s], ret=[%d]", key.c_str(), iRet);
        }
    }
    else
    {
        LOG_INFO("success to get nodes of key[%s]", SAFE_CSTR(key.c_str()));
        if (pWatcher != NULL)
        {
            updateWatcher(key, pWatcher);
        }
    }

    return iRet;
}

int ZkConftor::GetChildConfItem(const std::string &key, std::vector<std::string> &values, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag)
{
    WatcherInfo *pWatcher = NULL;
    if (pFnWatchConf != NULL)
    {
        pWatcher = new WatcherInfo();
        pWatcher->m_pFnWatchConf = pFnWatchConf;
        pWatcher->m_lpWatcher = lpWatcher;
        pWatcher->m_watchConfItemFlag = confItemFlag;
    }
    int iRet = getChildrenNodes(key, values, pWatcher);
    return iRet;
}

int ZkConftor::ExistConfItem(const std::string &key, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag)
{
    WatcherInfo *pWatcher = NULL;
    if (pFnWatchConf != NULL)
    {
        pWatcher = new WatcherInfo();
        pWatcher->m_pFnWatchConf = pFnWatchConf;
        pWatcher->m_lpWatcher = lpWatcher;
        pWatcher->m_watchConfItemFlag = confItemFlag;
    }
    if (NULL == m_zkClient)
    {
        LOG_ERROR("fail to get config nodes of key(%s), maybe the conftor is not  started, please to check it", SAFE_CSTR(key.c_str()));
        return GSE_ERROR;
    }

    bool bret = false;
    if (NULL == pWatcher)
    {
        bret = m_zkClient->ZkExists(key, NULL, NULL, NULL);
    }
    else
    {
        bret = m_zkClient->ZkExists(key, ZkConftor::existWatcher, this, NULL);
    }

    if (!bret)
    {
        if (m_zkClient->ZkGetError() == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("nodes of key:%s not exist", SAFE_CSTR(key.c_str()));
            if (pWatcher != NULL)
            {
                updateWatcher(key, pWatcher);
            }
        }
        else
        {
            LOG_ERROR("call zk exist node failed, key:%s, error:%s", SAFE_CSTR(key.c_str()));
        }
    }
    else
    {
        LOG_DEBUG("success to get nodes of key:%s", SAFE_CSTR(key.c_str()));
        if (pWatcher != NULL)
        {
            updateWatcher(key, pWatcher);
        }
    }

    return GSE_SUCCESS;
}

int ZkConftor::SetConfItem(const std::string &key, const std::string &value)
{
    int iRet = setNode(key, value);
    return iRet;
}

int ZkConftor::connectZkHost()
{
    closeZkHost();

    if (NULL == m_zkClient)
    {
        m_zkClient = new gse::discover::zkapi::ZkApi();
    }

    int flags = 0;
    int timeout = 30000; // 30s
    int ret = m_zkClient->ApiSetup();
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("Zkclient start failed");
        return ret;
    }
    ret = m_zkClient->ZkInit(m_zkParam.m_ZkHost, NULL, timeout, -1, NULL, this, flags, m_zkParam.m_password);
    if (GSE_SUCCESS != ret)
    {
        LOG_ERROR("fail to connect to the zk (%s), the timeout is (%d), please to check the password", m_zkParam.m_ZkHost.c_str(), timeout);
        return ret;
    }

    LOG_INFO("success to connect configure host[%s] of zk", SAFE_CSTR(m_zkParam.m_ZkHost.c_str()));

    return GSE_SUCCESS;
}

void ZkConftor::closeZkHost()
{
    if (m_zkClient != NULL && !m_external)
    {
        m_zkClient->ApiClose();
        m_zkClient->ApiJoin();
    }
}

int ZkConftor::createNode(const std::string &path, const std::string &value)
{
    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to create node[%s], because the conf client is NULL", SAFE_CSTR(path.c_str()));
        return GSE_ERROR;
    }

    string strrst;
    int ret = m_zkClient->ZkCreateNormal(path, value, strrst, m_acl);
    if (ret != GSE_SUCCESS)
    {
        if (GSE_ZK_NODE_EXIST == ret)
        {
            LOG_INFO("fail to create node[%s], because the node is exist", SAFE_CSTR(path.c_str()));
            ret = GSE_SUCCESS;
        }
        else
        {
            LOG_ERROR("fail to create node[%s], ret=[%d]", SAFE_CSTR(path.c_str()), ret);
        }

        return ret;
    }

    LOG_DEBUG("success to create node[%s] for value[%s]", SAFE_CSTR(path.c_str()), SAFE_CSTR(value.c_str()));

    return GSE_SUCCESS;
}

void ZkConftor::defaultWatcher(int type, int state, const char *path, void *wctx)
{
    LOG_INFO("zk conftor default watcher triggered, type=[%d], state=[%d], path=[%s]", type, state, path);

    if (ZK_SESSION_EVENT_DEF == type)
    {
        if (ZK_EXPIRED_SESSION_STATE_DEF == state)
        {
        }
    }
}

void ZkConftor::ZkEventHandle(int type, int state, const char *path, void *wctx)
{
    if (wctx == NULL)
    {
        return;
    }

    std::string key(path);
    ZkConftor *pZkConftor = (ZkConftor *)wctx;
    WatcherInfo *pWatcher = NULL;
    pZkConftor->m_mapKeyWatcher.Find(key, pWatcher);
    if (pWatcher == NULL)
    {
        LOG_WARN("no callback function for this key(%s), event:%d", key.c_str(), type);
        return;
    }
    pZkConftor->m_mapKeyWatcher.EraseByKey(key);
    std::string value;
    std::vector<std::string> nodes;
    WatchConfItem conf_item;
    conf_item.m_Key = std::string(path);
    conf_item.m_confItemFlag = pWatcher->m_watchConfItemFlag;

    switch (type)
    {
    case ZK_CREATED_EVENT_DEF:
        if (GSE_SUCCESS != pZkConftor->getNode(key, value, NULL))
        {
            LOG_WARN("fail to get config value of key[%s] when the config value wather triggered", path);
            delete pWatcher;
            return;
        }
        conf_item.m_valueType = CONFITEMVALUE_TYPE_CREATE;
        conf_item.m_Values.push_back(value);
        LOG_INFO("node create, trigger callback funtion[0x%x] when config nodes of key[%s] is triggered.", pWatcher->m_pFnWatchConf, path);
        pWatcher->m_pFnWatchConf(conf_item, pWatcher->m_lpWatcher);
        break;
    case ZK_DELETED_EVENT_DEF:
        conf_item.m_Values.push_back(value);
        conf_item.m_valueType = CONFITEMVALUE_TYPE_DELETE;
        LOG_INFO("node delete, trigger callback funtion[0x%x] when config nodes of key[%s] is triggered.", pWatcher->m_pFnWatchConf, path);
        pWatcher->m_pFnWatchConf(conf_item, pWatcher->m_lpWatcher);
        break;
    case ZK_CHANGED_EVENT_DEF:
        if (GSE_SUCCESS != pZkConftor->getNode(key, value, NULL))
        {
            LOG_WARN("fail to get config value of key[%s] when the config value wather triggered", path);
            delete pWatcher;
            return;
        }
        conf_item.m_valueType = CONFITEMVALUE_TYPE_VALUE;
        conf_item.m_Values.push_back(value);
        LOG_INFO("node value change, trigger callback funtion[0x%x] when config nodes of key[%s] is triggered. ", pWatcher->m_pFnWatchConf, path);
        pWatcher->m_pFnWatchConf(conf_item, pWatcher->m_lpWatcher);
        break;
    case ZK_CHILD_EVENT_DEF:
        if (GSE_SUCCESS != pZkConftor->getChildrenNodes(key, nodes, NULL))
        {
            LOG_WARN("fail to get config nodes of key[%s] when the config nodes wather triggered", path);
            delete pWatcher;
            return;
        }

        conf_item.m_Values.assign(nodes.begin(), nodes.end());
        conf_item.m_valueType = CONFITEMVALUE_TYPE_KEYS;
        LOG_INFO("trigger callback funtion[0x%x] when config nodes of key[%s] is triggered. the config nodes number is [%d]", pWatcher->m_pFnWatchConf, path, nodes.size());
        pWatcher->m_pFnWatchConf(conf_item, pWatcher->m_lpWatcher);
        break;
    default:
        LOG_ERROR("recv unkown event:%d, path:%s", type, path);
        delete pWatcher;
        return;
    }

    delete pWatcher;
}

void ZkConftor::existWatcher(int type, int state, const char *path, void *wctx)
{
    LOG_INFO("zk conftor exist watcher triggered, type[%d], state[%d], path[%s]", type, state, path);

    ZkEventHandle(type, state, path, wctx);
}

void ZkConftor::valueWatcher(int type, int state, const char *path, void *wctx)
{
    LOG_INFO("zk conftor value watcher triggered, type[%d], state[%d], path[%s]", type, state, path);

    ZkEventHandle(type, state, path, wctx);
}

void ZkConftor::childWatcher(int type, int state, const char *path, void *wctx)
{
    LOG_INFO("zk conftor child watcher triggered, type[%d], state[%d], path[%s]", type, state, path);
    ZkEventHandle(type, state, path, wctx);
}

int ZkConftor::getNode(const std::string &path, std::string &value, WatcherInfo *pWatcher)
{
    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to get config value of key[%s], because didn't create the conftor, check the conftor start or not");
        return GSE_ERROR;
    }

    int iRet = GSE_SUCCESS;
    if (NULL == pWatcher)
    {
        iRet = m_zkClient->ZkGet(path, value, NULL, this, NULL);
    }
    else
    {
        iRet = m_zkClient->ZkGet(path, value, ZkConftor::valueWatcher, this, NULL);
    }

    if (GSE_SUCCESS != iRet)
    {
        if (iRet == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("fail to get value of node:%s, node not exist", path.c_str());
        }
        else
        {
            LOG_ERROR("fail to get value of node:%s. iRet:%d", path.c_str(), iRet);
        }
    }
    else
    {
        LOG_INFO("success to get value of the node[%s], the value is [%s]", SAFE_CSTR(path.c_str()), SAFE_CSTR(value.c_str()));
        if (pWatcher != NULL)
        {
            updateWatcher(path, pWatcher);
        }
    }

    return iRet;
}

int ZkConftor::setNode(const std::string &path, const std::string &value)
{
    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to set node[%s], because didn't create the conftor, check the conftor start or not");
        return GSE_ERROR;
    }

    int iRet = m_zkClient->ZkSet(path, value, -1, NULL);
    if (GSE_SUCCESS != iRet)
    {
        if (iRet == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("fail to set value of node[%s]. node not exist", SAFE_CSTR(path.c_str()));
        }
        else
        {
            LOG_WARN("fail to set value of node[%s]. iRet=[%d]", SAFE_CSTR(path.c_str()), iRet);
        }
    }
    else
    {
        LOG_INFO("success to set value of the node [%s], the value is [%s]", SAFE_CSTR(path.c_str()), SAFE_CSTR(value.c_str()));
    }

    return iRet;
}

int ZkConftor::DeleteConfItem(const std::string &path)
{
    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to set node[%s], because didn't create the conftor, check the conftor start or not");
        return GSE_ERROR;
    }

    int iRet = m_zkClient->ZkDelete(path, -1);
    if (GSE_SUCCESS != iRet)
    {
        if (GSE_ZK_NODE_NOTEXIST == iRet)
        {
            LOG_INFO("failed to delete the node[%s], node not exist", SAFE_CSTR(path.c_str()));
        }
        else
        {
            LOG_WARN("failed to delete the node[%s]. ret=[%d]", SAFE_CSTR(path.c_str()), iRet);
        }
    }
    else
    {
        LOG_INFO("success to delete the node [%s]", SAFE_CSTR(path.c_str()));
    }

    return iRet;
}

void ZkConftor::updateWatcher(const std::string &key, WatcherInfo *pWatcher)
{
    WatcherInfo *pTmpWatcher = NULL;
    m_mapKeyWatcher.Find(key, pTmpWatcher);
    m_mapKeyWatcher.Push(key, pWatcher);
    if ((NULL == pTmpWatcher) || (pTmpWatcher != NULL && pTmpWatcher != pWatcher))
    {
        if (pTmpWatcher != NULL)
        {
            delete pTmpWatcher;
        }

        LOG_INFO("update wather[0x%x] for key[%s]", pWatcher, key.c_str());
    }
}

int ZkConftor::getChildrenNodes(const std::string &path, std::vector<std::string> &nodes, WatcherInfo *pWatcher)
{
    if (NULL == m_zkClient)
    {
        LOG_WARN("fail to get config nodes of key[%s], because didn't create the conftor, check the conftor start or not");
        return GSE_ERROR;
    }

    int iRet = GSE_SUCCESS;
    if (NULL == pWatcher)
    {
        iRet = m_zkClient->ZkGetChildren(path, NULL, NULL, nodes, NULL);
    }
    else
    {
        iRet = m_zkClient->ZkGetChildren(path, ZkConftor::childWatcher, this, nodes, NULL);
    }

    if (GSE_SUCCESS != iRet)
    {
        if (iRet == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to get nodes [%s], node not exist", path.c_str());
        }
        else
        {
            LOG_WARN("failed to get nodes [%s], ret=[%d]", path.c_str(), iRet);
        }
    }
    else
    {
        LOG_INFO("success to get nodes of key[%s]. the number of nodes is [%d]", SAFE_CSTR(path.c_str()), nodes.size());
        if (pWatcher != NULL)
        {
            updateWatcher(path, pWatcher);
        }
    }

    return iRet;
}
} // namespace data
} // namespace gse
