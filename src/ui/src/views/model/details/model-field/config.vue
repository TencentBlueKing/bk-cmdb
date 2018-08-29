<template>
    <ul class="label-wrapper clearfix">
        <li v-if="isEditableShow">
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <span class="cmdb-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</span>
                <input type="checkbox" v-model="localValue.editable">
            </label>
        </li>
        <li v-if="isRequiredShow">
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <span class="cmdb-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</span>
                <input type="checkbox" v-model="localValue.isrequired">
            </label>
        </li>
        <li v-if="isOnlyShow">
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <span class="cmdb-checkbox-text mr5">{{$t('ModelManagement["是否唯一"]')}}</span>
                <input type="checkbox" v-model="localValue.isonly">
            </label>
        </li>
    </ul>
</template>

<script>
    export default {
        props: {
            isReadOnly: {
                type: Boolean,
                default: false
            },
            type: {
                type: String,
                required: true
            },
            editable: {
                type: Boolean,
                default: true
            },
            isrequired: {
                type: Boolean,
                default: false
            },
            isonly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                editableMap: [
                    'singlechar',
                    'int',
                    'enum',
                    'date',
                    'time',
                    'longchar',
                    'singleasst',
                    'multiasst',
                    'objuser',
                    'timezone',
                    'bool'
                ],
                isrequiredMap: [
                    'singlechar',
                    'int',
                    'date',
                    'time',
                    'longchar',
                    'objuser',
                    'timezone'
                ],
                isonlyMap: [
                    'singlechar',
                    'int',
                    'longchar'
                ],
                localValue: {
                    editable: this.editable,
                    isrequired: this.isrequired,
                    isonly: this.isonly
                }
            }
        },
        computed: {
            isEditableShow () {
                return this.editableMap.indexOf(this.type) !== -1
            },
            isRequiredShow () {
                return this.isrequiredMap.indexOf(this.type) !== -1
            },
            isOnlyShow () {
                return this.isonlyMap.indexOf(this.type) !== -1
            }
        },
        methods: {
            getValue () {
                let {
                    editable,
                    isrequired,
                    isonly
                } = this.localValue
                return {
                    editable,
                    isrequired: this.isRequiredShow ? isrequired : false,
                    isonly: this.isOnlyShow ? isonly : false
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .label-wrapper {
        float: left;
        vertical-align: middle;
        li {
            float: left;
            height: 30px;
            label {
                line-height: 1;
            }
        }
    }
</style>
