<template>
    <bk-dialog
        ext-cls="error-content-dialog"
        :is-show="isModalShow"
        :title="' '"
        :width="'600'"
        padding="0 24px 40px 24px"
        :has-header="true"
        :has-footer="true"
        :quick-close="false"
        :close-icon="true"
        @cancel="onCloseDialog">
        <div class="permission-content" slot="content">
            <div class="permission-header">
                <span class="title-icon">
                    <img src="../../assets/images/lock-closed.svg" class="locked-icon" alt="locked-icon" />
                </span>
                <h3>{{i18n.permissionTitle}}</h3>
            </div>
            <cmdb-table
                :header="header"
                :list="list"
                :max-height="180"
                :empty-height="140"
                :visible="isModalShow"
                :sortable="false">
                <template slot="resource" slot-scope="{ item }">
                    <div class="resouce-list" v-html="item.resource"></div>
                </template>
            </cmdb-table>
        </div>
        <div class="permission-footer" slot="footer">
            <bk-button type="primary" @click="handleApplyPermission">{{ i18n.apply }}</bk-button>
            <bk-button type="default" @click="onCloseDialog">{{ i18n.cancel }}</bk-button>
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
                header: [{
                    id: 'scope',
                    name: this.$t('资源所属')
                }, {
                    id: 'resource',
                    name: this.$t('资源')
                }, {
                    id: 'action',
                    name: this.$t('需要申请的权限')
                }],
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
    .resouce-list {
        padding: 12px 0;
        word-break: break-all;
        white-space: normal;
    }
    /deep/ .bk-dialog-footer.bk-d-footer {
        height: 50px;
        line-height: 50px;
        .permission-footer {
            padding: 0 24px;
            text-align: right;
        }
        .bk-button {
            height: 32px;
            line-height: 30px;
        }
    }
    
</style>
