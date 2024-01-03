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
  <bk-select v-if="displayType === 'selector'"
    searchable
    v-model="localValue"
    v-bind="$attrs"
    :multiple="multiple"
    display-tag
    selected-style="checkbox"
    :loading="$loading(requestId)"
    :remote-method="searchArea"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
    <bk-option v-for="option in options"
      :key="option.bk_cloud_id"
      :id="option.bk_cloud_id"
      :name="`${option.bk_cloud_name}[${option.bk_cloud_id}]`">
    </bk-option>
  </bk-select>
  <span v-else>
    <slot name="info-prepend"></slot>
    {{info}}
  </span>
</template>

<script>
  import activeMixin from './mixins/active'
  import debounce from 'lodash.debounce'

  export default {
    name: 'cmdb-search-foreignkey',
    mixins: [activeMixin],
    props: {
      value: {
        type: [String, Array, Number],
        default: () => ([])
      },
      displayType: {
        type: String,
        default: 'selector',
        validator(type) {
          return ['selector', 'info'].includes(type)
        }
      }
    },
    data() {
      return {
        options: [],
        requestId: 'searchForeignkey'
      }
    },
    computed: {
      multiple() {
        return Array.isArray(this.value)
      },
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      },
      info() {
        const values = Array.isArray(this.value) ? this.value : [this.value]
        const info = []
        values.forEach((value) => {
          const data = this.options.find(data => data.bk_cloud_id === value)
          data && info.push(data.bk_cloud_name)
        })
        return info.join(' | ')
      }
    },
    created() {
      this.getCloudArea()
      this.searchArea =  debounce(this.getCloudArea, 300)
    },
    methods: {
      async getCloudArea(key) {
        const condition = {}
        if (key) {
          condition.$or = []
          condition.$or.push({
            bk_cloud_name: {
              $regex: key
            }
          })
          if (!isNaN(Number(key))) {
            condition.$or.push({
              bk_cloud_id: Number(key)
            })
          }
        }
        try {
          const { info } = await this.$store.dispatch('cloud/area/findMany', {
            params: {
              condition,
              page: {
                sort: 'bk_cloud_name'
              }
            },
            config: {
              requestId: this.requestId
            }
          })
          this.options = info
        } catch (error) {
          console.error(error)
        }
      },
    }
  }
</script>
