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
  import { computed, del, reactive, ref, set, watch } from 'vue'
  import { v4 as uuidv4 } from 'uuid'
  import cloneDeep from 'lodash/cloneDeep'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { UNIUQE_TYPES } from '@/dictionary/model-constants'
  import useUnique from './use-unique'

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
      type: Object
    }
  })

  const rules = ref([])

  const uniqueListLocal = ref(cloneDeep(props.uniqueList).map(item => ({
    type: item.keys.length > 1 ? UNIUQE_TYPES.UNION : UNIUQE_TYPES.SINGLE,
    ...item
  })))
  const fieldListLocal = ref(cloneDeep(props.fieldList))

  const { getUniqueByField } = useUnique(uniqueListLocal)
  const isFieldMode = computed(() => props.field !== undefined)

  const validateResult = reactive({})
  const SET_TYPE = Object.freeze({
    subset: Symbol('subset'),
    superset: Symbol('superset')
  })

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

  const newFieldID = -1
  const newDefaultRule = () => ({
    id: uuidv4(),
    keys: [!props.field.id ? newFieldID : props.field.id],
    type: UNIUQE_TYPES.UNION
  })

  watch(() => props.field, (field) => {
    if (isFieldMode.value) {
      const fieldId = !field.id ? newFieldID : field.id
      const currentField = fieldListLocal.value.find(item => item.id === fieldId)

      // 新创建字段时
      if (!currentField) {
        // 在列表头添加当前字段
        fieldListLocal.value.unshift({
          id: newFieldID,
          bk_property_name: field.bk_property_name,
          bk_property_type: field.bk_property_type
        })
      } else {
        currentField.bk_property_name = field.bk_property_name
        currentField.bk_property_type = field.bk_property_type
      }

      // 得到当前字段的唯一校验
      const { list: fieldUniqueList } = getUniqueByField({ ...field, id: fieldId })
      if (!fieldUniqueList.length) {
        fieldUniqueList.push(newDefaultRule())
      }

      uniqueListLocal.value = fieldUniqueList
    }
  }, { deep: true, immediate: true })

  const singleRuleProperties = computed(() => fieldListLocal.value
    .filter(field => singleRuleTypes.includes(field.bk_property_type)))
  const unionRuleProperties = computed(() => fieldListLocal.value
    .filter(field => unionRuleTypes.includes(field.bk_property_type)))

  const getDisplayProperties = unqiue => (unqiue.type === UNIUQE_TYPES.UNION
    ? unionRuleProperties.value
    : singleRuleProperties.value
  )

  const allUniqueListLocal = computed(() => {
    const all = props.uniqueList.slice()
    uniqueListLocal.value.forEach((item) => {
      if (!all.some(unique => unique.id === item.id)) {
        all.push(item)
      }
    })
    return all
  })

  // 判断当前规则是否是某些规则的子集，并返回已存在的超集
  const validateSubset = (rules, id) => {
    const superset = allUniqueListLocal.value
      .filter(existRule => existRule.id !== id && rules.every(id => existRule.keys.includes(id)))
    return superset.length ? superset : false
  }
  // 判断当前规则是否是某些规则的超集，并返回已存在的子集
  const validateSuperset = (rules, id) => {
    const subset = allUniqueListLocal.value
      .filter(existRule => existRule.id !== id && existRule.keys.every(id => rules.includes(id)))
    return subset.length ? subset : false
  }

  const getValidateRuleName = (rule) => {
    const name = rule.keys.map((id) => {
      const property = fieldListLocal.value.find(property => property.id === id)
      return property ? property.bk_property_name : id
    })
    return name.join('+')
  }

  const validate = (keys, id) => {
    del(validateResult, id)
    const subset = validateSubset(keys, id)
    if (subset) {
      set(validateResult, id, {
        type: SET_TYPE.subset,
        rules: subset
      })
      return false
    }
    const superset = validateSuperset(keys, id)
    if (superset) {
      set(validateResult, id, {
        type: SET_TYPE.superset,
        rules: superset
      })
      return false
    }

    return true
  }

  const validateAll = () => {
    const result = allUniqueListLocal.value.map(({ keys, id }) => validate(keys, id))
    return result.every(valid => valid)
  }

  const handleRuleKeysChange = (keys, id) => {
    // 更新数据
    const unique = uniqueListLocal.value.find(item => item.id === id)
    unique.keys = keys

    // 触发校验
    validate(keys, id)
  }

  const xd = computed(() => uniqueListLocal.value.map(item => item.keys))

  const handleAppendRule = (index) => {
    uniqueListLocal.value.splice(index + 1, 0, newDefaultRule())
  }
  const handleRemoveRule = (index) => {
    if (isFieldMode.value && uniqueListLocal.value.length === 1) {
      return
    }

    uniqueListLocal.value.splice(index, 1)

    // 重新检查校验状态
    validateAll()
  }

  defineExpose({
    validateAll,
    uniqueList: uniqueListLocal.value
  })
</script>

<template>
  <div class="unique-manage">
    <div class="unique-list">
      <div class="unique-item" v-for="(unique, uniqueIndex) in uniqueListLocal" :key="unique.id">
        <bk-compose-form-item :class="['compose-form', { 'is-field-mode': isFieldMode }]">
          <bk-select
            class="rule-type"
            :value="unique.keys.length > 1 ? UNIUQE_TYPES.UNION : UNIUQE_TYPES.SINGLE"
            :clearable="false"
            v-if="!isFieldMode">
            <bk-option :id="UNIUQE_TYPES.SINGLE" :name="$t('单独唯一')"></bk-option>
            <bk-option :id="UNIUQE_TYPES.UNION" :name="$t('联合唯一')"></bk-option>
          </bk-select>
          <bk-select class="rule-selector"
            :value="unique.keys"
            :searchable="getDisplayProperties(unique).length > 10"
            :multiple="true"
            :clearable="false"
            @change="(keys) => handleRuleKeysChange(keys, unique.id)">
            <bk-option v-for="property in getDisplayProperties(unique)"
              :disabled="property.id === field.id || property.id === newFieldID"
              :key="property.id"
              :id="property.id"
              :name="property.bk_property_name">
            </bk-option>
          </bk-select>
        </bk-compose-form-item>
        <div class="operation">
          <i :class="['icon-cc-zoom-in', 'icon-button', 'add']"
            @click="handleAppendRule(uniqueIndex)">
          </i>
          <i :class="['icon-cc-zoom-out', 'icon-button', 'remove',
                      { disabled: uniqueListLocal.length - 1 === 0 && isFieldMode }]"
            @click="handleRemoveRule(uniqueIndex)">
          </i>
        </div>
        <div class="rules-error" v-if="validateResult[unique.id]">
          <template v-if="validateResult[unique.id].type === SET_TYPE.superset">
            <i18n path="唯一校验子集提示">
              <template #name>
                <span class="rules-error-name"
                  v-for="(rule, index) in validateResult[unique.id].rules"
                  :key="rule.id">
                  {{getValidateRuleName(rule)}}
                  <template v-if="index !== (validateResult[unique.id].rules.length - 1)">、</template>
                </span>
              </template>
            </i18n>
          </template>
          <template v-else>{{$t('唯一校验超集提示')}}</template>
        </div>
      </div>
      <div>{{ xd }}</div>
    </div>

    <div class="removed-container" v-if="!isFieldMode">
      <div class="title">{{$t('本次删除的校验')}}</div>
      <div class="unique-list">
        <div class="unique-item">
          <bk-compose-form-item class="compose-form">
            <bk-select class="rule-type" :clearable="false">
              <bk-option :id="UNIUQE_TYPES.SINGLE" :name="$t('单独唯一')"></bk-option>
              <bk-option :id="UNIUQE_TYPES.UNION" :name="$t('联合唯一')"></bk-option>
            </bk-select>
            <bk-select class="rule-selector" v-model="rules"
              searchable
              :multiple="true">
              <bk-option v-for="property in unionRuleProperties"
                :key="property.id"
                :id="property.id"
                :name="property.bk_property_name">
              </bk-option>
            </bk-select>
          </bk-compose-form-item>
          <div class="operation">
            <a href="#">恢复</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .unique-list {
    display: flex;
    gap: 24px;
    flex-direction: column;

    .unique-item {
      position: relative;
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

        &.is-field-mode {
          .rule-selector {
            width: 100%;
          }
        }
      }

      .operation {
        display: flex;
        gap: 8px;

        .icon-button {
          cursor: pointer;

          &.disabled {
            color: $grayColor;
            cursor: not-allowed;
          }
        }
      }

      .rules-error {
        font-size: 12px;
        color: $dangerColor;
        position: absolute;
        top: 100%;
        .rules-error-name {
          color: $warningColor;
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
