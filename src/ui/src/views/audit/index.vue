<template>
    <div class="audit-wrapper">
        <div class="title-content clearfix">
            <div class="group-content group-content-business">
                <div class="selector-content selector-content-business">
                    <bk-selector
                        :list="business"
                        :selected.sync="filter.bizId"
                        :searchable="true"
                        :allow-clear="true"
                        display-key="bk_biz_name"
                        search-key="bk_biz_id"
                        setting-key="bk_biz_id"
                    ></bk-selector>
                </div>
            </div>
            <div class="group-content group-content-ip">
                <span class="title-name">IP</span>
                <div class="selector-content selector-content-ip">
                    <input class="cmdb-form-input" type="text" :placeholder="$t('OperationAudit[\'使用逗号分隔\']')" v-model.trim="filter.bkIP">
                </div>
            </div>
            <div class="group-content group-content-classify">
                <span class="title-name">{{$t('OperationAudit["模型"]')}}</span>
                <div class="selector-content selector-content-classify">
                    <bk-selector
                        :list="filterClassifications"
                        :selected.sync="filter.classify"
                        :has-children="true"
                        :searchable="true"
                    ></bk-selector>
                </div>
            </div>
            <div class="group-content group-content-type">
                <span class="title-name">{{$t('OperationAudit[\'类型\']')}}</span>
                <div class="selector-content selector-content-type">
                    <bk-selector
                        :list="operateTypeList"
                        :selected.sync="filter.bkOpType"
                    ></bk-selector>
                </div>
            </div>
            <div class="group-content group-content-time">
                <span class="title-name">{{$t('OperationAudit[\'时间\']')}}</span>
                <div class="selector-content selector-content-time">
                    <cmdb-form-date-range class="date-range" position="left" v-model="filter.bkCreateTime"></cmdb-form-date-range>
                </div>
            </div>
            <div class="group-content group-content-btn fr">
                <bk-button type="primary" :loading="$loading('getOperationLog')" @click="handlePageChange(1)">{{$t('OperationAudit[\'查询\']')}}</bk-button>
            </div>
        </div>
        <cmdb-table
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange"
            @handleRowClick="handleRowClick"
        ></cmdb-table>
        <cmdb-slider
            :isShow.sync="details.isShow"
            :title="$t('OperationAudit[\'操作详情\']')">
            <v-details :details="details.data" slot="content"></v-details>
        </cmdb-slider>
    </div>
</template>

<script>
    import vDetails from '@/components/audit-history/details'
    import { mapActions, mapGetters } from 'vuex'
    import moment from 'moment'
    export default {
        data () {
            return {
                filter: { // 查询筛选参数
                    bizId: '',
                    bkIP: '',
                    classify: '',
                    bkOpType: '',
                    bkCreateTime: []
                },
                operateTypeList: [{
                    id: '',
                    name: this.$t('OperationAudit["全部"]')
                }, {
                    id: 1,
                    name: this.$t('Common["新增"]')
                }, {
                    id: 2,
                    name: this.$t('Common["修改"]')
                }, {
                    id: 3,
                    name: this.$t('Common["删除"]')
                }, {
                    id: 100,
                    name: this.$t('OperationAudit["关系变更"]')
                }],
                table: {
                    header: [{
                        id: 'operator',
                        name: this.$t('OperationAudit["操作账号"]')
                    }, {
                        id: 'op_target',
                        name: this.$t('OperationAudit["对象"]')
                    }, {
                        id: 'op_desc',
                        name: this.$t('OperationAudit["描述"]')
                    }, {
                        id: 'bk_biz_name',
                        name: this.$t('OperationAudit["所属业务"]'),
                        sortKey: 'bk_biz_id'
                    }, {
                        id: 'ext_key',
                        name: 'IP'
                    }, {
                        id: 'op_type_name',
                        name: this.$t('OperationAudit["类型"]'),
                        sortKey: 'op_type'
                    }, {
                        id: 'op_time',
                        name: this.$t('OperationAudit["操作时间"]')
                    }],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time'
                },
                details: {
                    isShow: false,
                    data: null
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', [
                'business'
            ]),
            filterClassifications () {
                let classifications = []
                this.$classifications.map(classify => {
                    if (classify['bk_classification_id'] === 'bk_biz_topo') {
                        classifications.push({
                            name: classify['bk_classification_name'],
                            children: [{
                                id: 'set',
                                name: this.$t('Hosts["集群"]')
                            }, {
                                id: 'module',
                                name: this.$t('Hosts["模块"]')
                            }]
                        })
                    } else if (classify['bk_classification_id'] !== 'bk_host_manage') {
                        if (classify['bk_objects'].length) {
                            let children = []
                            classify['bk_objects'].map(({bk_obj_id: bkObjId, bk_obj_name: bkObjName}) => {
                                children.push({
                                    id: bkObjId,
                                    name: bkObjName
                                })
                            })
                            classifications.push({
                                name: classify['bk_classification_name'],
                                children
                            })
                        }
                    }
                })
                return classifications
            },
            params () {
                let opTime = []
                if (this.filter.bkCreateTime.length) {
                    opTime = [
                        this.filter.bkCreateTime[0] ? this.filter.bkCreateTime[0] + ' 00:00:00' : '',
                        this.filter.bkCreateTime[1] ? this.filter.bkCreateTime[1] + ' 23:59:59' : ''
                    ]
                }
                let params = {
                    condition: {
                        op_time: opTime
                    },
                    start: (this.table.pagination.current - 1) * this.table.pagination.size,
                    limit: this.table.pagination.size,
                    sort: this.table.sort
                }
                this.setParams(params.condition, 'bk_biz_id', this.filter.bizId)
                this.setParams(params.condition, 'op_type', this.filter.bkOpType)
                this.setParams(params.condition, 'op_target', this.filter.classify)
                if (this.filter.bkIP) { // 将IP分隔成查询数组
                    let ipArray = []
                    this.filter.bkIP.split(',').map((ip, index) => {
                        if (ip) {
                            ipArray.push(ip.trim())
                        }
                    })
                    this.setParams(params.condition, 'ext_key', {$in: ipArray})
                }
                return params
            },
            /* 业务ID与Name的mapping */
            applicationMap () {
                let applicationMap = {}
                this.business.forEach((application, index) => {
                    applicationMap[application['bk_biz_id']] = application['bk_biz_name']
                })
                return applicationMap
            },
            /* 操作类型map */
            operateTypeMap () {
                let operateTypeMap = {}
                this.operateTypeList.forEach((operateType, index) => {
                    operateTypeMap[operateType['id']] = operateType['name']
                })
                return operateTypeMap
            }
        },
        methods: {
            ...mapActions('operationAudit', ['getOperationLog']),
            setParams (obj, key, value) {
                if (value) {
                    obj[key] = value
                }
            },
            initDate () {
                this.filter.bkCreateTime = [
                    this.$tools.formatTime(moment().subtract(1, 'days'), 'YYYY-MM-DD'),
                    this.$tools.formatTime(moment(), 'YYYY-MM-DD')
                ]
                this.getTableData()
            },
            async getTableData () {
                const res = await this.getOperationLog({
                    params: this.params,
                    config: {
                        cancelPrevious: true,
                        requestId: 'getOperationLog'
                    }
                })
                this.initTableList(res.info)
                this.table.pagination.count = res.count
            },
            initTableList (list) {
                if (list) {
                    list.map(item => {
                        item['bk_biz_name'] = this.applicationMap[item['bk_biz_id']]
                        item['op_type_name'] = this.operateTypeMap[item['op_type']]
                        item['op_time'] = this.$tools.formatTime(moment(item['op_time']))
                    })
                    this.table.list = list
                }
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.getTableData()
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleRowClick (item) {
                this.details.data = item
                this.details.isShow = true
            }
        },
        mounted () {
            this.initDate()
        },
        components: {
            vDetails
        }
    }
</script>

<style lang="scss" scoped>
    .title-content{
        .group-content{
            float: left;
            margin: 0 1.9% 20px 0;
            white-space: nowrap;
            font-size: 0;
            &.group-content-business{
                width: calc(120 / (1020 - 348) * (100% - 348px));
                .selector-content-business{
                    width: 100%;
                }
            }
            &.group-content-ip{
                width: calc(145 / (1020 - 348) * (100% - 348px));
                .selector-content-ip{
                    width: calc(122 / 145 * 100%);
                }
            }
            &.group-content-classify{
                width: calc(165 / (1020 - 348) * (100% - 348px));
                .selector-content-classify{
                    width: calc(127 / 165 * 100%);
                }
            }
            &.group-content-type{
                width: calc(145 / (1020 - 348) * (100% - 348px));
                .selector-content-type{
                    width: calc(107 / 145 * 100%);
                }
            }
            &.group-content-time{
                width: 280px;
                margin-right: 0;
                .date-range{
                    width: 240px;
                }
            }
            &.group-content-btn{
                width: auto;
                margin-right: 0;
            }
            .search-btn{
                padding: 0 19px;
                height: 36px;
                line-height: 34px;
                font-size: 14px;
            }
            .title-name{
                display: inline-block;
                font-size: 14px;
                line-height: 36px;
                padding-right: 10px;
                vertical-align: top;
            }
            .selector-content{
                display: inline-block;
                font-size: 14px;
                position: relative;
            }
            &:nth-child(2){
                input{
                    height: 36px;
                    line-height: 36px;
                    width: 100%;
                    padding: 0 6px;
                }
            }
        }
    }
</style>
