<template>
    <div class="field-wrapper">
        <v-base-info ref="baseInfo" 
        :isReadOnly="isReadOnly" 
        :isEdit.sync="isEdit" 
        @createObject="createObject" 
        @confirm="confirm" 
        @cancel="cancel"></v-base-info>
        <div class="field-content clearfix" v-if="isEdit">
            <div class="create-field clearfix">
                <bk-button class="create-btn" :class="{'open': createForm.isShow}" v-if="!isReadOnly" type="primary" :title="$t('ModelManagement[\'新增字段\']')" @click="createForm.isShow = !createForm.isShow">
                    {{$t('ModelManagement["新增字段"]')}}<i class="bk-icon icon-angle-down"></i>
                    <i class="create-btn-mask"></i>
                </bk-button>
                <div class="btn-group">
                    <bk-button type="default" :loading="$loading('importObjectAttribute')" :title="$t('ModelManagement[\'导入\']')" :disabled="isReadOnly" class="btn mr10">
                        {{$t('ModelManagement["导入"]')}}
                        <input v-if="!isReadOnly" ref="fileInput" type="file" @change.prevent="handleFile">
                    </bk-button>
                    <form :action="exportUrl" method="POST" class="form">
                        <bk-button type="default submit" :title="$t('ModelManagement[\'导出\']')">
                            {{$t('ModelManagement["导出"]')}}
                        </bk-button>
                    </form>
                </div>
            </div>
            <div class="field-table">
                <ul class="table-title">
                    <li class="table-item">{{$t('ModelManagement["字段名"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["类型"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["唯一"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["必填字段"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["操作"]')}}</li>
                </ul>
                <ul class="table-content" v-bkloading="{isLoading: $loading(['initFieldList', 'deleteObjectAttribute'])}">
                    <li v-for="(field, index) in fieldList" :key="index">
                        <ul class="field-item clearfix" :class="{'disabled': field['ispre'] || isReadOnly}">
                            <li class="table-item" :title="`${field['bk_property_name']}(${field['bk_property_id']})`">{{field['bk_property_name']}}({{field['bk_property_id']}})</li>
                            <li class="table-item">{{fieldTypeMap[field['bk_property_type']]}}</li>
                            <li class="table-item"><i class="bk-icon icon-check-1" v-show="field['isonly']"></i></li>
                            <li class="table-item"><i class="bk-icon icon-check-1" v-show="field['isrequired']"></i></li>
                            <li class="table-item">
                                <div class="btn-contain">
                                    <span class="text-primary" @click="toggleDetailShow(field, index)">{{$t('Common["编辑"]')}}</span>
                                    <span class="text-danger" v-if="!field['ispre'] && !isReadOnly"
                                    @click.stop="showConfirmDialog(field, index)">
                                        {{$t('Common["删除"]')}}
                                    </span>
                                </div>
                            </li>
                        </ul>
                        <div class="field-detail" v-if="field.isShow">
                            <v-model-field
                                ref="modelField"
                                :isReadOnly="isReadOnly"
                                :isEditField="true"
                                :field="field"
                                @save="initFieldList"
                                @cancel="toggleDetailShow(field, index)">
                            </v-model-field>
                        </div>
                    </li>
                </ul>
            </div>
            <div class="field-mask"></div>
        </div>
        <transition name="slide">
            <div class="create-field-wrapper" v-if="createForm.isShow">
                <v-model-field
                    :isEditField="false"
                    @save="saveCreateForm"
                    @cancel="closeCreateForm"
                ></v-model-field>
            </div>
        </transition>
    </div>
</template>

<script>
    import vBaseInfo from './base-info'
    import vModelField from './model-field'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            vBaseInfo,
            vModelField
        },
        props: {
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                fieldTypeMap: {
                    'singlechar': this.$t('ModelManagement["短字符"]'),
                    'int': this.$t('ModelManagement["数字"]'),
                    'enum': this.$t('ModelManagement["枚举"]'),
                    'date': this.$t('ModelManagement["日期"]'),
                    'time': this.$t('ModelManagement["时间"]'),
                    'longchar': this.$t('ModelManagement["长字符"]'),
                    'singleasst': this.$t('ModelManagement["单关联"]'),
                    'multiasst': this.$t('ModelManagement["多关联"]'),
                    'objuser': this.$t('ModelManagement["用户"]'),
                    'timezone': this.$t('ModelManagement["时区"]'),
                    'bool': 'bool'
                },
                createForm: {
                    isShow: false
                },
                fieldList: []
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                return this.activeModel['bk_ispaused']
            },
            exportUrl () {
                return `${window.API_HOST}object/owner/${this.supplierAccount}/object/${this.activeModel['bk_obj_id']}/export`
            }
        },
        created () {
            this.initFieldList()
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute',
                'deleteObjectAttribute'
            ]),
            ...mapActions('objectBatch', [
                'importObjectAttribute'
            ]),
            createObject () {
                this.$emit('createObject')
                this.initFieldList()
            },
            isCloseConfirmShow () {
                if (this.$refs.baseInfo.isCloseConfirmShow()) {
                    return true
                }
                if (this.$refs.modelField && this.$refs.modelField.length && this.$refs.modelField[0].isCloseConfirmShow()) {
                    return true
                }
                return false
            },
            createField () {
                this.createForm.isShow = true
            },
            saveCreateForm (type) {
                if (['singleasst', 'multiasst'].indexOf(type) !== -1) {
                    this.$emit('updateTopo', true)
                }
                this.initFieldList()
                this.closeCreateForm()
            },
            closeCreateForm () {
                this.createForm.isShow = false
            },
            showConfirmDialog (field, index) {
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除字段？"]', field['bk_property_name'], {name: field['bk_property_name']}),
                    confirmFn: () => {
                        this.deleteField(field, index)
                    }
                })
            },
            async deleteField (field, index) {
                await this.deleteObjectAttribute({
                    id: field.id,
                    config: {
                        requestId: 'deleteObjectAttribute'
                    }
                }).then(() => {
                    this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                })
                if (field['bk_property_type'] === 'singleasst' || field['bk_property_type'] === 'multiasst') {
                    this.$emit('updateTopo', true)
                }
                this.fieldList.splice(index, 1)
            },
            toggleDetailShow (field, index) {
                if (!field.isShow) {
                    let fieldShow = this.fieldList.find(({isShow}) => isShow)
                    if (fieldShow) {
                        fieldShow.isShow = false
                    }
                }
                field.isShow = !field.isShow
            },
            async initFieldList () {
                const res = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: this.activeModel['bk_obj_id']
                    },
                    config: {
                        requestId: 'initFieldList'
                    }
                })
                let listHeader = []
                let listBody = []
                let listTail = []
                res.map(item => {
                    Object.assign(item, {isShow: false})
                    if (item['isonly'] && item['isrequired']) {
                        listHeader.push(item)
                    } else if (item['isonly']) {
                        listBody.push(item)
                    } else if (item['isrequired']) {
                        listBody.unshift(item)
                    } else {
                        listTail.push(item)
                    }
                })
                this.fieldList = [...listHeader, ...listBody, ...listTail]
            },
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
                            originalResponse: true
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
                            this.$success(this.$t('ModelManagement["导入成功"]'), 'success')
                            this.initFieldList()
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
            confirm () {
                this.$emit('confirm')
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .slide-enter-active, .slide-leave-active{
        transition: height .2s;
        overflow: hidden;
        height: 338px;
    }
    .slide-enter, .slide-leave-to{
        height: 0;
    }
    .field-wrapper {
        position: relative;
        padding: 0 10px;
        height: 100%;
        .field-content {
            height: calc(100% - 132px);
            padding-top: 20px;
            border-top: 1px solid $cmdbBorderLightColor;
            .field-mask {
                position: absolute;
                bottom: 57px;
                left: -20px;
                width: calc(100% + 40px);
                height: 52px;
                background: linear-gradient(to bottom, rgba(255, 255, 255, 0), rgba(255, 255, 255, .8));
                pointer-events: none;
            }
        }
        .create-field {
            .create-btn {
                position: relative;
                z-index: 2;
                &.open {
                    background: #fff;
                    color: #737987;
                    border: 1px solid #c3cdd7;
                    border-bottom-color: transparent;
                    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.2);
                    .icon-angle-down {
                        transform: rotate(180deg);
                    }
                    .create-btn-mask {
                        width: calc(100% + 2px);
                        height: 12px;
                        background: #fff;
                        display: inline-block;
                        position: absolute;
                        bottom: -12px;
                        left: -1px;
                        border-left: 1px solid #c3cdd7;
                        border-right: 1px solid #c3cdd7;
                    }
                }
                .icon-angle-down {
                    margin-left: 5px;
                    font-size: 12px;
                    font-weight: bold;
                    transition: all .2s linear;
                }
            }
            .btn-group {
                float: right;
                font-size: 0;
                .btn {
                    position: relative;
                    overflow: hidden;
                    input[type="file"] {
                        position: absolute;
                        left: -70px;
                        top: 0;
                        opacity: 0;
                        width: calc(100% + 70px);
                        height: 100%;
                        cursor: pointer;
                    }
                }
                .form {
                    display: inline-block;
                }
            }
        }
        .field-table {
            margin-top: 10px;
            height: calc(100% - 46px);
            .table-title {
                width: 100%;
                height: 40px;
            }
            .table-item {
                float: left;
                padding: 0 20px;
                height: 40px;
                line-height: 40px;
                background: $cmdbPrimaryColor;
                border: 1px solid $cmdbBorderLightColor;
                border-right: none;
                &:nth-child(1) {
                    width: 290px;
                }
                &:nth-child(2) {
                    width: 145px;
                }
                &:nth-child(3) {
                    width: 95px;
                }
                &:nth-child(4) {
                    width: 95px;
                }
                &:nth-child(5) {
                    width: 115px;
                    border-right: 1px solid $cmdbBorderLightColor;
                }
            }
            .table-content {
                padding-bottom: 20px;
                width: 760px;
                height: calc(100% - 40px);
                @include scrollbar;
                .field-item {
                    &.disabled {
                        color: #ccc;
                    }
                    >.table-item {
                        background: #fff;
                        border-left: none;
                        border-top: none;
                        &:first-child {
                            border-left: 1px solid $cmdbBorderLightColor;
                        }
                        .btn-contain {
                            font-size: 0;
                            >span {
                                cursor: pointer;
                                font-size: 14px;
                                &:first-child {
                                    margin-right: 10px;
                                }
                            }
                        }
                    }
                }
                .field-detail {
                    width: 740px;
                    border: 1px solid $cmdbBorderLightColor;
                    border-top: none;
                }
            }
        }
        .create-field-wrapper {
            position: absolute;
            z-index: 1;
            top: 143px;
            left: 10px;
            right: 10px;
            background: #fff;
            border: 1px solid #c3cdd7;
            box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.2);
            .title {
                width: 100%;
                height: 40px;
                background: linear-gradient(to bottom, #f9f9f9, #fff);
                cursor: pointer;
                text-align: center;
                img {
                    margin-top: 14px;
                }
            }
        }
    }
</style>
