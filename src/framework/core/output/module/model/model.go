package model

import (
	"configcenter/src/framework/common"
	"fmt"
)

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
	fmt.Println("test model")
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

func (cli *model) GetIcon() string {
	return cli.ObjIcon
}

func (cli *model) SetID(id string) {
	cli.ObjectID = id
}

func (cli *model) GetID() string {
	return cli.ObjectID
}

func (cli *model) SetName(name string) {
	cli.ObjectName = name
}
func (cli *model) GetName() string {
	return cli.ObjectName
}

func (cli *model) SetPaused() {
	cli.IsPaused = true
}

func (cli *model) SetNonPaused() {
	cli.IsPaused = false
}

func (cli *model) Paused() bool {
	return cli.IsPaused
}

func (cli *model) SetPosition(position string) {
	cli.Position = position
}

func (cli *model) GetPosition() string {
	return cli.Position
}

func (cli *model) SetSupplierAccount(ownerID string) {
	cli.OwnerID = ownerID
}
func (cli *model) GetSupplierAccount() string {
	return cli.OwnerID
}

func (cli *model) SetDescription(desc string) {
	cli.Description = desc
}
func (cli *model) GetDescription() string {
	return cli.Description
}
func (cli *model) SetCreator(creator string) {
	cli.Creator = creator
}
func (cli *model) GetCreator() string {
	return cli.Creator
}
func (cli *model) SetModifier(modifier string) {
	cli.Modifier = modifier
}
func (cli *model) GetModifier() string {
	return cli.Modifier
}
func (cli *model) CreateGroup() Group {
	g := &group{}
	return g
}

func (cli *model) FindAttributesLikeName(attributeName string) (AttributeIterator, error) {
	// TODO:按照名字正则查找
	return nil, nil
}
func (cli *model) FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error) {
	// TODO:按照条件查找
	return nil, nil
}
func (cli *model) FindGroupsLikeName(groupName string) (GroupIterator, error) {
	// TODO:按照名字正则查找
	return nil, nil
}
func (cli *model) FindGroupsByCondition(condition *common.Condition) (GroupIterator, error) {
	// TODO:按照条件查找
	return nil, nil
}
