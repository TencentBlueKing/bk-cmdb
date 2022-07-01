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
  import { $success } from '@/magicbox/index.js'
  import { OPERATION } from '@/dictionary/iam-auth'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import PropertyConfigDetails from '@/components/property-config/details.vue'
  import TemplateTree from './template-tree.vue'
  import useTemplateData from './use-template-data'
  import setTemplateService from '@/services/set-template'
  import { MENU_BUSINESS_SET_TEMPLATE_EDIT } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      GridLayout,
      GridItem,
      PropertyConfigDetails,
      TemplateTree
    },
    setup(props, { emit }) {
      const $this = getCurrentInstance()
      const $templateName = ref(null)
      const loading = ref(true)

      const createElement = h.bind($this)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => parseInt(router.app.$route.params.templateId, 10))

      const state = reactive({
        templateName: '',
        setProperties: [],
        setPropertyGroup: [],
        configProperties: [],
        propertyConfig: {},
      })

      // 模板编辑权限定义
      const auth = computed(() => ({
        type: OPERATION.U_SET_TEMPLATE,
        relation: [bizId.value, templateId.value]
      }))

      watchEffect(async () => {
        const {
          templateName,
          setProperties,
          setPropertyGroup,
          configProperties,
          propertyConfig
        } = await useTemplateData(bizId.value, templateId.value, true)

        loading.value = false

        state.templateName = templateName
        state.setProperties = setProperties
        state.setPropertyGroup = setPropertyGroup
        state.configProperties = configProperties

        state.propertyConfig = propertyConfig

        store.commit('setTitle', `${t('模板详情')}【${templateName}】`)
      })

      // 当前编辑属性
      const editState = ref({
        property: null,
        value: ''
      })

      const hasPropertyConfig = computed(() => Object.keys(state.propertyConfig).length > 0)

      // 在处理中的属性列表
      const loadingState = ref([])

      // 定义基础属性，用于标准化编辑状态展示
      const basicProperties = {
        templateName: { bk_property_id: 'templateName' }
      }

      // 设置编辑状态数据
      const setEditState = (property) => {
        let $component = null
        if (property === basicProperties.templateName) {
          editState.value.value = state.templateName
          $component = $templateName
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
        if (editState.value.property) {
          confirmSaveName()
        }
      }

      const confirmSaveName = async () => {
        const valid = await $this?.proxy?.$validator.validate('templateName')
        if (!valid) {
          return
        }

        if (state.templateName === editState.value.value) {
          resetEditState()
          return
        }

        saveName()
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

          await store.dispatch('setTemplate/updateSetTemplate', {
            bizId: bizId.value,
            setTemplateId: templateId.value,
            params: {
              name: templateName
            }
          })

          // 回显为保存后的值
          state.templateName = templateName

          // TODO保存成功后提示
        } finally {
          loadingState.value = loadingState.value.filter(item => item !== basicProperties.templateName)
        }
      }

      // 属性设置loaidng队列，元素为属性对象
      const propertyConfigLoadingState = ref([])

      // 显示同步提示的方法
      const showSyncInstanceTips = (text = '成功更新模板，您可以通过XXX') => {
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

      // 属性设置-保存
      const handleSavePropertyConfig = async ({ property, value }) => {
        let saveValue = value
        const { bk_property_type: propertyType } = property
        if (['int', 'float'].includes(propertyType)) {
          saveValue = Number(value)
        }
        try {
          propertyConfigLoadingState.value.push(property)
          const data = {
            id: templateId.value,
            bk_biz_id: bizId.value,
            attributes: [{
              bk_attribute_id: property.id,
              bk_property_value: saveValue
            }]
          }
          await setTemplateService.updateProperty(data)

          state.propertyConfig[property.id] = saveValue

          showSyncInstanceTips()
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
        await setTemplateService.deleteProperty(data)

        del(state.propertyConfig, property.id)

        showSyncInstanceTips()
      }

      const handleGoToEdit = () => {
        routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_EDIT,
          params: {
            templateId: templateId.value
          },
          history: true
        })
      }

      return {
        ...toRefs(state),
        templateId,
        auth,
        loading,
        editState,
        loadingState,
        basicProperties,
        propertyConfigLoadingState,
        $templateName,
        hasPropertyConfig,
        setEditState,
        handleSavePropertyConfig,
        handleDelPropertyConfig,
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
          <grid-layout mode="detail" :min-width="360" :max-width="560" class="form-content">
            <grid-item
              :label="$t('模板名称')"
              :label-width="160"
              :class="['cmdb-form-item', { 'is-error': errors.has('templateName') }]">
              <div class="editable-content">
                <div
                  :class="['basic-value', { 'is-loading': loadingState.includes(basicProperties.templateName) }]"
                  v-if="basicProperties.templateName !== editState.property">
                  {{templateName}}
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
                      v-model.trim="editState.value"
                      :data-vv-name="'templateName'"
                      v-validate="'required|singlechar|length:256'"
                      :placeholder="$t('请输入xx', { name: $t('模板名称') })"
                      @enter="handleSaveName"
                      @blur="handleSaveName">
                    </bk-input>
                    <p class="form-error">{{errors.first('templateName')}}</p>
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
              :properties="setProperties"
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
                        class="link"
                        :disabled="disabled"
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
        <cmdb-collapse :label="$t('集群拓扑')" arrow-type="filled">
          <div class="form-content">
            <template-tree
              :mode="'view'"
              :template-id="templateId">
            </template-tree>
          </div>
        </cmdb-collapse>
      </div>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="auth">
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
