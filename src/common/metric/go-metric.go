package metric

import (
	"fmt"
	"runtime"
)

func newGoMetricCollector() *Collector {
	golang := &golang{
		goRoutineMetric: goMetric{
			Name:    "go_goroutines",
			Help:    "Number of goroutines that currently exist.",
			GetFunc: func() float64 { return float64(runtime.NumGoroutine()) },
		},
		goProcessMetric: goMetric{
			Name: "go_threads",
			Help: "Number of OS threads created",
			GetFunc: func() float64 {
				n, _ := runtime.ThreadCreateProfile(nil)
				return float64(n)
			},
		},
		goCPUMetric: goMetric{
			Name: "go_cpu_used",
			Help: " the number of logical CPUs usable by the current process.",
			GetFunc: func() float64 { return float64(runtime.NumCPU() )},
		},
		goMemStateMetrics: []goMetric{
			{
				Name:    memstatNamespace("alloc_bytes"),
				Help:    "Number of bytes allocated and still in use.",
				MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.Alloc)},
			},
			{
				Name:    memstatNamespace("alloc_bytes_total"),
				Help:    "Total number of bytes allocated, even if freed.",
				MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.TotalAlloc)},
			},
			{
				Name:    memstatNamespace("sys_bytes"),
				Help:    "Number of bytes obtained from system.",
				MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.Sys)},
			},
			{
				Name:    memstatNamespace("mallocs_total"),
				Help:    "Total number of mallocs.",
				MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.Mallocs) },
			},
			{
				Name:    memstatNamespace("frees_total"),
				Help:    "Total number of frees.",
				MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.Frees)},
			},
			//{
			//	Name:    memstatNamespace("lookups_total"),
			//	Help:    "Total number of pointer lookups.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.Lookups)},
			//},
			//{
			//	Name:    memstatNamespace("heap_alloc_bytes"),
			//	Help:    "Number of heap bytes allocated and still in use.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.HeapAlloc},
			//},
			//{
			//	Name:    memstatNamespace("heap_sys_bytes"),
			//	Help:    "Number of heap bytes obtained from system.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.HeapSys)},
			//},
			//{
			//	Name:    memstatNamespace("heap_idle_bytes"),
			//	Help:    "Number of heap bytes waiting to be used.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.HeapIdle)},
			//},
			//{
			//	Name:    memstatNamespace("heap_inuse_bytes"),
			//	Help:    "Number of heap bytes that are in use.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.HeapInuse) },
			//},
			//{
			//	Name:    memstatNamespace("heap_released_bytes"),
			//	Help:    "Number of heap bytes released to OS.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.HeapReleased},
			//},
			//{
			//	Name:    memstatNamespace("heap_objects"),
			//	Help:    "Number of allocated objects.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.HeapObjects},
			//},
			//{
			//	Name:    memstatNamespace("stack_inuse_bytes"),
			//	Help:    "Number of bytes in use by the stack allocator.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 { return float64(ms.StackInuse},
			//},
			//{
			//	Name:    memstatNamespace("stack_sys_bytes"),
			//	Help:    "Number of bytes obtained from system for stack allocator.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.StackSys},
			//},
			//{
			//	Name:    memstatNamespace("mspan_inuse_bytes"),
			//	Help:    "Number of bytes in use by mspan structures.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.MSpanInuse},
			//},
			//{
			//	Name:    memstatNamespace("mspan_sys_bytes"),
			//	Help:    "Number of bytes used for mspan structures obtained from system.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.MSpanSys) },
			//},
			//{
			//	Name:    memstatNamespace("mcache_inuse_bytes"),
			//	Help:    "Number of bytes in use by mcache structures.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 { return float64(ms.MCacheInuse) },
			//},
			//{
			//	Name:    memstatNamespace("mcache_sys_bytes"),
			//	Help:    "Number of bytes used for mcache structures obtained from system.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.MCacheSys) },
			//},
			//{
			//	Name:    memstatNamespace("buck_hash_sys_bytes"),
			//	Help:    "Number of bytes used by the profiling bucket hash table.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 { return float64(ms.BuckHashSys},
			//},
			//{
			//	Name:    memstatNamespace("gc_sys_bytes"),
			//	Help:    "Number of bytes used for garbage collection system metadata.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 { return float64(ms.GCSys) },
			//},
			//{
			//	Name:    memstatNamespace("other_sys_bytes"),
			//	Help:    "Number of bytes used for other system allocations.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.OtherSys)},
			//},
			//{
			//	Name:    memstatNamespace("next_gc_bytes"),
			//	Help:    "Number of heap bytes when next garbage collection will take place.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.NextGC) },
			//},
			//{
			//	Name:    memstatNamespace("last_gc_time_seconds"),
			//	Help:    "Number of seconds since 1970 of last garbage collection.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return float64(ms.LastGC) / 1e9},
			//},
			//{
			//	Name:    memstatNamespace("gc_cpu_fraction"),
			//	Help:    "The fraction of this program's available CPU time used by the GC since the program started.",
			//	MemGetFunc: func(ms *runtime.MemStats) float64 {return ms.GCCPUFraction },
			//},
		},
	}

	return NewCollector("golang_metrics", golang)
}

func memstatNamespace(s string) string {
	return fmt.Sprintf("go_memstats_%s", s)
}

type golang struct {
	goRoutineMetric   goMetric
	goProcessMetric   goMetric
	goCPUMetric       goMetric
	goMemStateMetrics []goMetric
}

func (g golang) Collect()[]MetricInterf {
	m := make([]MetricInterf, 0)
	m = append(m, &g.goRoutineMetric)
	m = append(m, &g.goProcessMetric)

	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	for idx := range g.goMemStateMetrics {
		g.goMemStateMetrics[idx].MemStats = ms
		m = append(m, &g.goMemStateMetrics[idx])
	}
	return m
}

type goMetric struct {
	Name       string
	Help       string
	MemStats   *runtime.MemStats
	MemGetFunc func(stat *runtime.MemStats) float64
	GetFunc    func() float64
}

func (m goMetric)GetMeta() MetricMeta {
	return MetricMeta{
		Name: m.Name,
		Help: m.Help,
	}
}

func (m goMetric) GetValue() (*FloatOrString, error) {
	if m.MemStats != nil {
		return FormFloatOrString(m.MemGetFunc(m.MemStats))
	}
	return FormFloatOrString(m.GetFunc())
}

func (m goMetric) GetExtension() (*MetricExtension, error) {
	return nil, nil
}
