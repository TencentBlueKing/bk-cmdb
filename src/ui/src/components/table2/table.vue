<template>
    <div>
        <div :class="['table-wrapper', {'has-header': header.length}]" v-bkloading="{isLoading: isLoading}">
            <div ref="headWrapper" class="head-wrapper">
                <v-thead ref="head" :layout="layout"></v-thead>
            </div>
            <div ref="bodyWrapper" class="body-wrapper">
                <v-tbody ref="body" :layout="layout"></v-tbody>
            </div>
        </div>
        <div class="pagination-wrapper">
            <v-pagination></v-pagination>
        </div>
    </div>
</template>

<script>
    import vThead from './table-head'
    import vTbody from './table-body'
    import vPagination from './table-pagination'
    import TableLayout from './table-layout'
    import debounce from 'lodash.debounce'
    import throttle from 'lodash.throttle'
    import {addResizeListener, removeResizeListener} from '@/utils/resize-event.js'
    export default {
        props: {
            isLoading: {
                type: Boolean,
                default: false
            },
            header: {
                type: Array,
                default () {
                    return []
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
                        size: 10
                    }
                },
                validator (pagination) {
                    let keys = Object.keys(pagination)
                    return keys.includes('current') && keys.includes('count') && keys.includes('size')
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
            }

        },
        data () {
            const layout = new TableLayout({
                table: this
            })
            return {
                ready: false,
                layout,
                resizeListener: null
            }
        },
        watch: {
            header (header) {
                this.updateColumns()
            }
        },
        created () {
            this.updateColumns()
        },
        mounted () {
            this.init()
            this.layout.doLayout()
        },
        methods: {
            updateColumns () {
                let columns = []
                this.header.forEach(head => {
                    columns.push({
                        id: head[this.valueKey],
                        name: head[this.labelKey],
                        attr: head.attr || {},
                        type: head.type,
                        width: head.width || 70,
                        realWidth: head.width || 70,
                        flex: !head.width,
                        dragging: false
                    })
                })
                this.layout.columns = columns
            },
            init () {
                this.initScroll()
                this.initResize()
            },
            initScroll () {
                let {bodyWrapper, headWrapper} = this.$refs
                bodyWrapper.addEventListener('scroll', function () {
                    headWrapper.scrollLeft = this.scrollLeft
                })
            },
            initResize () {
                this.resizeListener = throttle(() => {
                    if (this.ready) {
                        this.layout.doLayout()
                    }
                }, 50)
                addResizeListener(this.$el, this.resizeListener)
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
        overflow: auto;
        border: 1px solid $tableBorderColor;
        @include scrollbar;
        position: relative;
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
        height: 300px;
        @include scrollbar;
    }
</style>