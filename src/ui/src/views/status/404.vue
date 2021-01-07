<template>
    <bk-exception class="exception-wrap-item" type="404">
        <p>{{$t('页面不存在')}}</p>
        <div class="text-wrap">
            <bk-link theme="primary" @click="handleRedirect">{{$t('返回首页')}}</bk-link>
        </div>
    </bk-exception>
</template>

<script>
    import { MENU_INDEX } from '@/dictionary/menu-symbol'
    export default {
        methods: {
            handleRedirect () {
                this.$routerActions.redirect({
                    name: MENU_INDEX
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .exception-wrap-item {
        justify-content: center;
    }
</style>
