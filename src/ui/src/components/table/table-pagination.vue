<template>
    <div class="table-pagination clearfix">
        <div class="pagination-info fl">
            <span class="mr10" v-if="hasCheckbox">{{$tc("Common['已选N行']", checked.length, {N: checked.length})}}</span>
            <span class="mr10">{{$tc('Common[\'页码\']', pagination.current, {current: pagination.current, count:pagination.count, total: totalPage})}}</span>
            <i18n path="Common['每页显示']">
                <div ref="paginationSize" place="page"
                    :class="['pagination-size', {'active': isShowSizeSetting}]"
                    @click="isShowSizeSetting = !isShowSizeSetting"
                    v-click-outside="closeSizeSetting">
                    <span>{{pagination.size}}</span>
                    <i class="bk-icon icon-angle-down"></i>
                    <transition :name="transformOrigin === 'top' ? 'toggle-slide-top' : 'toggle-slide-bottom'">
                        <ul :class="['pagination-size-setting', transformOrigin === 'top' ? 'bottom' : 'top']"
                            v-show="sizeTransitionReady">
                            <li v-for="(size, index) in sizeSetting"
                                :key="index"
                                :class="['size-setting-item', {'selected': pagination.size === size}]"
                                @click.stop="setSize(size)">{{size}}</li>
                        </ul>
                    </transition>
                </div>
            </i18n>
        </div>
        <bk-paging class="pagination-list fr"
            v-if="pagination.count && pagination.count > pagination.size"
            :type="'compact'"
            :totalPage="totalPage"
            :curPage="pagination.current"
            @page-change="setCurrent">
        </bk-paging>
    </div>
</template>

<script>
    export default {
        props: {
            layout: Object
        },
        data () {
            return {
                isShowSizeSetting: false,
                sizeTransitionReady: false,
                transformOrigin: 'top'
            }
        },
        computed: {
            table () {
                return this.layout.table
            },
            pagination () {
                return this.table.pagination
            },
            checked () {
                return this.table.checked
            },
            totalPage () {
                return Math.ceil(this.pagination.count / this.pagination.size)
            },
            sizeSetting () {
                return this.pagination.sizeSetting || [10, 20, 50, 100]
            },
            hasCheckbox () {
                return this.table.header.filter(({type}) => type === 'checkbox').length
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
                if (this.pagination.hasOwnProperty('sizeDirection') && this.pagination.sizeDirection !== 'auto') {
                    this.transformOrigin = this.pagination.sizeDirection
                } else {
                    const sizeSettingItemHeight = 32
                    const sizeSettingHeight = this.sizeSetting.length * sizeSettingItemHeight
                    const paginationSizeRect = this.$refs.paginationSize.getBoundingClientRect()
                    const bodyRect = document.body.getBoundingClientRect()
                    if (bodyRect.height - paginationSizeRect.y - paginationSizeRect.height > sizeSettingHeight) {
                        this.transformOrigin = 'bottom'
                    } else {
                        this.transformOrigin = 'top'
                    }
                }
                this.sizeTransitionReady = true
            },
            setSize (size) {
                this.table.$emit('update:pagination', Object.assign({}, this.pagination, {size, current: 1}))
                this.table.$emit('handleSizeChange', size)
                this.closeSizeSetting()
            },
            setCurrent (current) {
                if (current !== this.pagination.current) {
                    if (!this.table.crossPageCheck) {
                        this.table.$emit('update:checked', [])
                    }
                    this.table.$emit('update:pagination', Object.assign({}, this.pagination, {current}))
                    this.table.$emit('handlePageChange', current)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-pagination{
        height: 42px;
        line-height: 42px;
        padding: 0 20px;
        background-color: #fafbfd;
        font-size: 12px;
        color: #c3cdd7;
        border-top: 1px solid $tableBorderColor;
    }
    .pagination-info{
        &:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
        }
    }
    .pagination-size{
        display: inline-block;
        vertical-align: middle;
        position: relative;
        width: 54px;
        height: 24px;
        line-height: 22px;
        border: 1px solid #c3cdd7;
        border-radius: 2px;
        padding: 0 25px 0 6px;
        background-color: #fff;
        color: $textColor;
        font-size: 14px;
        cursor: pointer;
        &.active{
            border-color: #0082ff;
            .icon-angle-down{
                color: #0082ff;
                transform: rotate(180deg);
            }
        }
        .icon-angle-down{
            position: absolute;
            right: 6px;
            top: 6px;
            font-size: 12px;
            color: #c3cdd7;
            transition: transform .2s linear;
        }
        .pagination-size-setting{
            position: absolute;
            left: 0;
            width: 100%;
            background-color: #fff;
            box-shadow: 0 0 1px 1px rgba(0, 0, 0, 0.1);
            z-index: 100;
            &.top {
                top: 25px;
            }
            &.bottom {
                bottom: 25px;
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
    }
    .pagination-list{
        height: 32px;
        margin: 5px 0;
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
<style lang="scss">
    .table-pagination{
        .pagination-list{
            .page-item{
                height: 32px;
                line-height: 32px;
            }
        }
    }
</style>