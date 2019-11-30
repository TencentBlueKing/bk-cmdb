<template>
    <div class="field-view-layout">
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
                <span class="property-value">{{field.option.min || '--'}}</span>
            </div>
            <div class="property-item">
                <div class="property-name">
                    <span>{{$t('最大值')}}</span>：
                </div>
                <span class="property-value">{{field.option.max || '--'}}</span>
            </div>
        </template>
        <div class="property-item" v-else-if="['enum', 'list'].includes(field.bk_property_type)">
            <div class="property-name">
                <span>{{$t('枚举值')}}</span>：
            </div>
            <span class="property-value">{{getEnumValue()}}</span>
        </div>
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
            <span class="property-value">{{field.description || '--'}}</span>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            field: {
                type: Object,
                default: () => {}
            }
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
                    'list': this.$t('列表')
                }
            }
        },
        methods: {
            getEnumValue () {
                const value = this.field.option
                const type = this.field.bk_property_type
                if (Array.isArray(value)) {
                    if (type === 'enum') {
                        const arr = value.map(item => {
                            if (item.is_default) {
                                return `${item.name}(${this.$t('默认值')})`
                            }
                            return item.name
                        })
                        return arr.length ? arr.join(' / ') : '--'
                    } else if (type === 'list') {
                        return value.length ? value.join(' / ') : '--'
                    }
                }
                return '--'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .field-view-layout {
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        padding: 0 20px 20px;
        font-size: 14px;
        .property-item {
            min-width: 48%;
            padding-top: 20px;
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
</style>
