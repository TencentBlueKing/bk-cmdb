/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rpc

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/storage/types"
)

type Pool struct {
	sync.Mutex
	conns chan Client

	getServer types.GetServerFunc
	servers   []string
	lastIndex int

	enableTransaction bool
}

func NewClientPool(network string, getServer types.GetServerFunc, path string) (*Pool, error) {
	pool := &Pool{
		conns:     make(chan Client, 3),
		getServer: getServer,
	}
	conn, err := pool.new()
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	pool.put(conn)
	return pool, err
}

func (p *Pool) new() (Client, error) {
	var err error
	servers := []string{}
	for i := 3; i > 0; i-- {
		servers, err = p.getServer()
		if err != nil {
			blog.Infof("fetch tmserver address failed: %v, retry 2s later", err)
			time.Sleep(time.Second * 2)
			continue
		}
		if len(servers) <= 0 {
			err = fmt.Errorf("service discover returns 0 tmserver address")
			blog.Infof("fetch tmserver address failed: %v, retry 2s later", err)
			time.Sleep(time.Second * 2)
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}

	sort.Strings(servers)

	p.Lock()
	p.servers = servers

	p.lastIndex++
	if p.lastIndex >= len(p.servers) {
		p.lastIndex = 0
	}

	address, err := util.GetDailAddress(p.servers[p.lastIndex])
	if err != nil {
		p.Unlock()
		return nil, fmt.Errorf("GetDailAddress %s, failed: %v", p.servers[p.lastIndex], err)
	}
	p.Unlock()

	return DialHTTPPath("tcp", address, "/txn/v3/rpc")
}

func (p *Pool) pop() (conn Client) {
	conn, _ = <-p.conns
	return conn

}

func (p *Pool) put(conn Client) {
	select {
	case p.conns <- conn:
	default:
		go func() {
			for i := 3; i > 0; i-- {
				select {
				case p.conns <- conn:
					return
				default:
					time.Sleep(time.Second)
				}
			}
			// close the connection, because the idle connection is full
			blog.Warnf("idle connection is full, drop connection")
			conn.Close()
		}()
	}
}

func (p *Pool) Call(cmd string, input interface{}, result interface{}) (err error) {
	blog.V(4).Infof("calling %s, %v", cmd, input)
	defer blog.V(4).Infof("calling %s success", cmd)
	conn := p.pop()
	if conn != nil {
		err = conn.Call(cmd, input, result)
		if err == ErrRWTimeout {

		}
		if err != nil {
			if err != ErrRWTimeout {
				if pingErr := conn.Ping(); pingErr == nil {
					p.put(conn)
					blog.V(4).Infof("restore connection on error: %v", err)
					return err
				}
			}
			conn.Close()
		} else {
			p.put(conn)
			blog.V(4).Infof("restore connection on success")
			return
		}
	}

	blog.V(4).Infof("create new rpc connection")
	conn, err = p.new()
	if err != nil {
		return err
	}

	err = conn.Call(cmd, input, result)
	if err != nil {
		if pingErr := conn.Ping(); pingErr == nil {
			p.put(conn)
			return err
		}
		conn.Close()
	}
	p.put(conn)
	return nil
}

func (p *Pool) Ping() (err error) {
	conn := p.pop()
	if conn != nil {
		err = conn.Ping()
		if err != nil {
			conn.Close()
		} else {
			p.put(conn)
			return nil
		}
	}

	conn, err = p.new()
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		conn.Close()
	} else {
		p.put(conn)
		return nil
	}

	return err
}

func (p *Pool) Close() (err error) {
	for {
		conn, ok := <-p.conns
		if !ok {
			break
		}
		conn.Close()
	}
	return nil
}
