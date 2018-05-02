package v3

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
	// "encoding/json"
	// "errors"
	//"fmt"
	//"github.com/tidwall/gjson"
)

// CreateGroup create a group
func CreateGroup(data types.MapStr) (int, error) {
	return 0, nil
}

// DeleteGroup delete a group by condition
func (cli *Client) DeleteGroup(cond common.Condition) error {
	return nil
}

// UpdateGroup update a group by condition
func (cli *Client) UpdateGroup(data types.MapStr, cond common.Condition) error {
	return nil
}

// SearchGroups search some group by condition
func (cli *Client) SearchGroups(cond common.Condition) ([]types.MapStr, error) {
	return nil, nil
}
