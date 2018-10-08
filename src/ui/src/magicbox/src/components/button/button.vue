<template>
    <button
        class="bk-button"
        :title="title"
        :type="buttonType"
        :disabled="disabled"
        :class="['bk-' + themeType, 'bk-button-' + size, {'is-disabled': disabled, 'is-loading': loading}]"
        @click="handleClick">
        <i class="bk-icon" :class="['icon-' + icon]" v-if="icon"></i>
        <span><slot></slot></span>
    </button>
</template>

<script>
    /**
     * bk-button
     *
     * @module components/button
     * @desc 基础按钮
     *
     * @param type {string} [type=default] - 显示类型，接受 default primary info warning success danger
     * @param size {string} - 尺寸，接受 mini small normal large
     * @param title {string} - 提示信息
     * @param icon {string} - 显示 icon，使用蓝鲸 icon
     * @param disabled {boolean} [disabled=false] - 禁用
     * @param loading {boolean} [loading=true] - 加载中
     *
     * @example
     * <bk-button :type="primary" :size="large" :icon="bk" @click="btnClick">按钮</bk-button>
    */

    export default {
        name: 'bk-button',
        props: {
            // 'submit success abc'
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
            loading: Boolean
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
    @import '../../bk-magic-ui/src/button.scss'
</style>
