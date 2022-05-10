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
import ExportContent from './content.vue'
import RouterQuery from '@/router/query'
import $http from '@/api'
import { $info } from '@/magicbox'
const Component = Vue.extend({
  components: {
    ExportContent
  },
  data() {
    return {
      visible: false
    }
  },
  created() {
    this.unwatch = RouterQuery.watch('*', () => {
      this.close()
    })
  },
  beforeDestroy() {
    this.unwatch()
  },
  methods: {
    toggle(visible) {
      this.visible = visible
    },
    close() {
      document.body.removeChild(this.$el)
      this.$destroy()
    }
  },
  render() {
    return (
      <bk-dialog
        width={ 768 }
        value={ this.visible }
        draggable={ false }
        close-icon={ false }
        show-footer={ false }
        transfer={ false }
        on-change={ this.toggle }
        on-after-leave={ this.close }>
        <export-content { ...{ props: this.$options.attrs }} on-close={ () => this.toggle(false) }></export-content>
      </bk-dialog>
    )
  }
})

export default function (data = {}) {
  const props = {
    limit: 10000,
    ...data
  }
  if (props.count <= props.limit) {
    const options = props.options({
      start: 0,
      limit: props.limit
    });
    (async () => {
      const message = $info(i18n.t('正在导出'), 0)
      try {
        await $http.download(options)
      } catch (error) {
        console.error(error)
      } finally {
        message.close()
      }
    })()
  } else {
    const vm = new Component({
      store,
      i18n,
      attrs: props
    })
    vm.$mount()
    document.body.appendChild(vm.$el)
    vm.toggle(true)
  }
}
