<template>
    <div class="display-wrapper">
        <div class="display-setting">
            <label class="cmdb-form-checkbox">
                <input type="checkbox">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["显示模型名称"]')}}</span>
            </label>
            <label class="cmdb-form-checkbox">
                <input type="checkbox">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["显示关系名称"]')}}</span>
            </label>
        </div>
        <ul class="display-list">
            <li class="group-item">
                <p class="group-name">主机</p>
                <ul class="clearfix">
                    <li class="model-item">
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox">
                            <span class="cmdb-checkbox-text">{{$t('ModelManagement["主机"]')}}</span>
                            <span class="count">(65)</span>
                            <i class="bk-icon icon-angle-down"></i>
                        </label>
                        <div class="relation-detail">
                            asdf
                        </div>
                    </li>
                    <li class="model-item">
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox">
                            <span class="cmdb-checkbox-text">{{$t('ModelManagement["主机"]')}}</span>
                            <!-- <div class="count-box">
                                <span class="count">(65)</span>
                                <div class="relation-detail">
                                    asdf
                                </div>
                            </div> -->
                            <i class="bk-icon icon-angle-down"></i>
                        </label>
                    </li>
                </ul>
            </li>
        </ul>
    </div>
</template>

<style lang="scss" scoped>
    .display-wrapper {
        padding: 20px 30px;
    }
    .display-setting {
        margin-bottom: 30px;
        .cmdb-form-checkbox {
            width: 154px;
        }
    }
    .display-list {
        .group-item {
            position: relative;
            .group-name {
                margin-bottom: 15px;
                padding-left: 8px;
                border-left: 4px solid $cmdbBorderFocusColor;
                line-height: 16px;
            }
        }
        .model-item {
            float: left;
            width: 140px;
            .cmdb-form-checkbox {
                margin-right: 5px;
                width: 135px;
                font-size: 0;
                cursor: pointer;
            }
            .cmdb-checkbox-text {
                max-width: 65px;
                font-size: 14px;
                @include ellipsis;
            }
            .count {
                position: relative;
                font-size: 14px;
                vertical-align: middle;
                color: $cmdbBorderFocusColor;
                &:after {
                    position: absolute;
                    content: '';
                    border: 5px solid transparent;
                    border-bottom-color: #fff;
                    top: 17px;
                    left: calc(50% - 5px);
                    z-index: 1;
                }
                &:before {
                    position: absolute;
                    content: '';
                    border: 5px solid transparent;
                    border-bottom-color: $cmdbTableBorderColor;
                    top: 16px;
                    left: calc(50% - 5px);
                    z-index: 1;
                }
            }
            .relation-detail {
                position: absolute;
                width: 455px;
                padding: 20px;
                background: #fff;
                border: 1px solid $cmdbTableBorderColor;
            }
            .icon-angle-down {
                margin-left: 5px;
                font-size: 12px;
                color: $cmdbBorderColor;
            }
        }
    }
</style>
