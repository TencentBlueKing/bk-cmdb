<template>
    <div class="update-template-wrapper">
        <div class="basic-info">
            <h3>基本属性</h3>
            <div class="info-box clearfix">
                <div class="info-item fl">
                    <span class="item-title">
                        {{$t('ServiceManagement["模板名称"]')}}：
                    </span>
                    <span class="item-main">
                        Mysql_DBA模版
                    </span>
                </div>
                <div class="info-item fl">
                    <span class="info-title fl">
                        {{$t('ServiceManagement["服务分类"]')}}：
                    </span>
                    <div class="item-main fl" v-if="showEidtClassification">
                        <bk-selector
                            placeholder="请选择一级分类"
                            :list="[]"
                            :selected.sync="formData.primaryClassification">
                        </bk-selector>
                        <bk-selector
                            placeholder="请选择二级分类"
                            :list="[]"
                            :selected.sync="formData.secondaryClassification">
                        </bk-selector>
                        <div class="operation-btn">
                            <span class="text-primary" @click="handleEditSave">{{$t("Common['保存']")}}</span>
                            <span class="text-primary ml10" @click="handleEditCancel">{{$t("Common['取消']")}}</span>
                        </div>
                    </div>
                    <div class="item-main fl" v-else>
                        <span>数据库 / Mysql</span>
                        <i class="bk-icon icon-edit2" @click="handleEdit"></i>
                    </div>
                </div>
                <div class="info-item fl">
                    <span class="item-title">
                        {{$t('ServiceManagement["应用实例"]')}}：
                    </span>
                    <span style="color: #3a84ff;">
                        5
                    </span>
                </div>
            </div>
        </div>
        <div class="process-info">
            <h3>进程服务</h3>
            <div class="precess-box">
                <div class="process-create">
                    <bk-button type="primary" class="create-btn" @click="handleCreateProcess">
                        <span>{{$t("ProcessManagement['添加进程']")}}</span>
                    </bk-button>
                    <span class="create-tips">{{$t("ServiceManagement['新建进程提示']")}}</span>
                </div>
                <process-table></process-table>
            </div>
        </div>
        <cmdb-slider :is-show.sync="slider.show" :title="slider.title">
            <template slot="content">
                <process-form
                    :properties="properties"
                    :property-groups="propertyGroups"
                    :object-unique="objectUnique"
                    :inst="attribute.inst.details"
                    :type="attribute.type"
                    :save-disabled="true"
                    @on-submit="handleSliderSave"
                    @on-cancel="handleSliderCancel">
                </process-form>
            </template>
        </cmdb-slider>
    </div>
</template>

<script>
    import processForm from '@/components/service/process-form'
    import processTable from './process'
    import { mapActions } from 'vuex'
    export default {
        components: {
            processTable,
            processForm
        },
        data () {
            return {
                properties: [],
                propertyGroups: [],
                objectUnique: [],
                attribute: {
                    type: null,
                    inst: {
                        details: {},
                        edit: {}
                    }
                },
                slider: {
                    show: false,
                    title: ''
                },
                formData: {
                    primaryClassification: '',
                    secondaryClassification: '',
                    tempalteName: ''
                },
                showEidtClassification: false
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t("ServiceManagement['新建服务模版']"))
            this.reload()
        },
        methods: {
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectUnique', ['searchObjectUniqueConstraints']),
            async reload () {
                this.properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: 'process',
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_process`,
                        fromCache: false
                    }
                })
                this.getPropertyGroups()
                this.getObjectUnique()
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: 'process',
                    params: this.$injectMetadata(),
                    config: {
                        fromCache: false,
                        requestId: 'post_searchGroup_process'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            async getObjectUnique () {
                this.objectUnique = await this.searchObjectUniqueConstraints({
                    objId: 'process',
                    params: {},
                    config: {
                        requestId: 'searchObjectUniqueConstraints'
                    }
                })
            },
            handleEdit () {
                this.showEidtClassification = true
            },
            handleEditSave () {

            },
            handleEditCancel () {
                this.showEidtClassification = false
            },
            handleSliderSave () {

            },
            handleSliderCancel () {
                this.slider.show = false
            },
            handleCreateProcess () {
                this.slider.show = true
                this.slider.title = this.$t("ProcessManagement['添加进程']")
            }
        }
    }
</script>

<style lang="scss" scoped>
    .update-template-wrapper {
        h3 {
            color: #63656e;
            font-size: 14px;
            padding-bottom: 26px;
        }
        .basic-info {
            padding-bottom: 56px;
            .info-box {
                font-size: 14px;
                padding-left: 30px;
            }
            .info-item {
                line-height: 32px;
                padding-right: 80px;
            }
            .bk-selector {
                width: 200px;
                float: left;
                &:first-child {
                    margin-right: 10px;
                }
            }
            .icon-edit2 {
                margin-left: 10px;
                color: #3a84ff;
                cursor: pointer;
            }
            .operation-btn {
                float: left;
                padding-left: 12px;
                span:first-child {
                    margin-right: 5px;
                    position: relative;
                    &::after {
                        content: '';
                        position: absolute;
                        top: 25%;
                        right: -10px;
                        width: 2px;
                        height: 12px;
                        background-color: #dcdee5;
                    }
                }
            }
        }
        .process-create {
            padding-bottom: 20px;
            .create-tips {
                font-size: 12px;
                color: #979ba5;
                padding-left: 10px;
            }
        }
    }
</style>
