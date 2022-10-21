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

<script lang="ts">
  import { computed, defineComponent, onBeforeUnmount, reactive, ref, toRefs, watch, watchEffect } from 'vue'
  import isEqual from 'lodash/isEqual'
  import store from '@/store'
  import { t } from '@/i18n'
  import { formatValue } from '@/utils/tools'
  import { $bkInfo } from '@/magicbox/index.js'
  import { OPERATION } from '@/dictionary/iam-auth'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import PropertyConfig from '@/components/property-config/index.vue'
  import ProcessForm from './process-form.vue'
  import ProcessTable from './process.vue'
  import useTemplateData from './use-template-data'
  import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants'

  export default defineComponent({
    components: {
      GridLayout,
      GridItem,
      PropertyConfig,
      ProcessForm,
      ProcessTable,
    },
    props: {
      dataId: {
        type: Number
      },
      submitDisabled: Boolean,
      isClone: Boolean
    },
    setup(props, { emit }) {
      const $processForm = ref(null)
      const $propertyConfig = ref(null)

      const bizId = computed(() => store.getters['objectBiz/bizId'])
      const localProcessTemplate = computed(() => store.getters['serviceProcess/localProcessTemplate'])

      const templateId = computed(() => Number(props.dataId))

      const isEdit = computed(() => templateId.value > 0)

      const auth = computed(() => {
        if (props.isClone || !isEdit.value) {
          return {
            type: OPERATION.C_SERVICE_TEMPLATE,
            relation: [bizId.value]
          }
        }
        return {
          type: OPERATION.U_SERVICE_TEMPLATE,
          relation: [bizId.value, templateId.value]
        }
      })

      const state = reactive({
        moduleProperties: [],
        modulePropertyGroup: [],
        processProperties: [],
        processPropertyGroup: [],
        fullCategories: [],
        primaryCategories: [],
        secCategories: [],
        currentSecCategories: [],
        configProperties: [],
        isLocalProcessEdit: true, // 此组件均使用本地编辑模式，数据最后统一提交
        processSlider: {
          show: false,
          title: '',
          form: {
            inst: {},
            dataIndex: null,
            type: null
          }
        }
      })

      // 用于接收初始化时的表单数据
      let formDataCopy = null

      const formData = ref({
        id: null,
        templateName: '',
        primaryCategory: '',
        secCategory: '',
        propertyConfig: {},
        processList: localProcessTemplate.value,
      })

      const requestIds = {
        processList: Symbol()
      }

      const excludeModuleProperties = [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.MODULE].NAME]

      watchEffect(async () => {
        const {
          moduleProperties,
          modulePropertyGroup,
          processProperties,
          processPropertyGroup,
          primaryCategories,
          secCategories,
          basic,
          configProperties,
          propertyConfig,
          processList,
          formDataCopy: formDataCopyRaw
        } = await useTemplateData(bizId.value, templateId.value, isEdit.value)

        formDataCopy = formDataCopyRaw

        emit('data-loaded')

        // 属性数据主要提供给不同功能的组件使用
        state.moduleProperties = moduleProperties
        state.modulePropertyGroup = modulePropertyGroup
        state.processProperties = processProperties
        state.processPropertyGroup = processPropertyGroup

        state.primaryCategories = primaryCategories
        state.secCategories = secCategories

        // 编辑态数据初始化
        if (isEdit.value) {
          // 基础信息
          formData.value.id = basic.id
          formData.value.templateName = basic.templateName
          formData.value.primaryCategory = basic.primaryCategory
          formData.value.secCategory = basic.secCategory
          state.currentSecCategories = state.secCategories
            .filter(category => category.bk_parent_id === basic.primaryCategory)

          // 进程列表
          formData.value.processList = processList
          store.commit('serviceProcess/setLocalProcessTemplate', formatProcessSubmitData(processList))

          // 属性设置
          formData.value.propertyConfig = propertyConfig
          state.configProperties = configProperties

          if (!props.isClone) {
            store.commit('setTitle', `${t('编辑模板')}【${basic.templateName}】`)
          }
        }
      })

      const isFormDataChanged = computed(() => {
        const formDataLatest = {
          templateName: formData.value.templateName,
          primaryCategory: formData.value.primaryCategory,
          secCategory: formData.value.secCategory,
          propertyConfig: $propertyConfig.value?.getData(),
          processList: formData.value.processList
        }
        return !isEqual(formDataCopy, formDataLatest)
      })

      // 动态监测表单值是否变化，设置提交按钮的状态
      watch(isFormDataChanged, changed => emit('update:submitDisabled', !changed))

      const handleSelectPrimaryCategory = (id: number) => {
        state.currentSecCategories = state.secCategories.filter(classification => classification.bk_parent_id === id)
        if (!state.currentSecCategories?.length) {
          formData.value.secCategory = ''
        }
      }

      const handleCreateProcess = () => {
        state.processSlider.show = true
        state.processSlider.title = t('添加进程')
        state.processSlider.form.type = 'create'
        state.processSlider.form.inst = {}
      }

      const handleCancelProcess = () => {
        state.processSlider.show = false
      }

      const handleUpdateProcess = (template, index) => {
        state.processSlider.show = true
        state.processSlider.title = template.bk_func_name.value
        state.processSlider.form.type = 'update'
        state.processSlider.form.inst = template
        state.processSlider.form.dataIndex = index
      }

      const handleDeleteProcess = (template, index) => {
        $bkInfo({
          title: t('确认删除模板进程'),
          confirmFn: () => {
            store.commit('serviceProcess/deleteLocalProcessTemplate', { process: template, index })
            formData.value.processList = localProcessTemplate.value
          }
        })
      }

      const formatProcessSubmitData = (process) => {
        const processList = !Array.isArray(process) ? [process] : process
        processList.forEach((data) => {
          Object.keys(data).forEach((key) => {
            const property = state.processProperties.find(property => property.bk_property_id === key)
            if (property && property.bk_property_type === 'table') {
              (data[key].value || []).forEach((row) => {
                Object.keys(row).forEach((rowKey) => {
                  if (typeof row[rowKey] === 'object') {
                    const option = property.option || []
                    const columnProperty = option.find(columnProperty => columnProperty.bk_property_id === rowKey) || {}
                    row[rowKey].value = formatValue(row[rowKey].value, columnProperty)
                  }
                })
              })
            } else if (typeof data[key] === 'object') {
              data[key].value = formatValue(data[key].value, property)
            }
          })
        })

        return Array.isArray(process) ? processList : processList[0]
      }

      const handleProcessSliderBeforeClose = () => {
        const hasChanged = $processForm.value && $processForm.value.hasChange()
        if (hasChanged) {
          return new Promise((resolve) => {
            $bkInfo({
              title: t('确认退出'),
              subTitle: t('退出会导致未保存信息丢失'),
              extCls: 'bk-dialog-sub-header-center',
              confirmFn: () => {
                resolve(true)
              },
              cancelFn: () => {
                resolve(false)
              }
            })
          })
        }
        return true
      }

      onBeforeUnmount(() => {
        store.commit('serviceProcess/clearLocalProcessTemplate')
      })

      return {
        ...toRefs(state),
        formData,
        $processForm,
        $propertyConfig,
        requestIds,
        auth,
        excludeModuleProperties,
        isFormDataChanged,
        formatProcessSubmitData,
        handleSelectPrimaryCategory,
        handleCreateProcess,
        handleCancelProcess,
        handleUpdateProcess,
        handleDeleteProcess,
        handleProcessSliderBeforeClose
      }
    },
    methods: {
      getData() {
        const propertyConfigData = this.$refs.$propertyConfig.getData()

        const attributes = []
        for (const [key, value] of Object.entries(propertyConfigData)) {
          attributes.push({
            bk_attribute_id: Number(key),
            bk_property_value: value
          })
        }

        const processes = this.formData.processList.map((process) => {
          delete process.sign_id
          const data = {
            property: this.formatProcessSubmitData(process)
          }
          if (process.process_id) {
            data.id = process.process_id
          }

          return data
        })

        return {
          name: this.formData.templateName,
          service_category_id: this.formData.secCategory,
          processes,
          attributes
        }
      },
      async validateAll() {
        // 基础信息校验
        const basicValid = await this.$validator.validateAll()

        // 属性设置校验
        const configValid = await this.$refs.$propertyConfig.validateAll()

        const result = basicValid && configValid

        return result
      },
      async validate() {
        const basicResults = []
        for (const [key, val] of Object.entries(this.fields)) {
          if (val.dirty) {
            basicResults.push(await this.$validator.validate(key))
          }
        }

        const basicValid = basicResults.every(valid => valid)
        const configValid = this.$refs.$propertyConfig.validate()

        return basicValid && configValid
      },
      async changeHandler() {
        this.$nextTick(async () => {
          const valid = await this.validate()
          const { isFormDataChanged } = this
          this.$emit('update:submitDisabled', !valid || !isFormDataChanged)
        })
      },
      async handlePropertyConfigChange() {
        this.changeHandler()
      },
      async handleBasicChange() {
        this.changeHandler()
      }
    }
  })
</script>

<template>
  <div class="management-form">
    <div class="form-group">
      <cmdb-collapse :label="$t('基础信息')" arrow-type="filled">
        <grid-layout mode="form" :min-width="360" :max-width="560" :max-columns="2" :gap="24" class="form-content">
          <grid-item
            :label="$t('模板名称')"
            :label-width="120"
            required
            :class="['cmdb-form-item', { 'is-error': errors.has('templateName') }]">
            <bk-input type="text"
              name="templateName"
              font-size="normal"
              :placeholder="$t('模板名称将作为实例化后的模块名')"
              v-model.trim="formData.templateName"
              :data-vv-name="'templateName'"
              v-validate="'required|businessTopoInstNames|length:256'"
              @change="handleBasicChange">
            </bk-input>
            <p class="form-error">{{errors.first('templateName')}}</p>
          </grid-item>
          <grid-item
            :label="$t('服务分类')"
            :label-width="120"
            required>
            <div class="category-container">
              <div class="category-item" :class="['cmdb-form-item', { 'is-error': errors.has('primaryCategory') }]">
                <cmdb-selector
                  font-size="normal"
                  display-key="displayName"
                  :placeholder="$t('请选择一级分类')"
                  :searchable="true"
                  :auto-select="true"
                  :list="primaryCategories"
                  :popover-options="{
                    boundary: 'window'
                  }"
                  name="primaryCategory"
                  v-validate="'required'"
                  v-model="formData.primaryCategory"
                  @on-selected="handleSelectPrimaryCategory"
                  @change="handleBasicChange">
                  <template #default="{ name, id }">
                    <div class="bk-option-content-default" :title="`${name}（#${id}）`">
                      <div class="bk-option-name">
                        {{name}}<span class="category-id">（#{{id}}）</span>
                      </div>
                    </div>
                  </template>
                </cmdb-selector>
                <p class="form-error">{{errors.first('primaryCategory')}}</p>
              </div>
              <div class="category-item" :class="['cmdb-form-item', { 'is-error': errors.has('secCategory') }]">
                <cmdb-selector
                  font-size="normal"
                  display-key="displayName"
                  :placeholder="$t('请选择二级分类')"
                  :searchable="true"
                  :auto-select="true"
                  :list="currentSecCategories"
                  name="secCategory"
                  v-validate="'required'"
                  v-model="formData.secCategory"
                  @change="handleBasicChange">
                  <template #default="{ name, id }">
                    <div class="bk-option-content-default" :title="`${name}（#${id}）`">
                      <div class="bk-option-name">
                        {{name}}<span class="category-id">（#{{id}}）</span>
                      </div>
                    </div>
                  </template>
                </cmdb-selector>
                <p class="form-error">{{errors.first('secCategory')}}</p>
              </div>
            </div>
          </grid-item>
        </grid-layout>
      </cmdb-collapse>
    </div>
    <div class="form-group">
      <cmdb-collapse :label="$t('属性设置')" arrow-type="filled">
        <div class="form-content">
          <property-config
            ref="$propertyConfig"
            :properties="moduleProperties"
            :property-groups="modulePropertyGroup"
            :config="formData.propertyConfig"
            :selected="configProperties"
            :exclude="excludeModuleProperties"
            :max-columns="2"
            form-element-font-size="normal"
            @change="handlePropertyConfigChange">
            <template #tips>
              <div class="property-config-tips">{{$t('模板里定义的字段，在实例中将不可修改')}}</div>
            </template>
          </property-config>
        </div>
      </cmdb-collapse>
    </div>
    <div class="form-group">
      <cmdb-collapse :label="$t('服务进程')" arrow-type="filled">
        <div class="form-content">
          <div class="process-create-container">
            <cmdb-auth :auth="auth">
              <bk-button slot-scope="{ disabled }" v-test-id="'createProcess'"
                class="create-btn"
                theme="default"
                :disabled="disabled"
                @click="handleCreateProcess">
                <i class="bk-icon icon-plus"></i>
                <span>{{$t('新建进程')}}</span>
              </bk-button>
            </cmdb-auth>
            <span class="create-tips">{{$t('新建进程提示')}}</span>
          </div>
          <process-table
            v-if="formData.processList.length"
            :loading="$loading(requestIds.processList)"
            :properties="processProperties"
            :auth="auth"
            :show-operation="true"
            @on-edit="handleUpdateProcess"
            @on-delete="handleDeleteProcess"
            :list="formData.processList">
          </process-table>
        </div>
      </cmdb-collapse>
    </div>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="processSlider.show"
      :title="processSlider.title"
      :width="800"
      :before-close="handleProcessSliderBeforeClose">
      <template slot="content" v-if="processSlider.show">
        <process-form v-test-id.businessServiceTemplate="'processForm'"
          ref="$processForm"
          :auth="auth"
          :properties="processProperties"
          :property-groups="processPropertyGroup"
          :inst="processSlider.form.inst"
          :type="processSlider.form.type"
          :data-index="processSlider.form.dataIndex"
          :is-created-service="isLocalProcessEdit"
          :save-disabled="false"
          :submit-format="formatProcessSubmitData"
          @on-cancel="handleCancelProcess">
        </process-form>
      </template>
    </bk-sideslider>
  </div>
</template>

<style lang="scss" scoped>
  .management-form {
    .form-group {
      background: #fff;
      box-shadow: 0 2px 4px 0 rgba(25, 25, 41, 0.05);
      border-radius: 2px;
      padding: 16px 24px;

      & + .form-group {
        margin-top: 16px;
      }
    }

    .category-container {
      display: flex;
      .category-item {
        flex: 1;

        .bk-select {
          width: 100%;
        }

        & + .category-item {
          margin-left: 8px;
        }
      }
    }

    .form-content {
      padding: 24px 90px 12px 90px;
    }

    .process-create-container {
      display: flex;
      align-items: center;
      padding-bottom: 14px;

      .create-tips {
        color: #63656E;
        font-size: 12px;
        padding-left: 8px;
      }
    }

    .property-config-tips {
      color: #63656E;
      font-size: 12px;
      padding-left: 8px;
    }
  }
</style>
