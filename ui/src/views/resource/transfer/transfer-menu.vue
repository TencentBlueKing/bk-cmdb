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
  <div class="transfer-menu">
    <bk-dropdown-menu
      trigger="click"
      font-size="medium"
      :disabled="disabled"
      @show="handleMenuToggle(true)"
      @hide="handleMenuToggle(false)">
      <div class="dropdown-trigger-btn" style="padding-left: 19px;" slot="dropdown-trigger">
        <span>{{$t('转移到')}}</span>
        <i :class="['bk-icon icon-angle-down', { 'icon-flip': isShow }]"></i>
      </div>
      <ul class="bk-dropdown-list" slot="dropdown-content">
        <li>
          <a href="javascript:;" @click="transferToIdleModule">
            {{$store.state.globalConfig.config.idlePool.idle}}
          </a>
        </li>
        <li>
          <a href="javascript:;" @click="transferToBizModule">{{$t('业务模块')}}</a>
        </li>
        <li><a href="javascript:;" @click="transferToResourcePool">{{$t('主机池')}}</a></li>
        <li>
          <a
            v-bk-tooltips="{
              placement: 'right-start',
              content: $t('仅允许业务下空闲机跨业务转移', { idleModule: $store.state.globalConfig.config.idlePool.idle }),
              disabled: !transferToOtherDisabled
            }"
            :class="{ disabled: transferToOtherDisabled }"
            href="javascript:;" @click="transferToOtherBizModule">{{$t('其他业务')}}</a>
        </li>
      </ul>
    </bk-dropdown-menu>
    <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height">
      <component
        :is="dialog.component"
        v-bind="dialog.props"
        @cancel="handleDialogCancel"
        @confirm="handleDialogConfirm">
      </component>
    </cmdb-dialog>
  </div>
</template>

<script>
  import HostStore from './host-store'
  import ModuleSelector from '../../business-topology/host/module-selector'
  import MoveToResourceConfirm from '../../business-topology/host/move-to-resource-confirm'
  import AcrossBusinessModuleSelector from '@/views/business-topology/host/across-business-module-selector.vue'
  import { MENU_BUSINESS_TRANSFER_HOST } from '@/dictionary/menu-symbol'
  import RouterQuery from '@/router/query'
  import uniq from 'lodash.uniq'
  import { MULTI_TO_ONE } from '@/dictionary/host-transfer-type.js'

  export default {
    components: {
      [ModuleSelector.name]: ModuleSelector,
      [MoveToResourceConfirm.name]: MoveToResourceConfirm,
      [AcrossBusinessModuleSelector.name]: AcrossBusinessModuleSelector
    },
    data() {
      return {
        isShow: false,
        dialog: {
          width: 830,
          height: 600,
          show: false,
          component: null,
          props: {}
        },
        request: {
          moveToIdleModule: Symbol('moveToIdleModule'),
          moveToResource: Symbol('moveToResource'),
          moveToOtherBizModule: Symbol('moveToOtherBizModule')
        }
      }
    },
    computed: {
      disabled() {
        return !HostStore.isSelected
      },
      transferToOtherDisabled() {
        return !HostStore.isAllIdleSet || HostStore.hosts.some(host => host.biz[0].default === 1)
      },
    },
    methods: {
      handleMenuToggle(isShow) {
        this.isShow = isShow
      },
      validateSameBiz() {
        if (!HostStore.isSameBiz) {
          this.$error(this.$t('仅支持对相同业务下的主机进行操作'))
          return false
        }
        return true
      },
      transferToIdleModule() {
        const valid = this.validateSameBiz()
        if (!valid) {
          return false
        }
        if (HostStore.isAllResourceHost) {
          this.$error(this.$t('仅支持对业务下的主机进行操作'))
          return false
        }
        const props = {
          moduleType: 'idle',
          business: HostStore.uniqueBusiness,
          title: this.$t('转移主机到空闲模块', { idleSet: this.$store.state.globalConfig.config.set })
        }
        this.dialog.props = props
        this.dialog.width = 830
        this.dialog.height = 600
        this.dialog.component = ModuleSelector.name
        this.dialog.show = true
      },
      transferToBizModule() {
        const valid = this.validateSameBiz()
        if (!valid) {
          return false
        }
        if (HostStore.isAllResourceHost) {
          this.$error(this.$t('仅支持对业务下的主机进行操作'))
          return false
        }
        const props = {
          moduleType: 'business',
          business: HostStore.uniqueBusiness,
          title: this.$t('转移主机到业务模块')
        }
        const selection = HostStore.getSelected()
        const firstSelectionModules = selection[0].module.map(module => module.bk_module_id).sort()
        const firstSelectionModulesStr = firstSelectionModules.join(',')
        const allSame = selection.slice(1).every((item) => {
          const modules = item.module.map(module => module.bk_module_id).sort()
            .join(',')
          return modules === firstSelectionModulesStr
        })
        if (allSame) {
          props.previousModules = firstSelectionModules
        }
        this.dialog.props = props
        this.dialog.width = 830
        this.dialog.height = 600
        this.dialog.component = ModuleSelector.name
        this.dialog.show = true
      },
      transferToResourcePool() {
        const isSameBiz = this.validateSameBiz()
        if (!isSameBiz) {
          return false
        }
        if (HostStore.isAllResourceHost) {
          this.$error('所选主机已在主机池中')
          return false
        }
        const { isAllIdleSet } = HostStore
        if (!isAllIdleSet) {
          this.$error(this.$t('仅支持对空闲机池下的主机进行操作', { idleSet: this.$store.state.globalConfig.config.set }))
          return false
        }
        const [bizId] = HostStore.bizSet
        this.dialog.props = {
          count: HostStore.getSelected().length,
          bizId
        }
        this.dialog.width = 460
        this.dialog.height = 250
        this.dialog.component = MoveToResourceConfirm.name
        this.dialog.show = true
      },
      transferToOtherBizModule() {
        if (this.transferToOtherDisabled) return false
        const moduleBizs = uniq(HostStore.getSelected().map(host => host.biz[0].bk_biz_id))
        const props = {
          moduleType: 'business',
          business: moduleBizs,
          title: this.$t('转移空闲机到其他业务', { idleModule: this.$store.state.globalConfig.config.idlePool.idle }),
          type: MULTI_TO_ONE
        }
        this.dialog.props = props
        this.dialog.width = 830
        this.dialog.height = 600
        this.dialog.component = AcrossBusinessModuleSelector.name
        this.dialog.show = true
      },
      handleDialogCancel() {
        this.dialog.show = false
      },
      handleDialogConfirm(...args) {
        const theArgs = args
        this.dialog.show = false
        if (this.dialog.component === ModuleSelector.name) {
          if (this.dialog.props.moduleType === 'idle') {
            if (HostStore.isAllIdleSet) {
              this.transferDirectly(...theArgs)
            } else {
              this.gotoTransferPage(...theArgs)
            }
          } else {
            this.gotoTransferPage(...theArgs)
          }
        } else if (this.dialog.component === MoveToResourceConfirm.name) {
          this.moveHostToResource(...theArgs)
        } else if (this.dialog.component === AcrossBusinessModuleSelector.name) {
          this.moveHostToOtherBiz(...theArgs)
        }
      },
      async transferDirectly(modules) {
        try {
          const bizId = HostStore.uniqueBusiness.bk_biz_id
          const [internalModule] = modules
          await this.$http.post(`host/transfer_with_auto_clear_service_instance/bk_biz_id/${bizId}`, {
            bk_host_ids: HostStore.getSelected().map(data => data.host.bk_host_id),
            default_internal_module: internalModule.data.bk_inst_id,
            is_remove_from_all: true
          }, {
            requestId: this.request.moveToIdleModule
          })
          HostStore.clear()
          this.$success('转移成功')
          RouterQuery.set({
            _t: Date.now(),
            page: 1
          })
        } catch (e) {
          console.error(e)
        }
      },
      gotoTransferPage(modules) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_TRANSFER_HOST,
          params: {
            bizId: HostStore.uniqueBusiness.bk_biz_id,
            type: this.dialog.props.moduleType
          },
          query: {
            targetModules: modules.map(node => node.data.bk_inst_id).join(','),
            resources: HostStore.getSelected().map(item => item.host.bk_host_id)
              .join(',')
          },
          history: true
        })
        HostStore.clear()
      },
      async moveHostToResource(directoryId) {
        try {
          await this.$store.dispatch('hostRelation/transferHostToResourceModule', {
            params: {
              bk_biz_id: HostStore.uniqueBusiness.bk_biz_id,
              bk_host_id: HostStore.getSelected().map(item => item.host.bk_host_id),
              bk_module_id: directoryId
            },
            config: {
              requestId: this.request.moveToResource
            }
          })
          HostStore.clear()
          this.$success('转移成功')
          RouterQuery.set({
            _t: Date.now(),
            page: 1
          })
        } catch (e) {
          console.error(e)
        }
      },
      moveHostToOtherBiz(selectedNode, targetBizId) {
        const targetModuleId = selectedNode[0].data.bk_inst_id
        const resourceHosts = []

        HostStore.getSelected().forEach((selectedHost) => {
          const hostBizId = selectedHost.biz[0].bk_biz_id
          const hostId = selectedHost.host.bk_host_id
          const existedResourceHost = resourceHosts.find(resourceHost => resourceHost.src_bk_biz_id === hostBizId)

          if (existedResourceHost) {
            existedResourceHost.src_bk_host_ids.push(hostId)
          } else {
            resourceHosts.push({
              src_bk_biz_id: hostBizId,
              src_bk_host_ids: [hostId]
            })
          }
        })

        this.$store.dispatch('hostRelation/transferHostToOtherBizModule', {
          params: {
            resource_hosts: resourceHosts,
            dst_bk_biz_id: targetBizId,
            dst_bk_module_id: targetModuleId
          },
          config: {
            requestId: this.request.moveToOtherBizModule
          }
        })
          .then(() => {
            HostStore.clear()
            this.$success('转移成功')
            RouterQuery.set({
              _t: Date.now(),
              page: 1
            })
          })
      },
    }
  }
</script>

<style lang="scss" scoped>
    .transfer-menu {
        display: inline-block;
    }
    .dropdown-trigger-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #c4c6cc;
        height: 32px;
        min-width: 68px;
        border-radius: 2px;
        padding: 0 15px;
        color: #63656E;
        font-size: 14px;
    }
    .dropdown-trigger-btn.bk-icon {
        font-size: 18px;
    }
    .dropdown-trigger-btn .bk-icon {
        font-size: 22px;
    }
    .dropdown-trigger-btn:hover {
        cursor: pointer;
        border-color: #979ba5;
    }

    .bk-dropdown-list {
      font-size: 14px;
      color: $textColor;

      > li > a {
        &.disabled {
          color: $textDisabledColor;
          cursor: not-allowed;
          &:hover{
            background: inherit;
          }
        }
      }
  }
</style>
