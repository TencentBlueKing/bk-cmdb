<template>
    <div class="network-confirm-wrapper">
        <div class="filter-wrapper" :class="{'open': filter.isShow}">
            <bk-button type="default" @click="toggleFilter">
                {{$t('networkDiscovery["高级操作"]')}}
                <i class="bk-icon icon-angle-down"></i>
            </bk-button>
            <div class="filter-details clearfix" v-show="filter.isShow">
                <div class="details-left">
                    <bk-button type="default">
                        {{$t('networkDiscovery["忽略"]')}}
                    </bk-button>
                    <bk-button type="default">
                        {{$t('networkDiscovery["取消忽略"]')}}
                    </bk-button>
                    <label class="cmdb-form-checkbox">
                        <input type="checkbox">
                        <span class="cmdb-checkbox-text">{{$t('networkDiscovery["显示忽略"]')}}</span>
                    </label>
                </div>
                <div class="details-right clearfix">
                    <bk-selector
                        :list="changeInfo.list"
                        :selected="changeInfo.selected"
                        :placeholder="$t('networkDiscovery[\'全部变更\']')"
                    ></bk-selector>
                    <bk-selector
                        :list="typeInfo.list"
                        :selected="typeInfo.selected"
                        :placeholder="$t('networkDiscovery[\'全部类型\']')"
                    ></bk-selector>
                    <input type="text" class="cmdb-form-input" :placeholder="$t('networkDiscovery[\'请输入IP\']')">
                    <bk-button type="default">
                        {{$t('Common["查询"]')}}
                    </bk-button>
                </div>
            </div>
        </div>
        <cmdb-table
            class="confirm-table"
            :loading="$loading('searchSubscription')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <template v-if="header.id === 'operation'">
                    <div :key="index">
                        <span class="text-primary" @click.stop="showDetails">{{$t('NetworkDiscovery["详情"]')}}</span>
                        <span class="text-primary" @click.stop="">{{$t('NetworkDiscovery["忽略"]')}}</span>
                        <span class="text-primary" @click.stop="">{{$t('NetworkDiscovery["取消忽略"]')}}</span>
                    </div>
                </template>
                <template v-else>
                    {{item[header.id]}}
                </template>
            </template>
        </cmdb-table>
        <cmdb-slider
            :width="740"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <v-confirm-details slot="content"></v-confirm-details>
        </cmdb-slider>
        <footer class="footer">
            <bk-button type="primary" @click="changeConfirm">
                {{$t('NetworkDiscovery["确认变更"]')}}
            </bk-button>
        </footer>
        <bk-dialog
        :is-show.sync="resultInfo.isShow" 
        :has-header="false" 
        :has-footer="false" 
        :quick-close="false" 
        :width="448">
            <div slot="content">
                <div>
                    <h3>{{$t('NetworkDiscovery["执行结果"]')}}</h3>
                    <p>{{$t('NetworkDiscovery["属性变更成功"]')}}</p>
                    <p>{{$t('NetworkDiscovery["关联关系变更成功"]')}}</p>
                    <p>{{$t('NetworkDiscovery["属性变更失败"]')}}</p>
                    <p>{{$t('NetworkDiscovery["关联关系变更失败"]')}}</p>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import vConfirmDetails from './details'
    export default {
        components: {
            vConfirmDetails
        },
        data () {
            return {
                resultInfo: {
                    isShow: false
                },
                slider: {
                    title: '',
                    isShow: false
                },
                filter: {
                    isShow: false
                },
                changeInfo: {
                    selected: '',
                    list: [{
                        id: 'create',
                        name: this.$t("Common['新增']")
                    }, {
                        id: 'change',
                        name: this.$t("networkDiscovery['变更']")
                    }, {
                        id: 'delete',
                        name: this.$t("Common['删除']")
                    }]
                },
                typeInfo: {
                    selected: '',
                    list: [{
                        id: 'switch',
                        name: this.$t("networkDiscovery['交换机']")
                    }, {
                        id: 'host',
                        name: this.$t("Hosts['主机']")
                    }]
                },
                table: {
                    header: [{
                        id: 'change',
                        name: this.$t('Hosts["云区域"]')
                    }, {
                        id: 'type',
                        name: this.$t('ModelManagement["类型"]')
                    }, {
                        id: 'unique',
                        name: this.$t('NetworkDiscovery["唯一标识"]')
                    }, {
                        id: 'ip',
                        name: 'IP'
                    }, {
                        id: 'config',
                        name: this.$t('NetworkDiscovery["配置信息"]')
                    }, {
                        id: 'time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
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
        },
        methods: {
            toggleFilter () {
                this.filter.isShow = !this.filter.isShow
            },
            showDetails () {
                this.slider.title = 'asdf'
                this.slider.isShow = true
            },
            changeConfirm () {
                this.resultInfo.isShow = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .network-confirm-wrapper {
        background: #f5f6fa;
        .filter-wrapper {
            &.open {
                >.bk-button {
                    background: #fafbfd;
                    border-bottom-color: transparent !important;
                    position: relative;
                    z-index: 2;
                    i {
                        transform: rotate(180deg);
                    }
                }
                .filter-details {
                    position: relative;
                    z-index: 1;
                }
            }
            >.bk-button {
                &:hover {
                    border-color: $cmdbBorderColor;
                }
                i {
                    transition: all .2s linear;
                }
            }
            .filter-details {
                padding: 11px 20px;
                background: #fafbfd;
                border: 1px solid $cmdbBorderColor;
                border-radius: 0 0 2px 2px;
                margin-top: -1px;
            }
            .details-left {
                float: left;
                font-size: 0;
                .bk-button {
                    margin-right: 10px;
                }
            }
            .details-right {
                float: right;
                .bk-selector {
                    float: left;
                    margin-right: 10px;
                    width: 140px;
                }
                .cmdb-form-input {
                    float: left;
                    margin-right: 10px;
                    width: 180px;
                }
            }
        }
        .confirm-table {
            margin-top: 20px;
            background: #fff;
        }
        .footer {
            position: fixed;
            bottom: 0;
            left: 0;
            padding: 8px 20px;
            width: 100%;
            text-align: right;
            background: #fff;
            box-shadow: 0 -2px 5px 0 rgba(0, 0, 0, 0.05);
        }
    }
</style>
