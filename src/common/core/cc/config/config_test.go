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
 
package config

import (
	//"fmt"
	"testing"
)

func TestGetAddress(t *testing.T) {
	conf := NewCCAPIConfig()

	type Args struct {
		addrport string
	}

	type Wants struct {
		ip  string
		err error
	}

	tests := []struct {
		name string
		arg  Args
		want Wants
	}{
		{
			arg: Args{
				addrport: "127.0.0.1:8081",
			},
			want: Wants{
				ip: "127.0.0.1",
				//err: nil,
			},
		},
		{
			arg: Args{
				addrport: "127.0.0.1",
			},
			want: Wants{
				ip: "",
				//err: fmt.Errorf("the value of flag[AddrPort: 127.0.0.1] is wrong"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf.AddrPort = tt.arg.addrport
			ip, _ := conf.GetAddress()
			if ip != tt.want.ip {
				t.Errorf("GetAddress() ip:%s,  but want ip:%s", ip, tt.want.ip)
			}
		})
	}
}

func TestGetPort(t *testing.T) {
	conf := NewCCAPIConfig()

	type Args struct {
		addrport string
	}

	tests := []struct {
		name string
		arg  Args
		want uint
	}{
		{
			arg: Args{
				addrport: "127.0.0.1:8081",
			},
			want: 8081,
		},
		{
			arg: Args{
				addrport: "127.0.0.1",
			},
			want: 0,
		},
		{
			arg: Args{
				addrport: "127.0.0.1:abc",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf.AddrPort = tt.arg.addrport
			if port, _ := conf.GetPort(); port != tt.want {
				t.Errorf("GetPort() port:%d, but want: %d", port, tt.want)
			}
		})
	}
}
