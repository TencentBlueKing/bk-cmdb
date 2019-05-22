<template>
    <div class="create-template-wrapper">
        <div class="info-group">
            <h3>基本属性</h3>
            <div class="form-info clearfix">
                <label class="label-text fl" for="templateName">
                    {{$t('ServiceManagement["模板名称"]')}}
                    <span class="color-danger">*</span>
                </label>
                <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('fieldId') }">
                    <input type="text" class="cmdb-form-input" id="templateName"
                        name="fieldId"
                        :placeholder="$t('ServiceManagement[\'请输入模版名称\']')"
                        v-model.trim="formData.tempalteName"
                        v-validate="'required|fieldId'">
                    <p class="form-error">{{errors.first('fieldId')}}</p>
                </div>
            </div>
            <div class="form-info clearfix">
                <span class="label-text fl">
                    {{$t('ServiceManagement["服务分类"]')}}
                    <span class="color-danger">*</span>
                </span>
                <bk-selector
                    class="fl"
                    placeholder="请选择一级分类"
                    :list="[]"
                    :selected.sync="formData.primaryClassification">
                </bk-selector>
                <bk-selector
                    class="fl"
                    placeholder="请选择二级分类"
                    :list="[]"
                    :selected.sync="formData.secondaryClassification">
                </bk-selector>
            </div>
        </div>
        <div class="info-group">
            <h3>进程服务</h3>
            <div class="precess-box">
                <div class="process-create">
                    <bk-button class="create-btn">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t("ServiceManagement['新建进程']")}}</span>
                    </bk-button>
                    <span class="create-tips">{{$t("ServiceManagement['新建进程提示']")}}</span>
                </div>
                <div class="process-table">
                </div>
                <div class="btn-box">
                    <bk-button type="primary">{{$t("Common['确定']")}}</bk-button>
                    <bk-button>{{$t("Common['取消']")}}</bk-button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                formData: {
                    primaryClassification: '',
                    secondaryClassification: '',
                    tempalteName: ''
                }
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t("ServiceManagement['新建服务模版']"))
        }
    }
</script>

<style lang="scss" scoped>
    .create-template-wrapper {
        .info-group {
            h3 {
                color: #63656e;
                font-size: 14px;
                padding-bottom: 26px;
            }
            .form-info {
                font-size: 14px;
                padding-left: 30px;
                padding-bottom: 22px;
                .label-text {
                    line-height: 36px;
                    padding-right: 6px;
                }
                .cmdb-form-item {
                    width: 520px;
                }
                .bk-selector {
                    width: 254px;
                    margin-right: 10px;
                }
            }
            .precess-box {
                padding-left: 30px;
            }
            .process-create {
                .create-btn {
                    padding: 0 16px;
                    span {
                        margin-left: 0;
                    }
                }
                .icon-plus {
                    font-size: 12px;
                    line-height: normal;
                    font-weight: bold;
                }
                .create-tips {
                    color: #979Ba5;
                    font-size: 12px;
                    padding-left: 10px;
                }
            }
            .btn-box {
                padding-top: 30px;
            }
        }
    }
</style>
