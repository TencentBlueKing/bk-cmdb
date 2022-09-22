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

#include "gsedata_pkg_codec_v1.h"
#include "inner/dynamic_protocol/dynamicalProtocol.h"
#include "inner/protocol_head.h"
#include "tools/time.h"

#include "bbx/gse_errno.h"
#include "conf/configurator.h"
#include "ops/op_metric.h"
#include "tools/macros.h"

namespace gse {
namespace data {

GseDataPkgCodecV1::GseDataPkgCodecV1()
{
    m_typeNameDict = {
        {REPORT_EXT, "extend_protocol"},
        {MSG_DATA_REPORT, "v1_protocol"},
        {REPORT_DYNAMICAL_PROTOCOL_TYPE, "dynamic_protocol"},
        {MSG_TYPE_HEART, "heart"},
        {MSG_DATA_REPORT_OPS, "ops_report"}};

    m_decodeTypeFunc.insert(std::make_pair(REPORT_EXT, &GseDataPkgCodecV1::decodeOriginProtocolExt));
    m_decodeTypeFunc.insert(std::make_pair(MSG_DATA_REPORT, &GseDataPkgCodecV1::decodeOriginProtocol));
    m_decodeTypeFunc.insert(std::make_pair(REPORT_DYNAMICAL_PROTOCOL_TYPE, &GseDataPkgCodecV1::decodeDynamicProtocol));
    m_decodeTypeFunc.insert(std::make_pair(MSG_TYPE_HEART, &GseDataPkgCodecV1::decodeHeart));
    m_decodeTypeFunc.insert(std::make_pair(MSG_DATA_REPORT_OPS, &GseDataPkgCodecV1::decodeOPS));

    // 构造函数,初始话内部对象及属性
    m_ptrCurl = curl_easy_init();
    if (NULL == m_ptrCurl)
    {
        LOG_DEBUG("CURL INIT FAILE");
    }
}

GseDataPkgCodecV1::~GseDataPkgCodecV1()
{
}

// URL encode string
int GseDataPkgCodecV1::urlEncode(const std::string &src, std::string &result)
{
    if (NULL == m_ptrCurl)
    {
        m_ptrCurl = curl_easy_init();
    }

    if (NULL != m_ptrCurl)
    {
        char *output = curl_easy_escape(m_ptrCurl, src.c_str(), src.length());
        if (output)
        {
            result = output;
            curl_free(output);
            return GSE_SUCCESS;
        }
    }
    else
    {
        LOG_WARN("CURL INIT FAILE");
    }
    return GSE_ERROR;
}

std::string GseDataPkgCodecV1::TypeToName(int type)
{
    auto it = m_typeNameDict.find(type);
    if (it != m_typeNameDict.end())
    {
        return it->second;
    }

    return "unkown_type_" + gse::tools::strings::ToString(type);
}

int GseDataPkgCodecV1::decodeOriginProtocolExt(DataCell *pDataCell)
{
    GseReportMsg *ptr_head = (GseReportMsg *)(pDataCell->GetDataBuf());
    GseReportExt *ptr_head_ext = (GseReportExt *)(pDataCell->GetDataBuf() + sizeof(GseReportMsg));
    int body_length = ntohl(ptr_head_ext->bodylen);
    char *pDataBody = pDataCell->GetDataBuf() + sizeof(GseReportMsg) + sizeof(GseReportExt);
    if (body_length <= 0)
    {
        LOG_WARN("decode the data body failed, because its length[%d] is less then zero.", body_length);
        OPMetric::AddRecvMsgCounter(TypeToName(REPORT_EXT), kCodecV1, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return GSE_ERROR;
    }

    LOG_DEBUG("ds decode v1 ext data is %s", SAFE_CSTR(pDataBody));
    pDataCell->SetBizID(ntohl(ptr_head->lMsg.resv[1]));
    pDataCell->SetChannelID(ntohl(ptr_head->lMsg.dataid));
    pDataCell->SetCreationTimestamp(ntohl(ptr_head->lMsg.utctime));
    pDataCell->CopyData(pDataBody, body_length);
    OPMetric::AddRecvMsgCounter(TypeToName(REPORT_EXT), kCodecV1, MESSAGE_CODEC_SUCCESS, 1);
    OPMetric::AddRecvMsgBytesCounter(TypeToName(REPORT_EXT), kCodecV1, body_length);

    return GSE_SUCCESS;
}

int GseDataPkgCodecV1::decodeOriginProtocol(DataCell *pDataCell)
{
    GseReportMsg *ptr_head = (GseReportMsg *)(pDataCell->GetDataBuf());
    int headlen = sizeof(GseReportMsg);
    int body_length = ntohl(ptr_head->lMsg.bodylen);
    // decode real body
    char *pDataBody = pDataCell->GetDataBuf() + sizeof(GseReportMsg);
    if (body_length <= 0)
    {
        LOG_WARN("decode the data body failed, because its length[%d] is less then zero.", body_length);
        OPMetric::AddRecvMsgCounter(TypeToName(MSG_DATA_REPORT), kCodecV1, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return GSE_ERROR;
    }

    LOG_DEBUG("ds decode v1 data is %s", SAFE_CSTR(pDataBody));
    pDataCell->SetBizID(ntohl(ptr_head->lMsg.resv[1]));
    pDataCell->SetChannelID(ntohl(ptr_head->lMsg.dataid));
    pDataCell->SetCreationTimestamp(ntohl(ptr_head->lMsg.utctime));
    pDataCell->CopyData(pDataBody, body_length);
    OPMetric::AddRecvMsgCounter(TypeToName(MSG_DATA_REPORT), kCodecV1, MESSAGE_CODEC_SUCCESS, 1);
    OPMetric::AddRecvMsgBytesCounter(TypeToName(MSG_DATA_REPORT), kCodecV1, body_length);

    return GSE_SUCCESS;
}

int GseDataPkgCodecV1::decodeDynamicProtocol(DataCell *pDataCell)
{
    GseReportMsg *phead = (GseReportMsg *)(pDataCell->GetDataBuf());
    int bodyLen = ntohl(phead->lMsg.bodylen);
    if (bodyLen > TCP_MAX_MESSAGE_LEN)
    {
        LOG_ERROR("protocol error, message len invalid:%d", phead->lMsg.bodylen);
        OPMetric::AddRecvMsgCounter(TypeToName(REPORT_DYNAMICAL_PROTOCOL_TYPE), kCodecV1, MESSAGE_CODEC_BODYLEN_ERROR, 1);

        return GSE_ERROR;
    }

    // 提取通过动态协议上报的数据
    DynamicalProtocol dynamicalProtocol;
    // 加载metadata 数据
    int ret = dynamicalProtocol.loadData(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());
    if (0 != ret)
    {
        LOG_ERROR("protocol error, can not parse dynamic protocol");
        OPMetric::AddRecvMsgCounter(TypeToName(REPORT_DYNAMICAL_PROTOCOL_TYPE), kCodecV1, MESSAGE_CODEC_DYNAMIC_PROTOCOL_PARSE_ERROR, 1);
        return GSE_ERROR;
    }

    DynamicalProtocolHead dynamicalProtocolHead;
    ret = dynamicalProtocol.getDynamicalProtocol(dynamicalProtocolHead);
    if (ret != 0)
    {
        LOG_ERROR("protocol error, can not parse dynamic protocol");
        OPMetric::AddRecvMsgCounter(TypeToName(REPORT_DYNAMICAL_PROTOCOL_TYPE), kCodecV1, MESSAGE_CODEC_DYNAMIC_PROTOCOL_PARSE_ERROR, 1);
        return GSE_ERROR;
    }

    if (bodyLen < dynamicalProtocolHead.m_extHead.m_extMetaLen)
    {
        LOG_ERROR("protocol error, message meta len invalid:%d", phead->lMsg.bodylen);
        OPMetric::AddRecvMsgCounter(TypeToName(REPORT_DYNAMICAL_PROTOCOL_TYPE), kCodecV1, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return GSE_ERROR;
    }

    // 提取数据段数据
    char *pData = NULL;
    uint32_t dLen = 0;

    // get input monitor tag
    string input_tag = dynamicalProtocol.getMeta(DEFAULT_MONITOR_TAG_NAME);
    if (input_tag != "")
    {
        dynamicalProtocol.deleteMeta(DEFAULT_MONITOR_TAG_NAME);
    }

    dynamicalProtocol.getBody(pData, dLen);
    LOG_DEBUG("ds decode v1 dynamic data body is %s, monitor tag:%s", pData, input_tag.c_str());

    if (dynamicalProtocolHead.m_msgHead.lMsg.dataid <= 0)
    {
        LOG_DEBUG("codec v1 set default channelid recv ip:%s, port:%d, defaultchannelid:%d",
                  pDataCell->GetSourceIp().c_str(), pDataCell->GetSourcePort(), pDataCell->GetChannelID());
    }
    else
    {
        LOG_DEBUG("codec v1 set channelid recv ip:%s, port:%d, channelid:%d",
                  pDataCell->GetSourceIp().c_str(), pDataCell->GetSourcePort(), dynamicalProtocolHead.m_msgHead.lMsg.dataid);
        pDataCell->SetChannelID(dynamicalProtocolHead.m_msgHead.lMsg.dataid);
    }

    pDataCell->SetCreationTimestamp(dynamicalProtocolHead.m_msgHead.lMsg.utctime);
    pDataCell->SetBizID(dynamicalProtocolHead.m_msgHead.lMsg.resv[1]);
    pDataCell->SetInputTag(input_tag);
    pDataCell->CopyData(pData, dLen);

    OPMetric::AddRecvMsgCounter(TypeToName(REPORT_DYNAMICAL_PROTOCOL_TYPE), kCodecV1, MESSAGE_CODEC_SUCCESS, 1);
    OPMetric::AddRecvMsgBytesCounter(TypeToName(REPORT_DYNAMICAL_PROTOCOL_TYPE), kCodecV1, dLen);

    return GSE_SUCCESS;
}

int GseDataPkgCodecV1::decodeHeart(DataCell *pDataCell)
{
    LOG_DEBUG("heart beat");
    std::string transter = "transfer";
    pDataCell->SetErrorMsg(transter, OPS_ERROR_SKIP);
    OPMetric::AddRecvMsgCounter(TypeToName(MSG_TYPE_HEART), kCodecV1, MESSAGE_CODEC_SUCCESS, 1);
    return GSE_SUCCESS;
}
int GseDataPkgCodecV1::decodeOPS(DataCell *pDataCell)
{
    GseReportMsg *ptr_head = (GseReportMsg *)(pDataCell->GetDataBuf());
    char *pDataBody = pDataCell->GetDataBuf() + sizeof(GseReportMsg);
    int body_length = htonl(ptr_head->lMsg.bodylen);
    // 1.0 协议支持ops数据上报，dataid 是内部的消息号
    if (ntohl(ptr_head->lMsg.dataid) == COLLOCTOR_OP_TYPE_AGENT_FLOW)
    {
        std::string transter = "transfer";
        pDataCell->SetErrorMsg(transter, OPS_ERROR_SKIP);
        return GSE_SUCCESS;
    }

    std::string transter = "transfer";
    pDataCell->SetErrorMsg(transter, OPS_ERROR_TRANSFER);
    pDataCell->SetChannelID(ntohl(ptr_head->lMsg.dataid));
    pDataCell->CopyData(pDataBody, body_length);
    OPMetric::AddRecvMsgCounter(TypeToName(MSG_DATA_REPORT_OPS), kCodecV1, MESSAGE_CODEC_SUCCESS, 1);

    return GSE_SUCCESS;
}
int GseDataPkgCodecV1::DecodeMsg(DataCell *pDataCell)
{
    if (NULL == pDataCell)
    {
        LOG_ERROR("fail to decode gsedata package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }

    // 本方法内对Agent 发送来的数据进行提取并向kafka 写入
    GseReportMsg *ptr_head = (GseReportMsg *)(pDataCell->GetDataBuf());

    uint32_t agentTimestamp = ntohl(ptr_head->lMsg.utctime);
    //    uint32_t bodyLen = ntohl(ptr_head->lMsg.bodylen);
    // OPMetric::AddRecvMsgBytesCounter(msgType, kCodecV1, bodyLen);
    int msgType = ntohl(ptr_head->msgType);
    int delaySec = 0;
    TIME_DEALY_SECOND(agentTimestamp, delaySec);
    OPMetric::AddMsgDelayMsCounter(TypeToName(msgType), kRecvMessageDirection, kCodecV1, SECOND_TO_MS(delaySec));

    DecodeTypeFuncMap::iterator iter = m_decodeTypeFunc.find(msgType);
    if (iter != m_decodeTypeFunc.end())
    {
        DecodeFuncPtr ptr_callback = iter->second;
        return (this->*ptr_callback)(pDataCell);
    }

    return GSE_ERROR;
}

int GseDataPkgCodecV1::EncodeMsg(DataCell *pDataCell)
{
    return GSE_SUCCESS;
}
} // namespace data
} // namespace gse
