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
  <div
    class="model-management"
    v-bkloading="{ isLoading: $loading(requestIds.searchClassifications) }"
    :class="{
      'is-model-selectable': isModelSelectable,
      'is-tips-hidden': isTipsHidden
    }"
  >
    <div
      class="model-management-header clearfix"
    >
      <cmdb-tips
        tips-key="modelTips"
        @input="isVisible => isTipsHidden = !isVisible"
        more-link="https://bk.tencent.com/docs/markdown/配置平台/产品白皮书/产品功能/Model.md"
      >
        {{ $t("模型顶部提示") }}
      </cmdb-tips>

      <div class="model-management-options">
        <div class="model-export-label">{{$t('请选择需要导出的模型')}}</div>

        <div class="model-operation-options">
          <cmdb-auth :auth="{ type: $OPERATION.C_MODEL }">
            <template #default="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled || modelType === 'disabled'"
                @click="showModelDialog(UNCATEGORIZED_GROUP_ID)"
              >
                {{ $t("新建模型") }}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth :auth="{ type: $OPERATION.C_MODEL_GROUP }">
            <template #default="{ disabled }">
              <bk-button
                theme="default"
                :disabled="disabled || modelType === 'disabled'"
                @click="showGroupDialog(false)"
              >
                {{ $t("新建分组") }}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth :auth="[
            { type: $OPERATION.C_MODEL_GROUP },
            { type: $OPERATION.C_MODEL }
          ]">
            <template #default="{ disabled }">
              <bk-button
                theme="default"
                :disabled="disabled || modelType === 'disabled'"
                @click="importModel"
              >
                {{ $t("导入") }}
              </bk-button>
            </template>
          </cmdb-auth>
          <bk-button
            theme="default"
            :disabled="modelType === 'disabled'"
            @click="exportModel"
          >
            {{ $t("导出") }}
          </bk-button>
        </div>

        <div class="model-type-options">
          <bk-button
            class="model-type-button"
            :class="[{ 'is-active': modelType === '' }]"
            size="small"
            @click="modelType = ''"
          >
            {{ $t("全部") }}
          </bk-button>
          <bk-button
            class="model-type-button"
            :class="[{ 'is-active': modelType === 'enable' }]"
            size="small"
            @click="modelType = 'enable'"
          >
            {{ $t("启用中") }}
          </bk-button>
          <span
            class="inline-block-middle"
            style="outline: 0"
            v-bk-tooltips="disabledModelBtnText"
          >
            <bk-button
              class="model-type-button disabled"
              :class="[
                { 'is-active': modelType === 'disabled' }
              ]"
              size="small"
              :disabled="!disabledClassifications.length"
              @click="modelType = 'disabled'"
            >
              {{ $t("已停用") }}
            </bk-button>
          </span>
        </div>

        <div class="model-search-options">
          <bk-input
            class="model-search-input"
            :clearable="true"
            :right-icon="'bk-icon icon-search'"
            :placeholder="$t('请输入关键字')"
            v-model.trim="modelSearchKey"
          >
          </bk-input>
        </div>
      </div>
    </div>

    <div class="model-management-body">
      <ul
        class="group-list"
        :class="{
          'is-dragging': isDragging
        }"
      >
        <li
          class="group-item clearfix"
          ref="classification"
          v-for="(classification, classIndex) in currentClassifications"
          :key="classIndex"
        >
          <div class="group-header clearfix">
            <collapse-group-title
              :is-new-classify="classification.isNewClassify"
              :dropdown-menu="!isModelSelectable"
              :collapse=" classificationsCollapseState[classification.id]"
              :title="`${classification.bk_classification_name} ( ${ classification.bk_objects.length} )`"
              @click.native="toggleModelList(classification)"
              :commands="groupCommands(classification)"
              v-bk-tooltips="{
                disabled: !isBuiltinClass(classification),
                content: $t('内置模型组不支持删除和修改'),
                placement: 'right'
              }">
            </collapse-group-title>
            <bk-checkbox
              @change="handleModelGroupSelectionChange(classification, $event)"
              v-model="classificationSelectionState[classification.bk_classification_id]"
              :disabled="isGroupEmpty(classification)"
              class="full-selection-checkbox"
            >{{$t('全选')}}</bk-checkbox>
          </div>
          <bk-transition name="collapse" duration-type="ease">
            <draggable
              class="model-list clearfix"
              :class="{
                'is-empty': isGroupEmpty(classification)
              }"
              v-show="!classificationsCollapseState[classification.id]"
              tag="div"
              :sort="false"
              :animation="200"
              :disabled="isModelSelectable"
              draggable=".model-item"
              group="model-list"
              ghost-class="model-item-ghost"
              :data-group-id="classification.bk_classification_id"
              v-model="classification.bk_objects"
              @start="handleModelDragStart"
              :move="handleModelDragMove"
              @end="handleModelDragEnd"
              @add="handleModelDragAdd"
            >
              <div
                class="model-item bgc-white"
                v-for="model in classification['bk_objects']"
                :key="model.bk_obj_id"
                :data-group-id="model.bk_classification_id"
                :data-model-id="model.id"
                :class="{
                  'is-paused': model['bk_ispaused'],
                  'is-builtin': model.ispre
                }"
                @mouseenter="handleModelMouseEnterDebounce(model)"
              >
                <div
                  class="model-info"
                  :class="{
                    'no-instance-count': model.bk_ispaused || isNoInstanceModel(model.bk_obj_id)
                  }"
                  @click="handleModelClick(model, classification)"
                >
                  <div class="drag-icon"></div>
                  <div class="model-icon">
                    <i class="icon" :class="[model['bk_obj_icon']]"></i>
                  </div>
                  <div class="model-details">
                    <p class="model-name" :title="model['bk_obj_name']">
                      {{ model["bk_obj_name"] }}
                    </p>
                    <p class="model-id" :title="model['bk_obj_id']">
                      {{ model["bk_obj_id"] }}
                    </p>
                  </div>
                  <bk-checkbox
                    v-model="modelSelectionState[model.bk_obj_id]"
                    @change="handleModelSelectionChange(classification)"
                    :disabled="isBuiltinModel(model)"
                    v-bk-tooltips="{
                      content: $t('内置模型不允许导出'),
                      disabled: !isBuiltinModel(model)
                    }"
                    @click.stop.native
                    class="model-checkbox"
                  >
                  </bk-checkbox>
                </div>
                <div
                  v-if="!model.bk_ispaused && !isNoInstanceModel(model.bk_obj_id)"
                  class="model-instance-count"
                  @click="handleGoInstance(model)"
                >
                  <span class="count-number">
                    <cmdb-loading
                      :loading="!modelStatisticsSet[model.bk_obj_id] ||
                        $loading(requestIds.statistics[model.bk_obj_id])"
                    >
                      {{ modelStatisticsSet[model.bk_obj_id] | instanceCount }}
                    </cmdb-loading>
                  </span>
                </div>
              </div>
            </draggable>
          </bk-transition>
          <bk-transition name="collapse" class="group-empty-model"
            v-if="classification.bk_objects.length === 0"
            v-show="!classificationsCollapseState[classification.id]">
            <div>
              <i class="bk-icon icon-info-circle"></i>
              <i18n path="分组暂无模型提示">
                <template #btn>
                  <bk-button :text="true" title="primary" @click="showModelDialog(classification.bk_classification_id)">
                    {{$t('立即添加')}}
                  </bk-button>
                </template>
              </i18n>
            </div>
          </bk-transition>
        </li>
      </ul>
      <no-search-results
        v-if="!currentClassifications.length"
        :text="$t('搜不到相关模型')"
      />
    </div>

    <div class="model-management-footer" v-show="isModelSelectable">
      <div class="export-action-bar">
        <bk-checkbox
          class="full-selection"
          @change="toggleAllModelSelection"
          v-model="isAllSelected"
        >{{$t('全选')}}</bk-checkbox
        >
        <span class="selected-count"
        >{{$t('已选')}}：<em>{{ exportModelsLen }}</em></span
        >
        <bk-button @click="doubleCheckCancelExport" class="cancel-button">{{$t('取消')}}</bk-button>
        <bk-button
          @click="toNextExportStep"
          theme="primary"
          :disabled="exportModelsLen === 0" class="next-step-button">{{$t('下一步')}}</bk-button>
      </div>
    </div>

    <bk-dialog
      class="bk-dialog-no-padding bk-dialog-no-tools group-dialog dialog"
      :close-icon="false"
      :width="600"
      :mask-close="false"
      v-model="groupDialog.isShow"
    >
      <div class="dialog-content">
        <p class="title">{{ groupDialog.title }}</p>
        <div class="content">
          <label>
            <div class="label-title">
              {{ $t("唯一标识") }}<span class="color-danger">*</span>
            </div>
            <div
              class="cmdb-form-item"
              :class="{ 'is-error': errors.has('classifyId') }"
            >
              <bk-input
                type="text"
                class="cmdb-form-input"
                name="classifyId"
                :placeholder="$t('请填写英文开头，下划线，数字，英文的组合')"
                :disabled="groupDialog.isEdit"
                v-model.trim="groupDialog.data['bk_classification_id']"
                v-validate="'required|classifyId|length:128|reservedWord'">
              </bk-input>
              <p class="form-error" :title="errors.first('classifyId')">
                {{ errors.first("classifyId") }}
              </p>
            </div>
            <i
              class="bk-icon icon-info-circle"
              v-bk-tooltips="$t('请填写英文开头，下划线，数字，英文的组合')"
            ></i>
          </label>
          <label>
            <span class="label-title">
              {{ $t("名称") }}
            </span>
            <span class="color-danger">*</span>
            <div
              class="cmdb-form-item"
              :class="{ 'is-error': errors.has('classifyName') }"
            >
              <bk-input
                type="text"
                class="cmdb-form-input"
                name="classifyName"
                :placeholder="$t('请输入名称')"
                v-model.trim="groupDialog.data['bk_classification_name']"
                v-validate="'required|length:128'"
              >
              </bk-input>
              <p class="form-error" :title="errors.first('classifyName')">
                {{ errors.first("classifyName") }}
              </p>
            </div>
          </label>
        </div>
      </div>
      <div slot="footer" class="footer">
        <bk-button
          theme="primary"
          :loading="$loading(['updateClassification', 'createClassification'])"
          @click="saveGroup"
        >
          {{ groupDialog.isEdit ? $t("保存") : $t("提交") }}
        </bk-button>
        <bk-button theme="default" @click="hideGroupDialog">{{
          $t("取消")
        }}</bk-button>
      </div>
    </bk-dialog>

    <the-create-model
      :is-show.sync="modelDialog.isShow"
      :group-id.sync="modelDialog.groupId"
      :title="$t('新建模型')"
      :operating="$loading('createModel')"
      @confirm="saveModel">
    </the-create-model>

    <bk-dialog
      class="bk-dialog-no-padding"
      :width="400"
      :show-footer="false"
      :mask-close="true"
      v-model="modelCreatedDialogVisible"
    >
      <div class="success-content">
        <i class="bk-icon icon-check-1"></i>
        <p>{{ $t("模型创建成功") }}</p>
        <div class="btn-box">
          <bk-button
            theme="primary"
            @click="checkModelDetails(curCreateModel)"
          >{{ $t("配置字段") }}</bk-button
          >
          <bk-button @click="modelCreatedDialogVisible = false">{{
            $t("返回列表")
          }}</bk-button>
        </div>
      </div>
    </bk-dialog>

    <!-- 模型导出 -->
    <model-export-pane
      @cancel="doubleCheckCancelExport"
      @done="doneExport"
      :data="exportModels"
      v-if="exportPaneVisible">
    </model-export-pane>

    <!-- 模型导入 -->
    <model-import-pane
      v-if="importPaneVisible"
      @cancel="doubleCheckCancelImport"
      @done="doneImport">
    </model-import-pane>
  </div>
</template>

<script>
  import has from 'has'
  import theCreateModel from '@/components/model-manage/_create-model'
  import cmdbLoading from '@/components/loading/index.vue'
  import noSearchResults from '@/views/status/no-search-results.vue'
  import CollapseGroupTitle from './children/collapse-group-title.vue'
  import { mapGetters, mapMutations, mapActions } from 'vuex'
  import debounce from 'lodash.debounce'
  import Draggable from 'vuedraggable'
  import ModelExportPane from './children/model-import-export/model-export-pane/index.vue'
  import ModelImportPane from './children/model-import-export//model-import-pane/index.vue'
  import {
    MENU_RESOURCE_INSTANCE,
    MENU_MODEL_DETAILS,
  } from '@/dictionary/menu-symbol'
  import { BUILTIN_MODEL_RESOURCE_MENUS, UNCATEGORIZED_GROUP_ID } from '@/dictionary/model-constants.js'
  import Bus from '@/utils/bus'

  export default {
    name: 'ModelManagement',
    filters: {
      instanceCount(value) {
        if (!value) return
        if (value?.error) {
          return '--'
        }
        return value.inst_count > 999 ? '999+' : value.inst_count
      },
    },
    components: {
      theCreateModel,
      noSearchResults,
      cmdbLoading,
      CollapseGroupTitle,
      Draggable,
      ModelExportPane,
      ModelImportPane
    },
    data() {
      return {
        UNCATEGORIZED_GROUP_ID,
        modelType: '', // 模型启用、停用状态
        modelSearchKey: '', // 模型搜索关键字
        filterClassifications: [], // 过滤后的分组
        modelStatisticsSet: {}, // 模型实例数量统计
        curCreateModel: {}, // 当前创建的模型
        modelCreatedDialogVisible: false,

        // 分组表单弹窗
        groupDialog: {
          isShow: false,
          isEdit: false,
          title: this.$t('新建分组'),
          data: {
            bk_classification_id: '',
            bk_classification_name: '',
            id: '',
          },
        },

        // 新建模型弹窗
        modelDialog: {
          isShow: false,
          groupId: '',
        },
        requestIds: {
          statistics: [], // 模型实例数据加载请求 id
          searchClassifications: Symbol('searchClassifications') // 模型分组数据加载请求 id
        },
        isDragging: false, // 是否处于拖拽状态
        classificationsCollapseState: {}, // 分类折叠状态

        /**
         * 导出模型
         */
        exportModels: [], // 已选择的需要的导出模型
        isModelSelectable: false, // 是否在导出选择中
        exportPaneVisible: false, // 导出面板是否显示
        isAllSelected: false, // 是否全选所有模型
        classificationSelectionState: {}, // 导出模型组的选中状态
        modelSelectionState: {}, // 导出模型的选中状态

        /**
         * 导入模型
         */
        importPaneVisible: false, // 导入面板是否显示

        isTipsHidden: false // 模型管理提示是否隐藏
      }
    },
    computed: {
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('objectModelClassify', ['classifications']),
      allClassifications() {
        const allClassifications = []
        this.classifications
          .filter(classification => !classification?.bk_ishidden)
          .forEach((classification) => {
            allClassifications.push({
              ...classification,
              bk_objects: classification.bk_objects
                .filter(model => !model.bk_ishidden)
                .sort((a, b) => a.bk_ispaused - b.bk_ispaused),
            })
          })
        return allClassifications
      },
      enableClassifications() {
        const enableClassifications = []
        this.allClassifications.forEach((classification) => {
          enableClassifications.push({
            ...classification,
            bk_objects: classification.bk_objects.filter(model => !model.bk_ispaused),
          })
        })
        return enableClassifications.filter(item => item.bk_objects.length)
      },
      disabledClassifications() {
        const disabledClassifications = []

        this.allClassifications.forEach((classification) => {
          disabledClassifications.push({
            ...classification,
            bk_objects: classification.bk_objects.filter(model => model.bk_ispaused),
          })
        })

        return disabledClassifications.filter(item => item.bk_objects.length)
      },
      currentClassifications() {
        let currentClassifications = []

        if (!this.modelSearchKey && !this.modelType) {
          currentClassifications = this.allClassifications
        }

        if (this.modelType) {
          currentClassifications = this.modelType === 'enable' ? this.enableClassifications : this.disabledClassifications
        }

        if (this.modelSearchKey) {
          currentClassifications = this.filterClassifications
        }

        return currentClassifications.sort((a, b) => (b.bk_classification_id === UNCATEGORIZED_GROUP_ID ? -1 : 0))
      },
      disabledModelBtnText() {
        return this.disabledClassifications.length ? '' : this.$t('停用模型提示')
      },
      exportModelsLen() {
        return Object.values(this.modelSelectionState).filter(isSelected => isSelected).length
      },
      needLeaveConfirm() {
        return this.exportPaneVisible || this.isModelSelectable || this.importPaneVisible
      }
    },
    watch: {
      currentClassifications: {
        deep: true,
        handler(val) {
          const classificationsCollapseState = {}

          val.forEach(({ id }) => {
            classificationsCollapseState[id] = this.classificationsCollapseState?.[id] || false
          })

          this.classificationsCollapseState = classificationsCollapseState
        },
      },
      modelSearchKey(value) {
        if (!value) {
          return
        }
        const searchResult = []
        let currentClassifications = null

        if (!this.modelType) {
          currentClassifications = this.allClassifications
        } else if (this.modelType === 'enable') {
          currentClassifications = this.enableClassifications
        } else {
          currentClassifications = this.disabledClassifications
        }
        const classifications = this.$tools.clone(currentClassifications)
        const lowerCaseValue = value.toLowerCase()

        for (let i = 0; i < classifications.length; i++) {
          classifications[i].bk_objects = classifications[i].bk_objects.filter((model) => {
            const modelName = model.bk_obj_name.toLowerCase()
            const modelId = model.bk_obj_id.toLowerCase()
            return (
              modelName?.includes(lowerCaseValue)
              || modelId?.includes(lowerCaseValue)
            )
          })
          searchResult.push(classifications[i])
        }

        this.filterClassifications = searchResult.filter(item => item.bk_objects.length)
      },
      modelType() {
        this.modelSearchKey = ''
      },
    },
    async created() {
      window.onbeforeunload = () => {
        if (this.needLeaveConfirm) {
          return this.$t('系统不会保存您所做的修改，确认要离开？')
        }
      }

      this.handleModelMouseEnterDebounce = debounce(this.handleModelItemMouseEnter, 200)

      try {
        await this.loadAllModels()
      } catch (e) {
        this.$route.meta.view = 'error'
      }

      if (this.$route.query.modelSearchKey) {
        const { hash } = window.location
        this.modelSearchKey = this.$route.query.modelSearchKey
        window.location.hash = hash.substring(0, hash.indexOf('?'))
      }

      this.$on('model-selection-change', (model, classification) => {
        this.handleModelSelectionChange(classification)
      })
    },
    beforeDestroy() {
      this.$http.cancelRequest(this.requestIds.statistics)
    },
    beforeRouteLeave(to, from, next) {
      if (this.needLeaveConfirm) {
        this.$bkInfo({
          title: this.$t('确定离开页面？'),
          subTitle: this.$t('系统不会保存您所做的修改，确认要离开？'),
          confirmFn: () => {
            next()
          },
          cancelFn: () => {
            next(false)
          }
        })
      } else {
        next()
      }
    },
    methods: {
      ...mapMutations('objectModelClassify', [
        'updateClassify',
        'deleteClassify',
      ]),
      ...mapActions('objectModelClassify', [
        'searchClassificationsObjects',
        'getClassificationsObjectStatistics',
        'createClassification',
        'updateClassification',
        'deleteClassification',
      ]),
      ...mapActions('objectModel', ['createObject', 'updateObject']),
      loadAllModels() {
        return this.searchClassificationsObjects({
          params: {},
          config: {
            requestId: this.requestIds.searchClassifications
          }
        })
      },
      groupCommands(classification) {
        return [
          {
            text: this.$t('新建模型'),
            auth: {
              type: this.$OPERATION.C_MODEL,
              relation: [classification.id]
            },
            handler: () => this.showModelDialog(classification.bk_classification_id)
          },
          {
            text: this.$t('编辑分组'),
            visible: !this.isBuiltinClass(classification),
            auth: {
              type: this.$OPERATION.U_MODEL_GROUP,
              relation: [classification.id]
            },
            handler: () => this.showGroupDialog(true, classification)
          },
          {
            text: this.$t('删除分组'),
            visible: !this.isBuiltinClass(classification),
            disabled: Boolean(classification.bk_objects.length),
            disabledTooltips: this.$t('分组下有模型，不能删除'),
            auth: {
              type: this.$OPERATION.D_MODEL_GROUP,
              relation: [classification.id]
            },
            handler: () => this.deleteGroup(classification)
          }
        ]
      },
      toggleModelList(classification) {
        this.classificationsCollapseState[classification.id] = !this.classificationsCollapseState[classification.id]
      },
      handleModelDragStart() {
        this.isDragging = true
      },
      handleModelDragEnd() {
        this.isDragging = false
      },
      handleModelDragMove(event) {
        const draggedModel = event.draggedContext.element
        const targetGroupModels = event.relatedContext.list
        const { willInsertAfter } = event
        const isSameGroup = targetGroupModels
          .some(model => model.bk_classification_id === draggedModel?.bk_classification_id)
        if (isSameGroup && willInsertAfter) {
          return true
        }
        return !isSameGroup
      },
      handleModelDragAdd(event) {
        const { modelId } = event.item.dataset
        const newGroupId = event.to.dataset.groupId

        this.updateModelGroup({ modelId, newGroupId })
      },
      getExportModels() {
        return this.currentClassifications
          .map(classification => ({
            ...classification,
            bk_objects: classification.bk_objects
              .filter(model => this.modelSelectionState[model.bk_obj_id])
          }))
          .filter(classification => classification.bk_objects.length !== 0)
      },
      handleModelGroupSelectionChange(classification, value) {
        this.$set(this.classificationSelectionState, classification.bk_classification_id, value)
        classification.bk_objects.forEach((model) => {
          const isChecked = value && !this.isBuiltinModel(model)
          this.$set(this.modelSelectionState, model.bk_obj_id, isChecked)
        })

        this.generateAllModelSelection()
      },
      /**
       * 检查模型分组下的模型是否为空
       * @param {Object} classification 模型分组数据
       */
      isGroupEmpty(classification) {
        return classification.bk_objects.length === 0
      },
      handleModelSelectionChange(classification) {
        const isGroupChecked = classification.bk_objects.every(model => this.modelSelectionState[model.bk_obj_id])

        this.$set(this.classificationSelectionState, classification.bk_classification_id, isGroupChecked)

        this.generateAllModelSelection()
      },
      // 根据当前模型的选择状态生成全选所有 checkbox 的选择状态。
      generateAllModelSelection() {
        const isAllModelChecked = this.currentClassifications
          .every(classification => this.classificationSelectionState[classification.bk_classification_id])

        this.isAllSelected = isAllModelChecked
      },
      toggleAllModelSelection(value) {
        this.currentClassifications.forEach((classification) => {
          const isGroupChecked = !this.isGroupEmpty(classification) && value
          this.$set(this.classificationSelectionState, classification.bk_classification_id, isGroupChecked)

          classification.bk_objects.forEach((model) => {
            const isModelChecked = !this.isBuiltinModel(model) && value

            this.$set(this.modelSelectionState, model.bk_obj_id, isModelChecked)
          })
        })
      },
      exportModel() {
        this.isModelSelectable = true
      },
      doubleCheckCancelExport() {
        if (this.exportModelsLen > 0) {
          this.$bkInfo({
            title: this.$t('确定离开页面？'),
            subTitle: this.$t('系统不会保存您所做的修改，确认要离开？'),
            confirmFn: this.cancelExport
          })
        } else {
          this.cancelExport()
        }
      },
      cancelExport() {
        this.isModelSelectable = false
        this.exportPaneVisible = false
        this.modelSelectionState = {}
        this.classificationSelectionState = {}
        this.isAllSelected = false
        this.exportModels = this.getExportModels()
        Bus.$emit('disable-customize-breadcrumbs')
      },
      doneExport() {
        this.cancelExport()
        this.loadAllModels()
      },
      toNextExportStep() {
        if (this.exportModelsLen > 0) {
          this.exportModels = this.getExportModels()
          this.exportPaneVisible = true
          Bus.$emit('enable-customize-breadcrumbs', {
            title: this.$t('导出模型'),
            backward: () => {
              this.doubleCheckCancelExport()
            }
          })
        }
      },
      importModel() {
        this.importPaneVisible = true
        Bus.$emit('enable-customize-breadcrumbs', {
          title: this.$t('导入模型'),
          backward: () => {
            this.doubleCheckCancelImport()
          }
        })
      },
      cancelImport() {
        this.importPaneVisible = false
        Bus.$emit('disable-customize-breadcrumbs')
      },
      doneImport() {
        this.importPaneVisible = false
        this.loadAllModels()
      },
      doubleCheckCancelImport() {
        this.$bkInfo({
          title: this.$t('确定离开页面？'),
          subTitle: this.$t('系统不会保存您所做的修改，确认要离开？'),
          confirmFn: this.cancelImport
        })
      },
      isBuiltinClass(classification) {
        return classification.bk_classification_type === 'inner'
      },
      isBuiltinModel(model) {
        return model.ispre
      },
      showGroupDialog(isEdit, group) {
        if (isEdit) {
          this.groupDialog.data.id = group.id
          this.groupDialog.title = this.$t('编辑分组')
          this.groupDialog.data.bk_classification_id = group.bk_classification_id
          this.groupDialog.data.bk_classification_name = group.bk_classification_name
          this.groupDialog.data.id = group.id
        } else {
          this.groupDialog.title = this.$t('新建分组')
          this.groupDialog.data.bk_classification_id = ''
          this.groupDialog.data.bk_classification_name = ''
          this.groupDialog.data.id = ''
        }
        this.groupDialog.isEdit = isEdit
        this.groupDialog.isShow = true
      },
      hideGroupDialog() {
        this.groupDialog.isShow = false
        this.$validator.reset()
      },
      // 处理模型点击事件。在导出时改变模型选择状态；在非导出时跳转到模型详情。
      handleModelClick(model, classification) {
        if (this.isModelSelectable) {
          const isChecked = !this.modelSelectionState[model.bk_obj_id] && !this.isBuiltinModel(model)
          this.$set(this.modelSelectionState, model.bk_obj_id, isChecked)
          this.$emit('model-selection-change', model, classification)
        } else {
          this.checkModelDetails(model)
        }
      },
      checkModelDetails(model) {
        this.$store.commit('objectModel/setActiveModel', model)
        this.$routerActions.redirect({
          name: MENU_MODEL_DETAILS,
          params: {
            modelId: model.bk_obj_id,
          },
          history: true,
        })
      },
      handleGoInstance(model) {
        this.modelCreatedDialogVisible = false
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
              objId: model.bk_obj_id,
            },
          })
        }
      },
      isNoInstanceModel(modelId) {
        // 不能直接查看实例的模型
        const noInstanceModelIds = ['set', 'module']
        return noInstanceModelIds.includes(modelId)
      },
      handleModelItemMouseEnter(model) {
        if (!model) return

        const isDisabledModel = model.bk_ispaused && !model.bk_ishidden

        if (isDisabledModel || this.isNoInstanceModel(model.bk_obj_id)) return

        this.getModelInstanceCount(model.bk_obj_id)
      },
      async getModelInstanceCount(modelId) {
        // 存在则不再请求
        if (has(this.modelStatisticsSet, modelId)) {
          return
        }

        const requestId = `getModelInstanceCount_${modelId}`
        this.requestIds.statistics.push(requestId)

        // 取消上一个请求
        const currentIndex = this.requestIds.statistics.findIndex(rid => rid === requestId)
        const prevIndex = this.requestIds.statistics[currentIndex - 1]
        if (prevIndex) {
          this.$http.cancelRequest(prevIndex)
        }

        const result = await this.$store.dispatch(
          'objectCommonInst/searchInstanceCount',
          {
            params: {
              condition: { obj_ids: [modelId] },
            },
            config: {
              requestId,
              globalError: false,
            },
          }
        )

        const [data] = result

        this.$set(this.modelStatisticsSet, data.bk_obj_id, {
          error: data.error,
          inst_count: data.inst_count,
        })
      },
      async saveGroup() {
        try {
          const res = await Promise.all([
            this.$validator.validate('classifyId'),
            this.$validator.validate('classifyName')
          ])
          if (res.includes(false)) {
            return
          }
          const params = {
            bk_supplier_account: this.supplierAccount,
            bk_classification_id: this.groupDialog.data.bk_classification_id,
            bk_classification_name: this.groupDialog.data.bk_classification_name
          }
          if (this.groupDialog.isEdit) {
            // eslint-disable-next-line no-unused-vars
            const res = await this.updateClassification({
              id: this.groupDialog.data.id,
              params,
              config: {
                requestId: 'updateClassification'
              }
            })
            this.updateClassify({ ...params, ...{ id: this.groupDialog.data.id, isNewClassify: false } })
          } else {
            const res = await this.createClassification({
              params,
              config: { requestId: 'createClassification' }
            })
            this.updateClassify({ ...params, ...{ id: res.id, isNewClassify: true } })
            this.$success(this.$t('新建成功'))
          }
          this.hideGroupDialog()
          this.modelSearchKey = ''
          const classificationDomList = this.$refs.classification.at(-2)
          classificationDomList.scrollIntoView()
        } catch (error) {
          console.log(error)
        }
      },
      deleteGroup(group) {
        this.$bkInfo({
          title: this.$t('确认要删除此分组'),
          confirmFn: async () => {
            try {
              await this.deleteClassification({
                id: group.id
              })
              this.$store.commit('objectModelClassify/deleteClassify', group.bk_classification_id)
              this.modelSearchKey = ''
              this.$success(this.$t('删除成功'))
            } catch (error) {
              console.log(error)
            }
          }
        })
      },
      showModelDialog(groupId) {
        this.modelDialog.groupId = groupId || ''
        this.modelDialog.isShow = true
      },
      updateModelGroup({ modelId, newGroupId }) {
        return this.updateObject({
          id: modelId,
          params: {
            modifier: this.userName,
            bk_classification_id: newGroupId
          }
        })
          .then(() => {
            this.$success(this.$t('修改成功'))
            this.loadAllModels()
          })
      },
      async saveModel(data) {
        const params = {
          bk_supplier_account: this.supplierAccount,
          bk_obj_name: data.bk_obj_name,
          bk_obj_icon: data.bk_obj_icon,
          bk_classification_id: data.bk_classification_id,
          bk_obj_id: data.bk_obj_id,
          userName: this.userName,
        }
        try {
          const createModel = await this.createObject({ params, config: { requestId: 'createModel' } })
          this.curCreateModel = createModel
          this.modelCreatedDialogVisible = true
          this.$http.cancel('post_searchClassificationsObjects')
          this.getModelInstanceCount(params.bk_obj_id)
          this.loadAllModels()
          this.modelDialog.isShow = false
          this.modelDialog.groupId = ''
          this.modelSearchKey = ''
        } catch (error) {
          console.log(error)
        }
      }
    },
  }
</script>

<style lang="scss" scoped src="./style.scss"></style>
