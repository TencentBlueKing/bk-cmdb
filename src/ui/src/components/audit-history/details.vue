<template>
    <div class="history-details-wrapper">
        <template v-if="details">
            <div class="history-info">
                <div :class="['info-group', item.key]" v-for="(item, index) in informations" :key="index">
                    <label class="info-label">{{item.label}}：</label>
                    <span class="info-value">{{displayInfo[item.key]}}</span>
                </div>
            </div>
            <bk-table
                :data="displayList"
                :width="width || 700"
                :max-height="$APP.height - 300"
                :height="height"
                :row-border="true"
                :col-border="true"
                :cell-style="getCellStyle">
                <bk-table-column prop="bk_property_name"></bk-table-column>
                <bk-table-column v-if="!['create'].includes(details.action)"
                    prop="pre_data"
                    :label="$t('变更前')">
                    <template slot-scope="{ row }">
                        <span v-html="row.pre_data"></span>
                    </template>
                </bk-table-column>
                <bk-table-column v-if="!['delete'].includes(details.action)"
                    prop="cur_data"
                    :label="$t('变更后')">
                    <template slot-scope="{ row }" v-html="row.cur_data">
                        <span v-html="row.cur_data"></span>
                    </template>
                </bk-table-column>
            </bk-table>
            <p class="field-btn" @click="toggleFields" v-if="isShowToggle && !['create', 'delete'].includes(details.action)">
                {{isShowAllFields ? $t('收起') : $t('展开')}}
            </p>
        </template>
    </div>
</template>

<script>
    import { mapState, mapGetters } from 'vuex'
    export default {
        props: {
            details: Object,
            height: Number,
            width: Number,
            isShow: Boolean,
            showBusiness: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                isShowAllFields: false,
                informations: [{
                    label: this.$t('动作'),
                    key: 'action'
                }, {
                    label: this.$t('功能板块'),
                    key: 'resourceType'
                }, {
                    label: this.$t('操作实例'),
                    key: 'instance'
                }, {
                    label: this.$t('操作描述'),
                    key: 'desc'
                }, {
                    label: this.$t('操作时间'),
                    key: 'time'
                }, {
                    label: this.$t('操作账号'),
                    key: 'user'
                }],
                colWidth: [130, 280, 280],
                hostOperations: ['assign_host', 'unassign_host', 'transfer_host_module']
            }
        },
        computed: {
            ...mapState('operationAudit', ['funcActions']),
            ...mapGetters('objectBiz', ['authorizedBusiness']),
            operationDetail () {
                return this.details.operation_detail
            },
            bizSet () {
                const biz = {}
                this.authorizedBusiness.forEach(({ bk_biz_id: bkBizId, bk_biz_name: bkBizName }) => {
                    biz[bkBizId] = bkBizName
                })
                return biz
            },
            actionList () {
                const funcMap = []
                Object.keys(this.funcActions).forEach(key => {
                    funcMap.push(...this.funcActions[key])
                })
                return funcMap
            },
            funcModules () {
                const modules = {}
                this.actionList.forEach(item => {
                    modules[item.id] = this.$t(item.name)
                })
                return modules
            },
            actionSet () {
                const actionSet = {}
                const operations = this.actionList.reduce((acc, item) => acc.concat(item.operations), [])
                operations.forEach(action => {
                    actionSet[action.id] = this.$t(action.name)
                })
                return actionSet
            },
            displayInfo () {
                const info = {}
                info.action = this.getResourceAction()
                info.resourceType = this.getResourceType()
                info.instance = this.getResourceName()
                info.desc = `${info.action}"${info.instance}"`
                info.time = this.$tools.formatTime(this.details.operation_time)
                info.user = this.details.user
                if (this.showBusiness) {
                    const basicDetail = this.operationDetail.basic_detail
                    info.bizName = basicDetail && basicDetail.bk_biz_name
                }
                return info
            },
            tableList () {
                const list = []
                const type = this.details.resource_type
                if (type === 'instance_association') {
                    const keys = Object.keys(this.operationDetail.basic_asst_detail || {})
                    const attributes = [...keys, 'src_instance_id', 'src_instance_name', 'target_instance_id', 'target_instance_name']
                    const details = {
                        ...this.operationDetail,
                        ...this.operationDetail.basic_asst_detail
                    }
                    attributes.forEach(name => {
                        const data = details[name]
                        list.push({
                            'bk_property_name': name,
                            'pre_data': data,
                            'cur_data': data
                        })
                    })
                } else if (this.hostOperations.includes(this.details.action)) {
                    const content = this.operationDetail
                    const preBizId = content['pre_data']['bk_biz_id']
                    const curBizId = content['cur_data']['bk_biz_id']
                    const preSet = content['pre_data']['set'] || []
                    const curSet = content['cur_data']['set'] || []
                    const pre = []
                    const cur = []
                    preSet.forEach(set => {
                        pre.push(this.getTopoPath(preBizId, set))
                    })
                    curSet.forEach(set => {
                        cur.push(this.getTopoPath(curBizId, set))
                    })
                    const preData = pre.join('<br>')
                    const curData = cur.join('<br>')
                    list.push({
                        'bk_property_name': this.$t('关联关系'),
                        'pre_data': preData,
                        'cur_data': curData
                    })
                } else {
                    const attribute = this.operationDetail.basic_detail.details.properties
                    attribute.forEach(property => {
                        const preData = this.getCellValue(property, 'pre_data')
                        const curData = this.getCellValue(property, 'cur_data')
                        list.push({
                            'bk_property_name': property.bk_property_name,
                            'pre_data': preData,
                            'cur_data': curData
                        })
                    })
                }
                return list
            },
            changedList () {
                return this.tableList.filter(item => item.pre_data !== item.cur_data)
            },
            displayList () {
                return this.isShowAllFields ? this.tableList : this.changedList.length ? this.changedList : this.tableList
            },
            isShowToggle () {
                return this.tableList.length !== this.changedList.length && this.changedList.length > 0
            }
        },
        mounted () {
            this.isShowAllFields = ['create', 'delete'].includes(this.details.action)
            if (this.showBusiness) {
                this.informations.push({
                    label: this.$t('所属业务'),
                    key: 'bizName'
                })
            }
        },
        methods: {
            getResourceType () {
                if (this.details.label === null) {
                    const type = this.details.resource_type
                    if (type === 'model_instance') {
                        const model = this.$store.getters['objectModelClassify/getModelById'](this.operationDetail.bk_obj_id) || {}
                        return model.bk_obj_name || '--'
                    }
                    return this.funcModules[type] || '--'
                }
                const key = Object.keys(this.details.label)
                return this.funcModules[key[0]] || '--'
            },
            getResourceName () {
                const data = this.operationDetail
                if (this.hostOperations.includes(this.details.action)) {
                    return data.bk_host_innerip || '--'
                }
                if (['instance_association'].includes(this.details.resource_type)) {
                    return data.target_instance_name || '--'
                }
                return (data.basic_detail && data.basic_detail.resource_name) || '--'
            },
            getResourceAction () {
                const data = this.details
                if (data.label) {
                    const label = Object.keys(data.label)[0]
                    return this.actionSet[`${data.resource_type}-${data.action}-${label}`]
                }
                return this.actionSet[`${data.resource_type}-${data.action}`]
            },
            toggleFields () {
                this.isShowAllFields = !this.isShowAllFields
            },
            getCellValue (property, type) {
                const data = this.operationDetail.basic_detail.details[type]
                let value
                if (data) {
                    value = data[property.bk_property_id]
                }
                return [undefined, null, ''].includes(value) ? '--' : value
            },
            hasChanged (item) {
                const action = ['update', 'archive', 'recover'].concat(this.hostOperations)
                if (action.includes(this.details['action'])) {
                    return item['pre_data'] !== item['cur_data']
                }
                return false
            },
            getCellStyle ({ row, columnIndex }) {
                if (columnIndex > 0 && this.hasChanged(row)) {
                    return {
                        backgroundColor: '#e9faf0'
                    }
                }
                return {}
            },
            getTopoPath (bizId, set) {
                const path = [this.bizSet[bizId] || `业务ID：${bizId}`]
                const module = ((set.module || [])[0] || {}).bk_module_name
                if (set.bk_set_name) {
                    path.push(set.bk_set_name)
                }
                if (module) {
                    path.push(module)
                }
                return path.join('→')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-details-wrapper{
        padding: 32px 50px;
        height: 100%;
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
        }
        .info-value{
            padding-left: 4px;
            color: #333948;
            width: 220px;
        }
    }
    .field-btn{
        font-size: 14px;
        margin: 10px 0;
        text-align: right;
        color: #3c96ff;
        cursor: pointer;
        &:hover{
            color: #0082ff;
        }
    }
</style>
