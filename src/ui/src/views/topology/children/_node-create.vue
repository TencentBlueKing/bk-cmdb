<template>
    <div class="node-create-layout">
        <h2 class="node-create-title">{{$t('BusinessTopology["新增节点"]')}}</h2>
        <div class="node-create-path">{{topoPath}}</div>
        <div class="node-create-form">
            <div class="form-group"
                v-for="(property, index) in sortedProperties"
                :key="index">
                <label :class="['form-label', {
                    required: property['isrequired']
                }]">
                    {{property['bk_property_name']}}
                </label>
                <component :is="`cmdb-form-${property['bk_property_type']}`"
                    :data-vv-name="property['bk_property_id']"
                    :data-vv-as="property['bk_property_name']"
                    :options="property.option || []"
                    v-validate="getValidateRules(property)"
                    v-model.trim="values[property['bk_property_id']]">
                </component>
                <span class="form-error">{{errors.first(property['bk_property_id'])}}</span>
            </div>
        </div>
        <div class="node-create-options">
            <bk-button type="primary"
                :disabled="$loading() || errors.any()"
                @click="handleSave">
                {{$t('Common["保存"]')}}
            </bk-button>
            <bk-button type="default" @click="handleCancel">{{$t('Common["取消"]')}}</bk-button>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            state: {
                type: Object,
                required: true
            },
            properties: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                values: {}
            }
        },
        computed: {
            topoPath () {
                const path = this.getStatePath(this.state)
                return path.map(state => state.node['bk_inst_name']).join('-')
            },
            sortedProperties () {
                const sortedProperties = this.properties.filter(property => !['singleasst', 'multiasst'].includes(property['bk_property_type']))
                sortedProperties.sort((propertyA, propertyB) => {
                    return this.$tools.getPropertyPriority(propertyA) - this.$tools.getPropertyPriority(propertyB)
                })
                return sortedProperties
            }
        },
        watch: {
            properties () {
                this.initValues()
            }
        },
        created () {
            this.initValues()
        },
        methods: {
            initValues () {
                this.values = this.$tools.getInstFormValues(this.properties, {})
            },
            getStatePath (state) {
                let rootState = state
                let path = [state]
                if (state.parent) {
                    rootState = rootState.parent.state
                    path = [...this.getStatePath(rootState), ...path]
                }
                return path
            },
            getValidateRules (property) {
                const rules = {}
                const {
                    bk_property_type: propertyType,
                    option,
                    isrequired
                } = property
                if (isrequired) {
                    rules.required = true
                }
                if (option) {
                    if (propertyType === 'int') {
                        if (option.hasOwnProperty('min') && !['', null, undefined].includes(option.min)) {
                            rules['min_value'] = option.min
                        }
                        if (option.hasOwnProperty('max') && !['', null, undefined].includes(option.max)) {
                            rules['max_value'] = option.max
                        }
                    } else if (['singlechar', 'longchar'].includes(propertyType)) {
                        rules['regex'] = option
                    }
                }
                if (['singlechar', 'longchar'].includes(propertyType)) {
                    rules[propertyType] = true
                }
                if (propertyType === 'int') {
                    rules['numeric'] = true
                }
                return rules
            },
            handleSave () {
                this.$validator.validateAll().then(isValid => {
                    if (isValid) {
                        this.$emit('on-submit', this.values)
                    }
                })
            },
            handleCancel () {
                this.$emit('on-cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .node-create-layout {
        position: relative;
    }
    .node-create-title {
        position: absolute;
        top: -20px;
        left: 0;
        padding: 0 20px;
        line-height: 30px;
        font-size: 22px;
        color: #333948;
    }
    .node-create-path {
        padding: 20px;
        font-size: 14px;
    }
    .node-create-form {
        max-height: 400px;
        @include scrollbar-y;
        .form-group {
            position: relative;
            padding: 0 20px 16px;
            .form-label {
                position: relative;
                display: inline-block;
                padding: 0 10px 0 0;
                max-width: 100%;
                line-height: 24px;
                @include ellipsis;
                &.required:after {
                    position: absolute;
                    right: 0;
                    content: '*';
                    color: #ff5656;
                }
            }
            .form-error {
                position: absolute;
                bottom: -2px;
                left: 20px;
                font-size: 12px;
                color: #ff5656;
            }
        }
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #fafbfd;
    }
</style>