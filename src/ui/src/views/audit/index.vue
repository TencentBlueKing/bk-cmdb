<template>
    <div class="audit-layout">
        <bk-tab v-show="!$loading(request.list)"
            class="audit-tab"
            type="unborder-card"
            :active.sync="active">
            <bk-tab-panel
                v-for="panel in tabPanels"
                v-bind="panel"
                :key="panel.id">
            </bk-tab-panel>
        </bk-tab>
        <div class="audit-options" v-show="!$loading(request.list)">
            <component :is="optionsComponent" @condition-change="handleConditionChange"></component>
        </div>
        <bk-table
            v-show="!$loading(request.list)"
            v-bkloading="{ isLoading: $loading(request.list) }"
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
                :label="$t('操作对象')"
                :formatter="getResourceTypeName">
            </bk-table-column>
            <bk-table-column
                prop="action"
                :label="$t('动作')"
                :formatter="getActionName">
            </bk-table-column>
            <bk-table-column
                v-if="['host', 'business'].includes(active)"
                prop="bk_biz_name"
                :label="$t('所属业务')">
                <template slot-scope="{ row }">
                    <audit-business-selector type="info" :value="row.bk_biz_id"></audit-business-selector>
                </template>
            </bk-table-column>
            <bk-table-column
                prop="resource_name"
                :label="$t('操作实例')">
            </bk-table-column>
            <bk-table-column
                :label="$t('操作描述')">
                <template slot-scope="{ row }">{{`${getActionName(row)}${getResourceTypeName(row)}`}}</template>
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
    </div>
</template>

<script>
    import AuditBusinessOptions from './children/audit-business-options'
    import AuditResourceOptions from './children/audit-resource-options'
    import AuditOtherOptions from './children/audit-other-options'
    import AuditHostOptions from './children/audit-host-options'
    import AuditBusinessSelector from '@/components/audit-history/audit-business-selector'
    import RouterQuery from '@/router/query'
    import AuditDetails from '@/components/audit-history/details.js'
    export default {
        components: {
            [AuditBusinessOptions.name]: AuditBusinessOptions,
            [AuditResourceOptions.name]: AuditResourceOptions,
            [AuditOtherOptions.name]: AuditOtherOptions,
            [AuditHostOptions.name]: AuditHostOptions,
            AuditBusinessSelector
        },
        data () {
            return {
                active: RouterQuery.get('tab', 'host'),
                tabPanels: [
                    {
                        name: 'host',
                        label: this.$t('主机')
                    },
                    {
                        name: 'business',
                        label: this.$t('业务')
                    }, {
                        name: 'resource',
                        label: this.$t('资源')
                    }, {
                        name: 'other',
                        label: this.$t('其他')
                    }
                ],
                dictionary: [],
                condition: {},
                table: {
                    list: [],
                    pagination: this.$tools.getDefaultPaginationConfig({ 'limit-list': [20, 50, 100, 200] }),
                    sort: '-operation_time',
                    stuff: {
                        // 如果初始化时有查询参数，则认为处在查询模式
                        type: Object.keys(this.$route.query).length ? 'search' : 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                request: {
                    list: Symbol('list')
                }
            }
        },
        computed: {
            optionsComponent () {
                const componentMap = {
                    business: AuditBusinessOptions.name,
                    resource: AuditResourceOptions.name,
                    other: AuditOtherOptions.name,
                    host: AuditHostOptions.name
                }
                return componentMap[this.active]
            }
        },
        created () {
            this.setQueryWatcher()
            this.getAuditDictionary()
        },
        beforeDestroy () {
            this.teardownQueryWatcher()
        },
        methods: {
            setQueryWatcher () {
                this.unwatch = RouterQuery.watch('*', ({
                    page,
                    limit,
                    sort,
                    tab,
                    _e: isEvent
                }) => {
                    this.active = tab || 'host'
                    this.table.pagination.current = parseInt(page || this.table.pagination.current, 10)
                    this.table.pagination.limit = parseInt(limit || this.table.pagination.limit, 10)
                    this.table.sort = sort || this.table.sort
                    this.$nextTick(() => this.getAuditList(isEvent))
                })
            },
            teardownQueryWatcher () {
                this.unwatch && this.unwatch()
            },
            async getAuditDictionary () {
                try {
                    this.dictionary = await this.$store.dispatch('audit/getDictionary', {
                        fromCache: true,
                        globalPermission: false
                    })
                } catch (error) {
                    this.dictionary = []
                }
            },
            handleConditionChange (condition) {
                const usefulCondition = {}
                Object.keys(condition).forEach(key => {
                    const value = condition[key]
                    if (String(value).length) {
                        usefulCondition[key] = value
                    }
                })
                // 动态分组的ID是String, 其他的是Number, 区别转换
                if (usefulCondition.resource_id) {
                    usefulCondition.resource_id = usefulCondition.resource_type === 'dynamic_group'
                        ? usefulCondition.resource_id
                        : parseInt(usefulCondition.resource_id, 10)
                }
                // 转换时间范围为start/end的形式
                if (usefulCondition.operation_time) {
                    const [start, end] = usefulCondition.operation_time
                    usefulCondition.operation_time = {
                        start: start + ' 00:00:00',
                        end: end + ' 23:59:59'
                    }
                }
                this.condition = usefulCondition
            },
            async getAuditList (eventTrigger) {
                try {
                    const params = {
                        condition: this.condition,
                        page: {
                            ...this.$tools.getPageParams(this.table.pagination),
                            sort: this.table.sort
                        }
                    }
                    const { info, count } = await this.$store.dispatch('audit/getList', {
                        params,
                        config: {
                            requestId: this.request.list,
                            globalPermission: false
                        }
                    })

                    this.table.stuff.type = eventTrigger ? 'search' : 'default'
                    this.table.pagination.count = count
                    this.table.list = info
                } catch ({ permission }) {
                    this.$route.meta.view = 'permission'
                }
            },
            handlePageChange (current) {
                RouterQuery.set({
                    page: current,
                    _t: Date.now()
                })
            },
            handleSizeChange (size) {
                RouterQuery.set({
                    limit: size,
                    page: 1,
                    _t: Date.now()
                })
            },
            handleSortChange (sort) {
                RouterQuery.set({
                    page: 1,
                    sort: this.$tools.getSort(sort, 'operation_time'),
                    _t: Date.now()
                })
            },
            handleRowClick (row) {
                AuditDetails.show({
                    id: row.id
                })
            },
            getResourceTypeName (row) {
                const type = this.dictionary.find(type => type.id === row.resource_type)
                return type ? type.name : row.resource_type
            },
            getActionName (row) {
                const type = this.dictionary.find(type => type.id === row.resource_type)
                const operations = type ? type.operations : []
                const operation = operations.find(operation => operation.id === row.action)
                return operation ? operation.name : row.action
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
