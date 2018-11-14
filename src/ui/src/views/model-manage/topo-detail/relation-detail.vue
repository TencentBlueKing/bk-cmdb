<template>
    <div>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input" disabled :value="relationInfo['bk_obj_id']">
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input" disabled :value="relationInfo['bk_obj_id']">
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系描述"]')}}
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
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系约束"]')}}
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
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading('updateObjectAssociation')" @click="saveRelation">
                {{$t('ModelManagement["确定"]')}}
            </bk-button>
            <bk-button type="danger" @click="deleteRelation" :disabled="$loading('deleteObjectAssociation')">
                {{$t('ModelManagement["删除关系"]')}}
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
                type: Number
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
                'searchObjectAssociation',
                'updateObjectAssociation',
                'deleteObjectAssociation'
            ]),
            async initData () {
                let asstList = await this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_obj_id': this.objId
                        }
                    }
                })
                this.relationInfo = asstList.find(asst => asst.id === this.asstId)
            },
            async saveRelation () {
                await this.updateObjectAssociation({
                    id: this.asstId,
                    params: {
                        bk_obj_asst_name: this.relationInfo['bk_obj_asst_name']
                    },
                    config: {
                        requestId: 'updateObjectAssociation'
                    }
                })
                this.$emit('cancel')
            },
            deleteRelation () {
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除关联关系？"]', this.relationInfo['bk_property_name'], {name: this.relationInfo['bk_property_name']}),
                    confirmFn: async () => {
                        await this.deleteObjectAssociation({
                            id: this.asstId,
                            config: {
                                requestId: 'deleteObjectAssociation'
                            }
                        }).then(() => {
                            this.$http.cancel(`post_searchObjectAssociation_${this.objId}`)
                        })
                    }
                })
                this.$emit('save')
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
