<template>
    <div class="cloud-account-layout">
        <cmdb-tips class="cloud-account-tips">提示语</cmdb-tips>
        <div class="cloud-account-options">
            <bk-button theme="primary" @click="handleCreate">{{$t('新建')}}</bk-button>
        </div>
        <bk-table class="cloud-account-table"
            :data="list"
            :pagination="pagination"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @cell-click="handleCellClick">
            <bk-table-column :label="$t('账户名称')" prop="name" class-name="is-highlight" id="name"></bk-table-column>
            <bk-table-column :label="$t('账号类型')" prop="type"></bk-table-column>
            <bk-table-column :label="$t('修改人')" prop="updator"></bk-table-column>
            <bk-table-column :label="$t('修改时间')" prop="last_time"></bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <link-button class="mr10" @click="handleView(row)">{{$t('查看')}}</link-button>
                    <link-button
                        :disabled="row.pending"
                        v-bk-tooltips="{
                            disabled: !row.pending,
                            content: '发现任务正在进行中，不能删除'
                        }"
                        @click="handleDelete(row)">
                        {{$t('删除')}}
                    </link-button>
                </template>
            </bk-table-column>
        </bk-table>
        <account-sideslider ref="accountSideslider" @request-refresh="getData"></account-sideslider>
    </div>
</template>

<script>
    import AccountSideslider from './children/account-sideslider.vue'
    export default {
        components: {
            AccountSideslider
        },
        data () {
            return {
                list: [],
                pagination: {
                    ...this.$tools.getDefaultPaginationConfig(),
                    count: 2
                }
            }
        },
        created () {
            this.getData()
        },
        methods: {
            handleCreate () {
                this.$refs.accountSideslider.show({
                    type: 'form',
                    title: this.$t('新建账户'),
                    props: {
                        mode: 'create'
                    }
                })
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getData()
            },
            handleLimitChange (limit) {
                this.pagination.limit = limit
                this.pagination.current = 1
                this.getData()
            },
            handleCellClick (row, column) {
                if (column.property === 'name') {
                    this.handleView(row)
                }
            },
            handleView (row) {
                this.$refs.accountSideslider.show({
                    type: 'details',
                    title: `${this.$t('账户详情')} 【${row.name}】`,
                    props: {
                        id: row.id
                    }
                })
            },
            async handleDelete (row) {
                try {
                    await Promise.resolve()
                    this.$success('删除成功')
                    this.getData()
                } catch (e) {
                    console.error(e)
                }
            },
            async getData () {
                try {
                    const data = await Promise.resolve({
                        count: 2,
                        info: [{ id: 1 }, { id: 2 }]
                    })
                    if (data.count && !data.info.length) {
                        this.handlePageChange(this.pagination.current - 1)
                        return
                    }
                    this.list = data.info
                    this.pagination.count = data.count
                } catch (e) {
                    console.error(e)
                    this.list = []
                    this.pagination.count = 0
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cloud-account-layout {
        padding: 0 20px;
    }
    .cloud-account-tips {
        margin-top: 10px;
    }
    .cloud-account-options {
        margin-top: 10px;
    }
    .cloud-account-table {
        margin-top: 10px;
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
</style>
