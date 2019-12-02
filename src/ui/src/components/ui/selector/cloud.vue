<template>
    <bk-select
        class="cloud-selector"
        v-model="selected"
        :searchable="searchable"
        :clearable="allowClear"
        :disabled="disabled"
        :popover-options="popoverOptions"
    >
        <bk-option
            v-for="(option, index) in data"
            :key="index"
            :id="option.bk_cloud_id"
            :name="option.bk_cloud_name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-cloud-selector',
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            allowClear: {
                type: Boolean,
                default: false
            },
            popoverOptions: {
                type: Object,
                default () {
                    return {}
                }
            },
            requestConfig: {
                type: Object,
                default () {
                    return {}
                }
            }
        },
        data () {
            return {
                data: [],
                selected: ''
            }
        },
        computed: {
            searchable () {
                return this.data.length > 7
            }
        },
        watch: {
            value: {
                handler (value) {
                    if (value !== null) {
                        this.selected = value
                    }
                },
                immediate: true
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            }
        },
        created () {
            this.getData()
        },
        methods: {
            async getData () {
                const data = await this.$store.dispatch('cloudarea/getCloudarea', {
                    params: {},
                    config: { ...this.requestConfig, ...{ fromCache: true } }
                })
                this.data = data.info || []
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cloud-selector {
        width: 100%;
    }
</style>
