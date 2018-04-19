package metric

import (
	"configcenter/src/common/types"
	"errors"
	"fmt"
)

type MetricFamily struct {
	MetaData     *MetaData                   `json:"metaData"`
	MetricBundle map[CollectorName][]*Metric `json:"metricBundle"`
	ReportTimeMs int64                       `json:"reportTimeMs"`
}

func (m MetaData) Valid() error {
	var errs []error
	if len(m.Module) == 0 {
		errs = append(errs, errors.New("module is null."))
	}

	if len(m.IP) == 0 {
		errs = append(errs, errors.New("IPAddr is null"))
	}

	if len(errs) != 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

type Metric struct {
	MetricMeta `json:",inline"`
	Value      *FloatOrString   `json:"value"`
	Extension  *MetricExtension `json:"extension"`
}

func newMetric(m MetricInterf) (*Metric, error) {
	if m == nil {
		return nil, errors.New("metric is nil.")
	}
	meta := m.GetMeta()
	if len(meta.Name) == 0 {
		return nil, errors.New("metric name is null.")
	}

	if len(meta.Help) == 0 {
		return nil, errors.New("metric help is null")
	}

	val, err := m.GetValue()
	if nil != err {
		return nil, err
	}
	if nil == val {
		return nil, errors.New("metric value is nil")
	}

	extension, err := m.GetExtension()
	if nil != err {
		return nil, err
	}
	return &Metric{
		MetricMeta: meta,
		Value:      val,
		Extension:  extension,
	}, nil
}

type CollectorName string
type Collector struct {
	Name      CollectorName
	Collector CollectInter
}

type MetaData struct {
	Module     string            `json:"module"`
	IP         string            `json:"ip"`
	MetricPort uint              `json:"metricPort"`
	ClusterID  string            `json:"clusterID"`
	Labels     map[string]string `json:"label"`
}

type HealthResponse struct {
	Code    int        `json:"code"`
	OK      bool       `json:"ok"`
	Message string     `json:"message"`
	Data    HealthInfo `json:"data"`
	Result  bool
}

type HealthInfo struct {
	Module     string `json:"module"`
	Address    string `json:"address"`
	HealthMeta `json:",inline"`
	AtTime     types.Time `json:"at_time"`
}
