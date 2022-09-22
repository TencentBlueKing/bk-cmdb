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

#ifndef SOURCE_COMMON_DYNAMIC_PROTOCOL_DYNAMICALPROTOCOL_H_
#define SOURCE_COMMON_DYNAMIC_PROTOCOL_DYNAMICALPROTOCOL_H_
#include "inner/protocol_head.h"
#include <list>
#include <map>
#include <string>

namespace gse {
namespace data {

using namespace std;



class DynamicalProtocol
{
public:
    typedef string METADATA_KEY;
    typedef string METADATA_VALUE;
    typedef std::list<std::pair<METADATA_KEY, METADATA_VALUE>> METADATA;

public:
    DynamicalProtocol();
    ~DynamicalProtocol();

public:
    enum
    {
        GSE_REPORT_MSG_HEAD = sizeof(GseReportMsg),
        DYNAMICAL_PROTOCOL_HEAD = sizeof(DynamicalProtocolHead),
        DYNAMICAL_PROTOCOL_HEAD_EXT = sizeof(DynamicalProtocolHeadExt),
        DYNAMICAL_META_HEAD = sizeof(DynamicalMetaHead),
        DYNAMICAL_NEW_ALLOC_SIZE = (DYNAMICAL_META_HEAD + 128) * 3
    };

public:
    int setDynamicalProtocolHead(const DynamicalProtocolHead& dynamicalProtocolHead, const char* ptrData, int32_t dataLen);
    int loadData(char* ptrData, int32_t dataLen, bool isNeedCopy = false);
    int binaryData(char*& ptrData, int32_t& dataLen);
    int addMeta(const string& key, const string& value, bool isNeedDynamical = true);
    int addMeta(const char* key, const char* value, bool isNeedDynamical = true);
    int addMeta(const char* key, int keylen, const char* value, int valuelen, bool isNeedDynamical = true);
    void deleteMeta(const char* key);
    int updateMeta(const char* key, const char* value, bool isNeedDynamical = true);
    int getMeta(METADATA& metadata);
    DynamicalProtocol::METADATA_VALUE getMeta(const METADATA_KEY& key);
    bool CheckMetaDataValid(DynamicalMetaHead* pMetaHead);

public:
    int getDynamicalProtocol(DynamicalProtocolHead& head);

    uint32_t getBodyLength();
    int getBody(char*& ptrData, uint32_t& dataLen);

private:
    char* m_ptrData;
    int32_t m_dataLen;
    bool m_isAllocMemory;

private:
    DynamicalProtocolHeadExt* m_ptrDynamicalProtocolHeadExt;
    GseReportMsg* m_ptrGseReportMsg;
    DynamicalProtocolHead* m_ptrDynamicalProtocol;
};

} // namespace data
} // namespace gse
#endif /* SOURCE_COMMON_DYNAMIC_PROTOCOL_DYNAMICALPROTOCOL_H_ */
