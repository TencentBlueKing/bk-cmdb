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

package y3_13_202611021925

import (
	"regexp"
	"testing"
)

func TestHostIPv6Regex(t *testing.T) {
	re, err := regexp.Compile(hostIPv6Regex)
	if err != nil {
		t.Fatalf("host ipv6 option regex is invalid, err: %v", err)
	}

	matched := []string{
		// single ipv6 value, the value of the process bind info in the issue is one of them
		"2001:0db8:85a3:0000:0000:8a2e:0370:7331",
		"2001:db8::1",
		"2001:DB8::1",
		"::ffff:192.168.1.1",
		"fe80::1%eth0",
		// multiple ipv6 values, host ipv6 attributes are stored as arrays
		"2001:db8::1,2001:db8::2",
		"2001:db8::1,2001:db8::2,2001:db8::3,2001:db8::4,2001:db8::5,2001:db8::6",
	}

	for _, value := range matched {
		if !re.MatchString(value) {
			t.Errorf("%q should match the host ipv6 option", value)
		}
	}

	unmatched := []string{
		"1.2.3.4",
		// multiple ipv6 values without the comma between them
		"2001:db8::12001:db8::2",
		"2001:db8::1,",
		",2001:db8::1",
		"2001:db8::1,,2001:db8::2",
		"2001:db8::1, 2001:db8::2",
		"",
	}

	for _, value := range unmatched {
		if re.MatchString(value) {
			t.Errorf("%q should not match the host ipv6 option", value)
		}
	}
}
