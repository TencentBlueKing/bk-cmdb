<template>
    <button 
        class="bk-icon-button" 
        :title="title"
        :disabled="disabled"
        :class="['bk-' + type, 'bk-button-' + size, {'is-disabled': disabled, 'is-loading': loading}]"
        @click="handleClick"> 
        <i class="bk-icon" :class="['icon-' + (icon || 'bk')]"></i>
        <i class="bk-text" v-if="!hideText"><slot></slot></i>
    </button>
</template>

<script>
    /**
     * bk-icon-button
     * @module components/button-icon
     * @desc 图标按钮
     * @param type {string} [type=default] - 显示类型，接受 default primary info warning success danger
     * @param icon {string} - 显示icon，使用蓝鲸icon
     * @param size {string} - 尺寸，接受 mini small normal large
     * @param title {string} - 提示信息
     * @param disabled {boolean} [disabled=false] - 禁用
     * @param loading {boolean} [loading=true] - 加载中
     * @example
     * <bk-icon-button :type="primary" :size="large" :icon="bk" @click="btnClick">按钮</bk-icon-button>
    */
    export default{
        name: 'bk-icon-button',
        props: {
            icon: String,
            disabled: Boolean,
            loading: Boolean,
            hideText: false,
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
                        'danger'
                    ].indexOf(value) > -1
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
