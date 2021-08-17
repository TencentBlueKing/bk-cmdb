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
#include "protocol.h"
namespace gse { 
namespace dataserver {

GseDataPkgCodec::GseDataPkgCodec()
{
    //
}

GseDataPkgCodec::~GseDataPkgCodec()
{
    //
}

int GseDataPkgCodec::DecodeMsg(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataCell)
    {
        LOG_WARN("fail to decode gsedata package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }

    // decode
    DataMsgHead* pDataHead = (DataMsgHead*)(pDataCell->GetDataBuf());
    // decode head
    uint16_t msgType = ntohs(pDataHead->m_msgtype);
    uint32_t magicNum = ntohl(pDataHead->m_magic);
    uint32_t bodyLength = ntohl(pDataHead->m_msglen);
    LOG_DEBUG("decode data head, extType[%d], extVersion[%d], magicNum[0x%x], bodyLength[%d]", msgType, magicNum, bodyLength);

    if (bodyLength > MSG_MAX_LEN)
    {
        LOG_WARN("gsedata package message invalid, length(%d) > %d", bodyLength, MSG_MAX_LEN);
        return GSE_ERROR;
    }

    uint32_t channelID = ntohl(pDataHead->m_channelid);

    if (channelID <=0)
    {
        LOG_INFO("codec v1 set default channelid recv ip:%s, port:%d, defaultchannelid:%d",
            pDataCell->GetSourceIp().c_str(), pDataCell->GetSourcePort(), pDataCell->GetChannelID());
    }
    else
    {
        pDataCell->SetChannelID(channelID);
    }

    int msg_offset = 0;
    char* pDataBody = pDataCell->GetDataBuf() + sizeof(DataMsgHead);

    do
    {
        TagElement *ptr_tag = nullptr;
        ptr_tag = (TagElement *)pDataBody;
        if (ptr_tag->m_tag == enMsgContent)
        {
            pDataCell->CopyData(ptr_tag->m_value, ntohl(ptr_tag->m_len));
            break;
        }
        msg_offset += ptr_tag->m_len + sizeof(TagElement);

    }while(msg_offset < bodyLength);



    return iRet;
}

int GseDataPkgCodec::EncodeMsg(DataCell *pDataCell)
{    
    return GSE_SUCCESS;
}

}
}
