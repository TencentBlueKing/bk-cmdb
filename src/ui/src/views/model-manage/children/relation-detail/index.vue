<template>
    <div class="model-slider-content">
        <label class="form-label" v-if="isEdit">
            <span class="label-text">
                {{$t('唯一标识')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <bk-input type="text" class="cmdb-form-input" v-model.trim="objAsstId" disabled></bk-input>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('关联描述')}}
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstName') }">
                <bk-input type="text" class="cmdb-form-input"
                    name="asstName"
                    :placeholder="$t('请输入关联描述')"
                    :disabled="relationInfo.ispre || isReadOnly"
                    v-model.trim="relationInfo['bk_obj_asst_name']"
                    v-validate="'singlechar|length:256'">
                </bk-input>
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('源模型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('objId') }">
                <cmdb-selector
                    class="bk-select-full-width"
                    :disabled="relationInfo.ispre || isEdit"
                    :has-children="true"
                    :auto-select="false"
                    :list="asstList"
                    v-validate="'required'"
                    name="objId"
                    v-model="relationInfo['bk_obj_id']"
                ></cmdb-selector>
                <p class="form-error">{{errors.first('objId')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="form-label">
            <span class="label-text">
                {{$t('目标模型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstObjId') }">
                <cmdb-selector
                    class="bk-select-full-width"
                    :disabled="relationInfo.ispre || isEdit"
                    :has-children="true"
                    :auto-select="false"
                    :list="asstList"
                    v-validate="'required'"
                    name="asstObjId"
                    v-model="relationInfo['bk_asst_obj_id']"
                ></cmdb-selector>
                <p class="form-error">{{errors.first('asstObjId')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="form-label">
            <span class="label-text">
                {{$t('关联类型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstId') }">
                <cmdb-selector
                    class="bk-select-full-width"
                    :disabled="relationInfo.ispre || isReadOnly || isEdit"
                    :list="usefulRelationList"
                    v-validate="'required'"
                    name="asstId"
                    v-model="relationInfo['bk_asst_id']"
                ></cmdb-selector>
                <p class="form-error">{{errors.first('asstId')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="form-label">
            <span class="label-text">
                {{$t('源-目标约束')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('mapping') }">
                <cmdb-selector
                    class="bk-select-full-width"
                    :disabled="relationInfo.ispre || isEdit"
                    :list="mappingList"
                    v-validate="'required'"
                    name="mapping"
                    v-model="relationInfo.mapping"
                ></cmdb-selector>
                <p class="form-error">{{errors.first('mapping')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="btn-group">
            <bk-button theme="primary" :disabled="isReadOnly" :loading="$loading(['createObjectAssociation', 'updateObjectAssociation'])" @click="saveRelation">
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
            relationList: {
                type: Array
            }
        },
        data () {
            return {
                mappingList: [{
                    id: 'n:n',
                    name: 'N-N'
                }, {
                    id: '1:n',
                    name: '1-N'
                }, {
                    id: '1:1',
                    name: '1-1'
                }],
                relationInfo: {
                    ispre: false,
                    id: 0,
                    bk_obj_asst_id: '',
                    bk_obj_asst_name: '',
                    bk_obj_id: '',
                    bk_asst_obj_id: '',
                    bk_asst_id: '',
                    mapping: ''
                },
                specialModel: ['process', 'plat'],
                originRelationInfo: {}
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            ...mapGetters('objectModel', [
                'activeModel',
                'isInjectable'
            ]),
            usefulRelationList () {
                return this.relationList.filter(relation => relation.id !== 'bk_mainline')
            },
            objAsstId () {
                if (this.isEdit) {
                    return this.relation.bk_obj_asst_id
                }
                const relationInfo = this.relationInfo
                if (relationInfo['bk_obj_id'].length && relationInfo['bk_asst_id'].length && relationInfo['bk_asst_obj_id'].length) {
                    return `${relationInfo['bk_obj_id']}_${relationInfo['bk_asst_id']}_${relationInfo['bk_asst_obj_id']}`
                }
                return ''
            },
            createParams () {
                return {
                    bk_obj_asst_id: this.objAsstId,
                    bk_obj_asst_name: this.relationInfo['bk_obj_asst_name'],
                    bk_obj_id: this.relationInfo['bk_obj_id'],
                    bk_asst_obj_id: this.relationInfo['bk_asst_obj_id'],
                    bk_asst_id: this.relationInfo['bk_asst_id'],
                    mapping: this.relationInfo.mapping
                }
            },
            updateParams () {
                return {
                    bk_obj_asst_name: this.relationInfo['bk_obj_asst_name'],
                    bk_asst_id: this.relationInfo['bk_asst_id']
                }
            },
            asstList () {
                const asstList = []
                this.classifications.forEach(classify => {
                    if (classify['bk_objects'].length) {
                        const objects = []
                        classify['bk_objects'].forEach(({ bk_obj_id: objId, bk_obj_name: objName }) => {
                            if (!this.specialModel.includes(objId)) {
                                objects.push({
                                    id: objId,
                                    name: objName
                                })
                            }
                        })
                        if (objects.length) {
                            asstList.push({
                                name: classify['bk_classification_name'],
                                children: objects
                            })
                        }
                    }
                })
                return asstList
            },
            changedValues () {
                const changedValues = {}
                for (const propertyId in this.relationInfo) {
                    if (JSON.stringify(this.relationInfo[propertyId]) !== JSON.stringify(this.originRelationInfo[propertyId])) {
                        changedValues[propertyId] = this.relationInfo[propertyId]
                    }
                }
                return changedValues
            }
        },
        watch: {
            'relationInfo.bk_obj_id' (val) {
                if (val !== this.activeModel['bk_obj_id']) {
                    this.relationInfo['bk_asst_obj_id'] = this.activeModel['bk_obj_id']
                }
            },
            'relationInfo.bk_asst_obj_id' (val) {
                if (val !== this.activeModel['bk_obj_id']) {
                    this.relationInfo['bk_obj_id'] = this.activeModel['bk_obj_id']
                }
            }
        },
        created () {
            this.initData()
            this.$nextTick(() => {
                this.originRelationInfo = this.$tools.clone(this.relationInfo)
            })
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createObjectAssociation',
                'updateObjectAssociation'
            ]),
            async validateValue () {
                await this.$nextTick()
                this.$validator.validateAll()
            },
            initData () {
                if (this.isEdit) {
                    for (const key in this.relationInfo) {
                        this.relationInfo[key] = this.$tools.clone(this.relation[key])
                    }
                } else {
                    this.relationInfo['bk_obj_id'] = this.activeModel['bk_obj_id']
                }
            },
            async saveRelation () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                if (this.isEdit) {
                    await this.updateObjectAssociation({
                        id: this.relationInfo.id,
                        params: this.$injectMetadata(this.updateParams, {
                            clone: true,
                            inject: this.isInjectable
                        }),
                        config: {
                            requestId: 'updateObjectAssociation'
                        }
                    })
                } else {
                    await this.createObjectAssociation({
                        params: this.$injectMetadata(this.createParams, {
                            clone: true,
                            inject: this.isInjectable
                        }),
                        config: {
                            requestId: 'createObjectAssociation'
                        }
                    })
                }
                this.$emit('save')
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-relation-wrapper {
        padding: 20px;
    }
</style>
