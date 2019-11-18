<template>
    <div class="dialog-wrapper" v-transfer-dom v-show="showWrapper">
        <div ref="resizeTrigger">
            <transition name="dialog-fade">
                <div class="dialog-body" ref="body"
                    v-if="showBody"
                    :style="bodyStyle">
                    <div class="dialog-header" v-if="showHeader" ref="header">
                        <slot name="header"></slot>
                    </div>
                    <div class="dialog-content">
                        <slot></slot>
                    </div>
                    <div class="dialog-footer" v-if="showFooter" ref="footer">
                        <slot name="footer"></slot>
                    </div>
                    <i class="bk-icon icon-close" v-if="showCloseIcon" @click="handleCloseDialog"></i>
                </div>
            </transition>
        </div>
    </div>
</template>

<script>
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events.js'
    export default {
        name: 'cmdb-dialog',
        props: {
            value: Boolean,
            showHeader: {
                type: Boolean,
                default: true
            },
            showFooter: {
                type: Boolean,
                default: true
            },
            showCloseIcon: {
                type: Boolean,
                default: true
            },
            width: {
                type: Number,
                default: 720
            },
            height: Number,
            minHeight: Number
        },
        data () {
            return {
                timer: null,
                bodyHeight: 0,
                showWrapper: false,
                showBody: false
            }
        },
        computed: {
            autoResize () {
                return typeof this.height !== 'number'
            },
            bodyStyle () {
                const style = {
                    width: this.width + 'px',
                    '--height': (this.autoResize ? this.bodyHeight : this.height) + 'px'
                }
                if (!this.autoResize) {
                    style.height = this.height + 'px'
                    style.maxHeight = 'initial'
                }
                if (this.minHeight) {
                    style.minHeight = this.minHeight + 'px'
                }
                return style
            }
        },
        watch: {
            value: {
                immediate: true,
                handler (value) {
                    if (value) {
                        this.showWrapper = true
                        this.$nextTick(() => {
                            this.showBody = true
                        })
                    } else {
                        this.showBody = false
                        this.timer && clearTimeout(this.timer)
                        this.timer = setTimeout(() => {
                            this.showWrapper = false
                        }, 300)
                    }
                }
            }
        },
        mounted () {
            addResizeListener(this.$refs.resizeTrigger, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.resizeTrigger, this.resizeHandler)
        },
        methods: {
            resizeHandler () {
                this.$nextTick(() => {
                    this.bodyHeight = Math.min(this.$APP.height, this.$refs.resizeTrigger.offsetHeight)
                })
            },
            handleCloseDialog () {
                this.$emit('close')
                this.$emit('input', false)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .dialog-wrapper {
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        background-color: rgba(0, 0, 0, .6);
        z-index: 2000;
        @include scrollbar;
        .dialog-body {
            position: relative;
            margin: 0 auto;
            margin-top: calc((100vh - var(--height)) / 3);
            max-height: calc(100vh - 225px);
            min-height: 100px;
            border-radius: 2px;
            background-color: #fff;
            box-shadow:0px 4px 12px 0px rgba(0,0,0,0.2);
            @include scrollbar;
            .icon-close {
                position: absolute;
                top: 6px;
                right: 6px;
                width: 32px;
                height: 32px;
                line-height: 32px;
                font-size: 16px;
                text-align: center;
                color: #D8D8D8;
                cursor: pointer;
                &:hover {
                    color: #979BA5;
                }
            }
        }
    }
</style>

<style lang="scss">
    .dialog-fade-enter-active,
    .dialog-fade-leave-active {
        transition: opacity .3s ease;
    }
    .dialog-fade-enter,
    .dialog-fade-leave-to {
      opacity: 0;
    }
</style>
