<template>
    <div class="layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
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
                <bk-tab-panel name="serviceInstance" :label="$t('服务实例')" :disabled="!showServiceInstance">
                    <span slot="label"
                        v-bk-tooltips="{
                            content: $t('请选择业务模块'),
                            disabled: showServiceInstance
                        }">
                        {{$t('服务实例')}}
                    </span>
                    <service-instance ref="serviceInstance"></service-instance>
                </bk-tab-panel>
                <bk-tab-panel name="nodeInfo" :disabled="!showNodeInfo">
                    <span slot="label"
                        v-bk-tooltips="{
                            content: $t('请选择非内置节点'),
                            disabled: showNodeInfo
                        }">
                        {{$t('节点信息')}}
                    </span>
                    <service-node-info :active="activeTab === 'nodeInfo'" ref="nodeInfo"></service-node-info>
                </bk-tab-panel>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import TopologyTree from './children/topology-tree.vue'
    import HostList from './host/host-list.vue'
    import ServiceInstance from './children/service-instances.vue'
    import ServiceNodeInfo from './children/service-node-info.vue'
    import { mapGetters } from 'vuex'
    import Bus from '@/utils/bus.js'
    export default {
        components: {
            TopologyTree,
            HostList,
            ServiceNodeInfo,
            ServiceInstance
        },
        data () {
            return {
                activeTab: this.$route.query.tab || 'hostList',
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
            }
        },
        watch: {
            showServiceInstance (value) {
                if (!value && this.activeTab === 'serviceInstance') {
                    this.activeTab = 'hostList'
                }
            },
            showNodeInfo (value) {
                if (!value) {
                    this.activeTab = 'hostList'
                }
            },
            activeTab (tab) {
                const refresh = this.$refs[tab].refresh
                typeof refresh === 'function' && refresh(1)
            }
        },
        async created () {
            try {
                const topologyModels = await this.getTopologyModels()
                const properties = await this.getProperties(topologyModels)
                this.$store.commit('businessHost/setTopologyModels', topologyModels)
                this.$store.commit('businessHost/setPropertyMap', properties)
                this.$store.commit('businessHost/resolveCommonRequest')
            } catch (e) {
                console.error(e)
            }
        },
        beforeDestroy () {
            this.$store.commit('businessHost/clear')
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
                    params: this.$injectMetadata({
                        bk_obj_id: {
                            $in: models.map(model => model.bk_obj_id)
                        },
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: this.request.properties
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        border-top: 1px solid $cmdbLayoutBorderColor;
    }
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
            font-size: 12px;
            color: #fff;
            cursor: pointer;
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
</style>
