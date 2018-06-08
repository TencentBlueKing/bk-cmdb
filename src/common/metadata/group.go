package metadata

import (
	types "configcenter/src/common/mapstr"
)

// Group group metadata definition
type Group struct {
	GroupID    string `field:"bk_group_id"`
	GroupName  string `field:"bk_group_name"`
	GroupIndex int    `field:"bk_group_index"`
	ObjectID   string `field:"bk_obj_id"`
	OwnerID    string `field:"bk_supplier_account"`
	IsDefault  bool   `field:"bk_isdefault"`
	IsPre      bool   `field:"ispre"`
}

// Parse load the data from mapstr group into group instance
func (cli *Group) Parse(data types.MapStr) (*Group, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Group) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
