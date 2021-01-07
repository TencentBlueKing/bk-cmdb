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

package logics

import (
	"net/http"
	"time"

	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/esbserver"
)

type Logic struct {
	*backbone.Engine
}

type Logics struct {
	*backbone.Engine
	esbServ      esbserver.EsbClientInterface
	procHostInst *ProcHostInstConfig
	ErrHandle    errors.DefaultCCErrorIf
	header       http.Header
	rid          string
	ownerID      string
	user         string
	ccErr        errors.DefaultCCErrorIf
	ccLang       language.DefaultCCLanguageIf
}

// NewLogics get logic handle
func NewLogics(b *backbone.Engine, header http.Header, esbServ esbserver.EsbClientInterface, procHostInst *ProcHostInstConfig) *Logics {
	lang := util.GetLanguage(header)
	return &Logics{
		Engine:       b,
		header:       header,
		rid:          util.GetHTTPCCRequestID(header),
		ccErr:        b.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:       b.Language.CreateDefaultCCLanguageIf(lang),
		user:         util.GetUser(header),
		ownerID:      util.GetOwnerID(header),
		esbServ:      esbServ,
		procHostInst: procHostInst,
	}
}

// ProcHostInstConfig refresh process host instance number need config
type ProcHostInstConfig struct {
	MaxEventCount                int
	MaxRefreshModuleCount        int
	GetModuleIDInterval          time.Duration
	FetchGseOPProcResultInterval time.Duration
}
