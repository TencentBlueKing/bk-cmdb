package model

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	ClassificationID   string   `json:"bk_classification_id"`
	ClassificationName string   `json:"bk_classification_name"`
	ClassificationType string   `json:"bk_classification_type"`
	ClassificationIcon string   `json:"bk_classification_icon"`
	models             []*Model // all the models of this classification
}

func (cli *classification) GetID() string {
	return cli.ClassificationID
}

func (cli *classification) SetID(id string) {
	cli.ClassificationID = id
}

func (cli *classification) SetName(name string) {
	cli.ClassificationName = name
}

func (cli *classification) SetIcon(iconName string) {
	cli.ClassificationIcon = iconName
}

func (cli *classification) CreateModel() Model {
	m := &model{}
	m.ObjCls = cli.ClassificationID
	m.ObjIcon = cli.ClassificationIcon
	return m
}

func (cli *classification) FindModels(modelName string) (Iterator, error) {
	// TODO: 返回已经包含一定数量的Model数据的迭代器。
	return nil, nil
}
