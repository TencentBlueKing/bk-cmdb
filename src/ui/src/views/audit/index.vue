<template>
    <div class="audit-wrapper">
        <div class="title-content">
            <div class="group-content" v-if="isAdminView">
                <span class="title-name" :title="$t('业务')">{{$t('业务')}}</span>
                <div class="selector-content">
                    <bk-select v-model="filter.bizId" searchable font-size="medium">
                        <bk-option v-for="business in authorizedBusiness"
                            :key="business.bk_biz_id"
                            :id="business.bk_biz_id"
                            :name="business.bk_biz_name">
                        </bk-option>
                    </bk-select>
                </div>
            </div>
            <div class="group-content">
                <span class="title-name" title="IP">IP</span>
                <div class="selector-content">
                    <bk-input class="cmdb-form-input" type="text"
                        font-size="medium"
                        :placeholder="$t('使用逗号分隔')"
                        v-model.trim="filter.bkIP">
                    </bk-input>
                </div>
            </div>
            <div class="group-content">
                <span class="title-name" :title="$t('模型')">{{$t('模型')}}</span>
                <div class="selector-content">
                    <bk-select v-model="filter.classify" searchable font-size="medium">
                        <bk-option-group v-for="group in filterClassifications"
                            :key="group.id"
                            :name="group.name">
                            <bk-option v-for="classify in group.children"
                                :key="classify.id"
                                :id="classify.id"
                                :name="classify.name">
                            </bk-option>
                        </bk-option-group>
                    </bk-select>
                </div>
            </div>
            <div class="group-content">
                <span class="title-name" :title="$t('类型')">{{$t('类型')}}</span>
                <div class="selector-content">
                    <bk-select
                        font-size="medium"
                        v-model="filter.bkOpType"
                        :clearable="false">
                        <bk-option v-for="option in operateTypeList"
                            :key="option.id"
                            :id="option.id"
                            :name="option.name">
                        </bk-option>
                    </bk-select>
                </div>
            </div>
            <div class="group-content">
                <span class="title-name" :title="$t('时间')">{{$t('时间')}}</span>
                <div class="selector-content date-range-content">
                    <cmdb-form-date-range
                        class="date-range"
                        :clearable="false"
                        font-size="medium"
                        v-model="filter.bkCreateTime">
                    </cmdb-form-date-range>
                </div>
            </div>
            <div class="group-content button-group">
                <bk-button theme="primary" :loading="$loading('getOperationLog')" @click="handlePageChange(1, $event)">{{$t('查询')}}</bk-button>
            </div>
        </div>
        <bk-table
            v-bkloading="{ isLoading: $loading('getOperationLog') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 190"
            :row-style="{ cursor: 'pointer' }"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange"
            @row-click="handleRowClick">
            <bk-table-column
                sortable="custom"
                prop="operator"
                class-name="is-highlight"
                :label="$t('操作账号')">
            </bk-table-column>
            <bk-table-column
                sortable="custom"
                prop="op_target"
                :label="$t('对象')">
            </bk-table-column>
            <bk-table-column
                sortable="custom"
                prop="op_desc"
                :label="$t('描述')">
            </bk-table-column>
            <bk-table-column
                sortable="custom"
                prop="bk_biz_id"
                :label="$t('所属业务')">
                <template slot-scope="{ row }">{{row.bk_biz_name}}</template>
            </bk-table-column>
            <bk-table-column
                sortable="custom"
                prop="ext_key"
                label="IP">
            </bk-table-column>
            <bk-table-column
                sortable="custom"
                prop="op_type"
                :label="$t('类型')">
                <template slot-scope="{ row }">{{row.op_type_name}}</template>
            </bk-table-column>
            <bk-table-column
                sortable="custom"
                prop="op_time"
                :label="$t('操作时间')">
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :quick-close="true"
            :is-show.sync="details.isShow"
            :width="800"
            :title="$t('操作详情')">
            <v-details :details="details.data" slot="content" v-if="details.isShow"></v-details>
        </bk-sideslider>
    </div>
</template>

<script>
    import vDetails from '@/components/audit-history/details'
    import { mapActions, mapGetters } from 'vuex'
    import moment from 'moment'
    export default {
        components: {
            vDetails
        },
        data () {
            return {
                filter: { // 查询筛选参数
                    bizId: '',
                    bkIP: '',
                    classify: '',
                    bkOpType: 0,
                    bkCreateTime: []
                },
                operateTypeList: [{
                    id: 0,
                    name: this.$t('全部')
                }, {
                    id: 1,
                    name: this.$t('新增')
                }, {
                    id: 2,
                    name: this.$t('修改')
                }, {
                    id: 3,
                    name: this.$t('删除')
                }, {
                    id: 100,
                    name: this.$t('关系变更')
                }],
                table: {
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time',
                    stuff: {
                        // 如果初始化时有查询参数，则认为处在查询模式
                        type: Object.keys(this.$route.query).length ? 'search' : 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                details: {
                    isShow: false,
                    data: null
                }
            }
        },
        computed: {
            ...mapGetters(['isAdminView']),
            ...mapGetters('objectBiz', [
                'authorizedBusiness',
                'bizId'
            ]),
            ...mapGetters('objectModelClassify', ['classifications']),
            filterClassifications () {
                const classifications = []
                this.classifications.map(classify => {
                    if (classify['bk_classification_id'] === 'bk_biz_topo') {
                        classifications.push({
                            name: classify['bk_classification_name'],
                            children: [{
                                id: 'set',
                                name: this.$t('集群')
                            }, {
                                id: 'module',
                                name: this.$t('模块')
                            }]
                        })
                    } else if (classify['bk_classification_id'] !== 'bk_host_manage') {
                        if (classify['bk_objects'].length) {
                            const children = []
                            classify['bk_objects'].map(({ bk_obj_id: bkObjId, bk_obj_name: bkObjName }) => {
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
                const params = {
                    condition: {
                        op_time: opTime
                    },
                    start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                    limit: this.table.pagination.limit,
                    sort: this.table.sort
                }
                this.setParams(params.condition, 'bk_biz_id', this.isAdminView ? this.filter.bizId : this.bizId)
                this.setParams(params.condition, 'op_type', this.filter.bkOpType)
                this.setParams(params.condition, 'op_target', this.filter.classify)
                if (this.filter.bkIP) { // 将IP分隔成查询数组
                    const ipArray = []
                    this.filter.bkIP.split(',').map((ip, index) => {
                        if (ip) {
                            ipArray.push(ip.trim())
                        }
                    })
                    this.setParams(params.condition, 'ext_key', { $in: ipArray })
                }
                return params
            },
            /* 业务ID与Name的mapping */
            applicationMap () {
                const applicationMap = {}
                this.authorizedBusiness.forEach((application, index) => {
                    applicationMap[application['bk_biz_id']] = application['bk_biz_name']
                })
                return applicationMap
            },
            /* 操作类型map */
            operateTypeMap () {
                const operateTypeMap = {}
                this.operateTypeList.forEach((operateType, index) => {
                    operateTypeMap[operateType['id']] = operateType['name']
                })
                return operateTypeMap
            }
        },
        created () {
            this.$store.dispatch('objectBiz/getAuthorizedBusiness')
        },
        mounted () {
            this.initDate()
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
            async getTableData (event) {
                try {
                    const res = await this.getOperationLog({
                        params: this.params,
                        config: {
                            globalPermission: false,
                            cancelPrevious: true,
                            requestId: 'getOperationLog'
                        }
                    })
                    this.initTableList(res.info)
                    this.table.pagination.count = res.count
                    // 有传入event参数时认为来自用户搜索
                    if (event) {
                        this.table.stuff.type = 'search'
                    }
                } catch ({ permission }) {
                    this.table.list = []
                    // 从api调用层抛出的错误，仅当权限问题时会注入permission
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                }
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
            handlePageChange (current, event) {
                this.table.pagination.current = current
                this.getTableData(event)
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handleRowClick (item) {
                this.details.data = item
                this.details.isShow = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .audit-wrapper {
        padding: 0 20px;
    }
    .title-content{
        padding: 0 0 14px 0;
        display: flex;
        align-items: center;
        justify-content: flex-start;
        flex-direction: row;
        flex-wrap: nowrap;
        .group-content{
            flex: 1 1;
            margin: 0 1.5% 0 0;
            white-space: nowrap;
            font-size: 0;
            &.button-group {
                text-align: right;
                flex: 0 0 90px;
                margin: 0;
            }
            .selector-content {
                display: inline-block;
                vertical-align: middle;
                width: calc(100% - 40px);
                .bk-select {
                    width: 100%;
                }
            }
            .search-btn{
                padding: 0 19px;
                height: 36px;
                line-height: 34px;
                font-size: 14px;
            }
            .title-name{
                display: inline-block;
                vertical-align: middle;
                min-width: 40px;
                max-width: 48px;
                font-size: 14px;
                padding-right: 5px;
                @include ellipsis;
            }
            .date-range-content {
                width: calc(100% - 10px);
                .date-range {
                    width: 100%;
                }
            }
        }
    }
</style>
