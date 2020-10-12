<template>
    <div class="field-view-layout">
        <div class="field-view-list" ref="fieldList">
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('唯一标识')}}</span>：
                </div>
                <span class="property-value">{{field.bk_property_id}}</span>
            </div>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('字段名称')}}</span>：
                </div>
                <span class="property-value">{{field.bk_property_name}}</span>
            </div>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('字段类型')}}</span>：
                </div>
                <span class="property-value">{{fieldTypeMap[field.bk_property_type]}}</span>
            </div>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('是否可编辑')}}</span>：
                </div>
                <span class="property-value">{{field.editable ? $t('可编辑') : $t('不可编辑')}}</span>
            </div>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('是否必填')}}</span>：
                </div>
                <span class="property-value">{{field.isrequired ? $t('必填') : $t('非必填')}}</span>
            </div>
            <div class="property-item" v-if="['singlechar', 'longchar'].includes(field.bk_property_type)">
                <div class="property-name">
                    <span>{{$t('正则校验')}}</span>：
                </div>
                <span class="property-value">{{field.option || '--'}}</span>
            </div>
            <template v-else-if="['int', 'float'].includes(field.bk_property_type)">
                <div class="property-item">
                    <div class="property-name">
                        <span>{{$t('最小值')}}</span>：
                    </div>
                    <span class="property-value">{{field.option.min || (field.option.min === 0 ? 0 : '--')}}</span>
                </div>
                <div class="property-item">
                    <div class="property-name">
                        <span>{{$t('最大值')}}</span>：
                    </div>
                    <span class="property-value">{{field.option.max || (field.option.max === 0 ? 0 : '--')}}</span>
                </div>
            </template>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('单位')}}</span>：
                </div>
                <span class="property-value">{{field.unit || '--'}}</span>
            </div>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('用户提示')}}</span>：
                </div>
                <span class="property-value">{{field.placeholder || '--'}}</span>
            </div>
            <div class="property-item enum-list" v-if="['enum', 'list'].includes(field.bk_property_type)">
                <div class="property-name">
                    <span>{{$t('枚举值')}}</span>：
                </div>
                <span class="property-value" v-html="getEnumValue()"></span>
            </div>
        </div>
        <div class="btns" :class="{ 'sticky-layout': scrollbar }" v-if="canEdit">
            <bk-button class="mr10" theme="primary" @click="handleEdit">{{$t('编辑')}}</bk-button>
            <bk-button class="delete-btn" v-if="!field.ispre" @click="handleDelete">{{$t('删除')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        props: {
            field: {
                type: Object,
                default: () => {}
            },
            canEdit: Boolean
        },
        data () {
            return {
                fieldTypeMap: {
                    'singlechar': this.$t('短字符'),
                    'int': this.$t('数字'),
                    'float': this.$t('浮点'),
                    'enum': this.$t('枚举'),
                    'date': this.$t('日期'),
                    'time': this.$t('时间'),
                    'longchar': this.$t('长字符'),
                    'objuser': this.$t('用户'),
                    'timezone': this.$t('时区'),
                    'bool': 'bool',
                    'list': this.$t('列表'),
                    'organization': this.$t('组织')
                },
                scrollbar: false
            }
        },
        mounted () {
            addResizeListener(this.$refs.fieldList, this.handleScrollbar)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.fieldList, this.handleScrollbar)
        },
        methods: {
            getEnumValue () {
                const value = this.field.option
                const type = this.field.bk_property_type
                if (Array.isArray(value)) {
                    if (type === 'enum') {
                        const arr = value.map(item => {
                            if (item.is_default) {
                                return `${item.name}(${item.id}, ${this.$t('默认值')})`
                            }
                            return `${item.name}(${item.id})`
                        })
                        return arr.length ? arr.join('<br>') : '--'
                    } else if (type === 'list') {
                        return value.length ? value.join('<br>') : '--'
                    }
                }
                return '--'
            },
            handleEdit () {
                this.$emit('on-edit')
            },
            handleDelete () {
                this.$emit('on-delete')
            },
            handleScrollbar () {
                const el = this.$refs.fieldList
                this.scrollbar = el.scrollHeight !== el.offsetHeight
            }
        }
    }
</script>

<style lang="scss" scoped>
    .field-view-layout {
        height: 100%;
        overflow: hidden;
    }
    .field-view-list {
        max-height: calc(100% - 52px);
        @include scrollbar-y;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        padding: 0 20px 20px;
        font-size: 14px;
        .property-item {
            min-width: 48%;
            padding-top: 20px;
            &.enum-list {
                flex: 100%;
                .property-name,
                .property-value {
                    display: inline-block;
                    vertical-align: top;
                }
            }
        }
        .property-name {
            display: inline-block;
            color: #63656e;
            span {
                display: inline-block;
                min-width: 70px;
            }
        }
        .property-value {
            color: #313237;
        }
    }
    .btns {
        position: sticky;
        bottom: 0;
        left: 0;
        background: #ffffff;
        padding: 10px 20px;
        font-size: 0;
        .delete-btn:hover {
            color: #ffffff;
            background-color: #ff5656;
            border-color: #ff5656;
        }
        &.sticky-layout {
            border-top: 1px solid #dcdee5;
        }
    }
</style>
