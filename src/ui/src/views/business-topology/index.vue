<template>
    <div class="layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }" style="overflow: hidden;">
        <cmdb-resize-layout :class="['resize-layout fl', { 'is-collapse': layout.topologyCollapse }]"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480"
            :disabled="layout.topologyCollapse">
            <topology-tree ref="topologyTree" :active="activeTab"></topology-tree>
            <i class="topology-collapse-icon bk-icon icon-angle-left"
                @click="layout.topologyCollapse = !layout.topologyCollapse">
            </i>
        </cmdb-resize-layout>
        <div class="tab-layout">
            <bk-tab class="topology-tab" type="unborder-card"
                :active.sync="activeTab"
                :validate-active="false"
                :before-toggle="handleTabToggle">
                <bk-tab-panel name="hostList" :label="$t('主机列表')">
                    <host-list :active="activeTab === 'hostList'" ref="hostList"></host-list>
                </bk-tab-panel>
                <bk-tab-panel name="serviceInstance" :label="$t('服务实例')">
                    <div class="non-business-module" v-if="!showServiceInstance">
                        <div class="tips">
                            <i class="bk-cc-icon icon-cc-tips"></i>
                            <span>{{$t('非业务模块，无服务实例，请选择业务模块查看')}}</span>
                        </div>
                    </div>
                    <service-instance-view v-else-if="activeTab === 'serviceInstance'"></service-instance-view>
                </bk-tab-panel>
                <bk-tab-panel name="nodeInfo" :label="$t('节点信息')">
                    <div class="default-node-info" v-if="!showNodeInfo">
                        <div class="info-item">
                            <label class="name">{{$t('ID')}}:</label>
                            <span class="value">{{nodeId}}</span>
                        </div>
                        <div class="info-item">
                            <label class="name">{{$t('节点名称')}}</label>
                            <span class="value">{{nodeName}}</span>
                        </div>
                    </div>
                    <service-node-info v-else :active="activeTab === 'nodeInfo'" ref="nodeInfo"></service-node-info>
                </bk-tab-panel>
            </bk-tab>
        </div>
        <router-subview></router-subview>
    </div>
</template>

<script>
    import TopologyTree from './children/topology-tree.vue'
    import HostList from './host/host-list.vue'
    import ServiceNodeInfo from './children/service-node-info.vue'
    import { mapGetters } from 'vuex'
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    import ServiceInstanceView from './service-instance/view'
    export default {
        components: {
            TopologyTree,
            HostList,
            ServiceNodeInfo,
            ServiceInstanceView
        },
        data () {
            return {
                activeTab: RouterQuery.get('tab', 'hostList'),
                layout: {
                    topologyCollapse: false
                },
                request: {
                    mainline: Symbol('mainline'),
                    properties: Symbol('properties')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['selectedNode']),
            showServiceInstance () {
                return this.selectedNode && this.selectedNode.data.bk_obj_id === 'module' && this.selectedNode.data.default === 0
            },
            showNodeInfo () {
                return this.selectedNode && this.selectedNode.data.default === 0
            },
            nodeId () {
                return this.selectedNode ? this.selectedNode.data.bk_inst_id : '--'
            },
            nodeName () {
                return this.selectedNode && this.selectedNode.data.bk_inst_name
            }
        },
        watch: {
            activeTab (tab) {
                this.$nextTick(() => {
                    RouterQuery.set({
                        tab: tab,
                        _t: Date.now(),
                        page: '',
                        limit: ''
                    })
                })
            }
        },
        async created () {
            this.unwatch = RouterQuery.watch('tab', (value = 'hostList') => {
                this.activeTab = value
            })
            try {
                const topologyModels = await this.getTopologyModels()
                const properties = await this.getProperties(topologyModels)
                this.$store.commit('businessHost/setTopologyModels', topologyModels)
                this.$store.commit('businessHost/setPropertyMap', Object.freeze(properties))
                this.$store.commit('businessHost/resolveCommonRequest')
            } catch (e) {
                console.error(e)
            }
        },
        beforeDestroy () {
            this.$store.commit('businessHost/clear')
            this.unwatch()
        },
        methods: {
            handleTabToggle () {
                Bus.$emit('toggle-host-filter', false)
                return true
            },
            getTopologyModels () {
                return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                    config: {
                        requestId: this.request.mainline
                    }
                })
            },
            getProperties (models) {
                return this.$store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
                    injectId: 'host',
                    params: {
                        bk_biz_id: this.bizId,
                        bk_obj_id: {
                            $in: models.map(model => model.bk_obj_id)
                        },
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: this.request.properties
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .resize-layout {
        position: relative;
        width: 286px;
        height: 100%;
        padding-top: 10px;
        border-right: 1px solid $cmdbLayoutBorderColor;
        &.is-collapse {
            width: 0 !important;
            border-right: none;
            .topology-collapse-icon:before {
                display: inline-block;
                transform: rotate(180deg);
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
    .tab-layout {
        height: 100%;
        overflow: hidden;
        .topology-tab {
            /deep/ {
                .bk-tab-header {
                    padding: 0;
                    margin: 0 20px;
                }
            }
        }
    }

    .non-business-module {
        display: flex;
        height: 80%;
        justify-content: center;
        align-items: center;
        .tips {
            font-size: 14px;
            .bk-cc-icon {
                font-size: 16px;
                margin-top: -2px;
            }
        }
    }
    .default-node-info {
        padding: 20px 0 20px 36px;
        display: flex;
        .info-item {
            flex: auto;
            max-width: 400px;
            font-size: 14px;
            .name {
                color: #63656e;
            }
            .value {
                color: #313238;
            }
        }
    }
</style>
