<template>
    <bk-switcher
        class="form-bool"
        size="small"
        theme="primary"
        v-model="localValue"
        :disabled="disabled">
    </bk-switcher>
</template>

<script>
    export default {
        name: 'cmdb-form-bool',
        props: {
            value: {
                type: [String, Boolean], // String是为了兼容初始数据为空的情况
                default: false
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        computed: {
            localValue: {
                get () {
                    if (typeof this.value === 'boolean') {
                        return this.value
                    }
                    if (this.value === 'true') {
                        return true
                    }
                    return false
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
    .form-bool{
        display: inline-block;
        vertical-align: middle;
    }
</style>
