<template>
    <div class="breadcrumbs-layout clearfix">
        <h1 class="current fl">{{current}}</h1>
        <i class="breadcrumbs-split fl" v-show="showBreadcrumbs"></i>
        <div class="breadcrumbs fl" v-show="showBreadcrumbs">
            <template v-for="(item, index) in breadcrumbs">
                <a class="breadcrumbs-link" href="javascript:void(0)"
                    :class="{
                        'is-last': index === breadcrumbs.length - 1
                    }"
                    :key="index"
                    @click="handleBreadcrumbsClick(item, index)">
                    {{item.label}}
                </a>
                <span class="breadcrumbs-arrow"
                    v-if="index !== breadcrumbs.length - 1"
                    :key="index + 'arrow'">
                    &gt;
                </span>
            </template>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        computed: {
            ...mapGetters(['breadcrumbs', 'title']),
            showBreadcrumbs () {
                return this.breadcrumbs.length > 1
            },
            current () {
                const menuI18n = this.$route.meta.menu.i18n && this.$t(this.$route.meta.menu.i18n)
                return this.title || this.$route.meta.title || menuI18n
            }
        },
        methods: {
            handleBreadcrumbsClick (item, index) {
                const total = this.breadcrumbs.length
                if (index === total - 1) {
                    return false
                }
                if (item.hasOwnProperty('route')) {
                    this.$router.replace(item.route)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .breadcrumbs-layout {
        padding: 19px 20px;
        height: 58px;
        .current {
            font-size: 16px;
            line-height: 20px;
            color: #313238;
            font-weight: normal;
        }
        .breadcrumbs-split {
            width: 2px;
            height: 14px;
            margin: 3px 10px;
            background-color: $cmdbLayoutBorderColor;
        }
        .breadcrumbs {
            font-size: 0;
        }
        .breadcrumbs-link,
        .breadcrumbs-arrow {
            @include inlineBlock;
            font-size: 12px;
            line-height: 20px;
            color: #979BA5;
        }
        .breadcrumbs-link {
            &.is-last {
                cursor: default;
                color: #63656E;
            }
            &:not(.is-last):hover {
                color: #000;
            }
        }
        .breadcrumbs-arrow {
            margin: 0 5px;
        }
    }
</style>
