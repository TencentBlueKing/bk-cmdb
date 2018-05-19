package api

import (
	"configcenter/src/framework/core/output/module/inst"
)

// PlatIteratorWrapper the plat iterator wrapper
type PlatIteratorWrapper struct {
	plat inst.Iterator
}

// Next next the plat
func (cli *PlatIteratorWrapper) Next() (*PlatWrapper, error) {

	plat, err := cli.plat.Next()

	return &PlatWrapper{plat: plat}, err

}

// ForEach the foreach function
func (cli *PlatIteratorWrapper) ForEach(callback func(plat *PlatWrapper) error) error {

	return cli.plat.ForEach(func(item inst.Inst) error {
		return callback(&PlatWrapper{plat: item})
	})
}

// PlatWrapper the plat wrapper
type PlatWrapper struct {
	plat inst.Inst
}

// SetValue set the key value
func (cli *PlatWrapper) SetValue(key string, val interface{}) error {
	return cli.plat.SetValue(key, val)
}

// Save save the data
func (cli *PlatWrapper) Save() error {
	return cli.plat.Save()
}

// SetSupplierAccount set the supplier account code of the set
func (cli *PlatWrapper) SetSupplierAccount(supplierAccount string) error {
	return cli.plat.SetValue(fieldSupplierAccount, supplierAccount)
}

// GetSupplierAccount get the supplier account
func (cli *PlatWrapper) GetSupplierAccount() (string, error) {
	vals, err := cli.plat.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldSupplierAccount), nil
}

// GetID get the set id
func (cli *PlatWrapper) GetID() (int, error) {
	vals, err := cli.plat.GetValues()
	if nil != err {
		return 0, err
	}
	return vals.Int(fieldPlatID)
}

// SetName the name of the set
func (cli *PlatWrapper) SetName(name string) error {
	return cli.plat.SetValue(fieldPlatName, name)
}

// GetName get the set name
func (cli *PlatWrapper) GetName() (string, error) {
	vals, err := cli.plat.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldPlatName), nil
}
