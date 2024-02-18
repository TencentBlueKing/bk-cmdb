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

import isEqual from 'lodash/isEqual'
import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
import { isViewAuthFreeModel } from '@/service/auth'

const state = {
  authedList: []
}

const getters = {
  isViewAuthed: state => (auth) => {
    // 比较时展开...authKey是必要的，因state数据会包含__ob__对象，展开后会将其忽略
    const found = state.authedList.find(([authKey]) => isEqual({ ...authKey }, TRANSFORM_TO_INTERNAL(auth)[0]))
    return found?.[1]
  }
}

const mutations = {
  setAuthedList(state, list) {
    state.authedList = list
  }
}

const actions = {
  // 仅用于查询view级别的权限，并且假定viewAuthData为单个权限（非数组）
  async getViewAuth({ getters }, viewAuthData) {
    if (window.Site.authscheme !== 'iam') {
      return Promise.resolve(true)
    }

    if (viewAuthData.type === 'R_MODEL' && isViewAuthFreeModel({ id: viewAuthData.relation[0] })) {
      return Promise.resolve(true)
    }

    // 优先使用state中的权限状态判断，如果能查找到(非undefined)则可直接作为鉴权结果返回
    // 此行为是为了最大程度的减少请求鉴权的耗时，尽快进入页面
    const isViewAuthed = getters.isViewAuthed(viewAuthData)
    if (isViewAuthed !== undefined) {
      return Promise.resolve(isViewAuthed)
    }

    const result = await $http.post('auth/verify', {
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
  state,
  mutations,
  getters,
  actions
}
