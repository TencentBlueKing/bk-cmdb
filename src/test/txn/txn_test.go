package txn_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction Test", func() {
	Describe("test AutoRunTxn function", func() {
		It("test txn", func() {
			ctx := context.Background()
			h := test.GetHeader()
			h.Add(common.BKHTTPCCRequestID, "integration-test")
			objectID := "transaction"

			By("create transaction object with transaction")
			// use transaction to create transaction model and it's attributes.
			err := clientSet.CoreService().Txn().AutoRunTxn(ctx, h, func() error {
				// create model transaction
				inputTxn := metadata.CreateModel{
					Spec: metadata.Object{
						ObjCls:      "bk_network",
						ObjectID:    objectID,
						ObjectName:  "事务",
						IsPre:       false,
						IsPaused:    false,
						OwnerID:     "0",
						Description: "",
						Creator:     "cc_system",
					},
					Attributes: nil,
				}
				result, err := clientSet.CoreService().Model().CreateModel(ctx, h, &inputTxn)
				Expect(err).Should(BeNil())
				fmt.Printf("result: %v", result)
				Expect(result.Result).Should(BeTrue())

				return nil
			})

			Expect(err).Should(BeNil())

			By("pretend to delete start property in transaction object, and let the transaction failed.")
			// delete "start" attributes, and then failed.
			pretendErr := errors.New("pretend failed")
			pretendHeader := util.CCHeader(h)
			err = clientSet.CoreService().Txn().AutoRunTxn(ctx, pretendHeader, func() error {
				// pretend to delete a attribute
				opt := metadata.DeleteOption{
					Condition: mapstr.MapStr{
						common.BKObjIDField: objectID,
					},
				}
				result, err := clientSet.CoreService().Model().DeleteModel(ctx, pretendHeader, &opt)
				Expect(err).Should(BeNil())
				Expect(result.Result).Should(BeTrue())
				// now we really delete the attribute from db in this transaction

				// we return a err to pretend that an error is really occurred.
				return pretendErr
			})
			Expect(err).Should(Equal(pretendErr))

			By("check the *start* property in transaction object is exist or not, if exist then transaction is ok")
			// now we check if the attribute is still exist, if yes, then transaction is ok
			query := metadata.QueryCondition{
				Condition: mapstr.MapStr{
					common.BKObjIDField: objectID,
				},
			}
			// after we delete the transaction model with a pretend err, we do not really delete it.
			// so the model is still exist for now.
			attResult, err := clientSet.CoreService().Model().ReadModel(ctx, test.GetHeader(), &query)
			Expect(err).Should(BeNil())
			Expect(attResult.Result).Should(BeTrue())
			Expect(attResult.Data.Info[0].Spec.ObjectID).Should(Equal(objectID))

			// use commit and abort transaction to test transaction.
			ops := metadata.TxnOption{
				Timeout: time.Minute,
			}
			// create a new txnHeader to deliver the transaction info.
			txnHeader := util.CloneHeader(h)
			// create a new transaction
			txn, err := clientSet.CoreService().Txn().NewTransaction(txnHeader, ops)
			Expect(err).Should(BeNil())

			attributeTxn := metadata.CreateModelAttributes{
				Attributes: []metadata.Attribute{
					{
						OwnerID:       "0",
						ObjectID:      objectID,
						PropertyID:    "start",
						PropertyName:  "发起事务",
						PropertyGroup: "default",
						PropertyType:  "singlechar",
						Creator:       "cc_system",
					},
					{
						OwnerID:       "0",
						ObjectID:      objectID,
						PropertyID:    "commit",
						PropertyName:  "提交事务",
						PropertyGroup: "default",
						PropertyType:  "singlechar",
						Creator:       "cc_system",
					},
					{
						OwnerID:       "0",
						ObjectID:      objectID,
						PropertyID:    "abort",
						PropertyName:  "回滚事务",
						PropertyGroup: "default",
						PropertyType:  "singlechar",
						Creator:       "cc_system",
					},
				},
			}

			// do business logic.
			result, err := clientSet.CoreService().Model().CreateModelAttrs(ctx, txnHeader, objectID, &attributeTxn)
			Expect(err).Should(BeNil())
			// commit the transaction so that the attributes is committed.
			cmtErr := txn.CommitTransaction(context.Background(), txnHeader)
			Expect(cmtErr).Should(BeNil())
			Expect(result.Result).Should(BeTrue())
		})
	})

})
