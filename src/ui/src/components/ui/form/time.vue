<template>
    <bk-date-picker class="cmdb-time"
        v-model="time"
        type="datetime"
        transfer
        editable
        :clearable="clearable"
        :disabled="disabled"
        :placeholder="placeholder">
    </bk-date-picker>
</template>

<script>
    export default {
        name: 'cmdb-form-time',
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
            time: {
                get () {
                    if (!this.value) {
                        return ''
                    }
                    return new Date(this.value)
                },
                set (value) {
                    const previousValue = this.value
                    const currentValue = this.$tools.formatTime(value, 'YYYY-MM-DD HH:mm:ss')
                    this.$emit('input', currentValue)
                    this.$emit('change', currentValue, previousValue)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-time {
        width: 100%;
    }
</style>
