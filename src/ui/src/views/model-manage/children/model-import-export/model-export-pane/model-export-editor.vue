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
      <model-group-list
        :selected-model-id.sync="selectedModelId"
        @model-select="setCurrentModel"
        :groups="modelGroups">
        <template #model-append="{ modelGroup, model, modelIndex, groupIndex }">
          <bk-button
            v-bk-tooltips="{
              content: t('删除后，不导出'),
              delay: [1000, 0]
            }"
            class="model-item-delete-button"
            @click.stop="removeModel(modelGroup, model, modelIndex, groupIndex)"
            :text="true"
            icon="delete"
          ></bk-button>
        </template>
      </model-group-list>
    </template>
    <template #main>
      <model-relation
        :remote="true"
        :excluded-relations.sync="excludedRelations"
        :model="currentModel">
        <template #relation-target="{ relation }">
          <div
            class="model-relation-target"
            :class="{
              'is-excluded-model': !isModelInExportData(
                relation.targetModel.bk_obj_id
              )
            }"
          >
            <model-summary
              :data="relation.targetModel">
            </model-summary>
            <template v-if="!relation.targetModel.ispre">
              <i
                v-if="!isModelInExportData(relation.targetModel.bk_obj_id)"
                class="icon-cc-remind"
                v-bk-tooltips="{
                  content: t('不在本次导出的模型中')
                }"
              />
              <div
                class="add-button"
                @click="joinToExportData(relation.targetModel)"
                v-bk-tooltips="{
                  content: t('加选模型'),
                  delay: [1000, 0]
                }">
                <bk-icon
                  type="cc-plus"
                ></bk-icon>
              </div>
            </template>
            <router-link
              v-bk-tooltips="{
                content: t('查看模型'),
                delay: [1000, 0]
              }"
              class="jump-button"
              target="_blank"
              :to="genModelDetailRoute(relation.targetModel.bk_obj_id)">
              <bk-icon type="cc-jump-link"></bk-icon>
            </router-link>
          </div>
        </template>
      </model-relation>
    </template>
  </content-confirm-layout>
</template>
<script>
  import { MENU_MODEL_DETAILS } from '@/dictionary/menu-symbol'
  import { defineComponent, reactive, onMounted, ref, computed, watch } from '@vue/composition-api'
  import { t } from '@/i18n'
  import cloneDeep from 'lodash/cloneDeep'
  import ModelSummary from '../model-summary.vue'
  import store from '@/store'
  import ContentConfirmLayout from '../content-confirm-layout.vue'
  import ModelGroupList from '../model-group-list.vue'
  import ModelRelation from '../model-relation.vue'

  export default defineComponent({
    name: 'ModelExportEditor',
    components: { ModelSummary, ContentConfirmLayout, ModelGroupList, ModelRelation },
    model: {
      prop: 'value',
      event: 'value-change'
    },
    props: {
      // 需要导出的数据
      data: {
        type: Array,
        required: true,
        default: () => []
      },
      // 确认后的导出数据
      value: {
        type: Object,
        default: () => []
      }
    },
    setup({ data }, { emit }) {
      const modelGroups = ref(cloneDeep(data))
      const excludedRelations = ref({}) // 排除的关联关系
      const allModelGroups = computed(() => store.getters['objectModelClassify/classifications']) // 全量模型数据
      const selectedModelId = ref('') // 当前选中的模型的 bk_obj_id，用于手动控制选中的模型
      const currentModel = reactive({
        bk_obj_name: '',
        bk_obj_id: ''
      })

      /**
       * 生成模型详情页面的路由
       * @param {String} modelId 模型的 bk_obj_id
       */
      const genModelDetailRoute = modelId => ({
        name: MENU_MODEL_DETAILS,
        params: {
          modelId
        }
      })

      /**
       * 设置当前模型的数据，用于渲染模型相关的详情数据
       * @param {Object} model 模型数据
       */
      const setCurrentModel = async (model) => {
        currentModel.bk_obj_name = model.bk_obj_name
        currentModel.bk_obj_id = model.bk_obj_id
      }

      /**
       * 把不在导出数据中的模型加入到导出数据里来
       * @param {Object} model 模型数据
       */
      const joinToExportData = (model) => {
        let modelGroup = null

        modelGroup = modelGroups.value.find(group => group.bk_classification_id === model.bk_classification_id)

        if (modelGroup && !modelGroup.bk_objects.some(item => item.bk_obj_id === model.bk_obj_id)) {
          modelGroup.bk_objects.push(model)
        }

        if (!modelGroup) {
          modelGroup = cloneDeep(allModelGroups.value
            .find(group => group.bk_classification_id === model.bk_classification_id))
          modelGroup.bk_objects = [model]
          modelGroups.value.push(modelGroup)
        }
      }

      /**
       * 判断模型是否在导出数据中
       * @param {String} modelId 模型的 bk_obj_id
       */
      const isModelInExportData = modelId => modelGroups.value
        .some(modelGroup => modelGroup.bk_objects.some(model => model.bk_obj_id === modelId))

      /**
       * 删除导出列表中的模型
       * @param {Object} modelGroup 模型分组数据
       * @param {Object} model 模型数据
       * @param {Number} modelIndex 模型索引
       * @param {Number} groupIndex 模型分组索引
       */
      const removeModel = (modelGroup, model, modelIndex, groupIndex) => {
        let models = modelGroup.bk_objects
        let nextModel = null
        const modelsLen = modelGroups.value.reduce((acc, group) => acc + group.bk_objects.length, 0)

        if (modelsLen === 1) {
          return
        }

        models.splice(modelIndex, 1)

        if (model.bk_obj_id === currentModel.bk_obj_id) {
          if (models.length === 0) {
            modelGroups.value.splice(groupIndex, 1)
            models = modelGroups.value[Math.max(groupIndex - 1, 0)].bk_objects
          }
          nextModel = models[Math.min(models.length - 1, modelIndex)]
        }

        if (nextModel) {
          setCurrentModel(nextModel)
          selectedModelId.value = nextModel.bk_obj_id
        }
      }

      /**
       * 当模型分组数据和排除的关联关系有变更时，把数据更新到 value 中
       */
      watch(
        [modelGroups, excludedRelations], () => {
          emit('value-change', {
            modelGroups: modelGroups.value,
            excludedRelations
          })
        },
        {
          immediate: true,
          deep: true
        }
      )

      onMounted(async () => {
        // 默认选择第一个模型分组的第一个模型
        const firstModel = modelGroups.value[0].bk_objects[0]

        selectedModelId.value = firstModel.bk_obj_id
        setCurrentModel(firstModel)
      })

      return {
        selectedModelId,
        modelGroups,
        currentModel,
        joinToExportData,
        t,
        genModelDetailRoute,
        setCurrentModel,
        removeModel,
        excludedRelations,
        isModelInExportData,
      }
    }
  })
</script>

<style lang="scss" scoped>
.icon-cc-remind {
  margin-left: 5px;
  color: #ff9c01;
  font-size: 14px;
}

.model-item-delete-button {
  visibility: hidden;
  margin-left: auto;
  margin-right: 10px;
  color: #63656e;

  .model-item:hover & {
    visibility: visible;
  }
}

/deep/ .model-relation-table {
  height: calc(100% - 90px);
}

.model-relation-target {
    display: flex;
    align-items: center;

    &.is-excluded-model {
      .model-summary{
        opacity: 0.5;
      }

      .add-button,
      .jump-button{
        display: block;
      }
    }

    .add-button,
    .jump-button {
      cursor: pointer;
      display: none;
      color: #1768ef;
      margin-left: 5px;

       .bk-icon {
        vertical-align: top;
      }
    }

    .add-button {
      height: 18px;
      font-size: 18px;
    }

    .jump-button {
      font-size: 14px;
      height: 14px;
    }
  }
</style>
