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

<template>
  <bk-cascade
    v-model="selected"
    :list="cascadeList"
    :scroll-width="190"
    :disabled="disabled"
    trigger="hover"
    clearable
    ext-popover-cls="clipboard-cascade-popover"
    @change="handleChange">
    <template #trigger="{ isShow }">
      <bk-button :disabled="disabled" theme="default">
        {{$t('复制')}}
        <i :class="['bk-icon', 'icon-angle-down', { 'open': isShow }]"></i>
      </bk-button>
    </template>
  </bk-cascade>
</template>

<script>
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { IPWithCloudFields } from '@/dictionary/ip-with-cloud-symbol'

  export default {
    name: 'cmdb-clipboard-selector',
    props: {
      disabled: {
        type: Boolean,
        default: false
      },
      list: {
        type: Array,
        default() {
          return []
        }
      },
      idKey: {
        type: String,
        default: 'id'
      },
      labelKey: {
        type: String,
        default: 'name'
      }
    },
    data() {
      return {
        selected: []
      }
    },
    computed: {
      cascadeList() {
        const list = this.list
          .filter(item => item.bk_property_type !== PROPERTY_TYPES.INNER_TABLE)
          .map((item) => {
            const id = item[this.idKey]
            const name = item[this.labelKey]
            return {
              id,
              name,
              property: item
            }
          })

        const IPWithCloudKeys = Object.keys(IPWithCloudFields)
        const index = list.findIndex(item => IPWithCloudKeys.includes(item.id))
        if (index !== -1) {
          // 将IPWithCloudFields中的字段提取出来，再以子选项的形式插入
          const cloudChildren = list.splice(index, IPWithCloudKeys.length)
          list.splice(index, 0, {
            id: -1,
            name: `${this.$t('管控区域')}ID:IP`,
            children: cloudChildren
          })
        }

        return list
      }
    },
    methods: {
      handleChange(newValue, oldValue, selectList) {
        const selected = selectList.at(-1)
        this.$emit('on-copy', selected.property)
        this.selected = []
      }
    }
  }
</script>

<style lang="scss" scoped>
.icon-angle-down {
  transform-origin: center;
  transition: all .2s ease;
  &.open {
    transform: rotate(180deg);
  }
}
</style>
<style lang="scss">
.clipboard-cascade-popover {
  .bk-cascade-options .bk-option-content .bk-option-name {
    display: block;
    font-size: 14px;
  }
}
</style>
