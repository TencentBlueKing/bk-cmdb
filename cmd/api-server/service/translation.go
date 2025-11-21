/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// TranslationRequest translation request
type TranslationRequest struct {
	DefaultLang string `req:"default_lang,in:query"`
	LangPath    string `req:"lang_path,in:query"`
}

// ReloadTranslation for reload translation
func (s *Service) ReloadTranslation(kt *kit.Kit, req *TranslationRequest) (*rest.EmptyResp, error) {
	log.Info(kt, "handle ReloadTranslation")

	err := i18n.Reload(kt, &i18n.Options{
		LanguageDir: req.LangPath,
		DefaultLang: i18n.LanguageType(req.DefaultLang),
	})
	if err != nil {
		log.Error(kt, "reload language failed", log.E(err))
		return nil, err
	}

	return nil, nil
}
