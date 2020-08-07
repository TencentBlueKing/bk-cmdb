<template>
    <bk-select
        v-model="localValue"
        searchable
        font-size="medium"
        :clearable="false">
        <bk-option v-for="option in options"
            :key="option.bk_property_id"
            :id="option.bk_property_id"
            :name="option.bk_property_name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-property-selector',
        props: {
            properties: {
                type: Array,
                default: () => ([])
            },
            objectUnique: {
                type: Array,
                default: () => ([])
            },
            value: {
                type: [String, Number],
                default: ''
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
            },
            options () {
                const options = this.properties.filter(property => !!property.id)
                const uniqueIds = (this.objectUnique.find(unique => unique.must_check) || {}).keys || []
                options.sort((m, n) => {
                    const seedM = uniqueIds.some(key => key.key_id === m.id) ? m.id : Infinity
                    const seedN = uniqueIds.some(key => key.key_id === n.id) ? n.id : Infinity
                    const result = seedM - seedN
                    return isNaN(result) ? 0 : result
                })
                return options
            }
        }
    }
</script>
