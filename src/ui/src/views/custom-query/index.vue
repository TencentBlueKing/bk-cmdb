<template>
    <div class="api-wrapper">
        <feature-tips
            :feature-name="'customQuery'"
            :show-tips="showFeatureTips"
            :desc="$t('动态分组提示')"
            :more-href="'https://docs.bk.tencent.com/cmdb/Introduction.html#%EF%BC%886%EF%BC%89%E5%8A%A8%E6%80%81%E5%88%86%E7%BB%84'"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="filter-wrapper clearfix">
            <cmdb-auth class="inline-block-middle" :auth="$authResources({ type: $OPERATION.C_CUSTOM_QUERY })">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    class="api-btn"
                    :disabled="disabled"
                    @click="showUserAPISlider('create')">
                    {{$t('新建')}}
                </bk-button>
            </cmdb-auth>
            <div class="api-input fr">
                <bk-input type="text" class="cmdb-form-input"
                    right-icon="bk-icon icon-search"
                    clearable
                    v-model="filter.name"
                    font-size="medium"
                    :placeholder="$t('快速查询')"
                    @enter="getUserAPIList">
                </bk-input>
            </div>
        </div>
        <bk-table
            class="api-table"
            v-bkloading="{ isLoading: $loading('searchCustomQuery') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 229"
            :row-style="{ cursor: 'pointer' }"
            @row-click="showUserAPIDetails"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange">
            <!-- <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column> -->
            <bk-table-column prop="name" :label="$t('查询名称')" sortable="custom" fixed class-name="is-highlight"></bk-table-column>
            <bk-table-column prop="id" label="ID" sortable="custom" fixed></bk-table-column>
            <bk-table-column prop="create_user" :label="$t('创建用户')" sortable="custom"></bk-table-column>
            <bk-table-column prop="create_time" :label="$t('创建时间')" sortable="custom">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row['create_time'])}}
                </template>
            </bk-table-column>
            <bk-table-column prop="modify_user" :label="$t('修改人')" sortable="custom"></bk-table-column>
            <bk-table-column prop="last_time" :label="$t('修改时间')" sortable="custom">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row['last_time'])}}
                </template>
            </bk-table-column>
            <bk-table-column prop="operation" :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <bk-button class="mr10"
                        :text="true"
                        @click.stop="getUserAPIDetail(row)">
                        {{$t('预览')}}
                    </bk-button>
                    <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.U_CUSTOM_QUERY })">
                        <bk-button slot-scope="{ disabled }"
                            :disabled="disabled"
                            :text="true"
                            @click.stop="showUserAPIDetails(row)">
                            {{$t('编辑')}}
                        </bk-button>
                    </cmdb-auth>
                    <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.D_CUSTOM_QUERY })">
                        <bk-button slot-scope="{ disabled }"
                            :disabled="disabled"
                            :text="true"
                            @click.stop="deleteUserAPI(row)">
                            {{$t('删除')}}
                        </bk-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
            <cmdb-table-empty
                slot="empty"
                :stuff="table.stuff"
                :auth="$authResources({ type: $OPERATION.C_CUSTOM_QUERY })"
                @create="showUserAPISlider('create')"
            ></cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.isShow"
            :width="515"
            :title="slider.title"
            :before-close="handleSliderBeforeClose">
            <v-define slot="content"
                ref="define"
                v-if="slider.isShow"
                :id="slider.id"
                :biz-id="bizId"
                :type="slider.type"
                :object="object"
                @create="handleSuccess"
                @update="handleSuccess"
                @cancel="handleSliderBeforeClose">
            </v-define>
        </bk-sideslider>

        <!-- eslint-disable vue/space-infix-ops -->
        <cmdb-main-inject inject-type="prepend" v-transfer-dom>
            <v-preview ref="preview"
                v-if="isPreviewShow"
                :api-params="apiParams"
                :attribute="object"
                :table-header="previewHeader"
                @close="isPreviewShow = false">
            </v-preview>
        </cmdb-main-inject>
        <!-- eslint-disable end -->
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import vDefine from './define'
    import vPreview from './preview'
    import cmdbMainInject from '@/components/layout/main-inject'
    export default {
        components: {
            featureTips,
            vDefine,
            vPreview,
            cmdbMainInject
        },
        data () {
            return {
                showFeatureTips: false,
                isPreviewShow: false,
                previewHeader: ['bk_host_innerip', 'bk_set_name', 'bk_module_name', 'bk_biz_name', 'bk_cloud_id'],
                apiParams: {},
                filter: {
                    name: ''
                },
                table: {
                    list: [],
                    sort: '-last_time',
                    defaultSort: '-last_time',
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    stuff: {
                        type: 'default',
                        payload: {
                            resource: this.$t('动态分组')
                        }
                    }
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    type: 'create',
                    id: null,
                    title: this.$t('新建动态分组')
                },
                object: {
                    'host': {
                        id: 'host',
                        name: this.$t('主机'),
                        properties: [],
                        selected: []
                    },
                    'set': {
                        id: 'set',
                        name: this.$t('集群'),
                        properties: [],
                        selected: []
                    },
                    'module': {
                        id: 'module',
                        name: this.$t('模块'),
                        properties: [],
                        selected: []
                    },
                    'biz': {
                        id: 'biz',
                        name: this.$t('业务'),
                        properties: [],
                        selected: []
                    }
                }
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams', 'supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            searchParams () {
                const params = {
                    start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                    limit: this.table.pagination.limit,
                    sort: this.table.sort
                }
                this.filter.name ? params['condition'] = { 'name': this.filter.name } : void (0)
                return params
            }
        },
        async created () {
            this.showFeatureTips = this.featureTipsParams['customQuery']
            this.getUserAPIList()
            await this.initObjectProperties()
        },
        methods: {
            ...mapActions('hostCustomApi', [
                'searchCustomQuery',
                'getCustomQueryDetail'
            ]),
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            hideUserAPISlider () {
                this.slider.isShow = false
                this.slider.id = null
            },
            handleSliderBeforeClose () {
                if (this.$refs.define.isCloseConfirmShow()) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.hideUserAPISlider()
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.hideUserAPISlider()
                return true
            },
            handleSuccess (data) {
                this.hideUserAPISlider()
                this.handlePageChange(1)
            },
            async getUserAPIDetail (row) {
                const res = await this.getCustomQueryDetail({
                    id: row.id,
                    bizId: this.bizId,
                    config: {
                        requestId: 'getUserAPIDetail'
                    }
                })
                const properties = []
                const info = JSON.parse(res['info'])
                info.condition.forEach(condition => {
                    condition['condition'].forEach(property => {
                        const originalProperty = this.getOriginalProperty(property.field, condition['bk_obj_id'])
                        if (originalProperty) {
                            if (['time', 'date'].includes(originalProperty['bk_property_type']) && properties.some(({ propertyId }) => propertyId === originalProperty['bk_property_id'])) {
                                const repeatProperty = properties.find(({ propertyId }) => propertyId === originalProperty['bk_property_id'])
                                repeatProperty.value = [repeatProperty.value, property.value]
                            } else {
                                properties.push({
                                    'objId': originalProperty['bk_obj_id'],
                                    'objName': this.object[originalProperty['bk_obj_id']].name,
                                    'propertyType': originalProperty['bk_property_type'],
                                    'propertyName': originalProperty['bk_property_name'],
                                    'propertyId': originalProperty['bk_property_id'],
                                    'asstObjId': originalProperty['bk_asst_obj_id'],
                                    'operator': property.operator,
                                    'value': this.getUserPropertyValue(property, originalProperty)
                                })
                            }
                        }
                    })
                })
                this.apiParams = await this.getApiParams(row, properties)
                this.isPreviewShow = true
            },
            getOriginalProperty (bkPropertyId, bkObjId) {
                let property = null
                for (const objId in this.object) {
                    for (let i = 0; i < this.object[objId]['properties'].length; i++) {
                        const loopProperty = this.object[objId]['properties'][i]
                        if (loopProperty['bk_property_id'] === bkPropertyId && loopProperty['bk_obj_id'] === bkObjId) {
                            property = loopProperty
                            break
                        }
                    }
                    if (property) {
                        break
                    }
                }
                return property
            },
            getUserPropertyValue (property, originalProperty) {
                if (
                    property.operator === '$in'
                    && ['bk_module_name', 'bk_set_name'].includes(originalProperty['bk_property_id'])
                ) {
                    return property.value[property.value.length - 1]
                } else if (property.operator === '$multilike' && Array.isArray(property.value)) {
                    return property.value.join('\n')
                }
                return property.value
            },
            /* 生成保存自定义API的参数 */
            getApiParams (row, properties) {
                const paramsMap = [
                    { 'bk_obj_id': 'set', condition: [], fields: [] },
                    { 'bk_obj_id': 'module', condition: [], fields: [] },
                    {
                        'bk_obj_id': 'biz',
                        condition: [{
                            field: 'default', // 该参数表明查询非资源池下的主机
                            operator: '$ne',
                            value: 1
                        }],
                        fields: []
                    }, {
                        'bk_obj_id': 'host',
                        condition: [],
                        fields: this.previewHeader
                    }
                ]
                properties.forEach((property, index) => {
                    const param = paramsMap.find(({ bk_obj_id: objId }) => {
                        return objId === property.objId
                    })
                    if (property.value !== null && property.value !== undefined && String(property.value).length) {
                        if (property.propertyType === 'time' || property.propertyType === 'date') {
                            const value = property['value']
                            param['condition'].push({
                                field: property.propertyId,
                                operator: value[0] === value[1] ? '$eq' : '$gte',
                                value: value[0]
                            })
                            param['condition'].push({
                                field: property.propertyId,
                                operator: value[0] === value[1] ? '$eq' : '$lte',
                                value: value[1]
                            })
                        } else if (property.propertyType === 'bool' && ['true', 'false'].includes(property.value)) {
                            param['condition'].push({
                                field: property.propertyId,
                                operator: property.operator,
                                value: property.value === 'true'
                            })
                        } else if (property.operator === '$multilike') {
                            param.condition.push({
                                field: property.propertyId,
                                operator: property.operator,
                                value: property.value.split('\n').filter(str => str.trim().length).map(str => str.trim())
                            })
                        } else {
                            let operator = property.operator
                            let value = property.value
                            // 多模块与多集群查询
                            if (property.propertyId === 'bk_module_name' || property.propertyId === 'bk_set_name') {
                                operator = operator === '$regex' ? '$in' : operator
                                if (operator === '$in') {
                                    const arr = value.replace('，', ',').split(',')
                                    const isExist = arr.findIndex(val => {
                                        return val === value
                                    }) > -1
                                    value = isExist ? arr : [...arr, value]
                                }
                            }
                            param['condition'].push({
                                field: property.propertyId,
                                operator: operator,
                                value: value
                            })
                        }
                    }
                })
                const params = {
                    'bk_biz_id': this.bizId,
                    'info': {
                        condition: paramsMap
                    },
                    'name': row.name
                }
                return params
            },
            async getUserAPIList (value, event) {
                try {
                    const res = await this.searchCustomQuery({
                        bizId: this.bizId,
                        params: this.searchParams,
                        config: {
                            globalPermission: false,
                            requestId: 'searchCustomQuery'
                        }
                    })
                    if (res.count && !res.info.length) {
                        this.table.pagination.current -= 1
                        this.getUserAPIList()
                    }
                    if (res.count) {
                        this.table.list = res.info
                    } else {
                        this.table.list = []
                    }
                    this.table.pagination.count = res.count

                    if (event) {
                        this.table.stuff.type = 'search'
                    }
                } catch ({ permission }) {
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                }
            },
            showUserAPISlider (type) {
                this.slider.isShow = true
                this.slider.type = type
                this.slider.title = this.$t('新建动态分组')
            },
            /* 显示自定义API详情 */
            showUserAPIDetails (userAPI, event, column = {}) {
                if (column.property === 'operation') return
                this.slider.isShow = true
                this.slider.type = 'update'
                this.slider.id = userAPI['id']
                this.slider.title = this.$t('编辑动态分组')
            },
            deleteUserAPI (row) {
                this.$bkInfo({
                    title: this.$t('确定删除'),
                    subTitle: this.$t('确认要删除分组', { name: row.name }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        await this.$store.dispatch('hostCustomApi/deleteCustomQuery', {
                            bizId: this.bizId,
                            id: row.id,
                            config: {
                                requestId: 'deleteCustomQuery'
                            }
                        })
                        this.$success(this.$t('删除成功'))
                        this.getUserAPIList()
                        this.hideUserAPISlider()
                    }
                })
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.getUserAPIList()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.getUserAPIList()
            },
            async initObjectProperties () {
                const res = await Promise.all([
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_host',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'set',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_set',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'module',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_module',
                            fromCache: true
                        }
                    }),
                    this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'biz',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'post_searchObjectAttribute_biz',
                            fromCache: true
                        }
                    })
                ])
                this.object['host']['properties'] = res[0].filter(property => !property['bk_isapi'])
                this.object['set']['properties'] = res[1].filter(property => !property['bk_isapi'])
                this.object['module']['properties'] = res[2].filter(property => !property['bk_isapi'])
                this.object['biz']['properties'] = res[3].filter(property => !property['bk_isapi'])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .api-wrapper {
        padding: 0 20px;
        .filter-wrapper {
            .business-selector {
                float: left;
                width: 170px;
                margin-right: 10px;
            }
            .api-btn {
                float: left;
            }
            .api-input {
                float: right;
                width: 320px;
            }
        }
        .api-table {
            margin-top: 14px;
        }
    }
</style>
