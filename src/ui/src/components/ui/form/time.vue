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
        data () {
            return {
                localValue: this.value
            }
        },
        computed: {
            time: {
                get () {
                    if (!this.localValue) {
                        return ''
                    }
                    return new Date(this.localValue)
                },
                set (value) {
                    this.localValue = this.$tools.formatTime(value, 'YYYY-MM-DD HH:mm:ss')
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

<style lang="scss" scoped>
    .cmdb-time {
        width: 100%;
    }
</style>
