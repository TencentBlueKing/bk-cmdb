<template>
    <bk-tag-input ref="tagInput"
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
        },
        mounted () {
            this.addPasteEvent()
        },
        beforeDestroy () {
            this.removePasteEvent()
        },
        methods: {
            addPasteEvent () {
                this.$refs.tagInput.$refs.input.addEventListener('paste', this.handlePaste)
            },
            removePasteEvent () {
                this.$refs.tagInput.$refs.input.removeEventListener('paste', this.handlePaste)
            },
            handlePaste (event) {
                const text = event.clipboardData.getData('text')
                const values = text.split(/,|;|\n/).map(value => value.trim()).filter(value => value.length)
                const value = [...new Set([...this.localValue, ...values])]
                this.localValue = value
            }
        }
    }
</script>
