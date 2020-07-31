<template>
    <bk-table class="table-value" :data="list">
        <bk-table-column v-for="col in header"
            :key="col.bk_property_id"
            :prop="col.bk_property_id"
            :label="col.bk_property_name"
            :width="col.bk_property_type === 'bool' ? '90px' : ''">
            <template slot-scope="{ row }">
                <cmdb-property-value
                    v-bk-overflow-tips
                    :show-title="false"
                    :value="row[col['bk_property_id']]"
                    :property="col">
                </cmdb-property-value>
            </template>
        </bk-table-column>
        <div slot="empty">
            <span>{{$t('暂无数据')}}</span>
        </div>
    </bk-table>
</template>

<script>
    export default {
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            property: {
                type: Object,
                default: () => ({})
            }
        },
        data () {
            return {
                list: []
            }
        },
        computed: {
            header () {
                return this.property.option.map(option => option)
            }
        },
        watch: {
            value: {
                handler (value) {
                    this.list = value || []
                },
                immediate: true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-value {
        &.property-value {
            width: 100% !important;
            padding: 0 !important;
        }
    }
</style>
