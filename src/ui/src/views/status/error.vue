<template>
    <div class="error-layout">
        <div class="error-info">
            <img src="../../assets/images/error.png">
            <p class="text">系统出现异常, 请记录下错误场景并与开发人员联系, 谢谢!</p>
        </div>
    </div>
</template>

<style lang="scss" scoped>
    .error-layout{
        height: 100%;
        text-align: center;
        &:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
        }
        .error-info{
            display: inline-block;
            vertical-align: middle;
            .text{
                padding: 10px 0;
            }
        }
    }
</style>