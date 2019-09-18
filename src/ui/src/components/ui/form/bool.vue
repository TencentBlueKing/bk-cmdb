<template>
    <bk-switcher
        class="form-bool"
        size="small"
        theme="primary"
        :value="localValue"
        :disabled="disabled"
        @change="handleChange">
    </bk-switcher>
</template>

<script>
    export default {
        name: 'cmdb-form-bool',
        props: {
            value: {
                type: [Boolean, String], // String是为了兼容初始数据为空的情况
                default: false
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localValue: typeof this.value === 'boolean' ? this.value : false
            }
        },
        watch: {
            value (value) {
                this.localValue = typeof value === 'boolean' ? value : false
            }
        },
        methods: {
            handleChange (selected) {
                this.$emit('input', selected, this)
                this.$emit('change', selected, this)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-bool{
        display: inline-block;
        vertical-align: middle;
    }
</style>
