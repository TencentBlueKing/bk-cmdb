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
const cloudArea = {
  namespaced: true,
  actions: {
    findMany({ commit, state, dispatch }, { params, config }) {
      return $http.post('findmany/cloudarea', params, config)
    },
    delete(context, { id, config }) {
      return $http.delete(`delete/cloudarea/${id}`, config)
    },
    batchCreate(context, { params, config }) {
      return $http.post('createmany/cloudarea', params, config)
    },
    update(context, { id, params, config }) {
      return $http.put(`update/cloudarea/${id}`, params, config)
    },
    getHostCount(context, { params, config }) {
      return $http.post('findmany/cloudarea/hostcount', params, config)
    }
  }
}

const cloudAccount = {
  namespaced: true,
  actions: {
    findMany(context, { params, config }) {
      return $http.post('findmany/cloud/account', params, config)
    },
    async findOne(context, { id, config }) {
      return context.dispatch('findMany', {
        params: {
          condition: {
            bk_account_id: {
              $eq: id
            }
          }
        },
        config
      }).then(({ info }) => {
        if (!info.length) {
          return Promise.reject(new Error(`Can not find cloud account with id:${id}`))
        }
        return info[0]
      })
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
    },
    getStatus(context, { params, config }) {
      return $http.post('findmany/cloud/account/validity', params, config)
    }
  }
}

const cloudResource = {
  namespaced: true,
  actions: {
    createTask(context, { params, config }) {
      return $http.post('create/cloud/sync/task', params, config)
    },
    findTask(context, { params, config }) {
      return $http.post('findmany/cloud/sync/task', params, config)
    },
    findOneTask(context, { id, config }) {
      return context.dispatch('findTask', {
        params: {
          condition: {
            bk_task_id: {
              $eq: id
            }
          },
          latest_hostcount: true
        },
        config
      }).then(({ info }) => {
        if (!info.length) {
          return Promise.reject(new Error(`Can not find cloud task with id:${id}`))
        }
        return info[0]
      })
    },
    updateTask(context, { id, params, config }) {
      return $http.put(`update/cloud/sync/task/${id}`, params, config)
    },
    deleteTask(context, { id, config }) {
      return $http.delete(`delete/cloud/sync/task/${id}`, config)
    },
    findRegion({ state, commit }, { params, config }) {
      return $http.post('findmany/cloud/sync/region', params, config)
    },
    findVPC(context, { id, params, config }) {
      return $http.post(`findmany/cloud/account/vpc/${id}`, params, config)
    },
    findHistory(context, { params, config }) {
      return $http.post('findmany/cloud/sync/history', params, config)
    }
  }
}

export default {
  namespaced: true,
  modules: {
    area: cloudArea,
    account: cloudAccount,
    resource: cloudResource
  }
}
