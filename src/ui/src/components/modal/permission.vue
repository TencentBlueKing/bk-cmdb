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
    export default {
        name: 'permissionModal',
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
                const permission = this.permission
                const list = permission.map(datum => {
                    const scope = [datum.scope_type_name]
                    if (datum.scope_id) {
                        scope.push(datum.scope_name)
                    }
                    let resource
                    if (datum.resource_type_name) {
                        resource = datum.resource_type_name
                    } else {
                        resource = datum.resources.map(resource => {
                            const resourceInfo = resource.map(info => this.getPermissionText(info, 'resource_type_name', 'resource_name'))
                            return [...new Set(resourceInfo)].join('\n')
                        }).join('\n')
                    }
                    return {
                        scope: this.getPermissionText(datum, 'scope_type_name', datum.scope_type === 'system' ? null : 'scope_name'),
                        resource: resource,
                        action: datum.action_name
                    }
                })
                const uniqueList = []
                list.forEach(item => {
                    const exist = uniqueList.some(unique => {
                        return item.resource === unique.resource
                            && item.scope === unique.scope
                            && item.action === unique.action
                    })
                    if (!exist) {
                        uniqueList.push(item)
                    }
                })
                this.list = uniqueList
            },
            getPermissionText (data, necessaryKey, extraKey, split = '：') {
                const text = [data[necessaryKey]]
                if (extraKey && data[extraKey]) {
                    text.push(data[extraKey])
                }
                return text.join(split).trim()
            },
            async handleApplyPermission () {
                try {
                    const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
                        params: this.permission,
                        config: {
                            requestId: 'getSkipUrl'
                        }
                    })
                    window.open(skipUrl)
                } catch (e) {
                    console.error(e)
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
