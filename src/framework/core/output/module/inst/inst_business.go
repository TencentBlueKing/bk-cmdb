package inst

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ Inst = (*business)(nil)

type business struct {
	target model.Model
	datas  types.MapStr
}

func (cli *business) GetModel() model.Model {
	return cli.target
}

func (cli *business) IsMainLine() bool {
	// TODO：判断当前实例是否为主线实例
	return true
}

func (cli *business) GetAssociationModels() ([]model.Model, error) {
	// TODO:需要读取此实例关联的实例，所对应的所有模型
	return nil, nil
}

func (cli *business) GetInstID() int {
	return 0
}
func (cli *business) GetInstName() string {
	return ""
}

func (cli *business) GetValues() (types.MapStr, error) {
	return nil, nil
}

func (cli *business) GetAssociationsByModleID(modleID string) ([]Inst, error) {
	// TODO:获取当前实例所关联的特定模型的所有已关联的实例
	return nil, nil
}

func (cli *business) GetAllAssociations() (map[model.Model][]Inst, error) {
	// TODO:获取所有已关联的模型及对应的实例
	return nil, nil
}

func (cli *business) SetParent(parentInstID int) error {
	return nil
}

func (cli *business) GetParent() ([]Topo, error) {
	return nil, nil
}

func (cli *business) GetChildren() ([]Topo, error) {
	return nil, nil
}

func (cli *business) SetValue(key string, value interface{}) error {

	// TODO:需要根据model 的定义对输入的key 及value 进行校验

	cli.datas[key] = value

	return nil
}

func (cli *business) Save() error {
	return nil
}
