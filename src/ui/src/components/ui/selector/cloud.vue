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
  <bk-select
    class="cloud-selector"
    v-model="selected"
    :searchable="true"
    :clearable="allowClear"
    :disabled="disabled"
    :multiple="multiple"
    :display-tag="multiple"
    :selected-style="getSelectedStyle"
    :popover-options="popoverOptions"
  >
    <bk-option
      v-for="(option, index) in data"
      :key="index"
      :id="option.bk_cloud_id"
      :name="option.bk_cloud_name">
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    name: 'cmdb-cloud-selector',
    props: {
      value: {
        type: [Array, String, Number],
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean,
        default: false
      },
      allowClear: {
        type: Boolean,
        default: false
      },
      popoverOptions: {
        type: Object,
        default() {
          return {}
        }
      },
      requestConfig: {
        type: Object,
        default() {
          return {}
        }
      }
    },
    data() {
      return {
        data: [],
        selected: this.multiple ? [] : ''
      }
    },
    computed: {
      getSelectedStyle() {
        return this.multiple ? 'checkbox' : 'check'
      },
    },
    watch: {
      value: {
        handler(value) {
          if (value !== null) {
            this.selected = value
          }
        },
        immediate: true
      },
      selected(selected) {
        this.$emit('input', selected)
        this.$emit('on-selected', selected)
      }
    },
    created() {
      this.getData()
    },
    methods: {
      async getData() {
        const data = await this.$store.dispatch('cloud/area/findMany', {
          params: {},
          config: { ...this.requestConfig, ...{ fromCache: true } }
        })
        this.data = data.info || []
      }
    }
  }
</script>

<style lang="scss" scoped>
    .cloud-selector {
        width: 100%;
    }
</style>
