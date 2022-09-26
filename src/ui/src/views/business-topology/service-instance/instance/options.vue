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
  <div class="options">
    <div class="left">
      <cmdb-auth tag="div" :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
        <div slot-scope="{ disabled }"
          v-bk-tooltips="{
            content: () => $refs.tipsContent,
            onShow: () => $refs.tipsContent.style.display = 'block',
            onHide: () => $refs.tipsContent.style.display = 'none',
            theme: 'light',
            allowHtml: true,
            disabled: hasHost && !isEmptyServiceTemplate
          }">
          <bk-button theme="primary" v-test-id="'create'"
            :disabled="disabled || !hasHost || isEmptyServiceTemplate"
            @click="handleCreate">
            {{$t('新增')}}
          </bk-button>
        </div>
      </cmdb-auth>
      <bk-dropdown-menu class="ml10" trigger="click" font-size="medium">
        <bk-button slot="dropdown-trigger" v-test-id="'moreAction'">
          {{$t('实例操作')}}
          <i class="bk-icon icon-angle-down"></i>
        </bk-button>
        <ul class="menu-list" slot="dropdown-content" v-test-id="'moreAction'">
          <cmdb-auth tag="li" class="menu-item"
            :auth="{ type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }">
            <span class="menu-option" slot-scope="{ disabled: authDisabled }"
              :class="{ disabled: authDisabled || !selection.length }"
              @click="handleBatchDelete(authDisabled || !selection.length)">
              {{$t('批量删除')}}
            </span>
          </cmdb-auth>
          <li class="menu-item">
            <span :class="{ 'menu-option': true, disabled: !selection.length }" @click="handleCopy(!selection.length)">
              {{$t('复制IP')}}
            </span>
          </li>
          <cmdb-auth tag="li" class="menu-item"
            :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
            <span class="menu-option" slot-scope="{ disabled: authDisabled }"
              :class="{ disabled: authDisabled || !selection.length }"
              @click="handleBatchEditLabels(authDisabled || !selection.length)">
              {{$t('编辑标签')}}
            </span>
          </cmdb-auth>
        </ul>
      </bk-dropdown-menu>
      <cmdb-auth class="options-sync" v-if="withTemplate"
        :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
        <bk-button slot-scope="{ disabled: authDisabled }"
          :disabled="authDisabled || !hasDifference || isSyncing"
          @click="handleSyncTemplate">
          <span class="sync-wrapper">
            <i :class="['bk-icon', 'icon-refresh', { 'is-syncing': isSyncing }]"></i>
            {{$t(isSyncing ? '同步中' : '同步模板')}}
          </span>
          <span class="topo-status" v-show="hasDifference && !isSyncing"></span>
        </bk-button>
      </cmdb-auth>
    </div>
    <div class="right">
      <bk-checkbox class="options-expand-all" v-model="allExpanded" @change="handleExpandAll">
        {{$t('全部展开')}}
      </bk-checkbox>
      <bk-search-select class="options-search ml10"
        ref="searchSelect"
        :show-condition="false"
        :placeholder="$t('请输入实例名称或选择标签')"
        :data="searchMenuList"
        v-model="searchValue"
        @change="handleSearch">
      </bk-search-select>
      <view-switcher class="ml10" active="instance"></view-switcher>
    </div>
    <span class="tips-content" ref="tipsContent">
      <i18n path="当前模块主机为空，请先XXX" tag="p" v-if="!hasHost">
        <template #action>
          <a href="javascript:;"
            class="action-link"
            @click="handleClickAddHost">
            {{$t('添加主机')}}
          </a>
        </template>
      </i18n>
      <i18n path="当前模块所属的服务模板内进程信息为空，请先XXX" tag="p" v-else-if="isEmptyServiceTemplate">
        <template #action>
          <a href="javascript:;"
            class="action-link"
            @click="handleGoServiceTemplate">
            {{$t('前往补充')}}
          </a>
        </template>
      </i18n>
    </span>
  </div>
</template>

<script>
  import ViewSwitcher from '../common/view-switcher'
  import Bus from '../common/bus'
  import RouterQuery from '@/router/query'
  import { mapGetters } from 'vuex'
  import { Validator } from 'vee-validate'
  import { MULTIPLE_IP_REGEXP } from '@/dictionary/regexp.js'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'

  export default {
    components: {
      ViewSwitcher
    },
    data() {
      return {
        selection: [],
        hasDifference: false,
        hasProcessTemplate: false,
        allExpanded: false,
        historyLabels: {},
        searchValue: [],
        request: {
          label: Symbol('label'),
          diff: Symbol('diff')
        },
        syncStatusTimer: null,
        isSyncing: false
      }
    },
    computed: {
      ...mapGetters('businessHost', ['selectedNode']),
      bizId() {
        const { objectBiz, bizSet } = this.$store.state
        return objectBiz.bizId || bizSet.bizId
      },
      withTemplate() {
        return this.selectedNode && this.selectedNode.data.service_template_id
      },
      hasHost() {
        return this.selectedNode && this.selectedNode.data?.host_count > 0
      },
      nameFilterIndex() {
        return this.searchValue.findIndex(data => data.id === 'name')
      },
      searchMenuList() {
        const hasHistoryLables = Object.keys(this.historyLabels).length > 0
        const list = [{
          id: 'name',
          name: this.$t('服务实例名')
        }, {
          id: 'tagValue',
          name: this.$t('标签值'),
          conditions: Object.keys(this.historyLabels).map(key => ({ id: key, name: `${key}:` })),
          disabled: !hasHistoryLables
        }, {
          id: 'tagKey',
          name: this.$t('标签键'),
          disabled: !hasHistoryLables
        }]
        if (this.nameFilterIndex > -1) {
          return list.slice(1)
        }
        return list.slice(0)
      },
      isEmptyServiceTemplate() {
        return this.withTemplate && !this.hasProcessTemplate
      }
    },
    watch: {
      withTemplate: {
        immediate: true,
        handler(withTemplate) {
          if (withTemplate) {
            this.checkProcessTemplate()
            this.checkSyncStatus()
          }
        }
      },
      searchMenuList() {
        this.$nextTick(() => {
          const menu = this.$refs.searchSelect && this.$refs.searchSelect.menuInstance
          menu && (menu.list = this.$refs.searchSelect.data)
        })
      },
      selectedNode() {
        this.searchValue = []
        RouterQuery.set({
          node: this.selectedNode.id,
          instanceName: ''
        })
        Bus.$emit('filter-change', this.searchValue)
      }
    },
    created() {
      this.unwatch = RouterQuery.watch(['node', 'page'], () => {
        this.allExpanded = false
      })
      Bus.$on('instance-selection-change', this.handleInstanceSelectionChange)
      Bus.$on('update-labels', this.updateHistoryLabels)
      Bus.$on('delete-complete', this.checkDifference)
      this.setFilter()
      this.updateHistoryLabels()
    },
    beforeDestroy() {
      this.unwatch()
      Bus.$off('instance-selection-change', this.handleInstanceSelectionChange)
      Bus.$off('update-labels', this.updateHistoryLabels)
      Bus.$off('delete-complete', this.checkDifference)

      if (this.syncStatusTimer) {
        clearTimeout(this.syncStatusTimer)
      }
    },
    methods: {
      async updateHistoryLabels() {
        try {
          this.historyLabels = await this.$store.dispatch('instanceLabel/getHistoryLabel', {
            params: {
              bk_biz_id: this.bizId
            }
          })
        } catch (error) {
          console.error(error)
        }
      },
      setFilter() {
        const filterName = RouterQuery.get('instanceName')
        if (!filterName) {
          return false
        }
        this.searchValue.push({
          id: 'name',
          name: this.$t('服务实例名'),
          values: [{
            id: filterName,
            name: filterName
          }]
        })
        this.$nextTick(() => {
          // list.vue注册的监听晚于派发，因此nextTick再触发
          Bus.$emit('filter-change', this.searchValue)
        })
      },
      handleInstanceSelectionChange(selection) {
        this.selection = selection
      },
      handleCreate() {
        this.$routerActions.redirect({
          name: 'createServiceInstance',
          params: {
            moduleId: this.selectedNode.data.bk_inst_id,
            setId: this.selectedNode.parent.data.bk_inst_id
          },
          query: {
            title: this.selectedNode.data.bk_inst_name,
            node: this.selectedNode.id,
            tab: 'serviceInstance'
          },
          history: true
        })
      },
      async handleBatchDelete(disabled) {
        if (disabled) {
          return false
        }
        this.$bkInfo({
          title: this.$t('确定删除N个实例', { count: this.selection.length }),
          confirmLoading: true,
          confirmFn: async () => {
            const serviceInstanceIds = this.selection.map(row => row.id)
            try {
              await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                config: {
                  data: {
                    service_instance_ids: serviceInstanceIds,
                    bk_biz_id: this.bizId
                  }
                }
              })
              this.$success(this.$t('删除成功'))
              Bus.$emit('delete-complete')
              return true
            } catch (e) {
              console.error(e)
              return false
            }
          }
        })
      },
      async handleCopy(disabled) {
        if (disabled) {
          return false
        }
        try {
          const validator = new Validator()
          const validPromise = []
          this.selection.forEach((row) => {
            // eslint-disable-next-line prefer-destructuring
            const ip = row.name.split('_')[0]
            validPromise.push(new Promise(async (resolve) => {
              const { valid } = await validator.verify(ip, MULTIPLE_IP_REGEXP)
              resolve({ valid, ip })
            }))
          })
          const results = await Promise.all(validPromise)
          const validResult = results.filter(result => result.valid).map(result => result.ip)
          const unique = [...new Set(validResult)]
          if (unique.length) {
            await this.$copyText(unique.join('\n'))
            this.$success(this.$t('复制成功'))
          } else {
            this.$warn(this.$t('暂无可复制的IP'))
          }
        } catch (e) {
          console.error(e)
          this.$error(this.$t('复制失败'))
        }
      },
      handleEidtTag(disabled) {
        if (disabled) {
          return false
        }
      },
      async checkDifference() {
        if (!this.withTemplate) {
          return
        }

        try {
          // eslint-disable-next-line max-len
          const { modules: syncStatus = [] } = await this.$store.dispatch('serviceTemplate/getServiceTemplateSyncStatus', {
            bizId: this.bizId,
            params: {
              is_partial: false,
              bk_module_ids: [this.selectedNode.data.bk_inst_id],
            },
            config: {
              requestId: this.request.diff,
              cancelPrevious: true
            }
          })
          // eslint-disable-next-line max-len
          const difference = syncStatus.find(difference => difference.bk_module_id === this.selectedNode.data.bk_inst_id)
          this.hasDifference = !!difference && difference.need_sync
        } catch (error) {
          console.error(error)
        }
      },
      async checkProcessTemplate() {
        try {
          const processData = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
            params: {
              page: { start: 0, limit: 1 },
              bk_biz_id: this.bizId,
              service_template_id: this.selectedNode.data.service_template_id
            }
          })
          this.hasProcessTemplate = processData.info?.length > 0
        } catch (error) {
          console.error(error)
        }
      },
      async checkSyncStatus() {
        const [moduleSyncData] = await this.$store.dispatch('serviceTemplate/getServiceTemplateInstanceStatus', {
          bizId: this.bizId,
          params: {
            bk_module_ids: [this.selectedNode.data.bk_inst_id],
            service_template_id: this.selectedNode.data.service_template_id
          }
        })
        this.isSyncing = ['new', 'waiting', 'executing'].includes(moduleSyncData?.status)

        if (this.isSyncing) {
          if (this.syncStatusTimer) {
            clearTimeout(this.syncStatusTimer)
          }

          this.syncStatusTimer = setTimeout(() => {
            this.checkSyncStatus()
          }, 1000)
        } else {
          this.checkDifference()
        }
      },
      async handleSyncTemplate() {
        if (!this.hasDifference) {
          this.$info(this.$t('已是最新内容，无需同步'))
          return
        }

        this.$routerActions.redirect({
          name: 'syncServiceFromModule',
          params: {
            modules: String(this.selectedNode.data.bk_inst_id),
            template: this.selectedNode.data.service_template_id
          },
          history: true
        })
      },
      handleSearch(filters) {
        const transformedFilters = []
        filters.forEach((data) => {
          if (!data.values) {
            const nameIndex = transformedFilters.findIndex(filter => filter.id === 'name')
            if (nameIndex > -1) {
              transformedFilters[nameIndex].values = [{ ...data }]
            } else {
              transformedFilters.push({
                id: 'name',
                name: this.$t('服务实例名'),
                values: [{ ...data }]
              })
            }
          } else if (data.id === 'tagValue') {
            const [{ name }] = data.values
            if (data.condition && name) {
              transformedFilters.push(data)
            } else {
              this.$warn(this.$t('服务实例标签值搜索提示语'))
            }
          } else {
            transformedFilters.push(data)
          }
        })
        this.$nextTick(() => {
          this.searchValue = transformedFilters
          Bus.$emit('filter-change', this.searchValue)
        })
      },
      handleExpandAll(expanded) {
        Bus.$emit('expand-all-change', expanded)
      },
      handleBatchEditLabels(disabled) {
        if (disabled) {
          return false
        }
        Bus.$emit('batch-edit-labels')
      },
      handleClickAddHost() {
        RouterQuery.set({
          tab: 'hostList',
          _t: Date.now(),
          page: '',
          limit: ''
        })
      },
      handleGoServiceTemplate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            templateId: this.selectedNode.data.service_template_id
          },
          history: true
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .options {
        display: flex;
        justify-content: space-between;
        flex-wrap: wrap;
        .left,
        .right {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
    }
    .menu-list {
        .menu-item {
            line-height: 32px;
            .menu-option {
                display: block;
                padding: 0 10px;
                font-size: 14px;
                cursor: pointer;
                &:hover {
                    color: $primaryColor;
                }
                &.disabled {
                    color: $textDisabledColor;
                    cursor: not-allowed;;
                }
            }
        }
    }
    .options-sync {
        display: inline-block;
        position: relative;
        margin-left: 18px;
        padding: 0;
        &::before {
            content: '';
            position: absolute;
            top: 7px;
            left: -11px;;
            width: 1px;
            height: 20px;
            background-color: #dcdee5;
        }
        .sync-wrapper {
            display: flex;
            align-items: center;
        }
        .icon-refresh {
            font-size: 12px;
            margin-right: 4px;

            &.is-syncing {
              animation: spiner 1s linear infinite
            }
        }
        .topo-status {
            position: absolute;
            top: -4px;
            right: -4px;
            width: 8px;
            height: 8px;
            background-color: #ea3636;
            border-radius: 50%;
        }
    }
    .options-search {
        width: 200px;
    }

    .tips-content {
      display: none;
      .action-link {
        color: $primaryColor;
        margin-left: 2px;
      }
    }

    @keyframes spiner {
      0% {
        transform: rotate(0deg)
      }
      100% {
        transform: rotate(360deg)
      }
    }
</style>
