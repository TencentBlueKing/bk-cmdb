package model

import "configcenter/src/framework/core/types"

var _ Model = (*model)(nil)

// model the metadata structure definition of the model
type model struct {
	ObjCls      string `json:"bk_classification_id"`
	ObjIcon     string `json:"bk_obj_icon"`
	ObjectID    string `json:"bk_obj_id"`
	ObjectName  string `json:"bk_obj_name"`
	IsPre       bool   `json:"ispre"`
	IsPaused    bool   `json:"bk_ispaused"`
	Position    string `json:"position"`
	OwnerID     string `json:"bk_supplier_account"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Modifier    string `json:"modifier"`
}

func (cli *model) Save() error {
	return nil
}

func (cli *model) CreateAttribute() Attribute {
	attr := &attribute{}
	return attr
}

func (cli *model) SetClassification(class Classification) {
	cli.ObjCls = class.GetID()
}

func (cli *model) SetIcon(iconName string) {
	cli.ObjIcon = iconName
}

func (cli *model) SetID(id string) {
	cli.ObjectID = id
}

func (cli *model) SetName(name string) {
	cli.ObjectName = name
}

func (cli *model) SetPaused(isPaused bool) {
	cli.IsPaused = isPaused
}

func (cli *model) SetPosition(position string) {
	cli.Position = position
}

func (cli *model) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}

func (cli *model) SetDescription(desc string) {
	cli.Description = desc
}

func (cli *model) SetCreator(creator string) {
	cli.Creator = creator
}

func (cli *model) SetModifier(modifier string) {
	cli.Modifier = modifier
}

func (cli *model) CreateGroup() Group {
	g := &group{}
	return g
}

func (cli *model) FindAttributesLikeName(attributeName string) (AttributeIterator, error) {
	// TODO:按照名字正则查找
	return nil, nil
}
func (cli *model) FindAttributesByCondition(condition types.MapStr) (AttributeIterator, error) {
	// TODO:按照条件查找
	return nil, nil
}
func (cli *model) FindGroupsLikeName(groupName string) (GroupIterator, error) {
	// TODO:按照名字正则查找
	return nil, nil
}
func (cli *model) FindGroupsByCondition(condition types.MapStr) (GroupIterator, error) {
	// TODO:按照条件查找
	return nil, nil
}
func (cli *model) CreateInst() Inst {
	// TODO: 基于当前模型创建的实例，实例对象需要集成校验模式，能够知道自己包含的字段，以及字段的属性。
	tmp := &inst{}
	return tmp
}
