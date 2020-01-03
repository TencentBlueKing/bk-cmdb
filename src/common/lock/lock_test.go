package lock

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	redis "gopkg.in/redis.v5"
)

func testRedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:44446", //s.Addr(),
		Password: "xxxx",            // no password set
		DB:       0,                 // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}

func TestLock(t *testing.T) {
	client := testRedisClient()

	lock := NewLocker(client)

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	locked, err := lock.Lock(prefix+"lock1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"lock2", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"lock3", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"lock1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)
}

func TestUnlock(t *testing.T) {

	client := testRedisClient()

	lock := NewLocker(client)

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	locked, err := lock.Lock(prefix+"unlock", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	err = lock.Unlock()
	require.NoError(t, err)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"unlock", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

}

func TestUnlockLockErr(t *testing.T) {
	client := testRedisClient()

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	lockSucc := NewLocker(client)
	locked, err := lockSucc.Lock(prefix+"UnlockLockErr", time.Minute*2)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock := NewLocker(client)
	locked, err = lock.Lock(prefix+"UnlockLockErr", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	// lock failure, unlock no use
	err = lock.Unlock()
	require.NoError(t, err)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"UnlockLockErr", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	err = lockSucc.Unlock()
	require.NoError(t, err)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"UnlockLockErr", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

}

func TestUnlockLockExpire(t *testing.T) {
	client := testRedisClient()

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	lockSucc := NewLocker(client)
	locked, err := lockSucc.Lock(prefix+"UnlockLockExpire", time.Second*2)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock := NewLocker(client)
	locked, err = lock.Lock(prefix+"UnlockLockExpire", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	// lock failure, unlock no use
	err = lock.Unlock()
	require.NoError(t, err)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"UnlockLockExpire", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	time.Sleep(time.Second * 2)

	lock = NewLocker(client)
	locked, err = lock.Lock(prefix+"test1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

}
