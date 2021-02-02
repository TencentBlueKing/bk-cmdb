<template>
    <div class="permission-main">
        <div class="permission-content">
            <div class="permission-header">
                <bk-exception type="403" scene="part">
                    <h3>{{i18n.permissionTitle}}</h3>
                </bk-exception>
            </div>
            <bk-table ref="table"
                :data="list"
                :max-height="193"
                class="permission-table">
                <bk-table-column prop="name" :label="$t('需要申请的权限')" width="250"></bk-table-column>
                <bk-table-column prop="resource" :label="$t('关联的资源实例')">
                    <template slot-scope="{ row }">
                        <div v-if="row.relations.length" style="overflow: auto;">
                            <div class="permission-resource"
                                v-for="(relation, index) in row.relations"
                                v-bk-overflow-tips
                                :key="index">
                                <permission-resource-name :relations="relation" />
                            </div>
                        </div>
                        <span v-else>--</span>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="permission-footer">
            <template v-if="applied">
                <bk-button theme="primary" @click="handleRefresh">{{ i18n.applied }}</bk-button>
                <bk-button class="ml10" @click="handleClose">{{ i18n.close }}</bk-button>
            </template>
            <template v-else>
                <bk-button theme="primary"
                    :loading="$loading('getSkipUrl')"
                    @click="handleApply">
                    {{ i18n.apply }}
                </bk-button>
                <bk-button class="ml10" @click="handleClose">{{ i18n.cancel }}</bk-button>
            </template>
        </div>
    </div>
</template>
<script>
    import { IAM_ACTIONS, IAM_VIEWS_NAME } from '@/dictionary/iam-auth'
    import PermissionResourceName from './permission-resource-name.vue'
    export default {
        components: {
            PermissionResourceName
        },
        props: {
            permission: Object,
            applied: Boolean
        },
        data () {
            return {
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
            permission (v) {
                this.setList()
            }
        },
        created () {
            this.setList()
        },
        methods: {
            setList () {
                const languageIndex = this.$i18n.locale === 'en' ? 1 : 0
                this.list = this.permission.actions.map(action => {
                    const { id: actionId, related_resource_types: relatedResourceTypes = [] } = action
                    const definition = Object.values(IAM_ACTIONS).find(definition => definition.id === actionId)
                    const allRelationPath = []
                    relatedResourceTypes.forEach(({ type, instances = [] }) => {
                        instances.forEach(fullPaths => {
                            // 数据格式[type, id, label]
                            const topoPath = fullPaths.map(pathData => [pathData.type, pathData.id, IAM_VIEWS_NAME[pathData.type][languageIndex]])
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
            handleClose () {
                this.$emit('close')
            },
            handleApply () {
                this.$emit('apply')
            },
            handleRefresh () {
                this.$emit('refresh')
            },
            doTableLayout () {
                this.$refs.table.doLayout()
            }
        }
    }
</script>
<style lang="scss" scoped>
    .permission-content {
        margin-top: -26px;
        padding: 3px 24px 26px;
        .permission-header {
            padding-top: 16px;
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

            /deep/ {
                .bk-exception-img .exception-image {
                    height: 130px;
                }
            }
        }
    }
    .permission-footer {
        text-align: right;
        padding: 12px 24px;
        background-color: #fafbfd;
        border-top: 1px solid #dcdee5;
        border-radius: 2px;
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
