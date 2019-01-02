/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/types"
)

var (
	// GCommands global db operation command map definition
	GCommands = &commands{cmds: nil}
)

// Command db operation definition
type Command interface {
	Execute(ctx ContextParams, decoder rpc.Request) (*types.OPReply, error)
}

type commands struct {
	cmds map[types.OPCode]Command
}

func (c *commands) SetCommand(opCode types.OPCode, cmd Command) {
	if nil == c.cmds {
		c.cmds = make(map[types.OPCode]Command)
	}

	c.cmds[opCode] = cmd
}
