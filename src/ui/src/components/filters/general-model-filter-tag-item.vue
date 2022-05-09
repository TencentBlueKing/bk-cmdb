<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script>
  import Vue from 'vue'
  import i18n from '@/i18n'
  import store from '@/store'
  import Tippy from 'bk-magic-vue/lib/utils/tippy'
  import FilterTagItem from './filter-tag-item.vue'
  import GeneralModelFilterTagForm from './general-model-filter-tag-form.vue'
  import { clearOneSearchQuery } from './general-model-filter.js'

  export default {
    name: 'general-model-filter-tag-item',
    extends: FilterTagItem,
    inject: ['condition'],
    computed: {
      dynamicCondition() {
        return this.condition()
      }
    },
    methods: {
      handleClick() {
        if (this.tagFormInstance) {
          this.tagFormInstance.show()
        } else {
          const self = this
          this.tagFormViewModel = new Vue({
            i18n,
            store,
            render(h) {
              return h(GeneralModelFilterTagForm, {
                ref: 'filterTagForm',
                props: {
                  property: self.property,
                  condition: self.dynamicCondition
                },
                on: {
                  confirm: self.handleHideTagForm,
                  cancel: self.handleHideTagForm
                }
              })
            }
          })
          this.tagFormViewModel.$mount()
          this.tagFormInstance = this.$bkPopover(this.$el, {
            content: this.tagFormViewModel.$el,
            theme: 'light',
            allowHTML: true,
            placement: 'bottom',
            trigger: 'manual',
            interactive: true,
            arrow: true,
            zIndex: window.__bk_zIndex_manager.nextZIndex(), // eslint-disable-line no-underscore-dangle
            onHide: () => !this.tagFormViewModel.$refs.filterTagForm.active,
            onHidden: () => {
              this.tagFormViewModel.$refs.filterTagForm.resetCondition()
            }
          })
          this.tagFormInstance.show()
        }
        Tippy.hideAll({ exclude: this.tagFormInstance })
      },
      handleRemove() {
        clearOneSearchQuery(this.property, this.operator)
      }
    }
  }
</script>
