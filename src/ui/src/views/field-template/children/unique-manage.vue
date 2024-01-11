<!-- eslint-disable no-unused-vars -->
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
  import { computed, del, toRefs, reactive, ref, set, watch, getCurrentInstance, nextTick } from 'vue'
  import { v4 as uuidv4 } from 'uuid'
  import cloneDeep from 'lodash/cloneDeep'
  import { t } from '@/i18n'
  import { $error } from '@/magicbox'
  import { UNIUQE_TYPES } from '@/dictionary/model-constants'
  import MiniTag from '@/components/ui/other/mini-tag.vue'
  import useUnique, { wrapData, isUniqueExist, singleRuleTypes, unionRuleTypes } from './use-unique'

  const props = defineProps({
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    fieldList: {
      type: Array,
      default: () => ([])
    },
    beforeUniqueList: {
      type: Array,
      default: () => ([])
    },
    // 传入一个指定的字段，则表示以为指定的字段配置，用于创建字段模式
    field: {
      type: Object
    },
    type: {
      type: String,
      default: 'single'
    }
  })

  const emit = defineEmits(['update-unique'])

  const currentInstance = getCurrentInstance().proxy

  const { beforeUniqueList: oldUniqueList } = toRefs(props)

  const uniqueListLocal = ref(cloneDeep(props.uniqueList).map(wrapData))
  const fieldListLocal = ref(cloneDeep(props.fieldList))

  const { getUniqueByField, uniqueStatus, removedUniqueList } = useUnique(oldUniqueList, uniqueListLocal)
  const isFieldMode = computed(() => props.field !== undefined)

  const validateResult = reactive({})
  const SET_TYPE = Object.freeze({
    subset: Symbol('subset'),
    superset: Symbol('superset')
  })

  const newFieldID = -1
  const newDefaultRule = (autoSelect = true) => {
    const type = isFieldMode.value ? UNIUQE_TYPES.UNION : UNIUQE_TYPES.SINGLE
    const key = !props.field?.id ? newFieldID : props.field.id
    let keys = type === UNIUQE_TYPES.UNION ? [key] : key
    if (!autoSelect) {
      keys = type === UNIUQE_TYPES.UNION ? [] : ''
    }
    return {
      id: uuidv4(),
      type,
      keys
    }
  }

  watch(() => props.field, async (field) => {
    if (isFieldMode.value) {
      const fieldId = !field?.id ? newFieldID : field.id
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
    // 全局视图模式下所有的规则就是uniqueListLocal
    if (!isFieldMode.value) {
      return uniqueListLocal.value
    }

    // 字段模式下，uniqueListLocal只会是当前字段的列表，所以要计算出一个总的列表，用于校验
    const all = props.uniqueList.slice()
    uniqueListLocal.value.forEach((item) => {
      const index = all.findIndex(unique => unique.id === item.id)
      // 当修改的是已经存在的数据all中还是原本的老数据，这里替换为最新的数据
      if (~index) {
        all[index] = item
      } else {
        all.push(item)
      }
    })
    return all
  })

  // 判断当前规则是否是某些规则的子集，并返回已存在的超集
  const validateSubset = (rules, id) => {
    const superset = allUniqueListLocal.value
      .filter(existRule => existRule.id !== id
        && existRule.keys.length > 0
        && rules.length > 0
        && rules.every(id => existRule.keys.includes(id)))
    return superset.length ? superset : false
  }
  // 判断当前规则是否是某些规则的超集，并返回已存在的子集
  const validateSuperset = (rules, id) => {
    const subset = allUniqueListLocal.value
      .filter(existRule => existRule.id !== id
        && existRule.keys.length > 0
        && rules.length > 0
        && existRule.keys.every(id => rules.includes(id)))
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

  const isValid = async () => {
    const keysValid = Object.keys(validateResult).length === 0
    const selectValid = await currentInstance.$validator.validateAll()
    return keysValid && selectValid
  }

  const editDraft = reactive({})

  const getUniqueValue = (unique) => {
    if (unique.type === UNIUQE_TYPES.UNION || (isFieldMode.value && props.type === UNIUQE_TYPES.UNION)) {
      return unique.keys
    }
    return unique.keys?.[0]
  }

  const isMultiple = unique => (isFieldMode.value
    ? props.type === UNIUQE_TYPES.UNION : unique.type === UNIUQE_TYPES.UNION)

  const handleRuleKeysChange = (keys, id) => {
    // 更新数据
    const unique = uniqueListLocal.value.find(item => item.id === id)
    unique.keys = Array.isArray(keys) ? keys : [keys]

    // 触发校验
    validateAll()

    nextTick(async () => {
      emit('update-unique', uniqueListLocal.value, await isValid())
    })
  }

  const handleRuleTypeChange = (type, id) => {
    const unique = uniqueListLocal.value.find(item => item.id === id)
    const oldDraftId = `${unique.type}-${id}`
    const targetDraftId = `${type}-${id}`
    set(editDraft, oldDraftId, cloneDeep(unique))

    unique.type = type
    if (editDraft[targetDraftId]) {
      unique.keys = editDraft[targetDraftId].keys
    } else {
      unique.keys = type === UNIUQE_TYPES.UNION ? [] : ''
    }

    validateAll()

    nextTick(async () => {
      emit('update-unique', uniqueListLocal.value, await isValid())
    })
  }

  const handleAppendRule = async (index = -1) => {
    uniqueListLocal.value.splice(index + 1, 0, newDefaultRule(isFieldMode.value))

    // 字段模式中添加因为会默认选中需要即时验证一次
    if (isFieldMode.value) {
      validateAll()
    }

    emit('update-unique', uniqueListLocal.value, await isValid())
  }
  const handleRemoveRule = async (id) => {
    if (isFieldMode.value && uniqueListLocal.value.length === 1) {
      return
    }

    const index = uniqueListLocal.value.findIndex(item => item.id === id)
    uniqueListLocal.value.splice(index, 1)
    del(validateResult, id)

    // 重新检查校验状态
    validateAll()

    emit('update-unique', uniqueListLocal.value, await isValid())
  }

  const handleRecover = async (removedUnique) => {
    if (!validate(removedUnique.keys, removedUnique.id)) {
      return
    }
    if (isUniqueExist(removedUnique, uniqueListLocal.value)) {
      $error(t('唯一校验已在模板中存在，无法恢复'))
      return
    }
    const oriUnique = oldUniqueList.value.find(item => item.id === removedUnique.id)
    uniqueListLocal.value.push(wrapData(oriUnique))

    emit('update-unique', uniqueListLocal.value, await isValid())
  }

  defineExpose({
    isValid,
    getUniqueList: () => uniqueListLocal.value
  })
</script>

<template>
  <div class="unique-manage">
    <div class="unique-list">
      <div class="unique-item" v-for="(unique, uniqueIndex) in uniqueListLocal" :key="unique.id">
        <bk-compose-form-item :class="['compose-form', { 'is-field-mode': isFieldMode }]">
          <bk-select v-if="!isFieldMode"
            class="rule-type"
            :value="unique.type"
            :clearable="false"
            @change="(type) => handleRuleTypeChange(type, unique.id)">
            <bk-option :id="UNIUQE_TYPES.SINGLE" :name="$t('单独唯一')"></bk-option>
            <bk-option :id="UNIUQE_TYPES.UNION" :name="$t('联合唯一')"></bk-option>
          </bk-select>
          <bk-select class="rule-selector"
            :key="unique.type"
            :value="getUniqueValue(unique)"
            :searchable="getDisplayProperties(unique).length > 10"
            :multiple="isMultiple(unique)"
            :display-tag="isMultiple(unique)"
            :selected-style="isMultiple(unique) ? 'checkbox' : 'check'"
            :clearable="false"
            :data-vv-name="`ruleSelector${unique.id}`"
            data-vv-validate-on="change"
            v-validate="`minSelectLength:2,${$t('需要至少两个字段组成联合唯一')}`"
            @change="(keys) => handleRuleKeysChange(keys, unique.id)">
            <bk-option v-for="property in getDisplayProperties(unique)"
              :disabled="(isFieldMode && property.id === field.id) || property.id === newFieldID"
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
            @click="handleRemoveRule(unique.id)">
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
        <div class="rules-error" v-else-if="errors.has(`ruleSelector${unique.id}`)">
          {{errors.first(`ruleSelector${unique.id}`)}}
        </div>
        <template v-if="!isFieldMode">
          <mini-tag :text="$t('更新')" theme="changed" class="status-tag" v-if="uniqueStatus[unique.id].changed" />
          <mini-tag :text="$t('新增')" theme="new" class="status-tag" v-if="uniqueStatus[unique.id].new" />
        </template>
      </div>
    </div>

    <bk-exception type="empty" scene="part" class="empty" v-if="!isFieldMode && !uniqueListLocal.length">
      <i18n path="暂无XXX，LINK">
        <template #text>{{ $t('唯一校验') }}</template>
        <template #link>
          <bk-link theme="primary" @click="handleAppendRule">{{ $t('立即添加') }}</bk-link>
        </template>
      </i18n>
    </bk-exception>

    <div class="removed-container" v-if="!isFieldMode && removedUniqueList.length">
      <div class="title">{{$t('本次删除的校验')}}</div>
      <div class="unique-list">
        <div class="unique-item" v-for="unique in removedUniqueList" :key="unique.id">
          <bk-compose-form-item class="compose-form">
            <bk-select
              class="rule-type"
              :value="unique.type"
              :readonly="true"
              :clearable="false">
              <bk-option :id="UNIUQE_TYPES.SINGLE" :name="$t('单独唯一')"></bk-option>
              <bk-option :id="UNIUQE_TYPES.UNION" :name="$t('联合唯一')"></bk-option>
            </bk-select>
            <bk-select class="rule-selector"
              searchable
              :value="unique.keys"
              :readonly="true"
              :multiple="true">
              <bk-option v-for="property in getDisplayProperties(unique)"
                :key="property.id"
                :id="property.id"
                :name="property.bk_property_name">
              </bk-option>
            </bk-select>
          </bk-compose-form-item>
          <div class="operation">
            <bk-link class="link" theme="primary" @click="handleRecover(unique)">{{ $t('恢复') }}</bk-link>
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
          <mini-tag :text="$t('删除')" theme="removed" class="status-tag" v-if="uniqueStatus[unique.id].removed" />
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
        width: calc(100% - 50px);

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

      .status-tag {
        position: absolute;
        left: 0;
        top: 0;
        transform: translate(-40%, -50%);
        z-index: 1;
      }
    }
  }

  .removed-container {
    margin-top: 40px;
    .title {
      font-size: 12px;
      margin-bottom: 14px;
    }

    .unique-item {
      .compose-form {
        opacity: .5;
      }
      .operation {
        .link {
          :deep(.bk-link-text) {
            font-size: 12px;
          }
        }
      }
    }
  }
</style>
