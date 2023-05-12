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

<script setup>
  const props = defineProps({
    field: {
      type: Object
    }
  })

  console.log(props)
  const emit = defineEmits(['click-field'])

  const handleClickRemove = () => {}

  const handleClickField = (field) => {
    emit('click-field', field)
  }
</script>

<template>
  <div class="field-card" @click="handleClickField(field)">
    <span class="drag-icon"></span>
    <div class="field-icon">
      <i class="icon-cc icon-duanzifu"></i>
    </div>
    <div class="field-info">
      <div class="field-name-area">
        <span class="field-name" :title="field.bk_property_name">{{ field.bk_property_name }}</span>
        <span class="field-required">*</span>
        <i class="bk-icon icon-exclamation-circle-shape conflict-icon"></i>
      </div>
      <div class="field-id-area">
        <span class="field-id">{{ field.bk_property_id }}</span>
      </div>
    </div>
    <div class="tags">
      <!-- <span class="tag unique"><em class="tag-text">单独唯一</em></span> -->
      <span class="tag unique union"><em class="tag-text">联合唯一</em></span>
      <span class="tag template"><em class="tag-text">模板</em></span>
    </div>
    <bk-icon :title="$t('移除')" type="delete" class="del-icon"
      @click="handleClickRemove()" />
  </div>
</template>

<style lang="scss" scoped>
  .field-card {
    display: flex;
    align-items: center;
    position: relative;
    height: 60px;
    border: 1px solid transparent;
    background: #FFFFFF;
    box-shadow: 0 2px 4px 0 #1919290d;
    border-radius: 2px;
    padding: 0 12px 0 0;
    user-select: none;
    cursor: pointer;

    .field-icon {
      color: #989CA8;
      width: 24px;
      height: 24px;
      font-size: 22px;
      text-align: center;
      line-height: 24px;
    }

    .field-info {
      flex: 1;
      margin-left: 12px;
      width: 0;
    }

    .field-name-area {
      display: flex;
      gap: 4px;
      position: relative;
      top: 4px;
    }

    .field-name {
      font-size: 12px;
      @include ellipsis;
    }
    .field-id {
      font-size: 12px;
      color: #C4C6CC;
    }
    .field-required {
      font-size: 12px;
      color: #FF5656;
    }

    .tags {
      display: flex;
      position: absolute;
      gap: 2px;
      right: 0px;
      top: 0px;

      .tag {
        background: #E4FAF0;
        border-radius: 2px;
        padding: 1px 2px;
        height: 16px;
        line-height: 16px;

        .tag-text {
          display: block;
          font-size: 12px;
          font-style: normal;
          transform: scale(.875);
        }

        &.template {
          color: #14A568;
          background: #E4FAF0;
          border-radius: 2px;
        }

        &.unique {
          color: #3A84FF;
          background: #F0F5FF;
          border-radius: 2px;

          &.union {
            cursor: pointer;
            &:hover {
              background: #A3C5FD;
            }
          }
        }
      }
    }

    .drag-icon {
      @include dragIcon;
      visibility: hidden;
      margin: 0 4px;
    }

    .del-icon {
      visibility: hidden;
      cursor: pointer;
    }

    .conflict-icon {
      font-size: 14px;
      color: $dangerColor;
    }

    &:hover {
      background: #F0F5FF;
      border: 1px solid #3A84FF;
      .drag-icon {
        visibility: visible;
      }
      .del-icon {
        visibility: visible;
      }
    }
  }
</style>
