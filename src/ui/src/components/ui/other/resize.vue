<template>
    <div class="resize-layout" :class="localDirections">
        <slot></slot>
        <i v-for="(direction, index) in localDirections"
            :class="['resize-handler', direction]"
            @mousedown.left="handleMousedown($event, direction)">
        </i>
        <i :class="['resize-proxy', state.direction]" ref="resizeProxy"></i>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-resize-layout',
        props: {
            direction: {
                default () {
                    return ['bottom', 'right']
                },
                validator (val) {
                    const validDirections = ['bottom', 'right']
                    if (typeof val === 'string') {
                        return validDirections.includes(val)
                    } else if (val instanceof Array) {
                        return !val.some(direction => !validDirections.includes(direction))
                    }
                    return false
                }
            },
            min: {
                default () {
                    return {
                        bottom: 0,
                        right: 0
                    }
                },
                validator (val) {
                    return ['object', 'number'].includes(typeof val)
                }
            },
            max: {
                default () {
                    return {
                        bottom: Infinity,
                        right: Infinity
                    }
                },
                validator (val) {
                    return ['object', 'number'].includes(typeof val)
                }
            }
        },
        data () {
            return {
                state: {}
            }
        },
        computed: {
            localDirections () {
                if (typeof this.direction === 'string') {
                    return [this.direction]
                }
                return this.direction
            },
            localMin () {
                const min = {
                    bottom: 0,
                    right: 0
                }
                if (typeof this.min === 'number') {
                    min.bottom = this.min
                    min.right = this.min
                } else {
                    Object.assign(min, this.min)
                }
                return min
            },
            localMax () {
                const max = {
                    bottom: Infinity,
                    right: Infinity
                }
                if (typeof this.max === 'number') {
                    max.bottom = this.max
                    max.right = this.max
                } else {
                    Object.assign(max, this.max)
                }
                return max
            }
        },
        methods: {
            handleMousedown (event, direction) {
                const $handler = event.currentTarget
                const handlerRect = $handler.getBoundingClientRect()
                const $container = this.$el
                const containerRect = $container.getBoundingClientRect()
                const $resizeProxy = this.$refs.resizeProxy
                if (direction === 'right') {
                    this.state = {
                        direction,
                        startMouseLeft: event.clientX,
                        startLeft: handlerRect.right - containerRect.left
                    }
                    $resizeProxy.style.top = 0
                    $resizeProxy.style.left = this.state.startLeft + 'px'
                } else {
                    this.state = {
                        direction,
                        startMouseTop: event.clientY,
                        startTop: handlerRect.bottom - containerRect.top
                    }
                    $resizeProxy.style.left = 0
                    $resizeProxy.style.top = this.state.startTop + 'px'
                }
                $resizeProxy.style.visibility = 'visible'
                document.onselectstart = () => { return false }
                document.ondragstart = () => { return false }
                const handleMouseMove = (event) => {
                    if (direction === 'right') {
                        const deltaLeft = event.clientX - this.state.startMouseLeft
                        const proxyLeft = this.state.startLeft + deltaLeft
                        const maxLeft = this.localMax.right
                        const minLeft = this.localMin.right
                        $resizeProxy.style.left = Math.min(maxLeft, Math.max(minLeft, proxyLeft)) + 'px'
                        document.body.style.cursor = 'col-resize'
                    } else {
                        const deltaTop = event.clientY - this.state.startMouseTop
                        const proxyTop = this.state.startTop + deltaTop
                        const maxTop = this.localMax.bottom
                        const minTop = this.localMin.bottom
                        $resizeProxy.style.top = Math.min(maxTop, Math.max(minTop, proxyTop)) + 'px'
                        document.body.style.cursor = 'row-resize'
                    }
                }
                const handleMouseUp = (event) => {
                    if (direction === 'right') {
                        const finalLeft = parseInt($resizeProxy.style.left, 10)
                        this.$el.style.width = finalLeft + 'px'
                    } else {
                        const finalTop = parseInt($resizeProxy.style.top, 10)
                        this.$el.style.height = finalTop + 'px'
                    }
                    document.body.style.cursor = 'auto'
                    $resizeProxy.style.visibility = 'hidden'
                    document.removeEventListener('mousemove', handleMouseMove)
                    document.removeEventListener('mouseup', handleMouseUp)
                    document.onselectstart = null
                    document.ondragstart = null
                }
                document.addEventListener('mousemove', handleMouseMove)
                document.addEventListener('mouseup', handleMouseUp)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .resize-layout {
        position: relative;
        &.right{
            padding-right: 5px;
        }
        &.bottom{
            padding-bottom: 5px;
        }
    }
    .resize-handler {
        position: absolute;
        background-color: transparent;
        &.right {
            top: 0;
            left: 100%;
            width: 5px;
            height: 100%;
            margin: 0 0 0 -5px;
            cursor: col-resize;
        }
        &.bottom {
            top: 100%;
            left: 0;
            width: 100%;
            height: 5px;
            margin: -5px 0 0 0;
            cursor: row-resize;
        }
    }
    .resize-proxy{
        visibility: hidden;
        position: absolute;
        pointer-events: none;
        z-index: 9999;
        &.right {
            top: 0;
            height: 100%;
            border-left: 1px dashed #d1d5e0;
        }
        &.bottom {
            left: 0;
            width: 100%;
            border-top: 1px dashed #d1d5e0;
        }
    }
</style>