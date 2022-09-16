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

#ifndef _GSE_COMMON_CONFTOR_H_
#define _GSE_COMMON_CONFTOR_H_

#include <string>
#include <vector>
namespace gse {
namespace data {

typedef enum ConfItemValueTypeEnum_
{
    CONFITEMVALUE_TYPE_VALUE = 0,
    CONFITEMVALUE_TYPE_KEYS = 1,
    CONFITEMVALUE_TYPE_CREATE = 2,
    CONFITEMVALUE_TYPE_DELETE = 3
} ConfItemValueTypeEnum;

typedef struct WatchConfItem_
{
    std::string m_Key;
    std::vector<std::string> m_Values;
    ConfItemValueTypeEnum m_valueType;
    int m_confItemFlag;

    WatchConfItem_()
    {
        m_valueType = CONFITEMVALUE_TYPE_VALUE;
        m_confItemFlag = 0;
    }
    ~WatchConfItem_() {}

    WatchConfItem_& operator=(const WatchConfItem_& src)
    {
        this->m_Key = src.m_Key;
        this->m_Values.clear();
        this->m_Values.assign(src.m_Values.begin(), src.m_Values.end());
        this->m_valueType = src.m_valueType;
        this->m_confItemFlag = src.m_confItemFlag;
    }

    WatchConfItem_(const WatchConfItem_& src)
    {
        *this = src;
    }
} WatchConfItem;

typedef void (*FnWatchConf)(WatchConfItem& confItem, void* lpWatcher);

typedef struct WatcherInfo_
{
    FnWatchConf m_pFnWatchConf;
    int m_watchConfItemFlag;
    void* m_lpWatcher;

    WatcherInfo_()
    {
        m_lpWatcher = NULL;
        m_pFnWatchConf = NULL;
        m_watchConfItemFlag = 0;
    }
    ~WatcherInfo_() {}

    WatcherInfo_& operator=(const WatcherInfo_& src)
    {
        this->m_pFnWatchConf = src.m_pFnWatchConf;
        this->m_watchConfItemFlag = src.m_watchConfItemFlag;
        this->m_lpWatcher = src.m_lpWatcher;
        return *this;
    }
} WatcherInfo;

typedef void (*FnZkGetChildCallBack)(std::string& path, int rc, std::vector<std::string>& values, const void* data);
typedef void (*FnZkGetValueCallBack)(std::string& path, int rc, const char* value, int32_t value_len, const struct Stat* stat, const void* data);
//    typedef void(*ZK_DATA_FUN)(int32_t rc, const char *value, int32_t value_len, const struct Stat *stat, const void *data);
class Conftor
{
public:
    Conftor();
    virtual ~Conftor();

    virtual int Start() = 0;
    virtual int Stop() = 0;
    virtual int CreateConfItemWithParents(const std::string& key, std::string& value, bool isEphemeral = false) = 0;
    virtual int CreateConfItem(const std::string& key, std::string& value) = 0;
    virtual int GetConfItem(const std::string& key, std::string& value, FnWatchConf pFnWatchConf, void* lpWather, int confItemFlag) = 0;
    virtual int GetConfItemAsync(const std::string& key, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag, FnZkGetValueCallBack pFnGetValueCallBack, void* ptr_callback) = 0;
    virtual int GetChildConfItemAsync(const std::string& key, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag, FnZkGetChildCallBack pFnGetChildResultConf, void* ptr_callback) = 0;
    virtual int GetChildConfItem(const std::string& key, std::vector<std::string>& values, FnWatchConf pFnWatchConf, void* lpWather, int confItemFlag) = 0;
    virtual int SetConfItem(const std::string& key, const std::string& value) = 0;
    virtual int ExistConfItem(const std::string& key, FnWatchConf pFnWatchConf, void* lpWatcher, int confItemFlag) = 0;
    virtual int DeleteConfItem(const std::string& key) = 0;

protected:
private:
};

} // namespace data
} // namespace gse
#endif //_GSE_COMMON_CONFTOR_H_
