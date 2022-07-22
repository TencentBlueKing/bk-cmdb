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

const getUpdateParams = (payload) => {
  // eslint-disable-next-line
    const { updateType, objId, relation, id, value, params, multiple } = payload
  let newRelation = [...relation]
  if (multiple) {
    if (updateType === 'remove') {
      newRelation = newRelation.filter(associatedValue => associatedValue !== value)
    } else {
      newRelation.push(value)
    }
  } else {
    if (updateType === 'remove') {
      newRelation = []
    } else {
      newRelation = [value]
    }
  }
  const updateParams = { ...params }
  updateParams[id] = newRelation.join(',')
  return updateParams
}

const actions = {
  /**
     * 获取实例关联关系
     * @param {String} objId 模型ID
     * @param {String} instId 实例ID
     * @param {Object} config API请求配置
     * @return {Promise} promise 对象
     */
  getInstRelation({ commit, state, dispatch, rootGetters }, { objId, instId, params, config }) {
    return $http.post(`find/instassttopo/object/${objId}/inst/${instId}`, params, config)
  },

  getInstRelationTopo({ commit, state, dispatch, rootGetters }, { objId, instId, params, config }) {
    return $http.post(`/findmany/inst/association/object/${objId}/inst_id/${instId}/offset/0/limit/200/web`, params, config)
  },

  updateInstRelation({ commit, state, dispatch, rootGetters }, { params, config }) {
    const updateParams = getUpdateParams(params)
    let promise
    switch (params.objId) {
      case 'host':
        promise = $http.put('hosts/batch', updateParams)
        break
      case 'biz':
        promise = $http.put(`biz/${rootGetters.supplierAccount}/${params.bk_biz_id}`, updateParams)
        break
      default:
        promise = $http.put(`inst/${rootGetters.supplierAccount}/${params.objId}/${params.bk_inst_id}`, updateParams)
    }
    return promise
  }
}

export default {
  namespaced: true,
  actions
}
