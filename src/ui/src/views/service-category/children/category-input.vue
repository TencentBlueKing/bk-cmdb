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
  <div class="cagetory-input" v-click-outside="handleCancel">
    <bk-input type="text" v-test-id
      :ref="inputRef"
      :placeholder="placeholder"
      v-model="localValue"
      @enter="handleConfirm">
    </bk-input>
    <div class="operation">
      <span class="text-primary btn-confirm" v-test-id.businessServiceCategory="'btnConfirm'"
        @click.stop="handleConfirm">{{$t('确定')}}
      </span>
      <span class="text-primary" @click="handleCancel"
        v-test-id.businessServiceCategory="'btnCancel'">{{$t('取消')}}</span>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      placeholder: {
        type: String,
        default: ''
      },
      inputRef: {
        type: String,
        default: ''
      },
      editId: {
        type: Number,
        default: 0
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
      handleConfirm() {
        this.$emit('on-confirm', this.localValue, this.editId)
      },
      handleCancel() {
        this.$emit('on-cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .cagetory-input {
        @include space-between;
        width: 100%;
        font-weight: normal;
        .bk-form-control {
            flex: 1;
            font-size: 0;
            height: 32px;
            line-height: 32px;
            margin-right: 10px;
            /deep/ .bk-form-input {
                font-size: 14px;
                line-height: 32px;
            }
        }
        .text-primary {
            display: inline-block;
            line-height: normal;
            font-size: 12px;
            &.btn-confirm {
                position: relative;
                margin-right: 6px;
                &::after {
                    content: '';
                    position: absolute;
                    top: 3px;
                    right: -6px;
                    display: inline-block;
                    width: 1px;
                    height: 14px;
                    background-color: #dcdee5;
                }
            }
        }
    }
</style>
