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
  <div class="property-form-element">
    <component
      :is="`cmdb-form-${property.bk_property_type}`"
      :class="['form-element-item', property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
      :unit="property.unit"
      :options="property.option || []"
      :data-vv-name="property.bk_property_id"
      :data-vv-as="property.bk_property_name"
      :placeholder="$t('请输入xx', { name: property.bk_property_name })"
      :auto-check="false"
      :disabled="!property.host_apply_enabled || property.__extra__.ignore"
      v-bind="$tools.getValidateEvents(property)"
      v-validate="$tools.getValidateRules(property)"
      v-model.trim="property.__extra__.value"
    >
    </component>
    <div class="form-error"
      v-if="errors.has(property.bk_property_id)">
      {{errors.first(property.bk_property_id)}}
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      property: {
        type: Object,
        required: true,
        default: () => ({})
      }
    },
    watch: {
      // eslint-disable-next-line object-shorthand
      'property.__extra__.value': async function (value) {
        // eslint-disable-next-line no-underscore-dangle
        this.property.__extra__.valid = await this.$validator.validate()
        this.$emit('value-change', value)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .property-form-element {
        .form-error {
            margin-top: 2px;
            font-size: 12px;
            color: $cmdbDangerColor;
        }
        .form-element-item {
            &.cmdb-search-input {
                /deep/ .search-input-wrapper {
                    position: relative;
                }
            }
        }
    }
</style>
