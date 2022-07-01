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
  import { computed, defineComponent, del, reactive, ref, toRefs, watchEffect, getCurrentInstance, nextTick, h } from '@vue/composition-api'
  import { t } from '@/i18n'
  import router from '@/router/index.js'
  import store from '@/store'
  import routerActions from '@/router/actions'
  import { $bkInfo, $success } from '@/magicbox/index.js'
  import { formatValue } from '@/utils/tools'
  import { OPERATION } from '@/dictionary/iam-auth'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import ProcessForm from './process-form.vue'
  import ProcessTable from './process.vue'
  import PropertyConfigDetails from '@/components/property-config/details.vue'
  import useTemplateData, { templateDetailRequestId } from './use-template-data'
  import serviceTemplateService from '@/services/service-template'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_EDIT } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      GridLayout,
      GridItem,
      ProcessForm,
      ProcessTable,
      PropertyConfigDetails
    },
    setup(props, { emit }) {
      const $this = getCurrentInstance()
      const $processForm = ref(null)
      const $templateName = ref(null)
      const $secCategory = ref(null)
      const loading = ref(true)

      const createElement = h.bind($this)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => parseInt(router.app.$route.params.templateId, 10))

      const state = reactive({
        moduleProperties: [],
        modulePropertyGroup: [],
        processProperties: [],
        processPropertyGroup: [],
        fullCategories: [],
        primaryCategories: [],
        secCategories: [],
        currentSecCategories: [],
        basic: {} as any,
        propertyConfig: {},
        processList: [],
        processSlider: {
          show: false,
          title: '',
          form: {
            inst: {},
            dataIndex: null,
            type: ''
          }
        },
        requestIds: {
          processList: Symbol()
        }
      })

      // 服务模板编辑权限定义
      const auth = computed(() => ({
        type: OPERATION.U_SERVICE_TEMPLATE,
        relation: [bizId.value, templateId.value]
      }))

      watchEffect(async () => {
        const {
          moduleProperties,
          modulePropertyGroup,
          processProperties,
          processPropertyGroup,
          basic,
          propertyConfig,
          processList,
          primaryCategories,
          secCategories
        } = await useTemplateData(bizId.value, templateId.value, true)

        loading.value = false

        state.moduleProperties = moduleProperties
        state.modulePropertyGroup = modulePropertyGroup
        state.processProperties = processProperties
        state.processPropertyGroup = processPropertyGroup

        state.primaryCategories = primaryCategories
        state.secCategories = secCategories
        state.currentSecCategories = state.secCategories
          .filter(category => category.bk_parent_id === basic.primaryCategory)

        state.basic = basic
        state.basic.primaryCategoryCopy = basic.primaryCategory
        state.propertyConfig = propertyConfig
        state.processList = processList

        store.commit('setTitle', `${t('模板详情')}【${state.basic.templateName}】`)
      })

      // 获取进程列表
      const getProcessList = async () => {
        const data = await store.dispatch('processTemplate/getBatchProcessTemplate', {
          params: {
            bk_biz_id: bizId.value,
            service_template_id: templateId.value
          },
          config: {
            requestId: state.requestIds.processList
          }
        })
        state.processList = data.info.map(template => ({
          process_id: template.id,
          ...template.property
        })).sort((prev, next) => prev.process_id - next.process_id)
      }

      const serviceCategory = computed(() => {
        const primary = state.primaryCategories.find(category => category.id === state.basic.primaryCategory) || {}
        const second = state.secCategories.find(category => category.id === state.basic.secCategory) || {}
        return `${primary.name || '--'} / ${second.name || '--'}`
      })

      const hasPropertyConfig = computed(() => Object.keys(state.propertyConfig).length > 0)

      const handleChangePrimaryCategory = (id: number) => {
        state.currentSecCategories = state.secCategories.filter(category => category.bk_parent_id === id)
        if (editState.value.property === basicProperties.categorty) {
          editState.value.value = ''
        }
      }

      // 点击空白处取消分类编辑
      const handleCategoryClickOutside = () => {
        if (editState.value.property === basicProperties.categorty) {
          // 还原分类数据
          state.basic.primaryCategory = state.basic.primaryCategoryCopy
          state.currentSecCategories = state.secCategories
            .filter(category => category.bk_parent_id === state.basic.primaryCategoryCopy)

          resetEditState()
        }
      }
      // 防止点击分类下拉框时退出编辑
      const categoryClickOutsideMiddleware = event => (!event.target.closest('.bk-select-dropdown-content'))

      // 当前编辑属性
      const editState = ref({
        property: null,
        value: ''
      })

      // 在处理中的属性列表
      const loadingState = ref([])

      // 保存信息确认中状态
      const saveNameConfirming = ref(false)

      // 定义基础属性，用于标准化编辑状态展示
      const basicProperties = {
        templateName: { bk_property_id: 'templateName' },
        categorty: { bk_property_id: 'categorty' }
      }

      // 设置编辑状态数据
      const setEditState = (property) => {
        let $component = null
        if (property === basicProperties.templateName) {
          editState.value.value = state.basic.templateName
          $component = $templateName
        }
        if (property === basicProperties.categorty) {
          editState.value.value = state.basic.secCategory
          $component = $secCategory
        }
        editState.value.property = property

        nextTick(() => {
          $component?.value?.focus?.()
        })
      }

      // 重置编辑状态数据，用于取消编辑态
      const resetEditState = () => {
        editState.value.property = null
        editState.value.value = ''
      }

      // 名称回车或失焦事件回调，触发保存
      const handleSaveName = () => {
        if (editState.value.property && !saveNameConfirming.value) {
          confirmSaveName()
        }
      }

      const handleChangeSecCategory = () => {
        // 使用nextTick等待validate后执行，确保校验状态正确性
        nextTick(saveCategory)
      }

      const confirmSaveName = async () => {
        const valid = await $this?.proxy?.$validator.validate('templateName')
        if (!valid) {
          return
        }

        if (state.basic.templateName === editState.value.value) {
          resetEditState()
          return
        }

        $bkInfo({
          title: t('确认修改名称'),
          subTitle: t('确认修改名称提示'),
          width: 520,
          extCls: 'confirm-edit-service-template-name-infobox',
          async confirmFn() {
            saveName()
          },
          cancelFn() {
            resetEditState()
          },
          stateChangeFn(isShow) {
            saveNameConfirming.value = isShow
          },
          confirmLoading: true
        })
      }
      const saveName = async () => {
        try {
          const valid = await $this?.proxy?.$validator.validate('templateName')
          if (!valid) {
            return
          }

          // 先取出编辑后的值
          const { value: templateName } = editState.value

          // 重置编辑态，回到详情状态
          resetEditState()

          // 设置loading状态
          loadingState.value.push(basicProperties.templateName)

          await store.dispatch('serviceTemplate/updateServiceTemplate', {
            params: {
              bk_biz_id: bizId.value,
              id: templateId.value,
              name: templateName
            }
          })

          // 回显为保存后的值
          state.basic.templateName = templateName
        } finally {
          loadingState.value = loadingState.value.filter(item => item !== basicProperties.templateName)
        }
      }
      const saveCategory = async () => {
        try {
          const valid = await $this?.proxy?.$validator.validate('secCategory')
          if (!valid) {
            return
          }

          const { value: categoryId } = editState.value

          if (state.basic.secCategory === categoryId) {
            resetEditState()
            return
          }

          resetEditState()

          loadingState.value.push(basicProperties.categorty)

          await store.dispatch('serviceTemplate/updateServiceTemplate', {
            params: {
              bk_biz_id: bizId.value,
              id: templateId.value,
              service_category_id: categoryId
            }
          })

          state.basic.secCategory = categoryId
          state.basic.primaryCategoryCopy = state.basic.primaryCategory

          $success(t('修改服务分类成功提示'))
        } finally {
          loadingState.value = loadingState.value.filter(item => item !== basicProperties.categorty)
        }
      }

      // 显示同步提示的方法
      const showSyncInstanceTips = (text = '成功更新模板进程，您可以通过XXX') => {
        const link = createElement('bk-link', {
          slot: 'link',
          props: { theme: 'primary' },
          on: {
            click() {
              emit('active-change', 'instance')
            }
          }
        }, t('同步功能'))

        const message = createElement('i18n', {
          class: 'process-success-message',
          props: {
            path: text,
            tag: 'div',
          }
        }, [link])

        $success(message)
        emit('sync-change')
      }

      const saveProcessAfter = () => {
        state.processSlider.show = false
        showSyncInstanceTips()
      }

      // 属性设置loaidng队列，元素为属性对象
      const propertyConfigLoadingState = ref([])

      // 属性设置-保存
      const handleSavePropertyConfig = async ({ property, value }) => {
        try {
          propertyConfigLoadingState.value.push(property)
          const data = {
            id: templateId.value,
            bk_biz_id: bizId.value,
            attributes: [{
              bk_attribute_id: property.id,
              bk_property_value: value
            }]
          }
          await serviceTemplateService.updateProperty(data)

          state.propertyConfig[property.id] = value

          showSyncInstanceTips('成功更新模板，您可以通过XXX')
        } finally {
          propertyConfigLoadingState.value = propertyConfigLoadingState.value.filter(item => item !== property)
        }
      }

      // 属性设置-删除
      const handleDelPropertyConfig = async (property) => {
        const data = {
          id: templateId.value,
          bk_biz_id: bizId.value,
          bk_attribute_ids: [property.id]
        }
        await serviceTemplateService.deleteProperty(data)

        del(state.propertyConfig, property.id)

        showSyncInstanceTips('成功更新模板，您可以通过XXX')
      }

      const handleCreateProcess = () => {
        state.processSlider.show = true
        state.processSlider.title = t('添加进程')
        state.processSlider.form.type = 'create'
        state.processSlider.form.inst = {}
      }

      const handleSaveProcess = async (values, changedValues, type) => {
        const data = type === 'create' ? values : changedValues
        const processValues = formatProcessSubmitData(data)
        if (type === 'create') {
          await store.dispatch('processTemplate/createProcessTemplate', {
            params: {
              bk_biz_id: bizId.value,
              service_template_id: templateId.value,
              processes: [{
                spec: processValues
              }]
            }
          })
        } else {
          await store.dispatch('processTemplate/updateProcessTemplate', {
            params: {
              bk_biz_id: bizId.value,
              process_template_id: values.process_id,
              process_property: processValues
            }
          })
        }

        getProcessList()
        saveProcessAfter()
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

      const handleDeleteProcess = (template) => {
        $bkInfo({
          title: t('确认删除模板进程'),
          confirmFn: async () => {
            await store.dispatch('processTemplate/deleteProcessTemplate', {
              params: {
                data: {
                  bk_biz_id: bizId.value,
                  process_templates: [template.process_id]
                }
              }
            })

            getProcessList()
            showSyncInstanceTips()
          },
          confirmLoading: true
        })
      }

      const formatProcessSubmitData = (data = {}) => {
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
        return data
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

      const handleGoToEdit = () => {
        routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_EDIT,
          params: {
            templateId: templateId.value
          },
          history: true
        })
      }

      return {
        ...toRefs(state),
        bizId,
        templateId,
        auth,
        loading,
        editState,
        loadingState,
        basicProperties,
        templateDetailRequestId,
        serviceCategory,
        propertyConfigLoadingState,
        $templateName,
        $secCategory,
        $processForm,
        hasPropertyConfig,
        setEditState,
        categoryClickOutsideMiddleware,
        formatProcessSubmitData,
        handleChangeSecCategory,
        handleChangePrimaryCategory,
        handleCategoryClickOutside,
        handleSavePropertyConfig,
        handleDelPropertyConfig,
        handleCreateProcess,
        handleSaveProcess,
        handleCancelProcess,
        handleUpdateProcess,
        handleDeleteProcess,
        handleProcessSliderBeforeClose,
        handleSaveName,
        handleGoToEdit
      }
    }
  })
</script>

<template>
  <cmdb-sticky-layout class="details-sticky-layout">
    <div class="template-config" v-bkloading="{ isLoading: loading }">
      <div class="form-group">
        <cmdb-collapse :label="$t('基础信息')" arrow-type="filled">
          <grid-layout mode="detail" :min-width="360" :max-width="560" :gap="0" class="form-content">
            <grid-item
              :label="$t('模板名称')"
              :label-width="160"
              :class="['cmdb-form-item', { 'is-error': errors.has('templateName') }]">
              <div class="editable-content">
                <div
                  :class="['basic-value', { 'is-loading': loadingState.includes(basicProperties.templateName) }]"
                  v-if="basicProperties.templateName !== editState.property">
                  {{basic.templateName}}
                </div>
                <template v-if="!loadingState.includes(basicProperties.templateName)">
                  <cmdb-auth
                    v-show="basicProperties.templateName !== editState.property"
                    tag="i"
                    class="icon-cc-edit-shape property-edit-button"
                    :auth="auth"
                    @click="setEditState(basicProperties.templateName)">
                  </cmdb-auth>
                  <div class="property-form" v-if="basicProperties.templateName === editState.property">
                    <bk-input type="text"
                      ref="$templateName"
                      name="templateName"
                      size="small"
                      font-size="normal"
                      :placeholder="$t('模板名称将作为实例化后的模块名')"
                      v-model.trim="editState.value"
                      :data-vv-name="'templateName'"
                      v-validate="'required|businessTopoInstNames|length:256'"
                      @enter="handleSaveName"
                      @blur="handleSaveName">
                    </bk-input>
                    <p class="form-error">{{errors.first('templateName')}}</p>
                  </div>
                </template>
              </div>
            </grid-item>
            <grid-item
              label="服务分类"
              direction="row"
              :label-width="160">
              <div class="editable-content">
                <div
                  :class="['basic-value', { 'is-loading': loadingState.includes(basicProperties.categorty) }]"
                  v-if="basicProperties.categorty !== editState.property">
                  {{serviceCategory}}
                </div>
                <template v-if="!loadingState.includes(basicProperties.categorty)">
                  <cmdb-auth
                    v-show="basicProperties.categorty !== editState.property"
                    tag="i"
                    class="icon-cc-edit-shape property-edit-button"
                    :auth="auth"
                    @click="setEditState(basicProperties.categorty)">
                  </cmdb-auth>
                  <div class="category-container"
                    v-if="basicProperties.categorty === editState.property"
                    v-click-outside="{
                      handler: handleCategoryClickOutside,
                      middleware: categoryClickOutsideMiddleware
                    }">
                    <div :class="['category-item', 'cmdb-form-item', { 'is-error': errors.has('primaryCategory') }]">
                      <cmdb-selector
                        display-key="displayName"
                        size="small"
                        font-size="normal"
                        :placeholder="$t('请选择一级分类')"
                        :searchable="true"
                        :auto-select="false"
                        :list="primaryCategories"
                        :popover-options="{
                          boundary: 'window'
                        }"
                        name="primaryCategory"
                        v-validate="'required'"
                        v-model="basic.primaryCategory"
                        @change="handleChangePrimaryCategory">
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
                        ref="$secCategory"
                        display-key="displayName"
                        size="small"
                        font-size="normal"
                        :placeholder="$t('请选择二级分类')"
                        :searchable="true"
                        :auto-select="false"
                        :list="currentSecCategories"
                        name="secCategory"
                        v-validate="'required'"
                        v-model="editState.value"
                        @change="handleChangeSecCategory">
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
                </template>
              </div>
            </grid-item>
          </grid-layout>
        </cmdb-collapse>
      </div>
      <div class="form-group">
        <cmdb-collapse :label="$t('属性设置')" arrow-type="filled">
          <div class="form-content">
            <property-config-details v-if="hasPropertyConfig"
              :instance="propertyConfig"
              :properties="moduleProperties"
              :auth="auth"
              :loading-state="propertyConfigLoadingState"
              :max-columns="2"
              form-element-size="small"
              @save="handleSavePropertyConfig"
              @del="handleDelPropertyConfig">
            </property-config-details>
            <div class="property-config-empty" v-else-if="!loading">
              <i class="icon icon-cc-tips"></i>
              <cmdb-auth :auth="auth">
                <template #default="{ disabled }">
                  <i18n path="当前模板未配置提示">
                    <template #link>
                      <bk-link
                        theme="primary"
                        :disabled="disabled"
                        class="link"
                        @click="handleGoToEdit">
                        {{$t('立即配置')}}
                      </bk-link>
                    </template>
                  </i18n>
                </template>
              </cmdb-auth>
            </div>
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
              v-if="processList.length"
              :loading="$loading([templateDetailRequestId, requestIds.processList])"
              :properties="processProperties"
              :auth="auth"
              :show-operation="true"
              @on-edit="handleUpdateProcess"
              @on-delete="handleDeleteProcess"
              :list="processList">
            </process-table>
          </div>
        </cmdb-collapse>
      </div>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_TEMPLATE, relation: [bizId, templateId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="disabled"
            @click="handleGoToEdit">
            {{$t('编辑')}}
          </bk-button>
        </cmdb-auth>
      </div>
    </template>

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
          :is-created-service="false"
          :save-disabled="false"
          :submit-format="formatProcessSubmitData"
          @on-submit="handleSaveProcess"
          @on-cancel="handleCancelProcess">
        </process-form>
      </template>
    </bk-sideslider>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
.template-config {
  padding: 15px 20px 0 20px;

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

    .property-form {
      width: 100%;
    }
  }

  .editable-content {
    display: flex;
    align-items: center;

    &:hover {
      .property-edit-button {
        display: block;
      }
    }

    .property-edit-button {
      display: none;
      font-size: 16px;
      margin-left: 8px;
      cursor: pointer;

      &:hover {
        color: $primaryColor;
      }
    }
  }

  .basic-value {
    font-size: 12px;

    &.is-loading {
      font-size: 0;
      &:before {
        content: "";
        display: inline-block;
        width: 16px;
        height: 16px;
        margin: 2px 0;
        background-image: url("@/assets/images/icon/loading.svg");
      }
    }
  }

  .category-container {
    display: flex;
    width: 100%;
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

  .property-config-empty {
    font-size: 12px;
    display: flex;
    align-items: center;
    .icon {
      font-size: 14px;
      margin-right: 4px;
    }
    .link {
      line-height: normal;
      vertical-align: unset;
      ::v-deep .bk-link-text {
        font-size: 12px;
      }
    }
  }

}
.process-success-message {
  .bk-link {
    vertical-align: baseline;
  }
}

.details-sticky-layout {
  height: 100%;
  overflow-y: auto;

  .layout-footer {
    display: flex;
    align-items: center;
    height: 52px;
    padding: 0 20px;
    margin-top: 8px;
    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
    }
    .auth-box + .bk-button {
      margin-left: 8px;
    }
    &.is-sticky {
      background-color: #fff;
      border-top: 1px solid $borderColor;
    }
  }
}
</style>
<style lang="scss">
  .confirm-edit-service-template-name-infobox {
    .bk-dialog-sub-header {
      .bk-dialog-header-inner {
        text-align: left !important;
      }
    }
  }
</style>
