<template>
    <div>
        <bk-button class="create-btn" type="primary" @click="createRelation">
            {{$t('ModelManagement["新建关联关系"]')}}
        </bk-button>
        <cmdb-table
            class="relation-table"
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
        </cmdb-table>
        <cmdb-slider
            :width="514"
            :isShow.sync="slider.isShow">
            <the-relation-detail
                class="slider-content"
                slot="content"
                :isReadOnly="isReadOnly"
                :isEditField="slider.isEditField"
                :field="slider.curField"
                @save="saveRelation"
                @cancel="slider.isShow = false">
            </the-relation-detail>
        </cmdb-slider>
    </div>
</template>

<script>
    import theRelationDetail from './relation-detail'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            theRelationDetail
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    isEdit: false,
                    relation: {}
                },
                table: {
                    header: [{
                        id: 'name',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["关联类型"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["约束"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["源模型"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('ModelManagement["目标模型"]')
                    }, {
                        id: 'operation',
                        name: this.$t('ModelManagement["操作"]'),
                        sortable: false
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time'
                }
            }
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
            }
        },
        methods: {
            createRelation () {
                this.slider.isEdit = false
                this.slider.isReadOnly = false
                this.slider.relation = {}
                this.slider.title = this.$t('ModelManagement["新建关联关系"]')
                this.slider.isShow = true
            },
            saveRelation () {

            },
            handlePageChange (current) {
                this.pagination.current = current
                this.refresh()
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.sort = sort
                this.refresh()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
</style>
