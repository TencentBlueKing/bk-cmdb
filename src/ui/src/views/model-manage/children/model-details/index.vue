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
                :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
              >
                <template #append>
                  <bk-tag v-if="activeModel.bk_ispaused" size="small" theme="default">{{$t('已停用')}}</bk-tag>
                </template>
              </editable-field>
            </div>
            <div class="model-id" v-show="!modelNameIsEditing">
              {{activeModel['bk_obj_id'] || ''}}
            </div>
          </div>
          <div class="model-group-name">
            <span class="model-group-name-label">{{$t('所属分组')}}：</span>
            <editable-field
              v-model="activeModel.bk_classification_id"
              :label="modelClassificationName"
              :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
              validate="required"
              @confirm="handleModelGroupUpdateConfirm"
              type="enum"
              font-size="12px"
              :options="classifications
                .map(item => ({ id: item.bk_classification_id, name: item.bk_classification_name }))"
            >
            </editable-field>
          </div>
          <div class="instance-count"
            v-if="!activeModel['bk_ispaused'] && !isNoInstanceModel">
            <span>{{$t('实例数量')}}：</span>
            <span class="instance-count-text" @click="handleGoInstance">
              <cmdb-loading :loading="$loading(request.instanceCount)">
                {{modelInstanceCount || 0}}
              </cmdb-loading>
              <i class="icon-cc-share instance-count-link-icon"></i>
            </span>
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
            <template v-if="canBeImport">
              <cmdb-auth tag="label" class="label-btn"
                v-if="tab.active === 'field'"
                :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
                :class="{ 'disabled': isReadOnly }"
                @click="handleImportField">
                <i class="icon-cc-import"></i>
                <span class="label-btn-text">{{$t('导入')}}</span>
              </cmdb-auth>
              <label class="label-btn" @click="exportField">
                <i class="icon-cc-derivation"></i>
                <span class="label-btn-text">{{$t('导出')}}</span>
              </label>
            </template>
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
        <the-field-group ref="field" v-if="tab.active === 'field'"></the-field-group>
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
  import modelImportExportService from '@/service/model/import-export'
  import {
    MENU_MODEL_MANAGEMENT,
    MENU_RESOURCE_INSTANCE
  } from '@/dictionary/menu-symbol'
  import { BUILTIN_MODEL_RESOURCE_MENUS } from '@/dictionary/model-constants.js'
  import EditableField from './editable-field.vue'

  export default {
    name: 'ModelDetails',
    components: {
      theFieldGroup,
      theRelation,
      theVerification,
      theChooseIcon,
      cmdbImport,
      cmdbLoading,
      EditableField
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
          insert_failed: null,
          update_failed: null
        },
        request: {
          instanceCount: Symbol('instanceCount')
        },
        modelNameIsEditing: false
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
    },
    watch: {
      '$route.params.modelId'() {
        this.initObject()
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
          this.getModelInstanceCount()
        } else {
          this.$routerActions.redirect({ name: 'status404' })
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
          this.uploadResult.insert_failed = data.insert_failed
          this.uploadResult.update_failed = data.update_failed
        }
      },
      handleSliderHide() {
        this.uploadResult = {
          success: null,
          insert_failed: null,
          update_failed: null
        }
      },
      handleImportField() {
        if (this.isReadOnly) return
        this.importField.show = true
      }
    }
  }
</script>

<style lang="scss" scoped>
    .model-info {
        &-wrapper{
          padding: 20px 24px;
        }

        display: flex;
        height: 80px;
        background: #fff;
        font-size: 14px;
        box-shadow: 0px 2px 4px 0px rgba(25,25,41,0.05);
        align-items: center;

        .choose-icon-wrapper {
            position: relative;
            margin-left: 32px;
            .model-type {
                $builtinColor:#ffb23a;
                $customizeColor: #dcfde2;
                position: absolute;
                left: 30px;
                top: -16px;
                padding: 0 8px;
                border-radius: 4px;
                background-color: $customizeColor;
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
                    border-top: 8px solid $customizeColor;
                    border-right: 14px solid transparent;
                    transform: translateX(-50%);
                }
                &.is-builtin {
                    background-color: $builtinColor;
                    color: #fff;
                    &::after{
                      border-top-color: $builtinColor;
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
            $iconSize:56px;
            width: $iconSize;
            height: $iconSize;
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
                line-height: $iconSize;
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

          .model-name {
            font-weight: 700;

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
          }
        }

        .model-group-name {
          width: 250px;
          font-size: 12px;
          color: #63656e;
          display: flex;
          align-items: center;

          &-label {
            flex: 0 0 auto;
          }
        }

        .instance-count {
            display: flex;
            font-size: 12px;
            color: #63656e;
            &-text {
              color: #3a84ff;
              cursor: pointer;
              display: flex;
              align-items: center;
            }
            &-link-icon {
              margin-left: 6px;
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
        height: calc(100% - 120px);
        margin: 0 20px;
        background-color: #fff;
        border-radius: 2px;
        box-shadow: 0px 2px 4px 0px rgba(25,25,41,0.05);
        .bk-tab-header {
            padding: 0;
            margin: 0 10px;
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
</style>
