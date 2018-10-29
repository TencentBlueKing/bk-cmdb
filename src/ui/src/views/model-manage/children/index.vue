<template>
    <div class="model-detail-wrapper">
        <div class="model-info" v-if="activeModel">
            <div class="choose-icon-wrapper">
                <div class="icon-box" @click="isIconListShow = true">
                    <i class="icon" :class="activeModel ? activeModel['bk_obj_icon'] : 'icon-cc-default'"></i>
                    <p class="hover-text">{{$t('ModelManagement["点击切换"]')}}</p>
                </div>
                <div class="choose-icon-box" v-if="isIconListShow">
                    <the-choose-icon
                        :type="'update'"
                        v-model="modelInfo.objIcon"
                        @chooseIcon="chooseIcon"
                    ></the-choose-icon>
                </div>
            </div>
            <div class="model-text">
                <span>{{$t('ModelManagement["唯一标识"]')}}：</span>
                <span class="text-content name">{{activeModel ? activeModel['bk_obj_id'] : ''}}</span>
            </div>
            <div class="model-text">
                <span>{{$t('Hosts["名称"]')}}：</span>
                <template v-if="!isEditName">
                    <span class="text-content">{{activeModel ? activeModel['bk_obj_name'] : ''}}<i class="icon icon-cc-edit text-primary" v-if="!(isReadOnly || (activeModel && activeModel['ispre']))" @click="editModelName"></i></span>
                </template>
                <template v-else>
                    <input type="text" class="cmdb-form-input" v-model.trim="modelInfo.objName">
                    <span class="text-primary" @click="saveModel">{{$t("Common['保存']")}}</span>
                    <span class="text-primary" @click="isEditName = false">{{$t("Common['取消']")}}</span>
                </template>
            </div>
            <div class="btn-group">
                <form class="export-form" ref="submitForm" :action="exportUrl" method="POST">
                    <label class="label-btn" @click="exportField">
                        <i class="icon-cc-derivation"></i>
                        <span>{{$t('ModelManagement["导出"]')}}</span>
                    </label>
                </form>
                <template v-if="!activeModel['ispre']">
                    <label class="label-btn">
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
        </div>
        <div class="model-info" v-bkloading="{isLoading: $loading('searchObjects')}" v-else>
        </div>
        <bk-tab class="model-details-tab" :active-name.sync="tab.active">
            <bk-tabpanel name="field" :title="$t('ModelManagement[\'模型字段\']')">
                <the-field></the-field>
            </bk-tabpanel>
            <bk-tabpanel name="relation" :title="$t('ModelManagement[\'模型关系\']')">
                <the-relation></the-relation>
            </bk-tabpanel>
            <bk-tabpanel name="propertyGroup" :title="$t('ModelManagement[\'字段分组\']')">
                <the-property-group v-if="tab.active === 'propertyGroup'"></the-property-group>
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import thePropertyGroup from './_property-group.vue'
    import theField from './field'
    import theRelation from './relation'
    import theChooseIcon from '../_choose-icon'
    import { mapActions, mapGetters, mapMutations } from 'vuex'
    export default {
        components: {
            thePropertyGroup,
            theField,
            theRelation,
            theChooseIcon
        },
        props: {
            isReadOnly: {
                type: Boolean,
                default: false
            }
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
                isEditName: false
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount',
                'userName'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
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
            ...mapMutations('objectModel', [
                'setActiveModel'
            ]),
            chooseIcon () {
                this.isIconListShow = false
                this.saveModel()
            },
            editModelName () {
                this.modelInfo.objName = this.activeModel['bk_obj_name']
                this.isEditName = true
            },
            async saveModel () {
                await this.updateObject({
                    id: this.activeModel['id'],
                    params: this.modelParams
                }).then(() => {
                    this.$http.cancel('post_searchClassificationsObjects')
                })
                this.setActiveModel({...this.activeModel, ...this.modelParams})
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
                this.$store.commit('objectModel/setActiveModel', res[0])
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
                this.$route.push('/model-manage')
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
                this.$route.push('/model-manage')
                this.$http.cancel('post_searchClassificationsObjects')
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
                top: 100px;
                width: 395px;
                height: 262px;
                background: #fff;
                border: 1px solid #dde4e8;
                box-shadow: 0px 3px 6px 0px rgba(51, 60, 72, 0.13);
                z-index: 99;
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
            }
        }
        .model-text {
            float: left;
            margin: 32px 10px 32px 0;
            line-height: 36px;
            >span {
                display: inline-block;
                vertical-align: top;
            }
            .text-content {
                max-width: 110px;
                @include ellipsis;
                &.name {
                    width: 110px;
                }
                .icon {
                    margin-top: -4px;
                    margin: -4px 0 0 4px;
                }
            }
            .cmdb-form-input {
                display: inline-block;
                width: 200px;
                vertical-align: top;
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
            .export-form {
                display: inline-block;
            }
            .label-btn {
                margin-left: 10px;
                cursor: pointer;
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
