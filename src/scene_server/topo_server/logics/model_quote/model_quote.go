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

// Package modelquote defines model quote logics.
package modelquote

import (
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// ModelQuoteOperation model quote operation methods
type ModelQuoteOperation interface {
	GetQuotedObjID(kit *rest.Kit, objID string, propertyID string) (string, error)
}

// NewModelQuoteOperation create a new model quote operation instance
func NewModelQuoteOperation(client apimachinery.ClientSetInterface) ModelQuoteOperation {
	return &quote{
		clientSet: client,
	}
}

type quote struct {
	clientSet apimachinery.ClientSetInterface
}

// GetQuotedObjID get quoted object id by source object id & property id
func (q *quote) GetQuotedObjID(kit *rest.Kit, objID string, propertyID string) (string, error) {
	expr, err := filtertools.And(filtertools.GenAtomFilter(common.BKSrcModelField, filter.Equal, objID),
		filtertools.GenAtomFilter(common.BKPropertyIDField, filter.Equal, propertyID))
	if err != nil {
		return "", kit.CCError.New(common.CCErrCommParamsInvalid, err.Error())
	}

	opt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: expr},
		Page:               metadata.BasePage{Limit: 1},
		Fields:             []string{common.BKDestModelField},
	}

	res, err := q.clientSet.CoreService().ModelQuote().ListModelQuoteRelation(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("get quoted object id failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return "", err
	}

	if len(res.Info) != 1 {
		blog.Errorf("model quote relations length is invalid, obj: %s, attr: %s, rid: %s", objID, propertyID, kit.Rid)
		return "", kit.CCError.New(common.CCErrCommParamsInvalid, propertyID)
	}

	return res.Info[0].DestModel, nil
}
