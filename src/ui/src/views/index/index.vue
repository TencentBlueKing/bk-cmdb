<template>
    <div class="index-layout">
        <search-input :is-full-text-search="isFullTextSearch"></search-input>
        <the-map style="user-select: none;"></the-map>
        <div class="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. 腾讯蓝鲸 版权所有. {{site.buildVersion}}
        </div>
    </div>
</template>

<script>
    import searchInput from './children/search-input'
    import theMap from './children/map'
    import { mapGetters } from 'vuex'
    export default {
        name: 'index',
        components: {
            searchInput,
            theMap
        },
        data () {
            return {
                year: (new Date()).getFullYear()
            }
        },
        computed: {
            ...mapGetters(['site']),
            isFullTextSearch () {
                return this.site.fullTextSearch === 'on'
            }
        },
        methods: {}
    }
</script>

<style lang="scss" scoped>
    .index-layout {
        padding: 0 0 50px;
        background-color: #f5f6fa;
        position: relative;
        z-index: 1;
    }
    .copyright{
        position: absolute;
        width: calc(100% - 50px);
        height: 43px;
        left: 25px;
        bottom: 0;
        line-height: 42px;
        font-size: 12px;
        text-align: center;
        color: rgba(116, 120, 131, 0.5);
        border-top: 1px solid rgba(116, 120, 131, 0.2);
        background-color: #f5f6fa;
        z-index: 2;
    }
</style>
