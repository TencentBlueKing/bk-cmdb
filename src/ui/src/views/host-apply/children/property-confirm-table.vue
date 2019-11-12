<template>
    <div class="single-module-panel">
        <bk-table
            :data="confirmData"
        >
            <bk-table-column :label="$t('内网IP')" prop="bk_host_innerip"></bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="bk_cloud_id"></bk-table-column>
            <bk-table-column :label="$t('固资编号')" prop="bk_asset_id"></bk-table-column>
            <bk-table-column :label="$t('主机名称')" prop="bk_host_name"></bk-table-column>
            <bk-table-column :label="$t('修改值')">
                <template slot-scope="{ row }">
                    <span>{{row.diff_value}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <bk-button theme="primary" text @click="handlePropertyRowDel(row.id)">处理冲突</bk-button>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    export default {
        components: {
        },
        props: {
            confirmData: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
            }
        },
        computed: {
        },
        watch: {
        },
        created () {
        },
        methods: {

        }
    }
</script>

<style lang="scss" scoped>

</style>
