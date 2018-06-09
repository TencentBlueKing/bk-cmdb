package metadata

import (
	types "configcenter/src/common/mapstr"
)

// Classification the classification metadata definition
type Classification struct {
	ClassificationID   string `field:"bk_classification_id"`
	ClassificationName string `field:"bk_classification_name"`
	ClassificationType string `field:"bk_classification_type"`
	ClassificationIcon string `field:"bk_classification_icon"`
}

// Parse load the data from mapstr classification into classification instance
func (cli *Classification) Parse(data types.MapStr) (*Classification, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Classification) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
