<template>
    <div class="full-text-search-layout">
        <div class="full-text-sticky-layout" ref="topSticky">
            <div class="search-bar">
                <input id="fullTextSearch"
                    ref="searchInput"
                    autocomplete="off"
                    class="search-keywords"
                    type="text"
                    maxlength="32"
                    :placeholder="$t('请输入搜索内容')"
                    v-model.trim="query.queryString"
                    @input="handleInputSearch"
                    @keypress.enter="handleSearch">
                <label class="bk-icon icon-search" for="fullTextSearch" v-if="!query.queryString"></label>
                <i class="bk-icon icon-close-circle-shape" v-else @click="handleClearInput"></i>
            </div>
            <div class="classify" v-show="hasData">
                <span class="classify-item"
                    :class="['classify-item', { 'classify-active': -1 === currentClassify }]"
                    @click="toggleClassify(-1)">
                    {{$t('全部结果')}}（{{allSearchCount}}）
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
            v-bkloading="{ 'isLoading': $loading(requestId), 'afterLeave': loadingClose, opacity: 1 }">
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
                                <span class="desc-item" v-html="`${$t('主机ID')}：${source['bk_host_id']}`"> </span>
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
                                <span class="desc-item" v-html="`${$t('业务ID')}：${source['bk_biz_id']}`"> </span>
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
                <div class="pagination-info">
                    <span class="mr10">{{$tc('共计N条', pagination['total'], { N: pagination['total'] })}}</span>
                    <bk-pagination
                        size="small"
                        align="right"
                        :type="'compact'"
                        :current.sync="pagination['current']"
                        :limit="pagination['limit']"
                        :count="pagination['total']"
                        @limit-change="handleLimitChange"
                        @change="handlePageChange">
                    </bk-pagination>
                </div>
            </div>
        </div>
        <div class="no-data" v-show="!hasData && showNoData && !searching">
            <img src="../../../assets/images/full-text-search.png" alt="no-data">
            <p>{{$t('搜不到相关内容')}}</p>
        </div>
    </div>
</template>

<script>
    import {
        addMainScrollListener,
        removeMainScrollListener
    } from '@/utils/main-scroller'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                toggleTips: null,
                scrollTop: null,
                currentClassify: -1,
                requestId: 'fullTextSearch',
                searchTriggerType: 'input',
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
                debounceTimer: null,
                beforeSearchString: ''
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('objectModelClassify', ['models', 'getModelById']),
            params () {
                const notZhCn = this.query.queryString.replace(/\w\.?/g, '').length === 0
                const singleSpecial = /[!"#$%&'()\*,-\./:;<=>?@\[\\\]^_`{}\|~]{1}/
                const queryString = this.query.queryString.length === 1 ? this.query.queryString.replace(singleSpecial, '') : this.query.queryString
                return {
                    page: {
                        start: this.pagination.start,
                        limit: this.pagination.limit
                    },
                    bk_obj_id: this.query.objId,
                    bk_biz_id: this.bizId ? this.bizId.toString() : '',
                    query_string: notZhCn ? `*${queryString}*` : queryString,
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
                let delay = 0
                if (this.searchTriggerType === 'click') {
                    this.query.objId = this.query.objId
                } else {
                    this.query.objId = ''
                    delay = 600
                }
                if (queryString) {
                    this.pagination.start = 0
                    this.pagination.current = 1
                    this.handleSearch(delay)
                }
            }
        },
        created () {
            this.$store.commit('setHeaderStatus', {
                back: true
            })
            this.$store.commit('setHeaderTitle', this.$t('搜索结果'))
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
                if (index === this.currentClassify || this.searching) return
                this.searchTriggerType = 'click'
                this.currentClassify = index
                this.query.objId = model
                if (this.query.queryString) {
                    this.pagination.start = 0
                    this.pagination.current = 1
                    this.handleSearch()
                } else {
                    this.query.queryString = this.beforeSearchString
                }
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
                    window.location.hash = `${this.hash.substring(0, this.hash.search(/=/) + 1)}${this.query.queryString}&from=${this.$route.query.from}`
                    this.searching = true
                    this.$store.dispatch('fullTextSearch/search', {
                        params: this.params,
                        config: {
                            requestId: this.requestId,
                            cancelPrevious: true
                        }
                    }).then(async data => {
                        this.pagination.total = data.total
                        const hitsData = data.hits
                        const modelData = data.aggregations || []
                        if (data.total) {
                            if (!this.query.objId) {
                                this.currentClassify = -1
                                this.modelClassify = modelData.map(item => {
                                    return {
                                        ...this.getModelById(item.key),
                                        count: item.count > 999 ? '999+' : item.count
                                    }
                                }).sort((prev, next) => next.count - prev.count)
                                this.modelClassify.forEach(model => {
                                    this.modelClassifyName[model['bk_obj_id']] = model['bk_obj_name']
                                })
                                await this.processArray(modelData)
                                this.allSearchCount = data.total > 999 ? '999+' : data.total
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
                            this.beforeSearchString = this.query.queryString
                        } else {
                            this.hasData = false
                        }
                        this.$nextTick(() => {
                            this.initScrollListener(this.$refs.topSticky)
                            this.$parent.$refs.mainScroller.scrollTop = 0
                        })
                    }).catch((e) => {
                        console.error(e)
                        if (!e.hasOwnProperty('message')) {
                            this.searching = false
                            this.hasData = false
                        }
                    })
                }, wait)
            },
            loadingClose () {
                this.showNoData = true
                this.searching = false
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
                        name: 'resourceHostDetails',
                        params: {
                            id: source['bk_host_id']
                        },
                        query: {
                            from: this.$route.fullPath
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
                        name: 'generalModel',
                        params: {
                            objId: source['bk_obj_id'],
                            instId: source['bk_inst_id'].toString().replace(/(\<\/?em\>)/g, '')
                        },
                        query: {
                            from: this.$route.fullPath
                        }
                    })
                } else if (source['hitsType'] === 'biz') {
                    this.$router.push({
                        name: 'business',
                        params: {
                            bizName: source['bk_biz_name'].toString().replace(/(\<\/?em\>)/g, '')
                        },
                        query: {
                            from: this.$route.fullPath
                        }
                    })
                }
            },
            getShowPropertyText (property, source, thisProperty) {
                const cloneSource = this.$tools.clone(source)
                const reg = /\<em\>.+\<\/em\>/
                const propertyValue = cloneSource[thisProperty].toString()
                const isHeightLight = reg.test(propertyValue)
                cloneSource[thisProperty] = isHeightLight ? propertyValue.replace(/(\<\/?em\>)/g, '') : propertyValue
                const flatternedText = this.$tools.getPropertyText(property, cloneSource)
                return isHeightLight ? `<em>${flatternedText}</em>` : flatternedText
            },
            handleLimitChange (limit) {
                if (this.pagination.limit === limit) return
                this.pagination.limit = limit
                this.handlePageChange(1)
            },
            handlePageChange (current, type) {
                this.pagination.start = (current - 1) * this.pagination.limit
                this.pagination.current = current
                this.handleSearch()
            },
            handleInputSearch (el) {
                const target = el.srcElement
                this.searchTriggerType = 'input'
                if (target.value.length === 32) {
                    this.toggleTips && this.toggleTips.destroy()
                    this.toggleTips = this.$bkPopover(this.$refs.searchInput, {
                        theme: 'dark max-length-tips',
                        content: this.$t('最大支持搜索32个字符'),
                        zIndex: 9999,
                        trigger: 'manual',
                        boundary: 'window',
                        arrow: true
                    })
                    this.$nextTick(() => {
                        this.toggleTips.show()
                    })
                } else if (this.toggleTips) {
                    this.toggleTips.hide()
                }
            },
            handleClearInput () {
                this.query.queryString = ''
                if (this.toggleTips) {
                    this.toggleTips.hide()
                }
            }
        }
    }
</script>

<style lang="scss">
    .full-text-search-layout {
        position: relative;
        overflow: auto;
        padding: 0 0 20px !important;
        .full-text-tips {
            left: 248px !important;
        }
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
            .bk-icon {
                position: absolute;
                right: 20px;
                top: 13px;
                font-size: 18px;
                &.icon-close-circle-shape {
                    font-size: 16px;
                    right: 10px;
                    color: #c4c6cc;
                    cursor: pointer;
                    &:hover {
                        color: #979ba5;
                    }
                }
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
        .pagination-info {
            font-size: 14px;
            line-height: 30px;
            color: #737987;
            display: flex;
            .bk-page {
                flex: 1;
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
    .max-length-tips-theme {
        font-size: 12px;
        padding: 6px 12px;
        left: 248px !important;
    }
</style>
