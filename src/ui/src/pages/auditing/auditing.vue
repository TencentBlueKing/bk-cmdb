/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="auditing-content">
        <div class="right-content">
            <div class="record-content">
                <div class="title-content clearfix">
                    <div class="group-content group-content-business">
                        <div class="selector-content selector-content-business">
                            <bk-select :selected.sync="filter.bkBizId" :filterable="true" :showClear="true" :placeholder="$t('OperationAudit[\'请选择业务\']')">
                                <bk-select-option v-for="(option, index) in bkBizList"
                                    :key="option['bk_biz_id']"
                                    :value="option['bk_biz_id']"
                                    :label="option['bk_biz_name']"></bk-select-option>
                            </bk-select>
                        </div>
                    </div>
                    <div class="group-content group-content-ip">
                        <span class="title-name">IP</span>
                        <div class="selector-content selector-content-ip">
                            <input type="text" :placeholder="$t('OperationAudit[\'使用逗号分隔\']')" v-model.trim="filter.bkIP">
                        </div>
                    </div>
                    <div class="group-content group-content-classify">
                        <span class="title-name">{{$t('OperationAudit["模型"]')}}</span>
                        <div class="selector-content selector-content-classify">
                            <bk-select 
                                :showClear="true"
                                :selected.sync="filter.classify"
                                :filterable="true">   
                                <template v-for="(classifyGroup, groupIndex) in filterClassifications">
                                    <bk-option-group v-if="classifyGroup['bk_objects'].length"
                                            :label="classifyGroup['bk_classification_name']"
                                            :key="groupIndex">
                                            <bk-select-option v-for="(classify,classifyIndex) in classifyGroup['bk_objects']"
                                                :key="classifyIndex"
                                                :value="classify['bk_obj_id']"
                                                :label="classify['bk_obj_name']">
                                            </bk-select-option>
                                    </bk-option-group>
                                </template>
                            </bk-select>
                        </div>
                    </div>
                    <div class="group-content group-content-type">
                        <span class="title-name">{{$t('OperationAudit[\'类型\']')}}</span>
                        <div class="selector-content selector-content-type">
                            <bk-select
                                :selected.sync="filter.bkOpType"
                                :list="operateTypeList"
                            >
                                <bk-select-option v-for="(operateType, operateTypeIndex) in operateTypeList"
                                    :key="operateTypeIndex"
                                    :value="operateType.value"
                                    :label="$t(operateType.label)"
                                ></bk-select-option>
                            </bk-select>
                        </div>
                    </div>
                    <div class="group-content group-content-time">
                        <span class="title-name">{{$t('OperationAudit[\'时间\']')}}</span>
                        <div class="selector-content selector-content-time">
                            <bk-daterangepicker
                                :range-separator="'-'"
                                :quick-select="true"
                                :start-date="startDate"
                                :end-date="endDate"
                                :ranges="ranges"
                                @change="setFilterTime"
                            ></bk-daterangepicker>
                        </div>
                    </div>
                    <div class="group-content group-content-btn fr">
                        <bk-button type="primary" :loading="$loading('auditSearch')" class="" @click="setCurrentPage(1)">{{$t('OperationAudit[\'查询\']')}}</bk-button>
                    </div>
                </div>
                <div class="table-content">
                    <v-table ref="table"
                        :header="tableHeader"
                        :list="tableList"
                        :pagination.sync="pagination"
                        :loading="$loading('auditSearch')"
                        :defaultSort="defaultSort"
                        :wrapperMinusHeight="150"
                        @handlePageChange="setCurrentPage"
                        @handleSizeChange="setCurrentSize"
                        @handleSortChange="setCurrentSort"
                        @handleRowClick="showDetails">
                    </v-table>
                </div>
            </div>
        </div>
        <v-sideslider :isShow.sync="details.isShow" :title="{text: $t('OperationAudit[\'操作详情\']')}">
            <v-history-details :details="details.data" slot="content"></v-history-details>
        </v-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import moment from 'moment'
    import vSideslider from '@/components/slider/sideslider'
    import vHistoryDetails from '@/components/history/details'
    import vTable from '@/components/table/table'
    export default {
        components: {
            vTable,
            vSideslider,
            vHistoryDetails
        },
        data () {
            return {
                isShowClearIcon: {
                    'classify': false,
                    'biz': false
                },
                operateTypeList: [{
                    value: '',
                    label: 'OperationAudit["全部"]'
                }, {
                    value: 1,
                    label: 'Common["新增"]'
                }, {
                    value: 2,
                    label: 'Common["修改"]'
                }, {
                    value: 3,
                    label: 'Common["删除"]'
                }, {
                    value: 100,
                    label: 'OperationAudit["关系变更"]'
                }],
                ranges: {
                    昨天: [moment().subtract(1, 'days'), moment()],
                    最近一周: [moment().subtract(7, 'days'), moment()],
                    最近一个月: [moment().subtract(1, 'month'), moment()],
                    最近三个月: [moment().subtract(3, 'month'), moment()]
                },
                rangesForZh: {
                    昨天: [moment().subtract(1, 'days'), moment()],
                    最近一周: [moment().subtract(7, 'days'), moment()],
                    最近一个月: [moment().subtract(1, 'month'), moment()],
                    最近三个月: [moment().subtract(3, 'month'), moment()]
                },
                rangesForEn: {
                    'Yesterday': [moment().subtract(1, 'days'), moment()],
                    'Last Week': [moment().subtract(7, 'days'), moment()],
                    'Last Month': [moment().subtract(1, 'month'), moment()],
                    'Last Three Month': [moment().subtract(3, 'month'), moment()]
                },
                filter: {               // 查询筛选参数
                    bkBizId: '',
                    bkIP: '',
                    classify: '',
                    bkOpType: '',
                    bkCreateTime: []
                },
                pagination: {
                    current: 1,
                    count: 0,
                    size: 10
                },
                tableHeader: [{
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
                tableList: [],
                defaultSort: '-op_time',
                sort: '-op_time',
                details: {
                    isShow: false,
                    data: null
                }
            }
        },
        computed: {
            /* 从store中回去操作对象列表 */
            ...mapGetters([
                'bkBizList',
                'language'
            ]),
            ...mapGetters('navigation', ['activeClassifications']),
            filterClassifications () {
                return [{
                    'bk_classification_id': 'bk_biz_topo',
                    'bk_classification_name': this.$t('BusinessTopology["业务拓扑"]'),
                    'bk_objects': [{
                        'bk_obj_id': 'set',
                        'bk_obj_name': this.$t('Hosts["集群"]')
                    }, {
                        'bk_obj_id': 'module',
                        'bk_obj_name': this.$t('Hosts["模块"]')
                    }]
                }, ...this.activeClassifications]
            },
            /* 开始时间 */
            startDate () {
                return this.$formatTime(moment().subtract(1, 'days'), 'YYYY-MM-DD')
            },
            /* 结束时间 */
            endDate () {
                return this.$formatTime(moment(), 'YYYY-MM-DD')
            },
            /* 搜索参数 */
            searchParams () {
                if (!this.filter.bkCreateTime.length) {
                    this.setFilterTime(null, `${this.startDate} - ${this.endDate}`)
                }
                let params = {
                    condition: {
                        op_time: this.filter.bkCreateTime
                    },
                    start: (this.pagination.current - 1) * this.pagination.size,
                    limit: this.pagination.size,
                    sort: this.sort
                }
                this.setParams(params.condition, 'bk_biz_id', this.filter.bkBizId)
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
                this.bkBizList.forEach((application, index) => {
                    applicationMap[application['bk_biz_id']] = application['bk_biz_name']
                })
                return applicationMap
            },
            /* 操作类型map */
            operateTypeMap () {
                let operateTypeMap = {}
                this.operateTypeList.forEach((operateType, index) => {
                    operateTypeMap[operateType['value']] = this.$t(operateType['label'])
                })
                return operateTypeMap
            }
        },
        watch: {
            'filter.bkBizId' (newId, oldId) {
                this.isShowClearIcon.biz = Boolean(newId)
            },
            'filter.classify' (classify) {
                this.isShowClearIcon.classify = Boolean(classify)
            },
            'operateTypeMap' (val) {
                this.initTableList()
            },
            'language' (lang) {
                this.tableHeader = [{
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
                }]
            }
        },
        async created () {
            await this.getBkBizList()
            this.getTableList()
            // if (!this.bkBizList.length) {
            //     }
            if (this.language === 'en') {
                this.ranges = this.rangesForEn
            } else {
                this.ranges = this.rangesForZh
            }
        },
        methods: {
            ...mapActions(['getBkBizList']),
            /* 获取表格数据 */
            getTableList () {
                this.$axios.post('audit/search/', this.searchParams, {id: 'auditSearch'}).then((res) => {
                    if (res.result) {
                        this.initTableList(res.data.info)
                        this.pagination.count = res.data.count
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /* 根据返回的结果设置一些表格显示内容 */
            initTableList (list) {
                if (list) {
                    list.forEach((item, index) => {
                        item['bk_biz_name'] = this.applicationMap[item['bk_biz_id']]
                        item['op_type_name'] = this.operateTypeMap[item['op_type']]
                        item['op_time'] = this.$formatTime(moment(item['op_time']))
                    })
                    this.tableList = list
                } else {
                    this.tableList.forEach((item, index) => {
                        item['op_type_name'] = this.operateTypeMap[item['op_type']]
                    })
                }
            },
            /* 日期选择时设置筛选参数 */
            setFilterTime (oldValue, newValue) {
                this.filter.bkCreateTime = newValue.split(' - ').map((time, index) => {
                    if (index === 0) {
                        return time + ' 00:00:00'
                    } else {
                        return time + ' 23:59:59'
                    }
                })
            },
            /* 表格排序事件 */
            setCurrentSort (sort) {
                this.sort = sort
                this.setCurrentPage(1)
            },
            /* 翻页事件，设置搜索参数页码 */
            setCurrentPage (current) {
                this.pagination.current = current
                this.getTableList()
            },
            /* 设置每页显示数量 */
            setCurrentSize (size) {
                this.pagination.size = size
                this.setCurrentPage(1)
            },
            /* 设置查询参数，有则添加 */
            setParams (obj, key, value) {
                if (value) {
                    obj[key] = value
                }
            },
            showDetails (item) {
                this.details.data = item
                this.details.isShow = true
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #737987; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    .auditing-content{
        height: 100%;
        font-size: 14px;
        color: $primaryColor;
        .dn{
            display: none;
        }
        .left-tap-contain{
            float: left;
            width:188px;
            border-right:1px solid $borderColor;
            float:left;
            border-left: none;
            border-top: none;
            height: 100%;
            .list-tap{
                height: 100%;
                overflow-y: auto;
                @include scrollbar;
                ul{
                    padding-top:30px;
                    >li{
                        width: 100px;
                        height: 40px;
                        line-height: 40px;
                        padding: 0 30px 0 44px;
                        width: 100%;
                        cursor: pointer;
                        font-size: 14px;
                        color: $primaryColor;
                        font-size: 14px;
                        position: relative;
                        white-space:nowrap;
                        text-overflow:ellipsis;
                        -o-text-overflow:ellipsis;
                        overflow: hidden;
                        margin-bottom: 8px;
                        .icon-left{
                            margin-left: -12px;
                        }
                        &:hover{
                            color: #498fe0;
                            background: #f9f9f9;
                            border-right:4px solid #498fe0;
                        }
                        .text{
                            padding:0 3px 0 5px;
                            min-width:64px;
                            vertical-align: top;
                        }
                        &.active{
                            color: #498fe0;
                            background: #f9f9f9;
                            border-right:4px solid #498fe0;
                        }

                    }
                }
            }
        }
        .right-content{
            float: right;
            padding: 0 20px;
            height: calc(100% - 52px);
            width: 100%;
            .pd-conrtol{
                padding: 30px 50px 50px 50px;
            }
            /*记录*/
            .record-content{
                .title-content{
                    width: 100%;
                    position: relative;
                    padding-top: 20px;
                    z-index: 2;
                    .group-content{
                        float:left;
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
                            display:inline-block;
                            font-size: 14px;
                            color: $primaryColor;
                            line-height: 36px;
                            padding-right: 10px;
                        }
                        .selector-content{
                            display:inline-block;
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
            }
            /*用户*/
            .user-content{
                height: 100%;
                .title-content{
                    width: 100%;
                    .title-name{
                        display: inline-block;
                        padding-top: 15px;
                        font-weight: bold;
                    }
                    .buttom-content{
                        font-size: 0;
                        button{
                            padding: 0 20px;
                            float: left;
                            font-size: 14px;
                            height: 36px;
                            line-height: 36px;
                            min-width: 110px;
                            border-radius: 2px;

                        }
                    }
                }
                .detail-content{
                    padding-top: 10px;
                    width: 100%;
                    height: calc(100% - 36px);
                    .list-content{
                        width: 50%;
                        border:1px solid $borderColor;
                        float: left;
                        height: 100%;
                        .list-title{
                            width: 100%;
                            height: 50px;
                            line-height: 50px;
                            background: #fff;
                            padding-left: 20px;
                            border-bottom: 1px solid $borderColor;
                        }
                        .list-detail{
                            width: 100%;
                            height: calc(100% - 50px);
                            .role-choice{
                                .bk-form-radio{
                                    display: block;
                                    color: $primaryColor;
                                    padding: 3px 0;
                                }
                                .paas-host{
                                    color: #498fe0;
                                    cursor: pointer;
                                }
                            }
                        }
                    }
                    .user-list{
                        border-right: 0;
                    }
                }
            }
            /*导航*/
            .navigation-content{
                height: 100%;
                width: 100%;
                .detail-content{
                    width: 100%;
                    height: 100%;
                    .list-content{
                        width: 50%;
                        border:1px solid $borderColor;
                        float: left;
                        height: 100%;
                        .list-title{
                            width: 100%;
                            height: 50px;
                            line-height: 50px;
                            background: #fff;
                            padding-left: 20px;
                            border-bottom: 1px solid $borderColor;
                        }
                        .button-content{
                            font-size: 0;
                            button{
                                text-align: center;
                                min-width: 90px;
                                height: 30px;
                                line-height: 30px;
                                border-radius: 2px;
                                padding: 0 20px;
                                float: left;
                                font-size: 14px;
                            }
                        }
                    }
                    .custom-nav-detail{
                        width: 100%;
                        height: calc(100% - 50px);
                        .btn-content{
                            font-size: 0;
                            position: absolute;
                            right: 30px;
                            top: 30px;
                            cursor: pointer;
                            i{
                                font-size: 14px;
                                margin-right: 5px;

                            }
                            span{
                                font-size: 14px;
                                line-height: 0;
                                &:hover{
                                    i{
                                        color:#ef4c4c;
                                    }
                                    i.icon-cc-round-plus{
                                        background: #fff;
                                        color: #498fe0;
                                    }
                                }
                            }
                        }
                    }
                    .user-list{
                        border-right: 0;
                    }
                    .custom-nav-list{
                        position: relative;
                        padding: 30px 45px 45px 45px;
                    }
                }
            }
        }
    }
    .bk-date{
        width: 240px;
    }
</style>
