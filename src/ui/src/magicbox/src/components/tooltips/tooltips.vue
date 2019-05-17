<template>
    <transition name="v-tooltips-fade">
        <div v-show="visible" class="v-tooltips-container"
            :style="boxStyle"
            :class="boxClass"
            @mouseenter="showTooltips"
            @mouseleave="hiddenTooltips(true)">
            <div v-show="placement" class="v-tooltips-arrows" :class="placement" :style="arrowBox"></div>
            <span v-if="title" class="v-tooltips-title">{{title}}</span>
            <p v-if="content" class="v-tooltips-content" :style="contentHeight">
                {{content}}
            </p>
            <component v-if="customComponent" v-bind="customProps" v-on="customListeners" :is="customComponent"
                @hidden-tooltips="hiddenTooltips" @update-tooltips="updateTooltips"></component>
        </div>
    </transition>
</template>

<script>
    import {
        debounce,
        checkScrollable,
        getScrollContainer,
        computeArrowPos,
        computePlacementInfo,
        computeCoordinateBaseMid,
        computeCoordinateBaseEdge
    } from '../../util'

    // passive support check
    let supportsPassive = false
    document.addEventListener('passive-check', () => {}, {
        get passive () {
            supportsPassive = {
                passive: true
            }
        }
    })

    export default {
        name: 'bk-tooltips',
        props: {
            // 标题
            title: {
                type: String,
                default: ''
            },
            // 显示的内容
            content: {
                type: String,
                default: ''
            },
            // 工具函数调用时附加到自定义组件 props 上面的
            customProps: {
                type: Object,
                default () {
                    return {}
                }
            },
            // 对应 <component> 组件 is 属性
            customComponent: {
                type: [String, Function, Object],
                default: ''
            },
            // 用于监听自定义组件 emit 的事件
            customListeners: Object,
            // tip 绑定的目标元素
            target: null,
            // tip 的容器，默认插入 body 中
            container: null,
            // 用于限制 tip 展示的方向，优先级按顺序
            placements: {
                type: Array,
                default () {
                    return ['top', 'right', 'bottom', 'left']
                }
            },
            // tip 窗口多久后自动消失，<= -1 时不消失
            duration: {
                type: Number,
                default: 300
            },
            // 提示用的小箭头大小
            arrowsSize: {
                type: Number,
                default: 8
            },
            // 组件的宽度
            width: {
                type: [String, Number],
                default: 'auto'
            },
            // 内容的高度
            height: {
                type: [String, Number],
                default: 'auto'
            },
            // tip 的 z-index
            zIndex: {
                type: Number,
                default: 999
            },
            // 主题 dark light 默认为 dark
            theme: {
                type: String,
                default: 'dark'
            },
            // 自定义 class 的类名
            customClass: {
                type: String,
                default: ''
            },
            // 显示的回调函数
            onShow: {
                type: Function,
                default () {}
            },
            // 隐藏的回调函数
            onClose: {
                type: Function,
                default () {}
            }
        },
        data () {
            return {
                // tip 的展示方向（小箭头的方向）
                placement: '',
                visible: false,
                arrowsPos: {},
                containerNode: null,
                targetParentNode: null,
                visibleTimer: null
            }
        },
        computed: {
            arrowBox () {
                return Object.assign({borderWidth: `${this.arrowsSize}px`}, this.arrowsPos)
            },
            boxStyle () {
                const width = this.width
                return {
                    width: typeof width === 'string' ? width : `${width}px`,
                    zIndex: this.zIndex
                }
            },
            boxClass () {
                const {customClass, theme} = this
                return [customClass, theme]
            },
            contentHeight () {
                const height = this.height
                return {
                    height: typeof height === 'string' ? height : `${height}px`
                }
            }
        },
        watch: {
            visible (val) {
                if (val) {
                    this.onShow && this.onShow(this)
                }
                else {
                    this.onClose && this.onClose(this)
                }
            }
        },
        methods: {
            /**
             * 显示 tooltips
             */
            showTooltips () {
                clearTimeout(this.visibleTimer)
                this.visible = true
            },

            /**
             * 隐藏 tooltips
             *
             * @param {boolean} immediate 是否立即隐藏
             */
            hiddenTooltips (immediate) {
                // duration 设置为 <= -1 即不关闭
                if (this.duration <= -1) {
                    return
                }
                if (immediate) {
                    this.visible = false
                } else {
                    this.setVisible(false)
                }
            },

            /**
             * 更新 tooltips 位置
             */
            updateTooltips () {
                this.setContainerNode()
                this.showTooltips()
                this.$nextTick(this.setPosition)
            },

            /**
             * 设置 tooltips 的容器
             */
            setContainerNode () {
                const {$el, target, container, targetParentNode, containerNode: oldNode} = this

                // 目标元素的父级节点相同则不需要重新计算容器
                if (!target || target.parentNode === targetParentNode) {
                    return
                }

                this.targetParentNode = target.parentNode

                const newNode = container || getScrollContainer(target)
                if (newNode === oldNode) {
                    return
                }

                if ($el.parentNode !== newNode) {
                    newNode.appendChild($el)
                }

                const position = window.getComputedStyle(newNode, null).position
                if (!position || position === 'static') {
                    newNode.style.position = 'relative'
                }
                if (oldNode) {
                    oldNode.removeEventListener('scroll', this.scrollHandler, supportsPassive)
                }

                if (checkScrollable(newNode)) {
                    newNode.addEventListener('scroll', this.scrollHandler, supportsPassive)
                }
                this.containerNode = newNode
            },

            /**
             * 设置 tooltips 位置
             */
            setPosition () {
                const {$el, target, containerNode, placements, arrowsSize} = this
                if (!$el || !target || !containerNode) {
                    return
                }
                const placementInfo = computePlacementInfo(target, containerNode, $el, placements, arrowsSize)
                const coordinate = placementInfo.mod === 'mid'
                    ? computeCoordinateBaseMid(placementInfo, arrowsSize)
                    : computeCoordinateBaseEdge(placementInfo, arrowsSize)

                this.setArrowsPos(coordinate)
                this.placement = coordinate.placement

                const x = Math.round(coordinate.x + containerNode.scrollLeft)
                const y = Math.round(coordinate.y + containerNode.scrollTop)
                this.$el.style.transform = `translate3d(${x}px, ${y}px, 0)`
            },

            /**
             * 设置 tooltips 小三角形的位置
             */
            setArrowsPos ({placement, arrowsOffset}) {
                this.arrowsPos = computeArrowPos(placement, arrowsOffset, this.arrowsSize)
            },

            /**
             * 设置 tooltips 经过 duration ms 后的状态
             */
            setVisible (v) {
                clearTimeout(this.visibleTimer)
                this.visibleTimer = setTimeout(() => {
                    this.visible = v
                    this.visibleTimer = null
                }, this.duration)
            },

            /**
             * 元素父级容器发生滚动时的处理
             */
            scrollHandler: debounce(function () {
                this.setPosition()
            }, 200, true),

            /**
             * 清除 scroll 事件监听
             */
            clearScrollEvent () {
                if (this.containerNode) {
                    this.containerNode.removeEventListener('scroll', this.scrollHandler, supportsPassive)
                }
            },

            /**
             * 移除节点
             */
            removeParentNode () {
                if (this.$el.parentNode) {
                    this.$el.parentNode.removeChild(this.$el)
                }
            },

            /**
             * 销毁
             */
            destroy () {
                this.clearScrollEvent()
                this.removeParentNode()
                this.$destroy()
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/tooltips.scss'
</style>
