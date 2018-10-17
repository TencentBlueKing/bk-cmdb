<template>
    <div class="model-icon-list" v-show="isIconListShow">
        <span class="back" @click="isIconListShow = false">
            <i class="bk-icon icon-back2"></i>
        </span>
        <ul class="icon-box clearfix">
            <li class="icon" 
                v-for="(icon, index) in iconList" 
                :key="index">
                <i :class="icon.value"></i>
            </li>
        </ul>
        <div class="page">
            <input type="text" class="cmdb-form-input">
            <div class="page-btn">
                <bk-button type="default">
                    <i class="bk-icon icon-angle-left"></i>
                </bk-button>
                <bk-button type="default">
                    <i class="bk-icon icon-angle-right"></i>
                </bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import iconList from '@/assets/json/model-icon.json'
    export default {
        data () {
            return {
                iconList: []
            }
        },
        created () {
            this.iconList = iconList.slice(0, 28)
        }
    }
</script>

<style lang="scss" scoped>
    .icon-box {
        padding: 10px 20px;
        height: 236px;
        .icon {
            float: left;
            width: 80px;
            height: 49px;
            padding: 5px;
            font-size: 30px;
            text-align: center;
            margin-bottom: 10px;
            cursor: pointer;
            &:hover,
            &.active {
                background: #e2efff;
                color: #3c96ff;
            }
        }
        .page {
            height: 52px;
            padding: 10px 20px;
            .cmdb-form-input {
                float: left;
                width: 200px;
                height: 30px;
            }
            .page-btn {
                float: right;
                .bk-button {
                    padding: 0;
                    width: 30px;
                    height: 30px;
                    line-height: 1;
                    vertical-align: top;
                }
            }
        }
    }
</style>
