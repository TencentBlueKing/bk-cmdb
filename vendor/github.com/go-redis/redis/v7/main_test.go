package redis_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	redisPort          = "6380"
	redisAddr          = ":" + redisPort
	redisSecondaryPort = "6381"
)

const (
	ringShard1Port = "6390"
	ringShard2Port = "6391"
	ringShard3Port = "6392"
)

const (
	sentinelName       = "mymaster"
	sentinelMasterPort = "9123"
	sentinelSlave1Port = "9124"
	sentinelSlave2Port = "9125"
	sentinelPort       = "9126"
)

var (
	redisMain                                                *redisProcess
	ringShard1, ringShard2, ringShard3                       *redisProcess
	sentinelMaster, sentinelSlave1, sentinelSlave2, sentinel *redisProcess
)

var cluster = &clusterScenario{
	ports:     []string{"8220", "8221", "8222", "8223", "8224", "8225"},
	nodeIDs:   make([]string, 6),
	processes: make(map[string]*redisProcess, 6),
	clients:   make(map[string]*redis.Client, 6),
}

var _ = BeforeSuite(func() {
	var err error

	redisMain, err = startRedis(redisPort)
	Expect(err).NotTo(HaveOccurred())

	ringShard1, err = startRedis(ringShard1Port)
	Expect(err).NotTo(HaveOccurred())

	ringShard2, err = startRedis(ringShard2Port)
	Expect(err).NotTo(HaveOccurred())

	ringShard3, err = startRedis(ringShard3Port)
	Expect(err).NotTo(HaveOccurred())

	sentinelMaster, err = startRedis(sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())

	sentinel, err = startSentinel(sentinelPort, sentinelName, sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())

	sentinelSlave1, err = startRedis(
		sentinelSlave1Port, "--slaveof", "127.0.0.1", sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())

	sentinelSlave2, err = startRedis(
		sentinelSlave2Port, "--slaveof", "127.0.0.1", sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())

	Expect(startCluster(cluster)).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(redisMain.Close()).NotTo(HaveOccurred())

	Expect(ringShard1.Close()).NotTo(HaveOccurred())
	Expect(ringShard2.Close()).NotTo(HaveOccurred())
	Expect(ringShard3.Close()).NotTo(HaveOccurred())

	Expect(sentinel.Close()).NotTo(HaveOccurred())
	Expect(sentinelSlave1.Close()).NotTo(HaveOccurred())
	Expect(sentinelSlave2.Close()).NotTo(HaveOccurred())
	Expect(sentinelMaster.Close()).NotTo(HaveOccurred())

	Expect(stopCluster(cluster)).NotTo(HaveOccurred())
})

func TestGinkgoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-redis")
}

//------------------------------------------------------------------------------

func redisOptions() *redis.Options {
	return &redis.Options{
		Addr:               redisAddr,
		DB:                 15,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func redisClusterOptions() *redis.ClusterOptions {
	return &redis.ClusterOptions{
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func redisRingOptions() *redis.RingOptions {
	return &redis.RingOptions{
		Addrs: map[string]string{
			"ringShardOne": ":" + ringShard1Port,
			"ringShardTwo": ":" + ringShard2Port,
		},
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func performAsync(n int, cbs ...func(int)) *sync.WaitGroup {
	var wg sync.WaitGroup
	for _, cb := range cbs {
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(cb func(int), i int) {
				defer GinkgoRecover()
				defer wg.Done()

				cb(i)
			}(cb, i)
		}
	}
	return &wg
}

func perform(n int, cbs ...func(int)) {
	wg := performAsync(n, cbs...)
	wg.Wait()
}

func eventually(fn func() error, timeout time.Duration) error {
	errCh := make(chan error, 1)
	done := make(chan struct{})
	exit := make(chan struct{})

	go func() {
		for {
			err := fn()
			if err == nil {
				close(done)
				return
			}

			select {
			case errCh <- err:
			default:
			}

			select {
			case <-exit:
				return
			case <-time.After(timeout / 100):
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		close(exit)
		select {
		case err := <-errCh:
			return err
		default:
			return fmt.Errorf("timeout after %s without an error", timeout)
		}
	}
}

func execCmd(name string, args ...string) (*os.Process, error) {
	cmd := exec.Command(name, args...)
	if testing.Verbose() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Process, cmd.Start()
}

func connectTo(port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: ":" + port,
	})

	err := eventually(func() error {
		return client.Ping().Err()
	}, 30*time.Second)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type redisProcess struct {
	*os.Process
	*redis.Client
}

func (p *redisProcess) Close() error {
	if err := p.Kill(); err != nil {
		return err
	}

	err := eventually(func() error {
		if err := p.Client.Ping().Err(); err != nil {
			return nil
		}
		return errors.New("client is not shutdown")
	}, 10*time.Second)
	if err != nil {
		return err
	}

	p.Client.Close()
	return nil
}

var (
	redisServerBin, _  = filepath.Abs(filepath.Join("testdata", "redis", "src", "redis-server"))
	redisServerConf, _ = filepath.Abs(filepath.Join("testdata", "redis", "redis.conf"))
)

func redisDir(port string) (string, error) {
	dir, err := filepath.Abs(filepath.Join("testdata", "instances", port))
	if err != nil {
		return "", err
	}
	if err := os.RemoveAll(dir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0775); err != nil {
		return "", err
	}
	return dir, nil
}

func startRedis(port string, args ...string) (*redisProcess, error) {
	dir, err := redisDir(port)
	if err != nil {
		return nil, err
	}
	if err = exec.Command("cp", "-f", redisServerConf, dir).Run(); err != nil {
		return nil, err
	}

	baseArgs := []string{filepath.Join(dir, "redis.conf"), "--port", port, "--dir", dir}
	process, err := execCmd(redisServerBin, append(baseArgs, args...)...)
	if err != nil {
		return nil, err
	}

	client, err := connectTo(port)
	if err != nil {
		process.Kill()
		return nil, err
	}
	return &redisProcess{process, client}, err
}

func startSentinel(port, masterName, masterPort string) (*redisProcess, error) {
	dir, err := redisDir(port)
	if err != nil {
		return nil, err
	}
	process, err := execCmd(redisServerBin, os.DevNull, "--sentinel", "--port", port, "--dir", dir)
	if err != nil {
		return nil, err
	}
	client, err := connectTo(port)
	if err != nil {
		process.Kill()
		return nil, err
	}
	for _, cmd := range []*redis.StatusCmd{
		redis.NewStatusCmd("SENTINEL", "MONITOR", masterName, "127.0.0.1", masterPort, "1"),
		redis.NewStatusCmd("SENTINEL", "SET", masterName, "down-after-milliseconds", "500"),
		redis.NewStatusCmd("SENTINEL", "SET", masterName, "failover-timeout", "1000"),
		redis.NewStatusCmd("SENTINEL", "SET", masterName, "parallel-syncs", "1"),
	} {
		client.Process(cmd)
		if err := cmd.Err(); err != nil {
			process.Kill()
			return nil, err
		}
	}
	return &redisProcess{process, client}, nil
}

//------------------------------------------------------------------------------

type badConnError string

func (e badConnError) Error() string   { return string(e) }
func (e badConnError) Timeout() bool   { return false }
func (e badConnError) Temporary() bool { return false }

type badConn struct {
	net.TCPConn

	readDelay, writeDelay time.Duration
	readErr, writeErr     error
}

var _ net.Conn = &badConn{}

func (cn *badConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (cn *badConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (cn *badConn) Read([]byte) (int, error) {
	if cn.readDelay != 0 {
		time.Sleep(cn.readDelay)
	}
	if cn.readErr != nil {
		return 0, cn.readErr
	}
	return 0, badConnError("bad connection")
}

func (cn *badConn) Write([]byte) (int, error) {
	if cn.writeDelay != 0 {
		time.Sleep(cn.writeDelay)
	}
	if cn.writeErr != nil {
		return 0, cn.writeErr
	}
	return 0, badConnError("bad connection")
}

//------------------------------------------------------------------------------

type hook struct {
	beforeProcess func(ctx context.Context, cmd redis.Cmder) (context.Context, error)
	afterProcess  func(ctx context.Context, cmd redis.Cmder) error

	beforeProcessPipeline func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []redis.Cmder) error
}

func (h *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if h.beforeProcess != nil {
		return h.beforeProcess(ctx, cmd)
	}
	return ctx, nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if h.afterProcess != nil {
		return h.afterProcess(ctx, cmd)
	}
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if h.beforeProcessPipeline != nil {
		return h.beforeProcessPipeline(ctx, cmds)
	}
	return ctx, nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if h.afterProcessPipeline != nil {
		return h.afterProcessPipeline(ctx, cmds)
	}
	return nil
}
