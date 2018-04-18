package main

import (
	"configcenter/src/framework/api"
	"context"

	_ "configcenter/src/framework/plugins" // load all plugins
)

func main() {

	// init the framework

	for {
		select {
		case <-context.TODO().Done():
			api.UnInit()
			break
		}
	}

	// unint the framework
	//api.UnInit()
}
