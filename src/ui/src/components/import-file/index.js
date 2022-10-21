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

import Vue, { reactive } from 'vue'
import contentComponent from './import'
import store from '@/store'
import i18n from '@/i18n'
import useStep from './step'
import useFile from './file'
import { bkInfoBox } from 'bk-magic-vue'
let instance = null
const state = reactive({
  visible: false,
  title: null,
  submit: null,
  success: null,
  bk_obj_id: null,
  bk_biz_id: null,
  fileTips: null,
  template: null
})
const setState = (newState) => {
  Object.assign(state, newState)
}

const show = (props = {}) => {
  instance = createSideslider(props)
  instance.$mount()
  state.visible = true
}
const close = () => {
  setState({
    visible: false,
    title: null,
    submit: null,
    success: null,
    bk_biz_id: null,
    bk_obj_id: null,
    fileTips: null,
    template: null
  })
  setTimeout(() => {
    instance.$destroy()
    instance = null
    const [, { reset: resetStep }] = useStep()
    resetStep()
    const [, { clear: clearFile }] = useFile()
    clearFile()
  }, 200)
}

const beforeClose = () => {
  const [{ state }] = useFile()
  const [current] = useStep()
  if (current.value !== 2 || !state.value === 'pending') {
    return true
  }
  return new Promise((resolve) => {
    bkInfoBox({
      title: i18n.t('确认关闭'),
      subTitle: i18n.t('数据导入终止提示'),
      okText: i18n.t('确定'),
      cancelText: i18n.t('取消'),
      confirmFn: () => resolve(true),
      cancelFn: () => resolve(false)
    })
  })
}

/**
 * 生成导入侧滑组件实例
 */
const createSideslider = () => {
  const Component = Vue.extend({
    components: {
      contentComponent
    },
    mounted() {
      document.body.appendChild(this.$el)
    },
    beforeDestroy() {
      try {
        // document.body.removeChild(this.$el)
      } catch (error) {
        console.error(error)
      }
    },
    render() {
      return (
        <bk-sideslider
          is-show={ state.visible }
          width={ 700 }
          title={ state.title }
          transfer={ true }
          before-close={ beforeClose }
          { ...{ on: { 'update:isShow': close } } }>
          <content-component slot="content"></content-component>
        </bk-sideslider>
      )
    }
  })
  return new Component({ store, i18n })
}

export default function () {
  return [state, { close, show, setState }]
}
