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

// Package rest TODO
package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

// Action TODO
type Action struct {
	Verb    string
	Path    string
	Handler func(contexts *Contexts)
}

// RestfulConfig TODO
type RestfulConfig struct {
	RootPath string
}

// Config TODO
type Config struct {
	ErrorIf  errors.CCErrorIf
	Language language.CCLanguageIf
}

// NewRestUtility TODO
func NewRestUtility(conf Config) *RestUtility {
	once.Do(func() {
		initMetric()
	})

	return &RestUtility{
		Config:  conf,
		actions: make([]Action, 0),
	}
}

// RestUtility TODO
type RestUtility struct {
	Config
	actions []Action
}

// AddHandler TODO
func (r *RestUtility) AddHandler(action Action) {
	if r.actions == nil {
		r.actions = make([]Action, 0)
	}

	switch action.Verb {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
	default:
		panic(fmt.Sprintf("add http handler failed, inavlid http verb: %s.", action.Verb))
	}

	if len(action.Path) == 0 {
		panic("add http handler, but got empty http path.")
	}

	if action.Handler == nil {
		panic("add http handler, but got nil http handler")
	}

	r.actions = append(r.actions, action)
}

// AddToRestfulWebService TODO
func (r *RestUtility) AddToRestfulWebService(ws *restful.WebService) {

	for _, action := range r.actions {
		switch action.Verb {
		case http.MethodPost:
			ws.Route(ws.POST(action.Path).To(r.wrapperAction(action)))
		case http.MethodDelete:
			ws.Route(ws.DELETE(action.Path).To(r.wrapperAction(action)))
		case http.MethodPut:
			ws.Route(ws.PUT(action.Path).To(r.wrapperAction(action)))
		case http.MethodGet:
			ws.Route(ws.GET(action.Path).To(r.wrapperAction(action)))
		default:
			panic(fmt.Sprintf("rest utility add handler to webservice, but got unsupport verb: %s .", action.Verb))
		}

	}
	return
}

func (r *RestUtility) wrapperAction(action Action) func(req *restful.Request, resp *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		restContexts := new(Contexts)
		restContexts.Request = req
		restContexts.resp = resp
		restContexts.uri = action.Path

		header := req.Request.Header
		rid := util.GetHTTPCCRequestID(header)
		user := util.GetUser(header)
		owner := util.GetOwnerID(header)
		ctx := req.Request.Context()
		ctx = context.WithValue(ctx, common.ContextRequestIDField, rid)
		ctx = context.WithValue(ctx, common.ContextRequestUserField, user)
		ctx = context.WithValue(ctx, common.ContextRequestOwnerField, owner)

		// time out after 2 minutes, in case long request does not terminate, skip ui requests like import
		if header.Get(common.BKHTTPRequestFromWeb) != "true" {
			var cancel context.CancelFunc
			// task server has some task with 2 minutes' timeout, so we set the timeout of all servers to 2 minute
			ctx, cancel = context.WithTimeout(ctx, time.Minute*2)
			defer cancel()
		}

		if txnID := header.Get(common.TransactionIdHeader); len(txnID) != 0 {
			// we got a request with transaction info, which is only useful for coreservice.
			ctx = context.WithValue(ctx, common.TransactionIdHeader, txnID)
			ctx = context.WithValue(ctx, common.TransactionTimeoutHeader, header.Get(common.TransactionTimeoutHeader))
		}
		if mode := util.GetHTTPReadPreference(header); mode != common.NilMode {
			ctx = util.SetDBReadPreference(ctx, mode)
			header = util.SetHTTPReadPreference(header, mode)
		}

		restContexts.Kit = &Kit{
			Header:          header,
			Rid:             rid,
			Ctx:             ctx,
			User:            user,
			CCError:         r.ErrorIf.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header)),
			SupplierAccount: owner,
		}

		action.Handler(restContexts)
	}
}
