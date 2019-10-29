<template>
    <div class="layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <cmdb-resize-layout class="resize-layout fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480">
            <topology-tree ref="topologyTree"></topology-tree>
        </cmdb-resize-layout>
        <div class="tab-layout">
            <bk-tab class="topology-tab" type="unborder-card">
                <bk-tab-panel name="hostList" :label="$t('主机列表')">
                    <host-list></host-list>
                </bk-tab-panel>
                <bk-tab-panel name="serviceInstance" :label="$t('服务实例')"></bk-tab-panel>
                <bk-tab-panel name="nodeInfo" :label="$t('节点信息')"></bk-tab-panel>
            </bk-tab>
        </div>
    </div>
</template>

<script>
    import TopologyTree from './children/topology-tree.vue'
    import HostList from './children/host-list.vue'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            TopologyTree,
            HostList
        },
        data () {
            return {
                show: false,
                request: {
                    mainline: Symbol('mainline'),
                    properties: Symbol('properties')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId'])
        },
        async created () {
            try {
                const topologyModels = await this.getTopologyModels()
                const properties = await this.getProperties(topologyModels)
                this.$store.commit('businessHost/setTopologyModels', topologyModels)
                this.$store.commit('businessHost/setProperties', properties)
                this.$store.commit('businessHost/resolveCommonRequest')
            } catch (e) {
                console.error(e)
            }
        },
        beforeDestroy () {
            this.$store.commit('businessHost/clear')
        },
        methods: {
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
        width: 280px;
        height: 100%;
        padding-top: 10px;
        border-right: 1px solid $cmdbLayoutBorderColor;
    }
    .tab-layout {
        height: calc(100vh - 140px);
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
