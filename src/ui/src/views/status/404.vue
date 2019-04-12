<template>
    <div class="img-404">
        <img src="../../assets/images/404.png" alt="">
    </div>
</template>
<script>
    export default {
        created () {
            this.$store.commit('setHeaderTitle', '')
        }
    }
</script>
<style lang="scss" scoped>
    .img-404{
        position: relative;
        min-height: calc(100% - 50px);
        img{
            position: absolute;
            left: 50%;
            top: 50%;
            transform: translate(-50%, -50%);
        }
    }
</style>
