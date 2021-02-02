<template>
    <bk-input class="filter-fast-search"
        v-model.trim="value"
        :placeholder="$t('请输入IP或固资编号')"
        @enter="handleSearch"
        @paste="handlePaste">
    </bk-input>
</template>

<script>
    import FilterStore from './store'
    import IS_IP from 'validator/es/lib/isIP'
    export default {
        data () {
            return {
                value: ''
            }
        },
        methods: {
            async handleSearch () {
                this.dispatchFilter(this.value)
            },
            handlePaste (value, event) {
                event.preventDefault()
                const text = event.clipboardData.getData('text').trim()
                this.dispatchFilter(text)
            },
            dispatchFilter (currentValue) {
                const values = currentValue.trim().split(/,|;|\n/g).map(text => text.trim()).filter(text => text.length)
                const IP = []
                const assets = []
                values.forEach(text => {
                    IS_IP(text) ? IP.push(text) : assets.push(text)
                })
                if (IP.length) {
                    FilterStore.updateIP({
                        text: IP.join('\n'),
                        exact: true,
                        inner: true,
                        outer: true
                    })
                }
                if (assets.length) {
                    FilterStore.createOrUpdateCondition([{
                        field: 'bk_asset_id',
                        model: 'host',
                        operator: '$in',
                        value: assets
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
