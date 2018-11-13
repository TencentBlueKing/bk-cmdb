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
 
package utils

import (
	"testing"
	"net/url"
	"github.com/stretchr/testify/assert"
)

func TestSliceStrToInt(t *testing.T){
	strs := []string{"0","1","2","3"}
	expects := []int{0,1,2,3}
	ints,err := SliceStrToInt(strs)
	assert.Nil(t,err)

	for i,n :=range ints{
		assert.Equal(t,expects[i],n)
	}
}

func TestValidateFormData(t *testing.T){
	data := url.Values{
		"key1":make([]string,0),
	}

	ok,_ := ValidateFormData(data,[]string{"key1"})
	assert.Equal(t,false,ok)

	data = url.Values{
		"key1":[]string{"val1"},
	}
	ok,_ = ValidateFormData(data,[]string{"key1"})
	assert.Equal(t,true,ok)
}