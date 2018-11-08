<template>
    <div>
        <label class="form-label">
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
                {{$t('ModelManagement["别名"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstName')}">
                <input type="text" class="cmdb-form-input"
                name="asstName"
                :disabled="relationInfo.ispre"
                v-model.trim="relationInfo['bk_obj_asst_name']"
                v-validate="'required|singlechar'">
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["模型源"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('objId')}">
                <form-selector
                    :disabled="relationInfo.ispre"
                    :has-children="true"
                    :list="asstList"
                    v-validate="'required'"
                    name="objId"
                    v-model="relationInfo['bk_obj_id']"
                ></form-selector>
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
                <form-selector
                    :disabled="relationInfo.ispre"
                    :has-children="true"
                    :list="asstList"
                    v-validate="'required'"
                    name="asstObjId"
                    v-model="relationInfo['bk_asst_obj_id']"
                ></form-selector>
                <p class="form-error">{{errors.first('asstObjId')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <form-selector
                    :disabled="relationInfo.ispre"
                    :list="relationList"
                    v-validate="'required'"
                    name="asstId"
                    v-model="relationInfo['bk_asst_id']"
                ></form-selector>
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
                <form-selector
                    :disabled="relationInfo.ispre"
                    :list="mappingList"
                    v-validate="'required'"
                    name="mapping"
                    v-model="relationInfo.mapping"
                ></form-selector>
                <p class="form-error">{{errors.first('mapping')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="radio-box">
            <label class="label-text">
                {{$t('ModelManagement["联动删除"]')}}
            </label>
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" id="delete_dest" value="delete_dest" v-model="relationInfo['on_delete']"
                :disabled="relationInfo.ispre">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["源不存在联动删除目标"]')}}</span>
            </label>
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" id="delete_src" value="delete_src" v-model="relationInfo['on_delete']"
                :disabled="relationInfo.ispre">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["目标不存在联动删除源"]')}}</span>
            </label>
        </div>
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading(['createObjectAssociation', 'updateObjectAssociation'])" @click="saveRelation">
                {{$t('ModelManagement["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('ModelManagement["取消"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import formSelector from './form-selector'
    export default {
        components: {
            formSelector
        },
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
                    mapping: '',
                    on_delete: []
                },
                relationList: []
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            objAsstId () {
                let {
                    relationInfo
                } = this
                if (relationInfo['bk_obj_id'].length && relationInfo['bk_asst_id'].length && relationInfo['bk_asst_obj_id'].length) {
                    return `${relationInfo['bk_obj_id']}_${relationInfo['bk_asst_id']}_${relationInfo['bk_asst_obj_id']}`
                }
                return ''
            },
            params () {
                return {
                    bk_obj_asst_id: this.objAsstId,
                    bk_obj_asst_name: this.relationInfo['bk_obj_asst_name'],
                    bk_obj_id: this.relationInfo['bk_obj_id'],
                    bk_asst_obj_id: this.relationInfo['bk_asst_obj_id'],
                    bk_asst_id: this.relationInfo['bk_asst_id'],
                    mapping: this.relationInfo.mapping,
                    on_delete: this.relationInfo['on_delete'].length ? this.relationInfo['on_delete'].join('') : 'none'
                }
            },
            asstList () {
                let asstList = []
                this.classifications.forEach(classify => {
                    if (classify['bk_objects'].length) {
                        if (!this.relationInfo.ispre && classify['bk_classification_id'] === 'bk_biz_topo') {
                            return
                        }
                        let objects = []
                        classify['bk_objects'].map(({bk_obj_id: objId, bk_obj_name: objName}) => {
                            if (this.relationInfo.ispre || ['plat', 'process'].indexOf(objId) === -1) {
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
            },
            'relationInfo.on_delete' (newVal, oldVal) {
                if (newVal.length === 2) {
                    this.relationInfo['on_delete'] = newVal.filter(val => !oldVal.includes(val))
                }
            }
        },
        created () {
            this.initRelationList()
            this.initData()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createObjectAssociation',
                'updateObjectAssociation',
                'searchAssociationType'
            ]),
            async validateValue () {
                await this.$nextTick()
                this.$validator.validateAll()
            },
            initData () {
                if (this.isEdit) {
                    for (let key in this.relationInfo) {
                        if (key === 'on_delete') {
                            this.relationInfo[key] = this.relation[key] === 'none' ? [] : [this.relation[key]]
                        } else {
                            this.relationInfo[key] = this.$tools.clone(this.relation[key])
                        }
                    }
                }
            },
            async initRelationList () {
                const data = await this.searchAssociationType({})
                this.relationList = data.info.map(({bk_asst_id: asstId, bk_asst_name: asstName}) => {
                    if (asstName.length) {
                        return {
                            id: asstId,
                            name: `${asstId}(${asstName})`
                        }
                    }
                    return {
                        id: asstId,
                        name: asstId
                    }
                })
            },
            async saveRelation () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                if (this.isEdit) {
                    await this.updateObjectAssociation({
                        id: this.relationInfo.id,
                        params: this.params,
                        config: {
                            requestId: 'updateObjectAssociation'
                        }
                    })
                } else {
                    await this.createObjectAssociation({
                        params: this.params,
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
