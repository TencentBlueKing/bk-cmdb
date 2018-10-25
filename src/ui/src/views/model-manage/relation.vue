<template>
    <div class="relation-wrapper">
        <p class="operation-box">
            <bk-button type="primary" @click="slider.isShow = true">
                {{$t('ModelManagement["新增关联类型"]')}}
            </bk-button>
            <label class="search-input">
                <i class="bk-icon icon-search"></i>
                <input type="text" class="cmdb-form-input" v-model.trim="searchText" :placeholder="$t('ModelManagement[\'请输入关联类型名称\']')">
            </label>
        </p>
        <cmdb-table
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination">
        </cmdb-table>
        <cmdb-slider
            class="relation-slider"
            :width="514"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <the-relation slot="content" class="slider-content"></the-relation>
        </cmdb-slider>
    </div>
</template>

<script>
    import theRelation from './relation-type'
    import { mapActions } from 'vuex'
    export default {
        components: {
            theRelation
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    title: this.$t('ModelManagement["新增关联类型"]')
                },
                searchText: '',
                table: {
                    header: [{
                        id: 'bk_asst_name',
                        name: this.$t('Hosts["名称"]')
                    }, {
                        id: 'bk_asst_id',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'src_des',
                        name: this.$t('ModelManagement["源->目标描述"]')
                    }, {
                        id: 'dest_des',
                        name: this.$t('ModelManagement["目标描述->源"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["使用数"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Common["操作"]')
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
        created () {
            this.searchAssociationType()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType',
                'createAssociationType',
                'updateAssociationType',
                'deleteAssociationType'
            ])
        }
    }
</script>


<style lang="scss" scoped>
    .operation-box {
        margin: 20px 0;
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
