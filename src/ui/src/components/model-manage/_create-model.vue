<template>
    <bk-dialog
        class="model-dialog dialog bk-dialog-no-padding"
        width="600"
        :close-icon="false"
        :mask-close="false"
        v-model="isShow">
        <div class="dialog-content">
            <p class="title">{{title}}</p>
            <div class="content clearfix">
                <div class="content-left">
                    <div class="icon-wrapper" @click="modelDialog.isIconListShow = true">
                        <i :class="modelDialog.data['bk_obj_icon']"></i>
                    </div>
                    <div class="text">{{$t('选择图标')}}</div>
                </div>
                <div class="content-right">
                    <div class="label-item" v-if="!isMainLine">
                        <span class="label-title">{{$t('所属分组')}}</span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('modelGroup') }">
                            <bk-select style="width: 100%;" ref="groupSelector"
                                v-validate="'required'"
                                name="modelGroup"
                                :value="modelDialog.data.bk_classification_id"
                                :scroll-height="200">
                                <bk-option v-for="(option, index) in localClassifications"
                                    :key="index"
                                    :id="option.bk_classification_id"
                                    :name="option.bk_classification_name">
                                    <cmdb-auth class="group-auth" tag="div" style="display: block;"
                                        :auth="{ type: $OPERATION.C_MODEL, relation: [option.id] }"
                                        @click.native.stop
                                        @click="handleSelectGroup(option)">
                                        {{option.bk_classification_name}}
                                    </cmdb-auth>
                                </bk-option>
                            </bk-select>
                            <p class="form-error" :title="errors.first('modelGroup')">{{errors.first('modelGroup')}}</p>
                        </div>
                    </div>
                    <label>
                        <span class="label-title">{{$t('唯一标识')}}</span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('modelId') }">
                            <bk-input type="text" class="cmdb-form-input"
                                name="modelId"
                                :placeholder="$t('请输入唯一标识')"
                                v-model.trim="modelDialog.data['bk_obj_id']"
                                v-validate="'required|modelId|length:115'">
                            </bk-input>
                            <p class="form-error" :title="errors.first('modelId')">{{errors.first('modelId')}}</p>
                        </div>
                        <i class="icon-cc-exclamation-tips" v-bk-tooltips="$t('模型唯一标识提示语')"></i>
                    </label>
                    <label>
                        <span class="label-title">{{$t('名称')}}</span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('modelName') }">
                            <bk-input type="text" class="cmdb-form-input"
                                name="modelName"
                                :placeholder="$t('请输入名称')"
                                v-validate="'required|singlechar|length:128'"
                                v-model.trim="modelDialog.data['bk_obj_name']">
                            </bk-input>
                            <p class="form-error" :title="errors.first('modelName')">{{errors.first('modelName')}}</p>
                        </div>
                        <i class="icon-cc-exclamation-tips" v-bk-tooltips="$t('请填写模型名')"></i>
                    </label>
                </div>
            </div>
            <div class="model-icon-wrapper" v-if="modelDialog.isIconListShow">
                <the-choose-icon
                    v-model="modelDialog.data['bk_obj_icon']"
                    @chooseIcon="modelDialog.isIconListShow = false"
                    @close="modelDialog.isIconListShow = false"
                ></the-choose-icon>
            </div>
        </div>
        <div slot="footer" class="footer">
            <bk-button theme="primary" @click="confirm">{{$t('提交')}}</bk-button>
            <bk-button theme="default" @click="cancel">{{$t('取消')}}</bk-button>
        </div>
    </bk-dialog>
</template>

<script>
    import theChooseIcon from './choose-icon/_choose-icon'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            theChooseIcon
        },
        props: {
            title: {
                type: String,
                default: ''
            },
            isMainLine: {
                type: Boolean,
                default: false
            },
            isShow: {
                type: Boolean,
                default: false
            },
            groupId: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                modelDialog: {
                    isShow: false,
                    isIconListShow: false,
                    data: {
                        bk_classification_id: '',
                        bk_obj_icon: 'icon-cc-default',
                        bk_obj_id: '',
                        bk_obj_name: ''
                    }
                }
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            localClassifications () {
                const filterGroups = ['bk_biz_topo', 'bk_host_manage', 'bk_organization']
                return this.classifications.filter(group => !filterGroups.includes(group.bk_classification_id))
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.modelDialog.data['bk_classification_id'] = ''
                    this.modelDialog.data['bk_obj_icon'] = 'icon-cc-default'
                    this.modelDialog.data['bk_obj_id'] = ''
                    this.modelDialog.data['bk_obj_name'] = ''
                    this.$validator.reset()
                }
            },
            groupId (value) {
                this.modelDialog.data.bk_classification_id = value
            }
        },
        methods: {
            handleSelectGroup (group) {
                this.modelDialog.data.bk_classification_id = group.bk_classification_id
                this.$refs.groupSelector.close()
            },
            async confirm () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                this.$emit('confirm', this.modelDialog.data)
            },
            cancel () {
                this.$emit('update:isShow', false)
                this.$emit('update:groupId', '')
                this.$validator.reset()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .dialog {
        /deep/ .bk-dialog-tool {
            display: none;
        }
        .dialog-content {
            padding: 15px 15px 20px 28px;
        }
        .title {
            font-size: 20px;
            color: #333948;
            line-height: 1;
        }
        .label-item,
        label {
            display: block;
            margin-bottom: 10px;
            font-size: 0;
            &:last-child {
                margin: 0;
            }
            .color-danger {
                display: inline-block;
                font-size: 16px;
                width: 15px;
                text-align: center;
                vertical-align: middle;
            }
            .icon-cc-exclamation-tips {
                font-size: 18px;
                color: $cmdbBorderColor;
            }
            .label-title {
                font-size: 16px;
                line-height: 36px;
                vertical-align: middle;
                @include ellipsis;
            }
            .cmdb-form-item {
                display: inline-block;
                margin-right: 10px;
                width: 519px;
                vertical-align: middle;
            }
        }
        .footer {
            font-size: 0;
            text-align: right;
            .bk-primary {
                margin-right: 10px;
            }
        }
    }
    .model-dialog {
        text-align: left;
        .dialog-content {
            position: relative;
        }
        .content-left {
            padding-top: 20px;
            text-align: center;
            .icon-wrapper {
                margin: 0 auto;
                width: 85px;
                height: 85px;
                border: 1px solid $cmdbBorderColor;
                border-radius: 50%;
                font-size: 50px;
                cursor: pointer;
            }
            i {
                vertical-align: top;
                line-height: 83px;
                color: $cmdbBorderFocusColor;
            }
            .text {
                margin-top: 8px;
                font-size: 12px;
                line-height: 1;
            }
        }
        .model-icon-wrapper {
            position: absolute;
            left: 0;
            top:0;
            width: 100%;
            background: #fff;
            z-index: 99;
        }
    }
    .group-auth {
        margin: 0 -16px;
        padding: 0 16px;
        &.disabled {
            background-color: #fff;
            color: $textDisabledColor;
        }
    }
</style>
