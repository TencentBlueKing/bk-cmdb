<template>
  <div class="cmdb-form-innertable">
    <div class="innertable-container" v-bkloading="{ isLoading }">
      <data-row
        type="list"
        :mode="mode"
        :property="property"
        :readonly="readonly"
        :disabled="disabled"
        :disabled-tips="disabledTips"
        :obj-id="objId"
        :instance-id="instanceId"
        :immediate="immediate"
        :auth="auth"
        v-model="tableData"
        :adding="isShowAddRow"
        @add="handleClickAdd" />
      <div class="row-add" v-show="isShowAddRow">
        <data-row
          type="add"
          :mode="mode"
          :property="property"
          :readonly="readonly"
          :disabled="disabled"
          :disabled-tips="disabledTips"
          :show-header="false"
          :value="defaultRowData"
          :obj-id="objId"
          :instance-id="instanceId"
          :immediate="immediate"
          :auth="auth"
          :adding="isShowAddRow"
          @save="handleAddRow"
          @cancel="handleCancelAdd" />
      </div>
      <div class="row-append" v-if="!isShowAddRow && tableData.length > 0 && !readonly">
        <cmdb-auth :auth="auth">
          <template #default="authProps">
            <icon-text-button
              :text="$t('新增')"
              @click="handleClickAdd"
              :disabled="maxRowDisabled || disabled || authProps.disabled"
              :disabled-tips="authProps.disabled ? '' : (maxRowDisabled ? $t('最多添加50行') : disabledTips)" />
          </template>
        </cmdb-auth>
      </div>
    </div>
    <i class="title-copy icon-cc-details-copy" v-show="showCopyBtn" @click="handleCopyTable"></i>
  </div>
</template>
<script setup>
  import { computed, getCurrentInstance, ref, watch } from 'vue'
  import { t } from '@/i18n'
  import { clone, getPropertyDefaultValue } from '@/utils/tools'
  import { $success, $error } from '@/magicbox/index.js'
  import IconTextButton from '@/components/ui/button/icon-text-button.vue'
  import DataRow from './data-row.vue'
  import instanceTableService from '@/service/instance/table'

  const { proxy } = getCurrentInstance()
  const props = defineProps({
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
      type: String,
      default: ''
    },
    // 只读模式
    readonly: {
      type: Boolean,
      default: false
    },
    // 模型ID
    objId: {
      type: String,
      default: ''
    },
    // 实例ID（编辑的时候需要）
    instanceId: {
      type: [String, Number],
      default: ''
    },
    // 是否立即保存
    immediate: {
      type: Boolean,
      default: true
    },
    auth: {
      type: [Object, Array],
      default: () => ({})
    },
    // 新建或编辑或详情模式 update | info
    mode: {
      type: String,
      default: 'create'
    }
  })
  const emit = defineEmits(['input'])

  const defaultRowData = ref([])
  const newRowData = () => {
    const data = {}
    const header = props.property?.option?.header || []
    header.forEach((prop) => {
      data[prop.bk_property_id] = getPropertyDefaultValue(prop)
    })
    return data
  }

  const isShowAddRow = ref(false)

  const tableData = ref(clone(props.value))
  const watchOnce = watch(() => props.value, (value) => {
    tableData.value = clone(value)
    watchOnce()
  }, { deep: true })

  watch(tableData, () => {
    emit('input', tableData.value)
  })

  const maxRowDisabled = computed(() => tableData.value.length === 50)

  const exitAdd = () => {
    isShowAddRow.value = false
  }

  const handleAddRow = (row) => {
    tableData.value.push(row)
    exitAdd()
  }
  const handleCancelAdd = () => {
    exitAdd()
  }

  // 复制表格数据
  const showCopyBtn = ref(false)
  const handleCopyTable = () => {
    proxy.$copyText(JSON.stringify(tableData.value)).then(() => {
      $success(t('复制成功'))
    }, () => {
      $error(t('复制失败'))
    })
  }

  // 获取表格数据
  const isLoading = ref(false)
  const handleGetInstanceData = async () => {
    if (!props.instanceId) {
      // 无实例ID时，取默认值
      tableData.value = props.property?.option?.default || []
      return
    }

    isLoading.value = true
    const { info = [] } = await instanceTableService.find({
      bk_obj_id: props.objId,
      bk_property_id: props.property.bk_property_id,
      filter: {
        condition: 'OR',
        rules: [{
          field: 'bk_inst_id',
          operator: 'equal',
          value: props.instanceId
        }]
      },
      page: {
        limit: 50,
        start: 0
      }
    })
    tableData.value = info
    isLoading.value = false
  }

  const handleClickAdd = () => {
    defaultRowData.value = [newRowData()]
    isShowAddRow.value = true
  }

  watch(() => props.instanceId, handleGetInstanceData, { immediate: true })
</script>
<script>
  export default {
    name: 'cmdb-form-innertable'
  }
</script>
<style lang="scss" scoped>
.cmdb-form-innertable {
  position: relative;
  display: flex;
  width: 100%;
  gap: 5px;
}
.innertable-container {
  width: 100%;
}
.row-add {
  margin-top: -1px;
}
.row-append {
  border: 1px solid #dfe0e5;
  border-top: none;
  padding: 10px;
  background: #fff;
  font-size: 12px;
  &:hover {
    background-color: #f0f1f5;
  }
}

.icon-cc-details-copy {
  color: #3a84ff;
  cursor: pointer;
}
</style>
