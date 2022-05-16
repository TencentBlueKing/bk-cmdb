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
  <button class="bk-button bk-default icon-button" @click="handleClick($event)">
    <i :class="['icon-wrapper', icon]"></i>
  </button>
</template>

<script>
  export default {
    name: 'icon-button',
    props: {
      disabled: Boolean,
      icon: {
        type: String,
        required: true
      }
    },
    methods: {
      handleClick(event) {
        this.$emit('click', event)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .icon-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        line-height: 30px;
        padding: 0;
        font-size: 0px;
        .icon-wrapper {
            font-size: 0;
            &:before {
                font-size: 14px;
            }
        }
    }
</style>
