<template>
    <bk-date-picker
        type="datetimerange"
        v-model="localValue"
        v-bind="$attrs">
    </bk-date-picker>
</template>

<script>
    export default {
        name: 'cmdb-search-time',
        props: {
            value: {
                type: Array,
                default: () => ([])
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value.map(str => new Date(str))
                },
                set (values) {
                    const formattedValues = values.filter(value => !!value).map(date => this.$tools.formatTime(date, 'YYYY-MM-DD hh:mm:ss'))
                    this.$emit('input', formattedValues)
                    this.$emit('change', formattedValues)
                }
            }
        }
    }
</script>
