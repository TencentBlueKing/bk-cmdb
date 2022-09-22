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
  <div class="choose-icon">
    <bk-input type="text" class="system-icon-search" v-show="activeTab === 'system'"
      clearable
      right-icon="bk-icon icon-search"
      :placeholder="$t('请输入关键词')"
      v-model.trim="searchText">
    </bk-input>
    <bk-tab :active.sync="activeTab" type="unborder-card" class="icon-tab">
      <bk-tab-panel name="system" :label="$t('系统图标')">
        <icon-set v-model="curIcon" :icon-list="iconList" :filter-icon="searchText"></icon-set>
      </bk-tab-panel>
    </bk-tab>
    <div class="footer">
      <bk-button theme="primary" @click="handleConfirm">{{$t('确定')}}</bk-button>
      <bk-button @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import iconList from '@/assets/json/model-icon.json'
  import iconSet from './icon-set'
  export default {
    components: {
      iconSet
    },
    props: {
      value: {
        type: String,
        default: 'icon-cc-default'
      }
    },
    data() {
      return {
        iconList,
        activeTab: 'system',
        searchText: '',
        curIcon: this.value
      }
    },
    methods: {
      handleConfirm() {
        this.$emit('input', this.curIcon)
        this.$emit('chooseIcon')
      },
      handleCancel() {
        this.$emit('close')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .choose-icon {
        position: relative;
        height: 460px;
        overflow: hidden;
        .system-icon-search {
            position: absolute;
            top: 12px;
            right: 20px;
            width: 240px;
            z-index: 2;
        }
        .icon-tab {
            width: 100%;
            height: calc(100% - 58px);
            /deep/ .bk-tab-section {
                margin: 10px 0;
                height: calc(100% - 77px);
                @include scrollbar-y;
            }
            /deep/ .bk-tab-header {
                padding: 0;
                margin: 0 20px;
            }
        }
        .footer {
            height: 57px;
            line-height: 56px;
            text-align: right;
            font-size: 0;
            padding-right: 24px;
            background-color: #fafbfd;
            border-top: 1px solid #dcdee5;
            .bk-button {
                margin-left: 10px;
            }
        }
    }
</style>
