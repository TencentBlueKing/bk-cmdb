<template>
    <bk-select v-if="display === 'selector'"
        searchable
        :clearable="false"
        :readonly="readonly"
        :disabled="disabled"
        :placeholder="$t('请选择xx', { name: $t('地域') })"
        :loading="$loading(request)"
        v-model="selected">
        <bk-option v-for="region in regions"
            :key="region.bk_region"
            :name="region.bk_region_name"
            :id="region.bk_region">
            <div class="region-info">
                <span class="region-name" v-bk-overflow-tips>{{region.bk_region_name}}</span>
                <span class="region-host-count">
                    {{region.bk_host_count}}
                    <i class="icon-cc-host"></i>
                </span>
            </div>
        </bk-option>
    </bk-select>
    <span v-else>{{getRegionInfo()}}</span>
</template>

<script>
    import symbols from '../common/symbol'
    export default {
        name: 'task-region-selector',
        props: {
            account: Number,
            display: {
                type: String,
                default: 'selector'
            },
            readonly: Boolean,
            disabled: Boolean,
            value: {
                type: [String, Number]
            }
        },
        data () {
            return {
                regions: [],
                request: symbols.get(`taskRegionSelection-${this.account}`)
            }
        },
        computed: {
            selected: {
                get () {
                    return this.value
                },
                set (value, oldValue) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            }
        },
        created () {
            // 为0时是默认云区域，无地域信息
            this.account && this.getRegions()
        },
        methods: {
            async getRegions () {
                try {
                    const regions = await this.$store.dispatch('cloud/resource/findRegion', {
                        params: {
                            bk_account_id: this.account,
                            with_host_count: true
                        },
                        config: {
                            requestId: this.request,
                            fromCache: true
                        }
                    })
                    this.regions = regions
                    this.selected = regions.length ? regions[0].bk_region : ''
                } catch (e) {
                    console.error(e)
                    this.regions = []
                }
            },
            getRegionInfo () {
                const region = this.regions.find(region => region.bk_region === this.value)
                return region ? region.bk_region_name : '--'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .region-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 14px;
        .region-name {
            @include ellipsis;
        }
        .region-host-count {
            display: flex;
            margin-left: 15px;
            align-items: center;
            color: #C4C6CC;
            .icon-cc-host {
                margin-left: 6px;
            }
        }
    }
</style>
