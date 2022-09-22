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
  <div class="cmdb-input">
    <bk-input :class="['cmdb-form-input', { 'has-icon': !!icon }]" type="text"
      v-model="localValue"
      :placeholder="placeholder"
      v-bind="$attrs"
      @enter="handleEnter">
    </bk-input>
    <i :class="[icon, 'input-icon']" v-if="icon" @click="handleIconClick"></i>
  </div>
</template>

<script>
  export default {
    name: 'cmdb-input',
    props: {
      value: {
        type: String,
        default: ''
      },
      icon: {
        type: String,
        default: ''
      },
      placeholder: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        localValue: this.value
      }
    },
    watch: {
      value(value) {
        this.localValue = value
      },
      localValue(localValue) {
        this.$emit('input', localValue)
      }
    },
    methods: {
      handleEnter() {
        this.$emit('enter', this.localValue)
      },
      handleIconClick() {
        this.$emit('icon-click', this.localValue)
      },
      focus() {
        this.$el.querySelector('input').focus()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .cmdb-input {
        position: relative;
        @include inlineBlock;
        .input-icon {
            position: absolute;
            font-size: 14px;
            right: 11px;
            top: 10px;
            cursor: pointer;
        }
    }
</style>
