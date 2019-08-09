<template>
    <div class="details-wrapper">
        <div class="details-box">
            <div class="table-box">
                <p class="title clearfix">
                    <label class="label" @click="propertyTable.isShow = !propertyTable.isShow">
                        <i class="bk-icon icon-angle-down" :class="{ 'rotate': !propertyTable.isShow }"></i>
                        <span>{{$t('属性')}}</span>
                    </label>
                    <label class="cmdb-form-checkbox cmdb-checkbox-small">
                        <input type="checkbox" :disabled="ignore" v-model="propertyTable.isShowIgnore">
                        <span class="cmdb-checkbox-text">{{$t('显示忽略')}}</span>
                    </label>
                </p>
                <cmdb-collapse-transition>
                    <div v-show="propertyTable.isShow">
                        <cmdb-table
                            class="table"
                            :loading="$loading('searchNetcollectChangeDetail')"
                            :max-height="40 * propertyTableList.length + 40"
                            :header="propertyTable.header"
                            :list="propertyTableList"
                            :pagination.sync="propertyTable.pagination"
                            :default-sort="propertyTable.defaultSort"
                            @handleSortChange="propertyHandleSortChange">
                            <template v-for="(header, index) in propertyTable.header" :slot="header.id" slot-scope="{ item }">
                                <template v-if="header.id === 'isrequired'">
                                    <span :key="index" :class="{ 'disabled': item.method !== 'accept' }">
                                        {{item.isrequired ? $t('是') : $t('否')}}
                                    </span>
                                </template>
                                <template v-else-if="header.id === 'operation'">
                                    <span v-if="!item.isrequired" :key="index" class="text-primary" :class="{ 'disabled': ignore }" @click.stop="togglePropertyMethod(item)">{{item.method === 'accept' ? $t('忽略') : $t('取消忽略')}}</span>
                                </template>
                                <template v-else>
                                    <span :key="index" :class="{ 'disabled': item.method !== 'accept' }">{{item[header.id]}}</span>
                                </template>
                            </template>
                        </cmdb-table>
                    </div>
                </cmdb-collapse-transition>
            </div>
            <div class="table-box relation">
                <p class="title clearfix">
                    <label class="label" @click="relationTable.isShow = !relationTable.isShow">
                        <i class="bk-icon icon-angle-down" :class="{ 'rotate': !relationTable.isShow }"></i>
                        <span>{{$t('关系')}}</span>
                    </label>
                    <label class="cmdb-form-checkbox cmdb-checkbox-small">
                        <input type="checkbox" :disabled="ignore" v-model="relationTable.isShowIgnore">
                        <span class="cmdb-checkbox-text">{{$t('显示忽略')}}</span>
                    </label>
                </p>
                <cmdb-collapse-transition>
                    <div v-show="relationTable.isShow">
                        <cmdb-table
                            class="table"
                            :loading="$loading('searchNetcollectChangeDetail')"
                            :max-height="40 * relationTableList.length + 40"
                            :header="relationTable.header"
                            :list="relationTableList"
                            :pagination.sync="relationTable.pagination"
                            :default-sort="relationTable.defaultSort"
                            @handleSortChange="relationHandleSortChange">
                            <template v-for="(header, index) in relationTable.header" :slot="header.id" slot-scope="{ item }">
                                <template v-if="header.id === 'action'">
                                    <span :key="index" :class="{ 'color-danger': item.action === 'delete', 'disabled': item.method !== 'accept' }">{{actionMap[item.action]}}</span>
                                </template>
                                <template v-else-if="header.id === 'operation'">
                                    <span :key="index" class="text-primary" :class="{ 'disabled': ignore }" @click.stop="toggleRelationMethod(item)">{{item.method === 'accept' ? $t('忽略') : $t('取消忽略')}}</span>
                                </template>
                                <template v-else>
                                    <span :key="index" :class="{ 'disabled': item.method !== 'accept' }">{{item[header.id]}}</span>
                                </template>
                            </template>
                        </cmdb-table>
                    </div>
                </cmdb-collapse-transition>
            </div>
        </div>
        <div class="footer">
            <span>{{$t('忽略此实例')}}</span>
            <bk-switcher
                class="switcher"
                size="small"
                :show-text="false"
                :selected="ignore"
                @change="toggleSwitcher">
            </bk-switcher>
            <bk-button theme="default" :disabled="detailPage.prev" @click="updateView('prev')">
                {{$t('上一个')}}
            </bk-button>
            <bk-button theme="default" :disabled="detailPage.next" @click="updateView('next')">
                {{$t('下一个')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            attributes: {
                type: Array
            },
            associations: {
                type: Array
            },
            ignore: {
                type: Boolean
            },
            detailPage: {
                type: Object
            }
        },
        data () {
            return {
                isAccept: false,
                propertyTable: {
                    isShowIgnore: true,
                    isShow: true,
                    header: [{
                        id: 'bk_property_name',
                        name: this.$t('属性名')
                    }, {
                        id: 'isrequired',
                        name: this.$t('必须')
                    }, {
                        id: 'pre_value',
                        name: this.$t('原值')
                    }, {
                        id: 'value',
                        name: this.$t('新值')
                    }, {
                        id: 'operation',
                        name: this.$t('操作'),
                        sortable: false
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                },
                relationTable: {
                    isShowIgnore: true,
                    isShow: true,
                    header: [{
                        id: 'action',
                        name: this.$t('操作方式')
                    }, {
                        id: 'bk_asst_obj_name',
                        name: this.$t('模型')
                    }, {
                        id: 'configuration',
                        name: this.$t('配置信息')
                    }, {
                        id: 'operation',
                        name: this.$t('操作'),
                        sortable: false
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                },
                actionMap: {
                    'create': this.$t('新增关联'),
                    'delete': this.$t('删除关联')
                }
            }
        },
        computed: {
            propertyTableList () {
                return this.propertyTable.list.filter(item => {
                    if (!this.propertyTable.isShowIgnore && item.method !== 'accept') {
                        return false
                    }
                    return true
                })
            },
            relationTableList () {
                return this.relationTable.list.filter(item => {
                    if (!this.relationTable.isShowIgnore && item.method !== 'accept') {
                        return false
                    }
                    return true
                })
            }
        },
        created () {
            this.propertyTable.list = this.$tools.clone(this.attributes)
            this.relationTable.list = this.$tools.clone(this.associations)
        },
        methods: {
            updateView (type) {
                this.$emit('updateView', type)
                this.$nextTick(() => {
                    this.propertyTable.list = this.$tools.clone(this.attributes)
                    this.relationTable.list = this.$tools.clone(this.associations)
                })
            },
            toggleSwitcher (value) {
                this.$emit('toggleSwitcher', value)
            },
            togglePropertyMethod (item) {
                if (this.ignore) {
                    return
                }
                item.method = item.method === 'accept' ? 'reject' : 'accept'
                this.$emit('update:attributes', this.propertyTable.list)
            },
            toggleRelationMethod (item) {
                if (this.ignore) {
                    return
                }
                item.method = item.method === 'accept' ? 'reject' : 'accept'
                this.$emit('update:associations', this.relationTable.list)
            },
            propertyHandleSortChange (sort) {
                let key = sort
                if (sort[0] === '-') {
                    key = sort.substr(1, sort.length - 1)
                }
                this.propertyTable.list.sort((itemA, itemB) => {
                    if (itemA[key] === null) {
                        itemA[key] = ''
                    }
                    if (itemB[key] === null) {
                        itemB[key] = ''
                    }
                    return itemA[key].localeCompare(itemB[key])
                })
                if (sort[0] === '-') {
                    this.propertyTable.list.reverse()
                }
            },
            relationHandleSortChange (sort) {
                let key = sort
                if (sort[0] === '-') {
                    key = sort.substr(1, sort.length - 1)
                }
                this.relationTable.list.sort((itemA, itemB) => {
                    if (itemA[key] === null) {
                        itemA[key] = ''
                    }
                    if (itemB[key] === null) {
                        itemB[key] = ''
                    }
                    return itemA[key].localeCompare(itemB[key])
                })
                if (sort[0] === '-') {
                    this.relationTable.list.reverse()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-wrapper {
        padding: 15px 30px;
        height: calc(100% - 60px);
        .details-box {
            height: 100%;
            @include scrollbar;
        }
        .disabled {
            color: $cmdbBorderColor;
        }
        .table-box {
            &.relation {
                margin-top: 20px;
            }
            .title {
                line-height: 32px;
                >.cmdb-form-checkbox {
                    float: right;
                }
                .label {
                    cursor: pointer;
                }
                .icon-angle-down {
                    font-size: 12px;
                    font-weight: bold;
                    transition: all .2s;
                    &.rotate {
                        transform: rotate(180deg);
                    }
                }
            }
            .table {
                margin-top: 2px;
            }
        }
        .footer {
            position: absolute;
            bottom: 0;
            left: 0;
            padding: 12px 20px;
            width: 100%;
            font-size: 0;
            text-align: right;
            box-shadow: 0px -2px 5px 0px rgba(0, 0, 0, 0.05);
            >span {
                font-size: 14px;
                vertical-align: middle;
                margin-right: 8px;
            }
            .switcher {
                margin-right: 10px;
            }
            .bk-button {
                margin-left: 10px;
            }
        }
    }
</style>

<style lang="scss">
    .details-wrapper {
        .switcher {
            &.is-checked {
                background: $cmdbBorderFocusColor;
            }
        }
    }
</style>
