<template>
    <div class="model-detail-wrapper">
        <div class="model-info" v-bkloading="{ isLoading: $loading('searchObjects') }">
            <template v-if="activeModel !== null">
                <div class="choose-icon-wrapper">
                    <span class="model-type">{{getModelType()}}</span>
                    <template v-if="isEditable">
                        <div class="icon-box" @click="isIconListShow = true">
                            <i class="icon" :class="[activeModel ? activeModel['bk_obj_icon'] : 'icon-cc-default', { ispre: isPublicModel }]"></i>
                            <p class="hover-text">{{$t('点击切换')}}</p>
                        </div>
                        <div class="choose-icon-box" v-if="isIconListShow" v-click-outside="hideChooseBox">
                            <the-choose-icon
                                v-model="modelInfo.objIcon"
                                type="update"
                                @chooseIcon="chooseIcon">
                            </the-choose-icon>
                        </div>
                    </template>
                    <template v-else>
                        <div class="icon-box" style="cursor: default;">
                            <i class="icon" :class="[activeModel ? activeModel['bk_obj_icon'] : 'icon-cc-default', { ispre: isPublicModel }]"></i>
                        </div>
                    </template>
                </div>
                <div class="model-text">
                    <span>{{$t('唯一标识')}}：</span>
                    <span class="text-content id">{{activeModel ? activeModel['bk_obj_id'] : ''}}</span>
                </div>
                <div class="model-text">
                    <span>{{$t('名称')}}：</span>
                    <template v-if="!isEditName">
                        <span class="text-content">{{activeModel ? activeModel['bk_obj_name'] : ''}}
                        </span>
                        <i class="icon icon-cc-edit text-primary"
                            v-if="isEditable && !activeModel.ispre"
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
                <div class="btn-group">
                    <template v-if="canBeImport">
                        <label class="label-btn"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.U_MODEL),
                                auth: [$OPERATION.U_MODEL]
                            }"
                            v-if="tab.active === 'field'"
                            :class="{ 'disabled': isReadOnly }">
                            <i class="icon-cc-import"></i>
                            <span>{{$t('导入')}}</span>
                            <input v-if="!isReadOnly && $isAuthorized($OPERATION.U_MODEL)" ref="fileInput" type="file" @change.prevent="handleFile">
                        </label>
                        <label class="label-btn" @click="exportField">
                            <i class="icon-cc-derivation"></i>
                            <span>{{$t('导出')}}</span>
                        </label>
                    </template>
                    <template v-if="isShowOperationButton">
                        <label class="label-btn"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.U_MODEL),
                                auth: [$OPERATION.U_MODEL]
                            }"
                            v-if="!isMainLine"
                            v-bk-tooltips="$t('保留模型和相应实例，隐藏关联关系')">
                            <i class="bk-icon icon-minus-circle-shape"></i>
                            <span v-if="activeModel['bk_ispaused']" @click="dialogConfirm('restart')">
                                {{$t('启用')}}
                            </span>
                            <span v-else @click="dialogConfirm('stop')">
                                {{$t('停用')}}
                            </span>
                        </label>
                        <label class="label-btn"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION.D_MODEL),
                                auth: [$OPERATION.D_MODEL]
                            }"
                            v-bk-tooltips="$t('删除模型和其下所有实例，此动作不可逆，请谨慎操作')"
                            @click="dialogConfirm('delete')">
                            <i class="icon-cc-del"></i>
                            <span>{{$t('删除')}}</span>
                        </label>
                    </template>
                </div>
            </template>
        </div>
        <bk-tab class="model-details-tab" type="unborder-card" :active.sync="tab.active">
            <bk-tab-panel name="field" :label="$t('模型字段')">
                <the-field ref="field" v-if="tab.active === 'field'"></the-field>
            </bk-tab-panel>
            <bk-tab-panel name="relation" :label="$t('模型关联')" :visible="activeModel && !specialModel.includes(activeModel['bk_obj_id'])">
                <the-relation v-if="tab.active === 'relation'"></the-relation>
            </bk-tab-panel>
            <bk-tab-panel name="verification" :label="$t('唯一校验')">
                <the-verification v-if="tab.active === 'verification'"></the-verification>
            </bk-tab-panel>
            <bk-tab-panel name="propertyGroup" :label="$t('字段分组')">
                <the-property-group v-if="tab.active === 'propertyGroup'"></the-property-group>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import thePropertyGroup from './group.vue'
    import theField from './field'
    import theRelation from './relation'
    import theChooseIcon from '@/components/model-manage/_choose-icon'
    import theVerification from './verification'
    import { mapActions, mapGetters, mapMutations } from 'vuex'
    export default {
        components: {
            thePropertyGroup,
            theField,
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
                specialModel: ['process', 'plat']
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
                const updateAuth = this.$isAuthorized(this.$OPERATION.U_MODEL)
                if (!updateAuth) {
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
                return this.$isAuthorized(this.$OPERATION.U_MODEL)
                    && !this.isMainLine
                    && !cantImport.includes(this.$route.params.modelId)
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
                const model = this.$store.getters['objectModelClassify/getModelById'](this.$route.params.modelId)
                if (model) {
                    this.$store.commit('objectModel/setActiveModel', model)
                    this.$store.commit('setHeaderTitle', model['bk_obj_name'])
                    this.initModelInfo()
                } else {
                    this.$router.replace({ name: 'status404' })
                }
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
                if (type === 'delete') {
                    if (!this.$isAuthorized(this.$OPERATION.D_MODEL)) {
                        return false
                    }
                } else if (!this.$isAuthorized(this.$OPERATION.U_MODEL)) {
                    return false
                }
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
                }
                this.$http.cancel('post_searchClassificationsObjects')
                this.$router.replace({ name: 'model' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-detail-wrapper {
        padding: 0;
        height: 100%;
    }
    .model-details-tab {
        height: calc(100% - 100px) !important;
    }
    .model-info {
        padding: 0 24px 0 38px;
        height: 100px;
        background: rgba(235, 244, 255, .6);
        font-size: 14px;
        .choose-icon-wrapper {
            position: relative;
            float: left;
            margin: 14px 30px 0 0;
            .model-type {
                position: absolute;
                left: 58px;
                top: -8px;
                padding: 0 6px;
                border-radius: 4px;
                background-color: #ffb23a;
                font-size: 20px;
                line-height: 32px;
                color: #fff;
                white-space: nowrap;
                transform: scale(.5);
                transform-origin: left center;
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
                left: 0;
                top: 80px;
                width: 395px;
                height: 262px;
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
            padding-top: 20px;
            width: 72px;
            height: 72px;
            border: 1px solid #dde4eb;
            border-radius: 50%;
            background: #fff;
            text-align: center;
            font-size: 32px;
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
                height: 72px;
                width: 72px;
                font-size: 12px;
                line-height: 72px;
                border-radius: 50%;
                text-align: center;
                color: #fff;
            }
            .icon {
                vertical-align: top;
                &.ispre {
                    color: #868b97;
                }
            }
        }
        .model-text {
            float: left;
            margin: 32px 10px 32px 0;
            line-height: 36px;
            font-size: 0;
            &:hover .icon-cc-edit {
                display: inline-block;
            }
            >span {
                display: inline-block;
                vertical-align: middle;
                height: 36px;
                font-size: 14px;
            }
            .text-content {
                max-width: 200px;
                vertical-align: middle;
                @include ellipsis;
                &.id {
                    width: 110px;
                }
            }
            .icon-cc-edit {
                display: none;
                vertical-align: middle;
                font-size: 14px;
            }
            .cmdb-form-item {
                display: inline-block;
                width: 200px;
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
        .btn-group {
            float: right;
            height: 100px;
            line-height: 100px;
            .label-btn {
                outline: 0;
                position: relative;
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
            }
            .export-form {
                display: inline-block;
            }
            .label-btn {
                margin-left: 10px;
                cursor: pointer;
                &:hover {
                    color: $cmdbBorderFocusColor;
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
