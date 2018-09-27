<template>
    <div class="device-wrapper">
        <div class="title">
            <bk-button type="primary">
                {{$t('NetworkConfig["新增属性"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('Common["删除"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('ModelManagement["导入"]')}}
            </bk-button>
            <bk-button type="default">
                {{$t('ModelManagement["导出"]')}}
            </bk-button>
            <div class="filter">
                <bk-selector
                    class="search-selector"
                    :list="typeList"
                    :selected.sync="filter.type"
                ></bk-selector>
                <input class="cmdb-form-input" :placeholder="$t('Common[\'请输入\']')" type="text" id="SearchUserName" v-model.trim="filter.text" @keyup.enter="getRoleList">
                <i class="filter-search bk-icon icon-search"
                @click="getRoleList"></i>
            </div>
        </div>
        <cmdb-table
            class="device-table"
            :loading="$loading('searchSubscription')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort">
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                typeList: [],
                filter: {
                    type: '',
                    text: ''
                },
                table: {
                    header: [{
                        type: 'checkbox'
                    }, {
                        id: 'id',
                        name: 'ID'
                    }, {
                        id: 'device',
                        name: this.$t('NetworkConfig["所属设备"]')
                    }, {
                        id: 'type',
                        name: this.$t('NetworkConfig["计量单位"]')
                    }, {
                        id: 'type',
                        name: 'oid'
                    }, {
                        id: 'name',
                        name: this.$t('OperationAudit["模型"]')
                    }, {
                        id: 'model',
                        name: this.$t('NetworkConfig["模型属性"]')
                    }],
                    list: [{
                        change: 'asdf',
                        unique: 'asdf',
                        type: 'asdf',
                        ip: 'asdf',
                        config: 'aaa',
                        time: 'ddd'
                    }],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .device-wrapper {
        padding-top: 20px;
        .title {
            font-size: 0;
            >.bk-default {
                margin-left: 10px;
            }
            .filter {
                position: relative;
                float: right;
                height: 36px;
                line-height: 36px;
                .search-selector {
                    position: relative;
                    float: left;
                    width: 120px;
                    margin-right: -1px;
                    z-index: 1;
                }
                .cmdb-form-input {
                    position: relative;
                    width: 260px;
                    border-radius: 0 2px 2px 0;
                    &:focus {
                        z-index: 2;
                    }
                }
                .icon-search {
                    position: absolute;
                    right: 10px;
                    top: 11px;
                    cursor: pointer;
                    font-size: 14px;
                    z-index: 3;
                }
            }
        }
        .device-table {
            margin-top: 20px;
            background: #fff;
        }
    }
</style>
