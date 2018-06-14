package command

import (
	"configcenter/src/storage"
	"fmt"
	"os"
)

func export(position string, db storage.DI) error {
	stat, err := os.Stat(position)
	if nil != err {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("%s is not a file", position)
	}

	file, err := os.Create(position)
	if nil != err {
		return err
	}

	// db.GetMutilByCondition(cName, fields, condiction, result, sort, start, limit)

	return nil

}
