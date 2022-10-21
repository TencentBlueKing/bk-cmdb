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
  <div class="template-instance-layout">
    <div class="instance-main">
      <div class="options clearfix">
        <div class="fl">
          <cmdb-auth :auth="syncAuth">
            <template #default="{ disabled }">
              <bk-button theme="primary"
                :disabled="!table.selection.length || disabled"
                @click="handleBatchSync">
                {{$t('批量同步')}}
              </bk-button>
            </template>
          </cmdb-auth>
        </div>
        <div class="fr">
          <bk-input class="filter-item" right-icon="bk-icon icon-search"
            clearable
            :placeholder="$t('请输入拓扑路径')"
            v-model.trim="table.filter"
            @change="handleFilter">
          </bk-input>
        </div>
      </div>
      <bk-table class="instance-table"
        ref="instanceTable"
        v-bkloading="{ isLoading: instancesLoading || table.filtering }"
        :data="table.visibleList"
        :pagination="table.pagination"
        :max-height="$APP.height - 250"
        @page-limit-change="handleSizeChange"
        @page-change="handlePageChange">
        <BatchSelectionColumn
          width="60"
          :data="table.visibleList"
          :full-data="table.data"
          ref="batchSelectionColumn"
          row-key="bk_module_id"
          :selectable="row => !isSyncDisabled(row.status)"
          @selection-change="handleSelectionChange">
        </BatchSelectionColumn>
        <bk-table-column :label="$t('模块名称')" prop="bk_module_name" show-overflow-tooltip></bk-table-column>
        <bk-table-column :label="$t('拓扑路径')" sortable :sort-method="sortByPath" show-overflow-tooltip>
          <span slot-scope="{ row }" class="topo-path" @click="handlePathClick(row)">{{row.topoPath}}</span>
        </bk-table-column>
        <InstanceStatusColumn></InstanceStatusColumn>
        <bk-table-column :label="$t('上次同步时间')" sortable :sort-method="sortByTime" show-overflow-tooltip>
          <template slot-scope="{ row }">{{row.last_time | time}}</template>
        </bk-table-column>
        <bk-table-column :label="$t('操作')">
          <template slot-scope="{ row }">
            <cmdb-auth
              :auth="syncAuth"
              v-bk-tooltips="{
                content: isSyncing(row.status) ? $t('实例正在同步，无法操作') : $t('已是最新内容，无需同步'),
                disabled: !isSyncDisabled(row.status)
              }">
              <template #default="{ disabled }">
                <bk-button text
                  :disabled="isSyncDisabled(row.status) || disabled"
                  @click="handleSync(row)">
                  {{$t('去同步')}}
                </bk-button>
              </template>
            </cmdb-auth>
            <cmdb-auth class="ml10"
              :auth="delAuth"
              v-bk-tooltips="{
                content: row.hostCount ? $t('目标包含主机，不允许删除') : $t('由集群模板创建的模块无法删除'),
                disabled: !isDelDisabled(row)
              }">
              <template #default="{ disabled }">
                <bk-button text
                  :disabled="isDelDisabled(row) || disabled"
                  @click="handleDel(row)">
                  {{$t('删除')}}
                </bk-button>
              </template>
            </cmdb-auth>
          </template>
        </bk-table-column>
        <cmdb-table-empty slot="empty" :stuff="table.stuff">
          <div>
            <i18n path="空服务模板实例提示" tag="div">
              <template #link>
                <bk-button style="font-size: 14px;" text
                  @click="handleToCreatedInstance">
                  {{$t('创建服务实例')}}
                </bk-button>
              </template>
            </i18n>
          </div>
        </cmdb-table-empty>
      </bk-table>
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { time } from '@/filters/formatter'
  import debounce from 'lodash.debounce'
  import Bus from '@/utils/bus'
  import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
  import InstanceStatusColumn from './instance-status-column.vue'
  import { Polling } from '@/utils/polling'
  import BatchSelectionColumn from '@/components/batch-selection-column'
  import to from 'await-to-js'
  import cloneDeep from 'lodash/cloneDeep'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default {
    name: 'templateInstance',
    components: {
      BatchSelectionColumn,
      InstanceStatusColumn,
    },
    filters: {
      time
    },
    props: {
      active: Boolean,
    },
    data() {
      this.polling = new Polling(() => {
        const syncingModules = this.table.data
          .filter(theModule => this.isSyncing(theModule.status)).map(theModule => theModule.bk_module_id)
        if (syncingModules.length > 0) {
          return this.loadInstanceStatus(syncingModules)
        }
      }, 5000)
      return {
        table: {
          filter: '',
          filtering: false,
          selection: [],
          data: [],
          visibleList: [],
          backup: [],
          stuff: {
            type: 'default',
            payload: {}
          },
          pagination: {
            current: 1,
            count: 0,
            ...this.$tools.getDefaultPaginationConfig({
              'limit-list': [20, 50, 100]
            })
          }
        },
        handleFilter: null,
        instancesLoading: false,
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      serviceTemplateId() {
        return this.$route.params.templateId
      },
      syncAuth() {
        return {
          type: this.$OPERATION.U_SERVICE_TEMPLATE,
          relation: [this.bizId, Number(this.serviceTemplateId)]
        }
      },
      delAuth() {
        return {
          type: this.$OPERATION.D_TOPO,
          relation: [this.bizId]
        }
      }
    },
    watch: {
      active: {
        immediate: true,
        handler(active) {
          if (active) {
            this.loadIntegrationInstances()
            this.polling.start()
            this.$refs?.batchSelectionColumn?.clearSelection()
          } else {
            this.polling.stop()
          }
        }
      },
    },
    destroyed() {
      this.polling.stop()
    },
    created() {
      this.handleFilter = debounce(this.filterData, 300)
    },
    methods: {
      /**
       * 加载集成了状态、拓扑路径的实例
       */
      async loadIntegrationInstances() {
        this.instancesLoading = true
        const [loadInstanceErr, data] = await to(this.loadInstance())

        if (loadInstanceErr) {
          throw Error(`加载实例失败 ${loadInstanceErr}`)
        }

        if (data.count) {
          const [loadTopoPathErr, topoPath] = await to(this.loadTopoPath(data.info))

          if (loadTopoPathErr) {
            throw Error(`加载 topo path 失败 ${loadTopoPathErr}`)
          }

          const moduleHostCount = await this.getModuleHostCount(data.info)
          if (!moduleHostCount) {
            throw Error('获取实例主机数失败')
          }

          data.info.forEach((module) => {
            const { bk_module_id: moduleId } = module
            const topo = topoPath.nodes.find(topo => topo.topo_node.bk_inst_id === moduleId)
            module.topoPath = topo.topo_path.map(path => path.bk_inst_name).reverse()
              .join(' / ')
            module.hostCount = moduleHostCount.find(item => item.bk_inst_id === moduleId)?.host_count
          })
        }
        this.table.data = data.info
        this.table.backup = Object.freeze(cloneDeep(data.info))
        this.table.pagination.current = 1
        this.table.pagination.count = data.count
        if (data.count > 0 && data.info?.length > 0) {
          await to(this.loadInstanceStatus(this.table.data.map(i => i.bk_module_id)))
        }
        this.renderVisibleList()
        this.instancesLoading = false

        Bus.$emit('module-loaded', data.count)
      },
      /**
        * 加载实例
       */
      loadInstance() {
        return this.$store.dispatch('serviceTemplate/getServiceTemplateModules', {
          bizId: this.bizId,
          serviceTemplateId: this.serviceTemplateId,
        })
      },
      async getModuleHostCount(modules) {
        const moduleIds = modules.map(module => module.bk_module_id)
        try {
          const result = await this.$store.dispatch('objectMainLineModule/getTopoStatistics', {
            bizId: this.bizId,
            params: {
              condition: moduleIds.map(id => ({ bk_obj_id: BUILTIN_MODELS.MODULE, bk_inst_id: id }))
            }
          })
          return result
        } catch (err) {
          console.error(err)
        }
      },
      renderVisibleList() {
        const { limit, current } = this.table.pagination
        this.table.visibleList = this.table.data.slice((current - 1) * limit, current * limit)
      },
      /**
       * 加载实例状态
       */
      loadInstanceStatus(moduleIds) {
        return this.$store.dispatch('serviceTemplate/getServiceTemplateInstanceStatus', {
          bizId: this.bizId,
          params: {
            bk_module_ids: moduleIds,
            service_template_id: parseInt(this.serviceTemplateId, 10)
          }
        }).then((res) => {
          this.table.data.forEach((item) => {
            const instanceStatus = res.find(r => r.bk_inst_id === item.bk_module_id)
            if (instanceStatus) {
              this.$set(item, 'status', instanceStatus.status)
              this.$set(item, 'last_time', instanceStatus.last_time)
              this.$set(item, 'fail_tips', instanceStatus.fail_tips)
            }
          })
          this.table.data.sort(a => (a.status === 'need_sync' ? -1 : 0))

          this.$emit('sync-change')
        })
          .catch((err) => {
            console.log(err)
          })
      },
      /**
       * 判断实例是否正在同步中
       */
      isSyncing(status) {
        return ['new', 'waiting', 'executing'].includes(status)
      },
      /**
       * 实例是否置灰
       */
      isSyncDisabled(status) {
        return this.isSyncing(status) || status === 'finished'
      },
      isDelDisabled(row) {
        return Boolean(row.hostCount) || Boolean(row.set_template_id)
      },
      /**
       * 加载实例拓扑路径
       */
      loadTopoPath(modules) {
        return this.$store.dispatch('objectMainLineModule/getTopoPath', {
          bizId: this.bizId,
          params: {
            topo_nodes: modules.map(module => ({ bk_obj_id: 'module', bk_inst_id: module.bk_module_id }))
          }
        })
      },
      handleSync(row) {
        this.$routerActions.redirect({
          name: 'syncServiceFromTemplate',
          params: {
            modules: row.bk_module_id,
            template: this.serviceTemplateId
          },
          history: true
        })
      },
      handleSelectionChange(selection) {
        this.table.selection = selection
      },
      handleBatchSync() {
        this.$routerActions.redirect({
          name: 'syncServiceFromTemplate',
          params: {
            template: this.serviceTemplateId,
            modules: this.table.selection.map(row => row.bk_module_id).join(',')
          },
          history: true
        })
      },
      handleDel(row) {
        this.$bkInfo({
          title: `${this.$t('确定删除')} ${row.bk_module_name}?`,
          confirmLoading: true,
          confirmFn: async () => {
            try {
              await this.$store.dispatch('objectModule/deleteModule', {
                bizId: this.bizId,
                setId: row.bk_set_id,
                moduleId: row.bk_module_id
              })
              this.$success(this.$t('删除成功'))

              this.loadIntegrationInstances()
            } catch (e) {
              console.error(e)
            }
          }
        })
      },
      filterData() {
        this.table.filtering = true
        this.table.pagination.current = 1
        this.$refs.batchSelectionColumn.clearSelection()
        this.$nextTick(() => {
          if (this.table.filter) {
            this.table.visibleList = this.table.data.filter((row) => {
              const path = row.topoPath.replace(/\s*(\/)\s*/g, '$1')
              const filter = this.table.filter.replace(/\s*(\/)\s*/g, '$1')

              return path.indexOf(filter) > -1
            })
          } else {
            this.renderVisibleList()
          }
          this.table.stuff.type = this.table.filter ? 'search' : 'default'
          this.table.filtering = false
        })
      },
      sortByPath(rowA, rowB) {
        // eslint-disable-next-line no-underscore-dangle
        return rowA.topoPath.toLowerCase().localeCompare(rowB.topoPath.toLowerCase(), 'zh-Hans-CN', { sensitivity: 'accent' })
      },
      sortByTime(rowA, rowB) {
        const timeA = (new Date(rowA.last_time)).getTime()
        const timeB = (new Date(rowB.last_time)).getTime()
        return timeA - timeB
      },
      handleToCreatedInstance() {
        this.$routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
      },
      handlePathClick(row) {
        this.$routerActions.open({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          query: {
            node: `module-${row.bk_module_id}`
          }
        })
      },
      handleSizeChange(size) {
        this.table.pagination.limit = size
        this.table.pagination.current = 1
        this.renderVisibleList()
      },
      handlePageChange(page) {
        this.table.pagination.current = page
        this.renderVisibleList()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .template-instance-layout {
        height: 100%;
        padding: 0 20px;
    }
    .topo-path {
        cursor: pointer;
        &:hover {
            color: $primaryColor;
        }
    }
    .instance-main {
        .options {
            padding: 15px 0;
            font-size: 14px;
            color: #63656E;
            .filter-item {
                width: 230px;
            }
            .icon-cc-updating {
                color: #C4C6CC;
            }
        }
        .latest-sync {
            font-size: 12px;
            cursor: not-allowed;
            color: #DCDEE5;
        }

        .instance-table {
            .topo-path {
                cursor: pointer;
                &:hover {
                    color: $primaryColor;
                }
            }
        }
    }
</style>
