/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"regexp"
	"testing"
)

// TestPatternPodLabels test pattern for pod labels
func TestPatternPodLabels(t *testing.T) {

	testCases := []struct {
		input   string
		matched bool
	}{
		{
			input: `
		{
			"app": "test-1",
			"www.qq.com/lala": "test-2"
		}
		`,
			matched: true,
		},
		{
			input: `
		{}
		`,
			matched: true,
		},
		{
			input: `
		{
			"app": "test-1",
			"www.qq.com/lala": 2
		}
		`,
			matched: false,
		},
		{
			input: `
		{
			"app": "test-1",
		}
		`,
			matched: false,
		},
	}

	for i, test := range testCases {
		t.Logf("test %v", i)
		matched, err := regexp.Match(PatternPodLabels, []byte(test.input))
		if err != nil || matched != test.matched {
			t.Errorf("err %#v, matched %v", err, matched)
		}
	}

}
