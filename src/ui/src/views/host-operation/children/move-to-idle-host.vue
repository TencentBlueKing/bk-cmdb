<template>
    <section class="move-layout">
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
        <bk-table class="table" :data="list">
            <bk-table-column :label="$t('操作')">
                <!-- eslint-disable-next-line -->
                <template slot-scope="{ row }">{{$t('转移到空闲机的主机')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('IP')" prop="bk_host_innerip">
                <template slot-scope="{ row }">{{getHostValue(row, 'bk_host_innerip') | singlechar}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="bk_cloud_id">
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
        data () {
            return {
                data: [{
                    operation: this.$t('转移到空闲机')
                }]
            }
        },
        computed: {
            list () {
                return this.info.map(id => {
                    const target = this.$parent.hostInfo.find(target => target.host.bk_host_id === id)
                    return target || {}
                })
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
