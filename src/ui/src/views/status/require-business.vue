<template>
    <div>
        <span v-if="$i18n.locale === 'en'">
            You don't have any business permission, please contact with administrator.
        </span>
        <span v-else>
            当前并无任何业务权限，请联系管理员添加。
        </span>
    </div>
</template>
<script>
    export default {
        created () {
            this.$store.commit('setHeaderTitle', '')
        }
    }
</script>