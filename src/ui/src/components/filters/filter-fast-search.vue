<template>
    <bk-input class="filter-fast-search"
        v-model.trim="value"
        :placeholder="$t('请输入IP或固资编号')"
        @enter="handleSearch">
    </bk-input>
</template>

<script>
    import FilterStore from './store'
    export default {
        data () {
            return {
                value: ''
            }
        },
        methods: {
            async handleSearch () {
                const { valid: isIP } = await this.$validator.verify(this.value, 'ip')
                if (isIP) {
                    FilterStore.updateIP({
                        text: this.value,
                        exact: true,
                        inner: true,
                        outer: true
                    })
                } else {
                    FilterStore.createOrUpdateCondition([{
                        field: 'bk_asset_id',
                        model: 'host',
                        operator: '$in',
                        value: this.value
                    }])
                }
                this.value = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .filter-fast-search {
        display: inline-flex;
    }
</style>
