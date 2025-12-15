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
  import { ref, computed, reactive, watch, set, nextTick } from 'vue'
  import debounce from 'lodash/debounce'
  import { t } from '@/i18n'
  import RouterQuery from '@/router/query'
  import routerActions from '@/router/actions'
  import { OPERATION } from '@/dictionary/iam-auth'
  import {
    MENU_MODEL_FIELD_TEMPLATE_CREATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT,
    MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS,
    MENU_MODEL_FIELD_TEMPLATE_BIND,
    MENU_MODEL_DETAILS
  } from '@/dictionary/menu-symbol'
  import { $success, $error, $bkPopover } from '@/magicbox/index.js'
  import applyPermission from '@/utils/apply-permission.js'
  import { getDefaultPaginationConfig, getSort } from '@/utils/tools.js'
  import queryBuilderOperator, { QUERY_OPERATOR } from '@/utils/query-builder-operator'
  import fieldTemplateService from '@/service/field-template'
  import CmdbLoading from '@/components/loading/index.vue'
  import MiniTag from '@/components/ui/other/mini-tag.vue'
  import TemplateDetails from './details.vue'
  import ModelSyncStatus from './children/model-sync-status.vue'
  import SearchSelect from './children/search-select.vue'
  import CloneDialog from './children/clone-dialog.vue'
  import useModelSyncStatus from './children/use-model-sync-status'
  import useTemplate from './children/use-template'
  import { escapeRegexChar } from '@/utils/util'

  const requestIds = {
    list: Symbol('list'),
    fieldCount: Symbol('fieldCount'),
    modelCount: Symbol('modelCount'),
    modelList: Symbol('modelList')
  }

  // 响应式的query
  const query = computed(() => RouterQuery.getAll())

  const table = reactive({
    header: [],
    list: [],
    pagination: {
      count: 0,
      current: 1,
      ...getDefaultPaginationConfig()
    },
    sort: '',
    stuff: {
      type: 'default',
      payload: {
        resource: t('字段组合模板')
      }
    }
  })

  const templateList = computed(() => table.list)
  const { handleDelete: handleDeleteTemplate } = useTemplate(templateList)

  const isShowCloneDialog = ref(false)
  const cloneSourceTemplate = ref({})

  const boundModelPopoverContentRef = ref(null)
  const boundModelPopover = reactive({
    templateId: null,
    instance: null,
    modelList: [],
    show: false,
    clear: null
  })

  const filter = ref([])
  // 计算查询条件参数
  const searchParams = computed(() => {
    const params = {
      template_filter: {
        condition: 'AND',
        rules: []
      },
      object_filter: {
        condition: 'AND',
        rules: []
      },
      fields: [],
      page: {
        start: table.pagination.limit * (table.pagination.current - 1),
        limit: table.pagination.limit,
        sort: table.sort || '-last_time' // 默认按照最新更新时间倒序
      }
    }
    filter.value.forEach((item) => {
      const itemValue = item.value?.split(',')
      const operator = queryBuilderOperator(itemValue?.length > 1 ? QUERY_OPERATOR.IN : QUERY_OPERATOR.LIKE)
      const value = itemValue?.length > 1 ? itemValue : escapeRegexChar(itemValue[0])
      if (item.id === 'templateName') {
        params.template_filter.rules.push({
          field: 'name',
          operator,
          value
        })
      }
      if (item.id === 'modifier') {
        params.template_filter.rules.push({
          field: 'modifier',
          operator,
          value
        })
      }
      if (item.id === 'modelName') {
        params.object_filter.rules.push({
          field: 'bk_obj_name',
          operator,
          value
        })
      }
    })

    if (!params.template_filter.rules.length) {
      Reflect.deleteProperty(params, 'template_filter')
    }
    if (!params.object_filter.rules.length) {
      Reflect.deleteProperty(params, 'object_filter')
    }

    return params
  })

  const getList = async (options = {}) => {
    try {
      const { list, count } = await fieldTemplateService.find(searchParams.value, {
        requestId: requestIds.list,
        cancelPrevious: true,
        globalPermission: false
      })

      if (options.isDel && count && !list?.length) {
        RouterQuery.set({
          page: table.pagination.current - 1,
          _t: Date.now()
        })
        return
      }

      table.list = list
      table.pagination.count = count

      table.stuff.type = filter.value.length ? 'search' : 'default'
    } catch ({ permission }) {
      if (permission) {
        table.stuff = {
          type: 'permission',
          payload: { permission }
        }
      }
    }
  }

  const getFieldCount = async (templateIds) => {
    const data = await fieldTemplateService.getFieldCount({
      bk_template_ids: templateIds
    }, { requestId: requestIds.fieldCount })
    data.forEach((item) => {
      const row = table.list.find(row => row.id === item.bk_template_id)
      set(row, 'field_count', item.count)
    })
  }
  const getModelCount = async (templateIds) => {
    const data = await fieldTemplateService.getModelCount({
      bk_template_ids: templateIds
    }, { requestId: requestIds.modelCount })
    data.forEach((item) => {
      const row = table.list.find(row => row.id === item.bk_template_id)
      set(row, 'model_count', item.count)
    })
  }

  const detailsDrawer = reactive({
    open: false,
    template: {}
  })
  const openDetailsDrawer = (template) => {
    detailsDrawer.template  = template
    nextTick(() => {
      detailsDrawer.open = true
    })
  }

  // 监听查询参数触发查询
  watch(
    query,
    async (query) => {
      const {
        page = 1,
        limit = table.pagination.limit,
        action = '',
        id = '',
        templateName = '',
        modelName = '',
        modifier = '',
        _t = ''
      } = query

      // 在本页所做操作必须带上 _t 参数以此与外入动作区分
      if (action === 'view' && id && !_t) {
        const template = await fieldTemplateService.findById(Number(id))
        if (!template) {
          $error('no data found!')
          return
        }
        openDetailsDrawer(template)
      }

      table.pagination.current = parseInt(page, 10)
      table.pagination.limit = parseInt(limit, 10)

      const queryFilter = [
        { id: 'templateName', value: templateName },
        { id: 'modelName', value: modelName },
        { id: 'modifier', value: modifier },
      ]
      filter.value = queryFilter.filter(item => item.value.length)

      getList()
    },
    {
      immediate: true
    }
  )

  watch(() => table.list, async (list) => {
    const templateIds = list.map(item => item.id)
    if (templateIds.length) {
      getFieldCount(templateIds)
      getModelCount(templateIds)
    }
  })

  const handleBoundModelCountMouseenter = debounce(async (event, template) => {
    // 防止抖动，一个模板当前只执行一次
    if (template.id === boundModelPopover.templateId) {
      boundModelPopover.instance?.show()
      return
    }

    boundModelPopover.templateId = template.id
    boundModelPopover.instance?.destroy?.()
    boundModelPopover.clear?.()
    boundModelPopover.modelList = []

    boundModelPopover.instance = $bkPopover(event.target, {
      content: boundModelPopoverContentRef.value,
      delay: [300, 0],
      hideOnClick: true,
      interactive: true,
      placement: 'top',
      animateFill: false,
      sticky: true,
      theme: 'light bound-model-popover',
      boundary: 'window',
      trigger: 'mouseenter', // 'manual mouseenter',
      arrow: true,
      onShow: () => {
        boundModelPopover.show = true
      },
      onHidden: () => {
        boundModelPopover.show = false
        boundModelPopover.clear?.()
      }
    })

    boundModelPopover.instance.show()

    const modelList = await fieldTemplateService.getBindModel({
      bk_template_id: template.id
    }, { requestId: requestIds.modelList })
    boundModelPopover.modelList = modelList ?? []

    const activeModelIdList = modelList
      .filter(model => !model.bk_ispaused)
      .map(model => model.id)

    const { clear } = useModelSyncStatus(template.id, activeModelIdList)

    boundModelPopover.clear = clear
  }, 300)

  const handleCreate = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_CREATE,
      history: true
    })
  }
  const handleEdit = (id, type = 'edit') => {
    const name = type === 'clone' ? MENU_MODEL_FIELD_TEMPLATE_EDIT_FIELD_SETTINGS : MENU_MODEL_FIELD_TEMPLATE_EDIT
    routerActions.redirect({
      name,
      params: {
        id
      },
      history: false
    })
  }
  const handleBind = (id) => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_BIND,
      params: {
        id
      },
      history: true
    })
  }

  const handleSortChange = (sort) => {
    table.sort = getSort(sort)
    RouterQuery.refresh()
  }

  const handleSizeChange = (size) => {
    RouterQuery.set({
      limit: size,
      page: 1,
      _t: Date.now()
    })
  }

  const handlePageChange = (page) => {
    RouterQuery.set({
      page,
      _t: Date.now()
    })
  }

  const handleClone = (template) => {
    isShowCloneDialog.value = true
    cloneSourceTemplate.value = template
  }

  const handleCloneSuccess = (res) => {
    const { id } = res
    if (!id) return
    getList()
    $success(t('克隆成功'))
    isShowCloneDialog.value = false
    handleEdit(id, 'clone')
  }
  const handleCloneDone = () => {
    RouterQuery.refresh()
  }
  const handleCloneDialogToggle = (val) => {
    isShowCloneDialog.value = val
  }

  const handleDelete = (row) => {
    handleDeleteTemplate(row, () => {
      getList({ isDel: true })
      $success(t('删除成功'))
    })
  }
  const handleDeleteDone = () => {
    getList({ isDel: true })
  }

  const handleRowIDClick = (row) => {
    openDetailsDrawer(row)
  }
  const handleDetailsDrawerClose = () => {
    detailsDrawer.open = false
  }

  const handleClearFilter = () => {
    RouterQuery.clear()
  }

  const handleApplyPermission = async () => {
    try {
      await applyPermission({
        type: OPERATION.R_FIELD_TEMPLATE,
        relation: []
      })
    } catch (e) {
      console.error(e)
    }
  }
  const handleLinkToModel = (model) => {
    routerActions.open({
      name: MENU_MODEL_DETAILS,
      params: {
        modelId: model.bk_obj_id
      }
    })
  }
  const handleSearch = (filter) => {
    const query = {
      templateName: '',
      modelName: '',
      modifier: '',
      _t: Date.now(),
      page: 1
    }
    filter.forEach((item) => {
      const key = item.id === 'modifier' ? 'username' : 'name'
      query[item.id] = item.values.map(val => val[key]).join(',')
    })
    RouterQuery.set(query)
  }
  const handleBindModelChange = (templateId) => {
    getModelCount([templateId])
  }

  const handleUpdateTemplate = (id, val, dataKey) => {
    const row = table.list.find(row => row.id === id)
    if (row) {
      set(row, dataKey, val)
    }
  }
</script>

<template>
  <div class="field-template">
    <cmdb-tips class="mb10" tips-key="fieldTemplateTips">
      {{$t('字段组合模板功能提示')}}
    </cmdb-tips>
    <div class="toolbar">
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <span v-bk-tooltips.top="{
            content: $t('创建字段组合模板上限提示'),
            disabled: !disabled && table.pagination.count < 10000
          }">
            <bk-button
              theme="primary"
              :disabled="disabled"
              @click="handleCreate">
              {{$t('新建')}}
            </bk-button>
          </span>
        </template>
      </cmdb-auth>
      <search-select
        :default-filter="filter"
        @search="handleSearch">
      </search-select>
    </div>

    <bk-table class="data-table"
      v-bkloading="{ isLoading: $loading(requestIds.list) }"
      :data="table.list"
      :pagination="table.pagination"
      :max-height="$APP.height - 229"
      @sort-change="handleSortChange"
      @page-limit-change="handleSizeChange"
      @page-change="handlePageChange">
      <bk-table-column
        sortable="custom"
        prop="name"
        :label="$t('模板名称')"
        fixed="left"
        min-width="230"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <cmdb-auth tag="div" class="template-name-auth"
            :auth="{ type: $OPERATION.R_FIELD_TEMPLATE, relation: [row.id] }">
            <template #default="{ disabled }">
              <div :class="['cell-link-content', { disabled }]" @click.stop="handleRowIDClick(row)">
                {{ row.name }}
              </div>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <bk-table-column
        :label="$t('字段数量')"
        min-width="130"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{ row.field_count }}</div>
        </template>
      </bk-table-column>
      <bk-table-column
        min-width="150"
        :label="$t('绑定的模型')">
        <template slot-scope="{ row }">
          <cmdb-loading :loading="$loading(requestIds.modelCount)">
            <span v-if="row.model_count === 0" class="cell-unbind">0（{{$t('未绑定')}}）</span>
            <cmdb-auth v-else tag="div"
              :auth="{ type: $OPERATION.R_FIELD_TEMPLATE, relation: [row.id] }">
              <span class="cell-link-content"
                @mouseenter="(event) => handleBoundModelCountMouseenter(event, row)">
                {{ row.model_count ?? '--' }}
              </span>
            </cmdb-auth>
          </cmdb-loading>
        </template>
      </bk-table-column>
      <bk-table-column
        prop="description"
        :label="$t('描述')"
        min-width="290"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{ row.description || '--' }}</div>
        </template>
      </bk-table-column>
      <bk-table-column
        sortable="custom"
        prop="modifier"
        min-width="170"
        :label="$t('最近更新人')"
        show-overflow-tooltip>
        <template #default="{ row }">
          <cmdb-user-value :value="row.modifier" />
        </template>
      </bk-table-column>
      <bk-table-column
        sortable="custom"
        prop="last_time"
        :label="$t('最近更新时间')"
        min-width="190"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{row.last_time ? $tools.formatTime(row.last_time) : '--'}}</div>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" min-width="250" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [row.id] }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleBind(row.id)">
                {{$t('绑定新模型')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [row.id] }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleEdit(row.id)">
                {{$t('编辑')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth class="mr10" :ignore-passed-auth="true" :auth="[
            { type: $OPERATION.C_FIELD_TEMPLATE },
            { type: $OPERATION.R_FIELD_TEMPLATE, relation: [row.id] }
          ]">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleClone(row)">
                {{$t('克隆')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth :auth="{ type: $OPERATION.D_FIELD_TEMPLATE, relation: [row.id] }">
            <template #default="{ disabled }">
              <div v-bk-tooltips.top="{ content: $t('已被模型绑定，不能删除'), disabled: !row.model_count }">
                <bk-button
                  theme="primary"
                  :disabled="disabled || row.model_count > 0"
                  :text="true"
                  @click.stop="handleDelete(row)">
                  {{$t('删除')}}
                </bk-button>
              </div>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff"
        @clear="handleClearFilter">
        <bk-exception type="403" scene="part">
          <i18n path="字段组合模板列表提示语" class="table-empty-tips">
            <template #auth>
              <bk-link theme="primary"
                @click="handleApplyPermission">{{$t('申请查看权限')}}</bk-link>
            </template>
            <template #create>
              <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
                <bk-button slot-scope="{ disabled }" text
                  theme="primary"
                  class="text-btn"
                  :disabled="disabled"
                  @click="handleCreate">
                  {{$t('立即创建')}}
                </bk-button>
              </cmdb-auth>
            </template>
          </i18n>
        </bk-exception>
      </cmdb-table-empty>
    </bk-table>

    <template-details
      v-if="detailsDrawer.template.id"
      :open="detailsDrawer.open"
      :template="detailsDrawer.template"
      @bind-change="handleBindModelChange"
      @clone-done="handleCloneDone"
      @delete-done="handleDeleteDone"
      @update-template="handleUpdateTemplate"
      @close="handleDetailsDrawerClose">
    </template-details>

    <clone-dialog
      :show="isShowCloneDialog"
      :source-template="cloneSourceTemplate"
      @success="handleCloneSuccess"
      @toggle="handleCloneDialogToggle">
    </clone-dialog>

    <div ref="boundModelPopoverContentRef" v-show="boundModelPopover.show">
      <div class="bound-model-popover-content"
        v-bkloading="{ isLoading: $loading(requestIds.modelList), theme: 'primary', mode: 'spin', size: 'mini' }">
        <div class="bound-model"
          v-for="(item, index) in boundModelPopover.modelList"
          :key="`${item.id}_${index}`">
          <model-sync-status :model="item" :mini="true" class="status-icon" />
          <div :class="['model-item', { paused: item.bk_ispaused }]">
            <div class="model-icon">
              <i class="icon" :class="item.bk_obj_icon"></i>
            </div>
            <span class="model-name" :title="item.bk_obj_name">{{ item.bk_obj_name }}</span>
            <mini-tag theme="paused" v-if="item.bk_ispaused">{{ $t('已停用') }}</mini-tag>
          </div>
          <i class="link-icon icon-cc-share"
            @click="handleLinkToModel(item)">
          </i>
        </div>
      </div>
    </div>
  </div>
</template>
<style lang="scss" scoped>
  .field-template {
    padding: 15px 20px 0;

    .toolbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .data-table {
      margin-top: 14px;
    }

    .template-name-auth {
      width: 100%;
    }
    .cell-link-content {
      color: $primaryColor;
      cursor: pointer;
      @include ellipsis;

      &.disabled {
        color: #a3c5fd;
      }
    }
    .cell-unbind {
      color: #FF9C01;
    }

    .table-empty-tips {
      display: flex;
      align-items: center;
      justify-content: center;
      .text-btn {
        font-size: 14px;
      }
    }
  }

  .bound-model-popover-content {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 96px;
    min-height: 38px;
    max-height: 220px;
    width: 320px;
    padding: 12px 14px 12px 10px;
    @include scrollbar-y;

    .bound-model {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .model-item {
      display: flex;
      align-items: center;
      gap: 8px;
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
        flex: none;
        font-size: 12px;
        color: #63656E;
        max-width: 150px;
        @include ellipsis;
      }

      &.paused {
        .model-name {
          color: #C4C6CC;
        }
      }
    }
    .link-icon {
      font-size: 12px;
      color: $primaryColor;
      cursor: pointer;
      margin-left: auto;
      &:hover {
        opacity: .75;
      }
    }

    :deep(.status-icon) {
      min-width: 16px;
    }
  }
</style>
<style>
  .tippy-tooltip.bound-model-popover-theme {
    padding-left: 2px !important;
    padding-right: 2px !important;
  }
</style>
