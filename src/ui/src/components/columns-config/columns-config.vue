<template>
    <div class="config-layout clearfix">
        <div class="config-wrapper config-unselected fl">
            <div class="wrapper-header unselected-header">
                <label class="header-label">{{$t("Inst['隐藏属性']")}}</label>
                <input class="header-filter" type="text" :placeholder="$t('Inst[\'搜索属性\']')" v-model.trim="filter">
            </div>
            <ul class="property-list property-list-unselected">
                <li ref="unselectedPropertyItem" class="property-item" v-for="(property, index) in unselectedProperties" @click="selectProperty(property)">
                    <span class="property-name">{{property['bk_property_name']}}</span>
                    <i class="bk-icon icon-angle-right"></i>
                </li>
            </ul>
        </div>
        <div class="config-wrapper config-selected fl">
            <div class="wrapper-header selected-header">
                <label class="header-label">{{$t("Inst['已显示属性']")}}</label>
            </div>
            <div class="property-list-layout">
                <ul class="property-list property-list-selected">
                    <li class="property-item disabled"
                        v-for="(property, index) in undragbbleProperties">
                        <span class="property-name" :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                    </li>
                </ul>
                <vue-draggable element="ul" class="property-list property-list-selected"
                    v-model="drabbleProperties"
                    :options="{animation: 150}">
                    <li class="property-item"
                        v-for="(property, index) in drabbleProperties">
                        <i class="icon-triple-dot"></i>
                        <span class="property-name" :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                        <i class="bk-icon icon-eye-slash-shape"
                            v-tooltip="$t('Common[\'隐藏\']')"
                            @click="unselectProperty(property)">
                        </i>
                    </li>
                </vue-draggable>
            </div>
        </div>
        <div class="config-options clearfix">
            <bk-button class="config-button fl" type="primary" @click="handleApply">{{$t('Inst[\'应用\']')}}</bk-button>
            <bk-button class="config-button fl" type="default" @click="handleCancel">{{$t('Common[\'取消\']')}}</bk-button>
            <bk-button class="config-button fr" type="default" @click="handleReset">{{$t("Common['还原默认']")}}</bk-button>
        </div>
    </div>
</template>

<script>
    import vueDraggable from 'vuedraggable'
    export default {
        name: 'cmdb-columns-config',
        components: {
            vueDraggable
        },
        props: {
            properties: {
                type: Array,
                default () {
                    return []
                }
            },
            selected: {
                type: Array,
                default () {
                    return []
                }
            },
            disabledColumns: {
                type: Array,
                default () {
                    return []
                }
            },
            min: {
                type: Number,
                default: 1
            },
            max: {
                type: Number,
                default: 20
            }
        },
        data () {
            return {
                filter: '',
                localSelected: []
            }
        },
        computed: {
            sortedProperties () {
                return [...this.properties].sort((propertyA, propertyB) => {
                    return propertyA['bk_property_name'].localeCompare(propertyB['bk_property_name'], 'zh-Hans-CN', {sensitivity: 'accent'})
                })
            },
            unselectedProperties () {
                return this.sortedProperties.filter(property => {
                    const unselected = !this.localSelected.includes(property['bk_property_id'])
                    const includesFilter = property['bk_property_name'].toLowerCase().indexOf(this.filter.toLowerCase()) !== -1
                    return unselected && includesFilter
                })
            },
            undragbbleProperties () {
                const undragbbleProperties = []
                this.localSelected.forEach(propertyId => {
                    if (this.disabledColumns.includes(propertyId)) {
                        const property = this.properties.find(property => property['bk_property_id'] === propertyId)
                        if (property) {
                            undragbbleProperties.push(property)
                        }
                    }
                })
                return undragbbleProperties
            },
            drabbleProperties: {
                get () {
                    const drabbleProperties = []
                    this.localSelected.forEach(propertyId => {
                        if (!this.disabledColumns.includes(propertyId)) {
                            const property = this.properties.find(property => property['bk_property_id'] === propertyId)
                            if (property) {
                                drabbleProperties.push(property)
                            }
                        }
                    })
                    return drabbleProperties
                },
                set (drabbleProperties) {
                    this.localSelected = [...this.undragbbleProperties, ...drabbleProperties].map(property => property['bk_property_id'])
                }
            }
        },
        watch: {
            selected (selected) {
                this.initLocalSelected()
            }
        },
        created () {
            this.initLocalSelected()
        },
        methods: {
            initLocalSelected () {
                this.localSelected = this.selected.filter(propertyId => this.properties.some(property => property['bk_property_id'] === propertyId))
            },
            selectProperty (property) {
                if (this.localSelected.length < this.max) {
                    this.localSelected.push(property['bk_property_id'])
                } else {
                    this.$info(this.$t('Common["最多选择N项"]', {n: this.max}))
                }
            },
            unselectProperty (property) {
                if (this.localSelected.length > this.min) {
                    this.localSelected = this.localSelected.filter(propertyId => propertyId !== property['bk_property_id'])
                } else {
                    this.$info(this.$t('Common["至少选择N项"]', {n: this.min}))
                }
            },
            checkDisabled (property) {
                return this.disabledColumns.includes(property['bk_property_id'])
            },
            handleApply () {
                if (this.localSelected.length > this.max) {
                    this.$info(this.$t('Common["最多选择N项"]', {n: this.max}))
                } else if (this.localSelected.length < this.min) {
                    this.$info(this.$t('Common["至少选择N项"]', {n: this.min}))
                } else {
                    this.$emit('on-apply', [...this.undragbbleProperties, ...this.drabbleProperties])
                }
            },
            handleCancel () {
                this.$emit('on-cancel')
            },
            handleReset () {
                this.$bkInfo({
                    title: this.$t("Common['是否要还原回系统默认显示属性？']"),
                    confirmFn: () => {
                        this.$emit('on-reset')
                    }
                })
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
            line-height: 36px;
            border-top: 1px solid #e7e9ef;
            border-bottom: 1px solid #e7e9ef;
            .header-label{
                display: inline-block;
                vertical-align: middle;
                min-width: 120px;
            }
            .header-filter{
                display: inline-block;
                vertical-align: middle;
                width: 120px;
                height: 36px;
                padding: 0 15px;
                border: 1px solid $cmdbBorderColor;
                border-radius: 2px;
            }
        }
    }
    .property-list-layout {
        height: calc(100% - 78px);
        padding: 15px 0;
        @include scrollbar-y;
    }
    .property-list {
        &-selected{
            .property-item{
                cursor: move;
            }
        }
        &-unselected {
            height: calc(100% - 78px);
            @include scrollbar-y;
        }
        .property-item{
            position: relative;
            height: 42px;
            line-height: 42px;
            padding: 0 0 0 27px;
            cursor: pointer;
            &.disabled {
                cursor: not-allowed;
            }
            &:hover{
                background-color: #f9f9f9;
            }
            .property-name {
                display: inline-block;
                vertical-align: top;
                max-width: calc(100% - 50px);
                @include ellipsis;
            }
            .icon-triple-dot {
                position: absolute;
                left: 15px;
                top: 19px;
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