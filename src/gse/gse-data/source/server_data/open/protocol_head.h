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

#ifndef _PROTOCOL_PROTOCOL_HEAD_H_
#define _PROTOCOL_PROTOCOL_HEAD_H_

namespace gse {
namespace dataserver {


#pragma pack(1)

#define MGAIC_NUM 0xCAFECAFE
#define MSG_MAX_LEN (64 << 20)

#define OPS_ERROR_SKIP                  254
#define OPS_ERROR_TRANSFER              255


typedef struct DataMsgHead_
{
    unsigned int m_msgtype;
    unsigned int m_channelid;
    unsigned int m_magic;
    unsigned int m_msglen;
    char m_data[0];
}DataMsgHead;


typedef struct DataMsgHeadV1_
{
    unsigned int m_msgtype;
    unsigned int m_channelid;
    unsigned int m_msglen;
    char m_data[0];
}DataMsgHeadV1;

typedef struct TagElement_
{
    unsigned int m_tag;
    unsigned int m_len;
    char m_value[0];
}TagElement;

typedef struct TagLV_
{
    TagLV_()
    {
        m_len = 0;
    }
    unsigned int m_len;
    char m_value[0];
}TagLV;

enum tagMsgIE
{
    enExtentionId,
    enExterntionInfo,
    enMsgContent
};


#pragma pack()


}
}

#endif
