package metric

import (
	"configcenter/src/common/types"
	"errors"
	"net/http"
)

type MetricFamily struct {
	MetaData     *MetaData                   `json:"metaData"`
	MetricBundle map[CollectorName][]*Metric `json:"metricBundle"`
	ReportTimeMs int64                       `json:"reportTimeMs"`
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
	Module        string            `json:"module"`
	ServerAddress string            `json:"server_address"`
	ClusterID     string            `json:"clusterID"`
	Labels        map[string]string `json:"label"`
}

type HealthResponse struct {
	Code    int        `json:"code"`
	OK      bool       `json:"ok"`
	Message string     `json:"message"`
	Data    HealthInfo `json:"data"`
	Result  bool       `json:"result"`
}

type HealthInfo struct {
	Module     string `json:"module"`
	Address    string `json:"address"`
	HealthMeta `json:",inline"`
	AtTime     types.Time `json:"at_time"`
}

type Action struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}
