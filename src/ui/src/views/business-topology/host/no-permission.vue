<template>
    <permission-main class="no-permission" ref="main" :permission="permission" :applied="applied"
        @close="handleClose"
        @apply="handleApply"
        @refresh="handleRefresh" />
</template>
<script>
    import permissionMixins from '@/mixins/permission'
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
                applied: false
            }
        },
        methods: {
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
