<template>
    <bk-dialog
        ext-cls="permission-dialog"
        v-model="isModalShow"
        width="740"
        :z-index="2400"
        :close-icon="false"
        :mask-close="false"
        :show-footer="false"
        @cancel="onCloseDialog">
        <permission-main ref="main" :list="list" :applied="applied"
            @close="onCloseDialog"
            @apply="handleApply"
            @refresh="handleRefresh" />
    </bk-dialog>
</template>
<script>
    import permissionMixins from '@/mixins/permission'
    import { IAM_ACTIONS, IAM_VIEWS_NAME } from '@/dictionary/iam-auth'
    import PermissionMain from './permission-main.vue'
    export default {
        name: 'permissionModal',
        components: {
            PermissionMain
        },
        mixins: [permissionMixins],
        props: {},
        data () {
            return {
                applied: false,
                isModalShow: false,
                permission: [],
                list: []
            }
        },
        watch: {
            isModalShow (val) {
                if (val) {
                    setTimeout(() => {
                        this.$refs.main.doTableLayout()
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
    /deep/ .permission-dialog {
        .bk-dialog-body {
            padding: 0;
        }
    }
</style>
