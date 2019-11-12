<template>
    <bk-search-select
        :data="searchOptions"
        v-model="searchValue"
        placeholder="名称关键字 / 或字段 / 字段值"
    >
    </bk-search-select>
</template>
<script>
    export default {
        components: {

        },
        data () {
            return {
                searchOptions: [],
                searchValue: []
            }
        }
    }
</script>
<style lang="scss" scoped>

</style>
