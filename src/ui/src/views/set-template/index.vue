<template>
    <div class="template-layout">
        <div class="options clearfix">
            <div class="fl">
                <bk-button theme="primary">{{$t('新建')}}</bk-button>
                <bk-button theme="default" class="ml10">{{$t('批量删除')}}</bk-button>
            </div>
            <div class="fr">
                <bk-input :placeholder="$t('模板名称')" right-icon="icon-search"></bk-input>
            </div>
        </div>
        <bk-table class="template-table" :data="list">
            <bk-table-column type="selection"></bk-table-column>
            <bk-table-column :label="$t('模板名称')" prop="name"></bk-table-column>
            <bk-table-column :label="$t('应用数量')" prop="count"></bk-table-column>
            <bk-table-column :label="$t('修改人')" prop="modifier"></bk-table-column>
            <bk-table-column :label="$t('修改时间')" prop="last_time"></bk-table-column>
            <bk-table-column :label="$t('操作')" width="180">
                <template slot-scope="{ row }">
                    <bk-button text @click="handleEdit(row)">{{$t('编辑')}}</bk-button>
                    <bk-button text class="ml15" @click="handlePreview(row)">{{$t('预览')}}</bk-button>
                    <bk-button text class="ml15"
                        :disabled="!row.count"
                        @click="handleDelete(row)">
                        {{$t('删除')}}
                    </bk-button>
                </template>
            </bk-table-column>
            <template slot="empty">
                <i class="bk-table-empty-icon bk-icon icon-empty"></i>
                <i18n path="空集群模板提示" tag="div">
                    <bk-button text @click="handleCreate" place="link">{{$t('立即创建')}}</bk-button>
                </i18n>
            </template>
        </bk-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                list: []
            }
        },
        methods: {
            handleCreate () {},
            handleEdit (row) {},
            handlePreview (row) {},
            handleDelete (row) {}
        }
    }
</script>

<style lang="scss" scoped>
    .template-layout {
        padding: 0 20px;
    }
    .options {
        font-size: 0;
    }
    .template-table {
        margin-top: 16px;
    }
</style>
