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
  import { computed, ref } from 'vue'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'

  const props = defineProps({
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    fieldList: {
      type: Array,
      default: () => ([])
    },
    // 传入一个指定的字段，则表示以为指定的字段配置，用于创建字段模式
    field: {
      type: Object,
      default: () => ({})
    }
  })

  const type = ref('')
  const rules = ref([])

  const isUnion = computed(() => type.value === 'union')

  const singleRuleTypes = [
    PROPERTY_TYPES.SINGLECHAR,
    PROPERTY_TYPES.INT,
    PROPERTY_TYPES.FLOAT
  ]
  const unionRuleTypes = [
    PROPERTY_TYPES.SINGLECHAR,
    PROPERTY_TYPES.INT,
    PROPERTY_TYPES.FLOAT,
    PROPERTY_TYPES.ENUM,
    PROPERTY_TYPES.DATE,
    PROPERTY_TYPES.LIST
  ]
  const singleRuleProperties = computed(() => props.fieldList
    .filter(field => singleRuleTypes.includes(field.bk_property_type)))
  const unionRuleProperties = computed(() => props.fieldList
    .filter(field => unionRuleTypes.includes(field.bk_property_type)))

  console.log(singleRuleProperties, unionRuleProperties, props.field)
</script>

<template>
  <div class="unique-manage">
    <div class="unique-list">
      <div class="unique-item">
        <bk-compose-form-item class="compose-form">
          <bk-select class="rule-type" v-model="type" :clearable="false">
            <bk-option id="single" :name="$t('单独唯一')"></bk-option>
            <bk-option id="union" :name="$t('联合唯一')"></bk-option>
          </bk-select>
          <bk-select class="rule-selector" v-model="rules"
            searchable
            :multiple="isUnion">
            <bk-option v-for="property in unionRuleProperties"
              :key="property.id"
              :id="property.id"
              :name="property.bk_property_name">
            </bk-option>
          </bk-select>
        </bk-compose-form-item>
        <div class="operation">
          <i class="icon-cc-zoom-in icon-button add"></i>
          <i class="icon-cc-zoom-out icon-button remove"></i>
        </div>
      </div>
    </div>

    <div class="removed-container">
      <div class="title">{{$t('本次删除的校验')}}</div>
      <div class="unique-list">
        <div class="unique-item">
          <bk-compose-form-item class="compose-form">
            <bk-select class="rule-type" v-model="type" :clearable="false">
              <bk-option id="single" :name="$t('单独唯一')"></bk-option>
              <bk-option id="union" :name="$t('联合唯一')"></bk-option>
            </bk-select>
            <bk-select class="rule-selector" v-model="rules"
              searchable
              :multiple="isUnion">
              <bk-option v-for="property in unionRuleProperties"
                :key="property.id"
                :id="property.id"
                :name="property.bk_property_name">
              </bk-option>
            </bk-select>
          </bk-compose-form-item>
          <div class="operation">
            <i class="icon-cc-zoom-in icon-button add"></i>
            <i class="icon-cc-zoom-out icon-button remove"></i>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .unique-list {
    .unique-item {
      display: flex;
      align-items: center;
      gap: 10px;

      .compose-form {
        flex: 1;

        .rule-type {
          width: 96px
        }
        .rule-selector {
          width: calc(100% - 96px);
        }
      }

      .operation {
        display: flex;
        gap: 8px;

        .icon-button {
          cursor: pointer;
        }
      }
    }
  }

  .removed-container {
    margin-top: 40px;
    .title {
      font-size: 12px;
    }
  }
</style>
