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
  <div class="host-selector-layout">
    <div class="layout-content">
      <div class="module-host">
        <h2 class="title">{{title || $t('选择主机')}}</h2>
        <div class="host-list">
          <div class="top-bar">
            <div class="topo-path">已自动过滤<div class="path" :title="topoPath">{{topoPath}}</div>可选的主机：</div>
            <div class="search-input">
              <bk-input clearable right-icon="bk-icon icon-search" v-model.trim="keyword"></bk-input>
            </div>
          </div>
          <bk-table
            v-bkloading="{ isLoading: $loading(request.table) }"
            ref="table"
            :data="list"
            :pagination="pagination"
            :max-height="480"
            :outer-border="false"
            :header-border="false"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @select="handleToggleSelect"
            @select-all="handleToggleSelectAll">
            <bk-table-column type="selection" width="30" :selectable="getSelectable"></bk-table-column>
            <bk-table-column :label="$t('内网IP')">
              <template slot-scope="{ row }">
                {{row.host.bk_host_innerip}}
              </template>
            </bk-table-column>
            <bk-table-column :label="$t('云区域')" show-overflow-tooltip>
              <template slot-scope="{ row }">{{row.host.bk_cloud_id | foreignkey}}</template>
            </bk-table-column>
          </bk-table>
        </div>
      </div>
      <div class="layout-col result-preview">
        <div class="preview-title">{{$t('结果预览')}}</div>
        <div class="preview-content" v-show="selected.length > 0">
          <bk-collapse v-model="collapse.active" @item-click="handleCollapse">
            <bk-collapse-item name="host" hide-arrow>
              <div class="collapse-title">
                <div class="text">
                  <i :class="['bk-icon icon-angle-right', { expand: collapse.expanded.host }]"></i>
                  <i18n path="已选择N台主机">
                    <template #count><span class="count">{{selected.length}}</span></template>
                  </i18n>
                </div>
                <div class="more" @click.stop="handleClickMore">
                  <i class="bk-icon icon-more"></i>
                </div>
              </div>
              <ul slot="content" class="selected-host-list">
                <li class="host-item" v-for="(row, index) in selected" :key="index">
                  <div class="ip">
                    {{row.host.bk_host_innerip}}
                  </div>
                  <i class="bk-icon icon-close-line" @click="handleRemove(row)"></i>
                </li>
              </ul>
            </bk-collapse-item>
          </bk-collapse>
        </div>
      </div>
    </div>
    <div class="layout-footer">
      <div class="tips">注：已存在服务实例的主机，不能再勾选</div>
      <div>
        <bk-button class="mr10" theme="primary" :disabled="!selected.length" @click="handleNextStep">
          {{confirmText || $t('下一步')}}
        </bk-button>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </div>
    <ul class="more-menu" ref="moreMenu" v-show="more.show">
      <li class="menu-item" @click="handleCopyIp">{{$t('复制IP')}}</li>
      <li class="menu-item" @click="handleRemoveAll">{{$t('移除所有')}}</li>
    </ul>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import debounce from 'lodash.debounce'
  import { foreignkey } from '@/filters/formatter.js'
  export default {
    filters: {
      foreignkey
    },
    props: {
      moduleId: {
        type: Number,
        required: true
      },
      withTemplate: {
        type: Boolean,
        required: true
      },
      exist: {
        type: Array,
        default: () => ([])
      },
      confirmText: {
        type: String,
        default: ''
      },
      title: String
    },
    data() {
      return {
        list: [],
        noServiceInstanceHost: [],
        pagination: this.$tools.getDefaultPaginationConfig(),
        topoNodesPath: [],
        keyword: '',
        selected: this.exist || [],
        collapse: {
          expanded: { host: true },
          active: 'host'
        },
        more: {
          instance: null,
          show: false
        },
        request: {
          table: Symbol('table')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      topoPath() {
        const [nodePath = {}] = this.topoNodesPath
        // eslint-disable-next-line newline-per-chained-call
        return nodePath.topo_path?.reverse().map(node => node.bk_inst_name).join('/')
      }
    },
    watch: {
      list() {
        this.setChecked()
        this.getNoServiceInstanceHost()
      },
      keyword() {
        this.handleFilter()
      },
      selected() {
        this.setChecked()
      }
    },
    created() {
      this.getModuleHost()
      this.getModulePathInfo()
      this.handleFilter = debounce(this.getModuleHost, 300)
    },
    methods: {
      async getModuleHost() {
        try {
          const result = await this.$store.dispatch('hostSearch/searchHost', {
            params: this.getParams(),
            config: {
              requestId: this.request.table,
              cancelPrevious: true
            }
          })
          this.list = result.info
          this.pagination.count = result.count
        } catch {
          this.list = []
          this.pagination.count = 0
        }
      },
      async getModulePathInfo() {
        try {
          const result = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
            bizId: this.bizId,
            params: {
              topo_nodes: [{ bk_obj_id: 'module', bk_inst_id: this.moduleId }]
            }
          })
          this.topoNodesPath = result.nodes || []
        } catch (e) {
          console.error(e)
          this.topoNodesPath = []
        }
      },
      async getNoServiceInstanceHost() {
        if (!this.withTemplate) {
          return
        }

        try {
          const result = await this.$store.dispatch('serviceInstance/getNoServiceInstanceHost', {
            params: {
              bk_biz_id: this.bizId,
              bk_module_id: this.moduleId,
              bk_host_ids: this.list.map(item => item.host.bk_host_id)
            }
          })
          this.noServiceInstanceHost = result.bk_host_ids
        } catch (e) {
          console.error(e)
        }
      },
      getParams() {
        const params = {
          bk_biz_id: this.bizId,
          ip: {
            data: this.keyword ? [this.keyword] : [],
            exact: 0, // 模糊
            flag: 'bk_host_innerip|bk_host_outerip'
          },
          condition: [
            {
              bk_obj_id: 'host',
              fields: [
                'bk_host_id',
                'bk_host_innerip',
                'bk_cloud_id',
                'bk_host_outerip'
              ],
              condition: []
            },
            {
              bk_obj_id: 'module',
              fields: [
                'bk_module_name'
              ],
              condition: [
                {
                  field: 'bk_module_id',
                  operator: '$eq',
                  value: this.moduleId
                }
              ]
            }
          ],
          page: {
            ...this.$tools.getPageParams(this.pagination)
          }
        }

        return params
      },
      setChecked() {
        this.$nextTick(() => {
          const ids = [...new Set(this.selected.map(data => data.host.bk_host_id))]
          this.list.forEach((row) => {
            this.$refs.table.toggleRowSelection(row, ids.includes(row.host.bk_host_id))
          })
        })
      },
      getSelectable(row) {
        // 非模板创建不限制重复添加服务实例
        if (!this.withTemplate) {
          return true
        }

        return this.noServiceInstanceHost.includes(row.host.bk_host_id)
      },
      handleToggleSelect(selection) {
        this.handleSelectionChange(selection)
      },
      handleToggleSelectAll(selection) {
        this.handleSelectionChange(selection)
      },
      handleSelectionChange(selection) {
        const ids = [...new Set(selection.map(data => data.host.bk_host_id))]
        const unselected = this.list.filter(item => !ids.includes(item.host.bk_host_id))

        this.handleRemove(unselected)
        this.handleSelect(selection)
      },
      handleRemove(hosts) {
        const removeData = Array.isArray(hosts) ? hosts : [hosts]
        const ids = removeData.map(data => data.host.bk_host_id)
        this.selected = this.selected.filter(target => !ids.includes(target.host.bk_host_id))
      },
      handleSelect(hosts) {
        const selectData = Array.isArray(hosts) ? hosts : [hosts]
        selectData.forEach((data) => {
          const isExist = this.selected.find(target => target.host.bk_host_id === data.host.bk_host_id)
          if (!isExist) {
            this.selected.push(data)
          }
        })
      },
      handleClickMore(event) {
        this.more.instance && this.more.instance.destroy()
        this.more.instance = this.$bkPopover(event.target, {
          content: this.$refs.moreMenu,
          allowHTML: true,
          delay: 300,
          trigger: 'manual',
          boundary: 'window',
          placement: 'bottom-end',
          theme: 'light host-selector-popover',
          distance: 6,
          interactive: true
        })
        this.more.show = true
        this.$nextTick(() => {
          this.more.instance.show()
        })
      },
      handleCopyIp() {
        const ipList = this.selected.map(item => item.host.bk_host_innerip)
        this.$copyText(ipList.join('\n')).then(() => {
          this.$success(this.$t('复制成功'))
        }, () => {
          this.$error(this.$t('复制失败'))
        })
          .finally(() => {
            this.more.instance.hide()
          })
      },
      handleRemoveAll() {
        this.selected = []
        this.more.instance.hide()
      },
      handleCancel() {
        this.$emit('cancel')
      },
      handleNextStep() {
        this.$emit('confirm', this.selected)
      },
      handleCollapse(names) {
        Object.keys(this.collapse.expanded).forEach(key => (this.collapse.expanded[key] = false))
        names.forEach((name) => {
          this.$set(this.collapse.expanded, name, true)
        })
      },
      handlePageChange(current = 1) {
        this.pagination.current = current
        this.getModuleHost()
      },
      handleLimitChange(limit) {
        this.pagination.current = 1
        this.pagination.limit = limit
        this.getModuleHost()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .host-selector-layout {
        display: flex;
        flex-direction: column;
        height: 650px;
        padding: 0;

        .layout-content {
            display: flex;
            flex: auto;
            overflow: hidden;

            .module-host {
                padding: 20px 20px 0 20px;
                flex: auto;
                width: 0; // fix table宽度

                .title {
                  font-size: 20px;
                    font-weight: normal;
                    color: $textColor;
                    margin-bottom: 20px;
                }

                .host-list {
                    .top-bar {
                      display: flex;
                      align-items: center;
                      justify-content: space-between;
                      margin-bottom: 10px;

                      .topo-path {
                        display: flex;
                        font-size: 14px;
                        .path {
                          font-weight: bold;
                          margin: 0 2px;
                          max-width: 300px;
                          text-overflow: ellipsis;
                          overflow: hidden;
                          white-space: nowrap;
                        }
                      }
                      .search-input {
                        width: 300px;
                      }
                    }
                }
            }

            .result-preview {
                flex: none;
                width: 280px;
                border-left: 1px solid #dcdee5;
                background: #f5f6fa;

                .preview-title {
                    color: #313238;
                    font-size: 14px;
                    line-height: 22px;
                    padding: 10px 24px;
                }
                .preview-content {
                    height: calc(100% - 42px);
                    @include scrollbar;

                    /deep/ .bk-collapse-item .bk-collapse-item-header {
                        padding: 0;
                        height: 24px;
                        line-height: 24px;
                        &:hover {
                            color: #63656e;
                        }
                    }
                }
            }
        }

        .layout-footer {
            display: flex;
            align-items: center;
            justify-content: space-between;
            text-align: right;
            border-top: 1px solid #dcdee5;
            padding: 12px 20px;

            .tips {
              font-size: 12px;
            }
        }

        .result-preview {
            .collapse-title {
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 0 24px 0 18px;

                .icon-angle-right {
                    font-size: 24px;
                    transition: transform .2s ease-in-out;

                    &.expand {
                        transform: rotate(90deg);
                    }
                }
                .count {
                    font-weight: 700;
                    color: #3a84ff;
                    padding: 0 .1em;
                }
                .text {
                    display: flex;
                    align-items: center;
                }
                .more {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    width: 24px;
                    height: 24px;
                    border-radius: 2px;
                    .icon-more {
                        font-size: 18px;
                        outline: 0;
                    }
                    &:hover {
                        background: #e1ecff;
                        color: #3a84ff;
                    }
                }
            }
            .selected-host-list {
                padding: 0 14px;
                margin: 6px 0;
            }
            .host-item {
                display: flex;
                justify-content: space-between;
                align-items: center;
                height: 32px;
                line-height: 32px;
                background: #fff;
                padding: 0 12px;
                border-radius: 2px;
                box-shadow: 0 1px 2px 0 rgba(0,0,0,.06);
                margin-bottom: 2px;
                font-size: 12px;

                .ip {
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    word-break: break-all;

                    .repeat-tag {
                        height: 16px;
                        line-height: 16px;
                        background: #FFE8C3;
                        color: #FE9C00;
                        font-size: 12px;
                        padding: 0 4px;
                        border-radius: 2px;
                    }
                }
                .icon-close-line {
                    display: none;
                    cursor: pointer;
                    color: #3a84ff;
                    font-weight: 700;
                }

                &:hover {
                    .icon-close-line {
                        display: block;
                    }
                }
            }
        }
    }
</style>
<style lang="scss">
    .host-selector-popover-theme {
        padding: 0;
        .more-menu {
            font-size: 12px;
            padding: 6px 0;
            min-width: 84px;
            background: #fff;
            .menu-item {
                height: 32px;
                line-height: 32px;
                padding: 0 10px;
                cursor: pointer;
                &:hover {
                    background: #f5f6fa;
                    color: #3a84ff;
                }
            }
        }
    }
</style>
