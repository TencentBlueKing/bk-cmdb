<template>
    <div>
        <bk-button class="create-btn" type="primary" @click="createVerification">
            {{$t('ModelManagement["新建校验"]')}}
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
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <the-verification-detail
                class="slider-content"
                slot="content"
                :isReadOnly="isReadOnly"
                :isEditField="slider.isEdit"
                :field="slider.verification"
                @save="saveVerification"
                @cancel="slider.isShow = false">
            </the-verification-detail>
        </cmdb-slider>
    </div>
</template>

<script>
    import theVerificationDetail from './verification-detail'
    export default {
        components: {
            theVerificationDetail
        },
        props: {
            isReadOnly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    isEdit: false,
                    verification: {}
                },
                table: {
                    header: [{
                        id: 'name',
                        name: this.$t('ModelManagement["校验规则"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["优先级"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["是否显示为实例名称"]')
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
                    defaultSort: '-op_time',
                    sort: '-op_time'
                }
            }
        },
        methods: {
            createVerification () {
                this.slider.title = this.$t('ModelManagement["新增校验"]')
                this.slider.isEdit = false
                this.slider.isShow = true
            },
            editVerification (verification) {
                this.slider.title = this.$t('ModelManagement["编辑校验"]')
                this.slider.verification = verification
                this.slider.isEdit = true
                this.slider.isShow = true
            },
            saveVerification () {

            },
            searchVerification () {

            },
            handlePageChange (current) {
                this.pagination.current = current
                this.searchVerification()
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.sort = sort
                this.searchVerification()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
</style>

