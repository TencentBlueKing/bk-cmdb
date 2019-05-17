<template>
    <div class="index-wrapper">
        <v-search class="index-search"></v-search>
        <v-recently ref="recently"></v-recently>
        <v-classify></v-classify>
        <cmdb-main-inject class="copyright" ref="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. 腾讯蓝鲸 版权所有
        </cmdb-main-inject>
    </div>
</template>

<script>
    import cmdbMainInject from '@/components/layout/main-inject'
    import vSearch from './children/search'
    import vRecently from './children/recently'
    import vClassify from './children/classify'
    import { addMainResizeListener, removeMainResizeListener } from '@/utils/main-scroller'
    export default {
        components: {
            vSearch,
            vRecently,
            vClassify,
            cmdbMainInject
        },
        data () {
            const year = (new Date()).getFullYear()
            return {
                year,
                resizeHandler: null
            }
        },
        beforeRouteLeave (to, from, next) {
            this.$refs.recently.updateRecently(to)
            next()
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Index["首页"]'))
        },
        mounted () {
            const $copyright = this.$refs.copyright.$el
            this.resizeHandler = event => {
                const target = event.target
                if (target.offsetWidth < 1100) {
                    $copyright.style.bottom = '8px'
                } else {
                    $copyright.style.bottom = 0
                }
            }
            addMainResizeListener(this.resizeHandler)
            this.resizeHandler({target: document.querySelector('.main-scroller')})
        },
        beforeDestroy () {
            removeMainResizeListener(this.resizeHandler)
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
    }
</style>