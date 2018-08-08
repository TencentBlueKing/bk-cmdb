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
    <div class="screening-wrapper">
        <form>
            <div class="screening-group" v-if="isShowBiz">
                <label class="screening-group-label">{{$t('Hosts[\'选择业务\']')}}</label>
                <div class="screening-group-item screening-group-item-app">
                    <v-application-selector
                        :filterable="true"
                        @on-selected="bkBizSelected"
                        :selected.sync="bkBizId">
                    </v-application-selector>
                </div>
            </div>
            <div class="screening-group">
                <label class="screening-group-label">IP</label>
                <div class="screening-group-item screening-group-item-ip">
                    <textarea class="bk-form-textarea" 
                        v-model.trim="ip.text">
                    </textarea>
                    <label class="bk-form-checkbox">
                        <input type="checkbox" v-model="ip['bk_host_innerip']" :disabled="!ip['bk_host_outerip']">
                        <span>{{$t('HostResourcePool[\'内网\']')}}</span>
                    </label>
                    <label class="bk-form-checkbox">
                        <input type="checkbox" v-model="ip['bk_host_outerip']" :disabled="!ip['bk_host_innerip']">
                        <span>{{$t('HostResourcePool[\'外网\']')}}</span>
                    </label>
                    <label class="bk-form-checkbox">
                        <input type="checkbox" v-model="ip.exact" :true-value="1" :false-value="0">
                        <span>{{$t('HostResourcePool[\'精确\']')}}</span>
                    </label>
                </div>
            </div>
            <div class="screening-group screening-group-scope" v-if="isShowScope">
                <label class="screening-group-label">{{$t("Hosts['搜索范围']")}}</label>
                <div class="screening-group-item screening-group-item-scope">
                    <label class="bk-form-checkbox">
                        <input type="checkbox" checked disabled>
                        <span>{{$t("Hosts['未分配主机']")}}</span>
                    </label>
                    <label class="bk-form-checkbox">
                        <input type="checkbox" v-model="isShowAssigned">
                        <span>{{$t("Hosts['已分配主机']")}}</span>
                    </label>
                </div>
            </div>
            <template v-for="(column, index) in localQueryColumns">
                <div class="screening-group" v-if="column['bk_property_id'] !== 'bk_host_innerip' && column['bk_property_id'] !== 'bk_host_outerip'">
                    <label class="screening-group-label">{{getColumnLabel(column)}}</label>
                    <div class="screening-group-item clearfix">
                        <!-- 时间无条件选择 -->
                        <template v-if="column['bk_property_type'] === 'date' || column['bk_property_type'] === 'time'">
                            <bk-daterangepicker ref="dateRangePicker" :class="['screening-group-item',`screening-group-item-${column['bk_property_type']}`]"
                                style="width: 315px;"
                                :quick-select="false"
                                :range-separator="'-'"
                                :align="'right'"
                                :initDate="localQueryColumnData[column['bk_property_id']]['value'].join(' - ')"
                                @change="setQueryDate(...arguments, column)">
                            </bk-daterangepicker>
                        </template>
                        <template v-else>
                            <!-- 判断条件选择 -->
                            <template v-if="typeOfChar.indexOf(column['bk_property_type']) !== -1 || typeOfAsst.indexOf(column['bk_property_type']) !== -1">
                                <template v-if="column['bk_property_id'] === 'bk_module_name' || column['bk_property_id'] === 'bk_set_name'">
                                    <bk-select class="screening-group-item-operator" :selected.sync="localQueryColumnData[column['bk_property_id']]['operator']">
                                        <bk-select-option v-for="(operator, index) in operators['name']"
                                            :key="index"
                                            :label="operator.label"
                                            :value="operator.value">
                                        </bk-select-option>
                                    </bk-select>
                                </template>
                                <template v-else>
                                    <bk-select class="screening-group-item-operator" :selected.sync="localQueryColumnData[column['bk_property_id']]['operator']">
                                        <bk-select-option v-for="(operator, index) in operators['char']"
                                            :key="index"
                                            :label="operator.label"
                                            :value="operator.value">
                                        </bk-select-option>
                                    </bk-select>
                                </template>
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
                            <!-- 判断输入类型 -->
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
                                    <template v-if="column['option']">
                                        <bk-select-option v-for="(option, index) in column['option']"
                                            :key="index"
                                            :value="option.id"
                                            :label="option.name">
                                        </bk-select-option>
                                    </template>
                                </bk-select>
                            </template>
                            <template v-else>
                                <input type="text" class="bk-form-input screening-group-item-value" v-model.trim="localQueryColumnData[column['bk_property_id']]['value']">
                            </template>
                        </template>
                    </div>
                </div>
            </template>
            <div class="screening-btn" ref="screeningBtn">
                <bk-button type="primary" :loading="$loading('hostSearch')" @click.prevent="refresh">{{$t('HostResourcePool[\'刷新查询\']')}}</bk-button>
            </div>
        </form>
    </div>
</template>
<script>
    import vApplicationSelector from '@/components/common/selector/application'
    import vMemberSelector from '@/components/common/selector/member'
    import bus from '@/eventbus/bus'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            queryColumns: {
                type: Array,
                required: true
            },
            queryColumnData: {
                type: Object,
                default () {
                    return {}
                }
            },
            attribute: {
                type: Array,
                required: true
            },
            isShowBiz: {
                type: Boolean,
                default: true
            },
            isShowScope: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                bkBizId: '',
                ip: {
                    'text': '',
                    'bk_host_innerip': true,
                    'bk_host_outerip': true,
                    'exact': 0
                },
                isShowAssigned: false,
                localQueryColumnData: {},
                localQueryColumns: [],
                operators: {
                    'default': [{
                        value: '$eq',
                        label: this.$t('Common[\'等于\']')
                    }, {
                        value: '$ne',
                        label: this.$t('Common[\'不等于\']')
                    }],
                    'char': [{
                        value: '$regex',
                        label: this.$t('Common[\'包含\']')
                    }, {
                        value: '$eq',
                        label: this.$t('Common[\'等于\']')
                    }, {
                        value: '$ne',
                        label: this.$t('Common[\'不等于\']')
                    }],
                    'name': [{
                        value: '$regex',
                        label: 'IN'
                    }, {
                        value: '$eq',
                        label: this.$t('Common[\'等于\']')
                    }, {
                        value: '$ne',
                        label: this.$t('Common[\'不等于\']')
                    }],
                    'date': [{
                        value: '$in',
                        label: this.$t('Common[\'包含\']')
                    }]
                },
                typeOfChar: ['singlechar', 'longchar'],
                typeOfDate: ['date', 'time'],
                typeOfAsst: ['singleasst', 'multiasst'],
                specialObj: {
                    'host': 'bk_host_innerip',
                    'biz': 'bk_biz_name',
                    'plat': 'bk_cloud_name',
                    'module': 'bk_module_name',
                    'set': 'bk_set_name'
                }
            }
        },
        computed: {
            ...mapGetters(['hostSearch']),
            ipFlag () {
                let flag = []
                if (this['ip']['bk_host_innerip']) {
                    flag.push('bk_host_innerip')
                }
                if (this['ip']['bk_host_outerip']) {
                    flag.push('bk_host_outerip')
                }
                return flag.join('|')
            },
            ipData () {
                let ipData = []
                this['ip']['text'].split(/\n|;|；|,|，/).map(ip => {
                    if (ip) {
                        ipData.push(ip.trim())
                    }
                })
                return ipData
            },
            filter () {
                return this.calcSearchParams()
            }
        },
        watch: {
            queryColumns (queryColumns) {
                let localQueryColumns = []
                let localQueryColumnData = {}
                queryColumns.map(column => {
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
            filter (filter) {
                this.$emit('filterChange', this.$deepClone(filter))
            },
            queryColumnData ({bk_biz_id: bkBizId, ip, condition}) {
                if (bkBizId) {
                    this.bkBizId = bkBizId
                }
                if (ip) {
                    this['ip']['text'] = ip['data'].join(',')
                    this['ip']['bk_host_innerip'] = ip['bk_host_innerip']
                    this['ip']['bk_host_outerip'] = ip['bk_host_outerip']
                    this['ip']['exact'] = ip['exact']
                }
                if (condition) {
                    condition.map(queryCondition => {
                        queryCondition.condition.map(({field, operator, value}) => {
                            if (this.localQueryColumnData.hasOwnProperty(field)) {
                                this.localQueryColumnData[field]['operator'] = operator
                                this.localQueryColumnData[field]['value'] = value
                            }
                        })
                    })
                }
            },
            localQueryColumns () {
                this.$nextTick(() => {
                    this.calcRefreshPosition()
                })
            }
        },
        created () {
            this.setHostSearchParams()
            this.$emit('filterChange', this.$deepClone(this.filter))
        },
        methods: {
            setHostSearchParams () {
                this.ip.text = this.hostSearch.ip
                this.ip.exact = this.hostSearch.exact
                this.ip['bk_host_innerip'] = this.hostSearch.innerip
                this.ip['bk_host_outerip'] = this.hostSearch.outerip
                this.isShowAssigned = this.hostSearch.assigned
            },
            calcSearchParams () {
                let filter = {
                    bk_biz_id: this.bkBizId,
                    ip: {
                        flag: this.ipFlag,
                        exact: this.ip.exact,
                        data: this.ipData
                    },
                    condition: [{
                        'bk_obj_id': 'biz',
                        'condition': [],
                        fields: []
                    }]
                }
                let bizCondition = filter.condition[0]['condition']
                if (this.isShowScope) {
                    if (!this.isShowAssigned) {
                        bizCondition.push({
                            field: 'default',
                            operator: '$eq',
                            value: 1
                        })
                    }
                } else {
                    bizCondition.push({
                        field: 'default',
                        operator: '$ne',
                        value: 1
                    })
                }
                Object.keys(this.localQueryColumnData).map(columnPropertyId => {
                    let column = this.localQueryColumnData[columnPropertyId]
                    let value = column.value
                    if ((!Array.isArray(value) && value !== '') || (Array.isArray(value) && value.length)) {
                        let property = this.getColumnProperty(columnPropertyId, column['bk_obj_id'])
                        let condition = filter.condition.find(({bk_obj_id: bkObjId}) => bkObjId === column['bk_obj_id'])
                        if (!condition) {
                            condition = {
                                'bk_obj_id': column['bk_obj_id'],
                                fields: [],
                                condition: []
                            }
                            filter.condition.push(condition)
                        }
                        if (property['bk_property_type'] === 'bool') {
                            condition.condition.push({
                                field: column.field,
                                operator: column.operator,
                                value: ['true', 'false'].includes(value) ? value === 'true' : value
                            })
                        } else if (this.typeOfAsst.indexOf(property['bk_property_type']) !== -1) {
                            filter.condition.push({
                                'bk_obj_id': property['bk_asst_obj_id'],
                                fields: [],
                                condition: [{
                                    field: this.specialObj.hasOwnProperty(property['bk_asst_obj_id']) ? this.specialObj[property['bk_asst_obj_id']] : 'bk_inst_name',
                                    operator: column.operator,
                                    value: column.value
                                }]
                            })
                        } else if (this.typeOfDate.indexOf(property['bk_property_type']) === -1) {
                            let operator = column.operator
                            let value = column.value
                            if (property['bk_property_id'] === 'bk_module_name' || property['bk_property_id'] === 'bk_set_name') {
                                operator = operator === '$regex' ? '$in' : operator
                                if (operator === '$in') {
                                    let arr = value.replace('，', ',').split(',')
                                    let isExist = arr.findIndex(val => {
                                        return val === value
                                    }) > -1
                                    value = isExist ? arr : [...arr, value]
                                }
                            }
                            condition.condition.push({
                                field: column.field,
                                operator: operator,
                                value: value
                            })
                        } else {
                            condition.condition.push({
                                field: column.field,
                                operator: '$gte',
                                value: `${column.value[0]} 00:00:00`
                            })
                            condition.condition.push({
                                field: column.field,
                                operator: '$lte',
                                value: `${column.value[1]} 23:59:59`
                            })
                        }
                    }
                })
                let defaultObj = ['host', 'module', 'set', 'biz']
                defaultObj.forEach(id => {
                    if (!filter.condition.some(({bk_obj_id: bkObjId}) => bkObjId === id)) {
                        filter.condition.push({
                            'bk_obj_id': id,
                            fields: [],
                            condition: []
                        })
                    }
                })
                return filter
            },
            bkBizSelected (app) {
                this.$emit('bkBizSelected', app.value)
            },
            getColumnLabel (column) {
                let columnProperty = this.getColumnProperty(column['bk_property_id'], column['bk_obj_id'])
                return `${columnProperty['bk_obj_name']} - ${columnProperty['bk_property_name']}`
            },
            getColumnProperty (columnPropertyId, columnObjId) {
                let columnProperty = null
                this.attribute.map(({bk_obj_id: bkObjId, bk_obj_name: bkObjName, properties}) => {
                    properties.map(property => {
                        if (property['bk_property_id'] === columnPropertyId && bkObjId === columnObjId) {
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
                    'text': '',
                    'bk_host_innerip': true,
                    'bk_host_outerip': true,
                    'exact': 1
                }
                if (this.$refs.dateRangePicker && this.$refs.dateRangePicker.length) {
                    this.$refs.dateRangePicker.map(vDateRangePicker => {
                        vDateRangePicker.selectedDateView = ''
                        vDateRangePicker.selectedDateRange = []
                        vDateRangePicker.selectedDateRangeTmp = []
                    })
                }
                Object.keys(this.localQueryColumnData).map(columnPropertyId => {
                    let column = this.localQueryColumnData[columnPropertyId]
                    let columnProperty = this.getColumnProperty(column['field'], column['bk_obj_id'])
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
            },
            refresh () {
                this.$emit('refresh')
            },
            calcRefreshPosition () {
                let screenTabpanelContent = this.$parent.$parent.$parent.$refs.screeningTabpanel.$el.parentElement
                let offsetHeight = screenTabpanelContent.offsetHeight
                let scrollHeight = screenTabpanelContent.scrollHeight
                let screeningBtnClassList = this.$refs.screeningBtn.classList
                if (scrollHeight > offsetHeight) {
                    screeningBtnClassList.add('fixed')
                } else {
                    screeningBtnClassList.remove('fixed')
                }
            }
        },
        beforeDestroy () {
            this.resetQueryColumnData()
        },
        components: {
            vApplicationSelector,
            vMemberSelector
        }
    }
</script>
<style lang="scss" scoped>
    .screening-wrapper{
        padding: 0 2px 0 0;
    }
    .screening-group{
        padding: 20px 0 0 0;
        &.screening-group-scope{
            padding: 10px 0 0 0;
            .screening-group-label{
                padding: 0;
            }
        }
        .screening-group-label{
            display: block;
            font-size: 14px;
            // color: #6b7baa;
            padding: 0 0 10px 0;
        }
        .screening-group-item{
            display: block;
            font-size: 0;
            .screening-group-item-ip{
                width: 100%;
                min-height: 70px;
                padding: 10px;
                font-size: 14px;
            }
            .screening-group-item-date,
            .screening-group-item-time{
                z-index: 1;
            }
            .screening-group-item-operator{
                width: 77px;
                font-size: 14px;
                float: left;
            }
            .screening-group-item-value{
                width: 224px;
                font-size: 14px;
                float: right;
            }
            .bk-form-checkbox{
                margin-right: 15px;
                &:last-child{
                    margin-right: 0;
                }
            }
        }
    }
    .screening-btn{
        padding: 20px 0 0 0;
        &.fixed{
            position: absolute;
            bottom: 0;
            left: 0;
            padding: 20px;
        }
    }
</style>
<style lang="scss">
.screening-group-item-date,
.screening-group-item-time{
    .daterange-dropdown-panel{
        min-width: auto;
    }
    .date-picker.date-select-container{
        &.start-date{
            display: none;
        }
        &.end-date{
            float: none !important;
            margin: 0 auto;
        }
    }
}
.screening-group-item-time{
    input[name="date-select"]{
        font-size: 12px;
    }
}
</style>