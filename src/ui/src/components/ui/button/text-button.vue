<template>
    <component href="javascript:void(0)"
        :class="['text-button', { 'disabled': disabled }]"
        :is="tag"
        @click="handleClick">
        <slot></slot>
    </component>
</template>

<script>
    export default {
        name: 'text-button',
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
    .text-button {
        display: inline-block;
        vertical-align: middle;
        line-height: 32px;
        color: $primaryColor;
        &.disabled {
            color: $textDisabledColor;
        }
    }
</style>
