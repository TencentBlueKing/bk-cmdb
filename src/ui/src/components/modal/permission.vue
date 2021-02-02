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
        <permission-main ref="main" :permission="permission" :applied="applied"
            @close="onCloseDialog"
            @apply="handleApply"
            @refresh="handleRefresh" />
    </bk-dialog>
</template>
<script>
    import permissionMixins from '@/mixins/permission'
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
                permission: {
                    actions: []
                }
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
                this.applied = false
                this.isModalShow = true
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
