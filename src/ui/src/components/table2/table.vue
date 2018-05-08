<template>
    <div class="cc-table" v-bkloading="{isLoading: loading}" v-show="header.length">
        <div :class="['table-wrapper', {'has-header': header.length}]">
            <div ref="headWrapper" class="head-wrapper">
                <v-thead ref="head" :layout="layout"></v-thead>
            </div>
            <div ref="bodyWrapper" class="body-wrapper">
                <v-tbody ref="body" :layout="layout"></v-tbody>
            </div>
        </div>
        <div class="pagination-wrapper">
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
                    const checkboxHeader = header.filter(({type}) => type === 'checkbox')
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
                    let invalid = necessaryKeys.some(key => !pagination.hasOwnProperty(key))
                    if (pagination.hasOwnProperty('sizeSetting') && !Array.isArray(pagination.sizeSetting)) {
                        invalid = true
                    }
                    return !invalid
                }
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
                default: 6
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
        watch: {
            header (header) {
                this.updateColumns()
                this.throttleLayout()
            },
            list (list) {
                this.throttleLayout()
            },
            height (height) {
                this.initHeight()
                this.throttleLayout()
            },
            maxHeight (maxHeight) {
                this.initHeight()
                this.throttleLayout()
            }
        },
        created () {
            this.updateColumns()
        },
        mounted () {
            this.init()
        },
        beforeDestroy () {
            if (this.throttleLayout) {
                removeResizeListener(this.$el, this.throttleLayout)
            }
        },
        methods: {
            updateColumns () {
                let columns = []
                this.header.forEach(head => {
                    columns.push({
                        head: head,
                        id: head[this.valueKey],
                        name: head[this.labelKey],
                        attr: head.attr || {},
                        type: head.type || 'text',
                        width: head.width || 100,
                        realWidth: head.width || 100,
                        flex: !head.width,
                        sortable: head.hasOwnProperty('sortable') ? head.sortable : true,
                        sortKey: head.hasOwnProperty('sortKey') ? head.sortKey : head[this.valueKey],
                        dragging: false
                    })
                })
                this.layout.columns = columns
            },
            init () {
                this.initHeight()
                this.initScroll()
                this.initResize()
            },
            initHeight () {
                const height = parseInt(this.height, 10)
                const maxHeight = parseInt(this.maxHeight, 10)
                this.$el.style.height = height ? height + 'px' : 'auto'
                this.$el.style.maxHeight = maxHeight ? maxHeight + 'px' : 'none'

                let referenceHeight = height
                if (height && maxHeight) {
                    referenceHeight = Math.min(height, maxHeight)
                } else if (maxHeight) {
                    referenceHeight = maxHeight
                }
                this.$refs.bodyWrapper.style.maxHeight = referenceHeight ? referenceHeight - 82 + 'px' : 'none'
            },
            initScroll () {
                let {bodyWrapper, headWrapper} = this.$refs
                bodyWrapper.addEventListener('scroll', function () {
                    headWrapper.scrollLeft = this.scrollLeft
                })
            },
            initResize () {
                this.throttleLayout = throttle(() => {
                    this.layout.doLayout()
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
        position: relative;
        overflow: auto;
        border: 1px solid $tableBorderColor;
        @include scrollbar;
        &.has-header{
            border-top: none;
        }
    }
    .head-wrapper{
        overflow: hidden;
    }
    .body-wrapper{
        overflow: auto;
        border-top: none;
        @include scrollbar;
    }
</style>