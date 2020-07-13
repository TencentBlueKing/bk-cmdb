<template>
    <div class="vpc-selector clearfix">
        <div class="left-column fl">
            <h2 class="left-title">{{$t('添加VPC')}}</h2>
            <task-region-selector class="region-selector"
                :account="account"
                v-model="currentRegion">
            </task-region-selector>
            <div class="vpc-wrapper" v-bkloading="{ isLoading: $loading([request.vpc, request.region]) }">
                <template v-if="hasVpc">
                    <div class="vpc-options">
                        <span class="option-name">{{$t('VPC名称')}}</span>
                        <div class="option-right">
                            <span class="option-count">{{$t('主机数')}}</span>
                            <bk-checkbox class="option-checkbox"
                                :checked="selectState.all"
                                :indeterminate="selectState.indeterminate"
                                @change="handleToggleSelectAll">
                                {{$t('全选')}}
                            </bk-checkbox>
                        </div>
                    </div>
                    <ul class="vpc-list">
                        <li class="vpc-item"
                            v-for="vpc in regionVPC[currentRegion]"
                            :key="vpc.bk_vpc_id">
                            <span class="vpc-item-name" v-bk-overflow-tips>{{getVpcName(vpc)}}</span>
                            <div class="vpc-item-option">
                                <span :class="['vpc-item-count', $i18n.locale]">{{vpc.bk_host_count}}</span>
                                <bk-checkbox class="vpc-item-checkbox"
                                    :checked="selectionMap.hasOwnProperty(vpc.bk_vpc_id)"
                                    @change="handleVpcChange(vpc, ...arguments)">
                                </bk-checkbox>
                            </div>
                        </li>
                    </ul>
                </template>
                <p class="vpc-empty" v-else>{{$t('该地域暂无VPC')}}</p>
            </div>
        </div>
        <div class="right-column fl" v-bkloading="{ isLoading: $loading(request.region) }">
            <div class="right-options">
                <i18n class="selection-info" path="已选择">
                    <span class="selection-count" place="number">{{selection.length}}</span>
                </i18n>
                <bk-link class="selection-clear" theme="primary" @click="handleClear">{{$t('清空')}}</bk-link>
            </div>
            <bk-table class="selected-vpc-table"
                :height="368"
                :data="selection"
                :outer-border="false"
                :header-border="false"
                :row-class-name="getRowClass">
                <bk-table-column label="VPC" prop="bk_vpc_id" width="200">
                    <template slot-scope="{ row }">
                        <div class="vpc-info">
                            <span class="vpc-name" v-bk-overflow-tips>{{getVpcName(row)}}</span>
                            <span class="info-destroyed"
                                v-if="row.destroyed"
                                v-bk-tooltips="$t('VPC已销毁')">
                                {{$t('已失效')}}
                            </span>
                        </div>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('地域')" prop="bk_region_name" show-overflow-tooltip>
                    <task-region-selector slot-scope="{ row }"
                        display="info"
                        :value="row.bk_region"
                        :account="account">
                    </task-region-selector>
                </bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="bk_host_count" width="80" align="right"></bk-table-column>
                <bk-table-column :label="$t('操作')" width="80">
                    <template slot-scope="{ row }">
                        <bk-link theme="primary" @click="handleRemoveSelection(row)">{{$t('移除')}}</bk-link>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="clearfix"></div>
        <div class="bottom-options">
            <span v-bk-tooltips="{
                content: $t('请至少选择一个VPC'),
                disabled: !!selection.length
            }">
                <bk-button class="mr10" theme="primary"
                    :disabled="!selection.length"
                    @click="handleConfirm">
                    {{$t('确定')}}
                </bk-button>
            </span>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import TaskRegionSelector from './task-region-selector.vue'
    export default {
        name: 'task-vpc-selector',
        components: {
            [TaskRegionSelector.name]: TaskRegionSelector
        },
        props: {
            account: {
                type: Number,
                required: true
            },
            selected: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                currentRegion: '',
                regionVPC: {},
                selection: [...this.selected],
                request: {
                    region: Symbol('region'),
                    vpc: Symbol('region')
                }
            }
        },
        computed: {
            selectState () {
                const state = {}
                const vpcList = this.regionVPC[this.currentRegion]
                const selected = vpcList.filter(vpc => this.selectionMap.hasOwnProperty(vpc.bk_vpc_id))
                state.indeterminate = selected.length > 0 && selected.length < vpcList.length
                state.all = selected.length === vpcList.length
                return state
            },
            hasVpc () {
                return this.currentRegion && (this.regionVPC[this.currentRegion] || []).length
            },
            selectionMap () {
                const map = {}
                this.selection.forEach(vpc => {
                    map[vpc.bk_vpc_id] = vpc
                })
                return map
            }
        },
        watch: {
            currentRegion () {
                this.getVpcList()
            }
        },
        methods: {
            async getVpcList () {
                try {
                    if (this.regionVPC.hasOwnProperty(this.currentRegion)) {
                        return this.regionVPC[this.currentRegion]
                    }
                    const { info } = await this.$store.dispatch('cloud/resource/findVPC', {
                        id: this.account,
                        params: {
                            bk_account_id: this.account,
                            bk_region: this.currentRegion
                        },
                        config: {
                            requestId: this.request.vpc
                        }
                    })
                    this.$set(this.regionVPC, this.currentRegion, info)
                } catch (e) {
                    console.error(e)
                }
            },
            handleToggleSelectAll (checked) {
                const vpcList = this.regionVPC[this.currentRegion]
                if (checked) {
                    const appendVpc = vpcList.filter(vpc => !this.selectionMap.hasOwnProperty(vpc.bk_vpc_id))
                    this.selection.unshift(...appendVpc)
                } else {
                    this.selection = this.selection.filter(exist => !vpcList.some(vpc => exist.bk_vpc_id === vpc.bk_vpc_id))
                }
            },
            handleVpcChange (vpc, checked) {
                if (checked) {
                    this.selection.unshift(vpc)
                } else {
                    const index = this.selection.findIndex(exist => exist.bk_vpc_id === vpc.bk_vpc_id)
                    index > -1 && this.selection.splice(index, 1)
                }
            },
            handleRemoveSelection (row) {
                const index = this.selection.findIndex(vpc => vpc.bk_vpc_id === row.bk_vpc_id)
                index > -1 && this.selection.splice(index, 1)
            },
            handleClear () {
                this.selection = []
            },
            handleConfirm () {
                this.$emit('change', [...this.selection])
            },
            handleCancel () {
                this.$emit('cancel')
            },
            getVpcName (vpc) {
                const id = vpc.bk_vpc_id
                const name = vpc.bk_vpc_name
                if (id === name) {
                    return id
                }
                return `${id}(${name})`
            },
            getRowClass ({ row }) {
                if (row.destroyed) {
                    return 'is-destroyed'
                }
                return ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .vpc-selector {
        .left-column {
            width: 350px;
            height: 410px;
            border-right: 1px solid $borderColor;
            .left-title {
                height: 26px;
                margin: 15px 0 0 20px;
                font-size: 20px;
                line-height: 26px;
                font-weight: normal;
            }
            .region-selector {
                display: block;
                width: 310px;
                margin: 20px auto 0;
            }
        }
    }
    .vpc-wrapper {
        position: relative;
        height: 315px;
        @include scrollbar-y;
        .vpc-options {
            display: flex;
            align-items: center;
            justify-content: space-between;
            position: sticky;
            top: 0;
            left: 0;
            padding: 10px 20px 15px;
            background-color: #fff;
            z-index: 1;
            .option-name {
                font-size: 12px;
                font-weight: 700;
            }
            .option-right {
                display: flex;
                align-items: center;
                justify-content: space-between;
                .option-count {
                    margin-right: 10px;
                    font-size: 12px;
                    font-weight: 700;
                }
                .option-checkbox {
                    direction: rtl;
                    /deep/ {
                        .bk-checkbox-text {
                            font-size: 12px;
                            font-weight: 700;
                            margin-right: 6px;
                        }
                    }
                }
            }
        }
        .vpc-list {
            padding: 0 20px;
            .vpc-item {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin: 0 0 18px 0;
                font-size: 12px;
                line-height: 16px;
                .vpc-item-name {
                    margin-right: 20px;
                    @include ellipsis;
                }
                .vpc-item-option {
                    display: flex;
                    align-items: center;
                    justify-content: flex-end;
                }
                .vpc-item-count {
                    margin-right: 46px;
                    &.en {
                        margin-right: 78px;
                    }
                }
                .vpc-item-checkbox {
                    flex: 16px 0 0;
                }
            }
        }
        .vpc-empty {
            height: 100%;
            display: flex;
            justify-content: center;
            align-items: center;
            font-size: 14px;
        }
    }
    .right-column {
        width: 500px;
        height: 410px;
        .right-options {
            height: 42px;
            line-height: 42px;
            background-color: #FAFBFD;
            padding: 0 12px;
        }
        .selection-info {
            font-size: 12px;
            font-weight: 700;
            .selection-count {
                color: $primaryColor;
            }
        }
        .selection-clear {
            float: right;
            margin-top: 11px;
        }
    }
    .bottom-options {
        padding: 0 20px;
        height: 50px;
        display: flex;
        justify-content: flex-end;
        align-items: center;
        background-color: #FAFBFD;
        border-top: 1px solid $borderColor;
    }
    .selected-vpc-table {
        /deep/ {
            .bk-table-row.is-destroyed {
                color: #C4C6CC;
            }
        }
        .vpc-info {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            white-space: nowrap;
            .vpc-name {
                @include ellipsis;
            }
            .info-destroyed {
                margin-left: 4px;
                font-size: 12px;
                line-height: 18px;
                color: #EA3536;
                padding: 0 4px;
                border-radius: 2px;
                background-color: #FEDDDC;
            }
        }
    }
</style>
