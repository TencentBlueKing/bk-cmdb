package main

// 作为mongo故障检测用例，用于mongo切主、增加从节点等DB变动时，观察mongo操作的情况和db故障情况

import (
	"sync"

	"configcenter/src/common/blog"
	"configcenter/src/test/db/mongo/operator"
)

type metrics struct {
	writeCntTotal   int
	writeCntSuccess int
	writeCntFail    int
	readCntTotal    int
	readCntSuccess  int
	readCntFail     int
}

func main() {

	operator := operator.NewMongoOperator("cc_mongo_check")
	if err := operator.ClearData(); err != nil {
		blog.Errorf("ClearData failed, err:%s", err)
		return
	}

	metrics := metrics{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			err := operator.WriteWithTxn()
			metrics.writeCntTotal++
			if err != nil {
				metrics.writeCntFail++
				blog.Errorf("write failed, err:%s, metrics:%#v", err, metrics)
				continue
			}
			blog.Infof("write success, metrics:%#v", metrics)
			metrics.writeCntSuccess++
		}
	}()

	go func() {
		defer wg.Done()
		for {
			err := operator.ReadNoTxn()
			metrics.readCntTotal++
			if err != nil {
				metrics.readCntFail++
				blog.Errorf("read failed, err:%s, metrics:%#v", err, metrics)
				continue
			}
			metrics.readCntSuccess++
			blog.Infof("read success, metrics:%#v", metrics)
		}
	}()

	wg.Wait()

}
