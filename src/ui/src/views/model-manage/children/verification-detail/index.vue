<template>
    <div class="verification-detail-wrapper">
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["校验规则"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="verification-selector">
                <div class="text-content" @click="toggleSelector(true)" :class="{'open': attribute.isShow}">
                    <span>{{selectedName}}</span>
                    <i class="bk-icon icon-angle-down"></i>
                </div>
                <div class="verification-selector-mask" v-if="attribute.isShow"></div>
                <bk-selector
                    setting-key="bk_property_id"
                    display-key="bk_property_name"
                    ref="attrSelector"
                    :multi-select="true"
                    :selected.sync="verificationInfo.selected"
                    :list="attribute.list"
                    @visible-toggle="toggleSelector"
                ></bk-selector>
            </div>
        </label>
        <div class="radio-box">
            <label class="label-text">
                {{$t('ModelManagement["是否显示为实例名称"]')}}
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="isName" value="true" v-model="verificationInfo.isName">
                <span class="cmdb-radio-text">{{$t('ModelManagement["是"]')}}</span>
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="isName" value="false" v-model="verificationInfo.isName">
                <span class="cmdb-radio-text">{{$t('ModelManagement["否"]')}}</span>
            </label>
        </div>
        <div class="radio-box">
            <label class="label-text">
                {{$t('ModelManagement["是否为必须校验"]')}}
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="required" value="true" v-model="verificationInfo.required">
                <span class="cmdb-radio-text">{{$t('ModelManagement["是"]')}}</span>
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="required" value="false" v-model="verificationInfo.required">
                <span class="cmdb-radio-text">{{$t('ModelManagement["否"]')}}</span>
            </label>
        </div>
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading(['updateAssociationType', 'createAssociationType'])" @click="saveVerification">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('Common["取消"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            verification: {
                type: Object
            },
            isReadOnly: {
                type: Boolean,
                default: false
            },
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                verificationInfo: {
                    selected: [],
                    isName: false,
                    required: false
                },
                attribute: {
                    isShow: false,
                    list: []
                }
            }
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            selectedName () {
                let nameList = []
                this.verificationInfo.selected.map(propertyId => {
                    let attr = this.attribute.list.find(({bk_property_id: bkPropertyId}) => {
                        return bkPropertyId === propertyId
                    })
                    if (attr) {
                        nameList.push(attr['bk_property_name'])
                    }
                })
                return nameList.join(',')
            }
        },
        created () {
            this.initAttrList()
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            toggleSelector (isShow) {
                this.attribute.isShow = isShow
                this.$refs.attrSelector.open = isShow
            },
            async initAttrList () {
                const res = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: this.activeModel['bk_obj_id']
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`
                    }
                })
                this.attribute.list = res
            },
            saveVerification () {

            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .slider-content.verification-detail-wrapper {
        .form-label {
            .label-text {
                width: 130px;
            }
            .verification-selector {
                position: relative;
                display: inline-block;
                width: calc(100% - 145px);
                .bk-selector {
                    width: 100%;
                }
                .text-content {
                    border-radius: 2px;
                    border: 1px solid $cmdbBorderColor;
                    padding: 0 32px 0 10px;
                    height: 36px;
                    line-height: 34px;
                    font-size: 14px;
                    overflow: hidden;
                    &.open {
                        padding: 5px 28px 5px 16px;
                        height: auto;
                        line-height: 20px;
                        overflow: visible;
                        border-color: $cmdbBorderFocusColor;
                        .icon-angle-down {
                            color: $cmdbBorderFocusColor;
                            transform: rotate(180deg);
                        }
                    }
                }
                .verification-selector-mask {
                    position: absolute;
                    left: 0;
                    top: 0;
                    right: 0;
                    bottom: 0;
                    width: 100%;
                    height: 100%;
                }
                .icon-angle-down {
                    position: absolute;
                    right: 10px;
                    top: 12px;
                    font-size: 12px;
                    transition: transform .2s linear;
                }
            }
        }
        .radio-box {
            .label-text {
                width: 150px;
            }
            .cmdb-form-radio {
                width: 114px;
            }
        }
    }
</style>

<style lang="scss">
    .verification-detail-wrapper {
        .verification-selector {
            .bk-selector {
                position: absolute;
                left: 0;
                bottom: 32px;
            }
            .bk-selector-wrapper {
                display: none;
            }
            .bk-selector-list {
                top: 36px;
                left: 1px;
            }
        }
    }
</style>
