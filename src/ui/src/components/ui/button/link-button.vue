<template>
    <component href="javascript:void(0)"
        :class="['link-button', { 'disabled': disabled }]"
        :is="tag"
        @click="handleClick">
        <slot></slot>
    </component>
</template>

<script>
    export default {
        name: 'link-button',
        props: {
            disabled: Boolean,
            tag: {
                type: String,
                default: 'a'
            },
            href: {
                type: String,
                default: ''
            }
        },
        methods: {
            handleClick (event) {
                if (this.disabled) {
                    return false
                }
                this.$emit('click', event)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .link-button {
        display: inline-block;
        vertical-align: middle;
        line-height: 32px;
        color: $primaryColor;
        &.disabled {
            color: $textDisabledColor;
            cursor: not-allowed;
        }
    }
</style>
