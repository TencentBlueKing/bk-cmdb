<template>
    <div class="pagination clearfix">
        <div class="pagination-info fl">
            <span class="mr10">{{$tc('共计N条', pagination['total'], { N: pagination['total'] })}}</span>
            <i18n path="每页显示条数">
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
            v-if="pageCount > 1"
            :type="'compact'"
            :total-page="pageCount"
            :size="'small'"
            :cur-page="pagination['current']"
            @page-change="handlePageChange">
        </bk-paging>
    </div>
</template>

<script>
    export default {
        props: {
            pagination: {
                type: Object,
                default: () => {
                    return {
                        limit: 10,
                        start: 0,
                        current: 1,
                        total: 0
                    }
                }
            }
        },
        data () {
            return {
                transformOrigin: 'top',
                isShowSizeSetting: false,
                sizeTransitionReady: false,
                limitList: [10, 20, 50, 100]
            }
        },
        computed: {
            pageCount () {
                return Math.ceil(this.pagination.total / this.pagination.limit)
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
        methods: {
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
            handleLimitChange (limit) {
                this.$emit('limit-change', limit)
            },
            handlePageChange (current) {
                this.$emit('page-change', current)
            }
            // setPagination (total) {
            //     this.pagination.total = total
            //     this.pagination.count = Math.ceil(total / this.pagination.limit)
            // }
        }
    }
</script>

<style lang="scss" scoped>
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
</style>
