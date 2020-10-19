<template>
    <cmdb-form-objuser
        v-model="localValue"
        v-bind="$attrs">
    </cmdb-form-objuser>
</template>

<script>
    export default {
        name: 'cmdb-search-objuser',
        props: {
            value: {
                type: Array,
                default: () => ([])
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value.join(',')
                },
                set (value) {
                    const values = value.split(',')
                    this.$emit('input', values)
                    this.$emit('change', values)
                }
            }
        }
    }
</script>
