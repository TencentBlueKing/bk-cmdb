<template>
    <div v-bkloading="{isLoading: loading}"
        :class="['table-wrapper', {
            'has-header': hasHeader,
            'has-footer': hasFooter
        }]"
        :style="{height: bodyHeight}">
        <div class="table-layout">
            <div ref="headLayout" class="head-layout" v-if="hasHeader">
                <v-thead ref="head"
                    :layout="layout"
                    :style="{width: layout.bodyWidth ? layout.bodyWidth + 'px' : ''}">
                </v-thead>
            </div>
            <div ref="bodyLayout" class="body-layout" :style="{maxHeight: bodyLayoutMaxHeight}">
                <v-tbody ref="body"
                    :layout="layout"
                    :style="{width: bodyWidth}">
                </v-tbody>
            </div>
        </div>
        <div class="pagination-layout" v-if="showFooter && pagination.count > 0">
            <v-pagination :layout="layout"></v-pagination>
        </div>
    </div>
</template>

<script>
    import vThead from './table-head'
    import vTbody from './table-body'
    import vPagination from './table-pagination'
    import TableLayout from './table-layout'
    import throttle from 'lodash.throttle'
    import {addResizeListener, removeResizeListener} from '@/utils/resize-event.js'
    export default {
        props: {
            loading: {
                type: Boolean,
                default: false
            },
            header: {
                type: Array,
                default () {
                    return []
                },
                validator (header) {
                    const checkboxHeader = header.filter(({type}, index) => index === 0 && type === 'checkbox')
                    return checkboxHeader.length <= 1
                }
            },
            list: {
                type: Array,
                default () {
                    return []
                }
            },
            pagination: {
                type: Object,
                default () {
                    return {
                        current: 1,
                        count: 0,
                        size: 10,
                        sizeSetting: [10, 20, 50, 100]
                    }
                },
                validator (pagination) {
                    let necessaryKeys = ['current', 'count', 'size']
                    let invalid = necessaryKeys.some(key => typeof pagination[key] !== 'number')
                    if (pagination.hasOwnProperty('sizeSetting') && !Array.isArray(pagination.sizeSetting)) {
                        invalid = true
                    }
                    return !invalid
                }
            },
            showHeader: {
                type: Boolean,
                default: true
            },
            showFooter: {
                type: Boolean,
                default: true
            },
            sortable: {
                type: Boolean,
                default: true
            },
            defaultSort: {
                type: String,
                default: ''
            },
            checked: {
                type: Array,
                default () {
                    return []
                }
            },
            crossPageCheck: {
                type: Boolean,
                default: true
            },
            multipleCheck: {
                type: Boolean,
                default: true
            },
            valueKey: {
                type: String,
                default: 'id'
            },
            labelKey: {
                type: String,
                default: 'name'
            },
            gutterWidth: {
                type: Number,
                default: 8
            },
            rowBorder: {
                type: Boolean,
                default: false
            },
            colBorder: {
                type: Boolean,
                default: false
            },
            stripe: {
                type: Boolean,
                default: false
            },
            height: {
                type: [String, Number],
                validator (height) {
                    height = parseInt(height, 10)
                    return height && height > 122 // 40(header) + 40(1row) + 42(pagination)
                }
            },
            maxHeight: {
                type: [String, Number],
                validator (maxHeight) {
                    maxHeight = parseInt(maxHeight, 10)
                    return maxHeight && maxHeight > 122 // 40(header) + 40(1row) + 42(pagination)
                }
            }

        },
        data () {
            const layout = new TableLayout({
                table: this
            })
            return {
                layout,
                throttleLayout: null
            }
        },
        computed: {
            bodyWidth () {
                const { bodyWidth, scrollY } = this.layout
                return bodyWidth ? bodyWidth - (scrollY ? this.gutterWidth : 0) + 'px' : ''
            },
            hasHeader () {
                return this.showHeader && this.layout.columns.length
            },
            hasFooter () {
                return this.showFooter && this.pagination.count > 0
            },
            finalHeight () {
                const height = parseInt(this.height, 10)
                const maxHeight = parseInt(this.maxHeight, 10)
                if (height && maxHeight) {
                    return Math.min(height, maxHeight)
                } else if (maxHeight) {
                    return maxHeight
                } else if (height) {
                    return height
                }
                return null
            },
            bodyHeight () {
                if (!this.hasHeader && !this.hasFooter && !this.list.length) {
                    return this.finalHeight ? this.finalHeight + 'px' : '302px'
                }
                return 'auto'
            },
            bodyLayoutMaxHeight () {
                const bodyLayout = this.$refs.bodyLayout
                if (this.finalHeight) {
                    let minusHeight = 0
                    minusHeight += this.hasHeader ? 40 : 0
                    minusHeight += this.hasFooter ? 42 : 0
                    return this.finalHeight - minusHeight + 'px'
                }
                return 'none'
            }
        },
        watch: {
            header (header) {
                this.layout.updateColumns()
                this.throttleLayout()
            },
            list (list) {
                this.throttleLayout()
            },
            height (height) {
                this.throttleLayout()
            },
            maxHeight (maxHeight) {
                this.throttleLayout()
            }
        },
        created () {
            this.layout.updateColumns()
        },
        mounted () {
            this.bindEvents()
            this.doLayout()
        },
        beforeDestroy () {
            if (this.throttleLayout) {
                removeResizeListener(this.$el, this.throttleLayout)
            }
            this.$emit('update:checked', [])
            this.$emit('update:pagination', {
                current: 1,
                count: 0,
                size: 10,
                sizeSetting: [10, 20, 50, 100]
            })
        },
        methods: {
            bindEvents () {
                this.initScroll()
                this.initResize()
            },
            doLayout () {
                this.$nextTick(() => {
                    this.layout.checkScrollY()
                    this.layout.update()
                })
            },
            initScroll () {
                let {bodyLayout, headLayout} = this.$refs
                bodyLayout.addEventListener('scroll', function () {
                    headLayout ? headLayout.scrollLeft = this.scrollLeft : void 0
                })
            },
            initResize () {
                this.throttleLayout = throttle(() => {
                    this.doLayout()
                }, 100)
                addResizeListener(this.$el, this.throttleLayout)
            }
        },
        components: {
            vThead,
            vTbody,
            vPagination
        }
    }
</script>

<style lang="scss" scoped>
    .table-wrapper{
        border: 1px solid $tableBorderColor;
        &.has-header{
            border-top: none;
        }
    }
    .table-layout{
        position: relative;
        overflow: auto;
        @include scrollbar;
    }
    .head-layout{
        overflow: hidden;
    }
    .body-layout{
        overflow: auto;
        border-top: none;
        @include scrollbar;
    }
</style>