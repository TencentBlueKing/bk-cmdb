package inst

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

type proc struct {
	target model.Model
	datas  types.MapStr
}

func (cli *proc) GetModel() model.Model {
	return cli.target
}

func (cli *proc) IsMainLine() bool {
	// TODO：判断当前实例是否为主线实例
	return true
}

func (cli *proc) GetAssociationModels() ([]model.Model, error) {
	// TODO:需要读取此实例关联的实例，所对应的所有模型
	return nil, nil
}

func (cli *proc) GetInstID() int {
	return 0
}
func (cli *proc) GetInstName() string {
	return ""
}

func (cli *proc) GetValues() (types.MapStr, error) {
	return nil, nil
}

func (cli *proc) GetAssociationsByModleID(modleID string) ([]Inst, error) {
	// TODO:获取当前实例所关联的特定模型的所有已关联的实例
	return nil, nil
}

func (cli *proc) GetAllAssociations() (map[model.Model][]Inst, error) {
	// TODO:获取所有已关联的模型及对应的实例
	return nil, nil
}

func (cli *proc) SetParent(parentInstID int) error {
	return nil
}

func (cli *proc) GetParent() ([]Topo, error) {
	return nil, nil
}

func (cli *proc) GetChildren() ([]Topo, error) {
	return nil, nil
}

func (cli *proc) SetValue(key string, value interface{}) error {

	// TODO:需要根据model 的定义对输入的key 及value 进行校验

	cli.datas[key] = value

	return nil
}

func (cli *proc) Save() error {
	return nil
}
