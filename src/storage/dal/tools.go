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

package dal

import (
    "context"
    "net/http"

    "configcenter/src/common"
    "configcenter/src/common/util"
    "configcenter/src/storage/dal/types"
)

// GetDBContext returns a new context that contains JoinOption
func GetDBContext(parent context.Context, header http.Header) context.Context {
    rid := header.Get(common.BKHTTPCCRequestID)
    user := util.GetUser(header)
    owner := util.GetOwnerID(header)
    ctx := context.WithValue(parent, common.CCContextKeyJoinOption, types.JoinOption{
        RequestID:    rid,
        TxnID:        header.Get(common.BKHTTPCCTransactionID),
        TMAddr:       header.Get(common.BKHTTPCCTxnTMServerAddr),
        SessionID:    header.Get(common.BKHTTPCCTxnSessionID),
        SessionState: header.Get(common.BKHTTPCCTxnSessionState),
        TxnNumber:    header.Get(common.BKHTTPCCTransactionNumber),
    })
    ctx = context.WithValue(ctx, common.ContextRequestIDField, rid)
    ctx = context.WithValue(ctx, common.ContextRequestUserField, user)
    ctx = context.WithValue(ctx, common.ContextRequestOwnerField, owner)
    return ctx
}
