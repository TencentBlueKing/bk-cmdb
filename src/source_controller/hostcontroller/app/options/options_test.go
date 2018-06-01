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
 
package options

import (
    "testing"

    "github.com/spf13/pflag"
)

func TestNewServerOption(t *testing.T) {
    if so := NewServerOption(); so == nil || so.ServConf == nil {
        t.Fatal("server option is nil or ServerOption.ServConf is nil")
    }
}

func TestServerOption_AddFlags(t *testing.T) {
    so := NewServerOption()

    so.AddFlags(pflag.CommandLine)
    if so.ServConf.AddrPort != "127.0.0.1:50002" {
        t.Errorf("AddrPort not as expected: %s", so.ServConf.AddrPort)
    }
    if so.ServConf.RegDiscover != "" {
        t.Errorf("RegDiscover not as expected: %s", so.ServConf.RegDiscover)
    }
    if so.ServConf.ExConfig != "" {
        t.Errorf("ExConfig not as expected: %s", so.ServConf.ExConfig)
    }
}
