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
  import { t } from '@/i18n'
  import { useStore } from '@/store'
  import useGroupProperty from '@/hooks/utils/group-property'
  import useGroup from '@/hooks/model/group'
  import useProperty from '@/hooks/model/property'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import DiffBrand from './diff-brand.vue'
  import { diffFieldList } from './use-field'

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

  const previewShow = ref(false)

  const diffValue = ref('')

  const tipsName = computed(() => {
    const names = {
      create: t('新增字段'),
      update: t('字段配置更新'),
      conflict: t('字段冲突，该字段已经被其他模板绑定，请删除该模型或修改模板'),
      unbinded: t('不再纳管该字段'),
      unchanged: t('无变化')
    }
    return names[diffValue.value] || ''
  })
  const beforeDiffList = ref([])

  const afterDiffList = ref([])

  const beforeBindField = ref({})

  const afterBindField = ref({})

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

  const handleClickCard = (field) => {
    console.log(field)
    previewShow.value = true
    diffValue.value = getFieldCardClassNames(field) === 'new' ? 'create' : getFieldCardClassNames(field)
    beforeBindField.value = field
    afterBindField.value = props.templateFieldList.find(item => item.bk_property_id === field.bk_property_id)
    beforeDiffList.value = [...diffFieldList]
    afterDiffList.value = [...diffFieldList]
    getDiffLabel(field?.bk_property_type, beforeDiffList?.value, (value) => {
      beforeDiffList.value = value
    })
    getDiffLabel(afterBindField.value?.bk_property_type, afterDiffList?.value, (value) => {
      afterDiffList.value = value
    })
  }
  const getDiffLabel = (fieldType, diffList, setDiff) => {
    const diff = diffList.find(item => item.type === 'option')
    if (fieldType && ['int', 'float'].includes(fieldType)) {
      diffList.push({
        label: t('单位'),
        type: 'unit',
        value: ''
      })
    }
    switch (fieldType) {
      case 'singlechar':
      case 'longchar':
        diff.label = t('字段设置')
        break
      case 'int':
      case 'float':
        diff.label = t('数值范围')
        break
      case 'enum':
      case 'enummulti':
        diff.label = t('枚举值设置')
        break
      case 'list':
        diff.label = t('列表值设置')
        break
      default:
        setDiff([...diffFieldList].filter(item => item.type !== 'option'))
    }
  }

  const getDiffContent = (type, value) => {
    const bindField = type === 'before' ? beforeBindField.value : afterBindField.value
    if (bindField && !Object.keys(bindField)?.length) return
    const fieldValue = type === 'before' ? bindField[value] : bindField[value].value
    switch (value) {
      case 'isrequired':
      case 'editable':
        return fieldValue ? '是' : '否'
      case 'placeholder':
        return fieldValue || '--'
      case 'option':
        if (['int', 'float', 'enum', 'enummulti', 'list'].includes(bindField.bk_property_type)) return  ''
      case 'default':
        if (['enum', 'enummulti'].includes(bindField.bk_property_type)) {
          return bindField.option.filter(item => item.is_default).map(item => item.name)
            .join()
        }
        if (['bool'].includes(bindField.bk_property_type)) {
          return bindField.option ? 'true' : 'false'
        }
      default:
        return bindField[value] || '--'
    }
  }

  const getAfterBindFieldStyle = (diff) => {
    if (diffValue.value) {
      if (diffValue.value === 'update') {
        let value = ''
        const fieldData = props.diffs?.update.find(item => item.bk_property_id === beforeBindField.value.bk_property_id)
        value = fieldData?.update_data[diff.type] !== void 0 ? 'update' : 'unchange'
        return value
      }
      if (diffValue.value === 'create') {
        return 'new'
      }
      return diffValue.value
    }
  }

  const isShowFieldConfig = (type, diff, bindField, list) => {
    if (!bindField && !Object.keys(bindField)?.length) return
    if (type === 'before') {
      return diffValue.value !== 'create' && diff.type === 'option' && list.includes(bindField.bk_property_type)
    }
    return diff.type === 'option' && list.includes(bindField.bk_property_type)
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
          <diff-brand :count="counts.unbinded" :text="$t('解除纳管')" status="unbinded"
            :tooltips="$t('模板中删除了该字段，后续不再统一管理该字段')">
          </diff-brand>
          <diff-brand :count="counts.unchanged" :text="$t('无变化')" status="unchanged"></diff-brand>
          <span class="tips-content" id="field-template-field-diff-conflict-tooltips">
            <div>{{ $t('字段冲突的情况：') }}</div>
            <ul class="list-item">
              <li>{{$t('模板字段与模型字段 ID 类型一样，但已经被其他模板绑定')}}</li>
              <li>{{$t('模板字段与模型字段的 ID 一样，但字段类型不一致')}}</li>
              <li>{{$t('模板设置的唯一性校验与模型设置的冲突')}}</li>
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
            :is-template="isTemplate(field)"
            @click-field="handleClickCard(field)">
            <template #flag-append v-if="isConflict(field)">
              <i class="bk-icon icon-exclamation-circle-shape conflict-icon"></i>
            </template>
          </field-card>
        </div>
      </cmdb-collapse>
    </div>
    <bk-sideslider
      ref="sidesliderComp"
      v-transfer-dom
      :width="880"
      :title="$t('更新差异')"
      :is-show.sync="previewShow">
      <div slot="content">
        <div class="change-tips">
          <span class="tips-left">{{ $t('绑定变化') }}</span>:
          <span :class="['tips-right',diffValue]">{{ tipsName }}</span>
        </div>
        <div class="diff-table">
          <div class="table-head">
            <div class="col before-col">{{$t('绑定前')}}</div>
            <div class="col after-col">{{$t('绑定后')}}</div>
          </div>
          <div class="table-body">
            <div class="col before-col">
              <div :class="['diff-item',{
                     'big-opiton': isShowFieldConfig('before',diff,beforeBindField,
                                                     ['int','float','enum', 'enummulti','list'])
                   }]"
                v-for="diff of beforeDiffList" :key="diff.bk_obj_id">
                <span class="diff-item-label">
                  {{diff.label}}
                </span>
                <div class="diff-item-config"
                  v-if="isShowFieldConfig('before',diff,beforeBindField,['int','float'])">
                  <div>
                    <span>{{$t('最大值')}}</span>
                    <span>{{ beforeBindField.option.max || '--' }}</span>
                  </div>
                  <div>
                    <span>{{$t('最小值')}}</span>
                    <span>{{ beforeBindField.option.min || '--' }}</span>
                  </div>
                </div>
                <div class="diff-item-config"
                  v-else-if="isShowFieldConfig('before',diff,beforeBindField,['enum','enummulti'])">
                  <div v-for="option of beforeBindField.option" :key="option.id">
                    <span>{{option.id || '--'}}</span>
                    <span>{{ option.name || '--' }}</span>
                  </div>
                </div>
                <div class="diff-item-config"
                  v-else-if="isShowFieldConfig('before',diff,beforeBindField,['list'])">
                  <div v-for="option of beforeBindField.option" :key="option.id">
                    <p>{{ option || '--'}}</p>
                  </div>
                </div>
                <span v-if="diffValue === 'create'">--</span>
                <span v-else>
                  {{getDiffContent('before',diff.type)}}</span>
              </div>
            </div>
            <div class="col after-col">
              <div class="diff-item all-col" v-if="['unchanged','unbinded'].includes(diffValue)">
                <p>{{$t(tipsName)}}</p>
              </div>
              <div class="col after-col" v-else>
                <div :class="['diff-item',getAfterBindFieldStyle(diff),
                              { 'big-opiton': isShowFieldConfig('after',diff,afterBindField,
                                                                ['int','float','enum', 'enummulti','list']) }]"
                  v-for="diff of afterDiffList" :key="diff.bk_obj_id">
                  <span class="diff-item-label">{{diff.label}}</span>
                  <div class="diff-item-config"
                    v-if="isShowFieldConfig('after',diff,afterBindField,['int','float'])">
                    <div>
                      <span>{{$t('最大值')}}</span>
                      <span>{{ afterBindField.option.max || '--'}}</span>
                    </div>
                    <div>
                      <span>{{$t('最小值')}}</span>
                      <span>{{ afterBindField.option.min || '--' }}</span>
                    </div>
                  </div>
                  <div class="diff-item-config"
                    v-else-if="isShowFieldConfig('after',diff,afterBindField,['enum','enummulti'])">
                    <div v-for="option of afterBindField.option" :key="option.id">
                      <span>{{option.id || '--'}}</span>
                      <span>{{ option.name || '--' }}</span>
                    </div>
                  </div>
                  <div class="diff-item-config"
                    v-else-if="isShowFieldConfig('after',diff,afterBindField,['list'])">
                    <div v-for="option of afterBindField.option" :key="option.id">
                      <p>{{ option || '--' }}</p>
                    </div>
                  </div>
                  <span v-else>{{getDiffContent('after',diff.type) || '--' }}</span>
                  <bk-icon v-if="['editable','placeholder','isrequired']
                    .includes(diff.type) && afterBindField?.[diff.type]?.lock" type="lock" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </bk-sideslider>
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

  .change-tips{
    margin: 20px 40px;

    .tips-left{
      width: 70px;
      height: 19px;
      font-weight: 700;
      font-size: 14px;
      color: #63656E;
    }

    .tips-right{
      width: 42px;
      height: 19px;
      font-size: 14px;
        &.create {
          color: #2DCB56;
        }
        &.conflict {
          color: #EA3636;
        }
        &.update {
          color: #FF9C01;
        }
    }
  }
  .diff-table {
    display: grid;
    grid-template-rows: 42px auto;
    height: calc(100% - 52px);
    @include scrollbar-y;
    margin: 0 40px;
    border: 1px solid #DCDEE5;

    .table-head {
      display: grid;
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
      background: #FFF;
      box-shadow: inset 0 1px 0 0 #DCDEE5;

      .col {
        display: flex;
        flex-direction: column;
        gap: 8px;
        padding: 0 24px 0 16px;
        overflow: hidden;
        span:first-child {
              display: inline-block;
              position: relative;
              padding-right: 14px;
              &::after {
                position: absolute;
                right: 0;
                content: "：";
              }
            }
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
        &.conflict {
          color: #EA3636;
          background: #FEF2F2;
        }
        &.unchanged {
          background: #F5F7FA;
        }
        &.unbinded {
          background: #F0F1F5;
        }
        &.update {
          color: #FF9C01;
          background:#FFF9EF;
        }
        &-label {
          color: #63656E !important;
        }
        &.all-col {
          height: 100%;
          display: flex;
          justify-content: center;
          align-items: center;
        }
        &.big-opiton {
          min-height: 28px;
          height: auto;
          position: relative;
          display: flex;
          align-items: flex-start;
          .diff-item-label {
            margin-top: 10px;
          }
          .diff-item-config{
            margin: 15px 15px 15px 0;
            line-height: 20px;
          }
        }
      }
    }
  }
</style>
