<template>
    <div>
        <label class="form-label">
            <span class="label-text">
                {{$t('源模型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <bk-input type="text" class="cmdb-form-input" disabled :value="getModelName(relationInfo['bk_obj_id'])"></bk-input>
            </div>
        </label>
        <label class="form-label exchange-icon-wrapper">
            <span class="label-text">
                {{$t('目标模型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item">
                <bk-input type="text" class="cmdb-form-input" disabled :value="getModelName(relationInfo['bk_asst_obj_id'])"></bk-input>
            </div>
            <span class="exchange-icon" @click="exchangeObjAsst">
                <i class="bk-icon icon-sort"></i>
            </span>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('关联类型')}}
                <span class="color-danger">*</span>
            </span>
            <ul class="relation-label cmdb-form-item clearfix" :class="{ 'is-error': errors.has('asstId') }">
                <li :class="{ 'active': relationInfo['bk_asst_id'] === relation.id }"
                    v-for="(relation, relationIndex) in relationList"
                    :key="relationIndex"
                    @click="relationInfo['bk_asst_id'] = relation.id">
                    {{relation.name}}
                </li>
            </ul>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('关联描述')}}
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstName') }">
                <bk-input type="text" class="cmdb-form-input"
                    name="asstName"
                    v-validate="'singlechar|length:256'"
                    v-model.trim="relationInfo['bk_obj_asst_name']"
                    :placeholder="$t('请输入关联描述')">
                </bk-input>
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('源-目标约束')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('mapping') }">
                <cmdb-selector style="width: 100%;"
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
            <bk-button theme="primary" :loading="$loading('createObjectAssociation')" @click="saveRelation">
                {{$t('提交')}}
            </bk-button>
            <bk-button theme="default" @click="cancel">
                {{$t('取消')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
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
                    bk_obj_id: this.fromObjId,
                    bk_asst_obj_id: this.toObjId,
                    bk_asst_id: '',
                    bk_obj_asst_name: '',
                    mapping: ''
                }
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', ['models']),
            objAsstId () {
                const {
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
                const model = this.models.find(model => model['bk_obj_id'] === objId)
                if (model) {
                    return model['bk_obj_name']
                }
                return ''
            },
            exchangeObjAsst () {
                const {
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
                this.relationList = data.info.map(({ bk_asst_id: asstId, bk_asst_name: asstName }) => {
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
                }).filter(relation => {
                    return relation.id !== 'bk_mainline'
                })
                this.relationInfo['bk_asst_id'] = this.relationList[0].id
            },
            async searchModelRelationList () {
                const [source, dest] = await Promise.all([this.searchAsSource(), this.searchAsDest()])
                this.modelRelationList = [...source, ...dest]
            },
            searchAsSource () {
                return this.searchObjectAssociation({
                    params: this.$injectMetadata({
                        condition: {
                            'bk_obj_id': this.relationInfo['bk_obj_id']
                        }
                    })
                })
            },
            searchAsDest () {
                return this.searchObjectAssociation({
                    params: this.$injectMetadata({
                        condition: {
                            'bk_asst_obj_id': this.relationInfo['bk_obj_id']
                        }
                    })
                })
            },
            async saveRelation () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                const params = {
                    ...this.relationInfo,
                    ...{
                        bk_obj_asst_id: this.objAsstId
                    }
                }
                const res = await this.createObjectAssociation({
                    params: this.$injectMetadata(params),
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
