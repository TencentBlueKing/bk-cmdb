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
  <div class="export-relation" v-bkloading="{ isLoading: pending }">
    <bk-checkbox class="allow-export-checkbox" v-model="allowExport">{{$t('是否导出关联的模型实例')}}</bk-checkbox>
    <div v-show="allowExport">
      <div class="allow-export-model">
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
        :data="computedRelations"
        @select="handleSelect"
        @select-all="handleSelectAll">
        <bk-table-column type="selection" :selectable="() => allowExport"></bk-table-column>
        <bk-table-column :label="$t('关联的模型')" prop="model" width="250" :resizable="false">
          <template slot-scope="{ row }">
            <i :class="['model-icon', getModelIcon(row)]"></i>
            <span class="model-name">{{getModelName(row)}}</span>
            <span class="model-desc" v-if="row.bk_obj_id === row.bk_asst_obj_id">{{`(${$t('自关联')})`}}</span>
          </template>
        </bk-table-column>
        <bk-table-column :label="$t('唯一校验标识')" prop="identification" align="right" :resizable="false">
          <template slot-scope="{ row }">
            <bk-select class="unique-selector"
              :value="selectedRelations[row.relation_obj_id]"
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
  </div>
</template>

<script>
  import { computed, watch, toRef } from '@vue/composition-api'
  import useState from './state'
  import useModelAssociation from '@/hooks/model/association'
  import useBatchUniqueCheck from '@/hooks/unique-check/batch'
  import useUniqueCheck from '@/hooks/unique-check'
  import useProperty from '@/hooks/model/property'
  import usePending from '@/hooks/utils/pending'
  import { mapGetters } from 'vuex'
  export default {
    name: 'export-relation',
    setup(props, setupContext) {
      const [state, { setRelation, removeRelation }] = useState()
      const objectUniqueId = toRef(state, 'object_unique_id')
      const currentModelId = toRef(state, 'bk_obj_id')
      const selectedRelations = toRef(state, 'relations')
      // 获取当前模型的唯一校验，用于导出的参数object_unique_id
      const [{ uniqueChecks: modelUniqueChecks, pending: modelUniquePending }] = useUniqueCheck(currentModelId)

      // 获取当前模型的关联关系
      const [{ relations, pending: relationPending }] = useModelAssociation(currentModelId)

      // 加载关联模型的属性
      const relationModels = computed(() => {
        const modelSet = new Set()
        relations.value.forEach((item) => {
          const modelId = item.bk_obj_id
          const asstModelId = item.bk_asst_obj_id
          modelSet.add(modelId === currentModelId.value ? asstModelId : modelId)
        })
        return Array.from(modelSet)
      })
      const propertyOptions = computed(() => {
        if (!relationModels.value.length) {
          return {
            bk_obj_id: currentModelId.value
          }
        }
        const modelSet = new Set(relationModels.value)
        modelSet.add(currentModelId.value)
        return { bk_obj_id: { $in: Array.from(modelSet) } }
      })
      const [{ properties, pending: propertyPending }] = useProperty(propertyOptions)

      // 加载关联模型的唯一校验
      const [{ uniqueChecks: relationUniqueChecks, pending: uniqueCheckPending }] = useBatchUniqueCheck(relationModels)

      // 组合关联模型与唯一校验
      const uniqueRelations = computed(() => {
        const uniqueSet = new Set()
        const result = []
        relations.value.forEach((relation) => {
          const modelId = relation.bk_obj_id === currentModelId.value ? relation.bk_asst_obj_id : relation.bk_obj_id
          if (uniqueSet.has(modelId)) return
          uniqueSet.add(modelId)
          result.push(relation)
        })
        return result
      })
      const computedRelations = computed(() => uniqueRelations.value.map((relation) => {
        const modelId = relation.bk_obj_id === currentModelId.value ? relation.bk_asst_obj_id : relation.bk_obj_id
        const uniqueChecks = relationUniqueChecks.value.find((uniqueChecks) => {
          const isMatched = uniqueChecks.some(uniqueCheck => uniqueCheck.bk_obj_id === modelId)
          return isMatched
        })
        return { relation_obj_id: modelId, relation_unique_checks: uniqueChecks || [], ...relation }
      }))

      // 选择是否导出关联实例
      const allowExport = toRef(state, 'exportRelation')
      const initSelectedUniqueCheck = () => {
        const [modelUniqueCheck] = modelUniqueChecks.value
        state.object_unique_id.value = modelUniqueCheck?.id
      }
      const clearSelectedUniqueCheck = (clearSelf = true) => {
        setupContext.refs.table.clearSelection() // 这种方式在Vue3.0中不可使用
        clearSelf && (state.object_unique_id.value = '')
        computedRelations.value.forEach((relation) => {
          removeRelation(relation.relation_obj_id)
        })
      }
      watch(allowExport, (value) => {
        value ? initSelectedUniqueCheck() : clearSelectedUniqueCheck()
      })

      // 表格勾选操作
      const setSelectionUniqueCheck = (relation) => {
        const [firstUniqueCheck] = relation.relation_unique_checks
        const preset = relation.relation_unique_checks.find(uniqueCheck => uniqueCheck.ispre) || firstUniqueCheck
        setRelation(relation.relation_obj_id, preset.id)
      }
      const handleSelect = (selection, relation) => {
        const selected = selection.includes(relation)
        if (selected) {
          setSelectionUniqueCheck(relation)
        } else {
          removeRelation(relation.relation_obj_id)
        }
      }
      const handleSelectAll = (selection) => {
        if (!selection.length) {
          clearSelectedUniqueCheck(false)
        } else {
          selection.forEach(setSelectionUniqueCheck)
        }
      }
      const handleUniqueCheckChange = (relation, id) => {
        if (id) {
          setRelation(relation.relation_obj_id, id)
        } else {
          removeRelation(relation.relation_obj_id)
        }
      }
      // 总的请求状态
      const pending = usePending([modelUniquePending, relationPending, propertyPending, uniqueCheckPending], true)
      return {
        state,
        objectUniqueId,
        currentModelId,
        properties,
        computedRelations,
        selectedRelations,
        pending,
        allowExport,
        handleSelect,
        handleSelectAll,
        handleUniqueCheckChange,
        modelUniqueChecks
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById'])
    },
    methods: {
      getRelationModel(relation) {
        let modelId
        if (typeof relation === 'string') {
          modelId = relation
        } else {
          modelId = relation.bk_obj_id === this.currentModelId ? relation.bk_asst_obj_id : relation.bk_obj_id
        }
        return this.getModelById(modelId) || { bk_obj_id: modelId }
      },
      getModelIcon(relation) {
        const model = this.getRelationModel(relation)
        return model.bk_obj_icon || 'icon-cc-default'
      },
      getModelName(relation) {
        const model = this.getRelationModel(relation)
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
  .export-relation {
    padding: 20px 0 0 0;
    .allow-export-checkbox {
      display: block;
      margin: 0 16px;
    }
    .allow-export-model {
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
</style>
