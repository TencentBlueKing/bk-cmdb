<template>
    <div class="details-layout">
        <cmdb-host-info
            ref="info"
            @info-toggle="setInfoHeight">
        </cmdb-host-info>
        <bk-tab class="details-tab"
            type="unborder-card"
            :active.sync="active"
            :style="{
                '--infoHeight': infoHeight
            }">
            <bk-tab-panel name="property" :label="$t('主机属性')">
                <cmdb-host-property></cmdb-host-property>
            </bk-tab-panel>
            <bk-tab-panel name="association" :label="$t('关联')">
                <cmdb-host-association v-if="active === 'association'"></cmdb-host-association>
            </bk-tab-panel>
            <bk-tab-panel name="status" :label="$t('实时状态')">
                <cmdb-host-status v-if="active === 'status'"></cmdb-host-status>
            </bk-tab-panel>
            <bk-tab-panel name="service" :label="$t('服务列表')" :visible="!isAdminView">
                <cmdb-host-service v-if="active === 'service'"></cmdb-host-service>
            </bk-tab-panel>
            <bk-tab-panel name="history" :label="$t('变更记录')">
                <cmdb-host-history v-if="active === 'history'"></cmdb-host-history>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import { mapState, mapGetters } from 'vuex'
    import cmdbHostInfo from './children/info.vue'
    import cmdbHostAssociation from './children/association.vue'
    import cmdbHostProperty from './children/property.vue'
    import cmdbHostStatus from './children/status.vue'
    import cmdbHostHistory from './children/history.vue'
    import cmdbHostService from './children/service-list.vue'
    import {
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_RESOURCE_HOST
    } from '@/dictionary/menu-symbol'
    export default {
        components: {
            cmdbHostInfo,
            cmdbHostAssociation,
            cmdbHostProperty,
            cmdbHostStatus,
            cmdbHostHistory,
            cmdbHostService
        },
        data () {
            return {
                active: 'property',
                infoHeight: '81px'
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            ...mapGetters(['isAdminView']),
            id () {
                return parseInt(this.$route.params.id)
            },
            business () {
                const business = parseInt(this.$route.params.business)
                if (isNaN(business)) {
                    return -1
                }
                return business
            }
        },
        watch: {
            info (info) {
                const hostList = info.host.bk_host_innerip.split(',')
                const host = hostList.length > 1 ? `${hostList[0]}...` : hostList[0]
                this.setBreadcrumbs(host)
            },
            id () {
                this.getData()
            },
            business () {
                this.getData()
            },
            active (active) {
                if (active !== 'association') {
                    this.$store.commit('hostDetails/toggleExpandAll', false)
                }
            }
        },
        created () {
            this.getData()
        },
        methods: {
            setBreadcrumbs (ip) {
                const isFromBusiness = this.$route.query.from === 'business'
                this.$store.commit('setBreadcrumbs', [{
                    label: isFromBusiness ? this.$t('业务主机') : this.$t('主机'),
                    route: {
                        name: isFromBusiness ? MENU_BUSINESS_HOST_AND_SERVICE : MENU_RESOURCE_HOST,
                        query: {
                            node: isFromBusiness ? this.$route.query.node : undefined
                        }
                    }
                }, {
                    label: ip
                }])
            },
            getData () {
                this.getProperties()
                this.getPropertyGroups()
                this.getHostInfo()
            },
            async getHostInfo () {
                try {
                    const { info } = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.getSearchHostParams()
                    })
                    if (info.length) {
                        this.$store.commit('hostDetails/setHostInfo', info[0])
                    } else {
                        this.$router.replace({ name: 404 })
                    }
                } catch (e) {
                    console.error(e)
                    this.$store.commit('hostDetails/setHostInfo', null)
                }
            },
            getSearchHostParams () {
                const hostCondition = {
                    field: 'bk_host_id',
                    operator: '$eq',
                    value: this.id
                }
                return this.$injectMetadata({
                    bk_biz_id: this.business,
                    condition: ['biz', 'set', 'module', 'host'].map(model => {
                        return {
                            bk_obj_id: model,
                            condition: model === 'host' ? [hostCondition] : [],
                            fields: []
                        }
                    }),
                    ip: { flag: 'bk_host_innerip', exact: 1, data: [] }
                })
            },
            async getProperties () {
                try {
                    const properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: this.$injectMetadata({
                            bk_obj_id: 'host'
                        })
                    })
                    this.$store.commit('hostDetails/setHostProperties', properties)
                } catch (e) {
                    console.error(e)
                    this.$store.commit('hostDetails/setHostProperties', [])
                }
            },
            async getPropertyGroups () {
                try {
                    const propertyGroups = await this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                        objId: 'host',
                        params: this.$injectMetadata()
                    })
                    this.$store.commit('hostDetails/setHostPropertyGroups', propertyGroups)
                } catch (e) {
                    console.error(e)
                    this.$store.commit('hostDetails/setHostPropertyGroups', [])
                }
            },
            setInfoHeight (height) {
                this.infoHeight = height + 'px'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-layout {
        overflow: hidden;
        .details-tab {
            height: calc(100% - var(--infoHeight)) !important;
            min-height: 400px;
            /deep/ {
                .bk-tab-header {
                    padding: 0;
                    margin: 0 20px;
                }
                .bk-tab-section {
                    @include scrollbar-y;
                    padding-bottom: 10px;
                }
            }
        }
    }
</style>
