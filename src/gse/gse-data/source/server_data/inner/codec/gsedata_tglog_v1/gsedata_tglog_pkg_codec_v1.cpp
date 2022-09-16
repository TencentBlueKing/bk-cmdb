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

#include "gsedata_tglog_pkg_codec_v1.h"
#include "bbx/gse_errno.h"
#include "inner/protocol_head.h"
#include "ops/op_metric.h"
namespace gse {
namespace data {

GseDataTglogPkgCodecV1::GseDataTglogPkgCodecV1()
{
}

GseDataTglogPkgCodecV1::~GseDataTglogPkgCodecV1()
{
}

bool GseDataTglogPkgCodecV1::MsgValidCheck(char *msg)
{
    GseDataTglogMsgHead *phead = (GseDataTglogMsgHead *)(msg);
    uint32_t len = ntohl(phead->m_length);
    if (len > MSG_MAX_LEN)
    {
        LOG_WARN("tglog message invalid, length(%d) > %d", len, MSG_MAX_LEN);
        OPMetric::AddRecvMsgCounter(kCodecTGLOG, kCodecTGLOG, MESSAGE_CODEC_BODYLEN_ERROR, 1);
        return false;
    }

    if (ntohl(phead->m_magic) != MSG_CHECK_NUM)
    {
        LOG_WARN("tglog message invalid, magic error:%d", ntohl(phead->m_magic));
        OPMetric::AddRecvMsgCounter(kCodecTGLOG, kCodecTGLOG, MESSAGE_CODEC_MAGIC_NUMBER_ERROR, 1);
        return false;
    }
    return true;
}
int GseDataTglogPkgCodecV1::DecodeMsg(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataCell)
    {
        LOG_ERROR("fail to decode gse data tglog package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }

    if (!MsgValidCheck(pDataCell->GetDataBuf()))
    {
        LOG_ERROR("fail to decode gse data tglog package msg, because msg invalid");
        return GSE_ERROR;
    }

    GseDataTglogMsgHead *ptr_head = (GseDataTglogMsgHead *)(pDataCell->GetDataBuf());
    uint32_t data_length = ntohl(ptr_head->m_length);

    char *ptr_data_body = pDataCell->GetDataBuf() + sizeof(GseDataTglogMsgHead);
    uint32_t data_id = ntohl(ptr_head->m_destid);
    uint32_t magic_num = ntohl(ptr_head->m_magic);
    LOG_DEBUG("receive data id [%u], the data len is [%d], the magicnum is [0x%x]", data_id, data_length, ptr_head->m_magic);
    pDataCell->SetChannelID(data_id);
    pDataCell->CopyData(ptr_data_body, data_length);
    OPMetric::AddRecvMsgCounter(kCodecTGLOG, kCodecTGLOG, MESSAGE_CODEC_SUCCESS, 1);
    OPMetric::AddRecvMsgBytesCounter(kCodecTGLOG, kCodecTGLOG, data_length);

    return GSE_SUCCESS;
}

int GseDataTglogPkgCodecV1::EncodeMsg(DataCell *pDataCell)
{
    return GSE_SUCCESS;
}
} // namespace data
} // namespace gse
