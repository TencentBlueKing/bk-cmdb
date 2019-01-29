<template>
    <bk-dialog
        class="model-dialog dialog"
        :close-icon="false"
        :hasHeader="false"
        :width="600"
        :padding="0"
        :quick-close="false"
        :is-show.sync="isShow">
        <div slot="content" class="dialog-content">
            <p class="title">{{title}}</p>
            <div class="content clearfix">
                <div class="content-left">
                    <div class="icon-wrapper" @click="modelDialog.isIconListShow = true">
                        <i :class="modelDialog.data['bk_obj_icon']"></i>
                    </div>
                    <div class="text">{{$t('ModelManagement["选择图标"]')}}</div>
                </div>
                <div class="content-right">
                    <div class="label-item" v-if="!isMainLine">
                        <span class="label-title">{{$t('ModelManagement["所属分组"]')}}</span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('modelGroup')}">
                            <cmdb-selector
                                class="selector-box"
                                name="modelGroup"
                                setting-key="bk_classification_id"
                                display-key="bk_classification_name"
                                :content-max-height="200"
                                :selected.sync="modelDialog.data['bk_classification_id']"
                                :list="localClassifications"
                                :empty-text="isAdminView ? '' : $t('ModelManagement[\'无自定义分组\']')"
                                v-validate="'required'"
                                v-model="modelDialog.data['bk_classification_id']"
                            ></cmdb-selector>
                            <p class="form-error">{{errors.first('modelGroup')}}</p>
                        </div>
                    </div>
                    <label>
                        <span class="label-title">{{$t('ModelManagement["唯一标识"]')}}</span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('modelId')}">
                            <input type="text" class="cmdb-form-input"
                            name="modelId"
                            :placeholder="$t('ModelManagement[\'请输入唯一标识\']')"
                            v-model.trim="modelDialog.data['bk_obj_id']"
                            v-validate="'required|modelId'">
                            <p class="form-error">{{errors.first('modelId')}}</p>
                        </div>
                        <i class="bk-icon icon-info-circle" v-tooltip="$t('ModelManagement[\'下划线，数字，英文小写的组合\']')"></i>
                    </label>
                    <label>
                        <span class="label-title">{{$t('ModelManagement["名称"]')}}</span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('modelName')}">
                            <input type="text" class="cmdb-form-input"
                            name="modelName"
                            :placeholder="$t('ModelManagement[\'请输入名称\']')"
                            v-validate="'required|singlechar'"
                            v-model.trim="modelDialog.data['bk_obj_name']">
                            <p class="form-error">{{errors.first('modelName')}}</p>
                        </div>
                        <i class="bk-icon icon-info-circle" v-tooltip="$t('ModelManagement[\'请填写模型名\']')"></i>
                    </label>
                </div>
            </div>
            <div class="model-icon-wrapper" v-if="modelDialog.isIconListShow">
                <span class="back" @click="modelDialog.isIconListShow = false">
                    <i class="bk-icon icon-back2"></i>
                </span>
                <the-choose-icon
                    v-model="modelDialog.data['bk_obj_icon']"
                    @chooseIcon="modelDialog.isIconListShow = false"
                ></the-choose-icon>
            </div>
        </div>
        <div slot="footer" class="footer">
            <bk-button type="primary" @click="confirm">{{$t("Common['保存']")}}</bk-button>
            <bk-button type="default" @click="cancel">{{$t("Common['取消']")}}</bk-button>
        </div>
    </bk-dialog>
</template>

<script>
    import theChooseIcon from './_choose-icon'
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
            ...mapGetters(['isAdminView']),
            localClassifications () {
                let localClassifications = []
                this.classifications.forEach(classification => {
                    if (!['bk_biz_topo', 'bk_host_manage', 'bk_organization'].includes(classification['bk_classification_id'])) {
                        const localClassification = {
                            ...classification,
                            isModelShow: false
                        }
                        if (this.isAdminView) {
                            localClassifications.push(localClassification)
                        } else if (this.$tools.getMetadataBiz(classification)) {
                            localClassifications.push(localClassification)
                        }
                    }
                })
                return localClassifications
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
            }
        },
        methods: {
            async confirm () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                this.$emit('confirm', this.modelDialog.data)
            },
            cancel () {
                this.$emit('update:isShow', false)
                this.$validator.reset()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .dialog {
        .dialog-content {
            padding: 20px 15px 20px 28px;
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
            .icon-info-circle {
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
            padding: 0 24px;
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
    }
</style>

