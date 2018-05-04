package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client/v3"
	"configcenter/src/framework/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchHost(t *testing.T) {
	cli := v3.New(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition()
	rets, err := cli.Host().SearchHost(cond)
	assert.NoError(t, err)
	assert.NotEmpty(t, rets)
}

func TestDeleteHost(t *testing.T) {
	cli := v3.New(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_host_id").Eq("1")
	err := cli.Host().DeleteHost(cond)
	assert.NoError(t, err)
}
func TestUpdateHost(t *testing.T) {
	cli := v3.New(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_host_id").Eq("5")
	data := types.MapStr{"bk_host_name": "test_update"}
	err := cli.Host().UpdateHost(data, cond)
	assert.NoError(t, err)
}
