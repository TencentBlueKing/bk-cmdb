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

const actions = {
  search(context, { params, config }) {
    return $http.post('findmany/cloud/account', params, config)
  },
  async searchById(context, { id, config }) {
    try {
      const { info } = await context.dispatch('search', {
        params: {
          condition: { bk_account_id: id },
          page: {
            start: 0,
            limit: 1
          }
        },
        config
      })
      if (info.length) {
        return Promise.resolve(info[0])
      }
      throw new Error(`Can not find cloud account with id:${id}`)
    } catch (e) {
      return Promise.reject(e)
    }
  },
  verify(context, { params, config }) {
    return $http.post('cloud/account/verify', params, config)
  },
  create(context, { params, config }) {
    return $http.post('create/cloud/account', params, config)
  },
  update(context, { id, params, config }) {
    return $http.put(`update/cloud/account/${id}`, params, config)
  },
  delete(context, { id, config }) {
    return $http.delete(`delete/cloud/account/${id}`, config)
  }
}

export default {
  namespaced: true,
  actions
}
