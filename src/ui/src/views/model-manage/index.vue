<template>
    <div class="group-wrapper">
        <cmdb-main-inject
            inject-type="prepend"
            :class="['btn-group', 'clearfix', {sticky: !!scrollTop}]">
            <div class="fl">
                <bk-button type="primary"
                    v-if="isAdminView"
                    :disabled="!authority.includes('update') || modelType === 'disabled'"
                    @click="showModelDialog(false)">
                    {{$t('ModelManagement["新增模型"]')}}
                </bk-button>
                <bk-button type="primary"
                    v-else
                    v-tooltip="$t('ModelManagement[\'新增模型提示\']')"
                    :disabled="!authority.includes('update') || modelType === 'disabled'"
                    @click="showModelDialog(false)">
                    {{$t('ModelManagement["新增模型"]')}}
                </bk-button>
                <bk-button type="default"
                    :disabled="!authority.includes('update') || modelType === 'disabled'"
                    @click="showGroupDialog(false)">
                    {{$t('ModelManagement["新建分组"]')}}
                </bk-button>
            </div>
            <div class="model-type-options fr">
                <bk-button class="model-type-button enable"
                    size="mini"
                    :type="modelType === 'enable' ? 'primary' : 'default'"
                    @click="modelType = 'enable'">
                    {{$t('ModelManagement["启用模型"]')}}
                </bk-button>
                <bk-button class="model-type-button disabled"
                    size="mini"
                    :disabled="!disabledClassifications.length"
                    :type="modelType === 'disabled' ? 'primary' : 'default'"
                    @click="modelType = 'disabled'">
                    {{$t('ModelManagement["停用模型"]')}}
                </bk-button>
            </div>
        </cmdb-main-inject>
        <ul class="group-list">
            <li class="group-item clearfix"
                v-for="(classification, classIndex) in currentClassifications"
                :key="classIndex">
                <div class="group-title">
                    <span>{{classification['bk_classification_name']}}</span>
                    <span class="number">({{classification['bk_objects'].length}})</span>
                    <template v-if="authority.includes('update') && isEditable(classification)">
                        <i class="icon-cc-edit text-primary"
                        @click="showGroupDialog(true, classification)"></i>
                        <i class="icon-cc-del text-primary"
                        @click="deleteGroup(classification)"></i>
                    </template>
                </div>
                <ul class="model-list clearfix" >
                    <li class="model-item"
                    :class="{
                        'ispaused': model['bk_ispaused'],
                        'ispre': isInner(model)
                    }"
                    v-for="(model, modelIndex) in classification['bk_objects']"
                    :key="modelIndex"
                    @click="modelClick(model)">
                        <div class="icon-box">
                            <i class="icon" :class="[model['bk_obj_icon']]"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</p>
                            <p class="model-id" :title="model['bk_obj_id']">{{model['bk_obj_id']}}</p>
                        </div>
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
                        <div class="label-title">
                            {{$t('ModelManagement["唯一标识"]')}}<span class="color-danger">*</span>
                        </div>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('classifyId')}">
                            <input type="text" class="cmdb-form-input"
                            name="classifyId"
                            :placeholder="$t('ModelManagement[\'请输入唯一标识\']')"
                            :disabled="groupDialog.isEdit"
                            v-model.trim="groupDialog.data['bk_classification_id']"
                            v-validate="'required|classifyId'">
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
                            <input type="text" 
                            class="cmdb-form-input"
                            name="classifyName"
                            :placeholder="$t('ModelManagement[\'请输入名称\']')"
                            v-model.trim="groupDialog.data['bk_classification_name']"
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
        <the-create-model
            :is-show.sync="modelDialog.isShow"
            :title="$t('ModelManagement[\'新增模型\']')"
            @confirm="saveModel">
        </the-create-model>
    </div>
</template>

<script>
    import cmdbMainInject from '@/components/layout/main-inject'
    import theCreateModel from '@/components/model-manage/_create-model'
    import theModel from './children'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    import {addMainScrollListener, removeMainScrollListener} from '@/utils/main-scroller'
    export default {
        components: {
            theModel,
            theCreateModel,
            cmdbMainInject
        },
        data () {
            return {
                scrollHandler: null,
                scrollTop: 0,
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
                    isShow: false
                },
                modelType: 'enable'
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'admin', 'isAdminView', 'isBusinessSelected']),
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            enableClassifications () {
                const enableClassifications = []
                this.classifications.forEach(classification => {
                    enableClassifications.push({
                        ...classification,
                        'bk_objects': classification['bk_objects'].filter(model => {
                            return !model['bk_ispaused'] && !['process', 'plat'].includes(model['bk_obj_id'])
                        })
                    })
                })
                return enableClassifications
            },
            disabledClassifications () {
                const disabledClassifications = []
                this.classifications.forEach(classification => {
                    const disabledModels = classification['bk_objects'].filter(model => {
                        return model['bk_ispaused'] && !['process', 'plat'].includes(model['bk_obj_id'])
                    })
                    if (disabledModels.length) {
                        disabledClassifications.push({
                            ...classification,
                            'bk_objects': disabledModels
                        })
                    }
                })
                return disabledClassifications
            },
            currentClassifications () {
                return this.modelType === 'enable' ? this.enableClassifications : this.disabledClassifications
            },
            authority () {
                if (this.isAdminView || this.isBusinessSelected) {
                    return ['search', 'update', 'delete']
                }
                return []
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["模型"]'))
            this.scrollHandler = event => {
                this.scrollTop = event.target.scrollTop
            }
            addMainScrollListener(this.scrollHandler)
        },
        beforeDestroy () {
            removeMainScrollListener(this.scrollHandler)
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
            isEditable (classification) {
                if (classification['bk_classification_type'] === 'inner') {
                    return false
                }
                if (this.isAdminView) {
                    return true
                }
                return !!this.$tools.getMetadataBiz(classification)
            },
            isInner (model) {
                const metadata = model.metadata || {}
                const label = metadata.label || {}
                return !this.$tools.getMetadataBiz(model)
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
                this.$validator.reset()
            },
            async saveGroup () {
                const res = await Promise.all([
                    this.$validator.validate('classifyId'),
                    this.$validator.validate('classifyName')
                ])
                if (res.includes(false)) {
                    return
                }
                let params = this.$injectMetadata({
                    bk_supplier_account: this.supplierAccount,
                    bk_classification_id: this.groupDialog.data['bk_classification_id'],
                    bk_classification_name: this.groupDialog.data['bk_classification_name']
                })
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
                    const res = await this.createClassification({
                        params,
                        config: {requestId: 'createClassification'}
                    })
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
                this.modelDialog.isShow = true
            },
            async saveModel (data) {
                const params = this.$injectMetadata({
                    bk_supplier_account: this.supplierAccount,
                    bk_obj_name: data['bk_obj_name'],
                    bk_obj_icon: data['bk_obj_icon'],
                    bk_classification_id: data['bk_classification_id'],
                    bk_obj_id: data['bk_obj_id'],
                    userName: this.userName
                })
                await this.createObject({params, config: {requestId: 'createModel'}})
                this.$http.cancel('post_searchClassificationsObjects')
                this.searchClassificationsObjects({
                    params: this.$injectMetadata()
                })
                this.modelDialog.isShow = false
            },
            modelClick (model) {
                this.$store.commit('objectModel/setActiveModel', model)
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push({
                    name: 'modelDetails',
                    params: {
                        modelId: model['bk_obj_id']
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .group-wrapper {
        padding: 76px 20px 20px 0;
    }
    .btn-group {
        position: absolute;
        top: 0;
        left: 0;
        width: calc(100% - 8px);
        padding: 20px;
        font-size: 0;
        background-color: #fff;
        z-index: 100;
        .bk-primary {
            margin-right: 10px;
        }
        &.sticky {
            box-shadow: 0 0 8px 1px rgba(0, 0, 0, 0.03);
        }
    }
    .model-type-options {
        margin: 6px 0;
        font-size: 0;
        text-align: right;
        position: relative;
        z-index: 1;
        .model-type-button {
            position: relative;
            margin: 0;
            font-size: 12px;
            &.enable {
                border-radius: 2px 0 0 2px;
                z-index: 2;
            }
            &.disabled {
                border-radius: 0 2px 2px 0;
                margin-left: -1px;
                z-index: 1;
            }
            &:hover {
                z-index: 2;
            }
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
            padding: 0 40px 0 0;
            line-height: 21px;
            color: #333948;
            &:before {
                content: "";
                display: inline-block;
                width:4px;
                height:14px;
                margin: 0 10px 0 0;
                vertical-align: middle;
                background: $cmdbBorderColor;
            }
            >span {
                display: inline-block;
                vertical-align: middle;
            }
            .number {
                color: $cmdbBorderColor;
            }
            >.text-primary {
                display: none;
                vertical-align: middle;
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
            width: calc((100% - 10px * 4) / 5);
            height: 70px;
            border: 1px solid $cmdbTableBorderColor;
            border-radius: 4px;
            cursor: pointer;
            &:nth-child(5n) {
                margin-right: 0;
            }
            &.ispaused {
                background: #fcfdfe;
                border-color: #dde4eb;
                .icon-box {
                    color: #96c2f7;
                }
                .model-name {
                    color: #bfc7d2;
                }
            }
            &.ispre {
                .icon-box {
                    color: #798aad;
                }
            }
            &:hover {
                border-color: $cmdbBorderFocusColor;
                background: #ebf4ff;
            }
            .icon-box {
                float: left;
                width: 66px;
                text-align: center;
                font-size: 32px;
                color: $cmdbBorderFocusColor;
                .icon {
                    line-height: 68px;
                }
            }
            .model-details {
                padding: 0 10px 0 0;
                overflow: hidden;
            }
            .model-name {
                margin-top: 16px;
                line-height: 19px;
                font-size: 14px;
                @include ellipsis;
            }
            .model-id {
                line-height: 16px;
                font-size: 12px;
                color: #bfc7d2;
                @include ellipsis;
            }
        }
    }
    .dialog {
        .dialog-content {
            padding: 20px 15px 20px 28px;
        }
        .title {
            font-size: 20px;
            color: #333948;
            line-height: 1;
        }
        .label-item,
        label {
            display: block;
            margin-bottom: 10px;
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
            }
            .label-title {
                font-size: 16px;
                line-height: 36px;
                vertical-align: middle;
                @include ellipsis;
            }
            .cmdb-form-item {
                display: inline-block;
                margin-right: 10px;
                width: 519px;
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
</style>
