<template>
    <div class="template-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'serviceTemplate'"
            :show-tips="showFeatureTips"
            :desc="$t('ServiceManagement[\'功能提示\']')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="template-filter clearfix">
            <bk-button class="fl mr10" type="primary" @click="createTemplate">{{$t("Common['新建']")}}</bk-button>
            <bk-button class="fl mr10">{{$t("ServiceManagement['批量删除']")}}</bk-button>
            <div class="filter-text fr">
                <bk-selector
                    class="fl"
                    placeholder="请选择一级分类"
                    :list="[]"
                    :selected.sync="filter.primaryClassification">
                </bk-selector>
                <bk-selector
                    class="fl"
                    placeholder="请选择二级分类"
                    :list="[]"
                    :selected.sync="filter.secondaryClassification">
                </bk-selector>
                <div class="filter-search fl">
                    <input type="text"
                        class="bk-form-input"
                        :placeholder="$t('ServiceManagement[\'模板名称\']')">
                    <i class="bk-icon icon-search"></i>
                </div>
            </div>
        </div>
        <cmdb-table class="template-table" ref="table"
            :loading="$loading('get_proc_service_template')"
            :checked.sync="table.checked"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            :wrapper-minus-height="300">
            <template slot="operation" slot-scope="{ item }">
                <button class="text-primary mr10"
                    @click.stop="editeTemplate(item)">
                    {{$t('Common["编辑"]')}}
                </button>
                <span class="text-primary" style="color: #c4c6cc !important;" v-bktooltips.top="$t('ServiceManagement[\'不可删除\']')">{{$t('Common["删除"]')}}</span>
                <!-- <button class="text-primary"
                    @click.stop="deleteTemplate(item)">
                    {{$t('Common["删除"]')}}
                </button> -->
            </template>
        </cmdb-table>
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
                    primaryClassification: '',
                    secondaryClassification: ''
                },
                table: {
                    header: [
                        {
                            id: 'id',
                            type: 'checkbox',
                            width: 50
                        }, {
                            id: 'name',
                            name: this.$t("ServiceManagement['模板名称']")
                        }, {
                            id: 'modifier',
                            name: this.$t("ServiceManagement['修改人']")
                        }, {
                            id: 'last_time',
                            name: this.$t("ServiceManagement['修改时间']")
                        }, {
                            id: 'operation',
                            name: this.$t('Common["操作"]'),
                            sortable: false
                        }
                    ],
                    checked: [],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-id',
                    sort: '-id'
                }
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t("Nav['服务模板']"))
            this.showFeatureTips = this.featureTipsParams['serviceTemplate']
            this.getTableData()
        },
        methods: {
            ...mapActions('serviceTemplate', ['searchServiceTemplate', 'deleteServiceTemplate']),
            getTableData () {
                this.searchServiceTemplate({
                    params: this.$injectMetadata({
                        service_category_id: null
                    }),
                    config: {
                        requestId: 'get_proc_service_template',
                        cancelPrevious: true
                    }
                }).then(data => {
                    this.table.list = data.info
                    this.table.pagination.count = data.count
                })
            },
            createTemplate () {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push({ name: 'createTemplate' })
            },
            deleteTemplate (template) {
                this.$bkInfo({
                    title: this.$t("ServiceManagement['确认删除模版']"),
                    content: this.$tc("ServiceManagement['即将删除服务模版']", name, { name: template.name }),
                    confirmFn: async () => {
                        await this.deleteServiceTemplate({
                            params: this.$injectMetadata({
                                service_category_id: template.id
                            }),
                            config: {
                                requestId: 'deleteServiceTemplate'
                            }
                        })
                        this.getTableData()
                    }
                })
            },
            editeTemplate (template) {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push({ name: 'updateTemplate' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-wrapper {
        .filter-text {
            .bk-selector {
                width: 180px;
                margin-right: 10px;
            }
            .filter-search {
                width: 200px;
                position: relative;
                .bk-form-input {
                    padding-right: 30px;
                }
                .icon-search {
                    position: absolute;
                    right: 10px;
                    top: 11px;
                    cursor: pointer;
                }
            }
        }
        .template-table {
            margin-top: 15px;
        }
    }
</style>
