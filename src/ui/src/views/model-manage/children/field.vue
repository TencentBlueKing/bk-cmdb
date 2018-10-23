<template>
    <div class="model-field-wrapper">
        <bk-button class="create-btn" type="primary" @click="showSlider">
            {{$t('ModelManagement["新建字段"]')}}
        </bk-button>
        <cmdb-table
            class="field-table"
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handleSortChange="handleSortChange">
        </cmdb-table>
        <cmdb-slider
            :width="514"
            :isShow.sync="slider.isShow"
        >
            <the-field-detail class="slider-content" slot="content"></the-field-detail>
        </cmdb-slider>
    </div>
</template>

<script>
    import theFieldDetail from './field-detail'
    import { mapActions } from 'vuex'
    export default {
        components: {
            theFieldDetail
        },
        data () {
            return {
                slider: {
                    isShow: false
                },
                table: {
                    header: [{
                        id: 'name',
                        name: this.$t('ModelManagement["字段名称"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["类型"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["唯一"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["必填"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["字段描述"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('ModelManagement["创建时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('ModelManagement["操作"]')
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
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            showSlider () {
                this.slider.isShow = true
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
    .slider-content {
        padding: 20px;
    }
</style>
