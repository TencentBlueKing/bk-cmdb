<template>
    <div class="tips" :style="tipsStyle">
        <i :class="icon" v-if="icon" :style="iconStyle"></i>
        <p class="tips-content">
            <slot></slot>
        </p>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-tips',
        props: {
            tipsStyle: {
                type: Object,
                default: () => ({})
            },
            icon: {
                type: [String, Boolean],
                default: 'icon icon-cc-exclamation-tips'
            },
            iconStyle: {
                type: Object,
                default: () => ({})
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tips {
        display: flex;
        height: 30px;
        line-height: 30px;
        font-size: 14px;
        background: rgba(240, 248, 255, 1);
        border-radius: 2px;
        border: 1px solid rgba(163, 197, 253, 1);
        .icon {
            flex: 18px 0 0;
            font-size: 18px;
            line-height: 30px;
            color: #3A84FF;
            margin: 0px 5px;
        }
        .tips-content {
            @include ellipsis;
        }
    }
</style>
