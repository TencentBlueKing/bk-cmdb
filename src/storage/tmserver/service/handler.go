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

package service

import (
	"context"

	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

func (s *coreService) DBOperation(input rpc.Request) (interface{}, error) {

	var ctx core.ContextParams

	reply := types.OPReply{}
	err := input.Decode(&ctx)
	if nil != err {
		reply.Message = err.Error()
		return &reply, nil
	}
	ctx.Context = context.Background()
	ctx.ListenIP = s.listenIP

	return s.core.ExecuteCommand(ctx, input)

}

func (s *coreService) WatchTransaction(input rpc.Request, stream rpc.ServerStream) (err error) {
	ch := make(chan *types.Transaction, 100)
	s.core.Subscribe(ch)
	defer s.core.UnSubscribe(ch)
	for txn := range ch {
		if err = stream.Send(txn); err != nil {
			return err
		}
	}
	return nil
}
