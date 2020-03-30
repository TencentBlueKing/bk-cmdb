<template>
    <bk-select v-if="display === 'selector'"
        searchable
        :readonly="readonly"
        :disabled="disabled"
        :placeholder="$t('请选择xx', { name: $t('资源类型') })"
        v-model="selected">
        <bk-option id="host" :name="$t('主机')"></bk-option>
    </bk-select>
    <span v-else>{{getResourceInfo()}}</span>
</template>

<script>
    export default {
        name: 'task-resource-selector',
        props: {
            display: {
                type: String,
                default: 'selector'
            },
            readonly: Boolean,
            disabled: Boolean,
            value: {
                type: [String, Number]
            }
        },
        data () {
            return {
                resources: [{
                    id: 'host',
                    name: this.$t('主机')
                }]
            }
        },
        computed: {
            selected: {
                get () {
                    return this.value
                },
                set (value, oldValue) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            }
        },
        methods: {
            getResourceInfo () {
                const resource = this.resources.find(resource => resource.id === this.value)
                return resource ? resource.name : '--'
            }
        }
    }
</script>
