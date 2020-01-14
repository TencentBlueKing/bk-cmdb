<template>
    <div class="cloud-area-layout">
        <cmdb-tips class="cloud-area-tips">提示语</cmdb-tips>
        <bk-table class="cloud-area-table"
            :data="list"
            :pagination="pagination">
            <bk-table-column :label="$t('云区域名称')" prop="xxxx" class-name="is-highlight"></bk-table-column>
            <bk-table-column :label="$t('状态')" prop="xxxx">
                <div slot-scope="{ row }" class="row-status" v-bk-tooltips.right="'异常原因'">
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
                <text-button slot-scope="{ row }"
                    :disabled="row.host_count"
                    v-bk-tooltips="'主机不为空，不能删除'"
                    @click="handleFuck(row)">
                    {{$t('删除')}}
                </text-button>
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
        methods: {
            handleFuck (row) {
                console.log(row)
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
