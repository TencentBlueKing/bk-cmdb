<template>
    <div class="instance-details-wrapper">
        <cmdb-table class="details-table" ref="table"
            :header="header"
            :list="attributeList"
            :wrapper-minus-height="300">
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        props: {
            attributeList: {
                type: Array,
                default: () => {
                    return []
                }
            }
        },
        data () {
            return {
                header: [
                    {
                        id: 'property_name',
                        name: this.$t("BusinessSynchronous['属性名称']"),
                        sortable: false
                    }, {
                        id: 'property_value',
                        name: this.$t("BusinessSynchronous['变更前']"),
                        sortable: false
                    }, {
                        id: 'template_property_value',
                        name: this.$t("BusinessSynchronous['变更后']"),
                        sortable: false
                    }
                ]
            }
        }
    }
</script>

<style lang="scss" scoped>
    .instance-details-wrapper {
        padding: 20px;
    }
</style>
