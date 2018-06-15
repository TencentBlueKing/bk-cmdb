package logics

type Metric struct {
	MetricInterface
}

type MetricInterface interface {
	HandleMsg()
	Parse()
	Update()
}

