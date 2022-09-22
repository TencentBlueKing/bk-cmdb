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

const state = {
  localProcessTemplate: []
}

const getters = {
  localProcessTemplate: state => state.localProcessTemplate,
  // eslint-disable-next-line max-len
  hasProcessName: state => process => state.localProcessTemplate.find(template => template.bk_func_name.value === process.bk_func_name.value)
}

const actions = {}

const mutations = {
  setLocalProcessTemplate(state, processes) {
    state.localProcessTemplate = processes
  },
  addLocalProcessTemplate(state, process) {
    state.localProcessTemplate.push(process)
  },
  updateLocalProcessTemplate(state, { process, index }) {
    state.localProcessTemplate.splice(index, 1, process)
  },
  deleteLocalProcessTemplate(state, { index }) {
    state.localProcessTemplate.splice(index, 1)
  },
  clearLocalProcessTemplate(state) {
    state.localProcessTemplate = []
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
