<template>
    <div class="full-text-search-layout">
        <div class="full-text-sticky-layout" ref="topSticky">
            <div class="search-bar">
                <input id="fullTextSearch"
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
                <div class="pagination clearfix">
                    <div class="pagination-info fl">
                        <span class="mr10">{{$tc("Common['共计N条']", pagination['total'], { N: pagination['total'] })}}</span>
                        <i18n path="Common['每页显示条数']">
                            <div ref="paginationSize" place="num"
                                :class="['pagination-size', { 'active': isShowSizeSetting }]"
                                :sizeDirection="'auto'"
                                @click="isShowSizeSetting = !isShowSizeSetting"
                                v-click-outside="closeSizeSetting">
                                <span>{{pagination['limit']}}</span>
                                <i class="bk-icon icon-angle-down"></i>
                                <transition :name="transformOrigin === 'top' ? 'toggle-slide-top' : 'toggle-slide-bottom'">
                                    <ul :class="['pagination-size-setting', transformOrigin === 'top' ? 'bottom' : 'top']"
                                        v-show="sizeTransitionReady">
                                        <li v-for="(limit, index) in limitList"
                                            :key="index"
                                            :class="['size-setting-item', { 'selected': pagination['limit'] === limit }]"
                                            @click="handleLimitChange(limit)">
                                            {{limit}}
                                        </li>
                                    </ul>
                                </transition>
                            </div>
                        </i18n>
                    </div>
                    <bk-paging class="pagination-list fr"
                        v-if="pagination['count'] > 1"
                        :type="'compact'"
                        :total-page="pagination['count']"
                        :size="'small'"
                        :cur-page="pagination['current']"
                        @page-change="handlePageChange">
                    </bk-paging>
                </div>
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
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                reg: /(?<=keywords=).*/,
                scrollTop: null,
                isShowSizeSetting: false,
                sizeTransitionReady: false,
                transformOrigin: 'top',
                currentClassify: -1,
                requestId: 'fullTextSearch',
                // queryString: '',
                query: {
                    queryString: '',
                    objId: ''
                },
                properties: [],
                hasData: false,
                showNoData: false,
                searching: false,
                limitList: [10, 20, 50, 100],
                pagination: {
                    limit: 10,
                    start: 0,
                    count: 1,
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
            ...mapGetters('objectModelClassify', ['models', 'getModelById']),
            params () {
                return {
                    page: {
                        start: this.pagination.start,
                        limit: this.pagination.limit
                    },
                    bk_obj_id: this.query.objId,
                    query_string: `${this.query.queryString}`,
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
                window.location.hash = this.hash.replace(this.reg, this.query.queryString)
                this.query.objId = ''
                if (queryString) {
                    this.handleSearch(600)
                }
            }
        },
        created () {
            this.$store.commit('setHeaderStatus', {
                back: true
            })
            this.$store.commit('setHeaderTitle', this.$t('Common["搜索结果"]'))
            this.query.queryString = this.reg.exec(this.hash)[0]
        },
        mounted () {
            this.initScrollListener(this.$refs.topSticky)
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
            closeSizeSetting () {
                this.isShowSizeSetting = false
            },
            calcSizePosition () {
                const paginationSizeDom = this.$refs.paginationSize
                if (paginationSizeDom.getAttribute('sizeDirection') && paginationSizeDom.getAttribute('sizeDirection') !== 'auto') {
                    this.transformOrigin = this.pagination.sizeDirection
                } else {
                    const sizeSettingItemHeight = 32
                    const sizeSettingHeight = this.limitList.length * sizeSettingItemHeight
                    const paginationSizeRect = paginationSizeDom.getBoundingClientRect()
                    const bodyRect = document.body.getBoundingClientRect()
                    if (bodyRect.height - paginationSizeRect.y - paginationSizeRect.height > sizeSettingHeight) {
                        this.transformOrigin = 'bottom'
                    } else {
                        this.transformOrigin = 'top'
                    }
                }
                this.sizeTransitionReady = true
            },
            isPublicModel (objId) {
                const model = this.models.find(model => model['bk_obj_id'] === objId) || {}
                return !this.$tools.getMetadataBiz(model)
            },
            toggleClassify (index, model) {
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
                        this.setPagination(data.total)
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
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                if (source['hitsType'] === 'host') {
                    this.$router.push({
                        name: 'resource',
                        query: {
                            ip: source['bk_host_innerip'],
                            outer: false,
                            inner: true,
                            exact: 1,
                            assigned: true
                        }
                    })
                } else if (source['hitsType'] === 'object') {
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
            setPagination (total) {
                this.pagination.total = total
                this.pagination.count = Math.ceil(total / this.pagination.limit)
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
        .pagination {
            border-top: 1px solid #dde4eb;
            padding-top: 10px;
            color: $cmdbTextColor;
            font-size: 14px;
            line-height: 32px;
            .pagination-size {
                position: relative;
                display: inline-block;
                width: 54px;
                height: 32px;
                border: 1px solid #c3cdd7;
                padding: 0 25px 0 6px;
                margin-left: 6px;
                margin-right: 6px;
                vertical-align: middle;
                border-radius: 2px;
                cursor: pointer;
                .icon-angle-down {
                    position: absolute;
                    right: 6px;
                    top: 9px;
                    font-size: 12px;
                    color: #c3cdd7;
                    transition: transform .2s linear;
                }
                &.active {
                    border-color: #0082ff;
                    .icon-angle-down {
                        color: #0082ff;
                        transform: rotate(180deg);
                    }
                }
            }
        }
        .pagination-size-setting{
            position: absolute;
            left: 0;
            width: 100%;
            background-color: #fff;
            box-shadow: 0 0 1px 1px rgba(0, 0, 0, 0.1);
            z-index: 10001;
            &.top {
                top: 32px;
            }
            &.bottom {
                bottom: 32px;
            }
            .size-setting-item {
                height: 32px;
                line-height: 32px;
                padding: 0 12px;
                &.selected,
                &:hover{
                    background-color: #f5f5f5;
                    color: #0082ff;
                }
            }
        }
        .toggle-slide-bottom-enter-active,
        .toggle-slide-bottom-leave-active,
        .toggle-slide-top-enter-active,
        .toggle-slide-top-leave-active {
            transition: transform .3s cubic-bezier(.23, 1, .23, 1), opacity .5s cubic-bezier(.23, 1, .23, 1);
            transform-origin: center top;
        }
        .toggle-slide-top-enter-active,
        .toggle-slide-top-leave-active {
            transform-origin: center bottom;
        }
        .toggle-slide-bottom-enter,
        .toggle-slide-bottom-leave-active,
        .toggle-slide-top-enter,
        .toggle-slide-top-leave-active {
            transform: translateZ(0) scaleY(0);
            opacity: 0;
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
