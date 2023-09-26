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
  import { ref, watch, watchEffect, onBeforeUnmount, computed, set, del } from 'vue'
  import { useHttp } from '@/api'
  import { t } from '@/i18n'
  import Loading from '@/components/loading/index.vue'
  import MiniTag from '@/components/ui/other/mini-tag.vue'
  import SelectModelDialog from './select-model-dialog.vue'
  import FieldDiff from './field-diff.vue'
  import UniqueDiff from './unique-diff.vue'
  import CombineRequest from '@/api/combine-request.js'
  import fieldTemplateService from '@/service/field-template'

  const props = defineProps({
    modelList: {
      type: Array,
      default: () => ([])
    },
    // 已绑定的模型列表，绑定模型入口时使用此属性与编辑页中的modelList区分开
    bindedModelList: {
      type: Array,
      default: () => ([])
    },
    templateId: {
      type: Number,
    },
    // 模板字段列表
    fieldList: {
      type: Array,
      default: () => ([])
    },
    // 删除的模板字段列表
    removedFieldList: {
      type: Array,
      default: () => ([])
    },
    // 模板唯一校验列表
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    // 删除的模板唯一校验
    removedUniqueList: {
      type: Array,
      default: () => ([])
    },
    readonly: {
      type: Boolean,
      default: false
    },
    height: String
  })

  const emit = defineEmits(['update-diffs', 'update-model-auth'])

  const http = useHttp()

  const selectModelDialogRef = ref(null)
  const modelListLocal = ref([])
  const selectedModel = ref(null)
  const currentTab = ref('field')

  const fieldDiffs = ref({})
  const uniqueDiffs = ref({})
  const diffLoadingIds = ref({})
  const unmountCallbacks = []
  const hasDiffError = ref(false)

  const modelEditAuths = ref({})

  const fetchFieldDiff = async (modelList) => {
    const modelIds = modelList.filter(item => !item.bk_ispaused).map(item => item.id)
    const allResult = await CombineRequest.setup(Symbol(), (params) => {
      const [modelId] = params
      const requestId = Symbol(modelId)
      if (!diffLoadingIds.value[modelId]) {
        diffLoadingIds.value[modelId] = [requestId]
      } else {
        diffLoadingIds.value[modelId].push(requestId)
      }

      return fieldTemplateService.getFieldDifference({
        bk_template_id: props.templateId,
        object_id: modelId,
        attributes: props.fieldList
      }, {
        requestId,
        globalError: true
      })
    }, { segment: 1, concurrency: 5 }).add(modelIds)

    let groupIndex = 0
    for (const result of allResult) {
      // 一个分组的执行结果
      const results = await result
      for (let i = 0; i < results.length; i++) {
        // 分组中的每一个执行结果
        const { status, reason, value } = results[i]
        if (status === 'rejected') {
          console.error(reason?.message)
          hasDiffError.value = true
          continue
        }
        set(fieldDiffs.value, modelIds[(groupIndex * 5) + i], value ?? {})
      }
      groupIndex += 1
    }

    unmountCallbacks.push(() => allResult?.return())
  }

  const fetchUniqueDiff = async (modelList) => {
    const modelIds = modelList.filter(item => !item.bk_ispaused).map(item => item.id)
    const allResult = await CombineRequest.setup(Symbol(), (params) => {
      const [modelId] = params

      const requestId = Symbol(modelId)
      if (!diffLoadingIds.value[modelId]) {
        diffLoadingIds.value[modelId] = [requestId]
      } else {
        diffLoadingIds.value[modelId].push(requestId)
      }

      return fieldTemplateService.getUniqueDifference({
        bk_template_id: props.templateId,
        object_id: modelId,
        uniques: props.uniqueList
      }, {
        requestId,
        globalError: true
      })
    }, { segment: 1, concurrency: 5 }).add(modelIds)

    let groupIndex = 0
    for (const result of allResult) {
      // 一个分组的执行结果
      const results = await result
      for (let i = 0; i < results.length; i++) {
        // 分组中的每一个执行结果
        const { status, reason, value } = results[i]
        if (status === 'rejected') {
          console.error(reason?.message)
          hasDiffError.value = true
          continue
        }
        set(uniqueDiffs.value, modelIds[(groupIndex * 5) + i], value ?? {})
      }
      groupIndex += 1
    }

    unmountCallbacks.push(() => allResult?.return())
  }

  watchEffect(async () => {
    const initModelList = props.modelList.slice()
    if (initModelList.length) {
      fetchFieldDiff(initModelList)
      fetchUniqueDiff(initModelList)
    }
    modelListLocal.value = initModelList
  })

  // 查找匹配第1个可用的模型作为选中模型
  watch([modelListLocal, modelEditAuths], ([modelList, modelAuths]) => {
    // 当前存在选中的模型则不变更
    if (selectedModel.value) {
      return
    }

    let firstModel = null
    for (let i = 0; i < modelList.length; i++) {
      const model = modelList[i]
      // 未停用且有权限
      if (!model.bk_ispaused && modelAuths[model.id]) {
        firstModel = model
        break
      }
    }
    selectedModel.value = firstModel
  }, { deep: true })

  const fieldDiffCounts = computed(() => {
    const counts = {}
    Object.keys(fieldDiffs.value).forEach((modelId) => {
      const diffs = fieldDiffs.value[modelId]
      counts[modelId] = {
        total: (diffs?.create?.length ?? 0) + (diffs?.update?.length ?? 0) + (diffs?.conflict?.length ?? 0),
        conflict: diffs?.conflict?.length ?? 0
      }
    })
    return counts
  })
  const uniqueDiffCounts = computed(() => {
    const counts = {}
    Object.keys(uniqueDiffs.value).forEach((modelId) => {
      const diffs = uniqueDiffs.value[modelId]
      counts[modelId] = {
        total: (diffs?.create?.length ?? 0) + (diffs?.update?.length ?? 0) + (diffs?.conflict?.length ?? 0),
        conflict: diffs?.conflict?.length ?? 0
      }
    })
    return counts
  })

  watch([fieldDiffCounts, uniqueDiffCounts], (diffCounts) => {
    let hasDiffConflict = false
    diffCounts.forEach((diffItem) => {
      // 每一项分别是field和unique的diffCounts，只要有任何一个存在冲突则判定为冲突
      for (const [, value] of Object.entries(diffItem)) {
        if (value.conflict) {
          hasDiffConflict = true
          break
        }
      }
    })
    emit('update-diffs', hasDiffError.value, hasDiffConflict, diffCounts)
  })

  onBeforeUnmount(() => {
    unmountCallbacks.forEach(cb => cb?.())
    const allRquestIds = Object.values(diffLoadingIds.value).reduce((acc, cur) => acc.concat(cur), [])
    http.cancelRequest(allRquestIds)
  })

  const findAddDelete = (selectedModels) => {
    const addSelect = []
    const deleteSelect = []
    const modelListLocalSet = new Set()
    modelListLocal.value.forEach((selectModel) => {
      modelListLocalSet.add(selectModel.id)
    })
    // 添加的
    selectedModels.forEach((selectModel) => {
      if (!modelListLocalSet.has(selectModel.id)) {
        addSelect.push(selectModel)
      } else {
        modelListLocalSet.delete(selectModel.id)
      }
    })
    // 删除的
    modelListLocal.value.forEach((model) => {
      if (modelListLocalSet.has(model.id)) {
        deleteSelect.push(model)
      }
    })
    return [addSelect, deleteSelect]
  }

  const handleConfirmAddModel = (selectedModels) => {
    const [addSelect, deleteSelect] = findAddDelete(selectedModels)

    modelListLocal.value.push(...addSelect)

    deleteSelect.forEach((item) => {
      handleClickRemoveModel(item, modelListLocal.value.indexOf(item))
    })

    // 获取新选择模型diff
    fetchFieldDiff(addSelect)
    fetchUniqueDiff(addSelect)
  }

  const isSelected = model => model.id === selectedModel.value?.id

  const isFieldConflict = model => fieldDiffCounts.value[model.id] && fieldDiffCounts.value[model.id].conflict
  const isUniqueConflict = model => uniqueDiffCounts.value[model.id] && uniqueDiffCounts.value[model.id].conflict
  const isConflict = model => isFieldConflict(model) || isUniqueConflict(model)
  const getConflictTips = (model) => {
    if (isFieldConflict(model) && isUniqueConflict(model)) {
      return t('当前模型与模板绑定会存在字段冲突和唯一校验冲突')
    }
    if (isFieldConflict(model)) {
      return t('当前模型与模板绑定会存在字段冲突')
    }
    if (isUniqueConflict(model)) {
      return t('当前模型与模板绑定会存在唯一校验冲突')
    }
  }

  const getTotal = id => (fieldDiffCounts.value[id]?.total ?? 0) + (uniqueDiffCounts.value[id]?.total ?? 0)

  const handleClickAddModel = () => {
    selectModelDialogRef.value.show()
  }
  const handleClickRemoveModel = (model, modelIndex) => {
    modelListLocal.value.splice(modelIndex, 1)

    // 变更当前选择的模型
    if (selectedModel.value === model) {
      selectedModel.value = null
    }

    // 删除对应模型的diff数据
    del(fieldDiffs.value, model.id)
    del(uniqueDiffs.value, model.id)

    del(modelEditAuths.value, model.id)
  }
  const handleSelectModel = (model) => {
    if (model.bk_ispaused) {
      return
    }
    selectedModel.value = model
  }
  const handleToggleTab = (tab) => {
    currentTab.value = tab
  }

  const handleModelAuthUpdate = (model, isPass) => {
    set(modelEditAuths.value, model.id, isPass)
    emit('update-model-auth', model, isPass)
  }

  defineExpose({
    modelList: modelListLocal.value
  })
</script>

<template>
  <div class="bind-model">
    <bk-resize-layout :collapsible="true" :min="310" :border="false" :style="{ height }">
      <div slot="aside" :class="['aside', { 'empty': !modelListLocal.length }]">
        <div class="aside-top">
          <div class="total" v-if="modelListLocal.length > 1">
            <i18n path="共N个">
              <template #count><em class="count">{{modelListLocal.length}}</em></template>
            </i18n>
          </div>
          <bk-button theme="primary" :outline="true" class="add-model-button"
            @click="handleClickAddModel"
            v-if="!readonly">
            <bk-icon type="plus" />{{ $t('添加绑定模型') }}
          </bk-button>
        </div>
        <div class="model-list" v-if="modelListLocal.length">
          <template v-for="(model, modelIndex) in modelListLocal">
            <cmdb-auth
              :key="modelIndex"
              tag="div"
              :auth="[
                { type: $OPERATION.U_MODEL, relation: [model.id] },
                { type: $OPERATION.R_MODEL, relation: [model.id] }
              ]"
              @update-auth="isPass => handleModelAuthUpdate(model, isPass)">
              <template #default="{ disabled }">
                <div :key="modelIndex"
                  :class="['model-item', {
                    disabled,
                    readonly,
                    'is-paused': model.bk_ispaused,
                    'is-loading': $loading(diffLoadingIds[model.id]),
                    'is-conflict': isConflict(model),
                    'is-selected': isSelected(model)
                  }]"
                  @click="handleSelectModel(model)">
                  <div class="model-icon-container">
                    <i class="model-icon" :class="model.bk_obj_icon"></i>
                  </div>
                  <div class="model-info">
                    <div class="model-name-area" :title="model.bk_obj_name">
                      <div class="model-name" :title="model.bk_obj_name">{{ model.bk_obj_name }}</div>
                      <mini-tag theme="paused" v-if="model.bk_ispaused">{{ $t('已停用') }}</mini-tag>
                    </div>
                    <div class="model-id">{{ model.bk_obj_id }}</div>
                  </div>
                  <div class="tail" v-if="!model.bk_ispaused">
                    <bk-button
                      v-if="!readonly"
                      class="remove-button"
                      :text="true"
                      @click.stop="handleClickRemoveModel(model, modelIndex)">
                      <bk-icon class="button-icon" type="delete" />
                    </bk-button>
                    <loading :loading="$loading(diffLoadingIds[model.id])">
                      <i class="bk-icon icon-exclamation-circle-shape conflict-icon"
                        v-if="isConflict(model)"
                        v-bk-tooltips="{ content: getConflictTips(model) }">
                      </i>
                      <span class="count-tag" v-else>
                        {{ getTotal(model.id) }}
                      </span>
                    </loading>
                  </div>
                </div>
              </template>
            </cmdb-auth>
          </template>
        </div>
        <bk-exception type="empty" scene="part" class="empty" v-if="readonly && !modelListLocal.length">
          <div>{{$t('该模板暂未绑定任何模型')}}</div>
        </bk-exception>
      </div>
      <div slot="main" :class="['main', { empty: !selectedModel }]">
        <div class="main-content" v-if="selectedModel">
          <div class="main-top">
            <div class="tab-list">
              <div :class="['tab-item', {
                     active: currentTab === 'field',
                     'is-conflict': fieldDiffCounts[selectedModel.id] && fieldDiffCounts[selectedModel.id].conflict > 0
                   }]"
                @click="handleToggleTab('field')">
                {{ $t('字段对比') }}
                <span class="count-tag" v-if="fieldDiffCounts[selectedModel.id]">
                  {{ fieldDiffCounts[selectedModel.id].total }}
                </span>
              </div>
              <div :class="['tab-item', {
                     active: currentTab === 'unique',
                     'is-conflict': uniqueDiffCounts[selectedModel.id]
                       && uniqueDiffCounts[selectedModel.id].conflict > 0
                   }]"
                @click="handleToggleTab('unique')">
                {{ $t('唯一校验对比') }}
                <span class="count-tag" v-if="uniqueDiffCounts[selectedModel.id]">
                  {{ uniqueDiffCounts[selectedModel.id].total }}
                </span>
              </div>
            </div>
          </div>
          <div class="main-body">
            <field-diff
              v-if="currentTab === 'field'"
              :model="selectedModel"
              :diffs="fieldDiffs[selectedModel.id]"
              :template-field-list="fieldList"
              :template-removed-field-list="removedFieldList">
            </field-diff>
            <unique-diff
              v-else-if="currentTab === 'unique'"
              :model="selectedModel"
              :diffs="uniqueDiffs[selectedModel.id]"
              :template-field-list="fieldList"
              :template-unique-list="uniqueList"
              :template-removed-unique-list="removedUniqueList">
            </unique-diff>
          </div>
        </div>
        <bk-exception v-else type="empty" scene="part" class="empty">
          <div>{{$t('暂无对比，请先绑定模型')}}</div>
        </bk-exception>
      </div>
    </bk-resize-layout>
    <select-model-dialog
      ref="selectModelDialogRef"
      v-if="!readonly"
      :selected="modelListLocal"
      :binded="bindedModelList"
      @confirm="handleConfirmAddModel">
    </select-model-dialog>
  </div>
</template>

<style lang="scss" scoped>
  .bind-model {

    .aside {
      height: 100%;
      padding: 0 24px 18px 24px;
      background: #fff;
      @include scrollbar-y(6px, white);

      &.empty {
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .aside-top {
        position: sticky;
        top: 0;
        padding-top: 18px;
        background: #fff;
      }
    }

    .main {
      height: 100%;
      &.empty {
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .empty {
        :deep(.bk-exception-img.part-img) {
          width: 240px;
          height: 180px;
        }
      }

      .main-content {
        height: 100%;
      }

      .main-top {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 48px;
        background: #fff;
        box-shadow: 0 1px 0 0 #DCDEE5;

        .tab-list {
          display: flex;
          gap: 8px;

          .tab-item {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 4px;
            width: 162px;
            height: 32px;
            background: #F5F7FA;
            border-radius: 21px;
            font-size: 14px;
            cursor: pointer;
            &.active {
              background: #E1ECFF;
              color: #3A84FF;
            }

            &.is-conflict {
              .count-tag {
                color: #EA3636;
                background: #FFDDDD;
              }
            }
          }
        }
      }

      .main-body {
        padding: 0 12px 24px 18px;
        height: calc(100% - 48px);
      }
    }
  }

  .count-tag {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    height: 16px;
    background: #F0F1F5;
    border-radius: 8px;
    font-family: ArialMT, Arial;
    font-size: 12px;
    color: #979BA5;
    padding: 0 8px;
  }

  .total {
    background: #F0F1F5;
    border-radius: 2px;
    height: 32px;
    line-height: 32px;
    text-align: center;
    font-size: 14px;
    margin-bottom: 16px;
    .count {
      font-style: normal;
      font-weight: 700;
      padding: .2em;
    }
  }

  .add-model-button {
    width: 100%;
    margin-bottom: 12px;
    .icon-plus {
      font-size: 20px !important;
    }
  }

  .model-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .model-item {
    display: flex;
    align-items: center;
    background: #FFFFFF;
    border: 1px solid #DCDEE5;
    border-radius: 2px;
    height: 56px;
    padding: 0 12px 0 16px;
    cursor: pointer;
    .model-icon-container {
      width: 32px;
      height: 32px;
      background: #e1ecffcc;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      .model-icon {
        color: #3A84FF;
        font-size: 14px;
      }
    }

    .model-info {
      flex: 1;
      margin-left: 12px;
      width: 0;
    }

    .tail {
      display: flex;
      margin-left: auto;
      align-items: center;
    }

    .model-name-area {
      display: flex;
      align-items: center;
      gap: 4px;
      .model-name {
        font-weight: 700;
        font-size: 12px;
        color: #63656E;
        @include ellipsis;
      }
    }
    .model-id {
      font-size: 12px;
      color: #979BA5;
      @include ellipsis;
    }

    .remove-button {
      visibility: hidden;
      color: #63656e;
      &:hover {
        color: #3a84ff;
      }
      .button-icon {
        font-size: 14px;
      }
      &.is-disabled {
        color: #c4c6cc;
      }
    }

    .conflict-icon {
      font-size: 14px;
      color: $dangerColor;
      margin-left: 4px;
    }

    &.is-conflict {
      border: 1px solid #DCDEE5;
      background: #FFEEEE;
      &.is-selected {
        background: #FFEEEE;
        border: 1px solid #EA3636;
      }
    }

    &.is-selected {
      background: #F0F5FF;
      border: 1px solid #3A84FF;
      .count-tag {
        background: #E1ECFF;
        color: #3A84FF;
      }
    }

    &:hover {
      &:not(.is-conflict) {
        border: 1px solid #3A84FF;
      }

      &:not(.is-loading) {
        .remove-button {
          visibility: visible;
        }
      }

      &:not(.readonly) {
        .count-tag {
          display: none;
        }
      }
    }

    &.disabled {
      cursor: not-allowed;
      opacity: 0.5;
    }
    &.is-paused {
      cursor: not-allowed;
      opacity: 0.5;
      pointer-events: none;
    }
  }

  [bk-language="en"] {
    .bind-model .main .main-top .tab-list {
      .tab-item {
        width: 262px;
      }
    }
  }
</style>
