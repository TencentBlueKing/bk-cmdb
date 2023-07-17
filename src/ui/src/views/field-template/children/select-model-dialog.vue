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
  import { computed, reactive, ref, watchEffect } from 'vue'
  import { useStore } from '@/store'
  import { BUILTIN_MODELS, UNCATEGORIZED_GROUP_ID } from '@/dictionary/model-constants'

  const props = defineProps({
    selected: {
      type: Array,
      default: () => ([])
    },
    binded: {
      type: Array,
      default: () => ([])
    }
  })

  const emit = defineEmits(['confirm'])

  const isShow = ref(false)

  const store = useStore()

  const classifications = ref([])
  const selectedLocal = ref([])
  const filterWord = ref('')

  const dataEmpty = reactive({
    type: 'search'
  })

  watchEffect(async () => {
    classifications.value = await store.dispatch('objectModelClassify/searchClassificationsObjects', {
      fromCache: true
    })
  })
  watchEffect(() => {
    if (isShow.value) {
      selectedLocal.value = props.selected.slice()
    }
  })

  const dialogHeight = computed(() => `${Math.floor(Math.max(store.state.appHeight * 0.8, 200) - 110)}px`)
  const dialogPos = computed(() => ({
    top: `${Math.floor(Math.max(store.state.appHeight * 0.4 - (parseInt(dialogHeight.value, 10) / 2), 20))}`
  }))
  const dialogWidth = computed(() => (window.innerWidth > 1920 ? 1532 : 1132))

  const excludeModelIds = computed(() => store.getters['objectMainLineModule/mainLineModels']
    .filter(model => model.bk_obj_id !== BUILTIN_MODELS.HOST)
    .map(model => model.bk_obj_id))

  const modelGroupList = computed(() => {
    const list = []
    classifications.value.forEach((classification) => {
      list.push({
        ...classification,
        bk_objects: classification.bk_objects
          .filter(model => !model.bk_ispaused && !model.bk_ishidden && !excludeModelIds.value.includes(model.bk_obj_id))
      })
    })
    return list
      .filter(item => item.bk_objects.length > 0)
      .sort((a, b) => (b.bk_classification_id === UNCATEGORIZED_GROUP_ID ? -1 : 0))
  })

  const modelList = computed(() => modelGroupList.value.reduce((acc, cur) => acc.concat(cur.bk_objects), []))

  const displayModelGroupList = computed(() => {
    if (filterWord.value) {
      const reg = new RegExp(filterWord.value, 'i')
      const list = []
      modelGroupList.value.forEach((group) => {
        list.push({
          ...group,
          bk_objects: group.bk_objects.filter(model => reg.test(model.bk_obj_name))
        })
      })
      return list.filter(item => item.bk_objects.length > 0)
    }
    return modelGroupList.value
  })

  const selectedStatus = computed(() => {
    const status = {
      selected: {},
      binded: {}
    }
    modelList.value.forEach((model) => {
      status.selected[model.id] = selectedLocal.value.some(item => item.id === model.id)
      status.binded[model.id] = props.binded.some(item => item.id === model.id)
    })
    return status
  })

  const handleSelect = (model) => {
    if (selectedStatus.value.binded[model.id]) {
      return
    }
    if (selectedStatus.value.selected[model.id]) {
      const index = selectedLocal.value.findIndex(item => item.id === model.id)
      if (~index) {
        selectedLocal.value.splice(index, 1)
      }
    } else {
      selectedLocal.value.push(model)
    }
  }

  const handleSelectAll = (checked) => {
    if (checked) {
      modelList.value.forEach((model) => {
        if (!selectedStatus.value.binded[model.id]
          && !selectedStatus.value.selected[model.id]) {
          selectedLocal.value.push(model)
        }
      })
    } else {
      modelList.value.forEach((model) => {
        if (!selectedStatus.value.binded[model.id]
          && selectedStatus.value.selected[model.id]) {
          const index = selectedLocal.value.findIndex(item => item.id === model.id)
          if (~index) {
            selectedLocal.value.splice(index, 1)
          }
        }
      })
    }
  }


  const handleConfirm = () => {
    emit('confirm', selectedLocal.value)
    close()
  }
  const handleCancel = () => {
    close()
  }

  const handleClearFilter = () => {
    filterWord.value = ''
  }

  const show = () => {
    selectedLocal.value = []
    isShow.value = true
  }
  const close = () => isShow.value = false

  defineExpose({
    show
  })
</script>

<template>
  <bk-dialog
    :render-directive="'if'"
    v-model="isShow"
    :width="dialogWidth"
    :position="dialogPos"
    :close-icon="false"
    :draggable="false"
    ext-cls="custom-wrapper">
    <template #tools>
      <div class="dialog-top">
        <div class="title">{{$t('添加模型')}}</div>
        <bk-input
          class="search-input"
          v-model="filterWord"
          :placeholder="$t('请输入模型名称')"
          :clearable="true"
          :right-icon="'bk-icon icon-search'" />
      </div>
    </template>
    <div :class="['dialog-content', { empty: !displayModelGroupList.length }]">
      <cmdb-collapse
        v-for="(group, index) in displayModelGroupList"
        class="model-group"
        :key="index"
        :label="`${group.bk_classification_name}（${group.bk_objects.length}）`"
        arrow-type="filled">
        <div class="model-list">
          <template v-for="(model, modelIndex) in group.bk_objects">
            <cmdb-auth :key="modelIndex" tag="div" :auth="{ type: $OPERATION.U_MODEL, relation: [model.id] }">
              <template #default="{ disabled }">
                <div :class="['model-item', { 'is-builtin': model.ispre }]"
                  @click="handleSelect(model)">
                  <div class="model-icon">
                    <i class="icon" :class="model.bk_obj_icon"></i>
                  </div>
                  <div class="model-details">
                    <div class="model-name" :title="model.bk_obj_name">
                      {{ model["bk_obj_name"] }}
                    </div>
                    <div class="model-id" :title="model.bk_obj_id">
                      {{ model["bk_obj_id"] }}
                    </div>
                  </div>
                  <bk-checkbox class="model-checkbox"
                    v-bk-tooltips="{
                      disabled: !selectedStatus.binded[model.id],
                      content: $t('模型已绑定')
                    }"
                    :value="selectedStatus.selected[model.id]
                      || selectedStatus.binded[model.id]"
                    :disabled="disabled || selectedStatus.binded[model.id]">
                  </bk-checkbox>
                </div>
              </template>
            </cmdb-auth>
          </template>
        </div>
      </cmdb-collapse>
      <cmdb-data-empty
        v-show="!displayModelGroupList.length"
        :stuff="dataEmpty"
        @clear="handleClearFilter">
      </cmdb-data-empty>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <bk-checkbox @change="handleSelectAll" :value="selectedLocal.length === modelList.length" class="checkbox-all">
          <i18n path="全选N个模型">
            <template #count><em class="count">{{modelList.length}}</em></template>
          </i18n>
        </bk-checkbox>
        <div class="operation">
          <bk-button theme="primary"
            @click="handleConfirm">{{ $t('确定') }}</bk-button>
          <bk-button @click="handleCancel">{{ $t('取消') }}</bk-button>
        </div>
      </div>
    </template>
  </bk-dialog>
</template>

<style lang="scss" scoped>
  :deep(.custom-wrapper) {
    .bk-dialog-content {
      background: #F5F7FB;
    }
  }

  .dialog-top {
    display: flex;
    justify-content: space-between;
    padding: 20px 24px 0 24px;
    .title {
      color: #313238;
      font-size: 16px;
    }
    .search-input {
      width: 392px;
    }
  }

  .dialog-content {
    height: v-bind(dialogHeight);
    padding: 0 12px;
    @include scrollbar-y;

    &.empty {
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .operation {
      display: flex;
      gap: 8px;
    }

    .checkbox-all {
      .count {
        font-style: normal;
        padding: 0 .2em;
      }
    }
  }

  .model-group {
    margin-bottom: 24px;
    :deep(.collapse-trigger) {
      font-weight: 400;
    }
  }

  .model-list {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
    margin: 12px 0 0 0;
  }

  .model-item {
    display: flex;
    align-items: center;
    height: 60px;
    background-color: #fff;
    border-radius: 2px;
    box-shadow: 0px 2px 4px 0px rgba(25, 25, 41, 0.05);
    padding: 0 16px;
    cursor: pointer;

    &:hover {
      transition: all 200ms ease;
      box-shadow: 0px 2px 4px 0px rgba(25, 25, 41, 0.05),
        0px 2px 4px 0px rgba(0, 0, 0, 0.1);
    }

    &.is-builtin {
      .model-icon {
        background-color: #f5f7fa;
        transition: background-color 200ms ease;

        .icon {
          color: #798aad;
        }
      }

      &:hover .model-icon{
        background-color: #fff;
      }
    }

    .model-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background-color: #e1ecff;
      .icon {
        color: #3a84ff;
        font-size: 16px;
      }
    }

    .model-details {
      flex: 1;
      margin-left: 12px;
      margin-right: 4px;
      width: 0;
    }

    .model-name {
      line-height: 19px;
      font-size: 14px;
      @include ellipsis;
    }

    .model-id {
      line-height: 16px;
      font-size: 12px;
      color: #bfc7d2;
      @include ellipsis;
    }
    .model-checkbox {
      flex: none;
      width: 20px;
    }
  }

  @media screen and (min-width: 1920px) {
    .model-list {
      grid-template-columns: repeat(5, 1fr);
    }
  }
</style>
