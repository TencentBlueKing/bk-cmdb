<template>
    <bk-select
        multiple
        searchable
        display-tag
        v-bind="$attrs"
        v-model="localValue">
        <bk-option
            v-for="template in list"
            :key="template.id"
            :id="template.id"
            :name="template.name">
        </bk-option>
    </bk-select>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-search-service-template',
        props: {
            value: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                list: [],
                requestId: Symbol('serviceTemplate')
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        },
        created () {
            this.getServiceTemplate()
        },
        methods: {
            async getServiceTemplate () {
                try {
                    const { info } = await this.$store.dispatch('serviceTemplate/searchServiceTemplateWithoutDetails', {
                        params: {
                            bk_biz_id: this.bizId
                        },
                        config: {
                            requestId: this.requestId,
                            fromCache: true
                        }
                    })
                    this.list = info
                } catch (error) {
                    console.error(error)
                    this.list = []
                }
            }
        }
    }
</script>
