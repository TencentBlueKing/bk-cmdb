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
  <div class="host-apply-sidebar">
    <div class="tree-wrapper">
      <div class="config-mode">
        <bk-radio-group class="radio-group" v-model="configMode">
          <bk-radio-button
            :disabled="Boolean(actionMode) && configMode !== CONFIG_MODE.MODULE"
            :value="CONFIG_MODE.MODULE">{{$t('按业务拓扑')}}</bk-radio-button>
          <bk-radio-button
            :disabled="Boolean(actionMode) && configMode !== CONFIG_MODE.TEMPLATE"
            :value="CONFIG_MODE.TEMPLATE">{{$t('按服务模板')}}</bk-radio-button>
        </bk-radio-group>
      </div>
      <div class="searchbar">
        <div class="search-select">
          <search-select-mix :mode="configMode" />
        </div>
        <div class="action-menu" v-show="!actionMode">
          <bk-dropdown-menu
            @show="showBatchDropdown = true"
            @hide="showBatchDropdown = false"
            font-size="medium">
            <div :class="['dropdown-trigger', { selected: actionMode }]" slot="dropdown-trigger">
              <span>{{$t('批量操作')}}</span>
              <i :class="['bk-icon icon-angle-down', { 'icon-flip': showBatchDropdown }]"></i>
            </div>
            <ul class="bk-dropdown-list" slot="dropdown-content">
              <li>
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                  <a
                    href="javascript:;"
                    slot-scope="{ disabled }"
                    :class="{ disabled }"
                    @click="handleBatchAction('batch-edit')">
                    {{$t('批量编辑')}}
                  </a>
                </cmdb-auth>
              </li>
              <li>
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                  <a
                    href="javascript:;"
                    slot-scope="{ disabled }"
                    :class="{ disabled }"
                    @click="handleBatchAction('batch-del')">
                    {{$t('批量删除')}}
                  </a>
                </cmdb-auth>
              </li>
            </ul>
          </bk-dropdown-menu>
        </div>
      </div>
      <KeepAlive>
        <component
          :is="configTree.is"
          ref="topologyTree"
          v-bind="configTree.props"
          v-on="configTree.events" />
      </KeepAlive>
    </div>
    <div class="checked-list" v-show="showCheckPanel">
      <batch-check-panel
        v-bind="checkPanelProps"
        @clear="handleClearChecked"
        @edit="handleGoEdit"
        @cancel="handleCancelEdit">
        <template #title>
          <i18n path="已选择N个模块" v-if="isModuleMode">
            <template #count><em class="checked-num">{{checkedList.length}}</em></template>
          </i18n>
          <i18n path="已选择N个模板" v-else-if="isTemplateMode">
            <template #count><em class="checked-num">{{checkedTemplateList.length}}</em></template>
          </i18n>
        </template>
        <template #content>
          <dl class="module-list" v-if="isModuleMode">
            <div class="module-item" v-for="item in checkedList" :key="item.bk_inst_id">
              <dt class="module-name">{{item.bk_inst_name}}</dt>
              <dd class="module-path" :title="item.path.join(' / ')">{{item.path.join(' / ')}}</dd>
              <dd class="module-icon"><span>{{$i18n.locale === 'en' ? 'M' : '模'}}</span></dd>
              <dd class="action-icon">
                <a href="javascript:;" @click="handleRemoveChecked(item.bk_inst_id)">
                  <i class="bk-icon icon-close"></i>
                </a>
              </dd>
            </div>
          </dl>
          <ul class="template-list" v-else-if="isTemplateMode">
            <li class="template-item" v-for="item in checkedTemplateList" :key="item.id">
              <i class="template-icon"><span>{{$i18n.locale === 'en' ? 'M' : '模'}}</span></i>
              <div class="template-name">{{item.name}}</div>
              <a href="javascript:;" class="action-icon" @click="handleRemoveChecked(item.id)">
                <i class="bk-icon icon-close"></i>
              </a>
            </li>
          </ul>
        </template>
      </batch-check-panel>
    </div>
  </div>
</template>

<script>
  import searchSelectMix from './search-select-mix'
  import topologyTreeComponent from './topology-tree'
  import serviceTemplateListComponent from './service-template-list'
  import batchCheckPanel from './batch-check-panel.vue'
  import { MENU_BUSINESS_HOST_APPLY_EDIT } from '@/dictionary/menu-symbol'
  import { mapGetters } from 'vuex'
  import { CONFIG_MODE } from '@/services/service-template/index.js'
  export default {
    components: {
      searchSelectMix,
      batchCheckPanel
    },
    data() {
      return {
        treeOptions: {
          showCheckbox: false,
          selectable: true,
          checkOnClick: false,
          checkOnlyAvailableStrictly: false,
          displayMatchedNodeDescendants: true
        },
        actionMode: '',
        showCheckPanel: false,
        checkedList: [],
        checkedTemplateList: [],
        showBatchDropdown: false,
        CONFIG_MODE
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      configMode: {
        get() {
          return this.$route.params.mode || CONFIG_MODE.MODULE
        },
        set(mode) {
          this.$emit('mode-changed', mode)
        }
      },
      isModuleMode() {
        return this.configMode === CONFIG_MODE.MODULE
      },
      isTemplateMode() {
        return this.configMode === CONFIG_MODE.TEMPLATE
      },
      configTree() {
        const treeComponents = {
          [CONFIG_MODE.MODULE]: {
            is: topologyTreeComponent,
            props: {
              treeOptions: this.treeOptions,
              action: this.actionMode
            },
            events: {
              selected: this.handleTreeSelected,
              checked: this.handleTreeChecked
            }
          },
          [CONFIG_MODE.TEMPLATE]: {
            is: serviceTemplateListComponent,
            props: {
              options: this.treeOptions,
              action: this.actionMode
            },
            events: {
              selected: this.handleTreeSelected,
              checked: this.handleTemplateChecked
            }
          }
        }
        return treeComponents[this.configMode]
      },
      checkPanelProps() {
        const componentProps = {
          [CONFIG_MODE.MODULE]: {
            checkedList: this.checkedList,
            actionMode: this.actionMode
          },
          [CONFIG_MODE.TEMPLATE]: {
            checkedList: this.checkedTemplateList,
            actionMode: this.actionMode
          }
        }
        return componentProps[this.configMode]
      }
    },
    watch: {
      actionMode(value) {
        this.$emit('action-change', value)
      }
    },
    methods: {
      setApplyClosed(targetId, options) {
        this.$refs.topologyTree.updateNodeStatus(targetId, options)
      },
      removeChecked() {
        if (this.configMode === CONFIG_MODE.MODULE) {
          this.$refs.topologyTree.$refs.tree.removeChecked({ emitEvent: true })
        } else if (this.configMode === CONFIG_MODE.TEMPLATE) {
          this.$refs.topologyTree.removeChecked()
        }
      },
      async handleBatchAction(actionMode) {
        this.actionMode = actionMode
        this.showCheckPanel = true
        this.treeOptions.showCheckbox = true
        this.treeOptions.selectable = false
        this.treeOptions.checkOnClick = true
        this.treeOptions.checkOnlyAvailableStrictly = true
      },
      handleTreeSelected(node) {
        this.$emit('module-selected', node.data)
      },
      handleTreeChecked(ids) {
        const { treeData } = this.$refs.topologyTree
        const modules = []
        const findModuleNode = function (data, parent) {
          data.forEach((item) => {
            item.path = parent ? [...parent.path, item.bk_inst_name] : [item.bk_inst_name]
            if (item.bk_obj_id === 'module' && ids.includes(`module_${item.bk_inst_id}`)) {
              modules.push(item)
            }
            if (item.child) {
              findModuleNode(item.child, item)
            }
          })
        }
        findModuleNode(treeData)

        this.checkedList = modules
      },
      handleTemplateChecked(ids) {
        // 清空
        if (!ids?.length) {
          this.checkedTemplateList = []
          return
        }

        this.checkedTemplateList = this.$refs.topologyTree.displayList.filter(item => ids.includes(item.id))
      },
      handleRemoveChecked(id) {
        if (this.configMode === CONFIG_MODE.MODULE) {
          const { tree } = this.$refs.topologyTree.$refs
          const checkedIds = this.checkedList.filter(item => item.bk_inst_id !== id).map(item => `module_${item.bk_inst_id}`)
          tree.removeChecked({ emitEvent: true })
          tree.setChecked(checkedIds, { emitEvent: true, beforeCheck: true, checked: true })
        } else if (this.configMode === CONFIG_MODE.TEMPLATE) {
          this.$refs.topologyTree.removeChecked(id)
        }
      },
      handleClearChecked() {
        this.removeChecked()
      },
      handleGoEdit() {
        const checkedIds = {
          [CONFIG_MODE.MODULE]: this.checkedList.map(item => item.bk_inst_id),
          [CONFIG_MODE.TEMPLATE]: this.checkedTemplateList.map(item => item.id)
        }
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY_EDIT,
          query: {
            id: checkedIds[this.configMode].join(','),
            batch: 1,
            action: this.actionMode
          },
          history: true
        })
      },
      handleCancelEdit() {
        this.treeOptions.showCheckbox = false
        this.treeOptions.selectable = true
        this.treeOptions.checkOnClick = false
        this.showCheckPanel = false
        this.treeOptions.checkOnlyAvailableStrictly = false
        this.actionMode = ''
        this.removeChecked()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .host-apply-sidebar {
    position: relative;
    height: 100%;
    padding: 10px 0;

    .tree-wrapper {
      height: 100%;
    }
  }

  .config-mode {
    padding: 0 10px;
    margin-bottom: 12px;
    .radio-group {
      display: flex;
      .bk-form-radio-button {
        flex: 1;
        ::v-deep .bk-radio-button-text {
          display: block;
        }
        ::v-deep &:not(:first-child) .bk-radio-button-input:disabled + .bk-radio-button-text {
          border-left: none;
        }
        ::v-deep .bk-radio-button-input:disabled + .bk-radio-button-text {
          border: 1px solid #c4c6cc;
        }
      }
    }
  }

  .searchbar {
    display: flex;
    padding: 0 10px;

    .search-select {
      flex: 1;
    }
    .action-menu {
      flex: none;
      margin-left: 8px;

      .dropdown-trigger {
        border: 1px solid #c4c6cc;
        border-radius: 2px;
        padding: 0 8px;
        height: 32px;
        text-align: center;
        line-height: 32px;
        cursor: pointer;
        &:hover {
          border-color: #979ba5;
          color: #63656e;
        }
        &:active {
          border-color: #3a84ff;
          color: #3a84ff;
        }

        .icon-angle-down {
          font-size: 22px;
          margin: -3px -5px 0 -4px
        }
      }
    }
  }

  .checked-list {
    position: absolute;
    width: 290px;
    height: 100%;
    left: 100%;
    top: 0;
    z-index: 9999;
    border-left: 1px solid #dcdee5;
    border-right: 1px solid #dcdee5;

    .panel-hd,
    .panel-ft {
      background: #fafbfd;
    }
    .panel-hd {
      height: 52px;
      line-height: 52px;
      padding: 0 12px;
      border-bottom: 1px solid #dcdee5;
    }
    .panel-title {
      position: relative;
      font-size: 14px;
      color: #63656e;

      .checked-num {
        font-style: normal;
        font-weight: bold;
        color: #2dcb56;
        margin: .1em;
      }

      .clear-all {
          position: absolute;
        right: 0;
        top: 0;
        color: #3a84ff;
      }
    }
    .panel-bd {
      height: calc(100% - 52px - 60px);
      background: #fff;
      @include scrollbar-y;
    }
    .panel-ft {
      height: 60px;
      line-height: 58px;
      text-align: center;
      border-top: 1px solid #dcdee5;
      .bk-button {
        min-width: 86px;
        margin: 0 3px;
      }
    }

    .module-list {
      .module-item {
        position: relative;
        padding: 8px 42px;

        &:hover {
          background: #f0f1f5;
        }

        .module-name {
          font-size: 14px;
          color: #63656e;
        }
        .module-path {
          font-size: 12px;
          color: #c4c6cc;
          @include ellipsis;
        }
        .module-icon {
          position: absolute;
          left: 12px;
          top: 8px;
          font-size: 12px;
          border-radius: 50%;
          background-color: #c4c6cc;
          width: 22px;
          height: 22px;
          line-height: 21px;
          text-align: center;
          font-style: normal;
          color: #fff;
        }
        .action-icon {
          position: absolute;
          right: 8px;
          top: 10px;
          width: 28px;
          height: 28px;
          text-align: center;
          line-height: 28px;

          a {
            color: #c4c6cc;
            &:hover {
              color: #979ba5;
            }
          }
        }
      }
    }

    .template-list {
      .template-item {
        display: flex;
        align-items: center;
        padding: 8px 12px;
        &:hover {
          background: #f0f1f5;
        }

        .template-name {
          flex: 1;
          font-size: 14px;
          color: #63656e;
          margin: 0 8px;
          @include ellipsis;
        }
        .template-icon {
          font-size: 12px;
          border-radius: 50%;
          background-color: #c4c6cc;
          width: 22px;
          height: 22px;
          line-height: 21px;
          text-align: center;
          font-style: normal;
          color: #fff;
        }
        .action-icon {
          width: 28px;
          height: 28px;
          color: #c4c6cc;
          text-align: center;
          line-height: 26px;
          &:hover {
            color: #979ba5;
          }
        }
      }
    }
  }

  .bk-dropdown-list {
    .auth-box {
      width: 100%;
    }

    > li a {
      display: block;
      height: 32px;
      line-height: 33px;
      padding: 0 16px;
      color: #63656e;
      font-size: 14px;
      text-decoration: none;
      white-space: nowrap;

      &:hover {
        background-color: #eaf3ff;
        color: #3a84ff;
      }

      &.disabled {
        color: #c4c6cc;
      }
    }
  }
</style>
