<template>
    <div class="history-layout">
        <div class="history-options">
            <bk-date-picker class="options-date"
                type="datetimerange"
                :placeholder="$t('请选择xx', { name: '时间范围' })"
                v-model="timeRange">
            </bk-date-picker>
        </div>
        <div class="history-table">
            <bk-table ref="table" :data="histories" :max-height="$APP.height - 180">
                <bk-table-column type="expand" width="30" align="center">
                    <task-details-history-content
                        slot-scope="{ row }"
                        :details="row.bk_detail">
                    </task-details-history-content>
                </bk-table-column>
                <bk-table-column :label="$t('操作概要')" prop="bk_summary" show-overflow-tooltip></bk-table-column>
                <bk-table-column :label="$t('状态')" prop="bk_sync_status">
                    <div class="row-status" slot-scope="{ row }">
                        <i :class="['status', { 'is-error': row.bk_sync_status !== 1 }]"></i>
                        {{row.bk_sync_status === 1 ? $t('成功') : $t('失败')}}
                    </div>
                </bk-table-column>
                <bk-table-column :label="$t('时间')" prop="create_time">
                    <template slot-scope="{ row }">{{row.create_time | formatter('time')}}</template>
                </bk-table-column>
                <bk-table-column :label="$t('详情')">
                    <template slot-scope="{ row }">
                        <link-button @click="handleView(row)">{{$t('查看详情')}}</link-button>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
    </div>
</template>

<script>
    import TaskDetailsHistoryContent from './task-details-history-content.vue'
    export default {
        name: 'task-details-history',
        components: {
            [TaskDetailsHistoryContent.name]: TaskDetailsHistoryContent
        },
        props: {
            id: Number
        },
        data () {
            return {
                timeRange: [new Date(Date.now() - 8.64e7), new Date()],
                histories: [],
                pagination: this.$tools.getDefaultPaginationConfig()
            }
        },
        created () {
            this.getHistory()
        },
        methods: {
            async getHistory () {
                try {
                    const data = await this.$store.dispatch('cloud/resource/findHistory', {
                        params: {
                            bk_task_id: this.id,
                            start_time: this.$tools.formatTime(this.timeRange[0]),
                            end_time: this.$tools.formatTime(this.timeRange[0]),
                            page: this.$tools.getPageParams(this.pagination)
                        }
                    })
                    this.pagination.count = data.count
                    this.histories = data.info
                } catch (e) {
                    console.error(e)
                    this.histories = []
                }
            },
            handleView (row) {
                this.$refs.table.toggleRowExpansion(row)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-layout {
        .history-options {
            padding: 12px 30px;
            .options-date {
                width: 320px;
            }
        }
        .history-table {
            padding: 0 30px;
            /deep/ {
                .bk-table-expanded-cell {
                    padding: 0;
                }
            }
        }
        .row-status {
            display: inline-block;
            .status {
                display: inline-block;
                margin-right: 4px;
                width: 7px;
                height: 7px;
                border-radius: 50%;
                background-color: $successColor;
            }
        }
    }
</style>
