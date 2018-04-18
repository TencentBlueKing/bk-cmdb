/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"configcenter/src/common"
	"fmt"
)

func (cli *Client) ReForwardSelectMetaInstsTopo(callfunc func(url, method string) (string, error), ownerid, objid, instid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/topo/v1/inst/search/topo/owner/%s/object/%s/inst/%s", cli.address, ownerid, objid, instid), common.HTTPSelectPost)
	}
}
func (cli *Client) ReForwardSelectMetaInsts(callfunc func(url, method string) (string, error), ownerid, objid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/topo/v1/inst/search/%s/%s", cli.address, ownerid, objid), common.HTTPSelectPost)
	}
}

func (cli *Client) ReForwardSelectMetaInst(callfunc func(url, method string) (string, error), ownerid, objid, instid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/topo/v1/inst/search/%s/%s/%s", cli.address, ownerid, objid, instid), common.HTTPSelectPost)
	}
}

func (cli *Client) ReForwardSelectInstByAssociation(callfunc func(url, method string) (string, error), ownerid, objid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/topo/v1/inst/association/search/owner/%s/object/%s", cli.address, ownerid, objid), common.HTTPSelectPost)
	}
}
