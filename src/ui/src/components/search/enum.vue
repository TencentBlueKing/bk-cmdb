<template>
    <bk-select
        searchable
        v-model="localValue"
        v-bind="$attrs"
        :multiple="multiple">
        <bk-option v-for="option in options"
            :key="option.id"
            :id="option.id"
            :name="option.name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-search-enum',
        props: {
            value: {
                type: [String, Array],
                default: () => ([])
            },
            options: {
                type: Array,
                default: () => ([])
            }
        },
        computed: {
            multiple () {
                return Array.isArray(this.value)
            },
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        }
    }
</script>
