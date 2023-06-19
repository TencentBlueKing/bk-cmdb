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
  import { ref, computed, reactive, watch, set, nextTick } from 'vue'
  import { t } from '@/i18n'
  import RouterQuery from '@/router/query'
  import routerActions from '@/router/actions'
  import { OPERATION } from '@/dictionary/iam-auth'
  import {
    MENU_MODEL_FIELD_TEMPLATE_CREATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT,
    MENU_MODEL_FIELD_TEMPLATE_BIND
  } from '@/dictionary/menu-symbol'
  import { $bkInfo, $success, $error } from '@/magicbox/index.js'
  import applyPermission from '@/utils/apply-permission.js'
  import { getDefaultPaginationConfig, getSort } from '@/utils/tools.js'
  import fieldTemplateService from '@/service/field-template'
  import CmdbLoading from '@/components/loading/index.vue'
  import TemplateDetails from './details.vue'

  const requestIds = {
    list: Symbol('list'),
    fieldCount: Symbol('fieldCount'),
    modelCount: Symbol('modelCount'),
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

  const filterOptions = [
    {
      id: 'template_name',
      name: t('模板名称')
    },
    {
      id: 'model_name',
      name: t('模型名称')
    },
    {
      id: 'modifier',
      name: t('更新人')
    }
  ]
  const filter = ref([])

  // 计算查询条件参数
  const searchParams = computed(() => {
    const params = {
      // template_filter: {},
      // object_filter: {},
      fields: [],
      page: {
        start: table.pagination.limit * (table.pagination.current - 1),
        limit: table.pagination.limit,
        sort: table.sort
      }
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

      // table.stuff.type = filter.value.toString().length ? 'search' : 'default'
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

  const updateFilter = () => {}

  // 监听查询参数触发查询
  watch(
    query,
    async (query) => {
      const {
        page = 1,
        limit = table.pagination.limit,
        action = '',
        id = '',
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

      updateFilter()

      table.pagination.current = parseInt(page, 10)
      table.pagination.limit = parseInt(limit, 10)

      getList()
    },
    {
      immediate: true
    }
  )

  watch(() => table.list, async (list) => {
    const templateIds = list.map(item => item.id)
    getFieldCount(templateIds)
    getModelCount(templateIds)
  })

  const handleCreate = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_CREATE,
      history: true
    })
  }
  const handleEdit = (id) => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_EDIT,
      params: {
        id
      },
      history: true
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
    table.sort = getSort(sort, { prop: 'name', order: 'descending' })
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

  const handleDelete = (row) => {
    $bkInfo({
      title: t('确认要删除', { name: row.name }),
      confirmLoading: true,
      confirmFn: async () => {
        try {
          // TODO: 删除接口
          getList({ isDel: true })
          $success(t('删除成功'))
        } catch (error) {
          console.error(error)
          $error(t('删除失败'))
          return false
        }
      }
    })
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
</script>

<template>
  <div class="field-template">
    <cmdb-tips class="mb10" tips-key="fieldTemplateTips">
      {{$t('字段组合模板功能提示')}}
    </cmdb-tips>
    <div class="toolbar">
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleCreate">
            {{$t('新建')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <bk-search-select
        class="search-select"
        clearable
        :placeholder="$t('请输入模板名称/模型/更新人')"
        :show-popover-tag-change="true"
        :data="filterOptions"
        v-model="filter">
      </bk-search-select>
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
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div class="cell-link-content" @click.stop="handleRowIDClick(row)">{{ row.name }}</div>
        </template>
      </bk-table-column>
      <bk-table-column
        prop="name"
        :label="$t('字段数量')"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{ row.field_count }}</div>
        </template>
      </bk-table-column>
      <bk-table-column
        prop="name"
        :label="$t('绑定的模型')"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <cmdb-loading :loading="$loading(requestIds.modelCount)">
            <span v-if="row.model_count === 0" class="cell-unbind">0（{{$t('未绑定')}}）</span>
            <span v-else class="cell-link-content">{{ row.model_count ?? '--' }}</span>
          </cmdb-loading>
        </template>
      </bk-table-column>
      <bk-table-column
        prop="name"
        :label="$t('描述')"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{ row.description || '--' }}</div>
        </template>
      </bk-table-column>
      <bk-table-column
        sortable="custom"
        prop="name"
        :label="$t('最近更新人')"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{ row.modifier }}</div>
        </template>
      </bk-table-column>
      <bk-table-column
        sortable="custom"
        prop="name"
        :label="$t('最近更新时间')"
        show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div>{{ row.last_time }}</div>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_FIELD_TEMPLATE, relation: [row.id] }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleBind(row.id)">
                {{$t('绑定模型')}}
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
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.D_BUSINESS_SET, relation: [row.bk_biz_set_id] }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleCopy(row)">
                {{$t('克隆')}}
              </bk-button>
            </template>
          </cmdb-auth>
          <cmdb-auth
            :auth="{ type: $OPERATION.D_BUSINESS_SET, relation: [row.bk_biz_set_id] }"
            v-bk-tooltips.top="{ content: $t('已被模型绑定，不能删除'), disabled: false }">
            <template slot-scope="{ disabled }">
              <bk-button
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleDelete(row)">
                {{$t('删除')}}
              </bk-button>
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
      @close="handleDetailsDrawerClose">
    </template-details>
  </div>
</template>
<style lang="scss" scoped>
  .field-template {
    padding: 15px 20px 0;

    .toolbar {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .search-select {
        width: 480px;
      }
    }

    .data-table {
      margin-top: 14px;
    }

    .cell-link-content {
      color: $primaryColor;
      cursor: pointer;
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
</style>
