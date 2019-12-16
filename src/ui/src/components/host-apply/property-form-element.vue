<template>
    <div class="property-form-element">
        <component
            :is="`cmdb-form-${property.bk_property_type}`"
            :class="['form-element-item', property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
            :options="property.option || []"
            :data-vv-name="property.bk_property_id"
            :data-vv-as="property.bk_property_name"
            :placeholder="$t('请输入xx', { name: property.bk_property_name })"
            :auto-check="false"
            v-validate="$tools.getValidateRules(property)"
            v-model.trim="property.__extra__.value"
        >
        </component>
        <div class="form-error"
            v-if="errors.has(property.bk_property_id)">
            {{errors.first(property.bk_property_id)}}
        </div>
    </div>
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
            'property.__extra__.value': async function (value) {
                this.property.__extra__.valid = await this.$validator.validate()
                this.$emit('value-change', value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .property-form-element {
        .form-error {
            margin-top: 2px;
            font-size: 12px;
            color: $cmdbDangerColor;
        }
        .form-element-item {
            &.cmdb-search-input {
                /deep/ .search-input-wrapper {
                    position: relative;
                }
            }
        }
    }
</style>
