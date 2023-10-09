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
  <div class="search-value-container" v-if="displayType === 'info'">
    <div class="prepend"><slot name="info-prepend"></slot></div>
    <org-value :property="property" :value="localValue" show-on="search"></org-value>
  </div>
  <cmdb-form-organization
    v-else
    v-model="localValue"
    v-bind="$attrs"
    :multiple="multiple"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
  </cmdb-form-organization>
</template>

<script>
  import activeMixin from './mixins/active'
  import orgValue from '@/components/ui/other/org-value.vue'

  export default {
    name: 'cmdb-search-organization',
    components: {
      orgValue
    },
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String, Number],
        default: () => []
      },
      property: {
        type: Object,
        default: () => ({})
      },
      multiple: {
        type: Boolean,
        default: true
      },
      displayType: {
        type: String,
        default: 'selector',
        validator(type) {
          return ['selector', 'info'].includes(type)
        }
      }
    },
    computed: {
      localValue: {
        get() {
          if (this.multiple) {
            if (this.value && !Array.isArray(this.value)) {
              return [this.value]
            }
            return this.value || []
          }
          if (Array.isArray(this.value)) {
            return this.value[0] || ''
          }
          return this.value || ''
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .search-value-container {
    display: flex;
    .prepend {
      margin-right: 4px;
    }
  }
</style>
