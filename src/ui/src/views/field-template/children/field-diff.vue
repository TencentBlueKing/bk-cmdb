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
  import { ref, computed } from 'vue'
  import { useStore } from '@/store'
  import useGroupProperty from '@/hooks/utils/group-property'
  import useGroup from '@/hooks/model/group'
  import useProperty from '@/hooks/model/property'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import DiffBrand from './diff-brand.vue'

  const props = defineProps({
    model: {
      type: Object,
      default: () => ({})
    },
    // 对比的结果数据
    diffs: {
      type: Object,
      default: () => ({})
    },
    // 模板字段列表
    templateFieldList: {
      type: Array,
      default: () => ([])
    }
  })
  const store = useStore()

  const isOnlyShowTemplate = ref(false)

  // 查询模型的字段列表
  const propertyParams = computed(() => ({
    bk_obj_id: props.model.bk_obj_id,
    bk_supplier_account: store.getters.supplierAccount
  }))
  const [{ properties, pending }] = useProperty(propertyParams)
  const [{ groups }] = useGroup(propertyParams)
  const groupedPropertyies = useGroupProperty(groups, properties)

  const counts = computed(() => ({
    new: props.diffs?.create?.length ?? 0,
    update: props.diffs?.update?.length ?? 0,
    conflict: props.diffs?.conflict?.length ?? 0,
    unbinded: 0,
    unchanged: props.diffs?.unchanged?.length ?? 0,
  }))

  const newFieldList = computed(() => {
    const news = props.diffs.create ?? []
    return news.map(item => props.templateFieldList.find(field => field.bk_property_id === item.bk_property_id))
  })

  const displayFieldGroups = computed(() => {
    const displayFieldGroups = []
    groupedPropertyies.value.forEach((item) => {
      const data = {
        group: item.group,
        properties: item.properties.slice()
      }

      // 由模板此次新建的字段添加至默认分组的头部
      if (data.group.bk_group_id === 'default') {
        data.properties.unshift(...newFieldList.value)
      }

      displayFieldGroups.push(data)
    })
    displayFieldGroups.forEach((group) => {
      group.properties = group.properties.filter(field => (isOnlyShowTemplate.value ? isTemplate(field) : true))
    })
    return displayFieldGroups
  })

  const isConflict = field => props.diffs.conflict?.some(item => item.data.bk_property_id === field.bk_property_id)

  const getFieldCardClassNames = (field) => {
    // 新增：模型中没有，展示的是模板的字段
    // 更新：共同的字段，但是模板中有更新
    // 冲突：无法应用到模型的字段，因模型中的字段与模板当前的设置冲突
    // 解除：模型的字段在模板中已经找不到
    if (props.diffs.create?.some(item => item.bk_property_id === field.bk_property_id)) {
      return 'new'
    }
    if (props.diffs.update?.some(item => item.bk_property_id === field.bk_property_id)) {
      return 'update'
    }

    // 冲突使用模型数据中的字段id匹配
    if (isConflict(field)) {
      return 'conflict'
    }

    if (props.diffs.unchanged?.some(item => item.bk_property_id === field.bk_property_id)) {
      return 'unchanged'
    }
    if (!props.templateFieldList?.some(item => item.bk_property_id === field.bk_property_id)) {
      return 'unbinded'
    }
  }

  const isTemplate = field => props.templateFieldList.some(item => item.bk_property_id === field.bk_property_id)
</script>

<template>
  <div class="field-diff" v-bkloading="{ isLoading: pending }">
    <div class="status-bar">
      <div class="diff-summary">
        <div class="summary-title">{{$t('模板应用后的差异对比：')}}</div>
        <div class="summray-content">
          <diff-brand :count="counts.new" :text="$t('新增字段')" status="new"></diff-brand>
          <diff-brand :count="counts.update" :text="$t('更新覆盖')" status="update"></diff-brand>
          <diff-brand :count="counts.conflict" :text="$t('字段冲突')" status="conflict"
            :tooltips="'#field-template-field-diff-conflict-tooltips'">
          </diff-brand>
          <diff-brand :count="counts.unbinded" :text="$t('解除纳管')" status="unbinded"
            :tooltips="$t('模板中删除了该字段，后续不再统一管理该字段')">
          </diff-brand>
          <diff-brand :count="counts.unchanged" :text="$t('无变化')" status="unchanged"></diff-brand>
          <span class="tips-content" id="field-template-field-diff-conflict-tooltips">
            <div>{{ $t('字段冲突的情况：') }}</div>
            <ul class="list-item">
              <li>模板字段与模型字段 ID 类型一样，但已经被其他模板绑定</li>
              <li>模板字段与模型字段的 ID 一样，但字段类型不一致</li>
              <li>模板设置的唯一性校验与模型设置的冲突</li>
            </ul>
          </span>
        </div>
      </div>
      <bk-checkbox class="filter-checkbox" v-model="isOnlyShowTemplate">{{ $t('仅显示与模板相关字段') }}</bk-checkbox>
    </div>
    <div class="model-group-container">
      <cmdb-collapse
        v-for="({ group, properties: fieldList }) in displayFieldGroups"
        class="model-group"
        :key="group.id"
        :label="group.bk_group_name"
        arrow-type="filled">
        <div class="field-list">
          <field-card
            v-for="(field, index) in fieldList"
            :class="getFieldCardClassNames(field)"
            :key="index"
            :field="field"
            :sortable="false"
            :deletable="false"
            :is-template="isTemplate(field)">
            <template #flag-append v-if="isConflict(field)">
              <i class="bk-icon icon-exclamation-circle-shape conflict-icon"></i>
            </template>
          </field-card>
        </div>
      </cmdb-collapse>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .field-diff {
    height: 100%;
  }
  .status-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 52px;
    padding: 0 12px;

    .diff-summary {
      display: flex;
      .summary-title {
        font-size: 14px;
        font-weight: 700;
      }
      .summray-content {
        display: flex;
        align-items: center;
        gap: 24px;
      }
    }
    .filter-checkbox {
      font-size: 12px;
    }
  }

  .model-group-container {
    display: flex;
    flex-direction: column;
    gap: 24px;
    height: calc(100% - 52px);
    padding: 0 12px;
    @include scrollbar-y;

    .model-group {
      :deep(.collapse-trigger) {
        font-weight: 400;
      }
    }
  }

  .field-list {
    display: grid;
    gap: 16px;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    width: 100%;
    align-content: flex-start;
    margin-top: 12px;

    .field-card {
      &.new {
        background: #F2FFF4;
      }
      &.update {
        background: #FFF3E1;
      }
      &.conflict {
        background: #FFEEEE;
      }
      &.unchanged {
        background: #FFF;
      }
      &.unbinded {
        background: #F0F1F5;
      }

      .conflict-icon {
        font-size: 14px;
        color: $dangerColor;
      }
    }
  }

  .tips-content {
    font-size: 12px;
    .list-item {
      margin-left: 1em;
      li {
        list-style-type: disc;
      }
    }
  }
</style>
