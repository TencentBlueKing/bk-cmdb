<template>
    <permission-main class="no-permission" ref="main" :list="list" :applied="applied"
        @close="handleClose"
        @apply="handleApply"
        @refresh="handleRefresh" />
</template>
<script>
    import permissionMixins from '@/mixins/permission'
    import { IAM_ACTIONS, IAM_VIEWS_NAME } from '@/dictionary/iam-auth'
    import PermissionMain from '@/components/modal/permission-main'
    export default {
        components: {
            PermissionMain
        },
        mixins: [permissionMixins],
        props: {
            permission: Object
        },
        data () {
            return {
                list: [],
                applied: false
            }
        },
        watch: {
            permission () {
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
            handleClose () {
                this.$emit('cancel')
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
    .no-permission {
        height: var(--height, 600px);
        padding: 0 0 50px;
        position: relative;
        /deep/ .permission-content {
            padding: 16px 24px 0;
            margin: 0;
            height: 100%;
        }
        /deep/ .permission-footer {
            position: sticky;
            bottom: 0;
            left: 0;
            width: 100%;
            height: 50px;
            padding: 8px 20px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
            text-align: right;
            font-size: 0;
            z-index: 100;
        }
    }
</style>
