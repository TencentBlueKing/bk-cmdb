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
  <content-confirm-layout>
    <template #aside>
      <div class="view-tab">
        <div
          class="view-tab-item"
          @click="handleTabChange(viewTab)"
          v-for="viewTab in viewTabs"
          :key="viewTab.name"
          :class="{
            'is-active': activeViewTabName === viewTab.name
          }">
          {{viewTab.label}}
          <span class="item-count">{{viewTab.count()}}</span>
        </div>
      </div>

      <model-group-list
        v-show="activeViewTabName === MODEL_TAB_NAME"
        :selected-model-id.sync="selectedModelId"
        @model-select="setCurrentModel"
        :groups="modelGroups">
        <template #group-header-append="{ modelGroup }">
          <template v-if="!isGroupExistedInTargetEnv(modelGroup)">
            <i
              v-bk-tooltips="{
                content: t('该分组在目标环境中未存在，默认放置于未分类组')
              }"
              class="icon-cc-remind"
            />
            <dropdown-select
              v-model="importGroupCreationTypes[modelGroup.bk_classification_id]"
              class="model-group-dropdown"
              :options="[
                { value: JOIN_UNCATEGORIZED_GROUP, label: t('加入未分类分组') },
                { value: CREATE_GROUP, label: t('新建该分组') }]">
            </dropdown-select>
          </template>
        </template>
        <template #model-append="{ model }">
          <i
            v-if="isModelExistedInTargetEnv(model)"
            v-bk-tooltips="{
              content: t('模型在目标环境中已存在，不可再导入')
            }"
            class="icon-cc-remind"
          />
          <bk-switcher
            size="small"
            theme="primary"
            @click.native.stop
            @change="handleModelSwitchChange(model)"
            v-model="includedModels[model.bk_obj_id]"
            :disabled="isModelExistedInTargetEnv(model)">
          </bk-switcher>
        </template>
      </model-group-list>

      <div class="relation-type-list" v-show="activeViewTabName === RELATION_TAB_NAME">
        <template v-if="importRelationTypes.length > 0">
          <div
            class="relation-type-item"
            :class="{
              'is-active': relation.bk_asst_id === currentRelationType.bk_asst_id
            }"
            v-for="relation in importRelationTypes" :key="relation.bk_asst_id"
            @click="setCurrentRelationType(relation)">
            {{relation.bk_asst_name}}
            <span class="relation-type-id">（{{relation.bk_asst_id}}）</span>
            <bk-tag theme="warning" v-show="isRelationTypeConflict(relation)">{{t('冲突')}}</bk-tag>
          </div>
        </template>
        <bk-exception
          v-else
          type="empty" scene="part">{{t('当前导入内容没有关联类型')}}</bk-exception>
      </div>
    </template>
    <template #main>
      <div class="relation-type-diff" v-show="activeViewTabName === RELATION_TAB_NAME">
        <template v-if="importRelationTypes.length > 0">
          <div class="relation-type-diff-header">
            {{ currentRelationType.bk_asst_name }}（{{ currentRelationType.bk_asst_name }}）
            <bk-tag theme="warning" v-if="currentRelationType.isConflict">{{t('冲突')}}</bk-tag>
          </div>
          <bk-form class="conflict-resolve-form" v-show="currentRelationType.isConflict">
            <bk-form-item :label-width="140" required :label="t('关联类型冲突处理')">
              <bk-select
                :clearable="false"
                v-model="excludedRelationTypes[currentRelationType.bk_asst_id]" class="resolve-method-select">
                <bk-option :id="COVER" :name="t('覆盖更新')"></bk-option>
                <bk-option :id="SKIP" :name="t('跳过')"></bk-option>
              </bk-select>
            </bk-form-item>
          </bk-form>
          <bk-table :data="currentRelationType.relationDiffs" class="relation-type-relation-models">
            <bk-table-column v-if="currentRelationType.isConflict" width="100">
              <template #default="{ row }">
                <bk-tag v-if="row.isOld" theme="warning">{{t('原')}}</bk-tag>
                <bk-tag v-if="row.isNew" theme="success">{{t('新')}}</bk-tag>
              </template>
            </bk-table-column>
            <bk-table-column :label="t('名称')">
              <template #default="{ row }">
                <span
                  :class="{
                    'is-diff': row.isNew && currentRelationType.conflictProperties.includes('bk_asst_name')
                  }">{{row.bk_asst_name}}</span>
              </template>
            </bk-table-column>
            <bk-table-column :label="t('唯一标志')" property="bk_asst_id"></bk-table-column>
            <bk-table-column :label="t('源->目标描述')">
              <template #default="{ row }">
                <span :class="{
                  'is-diff': row.isNew && currentRelationType.conflictProperties.includes('src_des')
                }">{{row.src_des}}</span>
              </template>
            </bk-table-column>
            <bk-table-column :label="t('目标->源描述')">
              <template #default="{ row }">
                <span :class="{ 'is-diff': row.isNew && currentRelationType.conflictProperties.includes('dest_des')
                }">{{row.dest_des}}</span>
              </template>
            </bk-table-column>
          </bk-table>

          <div class="relation-model">
            <div class="relation-model-header">{{t('关联的模型')}}</div>
            <bk-table :data="currentRelationType.relationModels" :show-header="false">
              <bk-table-column>
                <template #default="{ row }">
                  <model-summary
                    :data="row">
                  </model-summary>
                </template>
              </bk-table-column>
            </bk-table>
          </div>
        </template>
        <bk-exception
          v-else
          type="empty" scene="part">{{t('当前导入内容没有关联类型')}}
        </bk-exception>
      </div>
      <model-relation
        v-show="activeViewTabName === MODEL_TAB_NAME"
        :excluded-relations.sync="excludedModelRelations"
        :relation-selectable="relation =>
          (
            isModelExistedInImportData(relation.bk_obj_id)
            || isModelExistedInTargetEnv({
              bk_obj_id: relation.bk_asst_obj_id,
              bk_obj_name: relation.bk_asst_obj_name
            })
          )
          && includedModels[currentModel.bk_obj_id]
        "
        :model="currentModel">
        <template #relation-target="{ relation }">
          <div
            class="model-relation-target"
            :class="{
              'is-excluded-model': !isModelExistedInImportData(relation.bk_asst_obj_id)
                && !isModelExistedInTargetEnv({
                  bk_obj_id: relation.bk_asst_obj_id,
                  bk_obj_name: relation.bk_asst_obj_name
                })
            }"
          >
            <model-summary
              :data="{
                bk_obj_id: relation.bk_asst_obj_id,
                bk_obj_name: relation.bk_asst_obj_name,
                bk_obj_icon: relation.bk_asst_obj_icon
              }">
            </model-summary>
            <i
              class="icon-cc-remind"
              v-if="!isModelExistedInImportData(relation.bk_asst_obj_id)
                && !isModelExistedInTargetEnv({
                  bk_obj_id: relation.bk_asst_obj_id,
                  bk_obj_name: relation.bk_asst_obj_name
                })"
              v-bk-tooltips="{
                content: t('模型不存在，无法导入')
              }">
            </i>
          </div>
        </template>
      </model-relation>
    </template>
  </content-confirm-layout>
</template>

<script>
  import { defineComponent, reactive, toRef, ref, computed, watch } from 'vue'
  import { t } from '@/i18n'
  import { UNCATEGORIZED_GROUP_ID } from '@/dictionary/model-constants.js'
  import cloneDeep from 'lodash/cloneDeep'
  import ModelSummary from '../model-summary.vue'
  import store from '@/store'
  import DropdownSelect from '../dropdown-select.vue'
  import ContentConfirmLayout from '../content-confirm-layout.vue'
  import ModelGroupList from '../model-group-list.vue'
  import ModelRelation from '../model-relation.vue'
  import { useAssociations } from '../../../use-associations.js'

  export default defineComponent({
    name: 'ModelImportEditor',
    components: { ContentConfirmLayout, ModelSummary, DropdownSelect, ModelGroupList, ModelRelation },
    model: {
      prop: 'value',
      event: 'value-change'
    },
    props: {
      data: {
        type: Object,
        required: true,
        default: () => ({
          import_object: [],
          import_asst: []
        })
      },
      /**
       * 确认好的导入模型数据，支持 v-model
       * @property {Object} value
       * @property {Object} value.confirmedModels 确认好的模型列表数据
       * @property {Object} value.confirmedRelationTypes 确认好的关联关系类型列表数据
       */
      value: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props, { emit }) {
      const { associations: allAssociations } = useAssociations()// 所有的模型关联类型
      const allModelGroups = computed(() => store.getters['objectModelClassify/classifications']) // 所有的模型分组

      const excludedModelRelations = ref({}) // 排除的模型的关联关系
      const excludedRelationTypes = reactive({}) // 排除的关联关系类型
      const includedModels = reactive({}) // 手动选择的包含的模型

      const importRelationTypes = ref([]) // 导入的关联关系类型

      const JOIN_UNCATEGORIZED_GROUP = 'join-uncategorized-group'
      const CREATE_GROUP = 'create-group'
      const importGroupCreationTypes = reactive({}) // 导入的分组的创建方式，新建分组或加入未分类分组

      const COVER = 'cover' // 覆盖关联关系类型导入
      const SKIP = 'skip' // 跳过关联关系类型导入

      const importData = toRef(props, 'data') // 导入的模型数据
      const modelGroups = ref([]) // 模型分组

      const selectedModelId = ref('') // 选中的模型 ID

      const MODEL_TAB_NAME = 'model' // 模型 TAB 名称
      const RELATION_TAB_NAME = 'relation' // 模型关联关系 TAB 名称
      const activeViewTabName = ref(MODEL_TAB_NAME) // 激活的 TAB 名称
      const viewTabs = ref([{
        name: MODEL_TAB_NAME,
        label: t('模型'),
        count: () => importData.value?.import_object?.length || 0
      }, {
        name: RELATION_TAB_NAME,
        label: t('关联类型'),
        count: () => importRelationTypes.value.length || 0
      }])

      // 当前查看的模型
      const currentModel = reactive({
        bk_obj_name: '',
        bk_obj_id: '',
        object_asst: []
      })

      // 当前查看的关联关系类型
      const currentRelationType = reactive({
        bk_asst_id: '',
        bk_asst_name: '',
        relationTypes: [], // 关联类型对比
        relationModels: [], // 关联的模型
        isConflict: false
      })

      const handleTabChange = (viewTab) => {
        activeViewTabName.value = viewTab.name
        emit('view-tab-change', viewTab)
      }

      const handleModelSwitchChange = (model) => {
        delete excludedModelRelations.value[model.bk_obj_id]
      }

      /**
       * 根据后端返回的模型导入数据生成适用于 UI 展示的模型分组数据
       * @param {Array} importModels 模型数据
       * @return {undefined}
       */
      const genModelGroups = (importModels) => {
        modelGroups.value = []

        importModels.forEach((model) => {
          let newGroup = modelGroups.value
            .find(group => group.bk_classification_id === model.bk_classification_id)

          if (newGroup) {
            newGroup.bk_objects.push(model)
          } else {
            newGroup = {
              bk_classification_id: model.bk_classification_id,
              bk_classification_name: model.bk_classification_name,
              bk_objects: [model]
            }

            if (!isGroupExistedInTargetEnv(newGroup)) {
              toRef(importGroupCreationTypes, newGroup.bk_classification_id).value = JOIN_UNCATEGORIZED_GROUP
            }

            modelGroups.value.push(newGroup)
          }

          toRef(includedModels, model.bk_obj_id).value = !isModelExistedInTargetEnv(model)
        })
      }

      /**
       *  设置当前关联关系类型
       * @param {Object} relationType 关联关系类型
       */
      const setCurrentRelationType = (relationType) => {
        currentRelationType.bk_asst_id = relationType.bk_asst_id
        currentRelationType.bk_asst_name = relationType.bk_asst_name
        currentRelationType.isConflict = isRelationTypeConflict(relationType)

        const oldRelation = allAssociations.value?.find(r => r.bk_asst_id === relationType.bk_asst_id) || {}
        const conflictProperties = Object.keys(oldRelation).filter(key => oldRelation[key] !== relationType[key])

        currentRelationType.conflictProperties = conflictProperties

        if (isRelationTypeConflict(relationType)) {
          currentRelationType.relationDiffs = [
            { ...oldRelation, isOld: true },
            { ...relationType, isNew: true },
          ]
        } else {
          currentRelationType.relationDiffs = [relationType]
        }

        currentRelationType.relationModels = importData.value.import_object
          .filter(model => model.object_asst?.some(asst => asst.bk_asst_id === relationType.bk_asst_id))
      }

      // 设置当前模型数据
      const setCurrentModel = async (model) => {
        currentModel.bk_obj_name = model.bk_obj_name
        currentModel.bk_obj_id = model.bk_obj_id
        currentModel.object_asst = model.object_asst
      }

      /**
       * 检查导入的模型分组是否与已存在于目标环境
       * @param {Object} importGroup 导入的模型分组
       * @return {Boolean} 是否存在
       */
      const isGroupExistedInTargetEnv = importGroup => allModelGroups.value
        .some(group => importGroup.bk_classification_id === group.bk_classification_id
          || importGroup.bk_classification_name === group.bk_classification_name)

      /**
       * 检查导入的模型是否已存在于目标环境
       * @param {Object} importModel 导入的模型
       * @return {Boolean} 是否冲突
       */
      const isModelExistedInTargetEnv = importModel => allModelGroups.value.some(group => group.bk_objects
        .some(model => model.bk_obj_id === importModel.bk_obj_id
          || model.bk_obj_name === importModel.bk_obj_name))

      /**
       * 检查模型是否存在于导入的数据中
       * @param {String} modelId 模型的 bk_obj_id
       * @return {Boolean} 是否存在
       */
      const isModelExistedInImportData = modelId => modelGroups.value
        .some(modelGroup => modelGroup.bk_objects.some(model => model.bk_obj_id === modelId))

      /**
       * 检查导入的模型关系类型是否与已有模型关系类型冲突
       * @param {Object} importRelationType 导入的模型类型
       * @return {Boolean} 是否冲突
       */
      const isRelationTypeConflict = (importRelationType) => {
        const isConflict = allAssociations.value.some((relationType) => {
          const isIdRepeat = relationType.bk_asst_id === importRelationType.bk_asst_id
          return isIdRepeat
        })

        return isConflict
      }

      // 默认选中第一个关联关系
      const initRelationTypeSelection = () => {
        const firstRelation = importRelationTypes.value[0]

        setCurrentRelationType(firstRelation)
      }

      // 默认选中第一个模型
      const initModelSelection = () => {
        const firstModel = modelGroups.value[0].bk_objects[0]

        selectedModelId.value = firstModel.bk_obj_id

        setCurrentModel(firstModel)
      }

      // 将导入的数据进行处理后传递出去
      const emitValue = () => {
        const importData = cloneDeep(props.data)

        // 过滤掉没有选择的模型的关联关系
        const confirmedModels = importData.import_object?.filter((model) => {
          const excluded = excludedModelRelations.value?.[model.bk_obj_id]

          if (model.object_asst && excluded) {
            model.object_asst = model.object_asst.filter(asst => !excluded
              .some(excludedAsst => excludedAsst.bk_obj_asst_id === asst.bk_obj_asst_id))
          }

          return includedModels?.[model.bk_obj_id]
        })

        // 处理模型分组创建方式，选择了未分组的则将模型的分组 ID 改为未分类模型分组 ID
        confirmedModels?.forEach((model) => {
          if (importGroupCreationTypes[model.bk_classification_id] === JOIN_UNCATEGORIZED_GROUP) {
            model.bk_classification_id = UNCATEGORIZED_GROUP_ID
          } else {
            model.bk_classification_id = importData.import_object
              .find(importModel => importModel.bk_obj_id === model.bk_obj_id)?.bk_classification_id
          }
        })

        // 过滤掉冲突中选择跳过的关联关系类型
        const confirmedRelationTypes = importData.import_asst
          ?.filter(asst => excludedRelationTypes?.[asst.bk_asst_id] !== SKIP
            && importRelationTypes.value.some(relationType => relationType.bk_asst_id === asst.bk_asst_id))

        emit('value-change', {
          confirmedModels,
          confirmedRelationTypes
        })
      }

      /**
       * 对导入数据进行初始化渲染 UI
       * @param {Object} newData 新导入的数据
       */
      const initValue = (newData) => {
        const clonedNewData = cloneDeep(newData)

        if (clonedNewData.import_object?.length > 0) {
          genModelGroups(clonedNewData.import_object)
          initModelSelection()
        }

        if (clonedNewData.import_asst?.length > 0) {
          importRelationTypes.value = clonedNewData.import_asst
          clonedNewData.import_asst.forEach((asst) => {
            if (isRelationTypeConflict(asst)) {
              toRef(excludedRelationTypes, asst.bk_asst_id).value = SKIP
            }
          })
          initRelationTypeSelection()
        }
      }

      const genImportRelationTypes = () => {
        if (!importData.value.import_object?.length) return

        const effectedRelations = []

        importData.value.import_object.forEach((model) => {
          if (includedModels[model.bk_obj_id] && model.object_asst) {
            model.object_asst.forEach((asst) => {
              if (!excludedModelRelations.value?.[model.bk_obj_id]
                ?.some(excludeItem => excludeItem.bk_asst_id === asst.bk_asst_id)
                && !effectedRelations.some(relation => relation.bk_asst_id === asst.bk_asst_id)) {
                effectedRelations.push(asst)
              }
            })
          }
        })

        importRelationTypes.value = cloneDeep(importData.value.import_asst)?.filter(asst => effectedRelations
          .some(relation => relation.bk_asst_id === asst.bk_asst_id)) || []

        if (importRelationTypes.value[0]) setCurrentRelationType(importRelationTypes.value[0])
      }

      // 监听各类型的数据变化，更新确认后的导入数据
      watch(includedModels, () => {
        genImportRelationTypes()
        emitValue()
      }, { immediate: true, deep: true })
      watch(excludedRelationTypes, () => {
        emitValue()
      }, { immediate: true, deep: true })
      watch(excludedModelRelations, () => {
        genImportRelationTypes()
        emitValue()
      }, { immediate: true, deep: true })
      watch(importGroupCreationTypes, emitValue, { immediate: true, deep: true })

      // 监听导入的数据初始化数据渲染 UI
      watch(importData, initValue, { immediate: true })

      return {
        MODEL_TAB_NAME,
        RELATION_TAB_NAME,
        JOIN_UNCATEGORIZED_GROUP,
        CREATE_GROUP,
        COVER,
        SKIP,
        importRelationTypes,
        excludedRelationTypes,
        includedModels,
        viewTabs,
        importGroupCreationTypes,
        activeViewTabName,
        handleTabChange,
        isGroupExistedInTargetEnv,
        isModelExistedInTargetEnv,
        isRelationTypeConflict,
        selectedModelId,
        excludedModelRelations,
        modelGroups,
        currentModel,
        currentRelationType,
        setCurrentRelationType,
        t,
        setCurrentModel,
        isModelExistedInImportData,
        handleModelSwitchChange
      }
    }
  })
</script>

<style lang="scss" scoped>
$viewTabHeight:40px;

.icon-cc-remind {
  margin-left: 5px;
  color: #ff9c01;
  font-size: 14px;
}

.view-tab {
  display: flex;
  height: $viewTabHeight;
  background: #fafbfd;
  border-right: 1px solid $borderColor;

  &-item {
    font-size: 12px;
    flex: 1 1 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border-bottom: 1px solid $borderColor;

    & + & {
      border-left: 1px solid $borderColor;
    }

    &.is-active {
      font-weight: 700;
      color: #3a84ff;
      background-color: #fff;
      cursor: default;
      border-bottom-color: #fff;

      .item-count {
        color: inherit;
        font-weight: normal;
        background-color: #e1ecff;
      }
    }

    .item-count {
      display: block;
      line-height: 16px;
      color: #979ba5;
      background: #f0f1f5;
      border-radius: 8px;
      padding: 0 8px;
      margin-left: 4px;
    }
  }
}

.model-group {
  &-list {
    overflow-y: auto;
    height: calc(100% - $viewTabHeight);
    font-size: 12px;
  }

  &-dropdown {
    margin-left: auto;
    margin-right: 7px;
    line-height: 22px;
  }

  &-conflict-icon {
    margin-left: 5px;
  }

  &-item {
    border-bottom: 1px solid $borderColor;
  }
}

.model {
  &-item {
    position: relative;
    display: flex;
    align-items: center;
    height: 40px;
    cursor: default;
    transition: background-color 200ms ease;

    .bk-switcher {
      margin-left: auto;
      margin-right: 16px;
    }

    &::before {
      content: "";
      display: block;
      position: absolute;
      top: 0;
      right: 0;
      left: 30px;
      height: 1px;
      background-color: $borderColor;
    }

    .model-summary {
      margin-left: 45px;
    }
  }
}

.model-relation {
  &-target {
    display: flex;
    align-items: center;

    &.is-excluded-model {
      .model-summary{
        opacity: 0.5;
      }
    }
  }

  &-name {
    display: inline-block;
    padding: 0 5px;
    background: #f0f1f5;
    border-radius: 2px;
    line-height: 18px;
  }
}

.relation-type {
  &-item {
    display: flex;
    align-items: center;
    padding-left: 16px;
    height: 40px;
    border-bottom: 1px solid $borderColor;
    border-right: 1px solid $borderColor;
    font-size: 12px;

    &:hover,
    &.is-active{
      background-color: #fff;
    }

    &.is-active{
      border-right: none;
    }
  }

  &-id {
    color: #979BA5;
  }
}

.relation-model {
  margin-top: 20px;
  &-header{
    height: 60px;
    display: flex;
    font-size: 16px;
    align-items: center;
    color: #000;
  }
}

.resolve-method-select {
  width: 270px;
  margin-bottom: 16px;
}

.relation-type-diff {
 .is-diff {
    color: #EA3636;
  }

  &-header {
    height: 60px;
    display: flex;
    align-items: center;
    font-size: 16px;
    color: #000;
  }
}
</style>
