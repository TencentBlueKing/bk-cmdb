<template>
    <div class="group-wrapper">
        <p class="btn-group">
            <bk-button type="primary" @click="createGroup">
                {{$t('ModelManagement["新建分组"]')}}
            </bk-button>
            <bk-button type="default" @click="createModel">
                {{$t('ModelManagement["新增模型"]')}}
            </bk-button>
        </p>
        <ul class="group-list">
            <li class="group-item clearfix">
                <p class="group-title">
                    <span>主机</span><span class="number">(8)</span>
                </p>
                <ul class="model-list">
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                </ul>
                <i class="bk-icon icon-angle-double-down"></i>
            </li>
            <li class="group-item clearfix">
                <p class="group-title">
                    <span>主机</span><span class="number">(8)</span>
                </p>
                <ul class="model-list">
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                    <li class="model-item">
                        <div class="icon-box">
                            <i class="icon icon-cc-host"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">数据库</p>
                            <p class="model-id">adsf</p>
                        </div>
                    </li>
                </ul>
            </li>
        </ul>
        <bk-dialog
            class="group-dialog"
            :close-icon="false"
            :hasHeader="false"
            :width="600"
            :padding="0"
            :is-show.sync="groupDialog.isShow">
            <div slot="content" class="dialog-content">
                <p class="title">{{$t('ModelManagement["新增模型"]')}}</p>
                <label for="">
                    <span>
                        {{$t('ModelManagement["唯一标识"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                </label>
                <label for="">
                    <span>
                        {{$t('ModelManagement["名称"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                </label>
            </div>
            <div slot="footer" class="footer">
                <bk-button type="primary">保存</bk-button>
                <bk-button type="default">取消</bk-button>
            </div>
        </bk-dialog>
        <bk-dialog
            class="model-dialog"
            :close-icon="false"
            :hasHeader="false"
            :width="600"
            :padding="0"
            :is-show.sync="modelDialog.isShow">
            <div slot="content" class="dialog-content">
                <p class="title">{{$t('ModelManagement["新增模型"]')}}</p>
                <div class="content clearfix">
                    <div class="content-left" @click="isIconListShow = true">
                        <div class="icon-wrapper">
                            <i class="icon-cc-host"></i>
                        </div>
                        <div class="text">{{$t('ModelManagement["点击切换"]')}}</div>
                    </div>
                    <ul class="content-right">
                        <li>
                            <label for="">{{$t('ModelManagement["所属分组"]')}}<span class="color-danger">*</span></label>
                            <bk-selector
                                class="selector-box"
                                :selected.sync="selected"
                                :list="list"
                            ></bk-selector>
                        </li>
                        <li>
                            <label for="">{{$t('ModelManagement["唯一标识"]')}}<span class="color-danger">*</span></label>
                            <input type="text" class="cmdb-form-input">
                            <i class="bk-icon icon-info-circle"></i>
                        </li>
                        <li>
                            <label for="">{{$t('ModelManagement["名称"]')}}<span class="color-danger">*</span></label>
                            <input type="text" class="cmdb-form-input">
                            <i class="bk-icon icon-info-circle"></i>
                        </li>
                    </ul>
                </div>
                <div v-if="isIconListShow">
                    <v-choose-icon
                    ></v-choose-icon>
                </div>
            </div>
            <div slot="footer" class="footer">
                <bk-button type="primary">保存</bk-button>
                <bk-button type="default">取消</bk-button>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import vChooseIcon from './_choose-icon'
    export default {
        components: {
            vChooseIcon
        },
        data () {
            return {
                groupDialog: {
                    isShow: false
                },
                modelDialog: {
                    isShow: false
                },
                isIconListShow: false,
                list: [],
                selected: '',
                iconList: []
            }
        },
        created () {
        },
        methods: {
            createModel () {
                this.modelDialog.isShow = true
            },
            createGroup () {
                this.groupDialog.isShow = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .btn-group {
        margin: 20px 0;
        font-size: 0;
        .bk-primary {
            margin-right: 10px;
        }
    }
    .group-list {
        .group-item {
            position: relative;
            padding: 10px 0 20px;
            >.icon-angle-double-down {
                position: absolute;
                left: 50%;
                bottom: -5px;
                margin-left: -5px;
                padding: 5px;
                font-size: 12px;
                cursor: pointer;
                transition: all .2s;
                &.rotate {
                    transform: translate(180deg);
                }
            }
        }
        .group-title {
            padding-left: 8px;
            border-left: 4px solid $cmdbBorderFocusColor;
            line-height: 16px;
            color: #333948;
            .number {
                color: $cmdbBorderColor;
            }
        }
    }
    .model-list {
        padding-left: 12px;
        .model-item {
            float: left;
            margin: 10px 10px 0 0;
            width: calc((100% - 30px) / 4);
            height: 70px;
            border: 1px solid $cmdbTableBorderColor;
            border-radius: 4px;
            cursor: pointer;
            &:nth-child(4n) {
                margin-right: 0;
            }
            &:hover {
                border-color: $cmdbBorderFocusColor;
                background: #ebf4ff;
            }
            .icon-box {
                float: left;
                .icon {
                    padding-left: 18px;
                    font-size: 32px;
                    line-height: 70px;
                    color: $cmdbBorderFocusColor;
                }
            }
            .model-details {
                float: left;
                line-height: 16px;
                margin-top: 20px;
                padding-left: 5px;
            }
            .model-name {
                font-size: 14px;
            }
            .model-id {
                font-size: 12px;
                color: $cmdbTableBorderColor;
            }
        }
    }
    .model-dialog {
        .dialog-content {
            position: relative;
            padding: 20px 10px;
        }
        .title {
            margin-bottom: 30px;
            font-size: 20px;
            color: #333948;
            line-height: 1;
        }
        .content {
            padding: 0 10px;
        }
        .content-left {
            float: left;
            width: 93px;
            height: 100px;
            border: 1px solid #dde4eb;
            border-radius: 4px 4px 0 0;
            cursor: pointer;
            .icon-wrapper {
                width: 100%;
                height: 68px;
                font-size: 38px;
                text-align: center;
                i {
                    vertical-align: top;
                    line-height: 68px;
                    color: $cmdbBorderFocusColor;
                }
            }
            .text {
                height: 30px;
                border-top: 1px solid #dde4eb;
                text-align: center;
                line-height: 30px;
                background: #ebf4ff;
            }
        }
        .content-right {
            float: right;
            width: 460px;
            font-size: 0;
            >li {
                margin-bottom: 20px;
                &:last-child {
                    margin: 0;
                }
                .icon-info-circle {
                    font-size: 18px;
                    color: $cmdbBorderColor;
                    line-height: 36px;
                    vertical-align: middle;
                }
            }
            label {
                display: inline-block;
                width: 100px;
                text-align: right;
                font-size: 16px;
                line-height: 36px;
                vertical-align: middle;
                @include ellipsis;
                .color-danger {
                    padding: 0 5px;
                }
            }
            .selector-box,
            .cmdb-form-input {
                display: inline-block;
                margin-right: 10px;
                width: 330px;
                vertical-align: middle;
            }
        }
        .model-icon-list {
            position: absolute;
            left: 0;
            top:0;
            width: 100%;
            height: calc(100% + 60px);
            background: #fff;
            .back {
                position: absolute;
                right: -47px;
                top: 0;
                width: 44px;
                height: 44px;
                padding: 7px;
                cursor: pointer;
                font-size: 18px;
                text-align: center;
                background: #2f2f2f;
                color: #fff;
            }
        }
        .footer {
            padding: 0 24px;
            font-size: 0;
            text-align: right;
            .bk-primary {
                margin-right: 10px;
            }
        }
    }
</style>
