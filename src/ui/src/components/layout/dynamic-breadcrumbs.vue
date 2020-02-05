<template>
    <div class="breadcrumbs-layout clearfix">
        <i class="icon icon-cc-arrow fl" v-if="previous" @click="handleClick"></i>
        <h1 class="current fl">{{current}}</h1>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        computed: {
            ...mapGetters(['title']),
            current () {
                const menuI18n = this.$route.meta.menu.i18n && this.$t(this.$route.meta.menu.i18n)
                return this.title || this.$route.meta.title || menuI18n
            },
            previous () {
                return this.$route.meta.layout && this.$route.meta.layout.previous
            }
        },
        methods: {
            async handleClick () {
                const config = typeof this.previous === 'function' ? await this.previous(this.$parent.$refs.view) : this.previous
                this.$router.replace(config)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .breadcrumbs-layout {
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
        }
    }
</style>
