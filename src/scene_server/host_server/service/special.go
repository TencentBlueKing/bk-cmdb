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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

/*
文件描述： 该文件文件中的接口是专用接口，对应解决某个问题API。 不对UI,第三方应用开放。

*/

// BKSystemInstall 蓝鲸组件机器安装agent，主机写入cmdb
// 描述: 1.  只能操作蓝鲸业务 2. 不能将主机转移到空闲机和故障机等内置模块
// 3. 不会删除主机已经存在的主机模块， 只会新加主机与模块。 4. 不存在的主机会新加， 规则通过内网IP和 cloud id 判断主机是否存在
// 4. 进程不存在不报错
func (s *Service) BKSystemInstall(ctx *rest.Contexts) {
	input := new(metadata.BkSystemInstallRequest)
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if input.SetName == "" {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKSetNameField))
		return
	}
	if input.ModuleName == "" {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKModuleNameField))
		return
	}
	if input.InnerIP == "" {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKHostInnerIPField))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logic.NewSpecial(ctx.Kit).BkSystemInstall(ctx.Kit.Ctx, common.BKAppName, input)
		if err != nil {
			blog.Errorf("BkSystemInstall handle err: %v, rid:%s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity("")
}

func (s *Service) FindSystemUserConfigBKSwitch(ctx *rest.Contexts) {

	// 没有权限校验
	data, err := s.Logic.CoreAPI.CoreService().System().GetUserConfig(ctx.Kit.Ctx, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("FindSystemUserConfig handle err: %v, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	canModify := false
	if data != nil {
		if data.BluekingModify.Flag == true && data.BluekingModify.ExpireAt > time.Now().Unix() {
			canModify = true
		}
	}

	ctx.RespEntity(canModify)
}