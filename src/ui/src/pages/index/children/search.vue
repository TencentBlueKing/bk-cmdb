<template>
    <div class="search-layout">
        <div class="search-box" v-click-outside="handleClickOutside">
            <input id="indexSearch" class="search-keyword" type="text" maxlength="40" :placeholder="$t('Index[\'开始查询\']')"
                v-model.trim="keyword"
                @focus="focus = true">
            <label class="bk-icon icon-search" for="indexSearch"></label>
            <div class="search-result" v-show="focus && keyword.length > 2">
                <div class="search-loading" v-bkloading="{isLoading: loading}" v-if="loading"></div>
                <div :class="['result-layout', {'result-layout-empty': !resultTabpanels.length}]" v-show="!loading">
                    <bk-tab class="result-tab" v-if="resultTabpanels.length"
                        :active-name.sync="resultTab.active"
                        :size="'small'"
                        :headStyle="resultTab.headStyle">
                        <bk-tabpanel v-for="(panel, index) in resultTabpanels"
                            :key="index"
                            :name="panel"
                            :title="getPanelTitle(panel)">
                            <v-search-item :list="resultTab.list[panel]" :model="panel"></v-search-item>
                        </bk-tabpanel>
                    </bk-tab>
                    <div class="result-empty" v-else>{{$t('Common["暂时没有数据"]')}}</div>
                    <div class="result-more" v-if="hasMore()" @click="showMore">{{$t('Index["查看更多结果"]')}}</div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import vSearchItem from './search-item'
    import Throttle from 'lodash.throttle'
    export default {
        components: {
            vSearchItem
        },
        data () {
            return {
                focus: false,
                keyword: '',
                searchParams: {
                    'page': {
                        'start': 0,
                        'limit': 15,
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
                cancelSource: null,
                resultTab: {
                    active: 'host',
                    headStyle: {
                        'height': '40px',
                        'color': '#3c96ff'
                    },
                    list: {},
                    count: {}
                }
            }
        },
        computed: {
            ...mapGetters('navigation', ['classifications']),
            allModels () {
                let allModels = []
                this.classifications.forEach(classify => {
                    classify['bk_objects'].forEach(model => {
                        allModels.push(model)
                    })
                })
                return allModels
            },
            resultTabpanels () {
                return Object.keys(this.resultTab.list)
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
                    this.resultTab.list = {}
                    this.resultTab.count = {}
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
                        if (res.result) {
                            this.addResultList('host', this.initSearchList(res.data.info))
                            this.addResultCount('host', res.data.count)
                            this.loading = false
                            this.cancelSource = null
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    }).catch((e) => {
                        if (!this.$Axios.isCancel(e)) {
                            this.cancelSource = null
                            this.loading = false
                        }
                    })
                }
            }, 500, {leading: false, trailing: true}),
            addResultList (model, list) {
                if (list.length) {
                    this.$set(this.resultTab.list, model, list)
                } else {
                    this.$delete(this.resultTab.list, model)
                }
            },
            addResultCount (model, count) {
                if (count) {
                    this.$set(this.resultTab.count, model, count)
                } else {
                    this.$delete(this.resultTab.count, model)
                }
            },
            getPanelTitle (panel) {
                const panelModel = this.allModels.find(model => model['bk_obj_id'] === panel)
                return panelModel ? `${panelModel['bk_obj_name']}(${this.resultTab.count[panel]})` : null
            },
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
            hasMore () {
                return this.resultTabpanels.length && this.resultTab.count[this.resultTab.active] > 15
            },
            showMore () {
                const funcMaps = {
                    'host': this.showMoreHost
                }
                const model = this.resultTab.active
                if (funcMaps.hasOwnProperty(model)) {
                    funcMaps[model]()
                }
                this.handleClickOutside()
            },
            showMoreHost () {
                this.$store.commit('setHostSearch', {
                    ip: this.keyword,
                    exact: 0,
                    innerip: true,
                    outerip: true,
                    assigned: true
                })
                this.$router.push('/resource')
            },
            handleClickOutside () {
                this.focus = false
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
            }
        }
    }
</script>
<style lang="scss" scoped>
    .search-layout{}
    .search-box{
        position: relative;
        height: 42px;
        overflow: visible;
        .search-keyword{
            width: 100%;
            height: 100%;
            border-radius: 2px;
            border: solid 1px #c3cdd7;
            padding: 0 56px 0 16px;
            font-size: 14px;
        }
        .icon-search{
            position: absolute;
            top: 12px;
            right: 16px;
            font-size: 20px;
        }
    }
    .search-result{
        position: absolute;
        top: 100%;
        left: 0;
        width: 100%;
        margin: 2px 0 0 0;
        background:rgba(255,255,255,1);
        box-shadow:0px 3px 6px 0px rgba(51,60,72,0.1);
        border-radius:2px;
        border:1px solid rgba(221,228,235,1);
        z-index: 2;
    }
    .search-loading {
        height: 40px;
    }
    .result-layout{
        height: 100%;
        font-size: 0;
        &-empty:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            height: 100%;
        }
        .result-tab{
            display: inline-block;
            width: 100%;
            border: none;
            padding: 0 20px;
        }
        .result-empty{
            height: 40px;
            line-height: 40px;
            width: 100%;
            display: inline-block;
            vertical-align: middle;
            text-align: center;
            font-size: 12px;
        }
        .result-more{
            height:33px;
            line-height: 32px;
            font-size: 12px;
            color: #3c96ff;
            text-align: center;
            background:rgba(250,251,253,1);
            border-top: 1px solid #ebf0f5;
            cursor: pointer;
        }
    }
</style>