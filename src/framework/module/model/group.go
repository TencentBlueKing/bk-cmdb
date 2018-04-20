package model

var _ Group = (*group)(nil)

type group struct {
	GroupID    string       `json:"bk_group_id"`
	GroupName  string       `json:"bk_group_name"`
	GroupIndex int          `json:"bk_group_index"`
	ObjectID   string       `json:"bk_obj_id"`
	OwnerID    string       `json:"bk_supplier_account"`
	IsDefault  bool         `json:"bk_isdefault"`
	IsPre      bool         `json:"ispre"`
	attrs      []*attribute // all the attributes of this group for a model
}

func (cli *group) SetID(id string) {
	cli.GroupID = id
}

func (cli *group) GetID() string {

	return cli.GroupID
}

func (cli *group) SetName(name string) {
	cli.GroupName = name
}

func (cli *group) SetIndex(idx int) {
	cli.GroupIndex = idx
}

func (cli *group) GetIndex() int {
	return cli.GroupIndex
}

func (cli *group) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}

func (cli *group) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *group) SetDefault() {
	cli.IsDefault = true
}
func (cli *group) SetNonDefault() {
	cli.IsDefault = false
}

func (cli *group) Default() bool {
	return cli.IsDefault
}

func (cli *group) CreateAttribute() Attribute {
	attr := &attribute{}
	return attr
}

func (cli *group) FindAttributes(attributeName string) (AttributeIterator, error) {
	// TODO: 返回已经包含一定数量的Atribute数据的迭代器。
	return nil, nil
}
