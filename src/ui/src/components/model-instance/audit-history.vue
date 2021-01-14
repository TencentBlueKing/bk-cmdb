<template>
    <div class="history">
        <div class="history-filter">
            <cmdb-form-date-range class="filter-item filter-range"
                v-model="condition.operation_time"
                :clearable="false"
                @input="handlePageChange(1)">
            </cmdb-form-date-range>
            <cmdb-form-objuser class="filter-item filter-user"
                v-model="condition.user"
                :exclude="false"
                :multiple="false"
                :palceholder="$t('操作账号')"
                @input="handlePageChange(1)">
            </cmdb-form-objuser>
        </div>
        <bk-table class="history-table"
            v-bkloading="{ isLoading: $loading(requestId) }"
            :data="history"
            :pagination="pagination"
            :max-height="$APP.height - 325"
            :row-style="{ cursor: 'pointer' }"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @row-click="handleRowClick">
            <bk-table-column :label="$t('操作描述')" :formatter="getFormatterDesc"></bk-table-column>
            <bk-table-column prop="user" :label="$t('操作账号')"></bk-table-column>
            <bk-table-column prop="operation_time" :label="$t('操作时间')">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row['operation_time'])}}
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import AuditDetails from '@/components/audit-history/details.js'
    export default {
        props: {
            category: {
                type: String,
                required: true
            },
            resourceId: {
                type: [Number, String]
            },
            resourceType: {
                type: String,
                default: ''
            }
        },
        data () {
            const today = this.$tools.formatTime(new Date(), 'YYYY-MM-DD')
            return {
                history: [],
                dictionary: [],
                condition: {
                    operation_time: [today, today],
                    user: '',
                    category: this.category,
                    resource_id: this.resourceId,
                    resource_type: this.resourceType
                },
                table: {
                    stuff: {
                        type: 'default',
                        payload: {}
                    }
                },
                pagination: {
                    count: 0,
                    current: 1,
                    limit: 10
                },
                requestId: Symbol('getList')
            }
        },
        created () {
            this.getAuditDictionary()
            this.getHistory()
        },
        methods: {
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
            async getHistory () {
                try {
                    const { info, count } = await this.$store.dispatch('audit/getList', {
                        params: {
                            condition: this.getUsefulConditon(),
                            page: {
                                ...this.$tools.getPageParams(this.pagination),
                                sort: '-operation_time'
                            }
                        },
                        config: {
                            requestId: this.requestId,
                            globalPermission: false
                        }
                    })
                    this.pagination.count = count
                    this.history = info
                } catch ({ permission }) {
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                    this.history = []
                }
            },
            getUsefulConditon () {
                const usefuleCondition = {}
                Object.keys(this.condition).forEach(key => {
                    const value = this.condition[key]
                    if (String(value).length) {
                        usefuleCondition[key] = value
                    }
                })
                if (usefuleCondition.operation_time) {
                    const [start, end] = usefuleCondition.operation_time
                    usefuleCondition.operation_time = {
                        start: start + ' 00:00:00',
                        end: end + ' 23:59:59'
                    }
                }
                return usefuleCondition
            },
            getFormatterDesc (row) {
                const type = this.dictionary.find(type => type.id === row.resource_type) || {}
                const action = (type.operations || []).find(action => action.id === row.action) || {}
                return `${action.name || row.action}${type.name || row.resource_type}`
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getHistory(true)
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.pagination.current = 1
                this.getHistory()
            },
            handleSortChange (sort) {
                this.sort = this.$tools.getSort(sort)
                this.getHistory()
            },
            handleRowClick (item) {
                AuditDetails.show({
                    id: item.id
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history {
        height: 100%;
    }
    .history-filter {
        padding: 14px 0;
        .filter-item {
            display: inline-block;
            vertical-align: middle;
            &.filter-range {
                width: 300px !important;
                margin: 0 5px 0 0;
            }
            &.filter-user {
                width: 240px;
                height: 32px;
            }
        }
    }
</style>
