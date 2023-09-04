<template>
  <bk-table
    :class="['data-row', { 'is-on-empty-add': !readonly && tableData.length === 0 && adding }]"
    :data="tableData"
    :max-height="420"
    :show-header="showHeader"
    v-bkloading="{ isLoading }">
    <bk-table-column
      v-for="prop in header"
      :key="prop.bk_property_id"
      :label="$tools.getHeaderPropertyName(prop)"
      :min-width="$tools.getHeaderPropertyMinWidth(prop, { min: 120 })">
      <template #default="{ row, $index }">
        <!-- 只读模式 -->
        <cmdb-property-value
          v-if="!editState.index.includes($index)"
          :value="row[prop.bk_property_id]"
          :show-unit="false"
          :property="prop"
          :is-show-overflow-tips="$tools.isShowOverflowTips(prop)" />
        <!-- 编辑模式 -->
        <property-form-element
          v-else
          :class="['detault-form-el', prop.bk_property_type]"
          :ref="`property-form-el-${$index}`"
          :property="prop"
          :disabled="checkDisabled(prop)"
          :disabled-tips="$t('系统限定不可修改')"
          :size="'small'"
          :font-size="'normal'"
          :row="1"
          error-display-type="tooltips"
          v-model="editState.row[$index][prop.bk_property_id]" />
      </template>
    </bk-table-column>
    <bk-table-column :label="$t('操作')" width="130" fixed="right" v-if="!readonly">
      <template #default="{ $index }">
        <template v-if="editState.index.includes($index)">
          <bk-button text @click="handleSaveRow($index)">{{ $t(immediate ? '保存' : '确定') }}</bk-button>
          <bk-button text class="ml10" @click="handleCancelEdit($index)">{{ $t('取消') }}</bk-button>
        </template>
        <template v-else>
          <template v-if="disabled">
            <bk-button text
              :disabled="disabled"
              @click="handleEditRow($index)">
              <span v-bk-tooltips="{ disabled: !disabled, allowHtml: true, content: disabledTips }">
                {{ $t('编辑') }}
              </span>
            </bk-button>
          </template>
          <template v-else>
            <cmdb-auth :auth="auth">
              <template #default="authProps">
                <bk-button text
                  :disabled="authProps.disabled"
                  @click="handleEditRow($index)">
                  {{ $t('编辑') }}
                </bk-button>
              </template>
            </cmdb-auth>
          </template>

          <template v-if="disabled">
            <bk-button text class="ml10"
              :disabled="disabled"
              @click="handleDeleteRow($index)">
              <span v-bk-tooltips="{ disabled: !disabled, allowHtml: true, content: disabledTips }">
                {{ $t('删除') }}
              </span>
            </bk-button>
          </template>
          <template v-else>
            <cmdb-auth :auth="auth">
              <template #default="authProps">
                <bk-button text class="ml10"
                  :disabled="authProps.disabled"
                  @click="handleDeleteRow($index)">
                  {{ $t('删除') }}
                </bk-button>
              </template>
            </cmdb-auth>
          </template>
        </template>
      </template>
    </bk-table-column>
    <template #empty v-if="!readonly">
      <template v-if="disabled">
        <icon-text-button
          class="table-empty-add-button"
          ref="tableEmptyAddButtonRef"
          :text="$t('新增')"
          :disabled="true"
          :disabled-tips="disabledTips"
          @click="handleClickAdd" />
      </template>
      <template v-else>
        <cmdb-auth :auth="auth">
          <template #default="authProps">
            <icon-text-button
              class="table-empty-add-button"
              ref="tableEmptyAddButtonRef"
              :text="$t('新增')"
              :disabled="authProps.disabled"
              @click="handleClickAdd" />
          </template>
        </cmdb-auth>
      </template>
    </template>
  </bk-table>
</template>
<script setup>
  import { $bkInfo } from '@/magicbox'
  import { nextTick, ref, set, watch, reactive, computed, getCurrentInstance, onMounted } from 'vue'
  import { t } from '@/i18n'
  import cloneDeep from 'lodash/cloneDeep'
  import { formatValues } from '@/utils/tools'
  import PropertyFormElement from '../property-form-element.vue'
  import IconTextButton from '@/components/ui/button/icon-text-button.vue'
  import { $success } from '@/magicbox/index.js'
  import instanceTableService from '@/service/instance/table'

  const props = defineProps({
    type: {
      type: String,
      default: 'list'
    },
    property: {
      type: Object,
      default: () => ({}),
      required: true
    },
    // 表格数据，支持v-model
    value: {
      type: Array,
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    },
    disabledTips: {
      type: String
    },
    // 只读没有操作入口
    readonly: {
      type: Boolean,
      default: false
    },
    showHeader: {
      type: Boolean,
      default: true
    },
    // 是否立即保存
    immediate: {
      type: Boolean,
      default: true
    },
    objId: {
      type: String,
      default: ''
    },
    // 实例ID
    instanceId: {
      type: [String, Number],
      default: ''
    },
    auth: {
      type: [Object, Array],
      default: () => ({})
    },
    adding: {
      type: Boolean,
      default: false
    },
    mode: {
      type: String,
      default: 'create'
    }
  })
  const emit = defineEmits(['input', 'cancel', 'save', 'delete', 'add'])

  const instacne = getCurrentInstance().proxy

  const header = computed(() => props.property?.option?.header || [])

  const isAddType = computed(() => props.type === 'add')

  const isLoading = ref(false)
  const tableEmptyAddButtonRef = ref(null)
  const tableData = ref(cloneDeep(props.value))
  watch(() => props.value, (value) => {
    tableData.value = cloneDeep(value)
  }, { deep: true })

  const editState = reactive({
    index: [],
    row: {}
  })

  const scrollAddButton = () => {
    if (tableEmptyAddButtonRef.value) {
      const emptyTextEl = tableEmptyAddButtonRef.value.$el?.closest('.bk-table-empty-text')
      const emptyTextWidth = emptyTextEl.offsetWidth
      if (emptyTextEl?.offsetParent) {
        emptyTextEl.offsetParent.scrollLeft = (emptyTextEl.parentElement.offsetWidth - emptyTextWidth) / 2
      }
    }
  }

  const checkDisabled = (property) => {
    if (props.mode === 'create') {
      return false
    }
    return !property.editable
  }

  watch(() => props.adding, (adding) => {
    // 进入新增状态
    if (adding && isAddType.value) {
      // 默认只有一行，索引为0
      const index = 0
      editState.index.push(index)
      set(editState.row, index, cloneDeep(tableData.value[index] || {}))
      focus()
    } else {
      nextTick(scrollAddButton)
    }
  })
  // 聚焦第一个输入框
  const focus = (index = 0) => {
    nextTick(() => {
      const component = instacne.$refs[`property-form-el-${index}`]?.[0]?.$refs?.[`component-${header.value?.[0].bk_property_id}`]
      component?.focus?.()
    })
  }

  const exitEdit = (index) => {
    const dataIndex = editState.index.findIndex(i => i === index)
    if (dataIndex !== -1) {
      editState.index.splice(dataIndex, 1)
      set(editState.row, index, {})
    }
  }

  // 编辑
  const handleEditRow = async (index) => {
    editState.index.push(index)
    set(editState.row, index, cloneDeep(tableData.value[index]))
    focus(index)
  }

  // 取消编辑
  const handleCancelEdit = (index) => {
    exitEdit(index)
    emit('cancel', tableData.value)

    if (isAddType.value) {
      tableData.value.splice(index, 1)
    }
  }

  // 删除
  const deleteRow = (index) => {
    const row = tableData.value.splice(index, 1)
    emit('input', tableData.value)
    emit('delete', row)

    if (tableData.value.length === 0) {
      nextTick(scrollAddButton)
    }
  }

  const validateAll = async (index) => {
    // 获得每一个表单元素的校验方法
    const validates = (instacne.$refs[`property-form-el-${index}`] || [])
      .map(formElement => formElement.$validator.validateAll())

    if (validates.length) {
      const results = await Promise.all(validates)
      return results.every(valid => valid)
    }

    return true
  }

  const handleDeleteRow = (index) => {
    if (props.immediate) {
      $bkInfo({
        title: t('确定删除'),
        extCls: 'bk-dialog-sub-header-center',
        confirmLoading: true,
        confirmFn: async () => {
          const row = tableData.value[index]
          try {
            await instanceTableService.deletemany({
              bk_obj_id: props.objId,
              bk_property_id: props.property.bk_property_id,
              ids: [row.id]
            })

            deleteRow(index)
            return true
          } catch (err) {
            return false
          }
        },
      })
    } else {
      deleteRow(index)
    }
  }

  // 保存
  const saveRow = (rowValues, index) => {
    set(tableData.value, index, rowValues)
    emit('input', tableData.value)
    emit('save', rowValues)
  }
  const handleSaveRow = async (index) => {
    if (!await validateAll(index)) {
      return
    }

    const row = cloneDeep(editState.row[index])
    const values = formatValues(row, header.value)

    if (props.immediate) {
      isLoading.value = true
      let req = null
      const baseParams = {
        bk_obj_id: props.objId,
        bk_property_id: props.property.bk_property_id
      }
      if (row.id) {
        req = instanceTableService.update({
          ...baseParams,
          ids: [row.id],
          data: values
        })
      } else {
        req = instanceTableService.create({
          ...baseParams,
          data: [{ ...values, bk_inst_id: props.instanceId }]
        })
      }
      try {
        await req
        $success(t('操作成功'))
        saveRow(values, index)
        exitEdit(index)
      } finally {
        isLoading.value = false
      }
    } else {
      saveRow(values, index)
      exitEdit(index)
    }
  }

  const handleClickAdd = () => {
    emit('add')
  }

  onMounted(() => {
    nextTick(scrollAddButton)
  })
</script>

<style lang="scss" scoped>
.data-row {
  &.is-on-empty-add {
    :deep(.bk-table-empty-block) {
      display: none;
    }
  }

  .table-empty-add-button {
    // 避免遮挡
    position: relative;
    z-index: 1;
  }

  &:focus-within {
    &.bk-table-scrollable-x,
    &.bk-table-scrollable-y {
      overflow: auto !important;
      :deep(.bk-table-body-wrapper) {
        overflow: auto !important;
      }
    }

    overflow: visible !important;
    :deep(.bk-table-body-wrapper) {
      overflow: visible !important;
    }
  }

  .detault-form-el {
    &:focus-within {
      &.longchar {
        position: absolute;
        left: -1px;
        top: 2px;
        z-index: 1;
        :deep(.bk-form-textarea) {
          min-height: 90px !important;
        }
        :deep(.control-icon) {
          display: none;
        }
      }
    }
  }
}
</style>
