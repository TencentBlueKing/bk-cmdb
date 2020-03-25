<template>
    <div class="cloud-area-layout">
        <cmdb-tips class="cloud-area-tips">提示语</cmdb-tips>
        <bk-table class="cloud-area-table"
            :data="list"
            :pagination="pagination"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange">
            <bk-table-column :label="$t('云区域名称')" prop="bk_cloud_name" class-name="is-highlight" show-overflow-tooltip></bk-table-column>
            <bk-table-column :label="$t('状态')" prop="bk_status">
                <div class="row-status" slot-scope="{ row }">
                    <i :class="['status', { 'is-error': row.bk_status !== '1' }]"></i>
                    {{row.bk_status === '1' ? $t('正常') : $t('异常')}}
                </div>
            </bk-table-column>
            <bk-table-column :label="$t('所属云厂商')" prop="bk_cloud_vendor">
                <template slot-scope="{ row }">{{row.bk_cloud_vendor | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('区域')" prop="bk_region">
                <template slot-scope="{ row }">{{row.bk_region | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('VPC')" prop="bk_vpc_name" show-overflow-tooltip>
                <template slot-scope="{ row }">{{getVpcInfo(row) | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('主机数量')" prop="host_count">
                <template slot-scope="{ row }">{{row.host_count | formatter('int')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('最近编辑')" prop="last_time">
                <template slot-scope="{ row }">{{row.last_time | formatter('time')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('编辑人')" prop="bk_last_editor"></bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <bk-popconfirm slot-scope="{ row }"
                    trigger="click"
                    :disabled="!!row.host_count"
                    :title="$t('确定删除该云区域')"
                    @confirm="handleDelete(row)">
                    <link-button
                        :disabled="!!row.host_count"
                        v-bk-tooltips="{
                            disabled: !row.host_count,
                            content: $t('主机不为空，不能删除')
                        }">
                        {{$t('删除')}}
                    </link-button>
                </bk-popconfirm>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                list: [{}, {}],
                pagination: this.$tools.getDefaultPaginationConfig()
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
                    await this.$store.dispatch('cloud/area/delete', { id: row.bk_cloud_id })
                    this.$success('删除成功')
                    this.getData()
                } catch (e) {
                    console.error(e)
                }
            },
            async getData () {
                try {
                    const data = await this.$store.dispatch('cloud/area/findMany', {
                        params: {
                            ...this.$tools.getPageParams(this.pagination),
                            host_count: true
                        }
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
            },
            getVpcInfo (row) {
                const id = row.bk_vpc_id
                const name = row.bk_vpc_name
                if (name && id !== name) {
                    return `${id}(${name})`
                }
                return id
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
