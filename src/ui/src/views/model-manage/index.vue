<template>
    <div class="group-wrapper" :style="{ 'padding-top': showFeatureTips ? '114px' : '72px' }">
        <cmdb-main-inject
            :style="{ 'padding-top': showFeatureTips ? '10px' : '' }"
            inject-type="prepend"
            :class="['btn-group', 'clearfix', { sticky: !!scrollTop }]">
            <feature-tips
                :feature-name="'model'"
                :show-tips="showFeatureTips"
                :desc="$t('模型顶部提示')"
                :more-href="'https://docs.bk.tencent.com/cmdb/Introduction.html#ModelManagement'"
                @close-tips="showFeatureTips = false">
            </feature-tips>
            <div class="fl">
                <span v-if="isAdminView" style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_MODEL),
                        auth: [$OPERATION.C_MODEL]
                    }">
                    <bk-button theme="primary"
                        :disabled="!$isAuthorized($OPERATION.C_MODEL) || modelType === 'disabled'"
                        @click="showModelDialog()">
                        {{createModelBtn}}
                    </bk-button>
                </span>
                <span v-else style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_MODEL),
                        auth: [$OPERATION.C_MODEL]
                    }">
                    <bk-button theme="primary"
                        v-bk-tooltips="$t('新增模型提示')"
                        :disabled="!$isAuthorized($OPERATION.C_MODEL) || modelType === 'disabled'"
                        @click="showModelDialog()">
                        {{createModelBtn}}
                    </bk-button>
                </span>
                <span style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_MODEL_GROUP),
                        auth: [$OPERATION.C_MODEL_GROUP]
                    }">
                    <bk-button theme="default"
                        :disabled="!$isAuthorized($OPERATION.C_MODEL_GROUP) || modelType === 'disabled'"
                        @click="showGroupDialog(false)">
                        {{createGroupBtn}}
                    </bk-button>
                </span>
            </div>
            <div class="model-type-options fr">
                <bk-button class="model-type-button enable"
                    size="small"
                    :theme="modelType === 'enable' ? 'primary' : 'default'"
                    @click="modelType = 'enable'">
                    {{$t('启用模型')}}
                </bk-button>
                <bk-popover
                    :content="$t('停用模型提示')"
                    placenment="bottom"
                    v-if="!disabledClassifications.length">
                    <bk-button class="model-type-button disabled"
                        v-bk-tooltips="$t('停用模型提示')"
                        size="small"
                        :disabled="!disabledClassifications.length"
                        :theme="modelType === 'disabled' ? 'primary' : 'default'"
                        @click="modelType = 'disabled'">
                        {{$t('停用模型')}}
                    </bk-button>
                </bk-popover>
                <bk-button class="model-type-button disabled"
                    v-else
                    size="small"
                    :disabled="!disabledClassifications.length"
                    :theme="modelType === 'disabled' ? 'primary' : 'default'"
                    @click="modelType = 'disabled'">
                    {{$t('停用模型')}}
                </bk-button>
            </div>
            <div class="model-search-options fr">
                <bk-input class="search-model"
                    :clearable="true"
                    :right-icon="'bk-icon icon-search'"
                    v-model.trim="searchModel">
                </bk-input>
            </div>
        </cmdb-main-inject>
        <ul class="group-list">
            <li class="group-item clearfix"
                v-for="(classification, classIndex) in currentClassifications"
                :key="classIndex">
                <div class="group-title" v-bk-tooltips="classification.bk_classification_type === 'inner' ? groupToolTips : ''">
                    <span>{{classification['bk_classification_name']}}</span>
                    <span class="number">({{classification['bk_objects'].length}})</span>
                    <template v-if="isEditable(classification)">
                        <i class="icon-cc-plus text-primary"
                            :style="{ 'margin': '0 6px', color: $isAuthorized($OPERATION.C_MODEL) ? '' : '#e6e6e6 !important' }"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.C_MODEL),
                                auth: [$OPERATION.C_MODEL]
                            }"
                            @click="showModelDialog(classification.bk_classification_id)">
                        </i>
                        <i class="icon-cc-edit text-primary"
                            :style="{ 'margin-right': '4px', color: $isAuthorized($OPERATION.U_MODEL_GROUP) ? '' : '#e6e6e6 !important' }"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.U_MODEL_GROUP),
                                auth: [$OPERATION.U_MODEL_GROUP]
                            }"
                            @click="showGroupDialog(true, classification)">
                        </i>
                        <i class="icon-cc-del text-primary"
                            :style="{ color: $isAuthorized($OPERATION.D_MODEL_GROUP) ? '' : '#e6e6e6 !important' }"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.D_MODEL_GROUP),
                                auth: [$OPERATION.D_MODEL_GROUP]
                            }"
                            @click="deleteGroup(classification)">
                        </i>
                    </template>
                </div>
                <ul class="model-list clearfix">
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
            class="bk-dialog-no-padding bk-dialog-no-tools group-dialog dialog"
            :close-icon="false"
            :width="600"
            :mask-close="false"
            v-model="groupDialog.isShow">
            <div class="dialog-content">
                <p class="title">{{groupDialog.title}}</p>
                <div class="content">
                    <label>
                        <div class="label-title">
                            {{$t('唯一标识')}}<span class="color-danger">*</span>
                        </div>
                        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('classifyId') }">
                            <bk-input type="text" class="cmdb-form-input"
                                name="classifyId"
                                :placeholder="$t('请输入唯一标识')"
                                :disabled="groupDialog.isEdit"
                                v-model.trim="groupDialog.data['bk_classification_id']"
                                v-validate="'required|classifyId'">
                            </bk-input>
                            <p class="form-error" :title="errors.first('classifyId')">{{errors.first('classifyId')}}</p>
                        </div>
                        <i class="bk-icon icon-info-circle" v-bk-tooltips="$t('下划线，数字，英文小写的组合')"></i>
                    </label>
                    <label>
                        <span class="label-title">
                            {{$t('名称')}}
                        </span>
                        <span class="color-danger">*</span>
                        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('classifyName') }">
                            <bk-input type="text"
                                class="cmdb-form-input"
                                name="classifyName"
                                :placeholder="$t('请输入名称')"
                                v-model.trim="groupDialog.data['bk_classification_name']"
                                v-validate="'required|classifyName'">
                            </bk-input>
                            <p class="form-error" :title="errors.first('classifyName')">{{errors.first('classifyName')}}</p>
                        </div>
                    </label>
                </div>
            </div>
            <div slot="footer" class="footer">
                <bk-button theme="primary" :loading="$loading(['updateClassification', 'createClassification'])" @click="saveGroup">{{$t('保存')}}</bk-button>
                <bk-button theme="default" @click="hideGroupDialog">{{$t('取消')}}</bk-button>
            </div>
        </bk-dialog>
        <the-create-model
            :is-show.sync="modelDialog.isShow"
            :group-id.sync="modelDialog.groupId"
            :title="$t('新建模型')"
            @confirm="saveModel">
        </the-create-model>
    </div>
</template>

<script>
    import cmdbMainInject from '@/components/layout/main-inject'
    import theCreateModel from '@/components/model-manage/_create-model'
    import featureTips from '@/components/feature-tips/index'
    // import theModel from './children'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    import { addMainScrollListener, removeMainScrollListener } from '@/utils/main-scroller'
    export default {
        components: {
            // theModel,
            theCreateModel,
            cmdbMainInject,
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
                scrollHandler: null,
                scrollTop: 0,
                groupDialog: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('新建分组'),
                    data: {
                        bk_classification_id: '',
                        bk_classification_name: '',
                        id: ''
                    }
                },
                modelDialog: {
                    isShow: false,
                    groupId: ''
                },
                modelType: 'enable',
                searchModel: '',
                filterClassifications: [],
                groupToolTips: {
                    content: this.$t('内置模型组不支持添加和修改'),
                    placement: 'right'
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'admin', 'isAdminView', 'isBusinessSelected', 'featureTipsParams']),
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
                if (!this.searchModel) {
                    return this.modelType === 'enable' ? this.enableClassifications : this.disabledClassifications
                } else {
                    return this.filterClassifications
                }
            },
            createGroupBtn () {
                return this.isAdminView ? this.$t('新建分组') : this.$t('新建业务分组')
            },
            createModelBtn () {
                return this.isAdminView ? this.$t('新建模型') : this.$t('新建业务模型')
            }
        },
        watch: {
            searchModel (value) {
                if (!value) {
                    return
                }
                const searchResult = []
                const reg = new RegExp(value, 'gi')
                const currentClassifications = this.modelType === 'enable' ? this.enableClassifications : this.disabledClassifications
                const classifications = this.$tools.clone(currentClassifications)
                for (let i = 0; i < classifications.length; i++) {
                    classifications[i].bk_objects = classifications[i].bk_objects.filter(model => reg.test(model.bk_obj_name) || reg.test(model.bk_obj_id))
                    searchResult.push(classifications[i])
                }
                this.filterClassifications = searchResult
            },
            modelType () {
                this.searchModel = ''
            }
        },
        created () {
            this.scrollHandler = event => {
                this.scrollTop = event.target.scrollTop
            }
            addMainScrollListener(this.scrollHandler)
            this.searchClassificationsObjects({
                params: this.$injectMetadata()
            })
            if (this.$route.query.searchModel) {
                const hash = window.location.hash
                this.searchModel = this.$route.query.searchModel
                window.location.hash = hash.substring(0, hash.indexOf('?'))
            }
            this.showFeatureTips = this.featureTipsParams['model']
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
                return !this.$tools.getMetadataBiz(model)
            },
            showGroupDialog (isEdit, group) {
                if (isEdit) {
                    if (!this.$isAuthorized(this.$OPERATION.U_MODEL_GROUP)) return
                    this.groupDialog.data.id = group.id
                    this.groupDialog.title = this.$t('编辑分组')
                    this.groupDialog.data.bk_classification_id = group['bk_classification_id']
                    this.groupDialog.data.bk_classification_name = group['bk_classification_name']
                    this.groupDialog.data.id = group.id
                } else {
                    if (!this.$isAuthorized(this.$OPERATION.C_MODEL_GROUP)) return
                    this.groupDialog.title = this.$t('新建分组')
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
                const params = this.$injectMetadata({
                    bk_supplier_account: this.supplierAccount,
                    bk_classification_id: this.groupDialog.data['bk_classification_id'],
                    bk_classification_name: this.groupDialog.data['bk_classification_name']
                })
                if (this.groupDialog.isEdit) {
                    // eslint-disable-next-line
                    const res = await this.updateClassification({
                        id: this.groupDialog.data.id,
                        params,
                        config: {
                            requestId: 'updateClassification'
                        }
                    })
                    this.updateClassify({ ...params, ...{ id: this.groupDialog.data.id } })
                } else {
                    const res = await this.createClassification({
                        params,
                        config: { requestId: 'createClassification' }
                    })
                    this.updateClassify({ ...params, ...{ id: res.id } })
                }
                this.hideGroupDialog()
                this.searchModel = ''
            },
            deleteGroup (group) {
                if (!this.$isAuthorized(this.$OPERATION.D_MODEL_GROUP)) return
                this.$bkInfo({
                    title: this.$t('确认要删除此分组'),
                    confirmFn: async () => {
                        await this.deleteClassification({
                            id: group.id,
                            config: {
                                data: this.$injectMetadata({}, {
                                    inject: !!this.$tools.getMetadataBiz(group)
                                })
                            }
                        })
                        this.$store.commit('objectModelClassify/deleteClassify', group['bk_classification_id'])
                        this.searchModel = ''
                    }
                })
            },
            showModelDialog (groupId) {
                this.modelDialog.groupId = groupId || ''
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
                await this.createObject({ params, config: { requestId: 'createModel' } })
                this.$http.cancel('post_searchClassificationsObjects')
                this.searchClassificationsObjects({
                    params: this.$injectMetadata()
                })
                this.modelDialog.isShow = false
                this.modelDialog.groupId = ''
                this.searchModel = ''
            },
            modelClick (model) {
                const fullPath = this.searchModel ? `${this.$route.fullPath}?searchModel=${this.searchModel}` : this.$route.fullPath
                this.$store.commit('objectModel/setActiveModel', model)
                this.$router.push({
                    name: 'modelDetails',
                    params: {
                        modelId: model['bk_obj_id']
                    },
                    query: {
                        from: fullPath
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .group-wrapper {
        padding: 72px 0 20px 0;
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
    .model-search-options {
        .search-model {
            width: 240px;
        }
    }
    .model-type-options {
        margin: 0 0 0 10px;
        font-size: 0;
        text-align: right;
        .model-type-button {
            position: relative;
            margin: 0;
            font-size: 12px;
            height: 32px;
            line-height: 30px;
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
            margin: 0 40px 0 0;
            line-height: 21px;
            color: #333948;
            outline: 0;
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
            .icon-cc-plus {
                border: 1px solid #3c96ff;
                border-radius: 2px;
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
