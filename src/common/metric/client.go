package metric

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/ssl"
)

var metricController *MetricController

func newMetricController(conf Config, healthFunc HealthFunc, collectors ...*Collector) error {
	metricController = new(MetricController)
	meta := MetaData{
		Module:     conf.ModuleName,
		IP:         conf.IP,
		MetricPort: conf.MetricPort,
		ClusterID:  conf.ClusterID,
		Labels:     conf.Labels,
	}
	if err := meta.Valid(); nil != err {
		return err
	}

	// set default golang metric.
	collectors = append(collectors, newGoMetricCollector())

	metricController.MetaData = &meta
	metricController.Collectors = make(map[CollectorName]CollectInter)
	for _, c := range collectors {
		if _, exist := metricController.Collectors[c.Name]; false == exist {
			metricController.Collectors[c.Name] = c.Collector
		} else {
			return fmt.Errorf("duplicate collector name: %s", c.Name)
		}
	}

	mux := http.NewServeMux()
	metricHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metric, err := metricController.PackMetrics()
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			info := fmt.Sprintf("get metrics failed. err: %v", err)
			w.Write([]byte(info))
			return
		}
		w.Write(*metric)
	})

	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h := healthFunc()
		info := HealthInfo{
			RunMode:    conf.RunMode,
			Module:     conf.ModuleName,
			ClusterID:  conf.ClusterID,
			IP:         conf.IP,
			HealthMeta: h,
			AtTime:     time.Now().Unix(),
		}
		js, err := json.Marshal(info)
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			info := fmt.Sprintf("get health info failed. err: %v", err)
			w.Write([]byte(info))
			return
		}
		w.Write(js)
	})

	mux.Handle("/metrics", metricHandler)
	mux.Handle("/healthz", healthHandler)

	if err := listenAndServe(conf, mux); err != nil {
		return fmt.Errorf("listen and serve failed, err: %v", err)
	}
	return nil
}

func listenAndServe(c Config, mux http.Handler) error {
	addr := fmt.Sprintf("%s:%d", c.IP, c.MetricPort)

	if c.SvrCertFile == "" && c.SvrKeyFile == "" {
		go func() {
			blog.Infof("started metric and listen insecure server on %s", addr)
			blog.Fatal(http.ListenAndServe(addr, mux))
		}()
		return nil
	}

	// user https
	ca, err := ioutil.ReadFile(c.SvrCaFile)
	if nil != err {
		return err
	}
	capool := x509.NewCertPool()
	capool.AppendCertsFromPEM(ca)
	tlsconfig, err := ssl.ServerTslConfVerityClient(c.SvrCaFile,
		c.SvrCertFile,
		c.SvrKeyFile,
		c.CertPasswd)
	if err != nil {
		return err
	}
	tlsconfig.BuildNameToCertificate()

	blog.Info("start metric secure serve on %s", addr)

	ln, err := net.Listen("tcp", net.JoinHostPort(c.IP, strconv.FormatUint(uint64(c.MetricPort), 10)))
	if err != nil {
		return err
	}
	listener := tls.NewListener(ln, tlsconfig)
	go func() {
		if err := http.Serve(listener, mux); nil != err {
			blog.Fatalf("server https failed. err: %v", err)
		}
	}()
	return nil
}

type MetricController struct {
	MetaData   *MetaData
	Collectors map[CollectorName]CollectInter
}

func (mc *MetricController) PackMetrics() (*[]byte, error) {
	mf := MetricFamily{
		MetaData:     mc.MetaData,
		MetricBundle: make(map[CollectorName][]*Metric),
	}

	for name, collector := range mc.Collectors {
		mf.MetricBundle[name] = make([]*Metric, 0)
		done := make(chan struct{}, 0)
		go func(c CollectInter) {
			for _, mc := range c.Collect() {
				metric, err := newMetric(mc)
				if nil != err {
					blog.Errorf("new metric failed. err: %v", err)
					continue
				}
				mf.MetricBundle[name] = append(mf.MetricBundle[name], metric)
			}
			done <- struct{}{}
		}(collector)

		select {
		case <-time.After(time.Duration(10 * time.Second)):
			blog.Errorf("get metric bundle: %s timeout, skip it.", name)
			continue
		case <-done:
			close(done)
		}
	}

	mf.ReportTimeMs = time.Now().Unix()
	js, err := json.Marshal(mf)
	if nil != err {
		return nil, err
	}
	return &js, nil
}
