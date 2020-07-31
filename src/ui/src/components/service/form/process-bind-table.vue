<template>
    <cmdb-form-table v-bkloading="{ isLoading: pending }"
        v-model="bindList"
        :mode="mode"
        :disabled-check="disabledCheck"
        :options="property.option || []">
        <template slot="col-ip" slot-scope="{ row, col, sizes, disabled, $index }">
            <cmdb-input-select
                name="ip"
                :disabled="disabled"
                :placeholder="$t('请选择或输入IP')"
                :options="IPList"
                :validate="IPRules"
                v-bind="sizes"
                v-model="row['ip'].value"
                v-bk-tooltips="{
                    disabled: !disabled,
                    allowHtml: true,
                    content: `#disabled-tips-${col['bk_property_id']}_${$index}`,
                    placement: 'top'
                }">
            </cmdb-input-select>
        </template>
        <template slot="disabled-tips" slot-scope="{ col, $index, disabled }">
            <span :id="`disabled-tips-${col['bk_property_id']}_${$index}`" v-show="disabled">
                <i18n path="进程表单锁定提示">
                    <bk-link theme="primary" @click="handleRedirect" place="link">{{$t('跳转服务模板')}}</bk-link>
                </i18n>
            </span>
        </template>
    </cmdb-form-table>
</template>

<script>
    export default {
        props: {
            serviceTemplateId: Number,
            processTemplateId: Number,
            hostId: Number,
            properties: Array,
            processTemplate: Object,
            list: Array
        },
        data () {
            return {
                bindList: [],
                IPList: [],
                pending: true
            }
        },
        computed: {
            property () {
                return this.properties.find(property => property.bk_property_id === 'bind_info')
            },
            mode () {
                return this.processTemplateId ? 'update' : 'input'
            },
            IPRules () {
                const IPProperty = this.property.option.find(property => property.bk_property_id === 'ip')
                if (!IPProperty) {
                    return {}
                }
                const rules = {}
                if (IPProperty.isrequired) {
                    rules.required = true
                }
                rules.regex = IPProperty.option
                return rules
            },
            bindedProperties () {
                const bindedProperties = {}
                if (this.processTemplate && Object.keys(this.processTemplate).length) {
                    const bindInfoList = this.processTemplate.property.bind_info.value
                    bindInfoList.forEach((row, index) => {
                        Object.keys(row).forEach(key => {
                            bindedProperties[`${key}_${index}`] = row[key].as_default_value
                        })
                    })
                }
                return bindedProperties
            }
        },
        watch: {
            list: {
                handler (value) {
                    const formattedValue = (value || []).map(item => {
                        const row = { ...item }
                        Object.keys(row).forEach(key => {
                            const field = row[key]
                            if (field !== null && typeof field === 'object') {
                                row[key] = { value: field.value }
                            } else {
                                row[key] = { value: field }
                            }
                        })
                        return row
                    })
                    this.bindList = formattedValue
                },
                immediate: true
            },
            bindList: {
                handler (value) {
                    this.$emit('change', value)
                },
                deep: true
            }
        },
        created () {
            try {
                this.getBindIPList()
            } catch (error) {
                console.error(error)
            } finally {
                this.pending = false
            }
        },
        methods: {
            async getBindIPList () {
                try {
                    const { options } = await this.$store.dispatch('serviceInstance/getInstanceIpByHost', {
                        hostId: this.hostId,
                        config: {
                            requestId: 'getInstanceIpByHost'
                        }
                    })
                    this.IPList = options.map(ip => ({ id: ip, name: ip }))
                } catch (error) {
                    this.IPList = []
                    console.error(error)
                }
            },
            disabledCheck (field, property, index) {
                return this.bindedProperties[`${property.bk_property_id}_${index}`]
            },
            handleRedirect () {
                this.$routerActions.redirect({
                    name: 'operationalTemplate',
                    params: {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        templateId: this.serviceTemplateId
                    },
                    history: true
                })
            }
        }
    }
</script>
