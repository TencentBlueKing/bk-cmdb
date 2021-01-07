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
                :max-height="193"
                class="permission-table">
                <bk-table-column prop="name" :label="$t('需要申请的权限')"></bk-table-column>
                <bk-table-column prop="resource" :label="$t('关联的资源实例')">
                    <template slot-scope="{ row }">
                        <div v-if="row.relations.length" style="overflow: auto;">
                            <div class="permission-resource"
                                v-for="(relation, index) in row.relations"
                                v-bk-overflow-tips
                                :key="index">
                                {{relation}}
                            </div>
                        </div>
                        <span v-else>--</span>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="permission-footer" slot="footer">
            <template v-if="applied">
                <bk-button theme="primary" @click="handleRefresh">{{ i18n.applied }}</bk-button>
                <bk-button class="ml10" @click="onCloseDialog">{{ i18n.close }}</bk-button>
            </template>
            <template v-else>
                <bk-button theme="primary"
                    :loading="$loading('getSkipUrl')"
                    @click="handleApply">
                    {{ i18n.apply }}
                </bk-button>
                <bk-button class="ml10" @click="onCloseDialog">{{ i18n.cancel }}</bk-button>
            </template>
        </div>
    </bk-dialog>
</template>
<script>
    import permissionMixins from '@/mixins/permission'
    import { IAM_ACTIONS, IAM_VIEWS_NAME } from '@/dictionary/iam-auth'
    export default {
        name: 'permissionModal',
        mixins: [permissionMixins],
        props: {},
        data () {
            return {
                applied: false,
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
                    applied: this.$t('已完成'),
                    cancel: this.$t('取消'),
                    close: this.$t('关闭')
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
                this.applied = false
                this.isModalShow = true
            },
            setList () {
                const languageIndex = this.$i18n.locale === 'en' ? 1 : 0
                this.list = this.permission.actions.map(action => {
                    const { id: actionId, related_resource_types: relatedResourceTypes = [] } = action
                    const definition = Object.values(IAM_ACTIONS).find(definition => definition.id === actionId)
                    const allRelationPath = []
                    relatedResourceTypes.forEach(({ type, instances = [] }) => {
                        instances.forEach(fullPaths => {
                            const topoPath = fullPaths.map(pathData => {
                                if (pathData.name) {
                                    return `${IAM_VIEWS_NAME[pathData.type][languageIndex]}：${pathData.name}`
                                }
                                return `${IAM_VIEWS_NAME[pathData.type][languageIndex]}ID：${pathData.id}`
                            }).join(' / ')
                            allRelationPath.push(topoPath)
                        })
                    })
                    return {
                        id: actionId,
                        name: definition.name[languageIndex],
                        relations: allRelationPath
                    }
                })
            },
            onCloseDialog () {
                this.isModalShow = false
            },
            async handleApply () {
                try {
                    await this.handleApplyPermission()
                    this.applied = true
                } catch (error) {}
            },
            handleRefresh () {
                window.location.reload()
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
    .permission-table {
        .permission-resource {
            line-height: 24px;
        }
        /deep/ {
            .bk-table-row {
                td.is-first {
                    vertical-align: top;
                    line-height: 42px;
                }
            }
        }
    }
</style>
