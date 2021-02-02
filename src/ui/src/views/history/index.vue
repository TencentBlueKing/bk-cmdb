<template>
    <div class="history-layout">
        <div class="history-options">
            <cmdb-form-date-range class="history-date-range"
                v-model="condition.operation_time"
                :clearable="false"
                @change="handlePageChange(1)">
            </cmdb-form-date-range>
            <bk-input class="history-host-filter ml10"
                v-if="isHost"
                right-icon="icon-search"
                clearable
                v-model="condition.resource_name"
                :placeholder="$t('请输入xx', { name: 'IP' })"
                @enter="handlePageChange(1)"
                @clear="handlePageChange(1)">
            </bk-input>
        </div>
        <bk-table class="history-table"
            v-bkloading="{ isLoading: $loading() }"
            :pagination="pagination"
            :data="history"
            :max-height="$APP.height - 190"
            :row-style="{ cursor: 'pointer' }"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @row-click="handleRowClick">
            <bk-table-column prop="resource_id" label="ID"></bk-table-column>
            <bk-table-column prop="resource_name" :label="isHost ? 'IP' : $t('资源')"></bk-table-column>
            <bk-table-column prop="operation_time" :label="$t('更新时间')">
                <template slot-scope="{ row }">{{$tools.formatTime(row.operation_time)}}</template>
            </bk-table-column>
            <bk-table-column prop="user" :label="$t('操作账号')">
                <template slot-scope="{ row }">
                    <cmdb-form-objuser :value="row.user" type="info"></cmdb-form-objuser>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import AuditDetails from '@/components/audit-history/details.js'
    export default {
        data () {
            const today = this.$tools.formatTime(new Date(), 'YYYY-MM-DD')
            return {
                dictionary: [],
                history: [],
                pagination: this.$tools.getDefaultPaginationConfig(),
                condition: {
                    operation_time: [today, today],
                    resource_name: '',
                    action: ['delete']
                },
                table: {
                    stuff: {
                        type: 'default',
                        payload: {}
                    }
                },
                requestId: Symbol('getHistory')
            }
        },
        computed: {
            objId () {
                if (this.$route.name === 'hostHistory') {
                    return 'host'
                }
                return this.$route.params.objId
            },
            isHost () {
                return this.objId === 'host'
            }
        },
        watch: {
            objId: {
                immediate: true,
                handler (objId) {
                    const model = this.$store.getters['objectModelClassify/getModelById'](objId) || {}
                    this.$store.commit('setTitle', `${model.bk_obj_name}${this.$t('删除历史')}`)
                }
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
                const usefuleCondition = {
                    category: this.isHost ? 'host' : 'resource',
                    resource_type: this.isHost ? 'host' : 'model_instance'
                }
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
            handleSizeChange (size) {
                this.pagination.limit = size
                this.handlePageChange(1)
            },
            handlePageChange (current) {
                this.pagination.current = current
                this.getHistory()
            },
            handleRowClick (row) {
                AuditDetails.show({
                    id: row.id
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-layout{
        padding: 15px 20px 0;
    }
    .history-options{
        font-size: 0px;
        .history-host-filter,
        .history-date-range {
            width: 260px !important;
            display: inline-block;
            vertical-align: top;
        }
    }
    .history-table{
        margin-top: 15px;
    }

</style>
