<template>
    <bk-tag-input
        allow-create
        allow-auto-match
        v-if="multiple"
        v-model="localValue"
        v-bind="$attrs"
        :list="[]">
    </bk-tag-input>
    <bk-input v-else
        v-model="localValue"
        v-bind="$attrs">
    </bk-input>
</template>

<script>
    export default {
        name: 'cmdb-search-singlechar',
        props: {
            value: {
                type: [String, Array],
                default: ''
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
