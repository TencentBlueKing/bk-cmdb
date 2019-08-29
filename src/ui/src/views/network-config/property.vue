<template>
    <div class="device-wrapper">
        <div class="title">
            <bk-button theme="primary" @click="showPropertyDialog('create')">
                {{$t('新增属性')}}
            </bk-button>
            <bk-button theme="default"
                :loading="$loading('deleteNetcollectProperty')"
                :disabled="!table.checked.length"
                @click="deleteProperty">
                {{$t('删除')}}
            </bk-button>
            <bk-button theme="default" @click="importSlider.isShow = true">
                {{$t('导入')}}
            </bk-button>
            <bk-button theme="default" form="exportForm" :disabled="!table.checked.length">
                {{$t('导出')}}
            </bk-button>
            <form id="exportForm" :action="url.export" method="POST" hidden>
                <input type="hidden" name="netcollect_property_id" :value="table.checked.join(',')">
            </form>
            <div class="filter">
                <bk-select class="search-selector"
                    v-model="filter.type">
                    <bk-option v-for="(option, index) in filter.typeList"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
                <bk-input class="cmdb-form-input"
                    type="text"
                    :placeholder="$t('请输入')"
                    v-model.trim="filter.text"
                    @enter="getTableData">
                </bk-input>
                <i class="bk-icon icon-search"
                    @click="getTableData"></i>
            </div>
        </div>
        <cmdb-table
            class="device-table"
            :loading="$loading('searchNetcollectProperty')"
            :header="table.header"
            :list="table.list"
            :checked.sync="table.checked"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange"
            @handleRowClick="handleRowClick">
        </cmdb-table>
        <bk-dialog
            class="create-dialog"
            v-model="propertyDialog.isShow"
            :title="propertyDialog.title"
            :show-footer="false"
            :close-icon="false"
            :width="424">
            <div>
                <div>
                    <label class="label first">
                        <span>{{$t('所属设备')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-select
                        searchable
                        v-model="propertyDialog.data.device_id">
                        <bk-option v-for="(option, index) in propertyDialog.deviceList"
                            :key="index"
                            :id="option.device_id"
                            :name="option.device_name">
                        </bk-option>
                    </bk-select>
                    <bk-input type="text" hidden name="device_id"
                        v-model="propertyDialog.data['device_id']"
                        v-validate="'required'">
                    </bk-input>
                    <div v-show="errors.has('device_id')" class="color-danger">{{ errors.first('device_id') }}</div>
                </div>
                <div>
                    <label class="label">
                        <span>oid<span class="color-danger">*</span></span>
                    </label>
                    <bk-input type="text" class="cmdb-form-input" name="oid"
                        v-model="propertyDialog.data.oid"
                        v-validate="'required|oid'">
                    </bk-input>
                    <div v-show="errors.has('oid')" class="color-danger">{{ errors.first('oid') }}</div>
                </div>
                <div>
                    <label class="label">
                        <span>{{$t('采集方式')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-select v-model="propertyDialog.data.action">
                        <bk-option v-for="(option, index) in propertyDialog.actionList"
                            :key="index"
                            :id="option.id"
                            :name="option.name">
                        </bk-option>
                    </bk-select>
                    <bk-input type="text" hidden name="action"
                        v-model="propertyDialog.data.action"
                        v-validate="'required'">
                    </bk-input>
                    <div v-show="errors.has('action')" class="color-danger">{{ errors.first('action') }}</div>
                </div>
                <div>
                    <label class="label">
                        <span>{{$t('模型属性')}}<span class="color-danger">*</span></span>
                    </label>
                    <bk-select
                        searchable
                        v-model="propertyDialog.data.bk_property_id"
                        :scroll-height="200">
                        <bk-option v-for="(option, index) in propertyDialog.attrList"
                            :key="index"
                            :id="option.bk_property_id"
                            :name="option.bk_property_name">
                        </bk-option>
                    </bk-select>
                    <bk-input type="text" hidden name="bk_property_id"
                        v-model="propertyDialog.data['bk_property_id']"
                        v-validate="'required'">
                    </bk-input>
                    <div v-show="errors.has('bk_property_id')" class="color-danger">{{ errors.first('bk_property_id') }}</div>
                </div>
                <div class="footer">
                    <bk-button theme="primary" @click="saveProperty" :loading="$loading(['createNetcollectProperty', 'updateNetcollectProperty'])">
                        {{$t('保存')}}
                    </bk-button>
                    <bk-button theme="default" @click="hidePropertyDialog">
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
                importSlider: {
                    isShow: false
                },
                propertyDialog: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('新增属性'),
                    deviceList: [],
                    attrList: [],
                    actionList: [{
                        id: 'get',
                        name: 'get'
                    }, {
                        id: 'getnext',
                        name: 'getnext'
                    }],
                    data: {
                        oid: '',
                        action: '',
                        device_id: '',
                        bk_property_id: '',
                        netcollect_property_id: ''
                    }
                },
                filter: {
                    typeList: [{
                        id: 'device_name',
                        name: this.$t('所属设备')
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('模型')
                    }, {
                        id: 'bk_property_name',
                        name: this.$t('模型属性')
                    }],
                    type: 'device_name',
                    text: ''
                },
                table: {
                    header: [{
                        id: 'netcollect_property_id',
                        type: 'checkbox'
                    }, {
                        id: 'netcollect_property_id',
                        name: 'ID'
                    }, {
                        id: 'device_name',
                        name: this.$t('所属设备')
                    }, {
                        id: 'unit',
                        name: this.$t('计量单位')
                    }, {
                        id: 'oid',
                        name: 'oid'
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('模型')
                    }, {
                        id: 'bk_property_name',
                        name: this.$t('模型属性')
                    }],
                    list: [],
                    checked: [],
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
                const prefix = `${window.API_HOST}netproperty/`
                return {
                    import: prefix + 'import',
                    export: prefix + 'export',
                    template: `${window.API_HOST}netcollect/importtemplate/netproperty`
                }
            }
        },
        watch: {
            'propertyDialog.data.device_id' (deviceId) {
                if (deviceId) {
                    const device = this.propertyDialog.deviceList.find(device => device.device_id === deviceId)
                    if (device) {
                        this.getObjAttr(device)
                    }
                }
            }
        },
        async created () {
            const res = await this.searchDevice({ params: {}, config: { requestId: 'searchDevice', fromCache: true } })
            this.propertyDialog.deviceList = res.info
            this.getTableData()
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('netCollectDevice', ['searchDevice']),
            ...mapActions('netCollectProperty', [
                'createNetcollectProperty',
                'updateNetcollectProperty',
                'searchNetcollectProperty',
                'deleteNetcollectProperty'
            ]),
            deleteProperty () {
                this.$bkInfo({
                    title: this.$t('确认删除属性'),
                    confirmFn: async () => {
                        const params = {
                            netcollect_property_id: this.table.checked
                        }
                        await this.deleteNetcollectProperty({ config: { data: params, requestId: 'deleteNetcollectProperty' } })
                        this.table.checked = []
                        this.handlePageChange(1)
                    }
                })
            },
            async showPropertyDialog (type) {
                if (type === 'create') {
                    this.propertyDialog.isEdit = false
                    this.propertyDialog.data['device_id'] = ''
                    this.propertyDialog.data['oid'] = ''
                    this.propertyDialog.data['action'] = ''
                    this.propertyDialog.data['bk_property_id'] = ''
                    this.propertyDialog.title = this.$t('新增属性')
                } else {
                    this.propertyDialog.isEdit = true
                    this.propertyDialog.title = this.$t('编辑属性')
                }
                this.propertyDialog.isShow = true
            },
            hidePropertyDialog () {
                this.propertyDialog.isShow = false
            },
            async getObjAttr (device) {
                this.propertyDialog.attrList = await this.searchObjectAttribute({
                    params: this.$injectMetadata({ bk_obj_id: device['bk_obj_id'] }),
                    config: {
                        requestId: `post_searchObjectAttribute_${device['bk_obj_id']}`,
                        fromCache: true
                    }
                })
            },
            async saveProperty () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                const params = {
                    device_id: this.propertyDialog.data['device_id'],
                    bk_property_id: this.propertyDialog.data['bk_property_id'],
                    oid: this.propertyDialog.data.oid,
                    action: this.propertyDialog.data.action
                }
                if (this.propertyDialog.isEdit) {
                    await this.updateNetcollectProperty({
                        propertyId: this.propertyDialog.data['netcollect_property_id'],
                        params,
                        config: { requestId: 'updateNetcollectProperty' }
                    })
                } else {
                    await this.createNetcollectProperty({ params, config: { requestId: 'createNetcollectProperty' } })
                }
                this.hidePropertyDialog()
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
                if (this.filter.text.length) {
                    Object.assign(params, { condition: [{ field: this.filter.type, operator: '$regex', value: this.filter.text }] })
                }
                const res = await this.searchNetcollectProperty({ params, config: { requestId: 'searchNetcollectProperty' } })
                this.table.pagination.count = res.count
                this.table.list = res.info
            },
            handleRowClick (item) {
                this.propertyDialog.data['device_id'] = item['device_id']
                this.propertyDialog.data['oid'] = item['oid']
                this.propertyDialog.data['action'] = item['action']
                this.propertyDialog.data['bk_property_id'] = item['bk_property_id']
                this.propertyDialog.data['netcollect_property_id'] = item['netcollect_property_id']
                this.showPropertyDialog('update')
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
            >.bk-default {
                margin-left: 10px;
            }
            .filter {
                position: relative;
                float: right;
                height: 36px;
                line-height: 36px;
                .search-selector {
                    position: relative;
                    float: left;
                    width: 120px;
                    margin-right: -1px;
                    z-index: 1;
                }
                .cmdb-form-input {
                    position: relative;
                    width: 260px;
                    border-radius: 0 2px 2px 0;
                    &:focus {
                        z-index: 2;
                    }
                }
                .icon-search {
                    position: absolute;
                    right: 10px;
                    top: 11px;
                    cursor: pointer;
                    font-size: 14px;
                    z-index: 3;
                }
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
