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
  import { computed, defineComponent, reactive, ref, toRefs, watch, watchEffect } from 'vue'
  import isEqual from 'lodash/isEqual'
  import store from '@/store'
  import { t } from '@/i18n'
  import { OPERATION } from '@/dictionary/iam-auth'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import PropertyConfig from '@/components/property-config/index.vue'
  import TemplateTree from './template-tree.vue'
  import useTemplateData from './use-template-data'
  import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants'

  export default defineComponent({
    components: {
      GridLayout,
      GridItem,
      PropertyConfig,
      TemplateTree
    },
    props: {
      dataId: {
        type: Number
      },
      submitDisabled: Boolean
    },
    setup(props, { emit }) {
      const $propertyConfig = ref(null)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => Number(props.dataId))

      const isEdit = computed(() => templateId.value > 0)

      const auth = computed(() => {
        if (!isEdit.value) {
          return {
            type: OPERATION.C_SET_TEMPLATE,
            relation: [bizId.value]
          }
        }
        return {
          type: OPERATION.U_SET_TEMPLATE,
          relation: [bizId.value, templateId.value]
        }
      })

      const state = reactive({
        setProperties: [],
        setPropertyGroup: [],
        configProperties: [],
        propertyConfig: {},
        treeMode: isEdit.value ? 'edit' : 'create'
      })

      // 集群拓扑服务模板是否变更
      const isServiceTemplateChanged = ref(false)

      // 用于接收初始化时的表单数据
      let formDataCopy = null

      // 表单填写的数据项
      const formData = ref({
        templateName: '',
        propertyConfig: {},
        serviceTemplates: []
      })

      const excludeModuleProperties = [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.SET].NAME]

      watchEffect(async () => {
        const {
          templateName,
          setProperties,
          setPropertyGroup,
          configProperties,
          propertyConfig,
          formDataCopy: formDataCopyRaw
        } = await useTemplateData(bizId.value, templateId.value, isEdit.value)

        formDataCopy = formDataCopyRaw

        emit('data-loaded')

        // 属性数据主要提供给不同功能的组件使用
        state.setProperties = setProperties
        state.setPropertyGroup = setPropertyGroup

        // 编辑态数据初始化
        if (isEdit.value) {
          formData.value.templateName = templateName

          // 属性设置
          formData.value.propertyConfig = propertyConfig
          state.configProperties = configProperties

          store.commit('setTitle', `${t('编辑模板')}【${templateName}】`)
        }
      })

      const getPropertyConfigData = () => {
        const propertyConfigData = $propertyConfig.value?.getData()

        const attributes = []
        for (const [key, value] of Object.entries(propertyConfigData)) {
          attributes.push({
            bk_attribute_id: Number(key),
            bk_property_value: value
          })
        }

        return attributes
      }

      const isFormDataChanged = computed(() => {
        const formDataLatest = {
          templateName: formData.value.templateName,
          propertyConfig: $propertyConfig.value?.getData()
        }
        return !isEqual(formDataCopy, formDataLatest) || isServiceTemplateChanged.value
      })

      // 动态监测表单值是否变化，设置提交按钮的状态
      watch(isFormDataChanged, changed => emit('update:submitDisabled', !changed))

      return {
        ...toRefs(state),
        formData,
        templateId,
        auth,
        excludeModuleProperties,
        $propertyConfig,
        getPropertyConfigData,
        isFormDataChanged,
        isServiceTemplateChanged
      }
    },
    methods: {
      getData() {
        return {
          name: this.formData.templateName,
          service_template_ids: this.formData.serviceTemplates.map(template => template.id),
          attributes: this.getPropertyConfigData()
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
        const configValid = await this.$refs.$propertyConfig.validate()

        return basicValid && configValid
      },
      async changeHandler() {
        this.$nextTick(async () => {
          const valid = await this.validate()
          const { isFormDataChanged } = this
          this.$emit('update:submitDisabled', !valid || !isFormDataChanged)
        })
      },
      handleServiceTemplateSelected(services) {
        this.formData.serviceTemplates = services
        this.changeHandler()
      },
      handleServiceTemplateChange(value) {
        this.isServiceTemplateChanged = value
        this.changeHandler()
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
        <grid-layout mode="form" :min-width="360" :max-width="560" class="form-content">
          <grid-item
            :label="$t('模板名称')"
            :label-width="120"
            required
            :class="['cmdb-form-item', { 'is-error': errors.has('templateName') }]">
            <bk-input type="text"
              name="templateName"
              font-size="normal"
              v-model.trim="formData.templateName"
              v-validate="'required|singlechar|length:256'"
              :placeholder="$t('请输入xx', { name: $t('模板名称') })"
              :data-vv-name="'templateName'"
              @change="handleBasicChange">
            </bk-input>
            <p class="form-error" v-if="errors.has('templateName')">{{errors.first('templateName')}}</p>
          </grid-item>
        </grid-layout>
      </cmdb-collapse>
    </div>
    <div class="form-group">
      <cmdb-collapse :label="$t('属性设置')" arrow-type="filled">
        <div class="form-content">
          <property-config
            ref="$propertyConfig"
            :properties="setProperties"
            :property-groups="setPropertyGroup"
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
      <cmdb-collapse :label="$t('集群拓扑')" arrow-type="filled">
        <div class="form-content">
          <div :class="['cmdb-form-item', { 'is-error': errors.has('serviceTemplates') }]">
            <template-tree
              :mode="treeMode"
              :template-id="templateId"
              @service-change="handleServiceTemplateChange"
              @service-selected="handleServiceTemplateSelected">
            </template-tree>
            <input type="hidden"
              :value="formData.serviceTemplates.length"
              v-validate="'min_value:1'"
              name="serviceTemplates">
            <p class="form-error" v-if="errors.has('serviceTemplates')">{{$t('请添加服务模板')}}</p>
          </div>
        </div>
      </cmdb-collapse>
    </div>
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

    .form-content {
      padding: 24px 90px 12px 90px;
    }

    .property-config-tips {
      color: #63656E;
      font-size: 12px;
      padding-left: 8px;
    }
  }
</style>
