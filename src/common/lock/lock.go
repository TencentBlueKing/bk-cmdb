package lock

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/xid"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

type lock struct {
	cache *redis.Client
	key   string
	// 是否需要释放key
	needUnlock bool
	isFirst    bool
}

// Locker redis atomic lock
type Locker interface {
	// Lock can lock one
	Lock(key StrFormat, expire time.Duration) (looked bool, err error)
	Unlock() error
}

func NewLocker(cache *redis.Client) Locker {

	return &lock{
		isFirst:    false,
		cache:      cache,
		key:        "",
		needUnlock: false,
	}
}

// Lock can lock one, key from GetLockKey function
func (l *lock) Lock(key StrFormat, expire time.Duration) (locked bool, err error) {
	if l.isFirst {
		return false, fmt.Errorf("repeat lock")
	}
	l.isFirst = true
	l.key = fmt.Sprintf("%s%s", common.BKCacheKeyV3Prefix, key)

	// 不能一样，一样的话，会提示设置成功
	uuid := xid.New().String()
	locked, err = l.cache.SetNX(l.key, uuid, expire).Result()
	// locked sucess , can unlock
	if locked {
		l.needUnlock = true
	}
	return locked, err
}

func (l *lock) Unlock() error {
	// locked sucess , can unlock
	if !l.needUnlock {
		return nil
	}
	return l.cache.Del(l.key).Err()
}

// TryGetLock 抢占锁
func TryGetLock(redisCli *redis.Client, lockName StrFormat, lockTimeOut time.Duration) (bool, error) {
	uuid := xid.New().String()
	if ok, err := redisCli.SetNX(string(lockName), uuid, lockTimeOut).Result(); err != nil && err != redis.Nil {
		return false, err
	} else if ok {
		return true, nil
	}
	return false, nil
}

// ReleaseLock 释放抢占到的锁通过抢占到的锁的Name
func ReleaseLock(redisCli *redis.Client, key StrFormat) error {
	value, err := redisCli.Get(string(key)).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("release lock failed: %v, key(%s)", err, key)
	}
	if value == "" {
		return nil
	}
	blog.Info("[ReleaseLock] key: %s", string(key))
	return redisCli.Del(string(key)).Err()
}

const (
	appendLockKeyToTranRecordError = "append lockKey to transaction record use redis"

	// appendLockKeyToTranRecordScript is used to use redis to record transaction get lockKeys, when transaction end,
	// to release these lockKeys
	// KEYS[1]: the key used to save the current transaction preemption lock
	// KEYS[2]: the get key of the current transaction
	// KEYS[3]: expire time
	// ARGV[1]: set failure error
	appendLockKeyToTranRecordScript = `
local value = redis.pcall('get', KEYS[1]);

if (value == false) then
	value = KEYS[2];
else
	value = string.format("%s,%s", value, KEYS[2]);
end;

local ok = redis.pcall("set", KEYS[1], value, "px", KEYS[3])

if ok['ok'] ~= 'OK' then
    return ARGV[1]
end;

return 
`
)

// AppendLockKeyToTransaction 往用于保存当前事务抢占到的锁列表，追加抢占到的锁
func AppendLockKeyToTransaction(redisCli *redis.Client, transactionKey StrFormat, key StrFormat, lockTimeOut time.Duration) error {
	lockTimeOutStr := strconv.Itoa(int(lockTimeOut))
	keys := []string{string(transactionKey), string(key), lockTimeOutStr}
	_, err := redisCli.Eval(appendLockKeyToTranRecordScript, keys, appendLockKeyToTranRecordError).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}
