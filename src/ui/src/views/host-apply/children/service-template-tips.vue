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

<template>
  <div class="service-template-tips">
    <bk-alert type="warning" class="tips-alert">
      <template slot="title">
        <i18n path="服务模板的自动应用配置将覆盖其关联模块已有的规则操作提示">
          <template #link>
            <bk-link theme="primary" class="view-link" @click="handleShow">{{$t('点击查看')}}</bk-link>
          </template>
        </i18n>
      </template>
    </bk-alert>
    <bk-dialog
      v-model="isShow"
      class="module-modal"
      header-position="left"
      :width="850"
      :mask-close="false"
      :draggable="false">
      <template #header>
        <div class="modal-header">
          <div class="title">{{$t('关联模块')}}</div>
          <bk-select
            v-model="currentTemplateId"
            ext-cls="template-selector"
            :clearable="false"
            searchable>
            <bk-option v-for="option in templateList"
              :key="option.id"
              :id="option.id"
              :name="option.name">
            </bk-option>
          </bk-select>
        </div>
      </template>

      <div class="module-list">
        <div class="searchbar">
          <i18n path="共N个关联模块" class="count">
            <template #count><em class="num">{{displayModuleList.length}}</em></template>
          </i18n>
          <bk-input
            class="search-input"
            type="text"
            :placeholder="$t('请输入业务拓扑搜索')"
            clearable
            right-icon="bk-icon icon-search"
            @input="handleFilter">
          </bk-input>
        </div>
        <bk-table class="module-list-table"
          :cell-class-name="cellClassName"
          ref="table"
          v-bkloading="{ isLoading: $loading(Object.values(requestIds)) }"
          :pagination="pagination"
          :data="displayModuleList"
          height="500"
          max-height="500"
          @page-change="handlePageChange"
          @page-limit-change="handlePageLimitChange">
          <bk-table-column
            show-overflow-tooltip
            :width="450"
            :label="$t('业务拓扑')"
            prop="topopath">
          </bk-table-column>
          <bk-table-column
            :width="220"
            :label="$t('是否配置自动应用')">
            <template slot-scope="{ row }">
              <bk-icon v-if="row.host_apply_enabled" class="enabled-icon" type="check-circle-shape" />
              <span v-else>--</span>
            </template>
          </bk-table-column>
          <bk-table-column
            :width="90"
            :label="$t('操作')">
            <template slot-scope="{ row }">
              <bk-button theme="primary" text @click="handleViewModuleHostApply(row)" v-if="row.host_apply_enabled">
                {{$t('跳转查看')}}
              </bk-button>
              <span v-else>--</span>
            </template>
          </bk-table-column>
          <bk-exception slot="empty" :type="isSearch ? 'search-empty' : 'empty'" scene="part">
            <span>{{$t('无模块')}}</span>
          </bk-exception>
        </bk-table>
      </div>

      <template #footer>
        <bk-button @click="handleClose">{{$t('关闭按钮')}}</bk-button>
      </template>
    </bk-dialog>
  </div>
</template>

<script>
  import { defineComponent, reactive, toRefs, watch, watchEffect } from 'vue'
  import store from '@/store'
  import routerActions from '@/router/actions'
  import { getDefaultPaginationConfig } from '@/utils/tools.js'
  import serviceTemplateService, { CONFIG_MODE } from '@/service/service-template/index.js'
  import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'

  export default defineComponent({
    props: {
      serviceTemplateIds: {
        type: Array,
        default: () => ([]),
        required: true
      }
    },
    setup(props) {
      const { serviceTemplateIds } = toRefs(props)
      const bizId = store.getters['objectBiz/bizId']

      const requestIds = {
        modules: Symbol(),
        topopath: Symbol()
      }

      const state = reactive({
        isShow: false,
        isSearch: false,
        templateList: [],
        currentTemplateId: null,
        pagination: getDefaultPaginationConfig({}, false),
        currentModuleList: [],
        displayModuleList: []
      })

      watchEffect(async () => {
        const list = await serviceTemplateService.findAllByIds(serviceTemplateIds.value, { bk_biz_id: bizId })
        state.templateList = list

        const [firstTemplate] = list
        state.currentTemplateId = firstTemplate.id
        state.displayModuleList = await getModulesById(state.currentTemplateId)
      })

      watch(() => state.currentTemplateId, async (id) => {
        state.displayModuleList = await getModulesById(id)
      })

      const getModulesById = async (id) => {
        const { info: modules } = await store.dispatch('serviceTemplate/getServiceTemplateModules', {
          bizId,
          serviceTemplateId: id,
          params: {
            fields: ['bk_module_id', 'host_apply_enabled']
          },
          config: {
            requestId: requestIds.modules
          }
        })

        if (!modules?.length) {
          return []
        }

        // 获取模块的拓扑路径
        const topopath = await store.dispatch('objectMainLineModule/getTopoPath', {
          bizId,
          params: {
            topo_nodes: modules.map(module => ({ bk_obj_id: 'module', bk_inst_id: module.bk_module_id }))
          },
          config: {
            requestId: requestIds.topopath
          }
        })

        // 将拓扑路径加入到模块数据中
        modules.forEach((module) => {
          const moduleNode = topopath.nodes.find(node => node.topo_node.bk_inst_id === module.bk_module_id)
          module.topopath = moduleNode.topo_path.map(path => path.bk_inst_name).reverse()
            .join(' / ')
        })

        return modules
      }

      const handleFilter = (keyword) => {
        // 备份一次全数据
        if (!state.currentModuleList.length) {
          state.currentModuleList = state.displayModuleList.slice()
        }

        if (!keyword?.length) {
          state.displayModuleList = state.currentModuleList.slice()
          return
        }

        state.displayModuleList = state.currentModuleList.filter(item => new RegExp(keyword, 'i').test(item.topopath.replaceAll(' / ', '/')))
      }

      const handleViewModuleHostApply = (module) => {
        routerActions.open({
          name: MENU_BUSINESS_HOST_APPLY,
          params: {
            mode: CONFIG_MODE.MODULE
          },
          query: { id: module.bk_module_id }
        })
      }

      const handlePageChange = (page) => {
        state.pagination.current = page
      }

      const handlePageLimitChange = (limit) => {
        state.pagination.current = 1
        state.pagination.limit = limit
      }

      const handleShow = () => {
        state.isShow = true
      }
      const handleClose = () => {
        state.isShow = false
      }

      const cellClassName = ({ column }) => {
        if (column?.property) {
          return column.property
        }
      }

      return {
        ...toRefs(state),
        requestIds,
        handleFilter,
        handlePageChange,
        handlePageLimitChange,
        handleViewModuleHostApply,
        handleShow,
        handleClose,
        cellClassName
      }
    }
  })
</script>

<style lang="scss" scoped>
  .tips-alert {
    .view-link {
      margin-top: -2px;
      ::v-deep .bk-link-text {
        font-size: 12px;
      }
    }
  }

  .module-list {
    .searchbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 8px;

      .count {
        color: #63656E;
        .num {
          margin: 0 .1em;
          font-weight: 700;
          font-style: normal;
        }
      }

      .search-input {
        width: 320px;
      }
    }
  }

  .module-list-table {
    .enabled-icon {
      font-size: 14px !important;
      color: #2DCB56;
    }
    .empty {
      ::v-deep {
        .bk-exception-img {
          width: 160px;
          height: 140px;
        }
        .bk-exception-text {
          font-size: 12px;
        }
      }
    }

    ::v-deep {
      .topopath .cell {
        direction: rtl;
        text-align: left;
        white-space: nowrap;
        display: block;
      }
    }
  }

  .module-modal {
    ::v-deep {
      .bk-dialog-header {
        padding-bottom: 12px;
      }
    }
    .modal-header {
      display: flex;
      align-items: center;

      .title {
        color: #313238;
      }

      .template-selector {
        width: 250px;
        margin-left: 10px;
        border: none;
        background: #f4f6fa;

        &.is-focus {
          box-shadow: none;
        }
      }
    }
  }
</style>
