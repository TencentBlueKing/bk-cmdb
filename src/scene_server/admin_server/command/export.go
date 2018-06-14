package command

import (
	"configcenter/src/storage"
	"encoding/json"
	"fmt"
	"os"
)

func export(db storage.DI, opt *option) error {
	file, err := os.Create(opt.position)
	if nil != err {
		return err
	}
	defer file.Close()
	defer file.Sync()

	topo, err := getBKTopo(db, opt)
	if nil != err {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(topo)
	if nil != err {
		return fmt.Errorf("encode topo error: %s", err.Error())
	}

	return nil
}
