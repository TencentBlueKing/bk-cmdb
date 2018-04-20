package model

var _ Inst = (*inst)(nil)

// FieldName the field name
type FieldName string
type inst struct {
	target  Model
	storage map[FieldName]interface{}
}

func (cli *inst) SetValue(key string, value interface{}) {

}

func (cli *inst) Save() error {
	return nil
}
