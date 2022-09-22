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

/* eslint-disable no-unused-vars */

import $http from '@/api'
export default {
  namespaced: true,
  actions: {
    create(context, { params, config }) {
      return $http.post('dynamicgroup', params, config)
    },
    update(context, { bizId, id, params, config }) {
      return $http.put(`dynamicgroup/${bizId}/${id}`, params, config)
    },
    delete(context, { bizId, id, config }) {
      return $http.delete(`dynamicgroup/${bizId}/${id}`, config)
    },
    details(context, { bizId, id, config }) {
      return $http.get(`dynamicgroup/${bizId}/${id}`, config)
    },
    preview(context, { bizId, id, params, config }) {
      return $http.post(`dynamicgroup/data/${bizId}/${id}`, params, config)
    },
    search(context, { bizId, params, config }) {
      return $http.post(`dynamicgroup/search/${bizId}`, params, config)
    }
  }
}
