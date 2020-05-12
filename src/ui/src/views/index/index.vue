<template>
    <div class="index-layout">
        <div class="search-layout" :style="{ paddingTop: inSearchPaddingTop + 'px' }">
            <div :class="['search-tab', { 'is-focus': isFocus }]" v-show="showSearchTab" v-if="isFullTextSearch">
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
                <search-input v-show="activeName === 'fullText'"
                    v-if="isFullTextSearch"
                    :is-full-text-search="true"
                    @search-status="handleSearchStatus"
                    @focus="handleFocus">
                </search-input>
            </div>
        </div>
        <the-map style="user-select: none;"></the-map>
        <the-footer></the-footer>
    </div>
</template>

<script>
    import hostSearch from './children/host-search'
    import searchInput from './children/search-input'
    import theMap from './children/map'
    import theFooter from './children/footer'
    import { mapGetters } from 'vuex'
    export default {
        name: 'index',
        components: {
            hostSearch,
            searchInput,
            theMap,
            theFooter
        },
        data () {
            return {
                activeName: 'host',
                inSearchPaddingTop: null,
                showSearchTab: true,
                isFocus: false
            }
        },
        computed: {
            ...mapGetters(['site']),
            isFullTextSearch () {
                return this.site.fullTextSearch === 'on'
            },
            paddingTop () {
                return parseInt((this.$APP.height - 58) / 3, 10)
            }
        },
        created () {
            this.inSearchPaddingTop = this.paddingTop
            const query = this.$route.query
            const showFullText = ['keywords', 'show'].every(key => query.hasOwnProperty(key))
            if (showFullText && this.isFullTextSearch) {
                this.activeName = 'fullText'
            }
        },
        methods: {
            handleChangeTab (name) {
                this.activeName = name
            },
            handleSearchStatus (status) {
                this.inSearchPaddingTop = status ? 0 : this.paddingTop
                this.showSearchTab = !status
            },
            handleFocus (status) {
                this.isFocus = status
            }
        }
    }
</script>

<style lang="scss" scoped>
    .index-layout {
        padding: 0 0 50px;
        background-color: #F5F6FA;
        position: relative;
        z-index: 1;
    }
    .search-layout {
        height: 100%;
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
</style>
