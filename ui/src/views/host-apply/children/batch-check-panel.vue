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
  <div class="batch-check-panel">
    <div class="panel-hd">
      <div class="panel-title">
        <slot name="title" />
        <a href="javascript:;" class="clear-all" @click="$emit('clear')">{{$t('清空')}}</a>
      </div>
    </div>
    <div class="panel-bd">
      <slot name="content" />
    </div>
    <div class="panel-ft">
      <bk-button theme="primary" :disabled="!checkedList.length" @click="$emit('edit')">
        {{$t(actionMode === 'batch-del' ? '去删除' : '去编辑')}}
      </bk-button>
      <bk-button theme="default" @click="$emit('cancel')">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import { defineComponent } from '@vue/composition-api'

  export default defineComponent({
    props: {
      checkedList: Array,
      actionMode: String
    }
  })
</script>

<style lang="scss" scoped>
  .batch-check-panel {
    height: 100%;

    .panel-hd,
    .panel-ft {
      background: #fafbfd;
    }
    .panel-hd {
      height: 52px;
      line-height: 52px;
      padding: 0 12px;
      border-bottom: 1px solid #dcdee5;
    }
    .panel-title {
      position: relative;
      font-size: 14px;
      color: #63656e;

      .checked-num {
        font-style: normal;
        font-weight: bold;
        color: #2dcb56;
        margin: .1em;
      }

      .clear-all {
          position: absolute;
        right: 0;
        top: 0;
        color: #3a84ff;
      }
    }
    .panel-bd {
      height: calc(100% - 52px - 60px);
      background: #fff;
      @include scrollbar-y;
    }
    .panel-ft {
      height: 60px;
      line-height: 58px;
      text-align: center;
      border-top: 1px solid #dcdee5;
      .bk-button {
        min-width: 86px;
        margin: 0 3px;
      }
    }

    .module-list {
      .module-item {
        position: relative;
        padding: 8px 42px;

        &:hover {
          background: #f0f1f5;
        }

        .module-name {
          font-size: 14px;
          color: #63656e;
        }
        .module-path {
          font-size: 12px;
          color: #c4c6cc;
          @include ellipsis;
        }
        .module-icon {
          position: absolute;
          left: 12px;
          top: 8px;
          font-size: 12px;
          border-radius: 50%;
          background-color: #c4c6cc;
          width: 22px;
          height: 22px;
          line-height: 21px;
          text-align: center;
          font-style: normal;
          color: #fff;
        }
        .action-icon {
          position: absolute;
          right: 8px;
          top: 10px;
          width: 28px;
          height: 28px;
          text-align: center;
          line-height: 28px;

          a {
            color: #c4c6cc;
            &:hover {
              color: #979ba5;
            }
          }
        }
      }
    }
  }
</style>
