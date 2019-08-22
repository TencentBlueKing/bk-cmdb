<template>
    <div class="history-details-wrapper">
        <template v-if="details">
            <div class="history-info">
                <div :class="['info-group', info.key]" v-for="(info, index) in informations" :key="index">
                    <label class="info-label">{{info.label}}：</label>
                    <span class="info-value">
                        <template v-if="info.key === 'op_time'">
                            {{$tools.formatTime(details[info.key])}}
                        </template>
                        <template v-else>
                            {{info.hasOwnProperty('optionKey') ? options[info.optionKey][details[info.key]] : details[info.key]}}
                        </template>
                    </span>
                </div>
            </div>
            <bk-table
                v-bkloading="{ isLoading: $loading(`post_searchObjectAttribute_${objId}`) }"
                :data="displayList"
                :width="width || 700"
                :max-height="$APP.height - 300"
                :height="height"
                :row-border="true"
                :col-border="true"
                :cell-style="getCellStyle">
                <bk-table-column prop="bk_property_name"></bk-table-column>
                <bk-table-column v-if="details.op_type !== 1"
                    prop="pre_data"
                    :label="$t('变更前')">
                    <template slot-scope="{ row }">
                        <span v-html="row.pre_data"></span>
                    </template>
                </bk-table-column>
                <bk-table-column v-if="details.op_type !== 3"
                    prop="cur_data"
                    :label="$t('变更后')">
                    <template slot-scope="{ row }" v-html="row.cur_data">
                        <span v-html="row.cur_data"></span>
                    </template>
                </bk-table-column>
            </bk-table>
            <p class="field-btn" @click="toggleFields" v-if="isShowToggle && details.op_type !== 1 && details.op_type !== 3">
                {{isShowAllFields ? $t('收起') : $t('展开')}}
            </p>
        </template>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            details: Object,
            height: Number,
            width: Number,
            isShow: Boolean
        },
        data () {
            return {
                attribute: [],
                isShowAllFields: false,
                informations: [{
                    label: this.$t('操作账号'),
                    key: 'operator'
                }, {
                    label: this.$t('所属业务'),
                    key: 'bk_biz_id',
                    optionKey: 'biz'
                }, {
                    label: this.$t('IP'),
                    key: 'ext_key'
                }, {
                    label: this.$t('类型'),
                    key: 'op_type',
                    optionKey: 'opType'
                }, {
                    label: this.$t('对象'),
                    key: 'op_target'
                }, {
                    label: this.$t('操作时间'),
                    key: 'op_time'
                }, {
                    label: this.$t('描述'),
                    key: 'op_desc'
                }],
                colWidth: [130, 280, 280]
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['authorizedBusiness']),
            objId () {
                return this.details ? this.details['op_target'] : null
            },
            options () {
                const biz = {}
                this.authorizedBusiness.forEach(({ bk_biz_id: bkBizId, bk_biz_name: bkBizName }) => {
                    biz[bkBizId] = bkBizName
                })
                const opType = {
                    1: this.$t('新增'),
                    2: this.$t('修改'),
                    3: this.$t('删除'),
                    100: this.$t('关系变更')
                }
                return {
                    biz,
                    opType
                }
            },
            tableList () {
                const list = []
                const attribute = (this.attribute || []).filter(({ bk_isapi: bkIsapi }) => !bkIsapi)
                if (this.details['op_type'] !== 100) {
                    attribute.forEach(property => {
                        const preData = this.getCellValue(property, 'pre_data')
                        const curData = this.getCellValue(property, 'cur_data')
                        list.push({
                            'bk_property_name': property['bk_property_name'],
                            'pre_data': preData,
                            'cur_data': curData
                        })
                    })
                } else {
                    const content = this.details.content
                    const preBizId = content['pre_data']['bk_biz_id']
                    const curBizId = content['cur_data']['bk_biz_id']
                    const preModule = content['pre_data']['module'] || []
                    const curModule = content['cur_data']['module'] || []
                    const pre = []
                    const cur = []
                    preModule.forEach(module => {
                        pre.push(this.getTopoPath(preBizId, module))
                    })
                    curModule.forEach(module => {
                        cur.push(this.getTopoPath(curBizId, module))
                    })
                    const preData = pre.join('<br>')
                    const curData = cur.join('<br>')
                    list.push({
                        'bk_property_name': this.$t('关联关系'),
                        'pre_data': preData,
                        'cur_data': curData
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
        watch: {
            objId (objId) {
                if (objId) {
                    this.getObjAttribute()
                }
            }
        },
        created () {
            this.getObjAttribute()
        },
        mounted () {
            this.isShowAllFields = this.details.op_type === 1 || this.details.op_type === 3
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            toggleFields () {
                this.isShowAllFields = !this.isShowAllFields
            },
            async getObjAttribute () {
                const res = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: this.objId
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_${this.objId}`,
                        fromCache: true
                    }
                })
                this.attribute = res
            },
            getCellValue (property, type) {
                const data = this.details.content[type]
                if (data) {
                    const {
                        bk_property_id: bkPropertyId,
                        bk_property_type: bkPropertyType,
                        // bk_property_name: bkPropertyName,
                        option
                    } = property
                    let value = data[bkPropertyId]
                    if (bkPropertyType === 'enum' && Array.isArray(option)) {
                        const targetOption = option.find(({ id }) => id === value)
                        value = targetOption ? targetOption.name : ''
                    } else if (bkPropertyType === 'date' || bkPropertyType === 'time') {
                        value = this.$tools.formatTime(value, bkPropertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
                    }
                    return value === '' ? null : value
                }
                return null
            },
            hasChanged (item) {
                if ([2, 100].includes(this.details['op_type'])) {
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
            getTopoPath (bizId, module) {
                const path = [this.options.biz[bizId] || `业务ID：${bizId}`]
                const set = ((module.set || [])[0] || {}).ref_name
                if (set) {
                    path.push(set)
                }
                if (module.ref_name) {
                    path.push(module.ref_name)
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
