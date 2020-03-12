<template>
    <div class="vpc-selector">
        <div class="content">
            <h2 class="title">{{$t('添加VPC')}}</h2>
            <bk-form form-type="vertical" class="mt10">
                <bk-form-item :label="$t('地域')">
                    <bk-select class="form-selector" v-model="form.region"
                        :clearable="false"
                        :loading="$loading(request.regions)"
                        :popover-options="{ boundary: 'window' }"
                        @change="handleRegionChange">
                        <bk-option
                            v-for="(region, index) in regions"
                            :key="index"
                            :id="region.bk_region"
                            :name="region.bk_region_name">
                        </bk-option>
                    </bk-select>
                </bk-form-item>
                <bk-form-item label="VPC">
                    <bk-select class="form-selector" v-model="form.vpc"
                        :loading="$loading(request.vpc)"
                        :popover-options="{ boundary: 'window' }"
                        :multiple="true">
                        <bk-option
                            v-for="vpc in VPCList"
                            :key="vpc.bk_vpc_id"
                            :id="vpc.bk_vpc_id"
                            :name="vpc.bk_vpc_name">
                        </bk-option>
                    </bk-select>
                </bk-form-item>
            </bk-form>
        </div>
        <div class="options">
            <bk-button theme="primary" class="mr10"
                :disabled="!form.region || !form.vpc.length"
                @click="handleConfirm">
                {{$t('确定')}}
            </bk-button>
            <bk-button @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'resource-vpc-selector',
        props: {
            accountId: Number,
            selected: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                regions: [],
                VPCList: [],
                form: {
                    region: '',
                    vpc: []
                },
                request: {
                    regions: Symbol('regions'),
                    vpc: Symbol('vpc')
                }
            }
        },
        created () {
            this.getRegions()
        },
        methods: {
            async getRegions () {
                try {
                    const { info: regions } = await this.$store.dispatch('cloud/resource/findRegion', {
                        params: {
                            bk_account_id: this.accountId
                        },
                        config: {
                            requestId: this.request.regions
                        }
                    })
                    this.regions = regions
                } catch (e) {
                    console.error(e)
                    this.regions = []
                }
            },
            handleRegionChange (region) {
                this.form.vpc = []
                if (!region) {
                    return
                }
                this.getVPCList(region)
            },
            async getVPCList (region) {
                try {
                    const { info: VPCList } = await this.$store.dispatch('cloud/resource/findVPC', {
                        id: this.accountId,
                        params: {
                            bk_region: region
                        },
                        config: {
                            requestId: this.request.vpc
                        }
                    })
                    this.VPCList = VPCList
                    this.form.vpc = this.selected.filter(vpc => vpc.bk_region === this.form.region).map(vpc => vpc.bk_vpc_id)
                } catch (e) {
                    console.error(e)
                    this.VPCList = []
                }
            },
            handleConfirm () {
                this.$emit(
                    'change',
                    this.form.vpc.map(id => {
                        const vpc = this.VPCList.find(vpc => vpc.bk_vpc_id === id)
                        const region = this.regions.find(region => region.bk_region === this.form.region)
                        return {
                            ...vpc,
                            bk_region_name: region.bk_region_name
                        }
                    }),
                    this.form.region
                )
            },
            handleCancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .vpc-selector {
        .content {
            padding: 18px 24px;
        }
        .title {
            line-height: 26px;
            font-size: 20px;
            font-weight: normal;
            color: #313238;
        }
        .form-selector {
            display: block;
        }
        .options {
            display: flex;
            justify-content: flex-end;
            padding: 8px 24px 9px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
        }
    }
</style>
