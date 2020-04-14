<template>
    <div class="instance-details-wrapper" v-bkloading="{ isLoading: $loading() }">
        <bk-table
            :data="detailsData"
            :max-height="$APP.height - 120">
            <bk-table-column prop="property_name" :label="$t('属性名称')" show-overflow-tooltip></bk-table-column>
            <bk-table-column prop="property_value" :label="$t('变更前')" show-overflow-tooltip>
                <template slot-scope="{ row }">{{getCellValue(row, 'before')}}</template>
            </bk-table-column>
            <bk-table-column prop="show_value" :label="$t('变更后')" show-overflow-tooltip>
                <template slot-scope="{ row }">{{getCellValue(row, 'after')}}</template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import formatter from '@/filters/formatter'
    export default {
        props: {
            module: Object,
            instance: Object,
            type: String,
            properties: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                detailsData: []
            }
        },
        created () {
            switch (this.type) {
                case 'added':
                    this.initTemplateData()
                    break
                case 'changed':
                    this.initChangedData()
                    break
                case 'removed':
                    this.initRemovedData()
                    break
                case 'others':
                    this.initOthersData()
            }
        },
        methods: {
            initChangedData () {
                this.detailsData = this.instance.changed_attributes
            },
            async initTemplateData () {
                try {
                    const { info } = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: this.$injectMetadata({
                            service_template_id: this.instance.service_instance.service_template_id
                        }, { injectBizId: true })
                    })
                    const process = info.find(process => process.id === this.module.process_template_id)
                    const details = []
                    Object.keys(process.property).forEach(key => {
                        const property = this.properties.find(property => property.bk_property_id === key)
                        if (property && !['', null].includes(process.property[key].value)) {
                            details.push({
                                property_id: key,
                                property_name: property.bk_property_name,
                                property_value: null,
                                template_property_value: {
                                    ...process.property[key]
                                }
                            })
                        }
                    })
                    this.detailsData = details
                } catch (e) {
                    console.error(e)
                }
            },
            initRemovedData () {
                const details = []
                Object.keys(this.instance.process).forEach(key => {
                    const property = this.properties.find(property => property.bk_property_id === key)
                    if (property && !['', null].includes(this.instance.process[key])) {
                        details.push({
                            property_id: key,
                            property_name: property.bk_property_name,
                            property_value: this.instance.process[key],
                            template_property_value: {
                                value: this.$t('该进程已删除')
                            }
                        })
                    }
                })
                this.detailsData = details
            },
            initOthersData () {
                this.detailsData = [...this.instance.changed_attributes]
            },
            getCellValue (row, type) {
                const propertyId = row.property_id
                let value = row.property_value
                if (type === 'after') {
                    value = typeof row.template_property_value === 'object' ? row.template_property_value.value : row.template_property_value
                }
                if (this.type !== 'others') {
                    const property = this.properties.find(property => property.bk_property_id === propertyId)
                    return formatter(value, property)
                }
                return formatter(value, 'singlechar')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .instance-details-wrapper {
        padding: 20px;
    }
</style>
