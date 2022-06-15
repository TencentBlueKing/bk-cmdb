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
  <div :class="['property-form-element', 'cmdb-form-item', { 'is-error': errors.has(property.bk_property_id) }]">
    <component
      :ref="`component-${property.bk_property_id}`"
      :is="`cmdb-form-${property.bk_property_type}`"
      :class="['form-element-item', property.bk_property_type]"
      :unit="property.unit"
      :options="property.option || []"
      :data-vv-name="property.bk_property_id"
      :data-vv-as="property.bk_property_name"
      :placeholder="getPlaceholder(property)"
      :auto-check="autoCheck"
      :disabled="disabled"
      v-bind="$tools.getValidateEvents(property)"
      v-validate="getValidateRules(property)"
      v-on="events"
      v-model.trim="localValue">
    </component>
    <div class="form-error"
      v-if="errors.has(property.bk_property_id)">
      {{errors.first(property.bk_property_id)}}
    </div>
  </div>
</template>

<script>
  import Utils from '@/components/filters/utils'

  export default {
    props: {
      property: {
        type: Object,
        required: true,
        default: () => ({})
      },
      value: {
        type: [String, Array, Boolean, Number],
        default: ''
      },
      autoCheck: {
        type: Boolean,
        default: true
      },
      disabled: {
        type: Boolean,
        default: false
      },
      required: {
        type: Boolean,
        default: false
      },
      events: {
        type: Object,
        default: () => ({})
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value
        },
        async set(value) {
          this.$emit('input', value)
          this.$emit('change', value, this.property)
        }
      }
    },
    methods: {
      getValidateRules(property) {
        const rules = this.$tools.getValidateRules(property)
        rules.required = this.required
        return rules
      },
      getPlaceholder(property) {
        return Utils.getPlaceholder(property)
      }
    }
  }
</script>

<style lang="scss" scoped>
  .property-form-element {
    // 重置 .cmdb-form-item display
    display: block;

    .form-element-item {
      &.cmdb-search-input {
        /deep/ .search-input-wrapper {
          position: relative;
        }
      }
    }
  }
</style>
