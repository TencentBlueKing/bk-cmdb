package input

import (
	"configcenter/src/framework/common"
)

// create a new inputer key
func makeInputerKey() InputerKey {
	return InputerKey(common.UUID())
}

// checkWorkerExists check whether the inputer exists
func inputerExists(target MapInputer, key InputerKey) bool {
	_, ok := target[key]
	return ok
}

// deleteInputer delete a inputer from MapInputer
func deleteInputer(target MapInputer, key InputerKey) bool {

	if inputerExists(target, key) {
		delete(target, key)
	}

	return true
}
