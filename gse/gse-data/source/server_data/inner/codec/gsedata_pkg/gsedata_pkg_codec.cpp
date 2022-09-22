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

#include "gsedata_pkg_codec.h"

#include "bbx/gse_errno.h"
#include "inner/protocol_head.h"
#include "ops/op_metric.h"

namespace gse {
namespace data {

GseDataPkgCodec::GseDataPkgCodec()
{
    //
}

GseDataPkgCodec::~GseDataPkgCodec()
{
    //
}

int GseDataPkgCodec::DecodeMsg(DataCell* pDataCell)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataCell)
    {
        LOG_WARN("fail to decode gsedata package msg, because the args of datacell is NULL");

        return GSE_ERROR;
    }

    // decode
    GseDataHead* pDataHead = (GseDataHead*)(pDataCell->GetDataBuf());
    // decode head
    uint16_t extType = ntohs(pDataHead->m_pubHead.m_extType);
    uint16_t extVersion = ntohs(pDataHead->m_pubHead.m_extVersion);
    uint32_t magicNum = ntohl(pDataHead->m_pubHead.m_magicNum);
    uint32_t bodyLength = ntohl(pDataHead->m_pubHead.m_bodyLength);
    LOG_DEBUG("decode data head, extType[%d], extVersion[%d], magicNum[0x%x], bodyLength[%d]", extType, extVersion, magicNum, bodyLength);

    OPMetric::AddRecvMsgBytesCounter(kCodecV2, kCodecV2, bodyLength);

    if (bodyLength > MSG_MAX_LEN)
    {
        LOG_WARN("gsedata package message invalid, length(%d) > %d", bodyLength, MSG_MAX_LEN);
        OPMetric::AddRecvMsgCounter(kCodecV2, kCodecV2, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return GSE_ERROR;
    }

    uint32_t channelID = ntohl(pDataHead->m_dataHeadExt.m_channelID);

    if (channelID <= 0)
    {
        LOG_INFO("codec v1 set default channelid recv ip:%s, port:%d, defaultchannelid:%d",
                 pDataCell->GetSourceIp().c_str(), pDataCell->GetSourcePort(), pDataCell->GetChannelID());
    }
    else
    {
        pDataCell->SetChannelID(channelID);
    }

    uint32_t timestamp = ntohl(pDataHead->m_dataHeadExt.m_timestamp);
    uint32_t dynamicDataBodyLength = ntohl(pDataHead->m_dataHeadExt.m_dynamicExtBodyLength);
    pDataCell->SetCreationTimestamp(timestamp);

    // decode dynamic data
    // decode real body
    char* pDataBody = pDataCell->GetDataBuf() + sizeof(GseDataHead) + dynamicDataBodyLength;
    int dataBodyLength = bodyLength - sizeof(GseDataDynamicExt) - dynamicDataBodyLength;
    if (dataBodyLength <= 0)
    {
        LOG_WARN("decode the data body failed, because its length[%d] is less then zero. package body length[%d], dynamic data body length[%d]", dataBodyLength, bodyLength, dynamicDataBodyLength);
        iRet = GSE_ERROR;
        OPMetric::AddRecvMsgCounter(kCodecV2, kCodecV2, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return iRet;
    }

    pDataCell->CopyData(pDataBody, dataBodyLength);
    OPMetric::AddRecvMsgCounter(kCodecV2, kCodecV2, MESSAGE_CODEC_SUCCESS, 1);
    OPMetric::AddRecvMsgBytesCounter(kCodecV2, kCodecV2, dataBodyLength);

    return iRet;
}

int GseDataPkgCodec::EncodeMsg(DataCell* pDataCell)
{
    return GSE_SUCCESS;
}
} // namespace data
} // namespace gse
