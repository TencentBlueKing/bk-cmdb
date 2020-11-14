<template>
    <bk-switcher class="cmdb-search-bool"
        size="small"
        v-bind="$attrs"
        v-model="localValue">
    </bk-switcher>
</template>

<script>
    import activeMixin from './mixin-active'
    export default {
        name: 'cmdb-search-bool',
        mixins: [activeMixin],
        props: {
            value: [String, Boolean]
        },
        computed: {
            localValue: {
                get () {
                    return !!this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    .cmdb-search-bool {
        flex: 28px 0 0 !important;
    }
</style>
