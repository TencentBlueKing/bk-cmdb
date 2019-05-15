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

package service

import (
	"net/http"
)

func (s *coreService) initDataSynchronize() {

	s.addAction(http.MethodPost, "/set/synchronize/instance", s.SynchronizeInstance, nil)
	s.addAction(http.MethodPost, "/set/synchronize/model", s.SynchronizeModel, nil)
	s.addAction(http.MethodPost, "/set/synchronize/association", s.SynchronizeAssociation, nil)
	s.addAction(http.MethodPost, "/read/synchronize", s.SynchronizeFind, nil)
	s.addAction(http.MethodDelete, "/clear/synchronize/data", s.SynchronizeClearData, nil)
	s.addAction(http.MethodPost, "/set/synchronize/identifier/flag", s.SetIdentifierFlag, nil)
}
