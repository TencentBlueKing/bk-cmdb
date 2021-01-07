<template>
    <bk-select
        multiple
        searchable
        v-model="localValue"
        v-bind="$attrs">
        <bk-option v-for="(option, index) in options"
            :key="index"
            :id="option"
            :name="option">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-search-list',
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            options: {
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
                    this.$emit('change', value)
                }
            }
        }
    }
</script>
