<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script setup>
  import { computed, reactive, ref, watch, onUnmounted } from 'vue'
  import debounce from 'lodash/debounce'
  import { useStore } from '@/store'
  import { t } from '@/i18n'
  import { $bkInfo, $success } from '@/magicbox/index.js'
  import fieldTemplateService from '@/service/field-template'
  import MiniTag from '@/components/ui/other/mini-tag.vue'
  import queryBuilderOperator, { QUERY_OPERATOR } from '@/utils/query-builder-operator'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE_BIND,
    MENU_MODEL_FIELD_TEMPLATE_SYNC_MODEL,
    MENU_MODEL_DETAILS
  } from '@/dictionary/menu-symbol'
  import ModelSyncStatus from './children/model-sync-status.vue'
  import useModelSyncStatus, { isSyncing, isSynced } from './children/use-model-sync-status'
  import { escapeRegexChar } from '@/utils/util'

  const props = defineProps({
    templateId: {
      type: Number
    }
  })

  const emit = defineEmits(['unbound', 'close'])

  const store = useStore()

  const requestIds = {
    rollList: Symbol('rollList')
  }
  const firstLoading = ref(true)

  const bindModelList = ref([])
  const searchName = ref('')
  const stuff = ref({ type: 'default', payload: { emptyText: t('bk.table.emptyText') } })

  const TABLE_ROW_HEIGHT = 43
  const tableMaxHeight = computed(() => store.state.appHeight - 272 - 50)
  const onePageMaxCount = computed(() => Math.ceil(tableMaxHeight.value / TABLE_ROW_HEIGHT))

  const activeModelIdList = computed(() => bindModelList.value
    .filter(model => !model.bk_ispaused)
    .map(model => model.id))

  const { statusMap, clear: clearModelSyncReq } = useModelSyncStatus(props.templateId, activeModelIdList)

  const scrollLadingOption = reactive({
    icon: '',
    size: 'small',
    theme: 'primary',
    isLoading: false
  })

  const queryParams = computed(() => {
    const params = {
      bk_template_id: props.templateId
    }

    if (searchName.value?.length) {
      stuff.value.type = 'search'
      params.filter = {
        condition: 'AND',
        rules: [{
          field: 'bk_obj_name',
          operator: queryBuilderOperator(QUERY_OPERATOR.LIKE),
          value: escapeRegexChar(searchName.value)
        }]
      }
    } else {
      stuff.value.type = 'default'
    }

    return params
  })

  let rollReq = null

  const getBindModel = debounce(async () => {
    rollReq = await fieldTemplateService.rollBindModel(
      queryParams.value,
      onePageMaxCount.value,
      {
        requestId: requestIds.rollList
      }
    )

    const { value, done } = rollReq.next()
    if (value) {
      firstLoading.value = true
      const result = await value()

      firstLoading.value = false
      bindModelList.value = result?.info ?? []
    } else if (done) {
      // 未找到数据
      bindModelList.value = []
      firstLoading.value = false
    }
  }, 200, { leading: true })

  watch([queryParams, onePageMaxCount], () => {
    getBindModel()
  }, { immediate: true })

  const handleScrollLoad = async () => {
    const { value, done } = rollReq.next()
    if (!done) {
      const result = await value()
      bindModelList.value.push(...result?.info ?? [])
    }
  }

  const handleGoBindModel = () => {
    emit('close')
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_BIND,
      params: {
        id: props.templateId
      }
    })
  }
  const handleGoSync = (row) => {
    routerActions.open({
      name: MENU_MODEL_FIELD_TEMPLATE_SYNC_MODEL,
      params: {
        id: props.templateId,
        modelId: row.id
      }
    })
  }
  const handleUnbind = (row) => {
    $bkInfo({
      type: 'warning',
      title: t('确认解绑该模板'),
      subTitle: t('解绑后，字段内容与唯一校验将会与模板脱离关系，不再受模板管理'),
      okText: t('解绑'),
      cancelText: t('取消'),
      confirmLoading: true,
      confirmFn: async () => {
        const params = {
          bk_template_id: props.templateId,
          object_id: row.id
        }
        await fieldTemplateService.unbind(params)
        $success(t('解绑成功'))
        getBindModel()
        emit('unbound')
        return true
      }
    })
  }

  const handleLinkToModel = (model) => {
    routerActions.open({
      name: MENU_MODEL_DETAILS,
      params: {
        modelId: model.bk_obj_id
      }
    })
  }
  const handleClearFilter = () => {
    searchName.value = ''
    stuff.value.type = 'default'
  }

  onUnmounted(() => {
    clearModelSyncReq()
  })
</script>

<template>
  <div class="details-model" v-bkloading="{ isLoading: firstLoading }">
    <div class="action-bar">
      <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :outline="true"
            :disabled="disabled"
            @click="handleGoBindModel">{{$t('绑定新模型')}}</bk-button>
        </template>
      </cmdb-auth>

      <bk-input
        class="search-input"
        v-model="searchName"
        :placeholder="$t('请输入模型名称')"
        :right-icon="'bk-icon icon-search'" />
    </div>
    <bk-table class="data-table"
      v-if="!firstLoading"
      :data="bindModelList"
      :max-height="tableMaxHeight"
      :scroll-loading="{ ...scrollLadingOption, isLoading: $loading(requestIds.rollList) }"
      @scroll-end="handleScrollLoad">
      <bk-table-column
        :label="$t('模型名称')"
        prop="bk_obj_name"
        show-overflow-tooltip>
        <template #default="{ row }">
          <div :class="['model-cell', { paused: row.bk_ispaused }]" @click="handleLinkToModel(row)">
            <div class="model-icon">
              <i class="icon" :class="row.bk_obj_icon"></i>
            </div>
            <span class="model-name">{{ row.bk_obj_name }}</span>
            <mini-tag theme="paused" v-if="row.bk_ispaused">{{ $t('已停用') }}</mini-tag>
          </div>
        </template>
      </bk-table-column>
      <bk-table-column
        :label="$t('同步状态')">
        <template #default="{ row }">
          <model-sync-status :model="row" />
        </template>
      </bk-table-column>
      <bk-table-column
        :label="$t('最近同步时间')"
        show-overflow-tooltip>
        <template #default="{ row }">
          <span v-if="statusMap[row.id] && !row.bk_ispaused">
            {{ statusMap[row.id].sync_time
              ? $tools.formatTime(statusMap[row.id].sync_time, 'YYYY-MM-DD HH:mm:ssZZ')
              : '--'
            }}
          </span>
          <span v-else>--</span>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')">
        <template #default="{ row }">
          <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }"
            v-bk-tooltips="{
              content: $t('操作正在进行中，无法解除'),
              disabled: row.bk_ispaused || (statusMap[row.id] && !isSyncing(statusMap[row.id].status))
            }">
            <template #default="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled || (!row.bk_ispaused && statusMap[row.id] && isSyncing(statusMap[row.id].status))"
                :text="true"
                @click.stop="handleUnbind(row)">
                {{$t('解除绑定')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth class="ml20"
            :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [templateId] }"
            v-if="!row.bk_ispaused && (statusMap[row.id] && !isSynced(statusMap[row.id].status))">
            <template #default="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleGoSync(row)">
                {{$t('去同步')}}
              </bk-button>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="stuff"
        @clear="handleClearFilter"
      ></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<style lang="scss" scoped>
.details-model {
  padding: 16px 8px;

  .action-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;

    .search-input {
      width: 580px;
    }
  }

  .data-table {
    .model-cell {
      display: flex;
      align-items: center;
      gap: 10px;
      cursor: pointer;
      .model-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 24px;
        height: 24px;
        border-radius: 50%;
        background-color: #F0F5FF;
        .icon {
          color: #3a84ff;
          font-size: 12px;
        }
      }
      .model-name {
        color: #63656E;
      }

      &.paused {
        .model-name {
          color: #C4C6CC;
        }
        .tag {
          background: #F5F7FA;
        }
      }
    }

    .action-button {
      :deep(.bk-link-text) {
        font-size: 12px;
      }
    }
  }
}
</style>
