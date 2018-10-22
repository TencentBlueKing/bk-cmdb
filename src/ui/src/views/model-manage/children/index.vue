<template>
    <div class="model-detail-wrapper">
        <div></div>
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
