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
        <!-- 编辑模式 -->
        <property-form-element
          ref="propertyFormEl"
          :property="prop"
          :size="'small'"
          :font-size="'normal'"
          :row="1"
          error-display-type="tooltips"
          v-model.trim="row[prop.bk_property_id]"
          v-if="editRowIndex === $index || isAddType" />
        <!-- 只读模式 -->
        <cmdb-property-value
          :value="row[prop.bk_property_id]"
          :show-unit="false"
          :property="prop"
          v-else />
      </template>
    </bk-table-column>
    <bk-table-column :label="$t('操作')" width="120" fixed="right" v-if="!readonly">
      <template #default="{ row, $index }">
        <template v-if="editRowIndex === $index || isAddType">
          <bk-button text @click="handleSaveRow(row, $index)">{{ $t(immediate ? '保存' : '确定') }}</bk-button>
          <bk-button text class="ml10" @click="handleCancelEdit($index)">{{ $t('取消') }}</bk-button>
        </template>
        <template v-else>
          <template v-if="!editable">
            <bk-button text
              @click="handleEditRow($index)">
              <span v-bk-tooltips="{ disabled: editable, content: $t('系统限定不可修改') }">{{ $t('编辑') }}</span>
            </bk-button>
          </template>
          <template v-else>
            <cmdb-auth :auth="auth">
              <template #default="{ disabled }">
                <bk-button text
                  :disabled="disabled"
                  @click="handleEditRow($index)">
                  {{ $t('编辑') }}
                </bk-button>
              </template>
            </cmdb-auth>
          </template>

          <template v-if="!editable">
            <bk-button text class="ml10"
              @click="handleDeleteRow($index)">
              <span v-bk-tooltips="{ disabled: editable, content: $t('系统限定不可修改') }">{{ $t('删除') }}</span>
            </bk-button>
          </template>
          <template v-else>
            <cmdb-auth :auth="auth">
              <template #default="{ disabled }">
                <bk-button text class="ml10"
                  :disabled="disabled"
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
      <template v-if="!editable">
        <icon-text-button
          :text="$t('新增')"
          :disabled="true"
          :disabled-tips="$t('系统限定不可修改')"
          @click="handleClickAdd" />
      </template>
      <template v-else>
        <cmdb-auth :auth="auth">
          <template #default="{ disabled }">
            <icon-text-button
              :text="$t('新增')"
              :disabled="disabled"
              @click="handleClickAdd" />
          </template>
        </cmdb-auth>
      </template>
    </template>
  </bk-table>
</template>
<script setup>
  import { $bkInfo } from '@/magicbox'
  import { nextTick, ref, set, watch, computed } from 'vue'
  import { t } from '@/i18n'
  import { clone } from '@/utils/tools'
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
    }
  })
  const emit = defineEmits(['input', 'cancel', 'save', 'delete', 'add'])

  const header = computed(() => props.property?.option?.header || [])

  const isAddType = computed(() => props.type === 'add')

  const propertyFormEl = ref(null)
  const isLoading = ref(false)
  const tableData = ref(clone(props.value))
  watch(() => props.value, (value) => {
    tableData.value = clone(value)
  }, { deep: true })

  const editRowIndex = ref(-1)

  const editable = computed(() => props.property.editable && !props.property.bk_isapi)

  watch(() => props.adding, (adding) => {
    if (adding) {
      if (props.type === 'list') {
        set(tableData.value, editRowIndex.value, clone(props.value[editRowIndex.value] || {}))// 还原初始值
        editRowIndex.value = -1
      }

      focus()
    }
  })
  // 聚焦第一个输入框
  const focus = () => {
    nextTick(() => {
      const component = propertyFormEl.value?.[0]?.$refs?.[`component-${header.value?.[0].bk_property_id}`]
      component?.focus?.()
    })
  }

  // 编辑
  const handleEditRow = async (index) => {
    handleCancelEdit(editRowIndex.value)// 取消上一次的编辑
    editRowIndex.value = index
    focus()
  }

  // 取消编辑
  const handleCancelEdit = (index) => {
    if (index <= -1) return
    set(tableData.value, index, clone(props.value[index] || {}))// 还原初始值
    editRowIndex.value = -1
    emit('cancel', tableData.value)
  }

  // 删除
  const deleteRow = (index) => {
    const row = tableData.value.splice(index, 1)
    emit('input', tableData.value)
    emit('delete', row)
  }

  const validateAll = async () => {
    // 获得每一个表单元素的校验方法
    const validates = (propertyFormEl.value || [])
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
  const saveRow = (row) => {
    editRowIndex.value = -1
    emit('input', tableData.value)
    emit('save', row)
  }
  const handleSaveRow = async (row) => {
    if (!await validateAll()) {
      return
    }

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
          data: row
        })
      } else {
        req = instanceTableService.create({
          ...baseParams,
          data: [{ ...row, bk_inst_id: props.instanceId }]
        })
      }
      try {
        await req
        $success(t('操作成功'))
        saveRow(row)
      } finally {
        isLoading.value = false
      }
    } else {
      saveRow(row)
    }
  }

  const handleClickAdd = () => {
    emit('add')
  }
</script>

<style lang="scss" scoped>
.data-row {
  &.is-on-empty-add {
    :deep(.bk-table-empty-block) {
      display: none;
    }
  }
}
</style>
