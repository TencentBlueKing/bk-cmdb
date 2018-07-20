<template>
    <div class="history-details-wrapper">
        <template v-if="details">
            <div class="history-info">
                <div :class="['info-group', info.key]" v-for="(info, index) in informations">
                    <label class="info-label">{{$t(info.label)}}：</label>
                    <span class="info-value">
                        {{info.hasOwnProperty('optionKey') ? options[info.optionKey][details[info.key]] : details[info.key]}}
                    </span>
                </div>
            </div>
            <v-table
                :loading="loadingAttribute"
                :sortable="false"
                :width="width? width : 700"
                :wrapperMinusHeight="300"
                :emptyHeight="230"
                :header="tableHeader"
                :list="tableList"
                :rowBorder="true"
                :colBorder="true"
                v-bind="height ? {height} : {}">
                <template slot="pre_data" slot-scope="{item}" v-html="item['pre_data']">
                    <div :class="['details-data', {'has-changed': hasChanged(item)}]" v-html="item['pre_data']"></div>
                </template>
                <template slot="cur_data" slot-scope="{item}">
                    <div :class="['details-data', {'has-changed': hasChanged(item)}]" v-html="item['cur_data']"></div>
                </template>
            </v-table>
            <p class="field-btn" @click="toggleFields" v-if="this.details.op_type !== 1 && this.details.op_type !== 3">
                {{isShowAllFields ? $t('EventPush["收起"]') : $t('EventPush["展开"]')}}
            </p>
        </template>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'
    import vTable from '@/components/table/table'
    export default {
        components: {
            vTable
        },
        props: {
            details: Object,
            height: Number,
            width: Number,
            isShow: Boolean
        },
        data () {
            return {
                isShowAllFields: false,
                loadingAttribute: true,
                informations: [{
                    label: 'OperationAudit[\'操作账号\']',
                    key: 'operator'
                }, {
                    label: 'OperationAudit[\'所属业务\']',
                    key: 'bk_biz_id',
                    optionKey: 'biz'
                }, {
                    label: 'OperationAudit[\'IP\']',
                    key: 'ext_key'
                }, {
                    label: 'OperationAudit[\'类型\']',
                    key: 'op_type',
                    optionKey: 'opType'
                }, {
                    label: 'OperationAudit[\'对象\']',
                    key: 'op_target'
                }, {
                    label: 'OperationAudit[\'操作时间\']',
                    key: 'op_time'
                }, {
                    label: 'OperationAudit[\'描述\']',
                    key: 'op_desc'
                }],
                colWidth: [130, 280, 280]
            }
        },
        computed: {
            ...mapGetters(['bkBizList']),
            ...mapGetters('object', ['attribute']),
            objId () {
                return this.details ? this.details['op_target'] : null
            },
            options () {
                let biz = {}
                this.bkBizList.forEach(({bk_biz_id: bkBizId, bk_biz_name: bkBizName}) => {
                    biz[bkBizId] = bkBizName
                })
                let opType = {
                    1: this.$t("Common['新增']"),
                    2: this.$t("Common['修改']"),
                    3: this.$t("Common['删除']"),
                    100: this.$t('OperationAudit["关系变更"]')
                }
                return {
                    biz,
                    opType
                }
            },
            tableHeader () {
                let header = [{
                    id: 'bk_property_name',
                    name: '',
                    width: 130
                }]
                const preDataHeader = {
                    id: 'pre_data',
                    name: this.$t("OperationAudit['变更前']")
                }
                const curDataHeader = {
                    id: 'cur_data',
                    name: this.$t("OperationAudit['变更后']")
                }
                if (this.details['op_type'] === 1) {
                    header.push(curDataHeader)
                } else if (this.details['op_type'] === 2 || this.details['op_type'] === 100) {
                    header.push(preDataHeader)
                    header.push(curDataHeader)
                } else if (this.details['op_type'] === 3) {
                    header.push(preDataHeader)
                }
                return header
            },
            tableList () {
                let list = []
                const attribute = (this.attribute[this.objId] || []).filter(({bk_isapi: bkIsapi}) => !bkIsapi)
                if (this.details['op_type'] !== 100) {
                    attribute.forEach(property => {
                        let preData = this.getCellValue(property, 'pre_data')
                        let curData = this.getCellValue(property, 'cur_data')
                        if (!this.isShowAllFields) {
                            if (preData !== curData) {
                                list.push({
                                    'bk_property_name': property['bk_property_name'],
                                    'pre_data': preData,
                                    'cur_data': curData
                                })
                            }
                        } else {
                            list.push({
                                'bk_property_name': property['bk_property_name'],
                                'pre_data': preData,
                                'cur_data': curData
                            })
                        }
                    })
                } else {
                    const content = this.details.content
                    const preBizId = content['pre_data']['bk_biz_id']
                    const curBizId = content['cur_data']['bk_biz_id']
                    const preModule = content['pre_data']['module'] || []
                    const curModule = content['cur_data']['module'] || []
                    let pre = []
                    let cur = []
                    preModule.forEach(module => {
                        pre.push(`${this.options.biz[preBizId]}→${module.set[0]['ref_name']}→${module['ref_name']}`)
                    })
                    curModule.forEach(module => {
                        cur.push(`${this.options.biz[curBizId]}→${module.set[0]['ref_name']}→${module['ref_name']}`)
                    })
                    let preData = pre.join('<br>')
                    let curData = cur.join('<br>')
                    if (!this.isShowAllFields) {
                        if (preData !== curData) {
                            list.push({
                                'bk_property_name': this.$t('Hosts["关联关系"]'),
                                'pre_data': preData,
                                'cur_data': curData
                            })
                        }
                    } else {
                        list.push({
                            'bk_property_name': this.$t('Hosts["关联关系"]'),
                            'pre_data': preData,
                            'cur_data': curData
                        })
                    }
                }
                return list
            }
        },
        watch: {
            async objId (objId) {
                if (objId && !this.attribute.hasOwnProperty(objId)) {
                    this.loadingAttribute = true
                    await this.$store.dispatch('object/getAttribute', {objId})
                    this.loadingAttribute = false
                } else {
                    this.loadingAttribute = false
                }
            },
            isShow (isShow) {
                this.isShowAllFields = this.details.op_type === 1 || this.details.op_type === 3
            }
        },
        methods: {
            toggleFields () {
                this.isShowAllFields = !this.isShowAllFields
            },
            getCellValue (property, type) {
                const data = this.details.content[type]
                if (data) {
                    const {
                        bk_property_id: bkPropertyId,
                        bk_property_type: bkPropertyType,
                        bk_property_name: bkPropertyName,
                        option
                    } = property
                    let value = data[bkPropertyId]
                    if (bkPropertyType === 'enum' && Array.isArray(option)) {
                        const targetOption = option.find(({id}) => id === value)
                        value = targetOption ? targetOption.name : ''
                    } else if (bkPropertyType === 'singleasst' || bkPropertyType === 'multiasst') {
                        let asstVal = [];
                        (Array.isArray(value) ? value : []).forEach(({bk_inst_name: bkInstName}) => {
                            if (bkInstName) {
                                asstVal.push(bkInstName)
                            }
                        })
                        value = asstVal.join(',')
                    } else if (bkPropertyType === 'date' || bkPropertyType === 'time') {
                        value = this.$formatTime(value, bkPropertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
                    }
                    return value
                }
                return null
            },
            hasChanged (item) {
                if ([2, 100].includes(this.details['op_type'])) {
                    return item['pre_data'] !== item['cur_data']
                }
                return false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-details-wrapper{
        padding: 32px 50px;
        height: calc(100% - 60px);
    }
    .field-btn{
        margin: 10px 0;
        text-align: right;
        color: #3c96ff;
        cursor: pointer;
        &:hover{
            color: #0082ff;
        }
    }
    .info-group{
        width: 50%;
        display: inline-block;
        white-space: nowrap;
        line-height: 26px;
        font-size: 12px;
        &.op_desc{
            width: 100%;
            .info-value{
                width: 450px;
            }
        }
        .info-label,
        .info-value{
            display: inline-block;
            @include ellipsis;
        }
        .info-label{
            text-align: right;
            width: 100px;
            color: $textColor;
        }
        .info-value{
            padding-left: 4px;
            color: #333948;
            width: 220px;
        }
    }
    .details-data{
        min-height: 100%;
        width: calc(100% + 32px);
        padding: 6px 16px;
        margin: 0 0 0 -16px;
        white-space: normal;
        &.has-changed{
            background-color: #e9faf0;
        }
    }
</style>