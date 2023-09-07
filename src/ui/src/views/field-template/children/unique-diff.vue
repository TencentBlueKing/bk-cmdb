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
  import { computed } from 'vue'
  import useUniqueCheck from '@/hooks/unique-check'
  import useProperty from '@/hooks/model/property'
  import MiniTag from '@/components/ui/other/mini-tag.vue'
  import DiffBrand from './diff-brand.vue'
  import { getUniqueName as getUniqueNameBase, DIFF_TYPES } from './use-unique'

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
    templateUniqueList: {
      type: Array,
      default: () => ([])
    },
    // 模板字段列表
    templateFieldList: {
      type: Array,
      default: () => ([])
    },
    // 删除的模板唯一校验
    templateRemovedUniqueList: {
      type: Array,
      default: () => ([])
    }
  })

  const currentModelId = computed(() => props.model.bk_obj_id)
  const propertyParams = computed(() => ({
    bk_obj_id: props.model.bk_obj_id
  }))
  const [{ uniqueChecks: modelUniqueList, pending }] = useUniqueCheck(currentModelId)
  const [{ properties }] = useProperty(propertyParams)

  const modelBeforeUniqueList = computed(() => modelUniqueList.value
    .filter(item => item.keys.every(key => properties.value.find(({ id }) => id === key.key_id))))

  const newFieldList = computed(() => {
    const news = props.diffs.create ?? []
    return news.map(item => props.templateUniqueList[item.index])
  })

  const modelAfterUniqueList = computed(() => {
    const afterUniqueList = modelBeforeUniqueList.value.slice()
    afterUniqueList.push(...newFieldList.value)
    return afterUniqueList
  })

  const isConflict = unique => props.diffs.conflict?.some(item => item.data.id === unique.id)
  const isNew = unique => props.templateUniqueList?.some(item => item.id === unique.id)
  const isUpdate = unique => props.diffs.update?.some(item => item.data.id === unique.id)
  const isUnchanged = unique => props.diffs.unchanged?.some(item => item.id === unique.id)
  const isUnbound = unique => props.templateRemovedUniqueList?.some(item => item.id === unique.bk_template_id)

  const unboundUniqueList = computed(() => modelBeforeUniqueList.value.filter(unique => isUnbound(unique)))

  const counts = computed(() => ({
    new: props.diffs?.create?.length ?? 0,
    update: props.diffs?.update?.length ?? 0,
    conflict: props.diffs?.conflict?.length ?? 0,
    unbound: unboundUniqueList.value?.length ?? 0,
    unchanged: props.diffs?.unchanged?.length ?? 0,
  }))

  const getUniqueDiffType = (unique) => {
    if (isNew(unique)) {
      return DIFF_TYPES.NEW
    }

    if (isUpdate(unique)) {
      return DIFF_TYPES.UPDATE
    }

    if (isConflict(unique)) {
      return DIFF_TYPES.CONFLICT
    }

    if (isUnbound(unique)) {
      return DIFF_TYPES.UNBOUND
    }

    if (isUnchanged(unique)) {
      return DIFF_TYPES.UNCHANGED
    }
  }

  const getUniqueName = (unique, isTemplate) => {
    const diffType = getUniqueDiffType(unique)
    const isTemplateUnqiue = isTemplate || diffType === DIFF_TYPES.NEW
    const fieldList = isTemplateUnqiue ? props.templateFieldList : properties.value
    return getUniqueNameBase(unique, fieldList, isTemplateUnqiue)
  }

  const getConflictUniqueName = (unique) => {
    const conflict = props.diffs.conflict?.find(item => item.data.id === unique.id)
    if (conflict) {
      return getUniqueName(props.templateUniqueList[conflict.index], true)
    }
    return 'unknown'
  }

  const getDiffClassName = unique => getUniqueDiffType(unique)

  const isTemplate = unique => props.templateUniqueList.some(item => item.id === unique.bk_template_id)
</script>

<template>
  <div class="unique-diff" v-bkloading="{ isLoading: pending }">
    <div class="status-bar">
      <div class="diff-summary">
        <div class="summary-title">{{$t('模板应用后的差异对比：')}}</div>
        <div class="summray-content">
          <diff-brand :count="counts.new" :text="$t('新增校验')" status="new"></diff-brand>
          <diff-brand :count="counts.update" :text="$t('更新覆盖')" status="update"></diff-brand>
          <diff-brand :count="counts.conflict" :text="$t('校验冲突')" status="conflict"
            :tooltips="$t('当前模板设置的唯一校验规则与模型已存在的规则冲突')">
          </diff-brand>
          <diff-brand :count="counts.unbound" :text="$t('解除纳管')" status="unbound"
            :tooltips="$t('模板中删除了该字段，后续不再统一管理该字段')">
          </diff-brand>
          <diff-brand :count="counts.unchanged" :text="$t('无变化')" status="unchanged"></diff-brand>
        </div>
      </div>
    </div>
    <div class="diff-table">
      <div class="table-head">
        <div class="col before-col">{{$t('绑定前的唯一校验')}}</div>
        <div class="col after-col">{{$t('绑定后的唯一校验')}}</div>
      </div>
      <div class="table-body">
        <div class="col before-col">
          <div class="diff-item" v-for="(unique, index) in modelBeforeUniqueList" :key="index">
            {{ getUniqueName(unique) }}
            <mini-tag :text="$t('模板')" v-if="isTemplate(unique)" />
          </div>
        </div>
        <div class="col after-col">
          <div v-for="(unique, index) in modelAfterUniqueList" :key="index"
            :class="['diff-item', getDiffClassName(unique)]">
            <span>{{ getUniqueName(unique) }}</span>
            <mini-tag :text="$t('模板')" v-if="isTemplate(unique)" />
            <template v-if="isConflict(unique)">
              <i class="bk-icon icon-exclamation-circle-shape conflict-icon"
                v-bk-tooltips="{ content: $t('绑定后的唯一性校验冲突提示语', { target: getConflictUniqueName(unique) }) }">
              </i>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .unique-diff {
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

  .diff-table {
    display: grid;
    grid-template-rows: 42px auto;
    padding: 0 12px;
    height: calc(100% - 52px);
    @include scrollbar-y;

    .table-head {
      display: grid;
      gap: 4px;
      grid-template-columns: 1fr 1fr;
      font-size: 12px;
      font-weight: 700;
      line-height: 42px;

      .col {
        padding-left: 24px;
        overflow: hidden;
      }
      .before-col {
        background: #F5F7FA;
      }
      .after-col {
        background: #F0F1F5;
      }
    }

    .table-body {
      display: grid;
      gap: 4px;
      grid-template-columns: 1fr 1fr;
      padding: 12px 0;
      font-size: 12px;
      background: #FAFBFD;
      box-shadow: inset 0 1px 0 0 #DCDEE5, inset 0 -1px 0 0 #DCDEE5;

      .col {
        display: flex;
        flex-direction: column;
        gap: 8px;
        padding: 0 24px 0 16px;
        overflow: hidden;
      }

      .diff-item {
        display: flex;
        align-items: center;
        gap: 4px;
        height: 28px;
        width: 100%;
        background: #F5F7FA;
        padding-left: 12px;
        @include ellipsis;

        &.new {
          color: #2DCB56;
          background: #F2FFF4;
        }
        &.update {
          color: #FF9C01;
          background: #FFF9EF;
        }
        &.conflict {
          color: #EA3636;
          background: #FEF2F2;

          &:hover {
            background: #FFDDDD;
          }
        }
        &.unchanged {
          background: #F5F7FA;
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
  }
</style>
