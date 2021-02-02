<template>
    <div class="group-wrapper"
        v-bkloading="{ isLoading: mainLoading }"
        :style="{ 'padding-top': topPadding + 'px' }">
        <cmdb-main-inject ref="mainInject"
            inject-type="prepend"
            :class="['btn-group', 'clearfix', { sticky: !!scrollTop }]">
            <cmdb-tips
                class="mb10"
                tips-key="modelTips"
                more-link="https://bk.tencent.com/docs/markdown/配置平台/产品白皮书/产品功能/Model.md">
                {{$t('模型顶部提示')}}
            </cmdb-tips>
            <div class="fl">
                <cmdb-auth :auth="{ type: $OPERATION.C_MODEL }">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        :disabled="disabled || modelType === 'disabled'"
                        @click="showModelDialog()">
                        {{$t('新建模型')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth :auth="{ type: $OPERATION.C_MODEL_GROUP }">
                    <bk-button slot-scope="{ disabled }"
                        theme="default"
                        :disabled="disabled || modelType === 'disabled'"
                        @click="showGroupDialog(false)">
                        {{$t('新建分组')}}
                    </bk-button>
                </cmdb-auth>
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
                <div class="group-title">
                    <div class="title-info"
                        v-bk-tooltips="{
                            disabled: isEditable(classification),
                            content: $t('内置模型组不支持添加和修改'),
                            placement: 'right'
                        }">
                        <span class="mr5">{{classification['bk_classification_name']}}</span>
                        <span class="number">({{classification['bk_objects'].length}})</span>
                    </div>
                    <template v-if="isEditable(classification) && modelType === 'enable'">
                        <cmdb-auth v-if="!mainLoading" class="group-btn ml5" :auth="{ type: $OPERATION.C_MODEL, relation: [classification.id] }">
                            <bk-button slot-scope="{ disabled }"
                                theme="primary"
                                text
                                :disabled="disabled"
                                @click="showModelDialog(classification.bk_classification_id)">
                                <i class="icon-cc-add-line"></i>
                            </bk-button>
                        </cmdb-auth>
                        <cmdb-auth v-if="!mainLoading" class="group-btn" :auth="{ type: $OPERATION.U_MODEL_GROUP, relation: [classification.id] }">
                            <bk-button slot-scope="{ disabled }"
                                theme="primary"
                                text
                                :disabled="disabled"
                                @click="showGroupDialog(true, classification)">
                                <i class="icon-cc-edit"></i>
                            </bk-button>
                        </cmdb-auth>
                        <cmdb-auth v-if="!mainLoading" class="group-btn" :auth="{ type: $OPERATION.D_MODEL_GROUP, relation: [classification.id] }">
                            <bk-button slot-scope="{ disabled }"
                                theme="primary"
                                text
                                :disabled="disabled"
                                @click="deleteGroup(classification)">
                                <i class="icon-cc-delete"></i>
                            </bk-button>
                        </cmdb-auth>
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
                                v-validate="'required|classifyId|length:128'">
                            </bk-input>
                            <p class="form-error" :title="errors.first('classifyId')">{{errors.first('classifyId')}}</p>
                        </div>
                        <i class="bk-icon icon-info-circle" v-bk-tooltips="$t('请填写英文开头，下划线，数字，英文的组合')"></i>
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
                                v-validate="'required|length:128'">
                            </bk-input>
                            <p class="form-error" :title="errors.first('classifyName')">{{errors.first('classifyName')}}</p>
                        </div>
                    </label>
                </div>
            </div>
            <div slot="footer" class="footer">
                <bk-button theme="primary"
                    :loading="$loading(['updateClassification', 'createClassification'])"
                    @click="saveGroup">
                    {{groupDialog.isEdit ? $t('保存') : $t('提交')}}
                </bk-button>
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
                    <bk-button theme="primary" @click="modelClick(curCreateModel)">{{$t('配置字段')}}</bk-button>
                    <bk-button @click="sucessDialog.isShow = false">{{$t('返回列表')}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import cmdbMainInject from '@/components/layout/main-inject'
    import theCreateModel from '@/components/model-manage/_create-model'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    import { addMainScrollListener, removeMainScrollListener } from '@/utils/main-scroller'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
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
            cmdbMainInject
        },
        data () {
            return {
                scrollHandler: null,
                scrollTop: 0,
                topPadding: 0,
                modelType: 'enable',
                searchModel: '',
                filterClassifications: [],
                modelStatisticsSet: {},
                curCreateModel: {},
                sucessDialog: {
                    isShow: false
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
                },
                request: {
                    statistics: Symbol('statistics'),
                    searchClassifications: Symbol('searchClassifications')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName']),
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            enableClassifications () {
                const enableClassifications = []
                this.classifications.forEach(classification => {
                    enableClassifications.push({
                        ...classification,
                        'bk_objects': classification['bk_objects'].filter(model => {
                            return !model.bk_ispaused && !model.bk_ishidden
                        })
                    })
                })
                return enableClassifications
            },
            disabledClassifications () {
                const disabledClassifications = []
                this.classifications.forEach(classification => {
                    const disabledModels = classification['bk_objects'].filter(model => {
                        return model.bk_ispaused && !model.bk_ishidden
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
            },
            mainLoading () {
                return this.$loading(Object.values(this.request))
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
                const lowerCaseValue = value.toLowerCase()
                for (let i = 0; i < classifications.length; i++) {
                    classifications[i].bk_objects = classifications[i].bk_objects.filter(model => {
                        const modelName = model.bk_obj_name.toLowerCase()
                        const modelId = model.bk_obj_id.toLowerCase()
                        return (modelName && modelName.indexOf(lowerCaseValue) !== -1) || (modelId && modelId.indexOf(lowerCaseValue) !== -1)
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
            try {
                await Promise.all([
                    this.getModelStatistics(),
                    this.searchClassificationsObjects({
                        params: {},
                        config: {
                            requestId: this.request.searchClassifications
                        }
                    })
                ])
            } catch (e) {
                this.$route.meta.view = 'error'
            }
            if (this.$route.query.searchModel) {
                const hash = window.location.hash
                this.searchModel = this.$route.query.searchModel
                window.location.hash = hash.substring(0, hash.indexOf('?'))
            }
        },
        mounted () {
            addResizeListener(this.$refs.mainInject.$el, this.handleSetPadding)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.mainInject.$el, this.handleSetPadding)
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
            handleSetPadding () {
                this.topPadding = this.$refs.mainInject.$el.offsetHeight
            },
            isEditable (classification) {
                return !['bk_biz_topo', 'bk_host_manage', 'bk_organization'].includes(classification.bk_classification_id)
            },
            showGroupDialog (isEdit, group) {
                if (isEdit) {
                    this.groupDialog.data.id = group.id
                    this.groupDialog.title = this.$t('编辑分组')
                    this.groupDialog.data.bk_classification_id = group['bk_classification_id']
                    this.groupDialog.data.bk_classification_name = group['bk_classification_name']
                    this.groupDialog.data.id = group.id
                } else {
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
                        requestId: this.request.statistics
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
                const params = {
                    bk_supplier_account: this.supplierAccount,
                    bk_classification_id: this.groupDialog.data['bk_classification_id'],
                    bk_classification_name: this.groupDialog.data['bk_classification_name']
                }
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
                this.$bkInfo({
                    title: this.$t('确认要删除此分组'),
                    confirmFn: async () => {
                        await this.deleteClassification({
                            id: group.id
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
                const params = {
                    bk_supplier_account: this.supplierAccount,
                    bk_obj_name: data['bk_obj_name'],
                    bk_obj_icon: data['bk_obj_icon'],
                    bk_classification_id: data['bk_classification_id'],
                    bk_obj_id: data['bk_obj_id'],
                    userName: this.userName
                }
                const createModel = await this.createObject({ params, config: { requestId: 'createModel' } })
                this.curCreateModel = createModel
                this.sucessDialog.isShow = true
                this.$http.cancel('post_searchClassificationsObjects')
                this.getModelStatistics()
                this.searchClassificationsObjects({
                    params: {}
                })
                this.modelDialog.isShow = false
                this.modelDialog.groupId = ''
                this.searchModel = ''
            },
            modelClick (model) {
                this.$store.commit('objectModel/setActiveModel', model)
                this.$routerActions.redirect({
                    name: 'modelDetails',
                    params: {
                        modelId: model['bk_obj_id']
                    },
                    history: true
                })
            },
            handleGoInstance (model) {
                this.sucessDialog.isShow = false
                const map = {
                    host: MENU_RESOURCE_HOST,
                    biz: MENU_RESOURCE_BUSINESS
                }
                if (map.hasOwnProperty(model.bk_obj_id)) {
                    const query = model.bk_obj_id === 'host' ? { scope: 'all' } : {}
                    this.$routerActions.redirect({
                        name: map[model.bk_obj_id],
                        query
                    })
                } else {
                    this.$routerActions.redirect({
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
        top: 53px;
        left: 0;
        width: calc(100% - 17px);
        padding: 15px 20px 20px;
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
            .title-info {
                @include inlineBlock;
                font-size: 0;
                > span {
                    @include inlineBlock;
                    font-size: 14px;
                    font-weight: 700;
                }
            }
            .number {
                color: $cmdbBorderColor;
            }
            .group-btn {
                display: none;
                vertical-align: middle;
                margin-right: 4px;
                .bk-button-text {
                    font-size: 16px;
                }
            }
            &:hover {
                .group-btn {
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
                color: #3a84ff;
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
            font-size: 50px;
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
