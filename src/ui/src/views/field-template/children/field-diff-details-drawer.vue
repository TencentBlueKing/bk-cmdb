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
  import { computed, ref, watchEffect } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { t } from '@/i18n'
  import { getPropertyDefaultValue } from '@/utils/tools'
  import FlexRow from '@/components/ui/other/flex-row.vue'
  import { DIFF_TYPES } from './use-field'
  import { PROPERTY_TYPES, PROPERTY_TYPE_NAMES } from '@/dictionary/property-constants'
  import routerActions from '@/router/actions'
  import fieldTemplateService from '@/service/field-template'
  import {
    MENU_MODEL_FIELD_TEMPLATE
  } from '@/dictionary/menu-symbol'

  const props = defineProps({
    show: {
      type: Boolean,
      default: false
    },
    diffType: {
      type: String,
      default: ''
    },
    fieldDiff: {
      type: Object,
      default: () => ({})
    },
    beforeField: {
      type: Object,
      default: () => ({})
    },
    afterField: {
      type: Object,
      default: () => ({})
    }
  })

  const emit = defineEmits(['toggle'])

  const numberLikes = [PROPERTY_TYPES.INT, PROPERTY_TYPES.FLOAT]
  const stringLikes = [PROPERTY_TYPES.SINGLECHAR, PROPERTY_TYPES.LONGCHAR]
  const listLikes = [PROPERTY_TYPES.LIST]
  const enumLikes = [PROPERTY_TYPES.ENUM, PROPERTY_TYPES.ENUMMULTI]
  const optionTypes = [...numberLikes, ...stringLikes, ...listLikes, ...enumLikes]

  const isShow = computed({
    get() {
      return props.show
    },
    set(val) {
      emit('toggle', val)
    }
  })

  const modelBeforeField = ref(cloneDeep(props.beforeField))
  const boundTemplate = ref({})
  const isTemplateBindConflict = computed(() => props.diffType === DIFF_TYPES.CONFLICT
    && props.fieldDiff?.data?.bk_template_id !== 0)

  watchEffect(async () => {
    modelBeforeField.value = props.diffType === DIFF_TYPES.NEW
      ? { bk_property_type: props.beforeField.bk_property_type }
      : cloneDeep(props.beforeField)
    if (isTemplateBindConflict?.value) {
      boundTemplate.value = await getBoundTemplate()
    }
  })

  const getBoundTemplate = async () => {
    const { bk_template_id: bkTemplateId, id, bk_obj_id: modelId } = modelBeforeField.value
    if (!id || !bkTemplateId) return {}
    const params = {
      bk_template_id: bkTemplateId,
      bk_attribute_id: id
    }
    return await fieldTemplateService.getFieldBindTemplate(params, {
      requestId: `${modelId}_${id}_${bkTemplateId}`,
      fromCache: true
    })
  }

  const isFieldIdSame = computed(() => props.afterField.bk_property_id === props.beforeField.bk_property_id)
  const isFieldNameSame = computed(() => props.afterField.bk_property_name === props.beforeField.bk_property_name)
  const isFieldTypeSame = computed(() => props.afterField.bk_property_type === props.beforeField.bk_property_type)

  const diffTitle = computed(() => {
    const titles = {
      new: t('新增字段'),
      update: t('字段配置更新'),
      conflict: isTemplateBindConflict.value ? t('字段冲突，该字段已经被其他模板绑定，请删除该模型或修改模板') : '',
      unbound: t('不再纳管该字段'),
      unchanged: t('无配置变化')
    }
    let title = titles[props.diffType]

    // 只有bk_template_id的key变更，使用指定title
    const updateDataKeys = Object.keys(props.fieldDiff?.update_data || {})
    if (props.diffType === DIFF_TYPES.UPDATE && updateDataKeys.length === 1 && updateDataKeys?.[0] === 'bk_template_id') {
      title = t('配置一致，配置被模板进行纳管')
    }

    if (props.diffType === DIFF_TYPES.CONFLICT && !isTemplateBindConflict.value) {
      if (isFieldIdSame.value && !isFieldTypeSame.value) {
        title = t('字段冲突，该字段与模板字段的字段类型不一致')
      }
      if (isFieldNameSame.value && !isFieldIdSame.value) {
        title = t('字段冲突，模板字段与模型字段的字段名称相同，但字段ID不一致')
      }
    }
    return title
  })

  const getOptionKeyName = (field) => {
    const { bk_property_type: fieldType } = field
    if (numberLikes.includes(fieldType)) {
      return t('数值范围')
    }
    if (stringLikes.includes(fieldType)) {
      return t('正则校验')
    }
    if (listLikes.includes(fieldType)) {
      return t('列表值')
    }
    if (enumLikes.includes(fieldType)) {
      return t('枚举值')
    }
  }

  const getDefaultValue = (field) => {
    const { bk_property_type: fieldType } = field
    if ([PROPERTY_TYPES.ENUMQUOTE].includes(fieldType)) {
      return  field.option.map(item => item.bk_inst_id)
    }

    if ([PROPERTY_TYPES.ENUM, PROPERTY_TYPES.ENUMMULTI].includes(fieldType)) {
      const defaultOptions = (field.option || []).filter(option => option.is_default)
      return defaultOptions.map(option => option.name).join(', ')
    }

    return getPropertyDefaultValue(field)
  }

  const getAfterDiffItemClass = (key, defaults = []) => {
    const classList = ['diff-item', ...defaults]

    const updateData = props.fieldDiff?.update_data || {}
    if (props.diffType === DIFF_TYPES.UPDATE && Object.keys(updateData).includes(key)) {
      classList.push('is-update')
    }

    if (props.diffType === DIFF_TYPES.CONFLICT && props.fieldDiff?.data.bk_template_id === 0) {
      if (key === 'bk_property_id' && isFieldNameSame.value && !isFieldIdSame.value) {
        classList.push('is-conflict')
      }
      if (key === 'bk_property_type' && isFieldIdSame.value && !isFieldTypeSame.value) {
        classList.push('is-conflict')
      }
    }

    return classList
  }

  const handleOtherTemplate = () => {
    const { id } = boundTemplate.value
    if (!id) return
    routerActions.open({
      name: MENU_MODEL_FIELD_TEMPLATE,
      query: {
        id,
        action: 'view',
        tab: 'model'
      }
    })
  }
</script>

<template>
  <bk-sideslider
    ref="sidesliderComp"
    v-transfer-dom
    :width="880"
    :title="`【${modelBeforeField?.bk_property_name ?? afterField.bk_property_name}】${$t('更新差异')}`"
    :is-show.sync="isShow">
    <div :class="['diff-details', diffType, { 'is-template-conflict': isTemplateBindConflict }]" slot="content">
      <div class="diff-top">
        <div class="top-label">{{ $t('绑定变化：') }}</div>
        <div :class="['top-content', diffType]">
          <template v-if="isTemplateBindConflict">
            <i18n path="字段冲突，该字段已经被其他模板绑定，请删除该模型或修改模板">
              <template #other>
                <bk-button
                  text
                  type="primary"
                  @click="handleOtherTemplate(modelBeforeField)">
                  「{{ boundTemplate?.name }} {{$t('模板')}}」
                </bk-button>
              </template>
            </i18n>
          </template>
          <template v-else>{{ diffTitle }}</template>
        </div>
      </div>
      <div class="diff-table">
        <div class="table-head">
          <div class="col before-col">{{$t('绑定前')}}</div>
          <div class="col after-col">{{$t('绑定后')}}</div>
        </div>
        <div class="table-body">
          <div class="col before-col">
            <div class="diff-item">
              <div class="field-key">{{$t('唯一标识')}}</div>
              <div class="field-value" v-bk-overflow-tips>{{ modelBeforeField.bk_property_id || '--' }}</div>
            </div>
            <div class="diff-item">
              <div class="field-key">{{$t('名称')}}</div>
              <div class="field-value" v-bk-overflow-tips>{{ modelBeforeField.bk_property_name || '--' }}</div>
            </div>
            <div class="diff-item">
              <div class="field-key">{{$t('字段类型')}}</div>
              <div class="field-value">
                <span v-if="diffType === DIFF_TYPES.NEW">--</span>
                <span v-else>{{ PROPERTY_TYPE_NAMES[modelBeforeField.bk_property_type] }}</span>
              </div>
            </div>
            <div class="diff-item is-option" v-if="optionTypes.includes(modelBeforeField.bk_property_type)">
              <div class="field-key">{{getOptionKeyName(modelBeforeField)}}</div>
              <div class="field-value">
                <template v-if="numberLikes.includes(modelBeforeField.bk_property_type)">
                  <flex-row :colon="true">
                    <div>{{ $t('最小值') }}</div>
                    <div>{{ modelBeforeField?.option?.min || (modelBeforeField.option?.min === 0 ? 0 : '--') }}</div>
                  </flex-row>
                  <flex-row :colon="true">
                    <div>{{ $t('最大值') }}</div>
                    <div>{{ modelBeforeField?.option?.max || (modelBeforeField.option?.max === 0 ? 0 : '--') }}</div>
                  </flex-row>
                </template>
                <template v-if="stringLikes.includes(modelBeforeField.bk_property_type)">
                  <div>{{ modelBeforeField.option || '--' }}</div>
                </template>
                <template v-if="listLikes.includes(modelBeforeField.bk_property_type)">
                  <div v-for="(option, index) in (modelBeforeField?.option || [])" :key="index">
                    {{ option || '--' }}
                  </div>
                  <span v-if="!modelBeforeField?.option">--</span>
                </template>
                <template v-if="enumLikes.includes(modelBeforeField.bk_property_type)">
                  <flex-row :colon="true" v-for="(option, index) in (modelBeforeField?.option || [])" :key="index">
                    <div>{{ option.id }}</div>
                    <div>{{ option.name }}</div>
                  </flex-row>
                  <span v-if="!modelBeforeField?.option">--</span>
                </template>
              </div>
            </div>
            <div class="diff-item" v-if="![PROPERTY_TYPES.INNER_TABLE].includes(modelBeforeField.bk_property_type)">
              <div class="field-key">{{ $t('默认值') }}</div>
              <div class="field-value">
                <cmdb-property-value
                  class="property-value"
                  v-if="[PROPERTY_TYPES.ORGANIZATION, PROPERTY_TYPES.ENUMQUOTE]
                    .includes(modelBeforeField.bk_property_type)"
                  :value="getDefaultValue(modelBeforeField)"
                  :property="modelBeforeField">
                </cmdb-property-value>
                <span v-else>{{ getDefaultValue(modelBeforeField) || '--' }}</span>
              </div>
            </div>
            <div class="diff-item">
              <div class="field-key">{{ $t('是否可编辑') }}</div>
              <div class="field-value">
                <span v-if="diffType === DIFF_TYPES.NEW">--</span>
                <span v-else>
                  {{ modelBeforeField.editable ? $t('是') : $t('否') }}
                </span>
              </div>
            </div>
            <div class="diff-item">
              <div class="field-key">{{ $t('是否必填') }}</div>
              <div class="field-value">
                <span v-if="diffType === DIFF_TYPES.NEW">--</span>
                <span v-else>
                  {{ modelBeforeField.isrequired ? $t('是') : $t('否') }}
                </span>
              </div>
            </div>
            <div class="diff-item">
              <div class="field-key">{{ $t('用户提示') }}</div>
              <div class="field-value">{{ modelBeforeField.placeholder || '--' }}</div>
            </div>
          </div>
          <div class="col after-col">
            <div class="diff-item is-full" v-if="[DIFF_TYPES.UNCHANGED, DIFF_TYPES.UNBOUND].includes(diffType)">
              {{ diffTitle }}
            </div>
            <template v-else>
              <div :class="getAfterDiffItemClass('bk_property_id')">
                <div class="field-key">{{$t('唯一标识')}}</div>
                <div class="field-value" v-bk-overflow-tips>{{afterField.bk_property_id}}</div>
              </div>
              <div :class="getAfterDiffItemClass('bk_property_name')">
                <div class="field-key">{{ $t('名称') }}</div>
                <div class="field-value" v-bk-overflow-tips>{{afterField.bk_property_name}}</div>
              </div>
              <div :class="getAfterDiffItemClass('bk_property_type')">
                <div class="field-key">{{ $t('字段类型') }}</div>
                <div class="field-value">{{PROPERTY_TYPE_NAMES[afterField.bk_property_type]}}</div>
              </div>
              <div :class="getAfterDiffItemClass('option', ['is-option'])"
                v-if="optionTypes.includes(afterField.bk_property_type)">
                <div class="field-key">{{ getOptionKeyName(afterField) }}</div>
                <div class="field-value">
                  <template v-if="numberLikes.includes(afterField.bk_property_type)">
                    <flex-row :colon="true">
                      <div>{{ $t('最小值') }}</div>
                      <div>{{ afterField.option?.min || (afterField.option?.min === 0 ? 0 : '--') }}</div>
                    </flex-row>
                    <flex-row :colon="true">
                      <div>{{ $t('最大值') }}</div>
                      <div>{{ afterField.option?.max || (afterField.option?.max === 0 ? 0 : '--') }}</div>
                    </flex-row>
                  </template>
                  <template v-if="stringLikes.includes(afterField.bk_property_type)">
                    <div>{{ afterField.option || '--' }}</div>
                  </template>
                  <template v-if="listLikes.includes(afterField.bk_property_type)">
                    <div v-for="(option, index) in afterField.option" :key="index">
                      {{ option || '--' }}
                    </div>
                  </template>
                  <template v-if="enumLikes.includes(afterField.bk_property_type)">
                    <flex-row :colon="true" v-for="(option, index) in afterField.option" :key="index">
                      <div>{{ option.id }}</div>
                      <div>{{ option.name }}</div>
                    </flex-row>
                  </template>
                </div>
              </div>
              <div :class="getAfterDiffItemClass('default')">
                <div class="field-key">{{ $t('默认值') }}</div>
                <div class="field-value">
                  <cmdb-property-value
                    class="property-value"
                    v-if="[PROPERTY_TYPES.ORGANIZATION, PROPERTY_TYPES.ENUMQUOTE].includes(afterField.bk_property_type)"
                    :value="getDefaultValue(afterField)"
                    :property="afterField">
                  </cmdb-property-value>
                  <span v-else>{{ getDefaultValue(afterField) || '--' }}</span>
                </div>
              </div>
              <div :class="getAfterDiffItemClass('editable')">
                <div class="field-key">{{ $t('是否可编辑') }}</div>
                <div class="field-value">
                  <span>{{ afterField.editable?.value ? $t('是') : $t('否') }}</span>
                </div>
              </div>
              <div :class="getAfterDiffItemClass('isrequired')">
                <div class="field-key">{{ $t('是否必填') }}</div>
                <div class="field-value">
                  <span>{{ afterField.isrequired?.value ? $t('是') : $t('否') }}</span>
                </div>
              </div>
              <div :class="getAfterDiffItemClass('placeholder')">
                <div class="field-key">{{ $t('用户提示') }}</div>
                <div class="field-value">
                  <span>{{ afterField.placeholder?.value || '--' }}</span>
                </div>
              </div>
              <div class="diff-item is-option">
                <div class="field-key">{{ $t('绑定模型后，在模型中可修改的选项') }}</div>
                <div class="field-value">
                  <flex-row :colon="true">
                    <div>{{ $t('在实例中可编辑') }}</div>
                    <div>{{ afterField.editable?.lock ? $t('不可修改') : $t('可修改') }}</div>
                  </flex-row>
                  <flex-row :colon="true">
                    <div>{{ $t('是否必填') }}</div>
                    <div>{{ afterField.isrequired?.lock ? $t('不可修改') : $t('可修改') }}</div>
                  </flex-row>
                  <flex-row :colon="true">
                    <div>{{ $t('用户提示') }}</div>
                    <div>{{ afterField.placeholder?.lock ? $t('不可修改') : $t('可修改') }}</div>
                  </flex-row>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </bk-sideslider>
</template>

<style lang="scss" scoped>
.diff-details {
  height: 100%;
  padding: 20px 40px;
  @include scrollbar-y;

  &.new {
    .diff-table {
      .after-col {
        .diff-item {
          background: #F0FCF4;
          .field-value {
            color: #2DCB56;
          }
        }
      }
    }
  }

  &.is-template-conflict {
    .diff-table {
      .after-col {
        .diff-item {
          background: #FEF2F2;
          .field-value {
            color: #EA3636;
          }
        }
      }
    }
  }

  &.unbound {
    .diff-table {
      .diff-item {
        background: #F0F1F5;
      }
    }
  }

  &.unchanged {
    .diff-table {
      .diff-item {
        background: #F5F7FA;
      }
    }
  }
}
.diff-top {
  display: flex;
  align-items: center;
  margin-bottom: 16px;

  .top-label {
    font-weight: 700;
    font-size: 14px;
    color: #63656E;
  }

  .top-content {
    font-size: 14px;
    &.new {
      color: #2DCB56;
    }
    &.conflict {
      color: #EA3636;
    }
    &.update {
      color: #FF9C01;
    }
  }
}
.diff-table {
  display: grid;
  grid-template-rows: 42px auto;
  max-height: calc(100% - 42px);
  border: 1px solid #DCDEE5;
  @include scrollbar-y;

  .table-head {
    display: grid;
    grid-template-columns: 1fr 1fr;
    font-size: 12px;
    font-weight: 700;
    line-height: 42px;
    .col {
      padding-left: 18px;
      overflow: hidden;
    }
    .before-col {
      background: #F5F7FA;
    }
    .after-col {
      background: #F0F1F5;
    }
  }

  .table-body {
    display: grid;
    gap: 16px;
    grid-template-columns: 1fr 1fr;
    padding: 12px 16px;
    font-size: 12px;
    background: #FFF;
    box-shadow: inset 0 1px 0 0 #DCDEE5;

    .col {
      display: flex;
      flex-direction: column;
      gap: 8px;
      overflow: hidden;
    }

  }

  .diff-item {
    display: flex;
    align-items: center;
    gap: 4px;
    height: 30px;
    width: 100%;
    background: #F5F7FA;
    padding: 0 12px;

    .field-key {
      position: relative;
      padding-right: 14px;
      &::after {
        position: absolute;
        right: 0;
        content: "：";
      }
    }

    .field-value {
      max-width: calc(100% - 70px);
      @include ellipsis;
    }

    &.is-option {
      flex-direction: column;
      align-items: flex-start;
      height: auto;
      padding: 12px;

      .field-value {
        padding-left: 56px;
      }
    }

    &.is-full {
      height: 100%;
      justify-content: center;
    }

    &.is-update {
      background:#FFF9EF;
      .field-value {
        color: #FF9C01;
      }
    }

    &.is-conflict {
      background: #FEF2F2;
      .field-value {
        color: #EA3636;
      }
      & ~ .is-conflict {
        background: #F5F7FA;
        .field-value {
          color: #63656E;
        }
      }
    }
  }
}
</style>
