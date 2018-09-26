<template>
    <div class="details-wrapper">
        <div class="table-box">
            <p class="title clearfix">
                <i class="bk-icon icon-angle-down"></i>
                <span>{{$t('Common["属性"]')}}</span>
                <label class="cmdb-form-checkbox cmdb-checkbox-small">
                    <input type="checkbox">
                    <span class="cmdb-checkbox-text">{{$t('networkDiscovery["显示忽略"]')}}</span>
                </label>
            </p>
            <cmdb-table
                class="table"
                :loading="$loading('searchSubscription')"
                :header="propertyTable.header"
                :list="propertyTable.list"
                :pagination.sync="propertyTable.pagination"
                :defaultSort="propertyTable.defaultSort">
            </cmdb-table>
        </div>
        <div class="table-box relation">
            <p class="title clearfix">
                <i class="bk-icon icon-angle-down"></i>
                <span>{{$t('networkDiscovery["关系"]')}}</span>
                <label class="cmdb-form-checkbox cmdb-checkbox-small">
                    <input type="checkbox">
                    <span class="cmdb-checkbox-text">{{$t('networkDiscovery["显示忽略"]')}}</span>
                </label>
            </p>
            <cmdb-table
                class="table"
                :loading="$loading('searchSubscription')"
                :header="relationTable.header"
                :list="relationTable.list"
                :pagination.sync="relationTable.pagination"
                :defaultSort="relationTable.defaultSort">
            </cmdb-table>
        </div>
        <footer class="footer">
            <span>{{$t('networkDiscovery["导入实例"]')}}</span>
            <bk-switcher
                class="switcher"
                size="small"
                :show-text="false">
            </bk-switcher>
            <bk-button type="default">
                {{$t('networkDiscovery["上一个"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('networkDiscovery["下一个"]')}}
            </bk-button>
        </footer>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                propertyTable: {
                    header: [{
                        id: 'bk_property_id',
                        name: this.$t('NetworkDiscovery["属性名"]')
                    }, {
                        id: 'isrequired',
                        name: this.$t('NetworkDiscovery["必须"]')
                    }, {
                        id: 'pre_value',
                        name: this.$t('NetworkDiscovery["原值"]')
                    }, {
                        id: 'new_value',
                        name: this.$t('NetworkDiscovery["新值"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [{
                        plat: 'asdf',
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
                },
                relationTable: {
                    header: [{
                        id: 'action',
                        name: this.$t('NetworkDiscovery["操作方式"]')
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('NetworkDiscovery["模型"]')
                    }, {
                        id: 'config',
                        name: this.$t('NetworkDiscovery["配置信息"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [{
                        action: 'asdf',
                        config: 'aaa',
                        last_time: 'ddd'
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
    .details-wrapper {
        padding: 15px 30px;
        .table-box {
            &.relation {
                margin-top: 20px;
            }
            .title {
                line-height: 32px;
                >label {
                    float: right;
                }
                .icon-angle-down {
                    font-size: 12px;
                    font-weight: bold;
                }
            }
            .table {
                margin-top: 2px;
            }
        }
        .footer {
            position: absolute;
            bottom: 0;
            left: 0;
            padding: 12px 20px;
            width: 100%;
            font-size: 0;
            text-align: right;
            box-shadow: 0px -2px 5px 0px rgba(0, 0, 0, 0.05);
            >span {
                font-size: 14px;
                vertical-align: middle;
                margin-right: 8px;
            }
            .switcher {
                margin-right: 10px;
            }
            .bk-button {
                margin-left: 10px;
            }
        }
    }
</style>

<style lang="scss">
    .details-wrapper {
        .switcher {
            &.is-checked {
                background: $cmdbBorderFocusColor;
            }
        }
    }
</style>
