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
            <table class="permission-table table-header">
                <thead>
                    <tr>
                        <th width="30%">{{i18n.system}}</th>
                        <th width="35%">{{i18n.resource}}</th>
                        <th width="35%">{{i18n.requiredPermissions}}</th>
                    </tr>
                </thead>
            </table>
            <div class="table-content">
                <table class="permission-table">
                    <tbody>
                        <tr>
                            <td width="30%">作业平台</td>
                            <td width="35%">作业：gse初始化作业作业：gse初始化作业作业：gse初始化作业</td>
                            <td width="35%">查看，执行</td>
                        </tr>
                        <tr>
                            <td>作业平台</td>
                            <td>快速脚本执行作业：gse初始化作业</td>
                            <td>执行</td>
                        </tr>
                        <tr>
                            <td>作业平台</td>
                            <td>快速脚本执行作业：gse初始化作业</td>
                            <td>执行</td>
                        </tr>
                        <tr>
                            <td>作业平台</td>
                            <td>快速脚本执行作业：gse初始化作业</td>
                            <td>执行</td>
                        </tr>
                        <tr v-if="false">
                            <td class="no-data" colspan="3">{{i18n.noData}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="permission-footer" slot="footer">
            <bk-button type="primary" @click="goToApply">{{ i18n.apply }}</bk-button>
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
            show (data) {
                this.isModalShow = true
                this.list = data
            },
            goToApply () {

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
        .permission-table {
            width: 100%;
            height: 100px;
            color: #63656e;
            border: 1px solid #dcdee5;
            border-collapse: collapse;
            table-layout: fixed;
            th,
            td {
                padding: 12px 18px;
                font-size: 12px;
                text-align: left;
                border-bottom: 1px solid #dcdee5;
                border-right: 1px solid #dcdee5;
            }
            th {
                color: #313238;
                background: rgb(250, 251, 253);
            }
        }
        .table-content {
            max-height: 180px;
            border-bottom: 1px solid #dcdee5;
            border-top: none;
            overflow: auto;
            .permission-table {
                border-top: none;
                border-bottom: none;
                td:last-child {
                    border-right: none;
                }
                tr:last-child td {
                    border-bottom: none;
                }
            }
            .no-data {
                text-align: center;
                color: #999999;
            }
        }
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
