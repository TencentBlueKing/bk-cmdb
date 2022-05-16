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

import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import FormPropertySelector from './form-property-selector.vue'
const Component = Vue.extend({
  components: {
    FormPropertySelector
  },
  created() {
    this.unwatch = RouterQuery.watch('*', () => {
      this.handleClose()
    })
  },
  beforeDestroy() {
    this.unwatch()
  },
  methods: {
    handleClose() {
      // magicbox实现相关，多个侧滑同时存在，后面的不会挂在到body中，而是挂在到popmanager中，此处不手动移出
      // document.body.removeChild(this.$el)
      this.$destroy()
    }
  },
  render() {
    return (<form-property-selector ref="selector" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></form-property-selector>)
  }
})

export default {
  show(data = {}, dynamicGroupForm) {
    const vm = new Component({
      store,
      i18n,
      attrs: data,
      provide: () => ({
        dynamicGroupForm
      })
    })
    vm.$mount()
    // magicbox实现相关，多个侧滑同时存在，后面的不会挂在到body中，而是挂在到popmanager中，此处不手动添加
    // document.body.appendChild(vm.$el)
    vm.$refs.selector.show()
  }
}
