<template>
    <div class="host-search-layout">
        <div class="search-bar">
            <bk-input class="search-input"
                ref="searchInput"
                type="textarea"
                :placeholder="$t('首页主机搜索提示语')"
                :rows="rows"
                v-model="searchContent"
                @focus="handleFocus"
                @blur="handleBlur"
                @keydown="handleKeydown">
            </bk-input>
            <bk-button theme="primary" class="search-btn"
                :loading="$loading(request.search)"
                @click="handleSearch">
                <i class="bk-icon icon-search"></i>
                {{$t('搜索')}}
            </bk-button>
            <span v-if="showEllipsis" class="search-text" @click="handleSearchInput">{{searchText}}</span>
        </div>
    </div>
</template>

<script>
    import { MENU_RESOURCE_HOST } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                rows: 1,
                searchText: '',
                searchContent: '',
                textarea: '',
                showEllipsis: false,
                textareaDom: null,
                request: {
                    search: Symbol('search')
                }
            }
        },
        watch: {
            searchContent () {
                this.$nextTick(this.setRows)
            }
        },
        mounted () {
            this.textareaDom = this.$refs.searchInput && this.$refs.searchInput.$refs.textarea
        },
        methods: {
            getSearchList () {
                const searchList = []
                this.searchContent.split('\n').forEach(text => {
                    const trimText = text.trim()
                    if (trimText.length) {
                        searchList.push(trimText)
                    }
                })
                return searchList
            },
            setRows () {
                const rows = this.searchContent.split('\n').length || 1
                this.rows = Math.min(10, rows)
            },
            handleFocus () {
                this.$emit('focus', true)
                this.setRows()
            },
            handleBlur () {
                this.textareaDom && this.textareaDom.blur()
                this.$emit('focus', false)
                const data = this.getSearchList()
                if (data.length) {
                    this.showEllipsis = true
                    this.searchText = data.join(',')
                } else {
                    this.searchContent = ''
                }
                this.$nextTick(() => {
                    this.rows = 1
                    this.textareaDom && (this.textareaDom.scrollTop = 0)
                })
            },
            handleKeydown (content, event) {
                const agent = window.navigator.userAgent.toLowerCase()
                const isMac = /macintosh|mac os x/i.test(agent)
                const modifierKey = isMac ? event.metaKey : event.ctrlKey
                if (modifierKey && event.code.toLowerCase() === 'enter') {
                    this.handleSearch()
                }
            },
            handleSearchInput () {
                this.showEllipsis = false
                this.textareaDom && this.textareaDom.focus()
            },
            async handleSearch () {
                const searchList = this.getSearchList()
                if (searchList.length > 500) {
                    this.$warn(this.$t('最多支持搜索500条数据'))
                } else if (searchList.length) {
                    try {
                        const validateQueue = searchList.map(text => this.$validator.verify(text, 'ip'))
                        const results = await Promise.all(validateQueue)
                        const ipList = []
                        const asstList = []
                        results.forEach((result, index) => {
                            (result.valid ? ipList : asstList).push(searchList[index])
                        })
                        if (ipList.length && asstList.length) {
                            this.$warn(this.$t('不支持混合搜索'))
                            return
                        } else if (ipList.length) {
                            this.navigateToResource(ipList, 'ip')
                        } else {
                            this.navigateToResource(asstList, 'bk_asset_id')
                        }
                    } catch (e) {
                        console.error(e)
                    }
                } else {
                    this.searchContent = ''
                    this.textareaDom && this.textareaDom.focus()
                }
            },
            navigateToResource (list, field) {
                this.$routerActions.redirect({
                    name: MENU_RESOURCE_HOST,
                    query: {
                        [field]: list.join(','),
                        scope: 'all',
                        exact: 1
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-search-layout {
        position: relative;
        width: 100%;
        max-width: 726px;
        height: 42px;
        margin: 0 auto;
    }
    .search-bar {
        position: absolute;
        width: 100%;
        height: 42px;
        z-index: 999;
        display: flex;
    }
    .search-input {
        flex: 1;
        /deep/ {
            .bk-textarea-wrapper {
                border: 0;
                border-radius: 2px 0 0 2px;
            }
            .bk-form-textarea {
                min-height: 42px;
                line-height: 30px;
                font-size: 14px;
                border: 1px solid #C4C6CC;
                padding: 5px 16px;
                border-radius: 2px 0 0 2px;
            }
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
    .search-text {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        max-width: 640px;
        height: 42px;
        line-height: 30px;
        font-size: 14px;
        color: #63656E;
        border: 1px solid #C4C6CC;
        background-color: #FFFFFF;
        padding: 5px 16px;
        z-index: 1;
        cursor: text;
        @include ellipsis;
    }
</style>
