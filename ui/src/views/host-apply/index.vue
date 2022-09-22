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
  <div class="host-apply">
    <div class="main-wrapper">
      <cmdb-resize-layout
        ref="resizeLayout"
        :class="['tree-layout fl', { 'is-collapse': layout.topologyCollapse }]"
        direction="right"
        :handler-offset="3"
        :width="layout.sidebarWidth"
        :min="310"
        :max="480"
        :disabled="layout.topologyCollapse">

        <sidebar ref="sidebar"
          @module-selected="handleSelectModule"
          @mode-changed="handleChangeMode"
          @action-change="handleActionChange">
        </sidebar>

        <i class="topology-collapse-icon bk-icon icon-angle-left"
          @click="layout.topologyCollapse = !layout.topologyCollapse">
        </i>
      </cmdb-resize-layout>
      <div class="main-content" v-bkloading="{ isLoading: $loading([requestIds.rules, requestIds.properties]) }">
        <config-details
          ref="details"
          v-show="!batchAction"
          :id="targetId"
          :biz-id="bizId"
          :rule-list="initRuleList"
          :has-rule="hasRule"
          :current-node="currentNode"
          :conflict-num="conflictNum"
          :checked-property-id-list="checkedPropertyIdList"
          @edit="handleEdit"
          @view-conflict="handleViewConflict"
          @close="handleCloseApply"
          @delete-rule="handleDeleteRule">
        </config-details>
      </div>
    </div>
  </div>
</template>

<script>
  /* eslint-disable no-underscore-dangle */
  import { mapGetters } from 'vuex'
  import sidebar from './children/sidebar.vue'
  import Bus from '@/utils/bus'
  import {
    MENU_BUSINESS_HOST_APPLY_EDIT,
    MENU_BUSINESS_HOST_APPLY_CONFLICT
  } from '@/dictionary/menu-symbol'
  import { CONFIG_MODE } from '@/service/service-template/index.js'
  import configDetails from './children/details.vue'

  export default {
    components: {
      sidebar,
      configDetails
    },
    data() {
      return {
        currentNode: {},
        initRuleList: [],
        checkedPropertyIdList: [],
        conflictNum: 0,
        clearRules: false,
        hasRule: false,
        layout: {
          topologyCollapse: false,
          sidebarWidth: undefined
        },
        requestIds: {
          rules: Symbol('rules'),
          properties: Symbol('properties'),
          conflictCount: Symbol('conflictCount'),
          setEnableStatus: Symbol('setEnableStatus'),
          del: Symbol('del')
        },
        keepAliveQueryIds: {
          [CONFIG_MODE.MODULE]: undefined,
          [CONFIG_MODE.TEMPLATE]: undefined,
        },
        batchAction: '',
        targetId: null
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      mode() {
        return this.$route.params.mode
      },
      isModuleMode() {
        return this.mode === CONFIG_MODE.MODULE
      },
      requestConfigs() {
        return {
          [this.requestIds.rules]: {
            [CONFIG_MODE.MODULE]: {
              action: 'hostApply/getRules',
              payload: {
                bizId: this.bizId,
                params: {
                  bk_module_ids: [this.targetId]
                }
              }
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/getTemplateRules',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  service_template_ids: [this.targetId]
                }
              }
            }
          },
          [this.requestIds.conflictCount]: {
            [CONFIG_MODE.MODULE]: {
              action: 'hostApply/getConflictCount',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  id: this.targetId
                }
              }
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/getTemplateConflictCount',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  id: this.targetId
                }
              }
            }
          },
          [this.requestIds.setEnableStatus]: {
            [CONFIG_MODE.MODULE]: {
              action: 'hostApply/setEnableStatus'
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/setTemplateEnableStatus'
            }
          },
          [this.requestIds.del]: {
            [CONFIG_MODE.MODULE]: {
              action: 'hostApply/deleteRules'
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/deleteTemplateRules'
            }
          }
        }
      },
      targetIdsKey() {
        const targetIdsKeys = {
          [CONFIG_MODE.MODULE]: 'bk_module_ids',
          [CONFIG_MODE.TEMPLATE]: 'service_template_ids'
        }
        return targetIdsKeys[this.mode]
      }
    },
    watch: {
      currentNode() {
        this.getData()
      }
    },
    created() {
      this.getHostPropertyList()
    },
    mounted() {
      const $resizeLayoutEl = this.$refs.resizeLayout.$el
      Bus.$on('topologyTree/expandChange', (lastNodeLevel) => {
        const { offsetWidth } = $resizeLayoutEl
        const visibleWidth = (lastNodeLevel + 2) * 30 + 100
        if (offsetWidth < visibleWidth) {
          this.layout.sidebarWidth = visibleWidth
        }
      })
    },
    methods: {
      async getData() {
        // 重置配置表格数据
        if (this.$refs.details) {
          this.$refs.details.reset()
        }

        // 业务拓扑模式且模板配置了自动应用，则不需要请求数据，但需要将规则数据重置
        if (this.isModuleMode && this.currentNode.service_template_host_apply_enabled) {
          this.initRuleList = []
          this.hasRule = false
          this.checkedPropertyIdList = []
          return
        }

        try {
          const ruleData = await this.getRules()

          this.initRuleList = ruleData.info || []
          this.hasRule = ruleData.count > 0
          this.checkedPropertyIdList = this.initRuleList.map(item => item.bk_attribute_id)

          if (this.currentNode.host_apply_enabled) {
            const { count } = await this.getConflictCount()
            this.conflictNum = count
          }
        } catch (e) {
          console.log(e)
        }
      },
      async getHostPropertyList() {
        try {
          const properties = await this.$store.dispatch('hostApply/getProperties', {
            params: { bk_biz_id: this.bizId },
            config: {
              requestId: this.requestIds.properties,
              fromCache: true
            }
          })
          this.$store.commit('hostApply/setPropertyList', properties)
        } catch (e) {
          console.error(e)
        }
      },
      getRules() {
        const requestConfig = this.requestConfigs[this.requestIds.rules][this.mode]
        return this.$store.dispatch(requestConfig.action, {
          config: {
            requestId: this.requestIds.rules
          },
          ...requestConfig.payload
        })
      },
      getConflictCount() {
        const requestConfig = this.requestConfigs[this.requestIds.conflictCount][this.mode]
        return this.$store.dispatch(requestConfig.action, requestConfig.payload)
      },
      emptyRules() {
        this.checkedPropertyIdList = []
        this.hasRule = false
      },
      handleCloseApply() {
        const h = this.$createElement
        this.$bkInfo({
          title: this.$t('确认关闭'),
          extCls: 'close-apply-confirm-modal',
          subHeader: h('div', { class: 'content' }, [
            h('p', { class: 'tips' }, this.$t('确认关闭当前模块的主机属性自动应用')),
            h('bk-checkbox', {
              props: {
                checked: true,
                trueValue: true,
                falseFalue: false
              },
              on: {
                change: value => (this.clearRules = !value)
              }
            }, this.$t('保留当前自动应用策略'))
          ]),
          confirmFn: async () => {
            const requestConfig = this.requestConfigs[this.requestIds.setEnableStatus][this.mode]
            try {
              await this.$store.dispatch(requestConfig.action, {
                bizId: this.bizId,
                params: {
                  ids: [this.targetId],
                  enabled: false,
                  clear_rules: this.clearRules
                }
              })

              this.$success(this.$t('关闭成功'))
              if (this.clearRules) {
                this.emptyRules()
              }
              this.$refs.sidebar.setApplyClosed(this.targetId, { isClose: true, isClear: this.clearRules })

              // 更新当前节点开启状态
              this.currentNode.host_apply_enabled = false

              Bus.$emit('host-apply-closed', this.mode, this.targetId, this.clearRules)
            } catch (e) {
              console.log(e)
            }
          }
        })
      },
      handleViewConflict() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY_CONFLICT,
          query: {
            id: this.targetId
          },
          history: true
        })
      },
      handleEdit() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY_EDIT,
          params: {
            mode: this.mode
          },
          query: {
            id: this.targetId,
            // 为统一UI交互全部使用batch模式
            batch: 1,
            action: 'batch-edit'
          },
          history: true
        })
      },
      handleSelectModule(data) {
        this.currentNode = { ...data }

        const modeConfig = {
          [CONFIG_MODE.MODULE]: {
            targetId: data.bk_inst_id
          },
          [CONFIG_MODE.TEMPLATE]: {
            targetId: data.id
          }
        }

        this.targetId = modeConfig[this.mode].targetId

        // 同步参数到路由中
        this.$router.push({
          query: { id: this.targetId }
        })

        // 按模式记录当前id，用于切换时还原
        this.keepAliveQueryIds[this.mode] = this.targetId
      },
      handleChangeMode(mode) {
        this.$router.push({
          params: {
            mode
          },
          query: {
            id: this.keepAliveQueryIds[mode] ?? undefined
          }
        })
      },
      handleActionChange(action) {
        this.batchAction = action
      },
      handleDeleteRule(property) {
        const deleteRule = async () => {
          const requestConfig = this.requestConfigs[this.requestIds.del][this.mode]
          try {
            await this.$store.dispatch(requestConfig.action, {
              bizId: this.bizId,
              params: {
                data: {
                  host_apply_rule_ids: [property.__extra__.ruleId],
                  [this.targetIdsKey]: [this.targetId]
                }
              }
            })

            this.$success(this.$t('删除成功'))

            // 从当前列表中移除
            const checkedIndex = this.checkedPropertyIdList.findIndex(id => id === property.id)
            this.checkedPropertyIdList.splice(checkedIndex, 1)
            const closed = this.checkedPropertyIdList.length === 0

            // 更新配置树状态
            this.$refs.sidebar.setApplyClosed(this.targetId, { isClose: closed, isClear: closed })

            // 更新当前节点开启状态
            this.currentNode.host_apply_enabled = !closed

            // 全部删除后重新获取一次数据
            if (closed) {
              this.getData()
            }
          } catch (error) {
            console.log(error)
          }
        }

        if (this.checkedPropertyIdList.length > 1) {
          deleteRule()
          return
        }

        this.$bkInfo({
          title: this.$t('确认删除自动应用字段？'),
          subTitle: this.$t('自动应用字段全部删除后将关闭自动应用'),
          confirmFn: deleteRule
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .main-wrapper {
    height: 100%;
  }

  .tree-layout {
    width: 310px;
    height: 100%;
    border-right: 1px solid $cmdbLayoutBorderColor;
    z-index: 9999;

    &.is-collapse {
      width: 0 !important;
      border-right: none;
      .topology-collapse-icon:before {
        display: inline-block;
        transform: rotate(180deg);
      }

      .host-apply-sidebar {
        display: none;
      }
    }
    .topology-collapse-icon {
      position: absolute;
      left: 100%;
      top: 50%;
      width: 16px;
      height: 100px;
      line-height: 100px;
      background: $cmdbLayoutBorderColor;
      border-radius: 0px 12px 12px 0px;
      transform: translateY(-50%);
      text-align: center;
      font-size: 20px;
      color: #fff;
      cursor: pointer;
      text-indent: -2px;
      &:hover {
        background: #699DF4;
      }
    }
  }

  .main-content {
    @include scrollbar-y;
    height: 100%;
    padding: 0 20px;
  }

  [bk-language="en"] {
    .tree-layout {
      width: 370px;
    }
  }
</style>
