<template>
  <bk-table :data="tableData" :max-height="420" :show-header="showHeader" v-bkloading="{ isLoading }">
    <bk-table-column
      v-for="item in header"
      :key="item.bk_property_id"
      :label="item.unit ? `${item.bk_property_name} (${item.unit})` : item.bk_property_name">
      <template #default="{ row, $index }">
        <!-- 编辑模式 -->
        <PropertyFormElement
          :ref="`com-${item.bk_property_id}-${$index}`"
          :property="item"
          :row="1"
          v-model.trim="row[item.bk_property_id]"
          v-if="editableRowIndex === $index" />
        <!-- 只读模式 -->
        <cmdb-property-value
          :value="row[item.bk_property_id]"
          :property="item"
          v-else />
      </template>
    </bk-table-column>
    <bk-table-column :label="$t('操作')" width="120" fixed="right" v-if="!readonly">
      <template #default="{ row, $index }">
        <template v-if="editableRowIndex === $index">
          <bk-button text @click="handleSaveRow(row, $index)">{{ $t('保存') }}</bk-button>
          <bk-button text class="ml10" @click="handleCancelEdit($index)">{{ $t('取消') }}</bk-button>
        </template>
        <template v-else>
          <bk-button text @click="handleEditRow($index)">{{ $t('编辑') }}</bk-button>
          <bk-button text class="ml10" @click="handleDeleteRow($index)">{{ $t('删除') }}</bk-button>
        </template>
      </template>
    </bk-table-column>
  </bk-table>
</template>
<script lang="ts" setup>
  /* eslint-disable camelcase */
  import { $bkInfo } from '@/magicbox'
  import {  getCurrentInstance, nextTick, PropType, ref, set, watch, computed } from 'vue'
  import { t } from '@/i18n'
  import { clone } from '@/utils/tools'
  import PropertyFormElement from '../property-form-element.vue'
  import { actions } from '@/store/modules/api/table-instance'
  import { $success } from '@/magicbox/index.js'

  export interface IHeader {
    id: string
    bk_property_name: string
    bk_property_id: string
    bk_property_type: string
    unit: string
    option: any[]
    ismultiple: boolean
  }

  export interface IProperty {
    bk_property_id: string
    option: {
      header: IHeader[]
      default: any[]
    }
  }

  const { proxy } = getCurrentInstance()
  const props = defineProps({
    property: {
      type: Object as PropType<IProperty>,
      default: () => ({}),
      required: true
    },
    // 表格数据，支持v-model
    value: {
      type: Array,
      default: () => []
    },
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
    // 初始化编辑行
    defaultEditRowIndex: {
      type: Number,
      default: -1
    },
    // 模型ID eg: host biz project biz_set
    objId: {
      type: String,
      default: ''
    },
    // 实例ID
    instanceId: {
      type: [String, Number],
      default: ''
    }
  })
  const emit = defineEmits(['input', 'cancel', 'save', 'delete'])

  const header = computed(() => props.property?.option?.header || [])

  const isLoading = ref(false)
  const tableData = ref(clone(props.value))
  watch(() => props.value, (value) => {
    tableData.value = clone(value)
  }, { deep: true })

  const editableRowIndex = ref<number>(props.defaultEditRowIndex)

  // 聚焦第一个输入框
  const focus = (index: number) => {
    const firstInputProp = header.value?.[0]?.bk_property_id
    const firstInputRef = proxy.$refs?.[`com-${firstInputProp}-${index}`]?.[0]
    firstInputRef?.focus()
  }

  // 编辑
  const handleEditRow = async (index: number) => {
    handleCancelEdit(editableRowIndex.value)// 取消上一次的编辑
    editableRowIndex.value = index
    await nextTick()
    focus(index)
  }

  // 取消编辑
  const handleCancelEdit = (index: number) => {
    if (index <= -1) return
    set(tableData.value, index, clone(props.value[index] || {}))// 还原初始值
    editableRowIndex.value = -1
    emit('cancel', tableData.value)
  }

  const handleUpdateTableData = async (data: any[]) => {
    if (!props.instanceId) {
      console.warn('instanceId is requierd when immediate prop is true')
      return false
    }
    let result = false
    const params = {
      bk_host_id: String(props.instanceId),
      [props.property.bk_property_id]: data
    }
    switch (props.objId) {
      case 'host':
        // 更新主机表格字段
        result = await actions.updateTableHostsBatch(null, { params })
        break
      case 'biz_set':
        // 更新业务集表格字段
        result = await actions.updateTableBizSet(null, { params })
        break
      case 'biz':
        // 更新集群表格字段
        result = await actions.updateTableBiz(null, { params })
        break
      case 'module':
        // 更新模块表格字段
        result = await actions.updateTableModule(null, { params })
        break
      default:
        // 通用更新逻辑
        result = await actions.updateTableInstance(null, { params })
    }
    return result
  }

  // 删除
  const deleteRow = (index: number) => {
    const row = tableData.value.splice(index, 1)
    emit('input', tableData.value)
    emit('delete', row)
  }
  const handleDeleteRow = (index) => {
    if (props.immediate) {
      $bkInfo({
        title: t('确定删除'),
        extCls: 'bk-dialog-sub-header-center',
        confirmLoading: true,
        confirmFn: async () => {
          const tmpData = clone(tableData.value)
          tmpData.splice(index, 1)
          const result = await handleUpdateTableData(tmpData)
          if (result) {
            deleteRow(index)
          }
        },
      })
    } else {
      deleteRow(index)
    }
  }

  // 保存
  const saveRow = (row) => {
    editableRowIndex.value = -1
    emit('input', tableData.value)
    emit('save', row)
  }
  const handleSaveRow = async (row, $index) => {
    // 校验当前行数据
    const validateArr = header.value.map((item) => {
      const comRef = proxy.$refs?.[`com-${item.bk_property_id}-${$index}`]?.[0]
      return comRef.$validator.validateAll()
    })
    const results =  await Promise.all(validateArr)
    if (results.some(result => !result)) return

    if (props.immediate) {
      isLoading.value = true
      const result = await handleUpdateTableData(tableData.value)
      isLoading.value = false
      if (result) {
        $success(t('操作成功'))
        saveRow(row)
      }
    } else {
      saveRow(row)
    }
  }
</script>
<script lang="ts">
  export default {
    name: 'inner-table'
  }
</script>
