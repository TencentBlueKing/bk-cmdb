<template>
    <section class="layout">
        <cmdb-tips
            :tips-style="{
                background: 'none',
                border: 'none',
                fontSize: '12px'
            }"
            :icon-style="{
                color: '#63656E',
                fontSize: '14px'
            }">
            {{$t('移动到空闲机的主机提示')}}
        </cmdb-tips>
        <bk-table class="table"
            :data="data"
            :cell-style="getCellStyle">
            <bk-table-column :label="$t('操作')" prop="operation"></bk-table-column>
            <bk-table-column :label="$t('IP')" prop="bk_host_innerip"></bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="bk_cloud_id"></bk-table-column>
        </bk-table>
    </section>
</template>

<script>
    export default {
        data () {
            return {
                data: [{
                    operation: this.$t('转移到空闲机')
                }]
            }
        },
        methods: {
            getCellStyle ({ columnIndex }) {
                if (columnIndex === 0) {
                    return {
                        color: '#3A84FF'
                    }
                }
                return {}
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table {
        margin-top: 8px;
    }
</style>
