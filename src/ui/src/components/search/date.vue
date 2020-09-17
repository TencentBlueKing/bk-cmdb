<template>
    <bk-date-picker
        type="daterange"
        v-model="localValue"
        v-bind="$attrs">
    </bk-date-picker>
</template>

<script>
    export default {
        name: 'cmdb-search-date',
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
                    const formattedValues = values.filter(value => !!value).map(date => this.$tools.formatTime(date, 'YYYY-MM-DD'))
                    this.$emit('input', formattedValues)
                    this.$emit('change', formattedValues)
                }
            }
        }
    }
</script>
