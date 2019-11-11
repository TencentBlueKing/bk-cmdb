<template>
    <bk-date-picker style="width: 100%"
        v-model="time"
        transfer
        :font-size="fontSize"
        :placeholder="placeholder || $t('选择日期范围')"
        :clearable="clearable"
        :type="timer ? 'datetimerange' : 'daterange'"
        :disabled="disabled">
    </bk-date-picker>
</template>

<script>
    export default {
        name: 'cmdb-form-date-range',
        props: {
            value: {
                type: [Array, String],
                default () {
                    return []
                }
            },
            disabled: {
                type: Boolean,
                default: false
            },
            timer: Boolean,
            clearable: {
                type: Boolean,
                default: true
            },
            placeholder: {
                type: String,
                default: ''
            },
            fontSize: {
                type: [String, Number],
                default: 'medium'
            }
        },
        data () {
            return {
                localValue: [...this.value]
            }
        },
        computed: {
            time: {
                get () {
                    return this.localValue.map(date => {
                        return date ? new Date(date) : ''
                    })
                },
                set (value) {
                    const localValue = value.map(date => this.$tools.formatTime(date, this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'))
                    this.localValue = localValue.filter(date => !!date)
                }
            }
        },
        watch: {
            value (value) {
                if ([...value].join('') !== this.localValue.join('')) {
                    this.localValue = [...value]
                }
            },
            localValue (value, oldValue) {
                if (value.join('') !== [...this.value].join('')) {
                    this.$emit('input', [...value])
                    this.$emit('change', [...value], [...oldValue])
                }
            }
        }
    }
</script>
