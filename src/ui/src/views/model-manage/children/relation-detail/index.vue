<template>
    <div>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["唯一标识"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input" v-model.trim="relationInfo['bk_obj_asst_id']">
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["别名"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input" v-model.trim="relationInfo['bk_obj_asst_name']">
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["模型源"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="relationInfo['bk_obj_id']"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="relationInfo['bk_asst_obj_id']"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="relationInfo['bk_asst_id']"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源-目标约束"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="mappingList"
                :selected.sync="relationInfo.mapping"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <div class="radio-box">
            <label class="label-text">
                {{$t('ModelManagement["联动删除"]')}}
            </label>
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" id="delete_dest" value="delete_dest" v-model="relationInfo['on_delete']">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["源不存在联动删除目标"]')}}</span>
            </label>
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" id="delete_src" value="delete_src" v-model="relationInfo['on_delete']">
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
                    bk_obj_asst_id: '',
                    bk_obj_asst_name: '',
                    bk_obj_id: '',
                    bk_asst_obj_id: '',
                    bk_asst_id: '',
                    mapping: '',
                    on_delete: []
                },
                relationList: [{
                    'id': 1,
                    'bk_asst_id': 'belong',
                    'bk_asst_name': '属于',
                    'src_des': '属于',
                    'dest_des': '被属于',
                    'direction': 'none'
                }],
                modelList: [],
                selected: ''
            }
        },
        created () {
            this.initRelationList()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createObjectAssociation',
                'updateObjectAssociation',
                'searchAssociationType'
            ]),
            async initRelationList () {
                this.relationList = await this.searchAssociationType({})
            },
            async saveRelation () {
                let params = this.relationInfo
                if (this.isEdit) {
                    await this.updateObjectAssociation({
                        params,
                        config: {
                            requestId: 'updateObjectAssociation'
                        }
                    })
                } else {
                    await this.createObjectAssociation({
                        params,
                        config: {
                            requestId: 'createObjectAssociation'
                        }
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-relation-wrapper {
        padding: 20px;
    }
</style>
