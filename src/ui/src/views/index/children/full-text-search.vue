<template>
    <div class="full-text-search-layout">
        <div class="results-wrapper" v-show="!showNoData">
            <div class="results-list">
                <div class="results-item"
                    v-for="(source, index) in searchData"
                    :key="index"
                    v-if="!['model'].includes(source['hitsType'])">
                    <template v-if="source['hitsType'] === 'object'">
                        <div class="results-title"
                            v-html="`${modelClassifyName[source['bk_obj_id']]} - ${source.bk_inst_name.toString()}`"
                            @click="jumpPage(source)"></div>
                        <div class="results-desc" v-if="propertyMap[source['bk_obj_id']]" @click="jumpPage(source)">
                            <span class="desc-item" v-html="`${$t('实例ID')}：${source['bk_inst_id']}`"> </span>
                            <span class="desc-item"
                                v-for="(property, childIndex) in propertyMap[source['bk_obj_id']]"
                                :key="childIndex"
                                v-if="source[property['bk_property_id']]"
                                v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                            </span>
                        </div>
                    </template>
                    <template v-else-if="source.hitsType === 'host'">
                        <div class="results-title"
                            v-html="`${modelClassifyName['host']} - ${source.bk_host_innerip.toString()}`"
                            @click="jumpPage(source)"></div>
                        <div class="results-desc" v-if="propertyMap['host']" @click="jumpPage(source)">
                            <span class="desc-item" v-html="`${$t('主机ID')}${source['bk_host_id']}`"> </span>
                            <span class="desc-item"
                                v-for="(property, childIndex) in propertyMap['host']"
                                :key="childIndex"
                                v-if="source[property['bk_property_id']]"
                                v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                            </span>
                        </div>
                    </template>
                    <template v-else-if="source.hitsType === 'biz'">
                        <div class="results-title"
                            v-html="`${modelClassifyName['biz']} - ${source.bk_biz_name.toString()}`"
                            @click="jumpPage(source)"></div>
                        <div class="results-desc" v-if="propertyMap['biz']" @click="jumpPage(source)">
                            <span class="desc-item"
                                v-for="(property, childIndex) in propertyMap['biz']"
                                :key="childIndex"
                                v-if="source[property['bk_property_id']]"
                                v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                            </span>
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
    import { MENU_RESOURCE_INSTANCE, MENU_RESOURCE_BUSINESS, MENU_RESOURCE_HOST_DETAILS } from '@/dictionary/menu-symbol'
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
                await this.processArray(modelData)
                this.searchData = hitsData.map(hits => {
                    const hit = {
                        ...hits.source,
                        ...hits.highlight,
                        hitsType: hits.type
                    }
                    if (hit.hasOwnProperty('bk_obj_id')) {
                        hit['bk_obj_id'] = hit['bk_obj_id'].toString().replace(/\<\/?em\>/g, '')
                    }
                    return hit
                })
            },
            isPublicModel (objId) {
                const model = this.models.find(model => model['bk_obj_id'] === objId) || {}
                return !this.$tools.getMetadataBiz(model)
            },
            async getPublicModelProperties (objId) {
                this.propertyMap = await this.batchSearchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: objId,
                        bk_supplier_account: this.supplierAccount
                    }, { inject: false }),
                    config: {
                        requestId: `post_batchSearchObjectAttribute_${objId['$in'].join('_')}`,
                        requestGroup: objId['$in'].map(id => `post_searchObjectAttribute_${id}`)
                    }
                })
            },
            async getProperties (objId) {
                const properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: objId,
                        bk_supplier_account: this.supplierAccount
                    }, { inject: true }),
                    config: {
                        requestId: `post_searchObjectAttribute_${objId}`,
                        fromCache: false
                    }
                })
                this.$set(this.propertyMap, objId, properties)
            },
            async processArray (data) {
                this.propertyMap = {}
                const publicObj = data.filter(aggregation => this.isPublicModel(aggregation.key)).map(model => model.key)
                const privateObj = data.filter(aggregation => !this.isPublicModel(aggregation.key)).map(model => model.key)
                if (publicObj.length) {
                    const objId = {
                        '$in': publicObj
                    }
                    await this.getPublicModelProperties(objId)
                }
                if (privateObj.length) {
                    for (let i = 0; i < privateObj.length; i++) {
                        await this.getProperties(privateObj[i])
                    }
                }
            },
            jumpPage (source) {
                if (source['hitsType'] === 'host') {
                    this.$router.push({
                        name: MENU_RESOURCE_HOST_DETAILS,
                        params: {
                            id: source['bk_host_id']
                        }
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
                    this.$router.push({
                        name: MENU_RESOURCE_INSTANCE,
                        params: {
                            objId: source['bk_obj_id'],
                            instId: source['bk_inst_id'].toString().replace(/(\<\/?em\>)/g, '')
                        }
                    })
                } else if (source['hitsType'] === 'biz') {
                    this.$router.push({
                        name: MENU_RESOURCE_BUSINESS,
                        params: {
                            bizName: source['bk_biz_name'].toString().replace(/(\<\/?em\>)/g, '')
                        }
                    })
                }
            },
            getShowPropertyText (property, source, thisProperty) {
                const cloneSource = this.$tools.clone(source)
                const reg = /\<em\>.+\<\/em\>/
                let propertyValue = cloneSource[thisProperty].toString()
                if (propertyValue === '[object Object]') {
                    propertyValue = cloneSource[thisProperty]
                }
                const isHeightLight = reg.test(propertyValue)
                cloneSource[thisProperty] = isHeightLight ? propertyValue.replace(/(\<\/?em\>)/g, '') : propertyValue
                const flatternedText = this.$tools.getPropertyText(property, cloneSource)
                return isHeightLight ? `<em>${flatternedText}</em>` : flatternedText
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
                            color: #3a84ff;
                            text-decoration: underline;
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
