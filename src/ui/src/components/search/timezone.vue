<template>
    <bk-select
        multiple
        searchable
        v-model="localValue"
        v-bind="$attrs">
        <bk-option v-for="timezone in timezones"
            :key="timezone"
            :id="timezone"
            :name="timezone">
        </bk-option>
    </bk-select>
</template>

<script>
    import TimeZones from '../ui/form/timezone.json'
    export default {
        name: 'cmdb-search-timezone',
        props: {
            value: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                timezones: Object.freeze(TimeZones)
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (values) {
                    this.$emit('input', values)
                    this.$emit('change', values)
                }
            }
        }
    }
</script>
