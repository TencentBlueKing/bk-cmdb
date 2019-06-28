<template>
    <div v-if="totalPage > 0" :class="[
        'bk-page',
        {
            'bk-page-compact': type === 'compact',
            'bk-page-small': size === 'small'
        }
    ]">
        <template v-if="paginationAble && location === 'left'">
            <bk-pagination
                :total-page="totalPage"
                :pagination-count.sync="paginationCountTmp"
                :pagination-list="paginationList">
            </bk-pagination>
        </template>
        <ul>
            <!-- 上一页 -->
            <li class="page-item" :class="{disabled: curPage === 1}" @click="prevPage">
                <a href="javascript:void(0);" class="page-button">
                    <i class="bk-icon icon-angle-left"></i>
                </a>
            </li>
            <!-- 第一页 -->
            <li class="page-item" v-show="renderList[0] > 1" @click="jumpToPage(1)">
                <a href="javascript:void(0);" class="page-button">1</a>
            </li>
            <!-- 展开较前的页码 -->
            <li v-show="renderList[0] > 2 && curPage > 3"
                :class="[
                    'page-item',
                    {
                        'page-omit': type !== 'compact'
                    }
                ]" @click="prevGroup">
                <span class="page-button">...</span>
            </li>
            <!-- 渲染主要页码列表 -->
            <li class="page-item" v-for="item in renderList" :class="{'cur-page': item === curPage}" @click="jumpToPage(item)">
                <a href="javascript:void(0);" class="page-button">{{ item }}</a>
            </li>
            <!-- 展开较后的页码 -->
            <li v-show="renderList[renderList.length - 1] < calcTotalPage - 1"
                :class="[
                    'page-item',
                    {
                        'page-omit': type !== 'compact'
                    }
                ]" @click="nextGroup">
                <span class="page-button">...</span>
            </li>
            <!-- 最后一页 -->
            <li class="page-item"
                v-show="renderList[renderList.length - 1] !== calcTotalPage"
                :class="{
                    'cur-page': curPage === calcTotalPage
                }"
                @click="jumpToPage(calcTotalPage)">
                <a href="javascript:void(0);" class="page-button">{{ calcTotalPage }}</a>
            </li>
            <!-- 下一页 -->
            <li class="page-item" :class="{disabled: curPage === calcTotalPage}" @click="nextPage">
                <a href="javascript:void(0);" class="page-button">
                    <i class="bk-icon icon-angle-right"></i>
                </a>
            </li>
        </ul>
        <template v-if="paginationAble && location === 'right'">
            <bk-pagination
                :total-page="totalPage"
                :pagination-count.sync="paginationCountTmp"
                :pagination-list="paginationList">
            </bk-pagination>  
        </template>
    </div>
</template>
<script>
    /**
     *  bk-paging
     *  @module components/paging
     *  @desc 分页组件
     *  @param type {String} - 组件的类型，可选default和compact，默认为default
     *  @param size {String} - 组件的尺寸，可选default和small，默认为default
     *  @param totalPage {Number} - 总页数，默认为5
     *  @param paginationCount {Number} - 每页显示数据量，默认为10
     *  @param paginationAble {Boolean} - 开启自定义页码
     *  @param location {Number} - 自定义页码位置，默认为left
     *  @param paginationList {Array} - 自定义页码数组，默认为[10, 20, 50, 100]
     *  @param curPage {Number} - 当前页数，默认为1，支持.sync修饰符
     *  @event page-change {Event} - 当改变页码时，广播给父组件的事件
     *  @event pagination-change {Event} - 当改变页码数时，广播给父组件的事件
     *  @example
     *  <bk-paging
         :cur-page.sync="curPage"
         :total-Count="totalCount"
         :type="'compact'"></bk-paging>
    */
    import bkPagination from '../pagination/pagination.vue'
    export default {
        name: 'bk-paging',
        components: {
            bkPagination
        },
        props: {
            type: {
                type: String,
                default: 'default',
                validator (value) {
                    return [
                        'default',
                        'compact'
                    ].indexOf(value) > -1
                }
            },
            size: {
                type: String,
                default: 'default',
                validator (value) {
                    return [
                        'default',
                        'small'
                    ].indexOf(value) > -1
                }
            },
            curPage: {
                type: Number,
                default: 1,
                required: true
            },
            totalPage: {
                type: Number,
                default: 5,
                validator (value) {
                    return value >= 0
                }
            },
            paginationCount: {
                type: Number,
                default: 10,
                validator (value) {
                    return value >= 0
                }
            },
            paginationAble: {
                type: Boolean,
                default: false
            },
            location: {
                type: String,
                default: 'left',
                validator (value) {
                    return [
                        'left',
                        'right'
                    ].indexOf(value) > -1
                }
            },
            paginationList: {
                type: Array,
                default: () => [10, 20, 50, 100]
            }
        },
        data () {
            return {
                pageSize: 5,
                renderList: [],
                curGroup: 1,
                paginationCountTmp: this.paginationCount
            }
        },
        computed: {
            calcTotalPage () {
                if (this.totalPage >= this.curPage) {
                    return this.totalPage
                }
                // console.error(this.t('paging.error'))
                this.$emit('update:curPage', this.totalPage)
                return this.curPage
            }
        },
        created () {
            this.calcPageList(this.curPage)
        },
        watch: {
            curPage (newVal) {
                this.calcPageList(newVal)
            },
            totalPage (newVal) {
                this.calcPageList(this.curPage)
            },
            paginationCountTmp (newVal) {
                this.$emit('pagination-change', newVal)
            },
            paginationCount (newVal) {
                this.paginationCountTmp = newVal
            }
            // paginationIndex (newVal) {
            //     this.totalPage = Math.ceil(this.totalCount / newVal)
            //     this.calcPageList(this.curPage)
            // },
            // paginationCount (newVal) {
            //     if (this.paginationList.includes(newVal)) {
            //         this.paginationIndex = newVal
            //     } else {
            //         this.paginationIndex = this.paginationList[0]
            //     }
            // }
        },
        methods: {
            _array (size) {
                return Array.apply(null, {
                    length: size
                })
            },
            /**
             *  计算被渲染的页码数组
             *  @param curPage {number} 当前页码 or 下一次渲染的数组中的中间位置的页码数
             */
            calcPageList (curPage) {
                let total = this.calcTotalPage
                let pageSize = this.pageSize
                let size = pageSize > total ? total : pageSize

                if (curPage >= size - 1) { // 当前页码大于 pageSize
                    if (total - curPage > Math.floor(size / 2)) {
                        this.renderList = this._array(size).map((v, i) => i + curPage - Math.ceil(size / 2) + 1)
                    } else { // 当到了最后的 pageSize / 2 长度的页码时
                        this.renderList = this._array(size).map((v, i) => total - i).reverse()
                    }
                } else {
                    this.renderList = this._array(size).map((v, i) => i + 1)
                }
            },
            /**
             *  点击左侧...按钮
             */
            prevGroup () {
                let pageSize = this.pageSize
                let middlePage = this.renderList[Math.ceil(this.renderList.length / 2)]

                if (middlePage - pageSize < 1) {
                    this.calcPageList(1)
                } else {
                    this.calcPageList(middlePage - pageSize)
                }
                this.jumpToPage(this.renderList[Math.floor(this.renderList.length / 2)])
            },
            /**
             *  点击右侧...按钮
             */
            nextGroup () {
                let pageSize = this.pageSize
                let totalPage = this.calcTotalPage
                let middlePage = this.renderList[Math.ceil(this.renderList.length / 2)]

                if (middlePage + pageSize > totalPage) {
                    this.calcPageList(totalPage)
                } else {
                    this.calcPageList(middlePage + pageSize)
                }
                this.jumpToPage(this.renderList[Math.floor(this.renderList.length / 2)])
            },
            // 上一页
            prevPage () {
                if (this.curPage !== 1) {
                    this.jumpToPage(this.curPage - 1)
                }
            },
            // 下一页
            nextPage () {
                if (this.curPage !== this.calcTotalPage) {
                    this.jumpToPage(this.curPage + 1)
                }
            },
            jumpToPage (page) {
                this.$emit('update:curPage', page)
                this.$emit('page-change', page)
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/paging.scss';
</style>
