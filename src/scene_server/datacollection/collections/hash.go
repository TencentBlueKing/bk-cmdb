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
	"net/url"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common/blog"

	"stathat.com/c/consistent"
)

const (
	// defaultUpdateInterval is default service node hash values update interval.
	defaultUpdateInterval = 3 * time.Second
)

// Hash is data hash that updates target nodes in dynamic mode,
// and calculates node base on hash key of data.
type Hash struct {
	// local datacollection server node hash value.
	localHashValue string

	// consistent hashring.
	consistent *consistent.Consistent

	// discovery is cc service discovery handler.
	discovery discovery.DiscoveryInterface

	// nodes records datacollection nodes infos, hash value -> address.
	nodes map[string]string
}

// NewHash creates a new hash object with local node hash value.
func NewHash(ip string, port uint, discovery discovery.DiscoveryInterface) *Hash {
	h := &Hash{
		localHashValue: fmt.Sprintf("%s:%d", ip, port),
		consistent:     consistent.New(),
		discovery:      discovery,
		nodes:          make(map[string]string),
	}
	go h.updateLoop()

	return h
}

// IsMatch matchs hash key base on dynamic hashring values, and
// return a bool that marks if the data hash is matched the local node.
func (h *Hash) IsMatch(hash string) bool {
	nodeHashValue, err := h.consistent.Get(hash)
	if err != nil {
		blog.Errorf("Hash| can't get target node hash, %+v", err)
		return false
	}

	if h.localHashValue != nodeHashValue {
		// not handled in this datacollection node.
		return false
	}

	// match local hash value, it would handled in this node.
	return true
}

// updateLoop keeps discovering datacollection instances and update local consistent.
func (h *Hash) updateLoop() {
	ticker := time.NewTicker(defaultUpdateInterval)
	defer ticker.Stop()

	isFirst := true

	for {
		if !isFirst {
			<-ticker.C
		}
		isFirst = false

		// discovery.
		servers, err := h.discovery.DataCollect().GetServers()
		if err != nil {
			blog.Errorf("Hash| update services hash values, %+v", err)
			continue
		}
		blog.V(4).Infof("Hash| discovery newest servers now, %+v", servers)

		// query.
		newest := make(map[string]string)

		for _, svr := range servers {
			// parse servers address, format: "scheme://ip:port".
			u, err := url.Parse(svr)
			if err != nil {
				blog.Errorf("Hash| update newest servers, node: %s, %+v", svr, err)
				continue
			}

			if len(u.Hostname()) == 0 || len(u.Port()) == 0 {
				continue
			}

			// update newest and current node records.
			hashValue := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
			newest[hashValue] = svr
		}

		// update.
		for hashValue, svr := range newest {
			if _, isExist := h.nodes[hashValue]; !isExist {
				// new node, add to consistent, do not add more replicas.
				blog.Infof("Hash| add new consistent hash node, %s", hashValue)
				h.consistent.Add(hashValue)
			}

			// upsert.
			h.nodes[hashValue] = svr
		}

		// delete.
		for hashValue := range h.nodes {
			if _, isExist := newest[hashValue]; !isExist {
				blog.Infof("Hash| remove old consistent hash node, %s", hashValue)
				h.consistent.Remove(hashValue)
				delete(h.nodes, hashValue)
			}
		}
		blog.V(4).Infof("Hash| sync consistent hash done, members %+v", h.consistent.Members())
	}
}
