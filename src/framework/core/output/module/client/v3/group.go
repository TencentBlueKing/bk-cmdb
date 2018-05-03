package v3

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
	//"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

// CreateGroup create a group
func (cli *Client) CreateGroup(data types.MapStr) (int, error) {

	targetURL := fmt.Sprintf("%s/api/v3/objectattr/group/new", cli.GetAddress())

	rst, err := cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return 0, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return 0, errors.New(gs.Get("bk_error_msg").String())
	}

	// parse id
	id := gs.Get("data.id").Int()

	return int(id), nil
}

// DeleteGroup delete a group by condition
func (cli *Client) DeleteGroup(cond common.Condition) error {

	data := cond.ToMapStr()
	id, err := data.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/objectattr/group/groupid/%d", cli.GetAddress(), id)

	rst, err := cli.httpCli.DELETE(targetURL, nil, nil)
	if nil != err {
		return err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return errors.New(gs.Get("bk_error_msg").String())
	}

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
