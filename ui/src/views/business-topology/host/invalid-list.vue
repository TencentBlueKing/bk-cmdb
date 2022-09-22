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
  <section class="invalid-list" v-if="list.length">
    <div class="header">
      <strong class="title">{{title}}</strong>
      <bk-button class="copy" size="small" text @click="handleCopy">{{$t('复制IP')}}</bk-button>
    </div>
    <ul class="list clearfix">
      <li class="item fl"
        v-for="(item, index) in list"
        v-bk-overflow-tips="{
          interactive: false
        }"
        :key="index">
        {{item}}
      </li>
    </ul>
  </section>
</template>

<script>
  export default {
    props: {
      title: {
        type: String,
        default: ''
      },
      list: {
        type: Array,
        default: () => ([])
      }
    },
    methods: {
      async handleCopy() {
        try {
          await this.$copyText(this.list.join('\n'))
          this.$success(this.$t('复制成功'))
        } catch (error) {
          this.$error(this.$t('复制失败'))
          console.error(error)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .invalid-list {
        padding: 10px 0;
        margin: 13px 20px 0;
        background-color: #f0f1f5;
        text-align: left;
        border-radius: 2px;
        .header {
            display: flex;
            align-items: center;
            padding: 0 15px;
            .title {
                font-size: 12px;
                font-weight: 700;
                color: #63656e;
                line-height: 16px;
            }
            .copy {
                margin-left: auto;
            }
        }
        .list {
            padding: 0 0 0 15px;
            max-height:  160px;
            line-height: 20px;
            @include scrollbar-y;
            .item {
                width: 100px;
                font-size: 12px;
                margin-right: 15px;
                @include ellipsis;
            }
        }
    }
</style>
