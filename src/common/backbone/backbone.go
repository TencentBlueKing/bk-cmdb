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

package backbone

import (
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
)

func NewBackbone(zkAddr string, c Config) (*Engine, error) {
	disc, err := NewServcieDiscovery(zkAddr)
	if err != nil {
		return nil, fmt.Errorf("new service discover failed, err:%v", err)
	}

	if err := ListenServer(c.Server); err != nil {
		return nil, err
	}

	return New(c, disc)
}

func New(c Config, disc ServiceDiscoverInterface) (*Engine, error) {
	if err := disc.Register(c.RegisterPath, c.RegisterInfo); err != nil {
		return nil, err
	}

	return &Engine{
		CoreAPI:  c.CoreAPI,
		SvcDisc:  disc,
		Language: c.Language,
		CCErr:    c.CCErr,
	}, nil
}

type Engine struct {
	CoreAPI  apimachinery.ClientSetInterface
	SvcDisc  ServiceDiscoverInterface
	Language language.CCLanguageIf
	CCErr    errors.CCErrorIf
}
