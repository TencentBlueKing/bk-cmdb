<template>
    <div class="host-filter-list clearfix">
        <div class="screening-group">
            <label class="screening-group-label">{{$t('Hosts[\'选择业务\']')}}</label>
            <div class="screening-group-item screening-group-item-app">
                <bk-select :filterable="true" :selected.sync="bkBizId" @on-selected="bkBizSelected">
                    <bk-select-option v-for="biz in bkBizList"
                        :key="biz['bk_biz_id']"
                        :value="biz['bk_biz_id']"
                        :label="biz['bk_biz_name']">
                    </bk-select-option>
                </bk-select>
            </div>
        </div>
        <div class="screening-group">
            <label class="screening-group-label">{{$t("Common['内网IP']")}}</label>
            <div class="screening-group-item screening-group-item-ip">
                <input class="bk-form-input" v-model.trim="ip.text"></input>
            </div>
        </div>
        <template v-for="(column, index) in localQueryColumns">
            <div class="screening-group clearfix" v-if="column['bk_property_id'] !== 'bk_host_innerip' && column['bk_property_id'] !== 'bk_host_outerip'">
                <label class="screening-group-label" :title="getColumnLabel(column)">{{getColumnLabel(column)}}</label>
                <div class="screening-group-item clearfix" :class="`screening-group-item-${column['bk_property_type']}`">
                    <!-- 时间无条件选择 -->
                    <template v-if="column['bk_property_type'] === 'date' || column['bk_property_type'] === 'time'">
                        <bk-daterangepicker ref="dateRangePicker" class="screening-group-item screening-group-item-date"
                            :quick-select="false"
                            :range-separator="'-'"
                            :align="'right'"
                            :initDate="localQueryColumnData[column['bk_property_id']]['value'].join(' - ')"
                            @change="setQueryDate(...arguments, column)">
                        </bk-daterangepicker>
                    </template>
                    <template v-else>
                        <!-- 判断条件选择 -->
                        <div class="operation-type">
                            <template v-if="typeOfChar.indexOf(column['bk_property_type']) !== -1 || typeOfAsst.indexOf(column['bk_property_type']) !== -1">
                                <bk-select class="screening-group-item-operator" :selected.sync="localQueryColumnData[column['bk_property_id']]['operator']">
                                    <bk-select-option v-for="(operator, index) in operators['char']"
                                        :key="index"
                                        :label="operator.label"
                                        :value="operator.value">
                                    </bk-select-option>
                                </bk-select>
                            </template>
                            <template v-else>
                                <bk-select class="screening-group-item-operator" :selected.sync="localQueryColumnData[column['bk_property_id']]['operator']">
                                    <bk-select-option v-for="(operator, index) in operators['default']"
                                        :key="index"
                                        :label="operator.label"
                                        :value="operator.value">
                                    </bk-select-option>
                                </bk-select>
                            </template>
                        </div>
                        <!-- 判断输入类型 -->
                        <div class="operation-value">
                            <template v-if="column['bk_property_type'] === 'int'">
                                <input type="text" maxlength="11" class="bk-form-input screening-group-item-value" v-model.number="localQueryColumnData[column['bk_property_id']]['value']">
                            </template>
                            <template v-else-if="column['bk_property_type'] === 'objuser'">
                                <v-member-selector class="screening-group-item-value"
                                    :exclude="true"
                                    :selected.sync="localQueryColumnData[column['bk_property_id']]['value']"
                                    :multiple="true">
                                </v-member-selector>
                            </template>
                            <template v-else-if="column['bk_property_type'] === 'enum'">
                                <bk-select class="screening-group-item-value" :selected.sync="localQueryColumnData[column['bk_property_id']]['value']">
                                    <bk-select-option v-for="(option, index) in column['bk_option']"
                                        :key="index"
                                        :value="option.id"
                                        :label="option.name">
                                    </bk-select-option>
                                </bk-select>
                            </template>
                            <template v-else>
                                <input type="text" class="bk-form-input screening-group-item-value" v-model.trim="localQueryColumnData[column['bk_property_id']]['value']">
                            </template>
                        </div>
                    </template>
                </div>
            </div>
        </template>
    </div>
</template>

<script>
    import vApplicationSelector from '@/components/common/selector/application'
    import vMemberSelector from '@/components/common/selector/member'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            queryColumns: {
                type: Array,
                required: true
            },
            attribute: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                bkBizId: -1,
                ip: {
                    'text': ''
                },
                localQueryColumnData: {},
                localQueryColumns: [],
                operators: {
                    'default': [{
                        value: '$eq',
                        label: this.$t("Common['等于']")
                    }, {
                        value: '$ne',
                        label: this.$t("Common['不等于']")
                    }],
                    'char': [{
                        value: '$regex',
                        label: this.$t("Common['包含']")
                    }, {
                        value: '$eq',
                        label: this.$t("Common['等于']")
                    }, {
                        value: '$ne',
                        label: this.$t("Common['不等于']")
                    }],
                    'date': [{
                        value: '$in',
                        label: this.$t("Common['包含']")
                    }]
                },
                typeOfChar: ['singlechar', 'longchar'],
                typeOfDate: ['date', 'time'],
                typeOfAsst: ['singleasst', 'multiasst']
            }
        },
        computed: {
            ...mapGetters(['bkBizList']),
            ipData () {
                let ipData = []
                this['ip']['text'].split(/\n|;|；|,|，/).map(ip => {
                    if (ip) {
                        ipData.push(ip)
                    }
                })
                return ipData
            },
            filter () {
                let filter = {
                    bk_biz_id: this.bkBizId,
                    ip: {
                        flag: 'bk_host_innerip',
                        exact: 0,
                        data: this.ipData
                    },
                    condition: [{
                        'bk_obj_id': 'host',
                        fields: [],
                        condition: []
                    }, {
                        'bk_obj_id': 'biz',
                        fields: [],
                        condition: [{
                            field: 'default',
                            operator: '$ne',
                            value: 1
                        }]
                    }, {
                        'bk_obj_id': 'module',
                        fields: [],
                        condition: []
                    }, {
                        'bk_obj_id': 'set',
                        fields: [],
                        condition: []
                    }]
                }
                Object.keys(this.localQueryColumnData).map(columnPropertyId => {
                    this.attribute.map(({bk_obj_id: bkObjId, properties}) => {
                        properties.map(property => {
                            if (columnPropertyId === property['bk_property_id']) {
                                let value = this.localQueryColumnData[columnPropertyId]['value']
                                if (this.typeOfAsst.indexOf(property['bk_property_type']) !== -1) {
                                    if (value) {
                                        filter.condition.push({
                                            'bk_obj_id': property['bk_asst_obj_id'],
                                            fields: [],
                                            condition: [{
                                                field: 'bk_inst_name',
                                                operator: this.localQueryColumnData[columnPropertyId]['operator'],
                                                value: value
                                            }]
                                        })
                                    }
                                } else {
                                    filter.condition.map(({bk_obj_id: filterBkObjId, condition}) => {
                                        if (filterBkObjId === bkObjId) {
                                            let isEmptyValue = false
                                            if (value === '' || (Array.isArray(value) && !value.length)) {
                                                isEmptyValue = true
                                            }
                                            if (!isEmptyValue) {
                                                condition.push({
                                                    field: columnPropertyId,
                                                    operator: this.localQueryColumnData[columnPropertyId]['operator'],
                                                    value: value
                                                })
                                            }
                                        }
                                    })
                                }
                            }
                        })
                    })
                })
                return filter
            }
        },
        watch: {
            queryColumns (queryColumns) {
                this.initLocalQuery()
            },
            filter (filter) {
                this.$emit('filterChange', Object.assign({}, filter))
            }
        },
        created () {
            this.initLocalQuery()
            if (this.bkBizList.length) {
                this.bkBizId = this.bkBizList[0]['bk_biz_id']
            }
        },
        methods: {
            initLocalQuery () {
                let localQueryColumns = []
                let localQueryColumnData = {}
                this.queryColumns.map(column => {
                    let {
                        bk_property_id: bkPropertyId,
                        bk_property_type: bkPropertyType,
                        bk_obj_id: bkObjId
                    } = column
                    let columnProperty = this.getColumnProperty(bkPropertyId, bkObjId)
                    if (columnProperty) {
                        localQueryColumns.push(column)
                        let operatorType = 'default'
                        if (this.typeOfChar.indexOf(bkPropertyType) !== -1) {
                            operatorType = 'char'
                        } else if (this.typeOfDate.indexOf(bkPropertyType) !== -1) {
                            operatorType = 'date'
                        }
                        localQueryColumnData[bkPropertyId] = {
                            field: bkPropertyId,
                            value: operatorType === 'date' ? [] : '',
                            operator: this.operators[operatorType][0]['value'],
                            'bk_obj_id': bkObjId
                        }
                    }
                })
                this.localQueryColumns = localQueryColumns
                this.localQueryColumnData = localQueryColumnData
            },
            bkBizSelected (app) {
                this.$emit('bkBizSelected', app.value)
            },
            getColumnLabel (column) {
                let columnProperty = this.getColumnProperty(column['bk_property_id'], column['bk_obj_id'])
                return `${columnProperty['bk_property_name']}`
            },
            getColumnProperty (columnPropertyId, columnObjId) {
                let columnProperty = null
                this.attribute.map(({bk_obj_id: bkObjId, bk_obj_name: bkObjName, properties}) => {
                    properties.map(property => {
                        if (property['bk_property_id'] === columnPropertyId && columnObjId === bkObjId) {
                            columnProperty = Object.assign({bk_obj_name: bkObjName}, property)
                        }
                    })
                })
                return columnProperty
            },
            setQueryDate (oldDate, newDate, column) {
                this.localQueryColumnData[column['bk_property_id']]['value'] = newDate.split(' - ')
            },
            resetQueryColumnData () {
                this.ip = {
                    'text': ''
                }
                Object.keys(this.localQueryColumnData).map(columnPropertyId => {
                    let column = this.localQueryColumnData[columnPropertyId]
                    let columnProperty = this.getColumnProperty(columnPropertyId, column['bk_obj_id'])
                    let {
                        bk_property_type: bkPropertyType,
                        bk_property_id: bkPropertyId
                    } = columnProperty
                    if (bkPropertyType === 'date' || bkPropertyType === 'time') {
                        this.localQueryColumnData[bkPropertyId]['value'] = []
                    } else {
                        this.localQueryColumnData[bkPropertyId]['value'] = ''
                    }
                    let operatorType = 'default'
                    if (this.typeOfChar.indexOf(bkPropertyType) !== -1) {
                        operatorType = 'char'
                    } else if (this.typeOfDate.indexOf(bkPropertyType) !== -1) {
                        operatorType = 'date'
                    }
                    this.localQueryColumnData[bkPropertyId]['operator'] = this.operators[operatorType][0]['value']
                })
            }
        },
        components: {
            vMemberSelector,
            vApplicationSelector
        }
    }
</script>

<style lang="scss" scoped>
    .host-filter-list{
        padding: 15px 0 10px;
    }
    .screening-group{
        float: left;
        margin-bottom: 10px;
        width: 300px;
        &:nth-child(2n){
            float: right;
        }
        .screening-group-label{
            float: left;
            width: 60px;
            line-height: 36px;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }
        .screening-group-item{
            float: left;
            &.screening-group-item-app{
                width: 240px;
            }
            &.screening-group-item-ip{
                width: 240px;
            }
            .operation-type{
                position: relative;
                float: left;
                width: 65px;
                &:hover{
                    z-index: 2;
                }
            }
            .operation-value{
                position: relative;
                float: left;
                width: 175px;
                margin-left: -1px;
                z-index: 1;
            }
        }
    }
</style>
