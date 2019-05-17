package plugin

import (
	"sync/atomic"

	"configcenter/src/common/metric"
)

//integer
func NewIntegerCounter(name, help string) *IntegerCounter {
	return &IntegerCounter{
		name:  name,
		help:  help,
		value: 0,
	}
}

//IntegerCounter counter for integer
type IntegerCounter struct {
	name  string
	help  string
	value int64
}

func (c *IntegerCounter) GetMeta() metric.MetricMeta {
	return metric.MetricMeta{
		Name: c.name,
		Help: c.help,
	}
}

func (c *IntegerCounter) GetValue() (*metric.FloatOrString, error) {
	return metric.FormFloatOrString(c.value)
}

func (c *IntegerCounter) GetExtension() (*metric.MetricExtension, error) {
	return nil, nil
}

func (c *IntegerCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

func (c *IntegerCounter) Inc(i int64) {
	atomic.AddInt64(&c.value, i)
}

func (c *IntegerCounter) Dec(i int64) {
	atomic.AddInt64(&c.value, -i)
}
