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
  <div class="form-bool-layout">
    <span class="default">{{$t('默认值')}}</span>
    <bk-switcher
      size="small"
      theme="primary"
      :value="localValue"
      :disabled="isReadonly"
      @change="handleChange">
    </bk-switcher>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: [String, Boolean],
        default: false
      },
      isReadonly: Boolean
    },
    data() {
      return {
        localValue: false
      }
    },
    watch: {
      value: {
        immediate: true,
        handler(value) {
          this.localValue = typeof value === 'boolean' ? value : false
          // 将空字符转为false
          this.handleChange(this.localValue)
        }
      }
    },
    methods: {
      handleChange(selected) {
        this.$emit('input', selected)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-bool-layout {
        .default {
            display: block;
            line-height: 36px;
            font-size: 14px;
        }
    }
</style>
