<template>
    <div class="tips clearfix">
        <i class="icon icon-cc-exclamation-tips fl"></i>
        <p class="tips-content">
            <slot></slot>
        </p>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-tips'
    }
</script>

<style lang="scss" scoped>
    .tips {
        height: 30px;
        line-height: 30px;
        font-size: 14px;
        background: rgba(240, 248, 255, 1);
        border-radius: 2px;
        border: 1px solid rgba(163, 197, 253, 1);
        .icon {
            font-size: 18px;
            color: #3A84FF;
            margin: 6px 5px 0 9px;
        }
        .tips-content {
            @include ellipsis;
        }
    }
</style>
