<template>
    <div class="map">
        <img src="../../../assets/images/map.svg" :style="imgStyles">
    </div>
</template>

<script>
    import {
        addResizeListener,
        removeResizeListener
    } from '@/utils/resize-events.js'
    export default {
        name: 'cmdb-index-map',
        data () {
            return {
                ratio: {
                    height: 404 / 857,
                    top: 181 / 767
                },
                imgStyles: {
                    width: '857px',
                    height: '404px',
                    visibility: 'hidden'
                },
                resizeHandler: null
            }
        },
        mounted () {
            this.initResizeEvent()
        },
        beforeDestroy () {
            removeResizeListener(this.$parent.$el, this.resizeHandler)
        },
        methods: {
            initResizeEvent () {
                this.resizeHandler = () => {
                    const parentRect = this.$parent.$el.getBoundingClientRect()
                    const imgStyles = {
                        width: Math.floor(parentRect.width * 0.66) + 'px',
                        height: Math.floor(parentRect.width * 0.66 * this.ratio.height) + 'px',
                        left: Math.floor(parentRect.width * 0.17 + parentRect.left) + 'px'
                    }
                    this.imgStyles = imgStyles
                }
                addResizeListener(this.$parent.$el, this.resizeHandler)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .map {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: -1;
        pointer-events: none;
        overflow: visible;
        svg, img {
            position: absolute;
            top: 181px;
        }
        img {
            opacity: 0.5745;
        }
    }
</style>
