<template>
    <button
        class="bk-icon-button"
        :title="title"
        :disabled="disabled"
        :type="buttonType"
        :class="['bk-' + themeType, 'bk-button-' + size, {'is-disabled': disabled, 'is-loading': loading}]"
        @click="handleClick">
        <i class="bk-icon" :class="['icon-' + (icon || 'bk')]"></i>
        <i class="bk-text" :class="{'is-disabled': disabled}" v-if="!hideText"><slot></slot></i>
    </button>
</template>

<script>
    /**
     * bk-icon-button
     *
     * @module components/button-icon
     * @desc 图标按钮
     *
     * @param type {string} [type=default] - 显示类型，接受 default primary info warning success danger
     * @param size {string} - 尺寸，接受 mini small normal large
     * @param title {string} - 提示信息
     * @param icon {string} - 显示 icon，使用蓝鲸 icon
     * @param disabled {boolean} [disabled=false] - 禁用
     * @param loading {boolean} [loading=true] - 加载中
     * @param hide-text {boolean} [hide-text=false] - 是否隐藏文字
     *
     * @example
     * <bk-icon-button :type="primary" :size="large" :icon="bk" @click="btnClick">按钮</bk-icon-button>
    */

    export default {
        name: 'bk-icon-button',
        props: {
            type: {
                type: String,
                default: 'default',
                validator (value) {
                    let types = value.split(' ') // ['submit', 'success', 'reset']
                    let buttons = [
                        'button',
                        'submit',
                        'reset'
                    ]
                    let thenme = [
                        'default',
                        'info',
                        'primary',
                        'warning',
                        'success',
                        'danger'
                    ]
                    let valid = true
                    // for (let type of types) {
                    //     if (buttons.indexOf(type) === -1 && thenme.indexOf(type) === -1) {
                    //         valid = false
                    //     }
                    // }
                    types.forEach(type => {
                        if (buttons.indexOf(type) === -1 && thenme.indexOf(type) === -1) {
                            valid = false
                        }
                    })
                    return valid
                }
            },
            size: {
                type: String,
                default: 'normal',
                validator (value) {
                    return [
                        'mini',
                        'small',
                        'normal',
                        'large'
                    ].indexOf(value) > -1
                }
            },
            title: {
                type: String,
                default: ''
            },
            icon: String,
            disabled: Boolean,
            loading: Boolean,
            hideText: false
        },
        computed: {
            buttonType () {
                let types = this.type.split(' ')
                return types.find((type) => type === 'submit' || type === 'button' || type === 'reset')
            },
            themeType () {
                let types = this.type.split(' ')
                return types.find((type) => type !== 'submit' && type !== 'button' && type !== 'reset')
            }
        },
        methods: {
            handleClick (e) {
                if (!this.disabled && !this.loading) {
                    this.$emit('click', e)
                }
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/icon-button.scss'
</style>
