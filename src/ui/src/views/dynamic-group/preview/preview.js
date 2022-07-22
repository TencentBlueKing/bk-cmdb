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
import DynamicGroupPreview from './preview.vue'
const Component = Vue.extend({
  components: {
    DynamicGroupPreview
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
      document.body.removeChild(this.$el)
      this.$destroy()
    }
  },
  render() {
    return (<dynamic-group-preview ref="preview" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></dynamic-group-preview>)
  }
})

export default {
  show(data = {}) {
    const vm = new Component({
      store,
      i18n,
      attrs: data
    })
    vm.$mount()
    document.body.appendChild(vm.$el)
    vm.$refs.preview.show()
  }
}
