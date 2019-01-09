<template>
    <div class="verification-detail-wrapper">
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["校验规则"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="verification-selector">
                <div class="text-content" @click="toggleSelector(true)" :class="{'open': attribute.isShow}">
                    <span>{{selectedName}}</span>
                    <i class="bk-icon icon-angle-down"></i>
                </div>
                <bk-selector
                    setting-key="id"
                    display-key="bk_property_name"
                    ref="attrSelector"
                    :multi-select="true"
                    :selected.sync="verificationInfo.selected"
                    :list="attribute.list"
                    @visible-toggle="toggleSelector"
                ></bk-selector>
            </div>
        </div>
        <div class="verification-selector-mask" v-if="attribute.isShow"></div>
        <div class="radio-box">
            <label class="label-text">
                {{$t('ModelManagement["是否为必须校验"]')}}
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="required" :value="true" :disabled="isReadOnly" v-model="verificationInfo['must_check']">
                <span class="cmdb-radio-text">{{$t('ModelManagement["是"]')}}</span>
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="required" :value="false" :disabled="isReadOnly" v-model="verificationInfo['must_check']">
                <span class="cmdb-radio-text">{{$t('ModelManagement["否"]')}}</span>
            </label>
        </div>
        <div class="btn-group">
            <bk-button type="primary"
            :disabled="isReadOnly"
            :loading="$loading(['createObjectUniqueConstraints', 'updateObjectUniqueConstraints'])"
            @click="saveVerification">
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
            },
            attributeList: {
                type: Array
            }
        },
        data () {
            return {
                verificationInfo: {
                    selected: [],
                    must_check: false
                },
                attribute: {
                    isShow: false,
                    list: this.attributeList
                }
            }
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            selectedName () {
                let nameList = []
                this.verificationInfo.selected.forEach(id => {
                    let attr = this.attribute.list.find(attr => attr.id === id)
                    if (attr) {
                        nameList.push(attr['bk_property_name'])
                    }
                })
                return nameList.join(',')
            },
            params () {
                let params = {
                    must_check: this.verificationInfo['must_check'],
                    keys: []
                }
                this.verificationInfo.selected.forEach(id => {
                    params.keys.push({
                        key_kind: 'property',
                        key_id: id
                    })
                })
                return params
            }
        },
        created () {
            if (this.isEdit) {
                this.initData()
            }
        },
        methods: {
            ...mapActions('objectUnique', [
                'createObjectUniqueConstraints',
                'updateObjectUniqueConstraints'
            ]),
            initData () {
                this.verificationInfo['must_check'] = this.verification['must_check']
                this.verification.keys.forEach(key => {
                    this.verificationInfo.selected.push(key['key_id'])
                })
            },
            toggleSelector (isShow) {
                if (!this.isReadOnly) {
                    this.$refs.attrSelector.open = isShow
                    this.attribute.isShow = isShow
                }
            },
            async saveVerification () {
                if (this.isEdit) {
                    await this.updateObjectUniqueConstraints({
                        id: this.verification.id,
                        objId: this.activeModel['bk_obj_id'],
                        params: this.params,
                        config: {
                            requestId: 'updateObjectUniqueConstraints'
                        }
                    })
                    this.$emit('save')
                } else {
                    await this.createObjectUniqueConstraints({
                        objId: this.activeModel['bk_obj_id'],
                        params: this.params,
                        config: {
                            requestId: 'createObjectUniqueConstraints'
                        }
                    })
                    this.$emit('save')
                }
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
                width: 100%;
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
                        padding: 5px 28px 5px 10px;
                        height: auto;
                        min-height: 36px;
                        line-height: 24px;
                        overflow: visible;
                        border-color: $cmdbBorderFocusColor;
                        .icon-angle-down {
                            color: $cmdbBorderFocusColor;
                            transform: rotate(180deg);
                        }
                    }
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
        .verification-selector-mask {
            position: absolute;
            left: 0;
            top: 0;
            right: 0;
            bottom: 0;
            width: 100%;
            height: 100%;
        }
        .radio-box {
            .label-text {
                width: 150px;
            }
            .cmdb-form-radio {
                width: 114px;
                vertical-align: middle;
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
