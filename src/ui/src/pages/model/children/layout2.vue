<template>
    <div class="tab-wrapper" v-bkloading="{isLoading: isLoading}">
        <div class="tab-box">
            <div class="hidden-list">
                <div class="hidden-list-title">
                    <i class="bk-icon icon-eye-slash-shape"></i>
                    隐藏字段
                </div>
                <ul>
                    <draggable class="content-left">
                        <li v-for="item in hideField">
                            <span class="hidden-list-icon">
                                <i></i><i></i><i></i>
                            </span>
                            <span class="hidden-list-text">
                                {{item['bk_property_name']}}
                            </span>
                        </li>
                    </draggable>
                </ul>
            </div>
            <div class="layout-list">
                <div class="layout-list-ul">
                    <div class="layout-list-title">
                        <span class="layout-title-text">sdaf</span>
                        <input type="text" class="layout-title-text border">
                        <i class="icon-cc-edit"></i>
                        <span class="layout-title-icon">
                            <i class="bk-icon icon-arrows-up"></i>
                            <i class="bk-icon icon-arrows-down"></i>
                            <i class="icon-cc-del f14 vm"></i>
                        </span>
                    </div>
                    <ul>
                        <draggable class="content-right" :options="{animation: 150, group:'field'}">
                            <li>
                                <span class="layout-list-icon">
                                    <i></i><i></i><i></i>
                                </span>
                                <span class="layout-list-text">asdf</span>
                                <i class="bk-icon icon-eye-slash-shape"></i>
                            </li>
                        </draggable>
                    </ul>
                </div>
                <div class="layout-list-add">
                    <span>
                        <i class="bk-icon icon-plus"></i>
                        <span>新增字段分组</span>
                    </span>
                </div>
            </div>
        </div>
        <div class="base-info">
            <bk-button class="btn main-btn" type="primary" title="确认">确认</bk-button>
            <bk-button class="btn vice-btn cancel-btn-sider" type="default" title="取消">取消</bk-button>
        </div>
    </div>
</template>

<script>
    import draggable from 'vuedraggable'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            isShow: {
                default: false
            },
            activeModel: {
                type: Object
            }
        },
        data () {
            return {
                isLoading: false,
                hideField: [],
                allField: []
            }
        },
        conputed: {
            ...mapGetters([
                'bkSupplierAccount'
            ])
        },
        methods: {
            async getFieldGroup () {
                try {
                    let fieldGroup = await this.$axios.post(`objectatt/group/property/owner/${this.bkSupplierAccount}/object/${this.activeModel['bk_obj_id']}`)
                    this.getAttr()
                } catch (e) {
                    console.error(e)
                    this.$alertMsg(e.data['bk_error_msg'] || e.message || e.statusText)
                }
            },
            async getAttr () {
                try {
                    let params = {
                        bk_obj_id: this.activeModel['bk_obj_id'],
                        bk_supplier_account: this.bkSupplierAccount
                    }
                    let res = await this.$axios.post('object/attr/search', params)
                } catch (e) {
                    console.error(e)
                    this.$alertMsg(e.data['bk_error_msg'] || e.message || e.statusText)
                }
            }
        },
        directives: {
            focus: {
                update (el, {value}) {
                    if (value) {
                        el.focus()
                    }
                }
            }
        },
        components: {
            draggable
        }
    }
</script>

<style lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #f9f9f9; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; // 主要颜色
    .tab-wrapper{
        height: 100%;
        padding: 30px 30px 85px 30px;
        .tab-box{
            width: 100%;
            height: 100%;
            overflow-y: auto;
            border:1px solid $borderColor;
            .hidden-list{
                width: 143px;
                height: 100%;
                float: left;
                border-right:1px solid $borderColor;
                text-align: center;
                &-title{
                    background: $primaryColor;
                    line-height: 42px;
                    height: 46px;
                    padding-top: 4px;
                    font-weight: bold;
                    border-bottom:1px solid $borderColor;
                }
                ul{
                    height: calc(100% - 46px);
                    overflow-y: auto;
                    &::-webkit-scrollbar{
                        width: 6px;
                        height: 5px;
                    }
                    &::-webkit-scrollbar-thumb{
                        border-radius: 20px;
                        background: #a5a5a5;
                    }
                    .content-left{
                        height: 100%;
                    }
                    li{
                        cursor: move;
                        border-bottom:1px solid $borderColor;
                        font-size: 0;
                        position: relative;
                        -webkit-transition: all .35s;
                        transition: all .35s;
                        .hidden-list-text,
                        .layout-list-text{
                            font-size: 12px;
                            width: 100%;
                            display: inline-block;
                            overflow: hidden;
                            white-space: nowrap;
                            text-overflow: ellipsis;
                            height: 40px;
                            line-height: 40px;
                            padding: 0 10px 0 15px;
                        }
                        .hidden-list-icon,
                        .layout-list-icon{
                            position: absolute;
                            left: 6px;
                            top: 14px;
                            width: 3px;
                            height: 15px;
                            overflow: hidden;
                            display: inline-block;
                            i{
                                display: inline-block;
                                width: 3px;
                                height: 3px;
                                background: $borderColor;
                                margin: 1px 0;
                            }
                        }
                        &:hover{
                            box-shadow: 0 0 6px #eeeeee;
                            .hidden-list-icon,
                            .layout-list-icon{
                                i{
                                    background: $primaryHoverColor;
                                }
                            }
                        }
                    }
                }
            }
            .layout-list{
                width: calc(100% - 143px);
                float: right;
                padding: 0 20px;
                height: 100%;
                overflow-y: auto;
                &::-webkit-scrollbar{
                    width: 6px;
                    height: 5px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 20px;
                    background: #a5a5a5;
                }
                .content-right{
                    min-height: 30px;
                }
                &-icon{
                    position: absolute;
                    left: 6px;
                    top: 7px;
                    width: 3px;
                    height: 15px;
                    overflow: hidden;
                    display: none;
                    i{
                        display: inline-block;
                        width: 3px;
                        height: 3px;
                        background: $primaryHoverColor;
                        margin: 1px 0;
                        float: left;
                    }
                }
                &-ul{
                    >ul{
                        width: 100%;
                        padding: 11px 0;
                        font-size: 0;
                        li{
                            position: relative;
                            width: 50%;
                            height: 30px;
                            display: inline-block;
                            overflow: hidden;
                            white-space: nowrap;
                            text-overflow: ellipsis;
                            cursor: move;
                            .icon-eye-slash-shape{
                                display: none;
                                font-size: 12px;
                                position: absolute;
                                right: 12px;
                                top: 9px;
                                cursor: pointer;
                            }
                            &:hover{
                                background-color: #f1f7ff;
                                .icon-eye-slash-shape,
                                .layout-list-icon{
                                    display: inline-block;
                                }
                            }
                        }
                    }
                }
                &-add{
                    width: 100%;
                    border-top:1px solid $borderColor;
                    text-align: center;
                    color: #498fe0;
                    line-height: 36px;
                    margin-top: 18px;
                    .icon-plus{
                        font-size: 12px;
                        position: relative;
                        top: -1px;
                        cursor: pointer;
                    }
                    span{
                        display: inline-block;
                        cursor: pointer;
                    }
                }
                &-title{
                    line-height: 42px;
                    height: 46px;
                    padding-top: 4px;
                    border-bottom:1px dashed $borderColor;
                    .layout-title-text{
                        color: #c3cdd7;
                        width: auto;
                        font-weight: bold;
                        display: inline-block;
                        line-height: 24px;
                        height: 26px;
                        border:1px solid $defaultColor;
                        /* padding: 0 10px; */
                        &.border{
                            border:1px solid $fnMainColor;
                            padding: 0 10px;
                        }
                    }
                    .icon-cc-edit{
                        display: none;
                        cursor: pointer;
                    }
                    .layout-title-icon{
                        float: right;
                        /* display: none; */
                        i{
                            color: #b4c1e8;
                            opacity: 0.5;
                            cursor: pointer;
                            -webkit-transition: all .35s;
                            transition: all .35s;
                            padding: 5px 0;
                            &:hover{
                                color: $primaryHoverColor;
                                &.icon-cc-del{
                                    color: #f05d5d;
                                }
                            }
                        }
                    }
                    &:hover{
                        .icon-cc-edit,
                        .layout-title-icon{
                            display: inline-block;
                            i{
                                opacity: 1;
                            }
                        }
                    }
                }
                &-text{
                    font-size: 14px;
                    width: 100%;
                    display: inline-block;
                    overflow: hidden;
                    white-space: nowrap;
                    text-overflow: ellipsis;
                    height: 30px;
                    line-height: 30px;
                    padding: 0 32px 0 15px;
                }
            }
        }
    }
    .base-info{
        width: 100%;
        position: absolute;
        left: 0;
        bottom: 0;
        padding: 14px 10px;
        background: #f9f9f9;
        button{
            height: 36px;
            line-height: 34px;
            border-radius: 2px;
            display: inline-block;
            min-width: 110px;
            margin-left: 10px;
            -webkit-transition: all .35s !important;
            transition: all .35s !important;
        }
    }
</style>
