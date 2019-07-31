<template>
    <div class="instance-details-wrapper">
        <cmdb-table class="details-table" ref="table"
            :header="header"
            :list="attributeList"
            :wrapper-minus-height="300">
            <template slot="before_value" slot-scope="{ item }">
                <span>{{item.before_value ? item.before_value : '--'}}</span>
            </template>
            <template slot="show_value" slot-scope="{ item }">
                <span>{{item.show_value ? item.show_value : '--'}}</span>
            </template>
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
                        id: 'before_value',
                        name: this.$t("BusinessSynchronous['变更前']"),
                        sortable: false
                    }, {
                        id: 'show_value',
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
