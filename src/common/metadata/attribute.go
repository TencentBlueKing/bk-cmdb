package metadata

import (
	types "configcenter/src/common/mapstr"
)

// Attribute attribute metadata definition
type Attribute struct {
	OwnerID       string      `field:"bk_supplier_account"`
	ObjectID      string      `field:"bk_obj_id"`
	PropertyID    string      `field:"bk_property_id"`
	PropertyName  string      `field:"bk_property_name"`
	PropertyGroup string      `field:"bk_property_group"`
	PropertyIndex int         `field:"bk_property_index"`
	Unit          string      `field:"unit"`
	Placeholder   string      `field:"placeholder"`
	IsEditable    bool        `field:"editable"`
	IsPre         bool        `field:"ispre"`
	IsRequired    bool        `field:"isrequired"`
	IsReadOnly    bool        `field:"isreadonly"`
	IsOnly        bool        `field:"isonly"`
	IsSystem      bool        `field:"bk_issystem"`
	IsAPI         bool        `field:"bk_isapi"`
	PropertyType  string      `field:"bk_property_type"`
	Option        interface{} `field:"option"`
	Description   string      `field:"description"`
	Creator       string      `field:"creator"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *Attribute) Parse(data types.MapStr) (*Attribute, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Attribute) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
