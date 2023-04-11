<template>
  <div class="cmdb-form-inner-table">
    <div class="inner-table-wrapper" v-bkloading="{ isLoading }">
      <InnerTable
        :property="property"
        :readonly="readonly"
        :obj-id="objId"
        :instance-id="instanceId"
        :immediate="immediate"
        v-model="tableData" />
      <div :class="[{ 'add-row-btn': !showInnerTable }, 'add-row']" v-if="!readonly">
        <InnerTable
          :property="property"
          :show-header="false"
          :default-edit-row-index="0"
          :value="[{}]"
          :obj-id="objId"
          :instance-id="instanceId"
          :immediate="immediate"
          v-if="showInnerTable"
          @save="handleAddRow"
          @cancel="handleCancelAdd" />
        <bk-button text size="small" v-else @click="showInnerTable = true">
          <span class="text-18px"><i class="bk-icon icon-plus"></i></span>
          {{ $t('新增') }}
        </bk-button>
      </div>
    </div>
    <i class="title-copy icon-cc-details-copy" v-show="showCopyBtn" @click="handleCopyTable"></i>
  </div>
</template>
<script lang="ts" setup>
  /* eslint-disable camelcase */
  import { getCurrentInstance, onBeforeMount, PropType, ref, watch } from 'vue'
  import { t } from '@/i18n'
  import { clone } from '@/utils/tools'
  import { $success, $error } from '@/magicbox/index.js'
  import InnerTable, { IProperty } from './inner-table.vue'
  import { actions } from '@/store/modules/api/table-instance'

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
    // 只读模式
    readonly: {
      type: Boolean,
      default: false
    },
    // 模型ID
    objId: {
      type: String,
      default: '',
      validator(value: string) {
        return !value || ['host', 'biz_set', 'biz', 'module'].includes(value)
      }
    },
    // 实例ID（编辑的时候需要）
    instanceId: {
      type: [String, Number],
      default: ''
    },
    // 是否立即保存
    immediate: {
      type: Boolean,
      default: false
    },
  })
  const emit = defineEmits(['input'])

  const tableData = ref(clone(props.value))
  const watchOnce = watch(() => props.value, (value) => {
    tableData.value = clone(value)
    watchOnce()
  }, { deep: true })

  watch(tableData, () => {
    emit('input', tableData.value)
  })

  // 添加数据
  const showInnerTable = ref(false)
  const handleAddRow = (row) => {
    tableData.value.unshift(row)
    handleCancelAdd()
  }
  const handleCancelAdd = () => {
    showInnerTable.value = false
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
    const { info = [] } = await actions.findmanyQuotedInstance(null, {
      params: {
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
      }
    })
    tableData.value = info
    isLoading.value = false
  }

  onBeforeMount(() => {
    handleGetInstanceData()
  })
</script>
<script lang="ts">
  export default {
    name: 'cmdb-form-innertable'
  }
</script>
<style lang="scss" scoped>
.cmdb-form-inner-table {
  position: relative;
  display: flex;
}
.inner-table-wrapper {
  max-width: 1200px;
  flex: 1;
  margin-right: 12px;
}
.add-row {
  margin-top: -1px;
  &-btn {
    display: flex;
    align-items: center;
    height: 42px;
    border: 1px solid #dfe0e5;
    &:hover {
      background-color: #f0f1f5;
    }
  }
}
.text-18px {
  font-size: 18px;
}

.icon-cc-details-copy {
  color: #3a84ff;
  cursor: pointer;
}
</style>
