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

#include "genneral_protocol_codec.h"
#include "bbx/gse_errno.h"
#include "inner/protocol_head.h"
#include "tools/endian.h"

#include "net/message.h"
#include "ops/op_metric.h"

namespace gse {
namespace data {

GeneralProtocalCodec::GeneralProtocalCodec()
{
    m_typeNameDict = {
        {DATA_GENNERAL_MSG_TYPE_COMMON_DATA, "common_report"},
        {DATA_GENNERAL_MSG_TYPE_HEART_ALARM, "heart_alarm"},
        {DATA_GENNERAL_MSG_TYPE_DATA_FLOW, "data_flow"},
        {DATA_GENNERAL_MSG_TYPE_CLUSTER_FLOW, "cluster_flow"},
        {DATA_GENNERAL_MSG_TYPE_FILE_FLOW, "file_flow"},
        {DATA_GENNERAL_MSG_TYPE_TASK_FLOW, "task_flow"},
    };
}
GeneralProtocalCodec::~GeneralProtocalCodec()
{
}

const char *GeneralProtocalCodec::FindTLV(const char *msg, int msgLen, int tag)
{
    TLVMsg *tagmsg = (TLVMsg *)msg;

    int offset = 0;
    while (offset < msgLen)
    {
        int taglen = gse::tools::endian::NetworkToHost32(tagmsg->m_tagHead.m_len);
        LOG_DEBUG("tag:%d in recv msg, tag len:%d, body len:%d, offset:%d", gse::tools::endian::NetworkToHost32(tagmsg->m_tagHead.m_tag), taglen, msgLen, offset);
        if (gse::tools::endian::NetworkToHost32(tagmsg->m_tagHead.m_tag) == tag)
        {
            LOG_DEBUG("find tag:%d, len:%d", tag, taglen);
            return (char *)tagmsg;
        }
        offset += taglen;
        tagmsg = (TLVMsg *)((char *)tagmsg + offset);
    }

    return NULL;
}

bool GeneralProtocalCodec::MsgValidCheck(char *msg)
{
    DataGenneralProtocol *ptrHead = (DataGenneralProtocol *)msg;

    if (ntohl(ptrHead->m_magic) != MSG_CHECK_NUM)
    {
        LOG_WARN("message invalid, magic error:%d", ntohl(ptrHead->m_magic));
        OPMetric::AddRecvMsgCounter(kDefaultMsgType, kGenneralCodec, MESSAGE_CODEC_MAGIC_NUMBER_ERROR, 1);
        return false;
    }

    uint32_t len = ntohl(ptrHead->m_msglen);
    if (len > MSG_MAX_LEN)
    {
        LOG_WARN("message invalid, length(%d) > %d", len, MSG_MAX_LEN);
        OPMetric::AddRecvMsgCounter(kDefaultMsgType, kGenneralCodec, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return false;
    }

    return true;
}

std::string GeneralProtocalCodec::TypeToStr(int type)
{
    auto it = m_typeNameDict.find(type);
    if (it != m_typeNameDict.end())
    {
        return it->second;
    }
    return "unkow_type_" + gse::tools::strings::ToString(type);
}

int GeneralProtocalCodec::DecodeMsg(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataCell)
    {
        LOG_ERROR("fail to decode gsedata package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }

    DataGenneralProtocol *ptrHead = (DataGenneralProtocol *)(pDataCell->GetDataBuf());
    if (!MsgValidCheck((char *)ptrHead))
    {
        LOG_WARN("fail to decode gsedata package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }
    uint32_t len = pDataCell->GetBufferLen();

    uint32_t timestamp = 0;

    bool bRequiredChannelId = true;
    int msgType = gse::tools::endian::NetworkToHost32(ptrHead->m_msgtype);

    OPMetric::AddRecvMsgBytesCounter(TypeToStr(msgType), kGenneralCodec, len);

    // 内部消息号这个号段的为内部数据上报
    if (msgType < DATA_GENNERAL_OPS_MSG_TYPE_END)
    {
        pDataCell->SetOpsMsg(true);
        pDataCell->SetOpsServiceId(ntohl(ptrHead->m_msgtype));
        bRequiredChannelId = false;
    }

    char *msgContent = ptrHead->m_data;
    int bodyLen = gse::tools::endian::NetworkToHost32(ptrHead->m_msglen);

    if (bodyLen > 0)
    {
        TLVChannelID *ptrChannleid = (TLVChannelID *)FindTLV(msgContent, bodyLen, TAG_CHANNEL_ID);
        if (ptrChannleid != NULL)
        {
            uint32_t channle_id = gse::tools::endian::NetworkToHost32(ptrChannleid->m_channel_id);
            pDataCell->SetChannelID(channle_id);
            LOG_DEBUG("filled in msg channel_id:%d", channle_id);
        }
        else if (bRequiredChannelId)
        {
            LOG_ERROR("required channelid tag in msg");
            OPMetric::AddRecvMsgCounter(TypeToStr(msgType), kGenneralCodec, MESSAGE_CODEC_CHANNELID_NOT_SET_ERROR, 1);
            return GSE_ERROR;
        }

        TLVTimeStamp *ptrTimestamp = (TLVTimeStamp *)FindTLV(msgContent, bodyLen, TAG_TIMESTAMP);
        if (ptrTimestamp != NULL)
        {
            timestamp = gse::tools::endian::NetworkToHost32(ptrTimestamp->m_timestamp);
            pDataCell->SetArrivedTimestamp(timestamp);
        }

        TLVDataContent *ptrMsgContent = (TLVDataContent *)FindTLV(msgContent, bodyLen, TAG_MSG_CONTENT);
        if (ptrMsgContent != NULL)
        {
            int len = gse::tools::endian::NetworkToHost32(ptrMsgContent->m_tagHead.m_len);

            pDataCell->CopyData(ptrMsgContent->m_data, len - TLV_HEAD_LEN);
            OPMetric::AddRecvMsgBytesCounter(TypeToStr(msgType), kGenneralCodec, len);
        }
    }
    else
    {
        LOG_DEBUG("not any tag in the message, msg len:%d", len);
    }

    int delaySec = 0;
    TIME_DEALY_SECOND(timestamp, delaySec);
    OPMetric::AddMsgDelayMsCounter(TypeToStr(msgType), kRecvMessageDirection, kGenneralCodec, SECOND_TO_MS(delaySec));

    OPMetric::AddRecvMsgCounter(TypeToStr(msgType), kGenneralCodec, MESSAGE_CODEC_SUCCESS, 1);

    return GSE_SUCCESS;
}

int GeneralProtocalCodec::EncodeMsg(DataCell *pDataCell)
{
    return GSE_SUCCESS;
}
} // namespace data
} // namespace gse
