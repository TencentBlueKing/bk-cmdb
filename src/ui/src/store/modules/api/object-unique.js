/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import $http from '@/api'

const state = {

}

const getters = {

}

const actions = {
    /**
     * 添加模型唯一约束
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型英文id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createObjectUniqueConstraints ({ commit, state, dispatch }, { objId, params, config }) {
        return $http.post(`object/${objId}/unique/action/create`, params, config)
    },
    /**
     * 编辑模型唯一约束
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型英文id
     * @param {Number} id 模型唯一约束的自增ID
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateObjectUniqueConstraints ({ commit, state, dispatch }, { objId, id, params, config }) {
        return $http.put(`object/${objId}/unique/${id}/action/update`, params, config)
    },
    /**
     * 删除模型唯一约束
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型英文id
     * @param {Number} id 模型唯一约束的自增ID
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteObjectUniqueConstraints ({ commit, state, dispatch }, { objId, id, config }) {
        return $http.delete(`object/${objId}/unique/${id}/action/delete`, config)
    },
    /**
     * 删除模型唯一约束
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} objId 模型英文id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchObjectUniqueConstraints ({ commit, state, dispatch }, { objId, config }) {
        return $http.get(`object/${objId}/unique/action/search`, config)
    }
}

const mutations = {

}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
