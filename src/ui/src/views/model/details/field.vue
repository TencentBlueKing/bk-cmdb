<template>
    <div class="field-wrapper">
        <v-base-info :isReadOnly="isReadOnly" :isEdit.sync="isEdit" @confirm="confirm" @cancel="cancel"></v-base-info>
        <div class="field-content clearfix" v-if="isEdit">
            <div class="create-field clearfix">
                <bk-button v-if="!isReadOnly" type="primary" :title="$t('ModelManagement[\'新增字段\']')" @click="createField">
                    {{$t('ModelManagement["新增字段"]')}}
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
                    <li class="table-item">{{$t('ModelManagement["唯一"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["必填字段"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["类型"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["字段名"]')}}</li>
                    <li class="table-item">{{$t('ModelManagement["操作"]')}}</li>
                </ul>
                <ul class="table-content" v-bkloading="{isLoading: $loading(['initFieldList', 'deleteObjectAttribute'])}">
                    <li v-for="(field, index) in fieldList" :key="index">
                        <ul class="field-item clearfix" :class="{'disabled': field['ispre'] || isReadOnly}" @click="toggleDetailShow(field, index)">
                            <li class="table-item"><i class="bk-icon icon-check-1" v-show="field['isonly']"></i></li>
                            <li class="table-item"><i class="bk-icon icon-check-1" v-show="field['isrequired']"></i></li>
                            <li class="table-item">{{fieldTypeMap[field['bk_property_type']]}}</li>
                            <li class="table-item" :title="`${field['bk_property_name']}(${field['bk_property_id']})`">{{field['bk_property_name']}}({{field['bk_property_id']}})</li>
                            <li class="table-item">
                                <div class="btn-contain">
                                    <i class="icon-cc-del"  v-if="!field['ispre'] && !isReadOnly"
                                        :class="{'editable':field['ispre'] || isReadOnly}"
                                        @click.stop="showConfirmDialog(field, index)"
                                    ></i>
                                </div>
                            </li>
                        </ul>
                        <div class="field-detail" v-if="field.isShow">
                            <v-model-field
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
        </div>
        <v-create-field v-if="createForm.isShow"
            @save="saveCreateForm"
            @cancel="closeCreateForm"
        ></v-create-field>
    </div>
</template>

<script>
    import vBaseInfo from './base-info'
    import vModelField from './model-field'
    import vCreateField from './create-field'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            vBaseInfo,
            vModelField,
            vCreateField
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
                'supplierAccount',
                'site'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                return this.activeModel['bk_ispaused']
            },
            exportUrl () {
                return `${this.site.url}object/owner/${this.supplierAccount}/object/${this.activeModel['bk_obj_id']}/export`
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
            createField () {
                this.createForm.isShow = true
            },
            saveCreateForm () {
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
                })
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
    .field-wrapper {
        position: relative;
        padding: 0 30px;
        height: 100%;
        .field-content {
            height: calc(100% - 140px);
            padding: 20px 0;
            border-top: 1px solid $cmdbBorderLightColor;
        }
        .create-field {
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
                height: 40px;
                line-height: 40px;
                background: $cmdbPrimaryColor;
                text-align: center;
                border: 1px solid $cmdbBorderLightColor;
                border-right: none;
                &:nth-child(1) {
                    width:80px;
                }
                &:nth-child(2) {
                    width:80px;
                }
                &:nth-child(3) {
                    width:145px;
                }
                &:nth-child(4) {
                    width:290px;
                }
                &:nth-child(5) {
                    width: 105px;
                    border-right: 1px solid $cmdbBorderLightColor;
                }
            }
            .table-content {
                width: 720px;
                height: calc(100% - 40px);
                @include scrollbar;
                .field-item {
                    cursor: pointer;
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
                        .icon-cc-del {
                            color: #ccc;
                        }
                    }
                }
                .field-detail {
                    width: 700px;
                    padding: 30px 14px 0;
                    border: 1px solid $cmdbBorderLightColor;
                    border-top: none;
                }
            }
        }
    }
</style>
