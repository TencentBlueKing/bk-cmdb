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
        <div class="form-label">
            <span class="label-text">
                {{$t('源模型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('objId') }">
                <bk-select class="bk-select-full-width"
                    :disabled="relationInfo.ispre || isEdit"
                    v-validate="'required'"
                    name="objId"
                    v-model="relationInfo.bk_obj_id">
                    <bk-option-group
                        v-for="(group, index) in asstList"
                        :key="index"
                        :name="group.name">
                        <cmdb-auth-option
                            v-for="model in group.children"
                            :key="model.bk_obj_id"
                            :id="model.bk_obj_id"
                            :name="model.bk_obj_name"
                            :auth="{ type: $OPERATION.U_MODEL, relation: [model.id] }"
                            :ignore="relationInfo.ispre || isEdit">
                        </cmdb-auth-option>
                    </bk-option-group>
                </bk-select>
                <p class="form-error">{{errors.first('objId')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
        <div class="form-label exchange-icon-wrapper">
            <span class="label-text">
                {{$t('目标模型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstObjId') }">
                <bk-select class="bk-select-full-width"
                    :disabled="relationInfo.ispre || isEdit"
                    v-validate="'required'"
                    name="asstObjId"
                    v-model="relationInfo.bk_asst_obj_id">
                    <bk-option-group
                        v-for="(group, index) in asstList"
                        :key="index"
                        :name="group.name">
                        <cmdb-auth-option
                            v-for="model in group.children"
                            :key="model.bk_obj_id"
                            :id="model.bk_obj_id"
                            :name="model.bk_obj_name"
                            :auth="{ type: $OPERATION.U_MODEL, relation: [model.id] }"
                            :ignore="relationInfo.ispre || isEdit">
                        </cmdb-auth-option>
                    </bk-option-group>
                </bk-select>
                <p class="form-error">{{errors.first('asstObjId')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
            <span class="exchange-icon" @click="exchangeObjAsst" v-if="!(relationInfo.ispre || isReadOnly || isEdit)">
                <i class="bk-icon icon-sort"></i>
            </span>
        </div>
        <div class="form-label">
            <span class="label-text">
                {{$t('关联类型')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstId') }">
                <cmdb-selector
                    class="bk-select-full-width"
                    :searchable="true"
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
        <label class="form-label">
            <span class="label-text">
                {{$t('关联描述')}}
            </span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstName') }">
                <bk-input type="textarea" class="cmdb-form-input"
                    name="asstName"
                    :maxlength="100"
                    :placeholder="$t('请输入关联描述')"
                    :disabled="relationInfo.ispre || isReadOnly"
                    v-model.trim="relationInfo['bk_obj_asst_name']"
                    v-validate="'singlechar'">
                </bk-input>
                <p class="form-error">{{errors.first('asstName')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
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
                        <span class="name" :title="relationName">{{relationName}}</span>
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
        <div class="btn-group" v-if="!isReadOnly">
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
                originRelationInfo: {}
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications',
                'getModelById'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
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
                        classify['bk_objects'].forEach(model => {
                            if (!model.bk_ishidden) {
                                objects.push(model)
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
            },
            exchangeObjAsst () {
                const { relationInfo } = this;
                [relationInfo['bk_obj_id'], relationInfo['bk_asst_obj_id']] = [relationInfo['bk_asst_obj_id'], relationInfo['bk_obj_id']]
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-relation-wrapper {
        padding: 20px;
    }
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
            padding: 10px 68px;
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
                        max-width: 120px;
                        font-size: 12px;
                        color: #868b97;
                        padding: 2px 8px;
                        background: #fff;
                        top: 50%;
                        transform: translate(-50%, -50%);
                        left: 50%;
                        white-space: nowrap;
                        text-align: center;
                        @include ellipsis;
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
