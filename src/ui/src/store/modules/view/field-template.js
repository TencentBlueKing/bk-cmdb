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

const defaultTemplateDraft = () => ({
  basic: {},
  fieldList: null,
  uniqueList: null
})
const sessionTemplateDraft = () => Object.assign(defaultTemplateDraft(), JSON.parse(sessionStorage.getItem('sessionTemplateDraft')) || {})
const setSessionTemplateDraft = val =>  sessionStorage.setItem('sessionTemplateDraft', val)

const state = {
  // 草稿数据，属性与唯一校验格式与接口参数一致
  templateDraft: sessionTemplateDraft()
}

const getters = {
  templateDraft: state => state.templateDraft
}

const mutations = {
  setTemplateDraft(state, templateDraft) {
    state.templateDraft = Object.assign(state.templateDraft, templateDraft)
    setSessionTemplateDraft(JSON.stringify(state.templateDraft))
  },
  clearTemplateDraft(state) {
    setSessionTemplateDraft(JSON.stringify(defaultTemplateDraft()))
    state.templateDraft = defaultTemplateDraft()
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations
}
