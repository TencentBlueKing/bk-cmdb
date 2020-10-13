<template>
    <cmdb-form-organization
        v-model="localValue"
        v-bind="$attrs">
    </cmdb-form-organization>
</template>

<script>
    export default {
        name: 'cmdb-search-organization',
        props: {
            value: {
                type: Array,
                default: () => ([])
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('input', value)
                }
            }
        }
    }
</script>
