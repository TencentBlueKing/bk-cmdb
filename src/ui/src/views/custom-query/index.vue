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
            <span class="inline-block-middle" v-cursor="{
                active: !$isAuthorized($OPERATION.C_CUSTOM_QUERY),
                auth: [$OPERATION.C_CUSTOM_QUERY]
            }">
                <bk-button theme="primary" class="api-btn"
                    :disabled="!$isAuthorized($OPERATION.C_CUSTOM_QUERY)"
                    @click="showUserAPISlider('create')">
                    {{$t('新建')}}
                </bk-button>
            </span>
            <div class="api-input fr">
                <bk-input type="text" class="cmdb-form-input"
                    v-model="filter.name"
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
            :max-height="$APP.height - 220"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange"
            @row-click="showUserAPIDetails">
            <!-- <bk-table-column type="selection" width="60" align="center" fixed class-name="bk-table-selection"></bk-table-column> -->
            <bk-table-column prop="id" label="ID" class-name="is-highlight" fixed></bk-table-column>
            <bk-table-column prop="name" :label="$t('查询名称')" sortable="custom" fixed></bk-table-column>
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
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_CUSTOM_QUERY),
                            auth: [$OPERATION.U_CUSTOM_QUERY]
                        }">
                        <bk-button class="mr10"
                            :disabled="!$isAuthorized($OPERATION.U_CUSTOM_QUERY)"
                            :text="true"
                            @click.stop="showUserAPIDetails(row)">
                            {{$t('编辑')}}
                        </bk-button>
                    </span>
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.D_CUSTOM_QUERY),
                            auth: [$OPERATION.D_CUSTOM_QUERY]
                        }">
                        <bk-button
                            :disabled="!$isAuthorized($OPERATION.D_CUSTOM_QUERY)"
                            :text="true"
                            @click.stop="deleteUserAPI(row)">
                            {{$t('删除')}}
                        </bk-button>
                    </span>
                </template>
            </bk-table-column>
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
                @create="handleCreate"
                @update="getUserAPIList"
                @cancel="handleSliderBeforeClose">
            </v-define>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import vDefine from './define'
    export default {
        components: {
            vDefine,
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
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
                        limit: 10
                    }
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    type: 'create',
                    id: null,
                    title: this.$t('新建查询')
                }
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams']),
            ...mapGetters('objectBiz', ['bizId']),
            searchParams () {
                const params = {
                    start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                    limit: this.table.pagination.limit,
                    sort: this.table.sort
                }
                this.filter.name ? params['condition'] = { 'name': this.filter.name } : void (0)
                return params
            },
            editable () {
                if (this.type === 'update') {
                    return this.$isAuthorized(this.$OPERATION.D_CUSTOM_QUERY)
                }
                return true
            }
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['customQuery']
            this.getUserAPIList()
        },
        methods: {
            ...mapActions('hostCustomApi', [
                'searchCustomQuery'
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
            handleCreate (data) {
                this.slider.id = data['id']
                this.slider.type = 'update'
                this.slider.title = this.$t('编辑查询')
                this.handlePageChange(1)
            },
            async getUserAPIList () {
                const res = await this.searchCustomQuery({
                    bizId: this.bizId,
                    params: this.searchParams,
                    config: {
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
            },
            showUserAPISlider (type) {
                this.slider.isShow = true
                this.slider.type = type
                this.slider.title = this.$t('新建查询')
            },
            /* 显示自定义API详情 */
            showUserAPIDetails (userAPI, event, column = {}) {
                if (column.property === 'operation') return
                this.slider.isShow = true
                this.slider.type = 'update'
                this.slider.id = userAPI['id']
                this.slider.title = this.$t('编辑查询')
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .api-wrapper {
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
