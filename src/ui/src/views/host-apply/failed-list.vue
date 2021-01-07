<template>
    <div class="failed-list">
        <div class="caption">
            <div class="title">{{$t('请确认以下主机应用信息')}}</div>
        </div>
        <bk-table
            :data="table.list"
            :pagination="table.pagination"
            :row-style="{ cursor: 'pointer' }"
            :max-height="$APP.height - 240"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @row-click="handleRowClick"
        >
            <bk-table-column :label="$t('内网IP')" prop="host.bk_host_innerip" class-name="is-highlight" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('云区域')" show-overflow-tooltip>
                <template slot-scope="{ row }">
                    {{getCloudName(row.host.bk_cloud_id)}}
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('固资编号')" prop="host.bk_asset_id" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('主机名称')" prop="host.bk_host_name" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('所属拓扑')" :formatter="getTopopath" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('失败原因')">
                <div class="fail-reason" slot-scope="{}">
                    网络中断
                </div>
            </bk-table-column>
        </bk-table>
        <div class="bottom-actionbar">
            <div class="actionbar-inner">
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                    <bk-button
                        theme="primary"
                        slot-scope="{ disabled }"
                        :disabled="disabled"
                        @click="handleRetry"
                    >
                        {{$t('重试')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button theme="default" @click="handleCopyIp">{{$t('复制IP')}}</bk-button>
                <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            </div>
        </div>
        <apply-status-modal
            ref="applyStatusModal"
            :request="applyRequest"
            @return="handleStatusModalBack"
            @view-host="handleViewHost"
            @view-failed="handleViewFailed"
        >
        </apply-status-modal>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :is-show.sync="slider.isShow"
            :title="slider.title"
            @hidden="handleSliderCancel"
        >
            <template slot="content">
                <cmdb-details
                    :show-options="false"
                    :inst="details.inst"
                    :properties="details.properties"
                    :property-groups="details.propertyGroups">
                </cmdb-details>
            </template>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapState, mapActions } from 'vuex'
    import applyStatusModal from './children/apply-status'
    import {
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_BUSINESS_HOST_APPLY,
        MENU_BUSINESS_HOST_APPLY_FAILED
    } from '@/dictionary/menu-symbol'
    export default {
        components: {
            applyStatusModal
        },
        data () {
            return {
                applyRequest: null,
                table: {
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    }
                },
                details: {
                    show: false,
                    title: '',
                    inst: {},
                    properties: [],
                    propertyGroups: []
                },
                slider: {
                    width: 514,
                    isShow: false,
                    content: '',
                    title: ''
                }
            }
        },
        computed: {
            ...mapState('hostApply', ['propertyConfig']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('hosts', ['propertyList']),
            hostIds () {
                return this.propertyConfig.bk_host_ids || []
            },
            isBatch () {
                return this.$route.query.batch === 1
            },
            moduleId () {
                const mid = this.$route.query.mid
                let moduleId
                if (mid) {
                    moduleId = Number(mid)
                }
                return moduleId
            }
        },
        created () {
            if (!Object.keys(this.propertyConfig).length) {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY
                })
            } else {
                this.setBreadcrumbs()
                this.getData()
            }
        },
        methods: {
            ...mapGetters('objectModelClassify', [
                'getModelById'
            ]),
            ...mapActions('hostApply', [
                'runApply'
            ]),
            getData () {
                this.getHostList()
            },
            setBreadcrumbs () {
                this.$store.commit('setTitle', this.$t('失败列表'))
            },
            async getHostList () {
                try {
                    const { info, count } = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.getSearchHostParams()
                    })
                    this.table.list = info
                    this.table.pagination.count = count
                } catch (e) {
                    console.error(e)
                }
            },
            getSearchHostParams () {
                const hostCondition = {
                    field: 'bk_host_id',
                    operator: '$in',
                    value: this.hostIds
                }
                return {
                    bk_biz_id: this.bizId,
                    condition: ['biz', 'set', 'module', 'host'].map(model => {
                        return {
                            bk_obj_id: model,
                            condition: model === 'host' ? [hostCondition] : [],
                            fields: []
                        }
                    }),
                    ip: { flag: 'bk_host_innerip', exact: 1, data: [] }
                }
            },
            getTopopath (info) {
                const topology = []
                const sets = info.set || []
                const modules = info.module || []
                const businesses = info.biz || []
                modules.forEach(module => {
                    const set = sets.find(set => set.bk_set_id === module.bk_set_id)
                    if (set) {
                        const business = businesses.find(business => business.bk_biz_id === set.bk_biz_id)
                        if (business) {
                            topology.push({
                                id: module.bk_module_id,
                                path: `${business.bk_biz_name} / ${set.bk_set_name} / ${module.bk_module_name}`,
                                isInternal: module.default !== 0
                            })
                        }
                    }
                })
                return topology.map(item => item.path).join('\n')
            },
            getPropertyGroups () {
                return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                    objId: 'host',
                    params: {
                        bk_biz_id: this.bizId
                    }
                })
            },
            getCloudName (cloud) {
                const names = cloud.map(item => item.bk_inst_name)
                return names.join(',')
            },
            setTableList () {
                const { start, limit } = this.$tools.getPageParams(this.table.pagination)
                this.table.list = this.table.list.slice(start, start + limit)
            },
            goBack () {
                this.$store.commit('hostApply/clearRuleDraft')
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY
                })
            },
            async handleRetry () {
                this.applyRequest = this.runApply({
                    bizId: this.bizId,
                    params: this.propertyConfig,
                    config: {
                        requestId: 'runHostApply'
                    }
                })
                this.$refs.applyStatusModal.show()
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.setTableList()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.setTableList()
            },
            handleRowClick (row) {
                this.handleShowDetails(row)
            },
            async handleShowDetails (row) {
                this.slider.title = `${this.$t('属性详情')}【${row.bk_host_innerip}】`
                this.slider.content = 'detail'
                const properties = this.propertyList
                const inst = row
                try {
                    const propertyGroups = await this.getPropertyGroups()
                    this.details.inst = inst
                    this.details.properties = properties
                    this.details.propertyGroups = propertyGroups
                    this.slider.isShow = true
                } catch (e) {
                    console.log(e)
                    this.details.inst = {}
                    this.details.properties = []
                    this.details.propertyGroups = []
                    this.slider.isShow = false
                }
            },
            handleSliderCancel () {
                this.slider.isShow = false
            },
            handleStatusModalBack () {
                this.goBack()
            },
            handleCancel () {
                this.goBack()
            },
            handleViewHost () {
                const query = {}
                if (!this.isBatch && this.moduleId) {
                    query.node = `module-${this.moduleId}`
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query,
                    history: true
                })
            },
            handleViewFailed () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY_FAILED,
                    query: this.$route.query,
                    history: true
                })
            },
            handleCopyIp () {
                const ips = this.table.list.map(item => item.host.bk_host_innerip)
                this.$copyText(ips.join('\n')).then(() => {
                    this.$success(this.$t('复制成功'))
                }, () => {
                    this.$error(this.$t('复制失败'))
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .failed-list {
        padding: 15px 20px 0;

        .caption {
            display: flex;
            margin-bottom: 14px;
            justify-content: space-between;
            align-items: center;
        }
    }

    .fail-reason {
        color: #ff5656;
    }

    .bottom-actionbar {
        position: absolute;
        width: 100%;
        height: 50px;
        border-top: 1px solid #dcdee5;
        bottom: 0;
        left: 0;

        .actionbar-inner {
            padding: 8px 0 0 20px;

            .bk-button {
                min-width: 86px;
            }
        }
    }
</style>
