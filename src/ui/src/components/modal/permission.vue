<template>
    <bk-dialog
        ext-cls="error-content-dialog"
        v-model="isModalShow"
        width="740"
        :z-index="2400"
        :close-icon="false"
        :mask-close="false"
        @cancel="onCloseDialog">
        <div class="permission-content">
            <div class="permission-header">
                <span class="title-icon">
                    <img src="../../assets/images/lock-closed02.svg" class="locked-icon" alt="locked-icon" />
                </span>
                <h3>{{i18n.permissionTitle}}</h3>
            </div>
            <bk-table ref="table"
                :data="list"
                :max-height="193">
                <bk-table-column prop="scope" :label="$t('资源所属')"></bk-table-column>
                <bk-table-column prop="resource" :label="$t('资源')">
                    <template slot-scope="{ row }">
                        <div v-html="row.resource"></div>
                    </template>
                </bk-table-column>
                <bk-table-column prop="action" :label="$t('需要申请的权限')"></bk-table-column>
            </bk-table>
        </div>
        <div class="permission-footer" slot="footer">
            <bk-button theme="primary"
                :loading="$loading('getSkipUrl')"
                @click="handleApplyPermission">
                {{ i18n.apply }}
            </bk-button>
            <bk-button theme="default" @click="onCloseDialog">{{ i18n.cancel }}</bk-button>
        </div>
    </bk-dialog>
</template>
<script>
    import permissionMixins from '@/mixins/permission'
    export default {
        name: 'permissionModal',
        mixins: [permissionMixins],
        props: {},
        data () {
            return {
                isModalShow: false,
                permission: [],
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
        watch: {
            isModalShow (val) {
                if (val) {
                    setTimeout(() => {
                        this.$refs.table.doLayout()
                    }, 0)
                }
            }
        },
        methods: {
            show (permission) {
                this.permission = permission
                this.setList()
                this.isModalShow = true
            },
            setList () {
                this.list = this.$getPermissionList(this.permission)
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
            padding-top: 34px;
            text-align: center;
            .locked-icon {
                height: 66px;
            }
            h3 {
                margin: 6px 0 30px;
                color: #63656e;
                font-size: 24px;
                font-weight: normal;
            }
        }
    }
</style>
