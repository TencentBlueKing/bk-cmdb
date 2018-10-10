<template>
    <div :class="classObject">
        <input type="checkbox" v-model="enabled" :disabled="disabled">
        <label class="switcher-label" v-show="showText">
            <span class="switcher-text on-text">{{onText}}</span>
            <span class="switcher-text off-text">{{offText}}</span>
        </label>
    </div>
</template>
<script>
    /**
     * bk-switcher
     * @module components/switcher
     * @desc 开关
     * @param {boolean} [selected=false] - 是否打开
     * @param {boolean} [show-text=true] - 是否显示ON/OFF
     * @param {boolean} [isDisabled=true] - 是否禁用
     * @example
     * <bk-switcher :selected="isSelected" :show-text="showText"></bk-switcher>
     */

    export default {
        name: 'bk-switcher',
        props: {
            disabled: {
                type: Boolean,
                default: false
            },
            showText: {
                type: Boolean,
                default: true
            },
            selected: {
                type: Boolean,
                default: false
            },
            onText: {
                type: String,
                default: 'ON'
            },
            offText: {
                type: String,
                default: 'OFF'
            },
            isOutline: {
                type: Boolean,
                default: false
            },
            isSquare: {
                type: Boolean,
                default: false
            },
            size: {
                type: String,
                default: 'normal',
                validator (value) {
                    return [
                        'normal',
                        'small'
                    ].indexOf(value) > -1
                }
            }
        },
        data () {
            return {
                label: this.selected ? this.onText : this.offText,
                enabled: !!this.selected
            }
        },
        watch: {
            enabled (val) {
                this.label = this.enabled ? this.onText : this.offText
                this.$emit('change', val)
            },
            selected (val) {
                this.enabled = !!val
            }
        },
        computed: {
            classObject () {
                const {
                    enabled,
                    disabled,
                    size,
                    showText,
                    isOutline,
                    isSquare
                } = this
                let style = {
                    'bk-switcher': true,
                    'bk-switcher-outline': isOutline,
                    'bk-switcher-square': isSquare,
                    'show-label': true,
                    'is-disabled': disabled,
                    'is-checked': enabled,
                    'is-unchecked': !enabled
                }
                if (size) {
                    let sizeStr = 'bk-switcher-' + size
                    style[sizeStr] = true
                }
                return style
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/switcher.scss'
</style>
