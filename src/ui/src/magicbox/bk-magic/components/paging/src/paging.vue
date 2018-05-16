<template>
    <div :class="['bk-page', {'bk-page-compact': type === 'compact', 'bk-page-small': size === 'small'}]" v-if="totalPage > 0">
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
            <li :class="['page-item',{'page-omit': type !== 'compact'}]" @click="prevGroup" v-show="renderList[0] > 1">
                <span class="page-button">...</span>
            </li>
            <!-- 渲染主要页码列表 -->
            <li class="page-item" :class="{'cur-page': item === curPage}" v-for="item in renderList" @click="jumpToPage(item)">
                <a href="javascript:void(0);" class="page-button">{{ item }}</a>
            </li>
            <!-- 展开较后的页码 -->
            <li :class="['page-item',{'page-omit': type !== 'compact'}]" @click="nextGroup" v-show="renderList[renderList.length - 1] < calcTotalPage - 1">
                <span class="page-button">...</span>
            </li>
            <!-- 最后一页 -->
            <li class="page-item" :class="{'cur-page': curPage === calcTotalPage}" @click="jumpToPage(calcTotalPage)" v-show="renderList[renderList.length - 1] !== calcTotalPage">
                <a href="javascript:void(0);" class="page-button">{{ calcTotalPage }}</a>
            </li>
            <!-- 下一页 -->
            <li class="page-item" :class="{disabled: curPage === calcTotalPage}" @click="nextPage">
                <a href="javascript:void(0);" class="page-button">
                    <i class="bk-icon icon-angle-right"></i>
                </a>
            </li>
        </ul>
    </div>
</template>
<script>
    /**
     *  bk-paging
     *  @module components/paging
     *  @desc 分页组件
     *  @param type {String} - 组件的类型，可选default和compact，默认为default
     *  @param size {String} - 组件的尺寸，可选default和small，默认为default
     *  @param totalPage {Number} - 总页数，默认为20
     *  @param curPage {Number} - 当前页数，默认为1，支持.sync修饰符
     *  @event page-change {Event} - 当改变页码时，广播给父组件的事件
     *  @example
     *  <bk-paging
         :cur-page.sync="curPage"
         :total-page="totalPage"
         :type="'compact'"></bk-paging>
    */
    export default {
        name: 'bk-paging',
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
            }
        },
        data () {
            return {
                pageSize: 5,
                renderList: [],
                curGroup: 1
            }
        },
        computed: {
            calcTotalPage () {
                if (this.totalPage >= this.curPage) {
                    return this.totalPage
                } else {
                    console.error(this.t('paging.error'))
                    return this.curPage
                }
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
            }
        },
        methods: {
            _array (size) {
                return Array.apply(null, {
                    length: size
                })
            },
            /**
             *  计算被渲染的页码数组
             *  @param curPage {Number} 当前页码or下一次渲染的数组中的中间位置的页码数
             */
            calcPageList (curPage) {
                let total = this.calcTotalPage
                let pageSize = this.pageSize
                let size = pageSize > total ? total : pageSize

                if (curPage >= size - 1) { // 当前页码大于pageSize
                    if (total - curPage > Math.floor(size / 2)) {
                        this.renderList = this._array(size).map((v, i) => i + curPage - Math.ceil(size / 2) + 1)
                    } else { // 当到了最后的pageSize / 2长度的页码时
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
                this.$emit('update:curPage', this.renderList[Math.ceil(this.renderList.length / 2)])
                this.$emit('page-change', this.renderList[Math.ceil(this.renderList.length / 2)])
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
                this.$emit('update:curPage', this.renderList[Math.ceil(this.renderList.length / 2)])
                this.$emit('page-change', this.renderList[Math.ceil(this.renderList.length / 2)])
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
