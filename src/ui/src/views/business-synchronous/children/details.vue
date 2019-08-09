<template>
    <div class="instance-details-wrapper">
        <bk-table
            :data="attributeList"
            :max-height="$APP.height - 300">
            <bk-table-column prop="property_name" :label="$t('属性名称')"></bk-table-column>
            <bk-table-column prop="before_value" :label="$t('变更前')">
                <template slot-scope="{ row }">
                    <span>{{row.before_value ? row.before_value : '--'}}</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="show_value" :label="$t('变更后')">
                <template slot-scope="{ row }">
                    <span>{{row.show_value ? row.show_value : '--'}}</span>
                </template>
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
