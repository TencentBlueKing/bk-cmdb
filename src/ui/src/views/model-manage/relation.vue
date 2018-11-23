<template>
    <div class="relation-wrapper">
        <p class="operation-box">
            <bk-button type="primary" @click="createRelation">
                {{$t('ModelManagement["新增关联类型"]')}}
            </bk-button>
            <label class="search-input">
                <i class="bk-icon icon-search" @click="searchRelation"></i>
                <input type="text" class="cmdb-form-input" v-model.trim="searchText" :placeholder="$t('ModelManagement[\'请输入关联类型名称\']')" @keyup.enter="searchRelation">
            </label>
        </p>
        <cmdb-table
            :loading="$loading('searchAssociationType')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
            <template slot="bk_asst_name" slot-scope="{ item }">
                {{item['bk_asst_name'] || '--'}}
            </template>
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary mr10" @click.stop="editRelation(item)">
                    {{$t('Common["编辑"]')}}
                </span>
                <span class="text-primary" v-if="!item.ispre && !isReadOnly" @click.stop="deleteRelation(item)">
                    {{$t('Common["删除"]')}}
                </span>
            </template>
        </cmdb-table>
        <cmdb-slider
            class="relation-slider"
            :width="514"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <the-relation
                slot="content"
                class="slider-content"
                :isEdit="slider.isEdit"
                :relation="slider.relation"
                @saved="saveRelation"
                @cancel="slider.isShow = false"
            ></the-relation>
        </cmdb-slider>
    </div>
</template>

<script>
    import theRelation from './relation-type'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            theRelation
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('ModelManagement["新增关联类型"]'),
                    relation: {}
                },
                searchText: '',
                table: {
                    header: [{
                        id: 'bk_asst_id',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'bk_asst_name',
                        name: this.$t('Hosts["名称"]')
                    }, {
                        id: 'src_des',
                        name: this.$t('ModelManagement["源->目标描述"]')
                    }, {
                        id: 'dest_des',
                        name: this.$t('ModelManagement["目标->源描述"]')
                    }, {
                        id: 'count',
                        name: this.$t('ModelManagement["使用数"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Common["操作"]'),
                        sortable: false
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    defaultSort: '-bk_asst_id',
                    sort: '-bk_asst_id'
                }
            }
        },
        created () {
            this.searchRelation()
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            searchParams () {
                let params = {
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    }
                }
                if (this.searchText.length) {
                    Object.assign(params, {
                        condition: {
                            bk_asst_name: {
                                '$regex': this.searchText
                            }
                        }
                    })
                }
                return params
            }
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType',
                'deleteAssociationType',
                'searchAssociationListWithAssociationKindList'
            ]),
            searchRelation () {
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
                let asstIds = []
                this.table.list.forEach(({bk_asst_id: asstId}) => asstIds.push(asstId))
                const res = await this.searchAssociationListWithAssociationKindList({
                    params: {
                        asst_ids: asstIds
                    }
                })
                this.table.list.forEach(item => {
                    let asst = res.associations.find(({bk_asst_id: asstId}) => asstId === item['bk_asst_id'])
                    if (asst) {
                        this.$set(item, 'count', asst.assts.length)
                    }
                })
                this.table.list.splice()
            },
            createRelation () {
                this.slider.title = this.$t('ModelManagement["新增关联类型"]')
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
                    title: this.$tc('ModelManagement["确定删除关联类型？"]', relation['bk_asst_name'], {name: relation['bk_asst_name']}),
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
                this.pagination.current = current
                this.searchRelation()
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.sort = sort
                this.searchRelation()
            }
        }
    }
</script>


<style lang="scss" scoped>
    .operation-box {
        margin: 0 0 20px 0;
        font-size: 0;
        .search-input {
            position: relative;
            display: inline-block;
            margin-left: 10px;
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
