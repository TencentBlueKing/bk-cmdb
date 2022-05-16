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

export default {
  namespaced: true,
  state: {
    /**
     * 业务集 ID
     * @type {number}
     */
    bizSetId: null,
    /**
     * 业务集名称
     * @type {string}
     */
    bizSetName: null,
    /**
     * 业务集列表
     * @type {Array}
     */
    bizSetList: [],
    /**
     * 当前所选节点的业务 ID
     * @type {number}
     */
    bizId: null,
  },
  mutations: {
    setBizSetId(state, bizSetId) {
      state.bizSetId = Number(bizSetId)
      state.bizSetName = state.bizSetList.find(item => item.bk_biz_set_id === state.bizSetId)?.bk_biz_set_name
    },
    setBizSetList(state, bizSetList) {
      state.bizSetList = bizSetList
    },
    setBizId(state, bizId) {
      state.bizId = Number(bizId)
    }
  }
}
