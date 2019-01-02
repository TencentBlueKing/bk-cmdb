<template>
    <div>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input" disabled :value="getModelName(relationInfo['bk_obj_id'])">
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input" disabled :value="getModelName(relationInfo['bk_asst_obj_id'])">
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关联描述"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstName')}">
                <input type="text" class="cmdb-form-input"
                name="asstName"
                :disabled="relationInfo.ispre || !isEdit"
                v-model.trim="relationInfo['bk_obj_asst_name']"
                v-validate="'required|singlechar'">
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源-目标约束"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <cmdb-selector
                    :disabled="true"
                    :list="mappingList"
                    v-validate="'required'"
                    name="mapping"
                    v-model="relationInfo.mapping"
                ></cmdb-selector>
                <p class="form-error">{{errors.first('asstId')}}</p>
            </div>
        </label>
        <div class="btn-group" v-if="isEdit">
            <bk-button type="primary" @click="saveRelation">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button type="danger" @click="deleteRelation" :disabled="relationInfo.ispre">
                {{$t('ModelManagement["删除关联"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            objId: {
                type: String
            },
            asstId: {
                required: true
            },
            asstInfo: {
                type: Object
            },
            isEdit: {
                type: Boolean
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
                }
            }
        },
        created () {
            this.initData()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchObjectAssociation'
            ]),
            async initData () {
                let asstList = await this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_obj_id': this.objId
                        }
                    }
                })
                if (this.asstId !== '') {
                    this.relationInfo = asstList.find(asst => asst.id === this.asstId)
                } else {
                    this.relationInfo = this.$tools.clone(this.asstInfo)
                }
            },
            getModelName (objId) {
                let model = this.$allModels.find(model => model['bk_obj_id'] === objId)
                if (model) {
                    return model['bk_obj_name']
                }
                return ''
            },
            async saveRelation () {
                this.$emit('save', {
                    type: 'update',
                    params: {
                        ...this.relationInfo
                    }
                })
                this.$emit('cancel')
            },
            deleteRelation () {
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除关联关系？"]', this.relationInfo['bk_property_name'], {name: this.relationInfo['bk_property_name']}),
                    confirmFn: async () => {
                        this.$emit('save', {
                            type: 'delete',
                            params: {
                                ...this.relationInfo,
                                id: this.asstId
                            }
                        })
                        this.$emit('cancel')
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .button-delete {
        background-color: #fff;
        color: #ff5656;
    }
</style>
