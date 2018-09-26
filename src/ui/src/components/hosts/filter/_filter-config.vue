<template>
    <div class="config-layout clearfix">
        <div class="config-wrapper config-unselected fl">
            <div class="wrapper-header unselected-header">
                <cmdb-selector class="header-selector" :list="modelOptions" v-model.trim="selectedModel"></cmdb-selector>
                <input class="header-filter" type="text" :placeholder="$t('Inst[\'搜索属性\']')" v-model.trim="filter">
            </div>
            <ul class="property-list">
                <li ref="unselectedPropertyItem" class="property-item" v-for="(property, index) in unselectedProperties" @click="selectProperty(property)">
                    <span>{{property['bk_property_name']}}</span>
                    <i class="bk-icon icon-angle-right"></i>
                </li>
            </ul>
        </div>
        <div class="config-wrapper config-selected fl">
            <div class="wrapper-header selected-header">
                <label class="header-label">{{$t("Inst['已显示属性']")}}</label>
            </div>
            <vue-draggable element="ul" class="property-list property-list-selected" v-model="selectedProperties" :options="{animation: 150}">
                <li class="property-item" v-for="(property, index) in selectedProperties">
                    <i class="icon-triple-dot"></i>
                    <span>{{property['bk_property_name']}}</span>
                    <i class="bk-icon icon-eye-slash-shape" @click="unselectProperty(property)" v-tooltip="$t('Common[\'隐藏\']')"></i>
                </li>
            </vue-draggable>
        </div>
        <div class="config-options clearfix">
            <bk-button class="config-button fl" type="primary" @click="handleApply">{{$t('Inst[\'应用\']')}}</bk-button>
            <bk-button class="config-button fl" type="default" @click="handleCancel">{{$t('Common[\'取消\']')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import vueDraggable from 'vuedraggable'
    export default {
        name: 'cmdb-filter-config',
        components: {
            vueDraggable
        },
        props: {
            properties: {
                type: Object,
                default () {
                    return {}
                }
            },
            selected: {
                type: Array,
                default () {
                    return []
                }
            },
            min: {
                type: Number,
                default: 0
            },
            max: {
                type: Number,
                default: 10
            }
        },
        data () {
            return {
                filter: '',
                selectedModel: '',
                localSelcted: []
            }
        },
        computed: {
            modelOptions () {
                return Object.keys(this.properties).map(objId => {
                    const model = this.$allModels.find(model => model['bk_obj_id'] === objId) || {}
                    return {
                        id: objId,
                        name: model['bk_obj_name']
                    }
                })
            },
            sortedProperties () {
                const properties = this.properties[this.selectedModel] || []
                return [...properties].sort((propertyA, propertyB) => {
                    return propertyA['bk_property_name'].localeCompare(propertyB['bk_property_name'], 'zh-Hans-CN', {sensitivity: 'accent'})
                })
            },
            unselectedProperties () {
                return this.sortedProperties.filter(property => {
                    return this.checkAvaliable(property) && !this.localSelcted.some(meta => meta['bk_property_id'] === property['bk_property_id'])
                })
            },
            selectedProperties: {
                get () {
                    return this.localSelcted.map(meta => {
                        return this.properties[meta['bk_obj_id']].find(property => property['bk_property_id'] === meta['bk_property_id'])
                    })
                },
                set (properties) {
                    this.localSelcted = properties.map(property => {
                        return {
                            'bk_property_id': property['bk_property_id'],
                            'bk_obj_id': property['bk_obj_id']
                        }
                    })
                }
            }
        },
        watch: {
            selected (selected) {
                this.initLocalSelected()
            },
            filter (filter) {
                this.unselectedProperties.forEach((property, index) => {
                    if (property['bk_property_name'].toLowerCase().indexOf(filter.toLowerCase()) !== -1) {
                        this.$refs.unselectedPropertyItem[index].style.display = 'block'
                    } else {
                        this.$refs.unselectedPropertyItem[index].style.display = 'none'
                    }
                })
            }
        },
        created () {
            this.initLocalSelected()
        },
        methods: {
            initLocalSelected () {
                this.localSelcted = this.selected.filter(selected => {
                    const properties = this.properties[selected['bk_obj_id']] || []
                    return properties.some(property => property['bk_property_id'] === selected['bk_property_id'] && this.checkAvaliable(property))
                })
            },
            checkAvaliable (property) {
                return !(['bk_host_innerip', 'bk_host_outerip'].includes(property['bk_property_id']) || property['bk_isapi'])
            },
            selectProperty (property) {
                if (this.localSelcted.length < this.max) {
                    this.localSelcted.push({
                        'bk_property_id': property['bk_property_id'],
                        'bk_obj_id': property['bk_obj_id']
                    })
                } else {
                    this.$info(this.$t('Common["最多选择N项"]', {n: this.max}))
                }
            },
            unselectProperty (property) {
                if (this.localSelcted.length > this.min) {
                    this.localSelcted = this.localSelcted.filter(selected => selected['bk_property_id'] !== property['bk_property_id'])
                } else {
                    this.$info(this.$t('Common["至少选择N项"]', {n: this.min}))
                }
            },
            handleApply () {
                if (this.localSelcted.length > this.max) {
                    this.$info(this.$t('Common["最多选择N项"]', {n: this.max}))
                } else if (this.localSelcted.length < this.min) {
                    this.$info(this.$t('Common["至少选择N项"]', {n: this.min}))
                } else {
                    this.$emit('on-apply', this.selectedProperties)
                }
            },
            handleCancel () {
                this.$emit('on-cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .config-layout{
        height: 100%;
    }
    .config-wrapper{
        width: 50%;
        height: calc(100% - 62px);
        border-right: 1px solid #e7e9ef;
        .wrapper-header{
            height: 78px;
            padding: 20px;
            border-top: 1px solid #e7e9ef;
            border-bottom: 1px solid #e7e9ef;
            font-size: 0;
            .header-selector{
                display: inline-block;
                vertical-align: middle;
                width: 120px;
                margin: 0 10px 0 0;
            }
            .header-label{
                display: inline-block;
                vertical-align: middle;
                line-height: 36px;
                min-width: 120px;
                font-size: 14px;
            }
            .header-filter{
                display: inline-block;
                vertical-align: middle;
                width: 120px;
                height: 36px;
                padding: 0 15px;
                border: 1px solid $cmdbBorderColor;
                border-radius: 2px;
                font-size: 14px;
            }
        }
    }
    .property-list{
        height: calc(100% - 78px);
        padding: 15px 0;
        @include scrollbar-y;
        &-selected{
            .property-item{
                cursor: move;
            }
        }
        .property-item{
            position: relative;
            height: 42px;
            line-height: 42px;
            padding: 0 0 0 27px;
            cursor: pointer;
            &:hover{
                background-color: #f9f9f9;
            }
            .icon-angle-right{
                position: absolute;
                top: 14px;
                right: 18px;
            }
            .icon-eye-slash-shape{
                position: absolute;
                top: 0;
                right: 0;
                width: 42px;
                height: 42px;
                line-height: 42px;
                text-align: center;
            }
            .icon-eye-slash-shape:hover{
                color: #f00;
            }
        }
    }
    .config-options{
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 62px;
        padding: 13px 20px;
        background-color: #f9f9f9;
        .config-button{
            width: 110px;
            margin: 0 0 0 10px;
            &:first-child{
                margin: 0;
            }
        }
    }
</style>