package metric

func NewMetricController(conf Config, healthFunc HealthFunc, collectors ...*Collector) error {
	return newMetricController(conf, healthFunc, collectors...)
}

type RunModeType string

// used when your module running with Master_Slave_Mode mode
type RoleType string

// metric const define
const (
	MetricPort = 60060
)

// Config define metric's define
type Config struct {
	// name of your module
	ModuleName string
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

// HealthFunc returns HealthMeta
type HealthFunc func() HealthMeta

// HealthMeta define the HealthMeta that shows whether this server healthy
type HealthMeta struct {
	// if this module is healthy
	IsHealthy bool `json:"healthy"`
	// messages which describes the health status
	Message string `json:"message"`

	Items []HealthItem `json:"items"`
}

type HealthItem struct {
	Name string `json:"name"`

	HealthMeta `json:",inline"`
}

// MetricMeta define the MetricMeta that shows the named metric
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
