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
package settemplate

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

type SetTemplateInterface interface {
	CreateSetTemplate(ctx context.Context, h http.Header, bizID int64, option metadata.CreateSetTemplateOption) (*metadata.SetTemplateResult, errors.CCErrorCoder)
	UpdateSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.UpdateSetTemplateOption) (*metadata.SetTemplateResult, errors.CCErrorCoder)
	DeleteSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder
	GetSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) (*metadata.SetTemplateResult, errors.CCErrorCoder)
	ListSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateOption) (*metadata.MultipleSetTemplateResult, errors.CCErrorCoder)
	ListSetTemplateWeb(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateOption) (*metadata.MultipleSetTemplateResult, errors.CCErrorCoder)
	ListSetTplRelatedSvcTpl(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder)
	ListSetTplRelatedSetsWeb(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.ListSetByTemplateOption) (*metadata.InstDataInfo, errors.CCErrorCoder)
	DiffSetTplWithInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) (*metadata.SetTplDiffResult, errors.CCErrorCoder)
}

func NewSetTemplateInterface(client rest.ClientInterface) SetTemplateInterface {
	return &SetTemplate{client: client}
}

type SetTemplate struct {
	client rest.ClientInterface
}
