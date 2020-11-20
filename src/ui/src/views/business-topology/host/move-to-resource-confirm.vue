<template>
    <div class="confirm-wrapper">
        <h2 class="title">{{$t('确认归还主机池')}}</h2>
        <i18n tag="p" path="确认归还主机池主机数量" class="content">
            <span class="count" place="count">{{count}}</span>
        </i18n>
        <p class="content">{{$t('归还主机池提示')}}</p>
        <div class="options">
            <bk-button class="mr10" theme="primary" @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-move-to-resource-confirm',
        props: {
            count: {
                type: Number,
                default: 0
            }
        },
        methods: {
            handleConfirm () {
                this.$emit('confirm')
            },
            handleCancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .confirm-wrapper {
        text-align: center;
    }
    .title {
        margin: 45px 0 27px;
        line-height: 32px;
        font-size:24px;
        font-weight: normal;
        color: #313238;
    }
    .content {
        line-height:20px;
        font-size:14px;
        color: $textColor;
        .count {
            font-weight: bold;
            color: #EA3536;
            padding: 0 4px;
        }
    }
    .options {
        margin: 25px 0 30px 0;
        font-size: 0;
    }
</style>
