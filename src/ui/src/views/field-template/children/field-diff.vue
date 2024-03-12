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
  import { ref, computed, reactive } from 'vue'
  import { useStore } from '@/store'
  import useGroupProperty from '@/hooks/utils/group-property'
  import useGroup from '@/hooks/model/group'
  import useProperty from '@/hooks/model/property'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import DiffBrand from './diff-brand.vue'
  import FieldDiffDetailsDrawer from './field-diff-details-drawer.vue'
  import { DIFF_TYPES } from './use-field'

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
    },
    // 删除的模板字段列表
    templateRemovedFieldList: {
      type: Array,
      default: () => ([])
    }
  })
  const store = useStore()

  const isOnlyShowTemplateRelated = ref(true)

  const diffDetails = reactive({
    show: false,
    title: '',
    diffType: '',
    beforeField: {},
    afterField: {}
  })

  // 查询模型的字段列表
  const propertyParams = computed(() => ({
    bk_obj_id: props.model.bk_obj_id,
    bk_supplier_account: store.getters.supplierAccount
  }))
  const [{ properties, pending }] = useProperty(propertyParams)
  const [{ groups }] = useGroup(propertyParams)
  const groupedPropertyies = useGroupProperty(groups, properties)

  const isConflict = field => props.diffs.conflict?.some(item => item.data?.bk_property_id === field.bk_property_id)
  const isNew = field => props.diffs.create?.some(item => item.bk_property_id === field.bk_property_id)
  const isUpdate = field => props.diffs.update?.some(item => item.bk_property_id === field.bk_property_id)
  const isUnchanged = field => props.diffs.unchanged?.some(item => item.bk_property_id === field.bk_property_id)
  const isUnbound = field => props.templateRemovedFieldList
    ?.some(({ field: item }) => item.id === field.bk_template_id && item.bk_property_id === field.bk_property_id)

  const unboundFieldList = computed(() => properties.value.filter(field => isUnbound(field)))
  const counts = computed(() => ({
    new: props.diffs?.create?.length ?? 0,
    update: props.diffs?.update?.length ?? 0,
    conflict: props.diffs?.conflict?.length ?? 0,
    unbound: unboundFieldList.value?.length ?? 0,
    unchanged: props.diffs?.unchanged?.length ?? 0
  }))

  // 新增的字段，来源于模板
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
      group.properties = group.properties
        .filter(field => (isOnlyShowTemplateRelated.value ? isTemplateRelated(field) : true))
      group.properties.forEach((field) => {
        if (isUpdate(field)) {
          setUpdateFields(field, ['bk_property_name', 'isrequired'])
        }
      })
    })
    const finalFieldGroups = displayFieldGroups.filter(group => group.properties.length > 0)

    return finalFieldGroups
  })

  // 传入的field是模型的字段
  const getFieldDiffType = (field) => {
    // 新增：模型中没有，展示的是模板的字段
    // 更新：共同的字段，但是模板中有更新
    // 冲突：无法应用到模型的字段，因模型中的字段与模板当前的设置冲突
    // 解除：模型的字段在模板中已经找不到
    if (isNew(field)) {
      return DIFF_TYPES.NEW
    }

    if (isUpdate(field)) {
      return DIFF_TYPES.UPDATE
    }

    // 冲突使用模型数据中的字段id匹配
    if (isConflict(field)) {
      return DIFF_TYPES.CONFLICT
    }

    // 由于解除绑定的字段仍然会存在于接口的unchanged数据中，所以这里要放在前面
    if (isUnbound(field)) {
      return DIFF_TYPES.UNBOUND
    }

    if (isUnchanged(field)) {
      return DIFF_TYPES.UNCHANGED
    }
  }

  const setUpdateFields = (field = {}, fields = []) => {
    const updateField = props.templateFieldList
      .find(item => item.bk_property_id === field.bk_property_id) ?? {}
    fields.forEach(item => field[`update_${item}`] = updateField[item] ?? field[item])
  }

  const getFieldCardClassName = field => getFieldDiffType(field)

  const isTemplate = (field) => {
    const diffType = getFieldDiffType(field)
    if (diffType === DIFF_TYPES.NEW) {
      return true
    }
    return props.templateFieldList.some(item => item.id === field.bk_template_id)
  }

  const isTemplateRelated = field => isNew(field)
    || isUpdate(field) || isConflict(field) || isUnbound(field) || isTemplate(field)

  const handleClickField = (field) => {
    const diffType = getFieldDiffType(field)
    diffDetails.diffType = diffType

    // 绑定前的字段即为展示的模型字段
    diffDetails.beforeField = field

    // 无变化/不再纳管不需要展示after，在模板字段中也找不到
    diffDetails.afterField = {}

    // 新增的就是模板的字段
    if (DIFF_TYPES.NEW === diffType) {
      diffDetails.afterField = props.templateFieldList.find(item => item.bk_property_id === field.bk_property_id)
    }

    // 冲突和更新会需要使用具体的diff数据
    diffDetails.fieldDiff = {}
    if (DIFF_TYPES.CONFLICT === diffType) {
      diffDetails.fieldDiff = props.diffs.conflict?.find(item => item.data.bk_property_id === field.bk_property_id)

      // 需通过模型字段找到diff数据再找到模板字段
      diffDetails.afterField = props.templateFieldList
        .find(item => item.bk_property_id === diffDetails.fieldDiff.bk_property_id) ?? {}
    }
    if (DIFF_TYPES.UPDATE === diffType) {
      diffDetails.fieldDiff = props.diffs.update?.find(item => item.data.bk_property_id === field.bk_property_id)
      diffDetails.afterField = props.templateFieldList
        .find(item => item.bk_property_id === diffDetails.fieldDiff.bk_property_id) ?? {}
    }

    diffDetails.show = true
  }

  const handelDetailsDrawerToggle = (val) => {
    diffDetails.show = val
  }
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
          <diff-brand :count="counts.unchanged" :text="$t('无配置变化')" status="unchanged"></diff-brand>
          <span class="tips-content" id="field-template-field-diff-conflict-tooltips">
            <div>{{ $t('字段冲突的情况：') }}</div>
            <ul class="list-item">
              <li>{{$t('模板字段与模型字段 ID 类型一样，但已经被其他模板绑定')}}</li>
              <li>{{$t('模板字段与模型字段的 ID 一样，但字段类型不一致')}}</li>
              <li>{{$t('模板字段与模型字段的字段名称相同，但字段 ID 不一致')}}</li>
            </ul>
          </span>
        </div>
      </div>
      <div class="diff-operate">
        <diff-brand :count="counts.unbound" :text="$t('解除纳管')" status="unbound"
          :tooltips="$t('模板中删除了该字段，后续不再统一管理该字段')">
        </diff-brand>
        <bk-checkbox class="filter-checkbox" v-model="isOnlyShowTemplateRelated">{{ $t('仅显示与模板相关字段') }}</bk-checkbox>
      </div>
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
            :class="getFieldCardClassName(field)"
            :key="index"
            :field="field"
            :sortable="false"
            :deletable="false"
            :is-template="isTemplate(field)"
            @click-field="handleClickField(field)">
            <template #flag-append v-if="isConflict(field)">
              <i class="bk-icon icon-exclamation-circle-shape conflict-icon"></i>
            </template>
          </field-card>
        </div>
      </cmdb-collapse>
    </div>

    <field-diff-details-drawer v-bind="diffDetails" @toggle="handelDetailsDrawerToggle" />
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

    .diff-operate {
      display: flex;
      gap: 24px;
    }

    .filter-checkbox {
      :deep(.bk-checkbox-text) {
        font-size: 12px !important;
      }
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
      &.unbound {
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
