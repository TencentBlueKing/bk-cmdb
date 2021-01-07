package host_server_test

import (
	"context"
	"fmt"

	"configcenter/src/common/metadata"
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
})
