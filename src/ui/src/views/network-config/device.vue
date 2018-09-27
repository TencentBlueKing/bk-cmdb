<template>
    <div class="device-wrapper">
        <div class="title">
            <bk-button type="primary">
                {{$t('NetworkConfig["新增设备"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('Common["删除"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('ModelManagement["导入"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('ModelManagement["导出"]')}}
            </bk-button>
        </div>
        <cmdb-table
            class="device-table"
            :loading="$loading('searchSubscription')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort">
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                table: {
                    header: [{
                        type: 'checkbox'
                    }, {
                        id: 'id',
                        name: 'ID'
                    }, {
                        id: 'type',
                        name: this.$t('NetworkConfig["设备型号"]')
                    }, {
                        id: 'name',
                        name: this.$t('NetworkConfig["设备名称"]')
                    }, {
                        id: 'model',
                        name: this.$t('NetworkConfig["对应模型"]')
                    }],
                    list: [{
                        change: 'asdf',
                        unique: 'asdf',
                        type: 'asdf',
                        ip: 'asdf',
                        config: 'aaa',
                        time: 'ddd'
                    }],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                }
            }
        }
    }
</script>


<style lang="scss" scoped>
    .device-wrapper {
        padding-top: 20px;
        .title {
            font-size: 0;
            .bk-default {
                margin-left: 10px;
            }
        }
        .device-table {
            margin-top: 20px;
            background: #fff;
        }
    }
</style>
