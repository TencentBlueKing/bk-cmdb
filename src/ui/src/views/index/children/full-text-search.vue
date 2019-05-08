<template>
    <div class="full-text-search-layout">
        <div class="full-text-sticky-layout" ref="stickyLayout">
            <the-search class="search-bar"></the-search>
            <div class="classify">
                <span class="classify-item"
                    v-for="(model, index) in classify"
                    :key="index">
                    {{model.title}}（{{model.num}}）
                </span>
            </div>
        </div>
        <div class="results-wrapper" ref="resultsWrapper">
            <div class="results-list">
                <div class="results-item"
                    v-for="(item, index) in result"
                    :key="index">
                    <div class="results-title" v-html="item.title"></div>
                    <div class="results-desc">
                        <span class="desc-item"
                            v-for="(desc, childIndex) in item.desc"
                            :key="childIndex" v-html="desc.field">
                        </span>
                    </div>
                </div>
            </div>
        </div>
        <div class="pagination clearfix">
            <div class="pagination-info fl">
                <span class="mr10">{{$tc("Common['共计N条']", 100, { N: 100 })}}</span>
                <i18n path="Common['每页显示条数']">
                    <div ref="paginationSize" place="num"
                        :class="['pagination-size', { 'active': isShowSizeSetting }]"
                        @click="isShowSizeSetting = !isShowSizeSetting"
                        v-click-outside="closeSizeSetting">
                        <span>{{pagination.limit}}</span>
                        <i class="bk-icon icon-angle-down"></i>
                        <transition :name="transformOrigin === 'top' ? 'toggle-slide-top' : 'toggle-slide-bottom'">
                            <ul :class="['pagination-size-setting', transformOrigin === 'top' ? 'bottom' : 'top']"
                                v-show="sizeTransitionReady">
                                <li v-for="(limit, index) in [10, 20, 50, 100]"
                                    :key="index"
                                    :class="['size-setting-item', { 'selected': pagination.limit === limit }]"
                                    @click.stop="setSize(limit)">{{limit}}</li>
                            </ul>
                        </transition>
                    </div>
                </i18n>
            </div>
            <bk-paging class="pagination-list fr"
                :type="'compact'"
                :total-page="80"
                :size="'small'"
                :cur-page="1">
            </bk-paging>
        </div>
    </div>
</template>

<script>
    import theSearch from './search'
    import {
        addMainScrollListener,
        removeMainScrollListener
    } from '@/utils/main-scroller'
    export default {
        components: {
            theSearch
        },
        data () {
            return {
                scrollTop: null,
                isShowSizeSetting: false,
                sizeTransitionReady: false,
                transformOrigin: 'top',
                pagination: {
                    limit: 10,
                    start: 0
                },
                classify: [
                    {
                        title: '全部结果',
                        num: 100
                    },
                    {
                        title: '应用',
                        num: 33
                    },
                    {
                        title: '交换机',
                        num: 10
                    },
                    {
                        title: '负载均衡',
                        num: 6
                    },
                    {
                        title: 'Role',
                        num: 24
                    },
                    {
                        title: 'Tomcat',
                        num: 6
                    }
                ],
                result: [
                    {
                        title: '应用 - <em>oauth</em>111',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>2222',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    },
                    {
                        title: '应用 - <em>oauth</em>',
                        desc: [
                            {
                                field: '业务组：快点击放大看',
                                fieldDesc: ''
                            },
                            {
                                field: '应用名称：<em>CMP</em>1',
                                fieldDesc: ''
                            }
                        ]
                    }
                ]
            }
        },
        watch: {
            isShowSizeSetting (isShowSizeSetting) {
                if (isShowSizeSetting) {
                    this.calcSizePosition()
                } else {
                    this.sizeTransitionReady = false
                }
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Common["搜索结果"]'))
            for (let i = 0; i < 10; i++) {
                this.result.push({
                    title: '应用 - <em>oauth</em>',
                    desc: [
                        {
                            field: '业务组：快点击放大看',
                            fieldDesc: ''
                        },
                        {
                            field: '应用名称：<em>CMP</em>1',
                            fieldDesc: ''
                        }
                    ]
                })
                this.classify.push({
                    title: '自定义模型',
                    num: 100
                })
            }
        },
        mounted () {
            this.initScrollListener()
        },
        destroyed () {
            removeMainScrollListener(this.scrollTop)
        },
        methods: {
            initScrollListener () {
                this.scrollTop = event => {
                    const target = event.target
                    this.$refs.stickyLayout.style.top = target.scrollTop + 'px'
                }
                this.$refs.resultsWrapper.style.paddingTop = this.$refs.stickyLayout.offsetHeight + 'px'
                addMainScrollListener(this.scrollTop)
            },
            closeSizeSetting () {
                this.isShowSizeSetting = false
            },
            calcSizePosition () {
                if (this.pagination.hasOwnProperty('sizeDirection') && this.pagination.sizeDirection !== 'auto') {
                    this.transformOrigin = this.pagination.sizeDirection
                } else {
                    const sizeSettingItemHeight = 32
                    const sizeSettingHeight = [10, 20, 50, 100].length * sizeSettingItemHeight
                    const paginationSizeRect = this.$refs.paginationSize.getBoundingClientRect()
                    const bodyRect = document.body.getBoundingClientRect()
                    if (bodyRect.height - paginationSizeRect.y - paginationSizeRect.height > sizeSettingHeight) {
                        this.transformOrigin = 'bottom'
                    } else {
                        this.transformOrigin = 'top'
                    }
                }
                this.sizeTransitionReady = true
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
            border-bottom: 1px solid #dde4eb;
            position: absolute;
            top: 0;
            left: 50%;
            background: #ffffff;
            transform: translateX(-50%);
        }
        .search-bar {
            padding-top: 50px;
            padding-bottom: 38px;
        }
        .classify {
            color: $cmdbTextColor;
            font-size: 14px;
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
            .results-list {
                padding-top: 14px;
                color: $cmdbTextColor;
                .results-item {
                    padding-bottom: 20px;
                    em {
                        color: #3a84ff !important;
                        font-style: normal !important;
                    }
                    .results-title {
                        font-size: 16px;
                        font-weight: bold;
                        margin-bottom: 4px;
                    }
                    .results-desc {
                        font-size: 12px;
                        .desc-item {
                            padding-right: 16px;
                        }
                    }
                }
            }
        }
        .pagination {
            width: 90%;
            margin: 0 auto;
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
            z-index: 100;
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
    }
</style>
