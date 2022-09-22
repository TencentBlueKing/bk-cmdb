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
  <bk-input class="cmdb-form-input" type="text" v-bind="$attrs" v-model.trim="localValue"></bk-input>
</template>

<script>
  export default {
    name: 'cmdb-form-bool-input',
    // eslint-disable-next-line
        props: ['value'],
    data() {
      return {
        localValue: ''
      }
    },
    watch: {
      value(value) {
        this.localValue = value
      },
      localValue(localValue) {
        if (['true', 'false'].includes(localValue)) {
          localValue = [true, false].find(value => value.toString() === localValue)
        }
        this.$emit('input', localValue)
        this.$emit('on-change', localValue)
      }
    },
    created() {
      this.localValue = this.value
    },
    methods: {
      focus() {
        this.$el.querySelector('input').focus()
      }
    }
  }
</script>
