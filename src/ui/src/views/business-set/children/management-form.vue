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
  <div class="management-form">
    <bk-sideslider
      v-transfer-dom
      :is-show.sync="isShow"
      :title="title"
      :quick-close="false"
      :width="800"
      :before-close="handleSliderBeforeClose">
      <cmdb-form slot="content" v-if="isShow"
        :properties="properties"
        :property-groups="propertyGroups"
        :inst="formData"
        :type="isEdit ? 'update' : 'create'"
        :save-auth="saveAuth"
        :submitting="submitting"
        @on-submit="handleSave"
        @on-cancel="handleSliderBeforeClose">
        <template #append>
          <div class="custom-group">
            <div class="group-title">{{$t('资源定义范围')}}</div>
            <ul class="property-list">
              <li class="property-item full-width">
                <div class="property-name">
                  <span class="property-name-text">{{$t('业务范围')}}</span>
                </div>
                <div class="property-value" v-bk-tooltips="{ content: $t('内置业务集不可编辑'), disabled: !isBuiltin }">
                  <business-scope-settings-form
                    class="form-component"
                    :disabled="isBuiltin"
                    :data="scopeSettingsFormData"
                    @change="handleScopeSettingsChange" />
                </div>
              </li>
            </ul>
          </div>
        </template>
        <template #side-options>
          <bk-button class="button-preview" @click="handlePreview">{{$t('预览')}}</bk-button>
        </template>
      </cmdb-form>
    </bk-sideslider>
    <business-scope-preview v-bind="previewProps" :show.sync="previewProps.show" />
  </div>
</template>

<script>
  import { computed, defineComponent, reactive, ref, toRefs, watch } from '@vue/composition-api'
  import cloneDeep from 'lodash/cloneDeep'
  import isEqual from 'lodash/isEqual'
  import store from '@/store'
  import { t } from '@/i18n'
  import { OPERATION } from '@/dictionary/iam-auth'
  import { $success } from '@/magicbox/index.js'
  import Utils from '@/components/filters/utils'
  import queryBuilderOperator from '@/utils/query-builder-operator'
  import businessScopeSettingsForm from '@/components/business-scope/settings-form.vue'
  import businessScopePreview from '@/components/business-scope/preview.vue'
  import businessSetService from '@/service/business-set/index.js'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default defineComponent({
    components: {
      businessScopeSettingsForm,
      businessScopePreview
    },
    props: {
      show: {
        type: Boolean,
        default: false
      },
      data: {
        type: Object,
        default: () => ({})
      },
      properties: {
        type: Array,
        required: true,
        default: () => ([])
      },
      propertyGroups: {
        type: Array,
        required: true,
        default: () => ([])
      }
    },
    setup(props, { emit }) {
      const {
        show: isShow,
        data: formData
      } = toRefs(props)

      const getModelById = store.getters['objectModelClassify/getModelById']
      const model = computed(() => getModelById(BUILTIN_MODELS.BUSINESS_SET) || {})
      const isEmptyRuleValue = value => value === null || value === undefined || !value.toString().length

      const submitting = ref(false)
      const isEdit = computed(() => Boolean(formData.value.bk_biz_set_id))
      const title = computed(() => (isEdit.value ? t('编辑') : `${t('创建')}${model.value.bk_obj_name}`))

      const isBuiltin = computed(() => formData.value?.default === 1)

      const previewProps = reactive({
        show: false,
        mode: 'before',
        payload: {}
      })

      // 格式化业务范围配置表单数据
      const scopeSettingsFormData = computed(() => {
        if (!isEdit.value) return {}

        const data = {
          selectedBusiness: [],
          condition: []
        }
        formData.value.bk_scope?.filter?.rules?.forEach((rule) => {
          if (rule.field === 'bk_biz_id') {
            data.selectedBusiness.push(...rule.value)
          } else {
            data.condition.push({
              field: rule.field,
              value: rule.value
            })
          }
        })

        return data
      })

      const saveAuth = computed(() => ({ type: isEdit.value ? OPERATION.U_BUSINESS_SET : OPERATION.C_BUSINESS_SET }))

      // 待保存的表单数据
      const defaultSaveData = () => ({
        bk_biz_set_attr: {},
        bk_scope: {
          match_all: true
        }
      })
      let saveData = defaultSaveData()

      let scopeCopy = null
      let scopeChanged = false

      const handleSave = async (values, changedValues, originalValues, type) => {
        try {
          submitting.value = true
          let result = null
          if (type === 'update') {
            // 编辑时模型属性中会存在bk_scope字段，这里删除掉使用saveData中的bk_scope
            Reflect.deleteProperty(changedValues, 'bk_scope')

            result = await businessSetService.update({
              bk_biz_set_ids: [formData.value.bk_biz_set_id],
              data: {
                ...saveData,
                bk_scope: scopeChanged ? saveData.bk_scope : undefined,
                bk_biz_set_attr: { ...changedValues }
              }
            })
            $success(t('编辑成功'))
          } else {
            saveData.bk_biz_set_attr = { ...values }
            result = await businessSetService.create(saveData)
            $success(t('创建成功'))
          }
          emit('save-success', result, type)
        } catch (err) {
          console.error(err)
        } finally {
          submitting.value = false
        }
      }

      // 业务范围数据变化时更新saveData
      const handleScopeSettingsChange = (data) => {
        const { condition, selectedBusiness } = data
        const rules = []
        for (let i = 0, item; item = condition[i]; i++) {
          const { field, value, property } = item
          // 忽略空值
          if (isEmptyRuleValue(value)) {
            continue
          }
          const { operator } = Utils.getDefaultData(property)
          rules.push({
            field,
            operator: queryBuilderOperator(operator),
            value
          })
        }

        // 业务ID
        if (selectedBusiness?.length) {
          rules.unshift({
            field: 'bk_biz_id',
            operator: 'in',
            value: [...selectedBusiness]
          })
        } else {
          const index = saveData.bk_scope?.filter?.rules?.findIndex(item => item.field === 'bk_biz_id')
          if (index !== undefined && ~index) {
            saveData.bk_scope.filter.rules.splice(index, 1)
          }
        }

        // 接口协议约定参数match_all与filter互斥
        const matchAll = !rules.length
        saveData.bk_scope.match_all = matchAll
        if (!matchAll) {
          saveData.bk_scope.filter = {
            condition: 'AND',
            rules
          }
        } else {
          Reflect.deleteProperty(saveData.bk_scope, 'filter')
        }

        if (!scopeCopy) {
          scopeCopy = cloneDeep(saveData.bk_scope)
        }

        scopeChanged = !isEqual(scopeCopy, saveData.bk_scope)
      }

      const resetData = () => {
        // 关闭时重置saveData
        saveData = defaultSaveData()

        scopeCopy = null
        scopeChanged = false
      }

      const handleSliderBeforeClose = () => {
        emit('update:show', false)
      }

      const handlePreview = async () => {
        previewProps.show = true
        previewProps.payload = { ...saveData }
      }

      watch(isShow, (show) => {
        if (!show) {
          resetData()
        }
      })

      return {
        isShow,
        isEdit,
        isBuiltin,
        title,
        formData,
        saveAuth,
        scopeSettingsFormData,
        submitting,
        handleSave,
        handleScopeSettingsChange,
        handleSliderBeforeClose,
        previewProps,
        handlePreview
      }
    }
  })
</script>

<style lang="scss" scoped>
  .custom-group {
    padding-left: 32px;
    margin: 8px 0 16px 0;

    .group-title {
      font-size: 14px;
      font-weight: 700;
    }
  }

  .property-list {
    padding: 4px 0;
    display: flex;
    flex-wrap: wrap;
    .property-item {
      width: 50%;
      margin: 12px 0 0;
      padding: 0 54px 0 0;
      font-size: 12px;
      flex: 0 0 50%;
      max-width: 50%;
      .property-name {
        margin: 2px 0 6px;
        color: $cmdbTextColor;
        line-height: 24px;
        font-size: 0;
      }
      .property-name-text {
        position: relative;
        display: inline-block;
        max-width: calc(100% - 20px);
        padding: 0 10px 0 0;
        vertical-align: middle;
        font-size: 14px;
        @include ellipsis;
        &.required:after {
          position: absolute;
          left: 100%;
          top: 0;
          margin: 0 0 0 -10px;
          content: "*";
          color: #ff5656;
        }
      }
      .property-value {
        font-size: 0;
        position: relative;
        display: flex;
        /deep/ .control-append-group {
          .bk-input-text {
            flex: 1;
          }
        }
      }
      .form-component:not(.form-bool) {
        flex: 1;
      }

      &.full-width {
        flex: 1;
        padding-right: 54px;
        width: 100%;
        max-width: unset;
      }
    }
  }

  .button-preview {
    min-width: 76px;
    margin: 4px;
  }
</style>
