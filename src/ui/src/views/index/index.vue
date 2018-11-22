<template>
    <div class="index-wrapper">
        <v-search class="index-search"></v-search>
        <v-recently ref="recently"></v-recently>
        <v-classify></v-classify>
    </div>
</template>

<script>
    import vSearch from './children/search'
    import vRecently from './children/recently'
    import vClassify from './children/classify'
    import getScrollbarWidth from '@/utils/scrollbar-width.js'
    export default {
        components: {
            vSearch,
            vRecently,
            vClassify
        },
        data () {
            const year = (new Date()).getFullYear()
            return {
                year,
                scrollHandler: null,
                resizeHandler: null
            }
        },
        beforeRouteLeave (to, from, next) {
            this.$refs.recently.updateRecently(to)
            next()
        }
    }
</script>

<style lang="scss" scoped>
    .index-wrapper{
        position: relative;
        background-color: #f5f6fa;
    }
    .index-search{
        width: 50%;
        margin: 0 auto;
        padding: 40px 0 50px;
    }
</style>