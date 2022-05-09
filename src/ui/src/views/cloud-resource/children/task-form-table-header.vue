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
  <div class="custom-header">
    {{data.column.label}}
    <bk-popover class="batch-trigger"
      trigger="manual" theme="light" placement="bottom" ref="popover"
      :tippy-options="{
        hideOnClick: false,
        onHidden: handleHidden
      }">
      <i class="icon-cc-batch-update"
        v-bk-tooltips="{
          content: $t('批量编辑'),
          disabled: disabled
        }"
        :class="{ 'is-disabled': disabled }"
        @click="handleClick">
      </i>
      <bk-form form-type="vertical" slot="content">
        <bk-form-item :label="$t('批量编辑')">
          <task-directory-selector class="directory-selector"
            v-model="selected">
          </task-directory-selector>
          <div class="selector-options">
            <link-button class="selector-link-button" @click="handleConfirm">{{$t('确定')}}</link-button>
            <link-button class="selector-link-button ml10" @click="handleCancel">{{$t('取消')}}</link-button>
          </div>
        </bk-form-item>
      </bk-form>
    </bk-popover>
  </div>
</template>

<script>
  import TaskDirectorySelector from './task-directory-selector.vue'
  export default {
    name: 'task-table-header',
    components: {
      TaskDirectorySelector
    },
    props: {
      data: {
        type: Object,
        required: true
      },
      batchSelectHandler: {
        type: Function,
        required: true
      },
      disabled: Boolean
    },
    data() {
      return {
        selected: ''
      }
    },
    methods: {
      handleClick() {
        if (this.disabled) {
          return false
        }
        this.$refs.popover.instance.show()
      },
      handleConfirm() {
        this.batchSelectHandler(this.selected)
        this.handleCancel()
      },
      handleCancel() {
        this.$refs.popover.instance.hide()
      },
      handleHidden() {
        this.selected = ''
      }
    }
  }
</script>

<style lang="scss" scoped>
    .batch-trigger {
        line-height: 1;
        display: inline-block;
        vertical-align: baseline;
        [class^="icon"] {
            font-size: 14px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                color: #1964E1;
            }
            &.is-disabled {
                color: #c4c6cc;
                cursor: not-allowed;
            }
        }
    }
    .directory-selector {
        width: 235px;
    }
    .selector-options {
        text-align: right;
        font-size: 0;
        .selector-link-button {
            font-size: 12px;
        }
    }
</style>
