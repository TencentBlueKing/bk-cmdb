/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful/v3"
)

type migrateCloudAccountCryptoOption struct {
	AccountIDs []int64          `json:"bk_account_ids"`
	Old        *oldCryptoOption `json:"old"`
	New        *cryptor.Config  `json:"new"`
}

type oldCryptoOption struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
}

// Validate migrateCloudAccountCryptoOption
func (m migrateCloudAccountCryptoOption) Validate(kit *rest.Kit) error {
	if len(m.AccountIDs) == 0 {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "bk_account_ids")
	}

	if len(m.AccountIDs) > common.BKMaxLimitSize {
		return kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "bk_account_ids", common.BKMaxLimitSize)
	}

	if m.Old == nil {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "old")
	}

	if m.New == nil {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "new")
	}

	if m.Old.Enabled && len(m.Old.Key) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "old.key")
	}

	if !m.Old.Enabled && !m.New.Enabled {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "old.enabled & new.enabled")
	}

	if err := m.New.Validate(); err != nil {
		return err
	}

	return nil
}

// migrateCloudAccountCrypto migrate cloud account crypto algorithm
func (s *Service) migrateCloudAccountCrypto(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	opt := new(migrateCloudAccountCryptoOption)
	if err := json.NewDecoder(req.Request.Body).Decode(opt); err != nil {
		blog.Errorf("decode request body failed, err: %v, rid: %s", err, kit.Rid)
		_ = resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if err := opt.Validate(kit); err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	// get cloud accounts from db
	cond := mapstr.MapStr{common.BKCloudAccountID: mapstr.MapStr{common.BKDBIN: opt.AccountIDs}}
	accounts := make([]metadata.CloudAccount, 0)
	err := s.db.Table(common.BKTableNameCloudAccount).Find(cond).Fields(common.BKCloudAccountID, "bk_secret_key").
		All(kit.Ctx, &accounts)
	if err != nil {
		blog.Errorf("get cloud accounts failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		_ = resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: kit.CCError.CCError(common.CCErrCommDBSelectFailed)})
		return
	}

	// migrate cloud account crypto algorithm
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		for _, account := range accounts {
			// decrypt old secret key
			secretKey := account.SecretKey
			if opt.Old.Enabled {
				secretKey, err = cryptor.NewAesEncrpytor(opt.Old.Key).Decrypt(account.SecretKey)
				if err != nil {
					blog.Errorf("decrypt secret key failed, err: %v, rid: %s", err, kit.Rid)
					return err
				}
			}

			// encrypt new secret key
			if opt.New.Enabled {
				crypto, err := cryptor.NewCrypto(opt.New)
				if err != nil {
					blog.Errorf("new crypto failed, err: %v, config: %+v, rid: %s", err, opt.New, kit.Rid)
					return err
				}

				secretKey, err = crypto.Encrypt(secretKey)
				if err != nil {
					blog.Errorf("encrypt secret key failed, err: %v, rid: %s", err, kit.Rid)
					return err
				}
			}

			updateCond := mapstr.MapStr{common.BKCloudAccountID: account.AccountID}
			updateData := mapstr.MapStr{"bk_secret_key": secretKey}
			err = s.db.Table(common.BKTableNameCloudAccount).Update(kit.Ctx, updateCond, updateData)
			if err != nil {
				blog.Errorf("update cloud account %d failed, err: %v, rid: %s", account.AccountID, err, kit.Rid)
				return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
			}
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}
