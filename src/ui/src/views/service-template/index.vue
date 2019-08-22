<template>
    <div class="template-wrapper" ref="templateWrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'serviceTemplate'"
            :show-tips="showFeatureTips"
            :desc="$t('服务模板功能提示')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="template-filter clearfix">
            <span class="fl mr10"
                v-cursor="{
                    active: !$isAuthorized($OPERATION.C_SERVICE_TEMPLATE),
                    auth: [$OPERATION.C_SERVICE_TEMPLATE]
                }">
                <bk-button
                    theme="primary"
                    :disabled="!$isAuthorized($OPERATION.C_SERVICE_TEMPLATE)"
                    @click="operationTemplate()">
                    {{$t('新建')}}
                </bk-button>
            </span>
            <div class="filter-text fr">
                <cmdb-selector
                    class="fl"
                    :placeholder="$t('所有一级分类')"
                    :auto-select="false"
                    :allow-clear="true"
                    :list="mainList"
                    v-model="filter['mainClassification']"
                    @on-selected="handleSelect">
                </cmdb-selector>
                <cmdb-selector
                    class="fl"
                    :placeholder="$t('所有二级分类')"
                    :auto-select="false"
                    :allow-clear="true"
                    :list="secondaryList"
                    v-model="filter['secondaryClassification']"
                    :empty-text="emptyText"
                    @on-selected="handleSelectSecondary">
                </cmdb-selector>
                <bk-input type="text"
                    class="filter-search fl"
                    :placeholder="$t('模板名称搜索')"
                    :right-icon="'bk-icon icon-search'"
                    v-model.trim="filter.templateName"
                    @enter="searchByTemplateName">
                </bk-input>
            </div>
        </div>
        <bk-table class="template-table"
            v-bkloading="{ isLoading: $loading('get_proc_service_template') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 210"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange">
            <bk-table-column prop="name" :label="$t('模板名称')"></bk-table-column>
            <bk-table-column prop="service_category" :label="$t('服务分类')"></bk-table-column>
            <bk-table-column prop="process_template_count" :label="$t('进程数量')"></bk-table-column>
            <bk-table-column prop="module_count" :label="$t('应用模块数')"></bk-table-column>
            <bk-table-column prop="modifier" :label="$t('修改人')"></bk-table-column>
            <bk-table-column prop="last_time" :label="$t('修改时间')">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm')}}
                </template>
            </bk-table-column>
            <bk-table-column prop="operation" :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_SERVICE_TEMPLATE),
                            auth: [$OPERATION.U_SERVICE_TEMPLATE]
                        }">
                        <bk-button class="mr10"
                            :disabled="!$isAuthorized($OPERATION.U_SERVICE_TEMPLATE)"
                            :text="true"
                            @click.stop="operationTemplate(row['id'])">
                            {{$t('编辑')}}
                        </bk-button>
                    </span>
                    <span class="text-primary"
                        style="color: #c4c6cc !important; cursor: not-allowed;"
                        v-if="row['module_count'] && $isAuthorized($OPERATION.D_SERVICE_TEMPLATE)"
                        v-bk-tooltips.top="$t('不可删除')">
                        {{$t('删除')}}
                    </span>
                    <span v-else
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.D_SERVICE_TEMPLATE),
                            auth: [$OPERATION.D_SERVICE_TEMPLATE]
                        }">
                        <bk-button
                            :disabled="!$isAuthorized($OPERATION.D_SERVICE_TEMPLATE)"
                            :text="true"
                            @click.stop="deleteTemplate(row)">
                            {{$t('删除')}}
                        </bk-button>
                    </span>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    export default {
        components: {
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
                filter: {
                    mainClassification: '',
                    secondaryClassification: '',
                    templateName: ''
                },
                table: {
                    height: 600,
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        limit: 10
                    },
                    defaultSort: '-last_time',
                    sort: '-id'
                },
                mainList: [],
                secondaryList: [],
                allSecondaryList: [],
                originTemplateData: [],
                maincategoryId: null,
                categoryId: null
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams']),
            params () {
                const id = this.categoryId
                    ? this.categoryId
                    : this.maincategoryId ? this.maincategoryId : 0
                return {
                    service_category_id: id,
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                        limit: this.table.pagination.limit,
                        sort: this.table.defaultSort
                    }
                }
            },
            emptyText () {
                return this.filter.mainClassification ? this.$t('没有二级分类') : this.$t('请选择一级分类')
            }
        },
        async created () {
            this.showFeatureTips = this.featureTipsParams['serviceTemplate']
            try {
                await this.getServiceClassification()
                await this.getTableData()
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('serviceTemplate', ['searchServiceTemplate', 'deleteServiceTemplate']),
            ...mapActions('serviceClassification', ['searchServiceCategory']),
            async getTableData () {
                const templateData = await this.getTemplateData()
                if (templateData.count && !templateData.info.length) {
                    this.table.pagination.current -= 1
                    this.getTableData()
                }
                this.table.pagination.count = templateData.count
                this.table.allList = templateData.info.map(template => {
                    const result = {
                        ...template,
                        ...template['service_template']
                    }
                    const secondaryCategory = this.allSecondaryList.find(classification => classification['id'] === result['service_category_id'])
                    const mainCategory = this.mainList.find(classification => secondaryCategory && classification['id'] === secondaryCategory['bk_parent_id'])
                    const secondaryCategoryName = secondaryCategory ? secondaryCategory['name'] : '--'
                    const mainCategoryName = mainCategory ? mainCategory['name'] : '--'
                    result['service_category'] = `${mainCategoryName} / ${secondaryCategoryName}`
                    return result
                })
                this.table.list = this.table.allList
            },
            getTemplateData () {
                return this.searchServiceTemplate({
                    params: this.$injectMetadata(this.params),
                    config: {
                        requestId: 'get_proc_service_template',
                        cancelPrevious: true
                    }
                })
            },
            async getServiceClassification () {
                const res = await this.searchServiceCategory({
                    params: this.$injectMetadata(),
                    config: {
                        requestId: 'get_proc_services_categories'
                    }
                })
                this.classificationList = res.info.map(item => item['category'])
                this.mainList = this.classificationList.filter(classification => !classification['bk_parent_id'])
                this.allSecondaryList = this.classificationList.filter(classification => classification['bk_parent_id'])
            },
            searchByTemplateName () {
                const reg = new RegExp(this.filter.templateName, 'gi')
                const filterList = this.table.allList.filter(template => reg.test(template['name']))
                this.table.list = this.filter.templateName ? filterList : this.table.allList
            },
            handleSelect (id, data) {
                this.secondaryList = this.allSecondaryList.filter(classification => classification['bk_parent_id'] === id)
                this.filter.secondaryClassification = ''
                this.maincategoryId = id
                this.getTableData()
            },
            handleSelectSecondary (id) {
                this.categoryId = id
                this.getTableData()
            },
            operationTemplate (id) {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: id
                    },
                    query: {
                        from: this.$route.fullPath
                    }
                })
            },
            deleteTemplate (template) {
                this.$bkInfo({
                    title: this.$t('确认删除模版'),
                    subTitle: this.$tc('即将删除服务模版', name, { name: template.name }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        await this.deleteServiceTemplate({
                            params: {
                                data: this.$injectMetadata({
                                    service_template_id: template.id
                                })
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
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
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
    .template-wrapper {
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
