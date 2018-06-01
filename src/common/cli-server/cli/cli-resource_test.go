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
 
package cli

import (
	"testing"
)

const (
	mysqlType  = "mysql"
	mongdbType = "mongdb"
	redisType  = "redis"
)

func TestGetDataCli(t *testing.T) {
	t.Log("begin to test GetDataCli")

	cli := NewCliResource()

	type args struct {
		conf  map[string]string
		dType string
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{ // mysql
			name: "test-mysql-client",
			args: args{
				conf: map[string]string{
					mysqlType + ".host":      "127.0.0.1",
					mysqlType + ".port":      "3306",
					mysqlType + ".usr":       "mysql",
					mysqlType + ".pwd":       "mysql",
					mysqlType + ".database":  "test",
					mysqlType + ".mechanism": "mechanism",
				},
				dType: mysqlType,
			},
			want: nil,
		},
		{ // mongdb
			name: "test-mongdb-client",
			args: args{
				conf: map[string]string{
					mongdbType + ".host":      "127.0.0.1",
					mongdbType + ".port":      "27017",
					mongdbType + ".usr":       "mongdb",
					mongdbType + ".pwd":       "mongdb",
					mongdbType + ".database":  "test",
					mongdbType + ".mechanism": "mechanism",
				},
				dType: mongdbType,
			},
			want: nil,
		},
		{ // redis
			name: "test-redis-client",
			args: args{
				conf: map[string]string{
					redisType + ".host":      "127.0.0.1",
					redisType + ".port":      "6379",
					redisType + ".usr":       "redis",
					redisType + ".pwd":       "redis",
					redisType + ".database":  "test",
					redisType + ".mechanism": "mechanism",
				},
				dType: redisType,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ret := cli.GetDataCli(tt.args.conf, tt.args.dType); ret != tt.want {
				t.Errorf("GetDataCli() = %v, want: %v", ret, tt.want)
			}
		})
	}
}
