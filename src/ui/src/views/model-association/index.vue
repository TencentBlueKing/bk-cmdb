<template>
    <div class="relation-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'association'"
            :show-tips="showFeatureTips"
            :desc="$t('关联关系提示')"
            :more-href="'https://docs.bk.tencent.com/cmdb/Introduction.html#%E6%A8%A1%E5%9E%8B%E5%85%B3%E8%81%94'"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <p class="operation-box">
            <span v-if="isAdminView" class="inline-block-middle"
                v-cursor="{
                    active: !$isAuthorized($OPERATION.C_RELATION),
                    auth: [$OPERATION.C_RELATION]
                }">
                <bk-button theme="primary"
                    class="create-btn"
                    :disabled="!$isAuthorized($OPERATION.C_RELATION)"
                    @click="createRelation">
                    {{$t('新建')}}
                </bk-button>
            </span>
            <label class="search-input">
                <!-- <i class="bk-icon icon-search" @click="searchRelation(true)"></i> -->
                <bk-input type="text" class="cmdb-form-input"
                    v-model.trim="searchText"
                    :right-icon="'bk-icon icon-search'"
                    :placeholder="$t('请输入关联类型名称')"
                    @enter="searchRelation(true)">
                </bk-input>
            </label>
        </p>
        <bk-table
            v-bkloading="{ isLoading: $loading('searchAssociationType') }"
            :data="table.list"
            :pagination="table.pagination"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @sort-change="handleSortChange">
            <bk-table-column prop="bk_asst_id" :label="$t('唯一标识')" sortable="custom" class-name="is-highlight">
                <template slot-scope="{ row }">
                    <div style="cursor: pointer; padding: 10px 0;" @click.stop="handleShowDetails(row)">
                        {{row['bk_asst_id']}}
                    </div>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_asst_name" :label="$t('名称')" sortable="custom">
                <template slot-scope="{ row }">
                    {{row['bk_asst_name'] || '--'}}
                </template>
            </bk-table-column>
            <bk-table-column prop="src_des" :label="$t('源->目标描述')" sortable="custom"></bk-table-column>
            <bk-table-column prop="dest_des" :label="$t('目标->源描述')" sortable="custom"></bk-table-column>
            <bk-table-column prop="count" :label="$t('使用数')"></bk-table-column>
            <bk-table-column v-if="isAdminView"
                fixed="right"
                :label="$t('操作')">
                <template slot-scope="{ row }">
                    <span class="text-primary disabled mr10"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_RELATION),
                            auth: [$OPERATION.U_RELATION]
                        }"
                        v-if="row.ispre || !$isAuthorized($OPERATION.U_RELATION)">
                        {{$t('编辑')}}
                    </span>
                    <span class="text-primary mr10"
                        v-else
                        @click.stop="editRelation(row)">
                        {{$t('编辑')}}
                    </span>
                    <span class="text-primary disabled"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.D_RELATION),
                            auth: [$OPERATION.D_RELATION]
                        }"
                        v-if="row.ispre || !$isAuthorized($OPERATION.D_RELATION)">
                        {{$t('删除')}}
                    </span>
                    <span class="text-primary"
                        v-else
                        @click.stop="deleteRelation(row)">
                        {{$t('删除')}}
                    </span>
                </template>
            </bk-table-column>
        </bk-table>
        <bk-sideslider
            class="relation-slider"
            :width="450"
            :title="slider.title"
            :is-show.sync="slider.isShow"
            :before-close="handleSliderBeforeClose">
            <the-relation
                ref="relationForm"
                slot="content"
                class="slider-content"
                v-if="slider.isShow"
                :is-edit="slider.isEdit"
                :is-read-only="slider.isReadOnly"
                :relation="slider.relation"
                @saved="saveRelation"
                @cancel="handleSliderBeforeClose">
            </the-relation>
        </bk-sideslider>
    </div>
</template>

<script>
    import featureTips from '@/components/feature-tips/index'
    import theRelation from './_detail'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            theRelation,
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
                slider: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('新建关联类型'),
                    relation: {},
                    isReadOnly: false
                },
                searchText: '',
                table: {
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        limit: 10
                    },
                    defaultSort: '-ispre',
                    sort: '-ispre'
                },
                sendSearchText: ''
            }
        },
        computed: {
            ...mapGetters(['isAdminView', 'featureTipsParams']),
            searchParams () {
                const params = {
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                        limit: this.table.pagination.limit,
                        sort: this.table.sort
                    }
                }
                if (this.sendSearchText.length) {
                    Object.assign(params, {
                        condition: {
                            bk_asst_name: {
                                '$regex': this.sendSearchText
                            }
                        }
                    })
                }
                return params
            }
        },
        created () {
            this.searchRelation()
            this.showFeatureTips = this.featureTipsParams['association']
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType',
                'deleteAssociationType',
                'searchAssociationListWithAssociationKindList'
            ]),
            searchRelation (fromClick) {
                if (fromClick) {
                    this.sendSearchText = this.searchText
                    this.table.pagination.current = 1
                }
                this.searchAssociationType({
                    params: this.searchParams,
                    config: {
                        requestId: 'searchAssociationType'
                    }
                }).then(data => {
                    if (data.count && !data.info.length) {
                        this.table.pagination.current -= 1
                        this.searchRelation()
                    }
                    this.table.list = data.info
                    this.searchUsageCount()
                    this.table.pagination.count = data.count
                    this.$http.cancel('post_searchAssociationType')
                })
            },
            async searchUsageCount () {
                const asstIds = []
                this.table.list.forEach(({ bk_asst_id: asstId }) => asstIds.push(asstId))
                const res = await this.searchAssociationListWithAssociationKindList({
                    params: {
                        asst_ids: asstIds
                    }
                })
                this.table.list.forEach(item => {
                    const asst = res.associations.find(({ bk_asst_id: asstId }) => asstId === item['bk_asst_id'])
                    if (asst) {
                        this.$set(item, 'count', asst.assts.length)
                    }
                })
                this.table.list.splice()
            },
            createRelation () {
                this.slider.title = this.$t('新建关联类型')
                this.slider.isReadOnly = false
                this.slider.isEdit = false
                this.slider.isShow = true
            },
            editRelation (relation) {
                this.slider.title = this.$t('编辑关联类型')
                this.slider.isReadOnly = false
                this.slider.relation = relation
                this.slider.isEdit = true
                this.slider.isShow = true
            },
            deleteRelation (relation) {
                this.$bkInfo({
                    title: this.$tc('确定删除关联类型？', relation['bk_asst_name'], { name: relation['bk_asst_name'] }),
                    confirmFn: async () => {
                        await this.deleteAssociationType({
                            id: relation.id,
                            config: {
                                requestId: 'deleteAssociationType'
                            }
                        })
                        this.searchRelation()
                    }
                })
            },
            saveRelation () {
                this.slider.isShow = false
                this.searchRelation()
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.searchRelation()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.searchRelation()
            },
            handleSliderBeforeClose () {
                const hasChanged = Object.keys(this.$refs.relationForm.changedValues).length
                if (hasChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.slider.isShow = false
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.slider.isShow = false
                return true
            },
            handleShowDetails (relation) {
                this.slider.title = this.$t('关联类型详情')
                this.slider.relation = relation
                this.slider.isReadOnly = true
                this.slider.isEdit = true
                this.slider.isShow = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .operation-box {
        margin: 0 0 14px 0;
        font-size: 0;
        .create-btn {
            margin: 0 10px 0 0;
        }
        .search-input {
            position: relative;
            display: inline-block;
            width: 300px;
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
