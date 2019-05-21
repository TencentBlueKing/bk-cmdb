<template>
    <div class="full-text-search-layout">
        <div class="full-text-sticky-layout" ref="topSticky">
            <div class="search-bar">
                <input id="fullTextSearch"
                    ref="searchInput"
                    autocomplete="off"
                    class="search-keywords"
                    type="text"
                    maxlength="40"
                    :placeholder="$t('Common[\'请输入搜索内容\']')"
                    v-model.trim="query.queryString"
                    @keypress.enter="handleSearch">
                <label class="bk-icon icon-search" for="fullTextSearch"></label>
            </div>
            <div class="classify" v-show="hasData">
                <span class="classify-item"
                    :class="['classify-item', { 'classify-active': -1 === currentClassify }]"
                    @click="toggleClassify(-1)">
                    {{$t("Index['全部结果']")}}（{{allSearchCount}}）
                </span>
                <span class="classify-item"
                    :class="['classify-item', { 'classify-active': index === currentClassify }]"
                    v-for="(model, index) in modelClassify"
                    :key="index"
                    @click="toggleClassify(index, model['bk_obj_id'])">
                    {{model['bk_obj_name']}}（{{model['count']}}）
                </span>
            </div>
        </div>
        <div class="results-wrapper" ref="resultsWrapper"
            :class="{ 'searching': searching }"
            v-bkloading="{ 'isLoading': searching }">
            <div v-show="hasData">
                <div class="results-list" ref="resultsList">
                    <div class="results-item"
                        v-for="(source, index) in searchData"
                        :key="index"
                        v-if="!['model'].includes(source['hitsType'])">
                        <template v-if="source['hitsType'] === 'object'">
                            <div class="results-title"
                                v-html="`${modelClassifyName[source['bk_obj_id']]} - ${source.bk_inst_name.toString()}`"
                                @click="jumpPage(source)"></div>
                            <div class="results-desc" v-if="propertyMap[source['bk_obj_id']]">
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
                            <div class="results-desc" v-if="propertyMap['host']">
                                <span class="desc-item"
                                    v-for="(property, childIndex) in propertyMap['host']"
                                    :key="childIndex"
                                    v-if="source[property['bk_property_id']]"
                                    v-html="`${property['bk_property_name']}：${getShowPropertyText(property, source, property['bk_property_id'])}`">
                                </span>
                            </div>
                        </template>
                    </div>
                </div>
                <full-text-pagination
                    :pagination="pagination"
                    @limit-change="handleLimitChange"
                    @page-change="handlePageChange">
                </full-text-pagination>
            </div>
        </div>
        <div class="no-data" v-show="!hasData && showNoData && !searching">
            <img src="../../../assets/images/full-text-search.png" alt="no-data">
            <p>{{$t("Index['搜不到相关内容']")}}</p>
        </div>
    </div>
</template>

<script>
    import {
        addMainScrollListener,
        removeMainScrollListener
    } from '@/utils/main-scroller'
    import fullTextPagination from './full-text-pagination.vue'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            fullTextPagination
        },
        data () {
            return {
                scrollTop: null,
                currentClassify: -1,
                requestId: 'fullTextSearch',
                query: {
                    queryString: '',
                    objId: ''
                },
                properties: [],
                hasData: false,
                showNoData: false,
                searching: false,
                pagination: {
                    limit: 10,
                    start: 0,
                    current: 1,
                    total: 0
                },
                allSearchCount: null,
                searchData: [],
                modelClassify: [],
                modelClassifyName: {},
                propertyMap: {},
                debounceTimer: null
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('objectModelClassify', ['models', 'getModelById']),
            params () {
                const notZhCn = this.query.queryString.replace(/\w\.?/g, '').length === 0
                return {
                    page: {
                        start: this.pagination.start,
                        limit: this.pagination.limit
                    },
                    bk_obj_id: this.query.objId,
                    bk_biz_id: this.bizId ? this.bizId.toString() : '',
                    query_string: notZhCn ? `*${this.query.queryString}*` : this.query.queryString,
                    filter: ['model']
                }
            },
            hash () {
                return decodeURIComponent(window.location.hash)
            }
        },
        watch: {
            isShowSizeSetting (isShowSizeSetting) {
                if (isShowSizeSetting) {
                    this.calcSizePosition()
                } else {
                    this.sizeTransitionReady = false
                }
            },
            'query.queryString' (queryString) {
                this.showNoData = false
                this.hasData = false
                window.location.hash = this.hash.substring(0, this.hash.search(/=/) + 1) + this.query.queryString
                this.query.objId = ''
                if (queryString) {
                    this.propertyMap = {}
                    this.pagination.start = 0
                    this.pagination.current = 1
                    this.handleSearch(600)
                }
            }
        },
        created () {
            this.$store.commit('setHeaderStatus', {
                back: true
            })
            this.$store.commit('setHeaderTitle', this.$t('Common["搜索结果"]'))
            this.query.queryString = this.$route.query.keywords
        },
        mounted () {
            this.initScrollListener(this.$refs.topSticky)
            this.$nextTick(() => {
                this.$refs.searchInput.focus()
            })
        },
        destroyed () {
            removeMainScrollListener(this.scrollTop)
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute', 'batchSearchObjectAttribute']),
            initScrollListener (domRefs) {
                this.scrollTop = event => {
                    const target = event.target
                    domRefs.style.top = target.scrollTop + 'px'
                }
                this.$refs.resultsWrapper.style.marginTop = domRefs.offsetHeight + 'px'
                addMainScrollListener(this.scrollTop)
            },
            isPublicModel (objId) {
                const model = this.models.find(model => model['bk_obj_id'] === objId) || {}
                return !this.$tools.getMetadataBiz(model)
            },
            toggleClassify (index, model) {
                if (index === this.currentClassify) return
                this.currentClassify = index
                this.query.objId = model
                this.pagination.start = 0
                this.pagination.current = 1
                this.handleSearch()
                // Object.keys(this.modelClassifyName).forEach(key => {
                //     this.modelClassifyName[key] = this.modelClassifyName[key].replace(/(\<\/?em\>)/g, '')
                // })
                // model && (this.modelClassifyName[model] = `<em>${this.modelClassifyName[model]}</em>`)
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
            handleSearch (wait = 0) {
                const debounceTimer = this.debounceTimer
                if (debounceTimer) {
                    clearTimeout(debounceTimer)
                    this.debounceTimer = null
                }
                this.debounceTimer = setTimeout(() => {
                    if (!this.query.queryString) return
                    this.searching = true
                    this.$store.dispatch('fullTextSearch/search', {
                        params: this.params,
                        config: {
                            requestId: this.requestId,
                            cancelPrevious: true
                        }
                    }).then(data => {
                        this.searching = false
                        this.pagination.total = data.total
                        const hitsData = data.hits
                        const modelData = data.aggregations
                        if (hitsData && modelData) {
                            if (!this.query.objId) {
                                this.currentClassify = -1
                                this.modelClassify = modelData.map(item => {
                                    return {
                                        ...this.getModelById(item.key),
                                        count: item.count
                                    }
                                })
                                this.modelClassify.forEach(model => {
                                    this.modelClassifyName[model['bk_obj_id']] = model['bk_obj_name']
                                })
                                this.processArray(modelData)
                                this.allSearchCount = data.total
                                this.hasData = data.total
                            }
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
                        }
                        this.$nextTick(() => {
                            this.initScrollListener(this.$refs.topSticky)
                        })
                    }).finally(() => {
                        this.showNoData = true
                    })
                }, wait)
            },
            async processArray (data) {
                const publicObj = data.filter(aggregation => this.isPublicModel(aggregation.key)).map(model => model.key)
                const privateObj = data.filter(aggregation => !this.isPublicModel(aggregation.key)).map(model => model.key)
                const getPublicAttribute = publicObj.filter(model => !this.propertyMap[model])
                const getPrivateAttribute = privateObj.filter(model => !this.propertyMap[model])
                if (getPublicAttribute.length) {
                    const objId = {
                        '$in': getPublicAttribute
                    }
                    await this.getPublicModelProperties(objId)
                }
                if (getPrivateAttribute.length) {
                    for (let i = 0; i < getPrivateAttribute.length; i++) {
                        await this.getProperties(getPrivateAttribute[i])
                    }
                }
            },
            jumpPage (source) {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                if (source['hitsType'] === 'host') {
                    this.$router.push({
                        name: 'resource',
                        query: {
                            ip: source['bk_host_innerip'].toString().replace(/(\<\/?em\>)/g, ''),
                            outer: false,
                            inner: true,
                            exact: 1,
                            assigned: true
                        }
                    })
                } else if (source['hitsType'] === 'object') {
                    const isPauserd = this.getModelById(source['bk_obj_id'])['bk_ispaused']
                    if (isPauserd) {
                        this.$bkMessage({
                            message: this.$t("Index['该模型已停用']"),
                            theme: 'warning'
                        })
                        return
                    }
                    this.$router.push({
                        name: 'generalModel',
                        params: {
                            objId: source['bk_obj_id'],
                            source: source
                        }
                    })
                }
            },
            getShowPropertyText (property, source, thisProperty) {
                const cloneSource = this.$tools.clone(source)
                const reg = /^\<em\>.+\<\/em\>$/
                const propertyValue = cloneSource[thisProperty].toString()
                const isHeightLight = reg.test(propertyValue)
                if (isHeightLight) {
                    cloneSource[thisProperty] = propertyValue.replace(/(\<\/?em\>)/g, '')
                }
                const flatternedText = this.$tools.getPropertyText(property, cloneSource)
                return isHeightLight ? `<em>${flatternedText}</em>` : flatternedText
            },
            handleLimitChange (limit) {
                if (this.pagination.limit === limit) return
                this.pagination.limit = limit
                this.handlePageChange(1, 'limitChange')
            },
            handlePageChange (current, type) {
                if (!type && this.pagination.current === current) return
                this.pagination.start = (current - 1) * this.pagination.limit
                this.pagination.current = current
                this.handleSearch()
            }
        }
    }
</script>

<style lang="scss">
    .full-text-search-layout {
        position: relative;
        overflow: auto;
        padding: 0 0 20px !important;
        .full-text-sticky-layout {
            width: 90%;
            margin: 0 auto;
            position: absolute;
            top: 0;
            left: 50%;
            background: #ffffff;
            transform: translateX(-50%);
            z-index: 10;
        }
        .search-bar {
            position: relative;
            width: 100%;
            max-width: 640px;
            line-height: 40px;
            margin: 50px auto 38px;
            .search-keywords {
                width: 100%;
                height: 42px;
                border-radius: 2px;
                border: 1px solid #c3cdd7;
                padding: 0 56px 0 16px;
                font-size: 14px;
            }
            .icon-search {
                position: absolute;
                right: 20px;
                top: 13px;
                font-size: 18px;
            }
        }
        .classify {
            color: $cmdbTextColor;
            font-size: 14px;
            border-bottom: 1px solid #dde4eb;
            .classify-item {
                display: inline-block;
                margin-right: 20px;
                margin-bottom: 12px;
                cursor: pointer;
                &.classify-active {
                    color: #3a84ff;
                    font-weight: bold;
                }
                &:hover {
                    color: #3a84ff;
                }
            }
        }
        .results-wrapper {
            width: 90%;
            margin: 0 auto;
            &.searching {
                min-height: 500px;
            }
            .bk-loading {
                z-index: 9 !important;
            }
            .results-list {
                padding-top: 14px;
                color: $cmdbTextColor;
                .results-item {
                    width: 60%;
                    padding-bottom: 14px;
                    em {
                        color: #3a84ff !important;
                        font-style: normal !important;
                    }
                    .results-title {
                        display: inline-block;
                        font-size: 16px;
                        font-weight: bold;
                        margin-bottom: 4px;
                        cursor: pointer;
                        &:hover {
                            color: #3a84ff;
                            text-decoration: underline;
                        }
                    }
                    .results-desc {
                        font-size: 12px;
                        .desc-item {
                            display: inline-block;
                            padding-right: 16px;
                            padding-bottom: 6px;
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
