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
	CreateSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.CreateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder)
	UpdateSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.UpdateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder)
	DeleteSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder
	GetSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) (metadata.SetTemplate, errors.CCErrorCoder)
	ListSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateOption) (*metadata.MultipleSetTemplateResult, errors.CCErrorCoder)
	CountSetTplInstances(ctx context.Context, header http.Header, bizID int64, option metadata.CountSetTplInstOption) (map[int64]int64, errors.CCErrorCoder)
	ListSetServiceTemplateRelations(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) ([]metadata.SetServiceTemplateRelation, errors.CCErrorCoder)
	ListSetTplRelatedSvcTpl(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder)
	UpdateSetTemplateSyncStatus(ctx context.Context, header http.Header, setID int64, syncStatus metadata.SetTemplateSyncStatus) errors.CCErrorCoder
	ListSetTemplateSyncStatus(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
	DeleteSetTemplateSyncStatus(ctx context.Context, header http.Header, bizID int64, setIDs []int64) errors.CCErrorCoder
	ListSetTemplateSyncHistory(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
	ModifySetTemplateSyncStatus(ctx context.Context, header http.Header, setID int64, syncStatus metadata.SyncStatus) errors.CCErrorCoder
}

func NewSetTemplateInterfaceClient(client rest.ClientInterface) SetTemplateInterface {
	return &setTemplate{client: client}
}

type setTemplate struct {
	client rest.ClientInterface
}
