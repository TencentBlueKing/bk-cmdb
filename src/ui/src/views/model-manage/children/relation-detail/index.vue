<template>
    <div class="slider-content">
        <label class="form-label" v-if="isEdit">
            <span class="label-text">
                {{$t('ModelManagement["唯一标识"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input" v-model.trim="objAsstId" disabled>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关联描述"]')}}
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstName')}">
                <input type="text" class="cmdb-form-input"
                name="asstName"
                :placeholder="$t('ModelManagement[\'请输入关联描述\']')"
                :disabled="relationInfo.ispre || isReadOnly"
                v-model.trim="relationInfo['bk_obj_asst_name']"
                v-validate="'singlechar'">
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('objId')}">
                <cmdb-selector
                    :disabled="relationInfo.ispre || isEdit"
                    :has-children="true"
                    :autoSelect="false"
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
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstObjId')}">
                <cmdb-selector
                    :disabled="relationInfo.ispre || isEdit"
                    :has-children="true"
                    :autoSelect="false"
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
                {{$t('ModelManagement["关联类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <cmdb-selector
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
                {{$t('ModelManagement["源-目标约束"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('mapping')}">
                <cmdb-selector
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
            <bk-button type="primary" :loading="$loading(['createObjectAssociation', 'updateObjectAssociation'])" @click="saveRelation">
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
                    id: '1:1',
                    name: '1-1'
                }, {
                    id: '1:n',
                    name: '1-N'
                }, {
                    id: 'n:n',
                    name: 'N-N'
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
                specialModel: ['process', 'plat']
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
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
                let asstList = []
                this.classifications.forEach(classify => {
                    if (classify['bk_objects'].length) {
                        let objects = []
                        classify['bk_objects'].forEach(({bk_obj_id: objId, bk_obj_name: objName}) => {
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
                    for (let key in this.relationInfo) {
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
                        params: this.updateParams,
                        config: {
                            requestId: 'updateObjectAssociation'
                        }
                    })
                } else {
                    await this.createObjectAssociation({
                        params: this.createParams,
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
