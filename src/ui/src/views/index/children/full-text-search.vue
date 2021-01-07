<template>
    <div class="full-text-search-layout">
        <div class="results-wrapper" v-show="!showNoData">
            <div class="results-list">
                <div class="results-item"
                    v-for="(source, index) in searchData"
                    :key="index">
                    <template v-if="source['hitsType'] === 'object'">
                        <div class="results-title" @click="jumpPage(source)">
                            <span v-html="`${modelClassifyName[source['bk_obj_id']]} - ${source.bk_inst_name.toString()}`"></span>
                        </div>
                        <div class="results-desc" v-if="propertyMap[source['bk_obj_id']]" @click="jumpPage(source)">
                            <span class="desc-item" v-html="`${$t('实例ID')}：${source['bk_inst_id']}`"> </span>
                            <template v-for="(property, childIndex) in propertyMap[source['bk_obj_id']]">
                                <span class="desc-item"
                                    :key="childIndex"
                                    v-if="source[property['bk_property_id']]"
                                    v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                                </span>
                            </template>
                        </div>
                    </template>
                    <template v-else-if="source.hitsType === 'host'">
                        <div class="results-title" @click="jumpPage(source)">
                            <span v-html="`${modelClassifyName['host']} - ${source.bk_host_innerip.toString()}`"></span>
                        </div>
                        <div class="results-desc" v-if="propertyMap['host']" @click="jumpPage(source)">
                            <span class="desc-item" v-html="`${$t('主机ID')}：${getHighlightValue(source['bk_host_id'], source, 'bk_host_id')}`"> </span>
                            <template v-for="(property, childIndex) in propertyMap['host']">
                                <span class="desc-item"
                                    v-if="source[property['bk_property_id']]"
                                    :key="childIndex"
                                    v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                                </span>
                            </template>
                        </div>
                    </template>
                    <template v-else-if="source.hitsType === 'biz'">
                        <div class="results-title" @click="jumpPage(source)">
                            <span v-html="`${modelClassifyName['biz']} - ${source.bk_biz_name.toString()}`"></span>
                            <i class="disabled-mark" v-if="source.bk_data_status === 'disabled'">{{$t('已归档')}}</i>
                        </div>
                        <div class="results-desc" v-if="propertyMap['biz']" @click="jumpPage(source)">
                            <template v-for="(property, childIndex) in propertyMap['biz']">
                                <span class="desc-item"
                                    :key="childIndex"
                                    v-if="source[property['bk_property_id']]"
                                    v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                                </span>
                            </template>
                        </div>
                    </template>
                </div>
            </div>
        </div>
        <div class="no-data" v-show="showNoData">
            <img src="../../../assets/images/full-text-search.png" alt="no-data">
            <p>{{$t('搜不到相关内容')}}</p>
        </div>
    </div>
</template>

<script>
    import {
        MENU_RESOURCE_INSTANCE_DETAILS,
        MENU_RESOURCE_BUSINESS_DETAILS,
        MENU_RESOURCE_HOST_DETAILS,
        MENU_RESOURCE_BUSINESS_HISTORY
    } from '@/dictionary/menu-symbol'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            queryData: {
                type: Object,
                default: () => {}
            },
            modelClassify: {
                type: Array,
                default: () => []
            }
        },
        data () {
            return {
                toggleTips: null,
                properties: [],
                showNoData: false,
                searchData: [],
                modelClassifyName: {},
                propertyMap: {}
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('objectModelClassify', ['models', 'getModelById'])
        },
        watch: {
            queryData () {
                if (!this.queryData.total) {
                    this.showNoData = true
                    return
                }
                this.showNoData = false
                this.initResult(this.queryData)
            }
        },
        created () {
            if (!this.queryData.total) {
                this.showNoData = true
                return
            }
            this.showNoData = false
            this.initResult(this.queryData)
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute', 'batchSearchObjectAttribute']),
            async initResult (data) {
                const hitsData = data.hits || []
                const modelData = data.aggregations || []
                this.modelClassify.forEach(model => {
                    this.modelClassifyName[model['bk_obj_id']] = model['bk_obj_name']
                })
                if (modelData.length) {
                    await this.getProperties({
                        $in: modelData.map(model => model.key)
                    })
                }
                this.searchData = hitsData.filter(hits => !['model'].includes(hits.type)).map(hits => {
                    const hit = {
                        ...hits.source,
                        hitsType: hits.type,
                        highlight: hits.highlight
                    }
                    if (hit.hasOwnProperty('bk_obj_id')) {
                        hit['bk_obj_id'] = hit['bk_obj_id'].toString().replace(/\<\/?em\>/g, '')
                    }
                    return hit
                })
            },
            async getProperties (objId) {
                this.propertyMap = await this.batchSearchObjectAttribute({
                    params: {
                        bk_obj_id: objId,
                        bk_supplier_account: this.supplierAccount
                    }
                })
            },
            jumpPage (source) {
                if (source['hitsType'] === 'host') {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_HOST_DETAILS,
                        params: {
                            id: source['bk_host_id'].toString().replace(/(\<\/?em\>)/g, '')
                        },
                        history: true
                    })
                } else if (source['hitsType'] === 'object') {
                    const model = this.getModelById(source['bk_obj_id'])
                    const isPauserd = this.getModelById(source['bk_obj_id'])['bk_ispaused']
                    if (model['bk_classification_id'] === 'bk_biz_topo') {
                        this.$bkMessage({
                            message: this.$t('主线模型无法查看'),
                            theme: 'warning'
                        })
                        return
                    } else if (isPauserd) {
                        this.$bkMessage({
                            message: this.$t('该模型已停用'),
                            theme: 'warning'
                        })
                        return
                    }
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_INSTANCE_DETAILS,
                        params: {
                            objId: source['bk_obj_id'],
                            instId: source['bk_inst_id'].toString().replace(/(\<\/?em\>)/g, '')
                        },
                        history: true
                    })
                } else if (source['hitsType'] === 'biz') {
                    const name = source.bk_data_status === 'disabled' ? MENU_RESOURCE_BUSINESS_HISTORY : MENU_RESOURCE_BUSINESS_DETAILS
                    this.$routerActions.redirect({
                        name: name,
                        params: {
                            bizId: source.bk_biz_id,
                            bizName: source['bk_biz_name'].toString().replace(/(\<\/?em\>)/g, '')
                        },
                        history: true
                    })
                }
            },
            getShowPropertyText (property, source, thisProperty) {
                let propertyValue = this.$tools.getPropertyText(property, source)

                if (!Object.keys(source.highlight).includes(thisProperty)) {
                    return propertyValue || '--'
                }

                // 对highlight属性值做高亮标签处理
                propertyValue = this.getHighlightValue(propertyValue, source, thisProperty)
                return propertyValue || '--'
            },
            getHighlightValue (value, source, thisProperty) {
                const highlightValue = source.highlight[thisProperty]
                if (!highlightValue) {
                    return value
                }
                let keyword = Array.isArray(highlightValue) ? highlightValue[0] : highlightValue
                keyword = keyword.match(/<em>(.+?)<\/em>/)[1]
                const reg = new RegExp(`(${keyword})`, 'g')
                return String(value).replace(reg, '<em>$1</em>')
            }
        }
    }
</script>

<style lang="scss">
    .full-text-search-layout {
        position: relative;
        .results-wrapper {
            width: 90%;
            margin: 0 auto;
            .results-list {
                padding-top: 14px;
                color: $cmdbTextColor;
                .results-item {
                    width: 65%;
                    padding-bottom: 35px;
                    color: #63656e;
                    em {
                        color: #3a84ff !important;
                        font-style: normal !important;
                        word-break: break-all;
                    }
                    .results-title {
                        display: inline-block;
                        font-size: 18px;
                        font-weight: bold;
                        margin-bottom: 4px;
                        cursor: pointer;
                        &:hover {
                            span {
                                color: #3a84ff;
                                text-decoration: underline;
                            }
                        }
                        .disabled-mark {
                            height: 18px;
                            line-height: 16px;
                            padding: 0 4px;
                            font-style: normal;
                            font-size: 12px;
                            color: #979BA5;
                            border: 1px solid #C4C6CC;
                            background-color: #FAFBFD;
                            border-radius: 2px;
                            margin-left: 4px;
                            text-decoration: none;
                        }
                    }
                    .results-desc {
                        font-size: 14px;
                        .desc-item {
                            display: inline-block;
                            padding-right: 16px;
                            padding-bottom: 6px;
                        }
                        &:hover {
                            color: #313238;
                            cursor: pointer;
                        }
                    }
                }
            }
        }
        .no-data {
            width: 90%;
            margin: 0 auto;
            padding-top: 240px;
            text-align: center;
            color: #63656E;
            font-size: 16px;
            img {
                width: 104px;
            }
        }
    }
</style>
