<template>
    <bk-dialog
        ext-cls="error-content-dialog"
        v-model="isModalShow"
        width="600"
        :mask-close="false"
        @cancel="onCloseDialog">
        <div class="permission-content">
            <div class="permission-header">
                <span class="title-icon">
                    <img src="../../assets/images/lock-closed.svg" class="locked-icon" alt="locked-icon" />
                </span>
                <h3>{{i18n.permissionTitle}}</h3>
            </div>
            <bk-table
                :data="list"
                :max-height="180"
                :border="true">
                <bk-table-column prop="scope" :label="$t('资源所属')"></bk-table-column>
                <bk-table-column prop="resource" :label="$t('资源')">
                    <template slot-scope="{ row }">
                        <div class="resource-list" v-html="row.resource"></div>
                    </template>
                </bk-table-column>
                <bk-table-column prop="action" :label="$t('需要申请的权限')"></bk-table-column>
            </bk-table>
        </div>
        <div class="permission-footer" slot="footer">
            <bk-button theme="primary" @click="handleApplyPermission">{{ i18n.apply }}</bk-button>
            <bk-button theme="default" @click="onCloseDialog">{{ i18n.cancel }}</bk-button>
        </div>
    </bk-dialog>
</template>
<script>
    export default {
        name: 'permissionModal',
        props: {},
        data () {
            return {
                isModalShow: false,
                list: [],
                i18n: {
                    permissionTitle: this.$t('没有权限访问或操作此资源'),
                    system: this.$t('系统'),
                    resource: this.$t('资源'),
                    requiredPermissions: this.$t('需要申请的权限'),
                    noData: this.$t('无数据'),
                    apply: this.$t('去申请'),
                    cancel: this.$t('取消')
                }
            }
        },
        methods: {
            show (list) {
                this.list = list
                this.isModalShow = true
            },
            handleApplyPermission () {
                const topWindow = window.top
                const isPaasConsole = topWindow !== window
                const authCenter = window.Site.authCenter || {}
                if (isPaasConsole) {
                    topWindow.postMessage(JSON.stringify({
                        action: 'open_other_app',
                        app_code: authCenter.appCode,
                        app_url: 'apply-by-system'
                    }), '*')
                } else {
                    window.open(authCenter.url)
                }
            },
            onCloseDialog () {
                this.isModalShow = false
            }
        }
    }
</script>
<style lang="scss" scoped>
    .permission-content {
        margin-top: -26px;
        .permission-header {
            text-align: center;
            .locked-icon {
                height: 60px;
            }
            h3 {
                margin: 10px 0 30px;
                color: #979ba5;
                font-size: 24px;
            }
        }
    }
    .resource-list {
        padding: 12px 0;
        word-break: break-all;
        white-space: normal;
    }
</style>
