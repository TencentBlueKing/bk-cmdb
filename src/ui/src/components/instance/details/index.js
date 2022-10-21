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

import useProperty from '@/hooks/model/property'
import useGroup from '@/hooks/model/group'
import usePending from '@/hooks/utils/pending'
import useInstance from '@/hooks/instance/instance'
import Vue, { reactive, computed, toRef } from 'vue'
import store from '@/store'
import i18n from '@/i18n'
const state = reactive({
  bk_obj_id: null,
  bk_inst_id: null,
  bk_biz_id: null,
  visible: false,
  title: null
})
const visible = toRef(state, 'visible')
const title = toRef(state, 'title')
const modelOptions = computed(() => ({ bk_obj_id: state.bk_obj_id, bk_biz_id: state.bk_biz_id }))
const instanceOptions = computed(() => ({ bk_inst_id: state.bk_inst_id, ...modelOptions.value }))
const [{ properties, pending: propertyPending }] = useProperty(modelOptions)
const [{ groups, pending: groupPneding }] = useGroup(modelOptions)
const [{ instance, pending: instancePending }] = useInstance(instanceOptions)
const pending = usePending([propertyPending, groupPneding, instancePending], true)

const createDetails = () => {
  const Component = Vue.extend({
    mounted() {
      document.body.appendChild(this.$el)
    },
    render() {
      const directives = [{ name: 'bkloading', value: { isLoading: pending.value } }]
      const close = () => {
        visible.value = false
        setTimeout(() => this.$destroy(), 200)
      }
      return (
        <bk-sideslider
          is-show={ visible.value }
          width={ 700 }
          title={ title.value }
          transfer={ true }
          { ...{ on: { 'update:isShow': close } } }>
          <cmdb-details slot="content"
            { ...{ directives } }
            show-options={ false }
            inst={ instance.value }
            properties={ properties.value }
            property-groups={ groups.value }>
          </cmdb-details>
        </bk-sideslider>
      )
    }
  })
  return new Component({ store, i18n })
}
export default function (options) {
  Object.assign(state, options)
  const details = createDetails()
  details.$mount()
  visible.value = true
}
