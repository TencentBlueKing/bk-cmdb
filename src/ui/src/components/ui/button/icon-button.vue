<template>
    <button class="bk-button bk-default icon-button" @click="handleClick($event)">
        <i :class="['icon-wrapper', icon]"></i>
    </button>
</template>

<script>
    export default {
        name: 'icon-button',
        props: {
            disabled: Boolean,
            icon: {
                type: String,
                required: true
            }
        },
        methods: {
            handleClick (event) {
                this.$emit('click', event)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .icon-button {
        width: 32px;
        height: 32px;
        line-height: 30px;
        padding: 0;
        font-size: 0px;
        &:before {
            content: "";
            width: 0;
            height: 100%;
            display: inline-block;
            vertical-align: middle;
        }
        .icon-wrapper {
            font-size: 14px;
        }
    }
</style>
