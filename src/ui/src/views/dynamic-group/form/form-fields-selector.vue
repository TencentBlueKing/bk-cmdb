<template>
    <bk-select v-model="localValue" v-bind="$attrs">
        <bk-option v-for="property in properties"
            :key="property.bk_property_id"
            :id="property.bk_property_id"
            :name="property.bk_property_name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: Array,
                default: () => ([])
            }
        },
        inject: ['dynamicGroupForm'],
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (values) {
                    this.$emit('input', values)
                }
            },
            target () {
                return this.dynamicGroupForm.formData.bk_obj_id
            },
            properties () {
                return this.dynamicGroupForm.propertyMap[this.target] || []
            }
        }
        
    }
</script>
