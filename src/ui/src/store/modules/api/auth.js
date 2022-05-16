/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
const actions = {
  async getViewAuth(context, viewAuthData) {
    if (window.Site.authscheme !== 'iam') {
      return Promise.resolve(true)
    }
    const result = await $http.post('auth/verify', {
      // eslint-disable-next-line new-cap
      resources: TRANSFORM_TO_INTERNAL(viewAuthData)
    })
    return Promise.resolve(result.every(data => data.is_pass))
  },
  async getSkipUrl(context, { params, config = {} }) {
    return $http.post('auth/skip_url', params, config)
  }
}

export default {
  namespaced: true,
  actions
}
