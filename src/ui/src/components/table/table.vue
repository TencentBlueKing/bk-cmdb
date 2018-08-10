<template>
    <div v-bkloading="{isLoading: loading}"
        :class="['table-wrapper', {
            'has-header': hasHeader,
            'has-footer': hasFooter
        }]"
        :style="{height: wrapperHeight, width: wrapperWidth}">
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
                        sizeSetting: [10, 20, 50, 100],
                        sizeDirection: 'auto'
                    }
                },
                validator (pagination) {
                    let necessaryKeys = ['current', 'count', 'size']
                    let invalid = necessaryKeys.some(key => typeof pagination[key] !== 'number')
                    if (pagination.hasOwnProperty('sizeSetting') && !Array.isArray(pagination.sizeSetting)) {
                        invalid = true
                    }
                    if (pagination.hasOwnProperty('sizeDirection') && !['top', 'bottom', 'auto'].includes(pagination.sizeDirection)) {
                        invalid = true
                    }
                    return !invalid
                }
            },
            renderCell: {
                type: Function
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
                default: true
            },
            colBorder: {
                type: Boolean,
                default: false
            },
            rowHoverColor: {
                type: [String, Boolean],
                default: '#f1f7ff'
            },
            stripe: {
                type: Boolean,
                default: false
            },
            height: {
                type: Number,
                validator (height) {
                    return height > 122 // 40(header) + 40(1row) + 42(pagination)
                }
            },
            maxHeight: {
                type: Number,
                validator (maxHeight) {
                    return maxHeight > 122 // 40(header) + 40(1row) + 42(pagination)
                }
            },
            emptyHeight: {
                type: Number,
                default: 302,
                validator (emptyHeight) {
                    return !!emptyHeight
                }
            },
            wrapperMinusHeight: {
                type: Number,
                default: 0
            },
            width: {
                type: Number,
                validator (width) {
                    return width > 102
                }
            },
            visible: {
                type: Boolean,
                default: true
            },
            selectedList: {
                type: Array
            },
            isCheckboxShow: {
                type: Boolean,
                default: true
            }
        },
        data () {
            const layout = new TableLayout({
                table: this
            })
            return {
                layout,
                wrapperMaxHeight: 0,
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
                const height = this.height
                const maxHeight = this.maxHeight
                if (height && maxHeight) {
                    return Math.min(height, maxHeight)
                } else if (maxHeight) {
                    return maxHeight
                } else if (height) {
                    return height
                } else {
                    return this.wrapperMaxHeight
                }
            },
            wrapperHeight () {
                if (!this.hasHeader && !this.hasFooter && !this.list.length) {
                    return this.finalHeight ? this.finalHeight + 'px' : '302px'
                }
                return 'auto'
            },
            wrapperWidth () {
                return this.width ? this.width + 'px' : 'auto'
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
            },
            visible (visible) {
                if (visible) {
                    this.throttleLayout()
                }
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
            updateWrapperMaxHeight () {
                const bodyRect = document.body.getBoundingClientRect()
                const wrapperRect = this.$el.getBoundingClientRect()
                this.wrapperMaxHeight = bodyRect.height - (this.wrapperMinusHeight ? this.wrapperMinusHeight : wrapperRect.top)
            },
            initScroll () {
                const $refs = this.$refs
                $refs.bodyLayout.addEventListener('scroll', function () {
                    $refs.headLayout ? $refs.headLayout.scrollLeft = this.scrollLeft : void 0
                })
            },
            initResize () {
                this.throttleLayout = throttle(() => {
                    this.doLayout()
                    this.updateWrapperMaxHeight()
                }, 50)
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