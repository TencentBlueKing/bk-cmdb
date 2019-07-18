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

package core

import (
	"context"
	"net/http"
	"time"

	"configcenter/src/common/errors"
	"configcenter/src/common/language"
)

// ContextParams the core function params
type ContextParams struct {
	context.Context
	Header          http.Header
	SupplierAccount string
	User            string
	ReqID           string
	Error           errors.DefaultCCErrorIf
	Lang            language.DefaultCCLanguageIf
}

// Deadline overwrite Context Deadline methods
func (c ContextParams) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

// Done overwrite Context Done methods
func (c ContextParams) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err overwrite Context Err methods
func (c ContextParams) Err() error {
	return c.Context.Err()
}

// Value overwrite Context Value methods
func (c ContextParams) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}
