<template>
    <div class="model-field-wrapper">
        <div class="form-content clearfix">
            <h3>{{$t('ModelManagement["字段配置"]')}}</h3>
            <div class="form-item has-right-content">
                <label class="form-label">{{$t('ModelManagement["中文名"]')}}<span class="color-danger"> * </span></label>
                <input type="text" class="cmdb-form-input">
            </div>
            <div class="form-item has-right-content">
                <label class="form-label">{{$t('ModelManagement["英文名"]')}}<span class="color-danger"> * </span></label>
                <input type="text" class="cmdb-form-input">
            </div>
            <div class="form-item">
                <label class="form-label">{{$t('ModelManagement["单位"]')}}</label>
                <input type="text" class="cmdb-form-input">
            </div>
            <div class="form-item block">
                <label class="form-label">{{$t('ModelManagement["提示语"]')}}</label>
                <input type="text" class="cmdb-form-input">
            </div>
        </div>
        <div class="form-content">
            <h3>{{$t('ModelManagement["选项"]')}}</h3>
            <div class="clearfix">
                <div class="form-item has-right-content">
                    <label class="form-label">{{$t('ModelManagement["类型"]')}}</label>
                    <bk-selector
                        class="form-selector bk-selector-small"
                        :list="fieldInfo.list"
                        :selected.sync="fieldInfo.type"
                    ></bk-selector>
                </div>
                <v-config :type="fieldInfo.type"></v-config>
            </div>
            <div class="field-config clearfix" v-if="isComponentShow">
                <component :is="`model-field-${fieldName}`"></component>
            </div>
        </div>
    </div>
</template>

<script>
    import modelFieldChar from './char'
    import modelFieldInt from './int'
    import modelFieldEnum from './enum'
    import modelFieldAsst from './asst'
    import vConfig from './config'
    export default {
        components: {
            modelFieldChar,
            modelFieldInt,
            modelFieldEnum,
            modelFieldAsst,
            vConfig
        },
        data () {
            return {
                fieldInfo: {
                    type: 'singlechar',
                    list: [{
                        id: 'singlechar',
                        name: this.$t('ModelManagement["短字符"]')
                    }, {
                        id: 'int',
                        name: this.$t('ModelManagement["数字"]')
                    }, {
                        id: 'enum',
                        name: this.$t('ModelManagement["枚举"]')
                    }, {
                        id: 'date',
                        name: this.$t('ModelManagement["日期"]')
                    }, {
                        id: 'time',
                        name: this.$t('ModelManagement["时间"]')
                    }, {
                        id: 'longchar',
                        name: this.$t('ModelManagement["长字符"]')
                    }, {
                        id: 'singleasst',
                        name: this.$t('ModelManagement["单关联"]')
                    }, {
                        id: 'multiasst',
                        name: this.$t('ModelManagement["多关联"]')
                    }, {
                        id: 'objuser',
                        name: this.$t('ModelManagement["用户"]')
                    }, {
                        id: 'timezone',
                        name: this.$t('ModelManagement["时区"]')
                    }, {
                        id: 'bool',
                        name: 'bool'
                    }]
                },
                charMap: ['singlechar', 'multichar'],
                asstMap: ['singleasst', 'multiasst']
            }
        },
        computed: {
            isComponentShow () {
                return ['singlechar', 'multichar', 'singleasst', 'multiasst', 'enum', 'int'].indexOf(this.fieldInfo.type) !== -1
            },
            fieldName () {
                let {
                    type
                } = this.fieldInfo
                if (this.charMap.indexOf(type) !== -1) {
                    return 'char'
                } else if (this.asstMap.indexOf(type) !== -1) {
                    return 'asst'
                }
                return type
            }
        }
    }
</script>


<style lang="scss" scoped>
    .model-field-wrapper {
        .form-content {
            margin-bottom: 30px;
            h3 {
                margin-bottom: 10px;
                padding-left: 4px;
                font-size: 14px;
                line-height: 1;
                border-left: 4px solid $cmdbTextColor;
            }
            .form-item {
                &.block {
                    margin-top: 20px;
                    .cmdb-form-input {
                        width: 590px;
                    }
                }
            }
            .field-config {
                margin-top: 20px;
            }
        }
    }
</style>
