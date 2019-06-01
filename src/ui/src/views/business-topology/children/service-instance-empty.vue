<template>
    <div class="empty">
        <div class="empty-content">
            <img src="../../../../static/svg/cc-empty.svg" alt="">
            <p class="empty-text">Agent模版尚未定义进程，无法创建服务</p>
            <p class="empty-tips">您可以先跳转模版添加进程或直接添加主机，后续再添加模版进程</p>
            <div class="empty-options">
                <bk-button class="empty-button" type="primary">跳转模板添加进程</bk-button>
                <bk-button class="empty-button" type="default">添加主机</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {}
        }
    }
</script>

<style lang="scss" scoped>
    .empty{
        height: 100%;
        text-align: center;
        &:before {
            content: "";
            width: 0;
            height: 100%;
            @include inlineBlock;
        }
        .empty-content {
            @include inlineBlock;
        }
        .empty-text {
            margin-top: 20px;
            line-height: 29px;
            font-size: 22px;
            color: #63656E;
        }
        .empty-tips {
            margin-top: 10px;
            font-size: 15px;
            line-height: 20px;
        }
        .empty-options {
            margin-top: 25px;
            .empty-button {
                height: 32px;
                margin: 0 4px;
                line-height: 30px;
            }
        }
    }
</style>
