<template>
    <div class="relation-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            :feature-name="'association'"
            :show-tips="showFeatureTips"
            :desc="$t('ModelManagement[\'关联关系提示\']')"
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
                    {{$t('Common["新建"]')}}
                </bk-button>
            </span>
            <label class="search-input">
                <i class="bk-icon icon-search" @click="searchRelation(true)"></i>
                <bk-input type="text" class="cmdb-form-input"
                    v-model.trim="searchText"
                    :placeholder="$t('ModelManagement[\'请输入关联类型名称\']')"
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
            <bk-table-column prop="bk_asst_id" :label="$t('ModelManagement[\'唯一标识\']')"></bk-table-column>
            <bk-table-column prop="bk_asst_name" :label="$t('Hosts[\'名称\']')">
                <template slot-scope="{ row }">
                    {{row['bk_asst_name'] || '--'}}
                </template>
            </bk-table-column>
            <bk-table-column prop="src_des" :label="$t('ModelManagement[\'源->目标描述\']')"></bk-table-column>
            <bk-table-column prop="dest_des" :label="$t('ModelManagement[\'目标->源描述\']')"></bk-table-column>
            <bk-table-column prop="count" :label="$t('ModelManagement[\'使用数\']')"></bk-table-column>
            <bk-table-column v-if="isAdminView"
                fixed="right"
                :label="$t('Common[\'操作\']')">
                <template slot-scope="{ row }">
                    <span class="text-primary disabled mr10"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_RELATION),
                            auth: [$OPERATION.U_RELATION]
                        }"
                        v-if="row.ispre || !$isAuthorized($OPERATION.U_RELATION)">
                        {{$t('Common["编辑"]')}}
                    </span>
                    <span class="text-primary mr10"
                        v-else
                        @click.stop="editRelation(row)">
                        {{$t('Common["编辑"]')}}
                    </span>
                    <span class="text-primary disabled"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.D_RELATION),
                            auth: [$OPERATION.D_RELATION]
                        }"
                        v-if="row.ispre || !$isAuthorized($OPERATION.D_RELATION)">
                        {{$t('Common["删除"]')}}
                    </span>
                    <span class="text-primary"
                        v-else
                        @click.stop="deleteRelation(row)">
                        {{$t('Common["删除"]')}}
                    </span>
                </template>
            </bk-table-column>
        </bk-table>
        <bk-sideslider
            class="relation-slider"
            :width="450"
            :title="slider.title"
            :is-show.sync="slider.isShow">
            <the-relation
                slot="content"
                class="slider-content"
                v-if="slider.isShow"
                :is-edit="slider.isEdit"
                :relation="slider.relation"
                @saved="saveRelation"
                @cancel="slider.isShow = false">
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
                    title: this.$t('ModelManagement["新建关联类型"]'),
                    relation: {}
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
                this.slider.title = this.$t('ModelManagement["新建关联类型"]')
                this.slider.isEdit = false
                this.slider.isShow = true
            },
            editRelation (relation) {
                this.slider.title = this.$t('ModelManagement["编辑关联类型"]')
                this.slider.relation = relation
                this.slider.isEdit = true
                this.slider.isShow = true
            },
            deleteRelation (relation) {
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除关联类型？"]', relation['bk_asst_name'], { name: relation['bk_asst_name'] }),
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .operation-box {
        margin: 0 0 20px 0;
        font-size: 0;
        .create-btn {
            margin: 0 10px 0 0;
        }
        .search-input {
            position: relative;
            display: inline-block;
            width: 300px;
            .icon-search {
                position: absolute;
                top: 9px;
                right: 10px;
                font-size: 18px;
                color: $cmdbBorderColor;
            }
            .cmdb-form-input {
                vertical-align: middle;
                padding-right: 36px;
            }
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
