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

package redisclient

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/storage"

	redis "gopkg.in/redis.v5"
)

type RedisConfig struct {
	Address    string
	User       string
	Password   string
	Database   string
	Port       string
	MasterName string
}

func NewFromConfig(cfg RedisConfig) (*redis.Client, error) {
	dbNum, err := strconv.Atoi(cfg.Database)
	if nil != err {
		return nil, err
	}
	var client *redis.Client
	if cfg.MasterName == "" {
		if !strings.Contains(cfg.Address, ":") {
			cfg.Address = cfg.Address + ":" + cfg.Port
		}
		option := &redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       dbNum,
			PoolSize: 100,
		}
		client = redis.NewClient(option)
	} else {
		hosts := strings.Split(cfg.Address, ",")
		option := &redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			SentinelAddrs: hosts,
			Password:      cfg.Password,
			DB:            dbNum,
			PoolSize:      100,
		}
		client = redis.NewFailoverClient(option)
	}

	err = client.Ping().Err()
	if err != nil {
		return nil, err
	}

	return client, err
}

type Redis struct {
	host    string
	port    string
	usr     string
	pwd     string
	dbName  string
	session *redis.Client
}

func NewRedis(host, port, usr, pwd, database string) (storage.DI, error) {
	RedisConfig := new(Redis)
	RedisConfig.host = host
	RedisConfig.port = port
	RedisConfig.usr = usr
	RedisConfig.pwd = pwd
	RedisConfig.dbName = database
	RedisConfig.session = nil

	return RedisConfig, nil
}

//open con
func (r *Redis) Open() error {
	dbNum, _ := strconv.Atoi(r.dbName)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     r.host + ":" + r.port,
		PoolSize: 100,
		Password: r.pwd,
		DB:       dbNum,
	})

	err := redisClient.Ping().Err()
	if err != nil {
		return err
	}

	r.session = redisClient

	return nil
}

// Ping will ping the db server
func (m *Redis) Ping() error {
	return m.session.Ping().Err()
}

// GetSession
func (r *Redis) GetSession() interface{} {
	return r.session
}

func (r *Redis) id(object interface{}) int {
	return 0
}

func (r *Redis) Insert(cName string, data interface{}) (int, error) {

	var err error
	mapData, ok := data.(common.KvMap)
	key, ok := mapData["key"].(string)
	if !ok {
		return 0, errors.New("params key can not be empty")
	}
	cName = strings.ToLower(cName)
	switch cName {
	case "set":
		value, ok := mapData["value"]
		exp, ok := mapData["expire"].(time.Duration)
		if !ok {
			exp = time.Duration(0)
		}
		err = r.session.Set(key, value, exp).Err()
	case "setnx":
		value, ok := mapData["value"]
		exp, ok := mapData["expire"]
		if !ok {
			exp = time.Duration(0)
		}
		err = r.session.SetNX(key, value, exp.(time.Duration)).Err()
	case "decr":
		err = r.session.Decr(key).Err()
	case "decrby":
		decr, _ := util.GetInt64ByInterface(mapData["decr"])
		err = r.session.DecrBy(key, decr).Err()
	case "incr":
		var result int64
		result, err = r.session.Incr(key).Result()
		return (int)(result), err
	case "incrby":
		incr, _ := util.GetInt64ByInterface(mapData["incr"])
		err = r.session.IncrBy(key, incr).Err()
	case "expire":
		exp, ok := mapData["expire"].(time.Duration)
		if !ok {
			return 0, errors.New("params Expire can not be empty")
		}
		cmd := r.session.Expire(key, exp)
		err = cmd.Err()
		if cmd.Val() {
			return 1, err
		} else {
			return 0, err
		}
	case "expireat":
		exp, ok := mapData["expire"].(time.Time)
		if !ok {
			return 0, errors.New("params Expire can not be empty")
		}
		cmd := r.session.ExpireAt(key, exp)
		err = cmd.Err()
		if cmd.Val() {
			return 1, err
		} else {
			return 0, err
		}

	case "hset":
		value, _ := mapData["value"]
		field, _ := mapData["field"].(string)
		err = r.session.HSet(key, field, value).Err()
	case "hmset":
		fields, _ := mapData["fields"].(map[string]string)
		err = r.session.HMSet(key, fields).Err()
	case "hsetnx":
		field, _ := mapData["field"].(string)
		value, _ := mapData["value"]
		err = r.session.HSetNX(key, field, value).Err()
	case "rpush":
		values, errConv := util.GetMapInterfaceByInerface(mapData["values"])
		if nil != errConv {
			err = errConv
		} else {
			err = r.session.RPush(key, values...).Err()
		}

	case "lset":
		index, _ := mapData["index"].(int64)
		value, _ := mapData["value"]
		err = r.session.LSet(key, index, value).Err()

	case "sadd":
		values, errConv := util.GetMapInterfaceByInerface(mapData["values"])
		if nil != errConv {
			err = errConv
		} else {
			err = r.session.SAdd(key, values...).Err()
		}
	default:
		err = errors.New("no support method")
	}
	if nil != err {
		return 0, err
	}
	return 1, nil
}

func (r *Redis) InsertMuti(cName string, data ...interface{}) error {
	return errors.New("no support method")

}

func (r *Redis) UpdateByCondition(cName string, oldobj, newobj interface{}) error {
	return errors.New("no support method")
}

func (r *Redis) GetOneByCondition(cName string, fields []string, selector, results interface{}) error {
	ret, _ := results.(*interface{})
	var err error
	cName = strings.ToLower(cName)
	mapData, _ := selector.(common.KvMap)
	key, _ := mapData["key"].(string)
	switch cName {
	case "get":
		cmd := r.session.Get(key)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "hget":
		field, _ := mapData["field"].(string)
		cmd := r.session.HGet(key, field)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "hgetall":
		cmd := r.session.HGetAll(key)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "getrange":
		start, _ := mapData["start"].(int64)
		end, _ := mapData["end"].(int64)
		cmd := r.session.GetRange(key, start, end)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "lrange":
		start, _ := util.GetInt64ByInterface(mapData["start"])
		end, _ := util.GetInt64ByInterface(mapData["end"])
		if 0 == end {
			err = errors.New("params end is requred")
		} else {
			cmd := r.session.LRange(key, start, end)
			err = cmd.Err()
			*ret, _ = cmd.Result()
		}

	case "blpop":
		timeout, _ := mapData["expire"].(time.Duration)
		keys, _ := mapData["key"].([]string)
		cmd := r.session.BLPop(time.Duration(timeout), keys...)
		err = cmd.Err()
		if err == nil {
			rett := results.(*[]string)
			*rett, _ = cmd.Result()
		}
	case "hexists":
		field, _ := mapData["field"].(string)
		cmd := r.session.HExists(key, field)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "hlen":
		cmd := r.session.HLen(key)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "exists":
		//key, _ := selector.(string)
		cmd := r.session.Exists(key)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "ttl":
		cmd := r.session.TTL(key)
		err = cmd.Err()
		*ret, _ = cmd.Result()
	case "smembers":
		cmd := r.session.SMembers(key)
		err = cmd.Err()
		if err == nil {
			rett := results.(*[]string)
			*rett, _ = cmd.Result()
		}
	default:
		err = errors.New("no support method")
	}
	if redis.Nil == err {
		err = nil
	}
	return err
}

func (r *Redis) GetMutilByCondition(cName string, fields []string, selector, results interface{}, sort string, skip, limit int) error {

	return errors.New("no support method")
}

func (r *Redis) GetCntByCondition(cName string, selector interface{}) (cnt int, err error) {

	return 0, errors.New("no support method")
}

func (m *Redis) GetIncID(cName string) (int64, error) {
	return 0, errors.New("not support method")
}

func (r *Redis) DelByCondition(cName string, delselector interface{}) error {
	var err error
	cName = strings.ToLower(cName)

	switch cName {
	case "del":
		keys, _ := delselector.([]string)
		//if ok {
		cmd := r.session.Del(keys...)
		err = cmd.Err()
		/*} else {
			key, ok := delselector.(string)
			if ok {
				cmd := r.session.Del(key)
				err = cmd.Err()
			}
		}*/
	case "hdel":
		mapData, _ := delselector.(common.KvMap)
		key, _ := mapData["key"].(string)
		fields, _ := mapData["fields"].([]string)
		cmd := r.session.HDel(key, fields...)
		err = cmd.Err()
	case "srem":
		mapData, _ := delselector.(common.KvMap)
		key, _ := mapData["key"].(string)
		fields, _ := mapData["values"].([]interface{})
		blog.Infof("srem key %v, values %v", key, fields)
		cmd := r.session.SRem(key, fields...)
		err = cmd.Err()
	default:
		err = errors.New("no support method")

	}
	return err
}

//判断表是否存在
func (r *Redis) HasTable(tableName string) (bool, error) {
	return false, errors.New("no support method")
}

//执行原始的sql语句，并不会返回数据， 只会是否执行出错
func (r *Redis) ExecSql(cmd interface{}) error {
	return errors.New("no support method")
}

func (r *Redis) CreateTable(sql string) error {
	return errors.New("no support method")
}

func (r *Redis) Index(tableName string, index *storage.Index) error {
	return errors.New("no support method")
}

func (r *Redis) DropTable(tableName string) error {
	return errors.New("no support method")
}

func (r *Redis) HasFields(tableName, field string) (bool, error) {
	return false, errors.New("no support method")
}

//新加字段， 表名，字段名,字段类型, 附加描述（是否为空， 默认值）
func (r *Redis) AddColumn(tableName string, column *storage.Column) error {
	return errors.New("no support method")
}

func (r *Redis) ModifyColumn(tableName, oldName, newColumn string) error {
	return errors.New("no support method")
}

func (r *Redis) DropColumn(tableName, field string) error {
	return errors.New("no support method")
}

//GetType 获取操作db的类
func (r *Redis) GetType() string {
	return storage.DI_REDIS
}

// Close close session
func (r *Redis) Close() {
	if r.session != nil {
		r.session.Close()
	}
}

// IsDuplicateErr returns whether err is duplicate error
func (r *Redis) IsDuplicateErr(err error) bool {
	return false
}

// IsNotFoundErr returns whether err is not found error
func (r *Redis) IsNotFoundErr(err error) bool {
	return redis.Nil == err
}
