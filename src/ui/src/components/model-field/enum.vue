<template>
    <ul class="form-enum-wrapper">
        <li class="form-item clearfix">
            <div class="enum-default">
                <input type="radio">
            </div>
            <div class="enum-id">
                <input type="text"
                    class="cmdb-form-input"
                    :placeholder="$t('ModelManagement[\'请输入ID\']')">
            </div>
            <div class="enum-name">
                <input type="text"
                    class="cmdb-form-input"
                    :placeholder="$t('ModelManagement[\'请输入名称英文数字\']')">
            </div>
            <button>
                <i class="icon-cc-del"></i>
            </button>
            <button>
                <i class="icon-cc-plus"></i>
            </button>
        </li>
    </ul>
</template>

<style lang="scss" scoped>
    .form-enum-wrapper {
        .form-item {
            .enum-default {
                float: left;
                width: 80px;
                height: 30px;
                padding-right: 5px;
                font-size: 16px;
                text-align: right;
            }
            .enum-id {
                float: left;
                width: 90px;
                margin-right: 10px;
                input {
                    width: 100%;
                }
            }
            .enum-name {
                float: left;
                width: 250px;
                input {
                    width: 100%;
                }
            }
        }
    }
</style>
