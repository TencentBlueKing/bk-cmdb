package host_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user operation test", func() {
	var bizId int64

	Describe("test preparation", func() {
		It("create business bk_biz_name = 'user_biz'", func() {
			input := map[string]interface{}{
				"bk_biz_name":       "user_biz",
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).To(ContainElement("user_biz"))
			bizId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("user custom test", func() {
		It("search default user custom", func() {
			rsp, err := hostServerClient.GetUserCustom(context.Background(), header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("save user custom", func() {
			input := map[string]interface{}{
				"index_v2_classify_navigation": []string{"bk_middleware"},
			}
			rsp, err := hostServerClient.SaveUserCustom(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search user custom", func() {
			rsp, err := hostServerClient.GetUserCustom(context.Background(), header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := rsp.Data.(map[string]interface{})
			Expect(data["index_v2_classify_navigation"].([]interface{})[0].(string)).To(Equal("bk_middleware"))
		})
	})

	Describe("user favorites test", func() {
		var favId string

		It("create user favorites", func() {
			input := &metadata.FavouriteParms{
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "[]",
				Name:        "123",
			}
			rsp, err := hostServerClient.AddHostFavourite(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			favId = rsp.Data.(map[string]interface{})["id"].(string)
		})

		It("create user favorites less name", func() {
			input := &metadata.FavouriteParms{
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "[]",
			}
			rsp, err := hostServerClient.AddHostFavourite(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create user favorites same name", func() {
			input := &metadata.FavouriteParms{
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "[]",
				Name:        "123",
			}
			rsp, err := hostServerClient.AddHostFavourite(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create user favorites invalid info", func() {
			input := &metadata.FavouriteParms{
				Info:        "abc",
				QueryParams: "[]",
				Name:        "1234",
			}
			rsp, err := hostServerClient.AddHostFavourite(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create user favorites invalid query_params", func() {
			input := &metadata.FavouriteParms{
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "abc",
				Name:        "12345",
			}
			rsp, err := hostServerClient.AddHostFavourite(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search user favorites", func() {
			input := &metadata.QueryInput{
				Start: 0,
				Limit: 10,
			}
			rsp, err := hostServerClient.GetHostFavourites(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(uint64(1)))
			Expect(rsp.Data.Info[0]["name"]).To(Equal("123"))
		})

		It("increase user favorites", func() {
			rsp, err := hostServerClient.IncrHostFavouritesCount(context.Background(), favId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("increase user favorites nonexist id", func() {
			rsp, err := hostServerClient.IncrHostFavouritesCount(context.Background(), "123456", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update user favorites", func() {
			input := &metadata.FavouriteParms{
				ID:          favId,
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "[]",
				Name:        "1234",
				Count:       2,
			}
			rsp, err := hostServerClient.UpdateHostFavouriteByID(context.Background(), favId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update user favorites nonexist id", func() {
			input := &metadata.FavouriteParms{
				ID:          "1000",
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "[]",
				Name:        "1234",
				Count:       2,
			}
			rsp, err := hostServerClient.UpdateHostFavouriteByID(context.Background(), "1000", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update user favorites less name", func() {
			input := &metadata.FavouriteParms{
				ID:          favId,
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "[]",
				Count:       2,
			}
			rsp, err := hostServerClient.UpdateHostFavouriteByID(context.Background(), favId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update user favorites invalid info", func() {
			input := &metadata.FavouriteParms{
				Info:        "abc",
				QueryParams: "[]",
				Name:        "1234",
			}
			rsp, err := hostServerClient.UpdateHostFavouriteByID(context.Background(), favId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update user favorites invalid query_params", func() {
			input := &metadata.FavouriteParms{
				Info:        fmt.Sprintf("{\"bk_biz_id\":%v,\"exact_search\":true,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}", bizId),
				QueryParams: "abc",
				Name:        "12345",
			}
			rsp, err := hostServerClient.UpdateHostFavouriteByID(context.Background(), favId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search user favorites", func() {
			input := &metadata.QueryInput{
				Start: 0,
				Limit: 10,
			}
			rsp, err := hostServerClient.GetHostFavourites(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(uint64(1)))
			Expect(rsp.Data.Info[0]["name"]).To(Equal("1234"))
		})

		It("delete user favorites", func() {
			rsp, err := hostServerClient.DeleteHostFavouriteByID(context.Background(), favId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete user favorites nonexist id", func() {
			rsp, err := hostServerClient.DeleteHostFavouriteByID(context.Background(), "123456", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search user favorites", func() {
			input := &metadata.QueryInput{
				Start: 0,
				Limit: 10,
			}
			rsp, err := hostServerClient.GetHostFavourites(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(uint64(0)))
		})
	})

	Describe("custom query test", func() {
		var queryId string
		var hostId int64

		It("create custom query", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"info":      "{\"condition\":[{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[],\"fields\":[\"bk_host_innerip\",\"bk_biz_name\",\"bk_set_name\",\"bk_module_name\",\"bk_cloud_id\"]}]}",
				"name":      "123",
			}
			rsp, err := hostServerClient.AddUserCustomQuery(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			queryId = rsp.Data.(map[string]interface{})["id"].(string)
		})

		It("create custom query less biz_id", func() {
			input := map[string]interface{}{
				"info": "{\"condition\":[{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[],\"fields\":[\"bk_host_innerip\",\"bk_biz_name\",\"bk_set_name\",\"bk_module_name\",\"bk_cloud_id\"]}]}",
				"name": "1234",
			}
			rsp, err := hostServerClient.AddUserCustomQuery(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create custom query less info", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"name":      "12345",
			}
			rsp, err := hostServerClient.AddUserCustomQuery(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create custom query less name", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"info":      "{\"condition\":[{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[],\"fields\":[\"bk_host_innerip\",\"bk_biz_name\",\"bk_set_name\",\"bk_module_name\",\"bk_cloud_id\"]}]}",
			}
			rsp, err := hostServerClient.AddUserCustomQuery(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create custom query same name", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"info":      "{\"condition\":[{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[],\"fields\":[\"bk_host_innerip\",\"bk_biz_name\",\"bk_set_name\",\"bk_module_name\",\"bk_cloud_id\"]}]}",
				"name":      "123",
			}
			rsp, err := hostServerClient.AddUserCustomQuery(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create custom query invalid info", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"info":      "abc",
				"name":      "abc",
			}
			rsp, err := hostServerClient.AddUserCustomQuery(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search custom query", func() {
			input := &metadata.QueryInput{
				Start: 0,
				Limit: 10,
			}
			rsp, err := hostServerClient.GetUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			count, err := commonutil.GetIntByInterface(rsp.Data.(map[string]interface{})["count"])
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("get custom query detail", func() {
			rsp, err := hostServerClient.GetUserCustomQueryDetail(context.Background(), strconv.FormatInt(bizId, 10), queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["info"]).To(Equal("{\"condition\":[{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[],\"fields\":[\"bk_host_innerip\",\"bk_biz_name\",\"bk_set_name\",\"bk_module_name\",\"bk_cloud_id\"]}]}"))
			Expect(rsp.Data["name"]).To(Equal("123"))
		})

		It("get nonexist custom query detail", func() {
			rsp, err := hostServerClient.GetUserCustomQueryDetail(context.Background(), strconv.FormatInt(bizId, 10), "100", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("get unmatching biz and custom query detail", func() {
			rsp, err := hostServerClient.GetUserCustomQueryDetail(context.Background(), "2", queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update custom query", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"info":      "{\"condition\":[{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[],\"fields\":[\"bk_host_innerip\",\"bk_biz_name\",\"bk_set_name\",\"bk_module_name\",\"bk_cloud_id\"]}]}",
				"name":      "1234",
				"id":        queryId,
			}
			rsp, err := hostServerClient.UpdateUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), queryId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update custom query invalid info", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"info":      "abc",
				"name":      "abc",
			}
			rsp, err := hostServerClient.UpdateUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), queryId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get custom query detail", func() {
			rsp, err := hostServerClient.GetUserCustomQueryDetail(context.Background(), strconv.FormatInt(bizId, 10), queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["name"]).To(Equal("1234"))
		})

		It("add host", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "3.0.0.1",
						"bk_asset_id":     "addhost_api_asset_1",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search host created using api", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"3.0.0.1"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("3.0.0.1"))
			hostId, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("get custom query data", func() {
			rsp, err := hostServerClient.GetUserCustomQueryResult(context.Background(), strconv.FormatInt(bizId, 10), queryId, "0", "10", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			data := metadata.SearchHostResult{}
			json.Unmarshal(j, &data)
			Expect(data.Data.Count).To(Equal(1))
			host := data.Data.Info[0]["host"].(map[string]interface{})
			Expect(host["bk_host_innerip"].(string)).To(Equal("3.0.0.1"))
			hostIdRes, err := commonutil.GetInt64ByInterface(host["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(hostIdRes).To(Equal(hostId))
		})

		It("get custom query data invalid biz_id", func() {
			rsp, err := hostServerClient.GetUserCustomQueryResult(context.Background(), "1000", queryId, "0", "10", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get custom query data invalid id", func() {
			rsp, err := hostServerClient.GetUserCustomQueryResult(context.Background(), strconv.FormatInt(bizId, 10), "123456", "0", "10", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get custom query data unmatching id and biz_id", func() {
			rsp, err := hostServerClient.GetUserCustomQueryResult(context.Background(), "2", queryId, "0", "10", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get custom query data invalid start", func() {
			rsp, err := hostServerClient.GetUserCustomQueryResult(context.Background(), strconv.FormatInt(bizId, 10), queryId, "erfre", "10", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get custom query data invalid limit", func() {
			rsp, err := hostServerClient.GetUserCustomQueryResult(context.Background(), strconv.FormatInt(bizId, 10), queryId, "0", "erfre", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete custom query", func() {
			rsp, err := hostServerClient.DeleteUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete custom query invalid biz_id", func() {
			rsp, err := hostServerClient.DeleteUserCustomQuery(context.Background(), "1234", queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete custom query invalid id", func() {
			rsp, err := hostServerClient.DeleteUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), "12345", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete custom query unmatching biz_id and id", func() {
			rsp, err := hostServerClient.DeleteUserCustomQuery(context.Background(), "2", queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete custom query twice", func() {
			rsp, err := hostServerClient.DeleteUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), queryId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search custom query", func() {
			input := &metadata.QueryInput{
				Start: 0,
				Limit: 10,
			}
			rsp, err := hostServerClient.GetUserCustomQuery(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			count, err := commonutil.GetIntByInterface(rsp.Data.(map[string]interface{})["count"])
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
