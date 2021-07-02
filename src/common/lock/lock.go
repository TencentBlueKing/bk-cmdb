package lock

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal/redis"

	"github.com/rs/xid"
)

type lock struct {
	cache redis.Client
	key   string
	// 是否需要释放key
	needUnlock bool
	isFirst    bool
}

type mlock struct {
	cache redis.Client
	keys  []string
	// if release lock or not
	needUnlock bool
	isFirst    bool
}

// Locker redis atomic lock
type Locker interface {
	// Lock can lock one
	Lock(key StrFormat, expire time.Duration) (locked bool, err error)
	Unlock() error
}

type MLocker interface {
	MLock(rid string, retrytimes int, expire time.Duration, values ...StrFormat) (locked bool, err error)
	MUnlock() error
}

func NewLocker(cache redis.Client) Locker {

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
	locked, err = l.cache.SetNX(context.Background(), l.key, uuid, expire).Result()
	// locked sucess , can unlock
	if locked {
		l.needUnlock = true
	}
	return locked, err
}

func NewMLocker(cache redis.Client) MLocker {

	return &mlock{
		isFirst:    false,
		cache:      cache,
		keys:       []string{},
		needUnlock: false,
	}
}

// parse setnx result，generate key-lockResult pair
// the return value like: " setnx cc:v3:61_felen_cx_run_felen_rt c3eqnlv6bt34kavtdotg: true "
// so we should split use space to get the key ,use ":" to get the result true or false
func getExecSetNxBoolResult(result string) (string, bool) {
	keySlice := strings.Split(result, " ")
	if len(keySlice) < 3 {
		return "", false
	}

	resultSlice := strings.Split(result, ":")
	if len(resultSlice) < 2 {
		return "", false
	}

	ResString := strings.TrimSpace(resultSlice[len(resultSlice)-1])
	if ResString == "true" {
		return keySlice[1], true
	}
	return keySlice[1], false
}

func (l *mlock) MLock(rid string, retry int, expire time.Duration, keys ...StrFormat) (locked bool, err error) {
	if l.isFirst {
		return false, errors.New("repeat lock")
	}

	var (
		bResultFlag bool
		delKeys     []string
	)
	pipeRes := make(map[bool][]string)
	l.isFirst = true

	for i := 0; i < retry; i++ {
		bPipeResultFlag := false
		delKeys = []string{}
		pipe := l.cache.TxPipeline(context.Background())
		l.keys = []string{}
		for _, k := range keys {
			// splice command
			key := fmt.Sprintf("%s%s", common.BKCacheKeyV3Prefix, k)
			uuid := xid.New().String()
			l.keys = append(l.keys, key)
			pipe.SetNX(key, uuid, 0)
			pipe.Expire(key, expire)
		}
		res, err := pipe.Exec()
		if err != nil {
			// exec error try it again
			continue
		}

		for k, r := range res {
			// mlock contain setnx and expire two commonds, you should only mark setnx if success or not
			if k%2 == 0 {
				key, bResult := getExecSetNxBoolResult(r.String())
				if !bResult {
					// obtain lock fail ones
					pipeRes[false] = append(pipeRes[false], key)
				} else {
					// obtain lock success ones
					bPipeResultFlag = true
					pipeRes[true] = append(pipeRes[true], key)
				}
			}
		}
		if bPipeResultFlag && len(pipeRes[false]) > 0 {
			// if some setnx fail, need release the success ones
			err := l.cache.Del(context.Background(), pipeRes[true]...).Err()
			if err != nil {
				// if del fail, need to del it when unlock
				for _, v := range pipeRes[true] {
					delKeys = append(delKeys, v)
				}
				blog.Errorf("delete key fail. the key: %v,rid: %s", pipeRes[true], rid)
			}
		} else {
			// obtain lock success
			bResultFlag = true
			break
		}
		for k := range pipeRes {
			// release the key-lockResult  pair
			delete(pipeRes, k)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if bResultFlag {
		// obtain lock success, you should release it
		l.needUnlock = true
		return true, nil
	}
	// delete fail,unlock retry delete it
	if len(delKeys) > 0 {
		l.needUnlock = true
		l.keys = delKeys
	}
	// set map nil for gc
	pipeRes = nil
	return false, errors.New("obtain lock fail")
}

func (l *lock) Unlock() error {
	// locked sucess , can unlock
	if !l.needUnlock {
		return nil
	}
	return l.cache.Del(context.Background(), l.key).Err()
}

func (l *mlock) MUnlock() error {
	// locked sucess , can unlock
	if !l.needUnlock {
		return nil
	}
	return l.cache.Del(context.Background(), l.keys...).Err()
}
