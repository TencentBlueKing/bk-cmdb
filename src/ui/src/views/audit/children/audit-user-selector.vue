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
  <cmdb-form-objuser
    v-model="userValue"
    :exclude="false"
    :multiple="true"
    :placeholder="$t('请输入xx', { name: $t('账号') })">
    <bk-select class="user-option" slot="prepend"
      :clearable="false"
      v-model="selectValue">
      <bk-option id="in" :name="$t('包含')"></bk-option>
      <bk-option id="not_in" :name="$t('不包含')"></bk-option>
    </bk-select>
  </cmdb-form-objuser>
</template>

<script>
  export default {
    props: {
      value: {
        type: Array,
        default: []
      }
    },
    computed: {
      userValue: {
        get() {
          return this.value[1]
        },
        set(values) {
          this.$emit('input', [this.selectValue, values])
        }
      },
      selectValue: {
        get() {
          return this.value[0]
        },
        set(value) {
          this.$emit('input', [value, this.userValue])
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .user-option {
    width: 90px;
    border-color: #c4c6cc;
    box-shadow: none;
  }
</style>
