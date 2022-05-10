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
import i18n from '@/i18n'
import store from '@/store'
import FilterForm from './filter-form.vue'
import FilterStore from './store'
export default {
  show(props = {}) {
    const exist = FilterStore.getComponent('FilterForm')
    if (exist) {
      exist.$refs.FilterForm.focusIP()
      return exist
    }
    const vm = new Vue({
      i18n,
      store,
      created() {
        FilterStore.setComponent('FilterForm', this)
      },
      beforeDestroy() {
        FilterStore.setComponent('FilterForm', null)
      },
      render(h) {
        return h(FilterForm, {
          ref: 'FilterForm',
          props,
          on: {
            closed: () => {
              this.$el && this.$el.parentElement && this.$el.parentElement.removeChild(this.$el)
              this.$destroy()
            }
          }
        })
      }
    })
    vm.$mount()
    document.body.appendChild(vm.$el)
    vm.$refs.FilterForm.open()
    vm.$watch(() => window.CMDB_APP.$route.name, vm.$refs.FilterForm.close)
    return vm
  }
}
