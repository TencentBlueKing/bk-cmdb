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
            <div ref="historyCompare" class="history-compare" @scroll="setHeader" v-bkloading="{isLoading: loadingAttribute}">
                <table ref="compareTableHeader" class="compare-table-header">
                    <thead>
                        <tr class="compare-header-row">
                            <td class="compare-header-cell" width="130"></td>
                            <td class="compare-header-cell" width="280">{{$t("OperationAudit['变更前']")}}</td>
                            <td class="compare-header-cell" width="280">{{$t("OperationAudit['变更后']")}}</td>
                        </tr>
                    </thead>
                </table>
                <table class="compare-table-body">
                    <tbody>
                        <tr :class="['compare-body-row', {changed: isChanged(header)}]" v-for="(header, index) in compareHeader" :key="index">
                            <td class="compare-body-cell header" width="130">{{header['bk_property_name']}}</td>
                            <td class="compare-body-cell pre" width="280">{{getCompareBodyCell(header, 'pre_data')}}</td>
                            <td class="compare-body-cell cur" width="280">{{getCompareBodyCell(header, 'cur_data')}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </template>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'
    export default {
        props: {
            details: Object
        },
        data () {
            return {
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
                    key: 'bk_host_innerip'
                }, {
                    label: 'OperationAudit[\'操作类型\']',
                    key: 'op_type',
                    optionKey: 'opType'
                }, {
                    label: 'OperationAudit[\'操作对象\']',
                    key: 'op_target'
                }, {
                    label: 'OperationAudit[\'操作时间\']',
                    key: 'op_time'
                }, {
                    label: 'OperationAudit[\'描述\']',
                    key: 'op_desc'
                }]
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
                    3: this.$t("Common['删除']")
                }
                return {
                    biz,
                    opType
                }
            },
            compareHeader () {
                return this.attribute[this.objId] || []
            }
        },
        watch: {
            async objId (objId) {
                if (objId && !this.attribute.hasOwnProperty(objId)) {
                    this.loadingAttribute = true
                    await this.$store.dispatch('object/getAttribute', objId)
                    this.loadingAttribute = false
                } else {
                    this.loadingAttribute = false
                }
            }
        },
        methods: {
            getCompareBodyCell (header, type) {
                let data = this.details.content[type]
                if (data) {
                    let cellText = data[header['bk_property_id']]
                    let bkPropertyType = header['bk_property_type']
                    if (Array.isArray(cellText)) {
                        cellText = cellText.map(({bk_inst_name: bkInstName}) => {
                            return bkInstName
                        }).join(',')
                    } else if (bkPropertyType === 'enum' && Array.isArray(header['options'])) {
                        let enumOption = header['option'].find(({id}) => id === cellText)
                        cellText = enumOption ? enumOption['name'] : ''
                    } else if (bkPropertyType === 'date' || bkPropertyType === 'time') {
                        cellText = this.$formatTime(cellText, bkPropertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
                    }
                    return !cellText ? null : cellText
                }
                return null
            },
            isChanged (header) {
                return this.getCompareBodyCell(header, 'pre_data') !== this.getCompareBodyCell(header, 'cur_data')
            },
            setHeader () {
                this.$refs.compareTableHeader.style.top = this.$refs.historyCompare.scrollTop + 'px'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-details-wrapper{
        padding: 32px 50px;
        height: calc(100% - 60px);
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
    .history-compare{
        position: relative;
        margin-top: 32px;
        padding-top: 43px;
        max-height: calc(100% - 136px - 32px);
        overflow-y: auto;
        overflow-x: hidden;
        @include scrollbar;
    }
    .compare-table-header{
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        .compare-header-row{
            height: 42px;
            background-color: #fafbfd;
            .compare-header-cell{
                border: 1px solid #dde4eb;
                padding: 0 20px;
            }
        }
    }
    .compare-table-body{
        width: 100%;
        .compare-body-row{
            &.changed{
                .compare-body-cell.pre,
                .compare-body-cell.cur{
                    background-color: #e9faf0;
                }
            }
            .compare-body-cell{
                line-height: 26px;
                padding: 8px 20px;
                border: 1px solid #dde4eb;
                border-top: none;
            }
        }
    }
</style>