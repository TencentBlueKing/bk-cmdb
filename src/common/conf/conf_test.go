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

package conf

import (
	"encoding/base64"
	"testing"
)

const (
	conf          = "W21vbmdvZGJdCmhvc3QgPSAxMjcuMC4wLjEKdXNyID0gY2MKcHdkID0gY2MKZGF0YWJhc2UgPSBjbWRiCnBvcnQgPSAyNzAxNwptYXhPcGVuQ29ubnMgPSAzMDAwCm1heElkbGVDb25ucyA9IDEwMDAKCltzbmFwLXJlZGlzXQpob3N0ID0gMTI3LjAuMC4xOjYzNzkKdXNyID0gY2MKcHdkID0gcmVkaXNhdXRoCmRhdGFiYXNlID0gMApjaGFuID0gM19zbmFwc2hvdAoKW3JlZGlzXQpob3N0ID0gMTI3LjAuMC4xOjYzNzkKdXNyID0gY2MKcHdkID0gcmVkaXNhdXRoCmRhdGFiYXNlID0gMAo="
	nodeMongo     = "mongodb"
	nodeSnapRedis = "snap-redis"
	nodeRedis     = "redis"
)

func TestParseConf(t *testing.T) {
	ctx, err := base64.StdEncoding.DecodeString(conf)
	if err != nil {
		t.Errorf("fail to decode config by base64")
		return
	}

	cfg := new(Config)
	cfg.ParseConf(ctx)

	// check result by Config.Read
	type Args struct {
		node string
		key  string
	}

	tests := []struct {
		name string
		args Args
		want string
	}{
		{ //mongdb.host
			args: Args{
				node: nodeMongo,
				key:  "host",
			},
			want: "127.0.0.1",
		},
		{ //mongdb.usr
			args: Args{
				node: nodeMongo,
				key:  "usr",
			},
			want: "cc",
		},
		{ //mongdb.pwd
			args: Args{
				node: nodeMongo,
				key:  "pwd",
			},
			want: "cc",
		},
		{ // snap-redis.host
			args: Args{
				node: nodeSnapRedis,
				key:  "host",
			},
			want: "127.0.0.1:6379",
		},
		{ // snap-redis.usr
			args: Args{
				node: nodeSnapRedis,
				key:  "usr",
			},
			want: "cc",
		},
		{ // snap-redis.pwd
			args: Args{
				node: nodeSnapRedis,
				key:  "pwd",
			},
			want: "redisauth",
		},
		{ // snap-redis.chan
			args: Args{
				node: nodeSnapRedis,
				key:  "chan",
			},
			want: "3_snapshot",
		},
		{ // redis.host
			args: Args{
				node: nodeRedis,
				key:  "host",
			},
			want: "127.0.0.1:6379",
		},
		{ // redis.usr
			args: Args{
				node: nodeRedis,
				key:  "usr",
			},
			want: "cc",
		},
		{ // redis.pwd
			args: Args{
				node: nodeRedis,
				key:  "pwd",
			},
			want: "redisauth"},
		{ // redis.database
			args: Args{
				node: nodeRedis,
				key:  "database",
			},
			want: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ret := cfg.Read(tt.args.node, tt.args.key); ret != tt.want {
				t.Errorf("conf[%s.%s] = %s, but want: %s", tt.args.node, tt.args.key, ret, tt.want)
			}
		})
	}
}

func TestConfig_InitConfig(t *testing.T) {
	type fields struct {
		Configmap map[string]string
		strcet    string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{map[string]string{}, "strcet_string"}, args{"./not_exist_file"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Configmap: tt.fields.Configmap,
				strcet:    tt.fields.strcet,
			}
			c.InitConfig(tt.args.path)
		})
	}
}
