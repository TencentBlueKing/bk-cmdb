<template>
    <div class="model-detail-wrapper">
        <div class="model-info" v-bkloading="{ isLoading: $loading('searchObjects') }">
            <cmdb-auth style="display: none;" :auth="$authResources({ resource_id: modelId, type: $OPERATION.U_MODEL })"
                @update-auth="handleReceiveAuth">
            </cmdb-auth>
            <template v-if="activeModel !== null">
                <div class="choose-icon-wrapper">
                    <span class="model-type">{{getModelType()}}</span>
                    <template v-if="isEditable">
                        <div class="icon-box" v-if="!activeModel['bk_ispaused']" @click="isIconListShow = true">
                            <i class="icon" :class="[activeModel['bk_obj_icon'] || 'icon-cc-default', { ispre: isPublicModel }]"></i>
                            <p class="hover-text">{{$t('点击切换')}}</p>
                        </div>
                        <div class="icon-box" v-else>
                            <i class="icon" :class="[activeModel['bk_obj_icon'] || 'icon-cc-default', { ispre: isPublicModel }]"></i>
                            <p class="hover-text is-paused">{{$t('已停用')}}</p>
                        </div>
                        <div class="choose-icon-box" v-if="isIconListShow" v-click-outside="hideChooseBox">
                            <the-choose-icon
                                v-model="modelInfo.objIcon"
                                @close="hideChooseBox"
                                @chooseIcon="chooseIcon">
                            </the-choose-icon>
                        </div>
                    </template>
                    <template v-else>
                        <div class="icon-box" style="cursor: default;">
                            <i class="icon" :class="[activeModel['bk_obj_icon'] || 'icon-cc-default', { ispre: isPublicModel }]"></i>
                        </div>
                    </template>
                </div>
                <div class="model-text">
                    <span>{{$t('唯一标识')}}：</span>
                    <span class="text-content id" :title="activeModel['bk_obj_id'] || ''">{{activeModel['bk_obj_id'] || ''}}</span>
                </div>
                <div class="model-text">
                    <span>{{$t('名称')}}：</span>
                    <template v-if="!isEditName">
                        <span class="text-content" :title="activeModel['bk_obj_name'] || ''">{{activeModel['bk_obj_name'] || ''}}</span>
                        <i class="icon icon-cc-edit text-primary"
                            v-if="isEditable && !activeModel.ispre && !activeModel.bk_ispaused"
                            @click="editModelName">
                        </i>
                    </template>
                    <template v-else>
                        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('modelName') }">
                            <bk-input type="text" class="cmdb-form-input"
                                name="modelName"
                                v-validate="'required|singlechar|length:256'"
                                v-model.trim="modelInfo.objName">
                            </bk-input>
                        </div>
                        <span class="text-primary" @click="saveModel">{{$t('保存')}}</span>
                        <span class="text-primary" @click="isEditName = false">{{$t('取消')}}</span>
                    </template>
                </div>
                <div class="model-text ml10" v-if="!activeModel['bk_ispaused'] && activeModel.bk_classification_id !== 'bk_biz_topo'">
                    <span>{{$t('实例数量')}}：</span>
                    <div class="text-content-count"
                        :title="modelStatisticsSet[activeModel['bk_obj_id']] || 0"
                        @click="handleGoInstance">
                        <span>{{modelStatisticsSet[activeModel['bk_obj_id']] || 0}}</span>
                        <i class="icon-cc-share"></i>
                    </div>
                </div>
                <cmdb-auth class="restart-btn"
                    v-if="!isMainLine && activeModel['bk_ispaused']"
                    v-bk-tooltips.right="$t('保留模型和相应实例，隐藏关联关系')"
                    :auth="$authResources({ resource_id: modelId, type: $OPERATION.U_MODEL })">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        :disabled="disabled"
                        @click="dialogConfirm('restart')">
                        {{$t('立即启用')}}
                    </bk-button>
                </cmdb-auth>
                <div class="btn-group">
                    <template v-if="canBeImport">
                        <label class="label-btn"
                            v-if="tab.active === 'field'"
                            :class="{ 'disabled': isReadOnly }">
                            <i class="icon-cc-import"></i>
                            <span>{{$t('导入')}}</span>
                            <input v-if="!isReadOnly && updateAuth" ref="fileInput" type="file" @change.prevent="handleFile">
                        </label>
                        <label class="label-btn" @click="exportField">
                            <i class="icon-cc-derivation"></i>
                            <span>{{$t('导出')}}</span>
                        </label>
                    </template>
                    <template v-if="isShowOperationButton">
                        <cmdb-auth class="label-btn"
                            v-if="!isMainLine && !activeModel['bk_ispaused']"
                            v-bk-tooltips="$t('保留模型和相应实例，隐藏关联关系')"
                            :auth="$authResources({ resource_id: modelId, type: $OPERATION.U_MODEL })">
                            <bk-button slot-scope="{ disabled }"
                                text
                                :disabled="disabled"
                                @click="dialogConfirm('stop')">
                                <i class="bk-icon icon-minus-circle-shape"></i>
                                <span>{{$t('停用')}}</span>
                            </bk-button>
                        </cmdb-auth>
                        <cmdb-auth class="label-btn"
                            v-bk-tooltips="$t('删除模型和其下所有实例，此动作不可逆，请谨慎操作')"
                            :auth="$authResources({ resource_id: modelId, type: $OPERATION.D_MODEL })">
                            <bk-button slot-scope="{ disabled }"
                                text
                                :disabled="disabled"
                                @click="dialogConfirm('delete')">
                                <i class="icon-cc-del"></i>
                                <span>{{$t('删除')}}</span>
                            </bk-button>
                        </cmdb-auth>
                    </template>
                </div>
            </template>
        </div>
        <bk-tab class="model-details-tab" type="unborder-card" :active.sync="tab.active">
            <bk-tab-panel name="field" :label="$t('模型字段')">
                <the-field-group ref="field" v-if="tab.active === 'field'"></the-field-group>
            </bk-tab-panel>
            <bk-tab-panel name="relation" :label="$t('模型关联')" :visible="activeModel && !specialModel.includes(activeModel['bk_obj_id'])">
                <the-relation v-if="tab.active === 'relation'" :model-id="modelId"></the-relation>
            </bk-tab-panel>
            <bk-tab-panel name="verification" :label="$t('唯一校验')">
                <the-verification v-if="tab.active === 'verification'" :model-id="modelId"></the-verification>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import theRelation from './relation'
    import theVerification from './verification'
    import theFieldGroup from '@/components/model-manage/field-group'
    import theChooseIcon from '@/components/model-manage/choose-icon/_choose-icon'
    import { mapActions, mapGetters, mapMutations } from 'vuex'
    import {
        MENU_MODEL_MANAGEMENT,
        MENU_RESOURCE_HOST,
        MENU_RESOURCE_BUSINESS,
        MENU_RESOURCE_INSTANCE,
        MENU_MODEL_BUSINESS_TOPOLOGY
    } from '@/dictionary/menu-symbol'
    export default {
        components: {
            theFieldGroup,
            theRelation,
            theVerification,
            theChooseIcon
        },
        data () {
            return {
                tab: {
                    active: 'field'
                },
                modelInfo: {
                    objName: '',
                    objIcon: ''
                },
                isIconListShow: false,
                isEditName: false,
                specialModel: ['process', 'plat'],
                modelStatisticsSet: {},
                updateAuth: false
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount',
                'userName',
                'admin',
                'isAdminView',
                'isBusinessSelected'
            ]),
            ...mapGetters('objectModel', [
                'activeModel',
                'isPublicModel',
                'isInjectable',
                'isMainLine'
            ]),
            ...mapGetters('objectModelClassify', ['models']),
            isShowOperationButton () {
                return (this.isAdminView || !this.isPublicModel)
                    && !this.activeModel['ispre']
            },
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            isEditable () {
                if (!this.updateAuth) {
                    return false
                }
                if (this.isAdminView) {
                    return !this.activeModel.ispre
                }
                return !this.isReadOnly && !this.isPublicModel
            },
            modelParams () {
                const {
                    objIcon,
                    objName
                } = this.modelInfo
                const params = {
                    modifier: this.userName
                }
                if (objIcon) {
                    Object.assign(params, { bk_obj_icon: objIcon })
                }
                if (objName.length && objName !== this.activeModel['bk_obj_name']) {
                    Object.assign(params, { bk_obj_name: objName })
                }
                return params
            },
            exportUrl () {
                return `${window.API_HOST}object/owner/${this.supplierAccount}/object/${this.activeModel['bk_obj_id']}/export`
            },
            canBeImport () {
                const cantImport = ['host', 'biz', 'process', 'plat']
                return this.updateAuth
                    && !this.isMainLine
                    && !cantImport.includes(this.$route.params.modelId)
            },
            modelId () {
                const model = this.$store.getters['objectModelClassify/getModelById'](this.$route.params.modelId)
                return model.id || null
            }
        },
        watch: {
            '$route.params.modelId' () {
                this.initObject()
            }
        },
        created () {
            this.initObject()
        },
        methods: {
            ...mapActions('objectModel', [
                'searchObjects',
                'updateObject',
                'deleteObject'
            ]),
            ...mapActions('objectBatch', [
                'importObjectAttribute',
                'exportObjectAttribute'
            ]),
            ...mapActions('objectMainLineModule', [
                'deleteMainlineObject'
            ]),
            ...mapMutations('objectModel', [
                'setActiveModel'
            ]),
            getModelType () {
                if (this.activeModel.ispre) {
                    return this.$t('内置')
                } else {
                    if (this.$tools.getMetadataBiz(this.activeModel)) {
                        return this.$t('自定义')
                    }
                    return this.$t('公共')
                }
            },
            async handleFile (e) {
                const files = e.target.files
                const formData = new FormData()
                formData.append('file', files[0])
                if (!this.isPublicModel) {
                    formData.append('metadata', JSON.stringify(this.$injectMetadata().metadata))
                }
                try {
                    const res = await this.importObjectAttribute({
                        params: formData,
                        objId: this.activeModel['bk_obj_id'],
                        config: {
                            requestId: 'importObjectAttribute',
                            globalError: false,
                            transformData: false
                        }
                    }).then(res => {
                        this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                        return res
                    })
                    if (res.result) {
                        const data = res.data[this.activeModel['bk_obj_id']]
                        if (data.hasOwnProperty('insert_failed')) {
                            this.$error(data['insert_failed'][0])
                        } else if (data.hasOwnProperty('update_failed')) {
                            this.$error(data['update_failed'][0])
                        } else {
                            this.$success(this.$t('导入成功'))
                            this.$refs.field.initFieldList()
                        }
                    } else {
                        this.$error(res['bk_error_msg'])
                    }
                } catch (e) {
                    this.$error(e.data['bk_error_msg'])
                } finally {
                    this.$refs.fileInput.value = ''
                }
            },
            checkModel () {
                return this.models.find(model => model['bk_obj_id'] === this.$route.params.modelId)
            },
            hideChooseBox () {
                this.isIconListShow = false
            },
            chooseIcon () {
                this.isIconListShow = false
                this.saveModel()
            },
            editModelName () {
                this.modelInfo.objName = this.activeModel['bk_obj_name']
                this.isEditName = true
            },
            async saveModel () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                await this.updateObject({
                    id: this.activeModel['id'],
                    params: this.$injectMetadata(this.modelParams, { clone: true })
                }).then(() => {
                    this.$http.cancel('post_searchClassificationsObjects')
                })
                this.setActiveModel({ ...this.activeModel, ...this.modelParams })
                this.isEditName = false
            },
            async initObject () {
                await this.getModelStatistics()
                const model = this.$store.getters['objectModelClassify/getModelById'](this.$route.params.modelId)
                if (model) {
                    this.$store.commit('objectModel/setActiveModel', model)
                    this.setBreadcrumbs(model)
                    this.initModelInfo()
                } else {
                    this.$router.replace({ name: 'status404' })
                }
            },
            setBreadcrumbs (model) {
                const breadcrumbs = [{
                    label: this.$t('模型管理'),
                    route: {
                        name: MENU_MODEL_MANAGEMENT
                    }
                }, {
                    label: model.bk_obj_name
                }]
                if (this.$route.query.from === 'business') {
                    breadcrumbs.splice(0, 1, {
                        label: this.$t('业务层级'),
                        route: {
                            name: MENU_MODEL_BUSINESS_TOPOLOGY
                        }
                    })
                }
                this.$store.commit('setBreadcrumbs', breadcrumbs)
            },
            async getModelStatistics () {
                const modelStatisticsSet = {}
                const res = await this.$store.dispatch('objectModelClassify/getClassificationsObjectStatistics', {
                    config: {
                        requestId: 'getClassificationsObjectStatistics'
                    }
                })
                res.forEach(item => {
                    modelStatisticsSet[item.bk_obj_id] = item.instance_count
                })
                this.modelStatisticsSet = modelStatisticsSet
            },
            initModelInfo () {
                this.modelInfo = {
                    objIcon: this.activeModel['bk_obj_icon'],
                    objName: this.activeModel['bk_obj_name']
                }
            },
            exportExcel (response) {
                const contentDisposition = response.headers['content-disposition']
                const fileName = contentDisposition.substring(contentDisposition.indexOf('filename') + 9)
                const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }))
                const link = document.createElement('a')
                link.style.display = 'none'
                link.href = url
                link.setAttribute('download', fileName)
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
            },
            async exportField () {
                const res = await this.exportObjectAttribute({
                    objId: this.activeModel['bk_obj_id'],
                    params: this.$injectMetadata({}, { inject: !this.isPublicModel }),
                    config: {
                        globalError: false,
                        originalResponse: true,
                        responseType: 'blob'
                    }
                })
                this.exportExcel(res)
            },
            dialogConfirm (type) {
                switch (type) {
                    case 'restart':
                        this.$bkInfo({
                            title: this.$t('确认要启用该模型？'),
                            confirmFn: () => {
                                this.updateModelObject(false)
                            }
                        })
                        break
                    case 'stop':
                        this.$bkInfo({
                            title: this.$t('确认要停用该模型？'),
                            confirmFn: () => {
                                this.updateModelObject(true)
                            }
                        })
                        break
                    case 'delete':
                        this.$bkInfo({
                            title: this.$t('确认要删除该模型？'),
                            confirmFn: () => {
                                this.deleteModel()
                            }
                        })
                        break
                    default:
                }
            },
            async updateModelObject (ispaused) {
                await this.updateObject({
                    id: this.activeModel['id'],
                    params: this.$injectMetadata({
                        bk_ispaused: ispaused
                    }),
                    config: {
                        requestId: 'updateModel'
                    }
                })
                this.$store.commit('objectModelClassify/updateModel', {
                    bk_ispaused: ispaused,
                    bk_obj_id: this.activeModel.bk_obj_id
                })
                this.setActiveModel({ ...this.activeModel, ...{ bk_ispaused: ispaused } })
            },
            async deleteModel () {
                if (this.isMainLine) {
                    await this.deleteMainlineObject({
                        bkObjId: this.activeModel['bk_obj_id'],
                        config: {
                            requestId: 'deleteModel'
                        }
                    })
                    const routerName = this.$route.query.from === 'business' ? MENU_MODEL_BUSINESS_TOPOLOGY : MENU_MODEL_MANAGEMENT
                    this.$router.replace({ name: routerName })
                } else {
                    await this.deleteObject({
                        id: this.activeModel['id'],
                        config: {
                            data: this.$injectMetadata({}, {
                                inject: this.isInjectable
                            }),
                            requestId: 'deleteModel'
                        }
                    })
                    this.$router.replace({ name: MENU_MODEL_MANAGEMENT })
                }
                this.$http.cancel('post_searchClassificationsObjects')
            },
            handleGoInstance () {
                const model = this.activeModel
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
            },
            handleReceiveAuth (auth) {
                this.updateAuth = auth
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-detail-wrapper {
        padding: 0;
    }
    .model-details-tab {
        height: calc(100% - 70px) !important;
        /deep/ {
            .bk-tab-header {
                padding: 0;
                margin: 0 20px;
            }
            .bk-tab-section {
                padding: 0;
            }
        }
    }
    .model-info {
        padding: 0 24px;
        height: 70px;
        background: #ebf4ff;
        font-size: 14px;
        border-bottom: 1px solid #dcdee5;
        border-top: 1px solid #dcdee5;
        .choose-icon-wrapper {
            position: relative;
            float: left;
            margin: 8px 30px 0 0;
            .model-type {
                position: absolute;
                left: 42px;
                top: -12px;
                padding: 0 6px;
                border-radius: 4px;
                background-color: #ffb23a;
                font-size: 20px;
                line-height: 32px;
                color: #fff;
                white-space: nowrap;
                transform: scale(.5);
                transform-origin: left center;
                z-index: 2;
                &:after {
                    content: "";
                    position: absolute;
                    top: 100%;
                    left: 10px;
                    width: 0;
                    height: 0;
                    border-top: 8px solid #ffb23a;
                    border-right: 10px solid transparent;
                    transform: skew(-15deg);
                    transform-origin: left top;
                }
            }
            .choose-icon-box {
                position: absolute;
                left: -12px;
                top: 62px;
                width: 600px;
                height: 460px;
                background: #fff;
                border: 1px solid #dde4e8;
                box-shadow: 0px 3px 6px 0px rgba(51, 60, 72, 0.13);
                z-index: 99;
                &:before {
                    position: absolute;
                    top: -13px;
                    left: 30px;
                    content: '';
                    border: 6px solid transparent;
                    border-bottom-color: rgba(51, 60, 72, 0.23);
                }
                &:after {
                    position: absolute;
                    top: -12px;
                    left: 30px;
                    content: '';
                    border: 6px solid transparent;
                    border-bottom-color: #fff;
                }
            }
        }
        .icon-box {
            padding-top: 16px;
            width: 54px;
            height: 54px;
            border: 1px solid #dde4eb;
            border-radius: 50%;
            background: #fff;
            text-align: center;
            font-size: 20px;
            color: $cmdbBorderFocusColor;
            cursor: pointer;
            &:hover {
                .hover-text {
                    background: rgba(0, 0, 0, .5);
                    display: block;
                }
            }
            .hover-text {
                display: none;
                position: absolute;
                top: 0;
                left: 0;
                width: 54px;
                height: 54px;
                line-height: 54px;
                font-size: 12px;
                border-radius: 50%;
                text-align: center;
                color: #fff;
                &.is-paused {
                    background: rgba(0, 0, 0, .5);
                    display: block !important;
                }
            }
            .icon {
                vertical-align: top;
                &.ispre {
                    color: #3a84ff;
                }
            }
        }
        .model-text {
            float: left;
            margin: 18px 10px 0 0;
            line-height: 36px;
            font-size: 0;
            >span {
                display: inline-block;
                vertical-align: middle;
                height: 36px;
                font-size: 14px;
                color: #737987;
            }
            .text-content {
                max-width: 110px;
                vertical-align: middle;
                color: #333948;
                @include ellipsis;
                &.id {
                    min-width: 50px;
                }
            }
            .text-content-count {
                display: inline-block;
                vertical-align: middle;
                color: #3a84ff;
                cursor: pointer;
                >span {
                    font-size: 14px;
                    vertical-align: middle;
                }
                .icon-cc-share {
                    font-size: 12px;
                    margin-left: 6px;
                    vertical-align: middle;
                }
            }
            .icon-cc-edit {
                vertical-align: middle;
                font-size: 14px;
            }
            .cmdb-form-item {
                display: inline-block;
                width: 156px;
                vertical-align: top;
                input {
                    vertical-align: top;
                }
            }
            .text-primary {
                cursor: pointer;
                margin-left: 5px;
            }
        }
        .restart-btn {
            display: inline-block;
            margin: 19px 0 0 20px;
        }
        .btn-group {
            float: right;
            height: 70px;
            line-height: 70px;
            display: flex;
            align-items: center;
            .label-btn {
                line-height: normal;
                outline: 0;
                position: relative;
                .bk-button-text {
                    color: #737987;
                    &:disabled {
                        color: #dcdee5 !important;
                        cursor: not-allowed;
                    }
                }
                &.disabled {
                    cursor: not-allowed;
                }
                input[type="file"] {
                    position: absolute;
                    left: 0;
                    top: 0;
                    opacity: 0;
                    width: 100%;
                    height: 100%;
                    cursor: pointer;
                }
                ::-webkit-file-upload-button {
                    cursor:pointer;
                }
            }
            .export-form {
                display: inline-block;
            }
            .label-btn {
                margin-left: 10px;
                cursor: pointer;
                &:hover {
                    color: $cmdbBorderFocusColor;
                    .bk-button-text {
                        color: $cmdbBorderFocusColor;
                    }
                }
            }
            i,
            span {
                vertical-align: middle;
            }
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
