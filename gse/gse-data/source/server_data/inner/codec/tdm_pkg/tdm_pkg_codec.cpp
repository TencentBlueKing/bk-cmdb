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

#include "tdm_pkg_codec.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "ops/op_metric.h"
#include "protocol.h"
#include "tools/macros.h"

namespace gse {
namespace data {

TdmPkgCodec::TdmPkgCodec()
{
}

TdmPkgCodec::~TdmPkgCodec()
{
}

int TdmPkgCodec::DecodeMsg(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataCell)
    {
        LOG_ERROR("fail to decode tdm package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }

    // decode
    GseDataHead *pDataHead = (GseDataHead *)(pDataCell->GetDataBuf());
    // decode head
    uint16_t extType = ntohs(pDataHead->m_pubHead.m_extType);
    uint16_t extVersion = ntohs(pDataHead->m_pubHead.m_extVersion);
    uint32_t magicNum = ntohl(pDataHead->m_pubHead.m_magicNum);
    uint32_t bodyLength = ntohl(pDataHead->m_pubHead.m_bodyLength);
    LOG_DEBUG("decode data head, extType[%d], extVersion[%d], magicNum[0x%x], bodyLength[%d]", extType, extVersion, magicNum, bodyLength);

    uint32_t channelID = ntohl(pDataHead->m_dataHeadExt.m_channelID);
    uint32_t timestamp = ntohl(pDataHead->m_dataHeadExt.m_timestamp);
    uint32_t dynamicDataBodyLength = ntohl(pDataHead->m_dataHeadExt.m_dynamicExtBodyLength);
    pDataCell->SetChannelID(channelID);
    pDataCell->SetCreationTimestamp(timestamp);
    LOG_DEBUG("decode data dynamic head, channelID[%d], timestamp[%d], dynamicDataBodyLength[%d]", channelID, timestamp, dynamicDataBodyLength);

    OPMetric::AddRecvMsgBytesCounter(kCodecTDM, kCodecTDM, bodyLength);
    int delaySec = 0;
    TIME_DEALY_SECOND(timestamp, delaySec);
    OPMetric::AddMsgDelayMsCounter(kCodecTDM, kRecvMessageDirection, kCodecTDM, SECOND_TO_MS(delaySec));

    do
    {

        if (dynamicDataBodyLength == 0)
        {
            LOG_DEBUG("break not set the dynamic for the chanenl id %d", channelID);
            break;
        }

        uint32_t lvHeadLen = sizeof(GseDataLVHead);
        // decode dynamic data
        int dynamicPos = sizeof(GseDataHead);
        //// source_id
        GseDataLVHead *pSourceIdHead = (GseDataLVHead *)(pDataCell->GetDataBuf() + dynamicPos);
        uint16_t sourceIdLength = ntohs(pSourceIdHead->m_valueLength);
        char *pSourceId = new (std::nothrow) char[sourceIdLength];
        if (NULL == pSourceId)
        {
            LOG_WARN("decode the source id of dynamic data failed. because fail to malloc memory for source id");
            iRet = GSE_ERROR;
            OPMetric::AddRecvMsgCounter(kCodecTDM, kCodecTDM, MESSAGE_CODEC_TDM_PROTOCOL_SOURCE_ID_ERROR, 1);
            return iRet;
        }

        dynamicPos += sizeof(GseDataLVHead);
        ::memcpy(pSourceId, pDataCell->GetDataBuf() + dynamicPos, sourceIdLength);
        std::string sourceId(pSourceId, sourceIdLength);
        pDataCell->PushExtension(sourceId);
        LOG_DEBUG("sourceid length is %u source id is %s", sourceIdLength, SAFE_CSTR(sourceId.c_str()));
        delete pSourceId;

        if (sourceIdLength + lvHeadLen == dynamicDataBodyLength)
        {
            LOG_DEBUG("break only sourceid, no others");
            break;
        }

        //// event_name
        dynamicPos += sourceIdLength;
        GseDataLVHead *pEventNameHead = (GseDataLVHead *)(pDataCell->GetDataBuf() + dynamicPos);
        uint16_t eventNameLength = ntohs(pEventNameHead->m_valueLength);
        char *pEventName = new (std::nothrow) char[eventNameLength];
        if (NULL == pEventName)
        {
            LOG_WARN("decode the event name of dynamic data failed. because fail to malloc memory for event name");
            OPMetric::AddRecvMsgCounter(kCodecTDM, kCodecTDM, MESSAGE_CODEC_TDM_PROTOCOL_EVENT_NAME_ERROR, 1);
            iRet = GSE_ERROR;
            return iRet;
        }
        dynamicPos += sizeof(GseDataLVHead);
        ::memcpy(pEventName, pDataCell->GetDataBuf() + dynamicPos, eventNameLength);
        std::string eventName(pEventName, eventNameLength);
        pDataCell->PushExtension(eventName);
        delete pEventName;

        if ((sourceIdLength + eventNameLength + lvHeadLen * 2) == dynamicDataBodyLength)
        {
            // TODO
            break;
        }
        // c_bussise_id
        dynamicPos += eventNameLength;
        GseDataLVHead *pBussiseIDHead = (GseDataLVHead *)(pDataCell->GetDataBuf() + dynamicPos);
        uint16_t bussiesIDLength = ntohs(pBussiseIDHead->m_valueLength);
        char *pBussiseID = new (std::nothrow) char[bussiesIDLength];
        if (NULL == pBussiseID)
        {
            LOG_ERROR("decode the bussise id of dynamic data failed. because fail to malloc memory for bussise id");
            OPMetric::AddRecvMsgCounter(kCodecTDM, kCodecTDM, MESSAGE_CODEC_TDM_PROTOCOL_BUSINESS_ID_ERROR, 1);
            iRet = GSE_ERROR;
            return iRet;
        }

        dynamicPos += sizeof(GseDataLVHead);
        ::memcpy(pBussiseID, pDataCell->GetDataBuf() + dynamicPos, bussiesIDLength);
        std::string bussiseID(pBussiseID, bussiesIDLength);
        pDataCell->PushExtension(bussiseID);
        delete pBussiseID;

        // if ((sourceIdLength + eventNameLength + bussiesIDLength + lvHeadLen * 3) == dynamicDataBodyLength)
        // {
        //     break;
        // }

        LOG_DEBUG("dynamic data. body len[%d] dynamic data len[%d] source id[%s], event name[%s]", bodyLength, dynamicDataBodyLength, SAFE_CSTR(sourceId.c_str()), SAFE_CSTR(eventName.c_str()));
    } while (false);

    // decode real body
    char *pDataBody = pDataCell->GetDataBuf() + sizeof(GseDataHead) + dynamicDataBodyLength;
    int dataBodyLength = bodyLength - sizeof(GseDataDynamicExt) - dynamicDataBodyLength;
    if (dataBodyLength <= 0 || dataBodyLength >= (1024 * 1024 * 64))
    {
        std::string client_ip = pDataCell->GetSourceIp();
        uint16_t client_port = pDataCell->GetSourcePort();
        LOG_ERROR("decode the data body failed, because its length[%d] is less then zero or over size 64M. package body length[%d], dynamic data body length[%d], data client address is:(%s:%d)", dataBodyLength, bodyLength, dynamicDataBodyLength, SAFE_CSTR(client_ip.c_str()), client_port);
        OPMetric::AddRecvMsgCounter(kCodecTDM, kCodecTDM, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        iRet = GSE_ERROR;
        return iRet;
    }

    pDataCell->CopyData(pDataBody, dataBodyLength);
    OPMetric::AddRecvMsgCounter(kCodecTDM, kCodecTDM, MESSAGE_CODEC_SUCCESS, 1);
    OPMetric::AddRecvMsgBytesCounter(kCodecTDM, kCodecTDM, dataBodyLength);

    return iRet;
}

int TdmPkgCodec::EncodeMsg(DataCell *pDataCell)
{
    return GSE_SUCCESS;
}
} // namespace data
} // namespace gse
