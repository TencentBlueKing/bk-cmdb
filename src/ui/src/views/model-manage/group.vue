<template>
    <div class="group-wrapper">
        <p class="btn-group">
            <bk-button type="primary" @click="showGroupDialog(false)">
                {{$t('ModelManagement["新建分组"]')}}
            </bk-button>
            <bk-button type="default" @click="showModelDialog(false)">
                {{$t('ModelManagement["新增模型"]')}}
            </bk-button>
        </p>
        <ul class="group-list">
            <li class="group-item clearfix" v-for="(classification, classIndex) in localClassifications" :key="classIndex">
                <p class="group-title">
                    <span>{{classification['bk_classification_name']}}</span>
                    <span class="number">({{classification['bk_objects'].length}})</span>
                    <i class="icon-cc-edit text-primary"
                    v-if="classification['bk_classification_type'] !== 'inner'"
                    @click="showGroupDialog(true, classification)"></i>
                    <i class="icon-cc-del text-primary"
                    v-if="classification['bk_classification_type'] !== 'inner'"
                    @click="deleteGroup(classification)"></i>
                </p>
                <ul class="model-list clearfix">
                    <li class="model-item"
                    v-for="(model, modelIndex) in classification['bk_objects']"
                    :key="modelIndex"
                    @click="modelClick(model)">
                        <div class="icon-box">
                            <i class="icon" :class="model['bk_obj_icon']"></i>
                        </div>
                        <div class="model-details">
                            <p class="model-name">{{model['bk_obj_name']}}</p>
                            <p class="model-id">{{model['bk_obj_id']}}</p>
                        </div>
                    </li>
                </ul>
                <i class="bk-icon icon-angle-double-down" v-if="classification['bk_objects'].length > 8"></i>
            </li>
        </ul>
        <bk-dialog
            class="group-dialog dialog"
            :close-icon="false"
            :hasHeader="false"
            :width="600"
            :padding="0"
            :is-show.sync="groupDialog.isShow">
            <div slot="content" class="dialog-content">
                <p class="title">{{groupDialog.title}}</p>
                <div class="content">
                    <label for="">
                        <span class="label-title">
                            {{$t('ModelManagement["唯一标识"]')}}
                        </span>
                        <span class="color-danger">*</span>
                        <input type="text" class="cmdb-form-input" v-model.trim="groupDialog.data['bk_classification_id']">
                        <i class="bk-icon icon-info-circle"></i>
                    </label>
                    <label for="">
                        <span class="label-title">
                            {{$t('ModelManagement["名称"]')}}
                        </span>
                        <span class="color-danger">*</span>
                        <input type="text" class="cmdb-form-input" v-model.trim="groupDialog.data['bk_classification_name']">
                        <i class="bk-icon icon-info-circle"></i>
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
            :is-show.sync="modelDialog.isShow">
            <div slot="content" class="dialog-content">
                <p class="title">{{$t('ModelManagement["新增模型"]')}}</p>
                <div class="content clearfix">
                    <div class="content-left" @click="isIconListShow = true">
                        <div class="icon-wrapper">
                            <i class="icon-cc-host"></i>
                        </div>
                        <div class="text">{{$t('ModelManagement["点击切换"]')}}</div>
                    </div>
                    <div class="content-right">
                        <label for="">
                            <span class="label-title">{{$t('ModelManagement["所属分组"]')}}</span>
                            <span class="color-danger">*</span>
                            <bk-selector
                                class="selector-box"
                                :selected.sync="selected"
                                :list="list"
                            ></bk-selector>
                        </label>
                        <label for="">
                            <span class="label-title">{{$t('ModelManagement["唯一标识"]')}}</span>
                            <span class="color-danger">*</span>
                            <input type="text" class="cmdb-form-input" v-model.trim="modelDialog.data['bk_obj_id']">
                            <i class="bk-icon icon-info-circle"></i>
                        </label>
                        <label for="">
                            <span class="label-title">{{$t('ModelManagement["名称"]')}}</span>
                            <span class="color-danger">*</span>
                            <input type="text" class="cmdb-form-input" v-model.trim="modelDialog.data['bk_obj_name']">
                            <i class="bk-icon icon-info-circle"></i>
                        </label>
                    </div>
                </div>
                <div class="model-icon-wrapper" v-if="isIconListShow">
                    <span class="back" @click="isIconListShow = false">
                        <i class="bk-icon icon-back2"></i>
                    </span>
                    <the-choose-icon
                    ></the-choose-icon>
                </div>
            </div>
            <div slot="footer" class="footer">
                <bk-button type="primary">{{$t("Common['保存']")}}</bk-button>
                <bk-button type="default">{{$t("Common['取消']")}}</bk-button>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import theChooseIcon from './_choose-icon'
    import theModel from './children'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    export default {
        components: {
            theChooseIcon,
            theModel
        },
        data () {
            return {
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
                modelDialog: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('ModelManagement["新增模型"]'),
                    data: {
                        bk_obj_icon: '',
                        bk_obj_id: '',
                        bk_obj_name: ''
                    }
                },
                isIconListShow: false,
                list: [],
                selected: '',
                iconList: [],
                localClassifications: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', [
                'classifications'
            ])
        },
        created () {
            this.classifications.forEach(classification => {
                this.localClassifications.push({...classification, ...{isModelShow: false}})
            })
        },
        methods: {
            ...mapMutations('objectModelClassify', [
                'updateClassify',
                'deleteClassify'
            ]),
            ...mapActions('objectModelClassify', [
                'createClassification',
                'updateClassification',
                'deleteClassification'
            ]),
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
                    title: this.$t('ModelManagement["确认要删除此分组？"]'),
                    confirmFn: async () => {
                        await this.deleteClassification({
                            id: group.id
                        })
                        this.$store.commit('objectModelClassify/deleteClassify', group['bk_classification_id'])
                    }
                })
            },
            showModelDialog () {

            },
            hideModelDialog () {

            },
            async saveModel () {

            },
            modelClick (model) {
                this.$router.push(`model-manage/${model['bk_obj_id']}`)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .btn-group {
        margin: 20px 0;
        font-size: 0;
        .bk-primary {
            margin-right: 10px;
        }
    }
    .group-list {
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
                    transform: translate(180deg);
                }
            }
        }
        .group-title {
            display: inline-block;
            padding-left: 8px;
            border-left: 4px solid $cmdbBorderFocusColor;
            line-height: 16px;
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
        .model-item {
            float: left;
            margin: 10px 10px 0 0;
            width: calc((100% - 30px) / 4);
            height: 70px;
            border: 1px solid $cmdbTableBorderColor;
            border-radius: 4px;
            cursor: pointer;
            &:nth-child(4n) {
                margin-right: 0;
            }
            &:hover {
                border-color: $cmdbBorderFocusColor;
                background: #ebf4ff;
            }
            .icon-box {
                float: left;
                .icon {
                    padding-left: 18px;
                    font-size: 32px;
                    line-height: 70px;
                    color: $cmdbBorderFocusColor;
                }
            }
            .model-details {
                float: left;
                line-height: 16px;
                margin-top: 20px;
                padding-left: 10px;
            }
            .model-name {
                font-size: 14px;
            }
            .model-id {
                font-size: 12px;
                color: $cmdbTableBorderColor;
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
                line-height: 36px;
                vertical-align: middle;
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
            .selector-box,
            .cmdb-form-input {
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
