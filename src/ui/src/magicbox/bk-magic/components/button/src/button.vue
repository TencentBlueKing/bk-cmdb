<template>
    <button
        class="bk-button"
        :type="btnType"
        :title="title"
        :disabled="disabled"
        :class="['bk-' + type, 'bk-button-' + size, {'is-disabled': disabled, 'is-loading': loading}]"
        @click="handleClick"> 
        <i class="bk-icon" :class="['icon-' + icon]" v-if="icon"></i>
        <span><slot></slot></span>
    </button>
</template>
<script>
    /**
     * bk-button
     * @module components/button
     * @desc 基础按钮
     * @param type {string} [type=default] - 显示类型，接受 default primary info warning success danger
     * @param btnType {string} [type=button] - 浏览器button的type属性 默认类型为 button
     * @param icon {string} - 显示icon，使用蓝鲸icon
     * @param size {string} - 尺寸，接受 mini small normal large
     * @param title {string} - 提示信息
     * @param disabled {boolean} [disabled=false] - 禁用
     * @param loading {boolean} [loading=true] - 加载中
     * @example
     * <bk-button :type="primary" :size="large" :icon="bk" @click="btnClick">按钮</bk-button>
    */
    export default {
        name: 'bk-button',
        props: {
            icon: String,
            disabled: Boolean,
            loading: Boolean,
            btnType: {
                type: String,
                default: 'button'
            },
            type: {
                type: String,
                default: 'default',
                validator (value) {
                    return [
                        'default',
                        'info',
                        'primary',
                        'warning',
                        'success',
                        'danger'].indexOf(value) > -1
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
                        'large'].indexOf(value) > -1
                }
            },
            title: {
                type: String,
                default: ''
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
