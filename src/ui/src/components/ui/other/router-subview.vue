<template>
    <router-view class="subview"></router-view>
</template>

<script>
    export default {
        name: 'router-subview'
    }
</script>

<style lang="scss" scoped>
    .subview {
        width: 100%;
        height: 100%;
        position: absolute;
        top: 0;
        left: 0;
        z-index: 100;
        background-color: #fff;
        transition: transform .5s cubic-bezier(0.4, 0, 0.2, 1);
    }
</style>
