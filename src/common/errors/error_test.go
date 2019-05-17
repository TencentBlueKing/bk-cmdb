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

package errors

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {

	ccerrmgr, err := New("./examples/errorres")
	if nil != err {
		t.Errorf("failed to create cc error manager, error info is %s", err.Error())
		return
	}

	fmt.Printf("\n[%d]%s\n", 0, ccerrmgr.Error("en", 0))
	fmt.Printf("\n[%d]%s\n", 0, ccerrmgr.Error("cn", 0))
	fmt.Printf("\n[%d]%s\n", 2000, ccerrmgr.Error("en", 20000))
	fmt.Printf("\n[%d]%s\n", 11000, ccerrmgr.Error("en", 11000))
	fmt.Printf("\n[%d]%s\n", 30000, ccerrmgr.Error("en", 30000))
	fmt.Printf("\n[%d]%s\n", 30000, ccerrmgr.Errorf("zn", 30000, "XXX"))
	fmt.Printf("\n[%d]%s\n\n", 10000, ccerrmgr.Errorf("cn", 10000, "XXX"))

	defaultErr := ccerrmgr.CreateDefaultCCErrorIf("cn")

	fmt.Printf("\ndefault[%d]%s\n", 0, defaultErr.Error(0))
	fmt.Printf("\ndefault[%d]%s\n", 2000, defaultErr.Error(20000))
	fmt.Printf("\ndefault[%d]%s\n", 30000, defaultErr.Error(30000))
	fmt.Printf("\ndefault[%d]%s\n", 30000, defaultErr.Errorf(30000, "XXX"))
	fmt.Printf("\ndefault[%d]%s\n\n", 10000, defaultErr.Errorf(10000, "XXX"))

	fmt.Println("test cc error code")
	var errd interface{}
	errd = defaultErr.Error(20000)
	switch e := errd.(type) {
	default:
		fmt.Println("default unkonw")
	case nil:
		fmt.Println("unknown")
	case CCErrorCoder:
		fmt.Printf("error code:%d, str:%s", e.GetCode(), e.Error())

	}

}
