<template>
    <bk-select class="audit-target-selector"
        v-bind="$attrs"
        v-model="localValue">
        <bk-option
            v-for="option in options"
            :key="option.id"
            :id="option.id"
            :name="option.name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            },
            category: {
                type: String,
                default: 'business',
                validator (category) {
                    return ['host', 'business', 'resource', 'other'].includes(category)
                }
            }
        },
        data () {
            return {
                dictionary: [],
                targetMap: Object.freeze({
                    host: new Set(['host']),
                    business: new Set([
                        'dynamic_group',
                        'set_template',
                        'service_template',
                        'service_category',
                        'module',
                        'set',
                        'mainline_instance',
                        'service_instance',
                        'process',
                        'service_instance_label',
                        'host_apply',
                        'custom_field'
                    ]),
                    resource: new Set([
                        'business',
                        'model_instance',
                        'instance_association',
                        'resource_directory',
                        'cloud_area',
                        'cloud_account',
                        'cloud_sync_task'
                    ]),
                    other: new Set([
                        'model_group',
                        'model',
                        'model_attribute',
                        'model_unique',
                        'model_association',
                        'model_attribute_group',
                        'event',
                        'association_kind'
                    ])
                })
            }
        },
        computed: {
            options () {
                return this.dictionary.filter(target => this.targetMap[this.category].has(target.id))
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
