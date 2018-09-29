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
            :loading="$loading('searchNetcollectList')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
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
            <bk-button type="primary" @click="resultDialog.isShow = true">
                {{$t('NetworkDiscovery["确认变更"]')}}
            </bk-button>
        </footer>
        <bk-dialog
        class="result-dialog"
        :is-show.sync="resultDialog.isShow" 
        :has-header="false"
        :has-footer="false"
        :quick-close="false"
        :close-icon="false"
        :width="448">
            <div slot="content">
                <h2>{{$t('NetworkDiscovery["执行结果"]')}}</h2>
                <div class="dialog-content">
                    <p>
                        <span class="info">{{$t('NetworkDiscovery["属性变更成功"]')}}</span>
                        <span class="number">22条</span>
                    </p>
                    <p>
                        <span class="info">{{$t('NetworkDiscovery["关联关系变更成功"]')}}</span>
                        <span class="number">22条</span>
                    </p>
                    <p class="fail">
                        <span class="info">{{$t('NetworkDiscovery["属性变更失败"]')}}</span>
                        <span class="number">22条</span>
                    </p>
                    <p class="fail">
                        <span class="info">{{$t('NetworkDiscovery["关联关系变更失败"]')}}</span>
                        <span class="number">22条</span>
                    </p>
                </div>
                <div class="dialog-details">
                    <p @click="toggleDialogDetails">
                        <i class="bk-icon icon-angle-down"></i>
                        <span>{{$t('NetworkDiscovery["展开详情"]')}}</span>
                    </p>
                    <transition name="toggle-slide">
                        <div class="detail-content-box" v-if="resultDialog.isDetailsShow">
                            <div class="detail-content">
                                Lorem ipsum dolor sit amet consectetur adipisicing elit. Officiis repellat sequi, eum fugiat, consectetur sunt omnis minus exercitationem in dolorum, asperiores hic nobis perspiciatis dignissimos dolorem non ipsam! Adipisci, facere!
                                Lorem ipsum dolor sit amet consectetur adipisicing elit. Officiis repellat sequi, eum fugiat, consectetur sunt omnis minus exercitationem in dolorum, asperiores hic nobis perspiciatis dignissimos dolorem non ipsam! Adipisci, facere!
                            </div>
                        </div>
                    </transition>
                </div>
                <footer class="footer">
                    <bk-button type="primary">
                        {{$t('Hosts["确认"]')}}
                    </bk-button>
                </footer>
            </div>
        </bk-dialog>
        <bk-dialog
        class="confirm-dialog"
        :is-show.sync="confirmDialog.isShow" 
        :title="$t('NetworkDiscovery[\'是否确认变更\']')"
        :has-footer="false"
        :quick-close="false"
        padding="0"
        :width="390">
            <div slot="content" class="dialog-content">
                <p>
                    {{$t('NetworkDiscovery["要在返回前确认变更吗？"]')}}
                </p>
                <footer class="footer">
                    <bk-button type="primary">
                        {{$t('NetworkDiscovery["确认变更"]')}}
                    </bk-button>
                    <bk-button type="default">
                        {{$t('NetworkDiscovery["丢弃"]')}}
                    </bk-button>
                    <bk-button type="default">
                        {{$t('Common["取消"]')}}
                    </bk-button>
                </footer>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import vConfirmDetails from './details'
    export default {
        components: {
            vConfirmDetails
        },
        data () {
            return {
                resultDialog: {
                    isShow: false,
                    isDetailsShow: true
                },
                confirmDialog: {
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
                        id: 'update',
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
                        id: 'action',
                        name: this.$t('NetworkDiscovery["变更方式"]')
                    }, {
                        id: 'device_type',
                        name: this.$t('ModelManagement["类型"]')
                    }, {
                        id: 'device_name',
                        name: this.$t('NetworkDiscovery["唯一标识"]')
                    }, {
                        id: 'bk_inner_ip',
                        name: 'IP'
                    }, {
                        id: 'device_attributes',
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
                        action: 'create',
                        device_type: 'switch',
                        device_name: 'asdf',
                        bk_inner_ip: '192.168.1.1',
                        device_attributes: '24个10/100M自适应RJ4端口',
                        last_time: '2018-04-17T15:00:49.274+08:00'
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
        computed: {
            params () {
                let params = {}
                return params
            }
        },
        methods: {
            ...mapActions('netDiscovery', [
                'searchNetcollectList'
            ]),
            toggleFilter () {
                this.filter.isShow = !this.filter.isShow
            },
            showDetails () {
                this.slider.title = 'asdf'
                this.slider.isShow = true
            },
            changeConfirm () {
                this.resultDialog.isShow = true
            },
            toggleDialogDetails () {
                this.resultDialog.isDetailsShow = !this.resultDialog.isDetailsShow
            },
            async getTableData () {
                const res = await this.searchNetcollectList({params: this.params, config: {requestId: 'searchNetcollectList'}})
                this.table.pagination.count = res.count
                this.table.list = res.info
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .toggle-slide-enter-active, .toggle-slide-leave-active{
        transition: height .2s;
        overflow: hidden;
        height: 190px;
    }
    .toggle-slide-enter, .toggle-slide-leave-to{
        height: 0 !important;
    }
    .network-confirm-wrapper {
        background: $cmdbBackgroundColor;
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
        >.footer {
            position: fixed;
            bottom: 0;
            left: 0;
            padding: 8px 20px;
            width: 100%;
            text-align: right;
            background: #fff;
            box-shadow: 0 -2px 5px 0 rgba(0, 0, 0, 0.05);
        }
        .result-dialog {
            h2 {
                margin-bottom: 10px;
                font-size: 22px;
                color: #333948;
            }
            .dialog-content {
                >p {
                    line-height: 26px;
                    span {
                        display: inline-block;
                    }
                    .info {
                        width: 155px;
                    }
                }
                .fail {
                    color: $cmdbDangerColor;
                }
            }
            .dialog-details {
                margin-top: 10px;
                >p {
                    font-weight: bold;
                    cursor: pointer;
                    .icon-angle-down {
                        font-size: 12px;
                        font-weight: bold;
                    }
                }
                .dialog-content-box {
                    height: 220px;
                }
                .detail-content {
                    margin-top: 10px;
                    padding: 15px 20px;
                    border: 1px dashed #dde4eb;
                    background: #fafbfd;
                    border-radius: 5px;
                    overflow-y: auto;
                    height: 190px;
                    @include scrollbar;
                }
            }
            .footer {
                border-top: 1px solid #e5e5e5;
                padding-right: 20px;
                margin: 25px -20px -20px;
                text-align: right;
                font-size: 0;
                background: #fafbfd;
                height: 54px;
                line-height: 54px;
            }
        }
        .confirm-dialog {
            .dialog-content {
                text-align: center;
                >p {
                    margin: 10px 0 20px;
                }
                .footer {
                    padding-bottom: 40px;
                    font-size: 0;
                    .bk-button.bk-default {
                        margin-left: 10px;
                    }
                }
            }
        }
    }
</style>
