<template>
    <div class="breadcumbs-layout">
        <h1 class="current">{{current}}</h1>
    </div>
</template>

<script>
    export default {
        computed: {
            current () {
                if (this.$route.meta.breadcumbs) {
                    return this.$t(this.$route.meta.breadcumbs)
                }
                if (this.$route.meta.menu) {
                    return this.$t(this.$route.meta.menu.i18n)
                }
                return ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .breadcumbs-layout {
        padding: 19px 20px;
        height: 58px;
        .current {
            font-size: 16px;
            line-height: 20px;
            color: #313238;
            font-weight: normal;
        }
    }
</style>
