<template>
    <div class="vpc-selector clearfix">
        <div class="left-column fl">
            <h2 class="left-title">{{$t('添加VPC')}}</h2>
            <bk-select class="region-selector"
                searchable
                v-model="currentRegion"
                :loading="$loading(request.region)">
                <bk-option v-for="region in regions"
                    :key="region.bk_region"
                    :id="region.bk_region"
                    :name="region.bk_region_name">
                </bk-option>
            </bk-select>
            <div class="vpc-wrapper" v-bkloading="{ isLoading: $loading([request.vpc, request.region]) }">
                <template v-if="hasVpc">
                    <div class="vpc-options clearfix">
                        <span class="option-name fl">{{$t('VPC名称')}}</span>
                        <bk-checkbox class="option-checkbox fr"
                            :checked="selectState.all"
                            :indeterminate="selectState.indeterminate"
                            @change="handleToggleSelectAll">
                            {{$t('全选')}}
                        </bk-checkbox>
                    </div>
                    <ul class="vpc-list">
                        <li class="vpc-item"
                            v-for="vpc in regionVPC[currentRegion]"
                            :key="vpc.bk_vpc_id">
                            {{getVpcName(vpc)}}
                            <bk-checkbox class="vpc-item-checkbox fr"
                                :checked="selectionMap.hasOwnProperty(vpc.bk_vpc_id)"
                                @change="handleVpcChange(vpc, ...arguments)">
                            </bk-checkbox>
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
            <bk-table
                :height="368"
                :data="selection"
                :outer-border="false"
                :header-border="false">
                <bk-table-column label="VPC" prop="bk_vpc_id" show-overflow-tooltip>
                    <template slot-scope="{ row }">{{getVpcName(row)}}</template>
                </bk-table-column>
                <bk-table-column :label="$t('地域')" prop="bk_region_name" show-overflow-tooltip>
                    <template slot-scope="{ row }">{{getRegionName(row)}}</template>
                </bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="bk_host_count"></bk-table-column>
                <bk-table-column :label="$t('操作')">
                    <template slot-scope="{ row }">
                        <bk-link theme="primary" @click="handleRemoveSelection(row)">{{$t('移除')}}</bk-link>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="clearfix"></div>
        <div class="bottom-options">
            <bk-button class="mr10" theme="primary" @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'task-vpc-selector',
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
                regions: [],
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
        created () {
            this.getRegions()
        },
        methods: {
            async getRegions () {
                try {
                    const regions = await this.$store.dispatch('cloud/resource/findRegion', {
                        params: {
                            bk_account_id: this.account,
                            with_host_count: false
                        },
                        config: {
                            requestId: this.request.region
                        }
                    })
                    this.regions = regions
                    this.currentRegion = regions[0] ? regions[0].bk_region : ''
                } catch (e) {
                    console.error(e)
                    this.regions = []
                }
            },
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
            getRegionName (vpc) {
                const region = this.regions.find(region => region.bk_region === vpc.bk_region)
                return region ? region.bk_region_name : '--'
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
            position: sticky;
            top: 0;
            left: 0;
            padding: 10px 20px 15px;
            background-color: #fff;
            z-index: 1;
            .option-name {
                font-size: 12px;
                font-weight: 700;
                line-height: 16px;
            }
            .option-checkbox {
                direction: rtl;
                /deep/ {
                    .bk-checkbox-text {
                        margin-right: 6px;
                    }
                }
            }
        }
        .vpc-list {
            padding: 0 20px;
            .vpc-item {
                margin: 0 0 18px 0;
                font-size: 12px;
                line-height: 16px;
                @include ellipsis;
            }
        }
        .vpc-empty {
            height: 100%;
            display: flex;
            justify-content: center;
            align-items: center;
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
</style>
