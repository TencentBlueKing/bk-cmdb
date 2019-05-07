<template>
    <div class="model-detail-wrapper">
        <div class="model-info" v-bkloading="{isLoading: $loading('searchObjects')}">
            <template v-if="activeModel !== null">
                <div class="choose-icon-wrapper">
                    <template v-if="authority.includes('update') && !isReadOnly">
                        <div class="icon-box" @click="isIconListShow = true">
                            <i class="icon" :class="[activeModel ? activeModel['bk_obj_icon'] : 'icon-cc-default', {ispre: activeModel['ispre']}]"></i>
                            <p class="hover-text">{{$t('ModelManagement["点击切换"]')}}</p>
                        </div>
                        <div class="choose-icon-box" v-if="isIconListShow" v-click-outside="hideChooseBox">
                            <the-choose-icon
                                :type="'update'"
                                v-model="modelInfo.objIcon"
                                @chooseIcon="chooseIcon"
                            ></the-choose-icon>
                        </div>
                    </template>
                    <template v-else>
                        <div class="icon-box" style="cursor: default;">
                            <i class="icon" :class="activeModel ? activeModel['bk_obj_icon'] : 'icon-cc-default'"></i>
                        </div>
                    </template>
                </div>
                <div class="model-text">
                    <span>{{$t('ModelManagement["唯一标识"]')}}：</span>
                    <span class="text-content id">{{activeModel ? activeModel['bk_obj_id'] : ''}}</span>
                </div>
                <div class="model-text">
                    <span>{{$t('Hosts["名称"]')}}：</span>
                    <template v-if="!isEditName">
                        <span class="text-content">{{activeModel ? activeModel['bk_obj_name'] : ''}}
                        </span>
                        <i class="icon icon-cc-edit text-primary"
                            v-if="!(isReadOnly || (activeModel && activeModel['ispre'])) && authority.includes('update')"
                            @click="editModelName">
                        </i>
                    </template>
                    <template v-else>
                        <div class="cmdb-form-item" :class="{'is-error': errors.has('modelName')}">
                            <input type="text" class="cmdb-form-input"
                            name="modelName"
                            v-validate="'required|singlechar'"
                            v-model.trim="modelInfo.objName">
                        </div>
                        <span class="text-primary" @click="saveModel">{{$t("Common['保存']")}}</span>
                        <span class="text-primary" @click="isEditName = false">{{$t("Common['取消']")}}</span>
                    </template>
                </div>
                <div class="btn-group">
                    <template v-if="canBeImport">
                        <label class="label-btn"
                            v-if="tab.active==='field' && authority.includes('update')"
                            :class="{'disabled': isReadOnly}">
                            <i class="icon-cc-import"></i>
                            <span>{{$t('ModelManagement["导入"]')}}</span>
                            <input v-if="!isReadOnly" ref="fileInput" type="file" @change.prevent="handleFile">
                        </label>
                        <form class="export-form" ref="submitForm" :action="exportUrl" method="POST" v-if="tab.active==='field'">
                            <label class="label-btn" @click="exportField">
                                <i class="icon-cc-derivation"></i>
                                <span>{{$t('ModelManagement["导出"]')}}</span>
                            </label>
                        </form>
                    </template>
                    <template v-if="!activeModel['ispre'] && authority.includes('update')">
                        <label class="label-btn"
                        v-if="!isMainLine"
                        v-tooltip="$t('ModelManagement[\'保留模型和相应实例，隐藏关联关系\']')">
                            <i class="bk-icon icon-minus-circle-shape"></i>
                            <span v-if="activeModel['bk_ispaused']" @click="dialogConfirm('restart')">
                                {{$t('ModelManagement["启用"]')}}
                            </span>
                            <span v-else @click="dialogConfirm('stop')">
                                {{$t('ModelManagement["停用"]')}}
                            </span>
                        </label>
                        <label class="label-btn"
                            v-tooltip="$t('ModelManagement[\'删除模型和其下所有实例，此动作不可逆，请谨慎操作\']')"
                            @click="dialogConfirm('delete')">
                            <i class="icon-cc-del"></i>
                            <span>{{$t("Common['删除']")}}</span>
                        </label>
                    </template>
                </div>
            </template>
        </div>
        <bk-tab class="model-details-tab" :active-name.sync="tab.active">
            <bk-tabpanel name="field" :title="$t('ModelManagement[\'模型字段\']')">
                <the-field ref="field" v-if="tab.active === 'field'"></the-field>
            </bk-tabpanel>
            <bk-tabpanel name="relation" :title="$t('ModelManagement[\'模型关联\']')" :show="activeModel && !specialModel.includes(activeModel['bk_obj_id'])">
                <the-relation v-if="tab.active === 'relation'"></the-relation>
            </bk-tabpanel>
            <bk-tabpanel name="verification" :title="$t('ModelManagement[\'唯一校验\']')">
                <the-verification v-if="tab.active === 'verification'"></the-verification>
            </bk-tabpanel>
            <bk-tabpanel name="propertyGroup" :title="$t('ModelManagement[\'字段分组\']')">
                <the-property-group v-if="tab.active === 'propertyGroup'"></the-property-group>
            </bk-tabpanel>
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
                'admin'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            isMainLine () {
                return this.activeModel['bk_classification_id'] === 'bk_biz_topo'
            },
            modelParams () {
                let {
                    objIcon,
                    objName
                } = this.modelInfo
                let params = {
                    modifier: this.userName
                }
                if (objIcon) {
                    Object.assign(params, {bk_obj_icon: objIcon})
                }
                if (objName.length && objName !== this.activeModel['bk_obj_name']) {
                    Object.assign(params, {bk_obj_name: objName})
                }
                return params
            },
            exportUrl () {
                return `${window.API_HOST}object/owner/${this.supplierAccount}/object/${this.activeModel['bk_obj_id']}/export`
            },
            authority () {
                return this.admin ? ['search', 'update', 'delete'] : []
            },
            canBeImport () {
                const cantImport = ['host', 'biz', 'process', 'plat']
                return this.authority.includes('update') && !this.isMainLine && !cantImport.includes(this.$route.params.modelId)
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
                'importObjectAttribute'
            ]),
            ...mapActions('objectMainLineModule', [
                'deleteMainlineObject'
            ]),
            ...mapMutations('objectModel', [
                'setActiveModel'
            ]),
            async handleFile (e) {
                let files = e.target.files
                let formData = new FormData()
                formData.append('file', files[0])
                try {
                    const res = await this.importObjectAttribute({
                        params: formData,
                        bkObjId: this.activeModel['bk_obj_id'],
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
                        let data = res.data[this.activeModel['bk_obj_id']]
                        if (data.hasOwnProperty('insert_failed')) {
                            this.$error(data['insert_failed'][0])
                        } else if (data.hasOwnProperty('update_failed')) {
                            this.$error(data['update_failed'][0])
                        } else {
                            this.$success(this.$t('ModelManagement["导入成功"]'))
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
                return this.$allModels.find(model => model['bk_obj_id'] === this.$route.params.modelId)
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
                    params: this.modelParams
                }).then(() => {
                    this.$http.cancel('post_searchClassificationsObjects')
                })
                this.setActiveModel({...this.activeModel, ...this.modelParams})
                this.$store.commit('setHeaderTitle', this.activeModel['bk_obj_name'])
                this.isEditName = false
            },
            async initObject () {
                const res = await this.searchObjects({
                    params: {
                        bk_obj_id: this.$route.params.modelId,
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: 'searchObjects'
                    }
                })
                if (res.length) {
                    this.$store.commit('objectModel/setActiveModel', res[0])
                    this.$store.commit('setHeaderTitle', this.activeModel['bk_obj_name'])
                    this.initModelInfo()
                } else {
                    this.$router.replace('/status-404')
                }
            },
            initModelInfo () {
                this.modelInfo = {
                    objIcon: this.activeModel['bk_obj_icon'],
                    objName: this.activeModel['bk_obj_name']
                }
            },
            exportField () {
                this.$refs.submitForm.submit()
            },
            dialogConfirm (type) {
                switch (type) {
                    case 'restart':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要启用该模型？"]'),
                            confirmFn: () => {
                                this.updateModelObject(false)
                            }
                        })
                        break
                    case 'stop':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要停用该模型？"]'),
                            confirmFn: () => {
                                this.updateModelObject(true)
                            }
                        })
                        break
                    case 'delete':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要删除该模型？"]'),
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
                    params: {
                        bk_ispaused: ispaused
                    },
                    config: {
                        requestId: 'updateModel'
                    }
                }).then(() => {
                    this.$http.cancel('post_searchClassificationsObjects')
                })
                this.initObject()
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
                            requestId: 'deleteModel'
                        }
                    })
                }
                this.$http.cancel('post_searchClassificationsObjects')
                this.$router.replace('/model')
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
            .choose-icon-box {
                position: absolute;
                left: 0;
                top: 95px;
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
            position: relative;
            float: left;
            margin-top: 14px;
            margin: 14px 30px 0 0;
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
