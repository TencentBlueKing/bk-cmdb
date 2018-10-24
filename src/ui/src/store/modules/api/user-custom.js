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
    classifyNavigationKey: 'index_v3_classify_navigation',
    recentlyKey: 'index_v3_recently',
    usercustom: {}
}

const getters = {
    classifyNavigationKey: state => state.classifyNavigationKey,
    recentlyKey: state => state.recentlyKey,
    usercustom: state => state.usercustom
}

const actions = {
    /**
     * 保存用户字段配置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    saveUsercustom ({ commit, state, dispatch }, usercustom = {}) {
        return $http.post(`usercustom`, usercustom, {cancelWhenRouteChange: false}).then(() => {
            $http.cancelCache('searchUserCustom')
            commit('setUsercustom', usercustom)
            return state.usercustom
        })
    },

    /**
     * 获取用户字段配置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    searchUsercustom ({ commit, state, dispatch }, {config}) {
        const mergedConfig = Object.assign({
            requestId: 'searchUserCustom'
        }, config)
        return $http.post(`usercustom/user/search`, {}, mergedConfig).then(usercustom => {
            commit('setUsercustom', usercustom)
            return usercustom
        })
    },

    /**
     * 获取默认字段配置
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @return {promises} promises 对象
     */
    getUserDefaultCustom ({ commit, state, dispatch }) {
        return $http.post(`usercustom/default/search`)
    }
}

const mutations = {
    setUsercustom (state, usercustom) {
        state.usercustom = Object.assign({}, state.usercustom, usercustom)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
