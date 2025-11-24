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

package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// Recoverer is a middleware that recovers from panics, logs the panic and stack trace
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if fatalErr := recover(); fatalErr != nil {
				w.WriteHeader(http.StatusInternalServerError)

				msg := fmt.Sprintf("panic err: %v", fatalErr)
				log.Error(r.Context(), msg, "stack_trace", debug.Stack())
			}
		}()

		next.ServeHTTP(w, r)
	})
}
