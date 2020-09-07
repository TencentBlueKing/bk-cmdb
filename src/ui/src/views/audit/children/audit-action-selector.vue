<template>
    <bk-select
        v-bind="$attrs"
        v-model="localValue"
        multiple>
        <bk-option
            v-for="action in actions"
            :key="action.id"
            :id="action.id"
            :name="action.name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            target: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                dictionary: []
            }
        },
        computed: {
            actions () {
                const target = this.dictionary.find(target => target.id === this.target)
                return target ? target.operations : []
            },
            localValue: {
                get () {
                    return this.value
                },
                set (values) {
                    this.$emit('input', values)
                    this.$emit('change', values)
                }
            }
        },
        watch: {
            target () {
                this.localValue = []
            }
        },
        created () {
            this.getAuditDictionary()
        },
        methods: {
            async getAuditDictionary () {
                try {
                    this.dictionary = await this.$store.dispatch('audit/getDictionary', {
                        fromCache: true
                    })
                } catch (error) {
                    this.dictionary = []
                }
            }
        }
    }
</script>
