<template>
    <bk-date-picker class="cmdb-date"
        v-model="date"
        transfer
        editable
        :clearable="clearable"
        :placeholder="placeholder"
        :disabled="disabled">
    </bk-date-picker>
</template>

<script>
    export default {
        name: 'cmdb-form-date',
        props: {
            value: {
                type: String,
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            clearable: {
                type: Boolean,
                default: true
            },
            placeholder: {
                type: String,
                default: ''
            }
        },
        computed: {
            date: {
                get () {
                    if (!this.value) {
                        return ''
                    }
                    return new Date(this.value)
                },
                set (value) {
                    const previousValue = this.value
                    const currentValue = this.$tools.formatTime(value, 'YYYY-MM-DD')
                    this.$emit('input', currentValue)
                    this.$emit('change', currentValue, previousValue)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-date {
        width: 100%;
    }
</style>
