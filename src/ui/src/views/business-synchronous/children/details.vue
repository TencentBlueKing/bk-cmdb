<template>
    <div class="instance-details-wrapper">
        <bk-table
            :data="attributeList"
            :max-height="$APP.height - 300">
            <bk-table-column prop="property_name" :label="$t('属性名称')" show-overflow-tooltip></bk-table-column>
            <bk-table-column prop="before_value" :label="$t('变更前')" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.before_value ? row.before_value : '--'}}</template>
            </bk-table-column>
            <bk-table-column prop="show_value" :label="$t('变更后')" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.show_value ? row.show_value : '--'}}</template>
            </bk-table-column>
        </bk-table>
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
        }
    }
</script>

<style lang="scss" scoped>
    .instance-details-wrapper {
        padding: 20px;
    }
</style>
