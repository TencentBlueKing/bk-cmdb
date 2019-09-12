<template>
    <div class="group-wrapper" :style="{ 'padding-top': showFeatureTips ? '94px' : '52px' }">
        <cmdb-main-inject
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
                <span style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_MODEL),
                        auth: [$OPERATION.C_MODEL]
                    }">
                    <bk-button theme="primary"
                        :disabled="!$isAuthorized($OPERATION.C_MODEL) || modelType === 'disabled'"
                        @click="showModelDialog()">
                        {{$t('新建模型')}}
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
                        {{$t('新建分组')}}
                    </bk-button>
                </span>
            </div>
            <div class="model-type-options fr">
                <bk-button class="model-type-button enable"
                    :class="[{ 'model-type-button-active': modelType === 'enable' }]"
                    size="small"
                    @click="modelType = 'enable'">
                    {{$t('启用中')}}
                </bk-button>
                <span class="inline-block-middle" style="outline: 0;" v-bk-tooltips="disabledModelBtnText">
                    <bk-button class="model-type-button disabled"
                        :class="[{ 'model-type-button-active': modelType === 'disabled' }]"
                        size="small"
                        :disabled="!disabledClassifications.length"
                        @click="modelType = 'disabled'">
                        {{$t('已停用')}}
                    </bk-button>
                </span>
            </div>
            <div class="model-search-options fr">
                <bk-input class="search-model"
                    :clearable="true"
                    :right-icon="'bk-icon icon-search'"
                    :placeholder="$t('请输入关键字')"
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
                        <i class="icon-cc-add-line text-primary"
                            :class="[{ 'disabled': !$isAuthorized($OPERATION.C_MODEL) }]"
                            :style="{ 'margin': '0 6px' }"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.C_MODEL),
                                auth: [$OPERATION.C_MODEL]
                            }"
                            @click="showModelDialog(classification.bk_classification_id)">
                        </i>
                        <i class="icon-cc-edit text-primary"
                            :class="[{ 'disabled': !$isAuthorized($OPERATION.U_MODEL_GROUP) }]"
                            :style="{ 'margin-right': '4px' }"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.U_MODEL_GROUP),
                                auth: [$OPERATION.U_MODEL_GROUP]
                            }"
                            @click="showGroupDialog(true, classification)">
                        </i>
                        <i class="icon-cc-delete text-primary"
                            :class="[{ 'disabled': !$isAuthorized($OPERATION.D_MODEL_GROUP) }]"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.D_MODEL_GROUP),
                                auth: [$OPERATION.D_MODEL_GROUP]
                            }"
                            @click="deleteGroup(classification)">
                        </i>
                    </template>
                </div>
                <ul class="model-list clearfix">
                    <li class="model-item bgc-white"
                        :class="{
                            'ispaused': model['bk_ispaused'],
                            'ispre': model.ispre
                        }"
                        v-for="(model, modelIndex) in classification['bk_objects']"
                        :key="modelIndex">
                        <div class="info-model"
                            :class="{ 'radius': modelType === 'disabled' || classification['bk_classification_id'] === 'bk_biz_topo' }"
                            @click="modelClick(model)">
                            <div class="icon-box">
                                <i class="icon" :class="[model['bk_obj_icon']]"></i>
                            </div>
                            <div class="model-details">
                                <p class="model-name" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</p>
                                <p class="model-id" :title="model['bk_obj_id']">{{model['bk_obj_id']}}</p>
                            </div>
                        </div>
                        <div v-if="modelType !== 'disabled' && model.bk_classification_id !== 'bk_biz_topo'"
                            class="info-instance"
                            @click="handleGoInstance(model)">
                            <i class="icon-cc-share"></i>
                            <p>{{modelStatisticsSet[model.bk_obj_id] | instanceCount}}</p>
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
        
        <bk-dialog
            class="bk-dialog-no-padding"
            :width="400"
            :show-footer="false"
            :mask-close="false"
            v-model="sucessDialog.isShow">
            <div class="success-content">
                <i class="bk-icon icon-check-1"></i>
                <p>{{$t('模型创建成功')}}</p>
                <div class="btn-box">
                    <bk-button theme="primary" @click="handleGoInstance(curCreateModel)">{{$t('添加实例')}}</bk-button>
                    <bk-button @click="modelClick(curCreateModel)">{{$t('配置字段')}}</bk-button>
                    <bk-button @click="sucessDialog.isShow = false">{{$t('返回列表')}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import cmdbMainInject from '@/components/layout/main-inject'
    import theCreateModel from '@/components/model-manage/_create-model'
    import featureTips from '@/components/feature-tips/index'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    import { addMainScrollListener, removeMainScrollListener } from '@/utils/main-scroller'
    import { MENU_RESOURCE_HOST, MENU_RESOURCE_BUSINESS, MENU_RESOURCE_INSTANCE } from '@/dictionary/menu-symbol'
    export default {
        filters: {
            instanceCount (value) {
                if ([null, undefined].includes(value)) {
                    return 0
                }
                return value > 999 ? '999+' : value
            }
        },
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
                modelType: 'enable',
                searchModel: '',
                filterClassifications: [],
                modelStatisticsSet: {},
                curCreateModel: {},
                sucessDialog: {
                    isShow: false
                },
                groupToolTips: {
                    content: this.$t('内置模型组不支持添加和修改'),
                    placement: 'right'
                },
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
            disabledModelBtnText () {
                return this.disabledClassifications.length ? '' : this.$t('停用模型提示')
            }
        },
        watch: {
            searchModel (value) {
                if (!value) {
                    return
                }
                const searchResult = []
                const currentClassifications = this.modelType === 'enable' ? this.enableClassifications : this.disabledClassifications
                const classifications = this.$tools.clone(currentClassifications)
                for (let i = 0; i < classifications.length; i++) {
                    classifications[i].bk_objects = classifications[i].bk_objects.filter(model => {
                        const modelName = model.bk_obj_name
                        const modelId = model.bk_obj_id
                        return (modelName && modelName.indexOf(value) !== -1) || (modelId && modelId.indexOf(value) !== -1)
                    })
                    searchResult.push(classifications[i])
                }
                this.filterClassifications = searchResult
            },
            modelType () {
                this.searchModel = ''
            }
        },
        async created () {
            this.scrollHandler = event => {
                this.scrollTop = event.target.scrollTop
            }
            addMainScrollListener(this.scrollHandler)
            this.getModelStatistics()
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
                'getClassificationsObjectStatistics',
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
            async getModelStatistics () {
                const modelStatisticsSet = {}
                const res = await this.getClassificationsObjectStatistics({
                    config: {
                        requestId: 'getClassificationsObjectStatistics'
                    }
                })
                res.forEach(item => {
                    modelStatisticsSet[item.bk_obj_id] = item.instance_count
                })
                this.modelStatisticsSet = modelStatisticsSet
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
                const createModel = await this.createObject({ params, config: { requestId: 'createModel' } })
                this.curCreateModel = createModel
                this.sucessDialog.isShow = true
                this.$http.cancel('post_searchClassificationsObjects')
                this.getModelStatistics()
                this.searchClassificationsObjects({
                    params: this.$injectMetadata()
                })
                this.modelDialog.isShow = false
                this.modelDialog.groupId = ''
                this.searchModel = ''
            },
            modelClick (model) {
                this.$store.commit('objectModel/setActiveModel', model)
                this.$router.push({
                    name: 'modelDetails',
                    params: {
                        modelId: model['bk_obj_id']
                    }
                })
            },
            handleGoInstance (model) {
                this.sucessDialog.isShow = false
                const map = {
                    host: MENU_RESOURCE_HOST,
                    biz: MENU_RESOURCE_BUSINESS
                }
                if (map.hasOwnProperty(model.bk_obj_id)) {
                    this.$router.push({
                        name: map[model.bk_obj_id]
                    })
                } else {
                    this.$router.push({
                        name: MENU_RESOURCE_INSTANCE,
                        params: {
                            objId: model.bk_obj_id
                        }
                    })
                }
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
        top: 58px;
        left: 0;
        width: calc(100% - 17px);
        padding: 0 20px 20px;
        font-size: 0;
        background-color: #fafbfd;
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
                border-right-color: #3a84ff;
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
            &-active {
                border-color: #3a84ff;
                color: #3a84ff;
            }
        }
    }
    .group-list {
        padding: 0 20px;
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
                &.disabled {
                    opacity: 1 !important;
                    color: #c4c6cc !important;
                }
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
            display: flex;
            position: relative;
            float: left;
            margin: 10px 10px 0 0;
            width: calc((100% - 10px * 4) / 5);
            height: 70px;
            border: 1px solid $cmdbTableBorderColor;
            border-radius: 4px;
            background-color: #ffffff;
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
                .info-instance {
                    display: block;
                }
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
                padding: 0 4px 0 0;
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
            .info-model {
                flex: 1;
                width: 0;
                border-radius: 4px 0 0 4px;
                &:hover {
                    background-color: #f0f5ff;
                }
                &.radius {
                    border-radius: 4px;
                }
            }
            .info-instance {
                display: none;
                width: 70px;
                padding: 0 8px 0 6px;
                text-align: center;
                color: #c3cdd7;
                border-radius: 0 4px 4px 0;
                &:hover {
                    background-color: #f0f5ff;
                    p {
                        color: #3a84ff;
                    }
                }
                .icon-cc-share {
                    font-size: 14px;
                    margin-top: 16px;
                    color: #3a84ff;
                }
                p {
                    font-size: 16px;
                    padding-top: 2px;
                    @include ellipsis;
                }
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
            padding-bottom: 14px;
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
            font-size: 0;
            text-align: right;
            .bk-primary {
                margin-right: 10px;
            }
        }
    }
    .success-content {
        text-align: center;
        padding-bottom: 46px;
        p {
            color: #444444;
            font-size: 24px;
            padding: 10px 0 20px;
        }
        .icon-check-1 {
            width: 58px;
            height: 58px;
            line-height: 58px;
            font-size: 30px;
            font-weight: bold;
            color: #fff;
            border-radius: 50%;
            background-color: #2dcb56;
            text-align: center;
        }
        .btn-box {
            font-size: 0;
            .bk-button {
                margin: 0 5px;
            }
        }
    }
</style>
