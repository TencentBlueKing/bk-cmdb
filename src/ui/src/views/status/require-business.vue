<template>
    <div>
        当前并无任何业务权限，请联系管理员添加
    </div>
</template>
<script>
    export default {
        created () {
            this.$store.commit('setHeaderTitle', '')
        }
    }
</script>