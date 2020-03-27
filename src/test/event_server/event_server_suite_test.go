package event_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"configcenter/src/common/metadata"
	"configcenter/src/test"
	"configcenter/src/test/reporter"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var eventServerClient = test.GetClientSet().EventServer()

func TestEventServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"eventserver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "EventServer Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})

var _ = Describe("event server test", func() {
	var _ = Describe("subscribe event test", func() {
		var subscriptionId1, subscriptionId2 string
		/*
			It("ping subscription", func() {
				input := map[string]interface{}{}
				rsp, err := eventServerClient.Ping(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("testing connection", func() {
				input := map[string]interface{}{}
				rsp, err := eventServerClient.Telnet(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})
		*/

		It("subscribe event bk_biz_id = 0", func() {
			input := &metadata.Subscription{
				SubscriptionName: "123",
				SystemName:       "cmdb",
				CallbackURL:      "http://127.0.0.1:8080",
				ConfirmMode:      "httpstatus",
				ConfirmPattern:   "200",
				SubscriptionForm: "hostdelete",
				TimeOutSeconds:   60,
			}
			rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			data := metadata.RspSubscriptionCreate{}
			json.Unmarshal(j, &data)
			subscriptionId1 = fmt.Sprintf("%d", data.Data.SubscriptionID)
		})

		It("subscribe event bk_biz_id = 0 and subscription_name = 'dwe'", func() {
			input := &metadata.Subscription{
				SubscriptionName: "dwe",
				SystemName:       "cmdb",
				CallbackURL:      "http://127.0.0.1:8080",
				ConfirmMode:      "httpstatus",
				ConfirmPattern:   "200",
				SubscriptionForm: "hostdelete",
				TimeOutSeconds:   60,
			}
			rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			data := metadata.RspSubscriptionCreate{}
			json.Unmarshal(j, &data)
			subscriptionId2 = fmt.Sprintf("%d", data.Data.SubscriptionID)
		})

		It("search subscribe bk_biz_id = 0", func() {
			input := metadata.ParamSubscriptionSearch{
				Page: metadata.BasePage{
					Sort:  "subscription_id",
					Limit: 10,
					Start: 0,
				},
			}
			rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.RspSubscriptionSearch{}
			json.Unmarshal(j, &data)
			// count is actually total, so it should be 2 instead of 1
			Expect(int(data.Count)).To(Equal(2))
			Expect(data.Info[0].SubscriptionName).To(Equal("123"))
			Expect(data.Info[1].SubscriptionName).To(Equal("dwe"))
		})

		It("search subscribe start = 1 and limit = 1", func() {
			input := metadata.ParamSubscriptionSearch{
				Page: metadata.BasePage{
					Sort:  "subscription_id",
					Limit: 1,
					Start: 1,
				},
			}
			rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.RspSubscriptionSearch{}
			json.Unmarshal(j, &data)
			Expect(int(data.Count)).To(Equal(2))
			Expect(data.Info[0].SubscriptionName).To(Equal("dwe"))
		})

		It("update event subscribe bk_biz_id = 0 and subscription_id = 1", func() {
			input := &metadata.Subscription{
				SubscriptionID:   1,
				SubscriptionName: "123d",
				SubscriptionForm: "hostcreate,hostdelete",
				SystemName:       "cmdb",
				CallbackURL:      "http://127.0.0.1:8080",
				ConfirmMode:      "httpstatus",
				ConfirmPattern:   "",
				TimeOutSeconds:   60,
			}
			rsp, err := eventServerClient.Rebook(context.Background(), "0", "0", subscriptionId1, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("unsubscribe event bk_biz_id = 0 and subscription_id = 2", func() {
			rsp, err := eventServerClient.UnSubscribe(context.Background(), "0", "0", subscriptionId2, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search subscribe bk_biz_id = 0", func() {
			input := metadata.ParamSubscriptionSearch{
				Page: metadata.BasePage{
					Sort:  "subscription_id",
					Limit: 10,
					Start: 0,
				},
			}
			rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.RspSubscriptionSearch{}
			json.Unmarshal(j, &data)
			Expect(int(data.Count)).To(Equal(1))
			Expect(data.Info[0].SubscriptionName).To(Equal("123d"))
			Expect(data.Info[0].SubscriptionForm).To(Equal("hostcreate,hostdelete"))
			Expect(data.Info[0].SystemName).To(Equal("cmdb"))
			Expect(data.Info[0].CallbackURL).To(Equal("http://127.0.0.1:8080"))
			Expect(data.Info[0].ConfirmMode).To(Equal("httpstatus"))
			Expect(data.Info[0].ConfirmPattern).To(Equal("200"))
			Expect(int(data.Info[0].TimeOutSeconds)).To(Equal(60))
		})

		It("unsubscribe event bk_biz_id = 0 and subscription_id = 1", func() {
			rsp, err := eventServerClient.UnSubscribe(context.Background(), "0", "0", subscriptionId1, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
		test.ClearDatabase()
	})

	var _ = Describe("subscribe event abnormal test", func() {
		var subscriptionId int64
		Describe("subscribe event missing parameters", func() {
			It("subscribe event missing callback url", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					ConfirmMode:      "httpstatus",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("subscribe event missing confirm mode", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("subscribe event missing subscription form", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmMode:      "httpstatus",
					ConfirmPattern:   "200",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})
		})

		Describe("subscribe event with duplicate subscription_name", func() {
			It("subscribe event bk_biz_id = 0 and subscription_name = 'dwe'", func() {
				input := &metadata.Subscription{
					SubscriptionName: "dwe",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmMode:      "httpstatus",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp)
				data := metadata.RspSubscriptionCreate{}
				json.Unmarshal(j, &data)
				subscriptionId = data.Data.SubscriptionID
			})

			It("search subscribe bk_biz_id = 0", func() {
				input := metadata.ParamSubscriptionSearch{
					Page: metadata.BasePage{
						Sort:  "subscription_id",
						Limit: 10,
						Start: 0,
					},
				}
				rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.RspSubscriptionSearch{}
				json.Unmarshal(j, &data)
				Expect(int(data.Count)).To(Equal(1))
				Expect(data.Info[0].SubscriptionName).To(Equal("dwe"))
				Expect(data.Info[0].SystemName).To(Equal("cmdb"))
				Expect(data.Info[0].CallbackURL).To(Equal("http://127.0.0.1:8080"))
				Expect(data.Info[0].ConfirmMode).To(Equal("httpstatus"))
				Expect(data.Info[0].ConfirmPattern).To(Equal("200"))
				Expect(data.Info[0].SubscriptionForm).To(Equal("hostdelete"))
				Expect(int(data.Info[0].TimeOutSeconds)).To(Equal(60))
			})

			It("subscribe event bk_biz_id = 0 and subscription_name = 'dwe'", func() {
				input := &metadata.Subscription{
					SubscriptionName: "dwe",
					SystemName:       "abc",
					CallbackURL:      "http://127.0.0.1:8080/callback",
					ConfirmMode:      "regular",
					ConfirmPattern:   ".*",
					SubscriptionForm: "hostcreate",
					TimeOutSeconds:   10,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("search subscribe bk_biz_id = 0", func() {
				input := metadata.ParamSubscriptionSearch{
					Page: metadata.BasePage{
						Sort:  "subscription_id",
						Limit: 10,
						Start: 0,
					},
				}
				rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.RspSubscriptionSearch{}
				json.Unmarshal(j, &data)
				Expect(int(data.Count)).To(Equal(1))
				Expect(data.Info[0].SubscriptionName).To(Equal("dwe"))
				Expect(data.Info[0].SystemName).To(Equal("cmdb"))
				Expect(data.Info[0].CallbackURL).To(Equal("http://127.0.0.1:8080"))
				Expect(data.Info[0].ConfirmMode).To(Equal("httpstatus"))
				Expect(data.Info[0].ConfirmPattern).To(Equal("200"))
				Expect(data.Info[0].SubscriptionForm).To(Equal("hostdelete"))
				Expect(int(data.Info[0].TimeOutSeconds)).To(Equal(60))
			})
		})

		Describe("subscribe event with invalid parameters", func() {
			It("subscribe event with invalid confirm mode", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmMode:      "123",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("subscribe event with empty subscription_name", func() {
				input := &metadata.Subscription{
					SubscriptionName: "",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmMode:      "httpstatus",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			/*
				It("subscribe event with mismatch parameter type", func() {
					input := &metadata.Subscription{
						SubscriptionName: 123,
						SystemName:       "cmdb",
						CallbackURL:      "http://127.0.0.1:8080",
						ConfirmMode:      "httpstatus",
						ConfirmPattern:   "200",
						SubscriptionForm: "hostdelete",
						TimeOutSeconds:   60,
					}
					rsp, err := eventServerClient.Subscribe(context.Background(), "0", "0", header, input)
					util.RegisterResponse(rsp)
					Expect(err).NotTo(HaveOccurred())
					Expect(rsp.Result).To(Equal(false))
				})
			*/
		})

		Describe("search subscribe event with invalid parameters", func() {
			It("search subscribe start = -1", func() {
				input := metadata.ParamSubscriptionSearch{
					Page: metadata.BasePage{
						Sort:  "subscription_id",
						Limit: 10,
						Start: -1,
					},
				}
				rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("search subscribe limit = -1", func() {
				input := metadata.ParamSubscriptionSearch{
					Page: metadata.BasePage{
						Sort:  "subscription_id",
						Limit: -1,
						Start: 0,
					},
				}
				rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.RspSubscriptionSearch{}
				json.Unmarshal(j, &data)
				Expect(int(data.Count)).To(Equal(1))
			})
		})

		Describe("update subscribe event missing parameters", func() {
			It("update subscribe event missing callback url", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					ConfirmMode:      "httpstatus",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Rebook(context.Background(), "0", "0", fmt.Sprintf("%d", subscriptionId), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("update subscribe event missing confirm mode", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmPattern:   "200",
					SubscriptionForm: "hostdelete",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Rebook(context.Background(), "0", "0", fmt.Sprintf("%d", subscriptionId), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("update subscribe event missing subscription form", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					CallbackURL:      "http://127.0.0.1:8080",
					ConfirmMode:      "httpstatus",
					ConfirmPattern:   "200",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Rebook(context.Background(), "0", "0", fmt.Sprintf("%d", subscriptionId), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})
		})

		Describe("update subscribe event with invalid parameters", func() {
			It("update event subscribe with invalid subscription_id", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "cmdb",
					ConfirmMode:      "regular",
					ConfirmPattern:   ".*",
					SubscriptionForm: "hostcreate",
					CallbackURL:      "http://127.0.0.1:8080",
					TimeOutSeconds:   60,
				}
				rsp, err := eventServerClient.Rebook(context.Background(), "0", "0", "1000", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})
		})

		Describe("update subscribe event with subscription_id = 1", func() {
			It("update event subscribe", func() {
				input := &metadata.Subscription{
					SubscriptionName: "123",
					SystemName:       "abc",
					ConfirmMode:      "regular",
					ConfirmPattern:   ".*",
					SubscriptionForm: "hostcreate",
					CallbackURL:      "http://127.0.0.1:8080/callback",
					TimeOutSeconds:   10,
				}
				rsp, err := eventServerClient.Rebook(context.Background(), "0", "0", fmt.Sprintf("%d", subscriptionId), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search subscribe bk_biz_id = 0", func() {
				input := metadata.ParamSubscriptionSearch{
					Page: metadata.BasePage{
						Sort:  "subscription_id",
						Limit: 10,
						Start: 0,
					},
				}
				rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.RspSubscriptionSearch{}
				json.Unmarshal(j, &data)
				Expect(int(data.Count)).To(Equal(1))
				Expect(data.Info[0].SubscriptionName).To(Equal("123"))
				Expect(data.Info[0].SystemName).To(Equal("abc"))
				Expect(data.Info[0].CallbackURL).To(Equal("http://127.0.0.1:8080/callback"))
				Expect(data.Info[0].ConfirmMode).To(Equal("regular"))
				Expect(data.Info[0].ConfirmPattern).To(Equal(".*"))
				Expect(data.Info[0].SubscriptionForm).To(Equal("hostcreate"))
				Expect(int(data.Info[0].TimeOutSeconds)).To(Equal(10))
			})
		})

		Describe("unsubscribe event with invalid parameters", func() {
			It("unsubscribe event with invalid subscription_id", func() {
				rsp, err := eventServerClient.UnSubscribe(context.Background(), "0", "0", "100", header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})
		})

		Describe("unsubscribe event twice", func() {
			It("unsubscribe event with subscription_id = 1", func() {
				rsp, err := eventServerClient.UnSubscribe(context.Background(), "0", "0", fmt.Sprintf("%d", subscriptionId), header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search subscribe bk_biz_id = 0", func() {
				input := metadata.ParamSubscriptionSearch{
					Page: metadata.BasePage{
						Sort:  "subscription_id",
						Limit: 10,
						Start: 0,
					},
				}
				rsp, err := eventServerClient.Query(context.Background(), "0", "0", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.RspSubscriptionSearch{}
				json.Unmarshal(j, &data)
				Expect(int(data.Count)).To(Equal(0))
			})

			It("unsubscribe event with subscription_id = 1", func() {
				rsp, err := eventServerClient.UnSubscribe(context.Background(), "0", "0", fmt.Sprintf("%d", subscriptionId), header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})
		})
		test.ClearDatabase()
	})
})
