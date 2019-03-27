<template>
    <div class="group-wrapper" @scroll="handleWrapperScroll">
        <p class="btn-group" :class="{sticky: wrapperScroll}">
            <bk-button type="primary"
                :disabled="!authority.includes('update')"
                @click="showModelDialog(false)">
                {{$t('ModelManagement["新建模型"]')}}
            </bk-button>
            <bk-button type="default"
                :disabled="!authority.includes('update')"
                @click="showGroupDialog(false)">
                {{$t('ModelManagement["新建分组"]')}}
            </bk-button>
        </p>
        <ul class="group-list">
            <li class="group-item clearfix"
                v-for="(classification, classIndex) in localClassifications"
                :key="classIndex">
                <p class="group-title">
                    <span>{{classification['bk_classification_name']}}</span>
                    <span class="number">({{classification['bk_objects'].length}})</span>
                    <template v-if="authority.includes('update')">
                        <i class="icon-cc-edit text-primary"
                        v-if="classification['bk_classification_type'] !== 'inner'"
                        @click="showGroupDialog(true, classification)"></i>
                        <i class="icon-cc-del text-primary"
                        v-if="classification['bk_classification_type'] !== 'inner'"
                        @click="deleteGroup(classification)"></i>
                    </template>
                </p>
                <ul class="model-list clearfix">
                    <li class="model-item"
                    :class="{'ispaused': model['bk_ispaused']}"
                    v-for="(model, modelIndex) in classification['bk_objects']"
                    :key="modelIndex"
                    @click="modelClick(model)">
                        <div class="icon-box">
                            <i class="icon" :class="[model['bk_obj_icon'], {ispre: model['ispre']}]"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</p>
                            <p class="model-id" :title="model['bk_obj_id']">{{model['bk_obj_id']}}</p>
                        </div>
                        <span class="paused-info" v-if="model['bk_ispaused']">
                            {{$t('ModelManagement["已停用"]')}}
                        </span>
                    </li>
                </ul>
            </li>
        </ul>
        <bk-dialog
            class="group-dialog dialog"
            :close-icon="false"
            :hasHeader="false"
            :width="600"
            :padding="0"
            :quick-close="false"
            :is-show.sync="groupDialog.isShow">
            <div slot="content" class="dialog-content">
                <p class="title">{{groupDialog.title}}</p>
                <div class="content">
                    <label>
                        <span class="label-title">
                            {{$t('ModelManagement["唯一标识"]')}}
                        </span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('classifyId')}">
                            <input type="text" class="cmdb-form-input"
                            v-model.trim="groupDialog.data['bk_classification_id']"
                            name="classifyId"
                            v-validate="'required|classifyId'"
                            :disabled="groupDialog.isEdit">
                            <p class="form-error">{{errors.first('classifyId')}}</p>
                        </div>
                        <i class="bk-icon icon-info-circle" v-tooltip="$t('ModelManagement[\'下划线，数字，英文小写的组合\']')"></i>
                    </label>
                    <label>
                        <span class="label-title">
                            {{$t('ModelManagement["名称"]')}}
                        </span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('classifyName')}">
                            <input type="text" class="cmdb-form-input"
                            v-model.trim="groupDialog.data['bk_classification_name']"
                            name="classifyName"
                            v-validate="'required|classifyName'">
                            <p class="form-error">{{errors.first('classifyName')}}</p>
                        </div>
                    </label>
                </div>
            </div>
            <div slot="footer" class="footer">
                <bk-button type="primary" :loading="$loading(['updateClassification', 'createClassification'])" @click="saveGroup">{{$t("Common['保存']")}}</bk-button>
                <bk-button type="default" @click="hideGroupDialog">{{$t("Common['取消']")}}</bk-button>
            </div>
        </bk-dialog>
        <bk-dialog
            class="model-dialog dialog"
            :close-icon="false"
            :hasHeader="false"
            :width="600"
            :padding="0"
            :quick-close="false"
            :is-show.sync="modelDialog.isShow">
            <div slot="content" class="dialog-content">
                <p class="title">{{$t('ModelManagement["新建模型"]')}}</p>
                <div class="content clearfix">
                    <div class="content-left" @click="modelDialog.isIconListShow = true">
                        <div class="icon-wrapper">
                            <i :class="modelDialog.data['bk_obj_icon']"></i>
                        </div>
                        <div class="text">{{$t('ModelManagement["点击切换"]')}}</div>
                    </div>
                    <div class="content-right">
                        <div class="label-item">
                            <span class="label-title">{{$t('ModelManagement["所属分组"]')}}</span>
                            <span class="color-danger">*</span>
                            <div class="cmdb-form-item" :class="{'is-error': errors.has('modelGroup')}">
                                <cmdb-selector
                                    class="selector-box"
                                    name="modelGroup"
                                    setting-key="bk_classification_id"
                                    display-key="bk_classification_name"
                                    :content-max-height="200"
                                    :selected.sync="modelDialog.data['bk_classification_id']"
                                    :list="modelDialog.classificationList"
                                    v-validate="'required'"
                                    v-model="modelDialog.data['bk_classification_id']"
                                ></cmdb-selector>
                                <p class="form-error">{{errors.first('modelGroup')}}</p>
                            </div>
                        </div>
                        <label>
                            <span class="label-title">{{$t('ModelManagement["唯一标识"]')}}</span>
                            <span class="color-danger">*</span>
                            <div class="cmdb-form-item" :class="{'is-error': errors.has('modelId')}">
                                <input type="text" class="cmdb-form-input"
                                name="modelId"
                                v-model.trim="modelDialog.data['bk_obj_id']"
                                v-validate="'required|modelId'">
                                <p class="form-error">{{errors.first('modelId')}}</p>
                            </div>
                            <i class="bk-icon icon-info-circle" v-tooltip="$t('ModelManagement[\'下划线，数字，英文小写的组合\']')"></i>
                        </label>
                        <label>
                            <span class="label-title">{{$t('ModelManagement["名称"]')}}</span>
                            <span class="color-danger">*</span>
                            <div class="cmdb-form-item" :class="{'is-error': errors.has('modelName')}">
                                <input type="text" class="cmdb-form-input"
                                name="modelName"
                                v-validate="'required|singlechar'"
                                v-model.trim="modelDialog.data['bk_obj_name']">
                                <p class="form-error">{{errors.first('modelName')}}</p>
                            </div>
                            <i class="bk-icon icon-info-circle" v-tooltip="$t('ModelManagement[\'请填写模型名\']')"></i>
                        </label>
                    </div>
                </div>
                <div class="model-icon-wrapper" v-if="modelDialog.isIconListShow">
                    <span class="back" @click="modelDialog.isIconListShow = false">
                        <i class="bk-icon icon-back2"></i>
                    </span>
                    <the-choose-icon
                        v-model="modelDialog.data['bk_obj_icon']"
                        @chooseIcon="modelDialog.isIconListShow = false"
                    ></the-choose-icon>
                </div>
            </div>
            <div slot="footer" class="footer">
                <bk-button type="primary" @click="saveModel">{{$t("Common['保存']")}}</bk-button>
                <bk-button type="default" @click="hideModelDialog">{{$t("Common['取消']")}}</bk-button>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import theChooseIcon from '@/components/model-manage/_choose-icon'
    import theModel from './children'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    export default {
        components: {
            theChooseIcon,
            theModel
        },
        data () {
            return {
                wrapperScroll: 0,
                groupDialog: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('ModelManagement["新建分组"]'),
                    data: {
                        bk_classification_id: '',
                        bk_classification_name: '',
                        id: ''
                    }
                },
                model: {
                    modelId: 'modelId'
                },
                modelDialog: {
                    isShow: false,
                    isEdit: false,
                    isIconListShow: false,
                    classificationList: [],
                    data: {
                        bk_classification_id: '',
                        bk_obj_icon: 'icon-cc-default',
                        bk_obj_id: '',
                        bk_obj_name: ''
                    }
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'admin']),
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            localClassifications () {
                let localClassifications = []
                this.classifications.forEach(classification => {
                    if (classification['bk_classification_id'] === 'bk_host_manage') {
                        const currentClassification = this.$tools.clone(classification)
                        currentClassification['bk_objects'] = classification['bk_objects'].filter(({bk_obj_id: objId}) => !['process', 'plat'].includes(objId))
                        localClassifications.push({...currentClassification, ...{isModelShow: false}})
                    } else {
                        localClassifications.push({...classification, ...{isModelShow: false}})
                    }
                })
                this.modelDialog.classificationList = localClassifications.filter(({bk_classification_id: classificationId}) => !['bk_biz_topo', 'bk_host_manage', 'bk_organization'].includes(classificationId))
                return localClassifications
            },
            authority () {
                return this.admin ? ['search', 'update', 'delete'] : []
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["模型"]'))
        },
        methods: {
            ...mapMutations('objectModelClassify', [
                'updateClassify',
                'deleteClassify'
            ]),
            ...mapActions('objectModelClassify', [
                'searchClassificationsObjects',
                'createClassification',
                'updateClassification',
                'deleteClassification'
            ]),
            ...mapActions('objectModel', [
                'createObject'
            ]),
            handleWrapperScroll () {
                this.wrapperScroll = this.$el.scrollTop
            },
            showGroupDialog (isEdit, group) {
                if (isEdit) {
                    this.groupDialog.data.id = group.id
                    this.groupDialog.title = this.$t('ModelManagement["编辑分组"]')
                    this.groupDialog.data.bk_classification_id = group['bk_classification_id']
                    this.groupDialog.data.bk_classification_name = group['bk_classification_name']
                    this.groupDialog.data.id = group.id
                } else {
                    this.groupDialog.title = this.$t('ModelManagement["新建分组"]')
                    this.groupDialog.data.bk_classification_id = ''
                    this.groupDialog.data.bk_classification_name = ''
                    this.groupDialog.data.id = ''
                }
                this.groupDialog.isEdit = isEdit
                this.groupDialog.isShow = true
            },
            hideGroupDialog () {
                this.groupDialog.isShow = false
            },
            async saveGroup () {
                const res = await Promise.all([
                    this.$validator.validate('classifyId'),
                    this.$validator.validate('classifyName')
                ])
                if (res.includes(false)) {
                    return
                }
                let params = {
                    bk_supplier_account: this.supplierAccount,
                    bk_classification_id: this.groupDialog.data['bk_classification_id'],
                    bk_classification_name: this.groupDialog.data['bk_classification_name']
                }
                if (this.groupDialog.isEdit) {
                    const res = await this.updateClassification({
                        id: this.groupDialog.data.id,
                        params,
                        config: {
                            requestId: 'updateClassification'
                        }
                    })
                    this.updateClassify({...params, ...{id: this.groupDialog.data.id}})
                } else {
                    const res = await this.createClassification({params, config: {requestId: 'createClassification'}})
                    this.updateClassify({...params, ...{id: res.id}})
                }
                this.hideGroupDialog()
            },
            deleteGroup (group) {
                this.$bkInfo({
                    title: this.$t('ModelManagement["确认要删除此分组"]'),
                    confirmFn: async () => {
                        await this.deleteClassification({
                            id: group.id
                        })
                        this.$store.commit('objectModelClassify/deleteClassify', group['bk_classification_id'])
                    }
                })
            },
            showModelDialog () {
                this.modelDialog.data['bk_obj_icon'] = 'icon-cc-default'
                this.modelDialog.data['bk_obj_id'] = ''
                this.modelDialog.data['bk_obj_name'] = ''
                this.modelDialog.data['bk_classification_id'] = ''
                this.$validator.reset()
                this.modelDialog.isShow = true
            },
            hideModelDialog () {
                this.modelDialog.isShow = false
            },
            async saveModel () {
                const res = await Promise.all([
                    this.$validator.validate('modelGroup'),
                    this.$validator.validate('modelId'),
                    this.$validator.validate('modelName')
                ])
                if (res.includes(false)) {
                    return
                }
                let params = {
                    bk_supplier_account: this.supplierAccount,
                    bk_obj_name: this.modelDialog.data['bk_obj_name'],
                    bk_obj_icon: this.modelDialog.data['bk_obj_icon'],
                    bk_classification_id: this.modelDialog.data['bk_classification_id'],
                    bk_obj_id: this.modelDialog.data['bk_obj_id'],
                    userName: this.userName
                }
                await this.createObject({params, config: {requestId: 'createModel'}})
                this.$http.cancel('post_searchClassificationsObjects')
                this.searchClassificationsObjects({})
                this.hideModelDialog()
            },
            modelClick (model) {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push(`model/details/${model['bk_obj_id']}`)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .group-wrapper {
        position: relative;
        height: 100%;
        padding: 0;
        overflow-y: auto;
    }
    .btn-group {
        position: sticky;
        top: 0;
        left: 0;
        padding: 20px;
        font-size: 0;
        z-index: 2;
        background-color: #fff;
        .bk-primary {
            margin-right: 10px;
        }
        &.sticky {
            box-shadow: 0 0 8px 1px rgba(0, 0, 0, 0.03);
        }
    }
    .group-list {
        padding: 0 20px 20px;
        .group-item {
            position: relative;
            padding: 10px 0 20px;
            >.icon-angle-double-down {
                position: absolute;
                left: 50%;
                bottom: -5px;
                margin-left: -5px;
                padding: 5px;
                font-size: 12px;
                cursor: pointer;
                transition: all .2s;
                &.rotate {
                    transform: rotate(180deg);
                }
            }
        }
        .group-title {
            display: inline-block;
            padding: 0 40px 0 8px;
            border-left: 4px solid $cmdbBorderColor;
            line-height: 14px;
            color: #333948;
            >span {
                display: inline-block;
            }
            .number {
                color: $cmdbBorderColor;
            }
            >.text-primary {
                display: none;
                vertical-align: top;
                cursor: pointer;
            }
            &:hover {
                >.text-primary {
                    display: inline-block;
                }
            }
        }
    }
    .model-list {
        padding-left: 12px;
        overflow: hidden;
        transition: height .2s;
        .model-item {
            position: relative;
            float: left;
            margin: 10px 10px 0 0;
            width: 260px;
            height: 70px;
            border: 1px solid $cmdbTableBorderColor;
            border-radius: 4px;
            cursor: pointer;
            &.ispaused {
                background: #fafbfd;
                opacity: .6;
                &:after {
                    content: '';
                    display: inline-block;
                    position: absolute;
                    top: -33px;
                    right: -33px;
                    border: 32px solid transparent;
                    border-bottom-color: $cmdbDangerColor;
                    transform: rotate(45deg);
                }
                .paused-info {
                    position: absolute;
                    right: -2px;
                    top: 7px;
                    font-size: 12px;
                    z-index: 1;
                    color: #fff;
                    transform: rotate(45deg) scale(.8);
                }
            }
            &:hover {
                border-color: $cmdbBorderFocusColor;
                background: #ebf4ff;
            }
            .icon-box {
                float: left;
                width: 50px;
                .icon {
                    padding-left: 18px;
                    font-size: 32px;
                    line-height: 70px;
                    color: $cmdbBorderFocusColor;
                    &.ispre {
                        color: #868b97;
                    }
                }
            }
            .model-details {
                float: left;
                width: 208px;
                line-height: 16px;
                margin-top: 20px;
                padding: 0 10px;
            }
            .model-name {
                font-size: 14px;
                @include ellipsis;
            }
            .model-id {
                font-size: 12px;
                color: #bfc7d2;
                @include ellipsis;
            }
        }
    }
    .dialog {
        .dialog-content {
            padding: 20px 10px;
            .content {
                padding: 0 10px;
            }
        }
        .title {
            margin-bottom: 30px;
            font-size: 20px;
            color: #333948;
            line-height: 1;
        }
        .label-item,
        label {
            display: block;
            margin-bottom: 20px;
            font-size: 0;
            &:last-child {
                margin: 0;
            }
            .color-danger {
                display: inline-block;
                font-size: 16px;
                width: 15px;
                text-align: center;
                vertical-align: middle;
            }
            .icon-info-circle {
                font-size: 18px;
                color: $cmdbBorderColor;
                // line-height: 36px;
                // vertical-align: middle;
            }
            .label-title {
                display: inline-block;
                width: 85px;
                text-align: right;
                font-size: 16px;
                line-height: 36px;
                vertical-align: middle;
                @include ellipsis;
            }
            .cmdb-form-item {
                display: inline-block;
                margin-right: 10px;
                width: calc(100% - 130px);
                vertical-align: middle;
            }
        }
        .footer {
            padding: 0 24px;
            font-size: 0;
            text-align: right;
            .bk-primary {
                margin-right: 10px;
            }
        }
    }
    .group-dialog {
        .dialog-content {
            .content {
                padding: 30px 10px 40px;
            }
        }
    }
    .model-dialog {
        .dialog-content {
            position: relative;
        }
        .content-left {
            float: left;
            width: 93px;
            height: 100px;
            border: 1px solid #dde4eb;
            border-radius: 4px 4px 0 0;
            cursor: pointer;
            .icon-wrapper {
                width: 100%;
                height: 68px;
                font-size: 38px;
                text-align: center;
                i {
                    vertical-align: top;
                    line-height: 68px;
                    color: $cmdbBorderFocusColor;
                }
            }
            .text {
                height: 30px;
                border-top: 1px solid #dde4eb;
                text-align: center;
                line-height: 30px;
                background: #ebf4ff;
            }
        }
        .content-right {
            float: right;
            width: 460px;
        }
        .model-icon-wrapper {
            position: absolute;
            left: 0;
            top:0;
            width: 100%;
            height: calc(100% + 60px);
            background: #fff;
            .back {
                position: absolute;
                right: -47px;
                top: 0;
                width: 44px;
                height: 44px;
                padding: 7px;
                cursor: pointer;
                font-size: 18px;
                text-align: center;
                background: #2f2f2f;
                color: #fff;
            }
        }
    }
</style>
