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
    searchHistory: []
}

const getters = {}

const actions = {
    /**
     * 全文检索
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
    search ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`find/full_text`, params, config)
    }
}

const mutations = {
    setSearchHistory (state, keywords) {
        const len = state.searchHistory.length
        !state.searchHistory.find(keyword => keyword === keywords) && state.searchHistory.unshift(keywords)
        if (len > 8) {
            state.searchHistory.pop(keywords)
        }
        localStorage.setItem('searchHistory', JSON.stringify(state.searchHistory))
    },
    getSearchHistory (state) {
        const history = JSON.parse(localStorage.getItem('searchHistory'))
        state.searchHistory = history || []
    },
    clearSearchHistory (state) {
        localStorage.setItem('searchHistory', JSON.stringify([]))
        state.searchHistory = []
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
