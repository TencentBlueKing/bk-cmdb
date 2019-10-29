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
            {{$t('删除服务实例提示')}}
        </cmdb-tips>
        <bk-table class="table"
            :data="data"
            :cell-style="getCellStyle">
            <bk-table-column :label="$t('操作')" prop="operation"></bk-table-column>
            <bk-table-column :label="$t('服务实例')" prop="name"></bk-table-column>
            <bk-table-column :label="$t('所属模块')" prop="module"></bk-table-column>
        </bk-table>
    </section>
</template>

<script>
    export default {
        data () {
            return {
                data: [{
                    operation: this.$t('删除')
                }, {
                    operation: this.$t('删除')
                }]
            }
        },
        methods: {
            getCellStyle ({ columnIndex }) {
                if (columnIndex === 0) {
                    return {
                        color: '#FF5656'
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
