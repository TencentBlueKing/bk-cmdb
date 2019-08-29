<template>
    <div class="device-wrapper">
        <div class="title">
            <bk-button theme="primary" @click="showDeviceDialog('create')">
                {{$t('新增设备')}}
            </bk-button>
            <bk-button theme="default"
                :loading="$loading('deleteDevice')"
                :disabled="!table.checked.length"
                @click="deleteDevices">
                {{$t('删除')}}
            </bk-button>
            <bk-button theme="default" @click="importSlider.isShow = true">
                {{$t('导入')}}
            </bk-button>
            <bk-button :disabled="!table.checked.length" type="submit" form="exportForm">
                {{$t('导出')}}
            </bk-button>
            <form id="exportForm" :action="url.export" method="POST" hidden>
                <input type="hidden" name="device_id" :value="table.checked.join(',')">
            </form>
        </div>
        <cmdb-table
            class="device-table"
            :loading="$loading('searchDevice')"
            :header="table.header"
            :list="table.list"
            :checked.sync="table.checked"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange"
            @handleRowClick="handleRowClick">
            <template slot="bk_obj_id" slot-scope="{ item }">
                <template>{{getObjName(item['bk_obj_id'])}}</template>
            </template>
        </cmdb-table>
        <bk-dialog
            class="create-dialog"
            v-model="deviceDialog.isShow"
            :title="deviceDialog.title"
            :show-footer="false"
            :close-icon="false"
            :width="424">
            <div>
                <div>
                    <label class="label first">
                        <span>{{$t('设备型号')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-input type="text"
                        class="cmdb-form-input"
                        name="device_model"
                        v-model.trim="deviceDialog.data['device_model']"
                        v-validate="'required|singlechar|length:256'">
                    </bk-input>
                    <div v-show="errors.has('device_model')" class="color-danger">{{ errors.first('device_model') }}</div>
                </div>
                <div>
                    <label class="label">
                        <span>{{$t('设备名称')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-input type="text" class="cmdb-form-input" name="device_name"
                        v-model.trim="deviceDialog.data['device_name']"
                        v-validate="'required|singlechar|length:256'">
                    </bk-input>
                    <div v-show="errors.has('device_name')" class="color-danger">{{ errors.first('device_name') }}</div>
                </div>
                <div>
                    <label class="label">
                        <span>{{$t('对应模型')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-select v-model="deviceDialog.data.bk_obj_id">
                        <bk-option v-for="(option, index) in netList"
                            :key="index"
                            :id="option.bk_obj_id"
                            :name="option.bk_obj_name">
                        </bk-option>
                    </bk-select>
                    <input type="text" hidden name="bk_obj_id" v-model.trim="deviceDialog.data['bk_obj_id']" v-validate="'required'">
                    <div v-show="errors.has('bk_obj_id')" class="color-danger">{{ errors.first('bk_obj_id') }}</div>
                </div>
                <div>
                    <label class="label">
                        <span>{{$t('厂商')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-input type="text" class="cmdb-form-input" name="bk_vendor"
                        v-model.trim="deviceDialog.data['bk_vendor']"
                        v-validate="'required|singlechar|length:256'">
                    </bk-input>
                    <div v-show="errors.has('bk_vendor')" class="color-danger">{{ errors.first('bk_vendor') }}</div>
                </div>
                <div class="footer">
                    <bk-button theme="primary" @click="saveDevice" :loading="$loading(['createDevice', 'updateDevice'])">
                        {{$t('保存')}}
                    </bk-button>
                    <bk-button theme="default" @click="hideDeviceDialog">
                        {{$t('取消')}}
                    </bk-button>
                </div>
            </div>
        </bk-dialog>
        <bk-sideslider
            :width="800"
            :is-show.sync="importSlider.isShow"
            :title="$t('批量导入')">
            <cmdb-import v-if="importSlider.isShow" slot="content"
                :template-url="url.template"
                :import-url="url.import"
                @success="handlePageChange(1)"
                @partialSuccess="handlePageChange(1)">
            </cmdb-import>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import cmdbImport from '@/components/import/import'
    export default {
        components: {
            cmdbImport
        },
        data () {
            return {
                deviceDialog: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('新增设备'),
                    data: {
                        device_model: '',
                        device_name: '',
                        bk_obj_id: '',
                        bk_vendor: '',
                        device_id: ''
                    }
                },
                netList: [],
                importSlider: {
                    isShow: false
                },
                table: {
                    header: [{
                        id: 'device_id',
                        type: 'checkbox'
                    }, {
                        id: 'device_id',
                        name: 'ID'
                    }, {
                        id: 'device_model',
                        name: this.$t('设备型号')
                    }, {
                        id: 'device_name',
                        name: this.$t('设备名称')
                    }, {
                        id: 'bk_obj_id',
                        name: this.$t('对应模型')
                    }, {
                        id: 'bk_vendor',
                        name: this.$t('厂商')
                    }],
                    checked: [],
                    list: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                }
            }
        },
        computed: {
            url () {
                const prefix = `${window.API_HOST}netdevice/`
                return {
                    import: prefix + 'import',
                    export: prefix + 'export',
                    template: `${window.API_HOST}netcollect/importtemplate/netdevice`
                }
            }
        },
        async created () {
            const res = await this.searchObjects({
                params: this.$injectMetadata({
                    bk_classification_id: 'bk_network'
                }),
                config: {
                    fromCache: true,
                    requestId: 'post_searchObjects_bk_network'
                }
            })
            this.netList = res
            this.getTableData()
        },
        methods: {
            ...mapActions('objectModel', [
                'searchObjects'
            ]),
            ...mapActions('netCollectDevice', [
                'createDevice',
                'updateDevice',
                'searchDevice',
                'deleteDevice'
            ]),
            deleteDevices () {
                this.$bkInfo({
                    title: this.$t('确认删除设备'),
                    confirmFn: async () => {
                        const params = {
                            device_id: this.table.checked
                        }
                        await this.deleteDevice({ config: { data: params, requestId: 'deleteDevice' } })
                        this.table.checked = []
                        this.handlePageChange(1)
                    }
                })
            },
            showDeviceDialog (type) {
                if (type === 'create') {
                    this.deviceDialog.isEdit = false
                    this.deviceDialog.data.device_model = ''
                    this.deviceDialog.data.device_name = ''
                    this.deviceDialog.data.bk_obj_id = ''
                    this.deviceDialog.data.bk_vendor = ''
                    this.deviceDialog.title = this.$t('新增设备')
                } else {
                    this.deviceDialog.isEdit = true
                    this.deviceDialog.title = this.$t('编辑设备')
                }
                this.deviceDialog.isShow = true
            },
            hideDeviceDialog () {
                this.deviceDialog.isShow = false
            },
            getObjName (bkObjId) {
                const obj = this.netList.find(({ bk_obj_id: objId }) => objId === bkObjId)
                if (obj) {
                    return obj['bk_obj_name']
                }
                return ''
            },
            async saveDevice () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                const params = {
                    device_model: this.deviceDialog.data['device_model'],
                    device_name: this.deviceDialog.data['device_name'],
                    bk_obj_id: this.deviceDialog.data['bk_obj_id'],
                    bk_vendor: this.deviceDialog.data['bk_vendor']
                }
                if (this.deviceDialog.isEdit) {
                    await this.updateDevice({
                        deviceId: this.deviceDialog.data['device_id'],
                        params,
                        config: { requestId: 'updateDevice' }
                    })
                } else {
                    await this.createDevice({ params, config: { requestId: 'createDevice' } })
                }
                this.hideDeviceDialog()
                this.handlePageChange(1)
            },
            async getTableData () {
                const pagination = this.table.pagination
                const params = {
                    page: {
                        start: (pagination.current - 1) * pagination.size,
                        limit: pagination.size,
                        sort: this.table.sort
                    }
                }
                const res = await this.searchDevice({ params, config: { requestId: 'searchDevice' } })
                this.table.pagination.count = res.count
                this.table.list = res.info
            },
            handleRowClick (item) {
                this.deviceDialog.data['device_model'] = item['device_model']
                this.deviceDialog.data['device_name'] = item['device_name']
                this.deviceDialog.data['bk_obj_id'] = item['bk_obj_id']
                this.deviceDialog.data['bk_vendor'] = item['bk_vendor']
                this.deviceDialog.data['device_id'] = item['device_id']
                this.showDeviceDialog('update')
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .device-wrapper {
        padding-top: 20px;
        .title {
            font-size: 0;
            .bk-default {
                margin-left: 10px;
            }
        }
        .device-table {
            margin-top: 20px;
            background: #fff;
        }
        .create-dialog {
            .label {
                &.first {
                    margin-top: 0;
                }
                display: block;
                margin: 15px 0 5px;
            }
            .footer {
                border-top: 1px solid #e5e5e5;
                padding-right: 20px;
                margin: 25px -20px -20px;
                text-align: right;
                font-size: 0;
                background: #fafbfd;
                height: 54px;
                line-height: 54px;
                .bk-default {
                    margin-left: 10px;
                }
            }
        }
    }
</style>
