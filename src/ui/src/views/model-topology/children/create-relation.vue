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
        <label class="form-label exchange-icon-wrapper">
            <span class="label-text">
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input" disabled :value="getModelName(relationInfo['bk_asst_obj_id'])">
            </div>
            <span class="exchange-icon" @click="exchangeObjAsst">
                <i class="bk-icon icon-sort"></i>
            </span>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关联类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <ul class="relation-label cmdb-form-item clearfix" :class="{'is-error': errors.has('asstId')}">
                <li :class="{'active': relationInfo['bk_asst_id'] === relation.id}"
                    v-for="(relation, relationIndex) in relationList"
                    :key="relationIndex"
                    @click="relationInfo['bk_asst_id'] = relation.id">
                    {{relation.name}}
                </li>
            </ul>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关联描述"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstName')}">
                <input type="text" class="cmdb-form-input"
                name="asstName"
                v-validate="'required|singlechar'"
                v-model.trim="relationInfo['bk_obj_asst_name']"
                :placeholder="$t('ModelManagement[\'请输入关联描述\']')">
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源-目标约束"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('mapping')}">
                <cmdb-selector
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
            <bk-button type="primary" :loading="$loading('createObjectAssociation')" @click="saveRelation">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('Common["取消"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            toObjId: {
                type: String
            },
            fromObjId: {
                type: String
            },
            topoModelList: {
                type: Array
            }
        },
        data () {
            return {
                relationList: [],
                modelRelationList: [],
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
                    bk_obj_id: this.fromObjId,
                    bk_asst_obj_id: this.toObjId,
                    bk_asst_id: '',
                    bk_obj_asst_name: '',
                    mapping: ''
                }
            }
        },
        computed: {
            objAsstId () {
                let {
                    relationInfo
                } = this
                if (relationInfo['bk_obj_id'].length && relationInfo['bk_asst_id'].length && relationInfo['bk_asst_obj_id'].length) {
                    return `${relationInfo['bk_obj_id']}_${relationInfo['bk_asst_id']}_${relationInfo['bk_asst_obj_id']}`
                }
                return ''
            }
        },
        created () {
            this.searchModelRelationList()
            this.initRelationList()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createObjectAssociation',
                'searchAssociationType',
                'searchObjectAssociation'
            ]),
            getModelName (objId) {
                let model = this.$allModels.find(model => model['bk_obj_id'] === objId)
                if (model) {
                    return model['bk_obj_name']
                }
                return ''
            },
            exchangeObjAsst () {
                let {
                    relationInfo
                } = this;
                [relationInfo['bk_obj_id'], relationInfo['bk_asst_obj_id']] = [relationInfo['bk_asst_obj_id'], relationInfo['bk_obj_id']]
            },
            async initRelationList () {
                const data = await this.searchAssociationType({
                    params: {},
                    config: {
                        requestId: 'post_searchAssociationType',
                        fromCache: true
                    }
                })
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
                }).filter(relation => relation.id !== 'bk_mainline')
                this.relationInfo['bk_asst_id'] = this.relationList[0].id
            },
            async searchModelRelationList () {
                const [source, dest] = await Promise.all([this.searchAsSource(), this.searchAsDest()])
                this.modelRelationList = [...source, ...dest]
            },
            searchAsSource () {
                return this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_obj_id': this.relationInfo['bk_obj_id']
                        }
                    }
                })
            },
            searchAsDest () {
                return this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_asst_obj_id': this.relationInfo['bk_obj_id']
                        }
                    }
                })
            },
            async saveRelation () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                let params = {
                    ...this.relationInfo,
                    ...{
                        bk_obj_asst_id: this.objAsstId
                    }
                }
                const res = await this.createObjectAssociation({
                    params,
                    config: {
                        requestId: 'createObjectAssociation'
                    }
                })
                this.$emit('save', res)
                this.$emit('cancel')
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .exchange-icon-wrapper {
        position: relative;
    }
    .exchange-icon {
        position: absolute;
        display: inline-block;
        right: 0;
        top: 0;
        padding-top: 2px;
        width: 20px;
        height: 20px;
        border: 1px solid $cmdbBorderFocusColor;
        border-radius: 50%;
        color: $cmdbBorderFocusColor;
        font-size: 12px;
        line-height: 1;
        text-align: center;
        cursor: pointer;
        &:hover {
            color: $cmdbMainBtnColor;
            border-color: $cmdbMainBtnColor;
        }
        i {
            transform: scale(.8);
            font-weight: bold;
        }
    }
    .relation-label {
        li {
            float: left;
            margin: 5px 7px 5px 0;
            padding: 0 8px;
            height: 26px;
            line-height: 24px;
            font-size: 12px;
            border: 1px solid $cmdbTableBorderColor;
            background: #f5f7f9;
            cursor: pointer;
            &:hover {
                background: #fafafa;
            }
            &.active {
                color: #fff;
                background: $cmdbBorderFocusColor;
                border-color: $cmdbBorderFocusColor;
            }
        }
    }
</style>
