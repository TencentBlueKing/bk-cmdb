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
  import { computed, ref, watch } from 'vue'
  import { useStore } from '@/store'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import ModelSelector from '@/components/model-instance/model-selector.vue'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'
  import useProperty from '@/hooks/model/property'
  import { excludeFieldType } from './use-field'

  const props = defineProps({
    templateFieldList: {
      type: Array,
      default: () => []
    }
  })
  const emit = defineEmits(['confirm', 'cancel'])

  const store = useStore()

  const excludeModelIds = computed(() => store.getters['objectMainLineModule/mainLineModels']
    .filter(model => model.bk_obj_id !== BUILTIN_MODELS.HOST)
    .map(model => model.bk_obj_id))

  const selectedModelId = ref('')
  const selected = ref([])

  watch(selectedModelId, () => {
    handleClearSelect()
  })

  const searchParams = computed(() => ({
    bk_obj_id: selectedModelId.value
  }))
  const [{ properties, pending }] = useProperty(searchParams)

  const fieldList = computed(() => properties.value
    .filter(item => !excludeFieldType.includes(item.bk_property_type) && !item.ispre))
  const selectedStatus = computed(() => {
    const status = {
      selected: {},
      disabled: {}
    }
    fieldList.value.forEach((prop) => {
      status.selected[prop.id] = selected.value.some(item => item.id === prop.id)
      status.disabled[prop.id] = props.templateFieldList.some(item => item.bk_property_id === prop.bk_property_id
        || item.bk_property_name === prop.bk_property_name)
    })
    return status
  })

  const handleSelect = (field) => {
    if (selectedStatus.value.disabled[field.id]) {
      return
    }

    if (selectedStatus.value.selected[field.id]) {
      const index = selected.value.findIndex(item => item.id === field.id)
      if (~index) {
        selected.value.splice(index, 1)
      }
    } else {
      selected.value.push(field)
    }
  }

  const handleSelectAll = () => {
    selected.value = fieldList.value.filter(field => !selectedStatus.value.disabled[field.id])
  }
  const handleClearSelect = () => {
    selected.value = []
  }

  const handleConfirm = () => {
    emit('confirm', selected.value)
  }
  const handleCancel = () => {
    emit('cancel')
  }
  defineExpose({
    selectedModelId
  })
</script>

<template>
  <cmdb-sticky-layout class="model-field-selector">
    <grid-layout class="main-content" mode="form" :gap="12" :font-size="'14px'" :max-columns="1">
      <grid-item required :label="$t('模型')">
        <model-selector
          class="model-selector"
          searchable
          name="refModel"
          :exclude="excludeModelIds"
          :placeholder="$t('请选择xx', { name: $t('模型') })"
          v-model="selectedModelId">
        </model-selector>
      </grid-item>
      <grid-item v-if="selectedModelId">
        <div class="select-toolbar">
          <div class="select-stat">
            {{$t('请选择导入的字段：')}}
            <i18n path="已选个数" v-show="selectedModelId">
              <template #count><span class="count">{{selected.length}}</span></template>
            </i18n>
          </div>
          <div class="select-action">
            <bk-link theme="primary" @click="handleSelectAll">{{ $t('全选') }}</bk-link>
            <bk-link theme="primary" @click="handleClearSelect">{{ $t('清空') }}</bk-link>
          </div>
        </div>
        <div :class="['field-list', { empty: !fieldList.length }]"
          v-bkloading="{ isLoading: pending }">
          <field-card
            v-for="(field, index) in fieldList"
            :key="index"
            :field="field"
            :sortable="false"
            :deletable="false"
            v-bk-tooltips="{
              disabled: !selectedStatus.disabled[field.id],
              content: $t('字段已在模板中存在，无法添加')
            }"
            @click-field="handleSelect">
            <template #action-append>
              <bk-checkbox
                @change="handleSelect(field)"
                :disabled="selectedStatus.disabled[field.id]"
                :value="selectedStatus.selected[field.id]">
              </bk-checkbox>
            </template>
          </field-card>
          <div class="field-list-empty" v-if="!fieldList.length">
            <div class="tips">
              <bk-icon type="info" />{{$t('无可用字段')}}
            </div>
          </div>
        </div>
      </grid-item>
    </grid-layout>
    <template slot="footer" slot-scope="{ sticky }">
      <div class="btn-group" :class="{ 'is-sticky': sticky }">
        <bk-button theme="primary" @click="handleConfirm">
          {{ $t('确定') }}
        </bk-button>
        <bk-button theme="default" @click="handleCancel">
          {{$t('取消')}}
        </bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
  .model-field-selector {
    height: 100%;
    padding: 0;
    @include scrollbar-y;
    .main-content {
      max-height: calc(100% - 52px);
      @include scrollbar-y;
      padding: 20px 24px;
    }

    .model-selector {
      width: 100%;
    }

    .field-list {
      display: grid;
      gap: 16px;
      grid-template-columns: repeat(2, 1fr);
      width: 100%;
      align-content: flex-start;
      padding: 16px;
      background: #F5F7FA;

      &.empty {
        display: block;
      }
    }

    .select-toolbar {
      display: flex;
      justify-content: space-between;
      margin-bottom: 6px;

      .select-stat {
        font-size: 14px;
        .count {
          color: $primaryColor;
          font-weight: bold;
          padding: 0 .2em;
        }
      }
      .select-action {
        display: flex;
        align-items: center;
        gap: 12px;
        :deep(.bk-link .bk-link-text) {
          font-size: 12px;
        }
      }
    }

    .unselected-model,
    .field-list-empty {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 60px;
      background: #F5F7FA;
      .tips {
        display: flex;
        align-items: center;
        gap: 3px;
        font-size: 12px;
      }
    }

    .btn-group {
      display: flex;
      gap: 8px;
      padding: 8px 24px;
      background: #fff;
      &.is-sticky {
        border-top: 1px solid #dcdee5;
      }
      .bk-button{
        width: 88px;
        height: 32px;
      }
    }
  }
</style>
