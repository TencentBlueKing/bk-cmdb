<template>
    <div class="relation-type-content">
        <div class="content-box">
            <label class="form-label">
                <span class="label-text">
                    {{$t('ModelManagement["唯一标识"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                    <input type="text" class="cmdb-form-input"
                    name="asstId"
                    :disabled="isEdit"
                    v-model.trim="relationInfo['bk_asst_id']"
                    v-validate="'required|associationId'"
                    :placeholder="$t('ModelManagement[\'请输入英文标识\']')">
                    <p class="form-error">{{errors.first('asstId')}}</p>
                </div>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('Hosts["名称"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{'is-error': errors.has('asstName')}">
                    <input type="text"
                    class="cmdb-form-input"
                    name="asstName"
                    :disabled="isEdit && relation.ispre"
                    v-validate="'required|singlechar'"
                    v-model.trim="relationInfo['bk_asst_name']" 
                    :placeholder="$t('ModelManagement[\'请输入名称\']')">
                    <p class="form-error">{{errors.first('asstName')}}</p>
                </div>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('ModelManagement["源->目标描述"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{'is-error': errors.has('srcDes')}">
                    <input type="text"
                    class="cmdb-form-input"
                    name="srcDes"
                    :disabled="isEdit && relation.ispre"
                    v-validate="'required|singlechar'"
                    v-model.trim="relationInfo['src_des']" 
                    :placeholder="$t('ModelManagement[\'请输入关联描述如：连接、运行\']')">
                    <p class="form-error">{{errors.first('srcDes')}}</p>
                </div>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('ModelManagement["目标->源描述"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{'is-error': errors.has('destDes')}">
                    <input type="text"
                    class="cmdb-form-input"
                    name="destDes"
                    :disabled="isEdit && relation.ispre"
                    v-validate="'required|singlechar'"
                    v-model.trim="relationInfo['dest_des']" 
                    :placeholder="$t('ModelManagement[\'请输入关联描述如：属于、上联\']')">
                    <p class="form-error">{{errors.first('destDes')}}</p>
                </div>
            </label>
            <div class="radio-box overflow">
                <div class="radio-box overflow">
                    <label class="label-text">
                        {{$t('ModelManagement["是否有方向"]')}}
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio" name="direction" value="src_to_dest" v-model="relationInfo.direction"
                        :disabled="isEdit && relation.ispre">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["有，源指向目标"]')}}</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio" name="direction" value="none" v-model="relationInfo.direction"
                        :disabled="isEdit && relation.ispre">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["无方向"]')}}</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio" name="direction" value="bidirectional" v-model="relationInfo.direction"
                        :disabled="isEdit && relation.ispre">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["双向"]')}}</span>
                    </label>
                </div>
            </div>
        </div>
        <slot name="operation">
            <div class="btn-group">
                <bk-button type="primary" :disabled="isEdit && relation.ispre" :loading="$loading(['updateAssociationType', 'createAssociationType'])" @click="saveRelation">
                    {{saveBtnText || $t('Common["确定"]')}}
                </bk-button>
                <bk-button type="default" @click="cancel">
                    {{$t('Common["取消"]')}}
                </bk-button>
            </div>
        </slot>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            relation: {
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
            saveBtnText: {
                type: String
            }
        },
        data () {
            return {
                relationInfo: {
                    id: null,
                    bk_asst_id: '',
                    bk_asst_name: '',
                    src_des: '',
                    dest_des: '',
                    direction: 'src_to_dest' // none, src_to_dest, bidirectional
                }
            }
        },
        computed: {
            createParams () {
                const createParams = {...this.relationInfo}
                delete createParams.id
                return createParams
            },
            updateParams () {
                const updateParams = {...this.relationInfo}
                delete updateParams['bk_asst_id']
                return updateParams
            }
        },
        created () {
            if (this.isEdit) {
                for (let key in this.relationInfo) {
                    this.relationInfo[key] = this.$tools.clone(this.relation[key])
                }
            }
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createAssociationType',
                'updateAssociationType'
            ]),
            async saveRelation () {
                try {
                    if (!await this.$validator.validateAll()) {
                        return
                    }
                    if (this.isEdit) {
                        await this.updateAssociationType({
                            id: this.relation.id,
                            params: this.updateParams,
                            config: {
                                requestId: 'updateAssociationType'
                            }
                        })
                    } else {
                        await this.createAssociationType({
                            params: this.createParams,
                            config: {
                                requestId: 'createAssociationType'
                            }
                        })
                    }
                    this.$emit('saved')
                } catch (e) {
                    console.error(e)
                }
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>


<style lang="scss" scoped>
    .relation-type-content {
        .content-box {
            .radio-box {
                .cmdb-form-radio {
                    width: auto;
                    padding-right: 20px;
                    vertical-align: middle;
                    &:last-child {
                        padding-right: 0;
                    }
                }
                .label-text {
                    padding-right: 2px;
                }
            }
        }
    }
</style>
