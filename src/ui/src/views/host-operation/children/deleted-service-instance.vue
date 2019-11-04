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
            :data="info"
            :cell-style="getCellStyle">
            <bk-table-column :label="$t('操作')">
                <!-- eslint-disable-next-line -->
                <template slot-scope="{ row }">{{$t('删除')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('服务实例')" prop="name"></bk-table-column>
            <bk-table-column :label="$t('所属模块')" prop="bk_module_id">
                <template slot-scope="{ row }">{{getModuleName(row)}}</template>
            </bk-table-column>
        </bk-table>
    </section>
</template>

<script>
    export default {
        props: {
            info: {
                type: Array,
                required: true
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
            },
            getModuleName (row) {
                const module = this.$parent.moduleInfo.find(module => module.bk_module_id === row.bk_module_id)
                if (module) {
                    return module.bk_module_name
                }
                return '--'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table {
        margin-top: 8px;
    }
</style>
