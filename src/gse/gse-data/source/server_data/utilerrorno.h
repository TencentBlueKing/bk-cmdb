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



#ifndef _DATA_UTIL_ERRNO_H_
#define  _DATA_UTIL_ERRNO_H_
namespace gse {
namespace dataserver {
// BK ERROR CODE:   Product_id + self-defined error_code. e.g: 1000115
// BK GSE PRODUCT ID 10
// -1 ==> 1000101
// 0 ==> 0
// other ==> 1000000 + other
#define BK_ERROR_CODE(errorcode)  ( -1 == (errorcode) ? 1000101 : (0 < (errorcode) && 1000000>(errorcode) ? (1000000+(errorcode)) : (errorcode)) )

/////////////////////////////////////////////////////////
#define GSE_DATA_ERROR  1200



///////////////////////////////////////////////
//              GSE DATA       ERROR             //
///////////////////////////////////////////////
#define GSE_NOSUCHDATAID                      (GSE_STATE_ERROR + 1)
#define GSE_NOSUCHSERVERID                    (GSE_STATE_ERROR + 2)
#define GSE_SERVERNOTSETUP                    (GSE_STATE_ERROR + 3)
#define GSE_DS_FD_INVALID                     ((GSE_STATE_ERROR + 4) * -1)
#define GSE_DS_IOCTL_ERROR                    ((GSE_STATE_ERROR + 5) * -1)
#define GSE_DS_CONFIG_FILE_LOST               ((GSE_STATE_ERROR + 6) * -1)
#define GSE_DS_CONFIG_INVALID                 ((GSE_STATE_ERROR + 7) * -1)
#define GSE_DS_PLUGIN_CONFIG_INVALID          ((GSE_STATE_ERROR + 8) * -1)
#define GSE_DS_SET_INOTIFY_FAILED             ((GSE_STATE_ERROR + 9) * -1)
#define GSE_DS_INVALID_DATAID (GSE_STATE_ERROR + 10) // 410
#define GSE_DS_INVALID_KAFKA (GSE_STATE_ERROR + 11) // 411
#define GSE_DS_INVALID_TOPIC (GSE_STATE_ERROR + 12) // 412


#define GSE_DATA_TGLOG_DECODE_INVALIDARGS  ((GSE_DATA_ERROR + 1) * -1)
#define GSE_DATA_TGLOG_DECODE_NOENOUGHBUF  ((GSE_DATA_ERROR + 2) * -1)
#define GSE_DATA_TGLOG_DECODE_FAILPARSE    ((GSE_DATA_ERROR + 3) * -1)
#define GSE_DATA_TGLOG_DECODE_FAIL         ((GSE_DATA_ERROR + 4) * -1)
#define GSE_DATA_TGLOG_ENCODE_INVALIDARGS  ((GSE_DATA_ERROR + 5) * -1)
#define GSE_DATA_TGLOG_ENCODE_NOENOUGHBUF  ((GSE_DATA_ERROR + 6) * -1)
#define GSE_DATA_TGLOG_ENCODE_FAILSERIALIZE ((GSE_DATA_ERROR + 7) * -1)
#define GSE_DATA_TGLOG_ENCODE_FAIL          ((GSE_DATA_ERROR + 8) * -1)
#define GSE_DATA_TGLOG_CONTENT_EMPTY        ((GSE_DATA_ERROR + 9) * -1)
#define GSE_DATA_TGLOG_UNKONWTYPE        ((GSE_DATA_ERROR + 10) * -1)

///////////////////////////////////////////////////////////////////////////////
#define GSE_AGENT_RSP 0xdeadbeef

}
}
#endif
//end
