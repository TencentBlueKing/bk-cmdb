package topo_server_test

import (
	"context"
	"encoding/json"
	"fmt"

	"configcenter/src/common/metadata"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("privilege test", func() {
	privilegeClient := topoServerClient.Privilege()

	It("search user group", func() {
		test.ClearDatabase()
		input := map[string]interface{}{
			"group_name": "",
			"user_list":  "",
		}
		rsp, err := privilegeClient.SearchUserGroup(context.Background(), "0", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("get role privilege bk_obj_id = 0 and bk_property_id = 'bk_biz_maintainer'", func() {
		rsp, err := privilegeClient.GetPrivilege(context.Background(), "0", "0", "bk_biz_maintainer", header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create user group", func() {
		input := map[string]interface{}{
			"group_name": "kkk",
			"user_list":  "admin",
		}
		rsp, err := privilegeClient.CreateUserGroup(context.Background(), "0", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	var groupId string
	It("search user group group_name = 'kkk' and user_list = 'admin'", func() {
		input := map[string]interface{}{
			"group_name": "kkk",
			"user_list":  "admin",
		}
		rsp, err := privilegeClient.SearchUserGroup(context.Background(), "0", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		j, err := json.Marshal(rsp.Data)
		data := []metadata.UserGroup{}
		json.Unmarshal(j, &data)
		Expect(data[0].GroupName).To(Equal("kkk"))
		Expect(data[0].UserList).To(Equal("admin"))
		Expect(data[0].SupplierAccount).To(Equal("0"))
		groupId = data[0].GroupID
	})

	It(fmt.Sprintf("search user group privilege group_id = %s", groupId), func() {
		rsp, err := privilegeClient.GetUserGroupPrivi(context.Background(), "0", groupId, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		j, err := json.Marshal(rsp.Data)
		data := metadata.GroupPrivilege{}
		json.Unmarshal(j, &data)
		Expect(data.GroupID).To(Equal(groupId))
		Expect(data.OwnerID).To(Equal("0"))
	})

	It(fmt.Sprintf("update user group group_id = %s", groupId), func() {
		input := map[string]interface{}{
			"group_name": "kkk2",
			"user_list":  "admin",
		}
		rsp, err := privilegeClient.UpdateUserGroup(context.Background(), "0", groupId, header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It(fmt.Sprintf("update user group privilege group_id = %s", groupId), func() {
		input := map[string]interface{}{
			"sys_config": map[string]interface{}{
				"global_busi": []string{
					"resource",
				},
			},
		}
		rsp, err := privilegeClient.UpdateUserGroupPrivi(context.Background(), "0", groupId, header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("bind role privilege bk_obj_id = 'biz' and bk_property_id = 'bk_biz_productor'", func() {
		data := map[string]interface{}{
			"origin": []string{
				"hosttrans",
			},
		}
		rsp, err := privilegeClient.CreatePrivilegeWithData(context.Background(), "0", "biz", "bk_biz_productor", header, data)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("get user privilege user_name = 'admin'", func() {
		rsp, err := privilegeClient.GetUserPrivi(context.Background(), "0", "admin", header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})
})

var _ = Describe("audit test", func() {
	instanceClient := topoServerClient.Instance()

	It("get user privilege user_name = 'admin'", func() {
		input := &metadata.QueryInput{
			Condition: map[string]interface{}{
				"op_time": []string{
					"2018-07-20 00:00:00",
					"2018-07-21 23:59:59",
				},
			},
			Start: 0,
			Limit: 10,
			Sort:  "-op_time",
		}
		rsp, err := instanceClient.QueryAuditLog(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})
})
