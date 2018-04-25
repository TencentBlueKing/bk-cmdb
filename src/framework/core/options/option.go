package option

import (
	"github.com/spf13/pflag"
)

// AddFlags add some flags
func (cli *Options) AddFlags(fs *pflag.FlagSet) {
	//fs.BoolVar(&cli.Stop, "stop", true, "need to stop")
}
