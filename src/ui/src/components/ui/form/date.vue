<template>
    <bk-date-picker style="width: 100%"
        v-model="date"
        transfer
        :clearable="clearable"
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
            }
        },
        data () {
            return {
                localValue: this.value
            }
        },
        computed: {
            date: {
                get () {
                    if (!this.localValue) {
                        return ''
                    }
                    return new Date(this.localValue)
                },
                set (value) {
                    this.localValue = this.$tools.formatTime(value, 'YYYY-MM-DD')
                }
            }
        },
        watch: {
            value (value) {
                if (value !== this.localValue) {
                    this.localValue = value
                }
            },
            localValue (value, oldValue) {
                if (value !== this.value) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            }
        }
    }
</script>
