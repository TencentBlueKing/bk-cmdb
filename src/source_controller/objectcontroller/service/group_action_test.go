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
	"testing"

	"github.com/emicklei/go-restful"
)

func TestService_CreateUserGroup(t *testing.T) {
	type args struct {
		req  *restful.Request
		resp *restful.Response
	}

	// get mock object
	svc, req, resp := NewRestfulTestCase(`{"k":"v"}`)

	tests := []struct {
		name    string
		service *Service
		args    args
	}{
		// TODO: Add test cases.
		{"", svc, args{req, resp}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.service.CreateUserGroup(tt.args.req, tt.args.resp)
			if tt.args.resp.StatusCode() != 200 {
				t.Fail()
			}
		})
	}
}

func TestService_UpdateUserGroup(t *testing.T) {
	type fields struct {
		Core     *backbone.Engine
		Instance storage.DI
		Cache    *redis.Client
	}
	type args struct {
		req  *restful.Request
		resp *restful.Response
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Service{
				Core:     tt.fields.Core,
				Instance: tt.fields.Instance,
				Cache:    tt.fields.Cache,
			}
			cli.UpdateUserGroup(tt.args.req, tt.args.resp)
		})
	}
}

func TestService_DeleteUserGroup(t *testing.T) {
	type fields struct {
		Core     *backbone.Engine
		Instance storage.DI
		Cache    *redis.Client
	}
	type args struct {
		req  *restful.Request
		resp *restful.Response
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Service{
				Core:     tt.fields.Core,
				Instance: tt.fields.Instance,
				Cache:    tt.fields.Cache,
			}
			cli.DeleteUserGroup(tt.args.req, tt.args.resp)
		})
	}
}

func TestService_SearchUserGroup(t *testing.T) {
	type fields struct {
		Core     *backbone.Engine
		Instance storage.DI
		Cache    *redis.Client
	}
	type args struct {
		req  *restful.Request
		resp *restful.Response
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Service{
				Core:     tt.fields.Core,
				Instance: tt.fields.Instance,
				Cache:    tt.fields.Cache,
			}
			cli.SearchUserGroup(tt.args.req, tt.args.resp)
		})
	}
}
