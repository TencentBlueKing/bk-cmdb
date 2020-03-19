<template>
    <bk-select v-if="display === 'selector'"
        searchable
        :readonly="readonly"
        :disabled="disabled"
        :placeholder="$t('请选择xx', { name: $t('地域') })"
        :loading="$loading(request)"
        v-model="selected">
        <bk-option v-for="region in regions"
            :key="region.bk_region"
            :name="region.bk_region_name"
            :id="region.bk_region">
        </bk-option>
    </bk-select>
    <span v-else>{{getRegionInfo()}}</span>
</template>

<script>
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
                request: `taskRegionSelection-${this.account}`
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
                            requestId: this.request,
                            fromCache: true,
                            cacheExpire: 'page'
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
