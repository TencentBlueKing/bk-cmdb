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

#include "dynamicalProtocol.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/endian.h"
#include "tools/macros.h"
#include <arpa/inet.h>
#include <string.h>
namespace gse {
namespace data {

DynamicalProtocol::DynamicalProtocol()
{
    m_ptrData = NULL;
    m_dataLen = 0;

    m_ptrDynamicalProtocol = NULL;
    m_ptrDynamicalProtocolHeadExt = NULL;
    m_ptrGseReportMsg = NULL;

    m_isAllocMemory = false;
}

DynamicalProtocol::~DynamicalProtocol()
{
    if (m_isAllocMemory && NULL != m_ptrData)
    {
        free(m_ptrData);
        m_ptrData = NULL;
        m_dataLen = 0;
    }
}

int DynamicalProtocol::setDynamicalProtocolHead(const DynamicalProtocolHead& dynamicalProtocolHead, const char* ptrData, int32_t dataLen)
{

    if (m_isAllocMemory && NULL != m_ptrData)
    {
        free(m_ptrData);
        m_ptrData = NULL;
        m_dataLen = 0;
    }

    int32_t metaMaxLen = dynamicalProtocolHead.m_extHead.m_extMetaMaxLen;

    if (metaMaxLen <= 0)
    {
        m_dataLen = DYNAMICAL_PROTOCOL_HEAD + DYNAMICAL_NEW_ALLOC_SIZE + dataLen;
        metaMaxLen = DYNAMICAL_NEW_ALLOC_SIZE;
    }
    else
    {
        m_dataLen = DYNAMICAL_PROTOCOL_HEAD + metaMaxLen + dataLen;
    }

    m_ptrData = (char*)malloc(m_dataLen);
    memset(m_ptrData, 0, m_dataLen);
    m_isAllocMemory = true;

    m_ptrGseReportMsg = (GseReportMsg*)m_ptrData;
    m_ptrDynamicalProtocol = (DynamicalProtocolHead*)m_ptrData;
    m_ptrDynamicalProtocolHeadExt = (DynamicalProtocolHeadExt*)(m_ptrData + GSE_REPORT_MSG_HEAD);

    m_ptrDynamicalProtocol->m_msgHead.msgType = htonl(dynamicalProtocolHead.m_msgHead.msgType);

    m_ptrDynamicalProtocol->m_msgHead.lMsg.dataid = htonl(dynamicalProtocolHead.m_msgHead.lMsg.dataid);
    m_ptrDynamicalProtocol->m_msgHead.lMsg.utctime = htonl(dynamicalProtocolHead.m_msgHead.lMsg.utctime);
    m_ptrDynamicalProtocol->m_msgHead.lMsg.bodylen = htonl(dataLen + DYNAMICAL_PROTOCOL_HEAD_EXT + metaMaxLen);
    m_ptrDynamicalProtocol->m_msgHead.lMsg.resv[0] = htonl(dynamicalProtocolHead.m_msgHead.lMsg.resv[0]);
    m_ptrDynamicalProtocol->m_msgHead.lMsg.resv[1] = htonl(dynamicalProtocolHead.m_msgHead.lMsg.resv[1]);

    m_ptrDynamicalProtocolHeadExt->m_extIndex = gse::tools::endian::Ntohl64(dynamicalProtocolHead.m_extHead.m_extIndex);
    m_ptrDynamicalProtocolHeadExt->m_extType = htonl(dynamicalProtocolHead.m_extHead.m_extType);
    m_ptrDynamicalProtocolHeadExt->m_extMetaLen = htonl(dynamicalProtocolHead.m_extHead.m_extMetaLen);
    m_ptrDynamicalProtocolHeadExt->m_extMetaMaxLen = htonl(metaMaxLen);
    m_ptrDynamicalProtocolHeadExt->m_extMetaCount = htonl(dynamicalProtocolHead.m_extHead.m_extMetaCount);

    memcpy(m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaMaxLen, ptrData, dataLen);

    return GSE_SUCCESS;
}

int DynamicalProtocol::loadData(char* ptrData, int32_t dataLen, bool isNeedCopy /*= false*/)
{
    if (m_isAllocMemory && NULL != m_ptrData)
    {
        free(m_ptrData);
        m_ptrData = NULL;
        m_dataLen = 0;
    }

    if (!isNeedCopy)
    {
        m_ptrData = ptrData;
        m_dataLen = dataLen;
    }
    else
    {
        m_ptrData = (char*)malloc(dataLen);
        memset(m_ptrData, 0, dataLen);
        m_dataLen = dataLen;
        m_isAllocMemory = true;
        memcpy(m_ptrData, ptrData, dataLen);
    }

    m_ptrGseReportMsg = (GseReportMsg*)m_ptrData;
    m_ptrDynamicalProtocol = (DynamicalProtocolHead*)m_ptrData;
    m_ptrDynamicalProtocolHeadExt = (DynamicalProtocolHeadExt*)(m_ptrData + GSE_REPORT_MSG_HEAD);

    return GSE_SUCCESS;
}

int DynamicalProtocol::binaryData(char*& ptrData, int32_t& dataLen)
{
    ptrData = m_ptrData;
    dataLen = m_dataLen;
    return dataLen;
}

int DynamicalProtocol::addMeta(const string& key, const string& value, bool isNeedDynamical /*= true*/)
{
    return addMeta(key.c_str(), key.size(), value.c_str(), value.size(), isNeedDynamical);
}

int DynamicalProtocol::addMeta(const char* key, const char* value, bool isNeedDynamical /*= true*/)
{
    int keylen = strlen(key);
    int valuelen = strlen(value);
    return addMeta(key, keylen, value, valuelen, isNeedDynamical);
}

int DynamicalProtocol::addMeta(const char* key, int keylen, const char* value, int valuelen, bool isNeedDynamical /*= true*/)
{
    if (NULL == m_ptrData || m_dataLen <= 0)
    {
        return -1;
    }

    int metaLen = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaLen);
    int metaMaxLen = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaMaxLen);
    int metaCount = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaCount);
    int totallen = metaLen + DYNAMICAL_META_HEAD + keylen + valuelen;
    int bodylen = ntohl(m_ptrGseReportMsg->lMsg.bodylen);
    int datalen = bodylen - DYNAMICAL_PROTOCOL_HEAD_EXT - metaMaxLen;

    if (0 > metaMaxLen || 0 > metaLen)
    {
        return -2;
    }

    if (totallen >= metaMaxLen)
    {
        if (!isNeedDynamical)
        {
            return -3;
        }

        // malloc new buffer
        int newbufsize = keylen + valuelen + DYNAMICAL_NEW_ALLOC_SIZE; // NEW ALLOC 3 OBJECT Buffer
        int bufsize = m_dataLen + newbufsize;
        char* ptrBuffer = (char*)malloc(bufsize);
        memset(ptrBuffer, 0, bufsize);
        m_isAllocMemory = true;

        // reset dynamical protocol
        m_ptrDynamicalProtocolHeadExt->m_extMetaLen = htonl(metaLen + keylen + valuelen + DYNAMICAL_META_HEAD);
        m_ptrDynamicalProtocolHeadExt->m_extMetaMaxLen = htonl(metaMaxLen + newbufsize);
        m_ptrDynamicalProtocolHeadExt->m_extMetaCount = htonl(metaCount + 1);
        m_ptrGseReportMsg->lMsg.bodylen = htonl(bodylen + newbufsize);

        // copy head
        memset(ptrBuffer, 0, bufsize);
        memcpy(ptrBuffer, m_ptrData, DYNAMICAL_PROTOCOL_HEAD);
        memcpy(ptrBuffer + DYNAMICAL_PROTOCOL_HEAD, m_ptrData + DYNAMICAL_PROTOCOL_HEAD, metaLen);

        // add metadata
        DynamicalMetaHead metaHead;
        metaHead.m_extMetaKeyLen = htonl(keylen);
        metaHead.m_extMetaValueLen = htonl(valuelen);

        memcpy(ptrBuffer + DYNAMICAL_PROTOCOL_HEAD + metaLen, &metaHead, DYNAMICAL_META_HEAD);
        memcpy(ptrBuffer + DYNAMICAL_PROTOCOL_HEAD + metaLen + DYNAMICAL_META_HEAD, key, keylen);
        memcpy(ptrBuffer + DYNAMICAL_PROTOCOL_HEAD + metaLen + DYNAMICAL_META_HEAD + keylen, value, valuelen);

        // copy origin data
        memcpy(ptrBuffer + DYNAMICAL_PROTOCOL_HEAD + metaMaxLen + newbufsize, m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaMaxLen, datalen);

        // reassigne data
        m_ptrData = ptrBuffer;
        m_dataLen = bufsize;

        m_ptrGseReportMsg = (GseReportMsg*)m_ptrData;
        m_ptrDynamicalProtocol = (DynamicalProtocolHead*)m_ptrData;
        m_ptrDynamicalProtocolHeadExt = (DynamicalProtocolHeadExt*)(m_ptrData + GSE_REPORT_MSG_HEAD);

        return GSE_SUCCESS;
    }
    else
    {
        // reset dynamical protocol
        m_ptrDynamicalProtocolHeadExt->m_extMetaLen = htonl(metaLen + keylen + valuelen + DYNAMICAL_META_HEAD);
        m_ptrDynamicalProtocolHeadExt->m_extMetaCount = htonl(metaCount + 1);

        // add metadata
        DynamicalMetaHead metaHead;
        metaHead.m_extMetaKeyLen = htonl(keylen);
        metaHead.m_extMetaValueLen = htonl(valuelen);

        memcpy(m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaLen, &metaHead, DYNAMICAL_META_HEAD);
        memcpy(m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaLen + DYNAMICAL_META_HEAD, key, keylen);
        memcpy(m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaLen + DYNAMICAL_META_HEAD + keylen, value, valuelen);
    }

    return GSE_SUCCESS;
}

void DynamicalProtocol::deleteMeta(const char* key)
{
    if (NULL == m_ptrData || m_dataLen <= 0)
    {
        return;
    }

    int metaCount = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaCount);
    int metaLen = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaLen);
    int metaMaxLen = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaMaxLen);
    char* metaPos = m_ptrData + DYNAMICAL_PROTOCOL_HEAD;

    // find key
    for (int metaIndex = 0; metaIndex < metaCount; ++metaIndex)
    {
        DynamicalMetaHead* metaHead = (DynamicalMetaHead*)metaPos;

        if (!CheckMetaDataValid(metaHead))
        {
            LOG_ERROR("failed delete meta from dynamic prottocol, meta invalid");
            return;
        }

        int keylen = ntohl(metaHead->m_extMetaKeyLen);
        int valuelen = ntohl(metaHead->m_extMetaValueLen);

        METADATA_KEY metaKey(metaPos + DYNAMICAL_META_HEAD, keylen);
        if (metaKey.compare(key) == 0)
        {
            // move data forward
            char* nextPos = metaPos + DYNAMICAL_META_HEAD + keylen + valuelen;
            char* metaEndPos = m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaMaxLen;
            // len=endpos-nextpos
            int remain = metaEndPos - nextPos;

            memcpy(metaPos, nextPos, remain);
            // reset dynamical protocol
            m_ptrDynamicalProtocolHeadExt->m_extMetaLen = htonl(metaLen - (DYNAMICAL_META_HEAD + keylen + valuelen));
            m_ptrDynamicalProtocolHeadExt->m_extMetaCount = htonl(metaCount - 1);

            return;
        }
        metaPos += DYNAMICAL_META_HEAD + keylen + valuelen;
    }

    // not found
    return;
}

int DynamicalProtocol::updateMeta(const char* key, const char* value, bool isNeedDynamical /*= true*/)
{
    // delete old one
    deleteMeta(key);
    // add new one
    return addMeta(key, value, isNeedDynamical);
}

bool DynamicalProtocol::CheckMetaDataValid(DynamicalMetaHead* pMetaHead)
{
    // m_extMetaLen
    int metaKeyLen = ntohl(pMetaHead->m_extMetaKeyLen);
    int metaValueLen = ntohl(pMetaHead->m_extMetaValueLen);
    int extMetaTotalLen = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaLen);
    if (metaKeyLen <= 0)
    {
        LOG_ERROR("recv invalid dynamical protocol message, meta key len invalid:%d", metaKeyLen);
        return false;
    }

    if (metaValueLen < 0)
    {
        LOG_ERROR("recv invalid dynamical protocol message, meta value len invalid:%d", metaValueLen);
        return false;
    }

    if (extMetaTotalLen < (metaKeyLen + metaValueLen))
    {
        LOG_ERROR("recv invalid dynamical protocol message, meta too long:%d, message meta total len:%d", (metaKeyLen + metaValueLen), extMetaTotalLen);
        return false;
    }

    return true;
}

int DynamicalProtocol::getMeta(DynamicalProtocol::METADATA& metadata)
{
    if (NULL == m_ptrData || NULL == m_ptrDynamicalProtocolHeadExt)
    {
        return -1;
    }

    int metaCount = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaCount);
    char* metaPos = m_ptrData + DYNAMICAL_PROTOCOL_HEAD;

    for (int metaIndex = 0; metaIndex < metaCount; ++metaIndex)
    {
        DynamicalMetaHead* metaHead = (DynamicalMetaHead*)metaPos;
        if (!CheckMetaDataValid(metaHead))
        {
            LOG_ERROR("failed get meta from dynamic prottocol, meta invalid");
            return -1;
        }

        int keylen = ntohl(metaHead->m_extMetaKeyLen);
        int valuelen = ntohl(metaHead->m_extMetaValueLen);

        METADATA_KEY key(metaPos + DYNAMICAL_META_HEAD, keylen);
        METADATA_VALUE value(metaPos + DYNAMICAL_META_HEAD + keylen, valuelen);

        // metadata[key] = value;
        metadata.push_back(std::make_pair(key, value));
        metaPos += DYNAMICAL_META_HEAD + keylen + valuelen;
    }

    return metadata.size();
}

// get value of specific key
DynamicalProtocol::METADATA_VALUE DynamicalProtocol::getMeta(const METADATA_KEY& key)
{
    // get meta data
    DynamicalProtocol::METADATA metaData;
    getMeta(metaData);

    // find key
    for (DynamicalProtocol::METADATA::iterator iter = metaData.begin(); iter != metaData.end(); ++iter)
    {
        if (iter->first.compare(key) == 0)
        {
            // find monitor tag
            return iter->second;
        }
    }
    return "";
}

int DynamicalProtocol::getDynamicalProtocol(DynamicalProtocolHead& head)
{
    if (NULL != m_ptrDynamicalProtocol)
    {
        head.m_msgHead.msgType = ntohl(m_ptrDynamicalProtocol->m_msgHead.msgType);
        head.m_msgHead.lMsg.dataid = ntohl(m_ptrDynamicalProtocol->m_msgHead.lMsg.dataid);
        head.m_msgHead.lMsg.utctime = ntohl(m_ptrDynamicalProtocol->m_msgHead.lMsg.utctime);
        head.m_msgHead.lMsg.bodylen = ntohl(m_ptrDynamicalProtocol->m_msgHead.lMsg.bodylen);
        head.m_msgHead.lMsg.resv[0] = ntohl(m_ptrDynamicalProtocol->m_msgHead.lMsg.resv[0]);
        head.m_msgHead.lMsg.resv[1] = ntohl(m_ptrDynamicalProtocol->m_msgHead.lMsg.resv[1]);

        head.m_extHead.m_extIndex = ntohl(m_ptrDynamicalProtocol->m_extHead.m_extIndex);
        head.m_extHead.m_extType = ntohl(m_ptrDynamicalProtocol->m_extHead.m_extType);
        head.m_extHead.m_extMetaLen = ntohl(m_ptrDynamicalProtocol->m_extHead.m_extMetaLen);
        head.m_extHead.m_extMetaMaxLen = ntohl(m_ptrDynamicalProtocol->m_extHead.m_extMetaMaxLen);
        head.m_extHead.m_extMetaCount = ntohl(m_ptrDynamicalProtocol->m_extHead.m_extMetaCount);

        return GSE_SUCCESS;
    }

    return GSE_ERROR;
}

int DynamicalProtocol::getBody(char*& ptrData, uint32_t& dataLen)
{
    int metaMaxLen = ntohl(m_ptrDynamicalProtocolHeadExt->m_extMetaMaxLen);
    dataLen = ntohl(m_ptrDynamicalProtocol->m_msgHead.lMsg.bodylen) - metaMaxLen - DYNAMICAL_PROTOCOL_HEAD_EXT;
    ptrData = m_ptrData + DYNAMICAL_PROTOCOL_HEAD + metaMaxLen;

    return dataLen;
}

uint32_t DynamicalProtocol::getBodyLength()
{
    char* data;
    uint32_t bodyLen;
    getBody(data, bodyLen);
    return bodyLen;
}

} // namespace data
} // namespace gse
