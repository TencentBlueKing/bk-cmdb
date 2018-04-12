package main

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/example/log"
	"configcenter/src/framework/example/sender"
	"configcenter/src/framework/example/worker"
	"fmt"
)

func main() {

	// create a new framework instance
	fr, err := api.CMDBFrameworkInit()

	if nil != err {
		fmt.Println("failed to create a framework:", err.Error())
		return
	}

	// TODO: reset the default implementation

	// replace the default loger
	fr.SetLoger(&log.CustomLog{})

	// register the custom output
	fr.Publisher.RegisterCustom(&sender.CustomSender{})

	// TODO: register the custom worker
	fr.AddWorker(&worker.WorkerOne{})
	fr.AddWorker(&worker.WorkerTwo{})

	// start the main loop, it will block here
	fr.Run()

	// destory the framework instance
	api.CMDBFrameworkUnInit(fr)
}
