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
  import Loading from '@/components/loading/index.vue'
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
    // 模板唯一校验列表
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    readonly: {
      type: Boolean,
      default: false
    },
    height: String
  })

  const emit = defineEmits(['update-diffs'])

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

  watchEffect(async () => {
    const initModelList = props.modelList.slice()
    selectedModel.value = initModelList?.[0] ?? null
    if (initModelList.length) {
      fetchFieldDiff(initModelList)
      fetchUniqueDiff(initModelList)
    }

    modelListLocal.value = initModelList
  })

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
  const displayDiffCounts = computed(() => {
    const countMap = {
      field: fieldDiffCounts.value,
      unique: uniqueDiffCounts.value,
    }
    return countMap[currentTab.value]
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

  const fetchFieldDiff = async (modelList) => {
    const modelIds = modelList.map(item => item.id)
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
        set(fieldDiffs.value, modelIds[i], value ?? {})
      }
    }

    unmountCallbacks.push(() => allResult?.return())
  }

  const fetchUniqueDiff = async (modelList) => {
    const modelIds = modelList.map(item => item.id)
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
        set(uniqueDiffs.value, modelIds[i], value ?? {})
      }
    }

    unmountCallbacks.push(() => allResult?.return())
  }

  onBeforeUnmount(() => {
    unmountCallbacks.forEach(cb => cb?.())
    const allRquestIds = Object.values(diffLoadingIds.value).reduce((acc, cur) => acc.concat(cur), [])
    http.cancelRequest(allRquestIds)
  })

  const handleConfirmAddModel = (selectedModels) => {
    modelListLocal.value.push(...selectedModels)
    if (!selectedModel.value) {
      selectedModel.value = modelListLocal.value?.[0]
    }

    // 获取新选择模型diff
    fetchFieldDiff(selectedModels)
    fetchUniqueDiff(selectedModels)
  }

  const isSelected = model => model.id === selectedModel.value.id
  const isConflict = model => displayDiffCounts.value[model.id] && displayDiffCounts.value[model.id].conflict
  const handleClickAddModel = () => {
    selectModelDialogRef.value.show()
  }
  const handleClickRemoveModel = (model, modelIndex) => {
    modelListLocal.value.splice(modelIndex, 1)

    // 变更当前选择的模型
    if (selectedModel.value === model) {
      selectedModel.value = modelListLocal.value[modelIndex + 1] ?? modelListLocal.value[0]
    }

    // 删除对应模型的diff数据
    del(fieldDiffs.value, model.id)
    del(uniqueDiffs.value, model.id)
  }
  const handleSelectModel = (model) => {
    selectedModel.value = model
  }
  const handleToggleTab = (tab) => {
    currentTab.value = tab
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
          <div class="total" v-if="modelListLocal.length > 0">
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
          <div v-for="(model, modelIndex) in modelListLocal" :key="modelIndex"
            :class="['model-item', {
              readonly,
              'is-loading': $loading(diffLoadingIds[model.id]),
              'is-conflict': isConflict(model),
              'is-selected': isSelected(model)
            }]"
            @click="handleSelectModel(model)">
            <div class="model-icon-container">
              <i class="model-icon" :class="model.bk_obj_icon"></i>
            </div>
            <div class="model-info">
              <div class="model-name" :title="model.bk_obj_name">{{ model.bk_obj_name }}</div>
              <div class="model-id">{{ model.bk_obj_id }}</div>
            </div>
            <div class="tail">
              <bk-button
                v-if="!readonly"
                class="remove-button"
                :text="true"
                @click.stop="handleClickRemoveModel(model, modelIndex)">
                <bk-icon class="button-icon" type="delete" />
              </bk-button>
              <loading :loading="$loading(diffLoadingIds[model.id])">
                <i class="bk-icon icon-exclamation-circle-shape conflict-icon" v-if="isConflict(model)"></i>
                <span class="count" v-else>{{ displayDiffCounts[model.id] && displayDiffCounts[model.id].total }}</span>
              </loading>
            </div>
          </div>
        </div>
        <bk-exception type="empty" scene="part" class="empty" v-if="readonly && !modelListLocal.length">
          <div>{{$t('暂无绑定模型')}}</div>
        </bk-exception>
      </div>
      <div slot="main" :class="['main', { empty: !selectedModel }]">
        <div class="main-content" v-if="selectedModel">
          <div class="main-top">
            <div class="tab-list">
              <div :class="['tab-item', { active: currentTab === 'field' }]"
                @click="handleToggleTab('field')">{{ $t('字段对比') }}</div>
              <div :class="['tab-item', { active: currentTab === 'unique' }]"
                @click="handleToggleTab('unique')">{{ $t('唯一校验对比') }}</div>
            </div>
          </div>
          <div class="main-body">
            <field-diff
              v-if="currentTab === 'field'"
              :model="selectedModel"
              :diffs="fieldDiffs[selectedModel.id]"
              :template-field-list="fieldList">
            </field-diff>
            <unique-diff
              v-else-if="currentTab === 'unique'"
              :model="selectedModel"
              :diffs="uniqueDiffs[selectedModel.id]"
              :template-unique-list="uniqueList"
              :template-field-list="fieldList">
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
      @include scrollbar-y;

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
            width: 162px;
            height: 32px;
            line-height: 32px;
            background: #F5F7FA;
            border-radius: 21px;
            font-size: 14px;
            text-align: center;
            cursor: pointer;
            &.active {
              background: #E1ECFF;
              color: #3A84FF;
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

    .model-name {
      font-weight: 700;
      font-size: 12px;
      color: #63656E;
      @include ellipsis;
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

    .count {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      height: 16px;
      background: #F0F1F5;
      border-radius: 8px;
      font-size: 12px;
      color: #979BA5;
      padding: 0 8px;
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
      .count {
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
        .count {
          display: none;
        }
      }
    }
  }
</style>
