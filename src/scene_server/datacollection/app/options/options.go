package options

import (
	"configcenter/src/common/core/cc/config"
	"configcenter/src/storage/mgoclient"
	"configcenter/src/storage/redisclient"

	"github.com/spf13/pflag"
)

type ServerOption struct {
	ServConf *config.CCAPIConfig
}

func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "127.0.0.1:50006", "The ip address and port for the serve on")

	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	fs.StringVar(&s.ServConf.ExConfig, "config", "", "The config path. e.g conf/api.conf")
}

type Config struct {
	MongoDB       mgoclient.MongoConfig
	CCRedis       redisclient.RedisConfig
	SnapRedis     redisclient.RedisConfig
	DiscoverRedis redisclient.RedisConfig
}
