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
            <div class="relation-label cmdb-form-item" :class="{ 'is-error': errors.has('asstId') }">
                <cmdb-selector
                    class="bk-select-full-width"
                    :list="relationList"
                    v-validate="'required'"
                    name="asstId"
                    v-model="relationInfo['bk_asst_id']"
                ></cmdb-selector>
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
        <label class="form-label">
            <span class="label-text">
                {{$t('关联描述')}}
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstName') }">
                <bk-input type="textarea" class="cmdb-form-input"
                    name="asstName"
                    :maxlength="100"
                    v-validate="'singlechar'"
                    v-model.trim="relationInfo['bk_obj_asst_name']"
                    :placeholder="$t('请输入关联描述')">
                </bk-input>
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
        </label>
        <div class="form-label topo-preview" v-show="sourceModel.id && targetModel.id">
            <span class="label-text">
                {{$t('效果示意')}}
            </span>
            <div class="topo-image">
                <div class="model-item" :class="{ 'ispre': sourceModel.ispre }">
                    <div class="model-icon">
                        <i :class="['icon', sourceModel['bk_obj_icon']]"></i>
                    </div>
                    <span class="model-name">{{sourceModel['bk_obj_name']}}</span>
                </div>
                <div class="model-edge">
                    <div class="connection">
                        <span class="name">{{relationName}}</span>
                    </div>
                </div>
                <div class="model-item" :class="{ 'ispre': targetModel.ispre }">
                    <div class="model-icon">
                        <i :class="['icon', targetModel['bk_obj_icon']]"></i>
                    </div>
                    <span class="model-name">{{targetModel['bk_obj_name']}}</span>
                </div>
            </div>
            <div class="topo-text">{{sourceModel['bk_obj_name']}} {{relationName}} {{targetModel['bk_obj_name']}}</div>
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
            ...mapGetters('objectModelClassify', ['models', 'getModelById']),
            isSelfRelation () {
                return this.relationInfo.bk_obj_id === this.relationInfo.bk_asst_obj_id
            },
            mappingList () {
                const mappingList = [{
                    id: 'n:n',
                    name: 'N-N'
                }, {
                    id: '1:n',
                    name: '1-N'
                }, {
                    id: '1:1',
                    name: '1-1'
                }]
                if (this.isSelfRelation) {
                    mappingList.splice(1, 1)
                }
                return mappingList
            },
            objAsstId () {
                const {
                    relationInfo
                } = this
                if (relationInfo['bk_obj_id'].length && relationInfo['bk_asst_id'].length && relationInfo['bk_asst_obj_id'].length) {
                    return `${relationInfo['bk_obj_id']}_${relationInfo['bk_asst_id']}_${relationInfo['bk_asst_obj_id']}`
                }
                return ''
            },
            sourceModel () {
                return this.getModelById(this.relationInfo['bk_obj_id']) || {}
            },
            targetModel () {
                return this.getModelById(this.relationInfo['bk_asst_obj_id']) || {}
            },
            relationName () {
                const asstId = this.relationInfo['bk_asst_id']
                return (this.relationList.find(relation => relation.id === asstId) || {}).name
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
                const params = {
                    ...this.relationInfo,
                    ...{
                        bk_obj_asst_id: this.objAsstId
                    }
                }
                const res = await this.createObjectAssociation({
                    params: params,
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

    .topo-preview {
        .topo-image {
            display: flex;
            justify-content: space-between;
            padding: 10px 88px;
            background: #f3f8ff;

            .model-item {
                display: flex;
                flex-direction: column;
                align-items: center;

                .model-icon {
                    width: 46px;
                    height: 46px;
                    line-height: 46px;
                    border-radius: 50%;
                    text-align: center;
                    background: #fff;
                    box-shadow: 0px 2px 4px 0px rgba(147, 147, 147, 0.5);

                    .icon {
                        color: #3a84ff;
                        font-size: 24px;
                    }
                }

                .model-name {
                    font-size: 12px;
                    color: #868b97;
                    margin-top: 2px;
                }

                &.ispre {
                    .model-icon {
                        .icon {
                            color: #798aad;
                        }
                    }
                }
            }

            .model-edge {
                flex: 1;
                margin: 0 2px;

                .connection {
                    height: 46px;
                    position: relative;

                    .name {
                        position: absolute;
                        font-size: 12px;
                        color: #868b97;
                        padding: 2px 8px;
                        background: #fff;
                        top: 50%;
                        transform: translate(-50%, -50%);
                        left: 50%;
                        white-space: nowrap;
                    }

                    &::before {
                        content: '';
                        position: absolute;
                        left: 0;
                        top: 50%;
                        width: 100%;
                        height: 1px;
                        background: #c4c6cc;
                        margin-top: -0.5px;
                    }

                    &::after {
                        content: '';
                        position: absolute;
                        top: 50%;
                        right: -5px;
                        width: 0;
                        height: 0;
                        border: 4px solid transparent;
                        border-left: 8px solid #c4c6cc;
                        transform: translateY(-50%);
                    }
                }
            }
        }

        .topo-text {
            font-size: 12px;
            color: #868b97;
            text-align: center;
            margin-top: 8px;
        }
    }
</style>
