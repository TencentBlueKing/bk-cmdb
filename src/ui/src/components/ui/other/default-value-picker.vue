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
  <div class="form-default default-value-picker">
    <span>
      默认值：
    </span>
    <cmdb-property-value
      class="form-default-text"
      tag="span"
      :is-show-overflow-tips="true"
      :ref="`default-value-${property.bk_property_id}`"
      :value="value"
      :property="property"
      :instance="instance">
    </cmdb-property-value>
    <span class="form-default-operate"
      @mousedown="handlePickDefault">
      填入
    </span>
  </div>
</template>

<script>
  export default {
    props: {
      property: {
        type: [Object, String],
        default: () => ({})
      },
      value: {
        type: [String, Number, Array, Boolean, Object],
        default: ''
      },
      instance: {
        type: Object,
        default: () => ({})
      }
    },
    data() {
      return {
      }
    },
    methods: {
      handlePickDefault() {
        this.$emit('pick-default', this.value)
      }
    }
  }
</script>

<style lang="scss" scoped>
.form-default {
  display: flex;
  position: absolute;
  top: 100%;
  left: 0;
  right: 80px;
  font-size: 12px;
  .form-default-text {
    flex: 1;
    @include ellipsis;
  }
  .form-default-operate {
    color: $primaryColor;
    cursor: pointer;
  }
}
</style>
