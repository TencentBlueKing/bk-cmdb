<template>
    <div class="model-detail-wrapper">
        <div class="model-info">
            <div class="icon-box">
                <i class="icon-cc-host"></i>
            </div>
            <div>
                <span>唯一标识：</span>
                <span>asdf</span>
            </div>
            <div>
                <span>名称：</span>
                <span>asdf</span>
            </div>
            <div class="btn-group">
                <label class="label-btn">
                    <i class="icon-cc-derivation"></i>
                    <span>{{$t('ModelManagement["导出"]')}}</span>
                </label>
                <label class="label-btn">
                    <i class="icon-cc-copy"></i>
                    <span>{{$t('Common["复制"]')}}</span>
                </label>
                <label class="label-btn">
                    <i class="bk-icon icon-minus-circle-shape"></i>
                    <span>{{$t('ModelManagement["停用"]')}}</span>
                </label>
                <label class="label-btn">
                    <i class="icon-cc-del"></i>
                    <span>{{$t("Common['删除']")}}</span>
                </label>
            </div>
        </div>
        <bk-tab :active-name.sync="tab.active">
            <bk-tabpanel name="field" :title="$t('ModelManagement[\'模型字段\']')">
            </bk-tabpanel>
            <bk-tabpanel name="relation" :title="$t('ModelManagement[\'模型关系\']')">
            </bk-tabpanel>
            <bk-tabpanel name="verification" :title="$t('ModelManagement[\'唯一校验\']')">
            </bk-tabpanel>
            <bk-tabpanel name="layout" :title="$t('ModelManagement[\'字段分组\']')">
            </bk-tabpanel>
            <bk-tabpanel name="history" :title="$t('ModelManagement[\'操作历史\']')">
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                tab: {
                    active: 'field'
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-detail-wrapper {
        padding: 0;
    }
    .model-info {
        height: 100px;
        background: rgba(235, 244, 255, .6);
        .btn-group {
            float: right;
            i,
            span {
                vertical-align: middle;
            }
        }
    }
</style>
