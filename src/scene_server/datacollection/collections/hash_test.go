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

package collections

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"stathat.com/c/consistent"
)

var (
	ipRanges = [][]int{{607649792, 608174079}, //36.56.0.0-36.63.255.255
		{1038614528, 1039007743},   //61.232.0.0-61.237.255.255
		{1783627776, 1784676351},   //106.80.0.0-106.95.255.255
		{2035023872, 2035154943},   //121.76.0.0-121.77.255.255
		{2078801920, 2079064063},   //123.232.0.0-123.235.255.255
		{-1950089216, -1948778497}, //139.196.0.0-139.215.255.255
		{-1425539072, -1425014785}, //171.8.0.0-171.15.255.255
		{-1236271104, -1235419137}, //182.80.0.0-182.92.255.255
		{-770113536, -768606209},   //210.25.0.0-210.47.255.255
		{-569376768, -564133889},   //222.16.0.0-222.95.255.255
	}

	maxClouid = 10

	dataCount = 10000
)

func numToIP(num int) string {
	arr := make([]int, 4)
	arr[0] = (num >> 24) & 0xff
	arr[1] = (num >> 16) & 0xff
	arr[2] = (num >> 8) & 0xff
	arr[3] = num & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func RandIP() string {
	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return numToIP(ipRanges[idx][0] + rand.Intn(ipRanges[idx][1]-ipRanges[idx][0]))
}

func RandCloudid() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(maxClouid)
}

func TestHashing(t *testing.T) {
	// data resource.
	dataHashs := []string{}
	for i := 0; i < dataCount; i++ {
		dataHashs = append(dataHashs, fmt.Sprintf("%d:%s", RandCloudid(), RandIP()))
	}

	// consistent hashring.
	c := consistent.New()

	stats1 := make(map[string]int)
	c.Add("127.0.0.1:80")
	c.Add("127.0.0.2:80")
	c.Add("127.0.0.3:80")

	t.Logf("\nNode [1, 2, 3]")
	for _, hash := range dataHashs {
		node, err := c.Get(hash)
		if err != nil {
			t.Fatal(err)
		}
		stats1[node]++
		t.Logf("%s => %s\n", hash, node)
	}

	stats2 := make(map[string]int)
	c.Add("127.0.0.4:80")
	c.Add("127.0.0.5:80")

	t.Logf("\nNode [1, 2, 3, 4, 5]")
	for _, hash := range dataHashs {
		node, err := c.Get(hash)
		if err != nil {
			t.Fatal(err)
		}
		stats2[node]++
		t.Logf("%s => %s\n", hash, node)
	}

	stats3 := make(map[string]int)
	c.Remove("127.0.0.3:80")

	t.Logf("\nNode [1, 2, 4, 5]")
	for _, hash := range dataHashs {
		node, err := c.Get(hash)
		if err != nil {
			t.Fatal(err)
		}
		stats3[node]++
		t.Logf("%s => %s\n", hash, node)
	}
	t.Logf("STAT:\nstats1:%+v\nstats2:%+v\nstats3:%+v", stats1, stats2, stats3)
}
