<template>
    <div class="breadcumbs-layout clearfix">
        <h1 class="current fl">{{current}}</h1>
        <i class="breadcumbs-split fl" v-show="breadcumbs.length"></i>
        <div class="breadcumbs fl" v-show="breadcumbs.length">
            <template v-for="(item, index) in breadcumbs">
                <a class="breadcumbs-link" href="javascript:void(0)"
                    :class="{
                        'is-last': index === breadcumbs.length - 1
                    }"
                    :key="index"
                    @click="handleBreadcumbsClick(item, index)">
                    {{item.i18n ? $t(item.i18n) : item.name}}
                </a>
                <span class="breadcumbs-arrow"
                    v-if="index !== breadcumbs.length - 1"
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
            ...mapGetters(['breadcumbs', 'title']),
            current () {
                if (this.$route.meta.menu.i18n) {
                    return this.$t(this.$route.meta.menu.i18n)
                }
                return this.title
            }
        },
        methods: {
            handleBreadcumbsClick (item, index) {
                if (index === this.breadcumbs.length - 1) {
                    return false
                }
                if (item.hasOwnProperty('route')) {
                    this.$router.push(item.route)
                }
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
        .breadcumbs-split {
            width: 2px;
            height: 14px;
            margin: 3px 10px;
            background-color: $cmdbLayoutBorderColor;
        }
        .breadcumbs {
            font-size: 0;
        }
        .breadcumbs-link,
        .breadcumbs-arrow {
            @include inlineBlock;
            font-size: 12px;
            line-height: 20px;
            color: #979BA5;
        }
        .breadcumbs-link {
            &.is-last {
                cursor: default;
                color: #63656E;
            }
            &:not(.is-last):hover {
                color: #000;
            }
        }
        .breadcumbs-arrow {
            margin: 0 5px;
        }
    }
</style>
