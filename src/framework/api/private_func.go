package api

import (
	"configcenter/src/framework/common"
)

// create a new worker key
func makeWorkerKey() WorkerKey {
	return WorkerKey(common.UUID())
}

// checkWorkerExists check whether the worker exists
func workerExists(target MapWorker, key WorkerKey) bool {
	_, ok := target[key]
	return ok
}

// deleteWorker delete a worker from MapWorker
func deleteWorker(target MapWorker, key WorkerKey) bool {

	if workerExists(target, key) {
		delete(target, key)
	}

	return true
}
