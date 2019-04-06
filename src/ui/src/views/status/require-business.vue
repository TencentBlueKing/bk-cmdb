<template>
    <div class="tips-wrapper">
        <div class="content-wrapper">
            <div class="title">
                <i class="icon icon-cc-no-authority"></i>
                <h2>{{$t("Common['无业务权限']")}}</h2>
            </div>
            <p>{{$t("Common['点击下方按钮申请']")}}</p>
            <div class="btns">
                <bk-button type="default">
                    {{$t("Common['申请业务权限']")}}
                </bk-button>
            </div>
        </div>
    </div>
</template>
<script>
    export default {
    }
</script>

<style lang="scss" scoped>
    .content-wrapper {
        margin-top: 100px;
        text-align: center;
        color: #63656E;
        font-size: 14px;
        .title {
            .icon {
                font-size: 56px;
                color: #979BA5;
            }
            h2 {
                margin-top: 10px;
                margin-bottom: 10px;
                font-size: 22px;
                color: #313238;
                font-weight: normal;
            }
        }
        .btns {
            margin-top: 24px;
            .bk-button {
                border-radius: 3px;
                padding-left: 10px;
                padding-right: 10px;
            }
        }
    }
</style>
