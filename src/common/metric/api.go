package metric

func NewMetricController(conf Config, healthFunc HealthFunc, collectors ...*Collector) error {
	return newMetricController(conf, healthFunc, collectors...)
}

type RunModeType string

// used when your module running with Master_Slave_Mode mode
type RoleType string

const (
	Master_Slave_Mode  RunModeType = "master-slave"
	Master_Master_Mode RunModeType = "master-master"
	MasterRole         RoleType    = "master"
	SlaveRole          RoleType    = "slave"
	UnknownRole        RoleType    = "unknown"
)

type Config struct {
	// name of your module
	ModuleName string
	// running mode of your module
	// could be one of Master_Slave_Mode or Master_Master_Mode
	RunMode RunModeType
	// ip address of this module running on
	IP string
	// port number of the metric's http handler depends on.
	MetricPort uint
	// cluster id of your module belongs to.
	ClusterID string
	// self defined info labeled on your metrics.
	Labels map[string]string
	// metric http server's ssl configuration
	SvrCaFile   string
	SvrCertFile string
	SvrKeyFile  string
	CertPasswd  string
}

type HealthFunc func() HealthMeta

type HealthMeta struct {
	// the running role of your module when you are running with Master_Slave_Mode.
	// must be not empty. if you set with an empty value, an error will be occurred.
	// when your module is running in Master_Master_Mode,  this filed should be set
	// with value of "Slave".
	CurrentRole RoleType `json:"current_role"`
	// if this module is healthy
	IsHealthy bool `json:"healthy"`
	// messages which describes the health status
	Message string `json:"message"`
}

type MetricMeta struct {
	// metric's name
	Name string `json:"name"`
	// metric's help info, which should be short and briefly.
	Help string `json:"help"`
}

type MetricInterf interface {
	GetMeta() MetricMeta
	GetValue() (*FloatOrString, error)
	GetExtension() (*MetricExtension, error)
}

type MetricExtension struct{}

type CollectInter interface {
	Collect() []MetricInterf
}

func NewCollector(name string, collector CollectInter) *Collector {
	return &Collector{
		Name:      CollectorName(name),
		Collector: collector,
	}
}
