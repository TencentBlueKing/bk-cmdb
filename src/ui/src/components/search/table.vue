<template>
    <bk-tag-input
        allow-create
        allow-auto-match
        v-if="multiple"
        v-model="localValue"
        v-bind="$attrs"
        :list="[]"
        @click.native="handleToggle(true)"
        @blur="handleToggle(false, ...arguments)">
    </bk-tag-input>
    <bk-input v-else
        v-model="localValue"
        v-bind="$attrs"
        @focus="handleToggle(true, ...arguments)"
        @blur="handleToggle(false, ...arguments)">
    </bk-input>
</template>

<script>
    import activeMixin from './mixin-active'
    export default {
        name: 'cmdb-search-table',
        mixins: [activeMixin],
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
