<template>
    <div class="search-container">
        <div class="search-box" v-click-outside="handleClickOutside">
            <input id="indexSearch" class="search-keyword" type="text" maxlength="40" :placeholder="$t('Index[\'开始查询\']')"
                v-model.trim="keyword"
                @focus="focus = true"
                @keydown="handleKeydow($event)">
            <label class="bk-icon icon-search" for="indexSearch"></label>
            <ul ref="searchList" class="search-list" v-show="focus && !loading && searchList.length">
                <li ref="searchItem" class="search-item clearfix" v-for="(result, index) in searchList" :key="index" @click="handleRoute(index)">
                    <search-item-match class="fl"
                        :result="result"
                        :keyword="keyword">
                    </search-item-match>
                    <span class="search-item-source fr">{{result['biz'][0]['bk_biz_name']}}</span>
                </li>
            </ul>
            <div class="search-loading" v-show="loading">
                <div v-bkloading="{isLoading: loading}" style="height: 100%;"></div>
            </div>
            <div class="search-empty" v-show="!loading && keyword.length > 2 && !searchList.length">{{$t("Common['暂时没有数据']")}}</div>
        </div>
    </div>
</template>

<script>
    import Throttle from 'lodash.throttle'
    export default {
        data () {
            return {
                focus: false,
                keyword: '',
                searchList: [],
                searchParams: {
                    'page': {
                        'start': 0,
                        'limit': 20,
                        'sort': 'bk_host_id'
                    },
                    'pattern': '',
                    'ip': {
                        'flag': 'bk_host_innerip|bk_host_outerip',
                        'exact': 0,
                        'data': []
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
                },
                loading: false,
                highlightIndex: -1,
                cancelSource: null
            }
        },
        watch: {
            keyword (keyword) {
                if (keyword.length > 2) {
                    this.searchParams.ip.data = [keyword]
                    this.loading = true
                    this.handleSearch(keyword)
                } else {
                    this.loading = false
                    this.cancelSource && this.cancelSource.cancel()
                    this.searchParams.ip.data = []
                    this.searchList = []
                }
                this.highlightIndex = -1
            },
            highlightIndex (highlightIndex) {
                let $searchItems = this.$refs.searchItem
                if ($searchItems && $searchItems.length) {
                    $searchItems.forEach(($searchItem, index) => {
                        index === highlightIndex ? $searchItem.classList.add('highlight') : $searchItem.classList.remove('highlight')
                    })
                }
                let scrollCount = 6
                let searchItemHeight = 48
                let $searchList = this.$refs.searchList
                if (highlightIndex !== -1 && (highlightIndex + 1) > scrollCount) {
                    $searchList.scrollTop = (highlightIndex - scrollCount + 1) * searchItemHeight
                } else {
                    $searchList.scrollTop = 0
                }
            }
        },
        methods: {
            // 函数节流，500ms发起一次主机查询
            handleSearch: Throttle(function () {
                this.cancelSource && this.cancelSource.cancel()
                if (this.keyword.length > 2) {
                    this.cancelSource = this.$Axios.CancelToken.source()
                    this.loading = true
                    this.$axios.post('hosts/search', this.searchParams, {cancelToken: this.cancelSource.token}).then(res => {
                        this.searchList = this.initSearchList(res.data.info)
                        this.loading = false
                        this.cancelSource = null
                    }).catch((e) => {
                        if (!this.$Axios.isCancel(e)) {
                            this.cancelSource = null
                            this.loading = false
                        }
                    })
                }
            }, 500, {leading: false, trailing: true}),
            initSearchList (data) {
                let list = []
                data.forEach(item => {
                    item.biz.forEach(biz => {
                        let uniqueItem = {...item}
                        uniqueItem['biz'] = [biz]
                        list.push(uniqueItem)
                    })
                })
                return list
            },
            // 点击搜索结果，进行路由跳转
            handleRoute (index) {
                this.highlightIndex = index
                this.handleEnter()
            },
            // 处理回车与上下箭头事件
            handleKeydow (event) {
                let keydownFunc = {
                    'Enter': this.handleEnter,
                    'ArrowDown': this.handleArrow,
                    'ArrowUp': this.handleArrow
                }
                let eventKey = event.key
                if (keydownFunc.hasOwnProperty(eventKey)) {
                    keydownFunc[event.key](event)
                }
            },
            // 跳转至选中的搜索结果的路由
            handleEnter () {
                if (this.highlightIndex >= 0) {
                    let targetResult = this.searchList[this.highlightIndex]
                    let targetInnerip = targetResult['host']['bk_host_innerip']
                    let targetDefault = targetResult['biz'][0]['default'] // 1为资源池下的主机，其他为业务下的主机
                    this.$store.commit('setQuickSearchParams', {
                        text: targetInnerip,
                        type: 'ip'
                    })
                    if (targetDefault === 1) {
                        this.$router.push('/resource')
                    } else {
                        this.$store.commit('setBkBizId', targetResult['biz'][0]['bk_biz_id'])
                        this.$router.push('/hosts')
                    }
                }
            },
            // 上下箭头事件，高亮搜索结果
            handleArrow (event) {
                event.preventDefault()
                let eventKey = event.key
                if (eventKey === 'ArrowDown' && this.highlightIndex < this.searchList.length - 1) {
                    this.highlightIndex++
                } else if (eventKey === 'ArrowUp' && this.highlightIndex >= 0) {
                    this.highlightIndex--
                }
            },
            handleClickOutside () {
                this.focus = false
            }
        },
        components: {
            'search-item-match': {
                props: {
                    result: {
                        required: true
                    },
                    keyword: {
                        type: String,
                        required: true
                    }
                },
                render (h) {
                    // todo 用render函数处理搜索结果的文本，高亮搜索结果中的关键词，需后端返回分词结果
                    let innerip = this.result['host']['bk_host_innerip']
                    return h('div', innerip)
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    .search-box{
        position: relative;
        .search-keyword{
            width: 100%;
            height: 44px;
            border-radius: 2px;
            border: solid 1px #c3cdd7;
            padding: 0 56px 0 16px;
            font-size: 16px;
        }
        .icon-search{
            position: absolute;
            top: 12px;
            right: 16px;
            font-size: 20px;
        }
    }
    .search-list{
        position: absolute;
        top: 100%;
        left: 0;
        width: 100%;
        max-height: 289px;
        overflow-x: hidden;
        overflow-y: auto;
        background-color: #fff;
        z-index: 1;
        box-shadow: 0px 1px 3px 0px rgba(0, 0, 0, 0.13);
        border-radius: 2px;
        border-top-right-radius: 0;
        border-top-left-radius: 0;
        border: solid 1px #c3cdd7;
        border-top: none;
        @include scrollbar;
        &::-webkit-scrollbar-thumb {
            background: #dde4eb;
        }
        .search-item{
            height: 48px;
            line-height: 48px;
            padding: 0 20px 0 17px;
            &.highlight,
            &:hover{
                background-color: #f2f7fe;
            }
            .search-item-source{
                font-size: 12px;
                color: #c3cdd7;
            }
        }
    }
    .search-loading,
    .search-empty{
        position: absolute;
        top: 100%;
        left: 0;
        width: 100%;
        height: 32px;
        line-height: 32px;
        padding: 0 20px 0 17px;
        background-color: #fff;
        border-radius: 2px;
        border-top-right-radius: 0;
        border-top-left-radius: 0;
        border: solid 1px #c3cdd7;
        border-top: none;
        font-size: 12px;
    }
</style>