<template>
    <div>
        <template v-if="!isCreateRelation">
            <label class="form-label exchange-icon-wrapper">
                <span class="label-text">
                    {{$t('ModelManagement["源模型"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item">
                    <input type="text" class="cmdb-form-input" disabled :value="getModelName(relationInfo['bk_obj_id'])">
                </div>
                <span class="exchange-icon" @click="exchangeObjAsst">
                    <i class="bk-icon icon-sort"></i>
                </span>
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
                    {{$t('ModelManagement["关系类型"]')}}
                    <span class="color-danger">*</span>
                </span>
                <ul class="relation-label cmdb-form-item clearfix" :class="{'is-error': errors.has('asstId')}">
                    <li :class="{'active': relationInfo['bk_asst_id'] === relation.id}"
                        v-for="(relation, relationIndex) in relationList"
                        :key="relationIndex"
                        @click="relationInfo['bk_asst_id'] = relation.id">
                        {{relation.name}}
                    </li>
                    <li @click="isCreateRelation = true">{{$t('ModelManagement["自定义关系"]')}}</li>
                </ul>
            </label>
            <label class="form-label">
                <span class="label-text">
                    {{$t('ModelManagement["关系描述"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item" :class="{'is-error': errors.has('asstName')}">
                    <input type="text" class="cmdb-form-input"
                    name="asstName"
                    v-validate="'required|singlechar'"
                    v-model.trim="relationInfo['bk_obj_asst_name']"
                    :placeholder="$t('ModelManagement[\'请输入关系描述\']')">
                    <p class="form-error">{{errors.first('asstName')}}</p>
                </div>
            </label>
            <div class="form-label">
                <span class="label-text">
                    {{$t('ModelManagement["关系约束"]')}}
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
                <bk-button type="primary" :loading="$loading(['updateAssociationType', 'createAssociationType'])" @click="saveRelation">
                    {{$t('ModelManagement["确定"]')}}
                </bk-button>
                <bk-button type="default" @click="cancel">
                    {{$t('ModelManagement["取消"]')}}
                </bk-button>
            </div>
        </template>
        <template v-else>
            <the-relation-type
                :saveBtnText="$t('ModelManagement[\'确定并返回\']')"
                @saved="saveRelationType"
                @cancel="isCreateRelation = false"
            >
            </the-relation-type>
        </template>
    </div>
</template>

<script>
    import theRelationType from '../relation-type'
    import { mapActions } from 'vuex'
    export default {
        components: {
            theRelationType
        },
        props: {
            toObjId: {
                type: String
            },
            fromObjId: {
                type: String
            }
        },
        data () {
            return {
                relationList: [],
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
                isCreateRelation: false,
                relationInfo: {
                    bk_obj_id: this.toObjId,
                    bk_asst_obj_id: this.fromObjId,
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
            this.initRelationList()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType'
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
                })
                this.relationInfo['bk_asst_id'] = this.relationList[0].id
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
                this.$emit('save', params)
                this.$emit('cancel')
            },
            cancel () {
                this.$emit('cancel')
            },
            saveRelationType () {
                this.$http.cancel('post_searchAssociationType')
                this.initRelationList()
                this.isCreateRelation = false
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
        top: 36px;
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
            &:last-child {
                margin-right: 0;
                padding: 0;
                border: none;
                color: $cmdbBorderFocusColor;
                background: #fff;
            }
        }
    }
</style>
