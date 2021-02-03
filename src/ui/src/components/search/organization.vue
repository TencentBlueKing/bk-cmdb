<template>
    <cmdb-form-organization
        v-model="localValue"
        v-bind="$attrs"
        @toggle="handleToggle">
    </cmdb-form-organization>
</template>

<script>
    import activeMixin from './mixins/active'
    export default {
        name: 'cmdb-search-organization',
        mixins: [activeMixin],
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
