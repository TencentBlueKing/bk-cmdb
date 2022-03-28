<template>
  <div class="model-relation">
    <div class="model-relation-header">
      {{ model.bk_obj_name }}（{{ model.bk_obj_id }}）{{ t("XX的关联关系") }}
    </div>
    <bk-table
      ref="tableRef"
      v-bkloading="{ isLoading: relationsLoading }"
      v-bind="tableAttrs"
      :data="relations"
      class="model-relation-table"
      row-key="id"
      @select-all="handleSelect"
      @select="handleSelect"
    >
      <bk-table-column width="40px" type="selection" :selectable="relationSelectable"></bk-table-column>
      <bk-table-column :label="t('关系描述')">
        <template #default="{ row: relation }">
          <slot name="relation-desc" v-bind="{ relation }">
            {{ relation.targetModel.bk_obj_name }}
            <span class="model-relation-name">{{ getRelationName(relation.bk_asst_id) }}</span>
            {{ model.bk_obj_name }}
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
  import { defineComponent, ref, watch, toRef, nextTick } from '@vue/composition-api'
  import { find as findModelRelations } from '@/service/model/association.js'
  import { find as findRelations } from '@/service/association/index.js'
  import store from '@/store'
  import { t } from '@/i18n'
  import to from 'await-to-js'

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
      const relationTypesLoading = ref(false)
      let relationTypes = []
      const relations = ref([])
      const relationsLoading = ref(false)
      const modelRef = toRef(props, 'model')
      const { model, excludedRelations } = props
      const tableRef = ref(null)

      const loadRelationTypes = () => {
        relationTypesLoading.value = true
        return findRelations().then(({ info }) => {
          relationTypes = info
        })
          .finally(() => {
            relationTypesLoading.value = false
          })
      }

      const findModelById = modelId => store.getters['objectModelClassify/getModelById'](modelId)

      const handleSelect = (selectedRelations) => {
        const unselectedRelations = relations.value
          .filter(relation => !selectedRelations.some(selectedRelation => relation.id === selectedRelation.id))

        if (unselectedRelations.length > 0) {
          excludedRelations[model.bk_obj_id] = unselectedRelations
        } else {
          delete excludedRelations[model.bk_obj_id]
        }

        emit('update:excludedRelations', excludedRelations)
      }

      const getRelationName = relationId => relationTypes
        .find(relation => relation.bk_asst_id === relationId)?.bk_asst_name || relationId


      const loadRelations = (modelId) => {
        if (!modelId) Promise.reject()

        relationsLoading.value = true

        return findModelRelations(modelId, 'target', {
          cancelPrevious: true
        })
          .then((targetRelations) => {
            relations.value = targetRelations.map(relation => ({
              ...relation,
              targetModel: findModelById(relation.bk_obj_id)
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

      loadRelationTypes()

      const genRelationSelection = () => {
        const modelExcludedRelations = excludedRelations?.[model.bk_obj_id]
        if (relations.value.length > 0) {
          nextTick(() => {
            relations.value.forEach((relation) => {
              const isSelected = !modelExcludedRelations
                ?.some(excludedRelation => excludedRelation.id === relation.id)
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
            relations.value = newModel?.bk_obj_asst || []
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
        relationTypesLoading,
        relationsLoading,
        getRelationName,
        handleSelect,
        relationTypes,
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
}
</style>
