<template>
    <div class="embed-permission">
        <div class="content">
            <bk-exception type="403">
                <div>
                    <i18n path="抱歉您没有查看权限">
                        <bk-button class="text-btn"
                            place="link"
                            text
                            theme="primary"
                            @click="handleApplyPermission"
                        >
                            {{$t('去申请')}}
                        </bk-button>
                    </i18n>
                </div>
            </bk-exception>
        </div>
    </div>
</template>
<script>
    import permissionMixins from '@/mixins/permission'
    export default {
        name: 'cmdb-permission',
        mixins: [permissionMixins],
        props: {
            permission: {
                type: Array,
                default: () => []
            }
        }
    }
</script>

<style lang="scss" scoped>
    .embed-permission {
        color: #63656e;
        font-size: 14px;
        text-align: center;
        .text-btn {
            font-size: 14px;
            height: auto;
        }
    }
</style>
