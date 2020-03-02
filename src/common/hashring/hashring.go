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

package hashring

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type HashFunc func(data []byte) uint32

type HashRing struct {
	hash     HashFunc
	replicas int               // 复制因子,单个节点的虚拟节点数
	keys     UInt32Slice       // 已排序的节点哈希值切片
	hashMap  map[uint32]string // 键是节点哈希值，值是节点Key
	mu       sync.RWMutex
}

// NewHashRing 创建哈希环
func NewHashRing(replicas int, fn HashFunc) *HashRing {
	m := &HashRing{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	// 默认使用CRC32算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty 判断哈希环是否为空
func (h *HashRing) IsEmpty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.keys) == 0
}

// Add 方法用来添加节点，参数为节点key
func (h *HashRing) Add(keys ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, key := range keys {
		// 结合复制因子计算所有虚拟节点的hash值，并存入h.keys中，同时在h.hashMap中保存哈希值和key的映射
		for i := 0; i < h.replicas; i++ {
			hash := h.hash([]byte(key + "_" + strconv.Itoa(i)))
			h.keys = append(h.keys, hash)
			h.hashMap[hash] = key
		}
	}
	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(h.keys)
}

// Del 方法用来删除节点，参数为节点key
func (h *HashRing) Del(keys ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	hashValues := make(map[uint32]bool)
	for _, key := range keys {
		// 删除h.hashMap中保存的哈希值和key映射
		for i := 0; i < h.replicas; i++ {
			hash := h.hash([]byte(key + "_" + strconv.Itoa(i)))
			hashValues[hash] = true
			delete(h.hashMap, hash)
		}
	}
	// 删除相应虚拟节点的哈希值
	var hashKeys UInt32Slice = make([]uint32, 0)
	for _, v := range h.keys {
		if _, ok := hashValues[v]; !ok {
			hashKeys = append(hashKeys, v)
		}
	}
	h.keys = hashKeys
}

// Clear 方法用来清空哈希环里的数据
func (h *HashRing) Clear(keys ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.keys = UInt32Slice{}
	h.hashMap = map[uint32]string{}
}

// Get 方法根据给定的对象获取最靠近它的那个节点key
func (h *HashRing) Get(key string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.IsEmpty() {
		return ""
	}

	hash := h.hash([]byte(key))

	// 通过二分查找获取第一个节点hash值大于对象hash值的节点
	idx := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= hash })

	// 如果查找结果大于节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，则放入第一个节点中
	if idx == len(h.keys) {
		idx = 0
	}

	return h.hashMap[h.keys[idx]]
}
