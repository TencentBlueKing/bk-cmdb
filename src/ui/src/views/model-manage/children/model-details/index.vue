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
  <div class="model-detail-wrapper">
    <div class="model-info-wrapper">
      <div class="model-info" v-bkloading="{ isLoading: $loading('getClassificationsObjectStatistics') }">
        <template v-if="activeModel !== null">
          <div class="choose-icon-wrapper">
            <span class="model-type" :class="{ 'is-builtin': activeModel.ispre }">{{getModelType()}}</span>
            <template v-if="isEditable">
              <cmdb-auth tag="div" class="icon-box"
                v-if="!activeModel.bk_ispaused"
                :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
                @click="isIconListShow = true">
                <i class="icon" :class="activeModel.bk_obj_icon || 'icon-cc-default'"></i>
                <p class="hover-text is-paused" v-if="activeModel.bk_ispaused">{{$t('已停用')}}</p>
                <p class="hover-text" v-else>{{$t('点击切换')}}</p>
              </cmdb-auth>
              <div class="choose-icon-box" v-if="isIconListShow" v-click-outside="hideChooseBox">
                <the-choose-icon
                  v-model="activeModel.bk_obj_icon"
                  @close="hideChooseBox"
                  @input="handleModelIconUpdateConfirm">
                </the-choose-icon>
              </div>
            </template>
            <template v-else>
              <div class="icon-box" :class="{ 'is-builtin': activeModel.ispre }" style="cursor: default;">
                <i class="icon" :class="activeModel.bk_obj_icon || 'icon-cc-default'"></i>
              </div>
            </template>
          </div>
          <div class="model-identity">
            <div class="model-name">
              <editable-field
                :editing.sync="modelNameIsEditing"
                v-model="activeModel.bk_obj_name"
                font-size="12px"
                @confirm="handleModelNameUpdateConfirm"
                :editable="isEditable"
                validate="required|singlechar|length:256|reservedWord"
                :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }">
                <template #append>
                  <bk-tag v-if="activeModel.bk_ispaused" size="small" theme="default">{{$t('已停用')}}</bk-tag>
                </template>
              </editable-field>
            </div>
            <div class="model-id" v-show="!modelNameIsEditing" v-bk-overflow-tips>
              {{activeModel['bk_obj_id'] || ''}}
            </div>
          </div>
          <div class="model-group-name">
            <span :class="['model-group-name-label', { 'model-group-name-label-editing': modelGroupIsEditing }]">
              {{$t('所属分组')}}
            </span>
            <editable-field
              :editing.sync="modelGroupIsEditing"
              class="model-group-name-edit"
              v-model="activeModel.bk_classification_id"
              :label="modelClassificationName"
              :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
              validate="required"
              @confirm="handleModelGroupUpdateConfirm"
              type="enum"
              font-size="12px"
              style="width: calc(100% - 60px)"
              :options="classifications
                .map(item => ({ id: item.bk_classification_id, name: item.bk_classification_name }))"
            >
            </editable-field>
          </div>
          <div class="instance-count"
            v-if="!activeModel['bk_ispaused'] && !isNoInstanceModel">
            <span class="instance-count-label">{{$t('实例数量')}}</span>
            <div>
              <span class="instance-count-text" @click="handleGoInstance">
                <cmdb-loading :loading="$loading(request.instanceCount)">
                  {{modelInstanceCount || 0}}
                </cmdb-loading>
              </span>
            </div>
          </div>
          <div class="field-template"
            v-if="!activeModel['bk_ispaused'] && !isNoInstanceModel">
            <span class="field-template-label">{{$t('绑定的字段组合模板')}}</span>
            <flex-tag
              v-if="templateList.length"
              class="field-template-tag"
              :max-width="'355px'"
              :list="templateList"
              :is-link-style="true"
              :popover-options="{
                boundary: 'scrollParent',
                appendTo: 'parent'
              }"
              @click-text="handleViewTemplate">
              <template #append="template">
                <cmdb-auth
                  tag="i"
                  class="unbind-icon icon-cc-unbind"
                  v-bk-tooltips="$t('解绑模版')"
                  :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
                  @click="handleUnbindTemplate(template)">
                </cmdb-auth>
              </template>
              <template #text-append="template">
                <i class="reddot"
                  v-if="templateDiffStatus[template.id] && templateDiffStatus[template.id].need_sync"
                  v-bk-tooltips="{
                    allowHTML: true,
                    theme: 'light template-diff-sync',
                    content: `#template-diff-sync-tooltips-${template.id}`
                  }">
                </i>
                <div :id="`template-diff-sync-tooltips-${template.id}`"
                  class="diff-sync-content"
                  v-if="templateDiffStatus[template.id] && templateDiffStatus[template.id].need_sync">
                  <i18n path="模型信息与模板信息有差异提示语" tag="div" class="content-tips">
                    <template #link>
                      <bk-link theme="primary" @click="handleGoSync(template)">{{ $t('去同步') }}</bk-link>
                    </template>
                  </i18n>
                </div>
              </template>
            </flex-tag>
            <div v-else>--</div>
          </div>
          <cmdb-auth class="restart-btn"
            v-if="!isMainLineModel && activeModel.bk_ispaused"
            :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }">
            <bk-button slot-scope="{ disabled }"
              theme="primary"
              :disabled="disabled"
              @click="dialogConfirm('restart')">
              {{$t('立即启用')}}
            </bk-button>
          </cmdb-auth>
          <div class="btn-group">
            <template v-if="isShowOperationButton">
              <cmdb-auth class="label-btn"
                v-if="!isMainLineModel && !activeModel['bk_ispaused']"
                v-bk-tooltips="$t('保留模型和相应实例，隐藏关联关系')"
                :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }">
                <bk-button slot-scope="{ disabled }"
                  text
                  :disabled="disabled"
                  @click="dialogConfirm('stop')">
                  <i
                    class="label-btn-icon bk-icon icon-minus-circle-shape">
                  </i>
                  <span class="label-btn-text">{{$t('停用')}}</span>
                </bk-button>
              </cmdb-auth>
              <cmdb-auth class="label-btn"
                v-bk-tooltips="$t('删除模型和其下所有实例，此动作不可逆，请谨慎操作')"
                :auth="{ type: $OPERATION.D_MODEL, relation: [modelId] }">
                <bk-button slot-scope="{ disabled }"
                  text
                  :disabled="disabled"
                  @click="dialogConfirm('delete')">
                  <i class="label-btn-icon icon-cc-del"></i>
                  <span class="label-btn-text">{{$t('删除')}}</span>
                </bk-button>
              </cmdb-auth>
            </template>
          </div>
        </template>
      </div>
    </div>
    <bk-tab class="model-details-tab" type="unborder-card"
      :active.sync="tab.active"
      @tab-change="handleTabChange">
      <bk-tab-panel name="field" :label="$t('模型字段')">
        <the-field-group
          ref="field"
          v-if="tab.active === 'field'"
          :is-read-only-import="isReadOnly"
          :can-be-import="canBeImport"
          :hide-import="hideImport"
          :import-auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
          @handleImportField="handleImportField"
          @exportField="exportField">
        </the-field-group>
      </bk-tab-panel>
      <bk-tab-panel name="relation" :label="$t('模型关联')" :visible="!!activeModel">
        <the-relation v-if="tab.active === 'relation'" :model-id="modelId"></the-relation>
      </bk-tab-panel>
      <bk-tab-panel name="verification" :label="$t('唯一校验')">
        <the-verification v-if="tab.active === 'verification'" :model-id="modelId"></the-verification>
      </bk-tab-panel>
    </bk-tab>

    <!-- 导入字段 -->
    <bk-sideslider
      v-transfer-dom
      :is-show.sync="importField.show"
      :width="800"
      :title="$t('导入字段')"
      @hidden="handleSliderHide"
    >
      <cmdb-import
        slot="content"
        v-if="importField.show"
        :template-url="importField.templateUrl"
        :import-url="importUrl"
        @upload-done="handleUploadDone"
      >
        <div slot="uploadResult">
          <div class="upload-details-success" v-if="uploadResult.success && uploadResult.success.length">
            <i class="bk-icon icon-check-circle-shape"></i>
            <span>{{$t('成功导入N个字段', { N: uploadResult.success.length })}}</span>
          </div>
          <div class="upload-details-fail" v-if="uploadResult.insert_failed && uploadResult.insert_failed.length">
            <div class="upload-details-fail-title">
              <i class="bk-icon icon-close-circle-shape"></i>
              <span>{{$t('新增失败列表')}}({{uploadResult.insert_failed.length}})</span>
            </div>
            <ul ref="failList" class="upload-details-fail-list">
              <li
                v-for="(fail, index) in uploadResult.insert_failed"
                :title="$t('第N行字段错误信息', { N: fail.row, field: fail.bk_property_id, info: fail.info })"
                :key="index">{{$t('第N行字段错误信息', { N: fail.row, field: fail.bk_property_id, info: fail.info })}}
              </li>
            </ul>
          </div>
          <div class="upload-details-fail" v-if="uploadResult.update_failed && uploadResult.update_failed.length">
            <div class="upload-details-fail-title">
              <i class="bk-icon icon-close-circle-shape"></i>
              <span>{{$t('更新失败列表')}}({{uploadResult.update_failed.length}})</span>
            </div>
            <ul ref="failList" class="upload-details-fail-list">
              <li
                v-for="(fail, index) in uploadResult.update_failed"
                :title="$t('第N行字段错误信息', { N: fail.row, field: fail.bk_property_id, info: fail.info })"
                :key="index">{{$t('第N行字段错误信息', { N: fail.row, field: fail.bk_property_id, info: fail.info })}}
              </li>
            </ul>
          </div>
          <div class="upload-details-fail" v-if="uploadResult.errors && uploadResult.errors.length">
            <div class="upload-details-fail-title">
              <i class="bk-icon icon-close-circle-shape"></i>
              <span>{{$t('上传失败列表')}}({{uploadResult.errors.length}})</span>
            </div>
            <ul ref="failList" class="upload-details-fail-list">
              <li
                v-for="(fail, index) in uploadResult.errors"
                :title="$t('第N行字段错误信息', { N: fail.row, field: fail.bk_property_id, info: fail.info })"
                :key="index">{{$t('第N行字段错误信息', { N: fail.row, field: fail.bk_property_id, info: fail.info })}}
              </li>
            </ul>
          </div>
        </div>
      </cmdb-import>
    </bk-sideslider>
    <!-- /导入字段 -->
  </div>
</template>

<script>
  import has from 'has'
  import theRelation from './relation'
  import theVerification from './verification'
  import cmdbLoading from '@/components/loading/index.vue'
  import theFieldGroup from '@/components/model-manage/field-group'
  import theChooseIcon from '@/components/model-manage/choose-icon/_choose-icon'
  import cmdbImport from '@/components/import/import'
  import { mapActions, mapGetters, mapMutations } from 'vuex'
  import RouterQuery from '@/router/query'
  import CombineRequest from '@/api/combine-request.js'
  import modelImportExportService from '@/service/model/import-export'
  import {
    MENU_MODEL_MANAGEMENT,
    MENU_RESOURCE_INSTANCE,
    MENU_MODEL_FIELD_TEMPLATE,
    MENU_MODEL_FIELD_TEMPLATE_SYNC_MODEL
  } from '@/dictionary/menu-symbol'
  import { BUILTIN_MODEL_RESOURCE_MENUS, BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import EditableField from '@/components/ui/details/editable-field.vue'
  import FlexTag from '@/components/ui/flex-tag'
  import fieldTemplateService from '@/service/field-template'

  export default {
    name: 'ModelDetails',
    components: {
      theFieldGroup,
      theRelation,
      theVerification,
      theChooseIcon,
      cmdbImport,
      cmdbLoading,
      EditableField,
      FlexTag
    },
    data() {
      return {
        tab: {
          active: RouterQuery.get('tab', 'field')
        },
        isIconListShow: false,
        isEditName: false,
        modelInstanceCount: null,
        isEditClassification: false,
        modelStatisticsSet: {},
        importField: {
          show: false,
          templateUrl: ''
        },
        uploadResult: {
          success: null,
          errors: null,
          insert_failed: null,
          update_failed: null
        },
        request: {
          instanceCount: Symbol('instanceCount')
        },
        modelNameIsEditing: false,
        modelGroupIsEditing: false,
        templateList: [],
        templateDiffStatus: {}
      }
    },
    computed: {
      ...mapGetters([
        'supplierAccount',
        'userName'
      ]),
      ...mapGetters('objectMainLineModule', ['isMainLine']),
      ...mapGetters('objectModelClassify', ['models', 'classifications']),
      activeModel: {
        get() {
          return this.$store.getters['objectModel/activeModel']
        },
        set(value) {
          this.setActiveModel(value)
        }
      },
      isMainLineModel() {
        return this.isMainLine(this.activeModel)
      },
      isShowOperationButton() {
        return this.activeModel && !this.activeModel.ispre
      },
      isReadOnly() {
        if (this.activeModel) {
          return this.activeModel.bk_ispaused
        }
        return false
      },
      isEditable() {
        if (this.activeModel) {
          return !this.activeModel.ispre && !this.activeModel.bk_ispaused
        }
        return false
      },
      modelClassificationName() {
        return this.classifications
          .find(item => item.bk_classification_id === this.activeModel.bk_classification_id)?.bk_classification_name || ''
      },
      importUrl() {
        return `${window.API_HOST}object/object/${this.activeModel.bk_obj_id}/import`
      },
      canBeImport() {
        return !this.isMainLineModel
      },
      modelId() {
        const model = this.$store.getters['objectModelClassify/getModelById'](this.$route.params.modelId)
        return model.id || null
      },
      isNoInstanceModel() {
        // 不能直接查看实例的模型
        const noInstanceModelIds = ['set', 'module']
        return noInstanceModelIds.includes(this.activeModel.bk_obj_id)
      },
      hideImport() {
        // 项目模型中隐藏导入按钮
        return this.tab.active === 'field' && this.$route.params.modelId !== BUILTIN_MODELS.PROJECT
      }
    },
    watch: {
      '$route.params.modelId'() {
        this.initObject()
      },
      async templateList(list) {
        if (!list?.length) {
          return
        }
        const templateIds = list.map(item => item.id)
        const allResult = await CombineRequest.setup(Symbol(), (params) => {
          const [templateId] = params
          return fieldTemplateService.getModelDiffStatus({
            bk_template_id: templateId,
            object_ids: [this.activeModel.id]
          })
        }, { segment: 1, concurrency: 5 }).add(templateIds)

        let groupIndex = 0
        for (const result of allResult) {
          const results = await result
          for (let i = 0; i < results.length; i++) {
            const { status, reason, value } = results[i]
            if (status === 'rejected') {
              console.error(reason?.message)
              continue
            }
            this.$set(this.templateDiffStatus, templateIds[(groupIndex * 5) + i], value?.[0] ?? {})
          }
          groupIndex += 1
        }
      }
    },
    created() {
      this.initObject()
    },
    beforeDestroy() {
      this.$http.cancelRequest(this.request.instanceCount)
    },
    methods: {
      handleTabChange(tab) {
        RouterQuery.set({ tab })
      },
      ...mapActions('objectModel', [
        'searchObjects',
        'updateObject',
        'deleteObject'
      ]),
      ...mapActions('objectBatch', [
        'importObjectAttribute',
        'exportObjectAttribute'
      ]),
      ...mapActions('objectMainLineModule', [
        'deleteMainlineObject',
        'searchMainlineObject'
      ]),
      ...mapMutations('objectModel', [
        'setActiveModel'
      ]),
      getModelType() {
        if (this.activeModel) {
          return this.activeModel.ispre ? this.$t('内置') : this.$t('自定义')
        }
        return ''
      },
      async handleFile(e) {
        const { files } = e.target
        const formData = new FormData()
        formData.append('file', files[0])
        try {
          const res = await this.importObjectAttribute({
            params: formData,
            objId: this.activeModel.bk_obj_id,
            config: {
              requestId: 'importObjectAttribute',
              globalError: false,
              transformData: false
            }
          }).then((res) => {
            this.$http.cancel(`post_searchObjectAttribute_${this.activeModel.bk_obj_id}`)
            return res
          })
          if (res.result) {
            const data = res.data[this.activeModel.bk_obj_id]
            if (has(data, 'insert_failed')) {
              this.$error(data.insert_failed[0])
            } else if (has(data, 'update_failed')) {
              this.$error(data.update_failed[0])
            } else {
              this.$success(this.$t('导入成功'))
              this.$refs.field && this.$refs.field.initFieldList()
            }
          } else {
            this.$error(res.bk_error_msg)
          }
        } catch (e) {
          this.$error(e.data.bk_error_msg)
        } finally {
          this.$refs.fileInput.value = ''
        }
      },
      checkModel() {
        return this.models.find(model => model.bk_obj_id === this.$route.params.modelId)
      },
      hideChooseBox() {
        this.isIconListShow = false
      },
      handleModelIconUpdateConfirm(value) {
        this.isIconListShow = false
        this.saveModel({ modelIcon: value })
      },
      handleModelNameUpdateConfirm({ value, confirm, stop }) {
        this.saveModel({
          modelName: value
        })
          .then(() => {
            confirm()
          })
          .catch((err) => {
            stop()
            console.log(err)
          })
      },
      handleModelGroupUpdateConfirm({ value, confirm, stop }) {
        this.saveModel({
          classificationId: value
        })
          .then(() => {
            confirm()
          })
          .catch((err) => {
            stop()
            console.log(err)
          })
      },
      async saveModel({ modelIcon, modelName, classificationId } = {}) {
        const params = {
          modifier: this.userName,
        }

        if (modelIcon) params.bk_obj_icon = modelIcon
        if (classificationId) params.bk_classification_id = classificationId
        if (modelName) params.bk_obj_name = modelName

        return this.updateObject({
          id: this.activeModel.id,
          params
        })
          .then(() => {
            this.$http.cancel('post_searchClassificationsObjects')
            this.$success(this.$t('修改成功'))
            this.activeModel = { ...this.activeModel, ...params }
          })
      },
      initObject() {
        const model = this.$store.getters['objectModelClassify/getModelById'](this.$route.params.modelId)
        if (model) {
          this.activeModel = model
          const menuI18n = this.$route.meta.menu.i18n && this.$t(this.$route.meta.menu.i18n)
          this.$store.commit('setTitle', `${menuI18n}【${this.activeModel.bk_obj_name}】`)
          this.getModelInstanceCount()
          this.getModelBindTemplate()
        } else {
          this.$routerActions.redirect({ name: 'status404' })
        }
      },
      async getModelBindTemplate() {
        const templateList = await fieldTemplateService.getModelBindTemplate({
          object_id: this.activeModel.id
        })
        if (templateList?.info?.length) {
          this.templateList = templateList.info.map(item => ({
            id: item.id,
            name: item.name
          }))
        } else {
          this.templateList = []
        }
      },
      async getModelInstanceCount() {
        const result = await this.$store.dispatch('objectCommonInst/searchInstanceCount', {
          params: {
            condition: { obj_ids: [this.activeModel.bk_obj_id] }
          },
          config: {
            requestId: this.request.instanceCount,
            globalError: false
          }
        })

        const [data] = result
        this.modelInstanceCount = data?.error ? '--' : data?.inst_count
      },
      exportField() {
        modelImportExportService.export(this.activeModel.bk_obj_id)
      },
      dialogConfirm(type) {
        switch (type) {
          case 'restart':
            this.$bkInfo({
              title: this.$t('确认要启用该模型？'),
              confirmFn: () => {
                this.updateModelObject(false)
              }
            })
            break
          case 'stop':
            this.$bkInfo({
              title: this.$t('确认要停用该模型？'),
              confirmFn: () => {
                this.updateModelObject(true)
              }
            })
            break
          case 'delete':
            this.$bkInfo({
              title: this.$t('确认要删除该模型？'),
              confirmFn: () => {
                this.deleteModel()
              }
            })
            break
          default:
        }
      },
      async updateModelObject(ispaused) {
        await this.updateObject({
          id: this.activeModel.id,
          params: {
            bk_ispaused: ispaused
          },
          config: {
            requestId: 'updateModel'
          }
        })
        this.$store.commit('objectModelClassify/updateModel', {
          bk_ispaused: ispaused,
          bk_obj_id: this.activeModel.bk_obj_id
        })
        this.activeModel = { ...this.activeModel, ...{ bk_ispaused: ispaused } }
      },
      async deleteModel() {
        if (this.isMainLineModel) {
          await this.deleteMainlineObject({
            bkObjId: this.activeModel.bk_obj_id,
            config: {
              requestId: 'deleteModel'
            }
          })
          this.$routerActions.back()

          // 更新主线模型
          this.searchMainlineObject()
        } else {
          await this.deleteObject({
            id: this.activeModel.id,
            config: {
              requestId: 'deleteModel'
            }
          })
          this.$routerActions.redirect({ name: MENU_MODEL_MANAGEMENT })
        }
        this.$http.cancel('post_searchClassificationsObjects')
      },
      handleGoInstance() {
        const model = this.activeModel
        if (has(BUILTIN_MODEL_RESOURCE_MENUS, model.bk_obj_id)) {
          const query = model.bk_obj_id === 'host' ? { scope: 'all' } : {}
          this.$routerActions.redirect({
            name: BUILTIN_MODEL_RESOURCE_MENUS[model.bk_obj_id],
            query
          })
        } else {
          this.$routerActions.redirect({
            name: MENU_RESOURCE_INSTANCE,
            params: {
              objId: model.bk_obj_id
            }
          })
        }
      },
      handleUploadDone(res) {
        const data = res.data[this.activeModel.bk_obj_id]
        if (res.result) {
          this.uploadResult.success = data.success
          this.$success(this.$t('导入成功'))
          this.$refs.field.resetData()
          this.importField.show = false
        } else {
          this.uploadResult.errors = data.errors
          this.uploadResult.insert_failed = data.insert_failed
          this.uploadResult.update_failed = data.update_failed
        }
      },
      handleSliderHide() {
        this.uploadResult = {
          success: null,
          errors: null,
          insert_failed: null,
          update_failed: null
        }
      },
      handleImportField() {
        if (this.isReadOnly) return
        this.importField.show = true
      },
      handleUnbindTemplate(template) {
        this.$bkInfo({
          type: 'warning',
          title: this.$t('确认解绑该模板'),
          subTitle: this.$t('解绑后，字段内容与唯一校验将会与模板脱离关系，不再受模板管理'),
          okText: this.$t('解绑'),
          cancelText: this.$t('取消'),
          confirmLoading: true,
          confirmFn: async () => {
            const params = {
              bk_template_id: template.id,
              object_id: this.activeModel.id
            }
            await fieldTemplateService.unbind(params)
            this.$success(this.$t('解绑成功'))
            this.getModelBindTemplate()
            return true
          }
        })
      },
      handleViewTemplate(template) {
        this.$routerActions.open({
          name: MENU_MODEL_FIELD_TEMPLATE,
          query: {
            id: template.id,
            action: 'view'
          }
        })
      },
      handleGoSync(template) {
        this.$routerActions.redirect({
          name: MENU_MODEL_FIELD_TEMPLATE_SYNC_MODEL,
          params: {
            id: template.id,
            modelId: this.activeModel.id
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .model-info {
        &-wrapper{
          padding: 0;
        }

        display: flex;
        height: 100px;
        background: #fff;
        font-size: 14px;
        box-shadow: 0px 2px 4px 0px rgba(25,25,41,0.05);
        align-items: center;

        .choose-icon-wrapper {
            position: relative;
            margin-left: 32px;
            .model-type {
                position: absolute;
                left: 30px;
                top: -16px;
                padding: 0 8px;
                border-radius: 4px;
                background-color: #dcfde2;
                font-size: 20px;
                line-height: 32px;
                color: #34ce5c;
                font-weight: 700;
                white-space: nowrap;
                transform: scale(.5);
                transform-origin: left center;
                z-index: 2;
                &::after {
                    content: "";
                    position: absolute;
                    top: 100%;
                    left: 50%;
                    width: 0;
                    height: 0;
                    border-top: 8px solid #dcfde2;
                    border-right: 14px solid transparent;
                    transform: translateX(-50%);
                }
                &.is-builtin {
                    background-color: #ffb23a;
                    color: #fff;
                    &::after{
                      border-top-color: #ffb23a;
                    }
                }
            }
            .choose-icon-box {
                position: absolute;
                left: -12px;
                top: 62px;
                width: 600px;
                height: 460px;
                background: #fff;
                border: 1px solid #dde4e8;
                box-shadow: 0px 3px 6px 0px rgba(51, 60, 72, 0.13);
                z-index: 99;
                &:before {
                    position: absolute;
                    top: -13px;
                    left: 30px;
                    content: '';
                    border: 6px solid transparent;
                    border-bottom-color: rgba(51, 60, 72, 0.23);
                }
                &:after {
                    position: absolute;
                    top: -12px;
                    left: 30px;
                    content: '';
                    border: 6px solid transparent;
                    border-bottom-color: #fff;
                }
            }
        }
        .icon-box {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 56px;
            height: 56px;
            border-radius: 50%;
            background: #e7f0ff;
            text-align: center;
            font-size: 20px;
            color: $cmdbBorderFocusColor;
            cursor: pointer;

            &.is-builtin {
              background: #f5f7fa;
              color: #798aad;
            }

            &:hover {
                .hover-text {
                    background: rgba(0, 0, 0, .5);
                    display: block;
                }
            }
            .hover-text {
                display: none;
                position: absolute;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                line-height: 56px;
                font-size: 12px;
                border-radius: 50%;
                text-align: center;
                color: #fff;
                &.is-paused {
                    background: rgba(0, 0, 0, .5);
                    display: block !important;
                }
            }
            .icon {
                &.ispre {
                    color: #3a84ff;
                }
            }
        }

        .model-identity {
          width: 225px;
          margin-left: 10px;
          margin-right: 10px;

          .model-name {
            font-weight: 700;
            color: #313238;
            line-height: 26px;

            .bk-tag {
              font-weight: normal;
              height: 18px;
              line-height: 18px;
              padding: 0 6px;
            }
          }

          .model-id {
            font-size: 12px;
            color: #979ba5;
            @include ellipsis;
          }
        }

        .model-group-name {
          width: 250px;
          font-size: 12px;
          color: #63656e;
          display: flex;
          flex-direction: column;

          &-label {
            flex: 0 0 auto;
            line-height: 26px;
            color: #979BA5;
          }
          &-label-editing {
            margin-top: 13px;
          }
        }

        .instance-count {
            width: 250px;
            display: flex;
            flex-wrap: wrap;
            flex-direction: column;
            font-size: 12px;
            color: #63656e;
            &-label {
              line-height: 26px;
              color: #979BA5;
            }
            &-text {
              color: #3a84ff;
              cursor: pointer;
              display: flex;
              align-items: center;
            }
         }
         .field-template {
            max-width: 400px;
            font-size: 12px;
            color: #63656e;
          }
          .field-template-label {
            line-height: 26px;
            color: #979BA5;
          }
          .field-template-tag {
            line-height: 26px;
            .unbind-icon {
              font-size: 12px !important;
              margin: 0 4px;
              padding: 0;
            }

            :deep(.tag-item-text) {
              position: relative;
              @include ellipsis;
              .reddot {
                position: relative;
                right: 0px;
                top: -6px;
                width: 6px;
                height: 6px;
                background: #EA3636;
                border-radius: 50%;
                display: inline-block;
              }
            }
          }
        .restart-btn {
            display: inline-block;
        }
        .btn-group {
            margin-left: auto;
            margin-right: 24px;
            display: flex;
            align-items: center;
            .label-btn {
                display: flex;
                align-items: center;
                outline: 0;
                position: relative;
                font-size: 12px;
                margin-left: 10px;
                cursor: pointer;
                &:hover {
                    color: $cmdbBorderFocusColor;
                    .bk-button-text {
                        color: $cmdbBorderFocusColor;
                    }
                }
                &-text {
                  margin-left: 4px;
                }
                &-icon {
                  display: inline-block;
                  vertical-align: middle;
                  line-height: 1;
                  height: 14px;
                 }
                .bk-button-text {
                  font-size: inherit;
                  color: #737987;
                    &:disabled {
                        color: #dcdee5 !important;
                        cursor: not-allowed;
                    }
                }
                &.disabled {
                    cursor: not-allowed;
                    opacity: 0.5;
                    &:hover {
                      color: inherit;
                    }
                }
                input[type="file"] {
                    position: absolute;
                    left: 0;
                    top: 0;
                    opacity: 0;
                    width: 100%;
                    height: 100%;
                    cursor: pointer;
                }
                ::-webkit-file-upload-button {
                    cursor:pointer;
                }
            }
            .export-form {
                display: inline-block;
            }
        }
    }
    /deep/ .model-details-tab {
      height: calc(100% - 100px);
      .bk-tab-header {
        padding: 0 18px;
        background: #fff;
        box-shadow: 0 2px 4px 0 #1919290d;
      }
      .bk-tab-section {
        padding: 0;
      }
    }
    .editable-field {
      width: 100%;
    }
</style>

<style lang="scss">
@import '@/assets/scss/model-manage.scss';

.template-diff-sync-theme {
  .diff-sync-content {
    .content-tips {
      display: flex;
      align-items: center;
    }
    .bk-link-text {
      font-size: 12px;
    }
  }
}
</style>
