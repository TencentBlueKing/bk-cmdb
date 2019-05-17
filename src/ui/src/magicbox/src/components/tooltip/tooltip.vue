<template>
    <div class="bk-tooltip" @mouseenter="handleShowPopper" @mouseleave="handleClosePopper">
        <div class="bk-tooltip-rel" ref="reference">
            <slot></slot>
        </div>
        <transition name="fade">
            <div
                class="bk-tooltip-popper"
                ref="popper"
                v-show="!disabled && (visible || always)"
                @mouseenter="handleShowPopper"
                @mouseleave="handleClosePopper"
                :data-transfer="transfer"
                v-transfer-dom>
                <div class="bk-tooltip-content">
                    <div class="bk-tooltip-arrow"></div>
                    <div class="bk-tooltip-inner" :style="{width: `${width}px`}"><slot name="content">{{content}}</slot></div>
                </div>
            </div>
        </transition>
    </div>
</template>
<script>
    import Popper from './popper'
    import TransferDom from '../../directives/transfer-dom'

    const oneOf = (value, validList) => {
        for (let i = 0; i < validList.length; i++) {
            if (value === validList[i]) {
                return true
            }
        }
        return false
    }

    export default {
        name: 'bk-tooltip',
        directives: {TransferDom},
        mixins: [Popper],
        props: {
            placement: {
                validator (value) {
                    return oneOf(
                        value,
                        [
                            'top', 'top-start', 'top-end', 'bottom', 'bottom-start', 'bottom-end',
                            'left', 'left-start', 'left-end', 'right', 'right-start', 'right-end'
                        ]
                    )
                },
                default: 'bottom'
            },
            content: {
                type: [String, Number],
                default: ''
            },
            delay: {
                type: Number,
                default: 100
            },
            width: {
                type: [String, Number],
                default: 'auto'
            },
            disabled: {
                type: Boolean,
                default: false
            },
            controlled: {
                type: Boolean,
                default: false
            },
            always: {
                type: Boolean,
                default: false
            },
            transfer: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
            }
        },
        methods: {
            handleShowPopper () {
                if (this.timeout) {
                    clearTimeout(this.timeout)
                }
                this.timeout = setTimeout(() => {
                    this.visible = true
                }, this.delay)
            },
            handleClosePopper () {
                if (this.timeout) {
                    clearTimeout(this.timeout)
                    if (!this.controlled) {
                        this.timeout = setTimeout(() => {
                            this.visible = false
                        }, 100)
                    }
                }
            }
        },
        mounted () {
            if (this.always) {
                this.updatePopper()
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/tooltip.scss'
</style>
