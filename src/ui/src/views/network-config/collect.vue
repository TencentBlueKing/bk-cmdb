<template>
    <div class="collect-wrapper">
        <div class="title">
            <bk-button type="primary">
                {{$t('NetworkConfig["执行发现"]')}}
            </bk-button>
            <div class="input-box">
                <input type="text" class="cmdb-form-input" :placeholder="$t('NetworkConfig[\'搜索IP、云区域\']')">
                <i class="bk-icon icon-search"></i>
            </div>
        </div>
        <cmdb-table
            class="collect-table"
            :loading="$loading('searchUserGroup')"
            :checked="table.checked"
            :header="table.header"
            :list="table.list"
            :wrapperMinusHeight="240">
            <template slot="status" slot-scope="{ item }">
                <div class="status-wrapper" @mouseover="setTooltip($event, item)" @mouseleave="removeTooltip">
                    asdf
                </div>
            </template>
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary" @click="showConfig(item)">{{$t('EventPush["配置"]')}}</span>
            </template>
        </cmdb-table>
        <bk-dialog
        class="config-dialog"
        :is-show.sync="configDialog.isShow" 
        :has-header="false"
        :has-footer="false"
        :quick-close="false"
        :close-icon="false"
        padding="0"
        :width="424">
            <div slot="content" class="dialog-content">
                <div class="content-box">
                    <h2 class="title">
                        {{$t('NetworkConfig["配置采集器"]')}}
                    </h2>
                    <label>
                        <span>{{$t('NetworkConfig["SNMP扫描范围"]')}}</span>
                        <span class="color-danger">*</span>
                        <i class="bk-icon icon-exclamation-circle" v-tooltip="{content: htmlEncode(), classes: 'collect-tooltip'}"></i>
                    </label>
                    <textarea name="" id="" cols="30" rows="10"></textarea>
                    <label>
                        <span>{{$t('NetworkConfig["采集频率"]')}}</span>
                        <span class="color-danger">*</span>
                    </label>
                    <bk-selector
                        :list="configDialog.periodList"
                        :selected.sync="configDialog.period"
                    ></bk-selector>
                    <label>
                        <span>{{$t('NetworkConfig["团体字"]')}}</span>
                        <span class="color-danger">*</span>
                        <i class="bk-icon icon-exclamation-circle" v-tooltip="'Community String'"></i>
                    </label>
                    <input type="text" class="cmdb-form-input">
                    <div class="info">
                        <i class="bk-icon icon-exclamation-circle"></i>
                        <span>下发配置失败，请重新下发</span>
                    </div>
                </div>
                <footer class="footer">
                    <bk-button type="primary">
                        {{$t('NetworkConfig["保存并下发"]')}}
                    </bk-button>
                    <bk-button type="default">
                        {{$t('Common["取消"]')}}
                    </bk-button>
                </footer>
            </div>
        </bk-dialog>
        <div class="status-tips" ref='tooltipContent' v-show="tooltip.id">
            <p class="tips-content">采集器状态：<span>正常</span></p>
            <p class="tips-content">配置状态：<span class="color-success">asdf</span></p>
            <p class="tips-content">上报状态：<span>完成</span></p>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                table: {
                    header: [{
                        id: 'bk_inner_ip',
                        type: 'checkbox'
                    }, {
                        id: 'bk_cloud_name',
                        name: this.$t('Hosts["云区域"]')
                    }, {
                        id: 'bk_inner_ip',
                        name: this.$t('Common["内网IP"]')
                    }, {
                        id: 'status',
                        name: this.$t('ProcessManagement["状态"]')
                    }, {
                        id: 'version',
                        name: `${this.$t('NetworkConfig["版本"]')}（最新1.3）`
                    }, {
                        id: 'period',
                        name: this.$t('NetworkConfig["采集频率"]')
                    }, {
                        id: 'deploy_time',
                        name: this.$t('NetworkConfig["采集统计"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [{
                        bk_inner_ip: '1',
                        status: 'a'
                    }, {
                        bk_inner_ip: '2',
                        status: '2adsf'
                    }],
                    checked: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                },
                configDialog: {
                    isShow: false,
                    period: '',
                    periodList: [{
                        id: '',
                        name: this.$t('NetworkConfig["手动"]')
                    }]
                },
                tooltip: {
                    instance: null,
                    id: ''
                }
            }
        },
        methods: {
            ...mapActions('netDataCollection', [
                'searchDataCollection',
                'collectDataCollection'
            ]),
            showConfig () {
                this.configDialog.isShow = true
            },
            removeTooltip () {
                this.tooltip.instance && this.tooltip.instance.destroy()
            },
            setTooltip (event, item) {
                this.tooltip.id = item['bk_inner_ip']
                this.tooltip.instance && this.tooltip.instance.destroy()
                this.tooltip.instance = this.$tooltips({
                    duration: -1,
                    theme: 'light',
                    zIndex: 9999,
                    container: document.body,
                    target: event.target
                })
                this.tooltip.instance.$el.append(this.$refs.tooltipContent)
            },
            htmlEncode () {
                let temp = document.createElement('div')
                temp.innerHTML = `填写格式&lt;/br&gt;指定IP：192.168.1.1&lt;/br&gt;IP范围：192.168.1.1-192.168.1.200&lt;/br&gt;cidr ip 范围：192.168.1.1/32`
                let output = temp.innerText
                temp = null
                return output
            },
            async getTableData () {
                let pagination = this.table.pagination
                let params = {
                    page: {
                        start: (pagination.current - 1) * pagination.size,
                        limit: pagination.size,
                        sort: this.table.sort
                    }
                }
                const res = await this.searchDataCollection({params, config: {requestId: 'searchDataCollection'}})
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
    .collect-wrapper {
        padding-top: 20px;
        >.title {
            .input-box {
                position: relative;
                float: right;
                input {
                    width: 260px;
                }
                .icon-search {
                    position: absolute;
                    top: 9px;
                    right: 8px;
                    font-size: 18px;
                    color: $cmdbTableBorderColor;
                }
            }
        }
        .collect-table {
            margin-top: 20px;
            background: #fff;
        }
        .config-dialog {
            .dialog-content {
                .content-box {
                    padding: 20px;
                    >h2 {
                        color: #333948;
                        font-size: 22px;
                        line-height: 1;
                    }
                    >label {
                        display: block;
                        margin: 15px 0 5px;
                        span,
                        i {
                            vertical-align: middle;
                        }
                    }
                    >textarea {
                        width: 100%;
                        height: 80px;
                        border-color: $cmdbBorderColor;
                        resize: none;
                        outline: none;
                        border-radius: 2px;
                    }
                    .info {
                        margin-top: 20px;
                        background: #fff3da;
                        border-radius: 2px;
                        width: 100%;
                        padding-left: 20px;
                        height: 42px;
                        line-height: 40px;
                        font-size: 0;
                        border: 1px solid #ffc947;
                        .bk-icon {
                            position: relative;
                            top: -1px;
                            margin-right: 10px;
                            color: #ffc947;
                            font-size: 20px;
                        }
                        span {
                            font-size: 14px;
                            vertical-align: middle;
                        }
                    }
                }
                .footer {
                    border-top: 1px solid #e5e5e5;
                    padding-right: 20px;
                    text-align: right;
                    font-size: 0;
                    background: #fafbfd;
                    height: 54px;
                    line-height: 54px;
                    .bk-default {
                        margin-left: 10px;
                    }
                }
            }
        }
    }
    .status-tips {
        padding: 5px 10px;
        font-size: 12px;
    }
</style>

<style lang="scss">
    .collect-tooltip {
        .tooltip-inner {
            max-width: 300px;
        }
    }
</style>
