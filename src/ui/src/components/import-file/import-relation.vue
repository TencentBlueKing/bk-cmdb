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
  <div class="import-relation">
    <bk-checkbox class="allow-import-checkbox" v-model="allowImport">{{$t('是否导入关联的模型实例')}}</bk-checkbox>
    <div v-show="allowImport">
      <div class="allow-import-model">
        <i :class="['model-icon', getModelIcon(currentModelId)]"></i>
        <span class="model-name">{{getModelName(currentModelId)}}</span>
        <span class="model-desc">{{`(${$t('模型本身')})`}}</span>
        <div class="unique-flag" v-if="modelUniqueChecks.length">
          <label class="unique-flag-label">{{$t('唯一校验标识')}}：</label>
          <bk-select class="unique-selector"
            v-model="objectUniqueId"
            :clearable="false">
            <bk-option v-for="uniqueCheck in modelUniqueChecks"
              :key="uniqueCheck.id"
              :id="uniqueCheck.id"
              :name="getUniqueCheckName(uniqueCheck)">
            </bk-option>
          </bk-select>
        </div>
      </div>
      <bk-table class="relation-table"
        ref="table"
        :outer-border="false"
        :max-height="$APP.height - 280"
        :data="computedList"
        @select="handleSelect"
        @select-all="handleSelectAll">
        <bk-table-column type="selection" :selectable="() => allowImport"></bk-table-column>
        <bk-table-column :label="$t('关联的模型')" prop="model" width="250" :resizable="false">
          <template slot-scope="{ row }">
            <i :class="['model-icon', getModelIcon(row.bk_obj_id)]"></i>
            <span class="model-name">{{getModelName(row.bk_obj_id)}}</span>
            <span class="model-desc" v-if="row.bk_obj_id === currentModelId">{{`(${$t('自关联')})`}}</span>
          </template>
        </bk-table-column>
        <bk-table-column :label="$t('唯一校验标识')" prop="identification" align="right" :resizable="false">
          <template slot-scope="{ row }">
            <bk-select class="unique-selector"
              :value="selectedRelations[row.bk_obj_id]"
              :clearable="false"
              :disabled="isUniqueCheckDisabled(row)"
              @change="handleUniqueCheckChange(row, ...arguments)">
              <bk-option v-for="uniqueCheck in row.relation_unique_checks"
                :key="uniqueCheck.id"
                :id="uniqueCheck.id"
                :name="getUniqueCheckName(uniqueCheck)">
              </bk-option>
            </bk-select>
          </template>
        </bk-table-column>
        <bk-exception slot="empty" type="empty" scene="part">{{$t('暂无关联模型，无需选择')}}</bk-exception>
      </bk-table>
    </div>
    <div class="options">
      <bk-button theme="default" @click="previousStep">{{$t('上一步')}}</bk-button>
      <bk-button theme="primary" class="ml10" :disabled="importDisabled" @click="startImport">{{$t('开始导入')}}</bk-button>
      <bk-button theme="default" class="ml10" @click="closeImport">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import { ref, watch, toRef, computed, set, del } from '@vue/composition-api'
  import useStep from './step'
  import useImport from './index'
  import useFile from './file'
  import useProperty from '@/hooks/model/property'
  import useBatchUniqueCheck from '@/hooks/unique-check/batch'
  import useUniqueCheck from '@/hooks/unique-check'
  import usePending from '@/hooks/utils/pending'
  import { mapGetters } from 'vuex'
  export default {
    name: 'import-relation',
    setup(props, setupContext) {
      const [currentStep, { previous: previousStep }] = useStep()
      const [importState, { close: closeImport }] = useImport()
      const [{ file, response }, {
        setState: setFileState,
        setError: setFileError,
        setResponse: setFileResponse
      }] = useFile()
      const currentModelId = toRef(importState, 'bk_obj_id')
      // 由接口解析出来的excel中的关联模型
      const computedRelationMap = computed(() => (response.value.association || {}))
      const computedRelationModels = computed(() => Object.keys(computedRelationMap.value))

      // 获取当前模型的唯一校验，用于导出的参数object_unique_id
      const [{ uniqueChecks: modelUniqueChecks, pending: modelUniquePending }] = useUniqueCheck(currentModelId)

      // 加载相关模型属性
      const [{ properties, pending: propertyPending }] = useProperty({
        bk_obj_id: { $in: [...computedRelationModels.value, currentModelId.value] }
      })

      // 加载关联模型的唯一校验
      const [{
        uniqueChecks: relationUniqueChecks,
        pending: uniqueCheckPending
      }] = useBatchUniqueCheck(computedRelationModels)
      // 计算显示列表
      const computedList = computed(() => {
        const list = computedRelationModels.value.map((modelId) => {
          const item = { bk_obj_id: modelId }
          const uniqueChecks = relationUniqueChecks.value.find((uniqueChecks) => {
            const isMatched = uniqueChecks.some(uniqueCheck => uniqueCheck.bk_obj_id === modelId)
            return isMatched
          })
          item.relation_unique_checks = uniqueChecks || []
          return item
        })
        return list
      })
      // 选择是否导出关联实例
      const allowImport = ref(false)
      const objectUniqueId = ref('')
      const selectedRelations = ref({})
      const initSelectedUniqueCheck = () => {
        const [modelUniqueCheck] = modelUniqueChecks.value
        objectUniqueId.value = modelUniqueCheck?.id
      }
      const clearSelectedUniqueCheck = (clearSelf = true) => {
        setupContext.refs.table.clearSelection() // 这种方式在Vue3.0中不可使用
        clearSelf && (objectUniqueId.value = '')
        computedList.value.forEach((row) => {
          del(selectedRelations.value, row.bk_obj_id)
        })
      }
      watch(allowImport, (value) => {
        value ? initSelectedUniqueCheck() : clearSelectedUniqueCheck()
      })

      // 表格勾选操作
      const setSelectionUniqueCheck = (row) => {
        const [firstUniqueCheck] = row.relation_unique_checks
        const preset = row.relation_unique_checks.find(uniqueCheck => uniqueCheck.ispre) || firstUniqueCheck
        set(selectedRelations.value, row.bk_obj_id, preset.id)
      }
      const handleSelect = (selection, row) => {
        const selected = selection.includes(row)
        if (selected) {
          setSelectionUniqueCheck(row)
        } else {
          del(selectedRelations.value, row.bk_obj_id)
        }
      }
      const handleSelectAll = (selection) => {
        if (!selection.length) {
          clearSelectedUniqueCheck(false)
        } else {
          selection.forEach(setSelectionUniqueCheck)
        }
      }
      const handleUniqueCheckChange = (row, id) => {
        if (id) {
          set(selectedRelations.value, row.bk_obj_id, id)
        } else {
          del(selectedRelations.value, row.bk_obj_id)
        }
      }
      const pending = usePending([modelUniquePending, propertyPending, uniqueCheckPending], true)
      const importDisabled = computed(() => {
        if (!allowImport.value) {
          return false
        }
        return Object.keys(selectedRelations.value).length === 0
      })
      return {
        currentStep,
        previousStep,
        closeImport,
        importState,
        file,
        response,
        setFileState,
        setFileError,
        setFileResponse,
        computedList,
        currentModelId,
        objectUniqueId,
        properties,
        allowImport,
        handleSelect,
        handleSelectAll,
        selectedRelations,
        handleUniqueCheckChange,
        modelUniqueChecks,
        pending,
        importDisabled
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById'])
    },
    methods: {
      getRelationModel(modelId) {
        return this.getModelById(modelId) || { bk_obj_id: modelId }
      },
      getModelIcon(modelId) {
        const model = this.getRelationModel(modelId)
        return model.bk_obj_icon || 'icon-cc-default'
      },
      getModelName(modelId) {
        const model = this.getRelationModel(modelId)
        return model.bk_obj_name || model.bk_obj_id
      },
      getUniqueCheckName({ keys }) {
        const idArray = keys.map(key => key.key_id)
        return idArray.map((id) => {
          const property = this.properties.find(property => property.id === id)
          return property ? property.bk_property_name : `${this.$t('未知属性')}(${id})`
        }).join(' + ')
      },
      isUniqueCheckDisabled(row) {
        return !this.$refs.table.selection.includes(row)
      },
      async startImport() {
        try {
          this.setFileState('pending')
          const response = await this.importState.submit({
            file: this.file,
            step: this.currentStep,
            importRelation: this.allowImport,
            object_unique_id: this.objectUniqueId,
            relations: this.selectedRelations,
            config: {
              transformData: false,
              globalError: false
            }
          })
          if (!response.result || response.data.error) {
            this.setFileState('error')
            const importError = response.data?.error ? response.data : { error: [response.bk_error_msg] }
            this.setFileError(importError)
          } else {
            this.setFileState('success')
            this.importState.success && this.importState.success()
          }
        } catch (error) {
          console.error(error)
          this.setFileState('error')
          this.setFileError(error)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  @mixin hoverUniqueSelector {
    background: #f0f1f5;
  }
  @mixin hoverUniqueSelectorDeep {
    .bk-select-angle {
      display: inline-block;
    }
    .bk-select-name {
      padding: 0 36px 0 10px;
    }
  }
  .import-relation {
    padding: 20px 0 0 0;
    .allow-import-checkbox {
      display: block;
      margin: 0 16px;
    }
    .allow-import-model {
      display: flex;
      align-items: center;
      height: 44px;
      border: 1px solid #dcdee5;
      border-radius: 3px;
      padding: 0 16px;
      margin: 13px 0 0 0;
      .unique-flag {
        margin-left: auto;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        .unique-flag-label {
          display: block;
          font-size: 12px;
          white-space: nowrap;
        }
      }
    }
    .model-icon {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      width: 26px;
      height: 26px;
      background: #f0f1f5;
      border-radius: 50%;
      font-size: 14px;
      margin-right: 10px;
    }
    .model-name {
      font-size: 12px;
      line-height: 16px;
    }
    .model-desc {
      font-size: 12px;
      color: #c4c6cc;
      line-height: 16px;
    }
    .unique-selector {
      border-color: transparent;
      background-color: transparent;
      text-align: right;
      max-width: 300px;
      &:before {
        left: auto;
        right: 10px;
      }
      /deep/ {
        .bk-select-angle {
          display: none;
        }
        .bk-select-name {
          padding: 0;
          font-size: 12px;
        }
      }
      &:hover {
        &:before {
          right: 30px;
        }
        @include hoverUniqueSelector;
        /deep/ {
          @include hoverUniqueSelectorDeep;
        }
      }
      &.is-focus {
        box-shadow: none;
        &:before {
          right: 30px;
        }
        @include hoverUniqueSelector;
        /deep/ {
          @include hoverUniqueSelectorDeep;
        }
      }
      &.is-disabled {
        text-align: right;
        background-color: transparent !important;
        border-color: transparent !important;
        pointer-events: none;
        cursor: not-allowed;
      }
    }
  }
  .relation-table {
    margin: 10px 0 0 0;
    /deep/ .bk-table-row {
      .unique-selector {
        min-width: 100px;
      }
      &.hover-row > td {
        background-color: #fff;
      }
      &.hover-row .unique-selector:not(.is-disabled) {
        &:before {
          right: 30px;
        }
        @include hoverUniqueSelector;
        @include hoverUniqueSelectorDeep;
      }
    }
  }
  .options {
    margin: 20px 0 0 0;
    display: flex;
  }
</style>
