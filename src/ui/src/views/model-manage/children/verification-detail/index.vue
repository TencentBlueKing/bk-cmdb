<template>
    <div class="model-slider-content verification-detail-wrapper">
        <div class="form-label">
            <span class="label-text">
                {{$t('校验规则')}}
                <span class="color-danger">*</span>
            </span>
            <div class="verification-selector">
                <bk-select
                    ref="attrSelector"
                    multiple
                    :disabled="isReadOnly"
                    :clearable="false"
                    v-model="verificationInfo.selected"
                    @toggle="toggleSelector">
                    <div class="text-content" slot="trigger" @click="toggleSelector(true)" :class="{ 'open': attribute.isShow }">
                        <span>{{selectedName}}</span>
                    </div>
                    <bk-option v-for="(option, index) in attribute.list"
                        :key="index"
                        :id="option.id"
                        :name="option.bk_property_name">
                    </bk-option>
                </bk-select>
            </div>
        </div>
        <div class="btn-group">
            <bk-button theme="primary"
                :disabled="isReadOnly || !verificationInfo.selected.length"
                :loading="$loading(['createObjectUniqueConstraints', 'updateObjectUniqueConstraints'])"
                @click="saveVerification">
                {{isEdit ? $t('保存') : $t('提交')}}
            </bk-button>
            <bk-button theme="default" @click="cancel">
                {{$t('取消')}}
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
                },
                originVerificationInfo: {}
            }
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            selectedName () {
                const nameList = []
                this.verificationInfo.selected.forEach(id => {
                    const attr = this.attribute.list.find(attr => attr.id === id)
                    if (attr) {
                        nameList.push(attr['bk_property_name'])
                    }
                })
                return nameList.join(',')
            },
            params () {
                const params = {
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
            },
            changedValues () {
                const changedValues = {}
                for (const propertyId in this.verificationInfo) {
                    if (JSON.stringify(this.verificationInfo[propertyId]) !== JSON.stringify(this.originVerificationInfo[propertyId])) {
                        changedValues[propertyId] = this.verificationInfo[propertyId]
                    }
                }
                return changedValues
            }
        },
        created () {
            if (this.isEdit) {
                this.initData()
            }
            this.originVerificationInfo = this.$tools.clone(this.verificationInfo)
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
                    // isShow ? this.$refs.attrSelector.show() : this.$refs.attrSelector.hide()
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
    .model-slider-content.verification-detail-wrapper {
        .form-label {
            .label-text {
                width: 130px;
            }
            .verification-selector {
                position: relative;
                display: inline-block;
                width: 100%;
                .bk-select {
                    width: 100%;
                    border: none;
                }
                .text-content {
                    border-radius: 2px;
                    border: 1px solid $cmdbBorderColor;
                    padding: 0 32px 0 10px;
                    height: 32px;
                    line-height: 30px;
                    overflow: hidden;
                    &.open {
                        padding: 3px 28px 3px 10px;
                        height: auto;
                        min-height: 32px;
                        line-height: 24px;
                        overflow: visible;
                        border-color: $cmdbBorderFocusColor;
                    }
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
