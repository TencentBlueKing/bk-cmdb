<template>
    <div class="details-wrapper">
        <div class="table-box">
            <p class="title clearfix">
                <label @click="propertyTable.isShow = !propertyTable.isShow">
                    <i class="bk-icon icon-angle-down" :class="{'hide': !propertyTable.isShow}"></i>
                    <span>{{$t('Common["属性"]')}}</span>
                </label>
                <label class="cmdb-form-checkbox cmdb-checkbox-small">
                    <input type="checkbox" :disabled="ignore">
                    <span class="cmdb-checkbox-text">{{$t('NetworkDiscovery["显示忽略"]')}}</span>
                </label>
            </p>
            <cmdb-collapse-transition>
                <div v-show="propertyTable.isShow">
                    <cmdb-table
                        class="table"
                        :loading="$loading('searchNetcollectChangeDetail')"
                        :header="propertyTable.header"
                        :list="propertyTable.list"
                        :pagination.sync="propertyTable.pagination"
                        :defaultSort="propertyTable.defaultSort">
                        <template v-for="(header, index) in propertyTable.header" :slot="header.id" slot-scope="{ item }">
                            <template v-if="header.id === 'isrequired'">
                                <span :key="index" :class="{'disabled': item.method !== 'accept'}">
                                    {{item.isrequired ? $t('NetworkDiscovery["是"]') : $t('NetworkDiscovery["否"]')}}
                                </span>
                            </template>
                            <template v-else-if="header.id === 'operation'">
                                <span :key="index" class="text-primary" :class="{'disabled': ignore}" @click.stop="togglePropertyMethod(item)">{{item.method === 'accept' ? $t('NetworkDiscovery["忽略"]') : $t('NetworkDiscovery["取消忽略"]')}}</span>
                            </template>
                            <template v-else>
                                <span :key="index" :class="{'disabled': item.method !== 'accept'}">{{item[header.id]}}</span>
                            </template>
                        </template>
                    </cmdb-table>
                </div>
            </cmdb-collapse-transition>
        </div>
        <div class="table-box relation">
            <p class="title clearfix">
                <label class="title" @click="relationTable.isShow = !relationTable.isShow">
                    <i class="bk-icon icon-angle-down"></i>
                    <span>{{$t('NetworkDiscovery["关系"]')}}</span>
                </label>
                <label class="cmdb-form-checkbox cmdb-checkbox-small">
                    <input type="checkbox" :disabled="ignore">
                    <span class="cmdb-checkbox-text">{{$t('NetworkDiscovery["显示忽略"]')}}</span>
                </label>
            </p>
            <cmdb-collapse-transition>
                <div v-show="relationTable.isShow">
                    <cmdb-table
                        class="table"
                        :loading="$loading('searchNetcollectChangeDetail')"
                        :header="relationTable.header"
                        :list="relationTable.list"
                        :pagination.sync="relationTable.pagination"
                        :defaultSort="relationTable.defaultSort">
                        <template v-for="(header, index) in relationTable.header" :slot="header.id" slot-scope="{ item }">
                            <template v-if="header.id === 'action'">
                                <span :key="index" :class="{'color-danger': item.action === 'delete', 'disabled': item.asst.method !== 'accept'}">{{actionMap[item.action]}}</span>
                            </template>
                            <template v-else-if="header.id === 'operation'">
                                <span :key="index" class="text-primary" :class="{'disabled': ignore}" @click.stop="toggleRelationMethod(item)">{{item.asst.method === 'accept' ? $t('NetworkDiscovery["忽略"]') : $t('NetworkDiscovery["取消忽略"]')}}</span>
                            </template>
                            <template v-else>
                                <span :key="index" :class="{'disabled': item.asst.method !== 'accept'}">{{item[header.id]}}</span>
                            </template>
                        </template>
                    </cmdb-table>
                </div>
            </cmdb-collapse-transition>
        </div>
        <footer class="footer">
            <span>{{$t('NetworkDiscovery["导入实例"]')}}</span>
            <bk-switcher
                class="switcher"
                size="small"
                :show-text="false"
                :selected="!ignore"
                @change="toggleSwitcher">
            </bk-switcher>
            <bk-button type="default">
                {{$t('NetworkDiscovery["上一个"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('NetworkDiscovery["下一个"]')}}
            </bk-button>
        </footer>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            attributes: {
                type: Array
            },
            associations: {
                type: Array
            },
            ignore: {
                type: Boolean
            }
        },
        data () {
            return {
                isAccept: false,
                propertyTable: {
                    isShow: true,
                    header: [{
                        id: 'bk_property_name',
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
                        name: this.$t('Association["操作"]')
                    }],
                    list: [{
                        bk_property_id: 'bk_inst_name',
                        bk_property_name: '实例名',
                        bk_obj_id: 'bk_switch',
                        isrequired: true,
                        new_value: 'ddd',
                        pre_value: 'asd'
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
                    isShow: true,
                    header: [{
                        id: 'action',
                        name: this.$t('NetworkDiscovery["操作方式"]')
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('OperationAudit["模型"]')
                    }, {
                        id: 'device_attributes',
                        name: this.$t('NetworkDiscovery["配置信息"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]')
                    }],
                    list: [{
                        action: 'update',
                        asst: {
                            bk_inst_id: 1,
                            bk_obj_id: 'bk_switch',
                            bk_obj_name: '交换机',
                            bk_asst_inst_id: 0,
                            bk_asst_obj_id: 'host',
                            bk_asst_obj_name: '主机'
                        }
                    }],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                },
                actionMap: {
                    'create': this.$t("Association['新增关联']"),
                    'delete': this.$t("Common['删除关联']")
                }
            }
        },
        created () {
            this.propertyTable.list = this.$tools.clone(this.attributes)
            this.relationTable.list = this.$tools.clone(this.associations)
        },
        methods: {
            toggleSwitcher (value) {
                this.$emit('toggleSwitcher', value)
            },
            togglePropertyMethod (item) {
                if (this.ignore) {
                    return
                }
                item.method = item.method === 'accept' ? 'reject' : 'accept'
                this.$emit('update:attributes', this.propertyTable.list)
            },
            toggleRelationMethod (item) {
                if (this.ignore) {
                    return
                }
                item.asst.method = item.asst.method === 'accept' ? 'reject' : 'accept'
                this.$emit('update:associations', this.propertyTable.list)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-wrapper {
        padding: 15px 30px;
        .disabled {
            color: $cmdbBorderColor;
        }
        .table-box {
            &.relation {
                margin-top: 20px;
            }
            .title {
                line-height: 32px;
                >.cmdb-form-checkbox {
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
