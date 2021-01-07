<template>
    <div class="form-layout">
        <div class="form-groups" v-if="hasAvaliableGroups">
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="groupedProperties[groupIndex].length">
                    <h3 class="group-name">{{group['bk_group_name']}}</h3>
                    <ul class="property-list clearfix">
                        <li class="property-item fl"
                            v-for="(property, propertyIndex) in groupedProperties[groupIndex]"
                            :key="propertyIndex">
                            <div class="property-name">
                                <cmdb-form-bool class="property-name-checkbox"
                                    :id="`property-name-${property['bk_property_id']}`"
                                    v-model="editable[property['bk_property_id']]">
                                </cmdb-form-bool>
                                <label class="property-name-text"
                                    :for="`property-name-${property['bk_property_id']}`"
                                    :class="{required: property['isrequired']}">
                                    {{property['bk_property_name']}}
                                </label>
                                <i class="property-name-tooltips bk-icon icon-info-circle-shape" v-if="property['placeholder']" v-tooltip="htmlEncode(property['placeholder'])"></i>
                            </div>
                            <div class="property-value">
                                <component class="form-component"
                                    v-if="property['bk_property_type'] === 'enum'"
                                    :is="`cmdb-form-${property['bk_property_type']}`"
                                    :class="{error: errors.has(property['bk_property_id'])}"
                                    :disabled="!editable[property['bk_property_id']]"
                                    :options="property.option || []"
                                    :data-vv-name="property['bk_property_name']"
                                    v-validate="getValidateRules(property)"
                                    v-model.trim="values[property['bk_property_id']]">
                                </component>
                                 <component class="form-component"
                                    v-else
                                    :is="`cmdb-form-${property['bk_property_type']}`"
                                    :class="{error: errors.has(property['bk_property_id'])}"
                                    :disabled="!editable[property['bk_property_id']]"
                                    :data-vv-name="property['bk_property_name']"
                                    v-validate="getValidateRules(property)"
                                    v-model.trim="values[property['bk_property_id']]">
                                </component>
                                <span class="form-error">{{errors.first(property['bk_property_name'])}}</span>
                            </div>
                        </li>
                    </ul>
                </div>
            </template>
        </div>
        <div class="form-empty" v-else>
            {{$t("Inst['暂无可批量更新的属性']")}}
        </div>
        <div class="form-options" :class="{sticky: scrollbar}">
            <slot name="details-options">
                <bk-button class="button-save" type="primary"
                    :disabled="!$authorized.update || !hasChange || $loading()"
                    @click="handleSave">
                    {{$t("Common['保存']")}}
                </bk-button>
                <bk-button class="button-cancel" @click="handleCancel">{{$t("Common['取消']")}}</bk-button>
            </slot>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        name: 'cmdb-form-multiple',
        mixins: [formMixins],
        data () {
            return {
                values: {},
                refrenceValues: {},
                editable: {},
                scrollbar: false
            }
        },
        computed: {
            changedValues () {
                const changedValues = {}
                for (let propertyId in this.values) {
                    if (this.values[propertyId] !== this.refrenceValues[propertyId]) {
                        changedValues[propertyId] = this.values[propertyId]
                    }
                }
                return changedValues
            },
            hasChange () {
                let hasChange = false
                for (let propertyId in this.editable) {
                    if (this.editable[propertyId] && this.changedValues.hasOwnProperty(propertyId)) {
                        hasChange = true
                        break
                    }
                }
                return hasChange
            },
            groupedProperties () {
                return this.$groupedProperties.map(properties => {
                    return properties.filter(property => {
                        const editable = property.editable
                        const isapi = property['bk_isapi']
                        const isonly = property.isonly
                        const isAsst = ['singleasst', 'multiasst'].includes(property['bk_property_type'])
                        return editable && !isapi && !isonly && !isAsst
                    })
                })
            },
            hasAvaliableGroups () {
                return this.groupedProperties.some(properties => !!properties.length)
            }
        },
        watch: {
            properties () {
                this.initValues()
                this.initEditableStatus()
            }
        },
        created () {
            this.initValues()
            this.initEditableStatus()
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$el, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$el, this.checkScrollbar)
        },
        methods: {
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            initValues () {
                this.values = this.$tools.getInstFormValues(this.properties, {})
                this.refrenceValues = this.$tools.clone(this.values)
            },
            initEditableStatus () {
                const editable = {}
                this.groupedProperties.forEach(properties => {
                    properties.forEach(property => {
                        editable[property['bk_property_id']] = false
                    })
                })
                this.editable = editable
            },
            htmlEncode (placeholder) {
                let temp = document.createElement('div')
                temp.innerHTML = placeholder
                let output = temp.innerText
                temp = null
                return output
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
            getMultipleValues () {
                const multipleValues = {}
                for (let propertyId in this.editable) {
                    if (this.editable[propertyId]) {
                        multipleValues[propertyId] = this.values[propertyId]
                    }
                }
                return multipleValues
            },
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        this.$emit('on-submit', this.getMultipleValues())
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
    .form-layout{
        height: 100%;
        @include scrollbar;
    }
    .form-groups{
        padding: 0 0 0 32px;
    }
    .property-group{
        padding: 7px 0 10px 0;
        &:first-child{
            padding: 28px 0 10px 0;
        }
    }
    .group-name{
        font-size: 14px;
        line-height: 14px;
        color: #333948;
        overflow: visible;
    }
    .property-list{
        padding: 4px 0;
        .property-item{
            width: 50%;
            margin: 12px 0 0;
            padding: 0 54px 0 0;
            font-size: 12px;
            .property-name{
                display: block;
                margin: 6px 0 9px;
                color: $cmdbTextColor;
                font-size: 0;
                line-height: 18px;
            }
            .property-name-checkbox{
                transform: scale(0.667);
                vertical-align: top;
                margin: 0 6px 0 0;
            }
            .property-name-text{
                position: relative;
                display: inline-block;
                max-width: calc(100% - 20px);
                padding: 0 10px 0 0;
                vertical-align: top;
                font-size: 12px;
                @include ellipsis;
                &.required:after{
                    position: absolute;
                    left: 100%;
                    top: 0;
                    margin: 0 0 0 -10px;
                    content: "*";
                    color: #ff5656;
                }
            }
            .property-name-tooltips{
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                font-size: 16px;
                color: #ffb400;
            }
            .property-value{
                height: 36px;
                line-height: 36px;
                font-size: 12px;
                position: relative;
            }
        }
    }
    .form-options{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 28px 32px 0;
        &.sticky {
            padding: 10px 32px;
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fff;
        }
        .button-save{
            min-width: 76px;
            margin-right: 4px;
        }
        .button-cancel{
            min-width: 76px;
            background-color: #fff;
        }
    }
    .form-error {
        position: absolute;
        top: 100%;
        left: 0;
        line-height: 14px;
        font-size: 12px;
        color: #ff5656;
    }
    .form-empty{
        height: 100%;
        text-align: center;
        &:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            height: 100%;
            width: 0;
        }
    }
</style>