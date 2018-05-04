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

	cond := common.CreateCondition().Field("bk_host_innerip").In([]string{"192.168.100.1"})
	rets, err := cli.Host().SearchHost(cond)
	t.Logf("search host result: %v", rets)
	assert.NoError(t, err)
	assert.NotEmpty(t, rets)
}

func TestDeleteHost(t *testing.T) {
	cli := v3.New(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	err := cli.Host().DeleteHostBatch("1")
	assert.NoError(t, err)
}
func TestUpdateHost(t *testing.T) {
	cli := v3.New(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	data := types.MapStr{"bk_host_name": "test_update"}
	err := cli.Host().UpdateHostBatch(data, "5")
	assert.NoError(t, err)
}
