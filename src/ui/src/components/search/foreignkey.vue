<template>
    <bk-select v-if="displayType === 'selector'"
        searchable
        v-model="localValue"
        v-bind="$attrs"
        :multiple="multiple"
        :loading="$loading(requestId)"
        @toggle="handleToggle">
        <bk-option v-for="option in options"
            :key="option.bk_cloud_id"
            :id="option.bk_cloud_id"
            :name="option.bk_cloud_name">
        </bk-option>
    </bk-select>
    <span v-else>
        <slot name="info-prepend"></slot>
        {{info}}
    </span>
</template>

<script>
    import activeMixin from './mixin-active'
    export default {
        name: 'cmdb-search-foreignkey',
        mixins: [activeMixin],
        props: {
            value: {
                type: [String, Array, Number],
                default: () => ([])
            },
            displayType: {
                type: String,
                default: 'selector',
                validator (type) {
                    return ['selector', 'info'].includes(type)
                }
            }
        },
        data () {
            return {
                options: [],
                requestId: 'searchForeignkey'
            }
        },
        computed: {
            multiple () {
                return Array.isArray(this.value)
            },
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            },
            info () {
                const values = Array.isArray(this.value) ? this.value : [this.value]
                const info = []
                values.forEach(value => {
                    const data = this.options.find(data => data.bk_cloud_id === value)
                    data && info.push(data.bk_cloud_name)
                })
                return info.join(' | ')
            }
        },
        async created () {
            try {
                const { info } = await this.$store.dispatch('cloud/area/findMany', {
                    params: {
                        page: {
                            sort: 'bk_cloud_name'
                        }
                    },
                    config: {
                        requestId: this.requestId,
                        fromCache: true
                    }
                })
                this.options = info
            } catch (error) {
                console.error(error)
            }
        }
    }
</script>
