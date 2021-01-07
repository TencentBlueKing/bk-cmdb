<template>
    <section class="move-layout">
        <cmdb-tips
            :tips-style="{
                background: 'none',
                border: 'none',
                fontSize: '12px',
                lineHeight: '30px',
                padding: 0
            }"
            :icon-style="{
                color: '#63656E',
                fontSize: '14px',
                lineHeight: '30px'
            }">
            {{$t('移动到空闲机的主机提示')}}
        </cmdb-tips>
        <bk-table class="table" :data="info">
            <bk-table-column :label="$t('操作')" show-overflow-tooltip>
                <!-- eslint-disable-next-line -->
                <template slot-scope="{ row }">{{$t('转移到空闲机')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('IP')" prop="bk_host_innerip" show-overflow-tooltip>
                <template slot-scope="{ row }">{{getHostValue(row, 'bk_host_innerip') | singlechar}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="bk_cloud_id" show-overflow-tooltip>
                <template slot-scope="{ row }">{{getHostValue(row, 'bk_cloud_id') | foreignkey}}</template>
            </bk-table-column>
        </bk-table>
    </section>
</template>

<script>
    import { foreignkey, singlechar } from '@/filters/formatter.js'
    export default {
        name: 'move-to-idle-host',
        filters: { foreignkey, singlechar },
        props: {
            info: {
                type: Array,
                required: true
            }
        },
        methods: {
            getHostValue (row, field) {
                const host = row.host
                if (host) {
                    return host[field]
                }
                return ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table {
        margin-top: 8px;
    }
</style>
