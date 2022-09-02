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

#ifndef _GSE_COMMON_ZKCONFTOR_H_
#define _GSE_COMMON_ZKCONFTOR_H_

#include "datastruct/safe_map.h"
#include "discover/zkapi/zk_api.h"
#include "discover/zkapi/zk_client.h"
#include <string>

#include "bbx/gse_errno.h"

#include "conftor.h"
namespace gse {
namespace data {
/* zookeeper state constants */
#define ZK_EXPIRED_SESSION_STATE_DEF -112
#define ZK_AUTH_FAILED_STATE_DEF -113
#define ZK_CONNECTING_STATE_DEF 1
#define ZK_ASSOCIATING_STATE_DEF 2
#define ZK_CONNECTED_STATE_DEF 3
#define ZK_NOTCONNECTED_STATE_DEF 999

/* zookeeper event type constants */
#define ZK_CREATED_EVENT_DEF 1
#define ZK_DELETED_EVENT_DEF 2
#define ZK_CHANGED_EVENT_DEF 3
#define ZK_CHILD_EVENT_DEF 4
#define ZK_SESSION_EVENT_DEF -1
#define ZK_NOTWATCHING_EVENT_DEF -2

typedef struct ZkConftorParam_
{
    std::string m_ZkHost;
    std::string m_ZkAuth;
    std::string m_HostIP;
    std::string m_BasePath;
    std::string m_password;

    ZkConftorParam_() {}
    ~ZkConftorParam_() {}

    ZkConftorParam_& operator=(const ZkConftorParam_& src)
    {
        this->m_ZkHost = src.m_ZkHost;
        this->m_ZkAuth = src.m_ZkAuth;
        this->m_HostIP = src.m_HostIP;
        this->m_password = src.m_password;
        return *this;
    }

    ZkConftorParam_(const ZkConftorParam_& src)
    {
        *this = src;
    }
} ZkConftorParam;

class ZkConftor : public Conftor
{
public:
    ZkConftor(const ZkConftorParam& zkParam);
    ZkConftor(gse::discover::zkapi::ZkApi* zkClient, bool acl);
    virtual ~ZkConftor();

    int Start();
    int Stop();
    int CreateConfItemWithParents(const std::string& key, std::string& value, bool isEphemeral = false);
    int CreateConfItem(const std::string& key, std::string& value);
    int GetConfItem(const std::string& key, std::string& value, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag);
    int GetChildConfItem(const std::string& key, std::vector<std::string>& values, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag);
    int ExistConfItem(const std::string& key, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag);
    int SetConfItem(const std::string& key, const std::string& value);

    int CreateEphemeralNode(const std::string& path, const std::string& value);
    int CreateEphemeralConfItemWithParents(const std::string& key, std::string& value);

    int GetConfItemAsync(const std::string& key, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag, FnZkGetValueCallBack pFnGetValueCallBack, void* ptr_callback);
    int GetChildConfItemAsync(const std::string& key, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag, FnZkGetChildCallBack pFnGetChildCallBack, void* ptr_callback);
    int DeleteConfItem(const std::string& path);

protected:
private:
    int connectZkHost();
    void closeZkHost();
    int createNode(const std::string& path, const std::string& value);
    int getNode(const std::string& path, std::string& value, WatcherInfo* pWatcher);
    int setNode(const std::string& path, const std::string& value);
    int getChildrenNodes(const std::string& path, std::vector<std::string>& nodes, WatcherInfo* pWatcher);
    void updateWatcher(const std::string& key, WatcherInfo* pWatcher);

private:
    //  static functions

    /**
     * type parameter:
     * CREATED_EVENT_DEF 1
       DELETED_EVENT_DEF 2
       CHANGED_EVENT_DEF 3
       CHILD_EVENT_DEF 4
       SESSION_EVENT_DEF -1
       NOTWATCHING_EVENT_DEF -2
     * **/
    static void defaultWatcher(int type, int state, const char* path, void* wctx);
    static void existWatcher(int type, int state, const char* path, void* wctx);
    static void childWatcher(int type, int state, const char* path, void* wctx);
    static void valueWatcher(int type, int state, const char* path, void* wctx);
    static void ZkEventHandle(int type, int state, const char* path, void* wctx);
    static void ZkSessionEventHandle(int type, int state, const char* path, void* wctx);

    static void GetValueCallback(int32_t rc, const char* value, int32_t value_len, const struct Stat* stat, const void* data);

private:
    // gse::discover::zkapi::ZkClient* m_zkClient;
    gse::discover::zkapi::ZkApi* m_zkClient;
    bool m_external;
    bool m_acl;
    std::string m_zkauth;
    ZkConftorParam m_zkParam;
    gse::datastruct::SafeMap<const std::string, WatcherInfo*> m_mapKeyWatcher;
};

} // namespace data
} // namespace gse
#endif //_GSE_COMMON_ZKCONFTOR_H_
