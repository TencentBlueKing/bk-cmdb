<template>
    <bk-sideslider
        :width="395"
        :is-show.sync="isShow"
        :title="$t('添加查询条件')"
        @hidden="hanldeHidden">
        <div class="property-selector-content" slot="content">
            <div class="property-selector-options">
                <bk-input class="options-filter"
                    v-model.trim="filter"
                    right-icon="icon-search"
                    clearable>
                </bk-input>
            </div>
            <div class="property-selector-group"
                v-for="model in models"
                v-show="isShowGroup(model)"
                :key="model.id">
                <label class="group-label">{{model.bk_obj_name}}</label>
                <div class="group-property-list">
                    <bk-checkbox class="group-property-item"
                        v-for="property in matchedPropertyMap[model.bk_obj_id]"
                        v-show="isShowProperty(property)"
                        :key="property.id"
                        :title="property.bk_property_name"
                        :checked="isChecked(property)"
                        @change="handleChange(property, ...arguments)">
                        {{property.bk_property_name}}
                    </bk-checkbox>
                </div>
            </div>
        </div>
        <div class="property-selector-footer" slot="footer">
            <bk-button class="mr10" theme="primary" @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" @click="handleConfirm">{{$t('取消')}}</bk-button>
        </div>
    </bk-sideslider>
</template>

<script>
    export default {
        props: {
            selected: {
                type: Array,
                default: () => ([])
            },
            handler: Function
        },
        inject: ['dynamicGroupForm'],
        data () {
            return {
                isShow: false,
                filter: '',
                localSelected: [...this.selected],
                matchedPropertyMap: this.dynamicGroupForm.propertyMap
            }
        },
        computed: {
            target () {
                return this.dynamicGroupForm.formData.bk_obj_id
            },
            models () {
                if (this.target === 'host') {
                    return this.dynamicGroupForm.availableModels
                }
                return this.dynamicGroupForm.availableModels.filter(model => model.bk_obj_id === this.target)
            },
            propertyMap () {
                return this.dynamicGroupForm.propertyMap
            }
        },
        watch: {
            filter (filter) {
                this.filterTimer && clearTimeout(this.filterTimer)
                this.filterTimer = setTimeout(() => this.handleFilter(filter), 500)
            }
        },
        methods: {
            handleFilter (filter) {
                if (!filter.length) {
                    this.matchedPropertyMap = this.propertyMap
                } else {
                    const matchedPropertyMap = {}
                    const lowerCaseFilter = filter.toLowerCase()
                    Object.keys(this.propertyMap).forEach(modelId => {
                        matchedPropertyMap[modelId] = this.propertyMap[modelId].filter(property => {
                            const lowerCaseName = property.bk_property_name.toLowerCase()
                            return lowerCaseName.indexOf(lowerCaseFilter) > -1
                        })
                    })
                    this.matchedPropertyMap = matchedPropertyMap
                }
            },
            isShowGroup (model) {
                return !!this.matchedPropertyMap[model.bk_obj_id].length
            },
            isShowProperty (property) {
                const modelId = property.bk_obj_id
                return this.matchedPropertyMap[modelId].some(target => target === property)
            },
            isChecked (property) {
                return this.localSelected.some(target => target.id === property.id)
            },
            handleChange (property, checked) {
                if (checked) {
                    this.localSelected.push(property)
                } else {
                    const index = this.localSelected.findIndex(target => target.id === property.id)
                    index > -1 && this.localSelected.splice(index, 1)
                }
            },
            handleConfirm () {
                this.handler && this.handler([...this.localSelected])
                this.isShow = false
            },
            handleCancel () {
                this.isShow = false
            },
            show () {
                this.isShow = true
            },
            hanldeHidden () {
                this.$emit('close')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .property-selector-content {
        height: 100%;
        padding: 10px 20px;
        @include scrollbar-y;
    }
    .property-selector-group {
        margin-top: 15px;
        .group-label {
            display: block;
            font-weight: bold;
            font-size: 14px;
            color: #313237;
        }
        .group-property-list {
            display: flex;
            flex-direction: row;
            flex-wrap: wrap;
            .group-property-item {
                display: inline-flex;
                align-items: center;
                flex: 50%;
                margin: 20px 0 0 0;
                /deep/ {
                    .bk-checkbox {
                        flex: 16px 0 0;
                    }
                    .bk-checkbox-text {
                        padding-right: 15px;
                        @include ellipsis;
                    }
                }
            }
        }
    }
    .property-selector-footer {
        display: flex;
        height: 100%;
        width: 100%;
        align-items: center;
        border-top: 1px solid $borderColor;
        padding: 0 20px;
    }
</style>
