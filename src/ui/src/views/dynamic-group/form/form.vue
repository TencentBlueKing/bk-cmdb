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
  <bk-sideslider
    :transfer="true"
    :width="1202"
    :title="title"
    :is-show.sync="isShow"
    :before-close="handleSliderBeforeClose"
    class="dynamic-slidebar"
    @hidden="handleHidden">
    <bk-resize-layout
      :collapsible="false"
      :initial-divide="412"
      :min="400"
      :max="500"
      slot="content"
      style="height: 100%;">
      <div slot="aside" class="dynamic-group-info">
        <cmdb-sticky-layout class="dynamic-sticky-layout">
          <template #default="{ sticky }">
            <bk-form
              class="dynamic-group-form"
              ref="form"
              form-type="vertical"
              v-bkloading="{ isLoading: $loading([request.mainline, request.property, request.details]) }">
              <h5 class="form-title">
                {{ $t('基础信息') }}
              </h5>
              <bk-form-item :label="$t('分组名称')" required>
                <bk-input class="form-item"
                  v-model.trim="formData.name"
                  v-validate="'required|length:256'"
                  data-vv-name="name"
                  :data-vv-as="$t('查询名称')"
                  :disabled="isPreviewProp"
                  :placeholder="$t('请输入xx', { name: $t('查询名称') })">
                </bk-input>
                <p class="form-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
              </bk-form-item>
              <bk-form-item :label="$t('查询对象')" required>
                <form-target class="form-item"
                  ref="formTarget"
                  v-model="formData.bk_obj_id"
                  :disabled="!isCreateMode"
                  :show-cancel-dialog="true"
                  @canShowDialog="handleCanShowDialog"
                  @change="handleModelChange">
                </form-target>
              </bk-form-item>

              <bk-form-item class="condition-form" v-for="item in conditionGroup" :key="item.id">
                <cmdb-collapse
                  arrow-type="filled"
                  :auto-expand="true"
                  :list="getList(item.type)"
                  @collapse-change="handleCollapseChange">
                  <template #title>
                    <span v-bk-tooltips.top="{
                      content: $t(item.tip)
                    }">
                      {{$t(item.name)}}
                    </span>
                  </template>
                  <form-property-list
                    ref="propertyList"
                    @remove="handleRemoveProperty"
                    @toggle="handleToggleProperty"
                    :disabled="isPreviewProp"
                    :condition-type="item.type">
                  </form-property-list>
                </cmdb-collapse>

                <condition-picker class="condition-picker"
                  :text="$t('添加')"
                  icon="icon-plus-circle"
                  :selected="selectedProperties"
                  :property-map="propertyMap"
                  :handler="handlePropertySelected"
                  :disabled="isPreviewProp"
                  :condition-type="item.type">
                </condition-picker>
              </bk-form-item>

              <div class="no-condition">
                <input type="hidden"
                  v-validate="'min_value:1'"
                  data-vv-name="condition"
                  data-vv-validate-on="submit"
                  :data-vv-as="$t('查询条件')"
                  v-model="selectedProperties.length">
                <p class="form-error" v-if="errors.has('condition')">{{$t('请添加查询条件')}}</p>
              </div>
              <div :class="['dynamic-group-options', { 'no-fixed': !sticky }]">
                <cmdb-auth :auth="saveAuth">
                  <bk-button class="mr10" slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="disabled"
                    :loading="$loading([request.create, request.update])"
                    @click="handleConfirm">
                    {{ $t(confirmText) }}
                  </bk-button>
                </cmdb-auth>
                <bk-button v-show="!isPreviewProp"
                  class="mr10" theme="default" @click="handlePreview" :disabled="!selectedProperties.length">
                  {{$t('预览')}}
                </bk-button>
                <bk-popconfirm
                  :content="$t('确定清空分组条件')"
                  width="280"
                  trigger="click"
                  :confirm-text="$t('确定')"
                  :cancel-text="$t('取消')"
                  @confirm="handleClearCondition">
                  <bk-button v-show="!isPreviewProp" class="mr10" theme="default"
                    :disabled="!selectedProperties.length">
                    {{$t('清空条件')}}
                  </bk-button>
                </bk-popconfirm>
                <bk-button v-show="!isPreviewProp"
                  class="btn-cancel" theme="default" @click="handleSliderBeforeClose('cancel')">
                  {{$t('取消')}}
                </bk-button>
              </div>
            </bk-form>
          </template>
        </cmdb-sticky-layout>
      </div>
      <div slot="main" class="dynamic-group-preview">
        <preview-result class="preview-result"
          :condition="previewCondition" :mode="bkObjId" :properties="propertyMap[bkObjId]">
        </preview-result>
      </div>
    </bk-resize-layout>
  </bk-sideslider>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { t } from '@/i18n'
  import FormPropertyList from './form-property-list.vue'
  import FormPropertySelector from './form-property-selector.js'
  import FormTarget from './form-target.vue'
  import RouterQuery from '@/router/query'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import useSideslider from '@/hooks/use-sideslider'
  import isEqual from 'lodash/isEqual'
  import PreviewResult from '../preview/preview-result.vue'
  import FilterStore from '../store'
  import { $success } from '@/magicbox'
  import ConditionPicker from '@/components/condition-picker'
  import { DYNAMIC_GROUP_COND_TYPES } from '@/dictionary/dynamic-group'

  const { IMMUTABLE, VARIABLE } = DYNAMIC_GROUP_COND_TYPES
  export default {
    components: {
      FormPropertyList,
      FormTarget,
      PreviewResult,
      ConditionPicker
    },
    props: {
      id: [String, Number],
      title: String,
      isPreview: {
        type: Boolean,
        value: false
      }
    },
    provide() {
      return {
        dynamicGroupForm: this
      }
    },
    data() {
      return {
        footerIsFixed: false,
        isPreviewData: false,
        bkObjId: 'host',
        previewCondition: {},
        isShow: false,
        details: null,
        formData: {
          name: '',
          bk_obj_id: 'host'
        },
        originFormData: {
          bk_obj_id: 'host',
          name: '',
        },
        selectedProperties: [],
        originProperties: [],
        request: Object.freeze({
          mainline: Symbol('mainline'),
          property: Symbol('property'),
          details: Symbol('details'),
          create: Symbol('create'),
          update: Symbol('update')
        }),
        availableModelIds: Object.freeze(['host', 'module', 'set']),
        availableModels: [],
        propertyMap: {},
        disabledPropertyMap: {},
        storageCondition: {},
        conditionGroup: [
          {
            id: 1,
            name: '可变条件',
            tip: '动态分组可变条件',
            type: VARIABLE
          }, {
            id: 2,
            name: '锁定条件',
            tip: '动态分组锁定条件',
            type: IMMUTABLE
          }
        ]
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectBiz', ['bizId']),
      isCreateMode() {
        return !this.id
      },
      searchTargetModels() {
        return this.availableModels.filter(model => ['host', 'set'].includes(model.bk_obj_id))
      },
      saveAuth() {
        if (this.id) {
          return { type: this.$OPERATION.U_CUSTOM_QUERY, relation: [this.bizId, this.id] }
        }
        return { type: this.$OPERATION.C_CUSTOM_QUERY, relation: [this.bizId] }
      },
      isPreviewProp() {
        return this.isPreviewData
      },
      confirmText() {
        let text = '保存'
        if (this.isPreviewProp) {
          text = '编辑'
        } else if (this.isCreateMode) {
          text = '提交'
        }
        return text
      }
    },
    watch: {
      selectedProperties: {
        deep: true,
        handler() {
          this.errors.remove('condition')
        }
      }
    },
    async created() {
      await this.getMainLineModels()
      await this.getModelProperties()
      if (this.id) {
        this.getDetails()
      }
      const { beforeClose, setChanged, setInfoData } = useSideslider()
      this.beforeClose = beforeClose
      this.setChanged = setChanged
      this.setInfoData = setInfoData
      this.isPreviewData = this.isPreview
    },
    methods: {
      getList(type) {
        return this.selectedProperties.filter(property => property.conditionType === type)
      },
      async getMainLineModels() {
        try {
          const models = await this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
            config: {
              requestId: this.request.mainline,
              fromCache: true
            }
          })
          // 业务调用方暂时只需要一下三种类型的查询
          // eslint-disable-next-line max-len
          const availableModels = this.availableModelIds.map(modelId => models.find(model => model.bk_obj_id === modelId))
          this.availableModels = Object.freeze(availableModels)
        } catch (error) {
          console.error(error)
        }
      },
      async getModelProperties() {
        try {
          const propertyMap = await this.$store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
            params: {
              bk_biz_id: this.bizId,
              bk_obj_id: { $in: this.availableModels.map(model => model.bk_obj_id) },
              bk_supplier_account: this.supplierAccount
            },
            config: {
              requestId: this.request.property,
              fromCache: true
            }
          })
          propertyMap.module.unshift(this.getServiceTemplateProperty())
          this.propertyMap = Object.freeze(propertyMap)

          Object.keys(this.propertyMap).forEach((objId) => {
            this.disabledPropertyMap[objId] = this.propertyMap[objId]
              .filter(item => item.bk_property_type === PROPERTY_TYPES.INNER_TABLE)
              .map(item => item.bk_property_id)
          })
        } catch (error) {
          console.error(error)
          this.propertyMap = {}
        }
      },
      getServiceTemplateProperty() {
        return {
          id: Date.now(),
          bk_obj_id: 'module',
          bk_property_id: 'service_template_id',
          bk_property_name: t('服务模板'),
          bk_property_index: -1,
          bk_property_type: 'service-template',
          isonly: true,
          ispre: true,
          bk_isapi: true,
          bk_issystem: true,
          isreadonly: true,
          editable: false,
          bk_property_group: null,
          _is_inject_: true
        }
      },
      async getDetails() {
        try {
          const details = await this.$store.dispatch('dynamicGroup/details', {
            bizId: this.bizId,
            id: this.id,
            config: {
              requestId: this.request.details
            }
          })
          const transformedDetails = this.transformDetails(details)
          const { name, bk_obj_id: modelId } = transformedDetails
          this.originFormData.name = name
          this.originFormData.bk_obj_id = modelId
          this.formData.name = name
          this.formData.bk_obj_id = modelId
          this.details = transformedDetails
          this.$nextTick(this.setDetailsSelectedProperties)
          setTimeout(() => {
            this.$refs.propertyList?.forEach(propertyList => propertyList?.setDetailsCondition())
            if (this.isPreview || this.id) {
              this.initPreviewParams()
            }
          })
        } catch (error) {
          console.error(error)
        }
      },
      transformDetails(details) {
        const { info } = details
        const transformedCondition = {
          condition: [],
          varCondition: []
        }
        Object.keys(info).forEach((type) => {
          info[type]?.forEach((data) => {
            const conditionType = type === IMMUTABLE
              ? IMMUTABLE : VARIABLE
            const realCondition = (data.condition || []).reduce((accumulator, current) => {
              current.conditionType = conditionType
              if (['$gte', '$lte'].includes(current.operator)) {
                // $gte和$lte，可能是单个field也可能是同一field的范围设置，如果是范围一个field会拆分为两条cond
                const isRange = data.condition.filter(cond => cond.field === current.field)?.length > 1

                // 将相同字段的$gte/$lte两个条件合并为一个range条件，用于表单组件渲染
                let index = accumulator.findIndex(exist => exist.field === current.field)
                if (index === -1) {
                  index = accumulator.push({
                    field: current.field,
                    operator: isRange ? '$range' : current.operator,
                    value: isRange ? [] : current.value,
                    conditionType
                  }) - 1
                }
                const range = accumulator[index]

                // 如果是范围并且确保field一致，需要组装为一个范围数组格式值
                if (isRange && current.field === range.field) {
                  range.value?.[current.operator === '$gte' ? 'unshift' : 'push'](current.value)
                }
              } else if (current.operator === '$eq') {
                // 将老数据的eq转换为当前支持的数据格式
                const transformType = ['singlechar', 'longchar', 'enum', 'objuser']
                const property = this.getConditionProperty(data.bk_obj_id, current.field)
                if (property && transformType.includes(property.bk_property_type)) {
                  accumulator.push({
                    field: current.field,
                    operator: '$in',
                    value: Array.isArray(current.value) ? current.value : [current.value],
                    conditionType
                  })
                } else {
                  accumulator.push(current)
                }
              } else {
                accumulator.push(current)
              }
              return accumulator
            }, [])
            if (data.time_condition) {
              data.time_condition.rules.forEach(({ field, start, end }) => {
                realCondition.push({
                  field,
                  operator: '$range',
                  value: [start, end],
                  conditionType
                })
              })
            }
            transformedCondition[conditionType].push({
              bk_obj_id: data.bk_obj_id,
              condition: realCondition
            })
          })
        })
        return {
          ...details,
          info: {
            ...transformedCondition
          }
        }
      },
      getConditionProperty(modelId, field) {
        const properties = this.propertyMap[modelId] || []
        return properties.find(property => property.bk_property_id === field)
      },
      setDetailsSelectedProperties() {
        const { condition, varCondition } = this.details.info
        const conditions = [...condition, ...varCondition]
        const properties = []
        conditions.forEach(({ bk_obj_id: modelId, condition }) => {
          condition.forEach(({ field, conditionType }) => {
            const property = this.propertyMap[modelId].find(property => property.bk_property_id === field)
            if (property) {
              property.conditionType = conditionType
              properties.push(property)
            }
          })
        })
        this.selectedProperties = this.$tools.clone(properties)
        this.originProperties = this.$tools.clone(properties)
        this.setFooterCls()
      },
      setFooterCls() {
        this.$nextTick(() => {
          // 根据选择的条件 展示不同的样式
          const el = document.querySelector('.dynamic-group-form')
          const { clientHeight, scrollHeight } = el
          // 是否出现了滚动条
          if (scrollHeight > clientHeight) {
            this.footerIsFixed = true
          } else {
            this.footerIsFixed = false
          }
        })
      },
      setProperty(property, conditionType) {
        const { bk_obj_id: modelId, bk_property_id: propertyId } = property
        const exchangeType = {
          [VARIABLE]: IMMUTABLE,
          [IMMUTABLE]: VARIABLE
        }
        property.conditionType = exchangeType[conditionType]
        this.getConditionProperty(modelId, propertyId).conditionType = exchangeType[conditionType]
      },
      setFormTarget(item) {
        this.$refs?.formTarget?.setSelected(item)
      },
      handleCollapseChange() {
        setTimeout(this.setFooterCls, 300)
      },
      handleCanShowDialog(item) {
        const data = {
          subTitle: '填写内容会清空，确认切换？',
          title: '确认切换？',
          okText: '确认'
        }
        const isChange = !isEqual(this.selectedProperties, this.originProperties)
        if (isChange) {
          this.setChanged(true)
          this.setInfoData(data)
          return this.beforeClose(() => {
            this.setInfoData()
            this.setFormTarget(item)
          }, () => {
            this.setInfoData()
          })
        }
        this.setFormTarget(item)
      },
      handleModelChange() {
        this.selectedProperties = []
      },
      handleClearCondition() {
        this.selectedProperties = []
      },
      handleShowPropertySelector(event) {
        this.formPropertySelector = FormPropertySelector.show({
          selected: this.selectedProperties,
          handler: this.handlePropertySelected
        }, this, event?.target)
      },
      handlePropertySelected(selected) {
        const { length } = selected
        if (!length) this.selectedProperties = []

        const addSelect = []
        const deleteSelect = []
        const selectedSet = new Set()

        selected.forEach(property => selectedSet.add(`${property.bk_property_id}-${property.bk_obj_id}`))
        this.selectedProperties.forEach((property) => {
          const { bk_property_id: propertyId, bk_obj_id: modelId } = property
          const key = `${propertyId}-${modelId}`
          if (selectedSet.has(key)) {
            selectedSet.delete(key)
          } else {
            deleteSelect.push(property)
          }
        })
        selected.forEach((property) => {
          const { bk_property_id: propertyId, bk_obj_id: modelId } = property
          const key = `${propertyId}-${modelId}`
          if (selectedSet.has(key)) {
            addSelect.push(property)
          }
        })

        deleteSelect.forEach(property => this.handleRemoveProperty(property))
        let start = 0
        const limit = 10
        while (start < addSelect.length) {
          setTimeout(() =>  this.selectedProperties.push(...addSelect.splice(0, limit)))
          start += limit
        }
        this.setFooterCls()
      },
      handleRemoveProperty(property) {
        const index = this.selectedProperties.findIndex(target => target.id === property.id)
        if (index > -1) {
          this.selectedProperties.splice(index, 1)
        }
        this.setFooterCls()
      },
      handleToggleProperty(property, conditionType) {
        this.handleRemoveProperty(property)
        this.setProperty(property, conditionType)
        this.selectedProperties.push(property)
      },
      async handlePreview() {
        const result = await this.validate()
        if (!result) {
          return
        }
        this.initPreviewParams()
      },
      validate() {
        return Promise.all(this.$refs.propertyList
          .map(propertyList => propertyList.$validator.validateAll()))
          .then(result => result.every(e => e))
      },
      getCondition() {
        return this.$refs.propertyList.reduce((current, prev) => {
          Object.assign(current, prev?.condition)
          return current
        }, {})
      },
      initPreviewParams() {
        this.bkObjId = this.formData.bk_obj_id
        FilterStore.setDynamicGroupModel(this.formData.bk_obj_id)
        const condition = this.getCondition()
        this.previewCondition = this.$tools.clone(condition)
      },
      async handleConfirm() {
        try {
          if (this.isPreviewProp) {
            this.isPreviewData = false
            return
          }
          const results = [
            await this.$validator.validateAll(),
            await this.validate()
          ]
          if (results.some(isValid => !isValid)) {
            return false
          }
          if (this.id) {
            await this.updateDynamicGroup()
            $success(t('保存成功'))
          } else {
            await this.createDynamicGroup()
            $success(t('新建成功'))
          }
          this.close('submit')
        } catch (error) {
          console.error(error)
        }
      },
      updateDynamicGroup() {
        return this.$store.dispatch('dynamicGroup/update', {
          bizId: this.bizId,
          id: this.id,
          params: {
            bk_biz_id: this.bizId,
            bk_obj_id: this.formData.bk_obj_id,
            name: this.formData.name,
            info: {
              ...this.getSubmitCondition()
            }
          },
          config: {
            requestId: this.request.update
          }
        })
      },
      createDynamicGroup() {
        return this.$store.dispatch('dynamicGroup/create', {
          params: {
            bk_biz_id: this.bizId,
            bk_obj_id: this.formData.bk_obj_id,
            name: this.formData.name,
            info: {
              ...this.getSubmitCondition()
            }
          },
          config: {
            requestId: this.request.create
          }
        })
      },
      getSubmitCondition() {
        const baseConditionMap = {
          [VARIABLE]: {},
          [IMMUTABLE]: {}
        }
        const timeConditionMap = {
          [VARIABLE]: {},
          [IMMUTABLE]: {}
        }
        const propertyCondition = this.getCondition()
        Object.values(propertyCondition).forEach(({ property, operator, value }) => {
          const type = property?.conditionType === IMMUTABLE
            ? IMMUTABLE : VARIABLE
          if (property.bk_property_type === 'time') { // 时间类型特殊处理
            const timeCondition = timeConditionMap[type][property.bk_obj_id] || { oper: 'and', rules: [] }
            const [start, end] = value
            timeCondition.rules.push({
              field: property.bk_property_id,
              start,
              end
            })
            timeConditionMap[type][property.bk_obj_id] = timeCondition
            return
          }
          const submitCondition = baseConditionMap[type][property.bk_obj_id] || []
          if (operator === '$range') {
            const [start, end] = value
            submitCondition.push({
              field: property.bk_property_id,
              operator: '$gte',
              value: start
            }, {
              field: property.bk_property_id,
              operator: '$lte',
              value: end
            })
          } else {
            submitCondition.push({
              field: property.bk_property_id,
              operator,
              value
            })
          }
          baseConditionMap[type][property.bk_obj_id] = submitCondition
        })
        const baseConditions = {}
        Object.keys(baseConditionMap).forEach((type) => {
          baseConditions[type] = Object.keys(baseConditionMap[type]).map(modelId => ({
            bk_obj_id: modelId,
            condition: baseConditionMap[type][modelId]
          }))
        })
        Object.keys(timeConditionMap).forEach((type) => {
          Object.keys(timeConditionMap[type]).forEach((modelId) => {
            const condition = baseConditions[type].find(condition => condition.bk_obj_id === modelId)
            if (condition) {
              condition.time_condition = timeConditionMap[type][modelId]
            } else {
              baseConditions[type].push({
                bk_obj_id: modelId,
                time_condition: timeConditionMap[type][modelId]
              })
            }
          })
        })
        baseConditions.variable_condition = baseConditions[VARIABLE]
        delete baseConditions[VARIABLE]
        return baseConditions
      },
      close(type) {
        this.isShow = false
        if (type !== 'normal') {
          RouterQuery.set({
            _t: Date.now(),
            action: ''
          })
        }
      },
      show() {
        this.isShow = true
      },
      handleSliderBeforeClose(type = 'normal') {
        const changedValues = !isEqual(this.formData, this.originFormData)
        const changedProperties =  !isEqual(this.selectedProperties, this.originProperties)
        if (changedValues || changedProperties) {
          this.setChanged(true)
          return this.beforeClose(() => {
            this.close(type)
          })
        }
        this.close(type)
        FormPropertySelector?.hide(this.formPropertySelector)
        return true
      },
      handleHidden() {
        this.$emit('close')
      }
    }
  }
</script>

<style lang="scss" scoped>
.dynamic-slidebar {
  :deep(.bk-sideslider-content) {
    overflow-x: hidden;
    .bk-resize-layout-border {
      border-top: 0;
    }
  }
}
.dynamic-group-info {
  width: 100%;
  float: left;
  height: 100%;

  .no-condition {
    height: 16px;
    .form-error {
      position: inherit;
    }
  }

  .condition-form {
    position: relative;
    margin-bottom: 24px;


    :deep(.collapse-trigger) {
      padding: 3px 8px;
      background: #F0F1F5;

      .collapse-text{
        >span {
          border-bottom: 1px dashed #63656E;
        }
      }
    }
    :deep(.collapse-content) {
      .bk-form {
        .bk-form-item {
          &:first-child {
            margin-top: 12px !important;
          }
          &:last-child {
            margin-bottom: 0 !important;
          }
        }
      }
    }

    .condition-picker {
      position: absolute;
      right: 8px;
      top: 0px;
      display: flex;
      align-items: center;
    }
  }
}
.dynamic-group-preview {
  width: 100%;
  float: right;
  height: 100%;
  background: #F5F7FA;
}
.dynamic-sticky-layout {
  @include scrollbar-y;
  height: 100%;
}
.dynamic-group-form {
  padding: 18px 24px 0;
  :nth-last-child(4) {
    margin-bottom: 0px !important;
  }
  .form-item {
    width: 100%;
  }
  .form-error {
    position: absolute;
    top: 100%;
    font-size: 12px;
    line-height: 14px;
    color: $dangerColor;
  }
  .form-title {
    font-weight: 700;
    font-size: 14px;
    color: #313238;
    line-height: 22px;
    margin-bottom: 10px;
  }
  :deep(.bk-form-item) {
    margin-bottom: 20px;
    margin-top: 0 !important;
  }
  :deep(.form-condition-button) {
    margin-top: 0 !important;
    height: 36px;

    > div {
      display: flex;
      align-items: center;
      .bk-icon {
        top: 0;
      }
    }
  }
}
.dynamic-group-options {
  display: flex;
  align-items: center;
  padding: 10px 24px;
  margin: 0 -24px;
  border-top: 1px solid $borderColor;
  background: #FAFBFD;
  position: sticky;
  bottom: 0;
  z-index: 999;
  left: 24px;
  right: 24px;

  :deep(.bk-button) {
    width: 88px;
    padding: 0 !important;
    &.btn-cancel {
      position: relative;

      &::before {
        content: '';
        width: 1px;
        height: 16px;
        background: #C4C6CC;
        display: inline-block;
        position: absolute;
        left: -6PX;
        top: 50%;
        transform: translateY(-50%);
      }
    }
  }
}
.no-fixed {
  position: static;
  border-top: 0;
  background: transparent;
}
:deep(.form-property-list) {
  .form-property-item {
    .item-value {
      margin: 0 !important;
    }
  }
}
</style>
