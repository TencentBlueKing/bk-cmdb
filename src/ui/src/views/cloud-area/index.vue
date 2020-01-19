<template>
    <div class="cloud-area-layout">
        <cmdb-tips class="cloud-area-tips">提示语</cmdb-tips>
        <bk-table class="cloud-area-table"
            :data="list"
            :pagination="pagination"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange">
            <bk-table-column :label="$t('云区域名称')" prop="xxxx" class-name="is-highlight"></bk-table-column>
            <bk-table-column :label="$t('状态')" prop="xxxx">
                <div class="row-status"
                    slot-scope="{ row }"
                    v-bk-tooltips.right="{
                        disabled: row.error,
                        content: '异常原因'
                    }">
                    <i :class="['status', { 'is-error': row.error }]"></i>
                    异常
                </div>
            </bk-table-column>
            <bk-table-column :label="$t('所属云厂商')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('国家')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('省份/州')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('VPC')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('主机数量')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('最近编辑')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('编辑人')" prop="xxxx"></bk-table-column>
            <bk-table-column :label="$t('操作')" prop="xxxx">
                <link-button slot-scope="{ row }"
                    :disabled="row.host_count"
                    v-bk-tooltips="{
                        disabled: !row.host_count,
                        content: '主机不为空，不能删除'
                    }"
                    @click="handleDelete(row)">
                    {{$t('删除')}}
                </link-button>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                list: [{}, {}],
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
            handlePageChange (page) {
                this.pagination.current = page
                this.getData()
            },
            handleLimitChange (limit) {
                this.pagination.limit = limit
                this.pagination.current = 1
                this.getData()
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
                        info: [{}, {}]
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
    .cloud-area-layout {
        padding: 0 20px;
    }
    .cloud-area-tips {
        margin-top: 10px;
    }
    .cloud-area-table {
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
