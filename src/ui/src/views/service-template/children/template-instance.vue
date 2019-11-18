<template>
    <div class="template-instance-layout">
        <div class="instance-main">
            <div class="options clearfix">
                <div class="fr">
                    <bk-input class="filter-item" right-icon="bk-icon icon-search"
                        clearable
                        :placeholder="$t('请输入模块名称搜索')">
                    </bk-input>
                </div>
            </div>
            <bk-table class="instance-table"
                ref="instanceTable"
                v-bkloading="{ isLoading: $loading('getSetInstanceData') }">
                <bk-table-column :label="$t('模块名称')" prop="bk_module_name"></bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="host_count"></bk-table-column>
            </bk-table>
        </div>
    </div>
</template>

<script>
    export default {
        props: {},
        data () {
            return {}
        }
    }
</script>

<style lang="scss" scoped>
    .template-instance-layout {
        height: 100%;
    }
    .instance-main {
        .options {
            padding: 15px 0;
            font-size: 14px;
            color: #63656E;
            .filter-item {
                width: 230px;
            }
            .icon-cc-updating {
                color: #C4C6CC;
            }
        }
    }
</style>
