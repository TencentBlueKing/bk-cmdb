<template>
    <cmdb-form-enum
        v-if="property.bk_property_type === 'enum'"
        :allow-clear="true"
        :options="property.option"
        v-model="property.__extra__.value"
    >
    </cmdb-form-enum>
    <cmdb-form-bool-input
        v-else-if="property.bk_property_type === 'bool'"
        v-model="property.__extra__.value"
    >
    </cmdb-form-bool-input>
    <cmdb-search-input class="form-element-item"
        v-else-if="['singlechar', 'longchar'].includes(property.bk_property_type)"
        v-model="property.__extra__.value"
    >
    </cmdb-search-input>
    <cmdb-form-date-range
        v-else-if="['date', 'time'].includes(property.bk_property_type)"
        v-model="property.__extra__.value"
    >
    </cmdb-form-date-range>
    <component
        v-else
        :is="`cmdb-form-${property.bk_property_type}`"
        v-model="property.__extra__.value"
    >
    </component>
</template>

<script>
    export default {
        props: {
            property: {
                type: Object,
                required: true,
                default: () => ({})
            }
        },
        watch: {
            'property.__extra__.value': function (value) {
                this.$emit('value-change', value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-element-item {
        &.cmdb-search-input {
            /deep/ .search-input-wrapper {
                position: relative;
            }
        }
    }
</style>
