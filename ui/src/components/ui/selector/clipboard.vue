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
  <bk-dropdown-menu trigger="click" :disabled="disabled" font-size="medium">
    <bk-button class="clipboard-trigger" theme="default" slot="dropdown-trigger" v-test-id="'copy'"
      :disabled="disabled">
      {{$t('复制')}}
      <i class="bk-icon icon-angle-down"></i>
    </bk-button>
    <ul class="clipboard-list" slot="dropdown-content" v-test-id="'copy'">
      <li v-for="(item, index) in list"
        class="clipboard-item"
        :key="index"
        @click="handleClick(item)">
        {{item[labelKey]}}
      </li>
    </ul>
  </bk-dropdown-menu>
</template>

<script>
  export default {
    name: 'cmdb-clipboard-selector',
    props: {
      disabled: {
        type: Boolean,
        default: false
      },
      list: {
        type: Array,
        default() {
          return []
        }
      },
      idKey: {
        type: String,
        default: 'id'
      },
      labelKey: {
        type: String,
        default: 'name'
      }
    },
    methods: {
      handleClick(item) {
        this.$emit('on-copy', item)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .clipboard-trigger{
        padding: 0 16px;
        .icon-angle-down {
            font-size: 20px;
            margin: 0 -4px;
        }
    }
    .clipboard-list{
        width: 100%;
        font-size: 14px;
        line-height: 32px;
        // 漏出半个 item，引导用户下拉
        max-height: calc(160px + (32px / 2));
        @include scrollbar-y;
        &::-webkit-scrollbar{
            width: 3px;
            height: 3px;
        }
        .clipboard-item{
            padding: 0 15px;
            cursor: pointer;
            @include ellipsis;
            &:hover{
                background-color: #ebf4ff;
                color: #3c96ff;
            }
        }
    }
</style>
