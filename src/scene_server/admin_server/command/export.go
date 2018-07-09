package command

import (
	"encoding/json"
	"fmt"
	"os"

	"configcenter/src/storage"
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
	encoder.SetIndent("", "    ")
	err = encoder.Encode(topo)
	if nil != err {
		return fmt.Errorf("encode topo error: %s", err.Error())
	}

	return nil
}
