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
import ColumnsConfig from './columns-config.vue'
import useSideslider from '@/hooks/use-sideslider'
export default {
  open({ props = {}, handler = {} }) {
    const vm = new Vue({
      i18n,
      store,
      data() {
        return {
          isShow: false,
          localSelected: []
        }
      },
      created() {
        const { beforeClose, setChanged } = useSideslider(this.localSelected)
        this.beforeClose = beforeClose
        this.setChanged = setChanged
      },
      render(h) {
        return h('bk-sideslider', {
          ref: 'sideslider',
          props: {
            title: i18n.t('列表显示属性配置'),
            width: 600,
            isShow: this.isShow,
            beforeClose: () => {
              const refColumns = this.$refs.cmdbColumnsConfig
              const { columnsChangedValues } = refColumns
              if (columnsChangedValues()) {
                this.setChanged(true)
                return this.beforeClose(() => {
                  this.isShow = false
                })
              }
              this.isShow = false
            }
          },
          on: {
            'update:isShow': (isShow) => {
              this.isShow = isShow
            },
            'animation-end': () => {
              this.$el && this.$el.parentElement && this.$el.parentElement.removeChild(this.$el)
              this.$destroy()
            }
          }
        }, [h(ColumnsConfig, {
          props,
          ref: 'cmdbColumnsConfig',
          slot: 'content',
          on: {
            cancel: () => {
              this.isShow = false
            },
            apply: (properties) => {
              this.isShow = false
              handler.apply && handler.apply(properties)
            },
            reset: () => {
              this.isShow = false
              handler.reset && handler.reset()
            }
          }
        })])
      }
    })
    vm.$mount()
    document.body.appendChild(vm.$el)
    vm.isShow = true
    return vm
  }
}
