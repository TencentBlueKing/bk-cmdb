<template>
    <div class="breadcrumb-wrapper">
        <ul class="breadcrumb-list">
            <li class="breadcrumb-item" v-for="(crumb, index) in breadcrumbs" :key="index">
                <router-link exact class="breadcrumb-prev"
                    v-if="(index !== count - 1) && crumb.hasOwnProperty('path')"
                    :to="crumb.path"
                    :title="crumb.name">
                    {{crumb.name}}
                </router-link>
                <span v-else :class="{'breadcrumb-last': index === count - 1}">{{crumb.name}}</span>
                <i class="breadcrumb-angle icon-cc-angle-right" v-if="index !== count - 1"></i>
            </li>
        </ul>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {}
        },
        computed: {
            ...mapGetters('main', ['breadcrumbs']),
            count () {
                return this.breadcrumbs.length
            }
        }
    }
</script>

<style lang="scss" scoped>
    .breadcrumb-list{
        font-size: 0;
        height: 36px;
        line-height: 36px;
        .breadcrumb-item{
            font-size: 14px;
            display: inline-block;
            a.breadcrumb-prev{
                color: #969696;
                letter-spacing: 0;
                &:hover{
                    color: #3c96ff;
                }
            }
            .breadcrumb-last{
                color: #535353;
            }
            .breadcrumb-angle{
                display: inline-block;
                font-size: 12px;
                color: #999;
                margin: 0 5px 0 2px;
            }
        }
    }
</style>