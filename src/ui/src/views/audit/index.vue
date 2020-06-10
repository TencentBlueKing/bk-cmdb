<template>
    <div class="audit-layout">
        <bk-tab
            class="audit-tab"
            type="unborder-card"
            :active.sync="active"
            @tab-change="handleTabChange">
            <bk-tab-panel
                v-for="panel in tabPanels"
                v-bind="panel"
                :key="panel.id">
            </bk-tab-panel>
        </bk-tab>
        <div class="filter">
            <div class="option" v-if="active === 'business'">
                <span class="name" :title="$t('业务')">{{$t('业务')}}</span>
                <div class="content">
                    <bk-select v-model="filter.bizId" searchable font-size="medium">
                        <bk-option v-for="business in authorizedBusiness"
                            :key="business.bk_biz_id"
                            :id="business.bk_biz_id"
                            :name="business.bk_biz_name">
                        </bk-option>
                    </bk-select>
                </div>
            </div>
            <div class="option action">
                <span class="name" :title="$t('功能板块动作')">{{$t('功能板块动作')}}</span>
                <div class="content">
                    <bk-select
                        searchable
                        multiple
                        v-model="filter.action"
                        font-size="medium"
                        :popover-min-width="160">
                        <bk-option-group
                            v-for="module in actionList"
                            :key="module.id"
                            :id="module.id"
                            :name="$t(module.name)">
                            <bk-option
                                v-for="option in module.operations"
                                :key="option.id"
                                :id="option.id"
                                :name="$t(option.name)">
                            </bk-option>
                        </bk-option-group>
                    </bk-select>
                </div>
            </div>
            <div class="option resource">
                <span class="name" :title="$t('操作对象')">{{$t('操作对象')}}</span>
                <div class="content">
                    <bk-input v-model="filter.resourceName" clearable></bk-input>
                </div>
            </div>
            <div class="option">
                <span class="name" :title="$t('时间')">{{$t('时间')}}</span>
                <div class="content">
                    <cmdb-form-date-range
                        class="date-range"
                        :clearable="false"
                        font-size="medium"
                        v-model="filter.time">
                    </cmdb-form-date-range>
                </div>
            </div>
            <div class="option operator">
                <span class="name" :title="$t('操作账号')">{{$t('操作账号')}}</span>
                <div class="content">
                    <cmdb-form-objuser
                        v-model="filter.operator"
                        :exclude="false"
                        :multiple="false"
                        :palceholder="$t('操作账号')">
                    </cmdb-form-objuser>
                </div>
            </div>
            <div class="option instance">
                <span class="name" :title="$t('实例ID')">{{$t('实例ID')}}</span>
                <div class="content">
                    <bk-select
                        searchable
                        v-model="filter.module"
                        font-size="medium"
                        :popover-min-width="120">
                        <bk-option v-for="module in operationTypes"
                            :key="module.id"
                            :id="module.id"
                            :name="module.name">
                        </bk-option>
                    </bk-select>
                    <bk-input v-model="filter.resourceId" :placeholder="$t('请输入xx', { name: 'ID' })" :clearable="true"></bk-input>
                </div>
            </div>
            <div class="option option-btn">
                <bk-button theme="primary" :loading="$loading('getOperationLog')" @click="handlePageChange(1, $event)">{{$t('查询')}}</bk-button>
                <bk-button @click="handleClearFilter">{{$t('清空')}}</bk-button>
            </div>
        </div>
        <bk-table
            v-bkloading="{ isLoading: $loading('getOperationLog') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 310"
            :row-style="{ cursor: 'pointer' }"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange"
            @row-click="handleRowClick">
            <bk-table-column
                prop="resource_type"
                class-name="is-highlight"
                :label="$t('功能板块')"
                :formatter="getResourceType">
            </bk-table-column>
            <bk-table-column
                prop="action"
                :label="$t('动作')"
                :formatter="getResourceAction">
            </bk-table-column>
            <bk-table-column
                v-if="active === 'business'"
                prop="bk_biz_name"
                :label="$t('所属业务')">
                <template slot-scope="{ row }">{{getBusinessName(row)}}</template>
            </bk-table-column>
            <bk-table-column
                prop="resource_name"
                :label="$t('操作实例')"
                :formatter="getResourceName">
            </bk-table-column>
            <bk-table-column
                :label="$t('操作描述')">
                <template slot-scope="{ row }">
                    {{`${getResourceAction(row)}"${getTargetName(row)}"`}}
                </template>
            </bk-table-column>
            <bk-table-column
                width="160"
                sortable="custom"
                prop="operation_time"
                :label="$t('操作时间')">
                <template slot-scope="{ row }">{{$tools.formatTime(row.operation_time)}}</template>
            </bk-table-column>
            <bk-table-column
                prop="user"
                :label="$t('操作账号')">
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :quick-close="true"
            :is-show.sync="details.isShow"
            :width="800"
            :title="$t('操作详情')">
            <template slot="content" v-if="details.isShow">
                <v-details
                    v-if="details.showDetailsList.includes(details.data.audit_type)"
                    :show-business="active === 'business'"
                    :details="details.data">
                </v-details>
                <v-json-details v-else :details="details.data"></v-json-details>
            </template>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapState } from 'vuex'
    import moment from 'moment'
    import vDetails from '@/components/audit-history/details'
    import vJsonDetails from '@/components/audit-history/details-json'
    export default {
        components: {
            vDetails,
            vJsonDetails
        },
        data () {
            // 受审计日志性能影响，限制最多显示100条
            const paginationConfig = this.$tools.getDefaultPaginationConfig()
            paginationConfig['limit-list'] = paginationConfig['limit-list'].filter(limit => limit !== 500)
            return {
                active: 'business',
                tabPanels: [
                    {
                        name: 'business',
                        label: this.$t('业务')
                    }, {
                        name: 'resource',
                        label: this.$t('资源')
                    }, {
                        name: 'other',
                        label: this.$t('其他'),
                        visible: false
                    }
                ],
                filter: {
                    bizId: '',
                    action: [],
                    resourceName: '',
                    operator: '',
                    instanceId: '',
                    module: '',
                    time: [
                        this.$tools.formatTime(moment().subtract(1, 'days'), 'YYYY-MM-DD'),
                        this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                    ]
                },
                table: {
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...paginationConfig
                    },
                    defaultSort: '-operation_time',
                    sort: '-operation_time',
                    stuff: {
                        // 如果初始化时有查询参数，则认为处在查询模式
                        type: Object.keys(this.$route.query).length ? 'search' : 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                details: {
                    isShow: false,
                    data: null,
                    showDetailsList: ['host', 'model_instance', 'business', 'cloud_area']
                }
            }
        },
        computed: {
            ...mapState('operationAudit', ['funcActions']),
            ...mapGetters('objectBiz', ['authorizedBusiness']),
            actionList () {
                return this.funcActions[this.active] || []
            },
            operationTypes () {
                const list = []
                const topoOperation = [{
                    id: 'module',
                    name: '模块'
                }, {
                    id: 'set',
                    name: '集群'
                }]
                this.actionList.forEach(item => {
                    if (item.id === 'biz_topology') {
                        list.push(...topoOperation)
                    } else {
                        list.push(item)
                    }
                })
                return list.map(item => ({
                    id: item.id,
                    name: this.$t(item.name)
                }))
            },
            funcModules () {
                const modules = {}
                this.actionList.forEach(item => {
                    modules[item.id] = this.$t(item.name)
                })
                return modules
            },
            actionSet () {
                const actionSet = {}
                const operations = this.actionList.reduce((acc, item) => acc.concat(item.operations), [])
                operations.forEach(action => {
                    actionSet[action.id] = this.$t(action.name)
                })
                return actionSet
            },
            params () {
                let time = []
                if (this.filter.time.length) {
                    time = [
                        this.filter.time[0] ? this.filter.time[0] + ' 00:00:00' : '',
                        this.filter.time[1] ? this.filter.time[1] + ' 23:59:59' : ''
                    ]
                }
                const params = {
                    condition: {
                        operation_time: time,
                        bk_biz_id: this.filter.bizId ? Number(this.filter.bizId) : null,
                        user: this.filter.operator,
                        resource_name: this.filter.resourceName,
                        resource_id: this.filter.resourceId ? Number(this.filter.resourceId) : null,
                        category: this.active
                    },
                    start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                    limit: this.table.pagination.limit,
                    sort: this.table.sort
                }
                const operations = this.filter.action
                const moduleName = this.filter.module
                if (operations.length) {
                    const types = []
                    const action = []
                    const label = []
                    operations.forEach(item => {
                        const resource = item.split('-')
                        types.push(resource[0])
                        action.push(resource[1])
                        if (resource[2] && !label.includes(resource[2])) {
                            label.push(resource[2])
                        }
                    })
                    if (moduleName && this.filter.resourceId && !types.includes(moduleName)) {
                        types.push(moduleName)
                    }
                    params.condition.resource_type = types
                    params.condition.action = action
                    label.length && (params.condition.label = label)
                } else if (moduleName && this.filter.resourceId) {
                    params.condition.resource_type = [moduleName]
                }
                return params
            }
        },
        async created () {
            this.$store.dispatch('objectBiz/getAuthorizedBusiness')
            await this.getTableData()
        },
        methods: {
            getResourceType (row) {
                if (row.label === null) {
                    const type = row.resource_type
                    if (type === 'model_instance') {
                        const objId = row.operation_detail.bk_obj_id
                        const model = this.$store.getters['objectModelClassify/getModelById'](objId) || {}
                        return model.bk_obj_name || '--'
                    }
                    return this.funcModules[type] || '--'
                }
                const key = Object.keys(row.label)
                return this.funcModules[key[0]] || '--'
            },
            getResourceName (row) {
                // 转移主机类的操作实例
                if (['assign_host', 'unassign_host', 'transfer_host_module'].includes(row.action)) {
                    return row.bk_host_innerip || row.operation_detail.bk_host_innerip || '--'
                }
                // 关联关系的操作实例
                if (['instance_association'].includes(row.resource_type)) {
                    return row.operation_detail.src_instance_name || '--'
                }
                // 自定义字段的操作实例
                if (['model_attribute'].includes(row.resource_type)) {
                    return `${row.operation_detail.bk_obj_name}/${row.operation_detail.resource_name}`
                }
                return this.$tools.getValue(row, 'operation_detail.resource_name') || '--'
            },
            getTargetName (row) {
                if (row.resource_type === 'instance_association') {
                    return row.operation_detail.target_instance_name
                }
                return this.getResourceName(row)
            },
            getResourceAction (row) {
                if (row.label) {
                    const label = Object.keys(row.label)[0]
                    return this.actionSet[`${row.resource_type}-${row.action}-${label}`]
                }
                return this.actionSet[`${row.resource_type}-${row.action}`]
            },
            getBusinessName (row) {
                return row.bk_biz_name
                    || this.$tools.getValue(row, 'operation_detail.bk_biz_name')
                    || '--'
            },
            async getTableData (event) {
                try {
                    const res = await this.$store.dispatch('operationAudit/getOperationLog', {
                        params: this.params,
                        config: {
                            globalPermission: false,
                            cancelPrevious: true,
                            requestId: 'getOperationLog'
                        }
                    })
                    this.table.list = res.info
                    this.table.pagination.count = res.count
                    // 有传入event参数时认为来自用户搜索
                    if (event) {
                        this.table.stuff.type = 'search'
                    }
                } catch ({ permission }) {
                    this.table.list = []
                    // 从api调用层抛出的错误，仅当权限问题时会注入permission
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                }
            },
            handleResetFilter () {
                this.filter.bizId = ''
                this.filter.action = []
                this.filter.resourceName = ''
                this.filter.operator = ''
                this.filter.resourceId = ''
                this.filter.module = ''
                this.filter.time = [
                    this.$tools.formatTime(moment().subtract(1, 'days'), 'YYYY-MM-DD'),
                    this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                ]
            },
            handleTabChange (tab) {
                this.table.list = []
                this.handleResetFilter()
                this.handlePageChange(1)
            },
            handlePageChange (current, event) {
                this.table.pagination.current = current
                this.getTableData(event)
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handleRowClick (item) {
                this.details.data = item
                this.details.isShow = true
            },
            handleClearFilter () {
                this.handleResetFilter()
                this.handlePageChange(1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .audit-layout{
        padding: 0 20px;
        .audit-tab {
            height: auto;
            /deep/ .bk-tab-header {
                padding: 0;
            }
        }
        .filter {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            flex-direction: row;
            flex-wrap: wrap;
            padding: 22px 0 10px 0;
            .option {
                flex: none;
                width: 27%;
                margin: 0 1.5% 12px 0;
                white-space: nowrap;
                font-size: 0;
                display: flex;
                align-items: center;
                &.instance {
                    .bk-select {
                        width: 40%;
                        margin-right: 5px;
                    }
                    .bk-form-control {
                        width: calc(60% - 5px);
                    }
                }

                &.action,
                &.operator {
                    .name {
                        width: 96px;
                        text-align: right;
                    }
                }

                &.resource,
                &.instance {
                    .name {
                        width: 70px;
                        text-align: right;
                    }
                }
            }
            .option-btn {
                width: auto;
                .bk-button + .bk-button {
                    margin-left: 8px;
                }
            }
            .name {
                font-size: 14px;
                padding-right: 10px;
                @include ellipsis;
            }
            .content {
                flex: 1;
                width: calc(100% - 96px);
                .bk-select {
                    width: 100%;
                }
            }
        }

    }

    [bk-language="en"] {
        .filter {
            .name {
                min-width: 70px;
                text-align: right;
            }

            .option {
                &.action,
                &.operator {
                    .name {
                        width: 146px;
                        text-align: right;
                    }
                }

                &.resource,
                &.instance {
                    .name {
                        width: 130px;
                        text-align: right;
                    }
                }
            }
        }
    }
</style>
