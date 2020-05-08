<template>
    <div class="relation-type-content">
        <div class="content-box">
            <label class="form-label">
                <span class="label-text">
                    {{$t('唯一标识')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstId') }">
                    <bk-input type="text" class="cmdb-form-input"
                        name="asstId"
                        :disabled="isReadOnly || isEdit"
                        v-model.trim="relationInfo['bk_asst_id']"
                        v-validate="'required|associationId|length:128'"
                        :placeholder="$t('请输入英文标识')">
                    </bk-input>
                    <p class="form-error" :title="errors.first('asstId')">{{errors.first('asstId')}}</p>
                </div>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('名称')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstName') }">
                    <bk-input type="text"
                        class="cmdb-form-input"
                        name="asstName"
                        :disabled="isReadOnly || isEdit && relation.ispre"
                        v-validate="'required|singlechar|length:20'"
                        v-model.trim="relationInfo['bk_asst_name']"
                        :placeholder="$t('请输入名称')">
                    </bk-input>
                    <p class="form-error">{{errors.first('asstName')}}</p>
                </div>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('源->目标描述')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{ 'is-error': errors.has('srcDes') }">
                    <input type="text"
                        class="cmdb-form-input"
                        name="srcDes"
                        :disabled="isReadOnly || isEdit && relation.ispre"
                        v-validate="'required|singlechar|length:256'"
                        v-model.trim="relationInfo['src_des']"
                        :placeholder="$t('请输入关联描述如：连接、运行')">
                    <p class="form-error">{{errors.first('srcDes')}}</p>
                </div>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('目标->源描述')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{ 'is-error': errors.has('destDes') }">
                    <bk-input type="text"
                        class="cmdb-form-input"
                        name="destDes"
                        :disabled="isReadOnly || isEdit && relation.ispre"
                        v-validate="'required|singlechar|length:256'"
                        v-model.trim="relationInfo['dest_des']"
                        :placeholder="$t('请输入关联描述如：属于、上联')">
                    </bk-input>
                    <p class="form-error">{{errors.first('destDes')}}</p>
                </div>
            </label>
            <div class="radio-box overflow">
                <div class="radio-box overflow">
                    <label class="label-text">
                        {{$t('是否有方向')}}
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio" name="direction" value="src_to_dest" v-model="relationInfo.direction"
                            :disabled="isReadOnly || isEdit && relation.ispre">
                        <span class="cmdb-radio-text">{{$t('有，源指向目标')}}</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio" name="direction" value="none" v-model="relationInfo.direction"
                            :disabled="isReadOnly || isEdit && relation.ispre">
                        <span class="cmdb-radio-text">{{$t('无方向')}}</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio" name="direction" value="bidirectional" v-model="relationInfo.direction"
                            :disabled="isReadOnly || isEdit && relation.ispre">
                        <span class="cmdb-radio-text">{{$t('双向')}}</span>
                    </label>
                </div>
            </div>
        </div>
        <slot name="operation">
            <div class="btn-group">
                <bk-button theme="primary" :disabled="isReadOnly || isEdit && relation.ispre" :loading="$loading(['updateAssociationType', 'createAssociationType'])" @click="saveRelation">
                    {{saveBtnText || $t('确定')}}
                </bk-button>
                <bk-button theme="default" @click="cancel">
                    {{$t('取消')}}
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
                },
                originInfo: {}
            }
        },
        computed: {
            createParams () {
                const createParams = { ...this.relationInfo }
                delete createParams.id
                return createParams
            },
            updateParams () {
                const updateParams = { ...this.relationInfo }
                delete updateParams['bk_asst_id']
                return updateParams
            },
            changedValues () {
                const values = {}
                Object.keys(this.relationInfo).forEach(key => {
                    if (this.relationInfo[key] !== this.originInfo[key]) {
                        values[key] = this.relationInfo[key]
                    }
                })
                return values
            }
        },
        created () {
            if (this.isEdit) {
                for (const key in this.relationInfo) {
                    this.relationInfo[key] = this.$tools.clone(this.relation[key])
                }
            }
            this.originInfo = this.$tools.clone(this.relationInfo)
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
            }
            .label-text {
                width: auto;
            }
        }
    }
</style>
