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
  <div class="model-relation">
    <div class="model-relation-header">
      <i18n path="X的关联关系">
        <template #model>
          <span>{{ model.bk_obj_name }}（{{ model.bk_obj_id }}）</span>
        </template>
      </i18n>
    </div>
    <bk-table
      ref="tableRef"
      v-bkloading="{ isLoading: relationsLoading }"
      v-bind="tableAttrs"
      :data="relations"
      class="model-relation-table"
      row-key="bk_obj_asst_id"
      @select-all="handleSelect"
      @select="handleSelect"
    >
      <bk-table-column width="40px" type="selection" :selectable="relationSelectable"></bk-table-column>
      <bk-table-column :label="t('关联类型')">
        <template #default="{ row: relation }">
          <slot name="relation-desc" v-bind="{ relation }">
            {{ getRelationName(relation.bk_asst_id) }}<span class="model-relation-id">（{{relation.bk_asst_id}}）</span>
          </slot>
        </template>
      </bk-table-column>
      <bk-table-column :label="t('关联对象')">
        <template #default="{ row: relation }">
          <slot name="relation-target" v-bind="{ relation }"></slot>
        </template>
      </bk-table-column>
    </bk-table>
  </div>
</template>

<script>
  import { defineComponent, ref, watch, toRef, nextTick } from 'vue'
  import { find as findModelRelations } from '@/service/model/association.js'
  import { useAssociations } from '../../use-associations.js'
  import store from '@/store'
  import { t } from '@/i18n'
  import to from 'await-to-js'
  import cloneDeep from 'lodash/cloneDeep'

  export default defineComponent({
    name: 'ModelRelation',
    props: {
      // 模型数据
      model: {
        type: Object,
        required: true
      },
      // 关联模型表格的属性
      tableAttrs: {
        type: Object,
        default: () => ({})
      },
      // 需要排除的关联关系
      excludedRelations: {
        type: Object,
        default: () => ({})
      },
      // 关联关系是否可选
      relationSelectable: {
        type: Function,
        default: () => true
      },
      // 是否远程加载模型关联关系，默认使用传入的关联关系数据
      remote: {
        type: Boolean,
        default: false
      }
    },
    setup(props, { emit }) {
      const { associations: allAssociations } = useAssociations()// 所有的模型关联类型
      const relations = ref([])
      const relationsLoading = ref(false)
      const modelRef = toRef(props, 'model')
      const excludedRelations = cloneDeep(props.excludedRelations)
      const tableRef = ref(null)

      const findModelById = modelId => store.getters['objectModelClassify/getModelById'](modelId)

      const handleSelect = (selectedRelations) => {
        const unselectedRelations = relations.value
          .filter(relation => !selectedRelations
            .some(selectedRelation => relation.bk_obj_asst_id === selectedRelation.bk_obj_asst_id))

        if (unselectedRelations.length > 0) {
          toRef(excludedRelations, modelRef.value.bk_obj_id).value = unselectedRelations
        } else {
          delete excludedRelations[modelRef.value.bk_obj_id]
        }

        emit('update:excludedRelations', cloneDeep(excludedRelations))
      }

      const getRelationName = relationId => allAssociations.value
        .find(relation => relation.bk_asst_id === relationId)?.bk_asst_name || relationId


      const loadRelations = (modelId) => {
        if (!modelId) return Promise.reject().catch(() => {
          console.log('modelId is missing')
        })

        relationsLoading.value = true

        return findModelRelations(modelId, 'source', {
          cancelPrevious: true
        })
          .then((targetRelations) => {
            relations.value = targetRelations.map(relation => ({
              ...relation,
              targetModel: findModelById(relation.bk_asst_obj_id)
            }))
          })
          .catch((err) => {
            relations.value = []
            console.log(err)
          })
          .finally(() => {
            relationsLoading.value = false
          })
      }

      const genRelationSelection = () => {
        const modelExcludedRelations = excludedRelations?.[modelRef.value.bk_obj_id]
        if (relations.value.length > 0) {
          nextTick(() => {
            relations.value.forEach((relation) => {
              const isSelected = !modelExcludedRelations
                ?.some(excludedRelation => excludedRelation.bk_obj_asst_id === relation.bk_obj_asst_id)
              tableRef.value.toggleRowSelection(relation, isSelected)
            })
          })
        }
      }

      watch(
        modelRef, async (newModel) => {
          if (props.remote) {
            await to(loadRelations(newModel.bk_obj_id))
          } else {
            relations.value = newModel?.object_asst || []
          }
          genRelationSelection()
        },
        {
          immediate: true,
          deep: true
        }
      )

      return {
        t,
        relationsLoading,
        getRelationName,
        handleSelect,
        tableRef,
        relations
      }
    },
  })
</script>

<style lang="scss" scoped>
.model-relation {
  height: 100%;

  &-header {
    height: 60px;
    display: flex;
    align-items: center;
    font-size: 16px;
    color: #000;
  }

  &-name {
    display: inline-block;
    padding: 0 5px;
    background: #f0f1f5;
    border-radius: 2px;
    line-height: 18px;
  }

  &-id {
    color: #979BA5;
  }
}
</style>
