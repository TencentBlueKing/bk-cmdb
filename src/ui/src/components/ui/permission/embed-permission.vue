<template>
    <div class="embed-permission">
        <div class="content">
            <img class="img-empty" src="../../../assets/images/no-authority.png" alt="">
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
        .img-empty {
            width: 90px;
        }
        .text-btn {
            font-size: 14px;
            height: auto;
        }
    }
</style>
