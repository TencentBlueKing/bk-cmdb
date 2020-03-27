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
	lastIndex int
}

func NewClientPool(network string, getServer types.GetServerFunc, path string) (*Pool, error) {
	pool := &Pool{
		conns:     make(chan Client, 40),
		getServer: getServer,
	}
	var err error
	var conn Client
	for i := 3; i > 0; i-- {
		conn, err = pool.new()
		if err != nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		err = conn.Ping()
		if err != nil {
			conn.Close()
			time.Sleep(time.Millisecond * 500)
			continue
		}
		pool.put(conn)
		return pool, nil
	}
	return nil, err
}

func (p *Pool) new() (Client, error) {
	var err error
	servers := []string{}
	for i := 3; i > 0; i-- {
		servers, err = p.getServer()
		if err != nil {
			blog.Infof("fetch tmserver address failed: %v, retry 2s later", err)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if len(servers) <= 0 {
			err = fmt.Errorf("service discover returns 0 tmserver address")
			blog.Infof("fetch tmserver address failed: %v, retry 2s later", err)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}

	sort.Strings(servers)

	p.Lock()
	p.lastIndex++
	if p.lastIndex >= len(servers) {
		p.lastIndex = 0
	}
	p.Unlock()

	address, err := util.GetDailAddress(servers[p.lastIndex])
	if err != nil {
		return nil, fmt.Errorf("GetDailAddress %s, failed: %v", servers[p.lastIndex], err)
	}

	return DialHTTPPath("tcp", address, "/txn/v3/rpc")
}

func (p *Pool) pop() Client {
	select {
	case conn, ok := <-p.conns:
		if !ok {
			return nil
		}
		return conn
	default:
		return nil
	}
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

// Call rpc call handler
func (p *Pool) Call(cmd string, input interface{}, result interface{}) (err error) {
	_, err = p.call(cmd, input, result)
	return
}

// CallInfo rpc call handler return client address information
func (p *Pool) CallInfo(cmd string, input interface{}, result interface{}) (addr string, err error) {
	return p.call(cmd, input, result)
}

func (p *Pool) call(cmd string, input interface{}, result interface{}) (addr string, err error) {

	var newConn bool

	conn := p.pop()

	if conn == nil {
		newConn = true
	} else {
		//  链接存在， 如果链接Ping 不同需要重新建立连接
		if pingErr := conn.Ping(); pingErr != nil {
			blog.Errorf("conn rpc connection ping err: %s", pingErr.Error())
			conn.Close()
			newConn = true
		}
	}
	if newConn {
		blog.V(5).Infof("create new rpc connection")
		conn, err = p.new()
		if err != nil {
			return "", err
		}
	}

	err = conn.Call(cmd, input, result)
	addr = conn.TargetID()
	if err != nil {
		if pingErr := conn.Ping(); pingErr == nil {
			p.put(conn)
			return "", err
		}
		conn.Close()
	}
	p.put(conn)
	return addr, nil
}

func (p *Pool) CallStream(cmd string, input interface{}) (*StreamMessage, error) {
	conn := p.pop()
	if conn != nil {
		stream, err := conn.CallStream(cmd, input)
		if err != nil {
			if err != ErrRWTimeout {
				if pingErr := conn.Ping(); pingErr == nil {
					p.put(conn)
					return nil, err
				}
			}
			conn.Close()
		} else {
			p.put(conn)
			return stream, nil
		}
	}

	blog.V(4).Infof("create new rpc connection")
	conn, err := p.new()
	if err != nil {
		return nil, err
	}

	stream, err := conn.CallStream(cmd, input)
	if err != nil {
		if pingErr := conn.Ping(); pingErr == nil {
			p.put(conn)
			return nil, err
		}
		conn.Close()
	}
	p.put(conn)
	return stream, nil
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

func (p *Pool) TargetID() string {
	conn := p.pop()
	if conn != nil {
		p.put(conn)
		return conn.TargetID()
	}
	return ""
}

func (p *Pool) Close() (err error) {
	for {
		select {
		case conn, ok := <-p.conns:
			if !ok {
				return nil
			}
			conn.Close()
		default:
			return nil
		}
	}

}
