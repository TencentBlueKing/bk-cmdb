<template>
    <div class="network-confirm-wrapper">
        <div class="filter-wrapper" :class="{'open': filter.isShow}">
            <bk-button type="default" @click="toggleFilter">
                {{$t('NetworkDiscovery["批量操作"]')}}
                <i class="bk-icon icon-angle-down"></i>
            </bk-button>
            <div class="filter-details clearfix" v-show="filter.isShow">
                <div class="details-left">
                    <bk-button type="default" @click="toggleIgnore(true)">
                        {{$t('NetworkDiscovery["忽略"]')}}
                    </bk-button>
                    <bk-button type="default" @click="toggleIgnore(false)">
                        {{$t('NetworkDiscovery["取消忽略"]')}}
                    </bk-button>
                    <label class="cmdb-form-checkbox">
                        <input type="checkbox" v-model="filter.isShowIgnore">
                        <span class="cmdb-checkbox-text">{{$t('NetworkDiscovery["显示忽略"]')}}</span>
                    </label>
                </div>
                <div class="details-right clearfix">
                    <bk-selector
                        :list="changeList"
                        :allow-clear="true"
                        :selected.sync="filterCopy.action"
                        :placeholder="$t('NetworkDiscovery[\'全部变更\']')"
                    ></bk-selector>
                    <bk-selector
                        :list="typeList"
                        :allow-clear="true"
                        :selected.sync="filterCopy['bk_obj_name']"
                        :placeholder="$t('NetworkDiscovery[\'全部类型\']')"
                    ></bk-selector>
                    <input type="text" class="cmdb-form-input" :placeholder="$t('NetworkDiscovery[\'请输入IP\']')">
                    <bk-button type="default" @click="search">
                        {{$t('Common["查询"]')}}
                    </bk-button>
                </div>
            </div>
        </div>
        <cmdb-table
            class="confirm-table"
            :loading="$loading('searchNetcollectList')"
            :header="table.header"
            :list="tableList"
            :wrapperMinusHeight="210"
            :defaultSort="table.defaultSort"
            :checked.sync="table.checked"
            @handleSortChange="handleSortChange"
            @handleCheckAll="handleCheckAll">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <label class="table-checkbox bk-form-checkbox bk-checkbox-small"
                    :key="index"
                    v-if="header.id === 'id'" 
                    @click.stop>
                    <input type="checkbox"
                        :value="item['bk_inst_key']"
                        v-model="table.checked">
                </label>
                <template v-else-if="header.id === 'operation'">
                    <div :key="index">
                        <span class="text-primary" @click.stop="showDetails(item)">{{$t('NetworkDiscovery["详情"]')}}</span>
                        <span class="text-primary" @click.stop="item.ignore = !item.ignore">{{item.ignore ? $t('NetworkDiscovery["取消忽略"]') : $t('NetworkDiscovery["忽略"]')}}</span>
                    </div>
                </template>
                <template v-else-if="header.id === 'action'">
                    <span :key="index" :class="{'ignore': item.ignore, 'color-warning': item.action === 'update', 'color-danger' : item.action === 'delete'}">{{actionMap[item.action]}}</span>
                </template>
                <template v-else-if="header.id === 'last_time'">
                    <span :key="index" :class="{'ignore': item.ignore}">{{$tools.formatTime(item['last_time'])}}</span>
                </template>
                <template v-else>
                    <span :key="index" :class="{'ignore': item.ignore}">{{item[header.id]}}</span>
                </template>
            </template>
        </cmdb-table>
        <cmdb-slider
            :width="740"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <v-confirm-details
            slot="content"
            :ignore="activeItem.ignore"
            :attributes.sync="activeItem.attributes"
            :associations.sync="activeItem.associations"
            :detailPage="detailPage"
            @toggleSwitcher="toggleSwitcher"
            @updateView="updateView"
            ></v-confirm-details>
        </cmdb-slider>
        <div class="footer">
            <bk-button type="primary" @click="showResultDialog">
                {{$t('NetworkDiscovery["确认变更"]')}}
            </bk-button>
        </div>
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
                        <span class="number">{{resultDialog.data['change_attribute_success']}}条</span>
                    </p>
                    <p>
                        <span class="info">{{$t('NetworkDiscovery["关联关系变更成功"]')}}</span>
                        <span class="number">{{resultDialog.data['change_associations_success']}}条</span>
                    </p>
                    <p class="fail">
                        <span class="info">{{$t('NetworkDiscovery["属性变更失败"]')}}</span>
                        <span class="number">{{resultDialog.data['change_attribute_failure']}}条</span>
                    </p>
                    <p class="fail">
                        <span class="info">{{$t('NetworkDiscovery["关联关系变更失败"]')}}</span>
                        <span class="number">{{resultDialog.data['change_associations_failure']}}条</span>
                    </p>
                </div>
                <div class="dialog-details" v-if="resultDialog.data.errors.length">
                    <p @click="toggleDialogDetails">
                        <i class="bk-icon icon-angle-down"></i>
                        <span>{{$t('NetworkDiscovery["展开详情"]')}}</span>
                    </p>
                    <transition name="toggle-slide">
                        <div class="detail-content-box" v-if="resultDialog.isDetailsShow">
                            <div class="detail-content">
                                <div v-for="(item, index) in resultDialog.data.errors" :key="index">
                                    {{item}}
                                </div>
                            </div>
                        </div>
                    </transition>
                </div>
                <div class="footer">
                    <bk-button type="primary" @click="resultDialog.isShow = false">
                        {{$t('Hosts["确认"]')}}
                    </bk-button>
                </div>
            </div>
        </bk-dialog>
        <bk-dialog
            class="confirm-dialog"
            :is-show.sync="confirmDialog.isShow" 
            :title="$t('NetworkDiscovery[\'退出确认\']')"
            :has-footer="false"
            :quick-close="false"
            padding="0"
            :width="390">
            <div slot="content" class="dialog-content">
                <p>
                    {{$t('NetworkDiscovery["当前改动尚未生效，是否放弃？"]')}}
                </p>
                <div class="footer">
                    <bk-button type="default" @click="routeToLeave">
                        {{$t('NetworkDiscovery["放弃改动"]')}}
                    </bk-button>
                    <bk-button type="default" @click="confirmDialog.isShow = false">
                        {{$t('Common["取消"]')}}
                    </bk-button>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import vConfirmDetails from './details'
    export default {
        components: {
            vConfirmDetails
        },
        data () {
            return {
                resultDialog: {
                    isShow: false,
                    isDetailsShow: true,
                    data: {
                        errors: []
                    }
                },
                confirmDialog: {
                    isShow: false,
                    isLeave: false,
                    leaveResolver: null
                },
                slider: {
                    title: '',
                    isShow: false
                },
                filter: {
                    isShow: false,
                    isShowIgnore: true,
                    action: '',
                    bk_obj_name: '',
                    bk_host_innerip: ''
                },
                filterCopy: {
                    action: '',
                    bk_obj_name: '',
                    bk_host_innerip: ''
                },
                changeList: [{
                    id: 'create',
                    name: this.$t("Common['新增']")
                }, {
                    id: 'update',
                    name: this.$t("NetworkDiscovery['变更']")
                }, {
                    id: 'delete',
                    name: this.$t("Common['删除']")
                }],
                typeList: [{
                    id: 'switch',
                    name: this.$t("NetworkDiscovery['交换机']")
                }, {
                    id: 'host',
                    name: this.$t("Hosts['主机']")
                }],
                table: {
                    header: [{
                        id: 'action',
                        name: this.$t('NetworkDiscovery["变更方式"]')
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('ModelManagement["类型"]')
                    }, {
                        id: 'bk_inst_key',
                        name: this.$t('NetworkDiscovery["唯一标识"]')
                    }, {
                        id: 'bk_host_innerip',
                        name: 'IP'
                    }, {
                        id: 'configuration',
                        name: this.$t('NetworkDiscovery["配置信息"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [],
                    listCopy: [],
                    checked: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                },
                actionMap: {
                    'create': this.$t("Common['新增']"),
                    'update': this.$t("NetworkDiscovery['变更']"),
                    'delete': this.$t("Common['删除']")
                },
                activeItem: {
                    index: 0,
                    ignore: false,
                    attributes: [],
                    associations: []
                }
            }
        },
        computed: {
            ...mapGetters('netDiscovery', ['cloudName']),
            tableList () {
                return this.table.list.filter(item => {
                    if (!this.filter.isShowIgnore && item.ignore) {
                        return false
                    }
                    if (this.filter['bk_obj_name'] !== '' && item['bk_obj_name'] !== this.filter['bk_obj_name']) {
                        return false
                    }
                    if (this.filter.action !== '' && item.action !== this.filter.action) {
                        return false
                    }
                    if (!item['bk_host_innerip'].includes(this.filter.bk_host_innerip)) {
                        return false
                    }
                    return true
                })
            },
            detailPage () {
                let index = this.tableList.findIndex(({bk_inst_key: instKey}) => instKey === this.activeItem['bk_inst_key'])
                return {
                    prev: index === 0,
                    next: index === this.tableList.length - 1
                }
            }
        },
        beforeRouteEnter (to, from, next) {
            next(vm => {
                if (vm.cloudName === null) {
                    vm.$router.push({name: 'networkDiscovery'})
                }
            })
        },
        async beforeRouteLeave (to, from, next) {
            if (this.cloudName === null) {
                next()
            } else {
                if (JSON.stringify(this.table.list) !== JSON.stringify(this.table.listCopy)) {
                    this.confirmDialog.isShow = true
                    await new Promise(async (resolve, reject) => {
                        this.confirmDialog.leaveResolver = () => {
                            resolve()
                        }
                    })
                    this.confirmDialog.isShow = false
                }
                next()
            }
        },
        created () {
            this.$route.meta.title = `${this.cloudName}${this.$t('NetworkDiscovery["变更确认"]')}`
            this.getTableData()
        },
        methods: {
            ...mapActions('netDiscovery', [
                'searchNetcollectList',
                'searchNetcollectChange',
                'confirmNetcollectChange'
            ]),
            updateView (type) {
                let index = this.tableList.findIndex(({bk_inst_key: instKey}) => instKey === this.activeItem['bk_inst_key'])
                if (type === 'prev') {
                    this.activeItem = this.tableList[index - 1]
                } else {
                    this.activeItem = this.tableList[index + 1]
                }
            },
            search () {
                this.filter['bk_obj_id'] = this.filterCopy['bk_obj_id']
                this.filter.action = this.filterCopy.action
                this.filter['bk_host_innerip'] = this.filterCopy['bk_host_innerip']
            },
            toggleIgnore (ignore) {
                this.table.checked.map(instKey => {
                    let item = this.table.list.find(({bk_inst_key: bkInstKey}) => bkInstKey === instKey)
                    if (item) {
                        item.ignore = ignore
                    }
                })
            },
            toggleSwitcher (value) {
                this.activeItem.ignore = value
            },
            routeToLeave () {
                if (this.confirmDialog.leaveResolver) {
                    this.confirmDialog.leaveResolver()
                }
            },
            toggleFilter () {
                if (this.filter.isShow) {
                    this.table.header.shift()
                } else {
                    this.table.header.unshift({
                        type: 'checkbox',
                        id: 'id'
                    })
                }
                this.filter.isShow = !this.filter.isShow
            },
            showDetails (item) {
                this.activeItem = item
                this.slider.title = item['bk_host_innerip']
                this.slider.isShow = true
            },
            async showResultDialog () {
                let params = {
                    reports: []
                }
                this.table.list.forEach(item => {
                    if (!item.ignore) {
                        let detail = {
                            bk_cloud_id: item['bk_cloud_id'],
                            bk_obj_id: item['bk_obj_id'],
                            bk_inst_key: item['bk_inst_key'],
                            action: item['action'],
                            configuration: item['configuration'],
                            bk_host_innerip: item['bk_host_innerip'],
                            last_time: item['last_time'],
                            attributes: [],
                            associations: []
                        }
                        item.attributes.forEach(attr => {
                            if (attr.method === 'accept') {
                                detail.attributes.push({
                                    bk_property_id: attr['bk_property_id'],
                                    bk_property_name: attr['bk_property_name'],
                                    value: attr['value'],
                                    method: 'accept'
                                })
                            }
                        })
                        item.associations.forEach(asst => {
                            if (asst.method === 'accept') {
                                detail.associations.push({
                                    bk_asst_inst_name: asst['bk_asst_inst_name'],
                                    bk_asst_obj_id: asst['bk_asst_obj_id'],
                                    bk_asst_obj_name: asst['bk_asst_obj_name'],
                                    bk_obj_asst_id: asst['bk_obj_asst_id'],
                                    bk_asst_property_id: asst['bk_asst_property_id'],
                                    method: 'accept'
                                })
                            }
                        })
                        params.reports.push(detail)
                    }
                })
                try {
                    const res = await this.confirmNetcollectChange({params, config: {globalError: false, requestId: 'confirmNetcollectChange', transformData: false}})
                    this.resultDialog.data = res.data
                } catch (e) {
                    this.$error(e.data['bk_error_msg'])
                }
                this.resultDialog.isShow = true
                this.getTableData()
            },
            toggleDialogDetails () {
                this.resultDialog.isDetailsShow = !this.resultDialog.isDetailsShow
            },
            async getTableData () {
                const res = await this.searchNetcollectList({params: {bk_cloud_id: Number(this.$route.params.cloudId)}, config: {requestId: 'searchNetcollectList'}})
                res.info.map(item => {
                    Object.assign(item, {ignore: false})
                    item.attributes.map(attr => Object.assign(attr, {method: 'accept'}))
                    item.associations.map(relation => Object.assign(relation, {method: 'accept'}))
                })
                this.table.list = res.info
                this.table.listCopy = this.$tools.clone(res.info)
            },
            handleSortChange (sort) {
                let key = sort
                if (sort[0] === '-') {
                    key = sort.substr(1, sort.length - 1)
                }
                this.table.list.sort((itemA, itemB) => {
                    if (itemA[key] === null) {
                        itemA[key] = ''
                    }
                    if (itemB[key] === null) {
                        itemB[key] = ''
                    }
                    return itemA[key].localeCompare(itemB[key])
                })
                if (sort[0] === '-') {
                    this.table.list.reverse()
                }
            },
            handleCheckAll () {
                this.table.checked = this.table.list.map(item => item['bk_inst_key'])
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
            .ignore {
                color: $cmdbBorderColor;
                &.color-danger {
                    color: $cmdbDangerColor;
                    opacity: .6;
                }
                &.color-warning {
                    color: $cmdbWarningColor;
                    opacity: .6;
                }
            }
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
