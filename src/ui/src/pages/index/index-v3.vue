<template>
    <div class="index-wrapper">
        <v-search class="index-search"></v-search>
        <v-recently></v-recently>
        <v-classify></v-classify>
         <p class="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. 腾讯蓝鲸 版权所有
        </p>
    </div>
</template>

<script>
    import vSearch from './children/search'
    import vRecently from './children/recently'
    import vClassify from './children/classify'
    export default {
        components: {
            vSearch,
            vRecently,
            vClassify
        },
        data () {
            const year = (new Date()).getFullYear()
            return {
                year
            }
        },
        async created () {
            await this.$store.dispatch('getBkBizList')
        }
    }
</script>

<style lang="scss" scoped>
    .index-wrapper{
        height: 100%;
        background-color: #f5f6fa;
        overflow-y: auto;
        overflow-x: hidden;
    }
    .index-search{
        width: 50%;
        margin: 0 auto;
        padding: 40px 0 50px;
    }
    .copyright{
        text-align: center;
        position: absolute;
        width: calc(100% - 50px);
        bottom: 0;
        left: 25px;
        height: 43px;
        margin: 0;
        line-height: 42px;
        font-size: 12px;
        color: rgba(116, 120, 131, 0.5);
        border-top: 1px solid rgba(116, 120, 131, 0.2);
        background-color: #f5f6fa;
    }
</style>