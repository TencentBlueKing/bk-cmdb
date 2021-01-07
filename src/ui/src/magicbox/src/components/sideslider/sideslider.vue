<template>
    <transition name="slide">
        <article class="bk-sideslider" v-show="isShow" @click.self="handleQuickClose">
            <section class="bk-sideslider-wrapper" :class="[{left: direction === 'left', right: direction === 'right'}]" :style="{width: width + 'px'}">
                <div class="bk-sideslider-header">
                    <div class="bk-sideslider-closer" @click="handleClose" :style="{float: calcDirection}">
                        <i class="bk-icon" :class="'icon-angle-' + direction"></i>
                    </div>
                    <div class="bk-sideslider-title" :style="{padding: this.calcDirection === 'left' ? '0 0 0 50px' : '0 50px 0 0', 'text-align': this.calcDirection}">
                        {{title || '标题'}}
                    </div>
                </div>
                <div class="bk-sideslider-content">
                    <slot name="content"></slot>
                </div>
            </section>
        </article>
    </transition>
</template>
<script>
    /**
     *  bk-sideslider
     *  @module components/side-slider
     *  @desc 滑动侧边栏组件
     *  @param isShow {Boolean} - 是否显示组件，默认为false；支持.sync修饰符
     *  @param title {String} - 组件标题
     *  @param quickClose {Boolean} - 是否支持点击遮罩关闭组件，默认为false
     *  @param width {String} - 组件的宽度，支持5% ~ 100%，步长为5%的值，默认为40%
     *  @param direction {String} - 组件滑出的方向，可选left，right；默认为right
     *  @example
     *  <bk-sideslider
          :is-show.sync="isShow"
          :title="'测试标题'"
          :width="'80%'"
          :direction="'left'">
        </bk-sideslider>
     */
    import {
        addClass,
        removeClass
    } from '../../util.js'

    export default {
        name: 'bk-sideslider',
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            title: {
                type: String,
                default: ''
            },
            quickClose: {
                type: Boolean,
                default: false
            },
            width: {
                default: 400
            },
            direction: {
                type: String,
                default: 'right',
                validator (value) {
                    return [
                        'left',
                        'right'
                    ].indexOf(value) > -1
                }
            }
        },
        watch: {
            isShow (val) {
                // 有动画效果，因此推迟发布事件
                let root = document.documentElement
                if (val) {
                    addClass(root, 'bk-sideslider-show')
                    if (this.isScrollY()) {
                        addClass(root, 'has-sideslider-padding')
                    }
                    setTimeout(() => {
                        this.$emit('shown')
                    }, 200)
                } else {
                    removeClass(root, 'bk-sideslider-show has-sideslider-padding')
                    setTimeout(() => {
                        this.$emit('hidden')
                    }, 200)
                }
            }
        },
        computed: {
            calcDirection () {
                return this.direction === 'left' ? 'right' : 'left'
            }
        },
        methods: {
            isScrollY () {
                return document.documentElement.offsetHeight > document.documentElement.clientHeight
            },
            show () {
                // 监听isShow，显示组件时，将页面垂直滚动条隐藏，关闭组件时再恢复
                let root = document.documentElement
                addClass(root, 'bk-sideslider-show')
                this.isShow = true
            },
            hide () {
                let root = document.querySelector('html')
                removeClass(root, 'bk-sideslider-show')
                this.isShow = false
            },
            handleClose () {
                this.$emit('update:isShow', false)
            },
            handleQuickClose () {
                if (this.quickClose) {
                    this.handleClose()
                }
            }
        },
        destroyed () {
            let root = document.querySelector('html')
            removeClass(root, 'bk-sideslider-show')
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/sideslider.scss'
</style>
