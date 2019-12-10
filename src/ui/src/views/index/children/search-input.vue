<template>
    <div class="search-layout" ref="searchLayout" :style="{ 'background-color': setStyle.backgroundColor }">
        <div :class="{ 'sticky-layout': result.isShow }" :style="{ 'padding-top': (setStyle.paddingTop) + 'px' }">
            <div class="search-bar">
                <div class="search-box" v-if="result.isShow">
                    <div :class="['search-tab', { 'is-focus': isFocus }]">
                        <span :class="['tab-item', { 'active': activeName === 'host' }]"
                            @click="handleChangeTab('host')">
                            {{$t('主机搜索')}}
                        </span>
                        <span :class="['tab-item', { 'active': activeName === 'fullText' }]"
                            @click="handleChangeTab('fullText')">
                            {{$t('全文检索')}}
                        </span>
                    </div>
                    <div class="tab-content">
                        <host-search v-show="activeName === 'host'" @focus="handleFocus"></host-search>
                        <div class="input-box" v-show="activeName === 'fullText'" v-click-outside="handleHideLenovo">
                            <bk-input ref="searchInput"
                                class="search-input"
                                autocomplete="off"
                                maxlength="32"
                                :placeholder="placeholder"
                                v-model.trim="keywords"
                                @input="handleInputSearch"
                                @focus="handleShowHistory"
                                @blur="isFocus = false"
                                @enter="handleShowResult">
                            </bk-input>
                            <i class="bk-icon search-clear icon-close-circle-shape" v-show="keywords" @click="handleClear"></i>
                            <bk-button theme="primary" class="search-btn" @click="handleShowResult">
                                <i class="bk-icon icon-search"></i>
                                {{$t('搜索')}}
                            </bk-button>
                        </div>
                    </div>
                </div>
                <div class="input-box" v-else v-click-outside="handleHideLenovo">
                    <bk-input ref="searchInput"
                        class="search-input"
                        autocomplete="off"
                        maxlength="32"
                        :placeholder="placeholder"
                        v-model.trim="keywords"
                        @input="handleInputSearch"
                        @focus="handleShowHistory"
                        @blur="$emit('focus', false)"
                        @enter="handleShowResult">
                    </bk-input>
                    <i class="bk-icon search-clear icon-close-circle-shape" v-show="keywords" @click="handleClear"></i>
                    <bk-button theme="primary" class="search-btn" @click="handleShowResult">
                        <i class="bk-icon icon-search"></i>
                        {{$t('搜索')}}
                    </bk-button>
                </div>
                <transition name="slide">
                    <div class="lenovo selectTips" :style="{ top: result.isShow ? '76px' : '47px' }" v-show="showLenovo && lenovoList.length">
                        <ul class="lenovo-result">
                            <template v-for="(item, index) in lenovoList">
                                <li :key="index" v-if="item.type === 'host'"
                                    :title="item.source.bk_host_innerip"
                                    @click="handleGoResource(item.source)">
                                    <span class="lenovo-name">{{item.source.bk_host_innerip}}</span>
                                    <span>({{$t('主机')}})</span>
                                </li>
                                <li :key="index" v-else-if="item.type === 'object'"
                                    :title="item.source.bk_inst_name"
                                    @click="handleGoInstace(item.source)">
                                    <span class="lenovo-name">{{item.source.bk_inst_name}}</span>
                                    <span>({{getShowModelName(item.source)}})</span>
                                </li>
                                <li :key="index" v-else-if="item.type === 'biz'"
                                    :title="item.source.bk_biz_name"
                                    @click="handleGoBusiness(item.source)">
                                    <span class="lenovo-name">{{item.source.bk_biz_name}}</span>
                                    <span>({{$t('业务')}})</span>
                                </li>
                            </template>
                        </ul>
                    </div>
                </transition>

                <transition name="slide">
                    <div class="history selectTips" :style="{ top: result.isShow ? '76px' : '47px' }" v-show="showHistory">
                        <div class="history-title clearfix" @click.stop>
                            <span class="fl">{{$t('搜索历史')}}</span>
                            <bk-button :text="true" class="clear-btn fr" @click="handlClearHistory">
                                <i class="bk-icon icon-cc-delete"></i>
                                {{$t('清空')}}
                            </bk-button>
                        </div>
                        <ul class="history-list">
                            <li v-for="(history, index) in historyList"
                                :key="index"
                                @click="handleHistorySearch(history)">
                                {{history}}
                            </li>
                        </ul>
                    </div>
                </transition>
            </div>

            <div class="classify" v-if="pagination['total'] && result.isShow">
                <span class="classify-item"
                    :class="['classify-item', { 'classify-active': -1 === currentClassify }]"
                    @click="toggleClassify(-1)">
                    {{$t('全部结果')}}（{{searchCount}}）
                </span>
                <span class="classify-item"
                    v-for="(model, index) in modelMap"
                    :key="index"
                    :class="['classify-item', { 'classify-active': index === currentClassify }]"
                    @click="toggleClassify(index, model['bk_obj_id'])">
                    {{model['bk_obj_name']}}（{{model['count']}}）
                </span>
            </div>
        </div>

        <div class="search-main"
            v-bkloading="{ 'isLoading': updating && $loading() }"
            v-if="result.isShow">
            <search-result :query-data="result.data" :model-classify="modelMap"></search-result>
            <div class="pagination-info" v-show="pagination['total']">
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
</template>

<script>
    import searchResult from './full-text-search'
    import hostSearch from './host-search'
    import { mapGetters } from 'vuex'
    import { MENU_INDEX, MENU_RESOURCE_INSTANCE, MENU_RESOURCE_BUSINESS, MENU_RESOURCE_HOST_DETAILS } from '@/dictionary/menu-symbol'
    export default {
        components: {
            searchResult,
            hostSearch
        },
        props: {
            isFullTextSearch: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                isFocus: false,
                activeName: 'fullText',
                toggleTips: null,
                beforeKeywords: '',
                keywords: '',
                debounceTimer: null,
                lenovoList: [],
                modelMap: [],
                curModelMap: [],
                searching: false,
                updating: false,
                showHistory: false,
                showLenovo: false,
                currentClassify: -1,
                searchCount: 0,
                setStyle: {
                    marginTop: 0,
                    backgroundColor: 'transparent'
                },
                query: {
                    trigger: 'input',
                    objId: '',
                    data: {}
                },
                result: {
                    _resolve: null,
                    _promise: null,
                    isShow: false,
                    data: {}
                },
                pagination: {
                    limit: 10,
                    start: 0,
                    current: 1,
                    total: 0
                },
                curPagination: {
                    limit: 10,
                    start: 0,
                    current: 1,
                    total: 0
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('objectModelClassify', ['models', 'getModelById']),
            placeholder () {
                return this.isFullTextSearch ? this.$t('请输入关键字，点击或回车搜索') : this.$t('请输入IP开始搜索')
            },
            params () {
                const keywords = this.keywords
                const notZhCn = keywords.replace(/\w\.?/g, '').length === 0
                const singleSpecial = /[!"#$%&'()\*,-\./:;<=>?@\[\\\]^_`{}\|~]{1}/
                const queryString = keywords.length === 1 ? keywords.replace(singleSpecial, '') : keywords
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
            searchParams () {
                return {
                    'page': {
                        'start': this.pagination.start,
                        'limit': this.pagination.limit,
                        'sort': 'bk_host_id'
                    },
                    'pattern': '',
                    'ip': {
                        'flag': 'bk_host_innerip|bk_host_outerip',
                        'exact': 0,
                        'data': [this.keywords]
                    },
                    'condition': [{
                        'bk_obj_id': 'host',
                        'fields': [],
                        'condition': []
                    }, {
                        'bk_obj_id': 'module',
                        'fields': [],
                        'condition': []
                    }, {
                        'bk_obj_id': 'set',
                        'fields': [],
                        'condition': []
                    }, {
                        'bk_obj_id': 'biz',
                        'fields': [],
                        'condition': []
                    }]
                }
            },
            historyList () {
                return this.$store.state.fullTextSearch.searchHistory
            },
            interfacePath () {
                return this.isFullTextSearch ? 'fullTextSearch/search' : 'hostSearch/searchHost'
            }
        },
        watch: {
            '$route' (to, from) {
                const queryLen = Object.keys(to.query).length
                if (to.path === '/index' && !queryLen) {
                    this.$emit('search-status', 0)
                    this.keywords = ''
                    this.setStyle.paddingTop = 0
                    this.setStyle.backgroundColor = 'transparent'
                    this.result.isShow = false
                }
            },
            keywords (keywords) {
                if (this.query.trigger === 'input') {
                    this.query.objId = ''
                }
                if (keywords) {
                    this.showHistory = false
                    this.showLenovo = true
                    this.pagination.start = 0
                    this.pagination.current = 1
                    this.getFullTextSearch(0)
                } else {
                    this.handleHideLenovo()
                }
            }
        },
        created () {
            this.$store.commit('fullTextSearch/getSearchHistory')
            const routerQuery = this.$route.query || {}
            if (routerQuery.show === 'true') {
                this.keywords = routerQuery.keywords
                this.$nextTick(() => {
                    this.handleShowResult()
                })
            }
        },
        methods: {
            handleFocus (status) {
                this.isFocus = status
            },
            handleChangeTab (name) {
                this.activeName = name
            },
            getShowModelName (source) {
                let modelName = ''
                try {
                    const model = this.getModelById(source['bk_obj_id'])
                    modelName = model['bk_obj_name']
                } catch (e) {
                    console.error(e)
                }
                return modelName
            },
            formatData (data) {
                const res = {
                    aggregations: [],
                    hits: [],
                    total: 0
                }
                if (data.count) {
                    res.total = data.count
                    res.aggregations = [{
                        count: data.count,
                        key: 'host'
                    }]
                    res.hits = data.info.map(item => {
                        return {
                            highlight: {},
                            type: 'host',
                            source: item.host
                        }
                    })
                }
                return res
            },
            getFullTextSearch (wait = 600) {
                let _resolve = null
                const _promise = new Promise((resolve, reject) => {
                    _resolve = resolve
                })
                if (this.debounceTimer) {
                    clearTimeout(this.debounceTimer)
                    this.debounceTimer = null
                }
                this.debounceTimer = setTimeout(async () => {
                    if (!this.keywords) return
                    this.searching = true
                    try {
                        let data = await this.$store.dispatch(this.interfacePath, {
                            params: this.isFullTextSearch ? this.params : this.searchParams,
                            config: {
                                requestId: 'search',
                                cancelPrevious: true
                            }
                        })
                        if (!this.isFullTextSearch) {
                            data = this.formatData(data)
                        }
                        if (this.query.trigger === 'click') {
                            this.result.data = data
                            this.pagination.total = data.total
                        } else {
                            this.query.data = data
                            this.curPagination.total = data.total
                        }
                        const hitData = data.hits || []
                        const modelData = data.aggregations || []
                        if (data.total) {
                            if (!this.query.objId) {
                                this.currentClassify = -1
                                this.curModelMap = modelData.sort((prev, next) => next.count - prev.count).map(item => {
                                    return {
                                        ...this.getModelById(item.key),
                                        count: item.count > 999 ? '999+' : item.count
                                    }
                                })
                            }
                        }
                        this.lenovoList = hitData.length > 8 ? hitData.slice(0, 8) : hitData
                        this.searching = false
                        this.query.trigger = 'input'
                        this.$nextTick(() => {
                            this.$refs.searchLayout.scrollTop = 0
                        })
                    } catch (e) {
                        console.error(e)
                        this.searching = false
                    }
                    _resolve && _resolve()
                }, wait)
                return _promise
            },
            handleHideLenovo () {
                this.showHistory = false
                this.showLenovo = false
                const timer = setTimeout(() => {
                    this.lenovoList = []
                    clearTimeout(timer)
                }, 300)
            },
            handleShowHistory () {
                if (!this.keywords) {
                    this.$emit('focus', true)
                }
                this.isFocus = true
                this.showHistory = !this.keywords && this.historyList.length
            },
            handlClearHistory () {
                this.$store.commit('fullTextSearch/clearSearchHistory')
                this.handleHideLenovo()
            },
            resetIndex () {
                // this.$refs.searchInput.$refs.input.focus()
                this.$router.replace({
                    name: MENU_INDEX
                })
            },
            async handleShowResult () {
                if (!this.keywords) {
                    this.resetIndex()
                    return
                }
                this.$emit('search-status', 1)
                this.query.trigger = 'input'
                this.updating = true
                await this.getFullTextSearch(0)
                this.setStyle.paddingTop = 50
                this.setStyle.backgroundColor = '#FAFBFD'
                this.showLenovo = false
                const total = this.curPagination.total
                this.searchCount = total > 999 ? '999+' : total
                this.pagination.total = total
                this.modelMap = this.$tools.clone(this.curModelMap)
                this.$set(this.result, 'data', this.$tools.clone(this.query.data))
                this.beforeKeywords = this.keywords
                this.$store.commit('fullTextSearch/setSearchHistory', this.keywords)
                this.$router.replace({
                    name: MENU_INDEX,
                    query: {
                        keywords: this.keywords,
                        show: true
                    }
                })
                this.result.isShow = true
                this.updating = false
                this.$nextTick(() => {
                    this.$refs.searchInput.$refs.input.focus()
                })
            },
            handleHistorySearch (keyword) {
                this.keywords = keyword
                this.$nextTick(() => {
                    this.handleShowResult()
                })
            },
            handleGoResource (host) {
                this.$router.push({
                    name: MENU_RESOURCE_HOST_DETAILS,
                    params: {
                        id: host.bk_host_id
                    }
                })
            },
            handleGoInstace (source) {
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
                        instId: source['bk_inst_id']
                    }
                })
            },
            handleGoBusiness (source) {
                this.$router.push({
                    name: MENU_RESOURCE_BUSINESS,
                    params: {
                        bizName: source['bk_biz_name']
                    }
                })
            },
            toggleClassify (index, objId) {
                if (index === this.currentClassify || this.searching) return
                this.query.trigger = 'click'
                this.currentClassify = index
                this.query.objId = objId
                this.keywords = this.beforeKeywords
                this.pagination.start = 0
                this.pagination.current = 1
                this.getFullTextSearch(0)
            },
            handleLimitChange (limit) {
                if (this.pagination.limit === limit) return
                this.pagination.limit = limit
                this.handlePageChange(1)
            },
            handlePageChange (current) {
                this.pagination.start = (current - 1) * this.pagination.limit
                this.pagination.current = current
                this.query.trigger = 'click'
                this.getFullTextSearch()
            },
            handleInputSearch (value) {
                if (value.length === 32) {
                    this.toggleTips && this.toggleTips.destroy()
                    this.toggleTips = this.$bkPopover(this.$refs.searchInput.$el, {
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
            handleClear () {
                this.$refs.searchInput.$refs.input.focus()
                this.keywords = ''
                this.toggleTips && this.toggleTips.hide()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-layout {
        position: relative;
        width: 100%;
        height: 100%;
        z-index: 3;
        overflow-y: auto;
        overflow-x: hidden;
        .sticky-layout {
            transition: all .3s;
            position: sticky;
            top: 0;
            padding: 50px 0 0;
            background-color: #FAFBFD;
            z-index: 4;
        }
        .search-bar {
            position: relative;
            width: 100%;
            max-width: 726px;
            margin: 0 auto 38px;
            .search-box {
                height: calc(100% + 50px);
                transition: all 0.4s;
                .search-tab {
                    max-width: 726px;
                    margin: 0 auto;
                    font-size: 0;
                    &.is-focus .tab-item.active {
                        border-color: #3A84FF;
                    }
                    .tab-item {
                        @include inlineBlock;
                        position: relative;
                        height: 30px;
                        line-height: 30px;
                        text-align: center;
                        padding: 0 14px;
                        margin: 0 4px -1px 0;
                        font-size: 14px;
                        color: #63656E;
                        background-color: #DCDEE5;
                        border: 1px solid #C4C6CC;
                        border-radius: 6px 6px 0 0;
                        transition: all 0.2s;
                        cursor: pointer;
                        &.active {
                            background-color: #FFFFFF;
                            border-bottom-color: #FFFFFF !important;
                            z-index: 1000;
                        }
                    }
                }
                .tab-content {
                    height: 100%;
                }
            }
            .input-box {
                width: 100%;
                display: flex;
            }
            .search-input {
                flex: 1;
                font-size: 0;
                /deep/ .bk-form-input {
                    font-size: 14px;
                    height: 42px;
                    line-height: 42px;
                    padding: 0 56px 0 16px;
                    border-radius: 2px 0 0 2px;
                }
            }
            .search-btn {
                width: 86px;
                height: 42px;
                line-height: 42px;
                padding: 0;
                border-radius: 0 2px 2px 0;
                .icon-search {
                    width: 18px;
                    height: 18px;
                    font-size: 18px;
                    margin: -2px 4px 0 0;
                }
            }
            .search-clear {
                position: absolute;
                right: 86px;
                bottom: 0;
                width: 50px;
                height: 42px;
                line-height: 42px;
                color: #C3CDD7;
                font-size: 18px;
                text-align: center;
                cursor: pointer;
                &:hover {
                    color: #979BA5;
                }
            }
            .selectTips {
                position: absolute;
                top: 47px;
                left: 0;
                width: calc(100% - 86px);
                background-color: #ffffff;
                box-shadow: 0px 2px 6px 0px rgba(0,0,0,0.15);
                border: 1px solid #DCDEE5;
                overflow: hidden;
                z-index: 99;
                ul li {
                    color: #63656E;
                    font-size: 14px;
                    padding: 0 20px;
                    line-height: 40px;
                    cursor: pointer;
                    &:hover {
                        color: #3A84FF;
                        background-color: #E1ECFF;
                    }
                }
            }
            .lenovo-result li{
                display: flex;
                .lenovo-name {
                    max-width: 86%;
                    @include ellipsis;
                }
                span:not(.lenovo-name) {
                    padding-left: 10px;
                    color: #C4C6CC;
                }
            }
            .history-title {
                font-size: 14px;
                line-height: 36px;
                color: #C4C6CC;
                padding: 5px 20px 0;
                &::after {
                    content: '';
                    display: block;
                    height: 1px;
                    background-color: #F0F1F5;
                }
                .clear-btn {
                    color: #C4C6CC;
                    &:hover {
                        color: #979BA5;
                    }
                    .icon-cc-delete {
                        margin-top: -2px;
                    }
                }
            }
            .history-list {
                margin-bottom: 5px;
            }
        }
        .classify {
            width: 90%;
            margin: 0 auto;
            color: $cmdbTextColor;
            background-color: #FAFBFD;
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
        .search-main {
            padding-bottom: 20px;
            /deep/ .bk-loading {
                background-color: #FAFBFD !important;
                z-index: 3 !important;
            }
        }
        .pagination-info {
            width: 90%;
            font-size: 14px;
            line-height: 30px;
            color: #737987;
            margin: 0 auto;
            display: flex;
            .bk-page {
                flex: 1;
            }
        }
        .slide-enter-active {
            transition: all .3s;
            max-height: 332px;
            opacity: 1;
            overflow: hidden;
        }
        .slide-enter, .slide-leave-to{
            opacity: 0;
            max-height: 0;
        }
    }
</style>

<style lang="scss">
    .max-length-tips-theme {
        font-size: 12px;
        padding: 6px 12px;
        left: 248px !important;
    }
</style>
