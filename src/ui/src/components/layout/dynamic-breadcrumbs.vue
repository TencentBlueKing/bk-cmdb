<template>
    <div class="breadcrumbs-layout clearfix">
        <i class="icon icon-cc-arrow fl" v-if="from && current" @click="handleClick"></i>
        <h1 class="current fl">{{current}}</h1>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { Base64 } from 'js-base64'
    export default {
        computed: {
            ...mapGetters(['title']),
            current () {
                const menuI18n = this.$route.meta.menu.i18n && this.$t(this.$route.meta.menu.i18n)
                return this.title || this.$route.meta.title || menuI18n
            },
            defaultFrom () {
                const menu = this.$route.meta.menu || {}
                if (menu.relative) {
                    return { name: menu.relative }
                }
                return null
            },
            latest () {
                let latest
                if (this.$route.query.hasOwnProperty('_f')) {
                    try {
                        const historyList = JSON.parse(window.sessionStorage.getItem('history'))
                        latest = historyList.pop()
                    } catch (e) {
                        // ignore
                    }
                }
                return latest
            },
            from () {
                if (this.latest) {
                    try {
                        return JSON.parse(Base64.decode(this.latest))
                    } catch (error) {
                        return this.defaultFrom
                    }
                }
                return this.defaultFrom
            }
        },
        methods: {
            async handleClick () {
                this.$routerActions.redirect({ ...this.from, back: true })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .breadcrumbs-layout {
        display: flex;
        align-items: center;
        padding: 14px 20px;
        height: 53px;
        border-bottom: 1px solid $borderColor;
        .icon-cc-arrow {
            display: block;
            width: 24px;
            height: 24px;
            line-height: 24px;
            font-size: 14px;
            text-align: center;
            margin-right: 3px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                color: #699df4;
            }
        }
        .current {
            font-size: 16px;
            line-height: 24px;
            color: #313238;
            font-weight: normal;
            @include ellipsis;
        }
    }
</style>
