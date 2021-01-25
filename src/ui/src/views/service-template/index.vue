<template>
    <div class="template-wrapper" ref="templateWrapper">
        <cmdb-tips class="mb10 top-tips" tips-key="serviceTemplateTips">
            <i18n path="服务模板功能提示">
                <a class="tips-link" href="javascript:void(0)" @click="handleTipsLinkClick" place="link">{{$t('业务拓扑')}}</a>
            </i18n>
        </cmdb-tips>
        <div class="template-filter clearfix">
            <cmdb-auth class="fl mr10" :auth="{ type: $OPERATION.C_SERVICE_TEMPLATE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="disabled"
                    @click="operationTemplate()">
                    {{$t('新建')}}
                </bk-button>
            </cmdb-auth>
            <div class="filter-text fr">
                <bk-select class="fl"
                    font-size="medium"
                    :placeholder="$t('所有一级分类')"
                    :allow-clear="true"
                    :searchable="true"
                    v-model="filter.mainClassification"
                    @selected="handleSelect"
                    @clear="() => handleSelect()">
                    <bk-option v-for="category in mainList"
                        :key="category.id"
                        :id="category.id"
                        :name="category.name">
                    </bk-option>
                </bk-select>
                <bk-select class="fl"
                    font-size="medium"
                    :placeholder="$t('所有二级分类')"
                    :allow-clear="true"
                    :searchable="true"
                    :empty-text="emptyText"
                    v-model="filter.secondaryClassification"
                    @selected="handleSelectSecondary"
                    @clear="() => handleSelectSecondary()">
                    <bk-option v-for="category in secondaryList"
                        :key="category.id"
                        :id="category.id"
                        :name="category.name">
                    </bk-option>
                </bk-select>
                <bk-input type="text"
                    class="filter-search fl"
                    :placeholder="$t('请输入xx', { name: $t('模板名称') })"
                    :right-icon="'bk-icon icon-search'"
                    clearable
                    font-size="medium"
                    v-model.trim="filter.templateName"
                    @enter="getTableData(true)"
                    @clear="handlePageChange(1)">
                </bk-input>
            </div>
        </div>
        <bk-table class="template-table"
            v-bkloading="{ isLoading: $loading(request.list) }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 229"
            :row-style="{ cursor: 'pointer' }"
            @row-click="handleRowClick"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange"
            @sort-change="handleSortChange">
            <bk-table-column prop="id" label="ID" class-name="is-highlight" show-overflow-tooltip sortable="custom"></bk-table-column>
            <bk-table-column prop="name" :label="$t('模板名称')" show-overflow-tooltip sortable="custom"></bk-table-column>
            <bk-table-column prop="service_category" :label="$t('服务分类')" show-overflow-tooltip></bk-table-column>
            <bk-table-column prop="process_template_count" :label="$t('进程数量')">
                <template slot-scope="{ row }">
                    <cmdb-loading :loading="$loading(request.count)">
                        <template v-if="row.process_template_count > 0">
                            {{row.process_template_count}}
                        </template>
                        <span style="color: #ff9c01" v-else>{{row.process_template_count}}（{{$t('未配置')}}）</span>
                    </cmdb-loading>
                </template>
            </bk-table-column>
            <bk-table-column prop="module_count" :label="$t('已应用模块数')">
                <template slot-scope="{ row }">
                    <cmdb-loading :loading="$loading(request.count)">{{row.module_count}}</cmdb-loading>
                </template>
            </bk-table-column>
            <bk-table-column prop="modifier" :label="$t('修改人')" sortable="custom"></bk-table-column>
            <bk-table-column prop="last_time" :label="$t('修改时间')" show-overflow-tooltip sortable="custom">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm')}}
                </template>
            </bk-table-column>
            <bk-table-column prop="operation" :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <cmdb-loading :loading="$loading(request.count)">
                        <!-- 与查询详情功能重复暂去掉 -->
                        <!-- <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SERVICE_TEMPLATE, relation: [bizId, row.id] }">
                            <bk-button slot-scope="{ disabled }"
                                theme="primary"
                                :disabled="disabled"
                                :text="true"
                                @click.stop="operationTemplate(row['id'], 'edit')">
                                {{$t('编辑')}}
                            </bk-button>
                        </cmdb-auth> -->
                        <cmdb-auth :auth="{ type: $OPERATION.D_SERVICE_TEMPLATE, relation: [bizId, row.id] }">
                            <template slot-scope="{ disabled }">
                                <span class="text-primary"
                                    style="color: #dcdee5 !important; cursor: not-allowed;"
                                    v-if="row['module_count'] && !disabled"
                                    v-bk-tooltips.top="$t('不可删除')">
                                    {{$t('删除')}}
                                </span>
                                <bk-button v-else
                                    theme="primary"
                                    :disabled="disabled"
                                    :text="true"
                                    @click.stop="deleteTemplate(row)">
                                    {{$t('删除')}}
                                </bk-button>
                            </template>
                        </cmdb-auth>
                    </cmdb-loading>
                </template>
            </bk-table-column>
            <cmdb-table-empty
                slot="empty"
                :stuff="table.stuff"
                :auth="{ type: $OPERATION.C_SERVICE_TEMPLATE, relation: [bizId] }"
                @create="operationTemplate"
            ></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import CmdbLoading from '@/components/loading/loading'
    export default {
        components: {
            CmdbLoading
        },
        data () {
            return {
                filter: {
                    mainClassification: '',
                    secondaryClassification: '',
                    templateName: ''
                },
                table: {
                    height: 600,
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    sort: '-id',
                    stuff: {
                        type: 'default',
                        payload: {
                            resource: this.$t('服务模板')
                        }
                    }
                },
                mainList: [],
                secondaryList: [],
                allSecondaryList: [],
                originTemplateData: [],
                maincategoryId: null,
                categoryId: null,
                request: {
                    list: Symbol('list'),
                    count: Symbol('count')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            params () {
                const id = this.categoryId
                    ? this.categoryId
                    : this.maincategoryId ? this.maincategoryId : 0
                return {
                    bk_biz_id: this.bizId,
                    service_category_id: id,
                    search: this.filter.templateName,
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                        limit: this.table.pagination.limit,
                        sort: this.table.sort
                    }
                }
            },
            emptyText () {
                return this.filter.mainClassification ? this.$t('没有二级分类') : this.$t('请选择一级分类')
            },
            hasFilter () {
                return Object.values(this.filter).some(value => !!value)
            }
        },
        async created () {
            try {
                await this.getServiceClassification()
                await this.getTableData()
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('serviceTemplate', ['searchServiceTemplate', 'deleteServiceTemplate']),
            ...mapActions('serviceClassification', ['searchServiceCategoryWithoutAmout']),
            async getTableData (event) {
                try {
                    const templateData = await this.getTemplateData()
                    if (templateData.count && !templateData.info.length) {
                        this.table.pagination.current -= 1
                        this.getTableData()
                    }
                    this.table.pagination.count = templateData.count
                    this.table.list = templateData.info.map(template => {
                        const secondaryCategory = this.allSecondaryList.find(classification => classification.id === template.service_category_id)
                        const mainCategory = this.mainList.find(classification => secondaryCategory && classification.id === secondaryCategory.bk_parent_id)
                        const secondaryCategoryName = secondaryCategory ? secondaryCategory.name : '--'
                        const mainCategoryName = mainCategory ? mainCategory.name : '--'
                        template.service_category = `${mainCategoryName} / ${secondaryCategoryName}`
                        return template
                    })
                    this.table.stuff.type = this.hasFilter ? 'search' : 'default'
                    this.table.list.length && this.getTemplateCount()
                } catch ({ permission }) {
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                }
            },
            async getTemplateCount () {
                try {
                    const data = await this.$store.dispatch('serviceTemplate/searchServiceTemplateCount', {
                        bizId: this.bizId,
                        params: {
                            service_template_ids: this.table.list.map(row => row.id)
                        },
                        config: {
                            requestId: this.request.count,
                            cancelPrevious: true
                        }
                    })
                    this.table.list.forEach(row => {
                        const counts = data.find(counts => counts.service_template_id === row.id) || {}
                        const {
                            module_count: moduleCount = '--',
                            process_template_count: processTemplateCount = '--'
                        } = counts
                        this.$set(row, 'module_count', moduleCount)
                        this.$set(row, 'process_template_count', processTemplateCount)
                    })
                } catch (error) {
                    console.error(error)
                    this.table.list.forEach(row => {
                        this.$set(row, 'module_count', '--')
                        this.$set(row, 'process_template_count', '--')
                    })
                }
            },
            getTemplateData () {
                return this.searchServiceTemplate({
                    params: this.params,
                    config: {
                        requestId: this.request.list,
                        cancelPrevious: true,
                        globalPermission: false
                    }
                })
            },
            async getServiceClassification () {
                const { info: categories } = await this.searchServiceCategoryWithoutAmout({
                    params: { bk_biz_id: this.bizId },
                    config: {
                        requestId: 'get_proc_services_categories'
                    }
                })
                this.classificationList = categories
                this.mainList = this.classificationList.filter(classification => !classification['bk_parent_id'])
                this.allSecondaryList = this.classificationList.filter(classification => classification['bk_parent_id'])
            },
            handleSelect (id = '') {
                this.secondaryList = this.allSecondaryList.filter(classification => classification['bk_parent_id'] === id)
                this.maincategoryId = id
                this.handleSelectSecondary()
            },
            handleSelectSecondary (id = '') {
                this.categoryId = id
                this.filter.secondaryClassification = id
                this.getTableData(true)
            },
            operationTemplate (id, type) {
                this.$routerActions.redirect({
                    name: 'operationalTemplate',
                    params: {
                        templateId: id,
                        isEdit: type === 'edit'
                    },
                    history: true
                })
            },
            deleteTemplate (template) {
                this.$bkInfo({
                    title: this.$t('确认删除模板'),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        await this.deleteServiceTemplate({
                            params: {
                                data: {
                                    bk_biz_id: this.bizId,
                                    service_template_id: template.id
                                }
                            },
                            config: {
                                requestId: 'delete_proc_service_template'
                            }
                        }).then(() => {
                            this.$success(this.$t('删除成功'))
                            this.getTableData()
                        })
                    }
                })
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort, '-id')
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            },
            handleRowClick (row, event, column) {
                if (column.property === 'operation') return
                this.operationTemplate(row.id)
            },
            handleTipsLinkClick () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-wrapper {
        padding: 15px 20px 0;
        .tips-link {
            color: #3A84FF;
            margin: 0;
        }
        .filter-text {
            .bk-select {
                width: 184px;
                margin-right: 10px;
            }
            .filter-search {
                width: 210px;
                position: relative;
            }
        }
        .template-table {
            margin-top: 14px;
        }
    }
</style>
